package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

	"holdtotype/internal/appid"
)

var llmAPIKey = func() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return appid.LLMAlias
	}
	return hex.EncodeToString(b)
}()

const maxResponseBytes = 8 << 20
const llmFile = "qwen2.5-1.5b-instruct-q4_k_m.gguf"

func llmInstalled(cfg *Config) bool {
	_, err := os.Stat(cfg.LLMModel)
	return err == nil
}

type llamaServer struct {
	baseURL string
	cmd     *exec.Cmd
	client  *http.Client
	job     uintptr

	mu     sync.Mutex
	exited bool
}

func startLlamaServer(cfg *Config, logw io.Writer) (*llamaServer, error) {
	s := &llamaServer{
		baseURL: fmt.Sprintf("http://127.0.0.1:%d", cfg.LLMPort),
		client:  &http.Client{Timeout: 3 * time.Minute},
	}
	exePath := cfg.LLMExe
	if !filepath.IsAbs(exePath) {
		if abs, err := filepath.Abs(exePath); err == nil {
			exePath = abs
		}
	}
	args := []string{
		"-m", cfg.LLMModel,
		"--host", "127.0.0.1",
		"--port", strconv.Itoa(cfg.LLMPort),
		"-t", strconv.Itoa(cfg.Threads),
		"-c", "4096",
		"--api-key", llmAPIKey,
	}
	cmd := exec.Command(exePath, args...)
	cmd.Stdout = logw
	cmd.Stderr = logw
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000,
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("запуск %s: %w", cfg.LLMExe, err)
	}
	s.cmd = cmd
	s.job = uintptr(attachProcessToJob(cmd.Process.Pid))
	go func() {
		_ = cmd.Wait()
		s.mu.Lock()
		s.exited = true
		s.mu.Unlock()
	}()
	return s, nil
}

func (s *llamaServer) alive() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return !s.exited
}

var healthClient = &http.Client{Timeout: 2 * time.Second}

func (s *llamaServer) waitReady(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if !s.alive() {
			return fmt.Errorf("llama-server завершился при старте (см. лог)")
		}
		resp, err := healthClient.Get(s.baseURL + "/health")
		if err == nil {
			ok := resp.StatusCode == http.StatusOK
			resp.Body.Close()
			if ok {
				return nil
			}
		}
		time.Sleep(400 * time.Millisecond)
	}
	return fmt.Errorf("llama-server не ответил за %s", timeout)
}

func (s *llamaServer) stop() {
	if s.cmd != nil && s.cmd.Process != nil {
		_ = s.cmd.Process.Kill()
	}
	if s.job != 0 {
		_ = windows.CloseHandle(windows.Handle(s.job))
		s.job = 0
	}
}

func (s *llamaServer) chat(ctx context.Context, system, user string) (string, error) {
	payload := map[string]any{
		"messages": []map[string]string{
			{"role": "system", "content": system},
			{"role": "user", "content": user},
		},
		"temperature": 0.2,
		"max_tokens":  2048,
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+llmAPIKey)
	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	if err != nil {
		return "", err
	}
	var parsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("llama-server (%d): %.200s", resp.StatusCode, string(raw))
	}
	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("llama-server: пустой ответ (%d): %.200s", resp.StatusCode, string(raw))
	}
	return parsed.Choices[0].Message.Content, nil
}

var (
	procGlobalMemoryStatusEx = kernel32.NewProc("GlobalMemoryStatusEx")
)

type memStatusEx struct {
	Length               uint32
	MemoryLoad           uint32
	TotalPhys            uint64
	AvailPhys            uint64
	TotalPageFile        uint64
	AvailPageFile        uint64
	TotalVirtual         uint64
	AvailVirtual         uint64
	AvailExtendedVirtual uint64
}

func ramMB() (total, avail int) {
	var ms memStatusEx
	ms.Length = uint32(unsafe.Sizeof(ms))
	procGlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&ms)))
	return int(ms.TotalPhys / (1024 * 1024)), int(ms.AvailPhys / (1024 * 1024))
}

func llmNeedMB(sizeMB int) int {
	return sizeMB + sizeMB/6 + 1500
}

func llmFit(sizeMB int) string {
	_, avail := ramMB()
	need := llmNeedMB(sizeMB)
	switch {
	case need < avail*7/10:
		return "ok"
	case need < avail*95/100:
		return "warn"
	default:
		return "bad"
	}
}

func (a *App) llmStatus() string {
	activeFile := filepath.Base(a.snapshot().LLMModel)
	type instRow struct {
		File   string `json:"file"`
		SizeMB int    `json:"size"`
		Active bool   `json:"active"`
	}
	var installed []instRow
	entries, _ := os.ReadDir("models")
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".gguf") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		installed = append(installed, instRow{File: name, SizeMB: int(info.Size() / (1024 * 1024)), Active: name == activeFile})
	}
	type dlRow struct {
		File string `json:"file"`
		Pct  int    `json:"pct"`
		Err  string `json:"err"`
	}
	var downloads []dlRow
	dlMu.Lock()
	for k, st := range dl {
		if !strings.HasPrefix(k, "llm-") {
			continue
		}
		if st.active {
			downloads = append(downloads, dlRow{File: strings.TrimPrefix(k, "llm-"), Pct: st.pct})
		} else if st.err != "" {
			downloads = append(downloads, dlRow{File: strings.TrimPrefix(k, "llm-"), Pct: -1, Err: st.err})
		}
	}
	dlMu.Unlock()
	total, avail := ramMB()
	out, _ := json.Marshal(map[string]any{
		"installed": installed,
		"downloads": downloads,
		"ram":       total,
		"ram_free":  avail,
	})
	return string(out)
}

func (a *App) llmDownloadFile(repo, file string) {
	if !repoOK(repo) || strings.Contains(file, "/") || strings.Contains(file, "\\") || !strings.HasSuffix(file, ".gguf") {
		return
	}
	u := "https://huggingface.co/" + repo + "/resolve/main/" + file
	a.startDownload("llm-"+file, file, u)
}

func (a *App) llmDelete(file string) string {
	if strings.Contains(file, "/") || strings.Contains(file, "\\") || !strings.HasSuffix(file, ".gguf") {
		return ""
	}
	cfg := a.snapshot()
	active := filepath.Base(cfg.LLMModel) == file
	if active {
		a.llmShutdown()
	}
	path := filepath.Join("models", file)
	var err error
	for i := 0; i < 10; i++ {
		if err = os.Remove(path); err == nil || os.IsNotExist(err) {
			err = nil
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	if err != nil {
		return err.Error()
	}
	if active {
		a.mu.Lock()
		c := *a.cfg
		c.LLMModel = ""
		a.cfg = &c
		a.mu.Unlock()
		if serr := saveConfig("config.json", &c); serr != nil {
			log.Printf("сохранение конфига: %v", serr)
		}
	}
	log.Printf("LLM-модель %s удалена", file)
	return tr("model.del.ok")
}

func repoOK(repo string) bool {
	if repo == "" || len(repo) > 200 || strings.Count(repo, "/") != 1 || strings.Contains(repo, "..") {
		return false
	}
	for _, r := range repo {
		ok := r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' ||
			r == '-' || r == '_' || r == '.' || r == '/'
		if !ok {
			return false
		}
	}
	return true
}

func (a *App) llmSearch(q string) string {
	q = strings.TrimSpace(q)
	if q == "" {
		return `{"repos":[]}`
	}
	client := &http.Client{Timeout: 15 * time.Second}
	u := "https://huggingface.co/api/models?search=" + url.QueryEscape(q) + "&filter=gguf&sort=downloads&direction=-1&limit=8&expand[]=lastModified&expand[]=downloads"
	resp, err := client.Get(u)
	if err != nil {
		out, _ := json.Marshal(map[string]any{"error": err.Error()})
		return string(out)
	}
	defer resp.Body.Close()
	var raw []struct {
		ID           string `json:"id"`
		Downloads    int    `json:"downloads"`
		LastModified string `json:"lastModified"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		out, _ := json.Marshal(map[string]any{"error": err.Error()})
		return string(out)
	}
	type repoRow struct {
		ID        string `json:"id"`
		Downloads int    `json:"downloads"`
		Updated   string `json:"updated"`
	}
	repos := make([]repoRow, 0, len(raw))
	for _, r := range raw {
		d := r.LastModified
		if len(d) > 10 {
			d = d[:10]
		}
		repos = append(repos, repoRow{ID: r.ID, Downloads: r.Downloads, Updated: d})
	}
	out, _ := json.Marshal(map[string]any{"repos": repos})
	return string(out)
}

func (a *App) llmRepoFiles(repo string) string {
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get("https://huggingface.co/api/models/" + repo + "?blobs=true")
	if err != nil {
		out, _ := json.Marshal(map[string]any{"error": err.Error()})
		return string(out)
	}
	defer resp.Body.Close()
	var raw struct {
		Siblings []struct {
			Rfilename string `json:"rfilename"`
			Size      int64  `json:"size"`
		} `json:"siblings"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		out, _ := json.Marshal(map[string]any{"error": err.Error()})
		return string(out)
	}
	type fileRow struct {
		File   string `json:"file"`
		SizeMB int    `json:"size"`
		Fit    string `json:"fit"`
		NeedMB int    `json:"need"`
	}
	var files []fileRow
	for _, s := range raw.Siblings {
		if !strings.HasSuffix(s.Rfilename, ".gguf") || strings.Contains(s.Rfilename, "/") {
			continue
		}
		mb := int(s.Size / (1024 * 1024))
		files = append(files, fileRow{File: s.Rfilename, SizeMB: mb, Fit: llmFit(mb), NeedMB: llmNeedMB(mb)})
		if len(files) >= 14 {
			break
		}
	}
	out, _ := json.Marshal(map[string]any{"files": files})
	return string(out)
}

func (a *App) llmShutdown() {
	a.mu.Lock()
	llm := a.llm
	a.llm = nil
	a.mu.Unlock()
	if llm != nil {
		llm.stop()
	}
}

var llmStartMu sync.Mutex

func (a *App) ensureLLM() (*llamaServer, error) {
	cfg := a.snapshot()
	if !llmInstalled(cfg) {
		return nil, fmt.Errorf("LLM-модель не установлена")
	}
	a.mu.Lock()
	llm := a.llm
	a.mu.Unlock()
	if llm != nil && llm.alive() {
		return llm, nil
	}

	llmStartMu.Lock()
	defer llmStartMu.Unlock()

	a.mu.Lock()
	llm = a.llm
	a.mu.Unlock()
	if llm != nil && llm.alive() {
		return llm, nil
	}
	if llm != nil {
		llm.stop()
	}
	log.Printf("запускаю llama-server")
	llm, err := startLlamaServer(cfg, logFile)
	if err != nil {
		return nil, err
	}
	if err := llm.waitReady(2 * time.Minute); err != nil {
		llm.stop()
		return nil, err
	}
	a.mu.Lock()
	if a.quitting {
		a.mu.Unlock()
		llm.stop()
		return nil, fmt.Errorf("завершение приложения")
	}
	a.llm = llm
	a.mu.Unlock()
	log.Printf("llama-server готов")
	return llm, nil
}
