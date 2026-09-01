package main

import (
	"strings"

	"fmt"
	"sync/atomic"

	"holdtotype/internal/appid"
)

var procGetUserDefaultUILanguage = kernel32.NewProc("GetUserDefaultUILanguage")

var curLang atomic.Value

var uiLangs = []string{"en", "ru", "uk", "de", "fr", "es", "it", "pl"}

func systemLang() string {
	langID, _, _ := procGetUserDefaultUILanguage.Call()
	switch langID & 0x3FF {
	case 0x19:
		return "ru"
	case 0x22:
		return "uk"
	case 0x07:
		return "de"
	case 0x0C:
		return "fr"
	case 0x0A:
		return "es"
	case 0x10:
		return "it"
	case 0x15:
		return "pl"
	}
	return "en"
}

func validUILang(l string) bool {
	for _, v := range uiLangs {
		if v == l {
			return true
		}
	}
	return false
}

func initLang(setting string) {
	if validUILang(setting) {
		curLang.Store(setting)
		return
	}
	curLang.Store(systemLang())
}

func lang() string {
	if v, ok := curLang.Load().(string); ok {
		return v
	}
	return "en"
}

var nameReplacer = strings.NewReplacer(
	"{app}", appid.Name,
	"{exe}", appid.Exe,
	"{setup}", appid.SetupExe,
	"{log}", appid.LogFile,
	"{zip}", appid.Portable,
)

func expandName(s string) string {
	if !strings.ContainsRune(s, '{') {
		return s
	}
	return nameReplacer.Replace(s)
}

func tr(key string) string {
	if s, ok := msgs[lang()][key]; ok {
		return expandName(s)
	}
	if s, ok := msgs["en"][key]; ok {
		return expandName(s)
	}
	return key
}

func trf(key string, args ...any) string {
	return fmt.Sprintf(tr(key), args...)
}

func strS(key string) string {
	if v := settingsStrings[lang()][key]; v != "" {
		return expandName(v)
	}
	return expandName(settingsStrings["en"][key])
}

var msgs = map[string]map[string]string{
	"ru": {
		"app.name":              "{app}",
		"already.running":       "Приложение уже запущено (иконка в трее).",
		"err.title":             "{app} — ошибка",
		"cfg.err.title":         "{app} — ошибка конфигурации",
		"err.details":           "\n\nПодробности в {log}",
		"err.hook":              "Не удалось установить хук клавиатуры: %s",
		"err.mic":               "Микрофон: %s",
		"err.hotkey.cfg":        "Сочетание клавиш в config.json: %s",
		"err.model.notfound":    "Файл модели не найден: %s\nСоберите проект (.\\build.ps1) или исправьте путь \"model\" в config.json",
		"err.server.repeat":     "whisper-server аварийно завершается раз за разом — подробности в {log}",
		"err.server.dead":       "сервер распознавания %s не отвечает (server_url/server_autostart в config.json)",
		"err.server.timeout":    "whisper-server не ответил за %s",
		"err.server.start":      "whisper-server завершился при старте (см. лог)",
		"err.srv.noaddr":       "удалённый сервер распознавания не настроен — задайте адрес в настройках",
		"err.webview":           "Для окна настроек нужен Microsoft WebView2 Runtime (входит в Windows 11).\nСейчас откроется страница загрузки — установите его и откройте настройки снова.",
		"status.loading":        "Загрузка модели…",
		"status.nomodel":        "Модель распознавания не скачана — выберите её в настройках",
		"state.loaded.none": "ничего не загружено",
		"state.week": "%d диктовок · %d знаков",
		"snd.ok": "уровень в норме",
		"snd.quiet": "тихо — говорите ближе к микрофону",
		"snd.clipped": "слишком громко, звук обрезался",
		"snd.silent": "тишина в записи",
		"status.parked": "Движок выгружен — нажмите сочетание, чтобы разбудить",
		"status.nomodel.lang": "Для языка %s не установлена модель %s — откройте «Языки и модели»",
		"status.ready":          "Готов — зажмите %s и говорите",
		"status.recording":      "Идёт запись…",
		"status.transcribing":   "Распознаю…",
		"status.disabled":       "Выключено",
		"status.server.restart": "Сервер распознавания упал, перезапускаю…",
		"status.cfg.err":        "Ошибка в config.json (см. лог)",
		"menu.settings":         "Настройки…",
		"menu.enable":           "Включить",
		"menu.disable":          "Выключить",
		"menu.open.config":      "Открыть config.json",
		"menu.open.log":         "Открыть лог",
		"menu.about":            "О приложении",
		"menu.quit":             "Выход",
		"menu.lastcopy":         "Копировать последний результат",
		"ov.copied":             "Скопировано в буфер обмена",
		"ov.kept":               "Отменено — текст в «Последнем результате»",
		"ov.llm.skipped":        "Вставлено без профиля «%s»",
		"fd.title":              "Фокус сменился — вставить?",
		"fd.here":               "Вставить сюда",
		"fd.copy":               "Копировать",
		"fd.keep":               "Оставить",
		"ov.speak":              "Говорите…",
		"err.port.busy":         "Порт %d занят другой программой. Смените порт в настройках.",
		"err.server.launch":     "Не удалось запустить %s — проверьте путь к серверу в разделе «Система»",
		"err.generic":           "Не получилось — подробности в логе",
		"err.cancelled":         "Отменено",
		"err.disk.full":         "На диске нет места",
		"err.file.denied":       "Нет доступа к файлу — закройте программу, которая его держит",
		"err.file.missing":      "Файл не найден",
		"err.answer":            "Сервер ответил непонятным образом — попробуйте позже",
		"err.net.cert":          "Соединение не защищено — проверьте дату и антивирус",
		"err.net.down":          "Не удалось соединиться — проверьте интернет",
		"err.net.timeout":       "Сервер не ответил вовремя — попробуйте ещё раз",
		"err.net.dns":           "Нет связи с %s — проверьте интернет",
		"ov.esc":                "1…9 · Enter · Esc — отмена",
		"ov.left":               "осталось %d с",
		"ov.transcribing":       "Распознаю",
		"ov.asking":             "Жду ответа",
		"ov.inserted":           "Вставлено: %d симв.",
		"ov.err.mic":            "Микрофон недоступен — проверьте устройство в настройках",
		"ov.err.recognize":      "Ошибка распознавания (см. лог)",
		"ov.err.paste":          "Не вставилось — текст в «Последнем результате»",
		"ov.moved":              "Окно сменилось — текст в буфере обмена",
		"copy.ok":               "Скопировано",
		"copy.none":             "Нечего копировать",
		"copy.fail":             "Не удалось скопировать: %s",
		"mic.busy":              "Идёт запись, проверка невозможна", "mic.check.ok": "Слышно хорошо: пик %.0f дБ, речь на %.0f%% записи",
		"mic.check.quiet": "Слишком тихо: пик %.0f дБ — прибавьте громкость микрофона в Windows или сядьте ближе", "mic.check.clipped": "Перегруз: обрезано %.1f%% отсчётов — убавьте громкость микрофона", "mic.check.silent": "Речи не слышно — проверьте, тот ли выбран микрофон и не выключен ли он",
		"ov.quiet": "Слишком тихо, почти ничего не слышно", "ov.clipped": "Перегруз — звук обрезан",
		"ov.cmd.cancelled":       "Отменено голосом",
		"ov.silence":             "Тишина — текст не распознан",
		"ov.tooshort":            "Слишком коротко — держите клавиши дольше",
		"ov.notranslate":         "Активная модель не переводит — вставлен исходный текст",
		"ov.engine.fallback":     "Второй движок не поднялся — распознаю текущим",
		"route.speech":           "Речь на %s",
		"route.other":            "Прочие языки",
		"route.translate":        "Перевод",
		"route.lang.auto":        "любом языке",
		"route.why.language":     "точнее на этом языке, с пунктуацией",
		"route.why.otherlang":    "99 языков",
		"route.why.translate":    "переводит только Whisper",
		"route.why.notinstalled": "русская модель не установлена",
		"route.why.unknownlang":  "язык не задан — только Whisper определит его сам",
		"route.why.forced":       "выбрано вручную в config.json",
		"status.line":            "Готов · %s · %.1f ГБ свободно", "state.ram.free": "свободно %d МБ",
		"ago.now":                "только что",
		"ago.min":                "%d мин назад",
		"ago.hour":               "%d ч назад",
		"chars":                  "%d символов",
		"inserted.into":          "вставлено в %s",
		"punct.prompt":           "Расставь знаки препинания и заглавные буквы. Не меняй слова, не переводи, не добавляй ничего от себя. Верни только исправленный текст.",
		"err.sherpa.notfound":    "Распознаватель sherpa не найден: %s",
		"err.sherpa.start":       "sherpa-server завершился при старте (см. лог)",
		"err.sherpa.translate":   "эта модель не умеет переводить",
		"err.sherpa.model":       "Файл модели не найден: %s — скачайте её в настройках или исправьте sherpa_model в config.json",
		"ov.server.loading":      "Сервер ещё загружается",
		"ov.cancelled":           "Отменено",
		"ov.editing":             "Редактирую: %s",
		"ov.translating":         "Перевожу",
		"ov.llm.needed":          "Для этого языка нужен LLM-модуль",
		"td.title":               "Переводить на:",
		"td.plain":               "Без перевода",
		"cap.title":              "{app} — сочетание клавиш",
		"cap.prompt":             "Нажмите новое сочетание клавиш\n\nсейчас: %s   ·   Esc — отмена",
		"cap.selected":           "Выбрано: %s",
		"cap.cancelled":          "Отменено",
		"hk.taken":               "Сочетание %s занято Windows: %s. Диктовка может не начаться",
		"hk.lock":                "блокировка компьютера",
		"hk.desktop":             "показать рабочий стол",
		"hk.explorer":            "проводник",
		"hk.run":                 "окно «Выполнить»",
		"hk.settings":            "параметры Windows",
		"hk.search":              "поиск",
		"hk.center":              "центр уведомлений",
		"hk.menu":                "меню опытного пользователя",
		"hk.clipboard":           "журнал буфера обмена",
		"hk.gamebar":             "игровая панель",
		"hk.voice":               "голосовой ввод Windows",
		"hk.project":             "проецирование на экран",
		"hk.tasks":               "представление задач",
		"hk.layout":              "смена раскладки",
		"hk.newdesktop":          "новый рабочий стол",
		"hk.closedesktop":        "закрыть рабочий стол",
		"hk.snip":                "снимок экрана",
		"hk.switch":              "переключение окон",
		"hk.close":               "закрыть окно",
		"hk.cycle":               "перебор окон",
		"hk.start":               "меню «Пуск»",
		"hk.taskmgr":             "диспетчер задач",
		"hk.secure":              "экран безопасности",
		"err.hotkey.dup":         "Сочетание «%s» назначено дважды — хоткеи не должны совпадать",
		"cfg.err.recovered":      "config.json повреждён (%s).\nФайл сохранён как %s, настройки сброшены к значениям по умолчанию.",
		"err.disk.space":         "мало места на диске: свободно %d МБ, нужно ~%d МБ",
		"err.save":               "не удалось сохранить настройки: %s — оставил прежние",
		"err.port":               "порт %d не подходит: нужен номер от 1024 до 65535",
		"err.nolangs":            "оставьте хотя бы один язык в списке для вопроса о переводе",
		"ov.mic.lost":            "Микрофон отключён — запись прервана",
		"err.hash":               "скачанный файл повреждён — попробуйте ещё раз",
		"models.check.ok":        "Проверено моделей: %d — все файлы целы",
		"models.check.none":      "Проверять нечего — нет установленных моделей с эталонными хешами",
		"models.check.bad":       "Повреждены файлы: %s — скачайте модель заново",
		"hist.insert.gone":       "запись не найдена",
		"ov.aim": "Щёлкните поле, куда вставить · Esc — отмена",
		"hist.aim.armed": "щёлкните поле, куда вставить",
		"hist.aim.busy": "уже жду щелчка",
		"hist.aim.off": "вставка отменена",
		"hist.insert.nowin":      "некуда вставлять — текст скопирован в буфер",
		"hist.insert.ok":         "вставлено в «%s»",
		"lists.bad":              "файл не подходит",
		"lists.saved":            "сохранено в %s",
		"lists.added":            "добавлено: %d, пропущено: %d",
		"lists.save.title":       "Куда сохранить списки",
		"lists.open.title":       "Откуда загрузить списки",
		"un.title":               "{app} — удаление",
		"un.confirm":             "Удалить {app} с этого компьютера?",
		"un.data":                "Удалить также настройки и скачанные модели?",
		"un.done":                "{app} удалён.",
		"model.switching":        "Переключаю модель — распознаватель перезапускается…",
		"srv.restarting":         "Перезапускаю распознаватель с новыми настройками…",
		"model.del.active":       "Нельзя удалить активную модель",
		"model.del.ok":           "Модель удалена",
		"about.text": "{app} %s\n\n" +
			"Голос → текст в позицию курсора.\n" +
			"Поставьте курсор в поле ввода, зажмите %s, скажите фразу, отпустите — текст вставится сам.\n\n" +
			"Распознавание: whisper.cpp, полностью локально и офлайн — звук не покидает компьютер.\n" +
			"Модель: %s (язык: %s)\n\n" +
			"Настройки: клик по иконке в трее или config.json.\n" +
			"Логи: {log} (не больше ~2 МБ).",
	},
	"en": {
		"app.name":              "{app}",
		"already.running":       "The application is already running (tray icon).",
		"err.title":             "{app} — error",
		"cfg.err.title":         "{app} — configuration error",
		"err.details":           "\n\nDetails in {log}",
		"err.hook":              "Failed to install keyboard hook: %s",
		"err.mic":               "Microphone: %s",
		"err.hotkey.cfg":        "Hotkey in config.json: %s",
		"err.model.notfound":    "Model file not found: %s\nBuild the project (.\\build.ps1) or fix the \"model\" path in config.json",
		"err.server.repeat":     "whisper-server keeps crashing — see {log}",
		"err.server.dead":       "recognition server %s is not responding (server_url/server_autostart in config.json)",
		"err.server.timeout":    "whisper-server did not respond within %s",
		"err.server.start":      "whisper-server exited during startup (see log)",
		"err.srv.noaddr":       "the remote recognition server is not set up — set its address in the settings",
		"err.webview":           "The settings window requires Microsoft WebView2 Runtime (bundled with Windows 11).\nThe download page will open now — install it and open Settings again.",
		"status.loading":        "Loading model…",
		"status.nomodel":        "No recognition model downloaded — pick one in Settings",
		"state.loaded.none": "nothing loaded",
		"state.week": "%d dictations · %d characters",
		"snd.ok": "sound level fine",
		"snd.quiet": "quiet — speak closer to the microphone",
		"snd.clipped": "too loud, the sound clipped",
		"snd.silent": "silence in the recording",
		"status.parked": "Engine unloaded — press the shortcut to wake it",
		"status.nomodel.lang": "For %s the %s model is not installed — open Languages & models",
		"status.ready":          "Ready — hold %s and speak",
		"status.recording":      "Recording…",
		"status.transcribing":   "Transcribing…",
		"status.disabled":       "Disabled",
		"status.server.restart": "Recognition server crashed, restarting…",
		"status.cfg.err":        "Error in config.json (see log)",
		"menu.settings":         "Settings…",
		"menu.enable":           "Enable",
		"menu.disable":          "Disable",
		"menu.open.config":      "Open config.json",
		"menu.open.log":         "Open log",
		"menu.about":            "About",
		"menu.quit":             "Quit",
		"menu.lastcopy":         "Copy last result",
		"ov.copied":             "Copied to clipboard",
		"ov.kept":               "Cancelled — text kept in Last Result",
		"ov.llm.skipped":        "Inserted without the \"%s\" profile",
		"fd.title":              "Focus changed — insert?",
		"fd.here":               "Insert here",
		"fd.copy":               "Copy",
		"fd.keep":               "Keep it",
		"ov.speak":              "Speak…",
		"err.port.busy":         "Port %d is busy: another program is listening on it. Change the port in the settings.",
		"err.server.launch":     "Could not start %s — check the server path under System",
		"err.generic":           "It did not work — details are in the log",
		"err.cancelled":         "Cancelled",
		"err.disk.full":         "The disk is full",
		"err.file.denied":       "No access to the file — close the program holding it",
		"err.file.missing":      "File not found",
		"err.answer":            "The server answered in a way we do not understand — try later",
		"err.net.cert":          "The connection is not secure — check the date and the antivirus",
		"err.net.down":          "Could not connect — check the internet",
		"err.net.timeout":       "The server did not answer in time — try again",
		"err.net.dns":           "No connection to %s — check the internet",
		"ov.esc":                "1…9 · Enter · Esc cancels",
		"ov.left":               "%d s left",
		"ov.transcribing":       "Transcribing",
		"ov.asking":             "Waiting for your answer",
		"ov.inserted":           "Inserted: %d chars",
		"ov.err.mic":            "Microphone unavailable — check the device in Settings",
		"ov.err.recognize":      "Recognition error (see log)",
		"ov.err.paste":          "Not pasted — the text is in Last Result",
		"ov.moved":              "The window changed — the text is on the clipboard",
		"copy.ok":               "Copied",
		"copy.none":             "Nothing to copy",
		"copy.fail":             "Could not copy: %s",
		"mic.busy":              "A dictation is in progress, cannot check now", "mic.check.ok": "Sounds good: peak %.0f dB, speech in %.0f%% of the recording",
		"mic.check.quiet": "Too quiet: peak %.0f dB — raise the microphone level in Windows or sit closer", "mic.check.clipped": "Clipping: %.1f%% of samples cut off — lower the microphone level", "mic.check.silent": "No speech heard — check that the right microphone is picked and not muted",
		"ov.quiet": "Too quiet, almost nothing was heard", "ov.clipped": "Clipping — the sound was cut off",
		"ov.cmd.cancelled":       "Cancelled by voice",
		"ov.silence":             "Silence — nothing recognized",
		"ov.tooshort":            "Too short — hold the keys longer",
		"ov.notranslate":         "The active model cannot translate — inserted as recognized",
		"ov.engine.fallback":     "The other engine did not start — using the current one",
		"route.speech":           "Speech in %s",
		"route.other":            "Other languages",
		"route.translate":        "Translation",
		"route.lang.auto":        "any language",
		"route.why.language":     "more accurate here, with punctuation",
		"route.why.otherlang":    "99 languages",
		"route.why.translate":    "only Whisper translates",
		"route.why.notinstalled": "the Russian model is not installed",
		"route.why.unknownlang":  "no language set — only Whisper detects it",
		"route.why.forced":       "forced in config.json",
		"status.line":            "Ready · %s · %.1f GB free", "state.ram.free": "%d MB free",
		"ago.now":                "just now",
		"ago.min":                "%d min ago",
		"ago.hour":               "%d h ago",
		"chars":                  "%d characters",
		"inserted.into":          "inserted into %s",
		"punct.prompt":           "Add punctuation and capital letters. Do not change the words, do not translate, do not add anything. Return only the corrected text.",
		"err.sherpa.notfound":    "sherpa recognizer not found: %s",
		"err.sherpa.start":       "sherpa-server exited during startup (see log)",
		"err.sherpa.translate":   "this model cannot translate",
		"err.sherpa.model":       "Model file not found: %s — download it in Settings or fix sherpa_model in config.json",
		"ov.server.loading":      "Server is still loading",
		"ov.cancelled":           "Cancelled",
		"ov.editing":             "Editing: %s",
		"ov.translating":         "Translating",
		"ov.llm.needed":          "This language requires the LLM module",
		"td.title":               "Translate to:",
		"td.plain":               "No translation",
		"cap.title":              "{app} — shortcut",
		"cap.prompt":             "Press a new key combination\n\nnow: %s   ·   Esc — cancel",
		"cap.selected":           "Selected: %s",
		"cap.cancelled":          "Cancelled",
		"hk.taken":               "%s is taken by Windows: %s. Dictation may never start",
		"hk.lock":                "lock the computer",
		"hk.desktop":             "show the desktop",
		"hk.explorer":            "File Explorer",
		"hk.run":                 "the Run box",
		"hk.settings":            "Windows Settings",
		"hk.search":              "search",
		"hk.center":              "notification centre",
		"hk.menu":                "the power user menu",
		"hk.clipboard":           "clipboard history",
		"hk.gamebar":             "Game Bar",
		"hk.voice":               "Windows voice typing",
		"hk.project":             "projecting to a screen",
		"hk.tasks":               "Task View",
		"hk.layout":              "switching the keyboard layout",
		"hk.newdesktop":          "a new virtual desktop",
		"hk.closedesktop":        "closing the virtual desktop",
		"hk.snip":                "the screenshot tool",
		"hk.switch":              "switching windows",
		"hk.close":               "closing the window",
		"hk.cycle":               "cycling windows",
		"hk.start":               "the Start menu",
		"hk.taskmgr":             "Task Manager",
		"hk.secure":              "the security screen",
		"err.hotkey.dup":         "The \"%s\" shortcut is assigned twice — hotkeys must be unique",
		"cfg.err.recovered":      "config.json is corrupted (%s).\nThe file was saved as %s and settings were reset to defaults.",
		"err.disk.space":         "low disk space: %d MB free, ~%d MB needed",
		"err.save":               "settings not saved: %s — the old ones are kept",
		"err.port":               "port %d will not do: pick a number between 1024 and 65535",
		"err.nolangs":            "leave at least one language for the translation question",
		"ov.mic.lost":            "The microphone is gone — recording stopped",
		"err.hash":               "the downloaded file is damaged — try again",
		"models.check.ok":        "Models checked: %d — all files intact",
		"models.check.none":      "Nothing to check — no installed model has a reference hash",
		"models.check.bad":       "Damaged files: %s — download the model again",
		"hist.insert.gone":       "entry not found",
		"ov.aim": "Click the field to paste into · Esc cancels",
		"hist.aim.armed": "click the field to paste into",
		"hist.aim.busy": "already waiting for a click",
		"hist.aim.off": "paste called off",
		"hist.insert.nowin":      "nowhere to paste — the text is on the clipboard",
		"hist.insert.ok":         "pasted into “%s”",
		"lists.bad":              "this file does not fit",
		"lists.saved":            "saved to %s",
		"lists.added":            "added: %d, skipped: %d",
		"lists.save.title":       "Where to save the lists",
		"lists.open.title":       "Which file to load",
		"un.title":               "{app} — Uninstall",
		"un.confirm":             "Remove {app} from this computer?",
		"un.data":                "Also delete settings and downloaded models?",
		"un.done":                "{app} has been removed.",
		"model.switching":        "Switching model — the recognizer is restarting…",
		"srv.restarting":         "Restarting the recognizer with the new settings…",
		"model.del.active":       "Cannot delete the active model",
		"model.del.ok":           "Model deleted",
		"about.text": "{app} %s\n\n" +
			"Voice → text at the cursor position.\n" +
			"Place the cursor in any input field, hold %s, say a phrase, release — the text is inserted automatically.\n\n" +
			"Recognition: whisper.cpp, fully local and offline — audio never leaves your computer.\n" +
			"Model: %s (language: %s)\n\n" +
			"Settings: click the tray icon, or edit config.json.\n" +
			"Logs: {log} (never exceeds ~2 MB).",
	},
}

var settingsStrings = map[string]map[string]string{
	"ru": {
		"S_TITLE":          "{app} — настройки",
		"S_TR_DEFAULT":     "Изменить язык текстового вывода",
		"S_TR_TARGET":      "Язык текстового вывода по умолчанию",
		"S_SRCLANG_SUB": "на нём вы говорите; он определяет модель распознавания",
		"S_TR_LANGS_SUB": "эти языки будут кнопками на плашке при вставке",
		"S_TR_UNAVAIL": "недоступно — %s не умеет переводить",
		"S_TR_LOCK": "%s нельзя убрать из списка — он выбран языком текстового вывода по умолчанию. Выберите другой язык по умолчанию, и тогда %s можно будет исключить.",
		"S_TR_LOCK_OK": "Понятно",
		"S_TR_ONE": "В списке отмечено несколько языков, но без вопроса перевод всегда пойдёт в один — %s (язык текстового вывода по умолчанию). Остальные останутся отмеченными, но будут отключены.",
		"S_TR_NOMODEL": "%s не умеет переводить. Если продолжить, перевод будет выключен и недоступен, пока работает эта модель.",
		"S_TR_CONFIRM": "Подтвердить",
		"S_TR_ASK":         "Спрашивать язык текстового вывода",
		"S_TR_ASK_NEVER":   "Не спрашивать — переводить сразу",
		"S_TR_ASK_ALWAYS":  "Спрашивать каждый раз",
		"S_TR_ASK_TIMEOUT": "Спрашивать, с таймаутом",
		"S_TR_SECONDS":     "Таймаут, сек",
		"S_TR_LANGS":       "Языки в диалоге",
		"S_DICT_HINT":      "Термины, имена и аббревиатуры через запятую — подсказка слуху, не команды. Работает для Whisper; на русскую речь через GigaAM не влияет. Набор по умолчанию меняется вместе с языком распознавания, пока вы не впишете своё.",
		"S_LLM_HINT":       "Отмеченные профили применяются по очереди, сверху вниз, при обычной диктовке. Ничего не отмечено — текст вставляется как есть.",
		"S_UPDATED":        "Дата последнего обновления модели",
		"S_FIT_OK":         "поместится",
		"S_FIT_WARN":       "впритык",
		"S_FIT_BAD":        "не хватит памяти",
		"S_RAM":            "Память компьютера:",
		"S_HF_PH":          "Название модели — например, qwen2.5 instruct",
		"S_NO_LLM":         "Пока не установлено ни одной модели — найдите и скачайте в поле поиска ниже.",
		"S_NO_LLM_PROF":    "Промпты станут доступны после установки модели — блок «Модель» выше на этой вкладке.",
		"S_PROF_EDIT":      "Редактировать",
		"S_PROF_CLOSE":     "Свернуть",
		"S_CONFIRM_DEL":    "Удалить модель «%s»? Её можно будет скачать заново.",
		"S_DEL_ACTIVE":     "Удалить активную модель «%s»? Распознавание остановится, пока вы не выберете другую — скачать её можно тут же.",
		"S_WIZ_NEED_MODEL": "Сначала скачайте модель — без неё распознавать нечем",
		"S_FREE":           "свободно",
		"S_SUB_PROMPTS":    "Промпты",
		"S_SUB_DICT":       "Словарь",
		"S_PROF_ADD":       "Добавить",
		"S_PROF_NAME":      "Имя",
		"S_PROF_PROMPT":    "Промпт",
		"S_PROF_TEST":      "Проверка",
		"S_HOTKEY":         "Сочетание клавиш",
		"S_CHANGE":         "Изменить…",
		"S_UILANG":         "Язык интерфейса",
		"S_AUTO":           "Как в системе",
		"S_BEEP":           "Звуковые сигналы записи",
		"S_SEC_SOUND":      "Звук",
		"S_SOUND":          "Звук сигналов",
		"S_SND_SPEECH":     "Голос Windows",
		"S_SND_CHIME":      "Колокольчик",
		"S_SND_SOFT":       "Мягкий",
		"S_SND_MARIMBA":    "Маримба",
		"S_SND_BLIP":       "Блип",
		"S_SND_POP":        "Поп",
		"S_AUTOENTER":      "Нажимать Enter после вставки (автоотправка)",
		"S_RESTORE":        "Восстанавливать буфер обмена после вставки",
		"S_OVERLAY":        "Показывать плашку", "S_OVERLAY_SUB": "во время диктовки на экране видно, что идёт запись",
		"S_NAV_HISTORY":    "История", "S_HIST_ON": "Хранить историю диктовок", "S_HIST_ON_SUB": "только текст, на этом компьютере; звук не сохраняется никогда",
		"S_HIST_DAYS": "Сколько дней хранить", "S_HIST_MAX": "Сколько записей хранить",
		"S_HIST_SKIP": "Не записывать из программ", "S_HIST_SKIP_SUB": "через запятую: keepass.exe, 1password.exe", "S_SKIP_ADD_DLG": "Добавление программы", "S_SKIP_EDIT_DLG": "Изменение программы", "S_SKIP_NAME": "Имя программы", "S_SKIP_NAME_SUB": "Файл без пути: keepass.exe. Звёздочка в конце ловит все версии: 1password*", "S_SKIP_OPEN": "Открытые сейчас программы", "S_SKIP_REFRESH": "Обновить список", "S_SKIP_PICKED": "Выбрано %d из %d", "S_SKIP_NONE": "Ничего не выбрано", "S_SKIP_EMPTY": "Список пуст — история пишется из всех программ", "S_SKIP_ADD_BTN": "Добавить программу", "S_SKIP_HINT": "Что надиктовано в эти программы, в историю не попадёт. Сама вставка работает как обычно.",
		"S_HIST_LIST": "Записи", "S_HIST_CLEAR": "Очистить", "S_HIST_TILL": "до %s", "S_HIST_TILL1": "до завтра", "S_HIST_TILL_FULL": "Удалится %s — срок хранения %s", "S_HIST_LIST_HINT": "Что было надиктовано: скопировать, вставить в любое окно или удалить.", "S_HIST_COPY": "Копировать",
		"S_HIST_KEEP":       "Сколько хранить",
		"S_UNIT_MIN":        "минут",
		"S_UNIT_HOUR":       "часов",
		"S_UNIT_DAY":        "дней",
		"S_HIST_FIND": "Найти в истории…", "S_HIST_EMPTY": "Истории пока нет", "S_HIST_ASK": "Удалить всю историю диктовок?",
		"S_SEC_CMD": "Голосовые команды", "S_CMD_HINT": "Сказанное вслух превращается в перенос строки, знак или отмену вместо того, чтобы попасть в текст. Проверяются целыми словами, применяются сверху вниз, уже после замен.",
		"S_CMD_ADD": "Добавить команду", "S_CMD_PRESET": "Добавить стандартные", "S_CMD_PH": "фраза, которую вы произнесёте",
		"S_CMD_NEWLINE": "перенос строки", "S_CMD_PARAGRAPH": "новый абзац", "S_CMD_TEXT": "подставить текст", "S_CMD_CANCEL": "отменить диктовку",
		"S_CMD_TEXT_PH": "что подставить", "S_CMD_EMPTY": "Команд пока нет", "S_CMD_DEL": "Удалить команду",
		"S_CMD_P_NEWLINE": "новая строка", "S_CMD_P_PARAGRAPH": "новый абзац", "S_CMD_P_CANCEL": "отмена",
		"S_SEC_REPLACE": "Замены после распознавания", "S_REPLACE_HINT": "Что услышано неправильно, заменяется на правильное — сразу после распознавания, до промптов и перевода. Порядок сверху вниз.",
		"S_REPL_WHOLE_FULL": "Только целые слова", "S_REPL_CASE_FULL": "Учитывать регистр", "S_CMD_ACTION": "Действие",
		"S_FM_ADD": "Добавить",
		"S_TIP_REPL_LANG": "Правило работает, только когда диктовка идёт на выбранном языке. «все языки» — работает всегда.",
		"S_TIP_REPL_CASE": "Большие и маленькие буквы различаются: «гит» и «Гит» — разные слова. Выключено — регистр не важен.",
		"S_TIP_REPL_WHOLE": "Замена срабатывает, только если текст стоит отдельным словом. Выключено — ищется и внутри других слов: «кот» найдётся в «котлета».",
		"S_TIP_CMD_ACTION": "Что произойдёт, когда вы произнесёте фразу: перенос строки, новый абзац, подстановка своего текста или отмена диктовки.",
		"S_LIST_FILTER_PH": "найти…",
		"S_REPL_DEL": "Удалить замену",
		"S_LIST_NOTHING": "Ничего не найдено: «%s»",
		"S_FM_T_REPL_ADD": "Добавление замены", "S_FM_T_REPL_EDIT": "Изменение замены",
		"S_FM_T_CMD_ADD": "Добавление команды", "S_FM_T_CMD_EDIT": "Изменение команды",
		"S_MT_DEL": "Удаление модели", "S_MT_DEL_PROMPT": "Удаление промпта", "S_MT_DL": "Загрузка модели",
		"S_MT_TR_OFF": "Отключение перевода", "S_MT_TR_ONE": "Перевод без вопроса", "S_MT_TR_LOCK": "Язык вывода по умолчанию",
		"S_MT_REMOTE": "Удалённый сервер", "S_MT_POST": "Внешний сервер", "S_MT_HIST": "Очистка истории",
		"S_MT_RESET": "Сброс настроек", "S_MT_EXE": "Путь к серверу",
		"S_DICT_ADD": "Добавить слово", "S_FM_T_DICT_ADD": "Добавление слова", "S_DICT_EMPTY": "Слов пока нет",
		"S_DICT_ADD_PH": "слово или несколько через запятую",
		"S_DICT_NOMODEL": "Текущая модель %s не поддерживает словарь — словарь читают только модели Whisper.",
		"S_OV_FREE": "Своё место", "S_OV_FREE_SUB": "плашку можно перетащить куда угодно",
		"S_OVPOS_DRAG_SUB": "тащите плашку мышью — она встанет где угодно",
		"S_OVMON_N": "Экран %d",
		"S_POST_ENABLE": "Включить постобработку",
		"S_API_SUM_URL": "адрес", "S_API_SUM_MODEL": "модель", "S_API_SUM_KEY": "ключ", "S_API_SUM_TIMEOUT": "ожидание",
		"S_API_SUM_STATE": "состояние", "S_API_NO_MODEL": "не указана",
		"S_API_NONE": "не настроен — постобработка идёт локально",
		"S_POSTAPI_SETUP": "Настроить", "S_API_EDIT": "Изменить", "S_API_KEY_DEL": "Удалить ключ", "S_API_DLG": "Внешний сервер",
		"S_LLM_CATALOG": "Каталог моделей", "S_LLM_BLOCK": "Установленные модели", "S_LLM_NONE_HINT": "Ни одной модели не установлено — скачайте найденную стрелкой, и она появится здесь", "S_LLM_IN_MEM": "в памяти", "S_LLM_ON_DISK": "на диске", "S_LLM_EJECT": "Выгрузить из памяти", "S_LLM_FOUND": "найдено %d", "S_LLM_NOSEARCH": "поиск не запускали", "S_LLM_SEARCH_HINT": "Введите название модели и нажмите «Найти»", "S_LLM_PICK_WAIT": "Будет доступен, когда модель скачается", "S_LLM_INSTALLED": "установлены",
		"S_LLM_SUM_MODEL": "модель", "S_LLM_SUM_SIZE": "размер", "S_LLM_SUM_COUNT": "установлено", "S_LLM_SUM_RAM": "память",
		"S_DLG_CLOSE": "Закрыть", "S_LLM_NOPICK": "не выбрана", "S_NO_PROMPTS": "Промптов пока нет", "S_PROF_DRAG": "перетащите, чтобы изменить порядок",
		"S_PROF_NAME_PH": "как назвать промпт", "S_PROF_TEST_PH": "напишите фразу для проверки",
		"S_PF_NEW": "Новый промпт", "S_PF_EDIT": "Изменение промпта",
		"S_POST_NO_MODEL": "включена, но модель не выбрана", "S_POST_NO_API": "включена, но сервер не настроен", "S_POST_BAD": "сервер не ответил: %s", "S_POST_NO_PROMPT": "включена, но не отмечен ни один промпт", "S_API_TEST": "Тест соединения", "S_API_TEST_RUN": "Проверяю…", "S_API_TEST_OK": "Сервер ответил", "S_API_CLEAR": "Очистить", "S_API_CLEAR_ASK": "Удалить адрес, модель и ключ внешнего сервера? Постобработка вернётся на локальную модель.", "S_RAM_AVAIL": "Доступно памяти %s ГБ из %s ГБ", "S_RAM_OF": "%s ГБ из %s ГБ",
		"S_REPL_ADD": "Добавить замену", "S_REPL_FROM_PH": "гит хаб", "S_REPL_TO_PH": "GitHub",
		"S_REPL_WHOLE": "целые слова", "S_REPL_CASE": "регистр", "S_REPL_EMPTY": "Замен пока нет",
		"S_PASTE_DELAY": "Задержка перед вставкой", "S_PASTE_DELAY_SUB": "если программа не успевает принять текст",
		"S_OVPOS": "Где показывать плашку", "S_OVPOS_SUB": "у курсора — рядом с местом ввода; если приложение его не показывает, рядом с указателем мыши",
		"S_OVPOS_CARET": "У курсора",
		"S_OVTEXT": "Показывать распознанный текст", "S_OVTEXT_SUB": "в плашке после вставки, вместо числа символов",
		"S_TYPEMODE":       "Вставлять посимвольно, а не через буфер",
		"S_RECLANG":        "Язык исходной речи",
		"S_RECAUTO":        "Автоопределение",
		"S_DL":             "Скачать",
		"S_DEL":            "Удалить",
		"S_M_BASE":         "быстрая, для слабых ПК",
		"S_M_SMALL":        "баланс скорости и точности",
		"S_M_MED":          "точнее, рекомендуем",
		"S_M_TURBO":        "максимум точности на CPU",
		"S_M_GIGAAM":       "точнее на русском, сама ставит знаки препинания",
		"S_M_PARAKEET":     "25 европейских языков, сама ставит знаки препинания",
		"S_PUNCT":          "Пунктуация и заглавные",
		"S_PUNCT_MODEL":    "из модели",
		"S_PUNCT_LLM":      "моделью-редактором",
		"S_PUNCT_OFF":      "убирать",
		"S_SEARCH":         "Найти настройку…",
		"S_HOTMODE":        "Режим",
		"S_HOTMODE_HOLD":   "удержание",
		"S_HOTMODE_TOGGLE": "фиксация",
		"S_SUB_HOTMODE":    "удерживать клавиши или включать и выключать нажатием",
		"S_GRP_GENERAL":    "Общее",
		"S_GRP_SPEECH":     "Обработка речи",
		"S_GRP_INFO":       "Сведения",
		"S_NAV_POST":       "Постобработка",
		"S_NAV_HELP":       "Справка",
		"S_HELP_TOC":      "На этой странице",
		"S_HELP_TOC_SHOW": "Показать оглавление — окно станет шире",
		"S_HELP_TOC_HIDE": "Скрыть оглавление и вернуть ширину окна",
		"S_NAV_CONTACTS":   "Контакты",
		"S_STATE_ACTIVE": "Распознаёт",
		"S_STATE_USED": "Задействованные модели",
		"S_STATE_INST": "Установлены локально",
		"S_STATE_INST_SUB": "модели на диске, готовые к назначению",
		"S_PRESETS": "Какая модель какому языку",
		"S_PRESETS_HINT": "Щёлкните язык — под ним раскроется выбор моделей для него. Языки без своей модели используют модель Автоопределения.",
		"S_MFOLDER": "Своя модель",
		"S_DICT_SAVE": "Сохранить",
		"S_OWNM_SUB": "Добавьте локальную модель распознавания речи",
		"S_OWNM_ONEFILE": "Один файл",
		"S_OWNM_FOLDERF": "Папка с файлами модели",
		"S_OWNM_S1": "Откройте папку моделей",
		"S_OWNM_S1S": "Папка назначения:",
		"S_OWNM_S2": "Скопируйте модель",
		"S_OWNM_S2S": "Выберите одну из поддерживаемых структур",
		"S_OWNM_S3": "Перезапустите приложение",
		"S_OWNM_S3S": "Модель появится для языков, которые она поддерживает",
		"S_AS_AUTO": "как Автоопределение",
		"S_REC_CHIP": "рекомендованный",
		"S_BACK_AUTO": "Вернуть как Автоопределение",
		"S_LANGS_COUNT": "языков: %d",
		"S_LANGS_UNKNOWN": "языки: неизвестны",
		"S_TR_EN": "переводит на английский",
		"S_TR_LIST": "переводит: %s",
		"S_DL_GOING": "скачивается:",
		"S_OPEN_FOLDER": "Открыть папку",
		"S_UNLOAD": "Выгрузить из памяти",
		"S_UNLOAD_SUB": "память освободится; следующая диктовка загрузит модель заново",
		"S_UNLOAD_GO": "Выгрузить",
		"S_UNLOADED": "Выгружено",
		"S_NOT_FOR_LANG": "%s не распознаёт этот язык",
		"S_MANUAL_NOTE": "Скачать из приложения нельзя — лицензия запрещает распространение. Скачайте архив по ссылке и распакуйте в models/moonshine-uk.",
		"S_MANUAL_LINK": "Скачать самому",
		"S_HF_FIT": "только подходящие этому компьютеру",
		"S_HF_HIDDEN": "скрыто: %s",
		"S_WIZ_SKIP_DL": "Скачать позже",
		"S_WIZ_SKIP_NOTE": "Без модели диктовка не заработает. Скачать можно в разделе «Языки и модели».",
		"S_M_GIGAAM2": "прежнее поколение русской модели: та же скорость, но без знаков препинания",
		"S_M_MOONUK": "украинская модель Moonshine: быстрая и лёгкая, без знаков препинания",
		"S_M_LOCAL": "найдена в папке models; характеристики неизвестны, поэтому полосок нет",
		"S_ALL_LANGS": "все языки",
		"S_OVPOS_SCHEME_SUB": "щёлкните по экрану — плашка встанет туда",
		"S_OVDRAG": "Перетащите, куда нужно",
		"S_OVMON": "Экран",
		"S_OVMON_SUB": "на каком мониторе показывать плашку",
		"S_OVMON_CURSOR": "Экран с курсором",
		"S_M_NEMOTRON": "печатает по ходу речи: текст появляется на плашке, пока вы говорите; 40 языков, знаки сама",
		"S_M_TINY": "самая маленькая и быстрая, для очень слабых машин; точность заметно ниже",
		"S_STATE_LOADED": "Сейчас в памяти",
		"S_STATE_LOADED_SUB": "модели выгружаются сами после простоя",
		"S_STATE_WEEK": "За неделю",
		"S_ST_SUMMARY": "Сводка", "S_ST_OVERLAY": "Плашка на экране", "S_ST_BEEP": "Звуковой сигнал", "S_ST_AUTORUN": "Запуск с Windows", "S_ST_POST": "Постобработка", "S_ST_LOCAL": "локально", "S_ST_CHECKED": "проверено %s", "S_ST_GB": "%s ГБ", "S_ST_ON_M": "включён", "S_ST_OFF_M": "выключен", "S_ST_MIC_OK": "сигнал в норме", "S_ST_MIC_BAD": "микрофон не отвечает", "S_ST_CHECK": "Проверить", "S_ST_RECOG": "распознаёт %s", "S_ST_VER": "Версия %s", "S_ST_LATEST": "последняя", "S_ST_OUTDATED": "не последняя", "S_ST_UPD_OK": "у вас последняя версия", "S_ST_UPD_DL": "Загружаю обновление…",
		"S_ST_QUICK": "Быстрые настройки",
		"S_ST_MODELS": "Модели",
		"S_ST_USAGE": "Использование за неделю",
		"S_ST_READY": "Готово к диктовке",
		"S_ST_OFF": "Выключено в трее",
		"S_ST_OFF_SUB": "сочетание не работает, пока не включите обратно",
		"S_ST_ENABLE": "Включить",
		"S_ST_GOTO": "Открыть настройку на её вкладке",
		"S_ST_HOTKEY_GO": "Сменить сочетание",
		"S_ST_UPD_LAST": "Версия %s — последняя",
		"S_ST_UPD_HAVE": "Доступна версия %s",
		"S_ST_MEM": "Свободно %s ГБ из %s",
		"S_ST_MEM_SUB": "в памяти %s · на диске %d моделей, %s ГБ",
		"S_ST_MEM_NONE": "ничего",
		"S_ST_LANG": "Язык речи",
		"S_ST_ASR": "Распознавание",
		"S_ST_ON": "включено",
		"S_ST_OFF_W": "выключено",
		"S_ST_ON_F": "включена",
		"S_ST_OFF_F": "выключена",
		"S_ST_ACTIVE": "активна",
		"S_ST_IDLE": "не запускается",
		"S_ST_DISK": "лежит на диске, %s",
		"S_ST_USAGE_SUB": "%d знаков · сегодня %d · в среднем %d знаков",
		"S_WEEK_OTHER": "прочие",
		"S_ST_NO_WEEK": "диктовок за неделю не было",
		"S_ST_AUTORUN_SUB": "приложение не поднимется само",
		"S_ST_OVERLAY_SUB": "видна при записи",
		"S_REPL_LANG": "Язык правила",
		"S_REPL_LANG_ALL": "все языки",
		"S_M_CANARY": "английский, немецкий, испанский, французский — и переводит между ними сама",
		"S_M_QWEN3": "около 30 языков, сама ставит знаки; самая тяжёлая и точная в каталоге",
		"S_POSTAPI": "Внешний сервер",
		"S_POST_HINT": "Правит распознанный текст по промптам: убирает слова-паразиты, чинит пунктуацию, меняет стиль. Выключено — текст вставляется как распознан.",
		"S_POST_MODEL": "Модель",
		"S_SRC_LOCAL": "Локальная",
		"S_SRC_USED": "используется",
		"S_HF_GO": "Найти",
		"S_POSTAPI_HINT": "По умолчанию пусто — вся постобработка идёт локально. Впишите адрес, и промпты будет выполнять внешний сервер: OpenAI, Groq, свой vLLM — любой с совместимым API.",
		"S_POSTAPI_URL": "Адрес",
		"S_POSTAPI_URL_SUB": "пусто = локальная модель; пример: https://api.openai.com/v1",
		"S_POSTAPI_MODEL": "Модель",
		"S_POSTAPI_KEY": "Ключ API",
		"S_POSTAPI_KEY_SET": "ключ сохранён (зашифрован Windows DPAPI)",
		"S_POSTAPI_KEY_NONE": "ключа нет",
		"S_POSTAPI_SAVE": "Сохранить ключ",
		"S_POSTAPI_TIMEOUT": "Ожидание ответа", "S_SEC_SHORT": "сек",
		"S_POSTAPI_WARN": "⚠ Распознанный текст диктовок будет уходить на этот адрес. Звук не уходит никогда. Ключ хранится зашифрованным.",
		"S_POSTAPI_ASK": "Отправлять распознанный текст на %s? Звук останется на компьютере, но текст будет уходить наружу.",
		"S_POSTAPI_BADGE": "внешний сервер",
		"S_HIST_ADD":       "Добавить",
		"S_CONTACT_MAIL":   "Почта",
		"S_CONTACT_TITLE":   "Связаться",
		"S_ABOUT_DEPS":      "Внешние зависимости",
		"S_ABOUT_DEPS_HINT": "Библиотеки и программы, из которых собрано приложение, и их лицензии. Адрес под каждой открывает страницу проекта.",
		"S_DEP_WHISPER":     "Запускает модели Whisper. Работает отдельной программой whisper-server.exe рядом с приложением: приложение отдаёт ей записанный звук и получает текст. Нужна всегда, когда язык диктовки ведёт Whisper-модель, а также для перевода речи на английский.",
		"S_DEP_LLAMA":       "Запускает модель-редактор в постобработке — ту, что переписывает распознанный текст по вашим промптам. Отдельная программа llama-server.exe; поднимается только когда постобработка включена и выбрана локальная модель.",
		"S_DEP_SHERPA":      "Второй движок распознавания: на нём работают GigaAM, Parakeet, Canary, Qwen и Moonshine. Именно он разбирает русскую речь по умолчанию. Живёт в файлах sherpa-server.exe и sherpa-online-server.exe.",
		"S_DEP_GGML":        "Библиотека вычислений, на которой построены whisper.cpp и llama.cpp. Отдельно не запускается — считает матрицы внутри обоих движков, в том числе на инструкциях AVX вашего процессора.",
		"S_DEP_ONNX":        "Исполняет модели формата ONNX внутри sherpa-onnx. Без неё не заработают GigaAM, Parakeet и остальные модели этого семейства.",
		"S_DEP_WEBVIEW":     "Мост между Go и WebView2: через него открывается и живёт это самое окно настроек, а кнопки страницы вызывают код приложения.",
		"S_DEP_WV2RT":       "Компонент Windows, который и рисует окно настроек. Он же обычно уже стоит в системе вместе с Edge; приложение только просит его показать свою страницу.",
		"S_DEP_MALGO":       "Захват звука с микрофона из Go: открывает выбранное устройство, отдаёт поток во время записи и показывает уровень голоса на плашке.",
		"S_DEP_MINIAUDIO":   "Звуковая библиотека внутри malgo. Ею же проигрываются сигналы начала и конца записи.",
		"S_DEP_WS":          "Связь по websocket с sherpa-server: приложение отправляет туда звук и получает распознанный текст.",
		"S_DEP_XSYS":        "Доступ к системным вызовам Windows из Go: глобальные горячие клавиши, буфер обмена, окна, реестр, значок в трее.",
		"S_DEP_WINLOADER":   "Загружает библиотеки WebView2 в память процесса. Служебная часть go-webview2, отдельно не используется.",
		"S_DEP_PLEX":        "Шрифт интерфейса. Вшит в приложение, поэтому окно выглядит одинаково на любой машине, даже если такого шрифта в системе нет.",
		"S_DEP_HF":          "Каталог, откуда качаются модели распознавания и постобработки. Обращение к нему происходит только когда вы сами нажимаете «Скачать».",
		"S_CONTACT_HINT":    "Ошибка, идея, вопрос по настройке — пишите на почту, если разговор личный, или заводите issue, если это ошибка.",
		"S_CONTACT_REPO":    "Репозиторий",
		"S_CONTACT_ISSUES":  "Ошибки и идеи",
		"S_CONTACT_WRITE":   "Написать письмо",
		"S_CONTACT_OPEN":    "Открыть",
		"S_DICT_MODEL":     "Модель распознавания",
		"S_LIB_ACC":        "точность",
		"S_LIB_SPD":        "скорость",
		"S_NAV_STATE":      "Состояние",
		"S_NAV_DICT":       "Управление и поведение",
		"S_NAV_MIC":        "Микрофон",
		"S_NAV_MODELS":     "Языки и модели",
		"S_NAV_TEXT":       "Правила",
		"S_NAV_TR":         "Перевод",
		"S_NAV_SYSTEM":     "Система",
		"S_NAV_ABOUT":      "О программе",
		"S_STATE_HINT":     "зажмите и говорите — текст появится там, где курсор",
		"S_STATE_PROC":     "Обработка",
		"S_STATE_MEM":      "Память",
		"S_STATE_MEM_SUB":  "модели держатся загруженными, первая фраза без задержки",
		"S_SUB_MINMS":      "отсекает случайные нажатия",
		"S_SUB_ENTER":      "отправляет сообщение сразу",
		"S_SUB_CLIP":       "картинки и файлы возвращаются как были",
		"S_SUB_TYPE":       "помогает там, где поле запрещает вставку из буфера",
		"S_SUB_THREADS":    "больше потоков — не всегда быстрее, проверьте на своей машине",
		"S_SUB_PUNCT":      "откуда берутся знаки препинания и заглавные",
		"S_SUB_TRTARGET":   "в него переводится текст; в диалоге на плашке он предложен первым",
		"S_SUB_AUTOSTART":  "выключите, если сервер поднимаете сами",
		"S_SUB_PORT":       "распознаватель перезапустится сам",
		"S_SUB_UPD":        "единственный сетевой запрос, кроме загрузки моделей",
		"S_STATE_LAST":     "Последняя диктовка",
		"S_STATE_COPY":     "Копировать",
		"S_SEC_OVERLAY":    "Плашка на экране",
		"S_SEC_SERVICE":    "Служебное",
		"S_SEC_LLM":        "Модель-редактор",
		"S_NOT_INSTALLED":  "не установлена",
		"S_CHANGE_MODEL":   "Сменить",
		"S_RETRY":          "Повторить", "S_BERR_OPEN": "Открыть настройки сервера",
		"S_PICK_MODEL":  "Подобрать",
		"S_MIC":         "Микрофон",
		"S_MIC_DEFAULT": "Системный по умолчанию",
		"S_MIC_CHECK":   "Проверить микрофон", "S_MIC_CHECK_SUB": "три секунды записи и разбор: громкость, перегруз, есть ли речь", "S_MIC_CHECKING": "Проверяю…",
		"S_MCHECK": "Проверить установленные модели", "S_MCHECK_SUB": "сверяет файлы моделей с эталонными хешами", "S_MCHECK_GO": "Проверить", "S_MCHECK_RUN": "Проверяю…",
		"S_HIST_INSERT": "Вставить",
		"S_MIC_REFRESH":     "Обновить список",
		"S_MIC_LEVEL":       "Уровень сигнала",
		"S_MIC_QUIET":       "тихо",
		"S_THREADS":         "Потоки CPU",
		"S_MINMS":           "Минимальная запись, мс",
		"S_MAXSEC":          "Максимальная запись, сек",
		"S_AUTOSTART":       "Запускать whisper-server автоматически",
		"S_PORT":            "Порт",
		"S_SERVEREXE":       "Путь к whisper-server",
		"S_THEME_PINK":      "Розовый",
		"S_THEME_BLUE":      "Синий",
		"S_THEME_AMBER":     "Янтарный",
		"S_THEME_GREEN":     "Зелёный",
		"S_THEME_SUB":       "цвет окна, плашки и значка в трее",
		"S_SKIN_TERMINAL": "Терминал",
		"S_SKIN_SOFT": "Мягкий",
		"S_SKIN_PAPER": "Документ",
		"S_SKIN_SUB": "шрифт, форма, эффекты и анимация",
		"S_SKIN": "Дизайн",
		"S_WND_CLOSE":       "Закрыть окно",
		"S_WND_MIN":         "Свернуть в трей",
		"S_WND_RESTORE":     "Вернуть прежний размер",
		"S_WND_MAX":         "Развернуть на весь экран",
		"S_THEME_NEON":      "Неон",
		"S_THEME_EDITOR":    "Редактор",
		"S_THEME":           "Цвет",
		"S_UPD_FOUND":       "Есть версия %s",
		"S_RELOAD_CFG_BTN":  "Перечитать",
		"S_RELOAD_CFG_SUB":  "если правили файл руками",
		"S_RELOAD_CFG":      "Перечитать config.json",
		"S_RESET_ALL_ASK":   "Вернуть все настройки к заводским? Модели, история и промпты останутся на месте.",
		"S_RESET_ALL_BTN":   "Сбросить",
		"S_RESET_ALL_SUB":   "вернуть всё, кроме моделей и истории, к заводскому виду",
		"S_RESET_ALL":       "Сбросить настройки",
		"S_EXE_WARN":        "Приложение само находит whisper-server рядом с собой. Если вписать путь вручную, после переноса папки распознавание перестанет запускаться. Изменить?",
		"S_EXE_RESET":       "Сбросить",
		"S_SERVEREXE_SUB":   "заполняется сам; менять стоит, только если сервер лежит в другом месте",
		"S_SERVERURL":       "Удалённый сервер распознавания (URL)",
		"S_URLHINT":         "Если задан — свой сервер не запускается",
		"S_REMOTE_WARN":     "Звук будет уходить на этот сервер. Локальный режим выключен.",
		"S_REMOTE_ASK":      "Аудио перестанет обрабатываться на этом компьютере и будет отправляться на %s. Включить удалённый режим?",
		"S_STT_SRV":         "Сервер распознавания",
		"S_STT_SRV_HINT":    "Whisper-модели крутит отдельная программа. Она может работать на этом компьютере или на другом — выберите, какую использовать.",
		"S_SRV_LOCAL":       "На этом компьютере",
		"S_SRV_REMOTE":      "На другом компьютере",
		"S_SRV_REMOTE_HINT": "Тот же whisper-server, поднятый где-то ещё: домашний сервер, машина с видеокартой, соседний компьютер в сети.",
		"S_SRV_K_AUTO":      "автозапуск",
		"S_SRV_K_FILE":      "файл",
		"S_SRV_K_ADDR":      "адрес",
		"S_SRV_K_CHECK":     "проверка",
		"S_SRV_NEAR":        "whisper-server.exe рядом с приложением",
		"S_SRV_NOADDR":      "не задан",
		"S_SRV_NOCHECK":     "не проверяли",
		"S_SRV_LOCAL_DLG":   "Локальный сервер распознавания",
		"S_SRV_ADDR":        "Адрес сервера",
		"S_SRV_ADDR_SUB":    "адрес whisper-server на другой машине, вместе с портом",
		"S_SRV_ON":          "включён",
		"S_SRV_OFF":         "выключен",
		"S_SRV_K_THREADS":   "потоки CPU",
		"S_SRV_K_PORT":      "порт",
		"S_SRV_DOWN":       "Распознавание недоступно",
		"S_SRV_DOWN_WHY":   "удалённый сервер распознавания не настроен — задайте адрес в настройках",
		"S_SRV_DOWN_GO":    "Открыть настройки сервера",
		"S_SRV_WARN_NOW":    "Диктовка сейчас не работает: выбран удалённый сервер, а его адрес не задан.",
		"S_SRV_WARN_LATER":  "Как только будет выбрана модель Whisper, распознавание работать не будет: адрес удалённого сервера не задан.",
		"S_REMOTE_BADGE":    "УДАЛЁННО",
		"S_REMOTE_ABOUT":    "Сейчас задан удалённый сервер: звук уходит на него, и обещание выше не действует.",
		"S_SAVED":           "Сохранено",
		"S_STATE_GET":       "Скачать",
		"S_OK":              "Да",
		"S_CANCEL":          "Отмена",
		"S_DL_ASK":          "Модель «%s» не загружена (%s). Начать загрузку?",
		"S_DL_START":        "Загрузить",
		"S_DL_CANCEL":       "Отменить загрузку",
		"S_NOT_FOUND":       "ничего",
		"S_UPD":             "Обновления",
		"S_UPD_CHECK":       "Проверить обновления",
		"S_UPD_AUTO":        "Проверять при запуске",
		"S_UPD_NONE":        "Установлена последняя версия",
		"S_BADGE_MODELS":    "Установленные модели",
		"S_BADGE_MISS":      "Модель не скачана",
		"S_BADGE_SYSTEM":    "Предупреждения — нужно внимание",
		"S_BADGE_HIST":      "Записей в истории",
		"S_LOG_OPEN":        "Открыть лог",
		"S_LOG":             "Журнал работы",
		"S_LOG_SUB":         "весь текст, который приложение пишет о себе",
		"S_UPD_AVAIL":       "Доступна версия %s.",
		"S_UPD_GO":          "Обновить",
		"S_UPD_ERR":         "Не удалось проверить обновления",
		"S_UPD_DL":          "Скачиваю обновление…",
		"S_AUTHOR_HTML": "<p style=\"font-size:15px;letter-spacing:2px\"><b>Vitalii Yemets</b></p>" +
			"<p>Автор и разработчик {app} — локального диктовщика для Windows: голос превращается в текст прямо в позиции курсора, без облаков и подписок.</p>" +
			"<p>Проект открыт: исходный код, сборка и свежие версии — на GitHub.</p>" +
			"<ul>" +
			"<li>Репозиторий: <span class=\"lnk\" onclick=\"appRepoLink()\">github.com/Vitalii-Yemets/holdtotype</span></li>" +
			"<li>Профиль автора: <span class=\"lnk\" onclick=\"appAuthorLink()\">github.com/Vitalii-Yemets</span></li>" +
			"</ul>" +
			"<p>Нашли ошибку или есть идея — создайте issue в репозитории.</p>",
		"S_ABOUT_HTML": "<p><b>Голос → текст в позицию курсора.</b></p>" +
			"<p>Поставьте курсор в поле ввода, зажмите сочетание клавиш, скажите фразу, отпустите — текст вставится сам.</p>" +
			"<p>Всё работает полностью локально и офлайн: звук и текст не покидают компьютер.</p>" +
			"<p class=\"wh\">Технический стек</p>" +
			"<ul>" +
			"<li><b>Go + WinAPI</b> — клиент: трей, глобальные хоткеи, оверлей, вставка текста;</li>" +
			"<li><b>WebView2</b> — окно настроек;</li>" +
			"<li><b>whisper.cpp</b> — распознавание речи (модели Whisper, формат GGML);</li>" +
			"<li><b>llama.cpp</b> — постобработка и перевод (LLM-модели, формат GGUF);</li>" +
			"<li><b>miniaudio</b> — захват звука с микрофона;</li>" +
			"<li><b>Hugging Face</b> — каталог, откуда скачиваются модели.</li>" +
			"</ul>" +
			"<p>Логи ограничены ~2 МБ на диске.</p>",
		"S_HELP_HTML": "<p class=\"wh\" id=\"s-what\">Что это такое</p>" +
			"<p>Вы говорите — слова появляются там, где мигает курсор. Как будто вы их напечатали, только быстрее.</p>" +
			"<p class=\"hdim\">Выглядит это так. Вы щёлкаете в поле, куда обычно печатаете: в письмо, в чат, в строку поиска, в документ. Зажимаете две клавиши и говорите фразу обычным голосом. Отпускаете. Через секунду фраза стоит в поле. Никуда переключаться не надо, ничего копировать не надо — приложение работает поверх любых программ и ждёт своих клавиш в фоне.</p>" +
			"<p class=\"hdim\">Считает всё ваш компьютер. Голос не сохраняется: он живёт несколько секунд в памяти, превращается в текст и исчезает. В интернет приложение выходит два раза за всю жизнь: когда вы скачиваете модель распознавания и когда оно проверяет, не вышла ли новая версия. Проверку можно выключить, и тогда сеть ему не нужна совсем.</p>" +
			"<p class=\"hdim\">Слово «модель» дальше встретится часто. Это просто файл, который умеет разбирать речь. Он скачивается один раз и лежит в папке рядом с приложением; чем он больше, тем точнее слышит и тем медленнее работает.</p>" +
			"<p class=\"wh\" id=\"s-flow\">Что происходит, когда вы отпустили клавиши</p>" +
			"<div class=\"hflow\">" +
			"<span>запись</span><em>→</em><span>распознавание</span><em>→</em><span>замены</span><em>→</em><span>команды</span><em>→</em><span>перевод</span><em>→</em><span>правка</span><em>→</em><span>вставка</span>" +
			"</div>" +
			"<p>Обычно это одна-две секунды, и разбираться в шагах незачем. Но если текст однажды выйдет не таким, как вы ожидали, по этой цепочке легко понять, кто его изменил.</p>" +
			"<ul>" +
			"<li><b>Запись.</b> Пока клавиши зажаты, приложение слушает микрофон. Слишком короткое нажатие оно пропускает мимо ушей, чтобы случайно задетая клавиша не начинала диктовку.</li>" +
			"<li><b>Распознавание.</b> Голос превращается в буквы. Этим занимается модель. Для русского обычно работает GigaAM, для остальных языков — Whisper; это просто названия, выбирать их вручную не обязательно.</li>" +
			"<li><b>Замены.</b> Ваш личный список «услышано → надо». Если модель упорно пишет «гит хаб» вместо «GitHub», вы поправляете это один раз, и дальше оно исправляется само.</li>" +
			"<li><b>Команды.</b> Несколько фраз работают как кнопки. Скажете «новая строка» — текст перейдёт на новую строку, а сами эти слова в поле не попадут.</li>" +
			"<li><b>Перевод.</b> Если вы его включили, вы говорите по-русски, а в поле встаёт английский текст.</li>" +
			"<li><b>Правка.</b> Необязательный шаг. Вторая модель причёсывает текст: убирает «эээ» и «ну как бы», расставляет запятые, делает формулировку деловой.</li>" +
			"<li><b>Вставка.</b> Текст попадает в поле. То, что лежало у вас в буфере обмена до диктовки, возвращается на место.</li>" +
			"</ul>" +
			"<p class=\"hdim\">Пока всё это идёт, внизу экрана видна узкая плашка с подписью текущего шага. По ней сразу понятно: приложение занято, а не зависло.</p>" +
			"<p class=\"wh\" id=\"s-first\">Первый запуск</p>" +
			"<p>При первом старте приложение задаёт пять вопросов. Это займёт минуту, и после них можно диктовать.</p>" +
			"<ul>" +
			"<li><b>На каком языке говорить с вами.</b> Язык окна настроек. Обычно уже угадан по Windows.</li>" +
			"<li><b>На каком языке вы будете диктовать.</b> Самый важный вопрос: под ответ приложение само найдёт и скачает подходящую модель. Файл весит от 150 МБ до гигабайта, качается в фоне, прогресс виден.</li>" +
			"<li><b>Какими клавишами запускать.</b> Просто нажмите то, что вам удобно держать одной рукой. По умолчанию это Ctrl и клавиша Windows.</li>" +
			"<li><b>Какой микрофон слушать.</b> Рядом живая полоска. Скажите что-нибудь: если она не шевелится, выбран не тот микрофон.</li>" +
			"<li><b>Проверка.</b> Поле, где можно надиктовать пробную фразу и убедиться, что всё работает.</li>" +
			"</ul>" +
			"<p class=\"hdim\">Любой шаг можно пропустить и вернуться к нему позже: всё это есть в настройках. Если захотите пройти мастер заново, запустите приложение с добавкой <code>-wizard</code>. При обновлении он больше не появляется.</p>" +
			"<p class=\"wh\" id=\"s-overlay\">Плашка внизу экрана</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hpill\"><span class=\"hdot\"></span><span>Говорите…</span><span class=\"hbars\"><i style=\"height:6px\"></i><i style=\"height:13px\"></i><i style=\"height:9px\"></i><i style=\"height:16px\"></i><i style=\"height:7px\"></i><i style=\"height:12px\"></i></span><span>✕</span></div>" +
			"<div class=\"hcap\">так приложение показывает, что слышит вас</div>" +
			"</div>" +
			"<p>Во время диктовки это единственное, что видно на экране. Плашка маленькая, мышь не перехватывает и работать не мешает — она просто рассказывает, что сейчас происходит с вашей фразой.</p>" +
			"<ul>" +
			"<li><b>«Говорите…»</b> Идёт запись. Столбики рядом дёргаются в такт голосу. Если они лежат неподвижно, вас не слышно, и говорить дальше бессмысленно.</li>" +
			"<li><b>«Распознаю…»</b> Вы отпустили клавиши, работает модель. Обычно секунда-две. Самая первая диктовка после запуска дольше: модель в этот момент загружается в память.</li>" +
			"<li><b>«Перевожу…»</b> Идёт перевод.</li>" +
			"<li><b>«Редактирую: Чистка (1 из 2)»</b> Текст правит вторая модель. Видно, какое именно указание выполняется и сколько их всего.</li>" +
			"<li><b>«Вставлено: 128 симв.»</b> Готово, текст уже в поле. Плашка погаснет сама через пару секунд.</li>" +
			"<li><b>Крестик.</b> Передумали — нажмите его на любом шаге. Запись выбрасывается, в поле ничего не попадёт.</li>" +
			"</ul>" +
			"<p class=\"hdim\">Плашку можно выключить, поставить в любой угол, привязать к курсору или просто перетащить мышкой туда, где она не мешает. Это на вкладке «Управление и поведение».</p>" +
			"<p class=\"wh\" id=\"s-state\">Экран «Состояние»</p>" +
			"<p>Этот экран отвечает на один вопрос: можно ли прямо сейчас диктовать. Если нет — говорит почему и даёт кнопку, которая ведёт туда, где это чинится.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hcard\"><span class=\"hled\"></span><div class=\"ht\"><b>Готово к диктовке</b><span>зажмите и говорите — текст появится там, где курсор</span></div><span class=\"hkey\">ctrl</span><span style=\"color:var(--dim)\">+</span><span class=\"hkey\">win</span></div>" +
			"<div class=\"hcap\">зелёный квадрат означает, что приложение слушает клавиши</div>" +
			"</div>" +
			"<p>Верхняя карточка бывает трёх видов. <b>Зелёная</b> — всё в порядке, рядом напоминание про клавиши. <b>Жёлтая</b> — вы выключили приложение через значок у часов, клавиши не работают, пока не нажмёте «Включить» тут же на карточке. <b>Красная</b> — настройки противоречат друг другу: выбрано распознавание на другом компьютере, а адрес не указан.</p>" +
			"<p>Ниже — то, на что обычно смотрят мельком:</p>" +
			"<ul>" +
			"<li><b>Использование за неделю.</b> Круговая диаграмма: в какие программы уходили ваши диктовки и сколько знаков вы наговорили.</li>" +
			"<li><b>Версия.</b> Номер сборки и кнопка проверки. Если вышла новее, приложение скажет об этом здесь же.</li>" +
			"<li><b>Память.</b> Сколько её свободно, какая модель сейчас в ней сидит и сколько все модели весят на диске. Кнопка «Выгрузить» освобождает память немедленно — например, перед тяжёлой игрой.</li>" +
			"<li><b>Быстрые настройки.</b> Пять карточек с самым ходовым: микрофон, язык, плашка, звук, автозапуск. Показывают, как настроено сейчас, а стрелка ведёт на вкладку, где это меняется. Переключателей здесь намеренно нет, чтобы ничего не задеть мимоходом.</li>" +
			"<li><b>Модели.</b> Что установлено, что работает прямо сейчас, за какие языки отвечает и что просто лежит на диске.</li>" +
			"</ul>" +
			"<p class=\"wh\" id=\"s-system\">Вкладка «Система»</p>" +
			"<p>Внешний вид приложения, его отношения с Windows и несколько служебных кнопок. В обычной жизни сюда заходят раз в год.</p>" +
			"<p class=\"hsub\">Как приложение выглядит</p>" +
			"<p>Три строки: язык интерфейса, дизайн и цвет. На работу они не влияют совсем.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Язык интерфейса</td><td>Язык окна настроек и подписей на плашке. К языку диктовки отношения не имеет: окно может быть русским, а говорить вы будете по-английски.</td></tr>" +
			"<tr><td>Дизайн</td><td>Общий облик. «Терминал» — зелёные буквы на чёрном, как здесь. «Документ» — светлый и деловой. «Мягкий» — светлый и округлый. Есть ещё «Редактор» и «Неон».</td></tr>" +
			"<tr><td>Цвет</td><td>Оттенок выбранного дизайна: зелёный, янтарный, синий, розовый. Тем же цветом будет плашка и значок у часов. У светлых дизайнов цвет свой, и строка прячется.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Запуск и обновления</p>" +
			"<p>Здесь решается, стартует ли приложение вместе с системой и следит ли оно за новыми версиями.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Запускать вместе с Windows</td><td>Приложение поднимется при входе в систему и будет ждать ваших клавиш. Это обычный ярлык в автозагрузке — такой же, какой вы могли бы положить туда сами.</td></tr>" +
			"<tr><td>Проверять при запуске</td><td>Один раз при старте приложение спрашивает, не вышла ли версия новее. Больше в интернет оно не ходит, если не считать скачивания моделей. Выключите — и о новых версиях будете узнавать сами.</td></tr>" +
			"<tr><td>Обновления</td><td>Кнопка проверяет прямо сейчас. Если версия новее, приложение скачает и поставит её поверх: настройки, модели и история останутся на месте.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Служебное</p>" +
			"<p>Три кнопки на случай, когда что-то пошло не так.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Журнал работы</td><td>Открывает файл, куда приложение записывает, что оно делало: какая модель запустилась, сколько заняло распознавание, где случилась ошибка. Причина обычно написана словами, а не кодами. Файл не разрастается — старое затирается новым.</td></tr>" +
			"<tr><td>Перечитать config.json</td><td>Все настройки лежат в одном текстовом файле рядом с приложением. Если вы правили его руками или принесли с другого компьютера, эта кнопка подхватит изменения без перезапуска.</td></tr>" +
			"<tr><td>Сбросить настройки</td><td>Возвращает всё как после установки. Скачанные модели и история остаются на месте — стираются только настройки.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Где работает распознавание</p>" +
			"<p>Речь разбирает отдельная маленькая программа. Обычно она лежит рядом с приложением, запускается сама, и вам про неё думать не нужно. Но её можно поставить и на другой компьютер: например, на домашний, где мощная видеокарта, а диктовать с лёгкого ноутбука. Выбор — переключателем на карточках.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hcard\" style=\"border-color:var(--hi);box-shadow:var(--higlow)\"><span class=\"hsw on\"><i></i></span><div class=\"ht\"><b>На этом компьютере</b><span>потоки 4 · автозапуск включён · порт 8910 · программа рядом с приложением</span></div><button class=\"hbtn\">Редактировать</button></div>" +
			"<div class=\"hcap\">включённая карточка подсвечена, вторая в это время притушена</div>" +
			"</div>" +
			"<p>Кнопка «Редактировать» открывает окно с четырьмя настройками. Все они со значениями по умолчанию, и почти всегда их можно не трогать.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Потоки</td><td>Сколько ядер процессора отдать распознаванию. Больше ядер — быстрее, но только до какого-то предела: дальше выигрыша нет, а остальные программы начинают подтормаживать.</td></tr>" +
			"<tr><td>Запускать автоматически</td><td>Приложение поднимает распознаватель само. Выключать это стоит в единственном случае: если вы запускаете его вручную и не хотите, чтобы приложение вмешивалось.</td></tr>" +
			"<tr><td>Порт</td><td>Номер, по которому приложение разговаривает с распознавателем внутри вашего же компьютера. Наружу он не виден. Менять нужно, только если этот номер заняла другая программа.</td></tr>" +
			"<tr><td>Путь к программе</td><td>Где лежит распознаватель. Заполняется сам. Вписывать своё стоит, только если вы сознательно положили его в другое место: иначе после переноса папки диктовка перестанет запускаться.</td></tr>" +
			"</table>" +
			"<p>У второй карточки настройка одна — адрес компьютера в сети, например <code>http://192.168.1.50:8080</code>. Учтите: в этом варианте записанный голос уходит на ту машину, и по обычному адресу он идёт незашифрованным. Для домашней сети это нормально, через интернет так лучше не работать.</p>" +
			"<div class=\"hbad\">Если выбрано «на другом компьютере», а адрес не вписан, и при этом ваш язык ведёт Whisper — диктовка работать не будет. Приложение скажет об этом красной карточкой на «Состоянии», красной меткой у «Системы» и строкой внизу окна.</div>" +
			"<p class=\"wh\" id=\"s-mic\">Вкладка «Микрофон»</p>" +
			"<p>Больше половины неудачных диктовок — это не модель и не настройки, а микрофон: слушается не тот или вас почти не слышно. На этой вкладке всё, чтобы проверить это за полминуты.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Микрофон</td><td>Какое устройство слушать. Вариант «системный по умолчанию» означает «то же, что выбрано в Windows» — удобно, если вы то в наушниках, то без них: приложение поедет за системой само.</td></tr>" +
			"<tr><td>Обновить список</td><td>Нужна, если гарнитуру подключили уже после запуска приложения: список перечитается.</td></tr>" +
			"<tr><td>Уровень сигнала</td><td>Живая полоска. Скажите что-нибудь — столбики должны заметно двигаться. Не двигаются: микрофон приглушён в Windows или выбран не тот.</td></tr>" +
			"<tr><td>Проверить микрофон</td><td>Три секунды записи и разбор простыми словами: достаточно ли громко, нет ли перегруза (когда звук «квакает» от слишком близкого микрофона) и слышна ли вообще речь.</td></tr>" +
			"<tr><td>Звуковые сигналы записи</td><td>Короткий звук в начале и в конце записи. Особенно полезно, если плашку вы выключили: тогда хотя бы слышно, что запись пошла.</td></tr>" +
			"<tr><td>Звук сигналов</td><td>Какой именно звук — голос Windows, колокольчик, маримба и ещё несколько. Кнопка рядом проигрывает выбранный, чтобы не гадать по названию.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-hist\">Вкладка «История»</p>" +
			"<p>История — это список того, что вы наговорили. Только текст: голос не сохраняется нигде и никогда. По умолчанию она выключена, и это осознанно — в диктовки часто попадает личное.</p>" +
			"<p>Включить её стоит, если вы хоть раз попадали в такую ситуацию: вставили текст не в то окно и он пропал; нужно повторить длинную фразу; хочется найти, что диктовали вчера, и отправить ещё раз.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hrow\"><label>Сколько хранить</label><span class=\"hnum\"><span class=\"hfld\">7</span><span class=\"hspin\"><i>▲</i><i>▼</i></span></span><span class=\"hfld\">дней</span></div>" +
			"<div class=\"hrow\"><label>Сколько записей хранить</label><span class=\"hnum\"><span class=\"hfld\">200</span><span class=\"hspin\"><i>▲</i><i>▼</i></span></span></div>" +
			"<div class=\"hcap\">срок можно задать в минутах, часах или днях — любым числом</div>" +
			"</div>" +
			"<table class=\"hopt\">" +
			"<tr><td>Хранить историю диктовок</td><td>Главный выключатель. Выключили — новые записи не появляются, а старые стираются.</td></tr>" +
			"<tr><td>Сколько хранить</td><td>Через сколько запись исчезнет сама. Можно поставить пять минут, если история нужна только на случай «вставилось не туда», а можно тридцать дней.</td></tr>" +
			"<tr><td>Сколько записей хранить</td><td>Второй предел, по количеству. Работает вместе с первым: что наступит раньше, то и сработает. При «7 дней и не больше 50» пятьдесят первая запись вытеснит самую старую, даже если неделя ещё не прошла.</td></tr>" +
			"<tr><td>Не записывать из программ</td><td>Список исключений. Всё, что вы диктуете в эти окна — банк, менеджер паролей, рабочий чат, — в историю не попадёт, а сама диктовка будет работать как обычно. Программу можно выбрать из списка открытых окон, не выискивая имя файла.</td></tr>" +
			"</table>" +
			"<p>С каждой записью можно сделать четыре вещи. <b>Скопировать</b> — текст ляжет в буфер обмена. <b>Вставить</b> — приложение попросит показать, куда именно: щёлкните по нужному полю в любой программе, и текст встанет туда. <b>Удалить</b> — уберёт одну запись. <b>Найти</b> — поиск идёт и по тексту, и по названию программы. Кнопка «Очистить» стирает всё сразу и поэтому спрашивает подтверждение.</p>" +
			"<p class=\"hdim\">Справа у каждой записи написано, когда она исчезнет: «до 7 сент.», а в последний день эта надпись становится жёлтой.</p>" +
			"<p class=\"wh\" id=\"s-dict\">Вкладка «Управление и поведение»</p>" +
			"<p>Как запускается диктовка и как текст попадает в поле. Сюда заходят, когда что-то работает не так, как хочется: клавиши неудобные, текст не встаёт, Enter нажимается сам.</p>" +
			"<p class=\"hsub\">Клавиши</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Сочетание клавиш</td><td>Нажмите «Изменить» и зажмите то, что вам удобно. Приложение перехватывает это сочетание во всей системе, поэтому не берите занятое — Ctrl+C, например. Реже всего заняты сочетания с клавишей Windows.</td></tr>" +
			"<tr><td>Режим</td><td>«Удержание» — пишет, пока клавиши зажаты: отпустили, и фраза пошла в работу. «Фиксация» — первое нажатие включает запись, второе выключает. Второй режим удобен для длинных монологов, когда держать клавиши устаёт рука.</td></tr>" +
			"<tr><td>Минимальная запись</td><td>Насколько короткое нажатие считать случайным. Задели клавиши на долю секунды — приложение промолчит, вместо того чтобы вставлять обрывок слова.</td></tr>" +
			"<tr><td>Максимальная запись</td><td>Страховка от забытой клавиши: запись оборвётся сама и не будет писать в пустоту полчаса.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Как вставляется текст</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Нажимать Enter после вставки</td><td>Текст не просто вставится, а сразу отправится: сообщение уйдёт, поиск выполнится. В мессенджере это экономит движение, в письме или документе — почти наверняка помешает.</td></tr>" +
			"<tr><td>Восстанавливать буфер обмена</td><td>Вставка идёт через буфер обмена, и приложение кладёт обратно то, что там лежало: картинку, файл, скопированный кусок текста. Оставьте включённым, если часто копируете что-то и не хотите это терять.</td></tr>" +
			"<tr><td>Вставлять посимвольно</td><td>Вместо вставки одним куском приложение наберёт текст, как будто вы печатаете. Это медленнее, зато выручает там, где вставка из буфера запрещена: некоторые банковские формы, терминалы, старые программы.</td></tr>" +
			"<tr><td>Задержка перед вставкой</td><td>Небольшая пауза перед тем, как отдать текст окну. Нужна редко: бывает, программа не успевает принять текст, и начало фразы теряется.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Плашка на экране</p>" +
			"<p>Здесь решается, показывать ли её вообще и где именно.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Показывать плашку</td><td>Выключите, если она мешает или вам хватает звуковых сигналов.</td></tr>" +
			"<tr><td>Где показывать</td><td>Схема экрана с девятью местами: щёлкните по тому углу, где плашка не будет мешать.</td></tr>" +
			"<tr><td>У курсора</td><td>Плашка появляется рядом с тем местом, куда вы диктуете. Если программа не сообщает, где курсор, — рядом с указателем мыши.</td></tr>" +
			"<tr><td>Своё место</td><td>Плашку можно просто перетащить мышью, и она запомнит точку.</td></tr>" +
			"<tr><td>Экран</td><td>На каком мониторе показывать, если их несколько.</td></tr>" +
			"<tr><td>Показывать распознанный текст</td><td>После вставки на плашке видно саму фразу, а не число символов. Удобно для проверки, но фразу увидит любой, кто смотрит на ваш экран.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-models\">Вкладка «Языки и модели»</p>" +
			"<p>Здесь вы говорите, на каком языке будете диктовать, а приложение решает, чем это распознавать. Моделей несколько: одни быстрые, другие точные, третьи знают редкие языки.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Язык исходной речи</td><td>Язык, на котором вы говорите. Если языков несколько, поставьте «Автоопределение»: модель разберётся сама, хотя чуть медленнее и чуть менее точно, чем когда язык назван прямо.</td></tr>" +
			"<tr><td>Какая модель какому языку</td><td>Щёлкните по языку — раскроется список моделей, которые его знают, с оценкой скорости, точности и весом файла. Так можно оставить для русского быструю модель, а для английского выбрать точную.</td></tr>" +
			"</table>" +
			"<div class=\"hshot\">" +
			"<div class=\"hrow\"><label>Русский<span class=\"hs\">GigaAM v3 · точность 5 из 5 · скорость 5 из 5 · 232 МБ</span></label><span class=\"hchip\"><i></i>активна</span></div>" +
			"<div class=\"hrow\"><label>English<span class=\"hs\">как Автоопределение · Whisper Medium</span></label><button class=\"hbtn\">Выбрать</button></div>" +
			"<div class=\"hcap\">видно, кто ведёт язык и чего это стоит по скорости и месту на диске</div>" +
			"</div>" +
			"<p class=\"hsub\">Перевод</p>" +
			"<p>Некоторые модели умеют не просто расшифровать речь, а сразу перевести её: вы говорите по-русски, в поле встаёт английский текст. Если выбранная модель так не умеет, переключатель будет недоступен, а рядом будет написано почему.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Изменить язык текстового вывода</td><td>Включает перевод.</td></tr>" +
			"<tr><td>Язык по умолчанию</td><td>На какой язык переводить, если вы не выбрали другой прямо во время диктовки.</td></tr>" +
			"<tr><td>Спрашивать язык</td><td>Можно не спрашивать и всегда переводить на язык по умолчанию. Можно спрашивать каждый раз: после диктовки на плашке появятся кнопки с языками. А можно спрашивать с таймаутом — те же кнопки, но если вы не успели нажать, возьмётся язык по умолчанию.</td></tr>" +
			"<tr><td>Языки в диалоге</td><td>Какие языки показывать этими кнопками. Оставьте два-три, которыми правда пользуетесь, иначе выбирать будет дольше, чем говорить.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Модели на диске</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Своя модель</td><td>Если у вас есть своя модель, положите её в папку моделей и перезапустите приложение — она появится в списках для тех языков, которые понимает. Кнопка рядом открывает нужную папку, чтобы не искать её в проводнике.</td></tr>" +
			"<tr><td>Выгрузить из памяти</td><td>Пока приложение работает, модель занимает оперативную память. Эта кнопка освобождает её немедленно — например, перед тяжёлой игрой. Следующая диктовка загрузит модель заново, и первая фраза будет дольше обычного.</td></tr>" +
			"<tr><td>Проверить установленные модели</td><td>Сверяет скачанные файлы с эталоном. Пригодится, если связь обрывалась во время загрузки и есть подозрение, что файл побился.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-rules\">Вкладка «Правила»</p>" +
			"<p>Три списка, которые правят текст автоматически, ещё до того, как он попадёт в окно. Все три необязательные: без них приложение прекрасно работает.</p>" +
			"<p class=\"hsub\">Пунктуация и заглавные</p>" +
			"<p>Откуда в тексте берутся точки, запятые и большие буквы. Вариантов три: <b>из модели</b> — как услышала, так и расставила, быстро и обычно достаточно; <b>моделью-редактором</b> — знаки наводит вторая модель, аккуратнее, но каждая диктовка становится на несколько секунд длиннее; <b>убирать</b> — текст сплошной строкой без знаков, это нужно редко, например для строки поиска.</p>" +
			"<p class=\"hsub\">Словарь</p>" +
			"<p>Список слов, которых модель обычно не знает: фамилии, названия компаний, профессиональные термины. Это не приказ «пиши так», а подсказка: услышав что-то похожее, модель охотнее выберет ваш вариант.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hchips\"><span class=\"hchip\">GitHub</span><span class=\"hchip\">Docker</span><span class=\"hchip\">Kubernetes</span><span class=\"hchip\">HoldToType</span></div>" +
			"<div class=\"hcap\">слова добавляются по одному и живут отдельными карточками</div>" +
			"</div>" +
			"<p class=\"hdim\">Словарь понимают модели Whisper. На русскую речь через GigaAM он не влияет — там подсказки устроены иначе. Стандартный набор слов меняется вместе с языком, пока вы не впишете своё.</p>" +
			"<p class=\"hsub\">Замены после распознавания</p>" +
			"<p>Тоже список слов, но работает иначе: это уже не подсказка, а гарантия. Что встретилось слева — станет тем, что справа. У каждой замены три переключателя: <b>только целые слова</b> (тогда «код» не тронет «кодировку» — обычно это и нужно), <b>учитывать регистр</b> (тогда «Го» и «го» считаются разными) и <b>язык</b> (замена сработает только когда вы диктуете на нём).</p>" +
			"<p class=\"hsub\">Голосовые команды</p>" +
			"<p>Фразы, которые не попадают в текст, а выполняют действие. Работают, только если фраза сказана отдельно и целиком: те же слова в середине предложения командой не считаются.</p>" +
			"<ul>" +
			"<li><b>«новая строка»</b> — перевести строку;</li>" +
			"<li><b>«новый абзац»</b> — оставить пустую строку между абзацами;</li>" +
			"<li><b>«отмена»</b> — выбросить фразу и ничего не вставлять;</li>" +
			"<li><b>своя команда</b> — вставить заготовку: подпись, адрес, дежурный ответ.</li>" +
			"</ul>" +
			"<p class=\"wh\" id=\"s-post\">Вкладка «Постобработка»</p>" +
			"<p>Постобработка — это вторая модель, которая переписывает распознанный текст по вашему указанию. Например: «убери слова-паразиты и повторы», «сделай из этого деловое письмо», «переведи в список». Указание пишется обычными словами и называется промптом.</p>" +
			"<p>По умолчанию она выключена, и вот почему: это единственный шаг, который заметно замедляет диктовку. Вместо секунды-двух вставка займёт несколько секунд, а на слабой машине и десяток.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Включить постобработку</td><td>Главный выключатель. Пока он выключен, остальное на вкладке не работает, и текст вставляется таким, как его распознали.</td></tr>" +
			"<tr><td>Модель · Локальная</td><td>Редактор, который работает на вашем компьютере. Никуда ничего не уходит, но нужна память: такие модели весят от гигабайта. Рядом каталог, где модель выбирается и качается в один щелчок.</td></tr>" +
			"<tr><td>Модель · Внешний сервер</td><td>Вместо своей модели текст отправляется в чужой сервис — OpenAI, Groq или ваш собственный сервер. Быстрее и умнее, но текст покидает компьютер. Ключ доступа хранится зашифрованным, а кнопка «Тест соединения» проверяет, что адрес и ключ рабочие.</td></tr>" +
			"<tr><td>Промпты</td><td>Ваши указания редактору. Отмеченные применяются по очереди сверху вниз: сначала, скажем, чистка, потом деловой стиль. Если не отмечено ничего, постобработка не делает ничего — приложение предупредит об этом жёлтой строкой.</td></tr>" +
			"</table>" +
			"<div class=\"hwarn\">Голос на внешний сервер не уходит никогда — только распознанный текст. Но и он уходит, так что включайте это осознанно.</div>" +
			"<p>У каждого промпта есть <b>название</b> (по нему он виден на плашке во время работы), <b>текст указания</b> обычными словами, необязательное <b>своё сочетание клавиш</b> — диктовка им прогонит фразу только через этот промпт, минуя остальные, — и <b>проверка</b>: поле, где можно вписать фразу и посмотреть, что получится, ничего не диктуя.</p>" +
			"<p class=\"wh\" id=\"s-about\">«О программе» и «Контакты»</p>" +
			"<p>На вкладке «О программе» — короткий рассказ, что это и как устроено, и список всех чужих библиотек, из которых собрано приложение, с их лицензиями. Названия кликабельны и ведут на страницы этих проектов. Там же ссылка на репозиторий, где лежат исходный код и все версии.</p>" +
			"<p>«Контакты» — три карточки: почта, репозиторий и страница для ошибок и идей. Щелчок по адресу почты откроет вашу почтовую программу с уже заполненным получателем, а иконка рядом просто скопирует адрес.</p>" +
			"<p class=\"wh\" id=\"s-files\">Что лежит в папке приложения</p>" +
			"<p>Приложение портативное: всё его хозяйство лежит в одной папке и никуда больше не расползается.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>config.json</td><td>Все ваши настройки одним текстовым файлом. Его можно скопировать на другой компьютер, чтобы не настраивать заново.</td></tr>" +
			"<tr><td>holdtotype.log</td><td>Журнал работы. Первое, куда стоит заглянуть, если приложение ведёт себя странно.</td></tr>" +
			"<tr><td>holdtotype-history.json</td><td>История диктовок, если вы её включили. Обычный текст: при желании файл можно просто удалить.</td></tr>" +
			"<tr><td>models</td><td>Папка со скачанными моделями — самая тяжёлая часть, несколько гигабайт. Её тоже можно перенести на другой компьютер, чтобы не качать заново.</td></tr>" +
			"<tr><td>whisper-server, sherpa-server, llama-server</td><td>Те самые маленькие программы, которые запускают модели. Без них модель — просто файл, читать его нечем.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-privacy\">Что уходит в интернет, а что нет</p>" +
			"<ul>" +
			"<li><b>Голос</b> не сохраняется и не отправляется — если вы сами не включили распознавание на другом компьютере.</li>" +
			"<li><b>Текст</b> не покидает машину — если вы сами не включили внешний сервер постобработки.</li>" +
			"<li><b>История</b> лежит файлом рядом с приложением и чистится по вашему сроку.</li>" +
			"<li><b>Сеть</b> нужна дважды: скачать модель и проверить обновление. Вторую проверку можно выключить.</li>" +
			"<li><b>Ключ</b> от внешнего сервиса хранится зашифрованным средствами Windows и никогда не попадает в журнал.</li>" +
			"</ul>" +
			"<p class=\"wh\" id=\"s-trouble\">Если что-то пошло не так</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Текст не появляется в поле</td><td>Программа не принимает вставку из буфера. Включите «Вставлять посимвольно» на вкладке «Управление и поведение» или добавьте небольшую задержку перед вставкой.</td></tr>" +
			"<tr><td>Диктовка не запускается</td><td>Сочетание клавиш занято другой программой. Задайте другое — лучше с клавишей Windows, их занимают реже.</td></tr>" +
			"<tr><td>Записалась тишина</td><td>Микрофон не тот или приглушён. Откройте вкладку «Микрофон», посмотрите на полоску уровня и нажмите «Проверить микрофон».</td></tr>" +
			"<tr><td>Первая фраза очень долгая</td><td>Так и задумано: модель загружается в память при первом обращении. Дальше распознавание занимает секунды.</td></tr>" +
			"<tr><td>Написано «Распознавание недоступно»</td><td>Выбрано распознавание на другом компьютере, а адрес не задан. Впишите адрес или вернитесь к варианту «на этом компьютере».</td></tr>" +
			"<tr><td>Постобработка ничего не меняет</td><td>Не отмечен ни один промпт, не выбрана модель или внешний сервер не отвечает. Причина написана на самой вкладке, рядом с выключателем.</td></tr>" +
			"<tr><td>Текст странный, слова заменены</td><td>Загляните в «Правила»: скорее всего сработала замена или голосовая команда.</td></tr>" +
			"<tr><td>Ничего из перечисленного</td><td>«Система» → «Открыть лог». Приложение пишет о себе всё, и причина обычно видна словами.</td></tr>" +
			"</table>" +
			"<p class=\"hdim\" style=\"margin-top:16px\">Не нашли ответа — напишите на почту со вкладки «Контакты» или заведите issue в репозитории. И то и другое читает автор.</p>",
		"S_WIZ_HELLO":      "Две минуты на настройку",
		"S_WIZ_HELLO_TEXT": "{app} превращает голос в текст прямо в позиции курсора: зажали сочетание клавиш, сказали фразу, отпустили — текст на месте. Всё считается на вашем компьютере, звук никуда не уходит.",
		"S_WIZ_LATER":      "Всё, что мы сейчас выберем, потом меняется в настройках.",
		"S_WIZ_T_MODEL":    "Язык и модель",
		"S_WIZ_MODEL_TEXT": "Скажите, на каком языке будете диктовать, — модель подберу сам. Русский распознаёт GigaAM, остальные языки — Whisper.",
		"S_WIZ_T_INPUT":    "Клавиши и микрофон",
		"S_WIZ_INPUT_TEXT": "Это сочетание вы будете держать во время речи. Скажите что-нибудь и проверьте, что полоска уровня шевелится.",
		"S_WIZ_T_TRY":      "Проба",
		"S_WIZ_TRY_PH":     "текст появится здесь",
		"S_WIZ_T_DONE":     "Готово",
		"S_WIZ_DONE_TEXT":  "{app} живёт в трее: левый клик по значку — настройки, правый — меню. Диктовать можно в любом окне, где есть курсор ввода.",
		"S_AUTORUN":        "Запускать вместе с Windows",
		"S_AUTORUN_SUB":    "Ярлык в автозагрузке текущего пользователя",
		"S_WIZ_SKIP":       "Пропустить",
		"S_WIZ_BACK":       "Назад",
		"S_WIZ_NEXT":       "Дальше",
		"S_WIZ_FINISH":     "Завершить",
		"S_WIZ_WAIT":       "Жду первую фразу…",
		"S_WIZ_HEARD":      "Услышал:",
		"S_WIZ_HAVE":       "Всё нужное уже скачано",
		"S_WIZ_TRY_TEXT":   "Поставьте курсор в поле ниже, зажмите %s, скажите фразу и отпустите.",
		"S_MODEL_READY":    "Модель загружена — выберите её, чтобы переключиться",
	},
	"en": {
		"S_TITLE":          "{app} — Settings",
		"S_TR_DEFAULT":     "Change the text output language",
		"S_TR_TARGET":      "Default text output language",
		"S_SRCLANG_SUB": "you speak it; it decides the recognition model",
		"S_TR_LANGS_SUB": "these languages become buttons on the plate at insertion",
		"S_TR_UNAVAIL": "unavailable — %s cannot translate",
		"S_TR_LOCK": "%s cannot be removed from the list — it is the default text output language. Pick another default language, and then %s can be excluded.",
		"S_TR_LOCK_OK": "Got it",
		"S_TR_ONE": "Several languages are checked, but without asking the translation will always go into one — %s (the default text output language). The others stay checked but get disabled.",
		"S_TR_NOMODEL": "%s cannot translate. If you continue, translation will be turned off and unavailable while this model works.",
		"S_TR_CONFIRM": "Confirm",
		"S_TR_ASK":         "Ask for the text output language",
		"S_TR_ASK_NEVER":   "Don’t ask — translate right away",
		"S_TR_ASK_ALWAYS":  "Ask every time",
		"S_TR_ASK_TIMEOUT": "Ask, with a timeout",
		"S_TR_SECONDS":     "Timeout, sec",
		"S_TR_LANGS":       "Languages in the dialog",
		"S_DICT_HINT":      "Comma-separated terms, names and abbreviations — a hint for the ear, not commands. It works for Whisper; Russian speech through GigaAM ignores it. The default set follows the recognition language until you write your own.",
		"S_LLM_HINT":       "Checked profiles apply one after another, top to bottom, on regular dictation. Nothing checked — text is inserted as is.",
		"S_UPDATED":        "Model last updated",
		"S_FIT_OK":         "fits",
		"S_FIT_WARN":       "tight",
		"S_FIT_BAD":        "not enough RAM",
		"S_RAM":            "Computer RAM:",
		"S_HF_PH":          "Model name — e.g. qwen2.5 instruct",
		"S_NO_LLM":         "No models installed yet — find and download one in the search field below.",
		"S_NO_LLM_PROF":    "Prompts become available once a model is installed — the Model block above on this tab.",
		"S_PROF_EDIT":      "Edit",
		"S_PROF_CLOSE":     "Collapse",
		"S_CONFIRM_DEL":    "Delete the \"%s\" model? It can be downloaded again.",
		"S_DEL_ACTIVE":     "Delete the active model \"%s\"? Recognition stops until you pick another one — you can download it right here.",
		"S_WIZ_NEED_MODEL": "Download a model first — without one there is nothing to recognise with",
		"S_FREE":           "free",
		"S_SUB_PROMPTS":    "Prompts",
		"S_SUB_DICT":       "Dictionary",
		"S_PROF_ADD":       "Add",
		"S_PROF_NAME":      "Name",
		"S_PROF_PROMPT":    "Prompt",
		"S_PROF_TEST":      "Test",
		"S_HOTKEY":         "Keyboard shortcut",
		"S_CHANGE":         "Change…",
		"S_UILANG":         "UI language",
		"S_AUTO":           "System default",
		"S_BEEP":           "Recording sound cues",
		"S_SEC_SOUND":      "Sound",
		"S_SOUND":          "Cue sound",
		"S_SND_SPEECH":     "Windows voice",
		"S_SND_CHIME":      "Chime",
		"S_SND_SOFT":       "Soft",
		"S_SND_MARIMBA":    "Marimba",
		"S_SND_BLIP":       "Blip",
		"S_SND_POP":        "Pop",
		"S_AUTOENTER":      "Press Enter after paste (auto-submit)",
		"S_RESTORE":        "Restore clipboard after paste",
		"S_OVERLAY":        "Show the plate", "S_OVERLAY_SUB": "during dictation the screen shows that recording is on",
		"S_NAV_HISTORY":    "History", "S_HIST_ON": "Keep a history of dictations", "S_HIST_ON_SUB": "text only, on this computer; audio is never kept",
		"S_HIST_DAYS": "How many days to keep", "S_HIST_MAX": "How many entries to keep",
		"S_HIST_SKIP": "Never record from these programs", "S_HIST_SKIP_SUB": "comma-separated: keepass.exe, 1password.exe", "S_SKIP_ADD_DLG": "Adding a program", "S_SKIP_EDIT_DLG": "Editing the program", "S_SKIP_NAME": "Program name", "S_SKIP_NAME_SUB": "The file name without a path: keepass.exe. A star at the end catches every version: 1password*", "S_SKIP_OPEN": "Programs open right now", "S_SKIP_REFRESH": "Refresh the list", "S_SKIP_PICKED": "%d of %d picked", "S_SKIP_NONE": "Nothing picked", "S_SKIP_EMPTY": "The list is empty — history is kept from every program", "S_SKIP_ADD_BTN": "Add a program", "S_SKIP_HINT": "What you dictate into these programs never reaches the history. The pasting itself works as usual.",
		"S_HIST_LIST": "Entries", "S_HIST_CLEAR": "Clear", "S_HIST_TILL": "till %s", "S_HIST_TILL1": "till tomorrow", "S_HIST_TILL_FULL": "Will be deleted on %s — kept for %s", "S_HIST_LIST_HINT": "What was dictated: copy it, paste it into any window, or delete it.", "S_HIST_COPY": "Copy",
		"S_HIST_KEEP":       "How long to keep",
		"S_UNIT_MIN":        "minutes",
		"S_UNIT_HOUR":       "hours",
		"S_UNIT_DAY":        "days",
		"S_HIST_FIND": "Search the history…", "S_HIST_EMPTY": "No history yet", "S_HIST_ASK": "Delete the whole dictation history?",
		"S_SEC_CMD": "Voice commands", "S_CMD_HINT": "What you say turns into a line break, a symbol or a cancel instead of landing in the text. Matched as whole words, applied top to bottom, after the replacements.",
		"S_CMD_ADD": "Add a command", "S_CMD_PRESET": "Add the usual ones", "S_CMD_PH": "the phrase you will say",
		"S_CMD_NEWLINE": "line break", "S_CMD_PARAGRAPH": "new paragraph", "S_CMD_TEXT": "insert text", "S_CMD_CANCEL": "cancel the dictation",
		"S_CMD_TEXT_PH": "what to insert", "S_CMD_EMPTY": "No commands yet", "S_CMD_DEL": "Delete the command",
		"S_CMD_P_NEWLINE": "new line", "S_CMD_P_PARAGRAPH": "new paragraph", "S_CMD_P_CANCEL": "cancel",
		"S_SEC_REPLACE": "Replacements after recognition", "S_REPLACE_HINT": "What was misheard becomes what you meant — right after recognition, before the prompts. Applied from top to bottom.",
		"S_REPL_WHOLE_FULL": "Whole words only", "S_REPL_CASE_FULL": "Match letter case", "S_CMD_ACTION": "Action",
		"S_FM_ADD": "Add",
		"S_TIP_REPL_LANG": "The rule fires only when you dictate in the chosen language. “all languages” — it always fires.",
		"S_TIP_REPL_CASE": "Capital and small letters matter: “git” and “Git” are different words. Off — case is ignored.",
		"S_TIP_REPL_WHOLE": "The replacement fires only when the text stands as a separate word. Off — it also matches inside other words: “cat” would match “catalog”.",
		"S_TIP_CMD_ACTION": "What happens when you say the phrase: a line break, a new paragraph, your own text, or cancelling the dictation.",
		"S_LIST_FILTER_PH": "find…",
		"S_REPL_DEL": "Delete the replacement",
		"S_LIST_NOTHING": "Nothing found: “%s”",
		"S_FM_T_REPL_ADD": "Adding a replacement", "S_FM_T_REPL_EDIT": "Editing the replacement",
		"S_FM_T_CMD_ADD": "Adding a command", "S_FM_T_CMD_EDIT": "Editing the command",
		"S_MT_DEL": "Deleting a model", "S_MT_DEL_PROMPT": "Deleting a prompt", "S_MT_DL": "Downloading a model",
		"S_MT_TR_OFF": "Turning translation off", "S_MT_TR_ONE": "Translating without asking", "S_MT_TR_LOCK": "Default output language",
		"S_MT_REMOTE": "Remote server", "S_MT_POST": "External server", "S_MT_HIST": "Clearing the history",
		"S_MT_RESET": "Resetting the settings", "S_MT_EXE": "Server path",
		"S_DICT_ADD": "Add a word", "S_FM_T_DICT_ADD": "Adding a word", "S_DICT_EMPTY": "No words yet",
		"S_DICT_ADD_PH": "a word, or several separated by commas",
		"S_DICT_NOMODEL": "The current model %s does not support the dictionary — only Whisper models read it.",
		"S_OV_FREE": "Your own spot", "S_OV_FREE_SUB": "the plate can be dragged anywhere",
		"S_OVPOS_DRAG_SUB": "drag the plate with the mouse — it lands anywhere",
		"S_OVMON_N": "Screen %d",
		"S_POST_ENABLE": "Turn on post-processing",
		"S_API_SUM_URL": "address", "S_API_SUM_MODEL": "model", "S_API_SUM_KEY": "key", "S_API_SUM_TIMEOUT": "timeout",
		"S_API_SUM_STATE": "state", "S_API_NO_MODEL": "not set",
		"S_API_NONE": "not set up — post-processing runs locally",
		"S_POSTAPI_SETUP": "Set up", "S_API_EDIT": "Change", "S_API_KEY_DEL": "Delete the key", "S_API_DLG": "External server",
		"S_LLM_CATALOG": "Model catalog", "S_LLM_BLOCK": "Installed models", "S_LLM_NONE_HINT": "No model is installed yet — download a found one with the arrow and it will appear here", "S_LLM_IN_MEM": "in memory", "S_LLM_ON_DISK": "on disk", "S_LLM_EJECT": "Unload from memory", "S_LLM_FOUND": "found %d", "S_LLM_NOSEARCH": "no search yet", "S_LLM_SEARCH_HINT": "Type a model name and press “Search”", "S_LLM_PICK_WAIT": "Available once the model is downloaded", "S_LLM_INSTALLED": "installed",
		"S_LLM_SUM_MODEL": "model", "S_LLM_SUM_SIZE": "size", "S_LLM_SUM_COUNT": "installed", "S_LLM_SUM_RAM": "memory",
		"S_DLG_CLOSE": "Close", "S_LLM_NOPICK": "not picked", "S_NO_PROMPTS": "No prompts yet", "S_PROF_DRAG": "drag to reorder",
		"S_PROF_NAME_PH": "what to call the prompt", "S_PROF_TEST_PH": "type a phrase to try it",
		"S_PF_NEW": "New prompt", "S_PF_EDIT": "Editing the prompt",
		"S_POST_NO_MODEL": "on, but no model is picked", "S_POST_NO_API": "on, but the server is not set up", "S_POST_BAD": "the server did not answer: %s", "S_POST_NO_PROMPT": "on, but no prompt is checked", "S_API_TEST": "Test connection", "S_API_TEST_RUN": "Checking…", "S_API_TEST_OK": "The server answered", "S_API_CLEAR": "Clear", "S_API_CLEAR_ASK": "Delete the address, the model and the key of the external server? Post-processing goes back to the local model.", "S_RAM_AVAIL": "Memory available: %s GB of %s GB", "S_RAM_OF": "%s GB of %s GB",
		"S_REPL_ADD": "Add a replacement", "S_REPL_FROM_PH": "git hub", "S_REPL_TO_PH": "GitHub",
		"S_REPL_WHOLE": "whole words", "S_REPL_CASE": "case", "S_REPL_EMPTY": "No replacements yet",
		"S_PASTE_DELAY": "Delay before inserting", "S_PASTE_DELAY_SUB": "when the program is not ready for the text yet",
		"S_OVPOS": "Where to show the plate", "S_OVPOS_SUB": "at the cursor — next to where you type; if the app hides it, next to the mouse pointer",
		"S_OVPOS_CARET": "At the cursor",
		"S_OVTEXT": "Show the recognised text", "S_OVTEXT_SUB": "on the plate after insertion, instead of the character count",
		"S_TYPEMODE":       "Insert character-by-character instead of the clipboard",
		"S_RECLANG":        "Source speech language",
		"S_RECAUTO":        "Auto-detect",
		"S_DL":             "Download",
		"S_DEL":            "Delete",
		"S_M_BASE":         "fast, for weak PCs",
		"S_M_SMALL":        "balanced speed and accuracy",
		"S_M_MED":          "more accurate, recommended",
		"S_M_TURBO":        "best accuracy on CPU",
		"S_M_GIGAAM":       "more accurate in Russian, punctuates by itself",
		"S_M_PARAKEET":     "25 European languages, punctuates by itself",
		"S_PUNCT":          "Punctuation and capitals",
		"S_PUNCT_MODEL":    "from the model",
		"S_PUNCT_LLM":      "by the editor model",
		"S_PUNCT_OFF":      "strip",
		"S_SEARCH":         "Find a setting…",
		"S_HOTMODE":        "Mode",
		"S_HOTMODE_HOLD":   "hold",
		"S_HOTMODE_TOGGLE": "toggle",
		"S_SUB_HOTMODE":    "hold the keys, or press once to start and once to stop",
		"S_GRP_GENERAL":    "General",
		"S_GRP_SPEECH":     "Speech processing",
		"S_GRP_INFO":       "Info",
		"S_NAV_POST":       "Post-processing",
		"S_NAV_HELP":       "Help",
		"S_HELP_TOC":      "On this page",
		"S_HELP_TOC_SHOW": "Show the contents — the window gets wider",
		"S_HELP_TOC_HIDE": "Hide the contents and restore the window width",
		"S_NAV_CONTACTS":   "Contacts",
		"S_STATE_ACTIVE": "Recognizing",
		"S_STATE_USED": "Models in use",
		"S_STATE_INST": "Installed locally",
		"S_STATE_INST_SUB": "models on disk, ready to be assigned",
		"S_PRESETS": "Which model serves which language",
		"S_PRESETS_HINT": "Click a language — the model choice for it unfolds below. Languages without their own model use the Auto-detect one.",
		"S_MFOLDER": "Your own model",
		"S_DICT_SAVE": "Save",
		"S_OWNM_SUB": "Add a local speech recognition model",
		"S_OWNM_ONEFILE": "One file",
		"S_OWNM_FOLDERF": "A folder with the model files",
		"S_OWNM_S1": "Open the models folder",
		"S_OWNM_S1S": "Destination folder:",
		"S_OWNM_S2": "Copy the model",
		"S_OWNM_S2S": "Pick one of the supported layouts",
		"S_OWNM_S3": "Restart the app",
		"S_OWNM_S3S": "The model appears for the languages it supports",
		"S_AS_AUTO": "as Auto-detect",
		"S_REC_CHIP": "recommended",
		"S_BACK_AUTO": "Back to Auto-detect",
		"S_LANGS_COUNT": "languages: %d",
		"S_LANGS_UNKNOWN": "languages: unknown",
		"S_TR_EN": "translates to English",
		"S_TR_LIST": "translates: %s",
		"S_DL_GOING": "downloading:",
		"S_OPEN_FOLDER": "Open the folder",
		"S_UNLOAD": "Unload from memory",
		"S_UNLOAD_SUB": "frees the memory; the next dictation loads the model again",
		"S_UNLOAD_GO": "Unload",
		"S_UNLOADED": "Unloaded",
		"S_NOT_FOR_LANG": "%s does not recognize this language",
		"S_MANUAL_NOTE": "Cannot be downloaded from the app — the licence forbids redistribution. Download the archive yourself and unpack it into models/moonshine-uk.",
		"S_MANUAL_LINK": "Download yourself",
		"S_HF_FIT": "only those that fit this computer",
		"S_HF_HIDDEN": "hidden: %s",
		"S_WIZ_SKIP_DL": "Download later",
		"S_WIZ_SKIP_NOTE": "Dictation will not work without a model. You can download one under Languages & models.",
		"S_M_GIGAAM2": "the previous generation of the Russian model: same speed, but no punctuation",
		"S_M_MOONUK": "Ukrainian Moonshine model: fast and light, no punctuation",
		"S_M_LOCAL": "found in the models folder; its properties are unknown, so there are no bars",
		"S_ALL_LANGS": "all languages",
		"S_OVPOS_SCHEME_SUB": "click the screen — the plate lands there",
		"S_OVDRAG": "Drag it where you want",
		"S_OVMON": "Screen",
		"S_OVMON_SUB": "which monitor shows the plate",
		"S_OVMON_CURSOR": "The screen with the cursor",
		"S_M_NEMOTRON": "types as you speak: the text shows on the plate while you talk; 40 languages, punctuates by itself",
		"S_M_TINY": "the smallest and fastest, for very weak machines; noticeably less accurate",
		"S_STATE_LOADED": "In memory right now",
		"S_STATE_LOADED_SUB": "models unload themselves after idling",
		"S_STATE_WEEK": "This week",
		"S_ST_SUMMARY": "Summary", "S_ST_OVERLAY": "Plate on screen", "S_ST_BEEP": "Sound cue", "S_ST_AUTORUN": "Start with Windows", "S_ST_POST": "Post-processing", "S_ST_LOCAL": "local", "S_ST_CHECKED": "checked %s", "S_ST_GB": "%s GB", "S_ST_ON_M": "on", "S_ST_OFF_M": "off", "S_ST_MIC_OK": "signal is fine", "S_ST_MIC_BAD": "the microphone is silent", "S_ST_CHECK": "Check", "S_ST_RECOG": "recognized by %s", "S_ST_VER": "Version %s", "S_ST_LATEST": "latest", "S_ST_OUTDATED": "outdated", "S_ST_UPD_OK": "you have the latest version", "S_ST_UPD_DL": "Downloading the update…",
		"S_ST_QUICK": "Quick settings",
		"S_ST_MODELS": "Models",
		"S_ST_USAGE": "This week",
		"S_ST_READY": "Ready to dictate",
		"S_ST_OFF": "Switched off in the tray",
		"S_ST_OFF_SUB": "the shortcut does nothing until you switch it back on",
		"S_ST_ENABLE": "Switch on",
		"S_ST_GOTO": "Open this setting on its own tab",
		"S_ST_HOTKEY_GO": "Change the shortcut",
		"S_ST_UPD_LAST": "Version %s — the latest",
		"S_ST_UPD_HAVE": "Version %s is available",
		"S_ST_MEM": "%s GB free of %s",
		"S_ST_MEM_SUB": "in memory: %s · on disk: %d models, %s GB",
		"S_ST_MEM_NONE": "nothing",
		"S_ST_LANG": "Speech language",
		"S_ST_ASR": "Recognition",
		"S_ST_ON": "on",
		"S_ST_OFF_W": "off",
		"S_ST_ON_F": "on",
		"S_ST_OFF_F": "off",
		"S_ST_ACTIVE": "active",
		"S_ST_IDLE": "not started",
		"S_ST_DISK": "on disk, %s",
		"S_ST_USAGE_SUB": "%d characters · %d today · %d characters on average",
		"S_WEEK_OTHER": "other",
		"S_ST_NO_WEEK": "no dictations this week",
		"S_ST_AUTORUN_SUB": "the app will not start by itself",
		"S_ST_OVERLAY_SUB": "shown while recording",
		"S_REPL_LANG": "Rule language",
		"S_REPL_LANG_ALL": "all languages",
		"S_M_CANARY": "English, German, Spanish, French — and translates between them by itself",
		"S_M_QWEN3": "about 30 languages, punctuates by itself; the heaviest and most accurate in the catalog",
		"S_POSTAPI": "External server",
		"S_POST_HINT": "Edits the recognized text by your prompts: removes filler words, fixes punctuation, changes style. Off — the text is inserted as recognized.",
		"S_POST_MODEL": "Model",
		"S_SRC_LOCAL": "Local",
		"S_SRC_USED": "in use",
		"S_HF_GO": "Search",
		"S_POSTAPI_HINT": "Empty by default — all post-processing runs locally. Enter an address and the prompts run on an external server: OpenAI, Groq, your own vLLM — anything with a compatible API.",
		"S_POSTAPI_URL": "Address",
		"S_POSTAPI_URL_SUB": "empty = the local model; example: https://api.openai.com/v1",
		"S_POSTAPI_MODEL": "Model",
		"S_POSTAPI_KEY": "API key",
		"S_POSTAPI_KEY_SET": "key saved (encrypted with Windows DPAPI)",
		"S_POSTAPI_KEY_NONE": "no key",
		"S_POSTAPI_SAVE": "Save key",
		"S_POSTAPI_TIMEOUT": "Response timeout", "S_SEC_SHORT": "s",
		"S_POSTAPI_WARN": "⚠ The recognized text of your dictations will go to this address. Audio never leaves. The key is stored encrypted.",
		"S_POSTAPI_ASK": "Send recognized text to %s? Audio stays on this computer, but the text will leave it.",
		"S_POSTAPI_BADGE": "external server",
		"S_HIST_ADD":       "Add",
		"S_CONTACT_MAIL":   "Email",
		"S_CONTACT_TITLE":   "Get in touch",
		"S_ABOUT_DEPS":      "External dependencies",
		"S_ABOUT_DEPS_HINT": "The libraries and programs this app is built from, and their licences. The address under each one opens the project page.",
		"S_DEP_WHISPER":     "Runs Whisper models. It is a separate program, whisper-server.exe, sitting next to the app: the app hands it the recorded audio and gets text back. Needed whenever a Whisper model serves your language, and for translating speech into English.",
		"S_DEP_LLAMA":       "Runs the editor model used by post-processing — the one that rewrites recognised text according to your prompts. A separate program, llama-server.exe; it only starts when post-processing is on and a local model is chosen.",
		"S_DEP_SHERPA":      "The second recognition engine: GigaAM, Parakeet, Canary, Qwen and Moonshine run on it. It is what handles Russian speech by default. Lives in sherpa-server.exe and sherpa-online-server.exe.",
		"S_DEP_GGML":        "The maths library whisper.cpp and llama.cpp are built on. It never runs on its own — it does the number crunching inside both engines, using your processor's AVX instructions.",
		"S_DEP_ONNX":        "Runs ONNX models inside sherpa-onnx. Without it GigaAM, Parakeet and the rest of that family would not work.",
		"S_DEP_WEBVIEW":     "The bridge between Go and WebView2: it opens this very settings window and lets the page call into the app.",
		"S_DEP_WV2RT":       "The Windows component that actually draws the settings window. It usually comes with Edge; the app only asks it to show its page.",
		"S_DEP_MALGO":       "Microphone capture from Go: opens the chosen device, streams audio while you speak and feeds the level meter on the overlay.",
		"S_DEP_MINIAUDIO":   "The audio library inside malgo. It also plays the start and stop cues.",
		"S_DEP_WS":          "The websocket link to sherpa-server: the app sends audio there and receives the recognised text.",
		"S_DEP_XSYS":        "Windows system calls from Go: global shortcuts, the clipboard, windows, the registry, the tray icon.",
		"S_DEP_WINLOADER":   "Loads the WebView2 libraries into the process. A helper part of go-webview2, never used on its own.",
		"S_DEP_PLEX":        "The interface font. It is embedded in the app, so the window looks the same on any machine even when the font is not installed.",
		"S_DEP_HF":          "The catalogue models are downloaded from. It is contacted only when you press Download yourself.",
		"S_CONTACT_HINT":    "A bug, an idea, a question about a setting — write an email if it is personal, or open an issue if it is a bug.",
		"S_CONTACT_REPO":    "Repository",
		"S_CONTACT_ISSUES":  "Bugs and ideas",
		"S_CONTACT_WRITE":   "Write an email",
		"S_CONTACT_OPEN":    "Open",
		"S_DICT_MODEL":     "Recognition model",
		"S_LIB_ACC":        "accuracy",
		"S_LIB_SPD":        "speed",
		"S_NAV_STATE":      "Status",
		"S_NAV_DICT":       "Controls & behaviour",
		"S_NAV_MIC":        "Microphone",
		"S_NAV_MODELS":     "Languages & models",
		"S_NAV_TEXT":       "Rules",
		"S_NAV_TR":         "Translation",
		"S_NAV_SYSTEM":     "System",
		"S_NAV_ABOUT":      "About",
		"S_STATE_HINT":     "hold it and speak — the text lands where the cursor is",
		"S_STATE_PROC":     "Post-processing",
		"S_STATE_MEM":      "Memory",
		"S_STATE_MEM_SUB":  "models stay loaded, the first phrase has no delay",
		"S_SUB_MINMS":      "ignores accidental key presses",
		"S_SUB_ENTER":      "sends the message right away",
		"S_SUB_CLIP":       "images and files come back as they were",
		"S_SUB_TYPE":       "helps where a field refuses a clipboard paste",
		"S_SUB_THREADS":    "more threads is not always faster — measure on your machine",
		"S_SUB_PUNCT":      "where punctuation and capitals come from",
		"S_SUB_TRTARGET":   "the text is translated into it; the plate dialog offers it first",
		"S_SUB_AUTOSTART":  "turn off if you start the server yourself",
		"S_SUB_PORT":       "the recognizer restarts itself",
		"S_SUB_UPD":        "the only network request besides model downloads",
		"S_STATE_LAST":     "Last dictation",
		"S_STATE_COPY":     "Copy",
		"S_SEC_OVERLAY":    "On-screen overlay",
		"S_SEC_SERVICE":    "Service",
		"S_SEC_LLM":        "Editor model",
		"S_NOT_INSTALLED":  "not installed",
		"S_CHANGE_MODEL":   "Change",
		"S_RETRY":          "Try again", "S_BERR_OPEN": "Open the server settings",
		"S_PICK_MODEL":  "Pick",
		"S_MIC":         "Microphone",
		"S_MIC_DEFAULT": "System default",
		"S_MIC_CHECK":   "Check the microphone", "S_MIC_CHECK_SUB": "three seconds of recording, then a verdict: level, clipping, whether there is speech", "S_MIC_CHECKING": "Checking…",
		"S_MCHECK": "Check installed models", "S_MCHECK_SUB": "compares model files with reference hashes", "S_MCHECK_GO": "Check", "S_MCHECK_RUN": "Checking…",
		"S_HIST_INSERT": "Paste",
		"S_MIC_REFRESH":     "Refresh list",
		"S_MIC_LEVEL":       "Input level",
		"S_MIC_QUIET":       "quiet",
		"S_THREADS":         "CPU threads",
		"S_MINMS":           "Min recording, ms",
		"S_MAXSEC":          "Max recording, s",
		"S_AUTOSTART":       "Start whisper-server automatically",
		"S_PORT":            "Port",
		"S_SERVEREXE":       "whisper-server path",
		"S_THEME_PINK":      "Pink",
		"S_THEME_BLUE":      "Blue",
		"S_THEME_AMBER":     "Amber",
		"S_THEME_GREEN":     "Green",
		"S_THEME_SUB":       "the colour of the window, the plate and the tray icon",
		"S_SKIN_TERMINAL": "Terminal",
		"S_SKIN_SOFT": "Soft",
		"S_SKIN_PAPER": "Document",
		"S_SKIN_SUB": "font, shape, effects and motion",
		"S_SKIN": "Design",
		"S_WND_CLOSE":       "Close the window",
		"S_WND_MIN":         "Hide to the tray",
		"S_WND_RESTORE":     "Back to the previous size",
		"S_WND_MAX":         "Fill the screen",
		"S_THEME_NEON":      "Neon",
		"S_THEME_EDITOR":    "Editor",
		"S_THEME":           "Colour",
		"S_UPD_FOUND":       "Version %s is out",
		"S_RELOAD_CFG_BTN":  "Re-read",
		"S_RELOAD_CFG_SUB":  "if you edited the file by hand",
		"S_RELOAD_CFG":      "Re-read config.json",
		"S_RESET_ALL_ASK":   "Put every setting back to factory? Models, history and prompts stay where they are.",
		"S_RESET_ALL_BTN":   "Reset",
		"S_RESET_ALL_SUB":   "everything back to factory, except models and history",
		"S_RESET_ALL":       "Reset the settings",
		"S_EXE_WARN":        "The app finds whisper-server next to itself. With a hand-written path, moving the folder stops recognition from starting. Change it?",
		"S_EXE_RESET":       "Reset",
		"S_SERVEREXE_SUB":   "filled in for you; change it only if the server lives somewhere else",
		"S_SERVERURL":       "Remote recognition server (URL)",
		"S_URLHINT":         "If set, the local server is not started",
		"S_REMOTE_WARN":     "Audio will be sent to this server. Local mode is off.",
		"S_REMOTE_ASK":      "Audio will stop being processed on this computer and will be sent to %s. Turn remote mode on?",
		"S_STT_SRV":         "Recognition server",
		"S_STT_SRV_HINT":    "Whisper models are run by a separate program. It can work on this computer or on another one — pick which one to use.",
		"S_SRV_LOCAL":       "On this computer",
		"S_SRV_REMOTE":      "On another computer",
		"S_SRV_REMOTE_HINT": "The same whisper-server, running somewhere else: a home server, a machine with a graphics card, the computer next door.",
		"S_SRV_K_AUTO":      "autostart",
		"S_SRV_K_FILE":      "file",
		"S_SRV_K_ADDR":      "address",
		"S_SRV_K_CHECK":     "last check",
		"S_SRV_NEAR":        "whisper-server.exe next to the app",
		"S_SRV_NOADDR":      "not set",
		"S_SRV_NOCHECK":     "never checked",
		"S_SRV_LOCAL_DLG":   "Local recognition server",
		"S_SRV_ADDR":        "Server address",
		"S_SRV_ADDR_SUB":    "the address of whisper-server on the other machine, port included",
		"S_SRV_ON":          "on",
		"S_SRV_OFF":         "off",
		"S_SRV_K_THREADS":   "CPU threads",
		"S_SRV_K_PORT":      "port",
		"S_SRV_DOWN":       "Recognition is unavailable",
		"S_SRV_DOWN_WHY":   "the remote recognition server is not set up — set its address in the settings",
		"S_SRV_DOWN_GO":    "Open the server settings",
		"S_SRV_WARN_NOW":    "Dictation does not work right now: the remote server is chosen, but its address is not set.",
		"S_SRV_WARN_LATER":  "As soon as a Whisper model is picked, recognition will not work: the remote server address is not set.",
		"S_REMOTE_BADGE":    "REMOTE",
		"S_REMOTE_ABOUT":    "A remote server is set: audio is sent to it, and the promise above does not hold while it is on.",
		"S_SAVED":           "Saved",
		"S_STATE_GET":       "Download",
		"S_OK":              "Yes",
		"S_CANCEL":          "Cancel",
		"S_DL_ASK":          "The \"%s\" model is not downloaded (%s). Start downloading?",
		"S_DL_START":        "Download",
		"S_DL_CANCEL":       "Cancel the download",
		"S_NOT_FOUND":       "none",
		"S_UPD":             "Updates",
		"S_UPD_CHECK":       "Check for updates",
		"S_UPD_AUTO":        "Check on startup",
		"S_UPD_NONE":        "You are on the latest version",
		"S_BADGE_MODELS":    "Installed models",
		"S_BADGE_MISS":      "A model is not downloaded",
		"S_BADGE_SYSTEM":    "Warnings need attention",
		"S_BADGE_HIST":      "Entries in history",
		"S_LOG_OPEN":        "Open the log",
		"S_LOG":             "Log",
		"S_LOG_SUB":         "everything the app writes about itself",
		"S_UPD_AVAIL":       "Version %s is available.",
		"S_UPD_GO":          "Update",
		"S_UPD_ERR":         "Update check failed",
		"S_UPD_DL":          "Downloading the update…",
		"S_AUTHOR_HTML": "<p style=\"font-size:15px;letter-spacing:2px\"><b>Vitalii Yemets</b></p>" +
			"<p>Author and developer of {app} — a local dictation tool for Windows: voice becomes text right at the cursor, no clouds, no subscriptions.</p>" +
			"<p>The project is open: source code, build pipeline and fresh releases live on GitHub.</p>" +
			"<ul>" +
			"<li>Repository: <span class=\"lnk\" onclick=\"appRepoLink()\">github.com/Vitalii-Yemets/holdtotype</span></li>" +
			"<li>Author profile: <span class=\"lnk\" onclick=\"appAuthorLink()\">github.com/Vitalii-Yemets</span></li>" +
			"</ul>" +
			"<p>Found a bug or have an idea — open an issue in the repository.</p>",
		"S_ABOUT_HTML": "<p><b>Voice → text at the cursor position.</b></p>" +
			"<p>Place the cursor in any input field, hold the shortcut, say a phrase, release — the text is inserted automatically.</p>" +
			"<p>Everything runs fully local and offline: audio and text never leave your computer.</p>" +
			"<p class=\"wh\">Tech stack</p>" +
			"<ul>" +
			"<li><b>Go + WinAPI</b> — the client: tray, global hotkeys, overlay, text insertion;</li>" +
			"<li><b>WebView2</b> — the settings window;</li>" +
			"<li><b>whisper.cpp</b> — speech recognition (Whisper models, GGML format);</li>" +
			"<li><b>llama.cpp</b> — post-processing and translation (LLM models, GGUF format);</li>" +
			"<li><b>miniaudio</b> — microphone capture;</li>" +
			"<li><b>Hugging Face</b> — the catalog models are downloaded from.</li>" +
			"</ul>" +
			"<p>Logs never exceed ~2 MB on disk.</p>",
		"S_HELP_HTML": "<p class=\"wh\">How it works</p>" +
			"<p>Hold the shortcut — recording starts (the overlay at the bottom of the screen shows your voice level). Release — the audio is transcribed, optionally translated and run through prompts, and the final text is inserted at the cursor. The ✕ on the overlay cancels at any stage.</p>" +
			"<p>Full pipeline: <b>record → recognise (GigaAM for Russian, Whisper for the rest) → translate (if enabled) → prompts (LLM) → paste</b>. Every stage is visible on the overlay.</p>" +
			"<p class=\"wh\">First run</p>" +
			"<p>The very first launch opens a five-step wizard: the interface language, the language you will dictate in (it picks and downloads the model for you), the shortcut and microphone with a live level bar, a field to try a dictation into, and — last — starting with Windows. You can skip it and everything still works; <b>{exe} -wizard</b> brings it back. Upgrades never see it.</p>" +
			"<p class=\"wh\">Overlay</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\"></span><span>Speak…</span><span class=\"mock-bars\"><i style=\"height:6px\"></i><i style=\"height:13px\"></i><i style=\"height:9px\"></i><i style=\"height:16px\"></i><i style=\"height:7px\"></i><i style=\"height:12px\"></i><i style=\"height:5px\"></i></span><span class=\"mock-x\">✕</span></div></div>" +
			"<ul>" +
			"<li><b>Speak…</b> — recording: a red dot and live voice-level bars.</li>" +
			"<li><b>Transcribing…</b> — Whisper is processing; while translating — \"Translating\", while running prompts — \"Editing: name (1/2)\".</li>" +
			"<li><b>Inserted: N chars</b> — done; on errors or silence a short reason is shown.</li>" +
			"<li>The ✕ on the right cancels at any stage; the overlay never steals input focus. The plate and its animation can be turned off in the Dictation section.</li>" +
			"<li>Where the plate appears — bottom of the screen, top, or at the cursor — and whether it shows the recognised text itself instead of a character count, is set on the Dictation section.</li>" +
			"<li>While the plate is asking something its top line says so — \"Waiting for your answer\" — and the dot stops pulsing. Every answer carries a number: 1…9 pick one, Enter takes the highlighted one, Esc cancels everything; the keys are spelled out at the right of the same row. Ten seconds before the recording limit an amber countdown runs on the plate.</li>" +
			"<li>The title bar carries three buttons: hide to the tray, fill the screen and close. The same button brings a filled window back to the size it had, and the size you set with the mouse is remembered — filling the screen does not replace it. The window never goes below 760×500, where the rows and cards stop fitting.</li>" +
			"<li>Long names — a device, a model, a file — are cut with an ellipsis on the Status cards so the cards line up; the whole name appears as a hint if the pointer rests on the card. The hints are drawn in the colours of the current skin, not the system ones.</li>" +
			"<li>The look comes from two lists in the System section. Design sets the font, the shape, the border width, the halo and the character of the motion; there are three — Terminal (green, the default), Editor (flat grey, no halo) and Neon (violet, rounded). Colour is offered to Terminal alone and changes nothing but the colour of the window, the plate and the tray icon: green, amber, blue, pink. The other designs bring their own colours. The choice applies at once, with no restart.</li>" +
			"</ul>" +
			"<p class=\"wh\">The translation question</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\" style=\"background:#ffb347;box-shadow:0 0 8px rgba(255,179,71,.8)\"></span><span>Transcribing…</span><span class=\"mock-x\">✕</span></div><div class=\"mock-ask\"><span>Translate to:</span><span class=\"mock-btn on mock-cd\">EN</span><span class=\"mock-btn\">DE</span><span class=\"mock-btn\">No translation</span></div></div>" +
			"<p>The plate itself asks, on a second line, the moment you let the shortcut go — in the \"always ask\" and \"ask with a countdown\" modes. The buttons come from \"Languages in the dialog\"; the target language is highlighted. With a countdown, a line under that button shrinks: when it runs out, the highlighted language is used. <b>No translation</b> inserts the text as it was heard; the ✕ on the plate cancels the whole operation. The keyboard works too: Enter takes the highlighted answer, 1…9 pick a button by number, Esc cancels.</p>" +
			"<p class=\"wh\">Safe insertion</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\" style=\"background:#ffb347;box-shadow:0 0 8px rgba(255,179,71,.8)\"></span><span>Transcribing…</span><span class=\"mock-x\">✕</span></div><div class=\"mock-ask\"><span style=\"color:var(--amber)\">Focus changed — insert?</span><span class=\"mock-btn on mock-cd\">Insert here</span><span class=\"mock-btn\">Copy</span></div></div>" +
			"<ul>" +
			"<li>The target window is captured the moment you press the hotkey. If focus changed while the speech was being processed, nothing is pasted — the plate asks on its second line: <b>Insert here</b> (into the current window), <b>Copy</b> (to the clipboard) or <b>Keep it</b>. Say nothing and after 30 seconds it picks Keep it by itself: nothing is inserted anywhere, and the text stays in Last Result and on the clipboard.</li>" +
			"<li>Enter after paste is pressed only when the target window has not changed.</li>" +
			"<li><b>Last Result</b> — the final text of every dictation is kept in memory until the next one; the tray menu has \"Copy last result\". A failed paste or focus change never loses a dictation.</li>" +
			"</ul>" +
			"<p class=\"wh\">Checking the microphone</p>" +
			"<p>The Test button on the Microphone section records three seconds and takes them apart: peak level in decibels, how much of the recording actually holds speech, and how many samples were clipped. The answer comes in words: sounds good, too quiet — raise the level in Windows, clipping — lower it, no speech heard — is the right microphone picked. The same numbers are measured after every dictation and written to the log; when recognition comes back empty the plate names the reason — too quiet, clipping or silence — instead of just saying it heard nothing.</p>" +
			"<p class=\"wh\">Pasting from the history</p>" +
			"<p>Every entry in the history has a Paste button: it brings back the window you opened the settings from and pastes the text there, like an ordinary dictation. When there is nowhere to go back to, the text simply lands on the clipboard and the program says so.</p>" +
			"<p class=\"wh\">The lists in one file</p>" +
			"<p>Replacements and voice commands can be saved into a single .json file and loaded on another computer — the buttons under the command list on the Text section. Loading overwrites nothing: only the lines that are not there yet are added, and the program says how many were added and how many were skipped.</p>" +
			"<p class=\"wh\">File integrity</p>" +
			"<p>Every model in the catalog has a known SHA-256 reference hash. After a download the file is compared against it: no match — the file is deleted and the download can be repeated. The Check button on the Models section compares the models already installed the same way, and when the program updates itself the downloaded installer is checked too — a foreign file will not be started.</p>" +
			"<p class=\"wh\">History of dictations</p>" +
			"<p>The History section in the left column keeps what you dictated: text only, on this computer only, audio is never kept. It is off by default and turns on with one switch in the same place. Entries are kept for a set number of days and up to a set count, older ones drop out on their own; \"Never record from these programs\" lists, separated by commas, the ones nothing should be saved from — password managers, banking apps. Search covers both the text and the program name, the button next to an entry puts it on the clipboard, and Clear removes everything at once along with the <b>{app}-history.json</b> file.</p>" +
			"<p class=\"wh\">Voice commands</p>" +
			"<p>Under the replacements on the Text section there is a list of commands: what you say turns into an action instead of words. \"New line\" and \"new paragraph\" put in a break — models never do; \"cancel\" throws the whole dictation away and inserts nothing; \"insert text\" drops in anything you like, an emoticon included. The button next to the list fills it with the usual phrases in the language of the interface. Commands are matched as whole words and run after the replacements, so the prompts and the translation already get the finished text. Spare spaces around the breaks are cleaned up. The field below tries both replacements and commands on any phrase: a break shows as ⏎.</p>" +
			"<p class=\"wh\">Replacements after recognition</p>" +
			"<p>On the Text section you can list what the model mishears and what it should become: \"git hub\" → GitHub, surnames, in-house terms. Replacements run right after recognition, before the prompts, so the editor already gets the right words. Translation into English happens inside recognition, so replacements see the translated text. By default they match whole words and ignore case; the two switches next to each row change that. Rules apply from top to bottom. The field at the bottom tries them on any phrase without dictating.</p>" +
			"<p class=\"wh\">Rules per application</p>" +
			"<p>On the Dictation section you can set rules for particular programs: what to insert with (the clipboard or character by character), whether to press Enter, how long to wait before inserting and which prompts to apply. A program is named by its file — <b>chrome.exe</b>; one rule can list several, separated by commas, and a trailing asterisk catches every name with that beginning. The first matching rule wins; with no rules, or none that match, everything works as in the general settings. The button next to the list fills in the program you last inserted into.</p>" +
			"<p class=\"wh\">The main settings</p>" +
			"<ul>" +
			"<li><b>Keyboard shortcut</b> — the main dictation hotkey. Any combination can be captured; left/right modifiers are distinguished. Dictation, translation and profile hotkeys must be unique — a duplicate blocks saving.</li>" +
			"<li><b>UI language</b> — switches instantly; \"Same as system\" follows Windows.</li>" +
			"<li><b>Recognition language</b> — a hint for Whisper; \"auto\" detects the language from speech.</li>" +
			"<li><b>Sound</b> — start/stop recording cues: several themes plus Windows system sounds, ▶ previews.</li>" +
			"<li><b>Enter after paste</b> — automatically sends the dictated text (handy in messengers).</li>" +
			"<li><b>Restore clipboard</b> — puts the previous clipboard content back in full after pasting, including images, files and rich text. When the content cannot be snapshotted, the clipboard is left untouched and the text is typed character by character.</li>" +
			"<li><b>Overlay and animation</b> — the status indicator at the bottom of the screen; animation can be turned off.</li>" +
			"<li><b>Type mode</b> — inserts text by simulating keystrokes instead of Ctrl+V, for fields where clipboard paste does not work.</li>" +
			"</ul>" +
			"<p class=\"wh\">Recognition</p>" +
			"<div class=\"mock\"><div class=\"mock-row\"><b>Auto-detect</b><span style=\"margin-left:auto\">Turbo (q5)</span></div>" +
			"<div class=\"mock-row\"><b>Russian</b><span style=\"margin-left:auto\">GigaAM v3</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-note\">English · as Auto-detect — Turbo (q5)</span><span style=\"margin-left:auto\">▼</span></div></div>" +
			"<ul>" +
			"<li><b>Which model serves which language</b> — the Languages & models tab is a list of languages. Click one — the models that can serve it unfold below: the assigned and the recommended first, the missing ones with a size and a download arrow. A click on a card is the choice; a missing model downloads itself and takes over once ready. Languages without a model of their own inherit the Auto-detect one and are drawn dim.</li>" +
			"<li><b>The catalogue</b> — Whisper: Base (fast, for weak PCs), Small (the balance), Medium and Turbo (more accurate and slower; “q5” is the quantized cut: a bit smaller and faster with almost no quality loss), and they also translate to English; GigaAM v3 is sharper on Russian and punctuates itself; Parakeet v3 covers 25 European languages; Nemotron 3.5 types as you speak. Downloads come from the official Hugging Face repositories, every file checked against its reference hash.</li>" +
			"<li><b>Your own model</b> — a Whisper single ggml-*.bin file or a sherpa-onnx model folder (encoder.onnx, decoder.onnx, tokens.txt) will do. Put it into the models folder next to the app and restart it — the model appears in the choice for matching languages; its powers are unknown, so it is shown honestly, without bars.</li>" +
			"<li>whisper-server keeps the model in memory between phrases — the first dictation after startup is slower (loading), afterwards recognition takes 1–3 seconds.</li>" +
			"<li><b>Dictionary</b> — comma-separated terms, names and abbreviations. A hint for Whisper's \"ear\" so rare words are recognized correctly; not commands.</li>" +
			"<li><b>Microphone and speed</b> — microphone selection with a live level meter (speak and the bar moves, so you know the device is heard), CPU threads (more = faster transcription), minimum recording length (filters accidental presses), maximum length (auto-stop). If the chosen device is unplugged the app falls back to the system default; a recording with no speech is never sent for recognition — it reports \"Silence\" instead.</li>" +
			"<li><b>Server</b> — whisper-server starts automatically and runs locally. You can change the port or point to an external server URL — then the local one is not used.</li>" +
			"<li><b>Translation</b> — all translation is done by Whisper: to English via its native translate mode, to other languages <b>experimentally</b>, by forcing the output language (quality depends on the language pair; major languages work best). The Turbo model is not trained for translation — the settings show a warning when it is active. \"Always translate to the target language\" makes every main-hotkey dictation translate to the chosen target with no questions. With the checkbox off, the ask mode applies: always or with a timeout — a language dialog appears before transcription and the target is used when time runs out. The separate translation hotkey translates once without affecting normal dictation. Settings that do not apply in the current mode are greyed out automatically.</li>" +
			"</ul>" +
			"<p class=\"wh\">Post-processing (LLM)</p>" +
			"<p>An optional second layer: a local language model (llama.cpp) edits the transcribed text according to your prompts — removes filler words, changes style, formats. Fully offline, CPU only.</p>" +
			"<div class=\"mock\"><div class=\"mock-row\" style=\"color:var(--dim)\">▸ Qwen/Qwen2.5-3B-Instruct-GGUF<span style=\"margin-left:auto;color:var(--faint)\">2024-09-20</span><span>↓303k</span><span>↗</span></div>" +
			"<div class=\"mock-row\" style=\"padding-left:18px\">q4_k_m.gguf<span style=\"margin-left:auto\">1.9 GB</span><span style=\"color:var(--green)\">● ≈3.7 GB</span><span style=\"color:var(--dim)\">⭳</span></div>" +
			"<div class=\"mock-row\" style=\"padding-left:18px\">q8_0.gguf<span style=\"margin-left:auto\">3.6 GB</span><span style=\"color:var(--amber)\">● ≈8.2 GB</span><span style=\"color:var(--dim)\">⭳</span></div></div>" +
			"<ul>" +
			"<li><b>Models</b> — installed LLM editing models; the radio selects the active one (applied at once), ✕ deletes (the active one too — post-processing then turns off). Download progress shows here as well.</li>" +
			"<li><b>Finding a model</b> — GGUF models on Hugging Face by name (e.g. \"qwen2.5 instruct\"). Each repository shows its last update date, download count and a ↗ link to the model page; clicking a row expands its quant files. The ● ≈N GB indicator is compared against the <b>free</b> RAM (shown above the list).</li>" +
			"<li><b>Picking a quant:</b> the number is bits per weight (Q4 — the sweet spot, Q8 — nearly uncompressed, Q3 — saves RAM at a quality cost); K_M beats K_S; IQ4 is the newer generation, better than classic quants at the same size. The ● ≈N GB indicator estimates the RAM needed (file plus context headroom): green fits, amber is tight, red won't fit.</li>" +
			"<li>A 1.5–3B model gives fast editing; 7–9B is noticeably smarter but each pass takes seconds on CPU. llama-server starts on first use and keeps the model warm in memory.</li>" +
			"</ul>" +
			"<p class=\"wh\">Prompts</p>" +
			"<div class=\"mock\"><div class=\"mock-row\"><span class=\"mock-cb on\">✓</span><b>Cleanup</b><span class=\"mock-note\">—</span><span style=\"margin-left:auto;color:var(--dim)\">✎ ✕</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-cb\"></span><b>Business style</b><span class=\"mock-note\">ctrl+alt+f</span><span style=\"margin-left:auto;color:var(--dim)\">✎ ✕</span></div></div>" +
			"<ul>" +
			"<li>A prompt is an instruction for the editing model. Presets ship out of the box: \"Cleanup\" (removes fillers, repetitions and false starts, fixes punctuation) and \"Business style\" (rewrites politely and formally); add your own freely.</li>" +
			"<li>Checked prompts apply to every dictation in order, top to bottom (as a chain: the output of one feeds the next); nothing checked — text is inserted as transcribed.</li>" +
			"<li>The ✎ pencil opens the prompt editor: name, text and a test field that runs a sample through the live model right from Settings. Drag a prompt by the handle on the left to change the order.</li>" +
			"<li>Tip: small models work much better with \"input → output\" examples inside the prompt — all presets are written that way.</li>" +
			"<li>If a profile fails (the model did not respond), the text is inserted without it: the overlay shows \"Inserted without the … profile\" and Enter after paste is not pressed in that case.</li>" +
			"</ul>" +
			"<p class=\"wh\">Dependencies</p>" +
			"<ul>" +
			"<li>Prompts require the editing model to be installed; translation does not depend on it — Whisper does it entirely.</li>" +
			"<li>The editing model loads into memory on first use and stays warm; large models are noticeably slower on CPU.</li>" +
			"<li>Check the memory indicator before downloading: a \"tight\" model can slow down the whole system.</li>" +
			"<li>Greyed-out controls indicate settings that do not participate in the current mode.</li>" +
			"</ul>" +
			"<p class=\"wh\">Install and portability</p>" +
			"<ul>" +
			"<li><b>{setup}</b> — the installer: no admin rights needed, Start Menu shortcut, optional autostart, clean removal via \"Apps & features\".</li>" +
			"<li>The installer downloads nothing by default: the wizard picks and fetches the model on first run. If a model is picked anyway — GigaAM v3 for Russian, Whisper for every other language — the download can be stopped with a button and the installation still finishes. The update-check switch sits there too, and the answer is written into the app's settings.</li>" +
			"<li><b>Portable</b> — just copy the whole folder with the exe (to a USB stick, another PC): settings, models and the log live next to the exe and travel with it. Nothing is written to the registry.</li>" +
			"<li>On first run without a recognition model the app opens the model catalog itself and waits for the download.</li>" +
			"<li>Requirements: Windows 10/11 x64, a CPU with AVX2 (~2013+), WebView2 Runtime for the settings window (included in Windows 11).</li>" +
			"</ul>" +
			"<p class=\"wh\">Tray and files</p>" +
			"<div class=\"mock\" style=\"max-width:290px\"><div class=\"mock-mi dim\">{app} — Ready…</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Settings…</div><div class=\"mock-mi\">Disable</div><div class=\"mock-mi\">Copy last result</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Open config.json</div><div class=\"mock-mi\">Open log</div><div class=\"mock-mi\">About</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Quit</div></div>" +
			"<ul>" +
			"<li>Left-click the tray icon — settings; right-click — the menu. Icon colors: green — ready, red — recording, orange — transcribing, grey — disabled or error.</li>" +
			"<li><b>config.json</b> — all settings; edits made by hand apply through <b>Re-read</b> in the System section. The log and <b>Reset the settings</b> live there too: a reset puts everything back to factory and leaves models, history and prompts alone.</li>" +
			"<li><b>{log}</b> — the log, automatically capped at ~2 MB.</li>" +
			"<li><b>models/</b> — downloaded Whisper and LLM models.</li>" +
			"</ul>",
		"S_WIZ_HELLO":      "Two minutes to set up",
		"S_WIZ_HELLO_TEXT": "{app} turns your voice into text right at the caret: hold the shortcut, say a phrase, let go — the text is there. Everything runs on your machine; the audio never leaves it.",
		"S_WIZ_LATER":      "Everything we pick now can be changed later in the settings.",
		"S_WIZ_T_MODEL":    "Language and model",
		"S_WIZ_MODEL_TEXT": "Tell me the language you will dictate in and I will pick the model. Russian goes to GigaAM, every other language to Whisper.",
		"S_WIZ_T_INPUT":    "Shortcut and microphone",
		"S_WIZ_INPUT_TEXT": "This is the combination you will hold while speaking. Say something and check that the level bar moves.",
		"S_WIZ_T_TRY":      "Try it",
		"S_WIZ_TRY_PH":     "the text will appear here",
		"S_WIZ_T_DONE":     "All set",
		"S_WIZ_DONE_TEXT":  "{app} lives in the tray: left-click the icon for the settings, right-click for the menu. You can dictate into any window that has a text caret.",
		"S_AUTORUN":        "Start with Windows",
		"S_AUTORUN_SUB":    "An entry in the current user's startup list",
		"S_WIZ_SKIP":       "Skip",
		"S_WIZ_BACK":       "Back",
		"S_WIZ_NEXT":       "Next",
		"S_WIZ_FINISH":     "Finish",
		"S_WIZ_WAIT":       "Waiting for the first phrase…",
		"S_WIZ_HEARD":      "Heard:",
		"S_WIZ_HAVE":       "Everything you need is already downloaded",
		"S_WIZ_TRY_TEXT":   "Put the caret in the field below, hold %s, say a phrase and let go.",
		"S_MODEL_READY":    "Model downloaded — pick it to switch",
	},
}
