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

	"holdtotype/internal/gpucard"
	"holdtotype/internal/subtitles"
	"holdtotype/internal/trmodel"
)

const (
	gpuExeName    = "llama-server-gpu.exe"
	trModelKey    = "trmodel"
	qualityNeedMB = 3300
	qualityThread = 4
)

type qualityFail struct {
	code   string
	detail string
}

func (f *qualityFail) Error() string { return f.detail }

var (
	cardOnce sync.Once
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

var trm struct {
	mu       sync.Mutex
	checking bool
	checked  bool
	size     int64
	mod      time.Time
	state    trmodel.State
	err      string
}

func forgetTrModelCheck() {
	trm.mu.Lock()
	trm.checked = false
	trm.mu.Unlock()
}

func trModelState() trmodel.State {
	if st, ok := trmodel.Known("models", trmodel.Fixed()); ok {
		return st
	}
	fi, err := os.Stat(filepath.Join("models", trmodel.File))
	if err != nil {
		return trmodel.Absent
	}
	trm.mu.Lock()
	defer trm.mu.Unlock()
	if trm.checked && trm.size == fi.Size() && trm.mod.Equal(fi.ModTime()) {
		if trm.err != "" {
			return "unreadable"
		}
		return trm.state
	}
	if !trm.checking {
		trm.checking = true
		size, mod := fi.Size(), fi.ModTime()
		go func() {
			started := time.Now()
			st, err := trmodel.Check("models", trmodel.Fixed())
			trm.mu.Lock()
			trm.checking = false
			trm.checked = true
			trm.size, trm.mod = size, mod
			trm.state = st
			trm.err = ""
			if err != nil {
				trm.err = humanError(err)
			}
			trm.mu.Unlock()
			if err != nil {
				log.Printf("translation model: the file cannot be checked: %v", err)
			} else {
				log.Printf("translation model: the file found in the models folder is checked in %.1f s: %s", time.Since(started).Seconds(), st)
			}
		}()
	}
	return "checking"
}

func trModelReady() bool { return trModelState() == trmodel.Ready }

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
	if st := trModelState(); st != trmodel.Ready {
		return engineSpec{}, &qualityFail{"model", fmt.Sprintf("the translation model is not ready (%s)", st)}
	}
	return gpuEngineFor(cfg, filepath.Join("models", trmodel.File), qualityNeedMB)
}

func gpuEngineFor(cfg *Config, model string, needMB int) (engineSpec, error) {
	card, fits := qualityCard()
	if !fits {
		return engineSpec{}, &qualityFail{"card", "no graphics card fits: " + cardWhy}
	}
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
	if dev.FreeMB < needMB {
		return engineSpec{}, &qualityFail{"memory", fmt.Sprintf("%d MB of the memory of %s are free, %d MB are needed", dev.FreeMB, card.Name, needMB)}
	}
	log.Printf("graphics card engine: device %s (%s), %d MB of %d MB free", dev.ID, dev.Name, dev.FreeMB, dev.TotalMB)
	return engineSpec{
		label: "llama-server-gpu",
		exe:   exe,
		port:  cfg.TranslatePort,
		args: []string{
			"-m", model,
			"--device", dev.ID,
			"-ngl", "999",
			"-t", strconv.Itoa(min(cfg.Threads, qualityThread)),
			"-c", "4096",
		},
		wait: 2 * time.Minute,
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

type trModelUI struct {
	State   string `json:"state"`
	Pct     int    `json:"pct"`
	Err     string `json:"err"`
	SizeMB  int    `json:"size"`
	Fits    bool   `json:"fits"`
	Card    string `json:"card"`
	Reason  string `json:"reason"`
	CardMB  int    `json:"card_mb"`
	NeedMB  int    `json:"need_mb"`
	Down    string `json:"down"`
	Running bool   `json:"running"`
}

func (a *App) trModelView() trModelUI {
	card, fits := qualityCard()
	v := trModelUI{SizeMB: trmodel.SizeMB(), Fits: fits, Card: card.Name, Reason: cardCode, CardMB: cardMB, NeedMB: gpucard.NeedMB}
	if f := qualityDown(); f != nil {
		v.Down = f.code
	}
	v.Running = a.gpuEng.live() != nil
	dlMu.Lock()
	st := dl[trModelKey]
	if st != nil {
		switch {
		case st.active && st.pct >= 100:
			v.State = "checking"
		case st.active:
			v.State, v.Pct = "downloading", st.pct
		case st.err != "":
			v.State, v.Err = "error", st.err
		}
	}
	dlMu.Unlock()
	if v.State != "" {
		return v
	}
	switch s := trModelState(); s {
	case trmodel.Ready:
		v.State = "ready"
	case trmodel.WrongSize, trmodel.WrongHash:
		v.State = "bad"
	case trmodel.Absent:
		v.State = "absent"
	case "unreadable":
		trm.mu.Lock()
		v.State, v.Err = "error", trm.err
		trm.mu.Unlock()
	default:
		v.State = "checking"
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

func (a *App) trModelDownloadSpec(spec trmodel.Spec, url string) {
	a.startDownloadWith(trModelKey, spec.Name, url, int(spec.Size/(1024*1024)), verifyTrModel(spec))
}

func (a *App) trModelDownload() {
	if st, ok := trmodel.Known("models", trmodel.Fixed()); ok && st == trmodel.Ready {
		return
	}
	log.Printf("translation model: download started")
	a.trModelDownloadSpec(trmodel.Fixed(), trmodel.URL)
}

func (a *App) trModelDelete() string {
	a.qualityShutdown()
	err := trmodel.Remove("models", trmodel.Fixed())
	forgetTrModelCheck()
	if err != nil {
		return humanError(err)
	}
	dlMu.Lock()
	delete(dl, trModelKey)
	dlMu.Unlock()
	log.Printf("translation model deleted")
	return tr("model.del.ok")
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
	app.trModelDownloadSpec(spec, args[0])
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
	need := int(fi.Size()/(1024*1024)) + 1000
	app := &App{cfg: cfg, enabled: true}
	started := time.Now()
	f := app.warmEngine(&app.gpuEng, func(c *Config) (engineSpec, error) { return gpuEngineFor(c, args[0], need) })
	<-f.done
	srv, err := asLlama(f.srv, f.err)
	if err != nil {
		fail := asQualityFail(err)
		log.Printf("gpucheck: did not start (%s): %s", fail.code, fail.detail)
		return
	}
	log.Printf("gpucheck: ready in %.1f s", time.Since(started).Seconds())
	primed := time.Now()
	if err := srv.prime("ru", ""); err != nil {
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
