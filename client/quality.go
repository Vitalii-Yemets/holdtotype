package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"holdtotype/internal/editormodel"
	"holdtotype/internal/gpucard"
	"holdtotype/internal/subtitles"
	"holdtotype/internal/trmodel"
)

const (
	gpuExeName    = "llama-server-gpu.exe"
	trModelKey    = "trmodel"
	qualityThread = 4
)

type qualityFail struct {
	code   string
	detail string
}

func (f *qualityFail) Error() string { return f.detail }

var (
	cardOnce sync.Once
	cardList []gpucard.Adapter
	cardErr  string
	cardBest gpucard.Adapter
	cardFits bool
	cardWhy  string
	cardCode string
	cardMB   int
)

func readCard() {
	cardOnce.Do(func() {
		list, err := gpucard.Adapters()
		if err != nil {
			cardErr = err.Error()
			cardWhy, cardCode = err.Error(), gpucard.ReasonUnknown
			log.Printf("graphics card: %s", cardWhy)
			return
		}
		var names []string
		for _, a := range list {
			if a.Hardware {
				kind := "own memory"
				if a.Integrated {
					kind = "integrated"
				}
				names = append(names, fmt.Sprintf("%s (%d MB, %s)", a.Name, a.DedicatedMB, kind))
			}
		}
		cardList = list
		best, why := gpucard.Pick(list)
		cardBest, cardFits, cardWhy = best, why == "", why
		cardCode, cardMB = gpucard.Reason(list)
		if cardFits {
			log.Printf("graphics card for translation: %s with %d MB of its own memory (adapters: %s)", best.Name, best.DedicatedMB, strings.Join(names, ", "))
		} else {
			log.Printf("graphics card for translation: none fits, %s (adapters: %s)", why, strings.Join(names, ", "))
		}
	})
}

func qualityCard() (gpucard.Adapter, bool) {
	readCard()
	return cardBest, cardFits
}

type cardFit struct {
	card   gpucard.Adapter
	fits   bool
	why    string
	code   string
	cardMB int
}

func cardFor(m trmodel.Model) cardFit {
	readCard()
	if cardErr != "" {
		return cardFit{why: cardErr, code: gpucard.ReasonUnknown}
	}
	best, why := gpucard.PickFor(cardList, m.CardMB)
	code, mb := gpucard.ReasonFor(cardList, m.CardMB)
	return cardFit{card: best, fits: why == "", why: why, code: code, cardMB: mb}
}

type trCheck struct {
	checking bool
	checked  bool
	size     int64
	mod      time.Time
	state    trmodel.State
	err      string
}

var trm struct {
	mu    sync.Mutex
	files map[string]*trCheck
}

func forgetTrModelCheck() {
	trm.mu.Lock()
	trm.files = nil
	trm.mu.Unlock()
}

func trModelStateOf(m trmodel.Model) trmodel.State {
	path := filepath.Join("models", m.File)
	if m.Own {
		if fi, err := os.Stat(path); err == nil && fi.Mode().IsRegular() {
			return trmodel.Ready
		}
		return trmodel.Absent
	}
	if st, ok := trmodel.Known("models", m.Spec()); ok {
		return st
	}
	fi, err := os.Stat(path)
	if err != nil {
		return trmodel.Absent
	}
	trm.mu.Lock()
	defer trm.mu.Unlock()
	if trm.files == nil {
		trm.files = map[string]*trCheck{}
	}
	c := trm.files[m.File]
	if c == nil {
		c = &trCheck{}
		trm.files[m.File] = c
	}
	if c.checked && c.size == fi.Size() && c.mod.Equal(fi.ModTime()) {
		if c.err != "" {
			return "unreadable"
		}
		return c.state
	}
	if !c.checking {
		c.checking = true
		size, mod := fi.Size(), fi.ModTime()
		go func() {
			started := time.Now()
			st, err := trmodel.Check("models", m.Spec())
			trm.mu.Lock()
			c.checking = false
			c.checked = true
			c.size, c.mod = size, mod
			c.state = st
			c.err = ""
			if err != nil {
				c.err = humanError(err)
			}
			trm.mu.Unlock()
			if err != nil {
				log.Printf("translation model %s: the file cannot be checked: %v", m.Name, err)
			} else {
				log.Printf("translation model %s: the file found in the models folder is checked in %.1f s: %s", m.Name, time.Since(started).Seconds(), st)
			}
		}()
	}
	return "checking"
}

func trCheckErr(m trmodel.Model) string {
	trm.mu.Lock()
	defer trm.mu.Unlock()
	if c := trm.files[m.File]; c != nil {
		return c.err
	}
	return ""
}

func trModelReadyOf(m trmodel.Model) bool { return trModelStateOf(m) == trmodel.Ready }

func checkTrModels() {
	for _, m := range trmodel.Catalog {
		trModelStateOf(m)
	}
}

func ownTrFiles() []trmodel.OwnFile {
	var out []trmodel.OwnFile
	for _, f := range editormodel.List("models") {
		out = append(out, trmodel.OwnFile{Name: f.Name, SizeMB: f.SizeMB})
	}
	return out
}

func currentTrModel(cfg *Config) (trmodel.Model, bool) {
	return trmodel.Resolve(cfg.SubtitlesModel, trModelReadyOf, ownTrFiles())
}

func shownTrModel(cfg *Config) trmodel.Model {
	return trmodel.Shown(cfg.SubtitlesModel, trModelReadyOf, ownTrFiles())
}

func trModelReady(cfg *Config) bool {
	_, ok := currentTrModel(cfg)
	return ok
}

var qualityFailed atomic.Pointer[qualityFail]

func qualityDown() *qualityFail { return qualityFailed.Load() }

func setQualityDown(f *qualityFail) { qualityFailed.Store(f) }

func gpuExePath(cfg *Config) string {
	return filepath.Join(filepath.Dir(absPath(cfg.LLMExe)), gpuExeName)
}

func listGPUDevices(exe string) ([]gpucard.Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, exe, "--list-devices")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	out, err := cmd.CombinedOutput()
	if ctx.Err() != nil {
		return nil, fmt.Errorf("listing the devices took longer than 30 s")
	}
	if err != nil {
		return nil, fmt.Errorf("%v: %.200s", err, strings.TrimSpace(string(out)))
	}
	return gpucard.ParseDevices(string(out)), nil
}

func gpuEngine(cfg *Config) (engineSpec, error) {
	m, ok := currentTrModel(cfg)
	if !ok {
		return engineSpec{}, &qualityFail{"model", "no translation model is ready"}
	}
	return gpuEngineFor(cfg, m, filepath.Join("models", m.File))
}

func gpuEngineFor(cfg *Config, m trmodel.Model, path string) (engineSpec, error) {
	fit := cardFor(m)
	if !fit.fits {
		return engineSpec{}, &qualityFail{"card", "no graphics card fits " + m.Name + ": " + fit.why}
	}
	card := fit.card
	exe := gpuExePath(cfg)
	if _, err := os.Stat(exe); err != nil {
		return engineSpec{}, &qualityFail{"exe", gpuExeName + " is missing next to the program"}
	}
	devs, err := listGPUDevices(exe)
	if err != nil {
		return engineSpec{}, &qualityFail{"vulkan", "the graphics card engine cannot list the devices: " + err.Error()}
	}
	if len(devs) == 0 {
		return engineSpec{}, &qualityFail{"vulkan", "the graphics card engine sees no Vulkan device: the graphics driver has no Vulkan runtime"}
	}
	dev, ok := gpucard.Match(devs, card.Name)
	if !ok {
		var seen []string
		for _, d := range devs {
			seen = append(seen, d.ID+" "+d.Name)
		}
		return engineSpec{}, &qualityFail{"device", fmt.Sprintf("the graphics card %s is not among the Vulkan devices (%s)", card.Name, strings.Join(seen, ", "))}
	}
	if dev.FreeMB < m.FreeMB {
		return engineSpec{}, &qualityFail{"memory", fmt.Sprintf("%d MB of the memory of %s are free, %d MB are needed for %s", dev.FreeMB, card.Name, m.FreeMB, m.Name)}
	}
	log.Printf("graphics card engine: device %s (%s), %d MB of %d MB free, model %s", dev.ID, dev.Name, dev.FreeMB, dev.TotalMB, m.Name)
	args := []string{
		"-m", path,
		"--device", dev.ID,
		"-ngl", "999",
		"-t", strconv.Itoa(min(cfg.Threads, qualityThread)),
		"-c", "4096",
	}
	args = append(args, m.ChatArgs...)
	return engineSpec{
		label: "llama-server-gpu",
		tag:   m.Key(),
		exe:   exe,
		port:  cfg.TranslatePort,
		args:  args,
		wait:  2 * time.Minute,
	}, nil
}

func asQualityFail(err error) *qualityFail {
	var f *qualityFail
	if errors.As(err, &f) {
		return f
	}
	code := "crash"
	switch {
	case errors.Is(err, errEngineTimeout):
		code = "timeout"
	case strings.Contains(logText(err), "is taken by another program"):
		code = "port"
	}
	return &qualityFail{code, logText(err)}
}

func (a *App) qualityWarmUp() *llmFlight { return a.warmEngine(&a.gpuEng, gpuEngine) }

func (a *App) qualityWait(ctx context.Context) (*llamaServer, error) {
	return asLlama(a.waitEngine(ctx, &a.gpuEng, gpuEngine, "the translation model is still starting"))
}

func (a *App) qualityShutdown() {
	if a.gpuEng.shutdown() {
		log.Printf("llama-server-gpu stopped, the graphics card memory is free again")
	}
}

type trModelRow struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	File    string `json:"file"`
	State   string `json:"state"`
	Pct     int    `json:"pct"`
	Err     string `json:"err"`
	SizeMB  int    `json:"size"`
	VramMB  int    `json:"vram"`
	CardMB  int    `json:"card_need"`
	Fits    bool   `json:"fits"`
	Reason  string `json:"reason"`
	Topic   bool   `json:"topic"`
	Terms   bool   `json:"terms"`
	Rec     bool   `json:"rec"`
	Acc     int    `json:"accuracy"`
	Speed   int    `json:"speed"`
	Own     bool   `json:"own"`
	Chosen  bool   `json:"chosen"`
	Running bool   `json:"running"`
	License string `json:"license"`
}

type trModelUI struct {
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	State   string       `json:"state"`
	Pct     int          `json:"pct"`
	Err     string       `json:"err"`
	SizeMB  int          `json:"size"`
	Fits    bool         `json:"fits"`
	Card    string       `json:"card"`
	Reason  string       `json:"reason"`
	CardMB  int          `json:"card_mb"`
	NeedMB  int          `json:"need_mb"`
	Down    string       `json:"down"`
	Running bool         `json:"running"`
	Models  []trModelRow `json:"models"`
}

func trDlKey(m trmodel.Model) string { return trModelKey + ":" + m.ID }

func trModelRowOf(a *App, m trmodel.Model, chosen string) trModelRow {
	fit := cardFor(m)
	row := trModelRow{ID: m.Key(), Name: m.Name, File: m.File, SizeMB: m.SizeMB(), VramMB: m.VramMB, CardMB: m.CardMB, Fits: fit.fits, Reason: fit.code,
		Topic: m.Topic, Terms: m.Terms, Rec: m.Recommended, Acc: m.Accuracy, Speed: m.Speed, Own: m.Own, Chosen: m.Key() == chosen, License: m.License}
	row.Running = row.Chosen && a.gpuEng.liveTag() == m.Key()
	dlMu.Lock()
	st := dl[trDlKey(m)]
	if st != nil {
		switch {
		case st.active && st.pct >= 100:
			row.State = "checking"
		case st.active:
			row.State, row.Pct = "downloading", st.pct
		case st.err != "":
			row.State, row.Err = "error", st.err
		}
	}
	dlMu.Unlock()
	if row.State != "" {
		return row
	}
	switch s := trModelStateOf(m); s {
	case trmodel.Ready:
		row.State = "ready"
	case trmodel.WrongSize, trmodel.WrongHash:
		row.State = "bad"
	case trmodel.Absent:
		row.State = "absent"
	case "unreadable":
		row.State, row.Err = "error", trCheckErr(m)
	default:
		row.State = "checking"
	}
	return row
}

func (a *App) trModelView(cfg *Config) trModelUI {
	shown := shownTrModel(cfg)
	chosen := ""
	if m, ok := currentTrModel(cfg); ok {
		chosen = m.Key()
	}
	var rows []trModelRow
	for _, m := range trmodel.Catalog {
		rows = append(rows, trModelRowOf(a, m, chosen))
	}
	for _, f := range trmodel.SortOwn(ownTrFiles()) {
		rows = append(rows, trModelRowOf(a, trmodel.Own(f.Name, f.SizeMB), chosen))
	}
	fit := cardFor(shown)
	v := trModelUI{ID: shown.Key(), Name: shown.Name, SizeMB: shown.SizeMB(), Fits: fit.fits, Card: fit.card.Name, Reason: fit.code, CardMB: fit.cardMB, NeedMB: shown.CardMB, Models: rows}
	if !fit.fits {
		v.Card = cardBest.Name
	}
	if f := qualityDown(); f != nil {
		v.Down = f.code
	}
	for _, r := range rows {
		if r.ID == shown.Key() {
			v.State, v.Pct, v.Err, v.Running = r.State, r.Pct, r.Err, r.Running
		}
	}
	return v
}

func verifyTrModel(spec trmodel.Spec) func(string) error {
	return func(string) error {
		st, err := trmodel.Check("models", spec)
		forgetTrModelCheck()
		if err != nil {
			return err
		}
		if st != trmodel.Ready {
			_ = trmodel.Remove("models", spec)
			log.Printf("%s: the downloaded file does not match its checksum (%s)", spec.Name, st)
			return uiErr("the downloaded file does not match its checksum", tr("err.hash"))
		}
		log.Printf("%s: hash matches", spec.Name)
		return nil
	}
}

func (a *App) trModelDownloadSpec(key string, spec trmodel.Spec, url string) {
	a.startDownloadWith(key, spec.Name, url, int(spec.Size/(1024*1024)), verifyTrModel(spec))
}

func trModelByKey(cfg *Config, id string) (trmodel.Model, bool) {
	if id == "" {
		return shownTrModel(cfg), true
	}
	return trmodel.ByID(id)
}

func (a *App) trModelDownload(id string) {
	m, ok := trModelByKey(a.snapshot(), id)
	if !ok || m.Own {
		return
	}
	if st, ok := trmodel.Known("models", m.Spec()); ok && st == trmodel.Ready {
		return
	}
	log.Printf("translation model %s: download started", m.Name)
	a.trModelDownloadSpec(trDlKey(m), m.Spec(), m.URL)
}

func (a *App) trModelCancel(id string) bool {
	m, ok := trModelByKey(a.snapshot(), id)
	if !ok {
		return false
	}
	return cancelDownload(trDlKey(m))
}

func (a *App) trModelDelete(id string) string {
	m, ok := trModelByKey(a.snapshot(), id)
	if !ok || m.Own {
		return ""
	}
	if a.gpuEng.liveTag() == m.Key() {
		a.qualityShutdown()
	}
	err := trmodel.Remove("models", m.Spec())
	forgetTrModelCheck()
	if err != nil {
		return humanError(err)
	}
	dlMu.Lock()
	delete(dl, trDlKey(m))
	dlMu.Unlock()
	log.Printf("translation model %s deleted", m.Name)
	return tr("model.del.ok")
}

func trModelTerms(cfg *Config, id string) string {
	m, ok := trModelByKey(cfg, id)
	if !ok || m.TermsURL == "" {
		return trmodel.TermsURL
	}
	return m.TermsURL
}

func runTrDownload(args []string) {
	if len(args) < 4 {
		log.Printf("trdownload: needs <url> <sha256> <bytes> <file name>")
		return
	}
	size, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil || size <= 0 || strings.ContainsAny(args[3], `/\`) {
		log.Printf("trdownload: bad arguments %q", args)
		return
	}
	spec := trmodel.Spec{Name: args[3], Size: size, Hash: args[1]}
	app := &App{cfg: defaultConfig(), enabled: true}
	app.trModelDownloadSpec(trModelKey, spec, args[0])
	for {
		time.Sleep(200 * time.Millisecond)
		dlMu.Lock()
		st := dl[trModelKey]
		active := st != nil && st.active
		dlMu.Unlock()
		if !active {
			break
		}
	}
	dlMu.Lock()
	st := dl[trModelKey]
	dlMu.Unlock()
	state, _ := trmodel.Check("models", spec)
	switch {
	case st != nil && st.err != "":
		log.Printf("trdownload: error state: %s, file %s", st.err, state)
	default:
		log.Printf("trdownload: done, file %s", state)
	}
}

func runGPUCheck(args []string) {
	if len(args) < 1 || !strings.HasSuffix(strings.ToLower(args[0]), ".gguf") {
		log.Printf("gpucheck: needs <model.gguf>")
		return
	}
	cfg, err := loadConfig("config.json")
	if err != nil {
		log.Printf("config: %v", err)
		return
	}
	fi, err := os.Stat(args[0])
	if err != nil {
		log.Printf("gpucheck: %v", err)
		return
	}
	m, known := trmodel.ByFile(args[0])
	if !known {
		m = trmodel.Own(filepath.Base(args[0]), int(fi.Size()/(1024*1024)))
	}
	app := &App{cfg: cfg, enabled: true}
	started := time.Now()
	f := app.warmEngine(&app.gpuEng, func(c *Config) (engineSpec, error) { return gpuEngineFor(c, m, args[0]) })
	<-f.done
	srv, err := asLlama(f.srv, f.err)
	if err != nil {
		fail := asQualityFail(err)
		log.Printf("gpucheck: did not start (%s): %s", fail.code, fail.detail)
		return
	}
	log.Printf("gpucheck: ready in %.1f s", time.Since(started).Seconds())
	primed := time.Now()
	if err := srv.prime("ru"); err != nil {
		log.Printf("gpucheck: the priming request failed: %v", err)
	} else {
		log.Printf("gpucheck: warm after %.1f s of priming, %.1f s after the start", time.Since(primed).Seconds(), time.Since(started).Seconds())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	asked := time.Now()
	var first time.Time
	out, err := srv.chatStream(ctx, subtitles.TranslateMessages("ru", nil, "So the first thing we noticed was the price."), 64, "", func(string) {
		if first.IsZero() {
			first = time.Now()
		}
	})
	cancel()
	if err != nil {
		log.Printf("gpucheck: the request failed: %v", err)
	} else {
		log.Printf("gpucheck: first word after %.2f s, answer after %.2f s: %q", first.Sub(asked).Seconds(), time.Since(asked).Seconds(), strings.TrimSpace(out))
	}
	app.qualityShutdown()
	log.Printf("gpucheck: done")
}
