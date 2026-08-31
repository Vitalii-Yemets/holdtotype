package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/sys/windows"

	"holdtotype/internal/logfilter"
)

type streamServer struct {
	addr    string
	dir     string
	cmd     *exec.Cmd
	doneCh  chan struct{}
	job     windows.Handle
	threads int

	mu       sync.Mutex
	exited   bool
	stopping bool
}

func streamModelArgs(dir string) ([]string, error) {
	pick := func(name string) string {
		p := filepath.Join(dir, name+".int8.onnx")
		if _, statErr := os.Stat(p); statErr != nil {
			p = filepath.Join(dir, name+".onnx")
		}
		return p
	}
	encoder, decoder, joiner := pick("encoder"), pick("decoder"), pick("joiner")
	tokens := filepath.Join(dir, "tokens.txt")
	for _, p := range []string{encoder, decoder, joiner, tokens} {
		if _, statErr := os.Stat(p); statErr != nil {
			return nil, uiErrModel(p)
		}
	}
	return []string{
		"--encoder=" + encoder,
		"--decoder=" + decoder,
		"--joiner=" + joiner,
		"--tokens=" + tokens,
	}, nil
}

func streamInstalled(cfg *Config) bool {
	_, err := streamModelArgs(cfg.StreamModel)
	return err == nil
}

func startStreamServer(cfg *Config, logw io.Writer) (*streamServer, error) {
	s := &streamServer{
		addr:    fmt.Sprintf("127.0.0.1:%d", cfg.StreamPort),
		dir:     cfg.StreamModel,
		doneCh:  make(chan struct{}),
		threads: cfg.SherpaThreads,
	}
	modelArgs, err := streamModelArgs(cfg.StreamModel)
	if err != nil {
		return nil, err
	}
	exePath := cfg.StreamExe
	if !filepath.IsAbs(exePath) {
		if abs, aerr := filepath.Abs(exePath); aerr == nil {
			exePath = abs
		}
	}
	if _, serr := os.Stat(exePath); serr != nil {
		return nil, uiErr(fmt.Sprintf("streaming server not found: %s", exePath), trf("err.sherpa.notfound", exePath))
	}
	args := append([]string{
		"--port=" + strconv.Itoa(cfg.StreamPort),
		"--num-threads=" + strconv.Itoa(cfg.SherpaThreads),
	}, modelArgs...)
	quiet := logfilter.New(logw,
		"handle_read_frame error",
		"handle_read_handshake error",
		"asio.system:10053",
		"asio.system:10054",
		"asio.system:10058")
	cmd := exec.Command(exePath, args...)
	cmd.Stdout = quiet
	cmd.Stderr = quiet
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	if err := cmd.Start(); err != nil {
		log.Printf("starting %s: %v", cfg.StreamExe, err)
		return nil, uiErr(fmt.Sprintf("streaming server did not start: %s", cfg.StreamExe), trf("err.server.launch", cfg.StreamExe))
	}
	s.cmd = cmd
	s.job = attachProcessToJob(cmd.Process.Pid)
	go func() {
		_ = cmd.Wait()
		_ = quiet.Flush()
		s.mu.Lock()
		s.exited = true
		stopping := s.stopping
		s.mu.Unlock()
		close(s.doneCh)
		if !stopping {
			log.Printf("the streaming recognizer crashed")
		}
	}()
	return s, nil
}

func (s *streamServer) engine() string { return engineStream }

func (s *streamServer) model() string { return filepath.ToSlash(s.dir) }

func (s *streamServer) external() bool { return false }

func (s *streamServer) done() <-chan struct{} { return s.doneCh }

func (s *streamServer) wasStopped() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stopping
}

func (s *streamServer) alive() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return !s.exited
}

func (s *streamServer) stop() {
	s.mu.Lock()
	s.stopping = true
	s.mu.Unlock()
	if s.cmd != nil && s.cmd.Process != nil {
		_ = s.cmd.Process.Kill()
	}
	if s.job != 0 {
		windows.CloseHandle(s.job)
	}
}

func (s *streamServer) waitReady(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if !s.alive() {
			return uiErr("streaming server exited during startup", tr("err.sherpa.start"))
		}
		if dialOK(s.addr) {
			return nil
		}
		time.Sleep(300 * time.Millisecond)
	}
	return uiErr(fmt.Sprintf("streaming server did not answer within %s", timeout), trf("err.server.timeout", timeout))
}

func floatBytes(pcm []byte) []byte {
	n := len(pcm) / 2
	out := make([]byte, n*4)
	for i := 0; i < n; i++ {
		v := int16(binary.LittleEndian.Uint16(pcm[i*2:]))
		binary.LittleEndian.PutUint32(out[i*4:], math.Float32bits(float32(v)/32768))
	}
	return out
}

type streamResult struct {
	Text    string `json:"text"`
	Segment int    `json:"segment"`
	Final   bool   `json:"is_final"`
}

// streamSession collects the text per segment: the online server closes a
// segment on every pause and starts the next one empty, so the phrase is
// the segments joined, not the last message.
type streamSession struct {
	conn    *websocket.Conn
	mu      sync.Mutex
	segs    map[int]string
	closed  bool
	partial func(text string)
	readErr chan error
}

func (s *streamServer) openSession(ctx context.Context, partial func(text string)) (*streamSession, error) {
	if !s.alive() {
		return nil, fmt.Errorf("the streaming recognizer is not running")
	}
	dialer := websocket.Dialer{HandshakeTimeout: 10 * time.Second}
	conn, _, err := dialer.DialContext(ctx, "ws://"+s.addr, nil)
	if err != nil {
		return nil, fmt.Errorf("connecting to the streaming recognizer: %w", err)
	}
	sess := &streamSession{conn: conn, partial: partial, segs: map[int]string{}, readErr: make(chan error, 1)}
	go sess.readLoop()
	return sess, nil
}

func (ss *streamSession) readLoop() {
	for {
		_, raw, err := ss.conn.ReadMessage()
		if err != nil {
			ss.readErr <- err
			return
		}
		var res streamResult
		if json.Unmarshal(raw, &res) != nil {
			continue
		}
		if strings.TrimSpace(res.Text) == "" {
			continue
		}
		ss.mu.Lock()
		ss.segs[res.Segment] = res.Text
		joined := ss.joinedLocked()
		cb := ss.partial
		ss.mu.Unlock()
		if cb != nil {
			cb(joined)
		}
	}
}

func (ss *streamSession) joinedLocked() string {
	if len(ss.segs) == 0 {
		return ""
	}
	maxSeg := 0
	for k := range ss.segs {
		if k > maxSeg {
			maxSeg = k
		}
	}
	var parts []string
	for i := 0; i <= maxSeg; i++ {
		if t := strings.TrimSpace(ss.segs[i]); t != "" {
			parts = append(parts, t)
		}
	}
	return strings.Join(parts, " ")
}

func (ss *streamSession) push(pcm []byte) error {
	if len(pcm) < 2 {
		return nil
	}
	ss.mu.Lock()
	closed := ss.closed
	ss.mu.Unlock()
	if closed {
		return nil
	}
	_ = ss.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return ss.conn.WriteMessage(websocket.BinaryMessage, floatBytes(pcm))
}

// finish tells the server the phrase is over and waits until the text stops
// changing — the online server keeps decoding what it already has.
func (ss *streamSession) finish(ctx context.Context) (string, error) {
	ss.mu.Lock()
	ss.closed = true
	ss.mu.Unlock()
	tail := make([]byte, sampleRate/2*2)
	_ = ss.conn.WriteMessage(websocket.BinaryMessage, floatBytes(tail))
	_ = ss.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	_ = ss.conn.WriteMessage(websocket.TextMessage, []byte("Done"))
	settle := time.NewTimer(600 * time.Millisecond)
	defer settle.Stop()
	lastSeen := ss.snapshot()
	for {
		select {
		case <-ctx.Done():
			ss.close()
			return ss.snapshot(), ctx.Err()
		case <-ss.readErr:
			ss.close()
			return ss.snapshot(), nil
		case <-settle.C:
			now := ss.snapshot()
			if now == lastSeen {
				ss.close()
				return now, nil
			}
			lastSeen = now
			settle.Reset(600 * time.Millisecond)
		}
	}
}

func (ss *streamSession) snapshot() string {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return strings.Join(strings.Fields(ss.joinedLocked()), " ")
}

func (ss *streamSession) close() {
	_ = ss.conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	_ = ss.conn.Close()
}

func (s *streamServer) transcribe(ctx context.Context, wav []byte, language, prompt string, translate bool) (string, error) {
	if translate {
		return "", uiErr("this streaming model cannot translate", tr("err.sherpa.translate"))
	}
	pcm := wav
	if len(wav) > 44 && string(wav[0:4]) == "RIFF" {
		pcm = wav[44:]
	}
	sess, err := s.openSession(ctx, nil)
	if err != nil {
		return "", err
	}
	chunk := sampleRate / 5 * 2
	for off := 0; off < len(pcm); off += chunk {
		end := off + chunk
		if end > len(pcm) {
			end = len(pcm)
		}
		if err := sess.push(pcm[off:end]); err != nil {
			sess.close()
			return "", fmt.Errorf("audio streaming: %w", err)
		}
	}
	return sess.finish(ctx)
}
