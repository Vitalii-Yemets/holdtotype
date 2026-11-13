package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
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
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"holdtotype/internal/appid"
	"holdtotype/internal/chatstream"
	"holdtotype/internal/editormodel"
	"holdtotype/internal/subtitles"
	"holdtotype/internal/trmodel"
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
	label   string
	baseURL string
	cmd     *exec.Cmd
	client  *http.Client
	job     uintptr

	mu     sync.Mutex
	exited bool

	primeMu sync.Mutex
	primed  atomic.Bool
}

func (s *llamaServer) prime(to string) error {
	s.primeMu.Lock()
	defer s.primeMu.Unlock()
	if s.primed.Load() {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), subtitles.TranslateWarmWait)
	defer cancel()
	for _, msgs := range subtitles.WarmMessages(to) {
		if _, err := s.chatStream(ctx, msgs, 1, "", nil); err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				return fmt.Errorf("%s %w: the first requests took more than %.0f s", s.label, errEngineTimeout, subtitles.TranslateWarmWait.Seconds())
			}
			return err
		}
	}
	s.primed.Store(true)
	return nil
}

func (a *App) translatorWarm(kind string) bool {
	var p engineProc
	switch kind {
	case editorFast:
		return a.fastWarm()
	case editorQuality:
		p = a.gpuEng.live()
	default:
		return true
	}
	s, ok := p.(*llamaServer)
	return ok && s != nil && s.primed.Load()
}

func needsStart(kind string) bool {
	return kind == editorQuality || kind == editorFast
}

func asLlama(p engineProc, err error) (*llamaServer, error) {
	if err != nil {
		return nil, err
	}
	s, ok := p.(*llamaServer)
	if !ok || s == nil {
		return nil, errors.New("the engine is not a language model server")
	}
	return s, nil
}

func (a *App) warmTranslator(kind, to string) error {
	var f *llmFlight
	switch kind {
	case editorQuality:
		f = a.qualityWarmUp()
	case editorFast:
		ctx, cancel := context.WithTimeout(context.Background(), fastLangWait+subtitles.TranslateWarmWait)
		defer cancel()
		return a.warmFast(ctx)
	default:
		return nil
	}
	<-f.done
	s, err := asLlama(f.srv, f.err)
	if err != nil {
		return err
	}
	return s.prime(to)
}

type engineProc interface {
	alive() bool
	stop()
	waitReady(time.Duration) error
}

type engineSpec struct {
	label  string
	exe    string
	port   int
	args   []string
	wait   time.Duration
	launch func(engineSpec, io.Writer) (engineProc, error)
}

var (
	errEngineExited  = errors.New("exited during startup (see the log)")
	errEngineTimeout = errors.New("did not answer in time")
)

func absPath(p string) string {
	if !filepath.IsAbs(p) {
		if abs, err := filepath.Abs(p); err == nil {
			return abs
		}
	}
	return p
}

func editorEngine(cfg *Config) (engineSpec, error) {
	if !llmInstalled(cfg) {
		return engineSpec{}, fmt.Errorf("no LLM model is installed")
	}
	return engineSpec{
		label: "llama-server",
		exe:   absPath(cfg.LLMExe),
		port:  cfg.LLMPort,
		args:  []string{"-m", cfg.LLMModel, "-t", strconv.Itoa(cfg.Threads), "-c", "4096"},
		wait:  2 * time.Minute,
	}, nil
}

func startLlamaServer(spec engineSpec, logw io.Writer) (*llamaServer, error) {
	s := &llamaServer{
		label:   spec.label,
		baseURL: fmt.Sprintf("http://127.0.0.1:%d", spec.port),
		client:  &http.Client{Timeout: 3 * time.Minute},
	}
	if !waitPortFree(spec.port, portFreeWait) {
		return nil, uiErr(fmt.Sprintf("%s: port %d is taken by another program", spec.label, spec.port), trf("err.port.busy", spec.port))
	}
	args := append(append([]string(nil), spec.args...),
		"--host", "127.0.0.1",
		"--port", strconv.Itoa(spec.port),
		"--api-key", llmAPIKey,
	)
	cmd := exec.Command(spec.exe, args...)
	cmd.Stdout = logw
	cmd.Stderr = logw
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000,
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("starting %s: %w", filepath.Base(spec.exe), err)
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
			return fmt.Errorf("%s %w", s.label, errEngineExited)
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
	return fmt.Errorf("%s %w: nothing within %s", s.label, errEngineTimeout, timeout)
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
	return s.chatWith(ctx, system, user, 2048, "")
}

func (s *llamaServer) chatWith(ctx context.Context, system, user string, maxTokens int, grammar string) (string, error) {
	payload := map[string]any{
		"messages": []map[string]string{
			{"role": "system", "content": system},
			{"role": "user", "content": user},
		},
		"temperature": 0.2,
		"max_tokens":  maxTokens,
	}
	if grammar != "" {
		payload["grammar"] = grammar
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
		return "", fmt.Errorf("llama-server: empty answer (%d): %.200s", resp.StatusCode, string(raw))
	}
	return parsed.Choices[0].Message.Content, nil
}

func (s *llamaServer) chatStream(ctx context.Context, msgs []subtitles.Message, maxTokens int, grammar string, onText func(string)) (string, error) {
	payload := map[string]any{
		"messages":    msgs,
		"temperature": 0.2,
		"max_tokens":  maxTokens,
		"stream":      true,
	}
	if grammar != "" {
		payload["grammar"] = grammar
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Authorization", "Bearer "+llmAPIKey)
	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return "", fmt.Errorf("llama-server (%d): %.200s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}
	out, err := chatstream.Read(io.LimitReader(resp.Body, maxResponseBytes), onText)
	if err != nil {
		return out, fmt.Errorf("llama-server: %w", err)
	}
	return out, nil
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
		Loaded bool   `json:"loaded"`
	}
	warm := a.editorEng.live() != nil
	var installed []instRow
	for _, f := range editormodel.List("models") {
		installed = append(installed, instRow{File: f.Name, SizeMB: f.SizeMB, Active: f.Name == activeFile, Loaded: warm && f.Name == activeFile})
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
	cfg := a.snapshot()
	out, _ := json.Marshal(map[string]any{
		"installed": installed,
		"downloads": downloads,
		"ram":       total,
		"ram_free":  avail,
		"external":  postAPIOn(cfg),
	})
	return string(out)
}

func (a *App) llmDownloadFile(repo, file string, sizeMB int) {
	if !repoOK(repo) || strings.Contains(file, "/") || strings.Contains(file, "\\") || !strings.HasSuffix(file, ".gguf") || strings.HasPrefix(file, "mmproj-") || trmodel.Is(file) {
		return
	}
	u := "https://huggingface.co/" + repo + "/resolve/main/" + file
	a.startDownload("llm-"+file, file, u, sizeMB)
}

func (a *App) llmDelete(file string) string {
	if strings.Contains(file, "/") || strings.Contains(file, "\\") || !strings.HasSuffix(file, ".gguf") || trmodel.Is(file) {
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
		return humanError(err)
	}
	if active {
		a.mu.Lock()
		c := *a.cfg
		c.LLMModel = ""
		a.cfg = &c
		a.mu.Unlock()
		if serr := saveConfig("config.json", &c); serr != nil {
			log.Printf("saving the config: %v", serr)
		}
	}
	log.Printf("LLM model %s deleted", file)
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
		out, _ := json.Marshal(map[string]any{"error": humanError(err)})
		return string(out)
	}
	defer resp.Body.Close()
	var raw []struct {
		ID           string `json:"id"`
		Downloads    int    `json:"downloads"`
		LastModified string `json:"lastModified"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		out, _ := json.Marshal(map[string]any{"error": humanError(err)})
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
		out, _ := json.Marshal(map[string]any{"error": humanError(err)})
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
		out, _ := json.Marshal(map[string]any{"error": humanError(err)})
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
		if !strings.HasSuffix(s.Rfilename, ".gguf") || strings.Contains(s.Rfilename, "/") || strings.HasPrefix(s.Rfilename, "mmproj-") || trmodel.Is(s.Rfilename) {
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

type engineSlot struct {
	mu       sync.Mutex
	cur      engineProc
	startMu  sync.Mutex
	flightMu sync.Mutex
	flight   *llmFlight
}

func (e *engineSlot) live() engineProc {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.cur != nil && e.cur.alive() {
		return e.cur
	}
	return nil
}

func (e *engineSlot) shutdown() bool {
	e.mu.Lock()
	s := e.cur
	e.cur = nil
	e.mu.Unlock()
	if s != nil {
		s.stop()
	}
	return s != nil
}

func (a *App) ensureEngine(e *engineSlot, prepare func(*Config) (engineSpec, error)) (engineProc, error) {
	if s := e.live(); s != nil {
		return s, nil
	}
	e.startMu.Lock()
	defer e.startMu.Unlock()
	if s := e.live(); s != nil {
		return s, nil
	}
	e.shutdown()
	spec, err := prepare(a.snapshot())
	if err != nil {
		return nil, err
	}
	log.Printf("starting %s", spec.label)
	launch := spec.launch
	if launch == nil {
		launch = func(spec engineSpec, w io.Writer) (engineProc, error) { return startLlamaServer(spec, w) }
	}
	s, err := launch(spec, logFile)
	if err != nil {
		return nil, err
	}
	if err := s.waitReady(spec.wait); err != nil {
		s.stop()
		return nil, err
	}
	a.mu.Lock()
	if a.quitting {
		a.mu.Unlock()
		s.stop()
		return nil, fmt.Errorf("the application is shutting down")
	}
	e.mu.Lock()
	e.cur = s
	e.mu.Unlock()
	a.mu.Unlock()
	log.Printf("%s ready", spec.label)
	return s, nil
}

type llmFlight struct {
	done chan struct{}
	srv  engineProc
	err  error
}

func (a *App) warmEngine(e *engineSlot, prepare func(*Config) (engineSpec, error)) *llmFlight {
	e.flightMu.Lock()
	defer e.flightMu.Unlock()
	if f := e.flight; f != nil {
		return f
	}
	f := &llmFlight{done: make(chan struct{})}
	e.flight = f
	go func() {
		f.srv, f.err = a.ensureEngine(e, prepare)
		e.flightMu.Lock()
		e.flight = nil
		e.flightMu.Unlock()
		close(f.done)
	}()
	return f
}

func (a *App) waitEngine(ctx context.Context, e *engineSlot, prepare func(*Config) (engineSpec, error), starting string) (engineProc, error) {
	if s := e.live(); s != nil {
		return s, nil
	}
	f := a.warmEngine(e, prepare)
	select {
	case <-f.done:
		return f.srv, f.err
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			return nil, errors.New(starting)
		}
		return nil, ctx.Err()
	}
}

func (a *App) llmShutdown() { a.editorEng.shutdown() }

func (a *App) ensureLLM() (*llamaServer, error) {
	return asLlama(a.ensureEngine(&a.editorEng, editorEngine))
}

const (
	editorAPI     = subtitles.KindServer
	editorQuality = subtitles.KindQuality
	editorFast    = subtitles.KindFast
)

type editorSource struct {
	kind    string
	name    string
	missing string
	detail  string
}

func (s editorSource) ok() bool { return s.kind != "" && s.missing == "" }

func (s editorSource) label() string {
	switch s.kind {
	case editorAPI:
		return s.name + " (post-processing server)"
	case editorQuality:
		return "the built-in translation model (graphics card " + s.name + ")"
	case editorFast:
		return "the light translation engine (" + s.name + ")"
	}
	return "none"
}

func translationSource(cfg *Config) editorSource {
	card, fits := qualityCard()
	fast := fastAvail(cfg)
	pick := subtitles.ChooseTranslator(cfg.SubtitlesTranslator, subtitles.Avail{
		ModelReady:  trModelReady(),
		CardFits:    fits,
		QualityDown: qualityDown() != nil,
		Server:      postAPISet(cfg),
		Fast:        fast.ok,
		FastMissing: fast.missing,
	})
	src := editorSource{kind: pick.Kind, missing: pick.Missing}
	switch pick.Kind {
	case editorAPI:
		src.name = postAPIName(cfg)
	case editorQuality:
		src.name = card.Name
	case editorFast:
		src.name = routeName(fast.lang, cfg.SubtitlesTranslate)
		src.detail = fast.detail
	}
	return src
}

var errNoTranslator = errors.New("this translator does not translate by a language model")

func (a *App) editorStream(ctx context.Context, cfg *Config, src editorSource, msgs []subtitles.Message, maxTokens int, onText func(string)) (string, error) {
	switch src.kind {
	case editorAPI:
		return externalChatStream(ctx, cfg, msgs, maxTokens, onText)
	case editorQuality:
		srv, err := a.qualityWait(ctx)
		if err != nil {
			return "", err
		}
		return srv.chatStream(ctx, msgs, maxTokens, "", onText)
	}
	return "", errNoTranslator
}
