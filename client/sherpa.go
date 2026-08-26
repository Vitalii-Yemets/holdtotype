package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	"holdtotype/internal/sherpaproto"
)

type sherpaServer struct {
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

func sherpaModelFiles(dir string) (encoder, decoder, joiner, tokens string, err error) {
	pick := func(name string) string {
		p := filepath.Join(dir, name+".int8.onnx")
		if _, statErr := os.Stat(p); statErr != nil {
			p = filepath.Join(dir, name+".onnx")
		}
		return p
	}
	encoder = pick("encoder")
	decoder = pick("decoder")
	joiner = pick("joiner")
	tokens = filepath.Join(dir, "tokens.txt")
	for _, p := range []string{encoder, decoder, joiner, tokens} {
		if _, statErr := os.Stat(p); statErr != nil {
			return "", "", "", "", fmt.Errorf("%s", trf("err.sherpa.model", p))
		}
	}
	return encoder, decoder, joiner, tokens, nil
}

func startSherpaServer(cfg *Config, logw io.Writer) (*sherpaServer, error) {
	s := &sherpaServer{
		addr:    fmt.Sprintf("127.0.0.1:%d", cfg.SherpaPort),
		dir:     cfg.SherpaModel,
		doneCh:  make(chan struct{}),
		threads: cfg.SherpaThreads,
	}
	encoder, decoder, joiner, tokens, err := sherpaModelFiles(cfg.SherpaModel)
	if err != nil {
		return nil, err
	}
	exePath := cfg.SherpaExe
	if !filepath.IsAbs(exePath) {
		if abs, aerr := filepath.Abs(exePath); aerr == nil {
			exePath = abs
		}
	}
	if _, serr := os.Stat(exePath); serr != nil {
		return nil, fmt.Errorf("%s", trf("err.sherpa.notfound", exePath))
	}
	args := []string{
		"--port=" + strconv.Itoa(cfg.SherpaPort),
		"--num-work-threads=2",
		"--num-io-threads=1",
		"--encoder=" + encoder,
		"--decoder=" + decoder,
		"--joiner=" + joiner,
		"--tokens=" + tokens,
		"--model-type=nemo_transducer",
		"--num-threads=" + strconv.Itoa(cfg.SherpaThreads),
	}
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
		log.Printf("запуск %s: %v", cfg.SherpaExe, err)
		return nil, fmt.Errorf("%s", trf("err.server.launch", cfg.SherpaExe))
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
			log.Printf("sherpa-server аварийно завершился")
		}
	}()
	return s, nil
}

func (s *sherpaServer) engine() string { return engineSherpa }

func (s *sherpaServer) external() bool { return false }

func (s *sherpaServer) done() <-chan struct{} { return s.doneCh }

func (s *sherpaServer) wasStopped() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stopping
}

func (s *sherpaServer) alive() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return !s.exited
}

func (s *sherpaServer) probe() (bool, error) {
	if !dialOK(s.addr) {
		return false, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	silence := make([]byte, sampleRate/5*2)
	if _, err := s.decode(ctx, silence, sampleRate); err != nil {
		return false, err
	}
	return true, nil
}

func (s *sherpaServer) waitReady(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		if !s.alive() {
			return fmt.Errorf("%s", tr("err.sherpa.start"))
		}
		ok, err := s.probe()
		if ok {
			if lastErr != nil {
				log.Printf("sherpa-server ответил после ошибки: %v", lastErr)
			}
			return nil
		}
		lastErr = err
		time.Sleep(300 * time.Millisecond)
	}
	if lastErr != nil {
		return fmt.Errorf("%s: %v", trf("err.server.timeout", timeout), lastErr)
	}
	return fmt.Errorf("%s", trf("err.server.timeout", timeout))
}

func (s *sherpaServer) stop() {
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

func (s *sherpaServer) transcribe(ctx context.Context, wav []byte, language, prompt string, translate bool) (string, error) {
	if !s.alive() {
		return "", fmt.Errorf("sherpa-server не запущен")
	}
	if translate {
		return "", fmt.Errorf("%s", tr("err.sherpa.translate"))
	}
	pcm, rate, err := sherpaproto.PCMFromWAV(wav)
	if err != nil {
		return "", err
	}
	if rate <= 0 {
		rate = sampleRate
	}
	return s.decode(ctx, pcm, rate)
}

func (s *sherpaServer) decode(ctx context.Context, pcm []byte, rate int) (string, error) {
	samples := sherpaproto.FromPCM16(pcm)

	dialer := websocket.Dialer{HandshakeTimeout: 10 * time.Second}
	conn, _, err := dialer.DialContext(ctx, "ws://"+s.addr, nil)
	if err != nil {
		return "", fmt.Errorf("подключение к sherpa-server: %w", err)
	}
	defer conn.Close()

	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetWriteDeadline(deadline)
		_ = conn.SetReadDeadline(deadline)
	} else {
		_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Minute))
		_ = conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
	}

	payload := append(sherpaproto.Header(rate, len(samples)), samples...)
	for _, part := range sherpaproto.Chunks(payload, sherpaproto.ChunkBytes) {
		if err := conn.WriteMessage(websocket.BinaryMessage, part); err != nil {
			return "", fmt.Errorf("передача звука: %w", err)
		}
	}

	_, raw, err := conn.ReadMessage()
	if err != nil {
		return "", fmt.Errorf("ответ sherpa-server: %w", err)
	}
	_ = conn.WriteMessage(websocket.TextMessage, []byte("Done"))
	_ = conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	for {
		if _, _, rerr := conn.ReadMessage(); rerr != nil {
			break
		}
	}

	text := string(raw)
	if len(raw) > 0 && raw[0] == '{' {
		var parsed struct {
			Text string `json:"text"`
		}
		if json.Unmarshal(raw, &parsed) == nil {
			text = parsed.Text
		}
	}
	return strings.TrimSpace(text), nil
}
