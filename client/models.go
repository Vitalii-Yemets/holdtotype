package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
}

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
	dl[key] = &dlState{active: true}
	dlMu.Unlock()

	go func() {
		err := a.doMultiDownload(key, m)
		dlMu.Lock()
		if err != nil {
			log.Printf("скачивание %s: %v", m.ID, err)
			dl[key] = &dlState{err: err.Error()}
		} else {
			log.Printf("модель %s скачана целиком", m.ID)
			dl[key] = &dlState{pct: 100}
		}
		dlMu.Unlock()
	}()
}

func (a *App) doMultiDownload(key string, m *modelInfo) error {
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
		written, err := downloadFile(m.BaseURL+f, filepath.Join(dir, f), func(n int64) {
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

func downloadFile(url, final string, progress func(written int64)) (int64, error) {
	tmp := final + ".part"
	client := &http.Client{Timeout: 0}
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	f, err := os.Create(tmp)
	if err != nil {
		return 0, err
	}
	var written int64
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
			os.Remove(tmp)
			return 0, rerr
		}
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return 0, err
	}
	if resp.ContentLength > 0 && written != resp.ContentLength {
		os.Remove(tmp)
		return 0, fmt.Errorf("файл получен не целиком: %d из %d байт", written, resp.ContentLength)
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
	dl[key] = &dlState{active: true}
	dlMu.Unlock()

	go func() {
		err := a.doDownload(key, file, url)
		dlMu.Lock()
		if err != nil {
			log.Printf("скачивание %s: %v", file, err)
			dl[key] = &dlState{err: err.Error()}
		} else {
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
			if info, ierr := e.Info(); ierr == nil {
				log.Printf("удаляю незавершённую загрузку %s (%d МБ)", p, info.Size()/(1024*1024))
			}
			_ = os.Remove(p)
		}
	}
}

func (a *App) doDownload(key, file, url string) error {
	if err := os.MkdirAll("models", 0o755); err != nil {
		return err
	}
	tmp := filepath.Join("models", file+".part")
	final := filepath.Join("models", file)

	client := &http.Client{Timeout: 0}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	total := resp.ContentLength
	if total > 0 {
		if free := freeDiskMB("models"); free >= 0 && int64(free) < total/(1024*1024)+512 {
			return fmt.Errorf("%s", trf("err.disk.space", free, total/(1024*1024)))
		}
	}

	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	var written int64
	buf := make([]byte, 256*1024)
	lastUpd := time.Now()
	for {
		n, rerr := resp.Body.Read(buf)
		if n > 0 {
			if _, werr := f.Write(buf[:n]); werr != nil {
				f.Close()
				os.Remove(tmp)
				return werr
			}
			written += int64(n)
			if total > 0 && time.Since(lastUpd) > 300*time.Millisecond {
				lastUpd = time.Now()
				dlMu.Lock()
				if st := dl[key]; st != nil {
					st.pct = int(written * 100 / total)
				}
				dlMu.Unlock()
			}
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			f.Close()
			os.Remove(tmp)
			return rerr
		}
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	os.Remove(final)
	return os.Rename(tmp, final)
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

type adviceOut struct {
	Primary   string `json:"primary"`
	Companion string `json:"companion"`
	Text      string `json:"text"`
	RAM       string `json:"ram"`
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
	out, _ := json.Marshal(adviceOut{
		Primary:   res.Primary,
		Companion: res.Companion,
		Text:      strings.Join(parts, " "),
		RAM:       trf("adv.ram", free),
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
		OtherModel: filepath.Base(cfg.Model),
		LLMOK:      llmInstalled(cfg),
		MicOK:      rec != nil,
		StatusLine: statusLine(cfg, ready, free),
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
