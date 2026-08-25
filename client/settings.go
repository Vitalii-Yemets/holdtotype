package main

import (
	"context"
	"encoding/json"
	"fmt"
	"holdtotype/internal/apprules"
	"holdtotype/internal/commands"
	"holdtotype/internal/replace"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"

	webview "github.com/jchv/go-webview2"
	"golang.org/x/sys/windows"

	"holdtotype/internal/hotkeys"

	"holdtotype/internal/appid"
	"holdtotype/internal/plexfont"
	"holdtotype/internal/theme"
)

var (
	settingsOpen atomic.Bool
	settingsHwnd atomic.Uintptr
)

func tryCreateWebView(width, height int) (w webview.WebView) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("createWebView: panic: %v", r)
			w = nil
		}
		if w != nil {
			hwnd := uintptr(w.Window())
			procSetWindowPos.Call(hwnd, 0, offscreenPos(), offscreenPos(), 0, 0, 0x0001|0x0004|0x0010)
			setDarkClientBackground(hwnd)
			applyDarkCaption(hwnd)
		}
	}()
	return webview.NewWithOptions(webview.WebViewOptions{
		DataPath:  filepath.Join(os.TempDir(), appid.TempDirName("webview", os.Getpid())),
		AutoFocus: true,
		WindowOptions: webview.WindowOptions{
			Title:  strS("S_TITLE"),
			Width:  uint(width),
			Height: uint(height),
			IconId: 1,
			Center: true,
		},
	})
}

func createWebView(width, height int) webview.WebView {
	for attempt := 1; attempt <= 3; attempt++ {
		if w := tryCreateWebView(width, height); w != nil {
			return w
		}
		log.Printf("createWebView: попытка %d не удалась", attempt)
		time.Sleep(800 * time.Millisecond)
	}
	return nil
}

func cleanupWebViewProfiles() {
	own := appid.TempDirName("webview", os.Getpid())
	stale := []string{
		appid.TempDirPrefix("webview"),
		appid.PrevTempPrefix("webview"),
		appid.PrevTempPrefix("setup"),
	}
	entries, err := os.ReadDir(os.TempDir())
	if err != nil {
		return
	}
	for _, e := range entries {
		name := e.Name()
		if !e.IsDir() || name == own {
			continue
		}
		for _, prefix := range stale {
			if strings.HasPrefix(name, prefix) {
				_ = os.RemoveAll(filepath.Join(os.TempDir(), name))
				break
			}
		}
	}
}

type settingsForm struct {
	Hotkey           string `json:"hotkey"`
	UILanguage       string `json:"ui_language"`
	Beep             bool   `json:"beep"`
	SoundTheme       string `json:"sound_theme"`
	AutoEnter        bool   `json:"auto_enter"`
	RestoreClipboard bool   `json:"restore_clipboard"`
	Overlay          bool   `json:"overlay"`
	OverlayPos       string `json:"overlay_position"`
	OverlayText      bool   `json:"overlay_text"`
	Animation        bool   `json:"animation"`
	TypeMode         bool   `json:"type_mode"`
	Language         string `json:"language"`
	ModelID          string `json:"model_id"`
	Threads          int    `json:"threads"`
	MinRecordMs      int    `json:"min_record_ms"`
	PasteDelayMs     int    `json:"paste_delay_ms"`
	HistoryOn        bool   `json:"history"`
	HistoryDays      int    `json:"history_days"`
	HistoryMax       int    `json:"history_max"`
	HistorySkip      string `json:"history_skip"`
	MaxRecordSeconds int    `json:"max_record_seconds"`
	ServerAutostart  bool   `json:"server_autostart"`
	CheckUpdates     bool   `json:"check_updates"`
	MicDevice        string `json:"mic_device"`
	MicDeviceName    string `json:"mic_device_name"`
	Punctuation      string `json:"punctuation"`
	HotkeyMode       string `json:"hotkey_mode"`
	UILevel          string `json:"ui_level"`
	Skin             string `json:"skin"`
	Theme            string `json:"theme"`
	ServerPort       int    `json:"server_port"`
	ServerExe        string `json:"server_exe"`
	ServerURL        string `json:"server_url"`

	WhisperPrompt       string             `json:"whisper_prompt"`
	TranslateHotkey     string             `json:"translate_hotkey"`
	PauseHotkey         string             `json:"pause_hotkey"`
	TranslateTarget     string             `json:"translate_target"`
	TranslateAsk        string             `json:"translate_ask"`
	TranslateAskSeconds int                `json:"translate_ask_seconds"`
	TranslateAskLangs   []string           `json:"translate_ask_langs"`
	TranslateDefault    bool               `json:"translate_default"`
	ActiveProfiles      []string           `json:"active_profiles"`
	AppRules            []apprules.Rule    `json:"app_rules"`
	Replacements        []replace.Rule     `json:"replacements"`
	Commands            []commands.Command `json:"commands"`
	LLMModelFile        string             `json:"llm_model_file"`
	Profiles            []Profile          `json:"profiles"`
}

func (a *App) openSettings(tab string) {
	log.Printf("openSettings: tab=%s", tab)
	if fg, _, _ := procGetForegroundWindow.Call(); fg != 0 && fg != settingsHwnd.Load() && !ownWindow(fg) {
		a.mu.Lock()
		a.settingsPrev = fg
		a.mu.Unlock()
		log.Printf("openSettings: вернуть текст можно в окно «%s»", windowTitle(fg))
	}
	if !settingsOpen.CompareAndSwap(false, true) {
		if hwnd := settingsHwnd.Load(); hwnd != 0 {
			log.Printf("openSettings: уже открыто, поднимаю на передний план")
			procShowWindow.Call(hwnd, 9)
			procSetForegroundWnd.Call(hwnd)
		}
		return
	}
	go a.settingsThread(tab, 1)
}

func (a *App) settingsThread(tab string, attempt int) {
	runtime.LockOSThread()
	if err := windows.CoInitializeEx(0, windows.COINIT_APARTMENTTHREADED); err != nil {
		log.Printf("openSettings: CoInitializeEx (попытка %d): %v", attempt, err)
		if attempt < 5 {
			go a.settingsThread(tab, attempt+1)
			return
		}
		settingsOpen.Store(false)
		msgBox(tr("err.title"), tr("err.webview"))
		return
	}
	func() {
		defer settingsOpen.Store(false)
		defer func() {
			if r := recover(); r != nil {
				log.Printf("openSettings: panic: %v\n%s", r, debug.Stack())
			}
		}()
		log.Printf("openSettings: создаю WebView2")
		winW, winH := settingsDefaultW, settingsDefaultH
		if c := a.snapshot(); c.SettingsW >= settingsMinW && c.SettingsH >= settingsMinH && c.SettingsW <= settingsMaxW && c.SettingsH <= settingsMaxH {
			winW, winH = c.SettingsW, c.SettingsH
		}
		lastWndW, lastWndH = 0, 0
		stopHider := hideWebViewWindowEarly(strS("S_TITLE"))
		w := createWebView(winW, winH)
		stopHider()
		if w == nil {
			log.Printf("WebView2 недоступен")
			if tab == "about" {
				showAbout(a.snapshot())
				return
			}
			msgBox(tr("err.title"), tr("err.webview"))
			shellOpenURL("https://developer.microsoft.com/en-us/microsoft-edge/webview2/")
			return
		}
		defer w.Destroy()

		hwnd := uintptr(w.Window())
		_ = w.Bind("appMin", func() {
			procShowWindow.Call(hwnd, 6)
		})
		_ = w.Bind("appMaxRestore", func() bool {
			return toggleMaximize(hwnd)
		})
		_ = w.Bind("appMaximized", func() bool {
			return windowMaximized(hwnd)
		})
		_ = w.Bind("appResizeEdge", func(edge int) {
			if e := uintptr(edge); validResizeEdge(e) {
				beginWindowResize(hwnd, e)
			}
		})
		_ = w.Bind("appClose", func() {
			procPostMessageW.Call(hwnd, wmClose, 0, 0)
		})
		_ = w.Bind("appDrag", func() {
			beginWindowDrag(hwnd)
		})
		_ = w.Bind("appSave", func(formJSON string) string {
			var f settingsForm
			if err := json.Unmarshal([]byte(formJSON), &f); err != nil {
				return jsonResult(saveResult{Severity: "error", Message: humanError(err)})
			}
			return jsonResult(a.applySettings(&f))
		})
		_ = w.Bind("appCapture", func() {
			go func() {
				a.changeHotkey()
				combo := a.snapshot().Hotkey
				warn := hotkeyWarning(combo)
				w.Dispatch(func() {
					w.Eval(fmt.Sprintf("setHotkey(%q, %q)", combo, warn))
				})
			}()
		})
		_ = w.Bind("appModels", func() string {
			return a.modelRows()
		})
		_ = w.Bind("appCopyLast", func() string {
			ok, msg := a.copyLastResult()
			return listsAnswer(listsReply{OK: ok, Text: msg})
		})
		_ = w.Bind("appState", func() string {
			return a.stateSnapshot()
		})
		_ = w.Bind("appRouting", func() string {
			out, _ := json.Marshal(routeRows(a.snapshot()))
			return string(out)
		})
		_ = w.Bind("appAdvise", func(lang, priority string, needTranslate bool) string {
			return adviseModel(lang, priority, needTranslate)
		})
		_ = w.Bind("appModelCancel", func(id string) bool {
			return cancelDownload(id)
		})
		_ = w.Bind("appModelDl", func(id string) {
			a.downloadModel(id)
		})
		_ = w.Bind("appPreviewSound", func(theme string) {
			previewTheme(theme)
		})
		_ = w.Bind("appLLM", func() string {
			return a.llmStatus()
		})
		_ = w.Bind("appLLMDlFile", func(repo, file string, sizeMB float64) {
			a.llmDownloadFile(repo, file, int(sizeMB))
		})
		_ = w.Bind("appLLMDel", func(file string) string {
			return a.llmDelete(file)
		})
		_ = w.Bind("appLLMSearch", func(q string) string {
			return a.llmSearch(q)
		})
		_ = w.Bind("appLLMFiles", func(repo string) string {
			return a.llmRepoFiles(repo)
		})
		_ = w.Bind("appHFPage", func(repo string) {
			if repoOK(repo) {
				shellOpenURL("https://huggingface.co/" + repo)
			}
		})
		_ = w.Bind("appHFHome", func() {
			shellOpenURL("https://huggingface.co/models?library=gguf&sort=downloads")
		})
		_ = w.Bind("appUpdateStatus", func() string {
			a.mu.Lock()
			v, u := a.updVer, a.updURL
			a.mu.Unlock()
			out, _ := json.Marshal(map[string]any{"current": appVersion, "latest": v, "url": u})
			return string(out)
		})
		_ = w.Bind("appCheckUpdate", func() string {
			tag, uurl, dig, err := fetchLatestRelease()
			if err != nil {
				out, _ := json.Marshal(map[string]any{"error": humanError(err)})
				return string(out)
			}
			newer := verNewer(tag, appVersion) && uurl != ""
			if newer {
				a.mu.Lock()
				a.updVer, a.updURL, a.updDigest = tag, uurl, dig
				a.mu.Unlock()
			}
			out, _ := json.Marshal(map[string]any{"current": appVersion, "latest": tag, "newer": newer})
			return string(out)
		})
		_ = w.Bind("appDoUpdate", func() {
			a.mu.Lock()
			uurl, udig := a.updURL, a.updDigest
			a.mu.Unlock()
			if uurl == "" {
				return
			}
			go func() {
				path, err := downloadSetup(uurl, udig, func(pct int) {
					w.Dispatch(func() { w.Eval(fmt.Sprintf("updProgress(%d)", pct)) })
				})
				if err == nil {
					err = launchUpdater(path)
				}
				if err != nil {
					enc, _ := json.Marshal(humanError(err))
					w.Dispatch(func() { w.Eval("updError(" + string(enc) + ")") })
					return
				}
				a.quitForUpdate()
			}()
		})
		_ = w.Bind("appRepoLink", func() {
			shellOpenURL(appid.RepoURL)
		})
		_ = w.Bind("appAuthorLink", func() {
			shellOpenURL(appid.AuthorURL)
		})
		_ = w.Bind("appCaptureCombo", func() {
			go func() {
				a.mu.Lock()
				hook := a.hook
				busy := a.capturing || hook == nil
				if !busy {
					a.capturing = true
				}
				a.mu.Unlock()
				if busy {
					return
				}
				defer func() {
					a.mu.Lock()
					a.capturing = false
					a.mu.Unlock()
				}()
				combo, ok := captureHotkeyDialog(hook, "—")
				if !ok {
					combo = ""
				}
				warn := hotkeyWarning(combo)
				if warn != "" {
					log.Printf("предупреждение: %s", warn)
				}
				w.Dispatch(func() {
					w.Eval(fmt.Sprintf("comboCaptured(%q, %q)", combo, warn))
				})
			}()
		})
		_ = w.Bind("appLLMTest", func(prompt, sample string) {
			go func() {
				out, err := a.llmProcess(context.Background(), prompt, sample)
				if err != nil {
					out = "⚠ " + humanError(err)
				}
				enc, _ := json.Marshal(out)
				w.Dispatch(func() {
					w.Eval("llmTestResult(" + string(enc) + ")")
				})
			}()
		})
		_ = w.Bind("appMics", func() string {
			a.mu.Lock()
			rec := a.rec
			a.mu.Unlock()
			if rec == nil {
				return "[]"
			}
			out, _ := json.Marshal(rec.devices())
			return string(out)
		})
		_ = w.Bind("appMicLevel", func() float64 {
			a.mu.Lock()
			rec := a.rec
			a.mu.Unlock()
			if rec == nil {
				return 0
			}
			rec.MonitorPing()
			return rec.Level()
		})
		_ = w.Bind("appMicSelect", func(id string) string {
			a.mu.Lock()
			rec := a.rec
			a.mu.Unlock()
			if rec == nil {
				return jsonResult(saveResult{Severity: "error", Message: tr("ov.err.mic")})
			}
			if err := rec.SetDevice(id); err != nil {
				log.Printf("выбор микрофона: %v", err)
				return jsonResult(saveResult{Severity: "error", Message: humanError(err)})
			}
			return jsonResult(saveResult{OK: true, Severity: "ok"})
		})
		_ = w.Bind("appJSError", func(msg string) {
			log.Printf("ошибка страницы настроек: %s", msg)
		})
		_ = w.Bind("appReload", func(tabName string) {
			w.Dispatch(func() {
				w.SetHtml(settingsHTML(a.snapshot(), tabName))
			})
		})
		_ = w.Bind("appModelDel", func(id string, force bool) string {
			return a.deleteModel(id, force)
		})
		_ = w.Bind("appRetryBackend", func() {
			log.Printf("настройки: запрошен повтор запуска распознавателя")
			a.requestServerRestart()
		})
		_ = w.Bind("appOpenLog", func() {
			log.Printf("настройки: открываю лог")
			shellOpen(appid.LogFile)
		})
		_ = w.Bind("appReloadConfig", func() {
			log.Printf("настройки: перечитываю config.json по кнопке")
			a.reloadConfig()
		})
		_ = w.Bind("appResetSettings", func() {
			a.resetSettings()
		})
		_ = w.Bind("appCheckModels", func() string {
			return a.verifyModels()
		})
		_ = w.Bind("appMicCheck", func() string {
			return a.micCheck()
		})
		_ = w.Bind("appHistory", func(query string) string {
			cfg := a.snapshot()
			a.enforceHistory(cfg)
			out, _ := json.Marshal(histStore.Search(query, cfg.HistoryMax))
			return string(out)
		})
		_ = w.Bind("appHistoryClear", func() {
			if err := histStore.Clear(); err != nil {
				log.Printf("история: очистка: %v", err)
			}
		})
		_ = w.Bind("appHistoryCopy", func(at float64) string {
			for _, it := range histStore.Items() {
				if it.At == int64(at) {
					if err := setClipboardText(it.Text); err != nil {
						log.Printf("копирование из истории: %v", err)
						return listsAnswer(listsReply{Text: trf("copy.fail", humanError(err))})
					}
					log.Printf("из истории скопировано: %d символов", len([]rune(it.Text)))
					return listsAnswer(listsReply{OK: true, Text: tr("copy.ok")})
				}
			}
			return listsAnswer(listsReply{Text: tr("hist.insert.gone")})
		})
		_ = w.Bind("appHistoryInsert", func(at float64) string {
			for _, it := range histStore.Items() {
				if it.At == int64(at) {
					return a.insertFromHistory(it.Text)
				}
			}
			out, _ := json.Marshal(map[string]any{"ok": false, "text": tr("hist.insert.gone")})
			return string(out)
		})
		_ = w.Bind("appListsExport", func(payload string) string {
			return exportLists(payload)
		})
		_ = w.Bind("appListsImport", func(payload string) string {
			return importLists(payload)
		})
		_ = w.Bind("appTestText", func(text string) string {
			cfg := a.snapshot()
			out := replace.Apply(cfg.Replacements, text)
			res := commands.Apply(cfg.Commands, out)
			shown := res.Text
			if res.Cancelled {
				shown = tr("ov.cmd.cancelled")
			}
			b, _ := json.Marshal(struct {
				Text      string `json:"text"`
				Cancelled bool   `json:"cancelled"`
			}{shown, res.Cancelled})
			return string(b)
		})
		_ = w.Bind("appAutorun", func() bool {
			return autorunEnabled()
		})
		_ = w.Bind("appSetAutorun", func(on bool) bool {
			if err := setAutorun(on); err != nil {
				log.Printf("автозапуск с Windows: %v", err)
			}
			return autorunEnabled()
		})
		_ = w.Bind("appWizardDone", func() {
			a.mu.Lock()
			c := *a.cfg
			c.WizardDone = true
			a.cfg = &c
			a.mu.Unlock()
			if err := saveConfig("config.json", &c); err != nil {
				log.Printf("сохранение конфига: %v", err)
			}
			log.Printf("мастер первого запуска пройден")
		})

		log.Printf("openSettings: WebView2 создан, открываю страницу")
		settingsHwnd.Store(hwnd)
		defer settingsHwnd.Store(0)
		setDarkClientBackground(hwnd)
		applyDarkCaption(hwnd)
		makeBorderless(hwnd)
		applyWindowLimits(hwnd, settingsMinW, settingsMinH, settingsMaxW, settingsMaxH)
		procSetWindowPos.Call(hwnd, 0, 0, 0, 0, 0, 0x0002|0x0001|0x0004|0x0020)
		var shown atomic.Bool
		reveal := func() {
			if !shown.CompareAndSwap(false, true) {
				return
			}
			w.Dispatch(func() { revealWindowCentered(hwnd, winW, winH) })
		}
		_ = w.Bind("appReady", reveal)
		go func() {
			time.Sleep(2 * time.Second)
			reveal()
		}()
		w.SetHtml(settingsHTML(a.snapshot(), tab))
		w.Run()
		log.Printf("openSettings: окно закрыто")
		savedDPI := dpiFor(settingsHwnd.Load())
		if savedDPI < 72 {
			savedDPI = 96
		}
		dipW, dipH := lastWndW*96/savedDPI, lastWndH*96/savedDPI
		if dipW >= settingsMinW && dipH >= settingsMinH && dipW <= settingsMaxW && dipH <= settingsMaxH {
			nw, nh := int(dipW), int(dipH)
			a.mu.Lock()
			changed := a.cfg.SettingsW != nw || a.cfg.SettingsH != nh
			if changed {
				c := *a.cfg
				c.SettingsW, c.SettingsH = nw, nh
				a.cfg = &c
			}
			cur := a.cfg
			a.mu.Unlock()
			if changed {
				if err := saveConfig("config.json", cur); err != nil {
					log.Printf("сохранение конфига: %v", err)
				}
			}
		}
	}()
}

type saveResult struct {
	OK       bool   `json:"ok"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

func (a *App) applySettings(f *settingsForm) saveResult {
	if _, err := parseHotkey(f.Hotkey); err != nil {
		return saveResult{Severity: "error", Message: err.Error()}
	}
	combos := []string{f.Hotkey, f.TranslateHotkey, f.PauseHotkey}
	for _, p := range f.Profiles {
		combos = append(combos, p.Hotkey)
	}
	if dup := hotkeys.FindDuplicate(combos); dup != "" {
		return saveResult{Severity: "error", Message: trf("err.hotkey.dup", dup)}
	}

	a.mu.Lock()
	old := a.cfg
	c := *old
	c.Hotkey = f.Hotkey
	c.UILanguage = f.UILanguage
	c.Beep = f.Beep
	if validSoundTheme(f.SoundTheme) {
		c.SoundTheme = f.SoundTheme
	}
	c.AutoEnter = f.AutoEnter
	c.RestoreClipboard = f.RestoreClipboard
	c.Overlay = f.Overlay
	if validOverlayPos(f.OverlayPos) {
		c.OverlayPos = f.OverlayPos
	}
	c.OverlayText = f.OverlayText
	setOverlayPos(c.OverlayPos)
	c.Animation = f.Animation
	if f.TypeMode {
		c.PasteMode = "type"
	} else {
		c.PasteMode = "clipboard"
	}
	c.Language = f.Language
	modelChanged := false
	if f.ModelID != "" && f.ModelID != "custom" {
		if m := findModel(f.ModelID); m != nil && m.installed() {
			before := primaryEngine(&c)
			switch m.Engine {
			case engineSherpa:
				if nd := "models/" + m.Dir; nd != c.SherpaModel {
					c.SherpaModel = nd
					modelChanged = true
				}
			default:
				if nm := "models/" + m.File; nm != c.Model {
					c.Model = nm
					modelChanged = true
				}
			}
			if primaryEngine(&c) != before {
				modelChanged = true
			}
		}
	}
	if f.Threads > 0 {
		c.Threads = f.Threads
	}
	if f.MinRecordMs >= 0 {
		c.MinRecordMs = f.MinRecordMs
	}
	if f.PasteDelayMs >= 0 && f.PasteDelayMs <= 5000 {
		c.PasteDelayMs = f.PasteDelayMs
	}
	c.HistoryOn = f.HistoryOn
	if f.HistoryDays > 0 {
		c.HistoryDays = f.HistoryDays
	}
	if f.HistoryMax > 0 {
		c.HistoryMax = f.HistoryMax
	}
	c.HistorySkip = strings.TrimSpace(f.HistorySkip)
	if f.MaxRecordSeconds > 0 {
		c.MaxRecordSeconds = f.MaxRecordSeconds
	}
	c.ServerAutostart = f.ServerAutostart
	c.CheckUpdates = f.CheckUpdates
	if validPunctuation(f.Punctuation) {
		c.Punctuation = f.Punctuation
	}
	if validHotkeyMode(f.HotkeyMode) {
		c.HotkeyMode = f.HotkeyMode
	}
	if validUILevel(f.UILevel) {
		c.UILevel = f.UILevel
	}
	themeChanged := false
	if theme.ValidSkin(f.Skin) && f.Skin != c.Skin {
		c.Skin = f.Skin
		themeChanged = true
	}
	if theme.ValidColour(f.Theme) && f.Theme != c.Theme {
		c.Theme = f.Theme
		themeChanged = true
	}
	c.MicDevice = f.MicDevice
	c.MicDeviceName = f.MicDeviceName
	if f.ServerPort < 1024 || f.ServerPort > 65535 {
		return saveResult{Severity: "error", Message: trf("err.port", f.ServerPort)}
	}
	c.ServerPort = f.ServerPort
	c.ServerExe = f.ServerExe
	c.ServerURL = strings.TrimSpace(f.ServerURL)
	c.WhisperPrompt = strings.TrimSpace(f.WhisperPrompt)
	if c.Language != old.Language {
		syncDictionary(&c)
	}
	if f.Profiles != nil {
		var cleaned []Profile
		for _, p := range f.Profiles {
			p.Name = strings.TrimSpace(p.Name)
			p.Prompt = strings.TrimSpace(p.Prompt)
			if p.ID == "" || p.Name == "" {
				continue
			}
			if p.Hotkey != "" {
				if _, err := parseHotkey(p.Hotkey); err != nil {
					log.Printf("хоткей профиля %s отброшен: %v", p.ID, err)
					p.Hotkey = ""
				}
			}
			cleaned = append(cleaned, p)
		}
		c.Profiles = cleaned
	}
	c.TranslateHotkey = f.TranslateHotkey
	if c.TranslateHotkey != "" {
		if _, err := parseHotkey(c.TranslateHotkey); err != nil {
			c.TranslateHotkey = ""
		}
	}
	c.PauseHotkey = f.PauseHotkey
	if c.PauseHotkey != "" {
		if _, err := parseHotkey(c.PauseHotkey); err != nil {
			c.PauseHotkey = ""
		}
	}
	if validTranslateLang(f.TranslateTarget) {
		c.TranslateTarget = f.TranslateTarget
	}
	switch f.TranslateAsk {
	case "always", "timeout", "never":
		c.TranslateAsk = f.TranslateAsk
	}
	if f.TranslateAskSeconds >= 1 && f.TranslateAskSeconds <= 10 {
		c.TranslateAskSeconds = f.TranslateAskSeconds
	}
	if f.TranslateAskLangs != nil {
		var langs []string
		for _, l := range f.TranslateAskLangs {
			if validTranslateLang(l) {
				langs = append(langs, l)
			}
		}
		if len(langs) == 0 {
			return saveResult{Severity: "error", Message: tr("err.nolangs")}
		}
		c.TranslateAskLangs = langs
	}
	c.TranslateDefault = f.TranslateDefault
	c.DefaultProfile = ""
	if f.ActiveProfiles != nil {
		var aps []string
		for _, id := range f.ActiveProfiles {
			if profileByID(&c, id) != nil {
				aps = append(aps, id)
			}
		}
		c.ActiveProfiles = aps
	}
	if f.AppRules != nil {
		c.AppRules = apprules.Clean(f.AppRules)
	}
	if f.Replacements != nil {
		c.Replacements = replace.Clean(f.Replacements)
	}
	if f.Commands != nil {
		c.Commands = commands.Clean(f.Commands)
	}
	llmChanged := false
	if f.LLMModelFile != "" && !strings.ContainsAny(f.LLMModelFile, "/\\") && strings.HasSuffix(f.LLMModelFile, ".gguf") {
		if _, err := os.Stat(filepath.Join("models", f.LLMModelFile)); err == nil {
			np := "models/" + f.LLMModelFile
			if np != c.LLMModel {
				c.LLMModel = np
				llmChanged = true
			}
		}
	}
	hook := a.hook
	a.mu.Unlock()

	if err := saveConfig("config.json", &c); err != nil {
		log.Printf("сохранение конфига: %v — настройки оставлены прежними", err)
		return saveResult{Severity: "error", Message: trf("err.save", humanError(err))}
	}

	a.mu.Lock()
	a.cfg = &c
	a.mu.Unlock()

	if themeChanged {
		applyTheme(c.Skin, c.Theme)
		refreshWindowChrome()
		trayReloadIcons(trayIdle)
		a.refreshIdleUI()
		overlayRefresh()
	}

	if llmChanged {
		a.llmShutdown()
		log.Printf("модель редактора переключена: %s", c.LLMModel)
	}

	if hook != nil {
		hook.SetCombos(buildCombos(&c))
	}
	if len(c.ActiveProfiles) > 0 && llmInstalled(&c) {
		go func() {
			if _, err := a.ensureLLM(); err != nil {
				log.Printf("прогрев LLM: %v", err)
			}
		}()
	}
	initLang(c.UILanguage)
	if c.HistoryDays != old.HistoryDays || c.HistoryMax != old.HistoryMax {
		a.enforceHistory(&c)
	}
	a.refreshIdleUI()
	if modelChanged {
		a.requestServerRestart()
	}

	serverChanged := c.ServerPort != old.ServerPort ||
		c.Threads != old.Threads ||
		c.ServerURL != old.ServerURL ||
		c.ServerExe != old.ServerExe ||
		c.ServerAutostart != old.ServerAutostart

	log.Printf("настройки сохранены: hotkey=%s ui=%s model=%s сервер=%v", c.Hotkey, c.UILanguage, c.Model, serverChanged)
	if serverChanged && !modelChanged {
		a.requestServerRestart()
	}
	if modelChanged {
		return saveResult{OK: true, Severity: "ok", Message: tr("model.switching")}
	}
	if serverChanged {
		return saveResult{OK: true, Severity: "ok", Message: tr("srv.restarting")}
	}
	return saveResult{OK: true, Severity: "ok", Message: strS("S_SAVED")}
}

func settingsHTML(cfg *Config, tab string) string {
	cfgMap := map[string]any{
		"hotkey":                  cfg.Hotkey,
		"ui_language":             cfg.UILanguage,
		"beep":                    cfg.Beep,
		"sound_theme":             cfg.SoundTheme,
		"auto_enter":              cfg.AutoEnter,
		"restore_clipboard":       cfg.RestoreClipboard,
		"overlay":                 cfg.Overlay,
		"overlay_position":        cfg.OverlayPos,
		"overlay_text":            cfg.OverlayText,
		"animation":               cfg.Animation,
		"type_mode":               cfg.PasteMode == "type",
		"language":                cfg.Language,
		"model":                   cfg.Model,
		"threads":                 cfg.Threads,
		"min_record_ms":           cfg.MinRecordMs,
		"paste_delay_ms":          cfg.PasteDelayMs,
		"history":                 cfg.HistoryOn,
		"history_days":            cfg.HistoryDays,
		"history_max":             cfg.HistoryMax,
		"history_skip":            cfg.HistorySkip,
		"max_record_seconds":      cfg.MaxRecordSeconds,
		"server_autostart":        cfg.ServerAutostart,
		"check_updates":           cfg.CheckUpdates,
		"mic_device":              cfg.MicDevice,
		"punctuation":             cfg.Punctuation,
		"hotkey_mode":             cfg.HotkeyMode,
		"ui_level":                cfg.UILevel,
		"skin":                    cfg.Skin,
		"theme":                   cfg.Theme,
		"server_port":             cfg.ServerPort,
		"server_exe":              cfg.ServerExe,
		"server_exe_path":         resolveServerExe(cfg.ServerExe),
		"server_exe_default":      defaultServerExe,
		"server_exe_path_default": resolveServerExe(defaultServerExe),
		"server_url":              cfg.ServerURL,
		"whisper_prompt":          cfg.WhisperPrompt,
		"translate_default":       cfg.TranslateDefault,
		"active_profiles":         cfg.ActiveProfiles,
		"app_rules":               cfg.AppRules,
		"replacements":            cfg.Replacements,
		"commands":                cfg.Commands,
		"translate_hotkey":        cfg.TranslateHotkey,
		"pause_hotkey":            cfg.PauseHotkey,
		"translate_target":        cfg.TranslateTarget,
		"translate_ask":           cfg.TranslateAsk,
		"translate_ask_seconds":   cfg.TranslateAskSeconds,
		"translate_ask_langs":     cfg.TranslateAskLangs,
		"llm_model":               filepath.Base(cfg.LLMModel),
		"profiles":                cfg.Profiles,
		"_version":                appVersion,
		"_tab":                    tab,
		"_cpus":                   runtime.NumCPU(),
		"_wizard":                 !cfg.WizardDone,
	}
	cfgJSON, _ := json.Marshal(cfgMap)

	str := strS
	lMap := map[string]string{"nohot": "—"}
	for jsKey, sKey := range map[string]string{
		"dl": "S_DL", "del": "S_DEL", "mdlready": "S_MODEL_READY", "add": "S_PROF_ADD",
		"pname": "S_PROF_NAME", "pprompt": "S_PROF_PROMPT", "phot": "S_PROF_HOTKEY",
		"pset": "S_PROF_SET", "pclr": "S_PROF_CLEAR", "ptest": "S_PROF_TEST",
		"fitok": "S_FIT_OK", "fitwarn": "S_FIT_WARN", "fitbad": "S_FIT_BAD",
		"ram": "S_RAM", "hfph": "S_HF_PH", "nollm": "S_NO_LLM", "nollmp": "S_NO_LLM_PROF",
		"upd": "S_UPDATED", "pedit": "S_PROF_EDIT", "pclose": "S_PROF_CLOSE",
		"confirmdel": "S_CONFIRM_DEL", "delactive": "S_DEL_ACTIVE", "wizneedmodel": "S_WIZ_NEED_MODEL",
		"free": "S_FREE", "updnone": "S_UPD_NONE",
		"skinterminal": "S_SKIN_TERMINAL", "skinname": "S_SKIN",
		"wndmax": "S_WND_MAX", "wndrestore": "S_WND_RESTORE", "wndmin": "S_WND_MIN", "wndclose": "S_WND_CLOSE",
		"themeeditor": "S_THEME_EDITOR", "themeneon": "S_THEME_NEON",
		"badgemodels": "S_BADGE_MODELS", "badgemiss": "S_BADGE_MISS", "badgesystem": "S_BADGE_SYSTEM", "badgehist": "S_BADGE_HIST",
		"exewarn": "S_EXE_WARN", "exeedit": "S_PROF_EDIT", "resetask": "S_RESET_ALL_ASK", "resetbtn": "S_RESET_ALL_BTN",
		"updfound":   "S_UPD_FOUND",
		"micdefault": "S_MIC_DEFAULT", "micquiet": "S_MIC_QUIET", "get": "S_STATE_GET", "change": "S_CHANGE_MODEL",
		"remotewarn": "S_REMOTE_WARN", "remoteask": "S_REMOTE_ASK", "remotebadge": "S_REMOTE_BADGE",
		"ok": "S_OK", "cancel": "S_CANCEL", "dlask": "S_DL_ASK", "dlstart": "S_DL_START", "dlcancel": "S_DL_CANCEL", "nofound": "S_NOT_FOUND",
		"advrolemain": "S_ADV_ROLE_MAIN", "advrolesecond": "S_ADV_ROLE_SECOND",
		"advprimary": "S_ADV_PRIMARY", "advcompanion": "S_ADV_COMPANION", "advhave": "S_ADV_HAVE", "advapply": "S_ADV_APPLY", "advask": "S_ADV_ASK",
		"more": "S_MORE", "less": "S_LESS",
		"pasteinh": "S_RULE_PASTE_INH", "enterinh": "S_RULE_ENTER_INH", "delaynone": "S_RULE_DELAY_NONE", "promptinh": "S_RULE_PROMPT_INH",
		"ruleclip": "S_RULE_CLIP", "ruletype": "S_RULE_TYPE",
		"ruleenteron": "S_RULE_ENTER_ON", "ruleenteroff": "S_RULE_ENTER_OFF", "rulenoprompt": "S_RULE_NOPROMPT",
		"rulelast": "S_RULE_LAST", "ruleempty": "S_RULE_EMPTY", "ruledel": "S_RULE_DEL",
		"ruleprompts": "S_RULE_PROMPTS", "ruleph": "S_RULE_PH",
		"replempty": "S_REPL_EMPTY", "repldel": "S_REPL_DEL", "replwhole": "S_REPL_WHOLE",
		"cmdempty": "S_CMD_EMPTY", "cmddel": "S_CMD_DEL", "cmdph": "S_CMD_PH",
		"cmdnewline": "S_CMD_NEWLINE", "cmdparagraph": "S_CMD_PARAGRAPH", "cmdtext": "S_CMD_TEXT", "cmdcancel": "S_CMD_CANCEL",
		"cmdtextph":   "S_CMD_TEXT_PH",
		"cmdpnewline": "S_CMD_P_NEWLINE", "cmdpparagraph": "S_CMD_P_PARAGRAPH", "cmdpcancel": "S_CMD_P_CANCEL",
		"histempty": "S_HIST_EMPTY", "histcopy": "S_HIST_COPY", "histask": "S_HIST_ASK", "histclear": "S_HIST_CLEAR",
		"micchecking": "S_MIC_CHECKING", "mchecking": "S_MCHECK_RUN", "histinsert": "S_HIST_INSERT",
		"slotru": "S_STATE_RU", "slotother": "S_STATE_OTHER",
		"replcase": "S_REPL_CASE", "replfromph": "S_REPL_FROM_PH", "repltoph": "S_REPL_TO_PH",
		"wiznext": "S_WIZ_NEXT", "wizfinish": "S_WIZ_FINISH", "wizwait": "S_WIZ_WAIT",
		"wizheard": "S_WIZ_HEARD", "wizhave": "S_WIZ_HAVE", "wiztry": "S_WIZ_TRY_TEXT",
		"updavail": "S_UPD_AVAIL", "updgo": "S_UPD_GO", "upderr": "S_UPD_ERR", "upddl": "S_UPD_DL",
	} {
		lMap[jsKey] = str(sKey)
	}
	lJSON, _ := json.Marshal(lMap)

	pairs := []string{"{{CFG}}", string(cfgJSON), "{{L_JSON}}", string(lJSON), "{{APP}}", appid.Name,
		"{{THEME_VARS}}", theme.Current(cfg.Skin, cfg.Theme).CSSVars(), "{{THEME_LIST}}", skinListJSON(), "{{FONT_FACE}}", plexfont.FaceCSS()}
	remote := strings.TrimSpace(cfg.ServerURL) != ""
	for k := range settingsStrings["en"] {
		v := str(k)
		if k == "S_ABOUT_HTML" && remote {
			v += `<p class="warn">` + str("S_REMOTE_ABOUT") + `</p>`
		}
		pairs = append(pairs, "{{"+k+"}}", v)
	}
	return strings.NewReplacer(pairs...).Replace(settingsPage)
}

const settingsPage = `<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>{{S_TITLE}}</title><style>{{FONT_FACE}}
:root{{{THEME_VARS}}}
*{box-sizing:border-box;margin:0;padding:0}
html,body{height:100%}
body{font:var(--fs)/1.45 var(--font);background:var(--bg);color:var(--green);user-select:none;display:flex;flex-direction:column;overflow:hidden;border:var(--wborder)}
body::after{content:"";position:fixed;inset:0;pointer-events:none;opacity:var(--scan);background:repeating-linear-gradient(transparent 0 2px,rgba(0,0,0,.10) 2px 3px)}
.content{flex:1;overflow-y:auto;overflow-x:hidden;min-height:0;scrollbar-gutter:stable both-edges}
::-webkit-scrollbar{width:10px}
::-webkit-scrollbar-track{background:var(--bg)}
::-webkit-scrollbar-thumb{background:var(--line);border:2px solid var(--bg)}
::-webkit-scrollbar-thumb:hover{background:var(--dim)}
.header{display:flex;align-items:center;gap:14px;padding:12px 12px 12px 20px;overflow:hidden;border-bottom:1px solid var(--line);background:var(--titlebg);box-shadow:0 1px 12px rgba(var(--rgb),.12);cursor:default}
.capbtns{display:flex;gap:6px;margin-left:10px;flex:none}
.rsz{position:fixed;z-index:60}
.rsz.t{top:0;left:9px;right:9px;height:5px;cursor:n-resize}
.rsz.b{bottom:0;left:9px;right:9px;height:5px;cursor:s-resize}
.rsz.l{left:0;top:9px;bottom:9px;width:5px;cursor:w-resize}
.rsz.r{right:0;top:9px;bottom:9px;width:5px;cursor:e-resize}
.rsz.tl{left:0;top:0;width:9px;height:9px;cursor:nw-resize}
.rsz.tr{right:0;top:0;width:9px;height:9px;cursor:ne-resize}
.rsz.bl{left:0;bottom:0;width:9px;height:9px;cursor:sw-resize}
.rsz.br{right:0;bottom:0;width:9px;height:9px;cursor:se-resize}
button.cap{width:36px;height:30px;background:none;border:1px solid var(--line);color:var(--dim);font:14px var(--font);cursor:pointer;padding:0;border-radius:calc(var(--r) * .5)}
button.cap:hover{background:var(--on);color:var(--green);box-shadow:var(--glow)}
button.cap.max{font-size:12px}
button.cap.close:hover{background:var(--badbg);color:var(--bad);border-color:var(--badline);box-shadow:var(--badglow)}
.logo{width:40px;height:40px;flex:none}
.logo svg{width:100%;height:100%;filter:var(--iconglow)}
.mk.mic{display:var(--markmic,block)}
.mk.face{display:var(--markface,none)}
.header h1{flex:0 1 auto;min-width:0;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;font-size:15px;font-weight:var(--wb);letter-spacing:var(--brandls);text-shadow:var(--glow);animation:var(--flicker);background:var(--brandbg);-webkit-background-clip:var(--brandclip);background-clip:var(--brandclip);-webkit-text-fill-color:var(--brandfill)}
.statusbar .ver{margin-left:auto;flex:none;color:var(--dim)}
@keyframes flicker{0%,93%,97%,100%{opacity:1}95%{opacity:.6}}
@keyframes pulse{0%,100%{opacity:.35;transform:scale(.94)}50%{opacity:1;transform:scale(1)}}
.wave{animation:pulse 1.6s infinite}
@media (prefers-reduced-motion:reduce){
 .header h1{animation:none}
 .wave{animation:none;opacity:.75}
 *{transition-duration:.01ms !important}
}
.tabs{display:flex;flex-wrap:wrap;gap:2px;padding:10px 16px 0;border-bottom:1px solid var(--line)}
.shell{display:flex;flex:1;min-height:0}
.snav{width:clamp(178px,24%,210px);flex:none;border-right:1px solid var(--line);background:var(--sidebg);padding:6px 0;display:flex;flex-direction:column;gap:1px;overflow-y:auto;overflow-x:hidden}
.snav .ngrp{font-size:9px;letter-spacing:.16em;color:var(--dim);padding:12px 12px 3px;text-transform:uppercase}
.nav{appearance:none;border:0;background:none;text-align:left;font:inherit;font-size:12.5px;color:var(--dim);padding:8px 12px;cursor:pointer;border-left:2px solid transparent}
.nav:hover{color:var(--green)}
.nav.active{color:var(--green);border-left-color:var(--hi);background:var(--navon);text-shadow:var(--glow)}
.nav{display:flex;align-items:center;gap:6px;min-width:0}
.nav .nlabel{flex:1;min-width:0;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.nbadge{flex:none;max-width:64px;font-size:9px;padding:1px 5px;border:1px solid var(--line);border-radius:calc(var(--r) * .4);color:var(--dim);white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.nbadge:empty{display:none}
.nbadge.warn{color:var(--amber);border-color:var(--amber)}
.nbadge.miss{color:var(--amber);border-color:var(--amber);box-shadow:var(--amberglow)}
.scard .led{width:6px;height:6px;border-radius:50%;background:var(--faint);display:inline-block;margin-right:6px;flex:none}
.scard .led.on{background:var(--ok);box-shadow:var(--higlow)}
.scard .led.warn{background:var(--amber)}
.scard .mini{align-self:flex-start;margin-top:auto}
.scard .v{white-space:nowrap;overflow:hidden;text-overflow:ellipsis;display:block}
.scard .v .led{vertical-align:middle;margin-right:6px}
.scard .miclevel{margin-top:auto}
.tip{position:fixed;z-index:200;max-width:320px;padding:5px 9px;border:var(--bw) solid var(--line);background:var(--panel);color:var(--green);font:11.5px var(--font);line-height:1.45;box-shadow:var(--shadow);pointer-events:none;white-space:pre-wrap;opacity:0;transition:opacity .12s;border-radius:calc(var(--r) * .6)}
.tip.on{opacity:1}
.scard .miclevel{flex:none}
.row .sub{display:block;font-size:10.5px;color:var(--dim);margin-top:2px;letter-spacing:0}
.row .lbl{flex:1;min-width:0}
.statusbar{border-top:1px solid var(--line);padding:6px 14px;display:flex;gap:12px;align-items:center;font-size:11px;color:var(--dim);flex-wrap:nowrap;white-space:nowrap}
.statusbar #st_main{min-width:0;overflow:hidden;text-overflow:ellipsis}
.statusbar .stsaved,.statusbar .stremote{flex:none}
.statusbar .stremote{color:var(--amber);border:1px solid var(--amber);border-radius:calc(var(--r) * .4);padding:0 5px;letter-spacing:.08em}
.statusbar .stremote:empty{display:none;border:0;padding:0}
.statusbar .led{width:6px;height:6px;border-radius:50%;background:var(--faint);flex:none}
.statusbar .led.on{background:var(--hi);box-shadow:var(--higlow)}
.statusbar .stsaved{color:var(--dim)}
.statusbar .stsaved.warn{color:var(--amber)}
.statusbar .stsaved.bad{color:var(--bad)}
.note.warn{font-size:11px;color:var(--amber);line-height:1.5;padding:2px 0 6px}
.note.warn:empty{display:none}
.omni{margin-left:auto;display:flex;align-items:center;gap:8px;flex:0 1 auto;min-width:0;width:min(320px,38%);border:1px solid var(--line);border-radius:calc(var(--r) * .7);background:var(--panel);padding:4px 9px}
.omni:focus-within{border-color:var(--dim);box-shadow:var(--glow)}
.omni input[type=text]{flex:1;min-width:0;width:auto;background:none;border:0;box-shadow:none;outline:none;color:var(--green);font:inherit;font-size:11.5px;padding:0;user-select:text;-webkit-user-select:text}
.omni input[type=text]:focus{border:0;box-shadow:none}
.omni input::placeholder{color:var(--dim)}
.omni svg{flex:none;color:var(--faint)}
.omni:focus-within svg{color:var(--dim)}
.omni .okey{flex:none;font-size:9px;letter-spacing:.06em;color:var(--dim);border:1px solid var(--line);border-radius:3px;padding:1px 5px}
.omni .ocount{flex:none;font-size:10px;color:var(--dim);white-space:nowrap}
.omni .ocount.none{color:var(--amber)}
.wiz{position:fixed;inset:0;z-index:30;background:var(--bg);display:none;flex-direction:column}
.wiz.on{display:flex}
.wizhead{display:flex;align-items:center;gap:14px;padding:12px 12px 12px 20px;border-bottom:1px solid var(--line);box-shadow:0 1px 12px rgba(var(--rgb),.12)}
.wizhead h2{font-size:14px;font-weight:var(--wb);letter-spacing:2px;text-shadow:var(--glow);flex:1;min-width:0;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.wizbody{flex:1;min-height:0;overflow-y:auto;padding:22px 26px}
.wizstep{display:none;flex-direction:column;gap:14px}
.wizstep.on{display:flex}
.wizh{font-size:15px;font-weight:var(--wb);letter-spacing:2px;text-shadow:var(--glow)}
.wiztext{color:var(--dim);font-size:13px;line-height:1.6;max-width:62ch}
.wizfoot{display:flex;align-items:center;gap:10px;padding:11px 20px;border-top:1px solid var(--line)}
.wizdots{display:flex;gap:6px;align-items:center}
.wizgap{flex:1}
.wizdots i{width:7px;height:7px;border:1px solid var(--line);display:block}
.wizdots i.on{background:var(--hi);border-color:var(--hi);box-shadow:var(--higlow)}
.wizrow{display:flex;align-items:center;gap:10px;flex-wrap:wrap}
.wizrow select{min-width:170px}
.wizrow select:focus{border-color:var(--dim);box-shadow:var(--glow);outline:none}
.wizkey{appearance:none;font:inherit;color:var(--green);border:1px solid var(--line);background:var(--keybg);border-radius:calc(var(--r) * .6);padding:4px 11px;letter-spacing:1px;text-shadow:var(--glow);white-space:nowrap;cursor:pointer}
.wizkey:hover{border-color:var(--dim);box-shadow:var(--glow)}
.wizplan{display:flex;flex-direction:column;gap:4px;border:1px solid var(--line);border-radius:var(--r);background:var(--panel);padding:10px 12px;min-height:38px}
.wizbar{height:14px;border:1px solid var(--line);border-radius:calc(var(--r) * .6);background:var(--panel);position:relative;overflow:hidden;flex:1;min-width:150px;max-width:430px}
.wizbar i{position:absolute;left:0;top:0;bottom:0;width:0;background:linear-gradient(90deg,var(--faint),var(--hi));box-shadow:var(--higlow);transition:width .15s linear}
.wizlvl{display:flex;align-items:flex-end;gap:2px;height:10px}
.wizlvl i{display:block;width:var(--lvlw,4px);height:10px;background:var(--soft);border:1px solid var(--line);box-sizing:border-box;border-radius:var(--lvlr,0);transition:background .12s linear}
.wizlvl i.on{background:var(--hi);border-color:var(--hi);box-shadow:var(--higlow)}
.wiztry{width:100%;min-height:76px;resize:none;background:var(--field);border:1px solid var(--line);border-radius:calc(var(--r) * .55);color:var(--green);font:inherit;font-size:var(--ctlfs);padding:var(--fieldpad);outline:none;user-select:text}
.wiztry:focus{border-color:var(--dim);box-shadow:var(--glow)}
.wizout{font-size:12.5px;color:var(--dim);min-height:18px}
.wizout.ok{color:var(--green);text-shadow:var(--glow)}
.wizbig{font-size:30px;text-shadow:var(--glow);line-height:1}
button.btn.ghost{border-color:var(--line);background:none;color:var(--dim);filter:none}
button.btn.ghost:hover{color:var(--green);border-color:var(--dim);background:none}
.lvlsw{display:flex;flex:none;gap:2px;border:1px solid var(--line);border-radius:calc(var(--r) * .55);background:var(--panel);padding:1px}
.lvlb{appearance:none;border:0;background:none;color:var(--dim);font:inherit;font-size:10px;letter-spacing:.12em;text-transform:uppercase;padding:4px 11px;cursor:pointer;border-radius:calc(var(--r) * .45)}
.lvlb:hover{color:var(--dim)}
.lvlb.on{background:var(--selbg);color:var(--selfg);text-shadow:none}
.row.hidden{display:none}
.page.advopen .row[data-adv]{border-left:2px solid var(--line);padding-left:9px;margin-left:-11px}
.row.hit{background:var(--navon);box-shadow:inset 2px 0 0 var(--hi)}
.moreb{appearance:none;background:none;border:1px dashed var(--line);border-radius:calc(var(--r) * .5);color:var(--dim);font:inherit;font-size:11px;padding:6px 10px;margin:6px 0 0;cursor:pointer}
.moreb:hover{color:var(--dim);border-color:var(--dim)}
.modal-bg{position:fixed;inset:0;background:var(--scrim);display:flex;align-items:center;justify-content:center;z-index:20}
.modal{background:var(--panel);border:1px solid var(--line);border-radius:var(--r);box-shadow:0 0 24px rgba(var(--rgb),.18),var(--shadow);padding:20px 22px;max-width:380px;display:flex;flex-direction:column;gap:16px}
.modal p{font-size:13px;line-height:1.55;color:var(--green)}
.modal-btns{display:flex;gap:10px;justify-content:flex-end}
.modal .btn{padding:7px 18px;border:1px solid var(--btnline);border-radius:calc(var(--r) * .5);background:var(--btnbg);color:var(--btnfg);font:inherit;font-size:12px;letter-spacing:var(--ls);text-transform:var(--caps);cursor:pointer}
.modal .btn:hover{filter:brightness(1.12);box-shadow:var(--glow)}
.modal .btn.ghost{border-color:var(--line);background:none;color:var(--dim);filter:none}
.modal .btn.ghost:hover{color:var(--green);border-color:var(--dim)}
.modal .btn:focus-visible{outline:1px solid var(--green);outline-offset:2px}
.hero{display:flex;align-items:center;gap:12px;border:1px solid var(--line);border-radius:var(--r);background:var(--panel);padding:12px 14px;margin-bottom:10px;flex-wrap:wrap}
.herokey{border:1px solid var(--line);background:var(--keybg);border-radius:calc(var(--r) * .6);color:var(--green);padding:5px 12px;font-size:14px;font-weight:var(--wb);letter-spacing:1px;text-shadow:var(--glow)}
.herotext{font-size:12px;color:var(--dim)}
.berr{display:flex;align-items:center;gap:10px;flex-wrap:wrap;border:1px solid var(--bad);border-radius:var(--r);background:var(--panel);padding:10px 12px;margin-bottom:10px}
.berr.upd{border-color:var(--amber)}
.berr.upd .berrtext{color:var(--amber)}
.berr .berrtext{flex:1 1 220px;color:var(--bad);min-width:0}
.cards{display:grid;grid-template-columns:repeat(auto-fit,minmax(150px,1fr));gap:8px;margin-bottom:12px}
.scard{border:1px solid var(--line);border-radius:var(--r);background:var(--panel);padding:9px 11px;display:flex;flex-direction:column;gap:5px}
.scard .k{font-size:9.5px;letter-spacing:.12em;color:var(--dim);text-transform:uppercase}
.scard .v{font-size:12.5px;color:var(--green)}
.lastres{display:block;max-width:100%;color:var(--dim);font-size:12px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
#state_last_meta{overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.tab{padding:9px 14px;border:1px solid transparent;border-bottom:none;border-radius:calc(var(--r) * .5) calc(var(--r) * .5) 0 0;background:none;font:inherit;color:var(--dim);cursor:pointer;letter-spacing:var(--ls);text-transform:var(--caps);font-size:12px}
.tab:hover{color:var(--green)}
.tab.active{color:var(--green);border-color:var(--line);background:var(--panel);text-shadow:var(--glow)}
.page{display:none;padding:14px 16px;max-width:900px}
.page.active{display:block}
.card{background:none;border:0;padding:0;margin:0 0 18px}
.row{display:flex;align-items:center;gap:10px;padding:9px 0;flex-wrap:wrap;border-bottom:1px solid var(--soft)}
.card .row:last-child{border-bottom:0}
.row label{flex:1;min-width:100px;color:var(--green)}
.row label .sub{display:block;font-size:10.5px;color:var(--dim);margin-top:2px;letter-spacing:0}
.row label .sub.warn{color:var(--amber)}
.row select{flex:0 1 auto;min-width:0;max-width:100%}
.row input[type=text]{flex:0 1 auto;min-width:0}
.row .hint{font-size:11px;color:var(--dim)}
select{color-scheme:var(--scheme,dark)}
option{background:var(--bg);color:var(--green)}
option:checked{background:linear-gradient(var(--on),var(--on));color:var(--green)}
select,::picker(select){appearance:base-select}
::picker(select){background:var(--bg);border:1px solid var(--line);border-radius:calc(var(--r) * .55);padding:2px;margin-top:2px;color:var(--green);box-shadow:var(--shadow)}
::picker(select) option{padding:6px 10px;background:none;color:var(--dim);border:0;border-radius:calc(var(--r) * .4);font:inherit;min-height:0}
::picker(select) option:hover,::picker(select) option:focus{background:var(--on);color:var(--green);outline:none}
::picker(select) option:checked{color:var(--green);text-shadow:var(--glow)}
option::checkmark{display:none}
select::picker-icon{color:var(--faint)}
select:open::picker-icon{transform:rotate(180deg)}
select:open{border-color:var(--dim)}
input[type=text],input[type=number],select{padding:var(--fieldpad);border:1px solid var(--line);border-radius:calc(var(--r) * .55);background:var(--field);color:var(--green);font:inherit;font-size:var(--ctlfs);line-height:1.2;outline:none}
input:focus,select:focus{border-color:var(--dim);box-shadow:var(--glow)}
input::placeholder{color:var(--dim)}
input:disabled,select:disabled{opacity:.35;cursor:default}
#trlangs label:has(input:disabled){opacity:.45}
input[type=text]{width:220px;max-width:100%}select{width:210px;max-width:100%}
input[type=range]{width:150px;accent-color:var(--dim);background:transparent}
input[type=checkbox]{appearance:none;-webkit-appearance:none;width:32px;height:17px;border:1px solid var(--dim);border-radius:calc(var(--r) * .8);position:relative;cursor:pointer;background:none;flex:none;padding:0;margin:0}
input[type=checkbox]::before,input[type=radio]::before{content:"";position:absolute;top:-11px;bottom:-11px;left:-9px;right:-9px}
input[type=checkbox]::after{content:"";position:absolute;top:2px;left:2px;width:11px;height:11px;border-radius:calc(var(--r) * .6);background:var(--dim);transition:.15s}
input[type=checkbox]:checked{border-color:var(--dim)}
input[type=checkbox]:checked::after{left:17px;background:var(--hi);box-shadow:var(--higlow)}
input[type=checkbox]:focus-visible{outline:1px solid var(--green);outline-offset:2px}
.row select,.row input[type=text]{border:1px solid var(--line);background:var(--field);color:var(--green)}
.row select{flex:0 0 auto;width:auto;min-width:118px;max-width:min(320px,100%)}
.row input[type=text]{flex:0 0 auto;width:min(230px,50%)}
.row .val{color:var(--dim);font-size:11.5px;min-width:44px;text-align:right}
button,select,input[type=checkbox],input[type=radio],input[type=range]{cursor:pointer}
button:disabled{cursor:default;opacity:.35}
.val{min-width:52px;text-align:right;color:var(--green);text-shadow:var(--glow);font-weight:var(--wb)}
.hotkey-val{appearance:none;font:inherit;font-weight:var(--wb);font-size:15px;color:var(--green);background:var(--keybg);border:1px solid var(--line);border-radius:calc(var(--r) * .6);padding:8px 14px;min-width:150px;text-align:center;text-shadow:var(--glow);letter-spacing:1px;cursor:pointer}
.hotkey-val:hover{border-color:var(--dim);box-shadow:var(--glow)}
button.btn{padding:8px 18px;border:1px solid var(--btnline);border-radius:calc(var(--r) * .5);background:var(--btnbg);color:var(--btnfg);font:inherit;cursor:pointer;letter-spacing:var(--ls);text-transform:var(--caps);font-size:12px}
button.btn:hover{filter:brightness(1.12);box-shadow:var(--glow)}
button.ghost{border-color:var(--line);background:none;color:var(--dim);filter:none}
button.ghost:hover{color:var(--green)}
.footer{flex:none;display:flex;gap:12px;align-items:center;padding:10px 16px;background:var(--panel);border-top:1px solid var(--line)}
.toast{color:var(--amber);font-size:13px;opacity:0;transition:opacity .3s;text-shadow:var(--amberglow)}
.toast.show{opacity:1}
.rulerow{display:flex;align-items:center;gap:8px;flex-wrap:wrap;padding:8px 0;border-bottom:1px solid var(--soft)}
.rulerow:last-child{border-bottom:none}
.rulerow input[type=text]{flex:1 1 190px;min-width:150px;width:auto}
.rulerow select{flex:0 0 auto;width:auto;min-width:118px}
.rdel{flex:none;border:1px solid var(--line);border-radius:calc(var(--r) * .5);background:none;color:var(--dim);font:inherit;font-size:12px;cursor:pointer;padding:4px 9px}
.rdel:hover{color:var(--bad);border-color:var(--badline)}
.histrow{display:flex;align-items:flex-start;gap:10px;padding:9px 2px;border-bottom:1px solid var(--soft)}
.histrow:last-child{border-bottom:none}
.histmeta{flex:none;width:150px;font-size:10.5px;color:var(--dim);line-height:1.5}
.histmeta b{display:block;color:var(--dim);font-weight:400}
.histtext{flex:1;min-width:0;font-size:12.5px;line-height:1.5;color:var(--green);user-select:text;overflow-wrap:anywhere}
.histrow .mini{flex:none}
.histfind{margin:8px 0 4px;max-width:320px}
.histempty{color:var(--dim);font-size:12px;padding:8px 0}
.micverdict{font-size:12.5px;line-height:1.5;color:var(--dim);padding:4px 0;min-height:18px}
.micverdict.ok{color:var(--green);text-shadow:var(--glow)}
.micverdict.bad{color:var(--amber)}
.sect .mini{margin-left:auto}
.replrow{display:flex;align-items:center;gap:8px;flex-wrap:wrap;padding:7px 0;border-bottom:1px solid var(--soft)}
.replrow:last-child{border-bottom:none}
.replrow input[type=text]{flex:1 1 160px;min-width:120px;width:auto}
.replrow .rarrow{flex:none;color:var(--dim)}
.replrow label{flex:none;display:flex;align-items:center;gap:6px;font-size:11px;color:var(--dim);white-space:nowrap}
.replrow label input[type=checkbox]{width:28px;height:15px}
.replcheck{display:flex;align-items:center;gap:10px;margin-top:12px;padding-top:10px;border-top:1px solid var(--soft);flex-wrap:wrap}
.replcheck input[type=text]{flex:1 1 220px;min-width:160px;width:auto}
.replout{flex:1 1 200px;min-width:0;font-size:12px;color:var(--green);text-shadow:var(--glow);overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.replout:empty{display:none}
.rulefoot{display:flex;align-items:center;gap:10px;margin-top:10px;flex-wrap:wrap}
.rulefoot .ghost{border-color:var(--line);color:var(--dim)}
.rulefoot .ghost:empty{display:none}
.ruleempty{color:var(--dim);font-size:12px;padding:6px 0}
.card>.hint{font-size:11.5px;color:var(--dim);margin-bottom:6px;line-height:1.5}
.mslot{font-size:9.5px;letter-spacing:.12em;text-transform:uppercase;color:var(--dim);padding:11px 2px 3px;border-bottom:1px solid var(--line);margin-bottom:2px}
.mslot.hidden{display:none}
.mrow{display:flex;align-items:center;gap:9px;padding:7px 2px;border-bottom:1px solid var(--soft);flex-wrap:nowrap}
.mrow .mdesc{min-width:0;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
@media (max-width:820px){.mrow .mram{display:none}}
.mrow:last-child{border-bottom:none}
input[type=radio]{appearance:none;-webkit-appearance:none;width:15px;height:15px;flex:none;margin:0;padding:0;border:1px solid var(--dim);border-radius:50%;background:var(--field);position:relative;cursor:pointer}
input[type=radio]:checked{border-color:var(--hi)}
input[type=radio]:checked::after{content:"";position:absolute;top:3px;left:3px;width:7px;height:7px;border-radius:50%;background:var(--hi);box-shadow:var(--higlow)}
input[type=radio]:disabled{opacity:.4;cursor:default}
input[type=radio]:focus-visible{outline:1px solid var(--green);outline-offset:2px}
button:focus-visible,input:focus-visible,select:focus-visible,textarea:focus-visible{outline:1px solid var(--green);outline-offset:2px}
.mrow .mname{width:132px;font-weight:var(--wb);white-space:nowrap}
.mrow .mdesc{flex:1;color:var(--dim);font-size:12px}
.mtag{font-size:9px;border:1px solid var(--line);border-radius:calc(var(--r) * .35);color:var(--dim);padding:0 4px;margin-left:6px;vertical-align:middle;letter-spacing:.06em}
.rrow{display:flex;align-items:center;gap:9px;padding:6px 2px;font-size:12px;color:var(--dim);flex-wrap:wrap}
.rrow .rcond{min-width:132px;color:var(--green)}
.rrow .rarr{color:var(--dim)}
.rrow .reng{border:1px solid var(--line);border-radius:calc(var(--r) * .4);padding:1px 7px;color:var(--green);font-size:11.5px}
.rrow .rwhy{margin-left:auto;color:var(--dim);font-size:11px}
#routing{border-bottom:1px solid var(--soft);margin-bottom:8px;padding-bottom:6px}
.advbar{display:flex;align-items:center;gap:8px;flex-wrap:wrap;margin-bottom:8px}
.fchips{display:flex;gap:5px;flex-wrap:wrap;margin-left:auto}
.fchip{appearance:none;background:none;border:1px solid var(--line);border-radius:calc(var(--r) * .5);color:var(--dim);font:inherit;font-size:10.5px;padding:2px 8px;cursor:pointer}
.fchip:hover{color:var(--dim)}
.fchip.on{background:var(--selbg);color:var(--selfg);border-color:transparent}
#advisor{border:1px solid var(--line);border-radius:var(--r);padding:10px 11px;margin-bottom:9px;display:flex;flex-direction:column;gap:9px}
.advrow{display:flex;align-items:center;gap:9px;font-size:12px;padding:3px 0}
.advrow .advrole{width:74px;color:var(--dim);text-transform:uppercase;font-size:10px;letter-spacing:.1em}
.advrow .advname{flex:1;color:var(--green)}
.advrow .advwhy{display:block;color:var(--dim);font-size:11px;letter-spacing:0;text-transform:none}
.advrow .advstate{color:var(--amber)}
.advrow .advstate.ok{color:var(--dim)}
.advq{display:flex;gap:10px;align-items:center;flex-wrap:wrap;font-size:11.5px;color:var(--dim)}
.advq select{margin-left:5px}
.advchk{display:flex;align-items:center;gap:5px}
.advout{font-size:12px;color:var(--green);line-height:1.5;min-height:1em}
.mrow.hidden{display:none}
.mram{color:var(--dim);font-size:10.5px;width:74px;text-align:right}
.mram.warn{color:var(--amber)}
.mram.bad{color:var(--bad)}
.mrow .msize{color:var(--dim);font-size:12px;width:70px;text-align:right}
.badge{font-size:11px;letter-spacing:1px;padding:4px 10px;border:1px solid var(--line);border-radius:calc(var(--r) * .5);color:var(--green);text-shadow:var(--glow);text-transform:var(--caps)}
button.mini{flex:none;padding:5px 12px;border:1px solid var(--line);border-radius:calc(var(--r) * .5);background:none;color:var(--dim);font:inherit;font-size:11.5px;cursor:pointer;letter-spacing:var(--ls);text-transform:var(--caps)}
button.mini:hover{color:var(--green);border-color:var(--dim);box-shadow:var(--glow)}
button.mini.danger:hover{color:var(--bad);border-color:var(--badline);background:var(--badbg);box-shadow:var(--badglow)}
.mpct{color:var(--amber);font-size:12px;min-width:44px;text-align:right;text-shadow:var(--amberglow)}
.sect{color:var(--dim);font-weight:400;font-size:11px;letter-spacing:.14em;text-transform:uppercase;margin:0 0 4px;display:flex;align-items:center;gap:8px}
.hfhome{margin-left:auto;cursor:pointer;color:var(--dim);font-size:11px;letter-spacing:1px;border:1px solid var(--line);border-radius:calc(var(--r) * .5);padding:3px 9px;text-transform:none}
.hfhome:hover{color:var(--green);border-color:var(--dim);box-shadow:var(--glow)}
.ramline{display:flex;align-items:center;flex-wrap:wrap;gap:6px;color:var(--dim);font-size:12px;margin:4px 0 10px}
.ramline b{color:var(--green);font-size:14px;font-weight:var(--wb);text-shadow:var(--glow);margin-right:4px}
.ramline .dot{margin-left:12px;font-size:10px}
.subhead{color:var(--dim);font-size:11px;letter-spacing:1px;text-transform:uppercase;margin:14px 0 2px;padding-top:10px;border-top:1px solid var(--soft)}
#hf_results{max-height:44vh;overflow-y:auto;overscroll-behavior:contain}
.miclevel{flex:none;display:flex;align-items:flex-end;gap:2px;height:10px;width:auto}
.miclevel i{display:block;width:var(--lvlw,4px);height:10px;background:var(--soft);border:1px solid var(--line);box-sizing:border-box;border-radius:var(--lvlr,0);transition:background .1s linear}
.miclevel i.on{background:var(--hi);border-color:var(--hi);box-shadow:var(--higlow)}
.miclevel.grow{width:100%;overflow:hidden}
.lvlrow .miclevel.grow{flex:0 0 auto;width:min(320px,100%);min-width:0}
#hf_clr{appearance:none;background:none;border:0;font:inherit;position:absolute;right:9px;top:50%;transform:translateY(-50%);color:var(--dim);cursor:pointer;display:none;font-size:13px;padding:2px 4px}
#hf_clr:hover{color:var(--green);text-shadow:var(--glow)}
#hf_go{appearance:none;background:none;border:0;font:inherit;position:absolute;left:9px;top:50%;transform:translateY(-50%);color:var(--dim);cursor:pointer;line-height:0;padding:3px}
#hf_go:hover{color:var(--green);filter:drop-shadow(0 0 4px rgba(var(--rgb),.6))}
.hfrepo{appearance:none;background:none;border:0;font:inherit;color:inherit;display:flex;align-items:center;gap:9px;flex:1;min-width:0;text-align:left;cursor:pointer;padding:2px 0}
.hfrepo .mdesc{flex:1;color:var(--dim);min-width:0;overflow:hidden;text-overflow:ellipsis}
.hfupd{color:var(--dim);font-size:12px;flex:none}
.hflink{appearance:none;background:none;border:0;font:inherit;color:var(--dim);cursor:pointer;padding:0 4px}
.hflink:hover{color:var(--green)}
button.iconbtn{border:none;background:none;padding:2px 5px;color:var(--dim);cursor:pointer;line-height:1;font:inherit;font-size:13px}
button.iconbtn:hover{color:var(--green);filter:drop-shadow(0 0 4px rgba(var(--rgb),.6))}
button.iconbtn.danger:hover{color:var(--bad);filter:var(--badfilter)}
.about p{margin:8px 0;line-height:1.55;user-select:text;color:var(--green);max-width:80ch}
.about li{max-width:78ch}
.toc{display:flex;flex-wrap:wrap;gap:6px 14px;margin:6px 0 4px;max-width:80ch}
.swatches{display:flex;gap:6px;margin-left:auto}
.swatch{width:16px;height:16px;border:1px solid var(--line);border-radius:calc(var(--r) * .35);cursor:pointer;padding:0;background:none;overflow:hidden}
.swatch i{display:block;width:100%;height:100%}
.swatch.on{border-color:var(--hi);box-shadow:var(--higlow)}
.swatch:focus-visible{outline:1px solid var(--green);outline-offset:2px}
.toc a{color:var(--dim);font-size:12px;text-decoration:none;border-bottom:1px dotted var(--line)}
.toc a:hover{color:var(--green);border-bottom-color:var(--dim)}
.toc a:focus-visible{outline:1px solid var(--green);outline-offset:2px}
.about .wh:target,.about .wh.hit{color:var(--green);text-shadow:var(--glow)}
.about p.hit,.about li.hit{background:var(--navon);box-shadow:inset 2px 0 0 var(--hi);padding-left:7px}
.about p.warn{color:var(--amber);border-left:2px solid var(--amber);padding-left:9px}
.about b{font-weight:var(--wb);text-shadow:var(--glow)}
.about .wh{color:var(--dim);font-size:12px;letter-spacing:1px;text-transform:uppercase;margin:16px 0 4px;border-bottom:1px solid var(--soft);padding-bottom:3px}
.about ul{margin:4px 0 10px 20px;padding:0}
.about li{margin:4px 0;line-height:1.55;color:var(--green);user-select:text}
.mock{background:var(--panel);border:1px solid var(--line);border-radius:var(--r);padding:10px 14px;margin:8px 0;max-width:420px}
.mock-pill{display:flex;align-items:center;gap:10px}
.mock-dot{width:11px;height:11px;border-radius:50%;background:var(--rec);box-shadow:0 0 8px var(--rec);flex:none}
.mock-bars{display:flex;gap:2px;align-items:center}
.mock-bars i{display:block;width:3px;border-radius:var(--barr,0);background:var(--hi)}
.mock-ask{display:flex;align-items:center;gap:8px;margin-top:9px;padding-top:9px;border-top:1px solid var(--line);flex-wrap:wrap;font-size:13px;color:var(--dim)}
.mock-x{margin-left:auto;color:var(--dim)}
.mock-btn{border:1px solid var(--line);border-radius:calc(var(--r) * .5);padding:5px 14px;color:var(--dim);font-size:13px}
.mock-btn.on{border-color:var(--green);color:var(--green);text-shadow:var(--glow)}
.mock-cd{position:relative}
.mock-cd::after{content:"";position:absolute;left:0;right:38%;bottom:-4px;height:2px;background:var(--hi);box-shadow:var(--higlow)}
.mock-mi{padding:4px 6px;color:var(--green);font-size:13px}
.mock-mi.dim{color:var(--faint)}
.mock-sep{border:none;border-top:1px solid var(--soft);margin:4px 0}
.mock-row{display:flex;align-items:center;gap:8px;padding:4px 0;font-size:13px;border-bottom:1px solid var(--soft)}
.mock-row:last-child{border-bottom:none}
.mock-radio{width:12px;height:12px;border-radius:50%;border:1px solid var(--dim);flex:none}
.mock-radio.on{background:var(--hi);box-shadow:var(--higlow)}
.mock-cb{width:13px;height:13px;border:1px solid var(--dim);flex:none;display:inline-flex;align-items:center;justify-content:center;font-size:10px;color:var(--bg)}
.mock-cb.on{background:var(--hi)}
.mock-note{color:var(--dim);font-size:12px}
.lnk{color:var(--green);text-decoration:underline;cursor:pointer}
.lnk:hover{text-shadow:var(--glow)}
</style></head><body>
<div class="header">
 <div class="logo"><svg viewBox="0 0 64 64">
  <rect x="2" y="2" width="60" height="60" rx="12" fill="var(--panel)" stroke="var(--line)" stroke-width="2"/>
  <g class="mk mic" stroke="var(--hi)" stroke-width="4" fill="none" stroke-linecap="round">
   <rect x="26" y="12" width="12" height="20" rx="6" fill="var(--hi)"/>
   <path d="M19 27a13 13 0 0 0 26 0"/>
   <line x1="32" y1="40" x2="32" y2="46"/>
   <line x1="24" y1="49" x2="40" y2="49"/>
  </g>
  <g class="mk face" fill="none" stroke-linecap="round">
   <circle cx="32" cy="35" r="19" fill="var(--on)" stroke="var(--selbg)" stroke-width="3"/>
   <path d="M19 19l4 8M45 19l-4 8" stroke="var(--selbg)" stroke-width="3.5"/>
   <circle cx="26" cy="33" r="2.8" fill="var(--hi)"/>
   <circle cx="38" cy="33" r="2.8" fill="var(--hi)"/>
   <path d="M28 41q4 3.4 8 0" stroke="var(--hi)" stroke-width="3"/>
  </g>
  <g stroke="var(--hi)" stroke-width="2.5" fill="none" stroke-linecap="round">
   <path class="wave" d="M13 20a17 17 0 0 0 0 14" style="animation-delay:.2s"/>
   <path class="wave" d="M51 20a17 17 0 0 1 0 14" style="animation-delay:.6s"/>
  </g>
 </svg></div>
 <h1>{{APP}}</h1>
 <span class="lvlsw">
  <button type="button" class="lvlb" data-l="simple">{{S_LVL_SIMPLE}}</button>
  <button type="button" class="lvlb" data-l="all">{{S_LVL_ALL}}</button>
 </span>
 <label class="omni"><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4"><circle cx="10.5" cy="10.5" r="6.5"/><line x1="15.5" y1="15.5" x2="21" y2="21"/></svg><input id="omni" type="text" placeholder="{{S_SEARCH}}" autocomplete="off"><span class="ocount" id="ocount"></span><span class="okey">Ctrl K</span></label>
 <div class="capbtns">
  <button class="cap" onclick="appMin()" title="{{S_WND_MIN}}">&#9472;</button>
  <button class="cap" id="cap_max" onclick="toggleMax()" title="{{S_WND_MAX}}">&#9744;</button>
  <button class="cap close" onclick="appClose()" title="{{S_WND_CLOSE}}">&#10005;</button>
 </div>
</div>
<div class="rsz t" data-edge="12"></div>
<div class="rsz b" data-edge="15"></div>
<div class="rsz l" data-edge="10"></div>
<div class="rsz r" data-edge="11"></div>
<div class="rsz tl" data-edge="13"></div>
<div class="rsz tr" data-edge="14"></div>
<div class="rsz bl" data-edge="16"></div>
<div class="rsz br" data-edge="17"></div>
<div class="tip" id="tip" role="tooltip"></div>
<div class="shell">
<nav class="snav" id="snav" role="tablist">
 <span class="ngrp">{{S_GRP_WORK}}</span>
 <button class="nav" role="tab" aria-selected="false" data-p="state"><span class="nlabel">{{S_NAV_STATE}}</span></button>
 <button class="nav" role="tab" aria-selected="false" data-p="dictation"><span class="nlabel">{{S_NAV_DICT}}</span></button>
 <button class="nav" role="tab" aria-selected="false" data-p="history"><span class="nlabel">{{S_NAV_HISTORY}}</span><span class="nbadge" id="badge_history"></span></button>
 <button class="nav" role="tab" aria-selected="false" data-p="mic"><span class="nlabel">{{S_NAV_MIC}}</span><span class="nbadge" id="badge_mic"></span></button>
 <span class="ngrp">{{S_GRP_REC}}</span>
 <button class="nav" role="tab" aria-selected="false" data-p="models"><span class="nlabel">{{S_NAV_MODELS}}</span><span class="nbadge" id="badge_models"></span></button>
 <button class="nav" role="tab" aria-selected="false" data-p="text"><span class="nlabel">{{S_NAV_TEXT}}</span></button>
 <button class="nav" role="tab" aria-selected="false" data-p="translate"><span class="nlabel">{{S_NAV_TR}}</span></button>
 <span class="ngrp">{{S_GRP_OTHER}}</span>
 <button class="nav" role="tab" aria-selected="false" data-p="system"><span class="nlabel">{{S_NAV_SYSTEM}}</span><span class="nbadge warn" id="badge_system"></span></button>
 <button class="nav" role="tab" aria-selected="false" data-p="about"><span class="nlabel">{{S_NAV_ABOUT}}</span></button>
</nav>

<div class="content">
<div class="page" role="tabpanel" aria-hidden="true" id="p-state">
 <div class="hero"><span class="herokey" id="state_hotkey"></span><span class="herotext">{{S_STATE_HINT}}</span></div>
 <div class="berr" id="state_backend" style="display:none">
  <span class="berrtext" id="state_backend_text"></span>
  <button type="button" class="mini" id="state_retry">{{S_RETRY}}</button>
  <button type="button" class="mini ghost" data-goto="system">{{S_BERR_OPEN}}</button>
 </div>
 <div class="berr upd" id="state_upd" style="display:none">
  <span class="berrtext" id="state_upd_text"></span>
  <button type="button" class="mini" data-goto="system">{{S_UPD}}</button>
 </div>
 <div class="cards">
  <div class="scard"><span class="k">{{S_NAV_MIC}}</span>
   <span class="v"><i class="led" id="state_mic_led"></i><span id="state_mic">—</span></span>
   <span class="miclevel grow" id="state_mic_bar"><i></i><i></i><i></i><i></i><i></i><i></i><i></i></span></div>
  <div class="scard"><span class="k">{{S_STATE_RU}}</span>
   <span class="v"><i class="led" id="state_ru_led"></i><span id="state_ru">—</span></span>
   <button class="mini" id="state_ru_btn" data-goto="models">{{S_CHANGE_MODEL}}</button></div>
  <div class="scard"><span class="k">{{S_STATE_OTHER}}</span>
   <span class="v"><i class="led" id="state_other_led"></i><span id="state_other">—</span></span>
   <button class="mini" id="state_other_btn" data-goto="models">{{S_CHANGE_MODEL}}</button></div>
  <div class="scard"><span class="k">{{S_STATE_PROC}}</span>
   <span class="v"><i class="led" id="state_llm_led"></i><span id="state_llm">—</span></span>
   <button class="mini" id="state_llm_btn" data-goto="models">{{S_PICK_MODEL}}</button></div>
 </div>
 <h2 class="sect">{{S_STATE_LAST}}</h2>
 <div class="row"><span class="lbl"><span class="lastres" id="state_last">—</span><span class="sub" id="state_last_meta"></span></span>
  <button class="mini" id="state_copy">{{S_STATE_COPY}}</button></div>
 <div class="row"><span class="lbl">{{S_STATE_MEM}}<span class="sub">{{S_STATE_MEM_SUB}}</span></span>
  <span class="val" id="state_ram">—</span></div>
</div>

<div class="page" role="tabpanel" aria-hidden="true" id="p-history">
 <div class="card">
  <div class="row"><label>{{S_HIST_ON}}<span class="sub">{{S_HIST_ON_SUB}}</span></label><input type="checkbox" id="history"></div>
  <div class="row" data-adv><label>{{S_HIST_DAYS}}</label>
   <select id="history_days"><option value="1">1</option><option value="3">3</option><option value="7">7</option><option value="30">30</option></select></div>
  <div class="row" data-adv><label>{{S_HIST_MAX}}</label>
   <select id="history_max"><option value="50">50</option><option value="100">100</option><option value="200">200</option><option value="500">500</option></select></div>
  <div class="row" id="hist_skip_row"><label>{{S_HIST_SKIP}}<span class="sub">{{S_HIST_SKIP_SUB}}</span></label><input type="text" id="history_skip"></div>
 </div>
 <div class="card" id="histcard">
  <h2 class="sect">{{S_HIST_LIST}}<button type="button" class="mini" id="hist_clear">{{S_HIST_CLEAR}}</button></h2>
  <label class="omni histfind"><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4"><circle cx="10.5" cy="10.5" r="6.5"/><line x1="15.5" y1="15.5" x2="21" y2="21"/></svg><input id="hist_find" type="text" placeholder="{{S_HIST_FIND}}" autocomplete="off"></label>
  <div id="histbody"></div>
 </div>
</div>
<div class="page" role="tabpanel" aria-hidden="true" id="p-dictation">
 <div class="card">
  <div class="row"><label>{{S_HOTKEY}}</label>
   <button type="button" class="hotkey-val" id="hotkey" onclick="appCapture()" title="{{S_CHANGE}}"></button></div>
  <div class="row"><label>{{S_HOTMODE}}<span class="sub">{{S_SUB_HOTMODE}}</span></label>
   <select id="hotkey_mode"><option value="hold">{{S_HOTMODE_HOLD}}</option><option value="toggle">{{S_HOTMODE_TOGGLE}}</option></select></div>
  <div class="row" data-adv><label>{{S_PAUSE}}<span class="sub">{{S_PAUSE_SUB}}</span></label>
   <button type="button" class="hotkey-val" id="pause_hotkey" style="min-width:110px" title="{{S_PROF_SET}}"></button>
   <button class="mini" id="pause_clear">{{S_PROF_CLEAR}}</button></div>
  <div class="row" data-adv><label>{{S_MINMS}}<span class="sub">{{S_SUB_MINMS}}</span></label><select id="min_record_ms"><option value="0">0 ms</option><option value="100">100 ms</option><option value="150">150 ms</option><option value="200">200 ms</option><option value="300">300 ms</option><option value="500">500 ms</option><option value="750">750 ms</option><option value="1000">1000 ms</option></select></div>
  <div class="row" data-adv><label>{{S_MAXSEC}}</label><select id="max_record_seconds"><option value="30">30 s</option><option value="60">60 s</option><option value="120">120 s</option><option value="180">180 s</option><option value="300">300 s</option></select></div>
 </div>
 <div class="card">
  <h2 class="sect">{{S_SEC_BEHAVIOR}}</h2>
  <div class="row"><label>{{S_AUTOENTER}}<span class="sub">{{S_SUB_ENTER}}</span></label><input type="checkbox" id="auto_enter"></div>
  <div class="row" data-adv><label>{{S_RESTORE}}<span class="sub">{{S_SUB_CLIP}}</span></label><input type="checkbox" id="restore_clipboard"></div>
  <div class="row"><label>{{S_TYPEMODE}}<span class="sub">{{S_SUB_TYPE}}</span></label><input type="checkbox" id="type_mode"></div>
  <div class="row" data-adv><label>{{S_PASTE_DELAY}}<span class="sub">{{S_PASTE_DELAY_SUB}}</span></label>
   <select id="paste_delay_ms"><option value="0">0 ms</option><option value="50">50 ms</option><option value="100">100 ms</option><option value="250">250 ms</option><option value="500">500 ms</option><option value="1000">1000 ms</option></select></div>
 </div>
 <div class="card">
  <h2 class="sect">{{S_SEC_RULES}}</h2>
  <div class="hint">{{S_RULES_HINT}}</div>
  <div id="rulesbody"></div>
  <div class="rulefoot">
   <button type="button" class="mini" id="rule_add">{{S_RULE_ADD}}</button>
   <button type="button" class="mini ghost" id="rule_last"></button>
  </div>
 </div>
 <div class="card">
  <h2 class="sect">{{S_SEC_OVERLAY}}</h2>
  <div class="row"><label>{{S_OVERLAY}}</label><input type="checkbox" id="overlay"></div>
  <div class="row"><label>{{S_OVPOS}}<span class="sub">{{S_OVPOS_SUB}}</span></label>
   <select id="overlay_position">
    <option value="bottom">{{S_OVPOS_BOTTOM}}</option>
    <option value="top">{{S_OVPOS_TOP}}</option>
    <option value="caret">{{S_OVPOS_CARET}}</option>
   </select></div>
  <div class="row"><label>{{S_OVTEXT}}<span class="sub">{{S_OVTEXT_SUB}}</span></label><input type="checkbox" id="overlay_text"></div>
  <div class="row" data-adv><label>{{S_ANIM}}</label><input type="checkbox" id="animation"></div>
 </div>
</div>

<div class="page" role="tabpanel" aria-hidden="true" id="p-mic">
 <div class="card">
  <div class="row"><label>{{S_MIC}}</label>
   <select id="mic_device"><option value="">{{S_MIC_DEFAULT}}</option></select>
   <button class="iconbtn" id="mic_refresh" title="{{S_MIC_REFRESH}}">&#8635;</button></div>
  <div class="row lvlrow"><label>{{S_MIC_LEVEL}}</label>
   <span class="mpct" id="mic_hint"></span>
   <span class="miclevel grow" id="mic_bar"><i></i><i></i><i></i><i></i><i></i><i></i><i></i></span></div>
  <div class="row"><label>{{S_MIC_CHECK}}<span class="sub">{{S_MIC_CHECK_SUB}}</span></label>
   <button type="button" class="mini" id="mic_check">{{S_PROF_TEST}}</button></div>
  <div class="micverdict" id="mic_verdict"></div>
  <div class="note warn" id="mic_err"></div>
 </div>
 <div class="card">
  <h2 class="sect">{{S_SEC_SOUND}}</h2>
  <div class="row"><label>{{S_BEEP}}</label><input type="checkbox" id="beep"></div>
  <div class="row" data-adv><label>{{S_SOUND}}</label>
   <button class="mini" onclick="appPreviewSound(document.getElementById('sound_theme').value)">&#9654;</button>
   <select id="sound_theme">
    <option value="speech">{{S_SND_SPEECH}}</option>
    <option value="chime">{{S_SND_CHIME}}</option>
    <option value="soft">{{S_SND_SOFT}}</option>
    <option value="marimba">{{S_SND_MARIMBA}}</option>
    <option value="blip">{{S_SND_BLIP}}</option>
    <option value="pop">{{S_SND_POP}}</option>
   </select></div>
 </div>
</div>

<div class="page" role="tabpanel" aria-hidden="true" id="p-models">
 <div class="card">
  <div id="routing"></div>
  <div class="advbar">
   <button type="button" class="mini" id="adv_open">{{S_ADV_TITLE}}</button>
   <span class="fchips">
    <button type="button" class="fchip on" data-f="all">{{S_F_ALL}}</button>
    <button type="button" class="fchip" data-f="multi">{{S_F_MULTI}}</button>
    <button type="button" class="fchip" data-f="punct">{{S_F_PUNCT}}</button>
    <button type="button" class="fchip" data-f="fit">{{S_F_FIT}}</button>
    <button type="button" class="fchip" data-f="ru">{{S_F_RU}}</button>
   </span>
  </div>
  <div id="advisor" style="display:none">
   <div class="advq">
    <label>{{S_ADV_LANGQ}}
     <select id="adv_lang"><option value="ru">RU</option><option value="en">EN</option><option value="multi">{{S_F_MULTI}}</option></select>
    </label>
    <label>{{S_ADV_PRIOQ}}
     <select id="adv_prio"><option value="balance">·</option><option value="accuracy">{{S_ADV_ACC}}</option><option value="speed">{{S_ADV_SPEED}}</option></select>
    </label>
    <label class="advchk"><input type="checkbox" id="adv_tr"> {{S_ADV_TRQ}}</label>
    <button type="button" class="mini" id="adv_go">{{S_ADV_GO}}</button>
   </div>
   <div id="adv_out" class="advout"></div>
  </div>
  <div id="models"></div>
  <div class="row"><label>{{S_MCHECK}}<span class="sub">{{S_MCHECK_SUB}}</span></label>
   <button type="button" class="mini" id="mcheck">{{S_MCHECK_GO}}</button></div>
  <div class="micverdict" id="mcheck_out"></div>
 </div>
 <div class="card">
  <div class="row"><label>{{S_RECLANG}}</label>
   <select id="language">
    <option value="auto">{{S_RECAUTO}}</option>
    <option value="ru">Русский</option><option value="en">English</option>
    <option value="uk">Українська</option><option value="de">Deutsch</option>
    <option value="fr">Français</option><option value="es">Español</option>
    <option value="it">Italiano</option><option value="pl">Polski</option>
   </select></div>
  <div class="row" data-adv><label>{{S_THREADS}}<span class="sub">{{S_SUB_THREADS}}</span></label><select id="threads"><option value="1">1</option><option value="2">2</option><option value="4">4</option><option value="6">6</option><option value="8">8</option><option value="12">12</option><option value="16">16</option></select></div>
 </div>
 <div class="card">
  <h2 class="sect">{{S_SEC_LLM}}<span class="hfhome" onclick="appHFHome()" title="huggingface.co">Hugging Face ↗</span></h2>
  <div id="proc-models" class="spage on"></div>
  <div id="proc-search" class="spage on"></div>
 </div>
</div>

<div class="page" role="tabpanel" aria-hidden="true" id="p-text">
 <div class="card">
  <div class="row"><label>{{S_PUNCT}}<span class="sub">{{S_SUB_PUNCT}}</span></label>
   <select id="punctuation">
    <option value="model">{{S_PUNCT_MODEL}}</option>
    <option value="llm">{{S_PUNCT_LLM}}</option>
    <option value="off">{{S_PUNCT_OFF}}</option>
   </select></div>
 </div>
 <div class="card">
  <h2 class="sect">{{S_SUB_DICT}}</h2>
  <div class="hint">{{S_DICT_HINT}}</div>
  <textarea id="whisper_prompt" rows="10" style="width:100%;min-height:150px;height:26vh;padding:8px 11px;border:1px solid var(--line);background:var(--field);color:var(--green);font:inherit;line-height:1.5;outline:none;resize:vertical"></textarea>
 </div>
 <div class="card" data-adv>
  <h2 class="sect">{{S_SEC_REPLACE}}</h2>
  <div class="hint">{{S_REPLACE_HINT}}</div>
  <div id="replbody"></div>
  <div class="rulefoot">
   <button type="button" class="mini" id="repl_add">{{S_REPL_ADD}}</button>
  </div>
 </div>
 <div class="card" data-adv>
  <h2 class="sect">{{S_SEC_CMD}}</h2>
  <div class="hint">{{S_CMD_HINT}}</div>
  <div id="cmdbody"></div>
  <div class="rulefoot">
   <button type="button" class="mini" id="cmd_add">{{S_CMD_ADD}}</button>
   <button type="button" class="mini ghost" id="cmd_preset">{{S_CMD_PRESET}}</button>
  </div>
  <div class="rulefoot">
   <span class="hint" style="margin:0;flex:1 1 auto">{{S_LISTS_HINT}}</span>
   <button type="button" class="mini ghost" id="lists_export">{{S_LISTS_EXPORT}}</button>
   <button type="button" class="mini ghost" id="lists_import">{{S_LISTS_IMPORT}}</button>
  </div>
  <div class="replcheck">
   <input type="text" id="repl_test" placeholder="{{S_REPL_TEST_PH}}">
   <span class="replout" id="repl_out"></span>
  </div>
 </div>
 <div class="card" data-adv>
  <h2 class="sect">{{S_SUB_PROMPTS}}</h2>
  <div class="hint">{{S_LLM_HINT}}</div>
  <div id="profbody"></div>
 </div>
</div>

<div class="page" role="tabpanel" aria-hidden="true" id="p-translate">
 <div class="card">
  <div class="hint">{{S_TR_HINT}}</div>
  <div id="tr_warn" style="display:none;color:var(--amber);font-size:12px;margin-bottom:6px">{{S_TR_TURBO}}</div>
  <div class="row"><label>{{S_TR_DEFAULT}}</label><input type="checkbox" id="tr_default"></div>
  <div class="row"><label>{{S_TR_TARGET}}<span class="sub">{{S_SUB_TRTARGET}}</span><span class="sub warn">{{S_TR_EXP}}</span></label>
   <select id="translate_target">
    <option value="en">English</option><option value="uk">Українська</option>
    <option value="de">Deutsch</option><option value="fr">Français</option>
    <option value="es">Español</option><option value="it">Italiano</option>
    <option value="pl">Polski</option><option value="ru">Русский</option>
   </select></div>
  <div class="row"><label>{{S_TR_ASK}}</label>
   <select id="translate_ask">
    <option value="never">{{S_TR_ASK_NEVER}}</option>
    <option value="always">{{S_TR_ASK_ALWAYS}}</option>
    <option value="timeout">{{S_TR_ASK_TIMEOUT}}</option>
   </select></div>
  <div class="row" data-adv><label>{{S_TR_SECONDS}}</label><select id="translate_ask_seconds"><option value="2">2 s</option><option value="3">3 s</option><option value="4">4 s</option><option value="5">5 s</option><option value="7">7 s</option><option value="10">10 s</option></select></div>
  <div class="row" data-adv><label>{{S_PROF_HOTKEY}}</label>
   <button type="button" class="hotkey-val" id="tr_hotkey" style="min-width:110px" title="{{S_PROF_SET}}"></button>
   <button class="mini" id="tr_clear">{{S_PROF_CLEAR}}</button></div>
  <div class="row" data-adv><label>{{S_TR_LANGS}}</label>
   <span id="trlangs" style="display:flex;gap:9px;flex-wrap:wrap">
    <label style="flex:none"><input type="checkbox" id="tl_en"> EN</label>
    <label style="flex:none"><input type="checkbox" id="tl_de"> DE</label>
    <label style="flex:none"><input type="checkbox" id="tl_fr"> FR</label>
    <label style="flex:none"><input type="checkbox" id="tl_es"> ES</label>
    <label style="flex:none"><input type="checkbox" id="tl_it"> IT</label>
    <label style="flex:none"><input type="checkbox" id="tl_pl"> PL</label>
    <label style="flex:none"><input type="checkbox" id="tl_ru"> RU</label>
    <label style="flex:none"><input type="checkbox" id="tl_uk"> UK</label>
   </span></div>
 </div>
</div>

<div class="page" role="tabpanel" aria-hidden="true" id="p-system">
 <div class="card">
  <div class="row"><label>{{S_UILANG}}</label>
   <select id="ui_language">
    <option value="auto">{{S_AUTO}}</option>
    <option value="en">English</option>
    <option value="ru">Русский</option>
    <option value="uk">Українська</option>
    <option value="de">Deutsch</option>
    <option value="fr">Français</option>
    <option value="es">Español</option>
    <option value="it">Italiano</option>
    <option value="pl">Polski</option>
   </select></div>
  <div class="row"><label>{{S_SKIN}}<span class="sub">{{S_SKIN_SUB}}</span></label>
   <select id="skin">
    <option value="terminal">{{S_SKIN_TERMINAL}}</option>
    <option value="editor">{{S_THEME_EDITOR}}</option>
    <option value="neon">{{S_THEME_NEON}}</option>
    <option value="soft">{{S_SKIN_SOFT}}</option>
    <option value="paper">{{S_SKIN_PAPER}}</option>
   </select></div>
  <div class="row" id="colour_row"><label>{{S_THEME}}<span class="sub">{{S_THEME_SUB}}</span></label>
   <span class="swatches" id="theme_swatches"></span>
   <select id="theme">
    <option value="green">{{S_THEME_GREEN}}</option>
    <option value="amber">{{S_THEME_AMBER}}</option>
    <option value="blue">{{S_THEME_BLUE}}</option>
    <option value="pink">{{S_THEME_PINK}}</option>
   </select></div>
  <div class="row"><label>{{S_AUTORUN}}<span class="sub">{{S_AUTORUN_SUB}}</span></label><input type="checkbox" id="autorun"></div>
  <div class="row"><label>{{S_UPD}}</label>
   <button class="mini" id="upd_check">{{S_UPD_CHECK}}</button></div>
  <div class="row" data-adv><label>{{S_UPD_AUTO}}<span class="sub">{{S_SUB_UPD}}</span></label><input type="checkbox" id="check_updates"></div>
  <div id="upd_out" style="font-size:12px;min-height:18px;color:var(--amber)"></div>
  <div class="row" data-adv><label>{{S_LOG}}<span class="sub">{{S_LOG_SUB}}</span></label>
   <button class="mini" id="log_open">{{S_LOG_OPEN}}</button></div>
  <div class="row" data-adv><label>{{S_RELOAD_CFG}}<span class="sub">{{S_RELOAD_CFG_SUB}}</span></label>
   <button class="mini" id="cfg_reload">{{S_RELOAD_CFG_BTN}}</button></div>
  <div class="row" data-adv><label>{{S_RESET_ALL}}<span class="sub">{{S_RESET_ALL_SUB}}</span></label>
   <button class="mini danger" id="cfg_reset">{{S_RESET_ALL_BTN}}</button></div>
 </div>
 <div class="card">
  <h2 class="sect" data-adv>{{S_SEC_SERVICE}}</h2>
  <div class="row" data-adv><label>{{S_AUTOSTART}}<span class="sub">{{S_SUB_AUTOSTART}}</span></label><input type="checkbox" id="server_autostart"></div>
  <div class="row" data-adv><label>{{S_PORT}}<span class="sub">{{S_SUB_PORT}}</span></label><input type="text" id="server_port" style="width:90px"></div>
  <div class="row" data-adv><label>{{S_SERVEREXE}}<span class="sub">{{S_SERVEREXE_SUB}}</span></label>
   <input type="text" id="server_exe" readonly>
   <button class="mini" id="exe_edit">{{S_PROF_EDIT}}</button>
   <button class="mini" id="exe_reset">{{S_EXE_RESET}}</button></div>
  <div class="row" data-adv><label>{{S_SERVERURL}}<div class="hint">{{S_URLHINT}}</div></label><input type="text" id="server_url"></div>
  <div class="note warn" id="remote_warn"></div>
 </div>
</div>

<div class="page about" role="tabpanel" aria-hidden="true" id="p-about">
 <div class="card">
  <p style="font-size:15px;letter-spacing:2px"><b>{{APP}}</b> <span id="ver2"></span></p>
  {{S_ABOUT_HTML}}
 </div>
 <div class="card">{{S_HELP_HTML}}</div>
 <div class="card">{{S_AUTHOR_HTML}}</div>
</div>
</div>
</div>

<div class="statusbar" id="statusbar">
 <span class="led" id="st_led"></span>
 <span id="st_main">—</span>
 <span class="stsaved" id="st_saved"></span>
 <span class="stremote" id="st_remote"></span>
 <span class="ver">v<span id="ver"></span></span>
</div>

<div class="wiz" id="wiz">
 <div class="wizhead" id="wizhead">
  <div class="logo"><svg viewBox="0 0 64 64">
   <rect x="2" y="2" width="60" height="60" rx="12" fill="var(--panel)" stroke="var(--line)" stroke-width="2"/>
   <g class="mk mic" stroke="var(--hi)" stroke-width="4" fill="none" stroke-linecap="round">
    <rect x="26" y="12" width="12" height="20" rx="6" fill="var(--hi)"/>
    <path d="M19 27a13 13 0 0 0 26 0"/>
    <line x1="32" y1="40" x2="32" y2="46"/>
    <line x1="24" y1="49" x2="40" y2="49"/>
   </g>
   <g class="mk face" fill="none" stroke-linecap="round">
    <circle cx="32" cy="35" r="19" fill="var(--on)" stroke="var(--selbg)" stroke-width="3"/>
    <path d="M19 19l4 8M45 19l-4 8" stroke="var(--selbg)" stroke-width="3.5"/>
    <circle cx="26" cy="33" r="2.8" fill="var(--hi)"/>
    <circle cx="38" cy="33" r="2.8" fill="var(--hi)"/>
    <path d="M28 41q4 3.4 8 0" stroke="var(--hi)" stroke-width="3"/>
   </g>
   <g stroke="var(--hi)" stroke-width="2.5" fill="none" stroke-linecap="round">
    <path class="wave" d="M13 20a17 17 0 0 0 0 14" style="animation-delay:.2s"/>
    <path class="wave" d="M51 20a17 17 0 0 1 0 14" style="animation-delay:.6s"/>
   </g>
  </svg></div>
  <h2>{{APP}}</h2>
  <div class="capbtns">
   <button class="cap" onclick="appMin()">&#9472;</button>
   <button class="cap close" onclick="appClose()">&#10005;</button>
  </div>
 </div>
 <div class="wizbody" id="wizbody">
  <div class="wizstep on" id="wz0">
   <div class="wizh">{{S_WIZ_HELLO}}</div>
   <div class="wiztext">{{S_WIZ_HELLO_TEXT}}</div>
   <div class="wizrow"><label for="wiz_ui">{{S_UILANG}}</label>
    <select id="wiz_ui">
     <option value="auto">{{S_AUTO}}</option>
     <option value="en">English</option>
     <option value="ru">Русский</option>
     <option value="uk">Українська</option>
     <option value="de">Deutsch</option>
     <option value="fr">Français</option>
     <option value="es">Español</option>
     <option value="it">Italiano</option>
     <option value="pl">Polski</option>
    </select></div>
   <div class="wiztext">{{S_WIZ_LATER}}</div>
  </div>
  <div class="wizstep" id="wz1">
   <div class="wizh">{{S_WIZ_T_MODEL}}</div>
   <div class="wiztext">{{S_WIZ_MODEL_TEXT}}</div>
   <div class="wizrow"><label for="wiz_lang">{{S_RECLANG}}</label>
    <select id="wiz_lang">
     <option value="ru">Русский</option><option value="en">English</option>
     <option value="uk">Українська</option><option value="de">Deutsch</option>
     <option value="fr">Français</option><option value="es">Español</option>
     <option value="it">Italiano</option><option value="pl">Polski</option>
     <option value="auto">{{S_RECAUTO}}</option>
    </select></div>
   <div class="wizplan" id="wiz_plan"></div>
   <div class="wizrow" id="wiz_dlrow" style="display:none"><span class="wizbar"><i id="wiz_dlbar"></i></span><span class="mpct" id="wiz_dlpct"></span></div>
   <div class="wizrow"><button type="button" class="btn" id="wiz_dl">{{S_DL}}</button><span class="wizout" id="wiz_dlout"></span></div>
  </div>
  <div class="wizstep" id="wz2">
   <div class="wizh">{{S_WIZ_T_INPUT}}</div>
   <div class="wiztext">{{S_WIZ_INPUT_TEXT}}</div>
   <div class="wizrow"><label>{{S_HOTKEY}}</label>
    <button type="button" class="wizkey" id="wiz_hot" title="{{S_CHANGE}}">—</button></div>
   <div class="wizrow"><label for="wiz_mic">{{S_MIC}}</label><select id="wiz_mic"></select></div>
   <div class="wizrow"><span class="wizlvl" id="wiz_micbar"><i></i><i></i><i></i><i></i><i></i><i></i><i></i></span><span class="wizout" id="wiz_michint"></span></div>
  </div>
  <div class="wizstep" id="wz3">
   <div class="wizh">{{S_WIZ_T_TRY}}</div>
   <div class="wiztext" id="wiz_trytext"></div>
   <textarea class="wiztry" id="wiz_try" placeholder="{{S_WIZ_TRY_PH}}"></textarea>
   <div class="wizout" id="wiz_tryout"></div>
  </div>
  <div class="wizstep" id="wz4">
   <div class="wizbig">&#10003;</div>
   <div class="wizh">{{S_WIZ_T_DONE}}</div>
   <div class="wiztext">{{S_WIZ_DONE_TEXT}}</div>
   <div class="row" style="max-width:430px"><label class="lbl" for="wiz_auto">{{S_AUTORUN}}<span class="sub">{{S_AUTORUN_SUB}}</span></label><input type="checkbox" id="wiz_auto"></div>
  </div>
 </div>
 <div class="wizfoot">
  <span class="wizdots" id="wizdots"></span>
  <button type="button" class="btn ghost" id="wiz_skip">{{S_WIZ_SKIP}}</button>
  <span class="wizgap"></span>
  <button type="button" class="btn ghost" id="wiz_back">{{S_WIZ_BACK}}</button>
  <button type="button" class="btn" id="wiz_next">{{S_WIZ_NEXT}}</button>
 </div>
</div>



<script>
window.onerror = function(m, s, l, c){ if(window.appJSError) appJSError(String(m) + " @line " + l + ":" + c); };
const CFG = {{CFG}};
const bools = ["beep","auto_enter","restore_clipboard","overlay","overlay_text","animation","type_mode","server_autostart","check_updates","history"];
const texts = ["history_skip"];
let exeStored = CFG.server_exe || "";
let exeUnlocked = false;
let remoteURL = (CFG.server_url || "").trim();
const nums  = ["threads","min_record_ms","max_record_seconds","translate_ask_seconds","server_port","paste_delay_ms","history_days","history_max"];
const sels  = ["ui_language","language","sound_theme","translate_target","translate_ask","hotkey_mode","overlay_position","theme","skin"];
const trAll = ["en","de","fr","es","it","pl","ru","uk"];
const L = {{L_JSON}};
const I_DL = '<svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2"><path d="M12 3v12"/><path d="M6 11l6 6 6-6"/><path d="M4 21h16"/></svg>';
const I_FIND = '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4"><circle cx="10.5" cy="10.5" r="6.5"/><line x1="15.5" y1="15.5" x2="21" y2="21"/></svg>';

let profiles = (CFG.profiles || []).map(p=>Object.assign({}, p));
let activeProfiles = (CFG.active_profiles || []).slice();
let translateDefault = !!CFG.translate_default;
let translateHotkey = CFG.translate_hotkey || "";
let pauseHotkey = CFG.pause_hotkey || "";
let expandedID = null;
let captureFor = null;

function esc(s){ const d=document.createElement("span"); d.textContent=s||""; return d.innerHTML; }
async function toggleMax(){
  const max = await appMaxRestore();
  paintMaxButton(max);
}
function paintMaxButton(max){
  const b = document.getElementById("cap_max");
  if(!b) return;
  b.innerHTML = max ? "&#10064;" : "&#9744;";
  b.title = max ? L.wndrestore : L.wndmax;
}
function initWindowButtons(){
  document.querySelectorAll(".rsz").forEach(strip=>{
    strip.addEventListener("mousedown", e=>{
      if(e.button === 0) appResizeEdge(Number(strip.dataset.edge));
    });
  });
  if(window.appMaximized) appMaximized().then(paintMaxButton);
}
const THEMES = {{THEME_LIST}};
function lookKey(skin, colour){
  return THEMES[skin + ":" + colour] ? skin + ":" + colour : skin;
}
function applyThemeVars(id){
  const p = THEMES[id] || THEMES["terminal:green"];
  if(!p) return;
  const root = document.documentElement.style;
  (p.vars || "").split(";").forEach(part=>{
    const at = part.indexOf(":");
    if(at > 0) root.setProperty(part.slice(0, at).trim(), part.slice(at + 1));
  });
  document.documentElement.dataset.skin = p.skin || id;
}
function paintLook(){
  const skin = document.getElementById("skin");
  const colour = document.getElementById("theme");
  const row = document.getElementById("colour_row");
  const box = document.getElementById("theme_swatches");
  if(!skin || !colour) return;
  const ownColours = !THEMES[skin.value + ":" + colour.value] && !THEMES[skin.value + ":green"];
  if(row) row.style.display = ownColours ? "none" : "";
  applyThemeVars(lookKey(skin.value, colour.value));
  if(box) box.querySelectorAll(".swatch").forEach(b=>b.classList.toggle("on", b.dataset.theme === colour.value));
}
function initTheme(){
  const skin = document.getElementById("skin");
  const colour = document.getElementById("theme");
  const box = document.getElementById("theme_swatches");
  if(!skin || !colour || !box) return;
  box.innerHTML = "";
  [...colour.options].forEach(o=>{
    const id = o.value;
    const look = THEMES["terminal:" + id];
    if(!look) return;
    const b = document.createElement("button");
    b.type = "button";
    b.className = "swatch" + (id === colour.value ? " on" : "");
    b.dataset.theme = id;
    b.title = o.textContent;
    const i = document.createElement("i");
    i.style.background = look.accent;
    b.appendChild(i);
    b.onclick = ()=>{ colour.value = id; colour.dispatchEvent(new Event("change", {bubbles:true})); };
    box.appendChild(b);
  });
  skin.addEventListener("change", paintLook);
  colour.addEventListener("change", paintLook);
  paintLook();
}
function initServerExe(){
  const box = document.getElementById("server_exe");
  const edit = document.getElementById("exe_edit");
  const reset = document.getElementById("exe_reset");
  if(!box || !edit || !reset) return;
  box.value = CFG.server_exe_path || CFG.server_exe || "";
  edit.onclick = async ()=>{
    if(exeUnlocked) return;
    if(!await askConfirm(L.exewarn, L.exeedit)) return;
    exeUnlocked = true;
    box.readOnly = false;
    box.focus();
    box.select();
  };
  reset.onclick = ()=>{
    exeStored = CFG.server_exe_default || "whisper-server.exe";
    exeUnlocked = false;
    box.readOnly = true;
    box.value = CFG.server_exe_path_default || CFG.server_exe_default || "";
    doSave();
    toast(L.upd, "ok");
  };
  box.onchange = ()=>{ if(exeUnlocked){ exeStored = box.value.trim(); doSave(); } };
}
function buildToc(){
  let card = null, best = 0;
  document.querySelectorAll("#p-about .card").forEach(c=>{
    const n = c.querySelectorAll(".wh").length;
    if(n > best){ best = n; card = c; }
  });
  if(!card || best < 3) return;
  const heads = [...card.querySelectorAll(".wh")];
  const toc = document.createElement("nav");
  toc.className = "toc";
  toc.setAttribute("aria-label", card.querySelector(".wh").textContent);
  heads.forEach((h, i)=>{
    if(!h.id) h.id = "wh" + i;
    const a = document.createElement("a");
    a.href = "#" + h.id;
    a.textContent = h.textContent;
    a.onclick = (e)=>{ e.preventDefault(); h.scrollIntoView({block:"start"}); h.classList.add("hit"); setTimeout(()=>h.classList.remove("hit"), 1500); };
    toc.appendChild(a);
  });
  card.insertBefore(toc, card.firstChild);
}
function labelPages(){
  document.querySelectorAll(".nav").forEach(b=>{
    const page = document.getElementById("p-" + b.dataset.p);
    const label = b.querySelector(".nlabel");
    if(!page || !label) return;
    page.setAttribute("aria-label", label.textContent.trim());
    b.setAttribute("aria-controls", page.id);
  });
}
function ariaFromTitle(root){
  const scope = root && root.querySelectorAll ? root : document;
  scope.querySelectorAll("button[title]:not([aria-label]),select[title]:not([aria-label])").forEach(b=>{
    const t = b.textContent.trim();
    if(!t || t.length <= 2 || b.classList.contains("iconbtn")) b.setAttribute("aria-label", b.getAttribute("title"));
  });
  scope.querySelectorAll("[title]").forEach(el=>{
    const text = el.getAttribute("title");
    if(!text){ el.removeAttribute("title"); return; }
    el.dataset.tip = text;
    el.removeAttribute("title");
  });
}
let tipTimer = null;
function hideTip(){
  clearTimeout(tipTimer);
  const tip = document.getElementById("tip");
  if(tip) tip.classList.remove("on");
}
function showTip(el){
  const tip = document.getElementById("tip");
  if(!tip) return;
  tip.textContent = el.dataset.tip;
  tip.classList.add("on");
  const box = el.getBoundingClientRect();
  const size = tip.getBoundingClientRect();
  let left = box.left + box.width / 2 - size.width / 2;
  left = Math.max(6, Math.min(left, window.innerWidth - size.width - 6));
  let top = box.bottom + 6;
  if(top + size.height > window.innerHeight - 6) top = box.top - size.height - 6;
  tip.style.left = Math.round(left) + "px";
  tip.style.top = Math.round(Math.max(6, top)) + "px";
}
function initTips(){
  document.addEventListener("mouseover", e=>{
    const el = e.target.closest ? e.target.closest("[data-tip]") : null;
    if(!el){ hideTip(); return; }
    clearTimeout(tipTimer);
    tipTimer = setTimeout(()=>showTip(el), 320);
  });
  document.addEventListener("mouseout", e=>{
    if(e.target.closest && e.target.closest("[data-tip]")) hideTip();
  });
  document.addEventListener("mousedown", hideTip);
  document.addEventListener("scroll", hideTip, true);
}
if(typeof MutationObserver === "function"){
  new MutationObserver(list=>{
    list.forEach(m=>m.addedNodes.forEach(n=>{ if(n.nodeType === 1) ariaFromTitle(n.parentNode || n); }));
  }).observe(document.documentElement, {childList:true, subtree:true});
}


let selLLM = null;
let hfRepos = [];
let hfOpenRepo = null;
let hfFiles = [];
let hfQuery = "";
function fitLabel(f, need){
  const gb = (need/1024).toFixed(1);
  const col = f==="ok" ? "var(--green)" : (f==="warn" ? "var(--amber)" : "var(--bad)");
  const tip = f==="ok" ? L.fitok : (f==="warn" ? L.fitwarn : L.fitbad);
  return '<span title="'+esc(tip)+'" style="color:'+col+';font-size:12px;white-space:nowrap">&#9679; &#8776;'+gb+' GB</span>';
}
async function refreshLLM(){
  const st = JSON.parse(await appLLM());
  const installed = st.installed || [];
  if(selLLM === null){
    const act = installed.find(m=>m.active);
    if(act) selLLM = act.file;
  }

  const body = document.getElementById("proc-models");
  body.innerHTML = "";
  if(!installed.length && !(st.downloads||[]).length){
    const empty = document.createElement("div");
    empty.style.cssText = "color:var(--dim);font-size:12px;padding:6px 2px";
    empty.textContent = L.nollm;
    body.appendChild(empty);
  }
  installed.forEach(m=>{
    const div = document.createElement("div");
    div.className = "mrow";
    const checked = selLLM === m.file ? " checked" : "";
    const right = '<button class="iconbtn danger" title="'+L.del+'" data-a="ldel" data-f="'+esc(m.file)+'">&#10005;</button>';
    div.innerHTML = '<input type="radio" name="llmmdl" value="'+esc(m.file)+'"'+checked+'>'+
      '<span class="mdesc" style="flex:1;color:var(--green)">'+esc(m.file)+'</span>'+
      '<span class="msize">'+m.size+' MB</span><span>'+right+'</span>';
    body.appendChild(div);
  });
  let busy = false;
  (st.downloads || []).forEach(d=>{
    const div = document.createElement("div");
    div.className = "mrow";
    if(d.pct >= 0){ busy = true;
      div.innerHTML = '<span class="mdesc" style="flex:1">'+esc(d.file)+'</span><span class="mpct">'+(d.pct>0?d.pct+"%":"…")+'</span>'+
        '<button class="iconbtn danger" title="'+L.dlcancel+'" data-a="lcancel" data-f="'+esc(d.file)+'">&#10005;</button>';
    } else {
      div.innerHTML = '<span class="mdesc" style="flex:1">'+esc(d.file)+'</span><span class="mpct">! '+esc(d.err)+'</span>';
    }
    body.appendChild(div);
  });
  body.querySelectorAll('input[name="llmmdl"]').forEach(r=>{
    r.onchange = async ()=>{ selLLM = r.value; await doSave(); refreshLLM(); };
  });
  body.querySelectorAll("button[data-a='lcancel']").forEach(b=>{
    b.onclick = async ()=>{
      await appModelCancel("llm-" + b.dataset.f);
      refreshLLM();
    };
  });
  body.querySelectorAll("button[data-a='ldel']").forEach(b=>{
    b.onclick = async ()=>{
      const f = b.dataset.f;
      if(!await askConfirm(L.confirmdel.replace("%s", f), L.del)) return;
      toast(await appLLMDel(f));
      if(selLLM === f){
        selLLM = null;
      }
      refreshLLM();
    };
  });

  const sbody = document.getElementById("proc-search");
  sbody.innerHTML = '<div class="row" style="padding-top:0">'+
    '<span style="position:relative;flex:1;display:flex;min-width:0">'+
    '<input type="text" id="hf_q" placeholder="'+L.hfph+'" style="flex:1;min-width:0;padding-left:34px;padding-right:30px">'+
    '<button type="button" id="hf_clr" title="'+L.cancel+'">&#10005;</button>'+
    '<button type="button" id="hf_go" title="'+L.hfph+'">'+I_FIND+'</button></span></div>'+
    '<div class="ramline">'+L.ram+' <b>'+((st.ram_free||st.ram)/1024).toFixed(1)+'</b> / '+(st.ram/1024).toFixed(0)+' GB '+L.free+
    ''+
      '<span class="dot" style="color:var(--green)">&#9679;</span>'+L.fitok+
      '<span class="dot" style="color:var(--amber)">&#9679;</span>'+L.fitwarn+
      '<span class="dot" style="color:var(--bad)">&#9679;</span>'+L.fitbad+'</div>'+
    '<div id="hf_results"></div>';
  const qEl = document.getElementById("hf_q");
  const clr = document.getElementById("hf_clr");
  const updClr = ()=>{ clr.style.display = qEl.value ? "block" : "none"; };
  qEl.value = hfQuery;
  updClr();
  qEl.oninput = ()=>{ hfQuery = qEl.value; updClr(); };
  clr.onclick = ()=>{ qEl.value = ""; hfQuery = ""; updClr(); qEl.focus(); };
  document.getElementById("hf_go").onclick = doHFSearch;
  qEl.onkeydown = e=>{ if(e.key === "Enter") doHFSearch(); };
  renderHF();

  renderProfiles(document.getElementById("profbody"), {state: installed.length ? "installed" : "absent"});
  if(busy) setTimeout(refreshLLM, 900);
}
async function doHFSearch(){
  const q = document.getElementById("hf_q").value;
  hfOpenRepo = null; hfFiles = [];
  const res = JSON.parse(await appLLMSearch(q));
  hfRepos = res.repos || [];
  if(res.error) toast(res.error);
  renderHF();
}
function renderHF(){
  const box = document.getElementById("hf_results");
  if(!box) return;
  box.innerHTML = "";
  hfRepos.forEach(r=>{
    const div = document.createElement("div");
    div.className = "mrow";
    div.innerHTML = '<button type="button" class="hfrepo" aria-expanded="'+(hfOpenRepo===r.id)+'">'+
      '<span class="mdesc">'+(hfOpenRepo===r.id?"▾ ":"▸ ")+esc(r.id)+'</span>'+
      '<span title="'+L.upd+'" class="hfupd">'+esc(r.updated||"")+'</span>'+
      '<span class="msize">↓'+(r.downloads>=1000000?(r.downloads/1000000).toFixed(1)+"M":(r.downloads/1000).toFixed(0)+"k")+'</span></button>'+
      '<button type="button" class="hflink" title="huggingface.co/'+esc(r.id)+'">↗</button>';
    div.querySelector(".hflink").onclick = e=>{ e.stopPropagation(); appHFPage(r.id); };
    div.querySelector(".hfrepo").onclick = async ()=>{
      if(hfOpenRepo === r.id){ hfOpenRepo = null; hfFiles = []; renderHF(); return; }
      hfOpenRepo = r.id; hfFiles = [];
      renderHF();
      const res = JSON.parse(await appLLMFiles(r.id));
      if(res.error){ toast(res.error); return; }
      hfFiles = res.files || [];
      renderHF();
    };
    box.appendChild(div);
    if(hfOpenRepo === r.id){
      hfFiles.forEach(f=>{
        const fd = document.createElement("div");
        fd.className = "mrow";
        fd.style.paddingLeft = "22px";
        fd.innerHTML = '<span class="mdesc" style="flex:1">'+esc(f.file)+'</span>'+
          '<span class="msize">'+(f.size>=1024?(f.size/1024).toFixed(1)+" GB":f.size+" MB")+'</span>'+
          '<span>'+fitLabel(f.fit, f.need)+'</span>'+
          '<button class="iconbtn" title="'+L.dl+'" data-repo="'+esc(r.id)+'" data-file="'+esc(f.file)+'" data-size="'+(f.size||0)+'">'+I_DL+'</button>';
        box.appendChild(fd);
      });
      box.querySelectorAll("button[data-repo]").forEach(b=>{
        b.onclick = async e=>{
          e.stopPropagation();
          await appLLMDlFile(b.dataset.repo, b.dataset.file, parseInt(b.dataset.size) || 0);
          refreshLLM();
        };
      });
    }
  });
}

function renderProfiles(body, st){
  body.innerHTML = "";
  const rows = st.state === "installed" ? profiles : [];
  rows.forEach(p=>{
    const div = document.createElement("div");
    div.className = "mrow";
    const checked = activeProfiles.includes(p.id) ? " checked" : "";
    let html = '<input type="checkbox" class="profcb" value="'+esc(p.id)+'"'+checked+'>';
    html += '<span class="mname" style="width:190px">'+esc(p.name)+'</span>';
    html += '<span class="mdesc">'+esc(p.hotkey||L.nohot)+'</span>';
    const on = expandedID === p.id ? ' style="color:var(--green);text-shadow:var(--glow)"' : '';
    html += '<span><button class="iconbtn" title="'+L.pedit+'" data-a="edit" data-id="'+esc(p.id)+'"'+on+'>&#9998;</button> ';
    html += '<button class="iconbtn danger" title="'+L.del+'" data-a="pdel" data-id="'+esc(p.id)+'">&#10005;</button></span>';
    div.innerHTML = html;
    body.appendChild(div);
    if(expandedID === p.id){
      body.appendChild(renderEditor(p));
    }
  });
  if(st.state !== "installed"){
    const note = document.createElement("div");
    note.style.cssText = "color:var(--dim);font-size:12px;padding:4px 2px";
    note.textContent = L.nollmp;
    body.appendChild(note);
  } else {
    const add = document.createElement("div");
    add.style.paddingTop = "8px";
    add.innerHTML = '<button class="mini" id="profadd">+ '+L.add+'</button>';
    body.appendChild(add);
    document.getElementById("profadd").onclick = ()=>{
      const id = "p" + Date.now();
      profiles.push({id:id, name:L.add, prompt:"", hotkey:""});
      expandedID = id;
      refreshLLM();
    };
  }
  body.querySelectorAll("input.profcb").forEach(r=>{
    r.onchange = ()=>{
      if(r.checked){
        if(!activeProfiles.includes(r.value)) activeProfiles.push(r.value);
      } else {
        activeProfiles = activeProfiles.filter(x=>x!==r.value);
      }
    };
  });
  body.querySelectorAll("button[data-a]").forEach(b=>{
    b.onclick = async ()=>{
      const id = b.dataset.id;
      if(b.dataset.a === "edit"){
        expandedID = expandedID === id ? null : id;
        refreshLLM();
        return;
      }
      const p = profiles.find(x=>x.id===id);
      if(!await askConfirm(L.confirmdel.replace("%s", p ? p.name : id), L.del)) return;
      profiles = profiles.filter(x=>x.id!==id);
      activeProfiles = activeProfiles.filter(x=>x!==id);
      if(expandedID === id) expandedID = null;
      refreshLLM();
      await doSave();
    };
  });
}

function renderEditor(p){
  const ed = document.createElement("div");
  ed.style.cssText = "border:1px solid var(--soft);padding:10px;margin:2px 0 8px";
  ed.innerHTML =
    '<div style="display:flex;justify-content:flex-end;margin:-4px -4px 2px 0"><button class="iconbtn" id="pf_close" title="'+L.pclose+'">&#9650;</button></div>'+
    '<div class="row" style="padding-top:0"><label>'+L.pname+'</label><input type="text" id="pf_name"></div>'+
    '<div class="row" style="align-items:flex-start"><label>'+L.pprompt+'</label></div>'+
    '<textarea id="pf_prompt" rows="4" style="width:100%;padding:7px 10px;border:1px solid var(--line);background:var(--field);color:var(--green);font:inherit;outline:none;resize:vertical"></textarea>'+
    '<div class="row"><label>'+L.phot+'</label><button type="button" class="hotkey-val" id="pf_hotkey" style="min-width:110px" title="'+L.pset+'"></button>'+
    '<button class="mini" id="pf_clear">'+L.pclr+'</button></div>'+
    '<div class="row"><label>'+L.ptest+'</label><input type="text" id="pf_sample" style="flex:1"><button class="iconbtn" id="pf_run">&#9654;</button></div>'+
    '<div id="pf_result" style="color:var(--amber);font-size:12px;min-height:16px;user-select:text"></div>';
  setTimeout(()=>{
    const name = document.getElementById("pf_name");
    const prompt = document.getElementById("pf_prompt");
    const hk = document.getElementById("pf_hotkey");
    name.value = p.name; prompt.value = p.prompt || "";
    hk.textContent = p.hotkey || L.nohot;
    name.oninput = ()=>{ p.name = name.value; };
    prompt.oninput = ()=>{ p.prompt = prompt.value; };
    document.getElementById("pf_close").onclick = ()=>{ expandedID = null; refreshLLM(); };
    document.getElementById("pf_hotkey").onclick = ()=>{ captureFor = p.id; appCaptureCombo(); };
    document.getElementById("pf_clear").onclick = ()=>{ p.hotkey=""; hk.textContent=L.nohot; doSave(); };
    document.getElementById("pf_run").onclick = ()=>{
      const s = document.getElementById("pf_sample").value;
      if(!s) return;
      document.getElementById("pf_result").textContent = "…";
      appLLMTest(p.prompt, s);
    };
  }, 0);
  return ed;
}

function updShow(latest, newer){
  const out = document.getElementById("upd_out");
  out.innerHTML = "";
  if(newer){
    out.appendChild(document.createTextNode(L.updavail.replace("%s", latest) + " "));
    const b = document.createElement("button");
    b.className = "mini";
    b.textContent = L.updgo;
    b.onclick = ()=>{ out.textContent = L.upddl; appDoUpdate(); };
    out.appendChild(b);
  } else {
    out.textContent = L.updnone;
  }
}
async function updCheck(){
  const out = document.getElementById("upd_out");
  out.textContent = "…";
  const r = JSON.parse(await appCheckUpdate());
  if(r.error){ out.textContent = L.upderr; return; }
  updShow(r.latest, r.newer);
}
function updProgress(p){ document.getElementById("upd_out").textContent = L.upddl + " " + p + "%"; }
function updError(e){ document.getElementById("upd_out").textContent = L.upderr + ": " + e; }

let micTimer = null;
async function refreshMics(){
  const sel = document.getElementById("mic_device");
  const mics = JSON.parse(await appMics());
  const chosen = sel.value || CFG.mic_device || "";
  sel.innerHTML = '<option value="">' + L.micdefault + '</option>';
  mics.forEach(m=>{
    const o = document.createElement("option");
    o.value = m.id;
    o.textContent = m.name;
    sel.appendChild(o);
  });
  sel.value = [...sel.options].some(o=>o.value===chosen) ? chosen : "";
  micChosen = sel.value;
}
const meterNow = {};
const MIC_FLOOR = 0.008;
const MIC_TOP = 0.35;
const MIC_DB_FLOOR = 20 * Math.log10(MIC_FLOOR);
const MIC_DB_TOP = 20 * Math.log10(MIC_TOP);
function meterPitch(box){
  const w = parseFloat(getComputedStyle(box).getPropertyValue("--lvlw")) || 4;
  return w + 2;
}
const METER_MIN_BARS = 7;
function micHeard(level){
  if(!(level > MIC_FLOOR)) return 0;
  const db = 20 * Math.log10(Math.min(1, level));
  const part = (db - MIC_DB_FLOOR) / (MIC_DB_TOP - MIC_DB_FLOOR);
  return Math.max(0, Math.min(1, part));
}
function meterFall(id, target){
  const prev = meterNow[id] || 0;
  const now = target > prev ? target : Math.max(target, prev - 0.12);
  meterNow[id] = now;
  return now;
}
function fitMeter(box){
  if(!box.classList.contains("grow")) return;
  const room = box.clientWidth;
  if(!room) return;
  const pitch = meterPitch(box);
  const want = Math.max(METER_MIN_BARS, Math.floor((room + 2) / pitch));
  if(want === box.children.length) return;
  box.innerHTML = new Array(want).fill("<i></i>").join("");
}
function paintMeter(box, level){
  fitMeter(box);
  const bars = box.querySelectorAll("i");
  if(!bars.length) return;
  const lit = Math.round(meterFall(box.id, micHeard(level)) * bars.length);
  bars.forEach((b, i)=>{ b.classList.toggle("on", i < lit); });
}
function startMeter(barId, pageId, hintId){
  return setInterval(async ()=>{
    const box = document.getElementById(barId);
    if(!box) return;
    const page = document.getElementById(pageId);
    if(!page || !page.classList.contains("active")){
      paintMeter(box, 0);
      return;
    }
    const lvl = await appMicLevel();
    paintMeter(box, lvl);
    if(hintId){
      const hint = document.getElementById(hintId);
      if(hint) hint.textContent = lvl > 0.02 ? "" : L.micquiet;
    }
  }, 120);
}
async function applyMic(){
  const micSel = document.getElementById("mic_device");
  const r = JSON.parse(await appMicSelect(micSel.value));
  const note = document.getElementById("mic_err");
  if(!r.ok){
    micSel.value = micChosen;
    if(note) note.textContent = r.message || "";
    toast(r.message, "error");
    return;
  }
  micChosen = micSel.value;
  if(note) note.textContent = "";
  doSave();
}
async function micCheck(){
  const btn = document.getElementById("mic_check");
  const out = document.getElementById("mic_verdict");
  if(!btn || !out) return;
  btn.disabled = true;
  out.className = "micverdict";
  out.textContent = L.micchecking;
  try {
    const r = JSON.parse(await appMicCheck());
    out.textContent = r.text;
    out.className = "micverdict " + (r.verdict === "ok" ? "ok" : "bad");
  } finally {
    btn.disabled = false;
  }
}
async function modelsCheck(){
  const btn = document.getElementById("mcheck");
  const out = document.getElementById("mcheck_out");
  if(!btn || !out) return;
  btn.disabled = true;
  out.className = "micverdict";
  out.textContent = L.mchecking;
  try {
    const r = JSON.parse(await appCheckModels());
    out.textContent = r.text;
    out.className = "micverdict " + (r.ok ? "ok" : "bad");
  } finally {
    btn.disabled = false;
  }
}
function startMicMeter(){
  if(micTimer) return;
  micTimer = startMeter("mic_bar", "p-mic", "mic_hint");
}

function updTrHotkey(){
  const el = document.getElementById("tr_hotkey");
  if(el) el.textContent = translateHotkey || L.nohot;
}
function updPauseHotkey(){
  const el = document.getElementById("pause_hotkey");
  if(el) el.textContent = pauseHotkey || L.nohot;
}
function comboCaptured(combo, warn){
  if(warn) toast(warn, "warn");
  if(!captureFor) return;
  if(captureFor === "__wt"){
    captureFor = null;
    if(combo) translateHotkey = combo;
    updTrHotkey();
    doSave();
    return;
  }
  if(captureFor === "__pause"){
    captureFor = null;
    if(combo) pauseHotkey = combo;
    updPauseHotkey();
    doSave();
    return;
  }
  const p = profiles.find(x=>x.id===captureFor);
  captureFor = null;
  if(p && combo){ p.hotkey = combo; }
  refreshLLM();
}

function llmTestResult(out){
  const el = document.getElementById("pf_result");
  if(el) el.textContent = out;
}

let micChosen = "";
let selModel = null;
let activeModelId = null;
let pendingDl = null;
async function refreshState(){
  const s = JSON.parse(await appState());
  const set = (id, v)=>{ const el = document.getElementById(id); if(el) el.textContent = v; };
  const setWithTip = (id, v)=>{
    const el = document.getElementById(id);
    if(!el) return;
    el.textContent = v;
    const card = el.closest(".scard") || el;
    if(v && v !== "—") card.dataset.tip = v; else delete card.dataset.tip;
  };
  const capLed = (ledId, btnId, state)=>{
    led(ledId, state === "ready" || state === "remote", state === "missing" || state === "downloading");
    const b = document.getElementById(btnId);
    if(b) b.textContent = state === "missing" ? L.get : L.change;
  };
  const led = (id, on, warn)=>{
    const el = document.getElementById(id);
    if(!el) return;
    el.classList.toggle("on", !!on);
    el.classList.toggle("warn", !!warn);
  };
  const berr = document.getElementById("state_backend");
  if(berr){
    berr.style.display = s.backend_err ? "" : "none";
    set("state_backend_text", s.backend_err || "");
  }
  set("state_hotkey", s.hotkey);
  setWithTip("state_mic", s.mic);
  setWithTip("state_ru", s.ru_model);
  setWithTip("state_other", s.other_model);
  setWithTip("state_llm", s.llm);
  set("state_ram", s.ram);
  set("state_last", s.last);
  const copyBtn = document.getElementById("state_copy");
  if(copyBtn) copyBtn.disabled = !s.last || s.last === "—";
  const lastEl = document.getElementById("state_last");
  if(lastEl) lastEl.title = s.last && s.last !== "—" ? s.last : "";
  set("state_last_meta", s.last_meta || "");
  const metaEl = document.getElementById("state_last_meta");
  if(metaEl) metaEl.title = s.last_meta || "";
  set("st_main", s.status_line || s.status);
  led("state_mic_led", s.mic_ok);
  capLed("state_ru_led", "state_ru_btn", s.ru_state);
  capLed("state_other_led", "state_other_btn", s.other_state);
  led("state_llm_led", s.llm_ok, !s.llm_ok);
  const llmBtn = document.getElementById("state_llm_btn");
  if(llmBtn) llmBtn.textContent = s.llm_ok ? L.change : L.get;
  led("st_led", s.ready);
  const badge = (id, v, full, cls)=>{
    const el = document.getElementById(id);
    if(!el) return;
    el.textContent = v || "";
    el.title = full || v || "";
    el.classList.toggle("miss", cls === "miss");
  };
  const missing = s.ru_state === "missing" || s.other_state === "missing";
  badge("badge_mic", s.badges && s.badges.mic, s.mic);
  badge("badge_models", s.badges && s.badges.models, missing ? L.badgemiss : L.badgemodels, missing ? "miss" : "");
  badge("badge_system", s.badges && s.badges.system, L.badgesystem);
  badge("badge_history", s.badges && s.badges.history, L.badgehist);
  const rem = document.getElementById("st_remote");
  if(rem) rem.textContent = s.remote ? L.remotebadge : "";
  const upd = document.getElementById("state_upd");
  if(upd){
    upd.style.display = s.upd_version ? "" : "none";
    const label = document.getElementById("state_upd_text");
    if(label) label.textContent = s.upd_version ? L.updfound.replace("%s", s.upd_version) : "";
  }
}
function initStateScreen(){
  const copy = document.getElementById("state_copy");
  if(copy) copy.onclick = async ()=>{ const r = JSON.parse(await appCopyLast()); toast(r.text, r.ok ? "ok" : "error"); };
  const retry = document.getElementById("state_retry");
  if(retry) retry.onclick = async ()=>{ retry.disabled = true; await appRetryBackend(); setTimeout(()=>{ retry.disabled = false; refreshState(); }, 1200); };
  document.querySelectorAll("[data-goto]").forEach(b=>{ b.onclick = ()=>show(b.dataset.goto); });
  setInterval(refreshState, 1500);
  startMeter("state_mic_bar", "p-state", null);
  refreshState();
}
let modelFilter = "all";
function modelPassesFilter(m){
  switch(modelFilter){
    case "ru":    return m.langs === "ru" || m.langs === "*";
    case "multi": return m.langs === "*";
    case "punct": return !!m.punct;
    case "fit":   return m.fit !== "bad";
    default:      return true;
  }
}
function initModelFilters(){
  document.querySelectorAll(".fchip").forEach(b=>{
    b.onclick = ()=>{
      document.querySelectorAll(".fchip").forEach(o=>o.classList.toggle("on", o === b));
      modelFilter = b.dataset.f;
      refreshModels();
    };
  });
  const open = document.getElementById("adv_open");
  const box = document.getElementById("advisor");
  if(open && box){
    open.onclick = ()=>{ box.style.display = box.style.display === "none" ? "flex" : "none"; };
  }
  const go = document.getElementById("adv_go");
  if(go){
    go.onclick = async ()=>{
      const r = JSON.parse(await appAdvise(
        document.getElementById("adv_lang").value,
        document.getElementById("adv_prio").value,
        document.getElementById("adv_tr").checked));
      renderAdvice(r);
    };
  }
}
function renderAdvice(r){
  const out = document.getElementById("adv_out");
  if(!out) return;
  out.innerHTML = "";
  const head = document.createElement("div");
  head.textContent = r.text + " (" + r.ram + ")";
  out.appendChild(head);
  const plan = r.plan || [];
  plan.forEach((p, i)=>{
    const row = document.createElement("div");
    row.className = "advrow";
    row.innerHTML = '<span class="advrole">'+(i === 0 ? L.advprimary : L.advcompanion)+'</span>'+
      '<span class="advname">'+esc(p.name)+'</span>'+
      '<span class="advstate'+(p.installed ? " ok" : "")+'">'+(p.installed ? L.advhave : (p.size ? p.size+" MB" : L.dlstart))+'</span>';
    out.appendChild(row);
  });
  const missing = plan.filter(p=>!p.installed);
  if(!plan.length) return;
  const apply = document.createElement("button");
  apply.type = "button";
  apply.className = "mini";
  apply.textContent = L.advapply;
  apply.onclick = async ()=>{
    if(missing.length && !await askConfirm(L.advask.replace("%s", missing.map(p=>p.name).join(", ")).replace("%s", r.need + " MB"), L.dlstart)) return;
    modelFilter = "all";
    document.querySelectorAll(".fchip").forEach(o=>o.classList.toggle("on", o.dataset.f === "all"));
    for(const p of missing) await appModelDl(p.id);
    if(missing.length) pendingDl = missing[missing.length-1].id;
    const first = plan[0];
    if(first.installed){
      selModel = first.id;
      await doSave();
    } else {
      selModel = first.id;
    }
    refreshModels();
  };
  out.appendChild(apply);
}
async function refreshRouting(){
  const host = document.getElementById("routing");
  if(!host) return;
  const rows = JSON.parse(await appRouting());
  host.innerHTML = "";
  rows.forEach(r=>{
    const div = document.createElement("div");
    div.className = "rrow";
    div.innerHTML = '<span class="rcond">'+r.cond+'</span><span class="rarr">&#8594;</span>'+
      '<span class="reng">'+r.engine+'</span><span class="rwhy">'+r.why+'</span>';
    host.appendChild(div);
  });
}
async function refreshModels(){
  const rows = JSON.parse(await appModels());
  refreshRouting();
  const activeRow = rows.find(m=>m.state==="active");
  activeModelId = activeRow ? activeRow.id : null;
  let applyDone = null;
  if(pendingDl){
    const row = rows.find(m=>m.id===pendingDl);
    if(!row){ pendingDl = null; }
    else if(row.state === "installed"){ applyDone = pendingDl; pendingDl = null; }
    else if(row.state === "active"){ pendingDl = null; toast(L.mdlready, "ok"); }
    else if(row.state === "absent" && row.err){ pendingDl = null; toast(row.err, "error"); }
  }
  const el = document.getElementById("models");
  el.innerHTML = "";
  let busy = false;
  const slotOf = m=>m.engine === "sherpa" ? "ru" : "other";
  const groups = [
    {slot: "ru", title: L.slotru, rows: rows.filter(m=>slotOf(m) === "ru")},
    {slot: "other", title: L.slotother, rows: rows.filter(m=>slotOf(m) !== "ru")},
  ];
  groups.forEach(g=>{
    if(!g.rows.length) return;
    const head = document.createElement("div");
    head.className = "mslot";
    head.dataset.slot = g.slot;
    head.textContent = g.title;
    el.appendChild(head);
  g.rows.forEach(m=>{
    const div = document.createElement("div");
    div.className = "mrow";
    div.dataset.slot = g.slot;
    const checked = (selModel === m.id || (!selModel && m.slot)) ? " checked" : "";
    const radio = '<input type="radio" name="mdl-'+g.slot+'" value="'+m.id+'"'+checked+(m.id==="custom"?" disabled":"")+'>';
    let right = "";
    if(m.state === "downloading"){ busy = true; right = '<span class="mpct">'+(m.pct>0?m.pct+"%":"…")+'</span><button class="iconbtn danger" title="'+L.dlcancel+'" data-a="cancel" data-id="'+m.id+'">&#10005;</button>'; }
    else if(m.state === "absent") right = '<button class="iconbtn" title="'+L.dl+'" data-a="dl" data-id="'+m.id+'">'+I_DL+'</button>';
    else if(m.state === "installed" || (m.state === "active" && m.id !== "custom")) right = '<button class="iconbtn danger" title="'+L.del+'" data-a="del" data-id="'+m.id+'" data-name="'+esc(m.name)+'">&#10005;</button>';
    const tag = m.engine === "sherpa" ? '<span class="mtag">RU</span>' : (m.langs === "*" ? '<span class="mtag">99</span>' : "");
    const ram = m.ram ? '<span class="mram '+(m.fit||"")+'">≈'+m.ram+' MB RAM</span>' : '<span class="mram"></span>';
    div.innerHTML = radio+'<span class="mname">'+m.name+tag+'</span><span class="mdesc">'+m.desc+'</span>'+ram+'<span class="msize">'+(m.size?m.size+" MB":"")+'</span><span>'+right+'</span>';
    if(!modelPassesFilter(m)) div.classList.add("hidden");
    el.appendChild(div);
  });
  });
  el.querySelectorAll(".mslot").forEach(h=>{
    const shown = [...el.querySelectorAll('.mrow[data-slot="'+h.dataset.slot+'"]')].some(r=>!r.classList.contains("hidden"));
    h.classList.toggle("hidden", !shown);
  });
  el.querySelectorAll('input[name^="mdl-"]').forEach(r=>{
    r.onchange = async ()=>{
      const id = r.value;
      const row = rows.find(m=>m.id===id);
      if(!row) return;
      if(row.state === "absent"){
        const size = row.size ? row.size + " MB" : "";
        if(!await askConfirm(L.dlask.replace("%s", row.name).replace("%s", size), L.dlstart)){
          refreshModels();
          return;
        }
        selModel = id;
        pendingDl = id;
        await appModelDl(id);
      } else {
        selModel = id;
        await doSave();
        selModel = null;
      }
      refreshModels();
    };
  });
  el.querySelectorAll("button[data-a]").forEach(b=>{
    b.onclick = async ()=>{
      if(b.dataset.a === "dl"){ await appModelDl(b.dataset.id); }
      else if(b.dataset.a === "cancel"){
        await appModelCancel(b.dataset.id);
        if(pendingDl === b.dataset.id) pendingDl = null;
        if(selModel === b.dataset.id) selModel = null;
      }
      else {
        const isActive = b.dataset.id === activeModelId;
        const mname = b.dataset.name || b.dataset.id;
        const ask = isActive ? L.delactive.replace("%s", mname) : L.confirmdel.replace("%s", mname);
        if(!await askConfirm(ask, L.del)) return;
        toast(await appModelDel(b.dataset.id, isActive));
        if(selModel === b.dataset.id) selModel = null;
      }
      refreshModels();
    };
  });
  if(applyDone){
    selModel = applyDone;
    await doSave();
    toast(L.mdlready, "ok");
    refreshState();
    return;
  }
  if(busy || pendingDl) setTimeout(refreshModels, 900);
}
function askConfirm(text, okText, cancelText){
  return new Promise(resolve=>{
    const bg = document.createElement("div");
    bg.className = "modal-bg";
    const box = document.createElement("div");
    box.className = "modal";
    const p = document.createElement("p");
    p.textContent = text;
    const row = document.createElement("div");
    row.className = "modal-btns";
    const yes = document.createElement("button");
    yes.type = "button";
    yes.className = "btn yes";
    yes.textContent = okText || L.ok;
    const no = document.createElement("button");
    no.type = "button";
    no.className = "btn ghost";
    no.textContent = cancelText || L.cancel;
    const from = document.activeElement;
    const done = v => {
      bg.remove();
      document.removeEventListener("keydown", onKey, true);
      document.querySelectorAll(".content, .snav, .header").forEach(el=>el.removeAttribute("inert"));
      if(from && from.focus) from.focus();
      resolve(v);
    };
    function onKey(e){
      if(e.key === "Escape"){ e.preventDefault(); done(false); return; }
      if(e.key === "Enter"){
        e.preventDefault();
        done(document.activeElement === yes);
        return;
      }
      if(e.key === "Tab"){
        e.preventDefault();
        (document.activeElement === no ? yes : no).focus();
      }
    }
    yes.onclick = ()=>done(true);
    no.onclick = ()=>done(false);
    bg.onclick = e => { if(e.target === bg) done(false); };
    document.addEventListener("keydown", onKey, true);
    row.appendChild(no);
    row.appendChild(yes);
    box.appendChild(p);
    box.appendChild(row);
    box.setAttribute("role", "dialog");
    box.setAttribute("aria-modal", "true");
    bg.appendChild(box);
    document.querySelectorAll(".content, .snav, .header").forEach(el=>el.setAttribute("inert", ""));
    document.body.appendChild(bg);
    no.focus();
  });
}
function toast(msg, severity){
  if(!msg) return;
  const t = document.getElementById("st_saved");
  if(!t) return;
  t.textContent = msg;
  t.className = "stsaved" + (severity === "error" ? " bad" : severity === "warn" ? " warn" : "");
  clearTimeout(toast._t);
  toast._t = setTimeout(()=>{ t.textContent = ""; t.className = "stsaved"; }, (severity === "error" || severity === "warn") ? 6000 : 2200);
}

const numSels = ["threads","min_record_ms","max_record_seconds","translate_ask_seconds","paste_delay_ms","history_days","history_max"];
function load(){
  document.getElementById("punctuation").value = CFG.punctuation || "model";
  document.getElementById("whisper_prompt").value = CFG.whisper_prompt || "";
  trAll.forEach(l=>{ document.getElementById("tl_"+l).checked = (CFG.translate_ask_langs||[]).includes(l); });
  const trd = document.getElementById("tr_default");
  trd.checked = translateDefault;
  trd.onchange = ()=>{ translateDefault = trd.checked; syncTrControls(); };
  document.getElementById("translate_ask").onchange = syncTrControls;
  document.getElementById("mic_device").onchange = applyMic;
  document.getElementById("mic_refresh").onclick = refreshMics;
  const micChk = document.getElementById("mic_check");
  if(micChk) micChk.onclick = micCheck;
  const mChk = document.getElementById("mcheck");
  if(mChk) mChk.onclick = modelsCheck;
  refreshMics();
  startMicMeter();
  if ((CFG.model || "").indexOf("turbo") >= 0) document.getElementById("tr_warn").style.display = "block";
  updTrHotkey();
  document.getElementById("tr_hotkey").onclick = ()=>{ captureFor = "__wt"; appCaptureCombo(); };
  document.getElementById("tr_clear").onclick = ()=>{ translateHotkey = ""; updTrHotkey(); doSave(); };
  const pset = document.getElementById("pause_hotkey");
  if(pset) pset.onclick = ()=>{ captureFor = "__pause"; appCaptureCombo(); };
  const pclr = document.getElementById("pause_clear");
  if(pclr) pclr.onclick = ()=>{ pauseHotkey = ""; updPauseHotkey(); doSave(); };
  updPauseHotkey();
  document.getElementById("hotkey").textContent = CFG.hotkey;
  document.getElementById("ver").textContent = CFG._version;
  document.getElementById("ver2").textContent = CFG._version;
  document.getElementById("threads").max = CFG._cpus || 16;
  bools.forEach(k=>document.getElementById(k).checked = !!CFG[k]);
  texts.forEach(k=>document.getElementById(k).value = CFG[k]||"");
  initServerExe();
  document.getElementById("server_url").value = remoteURL;
  nums.forEach(k=>document.getElementById(k).value = CFG[k]);
  sels.forEach(k=>{
    const el=document.getElementById(k), v=CFG[k]||"auto";
    if(![...el.options].some(o=>o.value===v)){const o=document.createElement("option");o.value=v;o.textContent=v;el.appendChild(o);}
    el.value=v;
  });
  initTheme();
  numSels.forEach(k=>{
    const el = document.getElementById(k), v = String(CFG[k]);
    if(![...el.options].some(o=>o.value===v)){
      const o = document.createElement("option"); o.value = v; o.textContent = v; el.insertBefore(o, el.firstChild);
    }
    el.value = v;
  });
  syncTrControls();
}
function syncRemoteWarn(){
  const el = document.getElementById("server_url");
  const note = document.getElementById("remote_warn");
  if(!el || !note) return;
  note.textContent = el.value.trim() ? L.remotewarn : "";
}
function initRemote(){
  const el = document.getElementById("server_url");
  if(!el) return;
  let known = el.value.trim();
  syncRemoteWarn();
  el.addEventListener("change", async e=>{
    e.stopPropagation();
    const url = el.value.trim();
    if(url === known){
      syncRemoteWarn();
      return;
    }
    if(url && !await askConfirm(L.remoteask.replace("%s", url))){
      el.value = known;
      syncRemoteWarn();
      return;
    }
    known = url;
    remoteURL = url;
    syncRemoteWarn();
    doSave();
  });
}
function syncTrControls(){
  const always = document.getElementById("tr_default").checked;
  const mode = document.getElementById("translate_ask").value;
  document.getElementById("translate_target").disabled = !always && mode !== "timeout";
  document.getElementById("translate_ask").disabled = always;
  document.getElementById("translate_ask_seconds").disabled = always || mode !== "timeout";
  trAll.forEach(l=>{ document.getElementById("tl_"+l).disabled = always || mode === "never"; });
}
function setHotkey(s, warn){
  CFG.hotkey=s;
  document.getElementById("hotkey").textContent=s;
  const w = document.getElementById("wiz_hot");
  if(w) w.textContent = s;
  if(warn) toast(warn, "warn");
}
async function doSave(){
  const micSel = document.getElementById("mic_device");
  const f={hotkey:CFG.hotkey, model_id:selModel||"",
    mic_device: micSel.value,
    punctuation: document.getElementById("punctuation").value,
    ui_level: uiLevel,
    mic_device_name: micSel.value ? micSel.options[micSel.selectedIndex].textContent : "",
    whisper_prompt: document.getElementById("whisper_prompt").value,
    translate_hotkey: translateHotkey,
    pause_hotkey: pauseHotkey,
    server_url: remoteURL,
    translate_ask_langs: trAll.filter(l=>document.getElementById("tl_"+l).checked),
    translate_default: translateDefault,
    active_profiles: activeProfiles,
    app_rules: rules,
    replacements: repls,
    commands: cmds,
    llm_model_file: selLLM||"",
    profiles: profiles};
  bools.forEach(k=>f[k]=document.getElementById(k).checked);
  texts.forEach(k=>f[k]=document.getElementById(k).value);
  f.server_exe = exeUnlocked ? document.getElementById("server_exe").value.trim() : exeStored;
  nums.forEach(k=>f[k]=parseInt(document.getElementById(k).value)||0);
  sels.forEach(k=>f[k]=document.getElementById(k).value);
  const langChanged = f.ui_language !== (CFG.ui_language || "auto");
  const r = JSON.parse(await appSave(JSON.stringify(f)));
  if(langChanged){ appReload(wizOn ? "wizard" : curTab); return; }
  toast(r.message, r.severity);
  refreshModels();
  refreshState();
}
let curTab = "state";
function show(p){
  curTab = p;
  document.querySelectorAll(".nav").forEach(b=>{
    const on = b.dataset.p === p;
    b.classList.toggle("active", on);
    b.setAttribute("aria-selected", on ? "true" : "false");
  });
  document.querySelectorAll(".page").forEach(el=>{
    const on = el.id === "p-" + p;
    el.classList.toggle("active", on);
    el.setAttribute("aria-hidden", on ? "false" : "true");
  });
  if(p==="state") refreshState();
  if(p==="history") refreshHistory();
  document.querySelector(".content").scrollTop = 0;
}
let uiLevel = CFG.ui_level || "simple";
const opened = {};
function advCount(page){ return page.querySelectorAll(".row[data-adv], .card[data-adv]").length; }
function bindLabels(){
  document.querySelectorAll(".row").forEach(row=>{
    const label = row.querySelector("label");
    if(!label || label.htmlFor) return;
    const field = row.querySelector('input[type=checkbox], input[type=text], input[type=number], select');
    if(field && field.id) label.htmlFor = field.id;
  });
}
function applyLevel(){
  document.querySelectorAll(".lvlb").forEach(b=>b.classList.toggle("on", b.dataset.l === uiLevel));
  document.querySelectorAll(".page").forEach(page=>{
    const show = uiLevel === "all" || opened[page.id];
    page.classList.toggle("advopen", uiLevel === "simple" && !!opened[page.id]);
    page.querySelectorAll("[data-adv]").forEach(el=>el.classList.toggle("hidden", !show));
    let btn = page.querySelector(".moreb");
    const n = advCount(page);
    if(uiLevel === "all" || n === 0){
      if(btn) btn.style.display = "none";
      return;
    }
    if(!btn){
      btn = document.createElement("button");
      btn.type = "button";
      btn.className = "moreb";
      btn.onclick = ()=>{ opened[page.id] = !opened[page.id]; applyLevel(); };
      page.appendChild(btn);
    }
    btn.style.display = "";
    btn.textContent = (opened[page.id] ? L.less : L.more).replace("%d", n);
  });
}
function setLevel(l){
  uiLevel = l;
  applyLevel();
  applyNow();
}
let hits = [];
let hitAt = -1;
let searchFrom = null;
function searchMatches(s){
  const items = [...document.querySelectorAll(".page .row, .page .mrow, .page .sect, .page .hint, .page .mslot, .page .wizh, .about p, .about li")];
  const seen = new Set();
  return items.filter(r=>{
    if(!r.textContent.toLowerCase().includes(s)) return false;
    if(seen.has(r)) return false;
    seen.add(r);
    return true;
  });
}
function showHit(i){
  document.querySelectorAll(".hit").forEach(r=>r.classList.remove("hit"));
  if(!hits.length) return;
  hitAt = (i + hits.length) % hits.length;
  const hit = hits[hitAt];
  const page = hit.closest(".page");
  show(page.id.slice(2));
  if(hit.hasAttribute("data-adv") && uiLevel === "simple"){
    opened[page.id] = true;
    applyLevel();
  }
  hit.classList.add("hit");
  if(hit.scrollIntoView) hit.scrollIntoView({block:"center"});
  updSearchCount();
}
function updSearchCount(){
  const el = document.getElementById("ocount");
  if(!el) return;
  const q = document.getElementById("omni");
  if(!q || q.value.trim().length < 2){ el.textContent = ""; el.classList.remove("none"); return; }
  el.classList.toggle("none", hits.length === 0);
  el.textContent = hits.length ? (hitAt + 1) + "/" + hits.length : L.nofound;
}
function searchSettings(q){
  const s = q.trim().toLowerCase();
  if(s.length < 2){
    document.querySelectorAll(".hit").forEach(r=>r.classList.remove("hit"));
    hits = [];
    hitAt = -1;
    updSearchCount();
    return;
  }
  hits = searchMatches(s);
  hitAt = -1;
  if(!hits.length){
    document.querySelectorAll(".hit").forEach(r=>r.classList.remove("hit"));
    updSearchCount();
    return;
  }
  showHit(0);
}
function searchStep(delta){
  if(hits.length) showHit(hitAt + delta);
}
const WIZ_N = 5;
let wizOn = false, wizStep = 0, wizBase = 0, wizPlan = [], wizDlIds = [];
function wizEl(id){ return document.getElementById(id); }
function wizShow(n){
  wizStep = Math.max(0, Math.min(WIZ_N - 1, n));
  document.querySelectorAll(".wizstep").forEach((s, i)=>s.classList.toggle("on", i === wizStep));
  const dots = wizEl("wizdots");
  dots.innerHTML = "";
  for(let i = 0; i < WIZ_N; i++){
    const d = document.createElement("i");
    if(i <= wizStep) d.classList.add("on");
    dots.appendChild(d);
  }
  wizEl("wiz_back").style.display = wizStep ? "" : "none";
  wizSyncNav();
  wizEl("wiz_skip").style.display = wizStep === WIZ_N - 1 ? "none" : "";
  wizEl("wiz_next").textContent = wizStep === WIZ_N - 1 ? L.wizfinish : L.wiznext;
  wizEl("wizbody").scrollTop = 0;
  if(wizStep === 1) wizAdvise();
  if(wizStep === 2) wizInput();
  if(wizStep === 3) wizTry();
  if(wizStep === 4) wizDone();
}
async function wizAdvise(){
  const v = wizEl("wiz_lang").value || "auto";
  const bucket = v === "ru" ? "ru" : (v === "en" ? "en" : "multi");
  const r = JSON.parse(await appAdvise(bucket, "balance", false));
  wizPlan = r.plan || [];
  const box = wizEl("wiz_plan");
  box.innerHTML = "";
  wizPlan.forEach((p, i)=>{
    const row = document.createElement("div");
    row.className = "advrow";
    row.innerHTML = '<span class="advrole">'+(i === 0 ? L.advprimary : L.advcompanion)+'</span>'+
      '<span class="advname">'+esc(p.name)+'<span class="advwhy">'+(i === 0 ? L.advrolemain : L.advrolesecond)+'</span></span>'+
      '<span class="advstate'+(p.installed ? " ok" : "")+'">'+(p.installed ? L.advhave : (p.size ? p.size+" MB" : ""))+'</span>';
    box.appendChild(row);
  });
  const missing = wizPlan.filter(p=>!p.installed);
  const btn = wizEl("wiz_dl");
  const out = wizEl("wiz_dlout");
  btn.style.display = missing.length ? "" : "none";
  btn.textContent = L.dl + (r.need ? " · " + r.need + " MB" : "");
  wizEl("wiz_dlrow").style.display = "none";
  out.textContent = missing.length ? "" : L.wizhave;
  out.classList.toggle("ok", !missing.length);
  wizDlIds = [];
  wizSyncNav();
}
async function wizApplyModel(){
  const ready = wizPlan.filter(p=>p.installed);
  if(!ready.length) return;
  for(const p of ready){
    selModel = p.id;
    await doSave();
  }
  refreshModels();
  wizSyncNav();
}
function wizSyncNav(){
  const next = wizEl("wiz_next");
  if(!next) return;
  const waiting = wizStep === 1 && (wizDlIds.length > 0 || !wizPlan.some(p=>p.installed));
  next.disabled = waiting;
  next.title = waiting ? L.wizneedmodel : "";
}
async function wizDownload(){
  const missing = wizPlan.filter(p=>!p.installed);
  if(!missing.length) return;
  wizEl("wiz_dl").style.display = "none";
  wizEl("wiz_dlrow").style.display = "";
  wizEl("wiz_dlout").textContent = "";
  wizDlIds = missing.map(p=>p.id);
  wizSyncNav();
  for(const p of missing) await appModelDl(p.id);
}
async function wizPollDl(){
  if(!wizDlIds.length) return;
  const rows = JSON.parse(await appModels());
  let sum = 0, done = 0, err = "";
  wizDlIds.forEach(id=>{
    const row = rows.find(m=>m.id === id);
    if(!row) return;
    if(row.state === "downloading"){ sum += row.pct; return; }
    if(row.state === "absent"){ if(row.err) err = row.err; return; }
    sum += 100;
    done++;
  });
  const pct = Math.round(sum / wizDlIds.length);
  wizEl("wiz_dlbar").style.width = pct + "%";
  wizEl("wiz_dlpct").textContent = pct + "%";
  if(err){
    wizDlIds = [];
    wizEl("wiz_dlrow").style.display = "none";
    wizEl("wiz_dl").style.display = "";
    const out = wizEl("wiz_dlout");
    out.textContent = err;
    out.classList.remove("ok");
    wizSyncNav();
    return;
  }
  if(done === wizDlIds.length){
    wizDlIds = [];
    wizPlan = wizPlan.map(p=>({...p, installed: true}));
    wizEl("wiz_dlrow").style.display = "none";
    const out = wizEl("wiz_dlout");
    out.textContent = L.wizhave;
    out.classList.add("ok");
    await wizApplyModel();
  }
}
function wizInput(){
  wizEl("wiz_hot").textContent = CFG.hotkey || L.nohot;
  const src = wizEl("mic_device"), dst = wizEl("wiz_mic");
  dst.innerHTML = src.innerHTML;
  dst.value = src.value;
}
async function wizTry(){
  const s = JSON.parse(await appState());
  wizBase = s.last_at || 0;
  wizEl("wiz_trytext").textContent = L.wiztry.replace("%s", CFG.hotkey || L.nohot);
  const out = wizEl("wiz_tryout");
  out.textContent = L.wizwait;
  out.classList.remove("ok");
  const ta = wizEl("wiz_try");
  ta.value = "";
  setTimeout(()=>ta.focus(), 60);
}
async function wizPollTry(){
  const s = JSON.parse(await appState());
  if(!s.last_at || s.last_at <= wizBase) return;
  const out = wizEl("wiz_tryout");
  out.textContent = L.wizheard + " " + s.last;
  out.classList.add("ok");
}
async function wizDone(){
  wizEl("wiz_auto").checked = true;
}
async function wizFinish(){
  const on = await appSetAutorun(wizEl("wiz_auto").checked);
  const box = wizEl("autorun");
  if(box) box.checked = on;
  await appWizardDone();
  wizClose();
}
async function wizSkip(){
  await appWizardDone();
  wizClose();
}
function wizClose(){
  wizOn = false;
  wizEl("wiz").classList.remove("on");
  show("state");
}
function initWizard(){
  wizEl("wiz_next").onclick = ()=>{ if(wizStep === WIZ_N - 1) wizFinish(); else wizShow(wizStep + 1); };
  wizEl("wiz_back").onclick = ()=>wizShow(wizStep - 1);
  wizEl("wiz_skip").onclick = wizSkip;
  wizEl("wiz_dl").onclick = wizDownload;
  wizEl("wiz_hot").onclick = ()=>appCapture();
  wizEl("wiz_ui").onchange = ()=>{
    wizEl("ui_language").value = wizEl("wiz_ui").value;
    doSave();
  };
  wizEl("wiz_lang").onchange = ()=>{
    wizEl("language").value = wizEl("wiz_lang").value;
    doSave();
    wizAdvise();
  };
  wizEl("wiz_mic").onchange = async ()=>{
    const src = wizEl("mic_device"), dst = wizEl("wiz_mic");
    src.value = dst.value;
    await applyMic();
    dst.value = src.value;
  };
  wizEl("wizhead").addEventListener("mousedown", e=>{
    if(e.target.closest("button, input, select")) return;
    if(e.button === 0) appDrag();
  });
  setInterval(async ()=>{
    if(!wizOn || wizStep !== 2) return;
    const lvl = await appMicLevel();
    paintMeter(wizEl("wiz_micbar"), lvl);
    wizEl("wiz_michint").textContent = lvl > 0.02 ? "" : L.micquiet;
  }, 120);
  setInterval(()=>{
    if(!wizOn) return;
    if(wizDlIds.length) wizPollDl();
    if(wizStep === 3) wizPollTry();
  }, 800);
}
function histWhen(ms){
  const d = new Date(ms);
  const pad = n=>String(n).padStart(2, "0");
  const today = new Date();
  const sameDay = d.getFullYear() === today.getFullYear() && d.getMonth() === today.getMonth() && d.getDate() === today.getDate();
  const time = pad(d.getHours()) + ":" + pad(d.getMinutes());
  return sameDay ? time : pad(d.getDate()) + "." + pad(d.getMonth() + 1) + " " + time;
}
async function refreshHistory(){
  const body = document.getElementById("histbody");
  if(!body) return;
  const q = (document.getElementById("hist_find") || {}).value || "";
  const items = JSON.parse(await appHistory(q));
  syncHistList(!items.length && !q);
  body.innerHTML = "";
  if(!items.length){
    const empty = document.createElement("div");
    empty.className = "histempty";
    empty.textContent = q ? L.nofound : L.histempty;
    body.appendChild(empty);
    return;
  }
  items.forEach(it=>{
    const row = document.createElement("div");
    row.className = "histrow";
    const meta = document.createElement("div");
    meta.className = "histmeta";
    meta.innerHTML = "<b>" + esc(histWhen(it.at)) + "</b>" + esc(it.app || "");
    row.appendChild(meta);
    const text = document.createElement("div");
    text.className = "histtext";
    text.textContent = it.text;
    row.appendChild(text);
    const copy = document.createElement("button");
    copy.type = "button";
    copy.className = "mini";
    copy.textContent = L.histcopy;
    copy.onclick = async ()=>{ const r = JSON.parse(await appHistoryCopy(it.at)); toast(r.text, r.ok ? "ok" : "error"); };
    row.appendChild(copy);
    const ins = document.createElement("button");
    ins.type = "button";
    ins.className = "mini ghost";
    ins.textContent = L.histinsert;
    ins.onclick = async ()=>{
      const r = JSON.parse(await appHistoryInsert(it.at));
      toast(r.text, r.ok ? "ok" : "error");
    };
    row.appendChild(ins);
    body.appendChild(row);
  });
}
function initHistory(){
  const find = document.getElementById("hist_find");
  if(!find) return;
  let t = null;
  find.addEventListener("input", ()=>{
    clearTimeout(t);
    t = setTimeout(refreshHistory, 200);
  });
  const clear = document.getElementById("hist_clear");
  if(clear){
    clear.onclick = async ()=>{
      if(!await askConfirm(L.histask, L.histclear)) return;
      await appHistoryClear();
      refreshHistory();
      refreshState();
    };
  }
  const on = document.getElementById("history");
  if(on) on.addEventListener("change", ()=>{ syncHistControls(); setTimeout(refreshHistory, 300); });
  syncHistControls();
}
function syncHistControls(){
  const on = document.getElementById("history");
  if(!on) return;
  ["history_days","history_max"].forEach(id=>{
    const el = document.getElementById(id);
    if(el) el.disabled = !on.checked;
  });
}
function syncHistList(empty){
  ["hist_clear","hist_find"].forEach(id=>{
    const el = document.getElementById(id);
    if(el && !(id === "hist_find" && el.value)) el.disabled = !!empty;
  });
}
let cmds = (CFG.commands || []).map(c=>({...c}));
function renderCmds(){
  const body = document.getElementById("cmdbody");
  if(!body) return;
  body.innerHTML = "";
  if(!cmds.length){
    const empty = document.createElement("div");
    empty.className = "ruleempty";
    empty.textContent = L.cmdempty;
    body.appendChild(empty);
  }
  cmds.forEach((c, i)=>{
    const row = document.createElement("div");
    row.className = "replrow";

    const phrase = document.createElement("input");
    phrase.type = "text";
    phrase.className = "cphrase";
    phrase.value = c.phrase || "";
    phrase.placeholder = L.cmdph;
    phrase.oninput = ()=>{ cmds[i].phrase = phrase.value; };
    phrase.onchange = ()=>{ cmds[i].phrase = phrase.value; applyNow(); replTest(); };
    row.appendChild(phrase);

    const arrow = document.createElement("span");
    arrow.className = "rarrow";
    arrow.textContent = "→";
    row.appendChild(arrow);

    const act = document.createElement("select");
    act.className = "caction";
    ruleOpts(act, [["newline", L.cmdnewline], ["paragraph", L.cmdparagraph], ["text", L.cmdtext], ["cancel", L.cmdcancel]], c.action || "newline");
    row.appendChild(act);

    const txt = document.createElement("input");
    txt.type = "text";
    txt.className = "ctext";
    txt.value = c.text || "";
    txt.placeholder = L.cmdtextph;
    txt.style.display = (c.action === "text") ? "" : "none";
    txt.oninput = ()=>{ cmds[i].text = txt.value; };
    txt.onchange = ()=>{ cmds[i].text = txt.value; applyNow(); replTest(); };
    row.appendChild(txt);

    act.onchange = ()=>{
      cmds[i].action = act.value;
      txt.style.display = (act.value === "text") ? "" : "none";
      applyNow();
      replTest();
    };

    const del = document.createElement("button");
    del.type = "button";
    del.className = "rdel";
    del.title = L.cmddel;
    del.textContent = "✕";
    del.onclick = ()=>{ cmds.splice(i, 1); renderCmds(); applyNow(); replTest(); };
    row.appendChild(del);

    body.appendChild(row);
  });
}
function initCmds(){
  const add = document.getElementById("cmd_add");
  if(!add) return;
  add.onclick = ()=>{
    cmds.push({id: "c" + Date.now(), phrase: "", action: "newline", text: ""});
    renderCmds();
    const rows = document.querySelectorAll("#cmdbody .cphrase");
    if(rows.length) rows[rows.length - 1].focus();
  };
  const preset = document.getElementById("cmd_preset");
  if(preset){
    preset.onclick = ()=>{
      const known = cmds.map(c=>(c.phrase || "").toLowerCase());
      [["newline", L.cmdpnewline], ["paragraph", L.cmdpparagraph], ["cancel", L.cmdpcancel]].forEach(([action, phrase], n)=>{
        if(!phrase || known.includes(phrase.toLowerCase())) return;
        cmds.push({id: "c" + Date.now() + n, phrase: phrase, action: action, text: ""});
      });
      renderCmds();
      applyNow();
    };
  }
  renderCmds();
}
let repls = (CFG.replacements || []).map(r=>({...r}));
function renderRepls(){
  const body = document.getElementById("replbody");
  if(!body) return;
  body.innerHTML = "";
  if(!repls.length){
    const empty = document.createElement("div");
    empty.className = "ruleempty";
    empty.textContent = L.replempty;
    body.appendChild(empty);
  }
  repls.forEach((r, i)=>{
    const row = document.createElement("div");
    row.className = "replrow";

    const from = document.createElement("input");
    from.type = "text";
    from.className = "rfrom";
    from.value = r.from || "";
    from.placeholder = L.replfromph;
    from.oninput = ()=>{ repls[i].from = from.value; };
    from.onchange = ()=>{ repls[i].from = from.value; applyNow(); replTest(); };
    row.appendChild(from);

    const arrow = document.createElement("span");
    arrow.className = "rarrow";
    arrow.textContent = "→";
    row.appendChild(arrow);

    const to = document.createElement("input");
    to.type = "text";
    to.className = "rto";
    to.value = r.to || "";
    to.placeholder = L.repltoph;
    to.oninput = ()=>{ repls[i].to = to.value; };
    to.onchange = ()=>{ repls[i].to = to.value; applyNow(); replTest(); };
    row.appendChild(to);

    const wholeLbl = document.createElement("label");
    const whole = document.createElement("input");
    whole.type = "checkbox";
    whole.className = "rwhole";
    whole.checked = r.whole !== false;
    whole.onchange = ()=>{ repls[i].whole = whole.checked; applyNow(); replTest(); };
    wholeLbl.appendChild(whole);
    wholeLbl.appendChild(document.createTextNode(L.replwhole));
    row.appendChild(wholeLbl);

    const caseLbl = document.createElement("label");
    const mcase = document.createElement("input");
    mcase.type = "checkbox";
    mcase.className = "rcase";
    mcase.checked = !!r.match_case;
    mcase.onchange = ()=>{ repls[i].match_case = mcase.checked; applyNow(); replTest(); };
    caseLbl.appendChild(mcase);
    caseLbl.appendChild(document.createTextNode(L.replcase));
    row.appendChild(caseLbl);

    const del = document.createElement("button");
    del.type = "button";
    del.className = "rdel";
    del.title = L.repldel;
    del.textContent = "✕";
    del.onclick = ()=>{ repls.splice(i, 1); renderRepls(); applyNow(); replTest(); };
    row.appendChild(del);

    body.appendChild(row);
  });
}
let replTimer = null;
async function replTest(){
  const input = document.getElementById("repl_test");
  const out = document.getElementById("repl_out");
  if(!input || !out) return;
  const text = input.value.trim();
  if(!text){ out.textContent = ""; return; }
  const r = JSON.parse(await appTestText(text));
  out.textContent = r.text.split("\n").join(" ⏎ ");
  out.classList.toggle("bad", !!r.cancelled);
}
function initRepls(){
  const add = document.getElementById("repl_add");
  if(!add) return;
  add.onclick = ()=>{
    repls.push({id: "x" + Date.now(), from: "", to: "", whole: true, match_case: false});
    renderRepls();
    const rows = document.querySelectorAll("#replbody .rfrom");
    if(rows.length) rows[rows.length - 1].focus();
  };
  const exp = document.getElementById("lists_export");
  if(exp) exp.onclick = async ()=>{
    const r = JSON.parse(await appListsExport(JSON.stringify({replacements: repls, commands: cmds})));
    if(r.cancelled) return;
    toast(r.text, r.ok ? "ok" : "error");
  };
  const imp = document.getElementById("lists_import");
  if(imp) imp.onclick = async ()=>{
    const r = JSON.parse(await appListsImport(JSON.stringify({replacements: repls, commands: cmds})));
    if(r.cancelled) return;
    if(r.ok){
      repls = r.replacements || repls;
      cmds = r.commands || cmds;
      renderRepls();
      renderCmds();
      doSave();
    }
    toast(r.text, r.ok ? "ok" : "error");
  };
  const test = document.getElementById("repl_test");
  if(test){
    test.addEventListener("input", ()=>{
      clearTimeout(replTimer);
      replTimer = setTimeout(replTest, 250);
    });
  }
  renderRepls();
}
let rules = (CFG.app_rules || []).map(r=>({...r}));
function ruleOpts(sel, items, value){
  sel.innerHTML = "";
  items.forEach(it=>{
    const o = document.createElement("option");
    o.value = it[0];
    o.textContent = it[1];
    sel.appendChild(o);
  });
  sel.value = [...sel.options].some(o=>o.value===value) ? value : items[0][0];
}
function ruleProfileValue(r){
  if(!r.use_profiles) return "";
  if(!r.profiles || !r.profiles.length) return "-";
  return r.profiles[0];
}
function renderRules(){
  const body = document.getElementById("rulesbody");
  if(!body) return;
  body.innerHTML = "";
  if(!rules.length){
    const empty = document.createElement("div");
    empty.className = "ruleempty";
    empty.textContent = L.ruleempty;
    body.appendChild(empty);
  }
  rules.forEach((r, i)=>{
    const row = document.createElement("div");
    row.className = "rulerow";

    const app = document.createElement("input");
    app.type = "text";
    app.className = "rmatch";
    app.value = r.match || "";
    app.placeholder = L.ruleph;
    app.oninput = ()=>{ rules[i].match = app.value; };
    app.onchange = ()=>{ rules[i].match = app.value; applyNow(); };
    row.appendChild(app);

    const paste = document.createElement("select");
    paste.className = "rpaste";
    ruleOpts(paste, [["", L.pasteinh], ["clipboard", L.ruleclip], ["type", L.ruletype]], r.paste_mode || "");
    paste.onchange = ()=>{ rules[i].paste_mode = paste.value; applyNow(); };
    row.appendChild(paste);

    const enter = document.createElement("select");
    enter.className = "renter";
    ruleOpts(enter, [["", L.enterinh], ["on", L.ruleenteron], ["off", L.ruleenteroff]], r.auto_enter || "");
    enter.onchange = ()=>{ rules[i].auto_enter = enter.value; applyNow(); };
    row.appendChild(enter);

    const delay = document.createElement("select");
    delay.className = "rdelay";
    ruleOpts(delay, [["0", L.delaynone], ["50", "50 ms"], ["100", "100 ms"], ["250", "250 ms"], ["500", "500 ms"], ["1000", "1000 ms"]], String(r.delay_ms || 0));
    delay.onchange = ()=>{ rules[i].delay_ms = parseInt(delay.value) || 0; applyNow(); };
    row.appendChild(delay);

    const prof = document.createElement("select");
    prof.className = "rprof";
    prof.title = L.ruleprompts;
    const items = [["", L.promptinh], ["-", L.rulenoprompt]];
    profiles.forEach(p=>items.push([p.id, p.name]));
    ruleOpts(prof, items, ruleProfileValue(r));
    prof.onchange = ()=>{
      const v = prof.value;
      rules[i].use_profiles = v !== "";
      rules[i].profiles = v === "" || v === "-" ? [] : [v];
      applyNow();
    };
    row.appendChild(prof);

    const del = document.createElement("button");
    del.type = "button";
    del.className = "rdel";
    del.title = L.ruledel;
    del.textContent = "✕";
    del.onclick = ()=>{ rules.splice(i, 1); renderRules(); applyNow(); };
    row.appendChild(del);

    body.appendChild(row);
  });
}
function initRules(){
  const add = document.getElementById("rule_add");
  if(!add) return;
  add.onclick = ()=>{
    rules.push({id: "r" + Date.now(), match: "", paste_mode: "", auto_enter: "", delay_ms: 0, use_profiles: false, profiles: []});
    renderRules();
    const last = document.querySelectorAll("#rulesbody .rmatch");
    if(last.length) last[last.length - 1].focus();
  };
  renderRules();
}
async function refreshLastApp(){
  const btn = document.getElementById("rule_last");
  if(!btn) return;
  const s = JSON.parse(await appState());
  const exe = s.last_app || "";
  btn.textContent = exe ? L.rulelast.replace("%s", exe) : "";
  btn.onclick = ()=>{
    if(!exe) return;
    const known = rules.some(r=>(r.match || "").toLowerCase().includes(exe.toLowerCase()));
    if(known) return;
    rules.push({id: "r" + Date.now(), match: exe, paste_mode: "", auto_enter: "", delay_ms: 0, use_profiles: false, profiles: []});
    renderRules();
    applyNow();
  };
}
async function initAutorun(){
  const box = document.getElementById("autorun");
  if(!box) return;
  box.checked = await appAutorun();
  box.onchange = async ()=>{
    box.checked = await appSetAutorun(box.checked);
    toast(L.upd, "ok");
  };
}
function wizStart(){
  wizOn = true;
  wizEl("wiz").classList.add("on");
  wizEl("wiz_ui").value = CFG.ui_language || "auto";
  wizEl("wiz_lang").value = CFG.language || "auto";
  if(!wizEl("wiz_lang").value) wizEl("wiz_lang").value = "auto";
  wizShow(0);
}
document.getElementById("upd_check").onclick = updCheck;
const logBtn = document.getElementById("log_open");
if(logBtn) logBtn.onclick = ()=>appOpenLog();
const cfgReloadBtn = document.getElementById("cfg_reload");
if(cfgReloadBtn) cfgReloadBtn.onclick = async ()=>{ await appReloadConfig(); toast(L.upd, "ok"); appReload(curTab); };
const cfgResetBtn = document.getElementById("cfg_reset");
if(cfgResetBtn) cfgResetBtn.onclick = async ()=>{
  if(!await askConfirm(L.resetask, L.resetbtn)) return;
  await appResetSettings();
  appReload(curTab);
};
let applyTimer = null;
function applyNow(){
  clearTimeout(applyTimer);
  applyTimer = setTimeout(()=>{ doSave(); }, 120);
}
document.querySelector(".content").addEventListener("change", e=>{
  if(e.target.closest("#advisor")) return;
  if(e.target.name === "mdl" || e.target.name === "llmmdl") return;
  if(e.target.id === "mic_device" || e.target.id === "autorun") return;
  applyNow();
});
(async()=>{ const s = JSON.parse(await appUpdateStatus()); if(s.latest && s.url) updShow(s.latest, true); })();
document.querySelectorAll(".nav").forEach(b=>b.onclick=()=>{
  const p = b.dataset.p;
  if(p !== curTab) show(p);
});
document.querySelector(".header").addEventListener("mousedown", e=>{
  if(e.target.closest("button, input, select, textarea, a, .omni, .lvlsw")) return;
  if(e.button===0) appDrag();
});
load();
(async ()=>{
  initModelFilters();
  initStateScreen();
  initRemote();
  initWizard();
  initAutorun();
  initRules();
  initRepls();
  initCmds();
  initHistory();
  refreshLastApp();
  bindLabels();
  applyLevel();
  ariaFromTitle(document);
  initTips();
  initWindowButtons();
  labelPages();
  buildToc();
  document.querySelectorAll(".lvlb").forEach(b=>b.onclick=()=>setLevel(b.dataset.l));
  const omni = document.getElementById("omni");
  if(omni){
    omni.addEventListener("input", ()=>searchSettings(omni.value));
    document.addEventListener("keydown", e=>{
      if(wizOn) return;
      if((e.ctrlKey || e.metaKey) && (e.code === "KeyK" || e.key.toLowerCase() === "k")){ e.preventDefault(); searchFrom = document.activeElement; omni.focus(); omni.select(); }
      if(e.key === "Escape" && document.activeElement === omni){ omni.value = ""; searchSettings(""); omni.blur(); if(searchFrom && searchFrom.focus) searchFrom.focus(); searchFrom = null; }
      if(e.key === "Enter" && document.activeElement === omni){ e.preventDefault(); searchStep(e.shiftKey ? -1 : 1); }
    });
  }
  await refreshModels();
  await refreshLLM();
  if(window.appReady) appReady();
})();
const tabAlias = {general:"state", rec:"models", proc:"text", server:"system", about:"about", state:"state", dictation:"dictation", history:"history", mic:"mic", models:"models", text:"text", translate:"translate", system:"system"};
show(tabAlias[CFG._tab] || "state");
if(CFG._wizard || CFG._tab === "wizard") wizStart();
setTimeout(()=>{ if(window.appReady) appReady(); }, 400);
</script>
</body></html>`

func jsonResult(v any) string {
	out, err := json.Marshal(v)
	if err != nil {
		return `{"ok":false,"severity":"error","message":"json"}`
	}
	return string(out)
}

const (
	settingsMinW     = 760
	settingsMinH     = 500
	settingsMaxW     = 1120
	settingsMaxH     = 940
	settingsDefaultW = 860
	settingsDefaultH = 620
)
