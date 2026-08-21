package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"

	webview "github.com/jchv/go-webview2"
	"golang.org/x/sys/windows"
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
	}()
	return webview.NewWithOptions(webview.WebViewOptions{
		DataPath:  filepath.Join(os.TempDir(), fmt.Sprintf("voxterminal-webview-%d", os.Getpid())),
		AutoFocus: true,
		WindowOptions: webview.WindowOptions{
			Title:  settingsStrings[lang()]["S_TITLE"],
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
	own := fmt.Sprintf("voxterminal-webview-%d", os.Getpid())
	entries, err := os.ReadDir(os.TempDir())
	if err != nil {
		return
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() && strings.HasPrefix(name, "voxterminal-webview") && name != own {
			_ = os.RemoveAll(filepath.Join(os.TempDir(), name))
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
	Animation        bool   `json:"animation"`
	TypeMode         bool   `json:"type_mode"`
	Language         string `json:"language"`
	ModelID          string `json:"model_id"`
	Threads          int    `json:"threads"`
	MinRecordMs      int    `json:"min_record_ms"`
	MaxRecordSeconds int    `json:"max_record_seconds"`
	ServerAutostart  bool   `json:"server_autostart"`
	CheckUpdates     bool   `json:"check_updates"`
	ServerPort       int    `json:"server_port"`
	ServerExe        string `json:"server_exe"`
	ServerURL        string `json:"server_url"`

	WhisperPrompt       string    `json:"whisper_prompt"`
	TranslateHotkey     string    `json:"translate_hotkey"`
	TranslateTarget     string    `json:"translate_target"`
	TranslateAsk        string    `json:"translate_ask"`
	TranslateAskSeconds int       `json:"translate_ask_seconds"`
	TranslateAskLangs   []string  `json:"translate_ask_langs"`
	TranslateDefault    bool      `json:"translate_default"`
	ActiveProfiles      []string  `json:"active_profiles"`
	LLMModelFile        string    `json:"llm_model_file"`
	Profiles            []Profile `json:"profiles"`
}

func (a *App) openSettings(tab string) {
	log.Printf("openSettings: tab=%s", tab)
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
		winW, winH := 640, 560
		if c := a.snapshot(); c.SettingsW >= 500 && c.SettingsH >= 400 {
			winW, winH = c.SettingsW, c.SettingsH
		}
		lastWndW, lastWndH = 0, 0
		w := createWebView(winW, winH)
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
		_ = w.Bind("appClose", func() {
			procPostMessageW.Call(hwnd, wmClose, 0, 0)
		})
		_ = w.Bind("appDrag", func() {
			beginWindowDrag(hwnd)
		})
		_ = w.Bind("appSave", func(formJSON string) string {
			var f settingsForm
			if err := json.Unmarshal([]byte(formJSON), &f); err != nil {
				return err.Error()
			}
			return a.applySettings(&f)
		})
		_ = w.Bind("appCapture", func() {
			go func() {
				a.changeHotkey()
				combo := a.snapshot().Hotkey
				w.Dispatch(func() {
					w.Eval(fmt.Sprintf("setHotkey(%q)", combo))
				})
			}()
		})
		_ = w.Bind("appModels", func() string {
			return a.modelRows()
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
		_ = w.Bind("appLLMDlFile", func(repo, file string) {
			a.llmDownloadFile(repo, file)
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
			tag, uurl, err := fetchLatestRelease()
			if err != nil {
				out, _ := json.Marshal(map[string]any{"error": err.Error()})
				return string(out)
			}
			newer := verNewer(tag, appVersion) && uurl != ""
			if newer {
				a.mu.Lock()
				a.updVer, a.updURL = tag, uurl
				a.mu.Unlock()
			}
			out, _ := json.Marshal(map[string]any{"current": appVersion, "latest": tag, "newer": newer})
			return string(out)
		})
		_ = w.Bind("appDoUpdate", func() {
			a.mu.Lock()
			uurl := a.updURL
			a.mu.Unlock()
			if uurl == "" {
				return
			}
			go func() {
				path, err := downloadSetup(uurl, func(pct int) {
					w.Dispatch(func() { w.Eval(fmt.Sprintf("updProgress(%d)", pct)) })
				})
				if err == nil {
					err = launchUpdater(path)
				}
				if err != nil {
					enc, _ := json.Marshal(err.Error())
					w.Dispatch(func() { w.Eval("updError(" + string(enc) + ")") })
					return
				}
				a.quitForUpdate()
			}()
		})
		_ = w.Bind("appRepoLink", func() {
			shellOpenURL("https://github.com/Vitalii-Yemets/vox-terminal")
		})
		_ = w.Bind("appAuthorLink", func() {
			shellOpenURL("https://github.com/Vitalii-Yemets")
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
				w.Dispatch(func() {
					w.Eval(fmt.Sprintf("comboCaptured(%q)", combo))
				})
			}()
		})
		_ = w.Bind("appLLMTest", func(prompt, sample string) {
			go func() {
				out, err := a.llmProcess(context.Background(), prompt, sample)
				if err != nil {
					out = "⚠ " + err.Error()
				}
				enc, _ := json.Marshal(out)
				w.Dispatch(func() {
					w.Eval("llmTestResult(" + string(enc) + ")")
				})
			}()
		})
		_ = w.Bind("appReload", func(tabName string) {
			w.Dispatch(func() {
				w.Navigate("data:text/html;charset=utf-8," + url.PathEscape(settingsHTML(a.snapshot(), tabName)))
			})
		})
		_ = w.Bind("appModelDel", func(id string) string {
			return a.deleteModel(id)
		})

		log.Printf("openSettings: WebView2 создан, открываю страницу")
		settingsHwnd.Store(hwnd)
		defer settingsHwnd.Store(0)
		applyDarkCaption(hwnd)
		makeBorderless(hwnd)
		applyMinSize(hwnd, 500, 400)
		procShowWindow.Call(hwnd, 5)
		procSetForegroundWnd.Call(hwnd)
		w.Navigate("data:text/html;charset=utf-8," + url.PathEscape(settingsHTML(a.snapshot(), tab)))
		w.Run()
		log.Printf("openSettings: окно закрыто")
		if lastWndW >= 500 && lastWndH >= 400 {
			nw, nh := int(lastWndW), int(lastWndH)
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

func (a *App) applySettings(f *settingsForm) string {
	if _, err := parseHotkey(f.Hotkey); err != nil {
		return err.Error()
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
	c.Animation = f.Animation
	if f.TypeMode {
		c.PasteMode = "type"
	} else {
		c.PasteMode = "clipboard"
	}
	c.Language = f.Language
	modelChanged := false
	if f.ModelID != "" && f.ModelID != "custom" {
		if m := findModel(f.ModelID); m != nil {
			if _, err := os.Stat(filepath.Join("models", m.File)); err == nil {
				nm := "models/" + m.File
				if nm != c.Model {
					c.Model = nm
					modelChanged = true
				}
			}
		}
	}
	if f.Threads > 0 {
		c.Threads = f.Threads
	}
	if f.MinRecordMs >= 0 {
		c.MinRecordMs = f.MinRecordMs
	}
	if f.MaxRecordSeconds > 0 {
		c.MaxRecordSeconds = f.MaxRecordSeconds
	}
	c.ServerAutostart = f.ServerAutostart
	c.CheckUpdates = f.CheckUpdates
	if f.ServerPort > 0 {
		c.ServerPort = f.ServerPort
	}
	c.ServerExe = f.ServerExe
	c.ServerURL = strings.TrimSpace(f.ServerURL)
	c.WhisperPrompt = strings.TrimSpace(f.WhisperPrompt)
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
		if len(langs) > 0 {
			c.TranslateAskLangs = langs
		}
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
	a.cfg = &c
	hook := a.hook
	a.mu.Unlock()
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
	if err := saveConfig("config.json", &c); err != nil {
		log.Printf("сохранение конфига: %v", err)
		return err.Error()
	}
	a.refreshIdleUI()
	if modelChanged {
		a.requestServerRestart()
	}

	restartNeeded := c.ServerPort != old.ServerPort ||
		c.Threads != old.Threads ||
		c.ServerURL != old.ServerURL ||
		c.ServerExe != old.ServerExe ||
		c.ServerAutostart != old.ServerAutostart
	log.Printf("настройки сохранены: hotkey=%s ui=%s model=%s restart=%v", c.Hotkey, c.UILanguage, c.Model, restartNeeded)
	switch {
	case restartNeeded:
		return settingsStrings[lang()]["S_RESTART"]
	case modelChanged:
		return tr("model.switching")
	default:
		return settingsStrings[lang()]["S_SAVED"]
	}
}

func settingsHTML(cfg *Config, tab string) string {
	cfgMap := map[string]any{
		"hotkey":             cfg.Hotkey,
		"ui_language":        cfg.UILanguage,
		"beep":               cfg.Beep,
		"sound_theme":        cfg.SoundTheme,
		"auto_enter":         cfg.AutoEnter,
		"restore_clipboard":  cfg.RestoreClipboard,
		"overlay":            cfg.Overlay,
		"animation":          cfg.Animation,
		"type_mode":          cfg.PasteMode == "type",
		"language":           cfg.Language,
		"threads":            cfg.Threads,
		"min_record_ms":      cfg.MinRecordMs,
		"max_record_seconds": cfg.MaxRecordSeconds,
		"server_autostart":   cfg.ServerAutostart,
		"check_updates":      cfg.CheckUpdates,
		"server_port":        cfg.ServerPort,
		"server_exe":         cfg.ServerExe,
		"server_url":         cfg.ServerURL,
		"whisper_prompt":     cfg.WhisperPrompt,
		"translate_default":  cfg.TranslateDefault,
		"active_profiles":    cfg.ActiveProfiles,
		"translate_hotkey":   cfg.TranslateHotkey,
		"translate_target":   cfg.TranslateTarget,
		"translate_ask":      cfg.TranslateAsk,
		"translate_ask_seconds": cfg.TranslateAskSeconds,
		"translate_ask_langs":   cfg.TranslateAskLangs,
		"llm_model":          filepath.Base(cfg.LLMModel),
		"profiles":           cfg.Profiles,
		"_version":           appVersion,
		"_tab":               tab,
		"_cpus":              runtime.NumCPU(),
	}
	cfgJSON, _ := json.Marshal(cfgMap)

	pairs := []string{"{{CFG}}", string(cfgJSON)}
	cur := settingsStrings[lang()]
	for k, enV := range settingsStrings["en"] {
		v := cur[k]
		if v == "" {
			v = enV
		}
		pairs = append(pairs, "{{"+k+"}}", v)
	}
	return strings.NewReplacer(pairs...).Replace(settingsPage)
}

const settingsPage = `<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>{{S_TITLE}}</title><style>
:root{--bg:#0b0f0c;--panel:#0e1410;--line:#1d4a2b;--green:#3cff6e;--dim:#20a34a;--faint:#14803a;--amber:#ffb347;--glow:0 0 7px rgba(60,255,110,.55)}
*{box-sizing:border-box;margin:0;padding:0}
html,body{height:100%}
body{font:14px Consolas,"Cascadia Mono",monospace;background:var(--bg);color:var(--green);user-select:none;display:flex;flex-direction:column;overflow:hidden}
body::after{content:"";position:fixed;inset:0;pointer-events:none;background:repeating-linear-gradient(transparent 0 2px,rgba(0,0,0,.18) 2px 3px)}
.content{flex:1;overflow-y:auto;overflow-x:hidden;min-height:0}
::-webkit-scrollbar{width:10px}
::-webkit-scrollbar-track{background:var(--bg)}
::-webkit-scrollbar-thumb{background:var(--line);border:2px solid var(--bg)}
::-webkit-scrollbar-thumb:hover{background:var(--dim)}
.header{display:flex;align-items:center;gap:14px;padding:12px 12px 12px 20px;border-bottom:1px solid var(--line);box-shadow:0 1px 12px rgba(60,255,110,.12);cursor:default}
.capbtns{display:flex;gap:6px;margin-left:10px}
button.cap{width:36px;height:30px;background:none;border:1px solid var(--line);color:var(--dim);font:14px Consolas,monospace;cursor:pointer;padding:0}
button.cap:hover{background:#123f22;color:var(--green);box-shadow:var(--glow)}
button.cap.close:hover{background:#3c1212;color:#ff7b6b;border-color:#7a2e2e;box-shadow:0 0 7px rgba(255,110,90,.5)}
.logo{width:40px;height:40px;flex:none}
.logo svg{width:100%;height:100%;filter:drop-shadow(0 0 5px rgba(60,255,110,.7))}
.header h1{font-size:15px;letter-spacing:2px;text-shadow:var(--glow);animation:flicker 6s infinite}
.header .ver{margin-left:auto;font-size:12px;color:var(--dim)}
@keyframes flicker{0%,93%,97%,100%{opacity:1}95%{opacity:.6}}
@keyframes pulse{0%,100%{opacity:.35;transform:scale(.94)}50%{opacity:1;transform:scale(1)}}
.wave{animation:pulse 1.6s infinite}
.tabs{display:flex;flex-wrap:wrap;gap:2px;padding:10px 16px 0;border-bottom:1px solid var(--line)}
.tab{padding:9px 14px;border:1px solid transparent;border-bottom:none;background:none;font:inherit;color:var(--dim);cursor:pointer;letter-spacing:1px;text-transform:uppercase;font-size:12px}
.tab:hover{color:var(--green)}
.tab.active{color:var(--green);border-color:var(--line);background:var(--panel);text-shadow:var(--glow)}
.page{display:none;padding:14px 16px}
.page.active{display:block}
.card{background:var(--panel);border:1px solid var(--line);box-shadow:inset 0 0 18px rgba(60,255,110,.05),0 0 9px rgba(60,255,110,.08);padding:10px 14px;margin-bottom:10px}
.row{display:flex;align-items:center;gap:10px;padding:6px 0;flex-wrap:wrap}
.row label{flex:1;min-width:100px;color:var(--green)}
.row select{flex:0 1 auto;min-width:0;max-width:100%}
.row input[type=text]{flex:0 1 auto;min-width:0}
.row .hint{font-size:11px;color:var(--faint)}
input[type=text],input[type=number],select{padding:7px 10px;border:1px solid var(--line);background:#08100b;color:var(--green);font:inherit;outline:none}
input:focus,select:focus{border-color:var(--dim);box-shadow:var(--glow)}
input::placeholder{color:var(--faint);opacity:.7}
input:disabled,select:disabled{opacity:.35;cursor:default}
#trlangs label:has(input:disabled){opacity:.45}
input[type=text]{width:220px;max-width:100%}select{width:210px;max-width:100%}
input[type=checkbox]{width:16px;height:16px;accent-color:var(--dim)}
input[type=range]{width:170px;accent-color:var(--dim);background:transparent}
button,select,input[type=checkbox],input[type=radio],input[type=range]{cursor:pointer}
button:disabled{cursor:default}
.val{min-width:52px;text-align:right;color:var(--green);text-shadow:var(--glow);font-weight:700}
.hotkey-box{display:flex;gap:8px;align-items:center;flex-wrap:wrap}
.hotkey-val{font-weight:700;font-size:15px;background:#08100b;border:1px solid var(--line);padding:8px 14px;min-width:150px;text-align:center;text-shadow:var(--glow);letter-spacing:1px}
button.btn{padding:8px 18px;border:1px solid var(--dim);background:#0d1a11;color:var(--green);font:inherit;cursor:pointer;letter-spacing:1px;text-transform:uppercase;font-size:12px}
button.btn:hover{background:#123f22;box-shadow:var(--glow)}
button.ghost{border-color:var(--line);color:var(--dim)}
button.ghost:hover{color:var(--green)}
.footer{flex:none;display:flex;gap:12px;align-items:center;padding:10px 16px;background:var(--panel);border-top:1px solid var(--line)}
.toast{color:var(--amber);font-size:13px;opacity:0;transition:opacity .3s;text-shadow:0 0 6px rgba(255,179,71,.5)}
.toast.show{opacity:1}
.mrow{display:flex;align-items:center;gap:9px;padding:7px 2px;border-bottom:1px solid #12241a;flex-wrap:wrap}
.mrow:last-child{border-bottom:none}
.mrow input[type=radio]{width:15px;height:15px;accent-color:var(--dim)}
.mrow .mname{width:104px;font-weight:700}
.mrow .mdesc{flex:1;color:var(--faint);font-size:12px}
.mrow .msize{color:var(--dim);font-size:12px;width:70px;text-align:right}
.badge{font-size:11px;letter-spacing:1px;padding:4px 10px;border:1px solid var(--dim);color:var(--green);text-shadow:var(--glow);text-transform:uppercase}
button.mini{padding:5px 12px;border:1px solid var(--line);background:none;color:var(--dim);font:12px Consolas,monospace;cursor:pointer;text-transform:uppercase}
button.mini:hover{color:var(--green);border-color:var(--dim);box-shadow:var(--glow)}
button.mini.danger:hover{color:#ff7b6b;border-color:#7a2e2e;box-shadow:0 0 7px rgba(255,110,90,.5)}
.mpct{color:var(--amber);font-size:12px;min-width:44px;text-align:right;text-shadow:0 0 6px rgba(255,179,71,.5)}
.sect{color:var(--dim);font-size:12px;letter-spacing:1px;text-transform:uppercase;margin-bottom:6px}
.hfhome{margin-left:auto;cursor:pointer;color:var(--dim);font-size:11px;letter-spacing:1px;border:1px solid var(--line);padding:3px 9px;text-transform:none}
.hfhome:hover{color:var(--green);border-color:var(--dim);box-shadow:var(--glow)}
.ramline{display:flex;align-items:center;flex-wrap:wrap;gap:6px;color:var(--faint);font-size:12px;margin:4px 0 10px}
.ramline b{color:var(--green);font-size:14px;text-shadow:var(--glow);margin-right:4px}
.ramline .dot{margin-left:12px;font-size:10px}
.subhead{color:var(--dim);font-size:11px;letter-spacing:1px;text-transform:uppercase;margin:14px 0 2px;padding-top:10px;border-top:1px solid #12241a}
#hf_results{max-height:44vh;overflow-y:auto;overscroll-behavior:contain}
.subtabs{display:flex;flex-wrap:wrap;gap:2px;align-items:flex-end;border-bottom:1px solid var(--line);margin-bottom:10px}
button.stab{padding:9px 11px;border:1px solid transparent;border-bottom:none;background:none;font:inherit;color:var(--dim);cursor:pointer;letter-spacing:1px;text-transform:uppercase;font-size:12px}
button.stab:hover{color:var(--green)}
button.stab.on{color:var(--green);border-color:var(--line);background:var(--bg);text-shadow:var(--glow)}
.subtabs .hfhome{margin-bottom:6px}
.spage{display:none}
.spage.on{display:block}
#hf_clr{position:absolute;right:9px;top:50%;transform:translateY(-50%);color:var(--dim);cursor:pointer;display:none;font-size:13px;padding:2px 4px}
#hf_clr:hover{color:var(--green);text-shadow:var(--glow)}
#hf_go{position:absolute;left:9px;top:50%;transform:translateY(-50%);color:var(--dim);cursor:pointer;line-height:0;padding:3px}
#hf_go:hover{color:var(--green);filter:drop-shadow(0 0 4px rgba(60,255,110,.6))}
button.iconbtn{border:none;background:none;padding:2px 5px;color:var(--dim);cursor:pointer;line-height:1;font:13px Consolas,monospace}
button.iconbtn:hover{color:var(--green);filter:drop-shadow(0 0 4px rgba(60,255,110,.6))}
button.iconbtn.danger:hover{color:#ff7b6b;filter:drop-shadow(0 0 4px rgba(255,110,90,.5))}
.modal-bg{position:fixed;inset:0;background:rgba(3,7,4,.75);display:none;align-items:center;justify-content:center;z-index:10}
.modal-bg.show{display:flex}
.modal{background:var(--panel);border:1px solid var(--line);box-shadow:0 0 24px rgba(60,255,110,.18);padding:20px 24px;max-width:360px;text-align:center}
.modal p{margin-bottom:16px;line-height:1.5}
.about p{margin:8px 0;line-height:1.55;user-select:text;color:var(--green)}
.about b{text-shadow:var(--glow)}
.about .wh{color:var(--dim);font-size:12px;letter-spacing:1px;text-transform:uppercase;margin:16px 0 4px;border-bottom:1px solid #12241a;padding-bottom:3px}
.about ul{margin:4px 0 10px 20px;padding:0}
.about li{margin:4px 0;line-height:1.55;color:var(--green);user-select:text}
.mock{background:#0a0e0b;border:1px solid var(--line);border-radius:8px;padding:10px 14px;margin:8px 0;max-width:420px;box-shadow:inset 0 0 14px rgba(60,255,110,.06)}
.mock-pill{display:flex;align-items:center;gap:10px}
.mock-dot{width:11px;height:11px;border-radius:50%;background:#ff5b4d;box-shadow:0 0 8px rgba(255,91,77,.8);flex:none}
.mock-bars{display:flex;gap:2px;align-items:center}
.mock-bars i{display:block;width:3px;background:var(--green);border-radius:1px}
.mock-x{margin-left:auto;color:var(--dim)}
.mock-btn{border:1px solid var(--line);padding:5px 14px;color:var(--dim);font-size:13px}
.mock-btn.on{border-color:var(--green);color:var(--green);text-shadow:var(--glow)}
.mock-mi{padding:4px 6px;color:var(--green);font-size:13px}
.mock-mi.dim{color:var(--faint)}
.mock-sep{border:none;border-top:1px solid #12241a;margin:4px 0}
.mock-row{display:flex;align-items:center;gap:8px;padding:4px 0;font-size:13px;border-bottom:1px solid #10201a}
.mock-row:last-child{border-bottom:none}
.mock-radio{width:12px;height:12px;border-radius:50%;border:1px solid var(--dim);flex:none}
.mock-radio.on{background:var(--green);box-shadow:var(--glow)}
.mock-cb{width:13px;height:13px;border:1px solid var(--dim);flex:none;display:inline-flex;align-items:center;justify-content:center;font-size:10px;color:#0b0f0c}
.mock-cb.on{background:var(--green)}
.mock-note{color:var(--faint);font-size:12px}
.lnk{color:var(--green);text-decoration:underline;cursor:pointer}
.lnk:hover{text-shadow:var(--glow)}
</style></head><body>
<div class="header">
 <div class="logo"><svg viewBox="0 0 64 64">
  <rect x="2" y="2" width="60" height="60" rx="12" fill="#0e1410" stroke="#1d4a2b" stroke-width="2"/>
  <g stroke="#3cff6e" stroke-width="4" fill="none" stroke-linecap="round">
   <rect x="26" y="12" width="12" height="20" rx="6" fill="#3cff6e"/>
   <path d="M19 27a13 13 0 0 0 26 0"/>
   <line x1="32" y1="40" x2="32" y2="46"/>
   <line x1="24" y1="49" x2="40" y2="49"/>
  </g>
  <g stroke="#3cff6e" stroke-width="2.5" fill="none" stroke-linecap="round">
   <path class="wave" d="M13 20a17 17 0 0 0 0 14" style="animation-delay:.2s"/>
   <path class="wave" d="M51 20a17 17 0 0 1 0 14" style="animation-delay:.6s"/>
  </g>
 </svg></div>
 <h1>VOX&nbsp;TERMINAL</h1>
 <span class="ver" style="margin-left:auto">v<span id="ver"></span></span>
 <div class="capbtns">
  <button class="cap" onclick="appMin()">&#9472;</button>
  <button class="cap close" onclick="guardClose()">&#10005;</button>
 </div>
</div>
<div class="tabs" id="tabbar">
 <button class="tab" data-p="general">{{S_TAB_GENERAL}}</button>
 <button class="tab" data-p="rec">{{S_TAB_REC}}</button>
 <button class="tab" data-p="proc">{{S_TAB_PROC}}</button>
 <button class="tab" data-p="about">{{S_TAB_ABOUT}}</button>
</div>

<div class="content">
<div class="page" id="p-general">
 <div class="card">
  <div class="row"><label>{{S_HOTKEY}}</label>
   <div class="hotkey-box"><span class="hotkey-val" id="hotkey"></span>
   <button class="btn" onclick="appCapture()">{{S_CHANGE}}</button></div></div>
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
  <div class="row"><label>{{S_RECLANG}}</label>
   <select id="language">
    <option value="auto">{{S_RECAUTO}}</option>
    <option value="ru">Русский</option><option value="en">English</option>
    <option value="uk">Українська</option><option value="de">Deutsch</option>
    <option value="fr">Français</option><option value="es">Español</option>
    <option value="pl">Polski</option>
   </select></div>
 </div>
 <div class="card">
  <div class="sect">{{S_SEC_SOUND}}</div>
  <div class="row"><label>{{S_BEEP}}</label><input type="checkbox" id="beep"></div>
  <div class="row"><label>{{S_SOUND}}</label>
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
 <div class="card">
  <div class="sect">{{S_SEC_BEHAVIOR}}</div>
  <div class="row"><label>{{S_AUTOENTER}}</label><input type="checkbox" id="auto_enter"></div>
  <div class="row"><label>{{S_RESTORE}}</label><input type="checkbox" id="restore_clipboard"></div>
  <div class="row"><label>{{S_OVERLAY}}</label><input type="checkbox" id="overlay"></div>
  <div class="row"><label>{{S_ANIM}}</label><input type="checkbox" id="animation"></div>
  <div class="row"><label>{{S_TYPEMODE}}</label><input type="checkbox" id="type_mode"></div>
 </div>
</div>

<div class="page" id="p-rec">
 <div class="card">
  <div class="subtabs">
   <button class="stab" data-s="models">{{S_SUB_MODELS}}</button>
   <button class="stab" data-s="dict">{{S_SUB_DICT}}</button>
   <button class="stab" data-s="params">{{S_SUB_PARAMS}}</button>
   <button class="stab" data-s="server">{{S_TAB_SERVER}}</button>
   <button class="stab" data-s="translate">{{S_SUB_TR}}</button>
  </div>
  <div id="rec-params" class="spage">
  <div class="row"><label>{{S_THREADS}}</label><span class="val" id="threads_v"></span><input type="range" id="threads" min="1" max="16" step="1"></div>
  <div class="row"><label>{{S_MINMS}}</label><span class="val" id="min_record_ms_v"></span><input type="range" id="min_record_ms" min="0" max="1000" step="50"></div>
  <div class="row"><label>{{S_MAXSEC}}</label><span class="val" id="max_record_seconds_v"></span><input type="range" id="max_record_seconds" min="10" max="300" step="10"></div>
  </div>
  <div id="rec-models" class="spage">
  <div id="models"></div>
  </div>
  <div id="rec-dict" class="spage">
  <div style="color:var(--faint);font-size:12px;margin-bottom:6px">{{S_DICT_HINT}}</div>
  <textarea id="whisper_prompt" rows="14" style="width:100%;min-height:180px;height:38vh;padding:8px 11px;border:1px solid var(--line);background:#08100b;color:var(--green);font:inherit;line-height:1.5;outline:none;resize:vertical"></textarea>
  </div>
  <div id="rec-translate" class="spage">
  <div style="color:var(--faint);font-size:12px;margin-bottom:6px">{{S_TR_HINT}}</div>
  <div class="row"><label>{{S_TR_DEFAULT}}</label><input type="checkbox" id="tr_default"></div>
  <div class="row"><label>{{S_TR_TARGET}}</label>
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
  <div class="row"><label>{{S_TR_SECONDS}}</label><span class="val" id="translate_ask_seconds_v"></span><input type="range" id="translate_ask_seconds" min="1" max="10" step="1"></div>
  <div class="row"><label>{{S_PROF_HOTKEY}}</label>
   <span class="hotkey-val" id="tr_hotkey" style="min-width:110px"></span>
   <button class="mini" id="tr_set">{{S_PROF_SET}}</button>
   <button class="mini" id="tr_clear">{{S_PROF_CLEAR}}</button></div>
  <div class="row"><label>{{S_TR_LANGS}}</label>
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
  <div id="rec-server" class="spage">
  <div class="row"><label>{{S_AUTOSTART}}</label><input type="checkbox" id="server_autostart"></div>
  <div class="row"><label>{{S_PORT}}</label><input type="text" id="server_port" style="width:90px"></div>
  <div class="row"><label>{{S_SERVEREXE}}</label><input type="text" id="server_exe"></div>
  <div class="row"><label>{{S_SERVERURL}}<div class="hint">{{S_URLHINT}}</div></label><input type="text" id="server_url"></div>
  </div>
 </div>
</div>

<div class="page" id="p-proc">
 <div style="color:var(--faint);font-size:12px;letter-spacing:1px;margin:0 2px 10px">{{S_PIPE}}</div>
 <div class="card">
  <div class="subtabs">
   <button class="stab" data-s="models">{{S_SUB_MODELS}}</button>
   <button class="stab" data-s="search">{{S_SUB_SEARCH}}</button>
   <button class="stab" data-s="prompts">{{S_SUB_PROMPTS}}</button>
   <span class="hfhome" onclick="appHFHome()" title="huggingface.co">Hugging Face ↗</span>
  </div>
  <div id="proc-models" class="spage"></div>
  <div id="proc-search" class="spage"></div>
  <div id="proc-prompts" class="spage">
   <div style="color:var(--faint);font-size:12px;margin-bottom:8px">{{S_LLM_HINT}}</div>
   <div id="profbody"></div>
  </div>
 </div>
</div>

<div class="page about" id="p-about">
 <div class="card">
  <div class="subtabs">
   <button class="stab" data-s="info">{{S_SUB_INFO}}</button>
   <button class="stab" data-s="help">{{S_SUB_HELP}}</button>
   <button class="stab" data-s="author">{{S_SUB_AUTHOR}}</button>
  </div>
  <div id="about-info" class="spage">
  <p style="font-size:15px;letter-spacing:2px"><b>VOX TERMINAL</b> <span id="ver2"></span></p>
  {{S_ABOUT_HTML}}
  <div class="row" style="border-top:1px solid #12241a;margin-top:10px;padding-top:12px">
   <label>{{S_UPD}}</label>
   <button class="mini" id="upd_check">{{S_UPD_CHECK}}</button></div>
  <div class="row"><label>{{S_UPD_AUTO}}</label><input type="checkbox" id="check_updates"></div>
  <div id="upd_out" style="font-size:12px;min-height:18px;color:var(--amber)"></div>
  </div>
  <div id="about-help" class="spage">
  {{S_HELP_HTML}}
  </div>
  <div id="about-author" class="spage">
  {{S_AUTHOR_HTML}}
  </div>
 </div>
</div>
</div>

<div class="footer">
 <button class="btn" onclick="save()">{{S_SAVE}}</button>
 <span class="toast" id="toast"></span>
</div>

<div class="modal-bg" id="modalbg">
 <div class="modal">
  <p>{{S_UNSAVED}}</p>
  <div style="display:flex;gap:10px;justify-content:center">
   <button class="btn" id="mYes">{{S_SAVE_YES}}</button>
   <button class="ghost btn" id="mNo">{{S_SAVE_NO}}</button>
  </div>
 </div>
</div>

<script>
const CFG = {{CFG}};
const bools = ["beep","auto_enter","restore_clipboard","overlay","animation","type_mode","server_autostart","check_updates"];
const texts = ["server_exe","server_url"];
const nums  = ["threads","min_record_ms","max_record_seconds","translate_ask_seconds","server_port"];
const sels  = ["ui_language","language","sound_theme","translate_target","translate_ask"];
const trAll = ["en","de","fr","es","it","pl","ru","uk"];
const L = {dl:"{{S_DL}}", del:"{{S_DEL}}", hint:"{{S_APPLY_HINT}}",
  add:"{{S_PROF_ADD}}", pname:"{{S_PROF_NAME}}", pprompt:"{{S_PROF_PROMPT}}", phot:"{{S_PROF_HOTKEY}}",
  pset:"{{S_PROF_SET}}", pclr:"{{S_PROF_CLEAR}}", ptest:"{{S_PROF_TEST}}",
  nohot:"—", fitok:"{{S_FIT_OK}}", fitwarn:"{{S_FIT_WARN}}", fitbad:"{{S_FIT_BAD}}",
  ram:"{{S_RAM}}", hfph:"{{S_HF_PH}}", nollm:"{{S_NO_LLM}}",
  nollmp:"{{S_NO_LLM_PROF}}", upd:"{{S_UPDATED}}", pedit:"{{S_PROF_EDIT}}", pclose:"{{S_PROF_CLOSE}}",
  updnone:"{{S_UPD_NONE}}", updavail:"{{S_UPD_AVAIL}}", updgo:"{{S_UPD_GO}}", upderr:"{{S_UPD_ERR}}", upddl:"{{S_UPD_DL}}"};
const I_DL = '<svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2"><path d="M12 3v12"/><path d="M6 11l6 6 6-6"/><path d="M4 21h16"/></svg>';
const I_FIND = '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4"><circle cx="10.5" cy="10.5" r="6.5"/><line x1="15.5" y1="15.5" x2="21" y2="21"/></svg>';

let profiles = (CFG.profiles || []).map(p=>Object.assign({}, p));
let activeProfiles = (CFG.active_profiles || []).slice();
let translateDefault = !!CFG.translate_default;
let translateHotkey = CFG.translate_hotkey || "";
let expandedID = null;
let captureFor = null;

function esc(s){ const d=document.createElement("span"); d.textContent=s||""; return d.innerHTML; }

let selLLM = null;
let hfRepos = [];
let hfOpenRepo = null;
let hfFiles = [];
let hfQuery = "";
function fitLabel(f, need){
  const gb = (need/1024).toFixed(1);
  const col = f==="ok" ? "var(--green)" : (f==="warn" ? "var(--amber)" : "#ff7b6b");
  const tip = f==="ok" ? L.fitok : (f==="warn" ? L.fitwarn : L.fitbad);
  return '<span title="'+esc(tip)+'" style="color:'+col+';font-size:12px;white-space:nowrap">&#9679; &#8776;'+gb+' GB</span>';
}
const curSubs = {};
function showSub(page, s){
  curSubs[page] = s;
  document.querySelectorAll("#p-"+page+" .stab").forEach(b=>b.classList.toggle("on", b.dataset.s===s));
  document.querySelectorAll("#p-"+page+" .spage").forEach(el=>el.classList.toggle("on", el.id===page+"-"+s));
}
async function refreshLLM(){
  const st = JSON.parse(await appLLM());
  const installed = st.installed || [];
  if(selLLM === null){
    const act = installed.find(m=>m.active);
    if(act) selLLM = act.file;
  }
  if(!curSubs.proc) showSub("proc", installed.length ? "models" : "search");

  const body = document.getElementById("proc-models");
  body.innerHTML = "";
  if(!installed.length && !(st.downloads||[]).length){
    const empty = document.createElement("div");
    empty.style.cssText = "color:var(--faint);font-size:12px;padding:6px 2px";
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
      div.innerHTML = '<span class="mdesc" style="flex:1">'+esc(d.file)+'</span><span class="mpct">'+(d.pct>0?d.pct+"%":"…")+'</span>';
    } else {
      div.innerHTML = '<span class="mdesc" style="flex:1">'+esc(d.file)+'</span><span class="mpct">! '+esc(d.err)+'</span>';
    }
    body.appendChild(div);
  });
  body.querySelectorAll('input[name="llmmdl"]').forEach(r=>{
    r.onchange = ()=>{ selLLM = r.value; toast(L.hint); refreshLLM(); };
  });
  body.querySelectorAll("button[data-a='ldel']").forEach(b=>{
    b.onclick = async ()=>{
      const f = b.dataset.f;
      toast(await appLLMDel(f));
      if(selLLM === f){
        selLLM = null;
        if(baseline){ const bl = JSON.parse(baseline); bl.lm = ""; baseline = JSON.stringify(bl); }
      }
      refreshLLM();
    };
  });

  const sbody = document.getElementById("proc-search");
  sbody.innerHTML = '<div class="row" style="padding-top:0">'+
    '<span style="position:relative;flex:1;display:flex;min-width:0">'+
    '<input type="text" id="hf_q" placeholder="'+L.hfph+'" style="flex:1;min-width:0;padding-left:34px;padding-right:30px">'+
    '<span id="hf_clr">&#10005;</span>'+
    '<span id="hf_go">'+I_FIND+'</span></span></div>'+
    '<div class="ramline">'+L.ram+' <b>'+(st.ram/1024).toFixed(0)+' GB</b>'+
      '<span class="dot" style="color:var(--green)">&#9679;</span>'+L.fitok+
      '<span class="dot" style="color:var(--amber)">&#9679;</span>'+L.fitwarn+
      '<span class="dot" style="color:#ff7b6b">&#9679;</span>'+L.fitbad+'</div>'+
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
    div.style.cursor = "pointer";
    div.innerHTML = '<span class="mdesc" style="flex:1;color:var(--dim)">'+(hfOpenRepo===r.id?"▾ ":"▸ ")+esc(r.id)+'</span>'+
      '<span title="'+L.upd+'" style="color:var(--dim);font-size:12px">'+esc(r.updated||"")+'</span>'+
      '<span class="msize">↓'+(r.downloads>=1000000?(r.downloads/1000000).toFixed(1)+"M":(r.downloads/1000).toFixed(0)+"k")+'</span>'+
      '<span class="hflink" title="huggingface.co/'+esc(r.id)+'" style="color:var(--dim);cursor:pointer;padding:0 4px">↗</span>';
    div.querySelector(".hflink").onclick = e=>{ e.stopPropagation(); appHFPage(r.id); };
    div.onclick = async ()=>{
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
          '<button class="iconbtn" title="'+L.dl+'" data-repo="'+esc(r.id)+'" data-file="'+esc(f.file)+'">'+I_DL+'</button>';
        box.appendChild(fd);
      });
      box.querySelectorAll("button[data-repo]").forEach(b=>{
        b.onclick = async e=>{
          e.stopPropagation();
          await appLLMDlFile(b.dataset.repo, b.dataset.file);
          showSub("proc", "models");
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
    note.style.cssText = "color:var(--faint);font-size:12px;padding:4px 2px";
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
    b.onclick = ()=>{
      const id = b.dataset.id;
      if(b.dataset.a === "edit"){
        expandedID = expandedID === id ? null : id;
        refreshLLM();
      } else {
        profiles = profiles.filter(x=>x.id!==id);
        activeProfiles = activeProfiles.filter(x=>x!==id);
        if(expandedID === id) expandedID = null;
        refreshLLM();
      }
    };
  });
}

function renderEditor(p){
  const ed = document.createElement("div");
  ed.style.cssText = "border:1px solid #12241a;padding:10px;margin:2px 0 8px";
  ed.innerHTML =
    '<div style="display:flex;justify-content:flex-end;margin:-4px -4px 2px 0"><button class="iconbtn" id="pf_close" title="'+L.pclose+'">&#9650;</button></div>'+
    '<div class="row" style="padding-top:0"><label>'+L.pname+'</label><input type="text" id="pf_name"></div>'+
    '<div class="row" style="align-items:flex-start"><label>'+L.pprompt+'</label></div>'+
    '<textarea id="pf_prompt" rows="4" style="width:100%;padding:7px 10px;border:1px solid var(--line);background:#08100b;color:var(--green);font:inherit;outline:none;resize:vertical"></textarea>'+
    '<div class="row"><label>'+L.phot+'</label><span class="hotkey-val" id="pf_hotkey" style="min-width:110px"></span>'+
    '<button class="mini" id="pf_set">'+L.pset+'</button><button class="mini" id="pf_clear">'+L.pclr+'</button></div>'+
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
    document.getElementById("pf_set").onclick = ()=>{ captureFor = p.id; appCaptureCombo(); };
    document.getElementById("pf_clear").onclick = ()=>{ p.hotkey=""; hk.textContent=L.nohot; };
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

function updTrHotkey(){
  const el = document.getElementById("tr_hotkey");
  if(el) el.textContent = translateHotkey || L.nohot;
}
function comboCaptured(combo){
  if(!captureFor) return;
  if(captureFor === "__wt"){
    captureFor = null;
    if(combo) translateHotkey = combo;
    updTrHotkey();
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

let selModel = null;
let activeModelId = null;
let pendingDl = null;
let baseline = null;
async function refreshModels(){
  const rows = JSON.parse(await appModels());
  const activeRow = rows.find(m=>m.state==="active");
  activeModelId = activeRow ? activeRow.id : null;
  if(pendingDl){
    const row = rows.find(m=>m.id===pendingDl);
    if(!row){ pendingDl = null; }
    else if(row.state === "installed"){ pendingDl = null; toast(L.hint); }
    else if(row.state === "active"){ pendingDl = null; }
    else if(row.state === "absent" && row.err){ pendingDl = null; toast(row.err); }
  }
  const checkedId = selModel || activeModelId;
  const el = document.getElementById("models");
  el.innerHTML = "";
  let busy = false;
  rows.forEach(m=>{
    const div = document.createElement("div");
    div.className = "mrow";
    const checked = checkedId === m.id ? " checked" : "";
    const radio = '<input type="radio" name="mdl" value="'+m.id+'"'+checked+(m.id==="custom"?" disabled":"")+'>';
    let right = "";
    if(m.state === "downloading"){ busy = true; right = '<span class="mpct">'+(m.pct>0?m.pct+"%":"…")+'</span>'; }
    else if(m.state === "absent") right = '<button class="iconbtn" title="'+L.dl+'" data-a="dl" data-id="'+m.id+'">'+I_DL+'</button>';
    else if(m.state === "installed") right = '<button class="iconbtn danger" title="'+L.del+'" data-a="del" data-id="'+m.id+'">&#10005;</button>';
    div.innerHTML = radio+'<span class="mname">'+m.name+'</span><span class="mdesc">'+m.desc+'</span><span class="msize">'+(m.size?m.size+" MB":"")+'</span><span>'+right+'</span>';
    el.appendChild(div);
  });
  el.querySelectorAll('input[name="mdl"]').forEach(r=>{
    r.onchange = async ()=>{
      const id = r.value;
      const row = rows.find(m=>m.id===id);
      if(!row) return;
      selModel = id;
      if(row.state === "absent"){
        pendingDl = id;
        await appModelDl(id);
      } else if(row.state === "installed"){
        toast(L.hint);
      }
      refreshModels();
    };
  });
  el.querySelectorAll("button[data-a]").forEach(b=>{
    b.onclick = async ()=>{
      if(b.dataset.a === "dl"){ await appModelDl(b.dataset.id); }
      else {
        toast(await appModelDel(b.dataset.id));
        if(selModel === b.dataset.id) selModel = null;
      }
      refreshModels();
    };
  });
  if(busy || pendingDl) setTimeout(refreshModels, 700);
}
function formState(){
  const o = {m: selModel || activeModelId};
  bools.forEach(k=>o[k]=document.getElementById(k).checked);
  texts.forEach(k=>o[k]=document.getElementById(k).value);
  nums.forEach(k=>o[k]=document.getElementById(k).value);
  sels.forEach(k=>o[k]=document.getElementById(k).value);
  o.hotkey = CFG.hotkey;
  o.wp = document.getElementById("whisper_prompt").value;
  o.td = translateDefault;
  o.ap = activeProfiles.join(",");
  o.lm = selLLM || "";
  o.wt = translateHotkey;
  o.tal = trAll.filter(l=>document.getElementById("tl_"+l).checked).join(",");
  o.prof = JSON.stringify(profiles);
  return JSON.stringify(o);
}
function dirty(){ return baseline !== null && formState() !== baseline; }
function revert(){
  if(baseline === null) return;
  const b = JSON.parse(baseline);
  bools.forEach(k=>document.getElementById(k).checked = b[k]);
  texts.forEach(k=>document.getElementById(k).value = b[k]);
  nums.forEach(k=>document.getElementById(k).value = b[k]);
  sels.forEach(k=>document.getElementById(k).value = b[k]);
  sliders.forEach(k=>document.getElementById(k).dispatchEvent(new Event("input")));
  selModel = b.m || null;
  document.getElementById("whisper_prompt").value = b.wp || "";
  translateDefault = !!b.td;
  activeProfiles = (b.ap||"") === "" ? [] : b.ap.split(",");
  selLLM = b.lm || null;
  translateHotkey = b.wt || "";
  document.getElementById("tr_default").checked = translateDefault;
  syncTrControls();
  updTrHotkey();
  const tal = (b.tal||"").split(",");
  trAll.forEach(l=>{ document.getElementById("tl_"+l).checked = tal.includes(l); });
  try { profiles = JSON.parse(b.prof) || []; } catch(e) {}
  expandedID = null;
  refreshModels();
  refreshLLM();
}
function askUnsaved(cb){
  const bg = document.getElementById("modalbg");
  bg.classList.add("show");
  const done = v => { bg.classList.remove("show"); if(v !== undefined) cb(v); };
  document.getElementById("mYes").onclick = ()=>done(true);
  document.getElementById("mNo").onclick = ()=>done(false);
  bg.onclick = e => { if(e.target === bg) done(undefined); };
}
function guardClose(){
  if(!dirty()) return appClose();
  askUnsaved(async s=>{ if(s) await doSave(); appClose(); });
}
function toast(msg){
  if(!msg) return;
  const t = document.getElementById("toast");
  t.textContent = msg; t.classList.add("show");
  setTimeout(()=>t.classList.remove("show"), 2500);
}

const sliders = ["threads","min_record_ms","max_record_seconds","translate_ask_seconds"];
function load(){
  document.getElementById("whisper_prompt").value = CFG.whisper_prompt || "";
  trAll.forEach(l=>{ document.getElementById("tl_"+l).checked = (CFG.translate_ask_langs||[]).includes(l); });
  const trd = document.getElementById("tr_default");
  trd.checked = translateDefault;
  trd.onchange = ()=>{ translateDefault = trd.checked; syncTrControls(); };
  document.getElementById("translate_ask").onchange = syncTrControls;
  updTrHotkey();
  document.getElementById("tr_set").onclick = ()=>{ captureFor = "__wt"; appCaptureCombo(); };
  document.getElementById("tr_clear").onclick = ()=>{ translateHotkey = ""; updTrHotkey(); };
  document.getElementById("hotkey").textContent = CFG.hotkey;
  document.getElementById("ver").textContent = CFG._version;
  document.getElementById("ver2").textContent = CFG._version;
  document.getElementById("threads").max = CFG._cpus || 16;
  bools.forEach(k=>document.getElementById(k).checked = !!CFG[k]);
  texts.forEach(k=>document.getElementById(k).value = CFG[k]||"");
  nums.forEach(k=>document.getElementById(k).value = CFG[k]);
  sels.forEach(k=>{
    const el=document.getElementById(k), v=CFG[k]||"auto";
    if(![...el.options].some(o=>o.value===v)){const o=document.createElement("option");o.value=v;o.textContent=v;el.appendChild(o);}
    el.value=v;
  });
  sliders.forEach(k=>{
    const el = document.getElementById(k), v = document.getElementById(k+"_v");
    const upd = ()=>{ v.textContent = el.value; };
    el.oninput = upd;
    upd();
  });
  syncTrControls();
}
function syncTrControls(){
  const always = document.getElementById("tr_default").checked;
  const mode = document.getElementById("translate_ask").value;
  document.getElementById("translate_target").disabled = !always;
  document.getElementById("translate_ask").disabled = always;
  document.getElementById("translate_ask_seconds").disabled = always || mode !== "timeout";
  trAll.forEach(l=>{ document.getElementById("tl_"+l).disabled = always || mode === "never"; });
}
function setHotkey(s){
  CFG.hotkey=s;
  document.getElementById("hotkey").textContent=s;
  if(baseline !== null){ const b=JSON.parse(baseline); b.hotkey=s; baseline=JSON.stringify(b); }
}
async function doSave(){
  const f={hotkey:CFG.hotkey, model_id:selModel||"",
    whisper_prompt: document.getElementById("whisper_prompt").value,
    translate_hotkey: translateHotkey,
    translate_ask_langs: trAll.filter(l=>document.getElementById("tl_"+l).checked),
    translate_default: translateDefault,
    active_profiles: activeProfiles,
    llm_model_file: selLLM||"",
    profiles: profiles};
  bools.forEach(k=>f[k]=document.getElementById(k).checked);
  texts.forEach(k=>f[k]=document.getElementById(k).value);
  nums.forEach(k=>f[k]=parseInt(document.getElementById(k).value)||0);
  sels.forEach(k=>f[k]=document.getElementById(k).value);
  const langChanged = f.ui_language !== (CFG.ui_language || "auto");
  const msg = await appSave(JSON.stringify(f));
  baseline = formState();
  if(langChanged){ appReload(curTab); return; }
  toast(msg);
  refreshModels();
}
function save(){ doSave(); }
let curTab = "general";
function show(p){
  curTab = p;
  document.querySelectorAll(".tab").forEach(b=>b.classList.toggle("active", b.dataset.p===p));
  document.querySelectorAll(".page").forEach(el=>el.classList.toggle("active", el.id==="p-"+p));
  document.querySelector(".footer").style.display = (p==="about") ? "none" : "flex";
}
document.querySelectorAll(".stab").forEach(b=>b.onclick=()=>showSub(b.closest(".page").id.slice(2), b.dataset.s));
document.getElementById("upd_check").onclick = updCheck;
(async()=>{ const s = JSON.parse(await appUpdateStatus()); if(s.latest && s.url) updShow(s.latest, true); })();
document.querySelectorAll(".tab").forEach(b=>b.onclick=()=>{
  const p = b.dataset.p;
  if(p === curTab) return;
  if(dirty()){
    askUnsaved(async s=>{ if(s){ await doSave(); } else { revert(); } show(p); });
  } else {
    show(p);
  }
});
document.querySelector(".header").addEventListener("mousedown", e=>{
  if(e.target.closest(".cap")) return;
  if(e.button===0) appDrag();
});
load();
(async ()=>{
  await refreshModels();
  await refreshLLM();
  baseline = formState();
})();
showSub("rec", "models");
showSub("about", "info");
show(["about","rec","proc","server"].includes(CFG._tab) ? (CFG._tab==="server" ? "rec" : CFG._tab) : "general");
if(CFG._tab === "server") showSub("rec", "server");
</script>
</body></html>`
