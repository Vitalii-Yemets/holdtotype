package main

import (
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
)

var appVersion = "0.7.3"

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
	srv       *whisperServer
	llm       *llamaServer
	hook      *hotkeyHook
	enabled   bool
	ready     bool
	capturing bool
	quitting  bool

	events chan ptEvent

	state atomic.Int32
	gen            int
	sessionCfg     *Config
	sessionProfile string
	sessionCancel  context.CancelFunc
	sessionTarget  uintptr

	lastResult string
	updVer     string
	updURL     string
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

	app := &App{
		cfg:     cfg,
		enabled: true,
		events:  make(chan ptEvent, 32),
	}
	app.startCore()
	go cleanupWebViewProfiles()
	go app.startupUpdateCheck()
	go cleanupStaleParts()
	args := os.Args[1:]
	for i, arg := range args {
		if arg == "-settings" {
			go func() {
				time.Sleep(2 * time.Second)
				app.openSettings("proc")
			}()
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
	runTray(app)
}

func setupLog() {
	logFile = newRotatingWriter("voxterminal.log", 1<<20)
	log.SetOutput(logFile)
	log.Printf("=== voxterminal %s запущен ===", appVersion)
}

func (a *App) snapshot() *Config {
	a.mu.Lock()
	defer a.mu.Unlock()
	c := *a.cfg
	return &c
}

func (a *App) post(ev ptEvent) {
	select {
	case a.events <- ev:
	default:
		log.Printf("очередь событий переполнена, событие %d пропущено", ev.kind)
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
	rec, err := NewRecorder()
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
			if _, err := os.Stat(cfg.Model); err != nil {
				a.mu.Lock()
				q := a.quitting
				a.mu.Unlock()
				if q {
					return
				}
				if !waiting {
					waiting = true
					log.Printf("модель %s не найдена — жду скачивания", cfg.Model)
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
		srv, err := startWhisperServer(cfg, logFile)
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

		timeout := 3 * time.Minute
		if srv.external() {
			timeout = 20 * time.Second
		}
		if err := srv.waitReady(timeout); err != nil {
			a.fatal(err.Error())
			return
		}

		a.mu.Lock()
		a.ready = true
		a.mu.Unlock()
		a.refreshIdleUI()
		log.Printf("готов: hotkey=%s model=%s lang=%s", cfg.Hotkey, cfg.Model, cfg.Language)

		if srv.external() {
			return
		}
		started := time.Now()
		<-srv.done

		a.mu.Lock()
		q := a.quitting
		a.ready = false
		a.mu.Unlock()
		if q {
			return
		}
		if srv.wasStopped() {
			attempts = 0
			log.Printf("перезапуск сервера по запросу")
			a.setStatus(tr("status.loading"))
			continue
		}
		if time.Since(started) > 5*time.Minute {
			attempts = 0
		}
		attempts++
		if attempts > 3 {
			a.fatal(tr("err.server.repeat"))
			return
		}
		log.Printf("whisper-server упал, перезапуск (попытка %d)", attempts)
		a.setStatus(tr("status.server.restart"))
		traySetIcon(trayOff)
	}
}

func (a *App) requestServerRestart() {
	a.mu.Lock()
	a.ready = false
	srv := a.srv
	a.mu.Unlock()
	if srv != nil && !srv.external() {
		srv.stop()
	}
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
	for ev := range a.events {
		switch ev.kind {
		case evDown:
			a.handleDown(ev.profile)
		case evUp:
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
		tdAbort()
		fdAbort()
		if a.sessionCancel != nil {
			a.sessionCancel()
		}
	}
}

func (a *App) handleDown(profileID string) {
	if a.state.Load() != stIdle {
		return
	}
	cfg := a.snapshot()
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
	if err := rec.Start(); err != nil {
		log.Printf("ошибка старта записи: %v", err)
		return
	}
	a.gen++
	a.sessionCfg = cfg
	a.sessionProfile = profileID
	a.sessionTarget, _, _ = procGetForegroundWindow.Call()
	a.state.Store(stRecording)
	traySetIcon(trayRecording)
	a.setStatus(tr("status.recording"))
	playCue(cfg.Beep, cfg.SoundTheme, cueStart)
	if cfg.Overlay {
		ovAnim.Store(cfg.Animation)
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
	if text == "" || onlyNoise.MatchString(text) {
		log.Printf("пустой результат (%q), вставлять нечего", text)
		if cfg.Overlay {
			overlaySet(ovFlashErr, tr("ov.silence"))
		}
		return
	}

	skipped := ""
	if chain := chainProfiles(cfg, profileID); len(chain) > 0 && llmInstalled(cfg) {
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

	a.insertResult(ctx, cfg, start, text, skipped, targetWnd)
}

func (a *App) insertResult(ctx context.Context, cfg *Config, start time.Time, text, skipped string, targetWnd uintptr) {
	a.mu.Lock()
	a.lastResult = text
	a.mu.Unlock()

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
		if skipped != "" {
			overlaySet(ovFlashErr, trf("ov.llm.skipped", skipped))
		} else {
			overlaySet(ovFlashOK, trf("ov.inserted", len([]rune(text))))
		}
	}
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

	if fresh.Model != old.Model {
		a.requestServerRestart()
	}
	restartNeeded := fresh.ServerPort != old.ServerPort ||
		fresh.Threads != old.Threads ||
		fresh.ServerURL != old.ServerURL ||
		fresh.ServerExe != old.ServerExe ||
		fresh.ServerAutostart != old.ServerAutostart

	if hook != nil {
		hook.SetCombos(buildCombos(fresh))
	} else {
		log.Printf("хук не установлен, новое сочетание применится после перезапуска")
	}

	if restartNeeded {
		a.setStatus(tr("status.restart.needed"))
		log.Printf("конфиг перечитан; модель/язык/сервер вступят в силу после перезапуска")
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
