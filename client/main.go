package main

import (
	"encoding/json"
	"holdtotype/internal/apprules"
	"holdtotype/internal/audiolevel"
	"holdtotype/internal/commands"
	"holdtotype/internal/evqueue"
	"holdtotype/internal/history"
	"holdtotype/internal/livetail"
	"holdtotype/internal/profiles"
	"holdtotype/internal/replace"
	"holdtotype/internal/segments"
	"holdtotype/internal/textshape"
	"unsafe"

	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"holdtotype/internal/appid"
)

var appVersion = appid.Version

type appState = int32

const (
	stIdle appState = iota
	stRecording
	stProcessing
)

const (
	evDown = iota
	evUp
	evTimeout
	evDone
	evCancel
	evMicLost
	evSilence
)

type ptEvent struct {
	kind    int
	gen     int
	profile string
	cut     int
	ep      *audiolevel.Endpoint
}

type stopReason int

const (
	stopKey stopReason = iota
	stopLimit
	stopSilence
)

type App struct {
	mu         sync.Mutex
	cfg        *Config
	rec        *Recorder
	srv        recognizer
	llm        *llamaServer
	hook       *hotkeyHook
	enabled    bool
	ready      bool
	capturing  bool
	parked     bool
	backendErr string
	quitting   bool

	evq     *evqueue.Queue[ptEvent]
	retryCh chan struct{}

	state          atomic.Int32
	gen            int
	sessionCfg     *Config
	sessionProfile string
	sessionCancel  context.CancelFunc
	sessionTarget  uintptr
	recCancel      context.CancelFunc
	silenceEv      ptEvent
	silenceStopAt  time.Time

	lastResult   string
	lastResultAt time.Time
	lastVerdict  string
	lastTarget   string
	lastWnd      uintptr
	settingsPrev uintptr
	lastProcess  string
	lastIns      *insertion
	armedShapes  []commands.Command
	armedAt      time.Time
	postErr      string
	postErrProf  string
	postErrAt    time.Time
	updVer       string
	updChecked   time.Time
	updURL       string
	updDigest    string

	altMu   sync.Mutex
	alt     recognizer
	altUsed time.Time

	liveMu  sync.Mutex
	live    *streamSession
	liveGen int
	liveFed int

	subMu sync.Mutex
	sub   *subRun
}

func (a *App) liveStream(gen int, cfg *Config) {
	a.mu.Lock()
	srv := a.srv
	ready := a.ready
	rec := a.rec
	a.mu.Unlock()
	ss, ok := srv.(*streamServer)
	if !ok || !ready || rec == nil {
		return
	}
	sess, err := ss.openSession(context.Background(), func(text string) {
		if a.state.Load() != stRecording {
			return
		}
		a.liveMu.Lock()
		mine := a.liveGen == gen
		a.liveMu.Unlock()
		if mine && cfg.Overlay {
			overlaySet(ovRecording, livetail.Tail(text, 600))
		}
	})
	if err != nil {
		log.Printf("live text did not start (%v) — recognizing after the recording", err)
		return
	}
	a.liveMu.Lock()
	if a.live != nil {
		a.live.close()
	}
	a.live = sess
	a.liveGen = gen
	a.liveFed = 0
	a.liveMu.Unlock()
	ovLive.Store(true)
	if cfg.Overlay && a.state.Load() == stRecording {
		overlaySet(ovRecording, "")
	}
	log.Printf("live text: the stream is open")
	offset := 0
	for {
		time.Sleep(180 * time.Millisecond)
		a.liveMu.Lock()
		mine := a.live == sess && a.liveGen == gen
		a.liveMu.Unlock()
		if !mine || a.state.Load() != stRecording {
			return
		}
		chunk, next := rec.TakeFrom(offset)
		if len(chunk) > 0 {
			if err := sess.push(chunk); err != nil {
				log.Printf("live text: audio streaming broke off (%v)", err)
				return
			}
			offset = next
			a.liveMu.Lock()
			if a.live == sess {
				a.liveFed = offset
			}
			a.liveMu.Unlock()
		}
	}
}

func (a *App) takeLiveSession(gen int) (*streamSession, int) {
	a.liveMu.Lock()
	defer a.liveMu.Unlock()
	if a.live == nil || a.liveGen != gen {
		return nil, 0
	}
	ovLive.Store(false)
	sess, fed := a.live, a.liveFed
	a.live = nil
	return sess, fed
}

func (a *App) dropLiveSession() {
	ovLive.Store(false)
	a.liveMu.Lock()
	sess := a.live
	a.live = nil
	a.liveMu.Unlock()
	if sess != nil {
		sess.close()
	}
}

func (a *App) dropLiveIf(gen int) {
	a.liveMu.Lock()
	var sess *streamSession
	if a.live != nil && a.liveGen == gen {
		ovLive.Store(false)
		sess = a.live
		a.live = nil
	}
	a.liveMu.Unlock()
	if sess != nil {
		sess.close()
	}
}

func waitReadyCtx(ctx context.Context, srv recognizer) error {
	done := make(chan error, 1)
	go func() { done <- srv.waitReady(engineReadyTimeout(srv)) }()
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		log.Printf("waiting for the engine was cancelled by the user")
		return ctx.Err()
	}
}

func (a *App) engineFor(ctx context.Context, cfg *Config, want, wantModel string) (recognizer, error) {
	matches := func(r recognizer) bool {
		if r.engine() != want {
			return false
		}
		return wantModel == "" || r.external() || r.model() == wantModel
	}
	a.mu.Lock()
	primary := a.srv
	ready := a.ready
	a.mu.Unlock()
	if primary != nil && ready && matches(primary) {
		return primary, nil
	}

	a.altMu.Lock()
	defer a.altMu.Unlock()
	if a.alt != nil && matches(a.alt) && a.alt.wasStopped() == false {
		select {
		case <-a.alt.done():
			a.alt = nil
		default:
			a.altUsed = time.Now()
			return a.alt, nil
		}
	}
	if a.alt != nil {
		a.alt.stop()
		a.alt = nil
	}
	started := time.Now()
	srv, err := startEngineInstance(cfg, want, wantModel, logFile)
	if err != nil {
		return nil, err
	}
	log.Printf("starting the second engine %s on port %d for this dictation", want, recognizerPort(srv))
	if err := waitReadyCtx(ctx, srv); err != nil {
		srv.stop()
		return nil, err
	}
	a.alt = srv
	a.altUsed = time.Now()
	log.Printf("second engine %s ready in %.1f s", want, time.Since(started).Seconds())
	return srv, nil
}

func (a *App) sweepOnce() bool {
	a.mu.Lock()
	idle := time.Duration(a.cfg.EngineIdleMin) * time.Minute
	a.mu.Unlock()
	if idle <= 0 {
		return false
	}
	a.altMu.Lock()
	defer a.altMu.Unlock()
	if a.alt == nil || time.Since(a.altUsed) <= idle {
		return false
	}
	name := a.alt.engine()
	a.alt.stop()
	a.alt = nil
	log.Printf("second engine %s unloaded after %d idle minutes", name, int(idle.Minutes()))
	return true
}

func (a *App) sweepIdleEngine() {
	for {
		time.Sleep(30 * time.Second)
		a.mu.Lock()
		quitting := a.quitting
		a.mu.Unlock()
		if quitting {
			return
		}
		a.sweepOnce()
	}
}

func (a *App) stopAltEngine() {
	a.altMu.Lock()
	defer a.altMu.Unlock()
	if a.alt != nil {
		a.alt.stop()
		a.alt = nil
	}
}

var logFile *rotatingWriter

func main() {
	exe, err := os.Executable()
	if err == nil {
		_ = os.Chdir(filepath.Dir(exe))
	}
	_ = os.Setenv("WEBVIEW2_DEFAULT_BACKGROUND_COLOR", "FF0B0F0C")
	initLang("auto")
	for _, arg := range os.Args[1:] {
		if arg == "-uninstall" {
			silent := false
			for _, a := range os.Args[1:] {
				if a == "-silent" {
					silent = true
				}
			}
			runUninstall(silent)
			return
		}
	}
	setupLog()

	for i, arg := range os.Args[1:] {
		if arg == "-listmics" {
			rec, rerr := NewRecorder("")
			if rerr != nil {
				log.Printf("microphones unavailable: %v", rerr)
				return
			}
			for _, d := range rec.devices() {
				log.Printf("microphone: %s (id=%s)", d.Name, d.ID)
			}
			rec.Close()
			return
		}
		if arg == "-transcribe" && i+1 < len(os.Args[1:]) {
			path := os.Args[1:][i+1]
			cfg, cerr := loadConfig("config.json")
			if cerr != nil {
				log.Printf("config: %v", cerr)
				return
			}
			if rest := os.Args[1:][i+2:]; len(rest) > 0 && validTranslateLang(rest[0]) {
				cfg.CanaryTarget = rest[0]
				log.Printf("transcribe: translation target %s", rest[0])
			}
			wav, rerr := os.ReadFile(path)
			if rerr != nil {
				log.Printf("file %s: %v", path, rerr)
				return
			}
			srv, serr := startRecognizer(cfg, logFile)
			if serr != nil {
				log.Printf("recognizer: %s", logText(serr))
				return
			}
			defer srv.stop()
			if werr := srv.waitReady(engineReadyTimeout(srv)); werr != nil {
				log.Printf("recognizer did not start: %s", logText(werr))
				return
			}
			started := time.Now()
			text, terr := srv.transcribe(context.Background(), wav, cfg.Language, cfg.WhisperPrompt, false)
			if terr != nil {
				log.Printf("recognition: %v", terr)
				return
			}
			log.Printf("engine=%s time=%.2f s text=%q", srv.engine(), time.Since(started).Seconds(), text)
			return
		}
		if arg == "-postcheck" && i+1 < len(os.Args[1:]) {
			sample := os.Args[1:][i+1]
			cfg, cerr := loadConfig("config.json")
			if cerr != nil {
				log.Printf("config: %v", cerr)
				return
			}
			if rest := os.Args[1:][i+2:]; len(rest) > 0 && rest[0] != "" {
				enc, kerr := protectKey(rest[0])
				if kerr != nil {
					log.Printf("postcheck: encrypting the key: %v", kerr)
					return
				}
				cfg.PostAPIKey = enc
				if serr := saveConfig("config.json", cfg); serr != nil {
					log.Printf("postcheck: saving: %v", serr)
					return
				}
				log.Printf("postcheck: the key is encrypted with DPAPI and saved (%d bytes)", len(enc))
			}
			if strings.TrimSpace(cfg.PostAPIURL) == "" {
				log.Printf("postcheck: post_api_url is empty")
				return
			}
			out, perr := externalChat(context.Background(), cfg, "You repeat the user's text in upper case. Return only the text.", sample)
			log.Printf("postcheck: error=%v answer=%q", perr, out)
			return
		}
		if arg == "-streamcheck" && i+1 < len(os.Args[1:]) {
			path := os.Args[1:][i+1]
			cfg, cerr := loadConfig("config.json")
			if cerr != nil {
				log.Printf("config: %v", cerr)
				return
			}
			wav, rerr := os.ReadFile(path)
			if rerr != nil {
				log.Printf("file %s: %v", path, rerr)
				return
			}
			pcm := wav
			if len(wav) > 44 && string(wav[0:4]) == "RIFF" {
				pcm = wav[44:]
			}
			srv, serr := startEngine(cfg, engineStream, logFile)
			if serr != nil {
				log.Printf("streaming recognizer: %s", logText(serr))
				return
			}
			defer srv.stop()
			if werr := srv.waitReady(engineReadyTimeout(srv)); werr != nil {
				log.Printf("streaming recognizer did not start: %s", logText(werr))
				return
			}
			partials := 0
			ss := srv.(*streamServer)
			sess, oerr := ss.openSession(context.Background(), func(text string) {
				partials++
				log.Printf("live text: %q", livetail.Tail(text, 90))
			})
			if oerr != nil {
				log.Printf("the stream did not open: %v", oerr)
				return
			}
			chunk := sampleRate / 5 * 2
			started := time.Now()
			for off := 0; off < len(pcm); off += chunk {
				end := off + chunk
				if end > len(pcm) {
					end = len(pcm)
				}
				if perr := sess.push(pcm[off:end]); perr != nil {
					log.Printf("audio streaming: %v", perr)
					return
				}
				time.Sleep(120 * time.Millisecond)
			}
			final, ferr := sess.finish(context.Background())
			log.Printf("streamcheck: partials=%d in %.1f s, error=%v, final=%q",
				partials, time.Since(started).Seconds(), ferr, final)
			return
		}
		if arg == "-routecheck" {
			cfg, cerr := loadConfig("config.json")
			if cerr != nil {
				log.Printf("config: %v", cerr)
				return
			}
			app := &App{cfg: cfg, enabled: true, evq: evqueue.New[ptEvent](8)}
			primary := primaryEngine(cfg)
			srv, serr := startEngine(cfg, primary, logFile)
			if serr != nil {
				log.Printf("primary engine: %s", logText(serr))
				return
			}
			if werr := srv.waitReady(engineReadyTimeout(srv)); werr != nil {
				log.Printf("primary engine did not start: %s", logText(werr))
				srv.stop()
				return
			}
			app.mu.Lock()
			app.srv = srv
			app.ready = true
			app.mu.Unlock()
			log.Printf("routecheck: primary engine %s started on port %d with model %s", primary, recognizerPort(srv), srv.model())

			other, otherModel := engineWhisper, ""
			if primary == engineWhisper {
				other = engineSherpa
			}
			altCfg := cfg
			if rest := os.Args[1:][i+1:]; len(rest) > 1 && rest[1] != "" {
				c := *cfg
				c.Language = rest[1]
				m := activeModel(&c)
				if m == nil {
					log.Printf("routecheck: no model resolves for language %s", c.Language)
					srv.stop()
					return
				}
				other, otherModel = m.Engine, m.modelPath()
				altCfg = &c
				log.Printf("routecheck: language %s resolves to %s (engine %s)", c.Language, m.ID, m.Engine)
			}
			alt, aerr := app.engineFor(context.Background(), altCfg, other, otherModel)
			if aerr != nil {
				log.Printf("routecheck: second engine %s did not start: %s", other, logText(aerr))
			} else {
				if alt == srv {
					log.Printf("routecheck: the primary engine serves this language too, no second engine")
				} else {
					log.Printf("routecheck: both engines alive — %s on port %d with model %s and %s on port %d with model %s",
						srv.engine(), recognizerPort(srv), srv.model(), alt.engine(), recognizerPort(alt), alt.model())
				}
				text, terr := alt.transcribe(context.Background(), mustReadWav(os.Args[1:], i), altCfg.Language, altCfg.WhisperPrompt, false)
				if terr != nil {
					log.Printf("routecheck: the second engine did not recognize: %v", terr)
				} else {
					log.Printf("routecheck: the second engine answered %q", text)
				}
			}
			log.Printf("routecheck: holding both engines for 20 seconds to measure memory")
			time.Sleep(20 * time.Second)
			app.altMu.Lock()
			app.altUsed = time.Now().Add(-24 * time.Hour)
			app.altMu.Unlock()
			if app.sweepOnce() {
				log.Printf("routecheck: idle unloading worked")
			} else {
				log.Printf("routecheck: idle unloading did NOT work")
			}
			app.stopAltEngine()
			srv.stop()
			log.Printf("routecheck: both engines stopped")
			return
		}
		if arg == "-subcheck" {
			runSubcheck(os.Args[1:][i+1:])
			return
		}
		if arg == "-dialogs" {
			cfg, _ := loadConfig("config.json")
			if cfg != nil {
				registerBundledFonts()
				initLang(cfg.UILanguage)
				applyTheme(cfg.Skin, cfg.Theme)
				setOverlayPlacement(cfg)
			}
			log.Printf("dialog demo: focus change")
			log.Printf("answer: %q", askFocusMismatch())
			if cfg != nil {
				log.Printf("dialog demo: translation language choice")
				log.Printf("answer: %q", askTranslateTarget(cfg))
			}
			return
		}
		if arg == "-replcheck" && i+1 < len(os.Args[1:]) {
			cfg, cerr := loadConfig("config.json")
			if cerr != nil {
				log.Printf("config: %v", cerr)
				return
			}
			in := os.Args[1:][i+1]
			log.Printf("replcheck: %d replacements in the config", len(cfg.Replacements))
			after := replace.Apply(replace.ForLang(cfg.Replacements, cfg.Language), in)
			log.Printf("replcheck: replacements: %q → %q", in, after)
			cmd := commands.Apply(cfg.Commands, after)
			log.Printf("replcheck: %d commands in the config, fired %v, cancelled=%v", len(cfg.Commands), cmd.Applied, cmd.Cancelled)
			log.Printf("replcheck: result: %q", cmd.Text)
			shaped := shapeText(cfg, cfg.Language, cmd.Text, cmd.Shapes)
			log.Printf("replcheck: undo=%v shapes=%v shaped=%q", cmd.Undo, shapeActions(cmd.Shapes), shaped)
			return
		}
		if arg == "-modelcheck" {
			app := &App{}
			var out struct {
				Rows []struct {
					Name string `json:"name"`
					OK   bool   `json:"ok"`
					Note string `json:"note"`
				} `json:"rows"`
				Text string `json:"text"`
			}
			if err := json.Unmarshal([]byte(app.verifyModels()), &out); err != nil {
				log.Printf("modelcheck: %v", err)
				return
			}
			for _, r := range out.Rows {
				state := "цел"
				if !r.OK {
					state = "повреждён: " + r.Note
				}
				log.Printf("modelcheck: %s — %s", r.Name, state)
			}
			log.Printf("modelcheck: %s", out.Text)
			return
		}
		if arg == "-dpi" {
			var pt point
			procGetCursorPosDPI.Call(uintptr(unsafe.Pointer(&pt)))
			log.Printf("dpi: GetDpiForSystem=%d dpiForCursor=%d cursor=%d,%d", dpiFor(0), dpiForCursor(), pt.X, pt.Y)
			mon, _, _ := procMonitorFromPoint.Call(uintptr(uint32(pt.X))|uintptr(uint32(pt.Y))<<32, 2)
			var dx, dy uint32
			r, _, _ := procGetDpiForMonitor.Call(mon, 0, uintptr(unsafe.Pointer(&dx)), uintptr(unsafe.Pointer(&dy)))
			log.Printf("dpi: monitor=%x GetDpiForMonitor rc=%d dx=%d", mon, r, dx)
			log.Printf("dpi: MonitorFromPoint.Find=%v GetDpiForMonitor.Find=%v", procMonitorFromPoint.Find(), procGetDpiForMonitor.Find())
			return
		}
		if arg == "-overlay" {
			cfg, _ := loadConfig("config.json")
			if cfg != nil {
				registerBundledFonts()
				initLang(cfg.UILanguage)
				applyTheme(cfg.Skin, cfg.Theme)
				setOverlayPlacement(cfg)
			}
			state := ovFlashErr
			text := tr("ov.err.mic")
			rest := os.Args[1:][i+1:]
			if len(rest) > 0 && !strings.HasPrefix(rest[0], "-") {
				text = rest[0]
			}
			if len(rest) > 1 && !strings.HasPrefix(rest[1], "-") {
				switch rest[1] {
				case "rec":
					state = ovRecording
				case "proc":
					state = ovProcessing
				case "ok":
					state = ovFlashOK
				}
			}
			if len(rest) > 1 && rest[1] == "live" {
				ovLive.Store(true)
				ovRecStart.Store(time.Now().UnixMilli())
				log.Printf("live plate demo")
				words := strings.Fields(text + " " + text + " " + text + " " + text)
				shown := ""
				overlaySet(ovRecording, "")
				for _, w := range words {
					shown += w + " "
					overlaySet(ovRecording, shown)
					time.Sleep(350 * time.Millisecond)
				}
				time.Sleep(2 * time.Second)
				ovLive.Store(false)
				overlayHide()
				return
			}
			log.Printf("plate demo: %q", text)
			overlaySet(state, text)
			time.Sleep(6 * time.Second)
			overlayHide()
			return
		}

		if arg == "-miclevel" {
			cfg, _ := loadConfig("config.json")
			device := ""
			if cfg != nil {
				device = cfg.MicDevice
			}
			rec, rerr := NewRecorder(device)
			if rerr != nil {
				log.Printf("microphone unavailable: %v", rerr)
				return
			}
			if err := rec.Start(5); err != nil {
				log.Printf("the recording for analysis did not start: %v", err)
			} else {
				time.Sleep(3 * time.Second)
				pcm := rec.Stop()
				rep := audiolevel.Analyze(pcm)
				log.Printf("analysis: peak %.0f dB, RMS %.0f dB, voice %.0f%%, clipped %.2f%% → %s",
					audiolevel.DBFS(rep.Peak), audiolevel.DBFS(rep.RMS), rep.VoiceRatio*100, rep.ClipRatio*100, audiolevel.Verdict(rep))
			}
			log.Printf("level measurement without dictation, 5 seconds")
			for i := 0; i < 25; i++ {
				rec.MonitorPing()
				time.Sleep(200 * time.Millisecond)
				log.Printf("level: %.3f", rec.Level())
			}
			rec.Close()
			return
		}
	}

	for i, arg := range os.Args[1:] {
		if arg == "-quit" {
			if quitRunningInstance() {
				log.Printf("quit command sent to the running instance")
			}
			return
		}
		if arg == "-subtitles" {
			action := ""
			if rest := os.Args[1:][i+1:]; len(rest) > 0 {
				action = rest[0]
			}
			if _, ok := subCommands[action]; !ok {
				log.Printf("subtitles: use -subtitles start, stop or toggle")
				return
			}
			if subtitlesInRunningInstance(action) {
				log.Printf("subtitles: %s sent to the running instance", action)
			} else {
				log.Printf("subtitles: the app is not running")
			}
			return
		}
	}

	if !acquireSingleInstance() {
		initLang(configUILanguage("config.json"))
		if openSettingsInRunningInstance() {
			log.Printf("already running — opened the settings of the running instance")
			return
		}
		msgBox(tr("app.name"), tr("already.running"))
		return
	}

	cfg, err := loadConfig("config.json")
	if err != nil {
		msgBox(tr("cfg.err.title"), humanError(err))
		return
	}
	registerBundledFonts()
	initLang(cfg.UILanguage)
	applyTheme(cfg.Skin, cfg.Theme)
	setOverlayPlacement(cfg)

	app := &App{
		cfg:     cfg,
		enabled: true,
		evq:     evqueue.New[ptEvent](256),
		retryCh: make(chan struct{}, 1),
	}
	app.startCore()
	go cleanupWebViewProfiles()
	go app.startupUpdateCheck()
	go cleanupStaleParts()
	go app.sweepIdleEngine()
	args := os.Args[1:]
	for i, arg := range args {
		if arg == "-settings" {
			tab := "proc"
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				tab = args[i+1]
			}
			go func() {
				time.Sleep(2 * time.Second)
				app.openSettings(tab)
			}()
		}
		if arg == "-wizard" {
			cfg.WizardDone = false
		}
		if arg == "-testpaste" && i+1 < len(args) {
			text := args[i+1]
			go func() {
				time.Sleep(5 * time.Second)
				tw, _, _ := procGetForegroundWindow.Call()
				if strings.HasPrefix(text, "@mismatch ") {
					text = strings.TrimPrefix(text, "@mismatch ")
					tw = 1
				}
				log.Printf("testpaste: target captured, pasting in 6 seconds")
				time.Sleep(6 * time.Second)
				app.insertResult(context.Background(), app.snapshot(), time.Now(), text, "", tw, false)
			}()
		}
		if arg == "-testundo" && i+1 < len(args) {
			text := args[i+1]
			go func() {
				time.Sleep(5 * time.Second)
				tw, _, _ := procGetForegroundWindow.Call()
				log.Printf("testundo: target captured, pasting now, undo in 3 seconds")
				app.insertResult(context.Background(), app.snapshot(), time.Now(), text, "", tw, false)
				time.Sleep(3 * time.Second)
				log.Printf("testundo: running the undo command")
				app.undoLast(context.Background(), app.snapshot())
			}()
		}
	}
	if !cfg.WizardDone {
		log.Printf("first run: opening the setup wizard")
		go func() {
			time.Sleep(1500 * time.Millisecond)
			app.openSettings("wizard")
		}()
	}
	runTray(app)
}

func setupLog() {
	logFile = newRotatingWriter(appid.LogFile, 1<<20)
	log.SetOutput(logFile)
	log.Printf("=== %s %s started ===", appid.Slug, appVersion)
}

func (a *App) snapshot() *Config {
	a.mu.Lock()
	defer a.mu.Unlock()
	c := *a.cfg
	return &c
}

func (a *App) post(ev ptEvent) {
	if !a.evq.Push(ev) {
		log.Printf("event queue is full, event %d skipped (%d in total)", ev.kind, a.evq.Dropped())
	}
}

func (a *App) startCore() {
	a.setStatus(tr("status.loading"))
	ovOnCancel = func() { a.post(ptEvent{kind: evCancel}) }
	go a.worker()

	if _, err := parseHotkey(a.cfg.Hotkey); err != nil {
		a.fatal(fmt.Sprintf("the shortcut in the config cannot be read: %v", err), trf("err.hotkey.cfg", err.Error()))
		return
	}
	hook, err := startHotkeyHook(buildCombos(a.cfg),
		func(id string) { a.post(ptEvent{kind: evDown, profile: id}) },
		func() { a.post(ptEvent{kind: evUp}) },
		func() bool {
			if a.state.Load() == stIdle {
				return false
			}
			a.post(ptEvent{kind: evCancel})
			return true
		},
	)
	if err != nil {
		a.fatal(fmt.Sprintf("the keyboard hook was not installed: %v", err), trf("err.hook", err.Error()))
		return
	}
	a.mu.Lock()
	a.hook = hook
	a.mu.Unlock()

	go a.initBackend()
}

func (a *App) initBackend() {
	rec, err := NewRecorder(a.snapshot().MicDevice)
	if err != nil {
		a.fatal(fmt.Sprintf("the microphone was not opened: %v", err), trf("err.mic", err.Error()))
		return
	}
	a.mu.Lock()
	if a.quitting {
		a.mu.Unlock()
		rec.Close()
		return
	}
	a.rec = rec
	a.mu.Unlock()
	rec.onLost = func() { a.post(ptEvent{kind: evMicLost}) }
	ovMu.Lock()
	ovRecorder = rec
	ovMu.Unlock()

	attempts := 0
	waiting := false
	for {
		cfg := a.snapshot()
		a.mu.Lock()
		parked := a.parked
		a.mu.Unlock()
		if parked {
			log.Printf("the engine was unloaded by hand — waiting for the next dictation")
			a.setStatus(tr("status.parked"))
			if !a.waitRetry() {
				return
			}
			continue
		}
		if cfg.ServerAutostart && cfg.SttSource != "remote" {
			if missing := missingModelPath(cfg); missing != "" {
				a.mu.Lock()
				q := a.quitting
				a.mu.Unlock()
				if q {
					return
				}
				if !waiting {
					waiting = true
					langName := langLabel(strings.ToLower(strings.TrimSpace(cfg.Language)))
					if langName == "" || strings.EqualFold(langName, "auto") {
						langName = tr("route.lang.auto")
					}
					msg := trf("status.nomodel.lang", langName, modelDisplayName(presetModelID(cfg, cfg.Language)))
					log.Printf("model %s not found — waiting for the download", missing)
					a.mu.Lock()
					a.backendErr = msg
					a.mu.Unlock()
					a.setStatus(msg)
					traySetIcon(trayError)
					go a.openSettings("rec")
				}
				time.Sleep(2 * time.Second)
				continue
			}
			if waiting {
				waiting = false
				a.mu.Lock()
				a.backendErr = ""
				a.mu.Unlock()
				a.setStatus(tr("status.loading"))
			}
		}
		srv, err := startRecognizer(cfg, logFile)
		if err != nil {
			a.backendFailed(logText(err), err.Error())
			if a.waitRetry() {
				continue
			}
			return
		}
		a.mu.Lock()
		if a.quitting {
			a.mu.Unlock()
			srv.stop()
			return
		}
		a.srv = srv
		a.mu.Unlock()

		if err := srv.waitReady(engineReadyTimeout(srv)); err != nil {
			srv.stop()
			a.backendFailed(logText(err), err.Error())
			if a.waitRetry() {
				continue
			}
			return
		}

		a.mu.Lock()
		a.ready = true
		a.backendErr = ""
		a.mu.Unlock()
		a.refreshIdleUI()
		log.Printf("ready: hotkey=%s engine=%s model=%s lang=%s", cfg.Hotkey, srv.engine(), activeModelPath(cfg), cfg.Language)

		started := time.Now()
		<-srv.done()

		a.mu.Lock()
		q := a.quitting
		a.ready = false
		a.mu.Unlock()
		if q {
			return
		}
		if srv.wasStopped() {
			attempts = 0
			log.Printf("restarting the recognizer on request")
			a.setStatus(tr("status.loading"))
			continue
		}
		if srv.external() {
			a.backendFailed("the recognition server is not answering", tr("err.server.dead"))
			if a.waitRetry() {
				continue
			}
			return
		}
		if time.Since(started) > 5*time.Minute {
			attempts = 0
		}
		attempts++
		if attempts > 3 {
			a.backendFailed("the recognition server keeps crashing", tr("err.server.repeat"))
			if a.waitRetry() {
				attempts = 0
				continue
			}
			return
		}
		log.Printf("the recognizer crashed, restarting (attempt %d)", attempts)
		a.setStatus(tr("status.server.restart"))
		traySetIcon(trayError)
	}
}

func (a *App) requestServerRestart() {
	a.mu.Lock()
	a.ready = false
	a.parked = false
	srv := a.srv
	a.mu.Unlock()
	if srv != nil {
		srv.stop()
	}
	a.stopAltEngine()
	a.signalRetry()
}

func (a *App) parkEngines() {
	a.mu.Lock()
	a.ready = false
	a.parked = true
	srv := a.srv
	a.mu.Unlock()
	if srv != nil {
		srv.stop()
	}
	a.stopAltEngine()
	log.Printf("engines unloaded by hand — memory released, the next dictation starts them again")
	a.setStatus(tr("status.parked"))
	traySetIcon(trayOff)
}

func (a *App) unparkEngines() bool {
	a.mu.Lock()
	parked := a.parked
	a.parked = false
	a.mu.Unlock()
	if parked {
		a.signalRetry()
	}
	return parked
}

func (a *App) fatal(logMsg, text string) {
	log.Printf("ERROR: %s", logMsg)
	a.setStatus(text)
	traySetIcon(trayError)
	msgBox(tr("err.title"), text+tr("err.details"))
}

func (a *App) setStatus(s string) {
	traySetTooltip(tr("app.name") + " — " + s)
}

func (a *App) refreshIdleUI() {
	if a.state.Load() != stIdle {
		return
	}
	a.mu.Lock()
	enabled, ready, hotkey := a.enabled, a.ready, a.cfg.Hotkey
	a.mu.Unlock()
	switch {
	case !ready:
		traySetIcon(trayError)
	case enabled:
		traySetIcon(trayIdle)
		a.setStatus(trf("status.ready", hotkey))
	default:
		traySetIcon(trayOff)
		a.setStatus(tr("status.disabled"))
	}
}

func (a *App) toggleEnabled() {
	a.mu.Lock()
	a.enabled = !a.enabled
	a.mu.Unlock()
	a.refreshIdleUI()
}

func buildCombos(cfg *Config) []comboDef {
	var combos []comboDef
	if groups, err := parseHotkey(cfg.Hotkey); err == nil {
		combos = append(combos, comboDef{id: "main", groups: groups})
	}
	if cfg.TranslateHotkey != "" {
		if groups, err := parseHotkey(cfg.TranslateHotkey); err == nil {
			combos = append(combos, comboDef{id: "wtranslate", groups: groups})
		}
	}
	combos = append(combos, langCombos(cfg)...)
	return combos
}

func (a *App) worker() {
	for range a.evq.Signal() {
		for {
			ev, ok := a.evq.Pop()
			if !ok {
				break
			}
			switch ev.kind {
			case evDown:
				a.handleDown(ev.profile)
			case evUp:
				if a.snapshot().HotkeyMode == hotkeyToggle {
					continue
				}
				a.handleStop(0, stopKey)
			case evTimeout:
				a.handleStop(ev.gen, stopLimit)
			case evDone:
				a.handleDone(ev.gen)
			case evCancel:
				a.handleCancel()
			case evMicLost:
				a.handleMicLost()
			case evSilence:
				a.handleSilence(ev)
			}
		}
	}
}

func (a *App) handleCancel() {
	a.disarmShapes()
	switch a.state.Load() {
	case stRecording:
		a.mu.Lock()
		rec := a.rec
		a.mu.Unlock()
		rec.Stop()
		a.stopSilenceWatch()
		a.dropLiveSession()
		overlayClearDeadline()
		a.gen++
		a.state.Store(stIdle)
		log.Printf("recording cancelled by the user")
		if a.sessionCfg != nil && a.sessionCfg.Overlay {
			overlaySet(ovFlashErr, tr("ov.cancelled"))
		}
		a.refreshIdleUI()
	case stProcessing:
		askAbort()
		if a.sessionCancel != nil {
			a.sessionCancel()
		}
	}
}


func (a *App) handleDown(profileID string) {
	cfg := a.snapshot()
	if cfg.HotkeyMode == hotkeyToggle && !a.silenceStopAt.IsZero() && time.Since(a.silenceStopAt) < time.Second {
		log.Printf("press ignored: 1 s after a silence stop")
		return
	}
	if a.state.Load() == stRecording && cfg.HotkeyMode == hotkeyToggle {
		log.Printf("latched: the second press stops the recording")
		a.handleStop(0, stopKey)
		return
	}
	if a.state.Load() != stIdle {
		return
	}
	if profileID == "main" {
		profileID = ""
	}
	langKey := ""
	if strings.HasPrefix(profileID, langProfilePrefix) {
		langKey = strings.TrimPrefix(profileID, langProfilePrefix)
		cfg.Language = langKey
		profileID = ""
	}
	a.mu.Lock()
	rec := a.rec
	enabled, ready, capturing, backendErr := a.enabled, a.ready, a.capturing, a.backendErr
	ok := enabled && ready && !capturing && rec != nil
	a.mu.Unlock()
	if !ok {
		if enabled && !ready && a.unparkEngines() {
			log.Printf("the press wakes the unloaded engine")
			a.setStatus(tr("status.loading"))
			if cfg.Overlay {
				overlaySet(ovFlashErr, tr("status.loading"))
			}
			return
		}
		a.explainIgnoredPress(cfg, enabled, capturing, backendErr, rec != nil)
		return
	}
	if err := rec.Start(cfg.MaxRecordSeconds); err != nil {
		log.Printf("recording start error: %v", err)
		if cfg.Overlay {
			overlaySet(ovFlashErr, tr("ov.err.mic"))
		}
		playCue(cfg.Beep, cfg.SoundTheme, cueError)
		return
	}
	a.gen++
	a.mu.Lock()
	a.mu.Unlock()
	a.sessionTarget, _, _ = procGetForegroundWindow.Call()
	a.sessionCfg = cfg
	a.sessionProfile = profileID
	a.state.Store(stRecording)
	traySetIcon(trayRecording)
	a.setStatus(tr("status.recording"))
	playCue(cfg.Beep, cfg.SoundTheme, cueStart)
	if langKey != "" {
		log.Printf("language shortcut: this dictation is in %s", langKey)
	}
	if cfg.Overlay {
		setOverlayPlacement(cfg)
		speak := tr("ov.speak")
		if langKey != "" {
			speak += " · " + strings.ToUpper(langLabel(langKey))
		}
		overlaySet(ovRecording, speak)
	}
	if m := activeModel(cfg); m != nil && m.Engine == engineStream {
		go a.liveStream(a.gen, cfg)
	}
	if cfg.HotkeyMode == hotkeyToggle && cfg.SilenceStop {
		a.stopSilenceWatch()
		ctx, cancel := context.WithCancel(context.Background())
		a.recCancel = cancel
		go a.watchSilence(ctx, a.gen, cfg, rec)
	}

	gen := a.gen
	overlaySetDeadline(time.Now().Add(time.Duration(cfg.MaxRecordSeconds) * time.Second))
	time.AfterFunc(time.Duration(cfg.MaxRecordSeconds)*time.Second, func() {
		a.post(ptEvent{kind: evTimeout, gen: gen})
	})
}

func (a *App) handleStop(expectGen int, why stopReason) {
	if a.state.Load() != stRecording {
		return
	}
	if expectGen != 0 && expectGen != a.gen {
		return
	}
	if why == stopLimit {
		log.Printf("max_record_seconds=%d reached, stopping", a.sessionCfg.MaxRecordSeconds)
	}
	a.stopSilenceWatch()
	a.mu.Lock()
	rec := a.rec
	a.mu.Unlock()
	pcm := rec.Stop()
	if why == stopSilence {
		a.silenceStopAt = time.Now()
		pcm = trimAfterSpeech(pcm, a.silenceEv)
		a.silenceEv = ptEvent{}
	}
	overlayClearDeadline()
	a.state.Store(stProcessing)
	playCue(a.sessionCfg.Beep, a.sessionCfg.SoundTheme, cueStop)
	traySetIcon(trayProcessing)
	a.setStatus(tr("status.transcribing"))
	if a.sessionCfg.Overlay {
		procText := ""
		switch {
		case why == stopSilence:
			procText = tr("ov.stop.silence")
		case a.sessionProfile == "wtranslate":
			procText = tr("ov.translating")
		}
		overlaySet(ovProcessing, procText)
	}
	ctx, cancel := context.WithCancel(context.Background())
	a.sessionCancel = cancel
	go a.process(ctx, pcm, a.gen, a.sessionCfg, a.sessionProfile, a.sessionTarget)
}

func (a *App) handleSilence(ev ptEvent) {
	if a.state.Load() != stRecording || ev.gen != a.gen {
		return
	}
	if c := a.sessionCfg; c == nil || c.HotkeyMode != hotkeyToggle || !c.SilenceStop {
		return
	}
	a.silenceEv = ev
	a.handleStop(ev.gen, stopSilence)
}

func (a *App) stopSilenceWatch() {
	if a.recCancel != nil {
		a.recCancel()
		a.recCancel = nil
	}
}

func (a *App) watchSilence(ctx context.Context, gen int, cfg *Config, rec *Recorder) {
	pauseMs := cfg.SilenceStopMs
	if !audiolevel.ValidPauseMs(pauseMs) {
		pauseMs = audiolevel.DefaultPauseMs
	}
	guard := 300 * time.Millisecond
	if cfg.Beep {
		guard = 800 * time.Millisecond
	}
	ep := audiolevel.NewEndpoint(sampleRate, time.Duration(pauseMs)*time.Millisecond, guard)
	log.Printf("stop on silence: on, pause %d ms", pauseMs)
	unheard := func() {
		if !ep.Heard() {
			log.Printf("stop on silence: no speech heard (floor %.0f dB)", ep.Floor())
		}
	}
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()
	offset := 0
	for {
		select {
		case <-ctx.Done():
			unheard()
			return
		case <-tick.C:
		}
		if a.state.Load() != stRecording {
			unheard()
			return
		}
		chunk, next := rec.TakeFrom(offset)
		offset = next
		if ep.Feed(chunk) {
			a.post(ptEvent{kind: evSilence, gen: gen, cut: ep.SpeechEnd() + sampleRate*2*600/1000, ep: ep})
			return
		}
	}
}

func trimAfterSpeech(pcm []byte, ev ptEvent) []byte {
	perSec := float64(sampleRate * 2)
	if ep := ev.ep; ep != nil {
		log.Printf("stop on silence: %.1f s quiet after %.1f s of speech (floor %.0f dB, threshold %.0f dB), stopping",
			float64(len(pcm)-ep.SpeechEnd())/perSec, float64(ep.SpeechMs())/1000, ep.Floor(), ep.Threshold())
	}
	cut := ev.cut &^ 1
	if cut <= 0 || cut >= len(pcm) {
		return pcm
	}
	log.Printf("stop on silence: %.1f s of trailing silence dropped", float64(len(pcm)-cut)/perSec)
	return pcm[:cut]
}

func (a *App) handleDone(gen int) {
	if a.state.Load() != stProcessing || gen != a.gen {
		return
	}
	a.state.Store(stIdle)
	a.refreshIdleUI()
}

var onlyNoise = regexp.MustCompile(`^[\s.,!?\*\[\]\(\)«»"'…·♪♫~-]*$|^\[[^\]]*\]$|^\([^\)]*\)$`)

func (a *App) process(ctx context.Context, pcm []byte, gen int, cfg *Config, profileID string, targetWnd uintptr) {
	defer a.post(ptEvent{kind: evDone, gen: gen})
	defer a.dropLiveIf(gen)

	minBytes := sampleRate * 2 * cfg.MinRecordMs / 1000
	if len(pcm) < minBytes {
		log.Printf("recording too short (%d ms), skipping", len(pcm)*1000/(sampleRate*2))
		if cfg.Overlay {
			overlaySet(ovFlashErr, tr("ov.tooshort"))
		}
		return
	}
	sound := audiolevel.Analyze(pcm)
	verdict := audiolevel.Verdict(sound)
	a.mu.Lock()
	a.lastVerdict = verdict
	a.mu.Unlock()
	log.Printf("audio: peak %.0f dB, voice %.0f%%, clipped %.1f%% — %s",
		audiolevel.DBFS(sound.Peak), sound.VoiceRatio*100, sound.ClipRatio*100, verdict)
	if verdict == audiolevel.VerdictSilent {
		log.Printf("silence in the recording — recognition skipped")
		if cfg.Overlay {
			overlaySet(ovFlashErr, tr("ov.silence"))
		}
		return
	}

	a.mu.Lock()
	srv := a.srv
	a.mu.Unlock()
	if srv == nil {
		log.Printf("the recognition server is not up yet")
		if cfg.Overlay {
			overlaySet(ovFlashErr, tr("ov.server.loading"))
		}
		return
	}

	start := time.Now()
	wantTranslate := profileID == "wtranslate" || (profileID == "" && (cfg.TranslateDefault || cfg.TranslateAsk != "never"))
	target := ""
	if wantTranslate {
		if cfg.TranslateDefault || cfg.TranslateAsk == "never" {
			target = cfg.TranslateTarget
		} else {
			target = askTranslateTarget(cfg)
		}
		if ctx.Err() != nil {
			log.Printf("recognition cancelled by the user")
			if cfg.Overlay {
				overlaySet(ovFlashErr, tr("ov.cancelled"))
			}
			return
		}
		switch {
		case target == "":
			log.Printf("translation: cancelled in the dialog, pasting as is")
		case target == cfg.Language:
			log.Printf("translation: target %s matches the recognition language — skipping", target)
			target = ""
		default:
			log.Printf("translation: target=%s by Whisper (mode %s)", target, cfg.TranslateAsk)
		}
	}
	active := activeModel(cfg)
	wantEngine, wantModel := primaryEngine(cfg), ""
	if active != nil {
		wantModel = active.modelPath()
	}
	viaModel := false
	if target != "" && modelTranslates(active, cfg.Language, target) {
		if active.TrLangs != "" {
			viaModel = true
			wantModel = active.modelPath() + "#" + target
			log.Printf("translation: %s translates on its own, target %s", active.NameKey, target)
		}
	} else if target != "" {
		name := ""
		if active != nil {
			name = active.NameKey
		}
		fallback := bestInstalledWhisper()
		switch {
		case fallback == nil:
			log.Printf("translation unavailable: %s does not translate and no Whisper is installed — pasting as is", name)
			if cfg.Overlay {
				overlayNote(tr("ov.notranslate"))
			}
			target = ""
		case askTranslateFallback(name, fallback.NameKey):
			log.Printf("translation: %s does not translate, the user picked %s", name, fallback.NameKey)
			wantEngine, wantModel = engineWhisper, fallback.modelPath()
		default:
			log.Printf("translation: %s does not translate, the user declined Whisper — pasting as is", name)
			target = ""
		}
		if ctx.Err() != nil {
			log.Printf("recognition cancelled by the user")
			if cfg.Overlay {
				overlaySet(ovFlashErr, tr("ov.cancelled"))
			}
			return
		}
	}
	if wantEngine != srv.engine() || (wantModel != "" && !srv.external() && srv.model() != wantModel) {
		alt, aerr := a.engineFor(ctx, cfg, wantEngine, wantModel)
		if aerr != nil {
			log.Printf("engine %s did not start (%v) — staying on %s", wantEngine, aerr, srv.engine())
			if cfg.Overlay {
				overlayNote(tr("ov.engine.fallback"))
			}
			viaModel = viaModel && srv.model() == wantModel
		} else {
			srv = alt
		}
	}
	log.Printf("preset: engine=%s model=%s language=%s translation=%v", srv.engine(), srv.model(), cfg.Language, target != "")
	if target != "" && !viaModel && !engineTranslates(srv.engine()) {
		log.Printf("translation unavailable: engine %s is active — pasting the recognized text as is", srv.engine())
		if cfg.Overlay {
			overlayNote(tr("ov.notranslate"))
		}
		target = ""
	}
	fastTranslate := target == "en" && !viaModel
	recLang := cfg.Language
	if target != "" && !fastTranslate && !viaModel {
		recLang = target
	}
	var text string
	var err error
	liveSess, liveFed := a.takeLiveSession(gen)
	if liveSess != nil && (srv.engine() != engineStream || target != "") {
		liveSess.close()
		liveSess = nil
	}
	if liveSess != nil {
		if liveFed < len(pcm) {
			_ = liveSess.push(pcm[liveFed:])
		}
		text, err = liveSess.finish(ctx)
		if err != nil && ctx.Err() == nil {
			log.Printf("live text did not reach the end (%v) — recognizing in batch", err)
			liveSess, err = nil, nil
		} else {
			log.Printf("live text: the final result arrived")
		}
	}
	if liveSess == nil && err == nil {
		prompt := cfg.WhisperPrompt
		if target != "" && !viaModel && isBuiltinDictionary(prompt) {
			prompt = builtinDictionary(target)
			log.Printf("dictionary: hinting the %s set for the translation", target)
		}
		text, err = srv.transcribe(ctx, wavFromPCM16(pcm, sampleRate), recLang, prompt, fastTranslate)
	}
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Printf("recognition cancelled by the user")
			if cfg.Overlay {
				overlaySet(ovFlashErr, tr("ov.cancelled"))
			}
			return
		}
		log.Printf("recognition: %v", err)
		playCue(cfg.Beep, cfg.SoundTheme, cueError)
		if cfg.Overlay {
			overlaySet(ovFlashErr, tr("ov.err.recognize"))
		}
		return
	}
	if ctx.Err() != nil {
		log.Printf("recognition cancelled by the user")
		if cfg.Overlay {
			overlaySet(ovFlashErr, tr("ov.cancelled"))
		}
		return
	}
	text = strings.TrimSpace(text)
	if joined := segments.Join(text); joined != text {
		log.Printf("recognition: %d line break(s) from the engine joined into spaces", strings.Count(text, "\n"))
		text = joined
	}
	outLang := cfg.Language
	if target != "" {
		outLang = target
	}
	if fixed := replace.Apply(replace.ForLang(cfg.Replacements, outLang), text); fixed != text {
		log.Printf("replacements: %d → %d characters", len([]rune(text)), len([]rune(fixed)))
		text = fixed
	}
	var shapes []commands.Command
	if cmd := commands.Apply(cfg.Commands, text); len(cmd.Applied) > 0 {
		log.Printf("commands %v: %d → %d characters (cancelled=%v)", cmd.Applied, len([]rune(text)), len([]rune(cmd.Text)), cmd.Cancelled)
		if cmd.Cancelled {
			a.disarmShapes()
			if cfg.Overlay {
				overlaySet(ovFlashErr, tr("ov.cmd.cancelled"))
			}
			return
		}
		if cmd.Undo {
			a.undoLast(ctx, cfg)
			return
		}
		if len(cmd.Shapes) > 0 && (strings.TrimSpace(cmd.Text) == "" || onlyNoise.MatchString(cmd.Text)) {
			a.armShapes(cfg, cmd.Shapes)
			return
		}
		text = cmd.Text
		shapes = cmd.Shapes
	}
	if text == "" || onlyNoise.MatchString(text) {
		log.Printf("empty result (%d characters), nothing to paste", len([]rune(text)))
		if cfg.Overlay {
			switch verdict {
			case audiolevel.VerdictClipped:
				overlaySet(ovFlashErr, tr("ov.clipped"))
			case audiolevel.VerdictQuiet:
				overlaySet(ovFlashErr, tr("ov.quiet"))
			default:
				overlaySet(ovFlashErr, tr("ov.silence"))
			}
		}
		return
	}

	skipped := ""
	if chain := punctChain(cfg, chainProfiles(cfg, profileID)); len(chain) > 0 && postReady(cfg) {
		for i, prof := range chain {
			label := prof.Name
			if len(chain) > 1 {
				label = fmt.Sprintf("%s (%d/%d)", prof.Name, i+1, len(chain))
			}
			if cfg.Overlay {
				overlaySet(ovProcessing, trf("ov.editing", label))
			}
			out, lerr := a.llmProcess(ctx, prof.Prompt, text)
			if lerr == nil {
				text = out
				a.setPostErr("", "")
				continue
			}
			if errors.Is(lerr, context.Canceled) || ctx.Err() != nil {
				log.Printf("editing cancelled by the user")
				if cfg.Overlay {
					overlaySet(ovFlashErr, tr("ov.cancelled"))
				}
				return
			}
			log.Printf("profile %s did not apply, continuing with the current text: %v", prof.ID, lerr)
			a.setPostErr(prof.Name, errText(lerr))
			if skipped == "" {
				skipped = prof.Name
			}
		}
	}

	if shapes = a.takeArmedShapes(shapes); len(shapes) > 0 {
		log.Printf("shapes: %s", strings.Join(shapeActions(shapes), ", "))
	}
	text = shapeText(cfg, outLang, text, shapes)
	if strings.TrimSpace(text) == "" {
		log.Printf("empty result after shaping, nothing to paste")
		return
	}
	a.insertResult(ctx, cfg, start, text, skipped, targetWnd, liveSess != nil)
}

const shapeArmTTL = 60 * time.Second

func shapeText(cfg *Config, lang, text string, shapes []commands.Command) string {
	if len(shapes) > 0 {
		return commands.Shape(text, shapes, commands.ShapeOptions{Lang: lang, StripPunct: cfg.Punctuation == punctOff})
	}
	if cfg.Punctuation == punctOff {
		return textshape.StripPunctuation(text)
	}
	return text
}

func shapeActions(shapes []commands.Command) []string {
	out := make([]string, 0, len(shapes))
	for _, s := range shapes {
		out = append(out, s.Action)
	}
	return out
}

func shapeLabel(action string) string {
	key := "cmd." + action
	if s := tr(key); s != key {
		return s
	}
	return action
}

func (a *App) armShapes(cfg *Config, shapes []commands.Command) {
	a.mu.Lock()
	var armed []commands.Command
	if len(a.armedShapes) > 0 && time.Since(a.armedAt) <= shapeArmTTL {
		for _, old := range a.armedShapes {
			kept := true
			for _, s := range shapes {
				if s.Action == old.Action {
					kept = false
				}
			}
			if kept {
				armed = append(armed, old)
			}
		}
	}
	armed = append(armed, shapes...)
	a.armedShapes = armed
	a.armedAt = time.Now()
	a.mu.Unlock()
	names := shapeActions(armed)
	log.Printf("shapes armed: %s", strings.Join(names, ", "))
	if cfg.Overlay {
		labels := make([]string, len(names))
		for i, n := range names {
			labels[i] = shapeLabel(n)
		}
		overlaySet(ovFlashOK, trf("ov.shape.next", strings.Join(labels, ", ")))
	}
}

func (a *App) takeArmedShapes(shapes []commands.Command) []commands.Command {
	a.mu.Lock()
	armed, at := a.armedShapes, a.armedAt
	a.armedShapes = nil
	a.mu.Unlock()
	if len(armed) == 0 {
		return shapes
	}
	if time.Since(at) > shapeArmTTL {
		log.Printf("shapes: the armed ones expired, not applied")
		return shapes
	}
	return append(append([]commands.Command{}, armed...), shapes...)
}

func (a *App) disarmShapes() {
	a.mu.Lock()
	had := len(a.armedShapes) > 0
	a.armedShapes = nil
	a.mu.Unlock()
	if had {
		log.Printf("shapes: the armed ones were dropped")
	}
}

type insertion struct {
	commands.UndoRecord
	histAt int64
}

func (a *App) clearUndo() {
	a.mu.Lock()
	a.lastIns = nil
	a.mu.Unlock()
}

func focusOf(wnd uintptr) uintptr {
	if wnd == 0 {
		return 0
	}
	tid, _, _ := procGetWindowThreadPID.Call(wnd, 0)
	if tid == 0 {
		return 0
	}
	var gui guiThreadInfo
	gui.Size = uint32(unsafe.Sizeof(gui))
	if ok, _, _ := procGetGUIThreadInfo.Call(tid, uintptr(unsafe.Pointer(&gui))); ok == 0 {
		return 0
	}
	return gui.Focus
}

func (a *App) recordInsertion(text string, wnd uintptr, mode string, enterSent bool, histAt int64) {
	if wnd == 0 || wnd == 1 {
		wnd, _, _ = procGetForegroundWindow.Call()
	}
	if wnd == 0 {
		a.clearUndo()
		return
	}
	a.mu.Lock()
	hook := a.hook
	a.mu.Unlock()
	ins := &insertion{
		UndoRecord: commands.UndoRecord{
			Text:      text,
			Wnd:       wnd,
			Focus:     focusOf(wnd),
			At:        time.Now(),
			Typed:     hook.Typed(),
			Clicked:   hook.Clicked(),
			EnterSent: enterSent,
			Mode:      mode,
		},
		histAt: histAt,
	}
	a.mu.Lock()
	a.lastIns = ins
	a.mu.Unlock()
}

func (a *App) undoLast(ctx context.Context, cfg *Config) {
	a.mu.Lock()
	ins := a.lastIns
	hook := a.hook
	a.mu.Unlock()
	var rec *commands.UndoRecord
	if ins != nil {
		rec = &ins.UndoRecord
	}
	now := time.Now()
	wnd, _, _ := procGetForegroundWindow.Call()
	if ok, reason := commands.CanUndo(rec, now, wnd, focusOf(wnd), hook.Typed(), hook.Clicked()); !ok {
		if rec != nil {
			log.Printf("undo: not done, %s (keys +%d, clicks +%d, %.1fs after insertion)", reason, hook.Typed()-rec.Typed, hook.Clicked()-rec.Clicked, now.Sub(rec.At).Seconds())
		} else {
			log.Printf("undo: not done, %s", reason)
		}
		if cfg.Overlay {
			overlaySet(ovFlashErr, tr("ov.undo."+reason))
		}
		return
	}
	age := now.Sub(rec.At)
	waitModifiersReleased(3 * time.Second)
	n := commands.Backspaces(rec.Text)
	done, err := eraseChars(ctx, n, rec.Wnd)
	if err != nil {
		a.mu.Lock()
		if a.lastIns == ins {
			a.lastIns = nil
		}
		a.mu.Unlock()
		switch {
		case errors.Is(err, errFocusMoved):
			log.Printf("undo: erasing stopped after %d of %d, the window changed", done, n)
			if cfg.Overlay {
				overlaySet(ovFlashErr, tr("ov.undo.moved"))
			}
		case errors.Is(err, context.Canceled):
			log.Printf("undo: erasing cancelled by the user after %d of %d", done, n)
			if cfg.Overlay {
				overlaySet(ovFlashErr, tr("ov.cancelled"))
			}
		default:
			log.Printf("undo: erasing stopped after %d of %d: %v", done, n, err)
			if cfg.Overlay {
				overlaySet(ovFlashErr, tr("ov.err.paste"))
			}
		}
		return
	}
	a.mu.Lock()
	if a.lastIns == ins {
		a.lastIns = nil
	}
	a.mu.Unlock()
	if ins.histAt != 0 {
		if herr := histStore.Delete(ins.histAt); herr != nil {
			log.Printf("history: %v", herr)
		}
	}
	log.Printf("undo: erased %d characters in %s (%s, %.1fs after insertion)", n, processNameOf(rec.Wnd), rec.Mode, age.Seconds())
	if cfg.Overlay {
		overlaySet(ovFlashOK, trf("ov.undo.done", n))
	}
}

func (a *App) insertResult(ctx context.Context, cfg *Config, start time.Time, text, skipped string, targetWnd uintptr, liveShown bool) {
	a.mu.Lock()
	a.lastResult = text
	a.lastResultAt = time.Now()
	a.lastWnd = targetWnd
	a.lastTarget = windowTitle(targetWnd)
	a.lastProcess = processNameOf(targetWnd)
	a.lastIns = nil
	targetApp := a.lastProcess
	a.mu.Unlock()
	histAt := a.rememberDictation(cfg, text, targetApp)

	allowEnter := cfg.AutoEnter && skipped == ""
	expect := targetWnd
	if targetWnd != 0 {
		cur, _, _ := procGetForegroundWindow.Call()
		if cur != targetWnd {
			log.Printf("focus changed during processing — asking")
			switch askFocusMismatch() {
			case "here":
				allowEnter = false
				expect, _, _ = procGetForegroundWindow.Call()
			case "copy":
				if err := setClipboardText(text); err != nil {
					log.Printf("copying: %v", err)
					return
				}
				log.Printf("the result was copied to the clipboard at the user request")
				if cfg.Overlay {
					overlaySet(ovFlashOK, tr("ov.copied"))
				}
				return
			default:
				_ = setClipboardText(text)
				log.Printf("paste cancelled, the text is kept as the last result and in the clipboard")
				if cfg.Overlay {
					overlaySet(ovFlashErr, tr("ov.kept"))
				}
				return
			}
			if ctx.Err() != nil {
				return
			}
		}
	}

	if ctx.Err() != nil {
		log.Printf("paste cancelled by the user before it started")
		return
	}
	mode, err := pasteTextMode(cfg, text, expect)
	if err != nil {
		playCue(cfg.Beep, cfg.SoundTheme, cueError)
		if errors.Is(err, errFocusMoved) {
			log.Printf("paste cancelled: the input window changed, the text is kept")
			_ = setClipboardText(text)
			if cfg.Overlay {
				overlaySet(ovFlashErr, tr("ov.moved"))
			}
			return
		}
		log.Printf("paste: %v", err)
		if cfg.Overlay {
			overlaySet(ovFlashErr, tr("ov.err.paste"))
		}
		return
	}
	enterSent := false
	if allowEnter && ctx.Err() == nil {
		time.Sleep(150 * time.Millisecond)
		if ctx.Err() != nil {
			log.Printf("auto-enter cancelled by the user")
		} else if !focusStillOn(expect) {
			log.Printf("auto-enter cancelled: the input window changed after pasting")
		} else {
			enterSent = true
			if err := pressEnter(); err != nil {
				log.Printf("auto-enter: %v", err)
			}
		}
	}
	a.recordInsertion(text, expect, mode, enterSent, histAt)
	log.Printf("done in %.1fs: %d characters", time.Since(start).Seconds(), len([]rune(text)))
	if cfg.Overlay {
		switch {
		case skipped != "":
			overlaySet(ovFlashErr, trf("ov.llm.skipped", skipped))
		case liveShown:
			overlayHide()
		case cfg.OverlayText && strings.TrimSpace(text) != "":
			overlaySet(ovFlashOK, oneLine(text))
		default:
			overlaySet(ovFlashOK, trf("ov.inserted", len([]rune(text))))
		}
	}
}

func oneLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

var histStore = history.Open(appid.HistoryFile)

func (a *App) rememberDictation(cfg *Config, text, app string) int64 {
	if !cfg.HistoryOn || strings.TrimSpace(text) == "" {
		return 0
	}
	if skip := strings.TrimSpace(cfg.HistorySkip); skip != "" && app != "" {
		if _, hit := apprules.Find([]apprules.Rule{{Match: skip}}, app); hit {
			log.Printf("history: %s is on the exclusion list, not recording", app)
			return 0
		}
	}
	now := time.Now()
	item := history.Item{At: now.UnixMilli(), Text: text, App: app}
	if err := histStore.Add(item, now.UnixMilli(), cfg.HistoryKeepMin, cfg.HistoryMax); err != nil {
		log.Printf("history: %v", err)
	}
	return item.At
}

func (a *App) changeHotkey() {
	a.mu.Lock()
	hook := a.hook
	if a.capturing || hook == nil {
		a.mu.Unlock()
		return
	}
	a.capturing = true
	a.mu.Unlock()
	defer func() {
		a.mu.Lock()
		a.capturing = false
		a.mu.Unlock()
	}()

	current := a.snapshot().Hotkey
	combo, ok := captureHotkeyDialog(hook, current)
	if !ok {
		log.Printf("hotkey change cancelled")
		return
	}
	if _, err := parseHotkey(combo); err != nil {
		log.Printf("captured hotkey %q could not be parsed: %v", combo, err)
		return
	}
	a.mu.Lock()
	c := *a.cfg
	c.Hotkey = combo
	a.cfg = &c
	a.mu.Unlock()
	hook.SetCombos(buildCombos(&c))
	if err := saveConfig("config.json", &c); err != nil {
		log.Printf("saving the config: %v", err)
	}
	a.refreshIdleUI()
	log.Printf("new hotkey: %s", combo)
	if warn := hotkeyWarning(combo); warn != "" {
		log.Printf("warning: %s", warn)
		overlayNote(warn)
	}
}

func chainProfiles(cfg *Config, profileID string) []*Profile {
	chain := profiles.Chain(cfg.Profiles, cfg.ActiveProfiles, profileID)
	out := make([]*Profile, 0, len(chain))
	for i := range chain {
		out = append(out, &chain[i])
	}
	return out
}

func (a *App) llmProcess(ctx context.Context, prompt, text string) (string, error) {
	if cfg := a.snapshot(); postAPIOn(cfg) {
		out, err := externalChat(ctx, cfg, prompt, text)
		if err != nil {
			return "", err
		}
		out = strings.TrimSpace(out)
		if out == "" {
			return "", errors.New("empty answer from the post-processing server")
		}
		return out, nil
	}
	llm, err := a.ensureLLM()
	if err != nil {
		return "", err
	}
	tctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	out, err := llm.chat(tctx, prompt, text)
	if err != nil {
		return "", err
	}
	out = strings.TrimSpace(out)
	if out == "" {
		return "", errors.New("empty answer from the LLM")
	}
	return out, nil
}

func (a *App) reloadConfig() {
	fresh, err := loadConfig("config.json")
	if err != nil {
		log.Printf("reloading the config: %v", err)
		a.setStatus(tr("status.cfg.err"))
		return
	}
	if _, err := parseHotkey(fresh.Hotkey); err != nil {
		log.Printf("reloading the config: %v", err)
		a.setStatus(tr("status.cfg.err"))
		return
	}

	a.mu.Lock()
	old := a.cfg
	a.cfg = fresh
	hook := a.hook
	a.mu.Unlock()
	initLang(fresh.UILanguage)

	canaryLangSwitch := false
	if fresh.Language != old.Language {
		if am := activeModel(fresh); am != nil && am.TrLangs != "" {
			canaryLangSwitch = true
		}
	}
	serverChanged := canaryLangSwitch ||
		fresh.Model != old.Model ||
		fresh.SherpaModel != old.SherpaModel ||
		primaryEngine(fresh) != primaryEngine(old) ||
		fresh.ServerPort != old.ServerPort ||
		fresh.Threads != old.Threads ||
		fresh.SherpaThreads != old.SherpaThreads ||
		fresh.ServerURL != old.ServerURL ||
		fresh.SttSource != old.SttSource ||
		fresh.ServerExe != old.ServerExe ||
		fresh.ServerAutostart != old.ServerAutostart
	if serverChanged {
		a.requestServerRestart()
	}

	if hook != nil {
		hook.SetCombos(buildCombos(fresh))
	} else {
		log.Printf("the hook is not installed, the new hotkey applies after a restart")
	}

	if serverChanged {
		a.setStatus(tr("status.loading"))
		log.Printf("config reloaded; the recognizer restarts with the new settings")
	} else {
		a.refreshIdleUI()
		log.Printf("config reloaded: hotkey=%s paste=%s", fresh.Hotkey, fresh.PasteMode)
	}
	a.subtitlesConfigChanged()
}

func (a *App) onExit() {
	a.stopSubtitles("the app is closing", subExitDrain)
	a.mu.Lock()
	a.quitting = true
	srv := a.srv
	rec := a.rec
	llm := a.llm
	hook := a.hook
	a.mu.Unlock()
	hook.release()
	if srv != nil {
		srv.stop()
	}
	a.stopAltEngine()
	if llm != nil {
		llm.stop()
	}
	if rec != nil {
		rec.Close()
	}
	log.Printf("exit")
	if logFile != nil {
		logFile.Close()
	}
}

func mustReadWav(args []string, i int) []byte {
	if i+1 < len(args) {
		if b, err := os.ReadFile(args[i+1]); err == nil {
			return b
		}
	}
	return wavFromPCM16(make([]byte, sampleRate/5*2), sampleRate)
}

func (a *App) setPostErr(profile, msg string) {
	a.mu.Lock()
	a.postErr = msg
	a.postErrProf = profile
	if msg == "" {
		a.postErrAt = time.Time{}
	} else {
		a.postErrAt = time.Now()
	}
	a.mu.Unlock()
}

func punctChain(cfg *Config, chain []*Profile) []*Profile {
	if cfg.Punctuation != punctByLLM {
		return chain
	}
	head := &Profile{ID: "punct", Name: strS("S_PUNCT"), Prompt: tr("punct.prompt")}
	return append([]*Profile{head}, chain...)
}

func (a *App) micCheck() string {
	type out struct {
		Verdict string  `json:"verdict"`
		Text    string  `json:"text"`
		PeakDB  float64 `json:"peak_db"`
		Voice   float64 `json:"voice"`
		Clip    float64 `json:"clip"`
	}
	fail := func(msg string) string {
		b, _ := json.Marshal(out{Verdict: "error", Text: msg})
		return string(b)
	}
	if a.state.Load() != stIdle {
		return fail(tr("mic.busy"))
	}
	a.mu.Lock()
	rec := a.rec
	a.mu.Unlock()
	if rec == nil {
		return fail(tr("ov.err.mic"))
	}
	if err := rec.Start(5); err != nil {
		log.Printf("microphone check: %v", err)
		return fail(humanError(err))
	}
	time.Sleep(3 * time.Second)
	pcm := rec.Stop()
	rep := audiolevel.Analyze(pcm)
	verdict := audiolevel.Verdict(rep)
	peak := audiolevel.DBFS(rep.Peak)
	log.Printf("microphone check: peak %.0f dB, voice %.0f%%, clipped %.1f%% — %s",
		peak, rep.VoiceRatio*100, rep.ClipRatio*100, verdict)
	text := ""
	switch verdict {
	case audiolevel.VerdictSilent:
		text = tr("mic.check.silent")
	case audiolevel.VerdictClipped:
		text = trf("mic.check.clipped", rep.ClipRatio*100)
	case audiolevel.VerdictQuiet:
		text = trf("mic.check.quiet", peak)
	default:
		text = trf("mic.check.ok", peak, rep.VoiceRatio*100)
	}
	b, _ := json.Marshal(out{Verdict: verdict, Text: text, PeakDB: peak, Voice: rep.VoiceRatio, Clip: rep.ClipRatio})
	return string(b)
}

func (a *App) backendFailed(logMsg, text string) {
	log.Printf("recognizer did not start: %s", logMsg)
	a.mu.Lock()
	a.ready = false
	a.backendErr = text
	a.srv = nil
	a.mu.Unlock()
	a.setStatus(text)
	traySetIcon(trayError)
	if cfg := a.snapshot(); cfg != nil && cfg.Overlay {
		overlaySet(ovFlashErr, text)
	}
}

func (a *App) waitRetry() bool {
	log.Printf("waiting for the settings to be fixed or for the retry button")
	for {
		select {
		case <-a.retryCh:
			log.Printf("trying to start the recognizer again")
			a.mu.Lock()
			a.backendErr = ""
			a.mu.Unlock()
			a.setStatus(tr("status.loading"))
			return true
		case <-time.After(time.Second):
			a.mu.Lock()
			q := a.quitting
			a.mu.Unlock()
			if q {
				return false
			}
		}
	}
}

func (a *App) signalRetry() {
	if a.retryCh == nil {
		return
	}
	select {
	case a.retryCh <- struct{}{}:
	default:
	}
}

func (a *App) handleMicLost() {
	a.mu.Lock()
	rec := a.rec
	cfg := a.sessionCfg
	a.mu.Unlock()
	if a.state.Load() == stRecording && rec != nil {
		rec.Stop()
		a.stopSilenceWatch()
		a.dropLiveSession()
		overlayClearDeadline()
		a.gen++
		a.state.Store(stIdle)
		log.Printf("microphone disconnected during the recording — recording aborted")
		if cfg != nil {
			playCue(cfg.Beep, cfg.SoundTheme, cueError)
			if cfg.Overlay {
				overlaySet(ovFlashErr, tr("ov.mic.lost"))
			}
		}
	} else {
		log.Printf("microphone disconnected")
	}
	if rec != nil {
		if err := rec.SetDevice(""); err != nil {
			log.Printf("falling back to the default microphone: %v", err)
		} else {
			log.Printf("switched to the default microphone")
			a.mu.Lock()
			if a.cfg != nil && a.cfg.MicDevice != "" {
				c := *a.cfg
				c.MicDevice = ""
				c.MicDeviceName = ""
				a.cfg = &c
				_ = saveConfig("config.json", &c)
			}
			a.mu.Unlock()
		}
	}
	a.refreshIdleUI()
}

func (a *App) explainIgnoredPress(cfg *Config, enabled, capturing bool, backendErr string, haveRec bool) {
	switch {
	case capturing:
		return
	case !enabled:
		log.Printf("press ignored: the app is switched off in the tray")
		return
	case backendErr != "":
		log.Printf("press ignored: the recognizer did not start")
		if cfg.Overlay {
			overlaySet(ovFlashErr, backendErr)
		}
		playCue(cfg.Beep, cfg.SoundTheme, cueError)
	case !haveRec:
		log.Printf("press ignored: the microphone is unavailable")
		if cfg.Overlay {
			overlaySet(ovFlashErr, tr("ov.err.mic"))
		}
		playCue(cfg.Beep, cfg.SoundTheme, cueError)
	default:
		log.Printf("press ignored: the recognizer is still getting ready")
		if cfg.Overlay {
			overlaySet(ovFlashErr, tr("status.loading"))
		}
	}
}

func (a *App) resetSettings() {
	a.mu.Lock()
	old := *a.cfg
	a.mu.Unlock()

	fresh := defaultConfig()
	fresh.Profiles = old.Profiles
	fresh.ActiveProfiles = old.ActiveProfiles
	fresh.Replacements = old.Replacements
	fresh.Commands = old.Commands
	fresh.Model = old.Model
	fresh.SherpaModel = old.SherpaModel
	fresh.LangModels = old.LangModels
	fresh.LLMModel = old.LLMModel
	fresh.WizardDone = true
	fresh.UILanguage = old.UILanguage
	syncDictionary(fresh)

	if err := saveConfig("config.json", fresh); err != nil {
		log.Printf("resetting the settings: %v", err)
		return
	}
	log.Printf("settings reset to factory defaults (models, history and prompts kept)")
	a.reloadConfig()
}
