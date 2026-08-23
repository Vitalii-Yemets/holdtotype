package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sys/windows"

	"holdtotype/internal/advisor"
)

type modelInfo struct {
	ID        string
	File      string
	SizeMB    int
	NameKey   string
	DescKey   string
	Engine    string
	Dir       string
	Files     []string
	BaseURL   string
	Langs     string
	Punct     bool
	Translate bool
	Speed     int
	Accuracy  int
}

func (m *modelInfo) ramEstimateMB() int {
	if m.Engine == engineSherpa {
		return m.SizeMB * 12 / 10
	}
	return m.SizeMB*15/10 + 60
}

func advisorCatalog() []advisor.Model {
	out := make([]advisor.Model, 0, len(modelCatalog))
	for i := range modelCatalog {
		m := &modelCatalog[i]
		langs := strings.Split(m.Langs, ",")
		out = append(out, advisor.Model{
			ID: m.ID, Engine: m.Engine, Langs: langs, SizeMB: m.SizeMB,
			RAMMB: m.ramEstimateMB(), Punct: m.Punct, Translate: m.Translate,
			Speed: m.Speed, Accuracy: m.Accuracy,
		})
	}
	return out
}

var modelCatalog = []modelInfo{
	{ID: "base", File: "ggml-base.bin", SizeMB: 142, NameKey: "Base", DescKey: "S_M_BASE",
		Engine: engineWhisper, Langs: "*", Translate: true, Speed: 5, Accuracy: 2},
	{ID: "small", File: "ggml-small.bin", SizeMB: 466, NameKey: "Small", DescKey: "S_M_SMALL",
		Engine: engineWhisper, Langs: "*", Translate: true, Speed: 3, Accuracy: 3},
	{ID: "medium-q5_0", File: "ggml-medium-q5_0.bin", SizeMB: 539, NameKey: "Medium (q5)", DescKey: "S_M_MED",
		Engine: engineWhisper, Langs: "*", Translate: true, Speed: 2, Accuracy: 4},
	{ID: "large-v3-turbo-q5_0", File: "ggml-large-v3-turbo-q5_0.bin", SizeMB: 574, NameKey: "Turbo (q5)", DescKey: "S_M_TURBO",
		Engine: engineWhisper, Langs: "*", Translate: false, Speed: 4, Accuracy: 4},
	{ID: "gigaam-v3", SizeMB: 232, NameKey: "GigaAM v3", DescKey: "S_M_GIGAAM",
		Engine: engineSherpa, Dir: "gigaam-v3", Langs: "ru", Punct: true, Speed: 5, Accuracy: 5,
		Files:   []string{"encoder.int8.onnx", "decoder.onnx", "joiner.onnx", "tokens.txt"},
		BaseURL: "https://huggingface.co/csukuangfj/sherpa-onnx-nemo-transducer-punct-giga-am-v3-russian-2025-12-16/resolve/main/"},
}

const modelBaseURL = "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/"

func (m *modelInfo) paths() []string {
	if m.Engine == engineSherpa {
		out := make([]string, 0, len(m.Files))
		for _, f := range m.Files {
			out = append(out, filepath.Join("models", m.Dir, f))
		}
		return out
	}
	return []string{filepath.Join("models", m.File)}
}

func (m *modelInfo) installed() bool {
	for _, p := range m.paths() {
		if _, err := os.Stat(p); err != nil {
			return false
		}
	}
	return true
}

func (m *modelInfo) slotMatches(cfg *Config) bool {
	if m.Engine == engineSherpa {
		return filepath.Base(filepath.Clean(cfg.SherpaModel)) == m.Dir
	}
	return filepath.Base(cfg.Model) == m.File
}

func (m *modelInfo) isActive(cfg *Config) bool {
	return m.slotMatches(cfg) && m.Engine == primaryEngine(cfg)
}

type dlState struct {
	active bool
	pct    int
	err    string
	cancel context.CancelFunc
}

const partKeep = 7 * 24 * time.Hour

var (
	dlMu sync.Mutex
	dl   = map[string]*dlState{}
)

func findModel(id string) *modelInfo {
	for i := range modelCatalog {
		if modelCatalog[i].ID == id {
			return &modelCatalog[i]
		}
	}
	return nil
}

type modelRow struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Size   int    `json:"size"`
	State  string `json:"state"`
	Pct    int    `json:"pct"`
	Err    string `json:"err"`
	Engine string `json:"engine"`
	Langs  string `json:"langs"`
	Punct  bool   `json:"punct"`
	Trans  bool   `json:"translate"`
	RAM    int    `json:"ram"`
	Fit    string `json:"fit"`
}

func ramFit(needMB, freeMB int) string {
	if freeMB <= 0 || needMB <= 0 {
		return ""
	}
	switch {
	case needMB*10 <= freeMB*6:
		return "ok"
	case needMB*10 <= freeMB*9:
		return "warn"
	default:
		return "bad"
	}
}

func (a *App) modelRows() string {
	cfg := a.snapshot()
	_, freeRAM := ramMB()
	var rows []modelRow
	known := false
	for i := range modelCatalog {
		m := &modelCatalog[i]
		row := modelRow{
			ID: m.ID, Name: m.NameKey, Desc: strS(m.DescKey), Size: m.SizeMB,
			Engine: m.Engine, Langs: m.Langs, Punct: m.Punct,
			Trans: m.Translate, RAM: m.ramEstimateMB(), Fit: ramFit(m.ramEstimateMB(), freeRAM),
		}
		have := m.installed()
		active := m.isActive(cfg)
		dlMu.Lock()
		st := dl[m.ID]
		dlMu.Unlock()
		switch {
		case st != nil && st.active:
			row.State = "downloading"
			row.Pct = st.pct
		case active && have:
			row.State = "active"
			known = true
		case have:
			row.State = "installed"
		default:
			if active {
				known = true
			}
			row.State = "absent"
			if st != nil && st.err != "" {
				row.Err = st.err
			}
		}
		rows = append(rows, row)
	}
	if !known {
		rows = append(rows, modelRow{
			ID: "custom", Name: filepath.Base(activeModelPath(cfg)),
			Desc: strS("S_M_CUSTOM"), State: "active", Engine: cfg.STTEngine, Langs: "*",
		})
	}
	out, _ := json.Marshal(rows)
	return string(out)
}

func (a *App) downloadModel(id string) {
	m := findModel(id)
	if m == nil {
		return
	}
	if m.Engine == engineSherpa {
		a.startMultiDownload(id, m)
		return
	}
	a.startDownload(id, m.File, modelBaseURL+m.File)
}

func (a *App) startMultiDownload(key string, m *modelInfo) {
	dlMu.Lock()
	if st := dl[key]; st != nil && st.active {
		dlMu.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	dl[key] = &dlState{active: true, cancel: cancel}
	dlMu.Unlock()

	go func() {
		defer cancel()
		err := a.doMultiDownload(ctx, key, m)
		dlMu.Lock()
		switch {
		case errors.Is(err, context.Canceled):
			log.Printf("скачивание %s отменено", m.ID)
			delete(dl, key)
		case err != nil:
			log.Printf("скачивание %s: %v", m.ID, err)
			dl[key] = &dlState{err: err.Error()}
		default:
			log.Printf("модель %s скачана целиком", m.ID)
			dl[key] = &dlState{pct: 100}
		}
		dlMu.Unlock()
	}()
}

func (a *App) doMultiDownload(ctx context.Context, key string, m *modelInfo) error {
	dir := filepath.Join("models", m.Dir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	if free := freeDiskMB("models"); free >= 0 && free < m.SizeMB+256 {
		return fmt.Errorf("%s", trf("err.disk.space", free, m.SizeMB))
	}
	total := int64(m.SizeMB) * 1024 * 1024
	var doneBytes int64
	for _, f := range m.Files {
		written, err := downloadFile(ctx, m.BaseURL+f, filepath.Join(dir, f), func(n int64) {
			dlMu.Lock()
			if st := dl[key]; st != nil && total > 0 {
				pct := int((doneBytes + n) * 100 / total)
				if pct > 99 {
					pct = 99
				}
				st.pct = pct
			}
			dlMu.Unlock()
		})
		if err != nil {
			return fmt.Errorf("%s: %w", f, err)
		}
		doneBytes += written
		log.Printf("модель %s: файл %s готов (%d МБ)", m.ID, f, written/(1024*1024))
	}
	return nil
}

func downloadFile(ctx context.Context, url, final string, progress func(written int64)) (int64, error) {
	tmp := final + ".part"
	have := int64(0)
	if fi, err := os.Stat(tmp); err == nil {
		have = fi.Size()
	}
	client := &http.Client{Timeout: 0}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	if have > 0 {
		req.Header.Set("Range", "bytes="+strconv.FormatInt(have, 10)+"-")
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	resumed := resp.StatusCode == http.StatusPartialContent
	if resp.StatusCode != http.StatusOK && !resumed {
		return 0, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	if !resumed {
		have = 0
	} else {
		log.Printf("%s: продолжаю с %d МБ", filepath.Base(final), have/(1024*1024))
	}
	var f *os.File
	if resumed {
		f, err = os.OpenFile(tmp, os.O_WRONLY|os.O_APPEND, 0o644)
	} else {
		f, err = os.Create(tmp)
	}
	if err != nil {
		return 0, err
	}
	written := have
	expected := resp.ContentLength
	if expected > 0 {
		expected += have
	}
	buf := make([]byte, 256*1024)
	lastUpd := time.Now()
	for {
		n, rerr := resp.Body.Read(buf)
		if n > 0 {
			if _, werr := f.Write(buf[:n]); werr != nil {
				f.Close()
				os.Remove(tmp)
				return 0, werr
			}
			written += int64(n)
			if progress != nil && time.Since(lastUpd) > 300*time.Millisecond {
				lastUpd = time.Now()
				progress(written)
			}
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			f.Close()
			if ctx.Err() != nil {
				return written, ctx.Err()
			}
			os.Remove(tmp)
			return 0, rerr
		}
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return 0, err
	}
	if expected > 0 && written != expected {
		os.Remove(tmp)
		return 0, fmt.Errorf("файл получен не целиком: %d из %d байт", written, expected)
	}
	os.Remove(final)
	if err := os.Rename(tmp, final); err != nil {
		return 0, err
	}
	if progress != nil {
		progress(written)
	}
	return written, nil
}

func (a *App) startDownload(key, file, url string) {
	dlMu.Lock()
	if st := dl[key]; st != nil && st.active {
		dlMu.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	dl[key] = &dlState{active: true, cancel: cancel}
	dlMu.Unlock()

	go func() {
		defer cancel()
		err := a.doDownload(ctx, key, file, url)
		dlMu.Lock()
		switch {
		case errors.Is(err, context.Canceled):
			log.Printf("скачивание %s отменено", file)
			delete(dl, key)
		case err != nil:
			log.Printf("скачивание %s: %v", file, err)
			dl[key] = &dlState{err: err.Error()}
		default:
			log.Printf("модель %s скачана", file)
			dl[key] = &dlState{pct: 100}
		}
		dlMu.Unlock()
	}()
}

func freeDiskMB(dir string) int {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return -1
	}
	var free uint64
	p, _ := windows.UTF16PtrFromString(abs)
	if err := windows.GetDiskFreeSpaceEx(p, &free, nil, nil); err != nil {
		return -1
	}
	return int(free / (1024 * 1024))
}

func cleanupStaleParts() {
	entries, err := os.ReadDir("models")
	if err != nil {
		return
	}
	dlMu.Lock()
	activeAny := false
	for _, st := range dl {
		if st.active {
			activeAny = true
		}
	}
	dlMu.Unlock()
	if activeAny {
		return
	}
	dirs := []string{"models"}
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, filepath.Join("models", e.Name()))
		}
	}
	for _, dir := range dirs {
		items, derr := os.ReadDir(dir)
		if derr != nil {
			continue
		}
		for _, e := range items {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".part") {
				continue
			}
			p := filepath.Join(dir, e.Name())
			info, ierr := e.Info()
			if ierr != nil {
				_ = os.Remove(p)
				continue
			}
			if time.Since(info.ModTime()) < partKeep {
				log.Printf("незавершённая загрузка %s (%d МБ) сохранена — можно продолжить", p, info.Size()/(1024*1024))
				continue
			}
			log.Printf("удаляю незавершённую загрузку %s (%d МБ)", p, info.Size()/(1024*1024))
			_ = os.Remove(p)
		}
	}
}

func (a *App) doDownload(ctx context.Context, key, file, url string) error {
	if err := os.MkdirAll("models", 0o755); err != nil {
		return err
	}
	final := filepath.Join("models", file)
	if m := findModel(key); m != nil && m.SizeMB > 0 {
		if free := freeDiskMB("models"); free >= 0 && free < m.SizeMB+512 {
			return fmt.Errorf("%s", trf("err.disk.space", free, m.SizeMB))
		}
	}
	total := int64(0)
	if m := findModel(key); m != nil {
		total = int64(m.SizeMB) * 1024 * 1024
	}
	_, err := downloadFile(ctx, url, final, func(written int64) {
		if total <= 0 {
			return
		}
		pct := int(written * 100 / total)
		if pct > 99 {
			pct = 99
		}
		dlMu.Lock()
		if st := dl[key]; st != nil {
			st.pct = pct
		}
		dlMu.Unlock()
	})
	return err
}

func (a *App) deleteModel(id string) string {
	m := findModel(id)
	if m == nil {
		return ""
	}
	if m.isActive(a.snapshot()) {
		return tr("model.del.active")
	}
	dlMu.Lock()
	if st := dl[id]; st != nil && st.active {
		dlMu.Unlock()
		return ""
	}
	dlMu.Unlock()
	if m.Engine == engineSherpa {
		dir := filepath.Join("models", m.Dir)
		if err := os.RemoveAll(dir); err != nil {
			return err.Error()
		}
		log.Printf("модель %s удалена", m.ID)
		return tr("model.del.ok")
	}
	if err := os.Remove(filepath.Join("models", m.File)); err != nil {
		return err.Error()
	}
	log.Printf("модель %s удалена", m.File)
	return tr("model.del.ok")
}

type advicePart struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	SizeMB    int    `json:"size"`
	Installed bool   `json:"installed"`
}


type adviceOut struct {
	Primary   string       `json:"primary"`
	Companion string       `json:"companion"`
	Text      string       `json:"text"`
	RAM       string       `json:"ram"`
	Plan      []advicePart `json:"plan"`
	NeedMB    int          `json:"need"`
}

func advicePartOf(id string) (advicePart, bool) {
	m := findModel(id)
	if m == nil {
		return advicePart{}, false
	}
	return advicePart{ID: m.ID, Name: m.NameKey, SizeMB: m.SizeMB, Installed: m.installed()}, true
}

func modelDisplayName(id string) string {
	if m := findModel(id); m != nil {
		return m.NameKey
	}
	return id
}

func adviseModel(lang, priority string, needTranslate bool) string {
	_, free := ramMB()
	res := advisor.Recommend(advisor.Input{
		Lang: lang, Priority: priority, RAMFreeMB: free, Translate: needTranslate,
	}, advisorCatalog())

	var parts []string
	if res.Primary == "" {
		parts = append(parts, strS("S_ADV_NONE"))
	} else {
		parts = append(parts, trf("adv.pick", modelDisplayName(res.Primary)))
		for _, why := range res.Why {
			switch why {
			case advisor.WhyLanguage:
				parts = append(parts, strS("S_ADV_LANG"))
			case advisor.WhyAccuracy:
				parts = append(parts, strS("S_ADV_ACC"))
			case advisor.WhySpeed:
				parts = append(parts, strS("S_ADV_SPEED"))
			case advisor.WhyRAM:
				parts = append(parts, strS("S_ADV_RAM"))
			}
		}
		if res.Companion != "" {
			parts = append(parts, trf("adv.companion", modelDisplayName(res.Companion)))
		}
	}
	var plan []advicePart
	need := 0
	for _, id := range []string{res.Primary, res.Companion} {
		if id == "" {
			continue
		}
		if p, ok := advicePartOf(id); ok {
			plan = append(plan, p)
			if !p.Installed {
				need += p.SizeMB
			}
		}
	}
	out, _ := json.Marshal(adviceOut{
		Primary:   res.Primary,
		Companion: res.Companion,
		Text:      strings.Join(parts, " "),
		RAM:       trf("adv.ram", free),
		Plan:      plan,
		NeedMB:    need,
	})
	return string(out)
}

type stateOut struct {
	Hotkey   string `json:"hotkey"`
	Mic      string `json:"mic"`
	Engine   string `json:"engine"`
	LLM      string `json:"llm"`
	RAM      string `json:"ram"`
	Last     string `json:"last"`
	LastMeta string `json:"last_meta"`
	Ready    bool   `json:"ready"`
	Status   string `json:"status"`

	RuModel    string `json:"ru_model"`
	OtherModel string `json:"other_model"`
	LLMOK      bool   `json:"llm_ok"`
	MicOK      bool   `json:"mic_ok"`
	StatusLine string `json:"status_line"`
	Remote      bool   `json:"remote"`
	RuState     string `json:"ru_state"`
	OtherState  string `json:"other_state"`
	Badges     struct {
		Mic    string `json:"mic"`
		Models string `json:"models"`
		System string `json:"system"`
	} `json:"badges"`
}

func agoLabel(d time.Duration) string {
	switch {
	case d < time.Minute:
		return tr("ago.now")
	case d < time.Hour:
		return trf("ago.min", int(d.Minutes()))
	default:
		return trf("ago.hour", int(d.Hours()))
	}
}

func installedModelCount() int {
	n := 0
	for i := range modelCatalog {
		if modelCatalog[i].installed() {
			n++
		}
	}
	return n
}

func systemWarnings(cfg *Config) int {
	n := 0
	if cfg.ServerURL != "" {
		n++
	}
	if !cfg.ServerAutostart {
		n++
	}
	return n
}

func modelStateFor(cfg *Config, engine string) string {
	var id string
	installed := false
	switch engine {
	case engineSherpa:
		installed = sherpaInstalled(cfg)
		dir := filepath.Base(filepath.Clean(cfg.SherpaModel))
		for i := range modelCatalog {
			if modelCatalog[i].Engine == engineSherpa && modelCatalog[i].Dir == dir {
				id = modelCatalog[i].ID
			}
		}
	default:
		if strings.TrimSpace(cfg.ServerURL) != "" {
			return "remote"
		}
		if _, err := os.Stat(cfg.Model); err == nil {
			installed = true
		}
		file := filepath.Base(cfg.Model)
		for i := range modelCatalog {
			if modelCatalog[i].Engine != engineSherpa && modelCatalog[i].File == file {
				id = modelCatalog[i].ID
			}
		}
	}
	if installed {
		return "ready"
	}
	if id != "" {
		dlMu.Lock()
		st := dl[id]
		dlMu.Unlock()
		if st != nil && st.active {
			return "downloading"
		}
	}
	return "missing"
}

func (a *App) stateSnapshot() string {
	cfg := a.snapshot()
	a.mu.Lock()
	ready := a.ready
	rec := a.rec
	last := a.lastResult
	a.mu.Unlock()

	mic := strS("S_MIC_DEFAULT")
	if cfg.MicDeviceName != "" {
		mic = cfg.MicDeviceName
	} else if rec != nil {
		for _, d := range rec.devices() {
			if d.System || d.Default {
				mic = d.Name
				if d.System {
					break
				}
			}
		}
	}
	llm := strS("S_NO_LLM")
	if llmInstalled(cfg) {
		llm = filepath.Base(cfg.LLMModel)
	}
	_, free := ramMB()
	status := tr("status.loading")
	if ready {
		status = trf("status.ready", cfg.Hotkey)
	}
	if last == "" {
		last = "—"
	}
	lastMeta := ""
	if a.lastResultAt.IsZero() {
		last = "—"
	} else {
		parts := []string{agoLabel(time.Since(a.lastResultAt)), trf("chars", len([]rune(a.lastResult)))}
		if a.lastTarget != "" {
			parts = append(parts, trf("inserted.into", a.lastTarget))
		}
		lastMeta = strings.Join(parts, " · ")
	}

	ruModel := filepath.Base(filepath.Clean(cfg.SherpaModel))
	if !sherpaInstalled(cfg) {
		ruModel = strS("S_NOT_INSTALLED")
	}
	otherModel := filepath.Base(cfg.Model)
	if modelStateFor(cfg, engineWhisper) == "missing" {
		otherModel = strS("S_NOT_INSTALLED")
	}
	st := stateOut{
		Hotkey:     cfg.Hotkey,
		Mic:        mic,
		Engine:     primaryEngine(cfg) + " · " + filepath.Base(filepath.Clean(activeModelPath(cfg))),
		LLM:        llm,
		RAM:        trf("adv.ram", free),
		Last:       last,
		LastMeta:   lastMeta,
		Ready:      ready,
		Status:     status,
		RuModel:    ruModel,
		OtherModel: otherModel,
		LLMOK:      llmInstalled(cfg),
		MicOK:      rec != nil,
		StatusLine: statusLine(cfg, ready, free),
		Remote:      strings.TrimSpace(cfg.ServerURL) != "",
		RuState:     modelStateFor(cfg, engineSherpa),
		OtherState:  modelStateFor(cfg, engineWhisper),
	}
	st.Badges.Mic = micBadge(mic)
	st.Badges.Models = itoaSafe(installedModelCount())
	if w := systemWarnings(cfg); w > 0 {
		st.Badges.System = itoaSafe(w)
	}
	out, _ := json.Marshal(st)
	return string(out)
}

func itoaSafe(n int) string {
	return fmt.Sprintf("%d", n)
}

func shortLabel(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max-1]) + "…"
}

var micNoise = map[string]bool{
	"microphone": true, "mic": true, "audio": true, "input": true, "device": true,
	"микрофон": true, "звук": true, "устройство": true, "system": true, "системный": true,
}

func micBadge(name string) string {
	clean := strings.Map(func(r rune) rune {
		if r == '(' || r == ')' || r == '-' || r == ',' {
			return ' '
		}
		return r
	}, name)
	for _, w := range strings.Fields(clean) {
		if len([]rune(w)) < 3 {
			continue
		}
		if micNoise[strings.ToLower(w)] {
			continue
		}
		return shortLabel(w, 10)
	}
	return shortLabel(strings.TrimSpace(name), 10)
}

func statusLine(cfg *Config, ready bool, freeMB int) string {
	if !ready {
		return tr("status.loading")
	}
	models := filepath.Base(filepath.Clean(activeModelPath(cfg)))
	if primaryEngine(cfg) == engineSherpa {
		models += " + " + filepath.Base(cfg.Model)
	}
	return trf("status.line", models, float64(freeMB)/1024)
}

func cancelDownload(id string) bool {
	dlMu.Lock()
	st := dl[id]
	dlMu.Unlock()
	if st == nil || !st.active || st.cancel == nil {
		return false
	}
	st.cancel()
	return true
}
