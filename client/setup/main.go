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
		"swon":      "ВКЛ",
		"swoff":     "ВЫКЛ",
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
		"retry":     "ПОВТОРИТЬ",
		"back":      "НАЗАД",
		"nodir":     "Укажите папку установки",
		"webview":   "Для установщика нужен Microsoft WebView2 Runtime (входит в Windows 11).\nСейчас откроется страница загрузки.",
	},
	"uk": {
		"swon":      "УВІМК",
		"swoff":     "ВИМК",
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
		"retry":     "ПОВТОРИТИ",
		"back":      "НАЗАД",
		"nodir":     "Вкажіть теку встановлення",
		"webview":   "Для інсталятора потрібен Microsoft WebView2 Runtime (входить до Windows 11).\nЗараз відкриється сторінка завантаження.",
	},
	"de": {
		"swon":      "EIN",
		"swoff":     "AUS",
		"title":     appid.Name + " — Installation",
		"tagline":   "Stimme → Text an der Eingabemarke. Vollständig lokal und offline.",
		"path":      "Installationsordner",
		"shortcut":  "Verknüpfung im Startmenü",
		"autorun":   "Mit Windows starten",
		"launch":    "Nach der Installation starten",
		"install":   "INSTALLIEREN",
		"update":    "AKTUALISIEREN",
		"updnote":   "Eine vorhandene Installation wurde gefunden — Einstellungen und geladene Modelle bleiben erhalten.",
		"model":     "Erkennungsmodell",
		"mnone":     "Nicht laden (beim ersten Start holen)",
		"browse":    "Installationsordner wählen",
		"prog":      "Installation…",
		"done":      "Installiert",
		"doneat":    "Die Anwendung wurde installiert nach:",
		"warnmodel": "Das Modell konnte nicht geladen werden — die App bietet es beim ersten Start an.",
		"finish":    "FERTIG",
		"err":       "Installationsfehler",
		"retry":     "WIEDERHOLEN",
		"back":      "ZURÜCK",
		"nodir":     "Wählen Sie den Installationsordner",
		"webview":   "Der Installer benötigt Microsoft WebView2 Runtime (in Windows 11 enthalten).\nDie Downloadseite wird jetzt geöffnet.",
	},
	"fr": {
		"swon":      "ON",
		"swoff":     "OFF",
		"title":     appid.Name + " — Installation",
		"tagline":   "Voix → texte à la position du curseur. Entièrement local et hors ligne.",
		"path":      "Dossier d'installation",
		"shortcut":  "Raccourci dans le menu Démarrer",
		"autorun":   "Démarrer avec Windows",
		"launch":    "Lancer après l'installation",
		"install":   "INSTALLER",
		"update":    "METTRE À JOUR",
		"updnote":   "Une installation existante a été trouvée — les réglages et les modèles téléchargés seront conservés.",
		"model":     "Modèle de reconnaissance",
		"mnone":     "Ne pas télécharger (au premier lancement)",
		"browse":    "Choisissez le dossier d'installation",
		"prog":      "Installation…",
		"done":      "Installé",
		"doneat":    "L'application est installée dans :",
		"warnmodel": "Le modèle n'a pas pu être téléchargé — l'application le proposera au premier lancement.",
		"finish":    "TERMINER",
		"err":       "Erreur d'installation",
		"retry":     "RÉESSAYER",
		"back":      "RETOUR",
		"nodir":     "Indiquez le dossier d'installation",
		"webview":   "L'installateur nécessite Microsoft WebView2 Runtime (fourni avec Windows 11).\nLa page de téléchargement va s'ouvrir.",
	},
	"es": {
		"swon":      "ON",
		"swoff":     "OFF",
		"title":     appid.Name + " — Instalación",
		"tagline":   "Voz → texto en la posición del cursor. Totalmente local y sin conexión.",
		"path":      "Carpeta de instalación",
		"shortcut":  "Acceso directo en el menú Inicio",
		"autorun":   "Iniciar con Windows",
		"launch":    "Ejecutar tras instalar",
		"install":   "INSTALAR",
		"update":    "ACTUALIZAR",
		"updnote":   "Se encontró una instalación existente: se conservarán los ajustes y los modelos descargados.",
		"model":     "Modelo de reconocimiento",
		"mnone":     "No descargar (se descarga al primer inicio)",
		"browse":    "Elige la carpeta de instalación",
		"prog":      "Instalando…",
		"done":      "Instalado",
		"doneat":    "La aplicación está instalada en:",
		"warnmodel": "No se pudo descargar el modelo: la aplicación lo ofrecerá al primer inicio.",
		"finish":    "FINALIZAR",
		"err":       "Error de instalación",
		"retry":     "REINTENTAR",
		"back":      "ATRÁS",
		"nodir":     "Indica la carpeta de instalación",
		"webview":   "El instalador necesita Microsoft WebView2 Runtime (incluido en Windows 11).\nSe abrirá la página de descarga.",
	},
	"it": {
		"swon":      "ON",
		"swoff":     "OFF",
		"title":     appid.Name + " — Installazione",
		"tagline":   "Voce → testo nella posizione del cursore. Tutto in locale e offline.",
		"path":      "Cartella di installazione",
		"shortcut":  "Collegamento nel menu Start",
		"autorun":   "Avvia con Windows",
		"launch":    "Avvia dopo l'installazione",
		"install":   "INSTALLA",
		"update":    "AGGIORNA",
		"updnote":   "È stata trovata un'installazione esistente: impostazioni e modelli scaricati verranno mantenuti.",
		"model":     "Modello di riconoscimento",
		"mnone":     "Non scaricare (lo scarichi al primo avvio)",
		"browse":    "Scegli la cartella di installazione",
		"prog":      "Installazione…",
		"done":      "Installato",
		"doneat":    "L'applicazione è installata in:",
		"warnmodel": "Il modello non è stato scaricato: l'app lo proporrà al primo avvio.",
		"finish":    "FINE",
		"err":       "Errore di installazione",
		"retry":     "RIPROVA",
		"back":      "INDIETRO",
		"nodir":     "Indica la cartella di installazione",
		"webview":   "L'installer richiede Microsoft WebView2 Runtime (incluso in Windows 11).\nSi aprirà la pagina di download.",
	},
	"pl": {
		"swon":      "ON",
		"swoff":     "OFF",
		"title":     appid.Name + " — Instalacja",
		"tagline":   "Głos → tekst w miejscu kursora. Całkowicie lokalnie i offline.",
		"path":      "Folder instalacji",
		"shortcut":  "Skrót w menu Start",
		"autorun":   "Uruchamiaj z Windows",
		"launch":    "Uruchom po instalacji",
		"install":   "ZAINSTALUJ",
		"update":    "ZAKTUALIZUJ",
		"updnote":   "Znaleziono istniejącą instalację — ustawienia i pobrane modele zostaną zachowane.",
		"model":     "Model rozpoznawania",
		"mnone":     "Nie pobieraj (pobierzesz przy pierwszym uruchomieniu)",
		"browse":    "Wybierz folder instalacji",
		"prog":      "Instalacja…",
		"done":      "Zainstalowano",
		"doneat":    "Aplikacja została zainstalowana w:",
		"warnmodel": "Nie udało się pobrać modelu — aplikacja zaproponuje go przy pierwszym uruchomieniu.",
		"finish":    "GOTOWE",
		"err":       "Błąd instalacji",
		"retry":     "PONÓW",
		"back":      "WSTECZ",
		"nodir":     "Wskaż folder instalacji",
		"webview":   "Instalator wymaga Microsoft WebView2 Runtime (dołączony do Windows 11).\nZa chwilę otworzy się strona pobierania.",
	},
	"en": {
		"swon":      "ON",
		"swoff":     "OFF",
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
		"retry":     "RETRY",
		"back":      "BACK",
		"nodir":     "Choose the installation folder",
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
		"{{SWON}}":         tr("swon"),
		"{{SWOFF}}":        tr("swoff"),
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
		"{{RETRY}}":        tr("retry"),
		"{{BACK}}":         tr("back"),
		"{{NODIR_JS}}":     template_jsstr(tr("nodir")),
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
