package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

type whisperServer struct {
	baseURL string
	cmd     *exec.Cmd
	client  *http.Client
	done    chan struct{}
	job     windows.Handle

	mu       sync.Mutex
	exited   bool
	stopping bool
}

func (s *whisperServer) wasStopped() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stopping
}

func (s *whisperServer) external() bool { return s.cmd == nil }

func startWhisperServer(cfg *Config, logw io.Writer) (*whisperServer, error) {
	s := &whisperServer{
		client: &http.Client{Timeout: 5 * time.Minute},
		done:   make(chan struct{}),
	}
	if cfg.ServerURL != "" {
		s.baseURL = strings.TrimRight(cfg.ServerURL, "/")
		return s, nil
	}
	s.baseURL = fmt.Sprintf("http://127.0.0.1:%d", cfg.ServerPort)
	if !cfg.ServerAutostart {
		return s, nil
	}

	if _, err := os.Stat(cfg.Model); err != nil {
		return nil, fmt.Errorf("%s", trf("err.model.notfound", cfg.Model))
	}

	args := []string{
		"-m", cfg.Model,
		"-l", cfg.Language,
		"-t", strconv.Itoa(cfg.Threads),
		"--host", "127.0.0.1",
		"--port", strconv.Itoa(cfg.ServerPort),
	}
	exePath := cfg.ServerExe
	if !filepath.IsAbs(exePath) {
		if abs, err := filepath.Abs(exePath); err == nil {
			exePath = abs
		}
	}
	cmd := exec.Command(exePath, args...)
	cmd.Stdout = logw
	cmd.Stderr = logw
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000,
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("запуск %s: %w", cfg.ServerExe, err)
	}
	s.cmd = cmd
	s.attachToJob()
	go func() {
		_ = cmd.Wait()
		s.mu.Lock()
		s.exited = true
		stopping := s.stopping
		s.mu.Unlock()
		close(s.done)
		if !stopping {
			log.Printf("whisper-server аварийно завершился")
		}
	}()
	return s, nil
}

func attachProcessToJob(pid int) windows.Handle {
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		log.Printf("Job Object не создан (%v) — при аварийном выходе дочерний процесс может остаться", err)
		return 0
	}
	info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{
		BasicLimitInformation: windows.JOBOBJECT_BASIC_LIMIT_INFORMATION{
			LimitFlags: windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE,
		},
	}
	if _, err := windows.SetInformationJobObject(job, windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)), uint32(unsafe.Sizeof(info))); err != nil {
		log.Printf("Job Object не настроен: %v", err)
		windows.CloseHandle(job)
		return 0
	}
	h, err := windows.OpenProcess(windows.PROCESS_SET_QUOTA|windows.PROCESS_TERMINATE, false, uint32(pid))
	if err != nil {
		log.Printf("процесс не открыт для Job Object: %v", err)
		windows.CloseHandle(job)
		return 0
	}
	defer windows.CloseHandle(h)
	if err := windows.AssignProcessToJobObject(job, h); err != nil {
		log.Printf("процесс не привязан к Job Object: %v", err)
		windows.CloseHandle(job)
		return 0
	}
	return job
}

func (s *whisperServer) attachToJob() {
	s.job = attachProcessToJob(s.cmd.Process.Pid)
}

func (s *whisperServer) hostPort() string {
	u, err := url.Parse(s.baseURL)
	if err != nil || u.Host == "" {
		return strings.TrimPrefix(s.baseURL, "http://")
	}
	host := u.Host
	if u.Port() == "" {
		switch u.Scheme {
		case "https":
			host += ":443"
		default:
			host += ":80"
		}
	}
	return host
}

func (s *whisperServer) waitReady(timeout time.Duration) error {
	addr := s.hostPort()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if !s.external() {
			s.mu.Lock()
			dead := s.exited
			s.mu.Unlock()
			if dead {
				return fmt.Errorf("%s", tr("err.server.start"))
			}
		}
		conn, err := net.DialTimeout("tcp", addr, time.Second)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(300 * time.Millisecond)
	}
	if s.external() {
		return fmt.Errorf("%s", trf("err.server.dead", s.baseURL))
	}
	return fmt.Errorf("%s", trf("err.server.timeout", timeout))
}

func (s *whisperServer) stop() {
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

func (s *whisperServer) transcribe(ctx context.Context, wav []byte, language, prompt string, translate bool) (string, error) {
	s.mu.Lock()
	dead := s.exited
	s.mu.Unlock()
	if dead {
		return "", fmt.Errorf("whisper-server не запущен")
	}

	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	fw, err := mw.CreateFormFile("file", "audio.wav")
	if err != nil {
		return "", err
	}
	if _, err := fw.Write(wav); err != nil {
		return "", err
	}
	_ = mw.WriteField("response_format", "json")
	_ = mw.WriteField("temperature", "0.0")
	if language != "" {
		_ = mw.WriteField("language", language)
	}
	if prompt != "" {
		_ = mw.WriteField("prompt", prompt)
	}
	if translate {
		_ = mw.WriteField("translate", "true")
	}
	if err := mw.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/inference", &body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("запрос к whisper-server: %w", err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var parsed struct {
		Text  string `json:"text"`
		Error string `json:"error"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("неожиданный ответ сервера (%d): %.200s", resp.StatusCode, string(raw))
	}
	if parsed.Error != "" {
		return "", fmt.Errorf("whisper-server: %s", parsed.Error)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("whisper-server: HTTP %d: %.200s", resp.StatusCode, string(raw))
	}
	return strings.TrimSpace(parsed.Text), nil
}
