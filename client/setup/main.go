package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	webview "github.com/jchv/go-webview2"
	"golang.org/x/sys/windows"

	"holdtotype/internal/appid"
)

func init() {
	runtime.LockOSThread()
}

var texts = map[string]map[string]string{
	"ru": {
		"title":     appid.Name + " — установка",
		"tagline":   "Голос → текст в позицию курсора. Полностью локально и офлайн.",
		"path":      "Папка установки",
		"shortcut":  "Ярлык в меню «Пуск»",
		"autorun":   "Автозапуск с Windows",
		"launch":    "Запустить после установки",
		"install":   "УСТАНОВИТЬ",
		"update":    "ОБНОВИТЬ",
		"updnote":   "Обнаружена установленная версия — настройки и скачанные модели будут сохранены.",
		"model":     "Модель распознавания",
		"mnone":     "Не скачивать (скачаете при первом запуске)",
		"browse":    "Выберите папку установки",
		"prog":      "Установка…",
		"done":      "Установлено",
		"doneat":    "Приложение установлено в:",
		"warnmodel": "Модель не скачалась — приложение предложит её при первом запуске.",
		"finish":    "ГОТОВО",
		"err":       "Ошибка установки",
		"webview":   "Для установщика нужен Microsoft WebView2 Runtime (входит в Windows 11).\nСейчас откроется страница загрузки.",
	},
	"uk": {
		"title":     appid.Name + " — встановлення",
		"tagline":   "Голос → текст у позицію курсора. Повністю локально й офлайн.",
		"path":      "Тека встановлення",
		"shortcut":  "Ярлик у меню «Пуск»",
		"autorun":   "Автозапуск із Windows",
		"launch":    "Запустити після встановлення",
		"install":   "ВСТАНОВИТИ",
		"update":    "ОНОВИТИ",
		"updnote":   "Виявлено встановлену версію — налаштування та завантажені моделі буде збережено.",
		"model":     "Модель розпізнавання",
		"mnone":     "Не завантажувати (завантажите при першому запуску)",
		"browse":    "Оберіть теку встановлення",
		"prog":      "Встановлення…",
		"done":      "Встановлено",
		"doneat":    "Застосунок встановлено в:",
		"warnmodel": "Модель не завантажилася — застосунок запропонує її при першому запуску.",
		"finish":    "ГОТОВО",
		"err":       "Помилка встановлення",
		"webview":   "Для інсталятора потрібен Microsoft WebView2 Runtime (входить до Windows 11).\nЗараз відкриється сторінка завантаження.",
	},
	"en": {
		"title":     appid.Name + " — Setup",
		"tagline":   "Voice → text at the cursor position. Fully local and offline.",
		"path":      "Install folder",
		"shortcut":  "Start Menu shortcut",
		"autorun":   "Start with Windows",
		"launch":    "Launch after install",
		"install":   "INSTALL",
		"update":    "UPDATE",
		"updnote":   "An existing installation was found — settings and downloaded models will be kept.",
		"model":     "Recognition model",
		"mnone":     "Don't download (get it on first run)",
		"browse":    "Select the install folder",
		"prog":      "Installing…",
		"done":      "Installed",
		"doneat":    "The application is installed to:",
		"warnmodel": "The model could not be downloaded — the app will offer it on first run.",
		"finish":    "FINISH",
		"err":       "Install error",
		"webview":   "The installer requires Microsoft WebView2 Runtime (bundled with Windows 11).\nThe download page will open now.",
	},
}

var langOverride string

func tr(key string) string {
	l := langOverride
	if _, ok := texts[l]; !ok {
		l = uiLang()
	}
	if s, ok := texts[l][key]; ok {
		return s
	}
	return texts["en"][key]
}

func page(updateDir string) string {
	opts := ""
	for _, m := range modelOpts {
		sel := ""
		if m.ID == "small" {
			sel = " selected"
		}
		opts += fmt.Sprintf(`<option value="%s"%s>%s (%d MB)</option>`, m.ID, sel, m.Name, m.SizeMB)
	}
	opts += `<option value="">` + tr("mnone") + `</option>`
	installLabel := tr("install")
	freshStyle := ""
	updNote := ""
	defDir := defaultInstallDir()
	updating := "false"
	if updateDir != "" {
		installLabel = tr("update")
		freshStyle = "display:none"
		updNote = `<div class="warn">` + tr("updnote") + `</div>`
		defDir = updateDir
		updating = "true"
	}
	repl := map[string]string{
		"{{TITLE}}":        tr("title"),
		"{{APP}}":          appid.Name,
		"{{TAGLINE}}":      tr("tagline"),
		"{{PATH}}":         tr("path"),
		"{{MODEL}}":        tr("model"),
		"{{MODELOPTS}}":    opts,
		"{{SHORTCUT}}":     tr("shortcut"),
		"{{AUTORUN}}":      tr("autorun"),
		"{{LAUNCH}}":       tr("launch"),
		"{{INSTALL}}":      installLabel,
		"{{FRESHSTYLE}}":   freshStyle,
		"{{UPDNOTE}}":      updNote,
		"{{UPDATING}}":     updating,
		"{{PROG}}":         tr("prog"),
		"{{DONE}}":         tr("done"),
		"{{DONEAT}}":       tr("doneat"),
		"{{WARNMODEL_JS}}": template_jsstr(tr("warnmodel")),
		"{{FINISH}}":       tr("finish"),
		"{{VERSION}}":      appVersion,
		"{{DEFDIR}}":       template_jsstr(defDir),
	}
	h := setupPage
	for k, v := range repl {
		h = strings.ReplaceAll(h, k, v)
	}
	return h
}

func template_jsstr(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func main() {
	_ = os.Setenv("WEBVIEW2_DEFAULT_BACKGROUND_COLOR", "FF0B0F0C")
	args := os.Args[1:]
	silent := false
	update := false
	dir := ""
	modelID := ""
	for i, a := range args {
		if a == "-silent" {
			silent = true
		}
		if a == "-update" {
			update = true
		}
		if a == "-dir" && i+1 < len(args) {
			dir = args[i+1]
		}
		if a == "-lang" && i+1 < len(args) {
			langOverride = args[i+1]
		}
		if a == "-model" && i+1 < len(args) {
			modelID = args[i+1]
		}
	}
	if update {
		if dir == "" {
			dir = existingInstall()
		}
		if dir == "" {
			os.Exit(1)
		}
		lnkExists := false
		if _, err := os.Stat(startMenuLnk()); err == nil {
			lnkExists = true
		}
		if _, err := install(dir, lnkExists, false, false, "", func(int, string) {}); err != nil {
			os.Exit(1)
		}
		exe := filepath.Join(dir, exeName)
		cmd := exec.Command(exe)
		cmd.Dir = dir
		_ = cmd.Start()
		return
	}
	if silent {
		if dir == "" {
			dir = defaultInstallDir()
		}
		if _, err := install(dir, true, false, true, modelID, func(int, string) {}); err != nil {
			os.Exit(1)
		}
		return
	}
	updateDir := existingInstall()
	if err := windows.CoInitializeEx(0, windows.COINIT_APARTMENTTHREADED); err != nil {
		_ = err
	}
	stopHider := hideWebViewWindowEarly(tr("title"))
	w := webview.NewWithOptions(webview.WebViewOptions{
		DataPath:  filepath.Join(os.TempDir(), appid.TempDirName("setup", os.Getpid())),
		AutoFocus: true,
		WindowOptions: webview.WindowOptions{
			Title:  tr("title"),
			Width:  560,
			Height: 560,
			IconId: 1,
			Center: true,
		},
	})
	stopHider()
	if w == nil {
		msgBox(tr("err"), tr("webview"))
		shellOpenURL("https://developer.microsoft.com/en-us/microsoft-edge/webview2/")
		return
	}
	defer w.Destroy()

	hwnd := uintptr(w.Window())
	setDarkClientBackground(hwnd)
	applyDarkCaption(hwnd)
	makeBorderless(hwnd)
	var shown int32
	reveal := func() {
		if !atomic.CompareAndSwapInt32(&shown, 0, 1) {
			return
		}
		w.Dispatch(func() { revealWindowCentered(hwnd, 560, 560) })
	}
	_ = w.Bind("appReady", reveal)
	go func() {
		time.Sleep(3 * time.Second)
		reveal()
	}()

	installedDir := ""
	_ = w.Bind("appDrag", func() { beginWindowDrag(hwnd) })
	_ = w.Bind("appClose", func() { w.Dispatch(func() { w.Terminate() }) })
	_ = w.Bind("appBrowse", func() string {
		return browseFolder(hwnd, tr("browse"))
	})
	_ = w.Bind("appInstall", func(dir string, shortcut, autorun bool, modelID string, updating bool) {
		go func() {
			touchAutorun := !updating
			if updating {
				shortcut = false
				if _, err := os.Stat(startMenuLnk()); err == nil {
					shortcut = true
				}
				autorun = false
				modelID = ""
			}
			warn, err := install(dir, shortcut, autorun, touchAutorun, modelID, func(pct int, name string) {
				enc, _ := json.Marshal(name)
				w.Dispatch(func() { w.Eval(fmt.Sprintf("setupProgress(%d,%s)", pct, enc)) })
			})
			msg := ""
			if err != nil {
				msg = err.Error()
			} else {
				installedDir = dir
			}
			enc, _ := json.Marshal(msg)
			encWarn, _ := json.Marshal(warn)
			encDir, _ := json.Marshal(dir)
			w.Dispatch(func() { w.Eval("setupDone(" + string(enc) + "," + string(encWarn) + "," + string(encDir) + ")") })
		}()
	})
	_ = w.Bind("appFinish", func(launch bool) {
		if launch && installedDir != "" {
			exe := filepath.Join(installedDir, exeName)
			cmd := exec.Command(exe)
			cmd.Dir = installedDir
			_ = cmd.Start()
		}
		w.Dispatch(func() { w.Terminate() })
	})

	w.SetHtml(page(updateDir))
	w.Run()
}
