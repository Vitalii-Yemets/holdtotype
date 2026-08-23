package main

import (
	"encoding/json"
	"holdtotype/internal/audiolevel"
	"holdtotype/internal/history"
	"holdtotype/internal/replace"
	"holdtotype/internal/apprules"
	"unsafe"
	"holdtotype/internal/evqueue"

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
)

type ptEvent struct {
	kind    int
	gen     int
	profile string
}

type App struct {
	mu        sync.Mutex
	cfg       *Config
	rec       *Recorder
	srv       recognizer
	llm       *llamaServer
	hook      *hotkeyHook
	enabled   bool
	ready     bool
	capturing bool
	quitting  bool

	evq *evqueue.Queue[ptEvent]

	state          atomic.Int32
	gen            int
	sessionCfg     *Config
	sessionProfile string
	sessionCancel  context.CancelFunc
	sessionTarget  uintptr

	lastResult   string
	lastResultAt time.Time
	lastTarget   string
	lastProcess  string
	updVer     string
	updURL     string

	altMu   sync.Mutex
	alt     recognizer
	altUsed time.Time
}

func (a *App) engineFor(cfg *Config, want string) (recognizer, error) {
	a.mu.Lock()
	primary := a.srv
	ready := a.ready
	a.mu.Unlock()
	if primary != nil && ready && primary.engine() == want {
		return primary, nil
	}

	a.altMu.Lock()
	defer a.altMu.Unlock()
	if a.alt != nil && a.alt.engine() == want && a.alt.wasStopped() == false {
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
	log.Printf("поднимаю второй движок %s под эту диктовку", want)
	started := time.Now()
	srv, err := startEngine(cfg, want, logFile)
	if err != nil {
		return nil, err
	}
	if err := srv.waitReady(engineReadyTimeout(srv)); err != nil {
		srv.stop()
		return nil, err
	}
	a.alt = srv
	a.altUsed = time.Now()
	log.Printf("второй движок %s готов за %.1f с", want, time.Since(started).Seconds())
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
	log.Printf("второй движок %s выгружен после %d минут простоя", name, int(idle.Minutes()))
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
				log.Printf("микрофоны недоступны: %v", rerr)
				return
			}
			for _, d := range rec.devices() {
				log.Printf("микрофон: %s (id=%s)", d.Name, d.ID)
			}
			rec.Close()
			return
		}
		if arg == "-transcribe" && i+1 < len(os.Args[1:]) {
			path := os.Args[1:][i+1]
			cfg, cerr := loadConfig("config.json")
			if cerr != nil {
				log.Printf("конфигурация: %v", cerr)
				return
			}
			wav, rerr := os.ReadFile(path)
			if rerr != nil {
				log.Printf("файл %s: %v", path, rerr)
				return
			}
			srv, serr := startRecognizer(cfg, logFile)
			if serr != nil {
				log.Printf("распознаватель: %v", serr)
				return
			}
			defer srv.stop()
			if werr := srv.waitReady(engineReadyTimeout(srv)); werr != nil {
				log.Printf("распознаватель не поднялся: %v", werr)
				return
			}
			started := time.Now()
			text, terr := srv.transcribe(context.Background(), wav, cfg.Language, cfg.WhisperPrompt, false)
			if terr != nil {
				log.Printf("распознавание: %v", terr)
				return
			}
			log.Printf("движок=%s время=%.2f c текст=%q", srv.engine(), time.Since(started).Seconds(), text)
			return
		}
		if arg == "-routecheck" {
			cfg, cerr := loadConfig("config.json")
			if cerr != nil {
				log.Printf("конфигурация: %v", cerr)
				return
			}
			app := &App{cfg: cfg, enabled: true, evq: evqueue.New[ptEvent](8)}
			primary := primaryEngine(cfg)
			srv, serr := startEngine(cfg, primary, logFile)
			if serr != nil {
				log.Printf("основной движок: %v", serr)
				return
			}
			if werr := srv.waitReady(engineReadyTimeout(srv)); werr != nil {
				log.Printf("основной движок не поднялся: %v", werr)
				srv.stop()
				return
			}
			app.mu.Lock()
			app.srv = srv
			app.ready = true
			app.mu.Unlock()
			log.Printf("routecheck: основной движок %s поднят", primary)

			other := engineWhisper
			if primary == engineWhisper {
				other = engineSherpa
			}
			alt, aerr := app.engineFor(cfg, other)
			if aerr != nil {
				log.Printf("routecheck: второй движок %s не поднялся: %v", other, aerr)
			} else {
				log.Printf("routecheck: оба движка живы — %s и %s", srv.engine(), alt.engine())
				text, terr := alt.transcribe(context.Background(), mustReadWav(os.Args[1:], i), cfg.Language, cfg.WhisperPrompt, false)
				if terr != nil {
					log.Printf("routecheck: второй движок не распознал: %v", terr)
				} else {
					log.Printf("routecheck: второй движок ответил %q", text)
				}
			}
			log.Printf("routecheck: держу оба движка 20 секунд для замера памяти")
			time.Sleep(20 * time.Second)
			app.altMu.Lock()
			app.altUsed = time.Now().Add(-24 * time.Hour)
			app.altMu.Unlock()
			if app.sweepOnce() {
				log.Printf("routecheck: выгрузка по простою сработала")
			} else {
				log.Printf("routecheck: выгрузка по простою НЕ сработала")
			}
			app.stopAltEngine()
			srv.stop()
			log.Printf("routecheck: оба движка остановлены")
			return
		}
		if arg == "-dialogs" {
			cfg, _ := loadConfig("config.json")
			if cfg != nil {
				initLang(cfg.UILanguage)
				setOverlayPos(cfg.OverlayPos)
			}
			log.Printf("демонстрация диалогов: смена фокуса")
			log.Printf("ответ: %q", askFocusMismatch())
			if cfg != nil {
				log.Printf("демонстрация диалогов: выбор языка перевода")
				log.Printf("ответ: %q", askTranslateTarget(cfg))
			}
			return
		}
		if arg == "-replcheck" && i+1 < len(os.Args[1:]) {
			cfg, cerr := loadConfig("config.json")
			if cerr != nil {
				log.Printf("конфигурация: %v", cerr)
				return
			}
			in := os.Args[1:][i+1]
			log.Printf("replcheck: замен в конфиге — %d", len(cfg.Replacements))
			log.Printf("replcheck: %q → %q", in, replace.Apply(cfg.Replacements, in))
			return
		}
		if arg == "-rulecheck" {
			cfg, cerr := loadConfig("config.json")
			if cerr != nil {
				log.Printf("конфигурация: %v", cerr)
				return
			}
			log.Printf("rulecheck: три секунды на переключение в нужное окно")
			time.Sleep(3 * time.Second)
			fg, _, _ := procGetForegroundWindow.Call()
			exe := processNameOf(fg)
			log.Printf("rulecheck: окно=%q процесс=%q", windowTitle(fg), exe)
			c := *cfg
			applyAppRule(&c, exe)
			log.Printf("rulecheck: вставка=%s enter=%v задержка=%d промпты=%v",
				c.PasteMode, c.AutoEnter, c.PasteDelayMs, c.ActiveProfiles)
			return
		}
		if arg == "-dpi" {
			var pt point
			procGetCursorPosDPI.Call(uintptr(unsafe.Pointer(&pt)))
			log.Printf("dpi: GetDpiForSystem=%d dpiForCursor=%d курсор=%d,%d", dpiFor(0), dpiForCursor(), pt.X, pt.Y)
			mon, _, _ := procMonitorFromPoint.Call(uintptr(uint32(pt.X))|uintptr(uint32(pt.Y))<<32, 2)
			var dx, dy uint32
			r, _, _ := procGetDpiForMonitor.Call(mon, 0, uintptr(unsafe.Pointer(&dx)), uintptr(unsafe.Pointer(&dy)))
			log.Printf("dpi: монитор=%x GetDpiForMonitor rc=%d dx=%d", mon, r, dx)
			log.Printf("dpi: MonitorFromPoint.Find=%v GetDpiForMonitor.Find=%v", procMonitorFromPoint.Find(), procGetDpiForMonitor.Find())
			return
		}
		if arg == "-overlay" {
			cfg, _ := loadConfig("config.json")
			if cfg != nil {
				initLang(cfg.UILanguage)
				setOverlayPos(cfg.OverlayPos)
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
			log.Printf("демонстрация плашки: %q", text)
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
				log.Printf("микрофон недоступен: %v", rerr)
				return
			}
			if err := rec.Start(5); err != nil {
				log.Printf("запись для разбора не началась: %v", err)
			} else {
				time.Sleep(3 * time.Second)
				pcm := rec.Stop()
				rep := audiolevel.Analyze(pcm)
				log.Printf("разбор: пик %.0f дБ, RMS %.0f дБ, речь %.0f%%, обрезано %.2f%% → %s",
					audiolevel.DBFS(rep.Peak), audiolevel.DBFS(rep.RMS), rep.VoiceRatio*100, rep.ClipRatio*100, audiolevel.Verdict(rep))
			}
			log.Printf("замер уровня без диктовки, 5 секунд")
			for i := 0; i < 25; i++ {
				rec.MonitorPing()
				time.Sleep(200 * time.Millisecond)
				log.Printf("уровень: %.3f", rec.Level())
			}
			rec.Close()
			return
		}
	}

	if !acquireSingleInstance() {
		msgBox(tr("app.name"), tr("already.running"))
		return
	}

	cfg, err := loadConfig("config.json")
	if err != nil {
		msgBox(tr("cfg.err.title"), err.Error())
		return
	}
	initLang(cfg.UILanguage)
	setOverlayPos(cfg.OverlayPos)

	app := &App{
		cfg:     cfg,
		enabled: true,
		evq:     evqueue.New[ptEvent](256),
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
				log.Printf("testpaste: цель захвачена, вставка через 6 секунд")
				time.Sleep(6 * time.Second)
				app.insertResult(context.Background(), app.snapshot(), time.Now(), text, "", tw)
			}()
		}
	}
	if !cfg.WizardDone {
		log.Printf("первый запуск: открываю мастер настройки")
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
	log.Printf("=== %s %s запущен ===", appid.Slug, appVersion)
}

func (a *App) snapshot() *Config {
	a.mu.Lock()
	defer a.mu.Unlock()
	c := *a.cfg
	return &c
}

func (a *App) post(ev ptEvent) {
	if !a.evq.Push(ev) {
		log.Printf("очередь событий переполнена, событие %d пропущено (всего %d)", ev.kind, a.evq.Dropped())
	}
}

func (a *App) startCore() {
	a.setStatus(tr("status.loading"))
	ovOnCancel = func() { a.post(ptEvent{kind: evCancel}) }
	go a.worker()

	if _, err := parseHotkey(a.cfg.Hotkey); err != nil {
		a.fatal(trf("err.hotkey.cfg", err.Error()))
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
		a.fatal(trf("err.hook", err.Error()))
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
		a.fatal(trf("err.mic", err.Error()))
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
	ovMu.Lock()
	ovRecorder = rec
	ovMu.Unlock()

	attempts := 0
	waiting := false
	for {
		cfg := a.snapshot()
		if cfg.ServerAutostart && cfg.ServerURL == "" {
			if missing := missingModelPath(cfg); missing != "" {
				a.mu.Lock()
				q := a.quitting
				a.mu.Unlock()
				if q {
					return
				}
				if !waiting {
					waiting = true
					log.Printf("модель %s не найдена — жду скачивания", missing)
					a.setStatus(tr("status.nomodel"))
					traySetIcon(trayOff)
					go a.openSettings("rec")
				}
				time.Sleep(2 * time.Second)
				continue
			}
			if waiting {
				waiting = false
				a.setStatus(tr("status.loading"))
			}
		}
		srv, err := startRecognizer(cfg, logFile)
		if err != nil {
			a.fatal(err.Error())
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
			a.fatal(err.Error())
			return
		}

		a.mu.Lock()
		a.ready = true
		a.mu.Unlock()
		a.refreshIdleUI()
		log.Printf("готов: hotkey=%s движок=%s модель=%s lang=%s", cfg.Hotkey, srv.engine(), activeModelPath(cfg), cfg.Language)

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
			log.Printf("перезапуск распознавателя по запросу")
			a.setStatus(tr("status.loading"))
			continue
		}
		if srv.external() {
			return
		}
		if time.Since(started) > 5*time.Minute {
			attempts = 0
		}
		attempts++
		if attempts > 3 {
			a.fatal(tr("err.server.repeat"))
			return
		}
		log.Printf("распознаватель упал, перезапуск (попытка %d)", attempts)
		a.setStatus(tr("status.server.restart"))
		traySetIcon(trayOff)
	}
}

func (a *App) requestServerRestart() {
	a.mu.Lock()
	a.ready = false
	srv := a.srv
	a.mu.Unlock()
	if srv != nil {
		srv.stop()
	}
	a.stopAltEngine()
}

func (a *App) fatal(text string) {
	log.Printf("ОШИБКА: %s", text)
	a.setStatus(text)
	traySetIcon(trayOff)
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
		traySetIcon(trayOff)
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
	for _, p := range cfg.Profiles {
		if p.Hotkey == "" {
			continue
		}
		if groups, err := parseHotkey(p.Hotkey); err == nil {
			combos = append(combos, comboDef{id: p.ID, groups: groups})
		} else {
			log.Printf("хоткей профиля %s не разобран: %v", p.ID, err)
		}
	}
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
				a.handleStop(0)
			case evTimeout:
				a.handleStop(ev.gen)
			case evDone:
				a.handleDone(ev.gen)
			case evCancel:
				a.handleCancel()
			}
		}
	}
}

func (a *App) handleCancel() {
	switch a.state.Load() {
	case stRecording:
		a.mu.Lock()
		rec := a.rec
		a.mu.Unlock()
		rec.Stop()
		a.gen++
		a.state.Store(stIdle)
		log.Printf("запись отменена пользователем")
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
	if a.state.Load() == stRecording && cfg.HotkeyMode == hotkeyToggle {
		log.Printf("фиксация: второе нажатие останавливает запись")
		a.handleStop(0)
		return
	}
	if a.state.Load() != stIdle {
		return
	}
	if profileID == "main" {
		profileID = ""
	}
	a.mu.Lock()
	rec := a.rec
	ok := a.enabled && a.ready && !a.capturing && rec != nil
	a.mu.Unlock()
	if !ok {
		return
	}
	if err := rec.Start(cfg.MaxRecordSeconds); err != nil {
		log.Printf("ошибка старта записи: %v", err)
		if cfg.Overlay {
			overlaySet(ovFlashErr, tr("ov.err.mic"))
		}
		playCue(cfg.Beep, cfg.SoundTheme, cueError)
		return
	}
	a.gen++
	a.sessionTarget, _, _ = procGetForegroundWindow.Call()
	applyAppRule(cfg, processNameOf(a.sessionTarget))
	a.sessionCfg = cfg
	a.sessionProfile = profileID
	a.state.Store(stRecording)
	traySetIcon(trayRecording)
	a.setStatus(tr("status.recording"))
	playCue(cfg.Beep, cfg.SoundTheme, cueStart)
	if cfg.Overlay {
		ovAnim.Store(cfg.Animation)
		setOverlayPos(cfg.OverlayPos)
		overlaySet(ovRecording, tr("ov.speak"))
	}

	gen := a.gen
	time.AfterFunc(time.Duration(cfg.MaxRecordSeconds)*time.Second, func() {
		a.post(ptEvent{kind: evTimeout, gen: gen})
	})
}

func (a *App) handleStop(expectGen int) {
	if a.state.Load() != stRecording {
		return
	}
	if expectGen != 0 && expectGen != a.gen {
		return
	}
	if expectGen != 0 {
		log.Printf("достигнут max_record_seconds=%d, останавливаю", a.sessionCfg.MaxRecordSeconds)
	}
	a.mu.Lock()
	rec := a.rec
	a.mu.Unlock()
	pcm := rec.Stop()
	a.state.Store(stProcessing)
	playCue(a.sessionCfg.Beep, a.sessionCfg.SoundTheme, cueStop)
	traySetIcon(trayProcessing)
	a.setStatus(tr("status.transcribing"))
	if a.sessionCfg.Overlay {
		procText := ""
		if a.sessionProfile == "wtranslate" {
			procText = tr("ov.translating")
		}
		overlaySet(ovProcessing, procText)
	}
	ctx, cancel := context.WithCancel(context.Background())
	a.sessionCancel = cancel
	go a.process(ctx, pcm, a.gen, a.sessionCfg, a.sessionProfile, a.sessionTarget)
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

	minBytes := sampleRate * 2 * cfg.MinRecordMs / 1000
	if len(pcm) < minBytes {
		log.Printf("запись слишком короткая (%d мс), пропускаю", len(pcm)*1000/(sampleRate*2))
		if cfg.Overlay {
			overlayHide()
		}
		return
	}
	sound := audiolevel.Analyze(pcm)
	verdict := audiolevel.Verdict(sound)
	log.Printf("звук: пик %.0f дБ, речь %.0f%%, обрезано %.1f%% — %s",
		audiolevel.DBFS(sound.Peak), sound.VoiceRatio*100, sound.ClipRatio*100, verdict)
	if verdict == audiolevel.VerdictSilent {
		log.Printf("тишина в записи — распознавание пропущено")
		if cfg.Overlay {
			overlaySet(ovFlashErr, tr("ov.silence"))
		}
		return
	}

	a.mu.Lock()
	srv := a.srv
	a.mu.Unlock()
	if srv == nil {
		log.Printf("сервер распознавания ещё не запущен")
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
			log.Printf("распознавание отменено пользователем")
			if cfg.Overlay {
				overlaySet(ovFlashErr, tr("ov.cancelled"))
			}
			return
		}
		switch {
		case target == "":
			log.Printf("перевод: отменён в диалоге, вставляю как есть")
		case target == cfg.Language:
			log.Printf("перевод: цель %s совпадает с языком распознавания — пропускаю", target)
			target = ""
		default:
			log.Printf("перевод: цель=%s силами Whisper (режим %s)", target, cfg.TranslateAsk)
		}
	}
	decision := pickEngine(cfg, target != "")
	if decision.Engine != srv.engine() {
		alt, aerr := a.engineFor(cfg, decision.Engine)
		if aerr != nil {
			log.Printf("движок %s не поднялся (%v) — остаюсь на %s", decision.Engine, aerr, srv.engine())
			if cfg.Overlay {
				overlaySet(ovFlashErr, tr("ov.engine.fallback"))
			}
		} else {
			srv = alt
		}
	}
	log.Printf("маршрутизация: движок=%s причина=%s язык=%s перевод=%v", srv.engine(), decision.Reason, cfg.Language, target != "")
	if target != "" && !engineTranslates(srv.engine()) {
		log.Printf("перевод недоступен: активен движок %s — вставляю распознанный текст как есть", srv.engine())
		if cfg.Overlay {
			overlaySet(ovFlashErr, tr("ov.notranslate"))
		}
		target = ""
	}
	fastTranslate := target == "en"
	if fastTranslate && strings.Contains(cfg.Model, "turbo") {
		log.Printf("предупреждение: модель turbo не обучена переводу — результат может оказаться транскрипцией")
	}
	recLang := cfg.Language
	if target != "" && !fastTranslate {
		recLang = target
	}
	text, err := srv.transcribe(ctx, wavFromPCM16(pcm, sampleRate), recLang, cfg.WhisperPrompt, fastTranslate)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Printf("распознавание отменено пользователем")
			if cfg.Overlay {
				overlaySet(ovFlashErr, tr("ov.cancelled"))
			}
			return
		}
		log.Printf("распознавание: %v", err)
		playCue(cfg.Beep, cfg.SoundTheme, cueError)
		if cfg.Overlay {
			overlaySet(ovFlashErr, tr("ov.err.recognize"))
		}
		return
	}
	if ctx.Err() != nil {
		log.Printf("распознавание отменено пользователем")
		if cfg.Overlay {
			overlaySet(ovFlashErr, tr("ov.cancelled"))
		}
		return
	}
	text = strings.TrimSpace(text)
	if fixed := replace.Apply(cfg.Replacements, text); fixed != text {
		log.Printf("замены: %q → %q", text, fixed)
		text = fixed
	}
	if text == "" || onlyNoise.MatchString(text) {
		log.Printf("пустой результат (%q), вставлять нечего", text)
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
	if chain := punctChain(cfg, chainProfiles(cfg, profileID)); len(chain) > 0 && llmInstalled(cfg) {
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
				continue
			}
			if errors.Is(lerr, context.Canceled) || ctx.Err() != nil {
				log.Printf("редактура отменена пользователем")
				if cfg.Overlay {
					overlaySet(ovFlashErr, tr("ov.cancelled"))
				}
				return
			}
			log.Printf("профиль %s не применился, продолжаю с текущим текстом: %v", prof.ID, lerr)
			if skipped == "" {
				skipped = prof.Name
			}
		}
	}

	if cfg.Punctuation == punctOff {
		text = stripPunctuation(text)
	}
	a.insertResult(ctx, cfg, start, text, skipped, targetWnd)
}

func (a *App) insertResult(ctx context.Context, cfg *Config, start time.Time, text, skipped string, targetWnd uintptr) {
	a.mu.Lock()
	a.lastResult = text
	a.lastResultAt = time.Now()
	a.lastTarget = windowTitle(targetWnd)
	a.lastProcess = processNameOf(targetWnd)
	targetApp := a.lastProcess
	a.mu.Unlock()
	a.rememberDictation(cfg, text, targetApp)

	allowEnter := cfg.AutoEnter && skipped == ""
	if targetWnd != 0 {
		cur, _, _ := procGetForegroundWindow.Call()
		if cur != targetWnd {
			log.Printf("фокус сменился во время обработки — спрашиваю")
			switch askFocusMismatch() {
			case "here":
				allowEnter = false
			case "copy":
				if err := setClipboardText(text); err != nil {
					log.Printf("копирование: %v", err)
					return
				}
				log.Printf("результат скопирован в буфер по выбору пользователя")
				if cfg.Overlay {
					overlaySet(ovFlashOK, tr("ov.copied"))
				}
				return
			default:
				log.Printf("вставка отменена, текст сохранён в последнем результате")
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

	if err := pasteText(cfg, text); err != nil {
		log.Printf("вставка: %v", err)
		playCue(cfg.Beep, cfg.SoundTheme, cueError)
		if cfg.Overlay {
			overlaySet(ovFlashErr, tr("ov.err.paste"))
		}
		return
	}
	if allowEnter {
		time.Sleep(150 * time.Millisecond)
		if err := pressEnter(); err != nil {
			log.Printf("auto-enter: %v", err)
		}
	}
	log.Printf("готово за %.1fс: %d символов", time.Since(start).Seconds(), len([]rune(text)))
	if cfg.Overlay {
		switch {
		case skipped != "":
			overlaySet(ovFlashErr, trf("ov.llm.skipped", skipped))
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

func (a *App) rememberDictation(cfg *Config, text, app string) {
	if !cfg.HistoryOn || strings.TrimSpace(text) == "" {
		return
	}
	if skip := strings.TrimSpace(cfg.HistorySkip); skip != "" && app != "" {
		if _, hit := apprules.Find([]apprules.Rule{{Match: skip}}, app); hit {
			log.Printf("история: %s в списке исключений, не записываю", app)
			return
		}
	}
	now := time.Now()
	item := history.Item{At: now.UnixMilli(), Text: text, App: app}
	if err := histStore.Add(item, now.UnixMilli(), cfg.HistoryDays, cfg.HistoryMax); err != nil {
		log.Printf("история: %v", err)
	}
}

func applyAppRule(cfg *Config, exe string) {
	if exe == "" || len(cfg.AppRules) == 0 {
		return
	}
	rule, ok := apprules.Find(cfg.AppRules, exe)
	if !ok {
		return
	}
	if rule.Paste == apprules.PasteClipboard || rule.Paste == apprules.PasteType {
		cfg.PasteMode = rule.Paste
	}
	switch rule.Enter {
	case apprules.EnterOn:
		cfg.AutoEnter = true
	case apprules.EnterOff:
		cfg.AutoEnter = false
	}
	if rule.DelayMs > 0 {
		cfg.PasteDelayMs = rule.DelayMs
	}
	if rule.UseProfiles {
		cfg.ActiveProfiles = append([]string(nil), rule.Profiles...)
	}
	log.Printf("правило для %s: вставка=%s enter=%v задержка=%d мс промпты=%v",
		exe, cfg.PasteMode, cfg.AutoEnter, cfg.PasteDelayMs, cfg.ActiveProfiles)
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
		log.Printf("смена сочетания отменена")
		return
	}
	if _, err := parseHotkey(combo); err != nil {
		log.Printf("захваченное сочетание %q не разобрано: %v", combo, err)
		return
	}
	a.mu.Lock()
	c := *a.cfg
	c.Hotkey = combo
	a.cfg = &c
	a.mu.Unlock()
	hook.SetCombos(buildCombos(&c))
	if err := saveConfig("config.json", &c); err != nil {
		log.Printf("сохранение конфига: %v", err)
	}
	a.refreshIdleUI()
	log.Printf("новое сочетание: %s", combo)
}

func chainProfiles(cfg *Config, profileID string) []*Profile {
	if profileID == "wtranslate" {
		return nil
	}
	var ids []string
	if profileID != "" {
		ids = []string{profileID}
	} else {
		active := map[string]bool{}
		for _, id := range cfg.ActiveProfiles {
			active[id] = true
		}
		for i := range cfg.Profiles {
			if active[cfg.Profiles[i].ID] {
				ids = append(ids, cfg.Profiles[i].ID)
			}
		}
	}
	var out []*Profile
	for _, id := range ids {
		if p := profileByID(cfg, id); p != nil && p.Prompt != "" {
			out = append(out, p)
		}
	}
	return out
}

func (a *App) llmProcess(ctx context.Context, prompt, text string) (string, error) {
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
		return "", errors.New("пустой ответ LLM")
	}
	return out, nil
}

func (a *App) reloadConfig() {
	fresh, err := loadConfig("config.json")
	if err != nil {
		log.Printf("перечитать конфиг: %v", err)
		a.setStatus(tr("status.cfg.err"))
		return
	}
	if _, err := parseHotkey(fresh.Hotkey); err != nil {
		log.Printf("перечитать конфиг: %v", err)
		a.setStatus(tr("status.cfg.err"))
		return
	}

	a.mu.Lock()
	old := a.cfg
	a.cfg = fresh
	hook := a.hook
	a.mu.Unlock()
	initLang(fresh.UILanguage)

	serverChanged := fresh.Model != old.Model ||
		fresh.ServerPort != old.ServerPort ||
		fresh.Threads != old.Threads ||
		fresh.ServerURL != old.ServerURL ||
		fresh.ServerExe != old.ServerExe ||
		fresh.ServerAutostart != old.ServerAutostart
	if serverChanged {
		a.requestServerRestart()
	}

	if hook != nil {
		hook.SetCombos(buildCombos(fresh))
	} else {
		log.Printf("хук не установлен, новое сочетание применится после перезапуска")
	}

	if serverChanged {
		a.setStatus(tr("status.loading"))
		log.Printf("конфиг перечитан; распознаватель перезапускается с новыми настройками")
	} else {
		a.refreshIdleUI()
		log.Printf("конфиг перечитан: hotkey=%s paste=%s", fresh.Hotkey, fresh.PasteMode)
	}
}

func (a *App) onExit() {
	a.mu.Lock()
	a.quitting = true
	srv := a.srv
	rec := a.rec
	llm := a.llm
	a.mu.Unlock()
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
	log.Printf("выход")
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

func punctChain(cfg *Config, chain []*Profile) []*Profile {
	if cfg.Punctuation != punctByLLM {
		return chain
	}
	head := &Profile{ID: "punct", Name: strS("S_PUNCT"), Prompt: tr("punct.prompt")}
	return append([]*Profile{head}, chain...)
}

var punctMarks = []string{".", ",", "!", "?", ";", ":", "…", "—", "–"}

func stripPunctuation(text string) string {
	for _, m := range punctMarks {
		text = strings.ReplaceAll(text, m, "")
	}
	text = strings.Join(strings.Fields(text), " ")
	return strings.ToLower(strings.TrimSpace(text))
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
		log.Printf("проверка микрофона: %v", err)
		return fail(err.Error())
	}
	time.Sleep(3 * time.Second)
	pcm := rec.Stop()
	rep := audiolevel.Analyze(pcm)
	verdict := audiolevel.Verdict(rep)
	peak := audiolevel.DBFS(rep.Peak)
	log.Printf("проверка микрофона: пик %.0f дБ, речь %.0f%%, обрезано %.1f%% — %s",
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
