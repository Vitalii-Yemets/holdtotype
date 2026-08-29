package main

import (
	"context"
	"encoding/json"
	"fmt"
	"holdtotype/internal/commands"
	"holdtotype/internal/ovplace"
	"holdtotype/internal/preset"
	"holdtotype/internal/replace"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
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
		log.Printf("createWebView: attempt %d failed", attempt)
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
	OverlayMonitor   string `json:"overlay_monitor"`
	OverlayXY        map[string]ovplace.Frac `json:"overlay_custom"`
	OverlayText      bool   `json:"overlay_text"`
	TypeMode         bool              `json:"type_mode"`
	Language         string            `json:"language"`
	LangModels       map[string]string `json:"lang_models"`
	Threads          int               `json:"threads"`
	MinRecordMs      int    `json:"min_record_ms"`
	PasteDelayMs     int    `json:"paste_delay_ms"`
	HistoryOn        bool   `json:"history"`
	HistoryDays      int    `json:"history_days"`
	HistoryMax       int    `json:"history_max"`
	HistorySkip      string `json:"history_skip"`
	PostEnabled      *bool  `json:"post_enabled"`
	PostSource       string `json:"post_source"`
	PostAPIURL       string `json:"post_api_url"`
	PostAPIModel     string `json:"post_api_model"`
	PostAPITimeout   int    `json:"post_api_timeout_s"`
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
	TranslateTarget     string             `json:"translate_target"`
	TranslateAsk        string             `json:"translate_ask"`
	TranslateAskSeconds int                `json:"translate_ask_seconds"`
	TranslateAskLangs   []string           `json:"translate_ask_langs"`
	TranslateDefault    bool               `json:"translate_default"`
	ActiveProfiles      []string           `json:"active_profiles"`
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
		log.Printf("openSettings: the text can be returned to the window [%s]", windowTitle(fg))
	}
	if !settingsOpen.CompareAndSwap(false, true) {
		if hwnd := settingsHwnd.Load(); hwnd != 0 {
			log.Printf("openSettings: already open, bringing it to the front")
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
		log.Printf("openSettings: CoInitializeEx (attempt %d): %v", attempt, err)
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
		log.Printf("openSettings: creating WebView2")
		winW, winH := settingsDefaultW, settingsDefaultH
		if c := a.snapshot(); c.SettingsW >= settingsMinW && c.SettingsH >= settingsMinH && c.SettingsW <= settingsMaxW && c.SettingsH <= settingsMaxH {
			winW, winH = c.SettingsW, c.SettingsH
		}
		lastWndW, lastWndH = 0, 0
		unpark := parkNewWebViewWindow()
		w := createWebView(winW, winH)
		unpark()
		if w == nil {
			log.Printf("WebView2 is unavailable")
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
		_ = w.Bind("appUnloadEngines", func() {
			a.parkEngines()
		})
		_ = w.Bind("appSetPostKey", func(key string) string {
			enc, err := protectKey(strings.TrimSpace(key))
			if err != nil {
				log.Printf("encrypting the API key: %v", err)
				return jsonResult(saveResult{Severity: "error", Message: humanError(err)})
			}
			a.mu.Lock()
			c := *a.cfg
			c.PostAPIKey = enc
			a.cfg = &c
			a.mu.Unlock()
			if err := saveConfig("config.json", &c); err != nil {
				return jsonResult(saveResult{Severity: "error", Message: humanError(err)})
			}
			if enc == "" {
				log.Printf("post-processing API key deleted")
			} else {
				log.Printf("post-processing API key saved (encrypted with DPAPI)")
			}
			return jsonResult(saveResult{OK: true, Severity: "ok", Message: strS("S_SAVED")})
		})
		_ = w.Bind("appPostTest", func(rawURL, model, key, timeout string) string {
			c := *a.snapshot()
			c.PostAPIURL = strings.TrimSpace(rawURL)
			c.PostAPIModel = strings.TrimSpace(model)
			if n, err := strconv.Atoi(strings.TrimSpace(timeout)); err == nil && n > 0 {
				c.PostAPITimeout = n
			}
			if k := strings.TrimSpace(key); k != "" {
				enc, err := protectKey(k)
				if err != nil {
					return jsonResult(saveResult{Severity: "error", Message: humanError(err)})
				}
				c.PostAPIKey = enc
			}
			if c.PostAPIURL == "" {
				return jsonResult(saveResult{Severity: "error", Message: strS("S_API_NONE")})
			}
			if _, err := externalChat(context.Background(), &c, postTestPrompt, "ping"); err != nil {
				msg := errText(err)
				log.Printf("post-processing server test: %v", err)
				a.setPostErr("", msg)
				return jsonResult(saveResult{Severity: "error", Message: msg})
			}
			log.Printf("post-processing server test: the server answered")
			a.setPostErr("", "")
			return jsonResult(saveResult{OK: true, Severity: "ok", Message: strS("S_API_TEST_OK")})
		})
		_ = w.Bind("appPostKeySet", func() bool {
			return strings.TrimSpace(a.snapshot().PostAPIKey) != ""
		})
		_ = w.Bind("appOpenModelsFolder", func() {
			if err := os.MkdirAll("models", 0o755); err != nil {
				log.Printf("models folder: %v", err)
			}
			if abs, err := filepath.Abs("models"); err == nil {
				shellOpenURL(abs)
			}
		})
		_ = w.Bind("appModelLink", func(id string) {
			if m := findModel(id); m != nil && m.LinkURL != "" {
				shellOpenURL(m.LinkURL)
			}
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
		_ = w.Bind("appLLMUnload", func() {
			a.llmShutdown()
			log.Printf("editor model unloaded from memory by hand")
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
					log.Printf("warning: %s", warn)
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
				log.Printf("microphone choice: %v", err)
				return jsonResult(saveResult{Severity: "error", Message: humanError(err)})
			}
			return jsonResult(saveResult{OK: true, Severity: "ok"})
		})
		_ = w.Bind("appJSError", func(msg string) {
			log.Printf("settings page error: %s", msg)
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
			log.Printf("settings: a recognizer restart was requested")
			a.requestServerRestart()
		})
		_ = w.Bind("appOpenLog", func() {
			log.Printf("settings: opening the log")
			shellOpen(appid.LogFile)
		})
		_ = w.Bind("appReloadConfig", func() {
			log.Printf("settings: reloading config.json on the button")
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
				log.Printf("history: clearing: %v", err)
			}
		})
		_ = w.Bind("appHistoryCopy", func(at float64) string {
			for _, it := range histStore.Items() {
				if it.At == int64(at) {
					if err := setClipboardText(it.Text); err != nil {
						log.Printf("copying from the history: %v", err)
						return listsAnswer(listsReply{Text: trf("copy.fail", humanError(err))})
					}
					log.Printf("copied from the history: %d characters", len([]rune(it.Text)))
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
		_ = w.Bind("appAutorun", func() bool {
			return autorunEnabled()
		})
		_ = w.Bind("appSetAutorun", func(on bool) bool {
			if err := setAutorun(on); err != nil {
				log.Printf("start with Windows: %v", err)
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
				log.Printf("saving the config: %v", err)
			}
			log.Printf("the first-run wizard is done")
		})

		log.Printf("openSettings: WebView2 created, opening the page")
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
		log.Printf("openSettings: the window is closed")
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
					log.Printf("saving the config: %v", err)
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
	combos := []string{f.Hotkey, f.TranslateHotkey}
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
	if validOverlayMonitor(f.OverlayMonitor) {
		c.OverlayMonitor = f.OverlayMonitor
	}
	if f.OverlayXY != nil {
		c.OverlayXY = ovplace.CleanCustom(f.OverlayXY)
	}
	setOverlayPlacement(&c)
	if f.TypeMode {
		c.PasteMode = "type"
	} else {
		c.PasteMode = "clipboard"
	}
	before := primaryEngine(old)
	c.Language = f.Language
	if f.LangModels != nil {
		cleaned, dropped := preset.Clean(f.LangModels, func(id string) *preset.Model {
			return presetView(findModel(id))
		})
		if len(dropped) > 0 {
			log.Printf("presets: dropped assignments %s", strings.Join(dropped, ", "))
		}
		c.LangModels = cleaned
	}
	if c.LangModels == nil {
		c.LangModels = map[string]string{}
	}
	modelChanged := false
	if applyPreset(&c) {
		modelChanged = true
	}
	if primaryEngine(&c) != before {
		modelChanged = true
	}
	if c.Language != old.Language {
		if am := activeModel(&c); am != nil && am.TrLangs != "" {
			modelChanged = true
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
	if f.PostEnabled != nil {
		c.PostEnabled = *f.PostEnabled
	}
	if f.PostSource == "local" || f.PostSource == "api" {
		c.PostSource = f.PostSource
	}
	if validPostAPIURL(f.PostAPIURL) {
		c.PostAPIURL = strings.TrimSpace(f.PostAPIURL)
	}
	c.PostAPIModel = strings.TrimSpace(f.PostAPIModel)
	if f.PostAPITimeout >= 5 && f.PostAPITimeout <= 120 {
		c.PostAPITimeout = f.PostAPITimeout
	}
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
		log.Printf("saving the config: %v — the settings were left unchanged", err)
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
		log.Printf("editor model switched: %s", c.LLMModel)
	}

	if hook != nil {
		hook.SetCombos(buildCombos(&c))
	}
	if len(c.ActiveProfiles) > 0 && !postAPIOn(&c) && llmInstalled(&c) {
		go func() {
			if _, err := a.ensureLLM(); err != nil {
				log.Printf("LLM warm-up: %v", err)
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

	log.Printf("settings saved: hotkey=%s ui=%s model=%s server=%v", c.Hotkey, c.UILanguage, c.Model, serverChanged)
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
		"overlay_monitor":         cfg.OverlayMonitor,
		"overlay_custom":          cfg.OverlayXY,
		"_monitors":               listMonitors(),
		"overlay_text":            cfg.OverlayText,
		"type_mode":               cfg.PasteMode == "type",
		"language":                cfg.Language,
		"model":                   cfg.Model,
		"lang_models":             cfg.LangModels,
		"threads":                 cfg.Threads,
		"min_record_ms":           cfg.MinRecordMs,
		"paste_delay_ms":          cfg.PasteDelayMs,
		"history":                 cfg.HistoryOn,
		"history_days":            cfg.HistoryDays,
		"history_max":             cfg.HistoryMax,
		"history_skip":            cfg.HistorySkip,
		"post_enabled":            cfg.PostEnabled,
		"post_source":              cfg.PostSource,
		"post_api_url":            cfg.PostAPIURL,
		"post_api_model":          cfg.PostAPIModel,
		"post_api_timeout_s":      cfg.PostAPITimeout,
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
		"replacements":            cfg.Replacements,
		"commands":                cfg.Commands,
		"translate_hotkey":        cfg.TranslateHotkey,
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
		"pname": "S_PROF_NAME", "pprompt": "S_PROF_PROMPT", "ptest": "S_PROF_TEST",
		"noprompts": "S_NO_PROMPTS", "pdrag": "S_PROF_DRAG", "pnameph": "S_PROF_NAME_PH",
		"ptestph": "S_PROF_TEST_PH", "pfnew": "S_PF_NEW", "pfedit": "S_PF_EDIT",
		"fitok": "S_FIT_OK", "fitwarn": "S_FIT_WARN", "fitbad": "S_FIT_BAD",
		"ram": "S_RAM", "ramavail": "S_RAM_AVAIL", "ramof": "S_RAM_OF", "hfph": "S_HF_PH", "nollm": "S_NO_LLM",
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
		"acc": "S_LIB_ACC", "spd": "S_LIB_SPD",
		"recauto": "S_RECAUTO",
		"asauto": "S_AS_AUTO", "recchip": "S_REC_CHIP",
		"backauto": "S_BACK_AUTO", "langsc": "S_LANGS_COUNT", "langsq": "S_LANGS_UNKNOWN",
		"tren": "S_TR_EN", "translist": "S_TR_LIST", "dlgoing": "S_DL_GOING",
		"trunavail": "S_TR_UNAVAIL", "trlockmsg": "S_TR_LOCK", "trlockok": "S_TR_LOCK_OK", "tronedlg": "S_TR_ONE", "trmodeldlg": "S_TR_NOMODEL", "trconfirm": "S_TR_CONFIRM", "fmsave": "S_DICT_SAVE",
		"srcused": "S_SRC_USED", "hfgo": "S_HF_GO",
		"notforlang": "S_NOT_FOR_LANG", "notinstalled": "S_NOT_INSTALLED",
		"manualnote": "S_MANUAL_NOTE", "manuallink": "S_MANUAL_LINK",
		"unload": "S_UNLOAD_GO", "unloaded": "S_UNLOADED",
		"hffit": "S_HF_FIT", "hfhidden": "S_HF_HIDDEN",
		"ovmoncursor": "S_OVMON_CURSOR", "ovmonn": "S_OVMON_N",
		"postwarn": "S_POSTAPI_WARN", "postask": "S_POSTAPI_ASK",
		"postkeyset": "S_POSTAPI_KEY_SET", "postkeynone": "S_POSTAPI_KEY_NONE",
		"apisumurl": "S_API_SUM_URL", "apisummodel": "S_API_SUM_MODEL", "apisumkey": "S_API_SUM_KEY",
		"apisumtimeout": "S_API_SUM_TIMEOUT", "apisumstate": "S_API_SUM_STATE", "apinomodel": "S_API_NO_MODEL",
		"apinone": "S_API_NONE", "apisetup": "S_POSTAPI_SETUP", "apiedit": "S_API_EDIT",
		"apitest": "S_API_TEST", "apitestrun": "S_API_TEST_RUN",
		"postnomodel": "S_POST_NO_MODEL", "postnoapi": "S_POST_NO_API", "postbad": "S_POST_BAD", "postnoprompt": "S_POST_NO_PROMPT",
		"llmcatalog": "S_LLM_CATALOG", "llminstalled": "S_LLM_INSTALLED", "llmsummodel": "S_LLM_SUM_MODEL",
		"llmblock": "S_LLM_BLOCK", "llmnonehint": "S_LLM_NONE_HINT", "llminmem": "S_LLM_IN_MEM", "llmondisk": "S_LLM_ON_DISK",
		"llmeject": "S_LLM_EJECT", "llmfound": "S_LLM_FOUND", "llmnosearch": "S_LLM_NOSEARCH", "llmsearchhint": "S_LLM_SEARCH_HINT", "llmpickwait": "S_LLM_PICK_WAIT",
		"llmsumsize": "S_LLM_SUM_SIZE", "llmsumcount": "S_LLM_SUM_COUNT", "llmsumram": "S_LLM_SUM_RAM",
		"dlgclose": "S_DLG_CLOSE", "llmnopick": "S_LLM_NOPICK",
		"apikeydel": "S_API_KEY_DEL", "apidlg": "S_API_DLG",
		"postapiurl": "S_POSTAPI_URL", "postapimodel": "S_POSTAPI_MODEL",
		"postapikey": "S_POSTAPI_KEY", "postapitimeout": "S_POSTAPI_TIMEOUT",
		"remotewarn": "S_REMOTE_WARN", "remoteask": "S_REMOTE_ASK", "remotebadge": "S_REMOTE_BADGE",
		"ok": "S_OK", "cancel": "S_CANCEL", "dlask": "S_DL_ASK", "dlstart": "S_DL_START", "dlcancel": "S_DL_CANCEL", "nofound": "S_NOT_FOUND",
		"replempty": "S_REPL_EMPTY", "repldel": "S_REPL_DEL", "replwhole": "S_REPL_WHOLE",
		"cmdempty": "S_CMD_EMPTY", "cmddel": "S_CMD_DEL", "cmdph": "S_CMD_PH",
		"cmdnewline": "S_CMD_NEWLINE", "cmdparagraph": "S_CMD_PARAGRAPH", "cmdtext": "S_CMD_TEXT", "cmdcancel": "S_CMD_CANCEL",
		"cmdtextph":   "S_CMD_TEXT_PH",
		"cmdpnewline": "S_CMD_P_NEWLINE", "cmdpparagraph": "S_CMD_P_PARAGRAPH", "cmdpcancel": "S_CMD_P_CANCEL", "cmdpreset": "S_CMD_PRESET",
		"histempty": "S_HIST_EMPTY", "histcopy": "S_HIST_COPY", "histask": "S_HIST_ASK", "histclear": "S_HIST_CLEAR",
		"micchecking": "S_MIC_CHECKING", "mchecking": "S_MCHECK_RUN", "histinsert": "S_HIST_INSERT",
		"replcase": "S_REPL_CASE", "replfromph": "S_REPL_FROM_PH", "repltoph": "S_REPL_TO_PH",
		"repllang": "S_REPL_LANG", "repllangall": "S_REPL_LANG_ALL", "listnothing": "S_LIST_NOTHING", "replwholefull": "S_REPL_WHOLE_FULL",
		"replcasefull": "S_REPL_CASE_FULL", "cmdaction": "S_CMD_ACTION", "fmadd": "S_FM_ADD",
		"fmtrepladd": "S_FM_T_REPL_ADD", "fmtrepledit": "S_FM_T_REPL_EDIT", "fmtcmdadd": "S_FM_T_CMD_ADD",
		"fmtcmdedit": "S_FM_T_CMD_EDIT", "mtdel": "S_MT_DEL", "mtdelprompt": "S_MT_DEL_PROMPT",
		"mtdl": "S_MT_DL", "mttroff": "S_MT_TR_OFF", "mttrone": "S_MT_TR_ONE", "mttrlock": "S_MT_TR_LOCK",
		"mtremote": "S_MT_REMOTE", "mtpost": "S_MT_POST", "mthist": "S_MT_HIST", "mtreset": "S_MT_RESET",
		"mtexe": "S_MT_EXE", "fmtdictadd": "S_FM_T_DICT_ADD", "dictempty": "S_DICT_EMPTY",
		"dictnomodel": "S_DICT_NOMODEL", "dictaddph": "S_DICT_ADD_PH",
		"ovposscheme": "S_OVPOS_SCHEME_SUB", "ovposdrag": "S_OVPOS_DRAG_SUB",
		"tiprepllang": "S_TIP_REPL_LANG", "tipreplcase": "S_TIP_REPL_CASE", "tipreplwhole": "S_TIP_REPL_WHOLE",
		"tipcmdaction": "S_TIP_CMD_ACTION",
		"wiznext": "S_WIZ_NEXT", "wizfinish": "S_WIZ_FINISH", "wizwait": "S_WIZ_WAIT",
		"wizheard": "S_WIZ_HEARD", "wizhave": "S_WIZ_HAVE", "wiztry": "S_WIZ_TRY_TEXT",
		"updavail": "S_UPD_AVAIL", "updgo": "S_UPD_GO", "upderr": "S_UPD_ERR", "upddl": "S_UPD_DL",
	} {
		lMap[jsKey] = str(sKey)
	}
	lJSON, _ := json.Marshal(lMap)

	pairs := []string{"{{CFG}}", string(cfgJSON), "{{L_JSON}}", string(lJSON), "{{APP}}", appid.Name,
		"{{SKIN}}", cfg.Skin,
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
<html data-skin="{{SKIN}}"><head><meta charset="utf-8"><title>{{S_TITLE}}</title><style>{{FONT_FACE}}
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
::-webkit-scrollbar-corner{background:var(--bg)}
textarea::-webkit-resizer{background-color:var(--field);background-image:linear-gradient(135deg,transparent 0 42%,var(--dim) 42% 52%,transparent 52% 66%,var(--dim) 66% 76%,transparent 76%)}
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
.nbadge{flex:none;max-width:64px;font-size:9px;padding:1px 5px;border:1px solid var(--line);border-radius:calc(var(--r) * .4);color:var(--lbl,var(--dim));text-shadow:var(--lblglow,none);white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
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
.note.bad{font-size:11px;color:var(--bad);line-height:1.5;padding:2px 0 6px;overflow-wrap:anywhere}
.note.bad:empty{display:none}
.note.ok{font-size:11px;color:var(--green);line-height:1.5;padding:2px 0 6px}
.note.ok:empty{display:none}
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
.row.hidden{display:none}
.row.hit{background:var(--navon);box-shadow:inset 2px 0 0 var(--hi)}
.modal-bg{position:fixed;inset:0;background:var(--scrim);display:flex;align-items:center;justify-content:center;z-index:20}
.modal{background:var(--panel);border:1px solid var(--line);border-radius:var(--r);box-shadow:0 0 24px rgba(var(--rgb),.18),var(--shadow);padding:20px 22px;max-width:380px;display:flex;flex-direction:column;gap:16px}
.modal p{font-size:13px;line-height:1.55;color:var(--green)}
.modal-btns{display:flex;gap:10px;justify-content:flex-end}
.modal .btn{padding:7px 18px;border:1px solid var(--btnline);border-radius:calc(var(--r) * .5);background:var(--btnbg);color:var(--btnfg);font:inherit;font-size:12px;letter-spacing:var(--ls);text-transform:var(--caps);cursor:pointer;white-space:nowrap;flex:none}
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
.card .row:last-child,.row.noline,#ovcard .row{border-bottom:0}
#ovcard,#behcard{border-top:1px solid var(--soft);padding-top:9px}
#behcard,#p-dictation .card:has(+#behcard){margin-bottom:0}
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
input[type=text],input[type=number],input[type=password],select{padding:var(--fieldpad);border:1px solid var(--line);border-radius:calc(var(--r) * .55);background:var(--field);color:var(--green);font:inherit;font-size:var(--ctlfs);line-height:1.2;outline:none}
input:focus,select:focus,textarea:focus{border-color:var(--dim);box-shadow:var(--glow)}
input::placeholder{color:var(--dim)}
input:disabled,select:disabled{opacity:.35;cursor:default}
#trlangs label:has(input:disabled){opacity:.45}
input[type=text],input[type=password]{width:220px;max-width:100%}select{width:210px;max-width:100%}
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
.row .val{color:var(--dim);font-size:11.5px;min-width:44px;text-align:right;flex:0 1 auto;overflow-wrap:anywhere}
.row .lbl{flex:1 1 45%}
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
.replrow .rlang{flex:none;width:auto;min-width:72px}
.replrow:last-child{border-bottom:none}
.replrow input[type=text]{flex:1 1 160px;min-width:120px;width:auto}
.replrow .rarrow{flex:none;color:var(--dim)}
.replrow label{flex:none;display:flex;align-items:center;gap:6px;font-size:11px;color:var(--dim);white-space:nowrap}
.replrow label input[type=checkbox]{width:28px;height:15px}
.rulefoot{display:flex;align-items:center;gap:10px;margin-top:10px;flex-wrap:wrap}
.rulefoot.right{justify-content:flex-end}
.rulefoot.nowrap{flex-wrap:nowrap}
.srchbox .clearx{position:static;transform:none;flex:none}
.rulefoot .ghost{border-color:var(--line);color:var(--dim)}
.rulefoot .ghost:empty{display:none}
.ruleempty{color:var(--dim);font-size:12px;padding:6px 0}
.card>.hint,.srccard>.hint,.blkhead .hint{font-size:11.5px;color:var(--dim);margin-bottom:6px;line-height:1.5}
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
button:focus-visible{outline:1px solid var(--green);outline-offset:2px}
.mrow .mname{width:132px;font-weight:var(--wb);white-space:nowrap}
.mrow .mdesc{flex:1;color:var(--dim);font-size:12px}
.mtag{font-size:9px;border:1px solid var(--line);border-radius:calc(var(--r) * .35);color:var(--dim);padding:0 4px;margin-left:6px;vertical-align:middle;letter-spacing:.06em}
.mrow.hidden{display:none}
.mram{color:var(--dim);font-size:10.5px;width:74px;text-align:right}
.mram.warn{color:var(--amber)}
.mram.bad{color:var(--bad)}
.mrow .msize{color:var(--dim);font-size:12px;width:70px;text-align:right}
.badge{font-size:11px;letter-spacing:1px;padding:4px 10px;border:1px solid var(--line);border-radius:calc(var(--r) * .5);color:var(--lbl,var(--green));text-shadow:var(--lblglow,var(--glow));text-transform:var(--caps)}
button.mini{flex:none;padding:5px 12px;border:1px solid var(--btnline);border-radius:calc(var(--r) * .5);background:none;color:var(--dim);font:inherit;font-size:11.5px;cursor:pointer;letter-spacing:var(--ls);text-transform:var(--caps)}
button.mini::before{content:var(--btnbo);color:var(--faint)}
button.mini::after{content:var(--btnbc);color:var(--faint)}
button.btn::before{content:var(--btnbo);color:var(--faint)}
button.btn::after{content:var(--btnbc);color:var(--faint)}
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
.mbar{display:flex;align-items:center;gap:7px;font-size:9.5px;letter-spacing:.06em;color:var(--dim)}
.mbar .mbl{width:64px;text-align:right;text-transform:uppercase}
.mtrack{flex:1;height:4px;background:var(--soft);border-radius:var(--barr,0);overflow:hidden}
.mtrack i{display:block;height:100%;background:var(--hi);border-radius:var(--barr,0)}
.skiplist{display:flex;gap:6px;flex-wrap:wrap;padding:4px 0 8px}
.skipchip{display:inline-flex;align-items:center;gap:6px;border:1px solid var(--line);border-radius:calc(var(--r) * .5);padding:3px 8px;font-size:12px;color:var(--lbl,var(--green));text-shadow:var(--lblglow,none);background:var(--field)}
.blklbl{display:block;color:var(--green);font-weight:var(--wb);font-size:13px;margin:0 0 4px}
.blkhead{display:flex;align-items:center;gap:16px}
.blkhead>div{flex:1;min-width:0}
.blkhead .hint{margin-bottom:0}
.plist{position:relative}
.prow{display:flex;align-items:center;gap:12px;border-bottom:1px solid var(--soft);padding:9px 0;font-size:13px;background:var(--bg)}
.prow:last-child{border-bottom:0}
.prow .pnm{flex:1;font-weight:var(--wb);color:var(--green)}
.prow .grip{flex:none;width:18px;color:var(--faint);cursor:grab;font-size:13px;line-height:1;user-select:none;text-align:center;letter-spacing:-2px}
.prow:hover .grip{color:var(--dim)}
.prow .grip:active{cursor:grabbing}
.prow.dragging{opacity:.55}
.prow.ghost{position:fixed;left:0;top:0;z-index:20;pointer-events:none;border:1px solid var(--dim);background:var(--panel);box-shadow:0 0 16px rgba(0,0,0,.55);padding:9px 12px;opacity:.95}
.dropline{position:absolute;left:0;right:0;height:2px;background:var(--hi);box-shadow:var(--higlow);pointer-events:none}
.dropline::before,.dropline::after{content:"";position:absolute;top:-3px;width:2px;height:8px;background:var(--hi)}
.dropline::before{left:0}
.dropline::after{right:0}
.plabel{display:block;font-size:12px;color:var(--green);margin-bottom:6px}
.formmodal textarea.pfprompt{width:100%;min-height:190px;border:1px solid var(--line);background:var(--field);color:var(--green);font:inherit;font-size:12.5px;line-height:1.55;padding:9px 11px;outline:none;resize:vertical}
.tryrow{display:flex;gap:10px;align-items:center}
.tryrow .clearwrap{flex:1}
.pfout{border:1px solid var(--soft);background:var(--field);padding:9px 11px;font-size:12px;color:var(--dim);min-height:38px;margin-top:8px;user-select:text;-webkit-user-select:text;white-space:pre-wrap}
.row>label.blklbl{flex:1;margin:0}
#cmdbody,#replbody,#dictbody{display:flex;flex-wrap:wrap;gap:8px;align-items:center;margin-top:12px}
#cmdbody .ruleempty,#replbody .ruleempty,#dictbody .ruleempty{flex:1 0 100%;text-align:center}
.cmdchip{display:inline-flex;align-items:center;gap:7px;border:1px solid var(--line);border-radius:calc(var(--r) * .5);padding:3px 9px;font-size:12px;background:var(--field);color:var(--green)}
.cmdchip .cpmeta{color:var(--dim);font-size:12px;text-shadow:none;line-height:1.4}
.cmdchip .cbtn{appearance:none;background:none;border:0;color:var(--dim);cursor:pointer;font:inherit;font-size:12px;line-height:1;padding:0;display:inline-flex;align-items:center}
.cmdchip .cbtn:hover{color:var(--green)}
.cmdchip .cbtn.rdel:hover{color:var(--bad)}
.skipchip .chipx{appearance:none;background:none;border:0;color:var(--dim);cursor:pointer;font:inherit;font-size:11px;padding:0}
.arow{display:flex;align-items:center;gap:10px;padding:7px 4px;border-bottom:1px solid var(--soft);cursor:pointer}
.arow:last-child{border-bottom:0}
.arow .aname{flex:none;min-width:150px}
.arow.on .aname{color:var(--green)}
.arow .alangs{flex:1;color:var(--dim);font-size:12px}
.arow .amiss{color:var(--bad);font-size:12px}
.lrow{display:flex;align-items:center;gap:12px;padding:8px 12px;margin:7px 0;border:1px solid var(--line);border-radius:var(--r);background:var(--field);cursor:pointer}
.lrow:hover{border-color:var(--dim)}
.lrow .plang{flex:none;width:170px;font-weight:var(--wb)}
.lrow .lmodel{flex:1;min-width:0;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.lrow .larr{flex:none;color:var(--faint);font-size:10px}
.lrow.dim .plang,.lrow.dim .lmodel{color:var(--dim);font-weight:400}
.lrow.open{border-color:var(--dim);border-bottom:none;border-bottom-left-radius:0;border-bottom-right-radius:0;margin-bottom:0;background:var(--navon)}
.lpick{border:1px solid var(--dim);border-top:none;border-radius:0 0 var(--r) var(--r);background:var(--navon);padding:9px 11px;margin:0 0 7px}
.lpick>.hint{margin-bottom:5px}
.pcard{display:flex;flex-direction:column;gap:7px;border:1px solid var(--line);border-radius:var(--r);background:var(--field);padding:10px 13px;margin:7px 0}
.pcard.cur{border-color:var(--hi);box-shadow:var(--higlow)}
.pcard.pickable{cursor:pointer}
.pcard.pickable:hover{border-color:var(--dim)}
.pcard .ptop{display:flex;gap:12px;align-items:flex-start;flex-wrap:wrap}
.pcard .ptop input.psw{flex:none;margin-top:2px}
#p-models .card+.card,#p-text .card+.card{margin-top:18px;border-top:1px solid var(--soft);padding-top:18px}
.pcard .pname{flex:1 1 200px;min-width:0;font-weight:var(--wb);font-size:13.5px;display:flex;align-items:center;flex-wrap:wrap;gap:4px 8px}
.pinfo{flex:none;width:15px;height:15px;border:1px solid var(--line);border-radius:calc(var(--r) * .4);display:inline-flex;align-items:center;justify-content:center;font-size:10px;font-weight:400;font-style:normal;color:var(--green);cursor:help}
.pcard .mbars{flex:none;width:190px;display:flex;flex-direction:column;gap:5px;padding-top:2px}
.pcard .pdesc{color:var(--dim);font-size:12px}
.pcard .pmeta{display:flex;align-items:center;gap:10px;flex-wrap:wrap;border-top:1px solid var(--soft);padding-top:7px;color:var(--dim);font-size:11px}
.pcard .pmt{border:1px solid var(--line);border-radius:calc(var(--r) * .4);padding:1px 8px;font-size:10px;letter-spacing:.06em;white-space:nowrap;color:var(--lbl,var(--dim));text-shadow:var(--lblglow,none)}
.pcard .pram{white-space:nowrap}
.pcard .pact{margin-left:auto;display:flex;align-items:center;gap:8px;white-space:nowrap}
.pcard .psize{color:var(--dim);font-size:12px}
.pcard .pcur{color:var(--hi);text-shadow:var(--higlow);font-size:14px;padding:0 5px}
.ownm{padding:7px 0 16px;border-bottom:1px solid var(--soft);margin-bottom:6px}
.ownm>label{display:block}
.ownm>label .sub{display:block;margin-top:2px}
.ownm .tl{margin-top:13px}
.tls{position:relative;display:grid;grid-template-columns:28px minmax(0,1fr) auto;gap:0 14px;padding-bottom:20px}
.tls:last-child{padding-bottom:0}
.tls:not(:last-child)::before{content:"";position:absolute;top:28px;bottom:0;left:13px;width:1px;background:var(--soft)}
.tls>i{position:relative;z-index:1;width:28px;height:28px;border:1px solid var(--line);border-radius:50%;display:inline-flex;align-items:center;justify-content:center;font-style:normal;color:var(--dim);font-size:12px;background:var(--panel)}
html[data-skin="terminal"] .tls>i{border-radius:calc(var(--r) * .4)}
.tls .tlc{min-width:0;padding-top:4px}
.tls .tlc>b{display:block;font-weight:var(--wb);font-size:13px}
.tls .tlsub{display:block;color:var(--dim);font-size:11.5px;margin-top:2px}
.tls .tlsub.tlok{color:var(--green)}
.tls .tla{padding-top:3px}
@media (max-width:640px){.tls{grid-template-columns:28px minmax(0,1fr)}.tls .tla{grid-column:2;padding-top:7px}}
.ownm .fmts{display:flex;gap:10px;flex-wrap:wrap;margin-top:9px}
.ownm .fmtc{flex:1 1 250px;min-width:230px;border:1px solid var(--line);border-radius:var(--r);background:var(--field);padding:9px 12px}
.ownm .fmtc .fr{display:flex;align-items:baseline;gap:12px;flex-wrap:wrap;margin-bottom:7px}
.ownm .fmtc b{font-weight:var(--wb);font-size:13px}
.ownm .fmtc .cnt{margin-left:auto;color:var(--dim);font-size:11px;white-space:nowrap}
.ownm .chips{display:flex;gap:7px;flex-wrap:wrap}
.code{background:var(--field);border:1px solid var(--line);border-radius:calc(var(--r) * .35);padding:0 6px;font-size:10.5px;color:var(--lbl,var(--green));text-shadow:var(--lblglow,none);white-space:nowrap}
.pchip{flex:none;font-size:9px;letter-spacing:.1em;text-transform:uppercase;border:1px solid var(--line);border-radius:calc(var(--r) * .4);padding:1px 6px;color:var(--lbl,var(--dim));text-shadow:var(--lblglow,none);vertical-align:middle;font-weight:400}
.pchip.on{color:var(--amber);border-color:var(--amber);text-shadow:none}
.dlline{color:var(--amber);font-size:11.5px;margin:2px 0 8px;text-shadow:var(--amberglow)}
input[type=text],input[type=number],input[type=password],select,button.mini{height:var(--ctlh,30px)}
button.mini{display:inline-flex;align-items:center;justify-content:center}
html[data-skin="terminal"] .row>button.mini,html[data-skin="terminal"] .tls .tla button.mini{min-width:142px}
html[data-skin="terminal"] #repl_add,html[data-skin="terminal"] #cmd_add,html[data-skin="terminal"] #dict_add{min-width:176px}
.pairrow{display:flex;align-items:center;gap:10px;padding:7px 2px;border-bottom:1px solid var(--soft);font-size:12.5px}
.pairrow:last-child{border-bottom:0}
.pairrow .psum{min-width:0;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;color:var(--green)}
.pairrow .pmeta2{margin-left:auto;color:var(--dim);font-size:11px;white-space:nowrap;flex:none}
.modal .mtitle{color:var(--green);font-weight:var(--wb);font-size:13px}
.formmodal{min-width:min(430px,92vw)}
.formmodal.wide{min-width:min(560px,92vw);max-width:92vw}
.formmodal .fmrow{display:flex;align-items:center;gap:14px;margin:10px 0;min-height:30px}
.formmodal .fmrow label .sub{display:block;font-size:10.5px;color:var(--dim);margin-top:2px;letter-spacing:0;line-height:1.5}
.clearwrap{position:relative;flex:1;min-width:0;display:flex}
.clearwrap input[type=text],.clearwrap input[type=password]{flex:1;width:100%;min-width:0;padding-right:28px}
.clearx{position:absolute;right:8px;top:50%;transform:translateY(-50%);appearance:none;background:none;border:0;color:var(--dim);cursor:pointer;font:inherit;font-size:11px;padding:0}
.clearx:hover{color:var(--bad)}
.formmodal .fmrow>label{flex:1;font-size:12px}
.formmodal .fmrow input[type=text]{flex:1;width:auto;min-width:0}
.formmodal .fmrow select{flex:0 0 auto;width:auto;min-width:176px}
.formmodal .fmrow input[type=checkbox]{margin-left:auto}
.row.dimmed>label{color:var(--dim)}
.row.dimmed>label .sub{color:var(--faint)}
.row.dimmed>label .sub.warn{color:var(--amber)}
.srccard{border:1px solid var(--line);border-radius:var(--r);background:var(--field);padding:11px 14px;margin:9px 0}
#src_api,#src_local{min-height:172px;display:flex;flex-direction:column}
#src_api .sum,#src_local .sum{flex:1;display:flex;flex-direction:column;gap:6px;margin:8px 0 10px}
.srccard .sumrow{display:flex;gap:10px;font-size:12.5px}
.srccard .sumk{color:var(--dim);flex:0 0 118px}
.srccard .sumv{color:var(--green);word-break:break-all}
.srccard .sumv.off{color:var(--faint)}
.srccard .acts{display:flex;justify-content:flex-end;gap:8px}
.llmdl{display:flex;align-items:center;gap:10px;font-size:11.5px;color:var(--dim);margin-bottom:8px}
.llmdl .nm{flex:0 1 auto;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.llmdl .mtrack{flex:1;min-width:60px}
.modal.llmcat{min-width:min(660px,94vw);max-height:86vh;overflow:auto;gap:18px}
.modal.llmcat #proc-search{display:flex;flex-direction:column;gap:12px}
.modal.llmcat #proc-models{display:flex;flex-direction:column;margin-top:16px;border-top:1px solid var(--soft);padding-top:14px}
.modal.llmcat .mrow{padding:0 4px 0 10px;gap:10px;height:32px;flex:none;border-bottom:1px solid var(--soft)}
.modal.llmcat #proc-models .mrow:first-of-type{border-top:1px solid var(--soft)}
.modal.llmcat #hf_results .mrow{padding:0 4px 0 10px}
.mrow .mstate{color:var(--faint);font-size:11px;white-space:nowrap}
.mrow.cur .mstate{color:var(--dim)}
.mrow.cur .mdesc{text-shadow:var(--glow)}
input.llmpick:disabled{border-color:var(--soft);cursor:default}
input.llmpick:disabled::after{background:var(--soft)}
.filtrow{display:flex;align-items:center;justify-content:space-between;gap:16px;flex-wrap:wrap}
.fcount{color:var(--dim);font-size:11.5px}
.searchempty{color:var(--faint);font-size:12px;text-align:center;padding:26px 10px}
.modal.llmcat .mrow.cur{border-color:var(--hi)}
.catsect{color:var(--dim);font-size:11px;letter-spacing:.12em;text-transform:uppercase;margin-top:4px}
input.llmpick{appearance:none;-webkit-appearance:none;width:32px;height:17px;border:1px solid var(--dim);border-radius:calc(var(--r) * .8);position:relative;cursor:pointer;background:none;flex:none;padding:0;margin:0}
input.llmpick::after{content:"";position:absolute;top:2px;left:2px;width:11px;height:11px;border-radius:calc(var(--r) * .6);background:var(--dim);transition:.15s}
input.llmpick:checked::after{left:17px;background:var(--hi);box-shadow:var(--higlow)}
.srccard .srchead{display:flex;align-items:center;gap:10px;font-size:12.5px;font-weight:var(--wb);color:var(--green);margin:0 0 7px}
.srccard.on{border-color:var(--hi);box-shadow:var(--higlow)}
.srccard.idle{cursor:pointer}
.srccard.idle>*{opacity:.38;pointer-events:none}
.srccard.idle>.acts{pointer-events:auto;opacity:.62}
.srccard.idle>.srchead{opacity:.62}
.srccard.idle>.srchead .srcpick{pointer-events:auto}
input.srcpick{appearance:none;-webkit-appearance:none;width:32px;height:17px;border:1px solid var(--dim);border-radius:calc(var(--r) * .8);position:relative;cursor:pointer;background:none;flex:none;padding:0;margin:0}
input.srcpick::after{content:"";position:absolute;top:2px;left:2px;width:11px;height:11px;border-radius:calc(var(--r) * .6);background:var(--dim);transition:.15s}
input.srcpick:checked::after{left:17px;background:var(--hi);box-shadow:var(--higlow)}
#p-post .card+.card{margin-top:18px;border-top:1px solid var(--soft);padding-top:18px}
.srccard.idle:hover{border-color:var(--dim)}
.card.offdim>:not(h2):not(.blklbl){opacity:.38;pointer-events:none}
.srchrow{display:flex;align-items:center;gap:8px}
.srchrow>button.mini{flex:none;height:var(--ctlh,30px)}
.fitleg{margin-top:0}
.srchbox{display:flex;align-items:center;gap:8px;flex:1;min-width:0;border:1px solid var(--line);border-radius:calc(var(--r) * .55);background:var(--field);height:var(--ctlh,30px);padding:0 5px 0 10px}
.srchbox:focus-within{border-color:var(--dim)}
.srchbox svg{color:var(--faint);flex:none}
.srchbox input[type=text]{flex:1;min-width:0;width:auto;height:auto;border:0;background:none;box-shadow:none;outline:none;padding:0}
.srchbox input[type=text]:focus{border:0;box-shadow:none;outline:none}
.srchbox .iconbtn{flex:none}
.srchbox button.mini{height:22px;font-size:10.5px;padding:0 8px;flex:none}
input::placeholder{color:var(--faint)}
.ovscheme{position:relative;flex:none;width:176px;touch-action:none}
.ovscheme.off{opacity:.35;pointer-events:none}
.ovcase{position:relative;padding:9px 10px 14px;background:linear-gradient(180deg,var(--soft) 0%,var(--panel) 100%);border:1px solid var(--line);border-radius:11px}
.ovcrt{position:relative;height:80px;background:var(--field);border-radius:16px/13px;overflow:hidden;box-shadow:inset 0 0 16px rgba(var(--rgb),.12),inset 0 0 4px rgba(var(--rgb),.30)}
.ovcrt::after{content:"";position:absolute;inset:0;pointer-events:none;background:radial-gradient(120% 90% at 28% 12%,rgba(var(--rgb),.10) 0%,rgba(0,0,0,0) 55%)}
.ovscan{position:absolute;inset:0;pointer-events:none;opacity:.55;z-index:1;background:repeating-linear-gradient(transparent 0 2px,rgba(0,0,0,.40) 2px 3px)}
.ovgrid{position:absolute;inset:0;z-index:1;pointer-events:none;opacity:0;transition:.15s;background:
  linear-gradient(90deg,transparent 33.2%,rgba(var(--rgb),.16) 33.2% 33.5%,transparent 33.5%),
  linear-gradient(90deg,transparent 66.4%,rgba(var(--rgb),.16) 66.4% 66.7%,transparent 66.7%),
  linear-gradient(180deg,transparent 33.2%,rgba(var(--rgb),.16) 33.2% 33.8%,transparent 33.8%),
  linear-gradient(180deg,transparent 66.4%,rgba(var(--rgb),.16) 66.4% 67%,transparent 67%)}
.ovcrt:not(.free):hover .ovgrid{opacity:1}
.ovzone{position:absolute;appearance:none;border:0;background:none;padding:0;cursor:pointer;z-index:2;width:33.333%;height:33.333%}
.ovzone:hover{background:rgba(var(--rgb),.10)}
.ovzone:focus-visible{outline:1px solid var(--green);outline-offset:-2px}
.ovcrt.free .ovzone{display:none}
.ovcrt.free{cursor:crosshair}
.ovled{position:absolute;right:11px;bottom:5px;width:3px;height:3px;border-radius:50%;background:var(--green);box-shadow:var(--higlow)}
.ovknob{position:absolute;right:21px;bottom:4px;width:5px;height:5px;border-radius:50%;border:1px solid var(--line)}
.ovvents{position:absolute;left:11px;bottom:5px;width:30px;height:3px;background:repeating-linear-gradient(90deg,var(--line) 0 1px,transparent 1px 3px)}
.ovneck{width:36px;height:10px;margin:0 auto;background:linear-gradient(180deg,var(--soft),var(--panel));border-left:1px solid var(--line);border-right:1px solid var(--line)}
.ovbase{width:76px;height:5px;margin:0 auto;background:var(--soft);border:1px solid var(--line);border-radius:0 0 3px 3px}
.ovmini{position:absolute;width:26px;height:7px;margin:-3.5px 0 0 -13px;border-radius:99px;background:var(--hi);box-shadow:var(--higlow);z-index:3;display:flex;align-items:center}
.ovcrt.free .ovmini{cursor:grab}
.ovcrt.free .ovmini:active{cursor:grabbing}
.ovmini::before{content:"";width:2px;height:2px;border-radius:50%;background:var(--bg);margin-left:3px}
.skipchip .chipx:hover{color:var(--bad)}
#hf_results{max-height:44vh;overflow-y:auto;overscroll-behavior:contain}
.miclevel{flex:none;display:flex;align-items:flex-end;gap:2px;height:10px;width:auto}
.miclevel i{display:block;width:var(--lvlw,4px);height:10px;background:var(--soft);border:1px solid var(--line);box-sizing:border-box;border-radius:var(--lvlr,0);transition:background .1s linear}
.miclevel i.on{background:var(--hi);border-color:var(--hi);box-shadow:var(--higlow)}
.miclevel.grow{width:100%;overflow:hidden}
.lvlrow .miclevel.grow{flex:0 0 auto;width:min(320px,100%);min-width:0}
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
 <span class="ngrp">{{S_GRP_GENERAL}}</span>
 <button class="nav" role="tab" aria-selected="false" data-p="state"><span class="nlabel">{{S_NAV_STATE}}</span></button>
 <button class="nav" role="tab" aria-selected="false" data-p="system"><span class="nlabel">{{S_NAV_SYSTEM}}</span><span class="nbadge warn" id="badge_system"></span></button>
 <button class="nav" role="tab" aria-selected="false" data-p="mic"><span class="nlabel">{{S_NAV_MIC}}</span><span class="nbadge" id="badge_mic"></span></button>
 <button class="nav" role="tab" aria-selected="false" data-p="history"><span class="nlabel">{{S_NAV_HISTORY}}</span><span class="nbadge" id="badge_history"></span></button>
 <span class="ngrp">{{S_GRP_SPEECH}}</span>
 <button class="nav" role="tab" aria-selected="false" data-p="dictation"><span class="nlabel">{{S_NAV_DICT}}</span></button>
 <button class="nav" role="tab" aria-selected="false" data-p="models"><span class="nlabel">{{S_NAV_MODELS}}</span><span class="nbadge" id="badge_models"></span></button>
 <button class="nav" role="tab" aria-selected="false" data-p="text"><span class="nlabel">{{S_NAV_TEXT}}</span></button>
 <button class="nav" role="tab" aria-selected="false" data-p="post"><span class="nlabel">{{S_NAV_POST}}</span></button>
 <span class="ngrp">{{S_GRP_INFO}}</span>
 <button class="nav" role="tab" aria-selected="false" data-p="help"><span class="nlabel">{{S_NAV_HELP}}</span></button>
 <button class="nav" role="tab" aria-selected="false" data-p="about"><span class="nlabel">{{S_NAV_ABOUT}}</span></button>
 <button class="nav" role="tab" aria-selected="false" data-p="contacts"><span class="nlabel">{{S_NAV_CONTACTS}}</span></button>
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
  <div class="scard"><span class="k">{{S_STATE_ACTIVE}} · <span id="state_active_lang"></span></span>
   <span class="v"><i class="led" id="state_active_led"></i><span id="state_active">—</span></span>
   <button class="mini" id="state_active_btn" data-goto="models">{{S_CHANGE_MODEL}}</button></div>
  <div class="scard"><span class="k">{{S_STATE_PROC}}</span>
   <span class="v"><i class="led" id="state_llm_led"></i><span id="state_llm">—</span></span>
   <button class="mini" id="state_llm_btn" data-goto="post">{{S_PICK_MODEL}}</button></div>
 </div>
 <h2 class="sect">{{S_STATE_USED}}</h2>
 <div id="state_assigned"></div>
 <div class="row"><span class="lbl">{{S_STATE_INST}}<span class="sub">{{S_STATE_INST_SUB}}</span></span>
  <span class="val" id="state_installed">—</span>
  <button class="mini" data-goto="models">{{S_CHANGE_MODEL}}</button></div>
 <h2 class="sect">{{S_STATE_LAST}}</h2>
 <div class="row"><span class="lbl"><span class="lastres" id="state_last">—</span><span class="sub" id="state_last_meta"></span></span>
  <button class="mini" id="state_copy">{{S_STATE_COPY}}</button></div>
 <div class="row"><span class="lbl">{{S_STATE_MEM}}<span class="sub">{{S_STATE_MEM_SUB}}</span></span>
  <span class="val" id="state_ram">—</span></div>
 <div class="row"><span class="lbl">{{S_STATE_LOADED}}<span class="sub">{{S_STATE_LOADED_SUB}}</span></span>
  <span class="val" id="state_loaded">—</span></div>
 <div class="row" id="state_week_row" style="display:none"><span class="lbl">{{S_STATE_WEEK}}</span>
  <span class="val" id="state_week">—</span></div>
</div>

<div class="page" role="tabpanel" aria-hidden="true" id="p-history">
 <div class="card">
  <div class="row"><label>{{S_HIST_ON}}<span class="sub">{{S_HIST_ON_SUB}}</span></label><input type="checkbox" id="history"></div>
  <div class="row" data-adv><label>{{S_HIST_DAYS}}</label>
   <select id="history_days"><option value="1">1</option><option value="3">3</option><option value="7">7</option><option value="30">30</option></select></div>
  <div class="row" data-adv><label>{{S_HIST_MAX}}</label>
   <select id="history_max"><option value="50">50</option><option value="100">100</option><option value="200">200</option><option value="500">500</option></select></div>
  <div class="row" id="hist_skip_row"><label>{{S_HIST_SKIP}}<span class="sub">{{S_HIST_SKIP_SUB}}</span></label>
   <input type="hidden" id="history_skip">
   <input type="text" id="hist_skip_new">
   <button type="button" class="mini" id="hist_skip_add">{{S_HIST_ADD}}</button></div>
  <div class="skiplist" id="hist_skip_list"></div>
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
  <div class="row" data-adv><label>{{S_MINMS}}<span class="sub">{{S_SUB_MINMS}}</span></label><select id="min_record_ms"><option value="0">0 ms</option><option value="100">100 ms</option><option value="150">150 ms</option><option value="200">200 ms</option><option value="300">300 ms</option><option value="500">500 ms</option><option value="750">750 ms</option><option value="1000">1000 ms</option></select></div>
  <div class="row" data-adv><label>{{S_MAXSEC}}</label><select id="max_record_seconds"><option value="30">30 s</option><option value="60">60 s</option><option value="120">120 s</option><option value="180">180 s</option><option value="300">300 s</option></select></div>
 </div>
 <div class="card" id="behcard">
  <div class="row"><label>{{S_AUTOENTER}}<span class="sub">{{S_SUB_ENTER}}</span></label><input type="checkbox" id="auto_enter"></div>
  <div class="row" data-adv><label>{{S_RESTORE}}<span class="sub">{{S_SUB_CLIP}}</span></label><input type="checkbox" id="restore_clipboard"></div>
  <div class="row"><label>{{S_TYPEMODE}}<span class="sub">{{S_SUB_TYPE}}</span></label><input type="checkbox" id="type_mode"></div>
  <div class="row" data-adv><label>{{S_PASTE_DELAY}}<span class="sub">{{S_PASTE_DELAY_SUB}}</span></label>
   <select id="paste_delay_ms"><option value="0">0 ms</option><option value="50">50 ms</option><option value="100">100 ms</option><option value="250">250 ms</option><option value="500">500 ms</option><option value="1000">1000 ms</option></select></div>
 </div>
 <div class="card" id="ovcard">
  <label class="blklbl">{{S_SEC_OVERLAY}}</label>
  <div class="hint">{{S_OVERLAY_SUB}}</div>
  <div class="row"><label>{{S_OVERLAY}}</label><input type="checkbox" id="overlay"></div>
  <div class="row noline"><label>{{S_OV_FREE}}<span class="sub">{{S_OV_FREE_SUB}}</span></label><input type="checkbox" id="ov_free"></div>
  <div class="row"><label>{{S_OVPOS}}<span class="sub" id="ovpos_sub">{{S_OVPOS_SCHEME_SUB}}</span></label>
   <input type="hidden" id="overlay_position">
   <div class="ovscheme" id="ovscheme">
    <div class="ovcase">
     <div class="ovcrt" id="ovcrt">
      <span class="ovgrid"></span>
      <span class="ovscan"></span>
      <button type="button" class="ovzone" data-pos="top-left"></button>
      <button type="button" class="ovzone" data-pos="top"></button>
      <button type="button" class="ovzone" data-pos="top-right"></button>
      <button type="button" class="ovzone" data-pos="left"></button>
      <button type="button" class="ovzone" data-pos="center"></button>
      <button type="button" class="ovzone" data-pos="right"></button>
      <button type="button" class="ovzone" data-pos="bottom-left"></button>
      <button type="button" class="ovzone" data-pos="bottom"></button>
      <button type="button" class="ovzone" data-pos="bottom-right"></button>
      <span class="ovmini" id="ovmini" title="{{S_OVDRAG}}"></span>
     </div>
     <i class="ovvents"></i><i class="ovknob"></i><i class="ovled"></i>
    </div>
    <div class="ovneck"></div><div class="ovbase"></div>
   </div></div>
  <div class="row"><label>{{S_OVPOS_CARET}}<span class="sub">{{S_OVPOS_SUB}}</span></label><input type="checkbox" id="ov_caret"></div>
  <div class="row" id="ovmon_row" style="display:none"><label>{{S_OVMON}}<span class="sub">{{S_OVMON_SUB}}</span></label>
   <select id="overlay_monitor"></select></div>
  <div class="row"><label>{{S_OVTEXT}}<span class="sub">{{S_OVTEXT_SUB}}</span></label><input type="checkbox" id="overlay_text"></div>
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
  <div class="row"><label>{{S_RECLANG}}<span class="sub">{{S_SRCLANG_SUB}}</span></label>
   <select id="language">
    <option value="auto">{{S_RECAUTO}}</option>
    <option value="de">Deutsch</option><option value="en">English</option>
    <option value="es">Español</option><option value="fr">Français</option>
    <option value="it">Italiano</option><option value="pl">Polski</option>
    <option value="uk">Українська</option><option value="ru">Русский</option>
   </select></div>
  <div class="row"><label>{{S_TR_DEFAULT}}<span class="sub warn" id="tr_unavail" style="display:none"></span></label><input type="checkbox" id="tr_default"></div>
  <div class="row"><label>{{S_TR_TARGET}}<span class="sub">{{S_SUB_TRTARGET}}</span></label>
   <select id="translate_target">
    <option value="de">Deutsch</option><option value="en">English</option>
    <option value="es">Español</option><option value="fr">Français</option>
    <option value="it">Italiano</option><option value="pl">Polski</option>
    <option value="uk">Українська</option><option value="ru">Русский</option>
   </select></div>
  <div class="row"><label>{{S_TR_ASK}}</label>
   <select id="translate_ask">
    <option value="never">{{S_TR_ASK_NEVER}}</option>
    <option value="always">{{S_TR_ASK_ALWAYS}}</option>
    <option value="timeout">{{S_TR_ASK_TIMEOUT}}</option>
   </select></div>
  <div class="row"><label>{{S_TR_SECONDS}}</label><select id="translate_ask_seconds"><option value="2">2 s</option><option value="3">3 s</option><option value="4">4 s</option><option value="5">5 s</option><option value="7">7 s</option><option value="10">10 s</option></select></div>
  <div class="row"><label>{{S_TR_LANGS}}<span class="sub">{{S_TR_LANGS_SUB}}</span></label>
   <span id="trlangs" style="display:grid;grid-template-columns:repeat(4,minmax(64px,max-content));gap:9px 20px">
    <label style="flex:none"><input type="checkbox" id="tl_de"> DE</label>
    <label style="flex:none"><input type="checkbox" id="tl_en"> EN</label>
    <label style="flex:none"><input type="checkbox" id="tl_es"> ES</label>
    <label style="flex:none"><input type="checkbox" id="tl_fr"> FR</label>
    <label style="flex:none"><input type="checkbox" id="tl_it"> IT</label>
    <label style="flex:none"><input type="checkbox" id="tl_pl"> PL</label>
    <label style="flex:none"><input type="checkbox" id="tl_uk"> UK</label>
    <label style="flex:none"><input type="checkbox" id="tl_ru"> RU</label>
   </span></div>
 </div>
 <div class="card">
  <label class="blklbl">{{S_PRESETS}}</label>
  <div class="hint">{{S_PRESETS_HINT}}</div>
  <div id="langlist"></div>
 </div>
 <div class="card">
  <div class="ownm">
   <label>{{S_MFOLDER}}<span class="sub">{{S_OWNM_SUB}}</span></label>
   <div class="tl">
    <div class="tls"><i>1</i>
     <div class="tlc"><b>{{S_OWNM_S1}}</b><span class="tlsub">{{S_OWNM_S1S}} <span class="code">models\</span></span></div>
     <div class="tla"><button type="button" class="mini" onclick="appOpenModelsFolder()">{{S_OPEN_FOLDER}}</button></div>
    </div>
    <div class="tls"><i>2</i>
     <div class="tlc"><b>{{S_OWNM_S2}}</b><span class="tlsub">{{S_OWNM_S2S}}</span>
      <div class="fmts">
       <div class="fmtc"><div class="fr"><b>Whisper (GGML)</b><span class="cnt">{{S_OWNM_ONEFILE}}</span></div>
        <div class="chips"><span class="code">ggml-*.bin</span></div></div>
       <div class="fmtc"><div class="fr"><b>Sherpa-ONNX</b><span class="cnt">{{S_OWNM_FOLDERF}}</span></div>
        <div class="chips"><span class="code">encoder.onnx</span><span class="code">decoder.onnx</span><span class="code">tokens.txt</span></div></div>
      </div>
     </div>
    </div>
    <div class="tls"><i>3</i>
     <div class="tlc"><b>{{S_OWNM_S3}}</b><span class="tlsub tlok">✓ {{S_OWNM_S3S}}</span></div>
    </div>
   </div>
  </div>
  <div class="row"><label>{{S_UNLOAD}}<span class="sub">{{S_UNLOAD_SUB}}</span></label>
   <button type="button" class="mini" id="munload">{{S_UNLOAD_GO}}</button></div>
  <div class="row" style="border-bottom:0"><label>{{S_MCHECK}}<span class="sub">{{S_MCHECK_SUB}}</span></label>
   <button type="button" class="mini" id="mcheck">{{S_MCHECK_GO}}</button></div>
  <div class="micverdict" id="mcheck_out"></div>
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
  <label class="blklbl">{{S_SUB_DICT}}</label>
  <div class="hint">{{S_DICT_HINT}}</div>
  <div class="note warn" id="dict_warn" style="display:none"></div>
  <div class="rulefoot nowrap">
   <label class="srchbox"><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4"><circle cx="10.5" cy="10.5" r="6.5"/><line x1="15.5" y1="15.5" x2="21" y2="21"/></svg><input type="text" id="dict_filter" placeholder="{{S_LIST_FILTER_PH}}" autocomplete="off"><button type="button" class="clearx" id="dict_filter_clear" tabindex="-1" style="display:none">✕</button></label>
   <button type="button" class="mini" id="dict_add">{{S_DICT_ADD}}</button>
  </div>
  <div id="dictbody"></div>
 </div>
 <div class="card" data-adv>
  <label class="blklbl">{{S_SEC_REPLACE}}</label>
  <div class="hint">{{S_REPLACE_HINT}}</div>
  <div class="rulefoot nowrap">
   <label class="srchbox"><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4"><circle cx="10.5" cy="10.5" r="6.5"/><line x1="15.5" y1="15.5" x2="21" y2="21"/></svg><input type="text" id="repl_filter" placeholder="{{S_LIST_FILTER_PH}}" autocomplete="off"><button type="button" class="clearx" id="repl_filter_clear" tabindex="-1" style="display:none">✕</button></label>
   <button type="button" class="mini" id="repl_add">{{S_REPL_ADD}}</button>
  </div>
  <div id="replbody"></div>
 </div>
 <div class="card" data-adv>
  <label class="blklbl">{{S_SEC_CMD}}</label>
  <div class="hint">{{S_CMD_HINT}}</div>
  <div class="rulefoot nowrap">
   <label class="srchbox"><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4"><circle cx="10.5" cy="10.5" r="6.5"/><line x1="15.5" y1="15.5" x2="21" y2="21"/></svg><input type="text" id="cmd_filter" placeholder="{{S_LIST_FILTER_PH}}" autocomplete="off"><button type="button" class="clearx" id="cmd_filter_clear" tabindex="-1" style="display:none">✕</button></label>
   <button type="button" class="mini" id="cmd_add">{{S_CMD_ADD}}</button>
  </div>
  <div id="cmdbody"></div>
 </div>
</div>

<div class="page" role="tabpanel" aria-hidden="true" id="p-post">
 <div class="card">
  <div class="row noline"><label class="blklbl">{{S_POST_ENABLE}}<span class="sub warn" id="post_warn" style="display:none"></span></label><input type="checkbox" id="post_enabled"></div>
  <div class="hint">{{S_POST_HINT}}</div>
 </div>
 <div class="card" id="post_model_card">
  <label class="blklbl">{{S_POST_MODEL}}</label>
  <div class="srccard" id="src_local">
   <h3 class="srchead"><input type="checkbox" class="srcpick" id="pick_local">{{S_SRC_LOCAL}}<span class="hfhome" onclick="appHFHome()" title="huggingface.co">Hugging Face ↗</span></h3>
   <div class="sum" id="llm_sum"></div>
   <div id="llm_dl"></div>
   <div class="acts"><button type="button" class="mini" id="llm_catalog">{{S_LLM_CATALOG}}</button></div>
  </div>
  <div class="srccard" id="src_api">
   <h3 class="srchead"><input type="checkbox" class="srcpick" id="pick_api">{{S_POSTAPI}}</h3>
   <div class="hint">{{S_POSTAPI_HINT}}</div>
   <div class="sum" id="api_sum"></div>
   <div class="note warn" id="postapi_warn"></div>
   <div class="note bad" id="postapi_err"></div>
   <div class="acts"><button type="button" class="mini" id="api_test">{{S_API_TEST}}</button><button type="button" class="mini" id="api_edit">{{S_POSTAPI_SETUP}}</button></div>
   <div hidden>
    <input type="text" id="post_api_url">
    <input type="text" id="post_api_model">
    <select id="post_api_timeout_s"><option value="10">10 s</option><option value="20">20 s</option><option value="30">30 s</option><option value="60">60 s</option><option value="120">120 s</option></select>
   </div>
  </div>
 </div>
 <div class="card" id="post_prompts_card">
  <div class="blkhead">
   <div>
    <label class="blklbl">{{S_SUB_PROMPTS}}</label>
    <div class="hint">{{S_LLM_HINT}}</div>
   </div>
   <button type="button" class="mini" id="profadd">{{S_PROF_ADD}}</button>
  </div>
  <div id="profbody"></div>
 </div>
</div>


<div class="page" role="tabpanel" aria-hidden="true" id="p-system">
 <div class="card">
  <div class="row"><label>{{S_UILANG}}</label>
   <select id="ui_language">
    <option value="auto">{{S_AUTO}}</option>
    <option value="de">Deutsch</option>
    <option value="en">English</option>
    <option value="es">Español</option>
    <option value="fr">Français</option>
    <option value="it">Italiano</option>
    <option value="pl">Polski</option>
    <option value="uk">Українська</option>
    <option value="ru">Русский</option>
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
  <div class="row" data-adv><label>{{S_THREADS}}<span class="sub">{{S_SUB_THREADS}}</span></label><select id="threads"><option value="1">1</option><option value="2">2</option><option value="4">4</option><option value="6">6</option><option value="8">8</option><option value="12">12</option><option value="16">16</option></select></div>
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
  <div class="row"><label>GitHub</label><button type="button" class="mini" onclick="appRepoLink()">github.com &#8599;</button></div>
 </div>
</div>

<div class="page about" role="tabpanel" aria-hidden="true" id="p-help">
 <div class="card">{{S_HELP_HTML}}</div>
</div>

<div class="page about" role="tabpanel" aria-hidden="true" id="p-contacts">
 <div class="card">{{S_AUTHOR_HTML}}</div>
 <div class="card"><div class="row"><label>{{S_CONTACT_MAIL}}</label><span class="val" style="user-select:text">holdtotype@outlook.com</span></div></div>
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
     <option value="de">Deutsch</option><option value="en">English</option>
     <option value="es">Español</option><option value="fr">Français</option>
     <option value="it">Italiano</option><option value="pl">Polski</option>
     <option value="uk">Українська</option><option value="ru">Русский</option>
     <option value="auto">{{S_RECAUTO}}</option>
    </select></div>
   <div class="wizplan" id="wiz_plan"></div>
   <div class="wizrow" id="wiz_dlrow" style="display:none"><span class="wizbar"><i id="wiz_dlbar"></i></span><span class="mpct" id="wiz_dlpct"></span></div>
   <div class="wizrow"><button type="button" class="btn" id="wiz_dl">{{S_DL}}</button>
    <button type="button" class="btn ghost" id="wiz_dl_skip">{{S_WIZ_SKIP_DL}}</button>
    <span class="wizout" id="wiz_dlout"></span></div>
   <div class="wiztext" id="wiz_skipnote" style="display:none">{{S_WIZ_SKIP_NOTE}}</div>
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
const bools = ["beep","auto_enter","restore_clipboard","overlay","overlay_text","type_mode","server_autostart","check_updates","history"];
const texts = ["history_skip","post_api_model"];
let exeStored = CFG.server_exe || "";
let exeUnlocked = false;
let remoteURL = (CFG.server_url || "").trim();
let postURL = (CFG.post_api_url || "").trim();
let postEnabled = CFG.post_enabled !== false;
let postSource = CFG.post_source === "api" ? "api" : "local";
const nums  = ["threads","min_record_ms","max_record_seconds","translate_ask_seconds","server_port","paste_delay_ms","history_days","history_max","post_api_timeout_s"];
const sels  = ["ui_language","language","sound_theme","translate_target","translate_ask","hotkey_mode","theme","skin"];
const trAll = ["de","en","es","fr","it","pl","uk","ru"];
const L = {{L_JSON}};
const I_DL = '<svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2"><path d="M12 3v12"/><path d="M6 11l6 6 6-6"/><path d="M4 21h16"/></svg>';
const I_EJECT = '<svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2"><path d="M12 4l7 9H5z"/><path d="M5 19h14"/></svg>';
const I_FIND = '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4"><circle cx="10.5" cy="10.5" r="6.5"/><line x1="15.5" y1="15.5" x2="21" y2="21"/></svg>';

let profiles = (CFG.profiles || []).map(p=>Object.assign({}, p));
let activeProfiles = (CFG.active_profiles || []).slice();
let translateDefault = !!CFG.translate_default;
let translateHotkey = CFG.translate_hotkey || "";
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
    if(!await askConfirm(L.exewarn, L.exeedit, null, L.mtexe)) return;
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
  document.querySelectorAll("#p-help .card").forEach(c=>{
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
  let top = box.top - size.height - 6;
  if(top < 6) top = box.bottom + 6;
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
let hfSearched = false;
function fitLabel(f, need){
  const gb = (need/1024).toFixed(1);
  const col = f==="ok" ? "var(--green)" : (f==="warn" ? "var(--amber)" : "var(--bad)");
  const tip = f==="ok" ? L.fitok : (f==="warn" ? L.fitwarn : L.fitbad);
  return '<span title="'+esc(tip)+'" style="color:'+col+';font-size:12px;white-space:nowrap">&#9679; &#8776;'+gb+' GB</span>';
}
let llmRam = null;
let postTestErr = "";
let llmPicked = false;
function punctMode(){
  const el = document.getElementById("punctuation");
  return (el && el.value) || CFG.punctuation || "model";
}
function updatePostWarn(){
  const el = document.getElementById("post_warn");
  if(!el) return;
  let text = "";
  if(postEnabled){
    if(postTestErr) text = L.postbad.replace("%s", postTestErr);
    else if(postSource === "api" && !postURL) text = L.postnoapi;
    else if(postSource !== "api" && !llmPicked) text = L.postnomodel;
    else if(!activeProfiles.length && punctMode() !== "llm") text = L.postnoprompt;
  }
  el.textContent = text;
  el.style.display = text ? "" : "none";
}
function llmSize(mb){
  return mb >= 1024 ? (mb / 1024).toFixed(1) + " GB" : mb + " MB";
}
function renderLLMSummary(st){
  const box = document.getElementById("llm_sum");
  if(!box) return;
  const installed = st.installed || [];
  box.innerHTML = "";
  const line = (k, v, off)=>{
    const r = document.createElement("div");
    r.className = "sumrow";
    const a = document.createElement("span");
    a.className = "sumk";
    a.textContent = k;
    const b = document.createElement("span");
    b.className = "sumv" + (off ? " off" : "");
    b.textContent = v;
    r.appendChild(a); r.appendChild(b);
    box.appendChild(r);
  };
  const cur = installed.find(m=>m.file === selLLM) || installed.find(m=>m.active);
  if(cur){
    line(L.llmsummodel, cur.file);
    line(L.llmsumsize, llmSize(cur.size));
    line(L.llmsumcount, String(installed.length));
  } else if(installed.length){
    line(L.llmsummodel, L.llmnopick, true);
    line(L.llmsumcount, String(installed.length));
  } else {
    line(L.apisumstate, L.nollm, true);
  }
  if(st.ram) line(L.llmsumram, L.ramof.replace("%s", ((st.ram_free || st.ram) / 1024).toFixed(1)).replace("%s", (st.ram / 1024).toFixed(0)));
  const dlBox = document.getElementById("llm_dl");
  if(!dlBox) return;
  dlBox.innerHTML = "";
  (st.downloads || []).filter(d=>d.pct >= 0).forEach(d=>{
    const row = document.createElement("div");
    row.className = "llmdl";
    const nm = document.createElement("span");
    nm.className = "nm";
    nm.textContent = L.dlgoing + " " + d.file;
    row.appendChild(nm);
    const track = document.createElement("span");
    track.className = "mtrack";
    const fill = document.createElement("i");
    fill.style.width = Math.max(2, d.pct) + "%";
    track.appendChild(fill);
    row.appendChild(track);
    const pct = document.createElement("span");
    pct.className = "mpct";
    pct.textContent = d.pct > 0 ? d.pct + "%" : "…";
    row.appendChild(pct);
    const stop = document.createElement("button");
    stop.type = "button";
    stop.className = "iconbtn danger";
    stop.title = L.dlcancel;
    stop.innerHTML = "&#10005;";
    stop.onclick = async ()=>{ await appModelCancel("llm-" + d.file); refreshLLM(); };
    row.appendChild(stop);
    dlBox.appendChild(row);
  });
}
function openLLMCatalog(){
  const bg = document.createElement("div");
  bg.className = "modal-bg";
  const box = document.createElement("div");
  box.className = "modal llmcat";
  mTitle(box, L.llmcatalog);
  const search = document.createElement("div");
  search.id = "proc-search";
  box.appendChild(search);
  const list = document.createElement("div");
  list.id = "proc-models";
  box.appendChild(list);
  const row = document.createElement("div");
  row.className = "modal-btns";
  const close = document.createElement("button");
  close.type = "button";
  close.className = "btn ghost";
  close.textContent = L.dlgclose;
  const done = ()=>{
    bg.remove();
    document.removeEventListener("keydown", onKey, true);
    unlockBg();
    refreshLLM();
  };
  const topmost = ()=>[...document.querySelectorAll(".modal-bg")].pop() === bg;
  function onKey(e){
    if(e.key !== "Escape" || e.defaultPrevented || !topmost()) return;
    e.preventDefault();
    done();
  }
  close.onclick = done;
  bg.onclick = e=>{ if(e.target === bg) done(); };
  document.addEventListener("keydown", onKey, true);
  row.appendChild(close);
  box.appendChild(row);
  box.setAttribute("role", "dialog");
  box.setAttribute("aria-modal", "true");
  bg.appendChild(box);
  document.body.appendChild(bg);
  lockBg();
  initHFBox();
  refreshLLM();
}
async function refreshLLM(){
  const st = JSON.parse(await appLLM());
  const installed = st.installed || [];
  if(selLLM === null){
    const act = installed.find(m=>m.active);
    if(act) selLLM = act.file;
  }

  const sc = document.querySelector(".content");
  const keepScroll = sc ? sc.scrollTop : 0;
  llmRam = st;
  llmPicked = !!(installed.find(m=>m.file === selLLM) || installed.find(m=>m.active));
  updatePostWarn();
  renderLLMSummary(st);
  let busy = (st.downloads || []).some(d=>d.pct >= 0);
  const body = document.getElementById("proc-models");
  if(body){
    body.innerHTML = "";
    const blk = document.createElement("label");
    blk.className = "blklbl";
    blk.textContent = L.llmblock;
    body.appendChild(blk);
    if(!installed.length && !(st.downloads||[]).length){
      const empty = document.createElement("div");
      empty.className = "hint";
      empty.textContent = L.llmnonehint;
      body.appendChild(empty);
    }
    installed.forEach(m=>{
      const cur = selLLM === m.file;
      const div = document.createElement("div");
      div.className = "mrow" + (cur ? " cur" : "");
      div.dataset.f = m.file;
      const sw = document.createElement("input");
      sw.type = "checkbox";
      sw.className = "llmpick";
      sw.checked = cur;
      sw.onchange = async ()=>{
        if(!sw.checked){ sw.checked = true; return; }
        selLLM = m.file;
        await doSave();
        refreshLLM();
      };
      div.appendChild(sw);
      const nm = document.createElement("span");
      nm.className = "mdesc";
      nm.style.flex = "1";
      nm.textContent = m.file;
      div.appendChild(nm);
      const state = document.createElement("span");
      state.className = "mstate";
      state.textContent = m.loaded ? L.llminmem : L.llmondisk;
      div.appendChild(state);
      const sz = document.createElement("span");
      sz.className = "msize";
      sz.textContent = llmSize(m.size);
      div.appendChild(sz);
      const eject = document.createElement("button");
      eject.type = "button";
      eject.className = "iconbtn";
      eject.title = L.llmeject;
      eject.disabled = !m.loaded;
      eject.innerHTML = I_EJECT;
      eject.onclick = async ()=>{ await appLLMUnload(); refreshLLM(); };
      div.appendChild(eject);
      const del = document.createElement("button");
      del.type = "button";
      del.className = "iconbtn danger";
      del.title = L.del;
      del.dataset.a = "ldel";
      del.dataset.f = m.file;
      del.innerHTML = "&#10005;";
      div.appendChild(del);
      body.appendChild(div);
    });
    (st.downloads || []).forEach(d=>{
      const div = document.createElement("div");
      div.className = "mrow";
      const wait = '<input type="checkbox" class="llmpick" title="'+esc(L.llmpickwait)+'" disabled>';
      if(d.pct >= 0){
        div.innerHTML = wait+'<span class="mdesc" style="flex:1">'+esc(d.file)+'</span><span class="mpct">'+(d.pct>0?d.pct+"%":"…")+'</span>'+
          '<button class="iconbtn danger" title="'+L.dlcancel+'" data-a="lcancel" data-f="'+esc(d.file)+'">&#10005;</button>';
      } else {
        div.innerHTML = wait+'<span class="mdesc" style="flex:1">'+esc(d.file)+'</span><span class="mpct">! '+esc(d.err)+'</span>';
      }
      body.appendChild(div);
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
        if(!await askConfirm(L.confirmdel.replace("%s", f), L.del, null, L.mtdel)) return;
        toast(await appLLMDel(f));
        if(selLLM === f){
          selLLM = null;
        }
        refreshLLM();
      };
    });
  }

  const ramline = document.getElementById("hf_ramline");
  if(ramline) ramline.innerHTML = L.ramavail
    .replace("%s", '<b>'+((st.ram_free||st.ram)/1024).toFixed(1)+'</b>')
    .replace("%s", (st.ram/1024).toFixed(0));
  const fitline = document.getElementById("hf_fitline");
  if(fitline) fitline.innerHTML =
    '<span class="dot" style="color:var(--green)">&#9679;</span>'+L.fitok+
    '<span class="dot" style="color:var(--amber)">&#9679;</span>'+L.fitwarn+
    '<span class="dot" style="color:var(--bad)">&#9679;</span>'+L.fitbad;

  renderProfiles(document.getElementById("profbody"), {state: (installed.length || st.external) ? "installed" : "absent"});
  if(sc) sc.scrollTop = keepScroll;
  if(busy) setTimeout(refreshLLM, 900);
}
function initHFBox(){
  const sbody = document.getElementById("proc-search");
  if(!sbody) return;
  sbody.innerHTML = '<div class="srchrow"><label class="srchbox">'+I_FIND+
    '<input type="text" id="hf_q" placeholder="'+L.hfph+'">'+
    '<button type="button" class="iconbtn" id="hf_clr" title="'+L.cancel+'" style="display:none">&#10005;</button></label>'+
    '<button type="button" class="mini" id="hf_go">'+L.hfgo+'</button></div>'+
    '<div class="ramline" id="hf_ramline"></div>'+
    '<div class="ramline fitleg" id="hf_fitline"></div>'+
    '<div class="filtrow"><label class="fitrow"><input type="checkbox" id="hf_fit" checked> '+L.hffit+'</label>'+
    '<span class="fcount" id="hf_count"></span></div>'+
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
  document.getElementById("hf_fit").onchange = renderHF;
  renderHF();
}
async function doHFSearch(){
  const q = document.getElementById("hf_q").value;
  hfOpenRepo = null; hfFiles = [];
  hfSearched = true;
  const res = JSON.parse(await appLLMSearch(q));
  hfRepos = res.repos || [];
  if(res.error) toast(res.error);
  renderHF();
}
function renderHF(){
  const box = document.getElementById("hf_results");
  if(!box) return;
  box.innerHTML = "";
  const count = document.getElementById("hf_count");
  if(count) count.textContent = hfSearched ? L.llmfound.replace("%d", hfRepos.length) : L.llmnosearch;
  if(!hfRepos.length){
    const none = document.createElement("div");
    none.className = "searchempty";
    none.textContent = hfSearched ? L.listnothing.replace("%s", hfQuery) : L.llmsearchhint;
    box.appendChild(none);
    return;
  }
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
      const fitBox = document.getElementById("hf_fit");
      const fitOnly = !fitBox || fitBox.checked;
      const shown = fitOnly ? hfFiles.filter(f=>f.fit !== "bad") : hfFiles;
      shown.forEach(f=>{
        const fd = document.createElement("div");
        fd.className = "mrow";
        fd.style.paddingLeft = "22px";
        fd.innerHTML = '<span class="mdesc" style="flex:1">'+esc(f.file)+'</span>'+
          '<span class="msize">'+(f.size>=1024?(f.size/1024).toFixed(1)+" GB":f.size+" MB")+'</span>'+
          '<span>'+fitLabel(f.fit, f.need)+'</span>'+
          '<button class="iconbtn" title="'+L.dl+'" data-repo="'+esc(r.id)+'" data-file="'+esc(f.file)+'" data-size="'+(f.size||0)+'">'+I_DL+'</button>';
        box.appendChild(fd);
      });
      if(shown.length < hfFiles.length){
        const hid = document.createElement("div");
        hid.className = "mrow";
        hid.style.cssText = "padding-left:22px;color:var(--dim);font-size:12px";
        hid.textContent = L.hfhidden.replace("%s", String(hfFiles.length - shown.length));
        box.appendChild(hid);
      }
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
  const rows = profiles;
  if(!rows.length){
    const empty = document.createElement("div");
    empty.className = "ruleempty";
    empty.textContent = L.noprompts;
    body.appendChild(empty);
  }
  const list = document.createElement("div");
  list.className = "plist";
  rows.forEach(p=>{
    const div = document.createElement("div");
    div.className = "prow";
    div.dataset.id = p.id;
    const checked = activeProfiles.includes(p.id) ? " checked" : "";
    let html = '<span class="grip" title="'+L.pdrag+'">&#8942;&#8942;</span>';
    html += '<input type="checkbox" class="profcb" value="'+esc(p.id)+'"'+checked+'>';
    html += '<span class="pnm">'+esc(p.name)+'</span>';
    html += '<button class="iconbtn" title="'+L.pedit+'" data-a="edit" data-id="'+esc(p.id)+'">&#9998;</button>';
    html += '<button class="iconbtn danger" title="'+L.del+'" data-a="pdel" data-id="'+esc(p.id)+'">&#10005;</button>';
    div.innerHTML = html;
    div.querySelector(".grip").addEventListener("pointerdown", e=>startPromptDrag(e, div, list));
    list.appendChild(div);
  });
  body.appendChild(list);
  body.querySelectorAll("input.profcb").forEach(r=>{
    r.onchange = ()=>{
      if(r.checked){
        if(!activeProfiles.includes(r.value)) activeProfiles.push(r.value);
      } else {
        activeProfiles = activeProfiles.filter(x=>x!==r.value);
      }
      updatePostWarn();
    };
  });
  body.querySelectorAll("button[data-a]").forEach(b=>{
    b.onclick = async ()=>{
      const id = b.dataset.id;
      if(b.dataset.a === "edit"){
        editPromptDialog(profiles.find(x=>x.id===id), false);
        return;
      }
      const p = profiles.find(x=>x.id===id);
      if(!await askConfirm(L.confirmdel.replace("%s", p ? p.name : id), L.del, null, L.mtdelprompt)) return;
      profiles = profiles.filter(x=>x.id!==id);
      activeProfiles = activeProfiles.filter(x=>x!==id);
      if(expandedID === id) expandedID = null;
      refreshLLM();
      await doSave();
    };
  });
}

function startPromptDrag(e, row, list){
  if(e.button !== undefined && e.button !== 0) return;
  e.preventDefault();
  const rows = [...list.querySelectorAll(".prow")];
  const from = rows.indexOf(row);
  const box = row.getBoundingClientRect();
  const grabY = e.clientY - box.top;
  const ghost = row.cloneNode(true);
  ghost.classList.add("ghost");
  ghost.style.width = box.width + "px";
  ghost.style.transform = "translate(" + box.left + "px," + box.top + "px)";
  document.body.appendChild(ghost);
  row.classList.add("dragging");
  const line = document.createElement("div");
  line.className = "dropline";
  list.appendChild(line);
  let to = from;
  const place = clientY=>{
    const lb = list.getBoundingClientRect();
    let idx = rows.length;
    for(let i = 0; i < rows.length; i++){
      const rb = rows[i].getBoundingClientRect();
      if(clientY < rb.top + rb.height / 2){ idx = i; break; }
    }
    to = idx;
    const anchor = idx < rows.length ? rows[idx].getBoundingClientRect().top : rows[rows.length - 1].getBoundingClientRect().bottom;
    line.style.top = (anchor - lb.top - 1) + "px";
  };
  place(e.clientY);
  const move = ev=>{
    ghost.style.transform = "translate(" + box.left + "px," + (ev.clientY - grabY) + "px)";
    place(ev.clientY);
  };
  const up = async ()=>{
    document.removeEventListener("pointermove", move);
    document.removeEventListener("pointerup", up);
    ghost.remove();
    line.remove();
    row.classList.remove("dragging");
    let target = to;
    if(target > from) target--;
    if(target !== from && target >= 0){
      const item = profiles.splice(from, 1)[0];
      profiles.splice(target, 0, item);
      await doSave();
    }
    refreshLLM();
  };
  document.addEventListener("pointermove", move);
  document.addEventListener("pointerup", up);
}
async function editPromptDialog(p, fresh){
  if(!p) return;
  const name = document.createElement("input");
  name.type = "text"; name.className = "pfname"; name.value = p.name || ""; name.placeholder = L.pnameph;
  const prompt = document.createElement("textarea");
  prompt.className = "pfprompt";
  prompt.value = p.prompt || "";
  const sample = document.createElement("input");
  sample.type = "text"; sample.className = "pfsample"; sample.placeholder = L.ptestph;
  const out = document.createElement("div");
  out.className = "pfout";
  out.id = "pf_result";
  const ok = await formModal(fresh ? L.fmadd : L.fmsave, body=>{
    body.appendChild(fmRow(clearWrap(name), L.pname));
    const pw = document.createElement("div");
    const pl = document.createElement("label");
    pl.className = "plabel";
    pl.textContent = L.pprompt;
    pw.appendChild(pl);
    pw.appendChild(prompt);
    body.appendChild(pw);
    const tw = document.createElement("div");
    const tl = document.createElement("label");
    tl.className = "plabel";
    tl.textContent = L.ptest;
    tw.appendChild(tl);
    const trow = document.createElement("div");
    trow.className = "tryrow";
    trow.appendChild(clearWrap(sample));
    const run = document.createElement("button");
    run.type = "button";
    run.className = "mini";
    run.innerHTML = "&#9654;";
    run.onclick = ()=>{
      const s = sample.value.trim();
      if(!s) return;
      out.textContent = "…";
      appLLMTest(prompt.value, s);
    };
    trow.appendChild(run);
    tw.appendChild(trow);
    tw.appendChild(out);
    body.appendChild(tw);
  }, null, fresh ? L.pfnew : L.pfedit, true);
  if(!ok) return;
  const nm = name.value.trim();
  if(!nm) return;
  p.name = nm;
  p.prompt = prompt.value;
  if(fresh) profiles.push(p);
  refreshLLM();
  await doSave();
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

function comboCaptured(combo, warn){
  if(warn) toast(warn, "warn");
  if(!captureFor) return;
  captureFor = null;
  refreshLLM();
}

function llmTestResult(out){
  const el = document.getElementById("pf_result");
  if(el) el.textContent = out;
  postTestErr = String(out || "").startsWith("⚠") ? String(out).replace(/^⚠\s*/, "") : "";
  updatePostWarn();
}

let micChosen = "";
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
  showPostErr(s.post_err || "");
  const berr = document.getElementById("state_backend");
  if(berr){
    berr.style.display = s.backend_err ? "" : "none";
    set("state_backend_text", s.backend_err || "");
  }
  set("state_hotkey", s.hotkey);
  setWithTip("state_mic", s.mic);
  setWithTip("state_active", s.active_model);
  set("state_active_lang", s.active_lang || "");
  setWithTip("state_llm", s.llm);
  set("state_ram", s.ram);
  set("state_loaded", s.loaded_now || "—");
  const wk = document.getElementById("state_week_row");
  if(wk){ wk.style.display = s.week_line ? "" : "none"; }
  set("state_week", s.week_line || "—");
  const abox = document.getElementById("state_assigned");
  if(abox){
    const sig = JSON.stringify(s.assigned || []);
    if(abox.dataset.sig !== sig){
      abox.dataset.sig = sig;
      abox.innerHTML = "";
      (s.assigned || []).forEach(r=>{
        const row = document.createElement("div");
        row.className = "arow" + (r.current ? " on" : "");
        row.innerHTML = '<i class="led'+(r.state==="ready"||r.state==="remote" ? " on" : " warn")+'"></i>'+
          '<span class="aname">'+esc(r.model)+'</span>'+
          '<span class="alangs">'+esc(r.langs)+'</span>'+
          (r.state==="missing" ? '<span class="amiss">'+L.notinstalled+'</span>' : "");
        row.onclick = ()=>show("models");
        abox.appendChild(row);
      });
    }
  }
  const inst = document.getElementById("state_installed");
  if(inst) inst.textContent = (s.installed_models || []).join(", ") || "—";
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
  capLed("state_active_led", "state_active_btn", s.active_state);
  led("state_llm_led", s.llm_ok, !s.llm_ok);
  const llmBtn = document.getElementById("state_llm_btn");
  if(llmBtn) llmBtn.textContent = s.llm_ok ? L.change : L.get;
  led("st_led", s.ready);
  const badge = (id, v, full, cls)=>{
    const el = document.getElementById(id);
    if(!el) return;
    el.textContent = v || "";
    if(full || v) el.dataset.tip = full || v; else delete el.dataset.tip;
    el.classList.toggle("miss", cls === "miss");
  };
  const missing = (s.assigned || []).some(r=>r.state === "missing");
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
let langModels = Object.assign({}, CFG.lang_models || {});
let modelRowsCache = [];
const presetLangs = [["auto", ""], ["de", "Deutsch"], ["en", "English"], ["es", "Español"], ["fr", "Français"], ["it", "Italiano"], ["pl", "Polski"], ["uk", "Українська"], ["ru", "Русский"]];
function rowById(id){ return modelRowsCache.find(m=>m.id===id) || null; }
function eligibleFor(m, lang){
  if(lang === "auto") return !!m.auto;
  return m.langs === "*" || (m.langs || "").split(",").includes(lang);
}
function effectiveFor(lang){
  const own = langModels[lang];
  if(own && rowById(own)) return own;
  if(lang !== "auto"){
    const uni = langModels["auto"];
    if(uni && rowById(uni)) return uni;
  }
  return "medium-q5_0";
}
function curLang(){
  const el = document.getElementById("language");
  return (el && el.value) || "auto";
}
async function assignModel(lang, id){
  const row = rowById(id);
  if(!row) return;
  if(!eligibleFor(row, lang)){
    toast(L.notforlang.replace("%s", row.name), "warn");
    return;
  }
  const trSw = document.getElementById("tr_default");
  if(trSw && trSw.checked && !trSw.disabled && !row.translate){
    if(!await askConfirm(L.trmodeldlg.replace("%s", row.name), L.trconfirm, null, L.mttroff)){
      refreshModels();
      return;
    }
    trSw.checked = false;
  }
  if(row.state === "absent" && !row.manual){
    const size = row.size ? row.size + " MB" : "";
    if(!await askConfirm(L.dlask.replace("%s", row.name).replace("%s", size), L.dlstart, null, L.mtdl)){
      refreshModels();
      return;
    }
    pendingDl = row.id;
    await appModelDl(row.id);
  }
  langModels[lang] = id;
  await doSave();
  refreshModels();
}
let openLang = null;
const REC_MODEL = {auto: "large-v3-turbo-q5_0", ru: "gigaam-v3", en: "medium-q5_0", uk: "parakeet-v3", de: "parakeet-v3", fr: "parakeet-v3", es: "parakeet-v3", it: "parakeet-v3", pl: "parakeet-v3"};
function langWord(m){
  if(m.custom) return L.langsq;
  if(m.langs === "*") return L.langsc.replace("%d", "99");
  const parts = (m.langs || "").split(",").filter(Boolean);
  if(parts.length > 1) return L.langsc.replace("%d", String(parts.length));
  const native = (presetLangs.find(p=>p[0]===m.langs) || [])[1];
  return native || (m.langs || "").toUpperCase();
}
function trWord(m){
  if(!m.translate) return "";
  return m.trlangs ? L.translist.replace("%s", m.trlangs.toUpperCase().split(",").join(", ")) : L.tren;
}
function renderLangs(){
  const box = document.getElementById("langlist");
  if(!box) return;
  const sc = document.querySelector(".content");
  const keepScroll = sc ? sc.scrollTop : 0;
  box.innerHTML = "";
  const dls = modelRowsCache.filter(m=>m.state === "downloading");
  if(dls.length){
    const d = document.createElement("div");
    d.className = "dlline";
    d.textContent = "⇩ " + L.dlgoing + " " + dls.map(m=>m.name + " · " + (m.pct > 0 ? m.pct + "%" : "…")).join(" · ");
    box.appendChild(d);
  }
  presetLangs.forEach(([lang, native])=>{
    const inherit = lang !== "auto" && !(langModels[lang] && rowById(langModels[lang]));
    const m = rowById(effectiveFor(lang));
    const row = document.createElement("div");
    row.className = "lrow" + (inherit ? " dim" : "") + (openLang === lang ? " open" : "");
    row.dataset.lang = lang;
    const pct = m && m.state === "downloading" ? ' <span class="mpct" data-id="' + m.id + '">' + (m.pct > 0 ? m.pct + "%" : "…") + '</span>' : "";
    row.innerHTML = '<span class="plang">' + (lang === "auto" ? L.recauto : native) + '</span>' +
      '<span class="lmodel">' + (inherit ? L.asauto + ' · ' : '') + esc(m ? m.name : "—") + pct + '</span>' +
      '<span class="larr">' + (openLang === lang ? "▲" : "▼") + '</span>';
    row.onclick = ()=>{ openLang = openLang === lang ? null : lang; renderLangs(); };
    box.appendChild(row);
    if(openLang === lang) box.appendChild(buildPicker(lang));
  });
  if(sc) sc.scrollTop = keepScroll;
}
function buildPicker(lang){
  const pick = document.createElement("div");
  pick.className = "lpick";
  const eff = effectiveFor(lang);
  const explicit = lang === "auto" || !!(langModels[lang] && rowById(langModels[lang]));
  const rec = REC_MODEL[lang];
  const list = modelRowsCache.filter(m=>eligibleFor(m, lang))
    .sort((a,b)=>a.name.localeCompare(b.name));
  list.forEach(m=>{
    const cur = m.id === eff && explicit;
    const card = document.createElement("div");
    card.className = "pcard" + (cur ? " cur" : "");
    card.dataset.id = m.id;
    const chip = m.id === rec ? '<span class="pchip">'+L.recchip+'</span>' : "";
    const bar = (label, v)=>'<span class="mbar"><span class="mbl">'+label+'</span><span class="mtrack"><i style="width:'+(v*20)+'%"></i></span></span>';
    const bars = (m.accuracy || m.speed) ? '<span class="mbars">'+bar(L.acc, m.accuracy||0)+bar(L.spd, m.speed||0)+'</span>' : "";
    let act = "";
    if(m.state === "downloading") act = '<span class="mpct" data-id="'+m.id+'">'+(m.pct>0?m.pct+"%":"…")+'</span><button class="iconbtn danger" title="'+L.dlcancel+'" data-a="cancel" data-id="'+m.id+'">&#10005;</button>';
    else if(m.state === "absent" && m.manual) act = '<button class="mini" data-a="link" data-id="'+m.id+'">'+L.manuallink+' &#8599;</button>';
    else if(m.state === "absent") act = '<span class="psize">'+(m.size?m.size+" MB":"")+'</span><button class="iconbtn" title="'+L.dl+'" data-a="dl" data-id="'+m.id+'">'+I_DL+'</button>';
    else if(cur) act = '<span class="pcur">&#10003;</span>';
    else act = (m.loaded ? '<button class="iconbtn" title="'+L.unload+'" data-a="unload" data-id="'+m.id+'">&#9167;</button>' : "")+
      '<button class="iconbtn danger" title="'+L.del+'" data-a="del" data-id="'+m.id+'" data-name="'+esc(m.name)+'">&#10005;</button>';
    const note = m.state === "absent" && m.manual ? '<span class="pdesc" style="color:var(--amber)">'+L.manualnote+'</span>' : "";
    const langNames = m.langs && m.langs !== "*" ? m.langs.split(",").map(c=>trLangName(c.trim())).join(", ") : "";
    const langChip = '<span class="pmt"'+(langNames && m.langs.includes(",") ? ' data-tip="'+esc(langNames)+'"' : '')+'>'+esc(langWord(m))+'</span>';
    const trNames = m.translate && m.trlangs ? m.trlangs.split(",").map(c=>trLangName(c.trim())).join(", ") : "";
    const trChip = m.translate ? '<span class="pmt"'+(trNames ? ' data-tip="'+esc(trNames)+'"' : '')+'>'+esc(trWord(m))+'</span>' : '';
    const meta = langChip + trChip +
      (m.ram ? '<span class="pram">≈'+m.ram+' MB RAM</span>' : '');
    const sw = '<input type="checkbox" class="psw"'+(cur?' checked':'')+((m.state === "absent" && m.manual)?' disabled':'')+' data-sw="'+m.id+'">';
    const info = m.desc ? '<span class="pinfo" data-tip="'+esc(m.desc)+'">i</span>' : "";
    card.innerHTML = '<span class="ptop">'+sw+'<span class="pname"><span class="pnm">'+esc(m.name)+'</span>'+chip+info+'</span>'+bars+'</span>'+
      note+
      '<span class="pmeta">'+meta+'<span class="pact">'+act+'</span></span>';
    pick.appendChild(card);
  });
  if(lang !== "auto" && langModels[lang]){
    const back = document.createElement("button");
    back.type = "button";
    back.className = "mini";
    back.textContent = L.backauto;
    back.onclick = async ()=>{ delete langModels[lang]; await doSave(); refreshModels(); };
    pick.appendChild(back);
  }
  wireModelActions(pick);
  pick.querySelectorAll("input.psw").forEach(sw=>{
    sw.addEventListener("click", e=>e.stopPropagation());
    sw.addEventListener("change", async e=>{
      e.stopPropagation();
      const id = sw.dataset.sw;
      if(sw.checked){ assignModel(lang, id); }
      else if(lang !== "auto" && langModels[lang] === id){
        delete langModels[lang];
        await doSave();
        refreshModels();
      } else {
        sw.checked = true;
      }
    });
  });
  return pick;
}
function wireModelActions(root){
  root.querySelectorAll("button[data-a]").forEach(b=>{
    b.onclick = async (e)=>{
      e.stopPropagation();
      if(b.dataset.a === "dl"){ await appModelDl(b.dataset.id); pendingDl = b.dataset.id; }
      else if(b.dataset.a === "link"){ appModelLink(b.dataset.id); return; }
      else if(b.dataset.a === "unload"){ appUnloadEngines(); toast(L.unloaded, "ok"); }
      else if(b.dataset.a === "cancel"){
        await appModelCancel(b.dataset.id);
        if(pendingDl === b.dataset.id) pendingDl = null;
      }
      else {
        const isActive = b.dataset.id === activeModelId;
        const mname = b.dataset.name || b.dataset.id;
        const ask = isActive ? L.delactive.replace("%s", mname) : L.confirmdel.replace("%s", mname);
        if(!await askConfirm(ask, L.del, null, L.mtdel)) return;
        toast(await appModelDel(b.dataset.id, isActive));
      }
      refreshModels();
    };
  });
}
function modelsSignature(rows){
  return rows.map(m=>[m.id, m.state, (m.serves||[]).join("+"), m.loaded?1:0, m.err||""].join(":")).join("|")
    + "@" + openLang + "@" + JSON.stringify(langModels) + "@" + curLang();
}
async function refreshModels(){
  const rows = JSON.parse(await appModels());
  modelRowsCache = rows;
  const activeRow = rows.find(m=>m.state==="active");
  activeModelId = activeRow ? activeRow.id : null;
  updateDictWarn();
  if(pendingDl){
    const row = rows.find(m=>m.id===pendingDl);
    if(!row){ pendingDl = null; }
    else if(row.state === "installed" || row.state === "active"){ pendingDl = null; toast(L.mdlready, "ok"); refreshState(); }
    else if(row.state === "absent" && row.err){ pendingDl = null; toast(row.err, "error"); }
  }
  const el = document.getElementById("langlist");
  const sig = modelsSignature(rows);
  let busy = rows.some(m=>m.state === "downloading");
  if(el.dataset.sig === sig && busy){
    rows.filter(m=>m.state === "downloading").forEach(m=>{
      el.querySelectorAll('.mpct[data-id="'+m.id+'"]').forEach(p=>{ p.textContent = m.pct > 0 ? m.pct + "%" : "…"; });
    });
    const dl = el.querySelector(".dlline");
    if(dl) dl.textContent = "⇩ " + L.dlgoing + " " + rows.filter(m=>m.state === "downloading").map(m=>m.name + " · " + (m.pct > 0 ? m.pct + "%" : "…")).join(" · ");
  } else {
    el.dataset.sig = sig;
    renderLangs();
  }
  syncTrControls();
  if(busy || pendingDl) setTimeout(refreshModels, 900);
}
function lockBg(){
  document.querySelectorAll(".content, .snav, .header").forEach(el=>el.setAttribute("inert", ""));
}
function unlockBg(){
  if(document.querySelector(".modal-bg")) return;
  document.querySelectorAll(".content, .snav, .header").forEach(el=>el.removeAttribute("inert"));
}
function mTitle(box, title){
  if(!title) return;
  const h = document.createElement("div");
  h.className = "mtitle";
  h.textContent = title;
  box.appendChild(h);
}
function askConfirm(text, okText, cancelText, title){
  return new Promise(resolve=>{
    const bg = document.createElement("div");
    bg.className = "modal-bg";
    const box = document.createElement("div");
    box.className = "modal";
    mTitle(box, title);
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
      unlockBg();
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
    lockBg();
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

const numSels = ["threads","min_record_ms","max_record_seconds","translate_ask_seconds","paste_delay_ms","history_days","history_max","post_api_timeout_s"];
function load(){
  document.getElementById("punctuation").value = CFG.punctuation || "model";
  document.getElementById("punctuation").addEventListener("change", updatePostWarn);
  trAll.forEach(l=>{ document.getElementById("tl_"+l).checked = (CFG.translate_ask_langs||[]).includes(l); });
  const trd = document.getElementById("tr_default");
  trd.checked = translateDefault || (CFG.translate_ask || "never") !== "never";
  trd.onchange = ()=>{ syncTrControls(); };
  const askSel = document.getElementById("translate_ask");
  let askPrev = askSel.value;
  askSel.addEventListener("change", async e=>{
    const tgt = document.getElementById("translate_target").value;
    const others = trAll.filter(l=>l !== tgt && document.getElementById("tl_"+l).checked);
    if(askSel.value === "never" && others.length > 0){
      e.stopPropagation();
      if(!await askConfirm(L.tronedlg.replace("%s", trLangName(tgt)), L.trconfirm, null, L.mttrone)){
        askSel.value = askPrev;
        syncTrControls();
        return;
      }
      askPrev = askSel.value;
      syncTrControls();
      doSave();
      return;
    }
    askPrev = askSel.value;
    syncTrControls();
  });
  const tgtSel = document.getElementById("translate_target");
  tgtSel.addEventListener("change", ()=>{
    const cb = document.getElementById("tl_"+tgtSel.value);
    if(cb && !cb.checked) cb.checked = true;
    syncTrControls();
  });
  trAll.forEach(l=>{
    const cb = document.getElementById("tl_"+l);
    cb.addEventListener("change", async e=>{
      if(!cb.checked && l === document.getElementById("translate_target").value){
        e.stopPropagation();
        cb.checked = true;
        await askAlert(L.trlockmsg.split("%s").join(trLangName(l)), L.trlockok, L.mttrlock);
        syncTrControls();
      }
    });
  });
  document.getElementById("language").addEventListener("change", ()=>{ syncTrControls(); });
  initDict();
  document.getElementById("mic_device").onchange = applyMic;
  document.getElementById("mic_refresh").onclick = refreshMics;
  const micChk = document.getElementById("mic_check");
  if(micChk) micChk.onclick = micCheck;
  const mChk = document.getElementById("mcheck");
  if(mChk) mChk.onclick = modelsCheck;
  const mUnload = document.getElementById("munload");
  if(mUnload) mUnload.onclick = async ()=>{ await appUnloadEngines(); toast(L.unloaded, "ok"); refreshModels(); };
  refreshMics();
  startMicMeter();
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
    if(url && !await askConfirm(L.remoteask.replace("%s", url), null, null, L.mtremote)){
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
async function runPostTest(btn, url, model, timeout, key){
  const was = btn.textContent;
  btn.disabled = true;
  btn.textContent = L.apitestrun;
  let r = {ok: false, message: "", severity: "error"};
  try {
    r = JSON.parse(await appPostTest(url || "", model || "", key || "", String(timeout || "")));
  } catch(e){
    r.message = String(e);
  }
  btn.disabled = false;
  btn.textContent = was;
  toast(r.message, r.severity);
  showPostErr(r.ok ? "" : r.message, true);
  return r;
}
function showPostErr(msg, force){
  const el = document.getElementById("postapi_err");
  if(el) el.textContent = (force || postSource === "api") ? msg : "";
}
function syncPostWarn(){
  const note = document.getElementById("postapi_warn");
  if(note) note.textContent = postURL ? L.postwarn : "";
}
async function updPostKeyState(){
  if(!window.appPostKeySet) return;
  apiKeySaved = !!(await appPostKeySet());
}
let apiKeySaved = false;
function renderApiSummary(){
  const box = document.getElementById("api_sum");
  const btn = document.getElementById("api_edit");
  if(!box) return;
  box.innerHTML = "";
  const line = (k, v, off)=>{
    const r = document.createElement("div");
    r.className = "sumrow";
    const a = document.createElement("span");
    a.className = "sumk";
    a.textContent = k;
    const b = document.createElement("span");
    b.className = "sumv" + (off ? " off" : "");
    b.textContent = v;
    r.appendChild(a); r.appendChild(b);
    box.appendChild(r);
  };
  const model = document.getElementById("post_api_model").value.trim();
  const timeout = document.getElementById("post_api_timeout_s").value;
  if(postURL){
    line(L.apisumurl, postURL);
    line(L.apisummodel, model || L.apinomodel, !model);
    line(L.apisumkey, apiKeySaved ? L.postkeyset : L.postkeynone, !apiKeySaved);
    line(L.apisumtimeout, timeout + " s");
  } else {
    line(L.apisumstate, L.apinone, true);
  }
  if(btn) btn.textContent = postURL ? L.apiedit : L.apisetup;
}
async function editPostAPI(){
  const urlEl = document.getElementById("post_api_url");
  const modelEl = document.getElementById("post_api_model");
  const toEl = document.getElementById("post_api_timeout_s");
  const url = document.createElement("input");
  url.type = "text"; url.className = "apiurl"; url.value = postURL; url.placeholder = "https://api.openai.com/v1";
  const model = document.createElement("input");
  model.type = "text"; model.className = "apimodel"; model.value = modelEl.value; model.placeholder = "gpt-4.1-mini";
  const key = document.createElement("input");
  key.type = "password"; key.className = "apikey"; key.autocomplete = "off";
  const keyMask = "•".repeat(12);
  let keyTouched = false;
  key.value = apiKeySaved ? keyMask : "";
  key.addEventListener("input", ()=>{ keyTouched = true; });
  key.addEventListener("focus", ()=>{
    if(key.value !== keyMask) return;
    key.value = "";
    key.dispatchEvent(new Event("input", {bubbles: true}));
    keyTouched = false;
  });
  key.addEventListener("blur", ()=>{
    if(keyTouched || !apiKeySaved || key.value) return;
    key.value = keyMask;
    key.dispatchEvent(new Event("input", {bubbles: true}));
    keyTouched = false;
  });
  const timeout = document.createElement("select");
  timeout.className = "apitimeout";
  [...toEl.options].forEach(o=>{
    const c = document.createElement("option");
    c.value = o.value; c.textContent = o.textContent;
    timeout.appendChild(c);
  });
  timeout.value = toEl.value;
  let dropKey = false;
  const extras = [];
  if(apiKeySaved) extras.push({text: L.apikeydel, onClick: async ()=>{
    dropKey = true;
    const r = JSON.parse(await appSetPostKey(""));
    toast(r.message, r.severity);
    apiKeySaved = false;
    renderApiSummary();
  }});
  const ok = await formModal(L.fmsave, body=>{
    const sized = (el, w)=>{ el.style.flex = "0 0 " + w + "px"; el.style.width = w + "px"; return el; };
    body.appendChild(fmRow(sized(clearWrap(url), 300), L.postapiurl));
    body.appendChild(fmRow(sized(clearWrap(model), 300), L.postapimodel));
    body.appendChild(fmRow(sized(clearWrap(key), 300), L.postapikey, apiKeySaved ? L.postkeyset : ""));
    body.appendChild(fmRow(sized(timeout, 110), L.postapitimeout));
  }, extras, L.apidlg, true);
  if(!ok || dropKey) return;
  const fresh = url.value.trim();
  if(fresh && fresh !== postURL && !await askConfirm(L.postask.replace("%s", fresh), null, null, L.mtpost)) return;
  postURL = fresh;
  urlEl.value = fresh;
  modelEl.value = model.value.trim();
  toEl.value = timeout.value;
  if(key.value && key.value !== keyMask){
    const r = JSON.parse(await appSetPostKey(key.value));
    toast(r.message, r.severity);
    await updPostKeyState();
  }
  const want = fresh ? "api" : "local";
  if(postSource !== want){
    postSource = want;
    applyPostState();
    refreshLLM();
  }
  syncPostWarn();
  renderApiSummary();
  updatePostWarn();
  doSave();
}
function initPostAPI(){
  const el = document.getElementById("post_api_url");
  if(!el) return;
  el.value = postURL;
  syncPostWarn();
  updPostKeyState().then(renderApiSummary);
  const edit = document.getElementById("api_edit");
  if(edit) edit.onclick = e=>{ e.stopPropagation(); editPostAPI(); };
  const test = document.getElementById("api_test");
  if(test) test.onclick = e=>{
    e.stopPropagation();
    const to = document.getElementById("post_api_timeout_s");
    const md = document.getElementById("post_api_model");
    runPostTest(test, postURL, md ? md.value : "", to ? to.value : "", "");
  };
  const master = document.getElementById("post_enabled");
  master.checked = CFG.post_enabled !== false;
  master.addEventListener("change", e=>{ e.stopPropagation(); postEnabled = master.checked; applyPostState(); doSave(); });
  document.querySelectorAll("input.srcpick").forEach(sw=>{
    sw.addEventListener("change", e=>{
      e.stopPropagation();
      if(!sw.checked){ sw.checked = true; return; }
      postSource = sw.id === "pick_api" ? "api" : "local";
      applyPostState();
      doSave();
      refreshLLM();
    });
  });
  applyPostState();
}
function applyPostState(){
  updatePostWarn();
  ["post_model_card", "post_prompts_card"].forEach(id=>{
    const c = document.getElementById(id);
    if(c) c.classList.toggle("offdim", !postEnabled);
  });
  const mark = (id, pick, on)=>{
    const c = document.getElementById(id);
    if(c){
      c.classList.toggle("on", on);
      c.classList.toggle("idle", !on);
    }
    const sw = document.getElementById(pick);
    if(sw){
      sw.checked = on;
      sw.disabled = !postEnabled;
    }
  };
  mark("src_local", "pick_local", postSource !== "api");
  mark("src_api", "pick_api", postSource === "api");
}
function trUnavailRow(){
  if(!modelRowsCache.length) return null;
  const m = rowById(effectiveFor(curLang()));
  return m && !m.translate ? m : null;
}
function trTargetsFor(){
  if(!modelRowsCache.length) return null;
  const m = rowById(effectiveFor(curLang()));
  if(!m || !m.translate) return [];
  return m.trlangs ? m.trlangs.split(",") : ["en"];
}
function syncTrControls(){
  const sw = document.getElementById("tr_default");
  const bad = trUnavailRow();
  if(bad && sw.checked) sw.checked = false;
  sw.disabled = !!bad;
  const on = sw.checked && !bad;
  const mode = document.getElementById("translate_ask").value;
  const tsel = document.getElementById("translate_target");
  const allowed = trTargetsFor();
  if(allowed !== null){
    [...tsel.options].forEach(o=>{ o.disabled = !allowed.includes(o.value); });
    if(on && allowed.length && !allowed.includes(tsel.value)){
      tsel.value = allowed.includes("en") ? "en" : allowed[0];
      const forced = document.getElementById("tl_"+tsel.value);
      if(forced && !forced.checked) forced.checked = true;
    }
  }
  const tgt = tsel.value;
  tsel.disabled = !on;
  document.getElementById("translate_ask").disabled = !on;
  document.getElementById("translate_ask_seconds").disabled = !on || mode !== "timeout";
  trAll.forEach(l=>{
    const na = allowed !== null && !allowed.includes(l);
    document.getElementById("tl_"+l).disabled = !on || na || (mode === "never" && l !== tgt);
  });
  const un = document.getElementById("tr_unavail");
  if(un){
    un.style.display = bad ? "" : "none";
    if(bad) un.textContent = L.trunavail.replace("%s", bad.name);
  }
  const dimRow = (el, dim)=>{ const r = el && el.closest(".row"); if(r) r.classList.toggle("dimmed", dim); };
  dimRow(sw, !!bad);
  dimRow(document.getElementById("translate_target"), !on);
  dimRow(document.getElementById("translate_ask"), !on);
  dimRow(document.getElementById("translate_ask_seconds"), !on || mode !== "timeout");
  dimRow(document.getElementById("trlangs"), !on);
}
function askAlert(text, okText, title){
  return new Promise(resolve=>{
    const bg = document.createElement("div");
    bg.className = "modal-bg";
    const box = document.createElement("div");
    box.className = "modal";
    mTitle(box, title);
    const p = document.createElement("p");
    p.textContent = text;
    const row = document.createElement("div");
    row.className = "modal-btns";
    const yes = document.createElement("button");
    yes.type = "button";
    yes.className = "btn yes";
    yes.textContent = okText || L.ok;
    const done = ()=>{
      bg.remove();
      document.removeEventListener("keydown", onKey, true);
      unlockBg();
      resolve(true);
    };
    function onKey(e){
      if(e.key === "Escape" || e.key === "Enter"){ e.preventDefault(); done(); }
    }
    yes.onclick = done;
    bg.onclick = e=>{ if(e.target === bg) done(); };
    document.addEventListener("keydown", onKey, true);
    row.appendChild(yes);
    box.appendChild(p);
    box.appendChild(row);
    box.setAttribute("role", "dialog");
    box.setAttribute("aria-modal", "true");
    bg.appendChild(box);
    document.body.appendChild(bg);
    lockBg();
    yes.focus();
  });
}
function trLangName(code){
  const hit = presetLangs.find(p=>p[0]===code);
  return hit && hit[1] ? hit[1] : (code || "").toUpperCase();
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
  const f={hotkey:CFG.hotkey, lang_models: langModels,
    overlay_position: ovPos,
    overlay_monitor: (document.getElementById("overlay_monitor")||{}).value || "",
    overlay_custom: ovCustom,
    post_api_url: postURL,
    post_enabled: postEnabled,
    post_source: postSource,
    mic_device: micSel.value,
    punctuation: document.getElementById("punctuation").value,
    ui_level: "all",
    mic_device_name: micSel.value ? micSel.options[micSel.selectedIndex].textContent : "",
    whisper_prompt: dictWords.join(", "),
    translate_hotkey: translateHotkey,
    server_url: remoteURL,
    translate_ask_langs: trAll.filter(l=>document.getElementById("tl_"+l).checked),
    active_profiles: activeProfiles,
    replacements: repls,
    commands: cmds,
    llm_model_file: selLLM||"",
    profiles: profiles};
  bools.forEach(k=>f[k]=document.getElementById(k).checked);
  texts.forEach(k=>f[k]=document.getElementById(k).value);
  f.server_exe = exeUnlocked ? document.getElementById("server_exe").value.trim() : exeStored;
  nums.forEach(k=>f[k]=parseInt(document.getElementById(k).value)||0);
  sels.forEach(k=>f[k]=document.getElementById(k).value);
  const trSw = document.getElementById("tr_default");
  const trOnNow = trSw.checked && !trSw.disabled;
  f.translate_default = trOnNow && f.translate_ask === "never";
  if(!trOnNow) f.translate_ask = "never";
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
function bindLabels(){
  document.querySelectorAll(".row").forEach(row=>{
    const label = row.querySelector("label");
    if(!label || label.htmlFor) return;
    const field = row.querySelector('input[type=checkbox], input[type=text], input[type=number], select');
    if(field && field.id) label.htmlFor = field.id;
  });
}
let hits = [];
let hitAt = -1;
let searchFrom = null;
function searchMatches(s){
  const items = [...document.querySelectorAll(".page .row, .page .mrow, .page .sect, .page .blklbl, .page .hint, .page .mslot, .page .wizh, .about p, .about li")];
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
let wizOn = false, wizStep = 0, wizBase = 0, wizPlan = [], wizDlIds = [], wizSkippedDl = false;
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
  const rows = JSON.parse(await appModels());
  const m = rows.find(r=>r.id === "medium-q5_0");
  wizPlan = m ? [{id: m.id, name: m.name, size: m.size, installed: m.state !== "absent"}] : [];
  const box = wizEl("wiz_plan");
  box.innerHTML = "";
  wizPlan.forEach(p=>{
    const row = document.createElement("div");
    row.className = "advrow";
    row.innerHTML = '<span class="advrole">'+L.advprimary+'</span>'+
      '<span class="advname">'+esc(p.name)+'<span class="advwhy">'+L.advrolemain+'</span></span>'+
      '<span class="advstate'+(p.installed ? " ok" : "")+'">'+(p.installed ? L.advhave : (p.size ? p.size+" MB" : ""))+'</span>';
    box.appendChild(row);
  });
  const missing = wizPlan.filter(p=>!p.installed);
  const btn = wizEl("wiz_dl");
  const out = wizEl("wiz_dlout");
  btn.style.display = missing.length ? "" : "none";
  btn.textContent = L.dl + (missing.length && missing[0].size ? " · " + missing[0].size + " MB" : "");
  wizEl("wiz_dl_skip").style.display = missing.length ? "" : "none";
  wizEl("wiz_skipnote").style.display = wizSkippedDl ? "" : "none";
  wizEl("wiz_dlrow").style.display = "none";
  out.textContent = missing.length ? "" : L.wizhave;
  out.classList.toggle("ok", !missing.length);
  wizDlIds = [];
  wizSyncNav();
}
async function wizApplyModel(){
  refreshModels();
  wizSyncNav();
}
function wizSyncNav(){
  const next = wizEl("wiz_next");
  if(!next) return;
  const waiting = wizStep === 1 && !wizSkippedDl && (wizDlIds.length > 0 || !wizPlan.some(p=>p.installed));
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
  wizEl("wiz_dl_skip").onclick = ()=>{
    wizSkippedDl = true;
    wizEl("wiz_skipnote").style.display = "";
    wizSyncNav();
  };
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
      if(!await askConfirm(L.histask, L.histclear, null, L.mthist)) return;
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
function formModal(okText, build, extra, title, wide){
  return new Promise(resolve=>{
    const bg = document.createElement("div");
    bg.className = "modal-bg";
    const box = document.createElement("div");
    const extras = Array.isArray(extra) ? extra : (extra ? [extra] : []);
    box.className = "modal formmodal" + (extras.length || wide ? " wide" : "");
    mTitle(box, title);
    const body = document.createElement("div");
    box.appendChild(body);
    const row = document.createElement("div");
    row.className = "modal-btns";
    const no = document.createElement("button");
    no.type = "button"; no.className = "btn ghost"; no.textContent = L.cancel;
    const yes = document.createElement("button");
    yes.type = "button"; yes.className = "btn yes"; yes.textContent = okText;
    const done = v=>{
      bg.remove();
      document.removeEventListener("keydown", onKey, true);
      unlockBg();
      resolve(v);
    };
    function onKey(e){
      if(e.key === "Escape"){ e.preventDefault(); done(false); }
      else if(e.key === "Enter" && e.target.tagName !== "SELECT"){ e.preventDefault(); done(true); }
    }
    yes.onclick = ()=>done(true);
    no.onclick = ()=>done(false);
    bg.onclick = e=>{ if(e.target === bg) done(false); };
    document.addEventListener("keydown", onKey, true);
    row.appendChild(no);
    extras.forEach(x=>{
      const ex = document.createElement("button");
      ex.type = "button"; ex.className = "btn ghost"; ex.textContent = x.text;
      ex.id = x.id || (row.querySelector("#fm_extra") ? "" : "fm_extra");
      ex.onclick = ()=>{ x.onClick(ex); if(!x.keep) done(false); };
      row.appendChild(ex);
    });
    row.appendChild(yes);
    box.appendChild(row);
    box.setAttribute("role", "dialog");
    box.setAttribute("aria-modal", "true");
    bg.appendChild(box);
    document.body.appendChild(bg);
    lockBg();
    build(body);
    const f = body.querySelector('input[type=text]') || body.querySelector("input,select");
    if(f) f.focus();
  });
}
function clearWrap(inp){
  const w = document.createElement("span");
  w.className = "clearwrap";
  w.appendChild(inp);
  const x = document.createElement("button");
  x.type = "button"; x.className = "clearx"; x.tabIndex = -1; x.title = L.del; x.textContent = "✕";
  const upd = ()=>{ x.style.display = inp.value ? "" : "none"; };
  inp.addEventListener("input", upd);
  x.onclick = ()=>{ inp.value = ""; upd(); inp.focus(); inp.dispatchEvent(new Event("input", {bubbles: true})); };
  upd();
  w.appendChild(x);
  return w;
}
function fmRow(ctl, labelText, tip){
  const r = document.createElement("div");
  r.className = "fmrow";
  if(labelText){
    const l = document.createElement("label");
    l.textContent = labelText;
    if(tip){
      const s = document.createElement("span");
      s.className = "sub";
      s.textContent = tip;
      l.appendChild(s);
    }
    r.appendChild(l);
  }
  r.appendChild(ctl);
  return r;
}
function pairRow(sum, meta, onEdit, onDel){
  const row = document.createElement("div");
  row.className = "pairrow";
  const s = document.createElement("span");
  s.className = "psum";
  s.textContent = sum;
  row.appendChild(s);
  const m = document.createElement("span");
  m.className = "pmeta2";
  m.textContent = meta;
  row.appendChild(m);
  const ed = document.createElement("button");
  ed.type = "button"; ed.className = "iconbtn redit"; ed.title = L.pedit; ed.innerHTML = "&#9998;";
  ed.onclick = onEdit;
  row.appendChild(ed);
  const del = document.createElement("button");
  del.type = "button"; del.className = "iconbtn danger rdel"; del.title = L.del; del.textContent = "✕";
  del.onclick = onDel;
  row.appendChild(del);
  return row;
}
let replFilterVal = ()=>"", cmdFilterVal = ()=>"";
function wireFilter(inputId, clearId, onChange){
  const inp = document.getElementById(inputId);
  const clr = document.getElementById(clearId);
  if(!inp) return ()=>"" ;
  const upd = ()=>{ if(clr) clr.style.display = inp.value ? "" : "none"; };
  inp.addEventListener("input", ()=>{ upd(); onChange(); });
  if(clr) clr.onclick = ()=>{ inp.value = ""; upd(); inp.focus(); onChange(); };
  upd();
  return ()=>inp.value.trim().toLowerCase();
}
function listChip(text, meta, onEdit, onDel){
  const chip = document.createElement("span");
  chip.className = "cmdchip";
  const t = document.createElement("span");
  t.textContent = text;
  chip.appendChild(t);
  if(meta){
    const m = document.createElement("span");
    m.className = "cpmeta";
    m.textContent = "/ " + meta;
    chip.appendChild(m);
  }
  if(onEdit){
    const ed = document.createElement("button");
    ed.type = "button"; ed.className = "cbtn redit"; ed.title = L.pedit; ed.innerHTML = "&#9998;";
    ed.onclick = onEdit;
    chip.appendChild(ed);
  }
  const del = document.createElement("button");
  del.type = "button"; del.className = "cbtn rdel"; del.title = L.del; del.textContent = "✕";
  del.onclick = onDel;
  chip.appendChild(del);
  return chip;
}
function cmdActionLabel(a){
  return a === "paragraph" ? L.cmdparagraph : a === "text" ? L.cmdtext : a === "cancel" ? L.cmdcancel : L.cmdnewline;
}
async function editCmd(i){
  const c = i >= 0 ? {...cmds[i]} : {id: "c" + Date.now(), phrase: "", action: "newline", text: ""};
  const phrase = document.createElement("input");
  phrase.type = "text"; phrase.className = "cphrase"; phrase.value = c.phrase || ""; phrase.placeholder = L.cmdph;
  const act = document.createElement("select");
  act.className = "caction";
  ruleOpts(act, [["newline", L.cmdnewline], ["paragraph", L.cmdparagraph], ["text", L.cmdtext], ["cancel", L.cmdcancel]], c.action || "newline");
  const txt = document.createElement("input");
  txt.type = "text"; txt.className = "ctext"; txt.value = c.text || ""; txt.placeholder = L.cmdtextph;
  const txtRow = fmRow(clearWrap(txt));
  txtRow.style.display = (c.action === "text") ? "" : "none";
  act.onchange = ()=>{ txtRow.style.display = (act.value === "text") ? "" : "none"; };
  const ok = await formModal(i >= 0 ? L.fmsave : L.fmadd, body=>{
    body.appendChild(fmRow(act, L.cmdaction, L.tipcmdaction));
    body.appendChild(fmRow(clearWrap(phrase)));
    body.appendChild(txtRow);
  }, i >= 0 ? null : {text: L.cmdpreset, onClick: addPresetCmds}, i >= 0 ? L.fmtcmdedit : L.fmtcmdadd);
  if(!ok) return;
  c.phrase = phrase.value.trim();
  c.action = act.value;
  c.text = act.value === "text" ? txt.value : "";
  if(!c.phrase) return;
  if(i >= 0) cmds[i] = c; else cmds.push(c);
  renderCmds();
  applyNow();
}
function listEmpty(body, inputId, q, noneText){
  const empty = document.createElement("div");
  empty.className = "ruleempty";
  const raw = ((document.getElementById(inputId) || {}).value || "").trim();
  empty.textContent = q ? L.listnothing.replace("%s", raw) : noneText;
  body.appendChild(empty);
}
function renderCmds(){
  const body = document.getElementById("cmdbody");
  if(!body) return;
  body.innerHTML = "";
  const q = cmdFilterVal();
  let shown = 0;
  cmds.forEach((c, i)=>{
    const meta = cmdActionLabel(c.action) + (c.action === "text" && c.text ? " · " + c.text : "");
    if(q && !((c.phrase || "") + " " + meta).toLowerCase().includes(q)) return;
    shown++;
    body.appendChild(listChip(c.phrase || "—", meta,
      ()=>{ editCmd(i); },
      ()=>{ cmds.splice(i, 1); renderCmds(); applyNow(); }));
  });
  if(!shown) listEmpty(body, "cmd_filter", q, L.cmdempty);
}
function addPresetCmds(){
  const known = cmds.map(c=>(c.phrase || "").toLowerCase());
  [["newline", L.cmdpnewline], ["paragraph", L.cmdpparagraph], ["cancel", L.cmdpcancel]].forEach(([action, phrase], n)=>{
    if(!phrase || known.includes(phrase.toLowerCase())) return;
    cmds.push({id: "c" + Date.now() + n, phrase: phrase, action: action, text: ""});
  });
  renderCmds();
  applyNow();
}
function initCmds(){
  const add = document.getElementById("cmd_add");
  if(!add) return;
  add.onclick = ()=>{ editCmd(-1); };
  cmdFilterVal = wireFilter("cmd_filter", "cmd_filter_clear", renderCmds);
  renderCmds();
}
let repls = (CFG.replacements || []).map(r=>({...r}));
async function editRepl(i){
  const r = i >= 0 ? {...repls[i]} : {id: "x" + Date.now(), from: "", to: "", whole: true, match_case: false, lang: ""};
  const from = document.createElement("input");
  from.type = "text"; from.className = "rfrom"; from.value = r.from || ""; from.placeholder = L.replfromph;
  const to = document.createElement("input");
  to.type = "text"; to.className = "rto"; to.value = r.to || ""; to.placeholder = L.repltoph;
  const langSel = document.createElement("select");
  langSel.className = "rlang";
  langSel.title = L.repllang;
  [["", L.repllangall], ["de", trLangName("de")], ["en", trLangName("en")], ["es", trLangName("es")], ["fr", trLangName("fr")], ["it", trLangName("it")], ["pl", trLangName("pl")], ["uk", trLangName("uk")], ["ru", trLangName("ru")]].forEach(([v, t])=>{
    const o = document.createElement("option");
    o.value = v; o.textContent = t;
    langSel.appendChild(o);
  });
  langSel.value = [...langSel.options].some(o=>o.value===(r.lang||"")) ? (r.lang||"") : "";
  const whole = document.createElement("input");
  whole.type = "checkbox"; whole.className = "rwhole"; whole.checked = r.whole !== false;
  const mcase = document.createElement("input");
  mcase.type = "checkbox"; mcase.className = "rcase"; mcase.checked = !!r.match_case;
  const ok = await formModal(i >= 0 ? L.fmsave : L.fmadd, body=>{
    body.appendChild(fmRow(clearWrap(from)));
    body.appendChild(fmRow(clearWrap(to)));
    body.appendChild(fmRow(langSel, L.repllang, L.tiprepllang));
    body.appendChild(fmRow(mcase, L.replcasefull, L.tipreplcase));
    body.appendChild(fmRow(whole, L.replwholefull, L.tipreplwhole));
  }, null, i >= 0 ? L.fmtrepledit : L.fmtrepladd);
  if(!ok) return;
  r.from = from.value;
  r.to = to.value;
  r.lang = langSel.value;
  r.whole = whole.checked;
  r.match_case = mcase.checked;
  if(!r.from.trim()) return;
  if(i >= 0) repls[i] = r; else repls.push(r);
  renderRepls();
  applyNow();
}
function renderRepls(){
  const body = document.getElementById("replbody");
  if(!body) return;
  body.innerHTML = "";
  const q = replFilterVal();
  let shown = 0;
  repls.forEach((r, i)=>{
    const text = (r.from || "—") + " → " + (r.to || "");
    const meta = r.lang ? r.lang.toUpperCase() : L.repllangall;
    if(q && !(text + " " + meta).toLowerCase().includes(q)) return;
    shown++;
    body.appendChild(listChip(text, meta,
      ()=>{ editRepl(i); },
      ()=>{ repls.splice(i, 1); renderRepls(); applyNow(); }));
  });
  if(!shown) listEmpty(body, "repl_filter", q, L.replempty);
}
let dictWords = (CFG.whisper_prompt || "").split(",").map(s=>s.trim()).filter(Boolean);
let dictFilterVal = ()=>"";
function renderDict(){
  const body = document.getElementById("dictbody");
  if(!body) return;
  body.innerHTML = "";
  const q = dictFilterVal();
  let shown = 0;
  dictWords.forEach((word, i)=>{
    if(q && !word.toLowerCase().includes(q)) return;
    shown++;
    body.appendChild(listChip(word, "",
      null,
      ()=>{ dictWords.splice(i, 1); renderDict(); doSave(); }));
  });
  if(!shown) listEmpty(body, "dict_filter", q, L.dictempty);
}
function updateDictWarn(){
  const warn = document.getElementById("dict_warn");
  if(!warn) return;
  const row = rowById(activeModelId);
  const bad = row && row.engine !== "whisper";
  warn.style.display = bad ? "" : "none";
  if(bad) warn.textContent = L.dictnomodel.replace("%s", row.name);
}
async function addDictWords(){
  const inp = document.createElement("input");
  inp.type = "text"; inp.className = "dword"; inp.placeholder = L.dictaddph;
  const ok = await formModal(L.fmadd, body=>{
    body.appendChild(fmRow(clearWrap(inp)));
  }, null, L.fmtdictadd);
  if(!ok) return;
  const fresh = inp.value.split(",").map(s=>s.trim()).filter(Boolean);
  const known = dictWords.map(w=>w.toLowerCase());
  fresh.forEach(w=>{
    if(known.includes(w.toLowerCase())) return;
    dictWords.push(w);
    known.push(w.toLowerCase());
  });
  renderDict();
  doSave();
}
function initDict(){
  const add = document.getElementById("dict_add");
  if(!add) return;
  add.onclick = addDictWords;
  dictFilterVal = wireFilter("dict_filter", "dict_filter_clear", renderDict);
  renderDict();
  updateDictWarn();
}
function initRepls(){
  const add = document.getElementById("repl_add");
  if(!add) return;
  add.onclick = ()=>{ editRepl(-1); };
  replFilterVal = wireFilter("repl_filter", "repl_filter_clear", renderRepls);
  renderRepls();
}
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
let ovPos = CFG.overlay_position || "bottom";
let ovCustom = Object.assign({}, CFG.overlay_custom || {});
const OV_ANCHORS = {"top-left":[.12,.15],"top":[.5,.15],"top-right":[.88,.15],"left":[.12,.5],"center":[.5,.5],"right":[.88,.5],"bottom-left":[.12,.85],"bottom":[.5,.85],"bottom-right":[.88,.85]};
const OV_CELLS = {"top-left":[0,0],"top":[1,0],"top-right":[2,0],"left":[0,1],"center":[1,1],"right":[2,1],"bottom-left":[0,2],"bottom":[1,2],"bottom-right":[2,2]};
let ovAnchorLast = OV_ANCHORS[CFG.overlay_position] ? CFG.overlay_position : "bottom";
function ovMonKey(){
  const mons = CFG._monitors || [];
  const sel = document.getElementById("overlay_monitor");
  const v = sel ? sel.value : "";
  if(v !== "" && v !== "cursor"){
    const m = mons[parseInt(v)];
    if(m) return m.w + "x" + m.h;
  }
  const prim = mons.find(m=>m.primary) || mons[0];
  return prim ? prim.w + "x" + prim.h : "1920x1080";
}
function syncOverlayCard(){
  const master = document.getElementById("overlay");
  const card = document.getElementById("ovcard");
  if(!master || !card) return;
  const on = master.checked;
  [...card.querySelectorAll(".row")].slice(1).forEach(r=>{
    r.classList.toggle("dimmed", !on);
    r.querySelectorAll("input,select,button").forEach(el=>{ el.disabled = !on; });
  });
  paintOvScheme();
}
function paintOvScheme(){
  const box = document.getElementById("ovscheme");
  const mini = document.getElementById("ovmini");
  if(!box || !mini) return;
  const master = document.getElementById("overlay");
  const shown = !master || master.checked;
  const store = document.getElementById("overlay_position");
  if(store) store.value = ovPos;
  const caret = document.getElementById("ov_caret");
  if(caret) caret.checked = ovPos === "caret";
  box.classList.toggle("off", ovPos === "caret" || !shown);
  const free = ovPos === "custom";
  const crt = document.getElementById("ovcrt");
  if(crt) crt.classList.toggle("free", free);
  const freeBox = document.getElementById("ov_free");
  if(freeBox) freeBox.checked = free;
  const sub = document.getElementById("ovpos_sub");
  if(sub) sub.textContent = free ? L.ovposdrag : L.ovposscheme;
  box.querySelectorAll(".ovzone").forEach(z=>{
    const cell = OV_CELLS[z.dataset.pos];
    z.style.left = (cell[0] * 33.333) + "%";
    z.style.top = (cell[1] * 33.333) + "%";
  });
  let f = OV_ANCHORS[ovPos] || [.5,.85];
  if(free){
    const c = ovCustom[ovMonKey()];
    if(c) f = [c.x, c.y];
  }
  mini.style.left = (f[0]*100) + "%";
  mini.style.top = (f[1]*100) + "%";
  mini.style.display = ovPos === "caret" ? "none" : "";
}
function initOverlayScheme(){
  const box = document.getElementById("ovscheme");
  const mini = document.getElementById("ovmini");
  const caret = document.getElementById("ov_caret");
  const monRow = document.getElementById("ovmon_row");
  const monSel = document.getElementById("overlay_monitor");
  if(!box || !mini || !caret || !monSel) return;
  const mons = CFG._monitors || [];
  const o = document.createElement("option");
  o.value = ""; o.textContent = L.ovmoncursor;
  monSel.appendChild(o);
  mons.forEach((m, i)=>{
    const op = document.createElement("option");
    op.value = String(i);
    const label = m.name || L.ovmonn.replace("%d", i + 1);
    op.textContent = label + " (" + m.w + "×" + m.h + ")" + (m.primary ? " ★" : "");
    monSel.appendChild(op);
  });
  monSel.value = [...monSel.options].some(x=>x.value===(CFG.overlay_monitor||"")) ? (CFG.overlay_monitor||"") : "";
  if(monRow && mons.length > 1) monRow.style.display = "";
  monSel.addEventListener("change", paintOvScheme);
  box.querySelectorAll(".ovzone").forEach(z=>{
    z.onclick = ()=>{ ovPos = z.dataset.pos; ovAnchorLast = ovPos; paintOvScheme(); doSave(); };
  });
  caret.onchange = e=>{ e.stopPropagation(); ovPos = caret.checked ? "caret" : ovAnchorLast; paintOvScheme(); doSave(); };
  const freeBox = document.getElementById("ov_free");
  if(freeBox) freeBox.onchange = e=>{
    e.stopPropagation();
    if(freeBox.checked){
      const key = ovMonKey();
      if(!ovCustom[key]){
        const a = OV_ANCHORS[ovAnchorLast] || [.5,.85];
        ovCustom[key] = {x: a[0], y: a[1]};
      }
      ovPos = "custom";
    } else {
      ovPos = ovAnchorLast;
    }
    paintOvScheme();
    doSave();
  };
  const crt = document.getElementById("ovcrt");
  let dragging = false;
  const dragTo = e=>{
    const r = crt.getBoundingClientRect();
    if(!r.width || !r.height) return;
    const fx = Math.max(.07, Math.min(.93, (e.clientX - r.left) / r.width));
    const fy = Math.max(.07, Math.min(.93, (e.clientY - r.top) / r.height));
    ovCustom[ovMonKey()] = {x: +fx.toFixed(3), y: +fy.toFixed(3)};
    paintOvScheme();
  };
  if(crt) crt.addEventListener("pointerdown", e=>{
    if(ovPos !== "custom") return;
    dragging = true;
    if(crt.setPointerCapture && e.pointerId !== undefined) crt.setPointerCapture(e.pointerId);
    e.preventDefault();
    dragTo(e);
  });
  if(crt) crt.addEventListener("pointermove", e=>{
    if(!dragging) return;
    dragTo(e);
  });
  if(crt) crt.addEventListener("pointerup", ()=>{
    if(!dragging) return;
    dragging = false;
    doSave();
  });
  const master = document.getElementById("overlay");
  if(master) master.addEventListener("change", syncOverlayCard);
  syncOverlayCard();
}
function renderSkip(){
  const box = document.getElementById("hist_skip_list");
  const store = document.getElementById("history_skip");
  if(!box || !store) return;
  const items = store.value.split(",").map(s=>s.trim()).filter(Boolean);
  box.innerHTML = "";
  items.forEach((p, i)=>{
    const chip = document.createElement("span");
    chip.className = "skipchip";
    chip.textContent = p;
    const x = document.createElement("button");
    x.type = "button"; x.className = "chipx"; x.textContent = "✕";
    x.onclick = ()=>{ items.splice(i, 1); store.value = items.join(", "); renderSkip(); doSave(); };
    chip.appendChild(x);
    box.appendChild(chip);
  });
}
function initHistSkip(){
  const add = document.getElementById("hist_skip_add");
  const inp = document.getElementById("hist_skip_new");
  const store = document.getElementById("history_skip");
  if(!add || !inp || !store) return;
  const commit = ()=>{
    const v = inp.value.trim();
    if(!v) return;
    const items = store.value.split(",").map(s=>s.trim()).filter(Boolean);
    if(!items.includes(v)) items.push(v);
    store.value = items.join(", ");
    inp.value = "";
    renderSkip();
    doSave();
  };
  add.onclick = commit;
  inp.addEventListener("keydown", e=>{ if(e.key === "Enter"){ e.preventDefault(); commit(); } });
  renderSkip();
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
  if(!await askConfirm(L.resetask, L.resetbtn, null, L.mtreset)) return;
  await appResetSettings();
  appReload(curTab);
};
let applyTimer = null;
function applyNow(){
  clearTimeout(applyTimer);
  applyTimer = setTimeout(()=>{ doSave(); }, 120);
}
document.querySelector(".content").addEventListener("change", e=>{
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
  if(e.target.closest("button, input, select, textarea, a, .omni")) return;
  if(e.button===0) appDrag();
});
load();
(async ()=>{
  initStateScreen();
  initRemote();
  initPostAPI();
  initWizard();
  initAutorun();
  initHistSkip();
  initOverlayScheme();
  initRepls();
  initCmds();
  initHistory();
  bindLabels();
  ariaFromTitle(document);
  initTips();
  initWindowButtons();
  labelPages();
  buildToc();
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
  const llmBtn = document.getElementById("llm_catalog");
  if(llmBtn) llmBtn.onclick = e=>{ e.stopPropagation(); openLLMCatalog(); };
  const profAdd = document.getElementById("profadd");
  if(profAdd) profAdd.onclick = e=>{
    e.stopPropagation();
    editPromptDialog({id: "p" + Date.now(), name: "", prompt: ""}, true);
  };
  await refreshModels();
  await refreshLLM();
  if(window.appReady) appReady();
})();
const tabAlias = {general:"state", rec:"models", proc:"text", server:"system", about:"about", state:"state", dictation:"dictation", history:"history", mic:"mic", models:"models", text:"text", translate:"dictation", system:"system", post:"post", help:"help", contacts:"contacts"};
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
