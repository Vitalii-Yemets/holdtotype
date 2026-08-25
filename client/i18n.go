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
		"err.webview":           "Для окна настроек нужен Microsoft WebView2 Runtime (входит в Windows 11).\nСейчас откроется страница загрузки — установите его и откройте настройки снова.",
		"status.loading":        "Загрузка модели…",
		"status.nomodel":        "Модель распознавания не скачана — выберите её в настройках",
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
		"adv.pick":               "Рекомендую %s.",
		"adv.companion":          "Рядом пригодится %s — для остальных языков и перевода.",
		"adv.ram":                "свободно %d МБ",
		"status.line":            "Готов · %s · %.1f ГБ свободно",
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
		"ov.paused":              "Пауза",
		"status.paused":          "Пауза — запись ждёт",
		"hist.insert.gone":       "запись не найдена",
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
		"err.webview":           "The settings window requires Microsoft WebView2 Runtime (bundled with Windows 11).\nThe download page will open now — install it and open Settings again.",
		"status.loading":        "Loading model…",
		"status.nomodel":        "No recognition model downloaded — pick one in Settings",
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
		"adv.pick":               "I recommend %s.",
		"adv.companion":          "%s makes a good companion — for other languages and translation.",
		"adv.ram":                "%d MB free",
		"status.line":            "Ready · %s · %.1f GB free",
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
		"ov.paused":              "Paused",
		"status.paused":          "Paused — recording is waiting",
		"hist.insert.gone":       "entry not found",
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
		"S_TR_HINT":        "Перевод выполняет Whisper: на английский — штатным режимом, на остальные языки — экспериментально, принудительным языком вывода (качество зависит от пары языков).",
		"S_TR_TURBO":       "⚠ Активная модель Turbo не обучена переводу на английский — для перевода выберите другую модель в разделе «Модели».",
		"S_TR_DEFAULT":     "Всегда переводить на целевой язык",
		"S_TR_TARGET":      "Целевой язык",
		"S_TR_ASK":         "Выбор языка",
		"S_TR_ASK_NEVER":   "Не спрашивать",
		"S_TR_ASK_ALWAYS":  "Спрашивать каждый раз",
		"S_TR_ASK_TIMEOUT": "Спрашивать с таймаутом",
		"S_TR_SECONDS":     "Таймаут, сек",
		"S_TR_LANGS":       "Языки в диалоге",
		"S_DICT_HINT":      "Термины, имена и аббревиатуры через запятую — подсказка слуху, не команды. Работает для Whisper; на русскую речь через GigaAM не влияет. Набор по умолчанию меняется вместе с языком распознавания, пока вы не впишете своё.",
		"S_LLM_HINT":       "Отмеченные профили применяются по очереди, сверху вниз, при обычной диктовке. Ничего не отмечено — текст вставляется как есть. Хоткей профиля применяет разово только его.",
		"S_UPDATED":        "Дата последнего обновления модели",
		"S_FIT_OK":         "поместится",
		"S_FIT_WARN":       "впритык",
		"S_FIT_BAD":        "не хватит памяти",
		"S_RAM":            "Память компьютера:",
		"S_HF_PH":          "Название модели — например, qwen2.5 instruct",
		"S_NO_LLM":         "Пока не установлено ни одной модели — найдите и скачайте в поле поиска ниже.",
		"S_NO_LLM_PROF":    "Промпты станут доступны после установки модели — раздел «Модели», список модели-редактора внизу.",
		"S_PROF_EDIT":      "Редактировать",
		"S_PROF_CLOSE":     "Свернуть",
		"S_CONFIRM_DEL":    "Удалить модель «%s»? Её можно будет скачать заново.",
		"S_DEL_ACTIVE":     "Удалить активную модель «%s»? Распознавание остановится, пока вы не выберете другую — скачать её можно тут же.",
		"S_WIZ_NEED_MODEL": "Сначала скачайте модель — без неё распознавать нечем",
		"S_FREE":           "свободно",
		"S_SUB_PROMPTS":    "Промпты",
		"S_SUB_DICT":       "Словарь",
		"S_PROF_ADD":       "Добавить профиль",
		"S_PROF_NAME":      "Имя",
		"S_PROF_PROMPT":    "Промпт",
		"S_PROF_HOTKEY":    "Хоткей",
		"S_PROF_SET":       "Задать…",
		"S_PROF_CLEAR":     "Сброс",
		"S_PROF_TEST":      "Проверка",
		"S_HOTKEY":         "Сочетание клавиш",
		"S_CHANGE":         "Изменить…",
		"S_UILANG":         "Язык интерфейса",
		"S_AUTO":           "Как в системе",
		"S_BEEP":           "Звуковые сигналы записи",
		"S_SEC_SOUND":      "Звук",
		"S_SEC_BEHAVIOR":   "Поведение",
		"S_SOUND":          "Звук сигналов",
		"S_SND_SPEECH":     "Голос Windows",
		"S_SND_CHIME":      "Колокольчик",
		"S_SND_SOFT":       "Мягкий",
		"S_SND_MARIMBA":    "Маримба",
		"S_SND_BLIP":       "Блип",
		"S_SND_POP":        "Поп",
		"S_AUTOENTER":      "Нажимать Enter после вставки (автоотправка)",
		"S_RESTORE":        "Восстанавливать буфер обмена после вставки",
		"S_OVERLAY":        "Индикатор внизу экрана",
		"S_NAV_HISTORY":    "История", "S_HIST_ON": "Хранить историю диктовок", "S_HIST_ON_SUB": "только текст, на этом компьютере; звук не сохраняется никогда",
		"S_HIST_DAYS": "Сколько дней хранить", "S_HIST_MAX": "Сколько записей хранить",
		"S_HIST_SKIP": "Не записывать из программ", "S_HIST_SKIP_SUB": "через запятую: keepass.exe, 1password.exe",
		"S_HIST_LIST": "Записи", "S_HIST_CLEAR": "Очистить", "S_HIST_COPY": "Копировать",
		"S_HIST_FIND": "Найти в истории…", "S_HIST_EMPTY": "Истории пока нет", "S_HIST_ASK": "Удалить всю историю диктовок?",
		"S_SEC_CMD": "Голосовые команды", "S_CMD_HINT": "Сказанное вслух превращается в перенос строки, знак или отмену вместо того, чтобы попасть в текст. Проверяются целыми словами, применяются сверху вниз, уже после замен.",
		"S_CMD_ADD": "Добавить команду", "S_CMD_PRESET": "Добавить стандартные", "S_CMD_PH": "новая строка",
		"S_CMD_NEWLINE": "перенос строки", "S_CMD_PARAGRAPH": "новый абзац", "S_CMD_TEXT": "подставить текст", "S_CMD_CANCEL": "отменить диктовку",
		"S_CMD_TEXT_PH": "что подставить", "S_CMD_EMPTY": "Команд пока нет", "S_CMD_DEL": "Удалить команду",
		"S_CMD_P_NEWLINE": "новая строка", "S_CMD_P_PARAGRAPH": "новый абзац", "S_CMD_P_CANCEL": "отмена",
		"S_SEC_REPLACE": "Замены после распознавания", "S_REPLACE_HINT": "Что услышано неправильно, заменяется на правильное — сразу после распознавания, до промптов и перевода. Порядок сверху вниз.",
		"S_REPL_ADD": "Добавить замену", "S_REPL_FROM_PH": "гит хаб", "S_REPL_TO_PH": "GitHub",
		"S_REPL_WHOLE": "целые слова", "S_REPL_CASE": "регистр", "S_REPL_EMPTY": "Замен пока нет",
		"S_REPL_DEL": "Удалить замену", "S_REPL_TEST_PH": "напишите фразу, чтобы проверить замены и команды",
		"S_SEC_RULES": "Правила по приложениям", "S_RULES_HINT": "Для отдельных программ вставка может работать иначе. Выигрывает первое подходящее правило.",
		"S_RULE_ADD": "Добавить правило", "S_RULE_PH": "chrome.exe, msedge.exe",
		"S_RULE_PASTE_INH": "вставка: как везде", "S_RULE_ENTER_INH": "Enter: как везде", "S_RULE_DELAY_NONE": "без задержки", "S_RULE_PROMPT_INH": "промпты: как везде",
		"S_RULE_CLIP": "буфер обмена", "S_RULE_TYPE": "посимвольно", "S_RULE_ENTER_ON": "с Enter", "S_RULE_ENTER_OFF": "без Enter",
		"S_RULE_NOPROMPT": "без промптов", "S_RULE_LAST": "последняя вставка: %s", "S_RULE_EMPTY": "Правил пока нет",
		"S_RULE_DEL": "Удалить правило", "S_RULE_PROMPTS": "Промпты",
		"S_PASTE_DELAY": "Задержка перед вставкой", "S_PASTE_DELAY_SUB": "если программа не успевает принять текст",
		"S_OVPOS": "Где показывать плашку", "S_OVPOS_SUB": "у курсора — рядом с местом ввода; если приложение его не показывает, рядом с указателем мыши",
		"S_OVPOS_BOTTOM": "Внизу экрана", "S_OVPOS_TOP": "Вверху экрана", "S_OVPOS_CARET": "У курсора",
		"S_OVTEXT": "Показывать распознанный текст", "S_OVTEXT_SUB": "в плашке после вставки, вместо числа символов",
		"S_ANIM":           "Анимация записи и распознавания",
		"S_TYPEMODE":       "Вставлять посимвольно, а не через буфер",
		"S_RECLANG":        "Язык распознавания",
		"S_RECAUTO":        "Автоопределение",
		"S_DL":             "Скачать",
		"S_DEL":            "Удалить",
		"S_M_BASE":         "быстрая, для слабых ПК",
		"S_M_SMALL":        "баланс скорости и точности",
		"S_M_MED":          "точнее, рекомендуем",
		"S_M_TURBO":        "максимум точности на CPU",
		"S_M_CUSTOM":       "пользовательская (из config.json)",
		"S_M_GIGAAM":       "точнее на русском, сама ставит знаки препинания",
		"S_ADV_TITLE":      "Подобрать модель",
		"S_ADV_LANGQ":      "На каком языке диктуете",
		"S_ADV_PRIOQ":      "Что важнее",
		"S_ADV_TRQ":        "Нужен перевод",
		"S_ADV_GO":         "Подобрать",
		"S_ADV_LANG":       "Она заточена под ваш язык.",
		"S_ADV_ACC":        "Выбрана по точности.",
		"S_ADV_SPEED":      "Выбрана по скорости.",
		"S_ADV_RAM":        "Более крупные не поместились бы в свободную память.",
		"S_ADV_NONE":       "Свободной памяти не хватает ни на одну модель — закройте лишние программы.",
		"S_F_ALL":          "все",
		"S_F_RU":           "русский",
		"S_F_MULTI":        "много языков",
		"S_F_PUNCT":        "с пунктуацией",
		"S_F_FIT":          "влезает в память",
		"S_PUNCT":          "Пунктуация и заглавные",
		"S_PUNCT_MODEL":    "из модели",
		"S_PUNCT_LLM":      "моделью-редактором",
		"S_PUNCT_OFF":      "убирать",
		"S_SEARCH":         "Найти настройку…",
		"S_LVL_SIMPLE":     "просто",
		"S_LVL_ALL":        "всё",
		"S_MORE":           "Ещё %d настроек",
		"S_LESS":           "Свернуть %d настроек",
		"S_HOTMODE":        "Режим",
		"S_HOTMODE_HOLD":   "удержание",
		"S_HOTMODE_TOGGLE": "фиксация",
		"S_SUB_HOTMODE":    "удерживать клавиши или включать и выключать нажатием",
		"S_GRP_GENERAL":    "Общее",
		"S_GRP_SPEECH":     "Речь",
		"S_GRP_TEXT":       "Текст",
		"S_GRP_INFO":       "Сведения",
		"S_NAV_POST":       "Пост-обработка",
		"S_NAV_HELP":       "Справка",
		"S_NAV_CONTACTS":   "Контакты",
		"S_NAV_STATE":      "Сводка",
		"S_NAV_DICT":       "Диктовка",
		"S_NAV_MIC":        "Микрофон",
		"S_NAV_MODELS":     "Модели",
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
		"S_SUB_TRTARGET":   "на английский переводит Whisper, остальные — экспериментально",
		"S_TR_EXP":         "кроме английского приложение принудительно задаёт язык вывода, а не переводит — текст может остаться на языке речи",
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
		"S_STATE_RU":    "Русская речь",
		"S_STATE_OTHER": "Прочие языки",
		"S_MIC":         "Микрофон",
		"S_MIC_DEFAULT": "Системный по умолчанию",
		"S_MIC_CHECK":   "Проверить микрофон", "S_MIC_CHECK_SUB": "три секунды записи и разбор: громкость, перегруз, есть ли речь", "S_MIC_CHECKING": "Проверяю…",
		"S_PAUSE": "Пауза в записи", "S_PAUSE_SUB": "работает в режиме фиксации: нажали — запись замерла, нажали ещё раз — пошла дальше",
		"S_MCHECK": "Проверить установленные модели", "S_MCHECK_SUB": "сверяет файлы моделей с эталонными хешами", "S_MCHECK_GO": "Проверить", "S_MCHECK_RUN": "Проверяю…",
		"S_HIST_INSERT": "Вставить",
		"S_LISTS_HINT":  "Замены и команды одним файлом — перенести на другой компьютер", "S_LISTS_EXPORT": "Сохранить в файл", "S_LISTS_IMPORT": "Загрузить из файла",
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
		"S_ADV_PRIMARY":     "основная",
		"S_ADV_COMPANION":   "вторая",
		"S_ADV_HAVE":        "уже есть",
		"S_ADV_APPLY":       "Применить",
		"S_ADV_ASK":         "Будут загружены: %s — всего %s. Начать?",
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
		"S_ADV_ROLE_MAIN":   "для выбранного языка",
		"S_ADV_ROLE_SECOND": "для остальных языков и перевода",
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
		"S_HELP_HTML": "<p class=\"wh\">Как это работает</p>" +
			"<p>Зажмите сочетание клавиш — начнётся запись с микрофона (оверлей внизу экрана показывает уровень голоса). Отпустите — звук распознаётся, затем при необходимости переводится и прогоняется через промпты, и готовый текст вставляется в позицию курсора. Крестик на оверлее отменяет операцию на любой стадии.</p>" +
			"<p>Полный конвейер: <b>запись → распознавание (русский — GigaAM, остальные языки — Whisper) → перевод (если включён) → промпты (LLM) → вставка</b>. Каждая стадия видна на оверлее.</p>" +
			"<p class=\"wh\">Первый запуск</p>" +
			"<p>При самом первом запуске открывается мастер из пяти шагов: язык интерфейса, язык диктовки (модель он подберёт и скачает сам), сочетание клавиш и микрофон с живой полоской уровня, поле для пробной диктовки и, последним, запуск вместе с Windows. Мастер можно пропустить — всё работает и без него; вернуть его — запуском <b>{exe} -wizard</b>. При обновлении он не появляется.</p>" +
			"<p class=\"wh\">Оверлей</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\"></span><span>Говорите…</span><span class=\"mock-bars\"><i style=\"height:6px\"></i><i style=\"height:13px\"></i><i style=\"height:9px\"></i><i style=\"height:16px\"></i><i style=\"height:7px\"></i><i style=\"height:12px\"></i><i style=\"height:5px\"></i></span><span class=\"mock-x\">✕</span></div></div>" +
			"<ul>" +
			"<li><b>Говорите…</b> — идёт запись: красная точка и живые столбики уровня голоса.</li>" +
			"<li><b>Распознаю…</b> — Whisper обрабатывает звук; при переводе — «Перевожу», при промптах — «Редактирую: имя (1/2)».</li>" +
			"<li><b>Вставлено: N симв.</b> — готово; при ошибке или тишине — короткое сообщение о причине.</li>" +
			"<li>Крестик ✕ справа отменяет операцию на любой стадии; фокус ввода оверлей не забирает. Показывать плашку и её анимацию можно отключить на «Основных».</li>" +
			"<li>Где показывать плашку — внизу экрана, вверху или у курсора — и показывать ли в ней сам распознанный текст вместо числа символов, настраивается на «Диктовке».</li>" +
			"<li>Пока плашка о чём-то спрашивает, верхняя строка так и говорит — «Жду ответа», и точка перестаёт мигать. У каждой кнопки свой номер: 1…9 выбирают ответ, Enter берёт подсвеченный, Esc отменяет всё; клавиши подписаны справа в той же строке. За десять секунд до предела записи на плашке идёт янтарный обратный отсчёт.</li>" +
			"<li>В заголовке окна три кнопки: свернуть в трей, развернуть на весь экран и закрыть. Развёрнутое окно возвращается к прежнему размеру той же кнопкой, а размер, который вы задали мышью, запоминается — на весь экран он не заменяется. Меньше 760×500 окно не уменьшается: на меньшем строки и карточки перестают помещаться.</li>" +
			"<li>Длинные имена — устройства, модели, файла — на карточках «Состояния» обрезаются многоточием, чтобы карточки стояли ровно; полное имя показывается подсказкой, если задержать на карточке указатель. Подсказки нарисованы в цветах текущего облика, а не системные.</li>" +
			"<li>Внешний вид задаётся двумя списками в разделе «Система». «Дизайн» — это шрифт, форма, толщина рамок, свечение и характер анимации; дизайнов три: «Терминал» (зелёный, по умолчанию), «Редактор» (плоский серый, без свечения) и «Неон» (фиолетовый, со скруглениями). «Цвет» предлагается только «Терминалу» и меняет только цвет окна, плашки и значка в трее: зелёный, янтарный, синий, розовый. У остальных дизайнов цвета свои. Выбор применяется сразу, перезапуск не нужен.</li>" +
			"</ul>" +
			"<p class=\"wh\">Вопрос о языке перевода</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\" style=\"background:#ffb347;box-shadow:0 0 8px rgba(255,179,71,.8)\"></span><span>Распознаю…</span><span class=\"mock-x\">✕</span></div><div class=\"mock-ask\"><span>Переводить на:</span><span class=\"mock-btn on mock-cd\">EN</span><span class=\"mock-btn\">DE</span><span class=\"mock-btn\">Без перевода</span></div></div>" +
			"<p>Спрашивает сама плашка второй строкой, как только вы отпустили сочетание, — в режимах «Спрашивать всегда» и «Спрашивать с таймаутом». Набор кнопок — «Языки в диалоге», подсвечен целевой язык. В режиме с таймаутом под этой кнопкой укорачивается полоска: когда она кончится, применится подсвеченный язык. <b>Без перевода</b> вставляет как услышано; крестик плашки отменяет операцию целиком. Клавиатура тоже работает: Enter — подсвеченный ответ, 1…9 — кнопка по номеру, Esc — отмена.</p>" +
			"<p class=\"wh\">Безопасная вставка</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\" style=\"background:#ffb347;box-shadow:0 0 8px rgba(255,179,71,.8)\"></span><span>Распознаю…</span><span class=\"mock-x\">✕</span></div><div class=\"mock-ask\"><span style=\"color:var(--amber)\">Фокус сменился — вставить?</span><span class=\"mock-btn on mock-cd\">Вставить сюда</span><span class=\"mock-btn\">Копировать</span></div></div>" +
			"<ul>" +
			"<li>Окно-цель запоминается в момент нажатия хоткея. Если за время обработки фокус сменился, текст не вставляется — плашка спрашивает второй строкой: <b>Вставить сюда</b> (в текущее окно), <b>Копировать</b> (в буфер обмена) или <b>Оставить</b>. Если промолчать, через 30 секунд сам выбирается «Оставить»: никуда ничего не вставляется, текст лежит в «Последнем результате» и в буфере обмена.</li>" +
			"<li>Enter после вставки нажимается только если окно-цель не менялось.</li>" +
			"<li><b>Последний результат</b> — финальный текст каждой диктовки хранится в памяти до следующей; в меню трея есть пункт «Копировать последний результат». Ошибка вставки или смена фокуса не теряют надиктованное.</li>" +
			"</ul>" +
			"<p class=\"wh\">Проверка микрофона</p>" +
			"<p>Кнопка «Проверка» на «Микрофоне» записывает три секунды и разбирает их: пиковая громкость в децибелах, доля записи, в которой действительно есть речь, и доля обрезанных отсчётов. Ответ приходит словами: «слышно хорошо», «слишком тихо — прибавьте громкость в Windows», «перегруз — убавьте громкость», «речи не слышно — тот ли микрофон выбран». То же самое считается после каждой диктовки и пишется в лог; если распознать не удалось, плашка скажет причину — тихо, перегруз или тишина, — а не просто «ничего не услышал».</p>" +
			"<p class=\"wh\">Пауза в записи</p>" +
			"<p>В режиме фиксации (нажали — пишет, нажали ещё раз — остановилась) можно задать отдельное сочетание для паузы: на «Диктовке», строка «Пауза в записи». Нажали — запись замерла, плашка показывает «Пауза», и ничего не записывается; нажали ещё раз — запись пошла дальше, а всё сказанное до паузы сохраняется. Предел длины записи на паузе не срабатывает.</p>" +
			"<p class=\"wh\">Повторная вставка из истории</p>" +
			"<p>У каждой записи в истории есть кнопка «Вставить»: она возвращает окно, из которого вы открыли настройки, и вставляет текст туда — как обычная диктовка. Если возвращаться некуда, текст просто кладётся в буфер обмена, и программа об этом скажет.</p>" +
			"<p class=\"wh\">Списки одним файлом</p>" +
			"<p>Замены и голосовые команды можно сохранить в один файл .json и загрузить на другом компьютере — кнопки под списком команд на «Тексте». Загрузка ничего не затирает: добавляются только те строки, которых ещё нет, а сколько добавлено и сколько пропущено, программа скажет.</p>" +
			"<p class=\"wh\">Целостность файлов</p>" +
			"<p>Для каждой модели из каталога известен эталонный хеш SHA-256. После скачивания файл сверяется с ним: не сошлось — файл удаляется, и загрузку можно повторить. Кнопка «Проверить» на «Моделях» так же сверяет уже установленные модели, а при обновлении программы сверяется и скачанный установщик — чужой файл не запустится.</p>" +
			"<p class=\"wh\">История диктовок</p>" +
			"<p>Раздел «История» в левом столбце хранит то, что вы надиктовали: только текст, только на этом компьютере, звук не сохраняется никогда. По умолчанию выключено — включается одним переключателем там же. Хранится заданное число дней и записей, старое удаляется само; поле «Не записывать из программ» перечисляет через запятую те, из которых не нужно сохранять ничего — менеджеры паролей, банк-клиент. Поиск ищет и по тексту, и по названию программы, кнопка рядом с записью кладёт её в буфер обмена, «Очистить» удаляет всё сразу вместе с файлом <b>{app}-history.json</b>.</p>" +
			"<p class=\"wh\">Голосовые команды</p>" +
			"<p>Под заменами на «Тексте» — список команд: сказанное вслух превращается не в слова, а в действие. «Новая строка» и «новый абзац» ставят перенос — модели их сами не ставят никогда; «отмена» выбрасывает всю диктовку, ничего не вставляя; «подставить текст» кладёт что угодно, хоть смайлик. Кнопка «Добавить стандартные» заполняет список привычными фразами на языке интерфейса. Команды ищутся целыми словами и применяются после замен, поэтому в промпты и в перевод уходит уже готовый текст. Лишние пробелы вокруг переносов убираются сами. Поле внизу проверяет и замены, и команды на любой фразе: перенос показывается значком ⏎.</p>" +
			"<p class=\"wh\">Замены после распознавания</p>" +
			"<p>На «Тексте» можно перечислить, что модель слышит неправильно и во что это превращать: «гит хаб» → GitHub, фамилии, внутренние термины. Замены применяются сразу после распознавания, до промптов, поэтому редактор получает уже правильные слова. Перевод на английский делает сам распознаватель, поэтому замены видят уже переведённый текст. По умолчанию ищутся целые слова и без учёта регистра, обе галочки рядом это меняют. Правила применяются сверху вниз. Поле внизу проверяет их на любой фразе, не диктуя.</p>" +
			"<p class=\"wh\">Правила по приложениям</p>" +
			"<p>На «Диктовке» можно задать правила для отдельных программ: чем вставлять (буфером или посимвольно), нажимать ли Enter, сколько ждать перед вставкой и какие промпты применять. Программа указывается именем файла — <b>chrome.exe</b>; в одном правиле их можно перечислить через запятую, а звёздочка в конце ловит все имена с таким началом. Выигрывает первое подходящее правило; если правил нет или ни одно не подошло, всё работает как в общих настройках. Кнопка рядом со списком подставляет программу, в которую вставляли в последний раз.</p>" +
			"<p class=\"wh\">Основные настройки</p>" +
			"<ul>" +
			"<li><b>Сочетание клавиш</b> — главный хоткей диктовки. Захватывается любая комбинация, различаются левые/правые модификаторы. Хоткеи диктовки, перевода и профилей не должны совпадать — дубликат не даст сохранить настройки.</li>" +
			"<li><b>Язык интерфейса</b> — переключается мгновенно, «как в системе» берёт язык Windows.</li>" +
			"<li><b>Язык распознавания</b> — подсказка Whisper; «авто» определяет язык по речи.</li>" +
			"<li><b>Звук</b> — сигналы начала/конца записи: несколько тем + системные звуки Windows, кнопка ▶ — предпрослушивание.</li>" +
			"<li><b>Enter после вставки</b> — автоматически отправляет надиктованное (удобно в мессенджерах).</li>" +
			"<li><b>Восстанавливать буфер обмена</b> — после вставки возвращает прежнее содержимое буфера целиком, включая картинки, файлы и форматированный текст. Если содержимое сохранить нельзя, буфер не трогается, а текст вводится посимвольно.</li>" +
			"<li><b>Оверлей и анимация</b> — индикатор состояния внизу экрана; анимацию можно отключить.</li>" +
			"<li><b>Посимвольный ввод</b> — вставка имитацией клавиатуры вместо Ctrl+V, для полей, где вставка из буфера не работает.</li>" +
			"</ul>" +
			"<p class=\"wh\">Распознавание</p>" +
			"<div class=\"mock\"><div class=\"mock-row\"><span class=\"mock-radio on\"></span><b>Small</b><span class=\"mock-note\">баланс скорости и точности</span><span style=\"margin-left:auto\">466 MB</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-radio\"></span><b>Medium (q5)</b><span class=\"mock-note\">точнее, рекомендуем</span><span style=\"margin-left:auto\">539 MB ⭳</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-radio\"></span><b>Turbo (q5)</b><span class=\"mock-note\">максимум точности на CPU</span><span style=\"margin-left:auto\">574 MB ⭳</span></div></div>" +
			"<ul>" +
			"<li><b>Два слота</b> — одна модель для русской речи (GigaAM v3: точнее на русском и сама ставит знаки препинания) и одна для остальных языков и перевода. На экране «Состояние» и в таблице маршрутизации над каталогом видно, куда уходит каждый случай; советчик отвечает моделью и причиной выбора.</li>" +
			"<li><b>Модели</b> — каталог моделей Whisper: Base (быстрая, для слабых ПК), Small (баланс), Medium и Turbo (точнее и медленнее; «q5» — квантованная версия: чуть меньше и быстрее почти без потери качества). Радиокнопка выбирает активную (применяется сразу, распознаватель перезапустится); клик по радио нескачанной модели сразу загружает её с официального репозитория whisper.cpp на Hugging Face. Одна модель обслуживает и распознавание, и перевод.</li>" +
			"<li>Модель держится в памяти сервером whisper-server между фразами — поэтому первая диктовка после запуска чуть дольше (загрузка), дальше распознавание занимает 1–3 секунды.</li>" +
			"<li><b>Словарь</b> — термины, имена и аббревиатуры через запятую. Это подсказка «слуху» Whisper, чтобы редкие слова распознавались правильно; это не команды.</li>" +
			"<li><b>Микрофон и скорость</b> — выбор микрофона со шкалой уровня (говорите — полоса двигается, значит устройство слышит), потоки CPU (больше — быстрее распознавание), минимальная длительность записи (отсекает случайные нажатия), максимальная (автостоп записи). Если выбранное устройство отключить, приложение само переключится на системное; запись без речи не отправляется на распознавание — покажет «Тишина».</li>" +
			"<li><b>Сервер</b> — whisper-server запускается автоматически и работает локально. Можно сменить порт или указать URL внешнего сервера — тогда локальный не используется.</li>" +
			"<li><b>Перевод</b> — весь перевод выполняет Whisper: на английский — встроенным режимом перевода, на остальные языки — <b>экспериментально</b>, принудительным языком вывода (качество зависит от пары языков; на крупные языки лучше). Модель Turbo переводу не обучена — при ней настройки покажут предупреждение. «Всегда переводить на целевой язык» — каждая диктовка основным сочетанием переводится на выбранную цель без вопросов. Со снятым чекбоксом работает режим «Спрашивать»: всегда или с таймаутом — перед распознаванием появляется диалог выбора языка, по истечении секунд берётся целевой. Отдельный хоткей перевода переводит разово, не трогая обычную диктовку. Неприменимые при текущем режиме настройки автоматически гаснут серым.</li>" +
			"</ul>" +
			"<p class=\"wh\">Постобработка (LLM)</p>" +
			"<p>Второй, необязательный слой: локальная языковая модель (llama.cpp) редактирует уже распознанный текст по вашим промптам — чистит речь от слов-паразитов, меняет стиль, форматирует. Работает полностью офлайн на CPU.</p>" +
			"<div class=\"mock\"><div class=\"mock-row\" style=\"color:var(--dim)\">▸ Qwen/Qwen2.5-3B-Instruct-GGUF<span style=\"margin-left:auto;color:var(--faint)\">2024-09-20</span><span>↓303k</span><span>↗</span></div>" +
			"<div class=\"mock-row\" style=\"padding-left:18px\">q4_k_m.gguf<span style=\"margin-left:auto\">1.9 GB</span><span style=\"color:var(--green)\">● ≈3.7 GB</span><span style=\"color:var(--dim)\">⭳</span></div>" +
			"<div class=\"mock-row\" style=\"padding-left:18px\">q8_0.gguf<span style=\"margin-left:auto\">3.6 GB</span><span style=\"color:var(--amber)\">● ≈8.2 GB</span><span style=\"color:var(--dim)\">⭳</span></div></div>" +
			"<ul>" +
			"<li><b>Модели</b> — установленные LLM-модели редактора; радиокнопка выбирает активную (применяется сразу), крестик удаляет (можно и активную — тогда постобработка отключится). Здесь же виден прогресс загрузок.</li>" +
			"<li><b>Поиск модели</b> — GGUF-модели на Hugging Face по имени (например, «qwen2.5 instruct»). У репозитория: дата последнего обновления, число загрузок и ссылка ↗ на страницу модели; клик по строке раскрывает файлы-кванты. Индикатор ● ≈N GB сравнивается со <b>свободной</b> оперативной памятью (она показана над списком).</li>" +
			"<li><b>Как выбрать квант:</b> цифра — сколько бит на вес (Q4 — золотая середина, Q8 — почти без сжатия, Q3 — экономия памяти ценой качества); K_M точнее K_S; IQ4 — новое поколение, лучше классических при том же размере. Индикатор ● ≈N GB — оценка нужной оперативной памяти (файл + запас на контекст): зелёный — помещается, жёлтый — впритык, красный — не хватит.</li>" +
			"<li>Модель 1.5–3B — быстрая редактура; 7–9B — заметно умнее, но на CPU каждая обработка занимает секунды. llama-server поднимается при первом использовании и держит модель в памяти наготове.</li>" +
			"</ul>" +
			"<p class=\"wh\">Промпты</p>" +
			"<div class=\"mock\"><div class=\"mock-row\"><span class=\"mock-cb on\">✓</span><b>Чистка речи</b><span class=\"mock-note\">—</span><span style=\"margin-left:auto;color:var(--dim)\">✎ ✕</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-cb\"></span><b>Деловой стиль</b><span class=\"mock-note\">ctrl+alt+f</span><span style=\"margin-left:auto;color:var(--dim)\">✎ ✕</span></div></div>" +
			"<ul>" +
			"<li>Промпт — инструкция для модели-редактора. Из коробки есть пресеты: «Чистка речи» (убирает слова-паразиты, повторы и фальстарты, чинит пунктуацию) и «Деловой стиль» (переписывает вежливо и формально); можно добавлять свои.</li>" +
			"<li>Отмеченные чекбоксами применяются к каждой диктовке по очереди, сверху вниз (цепочкой: результат первого идёт на вход второму); ничего не отмечено — текст вставляется как распознан.</li>" +
			"<li>У промпта может быть личный хоткей: диктовка этим хоткеем применяет только его, разово. Карандаш ✎ открывает редактор: имя, текст промпта, хоткей и поле проверки ▶ — прогон примера через живую модель прямо из настроек.</li>" +
			"<li>Совет: маленьким моделям помогают примеры «вход → выход» прямо в тексте промпта — все пресеты так и написаны.</li>" +
			"<li>Если профиль не сработал (модель не ответила), текст вставляется без него: оверлей покажет «Вставлено без профиля …», а Enter после вставки в этом случае не нажимается.</li>" +
			"</ul>" +
			"<p class=\"wh\">Зависимости</p>" +
			"<ul>" +
			"<li>Промпты работают только при установленной модели-редакторе; перевод от неё не зависит — его целиком выполняет Whisper.</li>" +
			"<li>Модель-редактор загружается в память при первом использовании и остаётся готовой; крупные модели заметно медленнее на процессоре.</li>" +
			"<li>Проверьте индикатор памяти перед скачиванием: модель «впритык» может замедлить всю систему.</li>" +
			"<li>Серые (неактивные) элементы интерфейса показывают настройки, которые не участвуют при текущем режиме.</li>" +
			"</ul>" +
			"<p class=\"wh\">Установка и переносимость</p>" +
			"<ul>" +
			"<li><b>{setup}</b> — инсталлер: ставит без прав администратора, ярлык в Пуске, опциональный автозапуск, корректное удаление через «Установку и удаление программ».</li>" +
			"<li>Установщик ничего не качает по умолчанию: модель подберёт и скачает мастер при первом запуске. Если модель всё же выбрана — для русского это GigaAM v3, для остальных языков Whisper, — загрузку можно остановить кнопкой, и установка всё равно дойдёт до конца. Там же галочка «Проверять обновления»: ответ записывается в настройки приложения.</li>" +
			"<li><b>Портативный вариант</b> — просто скопируйте папку с exe целиком (на флешку, другой ПК): настройки, модели и лог живут рядом с exe и переезжают вместе с ним. В реестр ничего не пишется.</li>" +
			"<li>При первом запуске без модели распознавания приложение само откроет каталог моделей и дождётся скачивания.</li>" +
			"<li>Требования: Windows 10/11 x64, CPU с AVX2 (~2013+), WebView2 Runtime для окна настроек (в Windows 11 есть).</li>" +
			"</ul>" +
			"<p class=\"wh\">Трей и файлы</p>" +
			"<div class=\"mock\" style=\"max-width:290px\"><div class=\"mock-mi dim\">{app} — Готов…</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Настройки…</div><div class=\"mock-mi\">Выключить</div><div class=\"mock-mi\">Копировать последний результат</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Открыть config.json</div><div class=\"mock-mi\">Открыть лог</div><div class=\"mock-mi\">О приложении</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Выход</div></div>" +
			"<ul>" +
			"<li>Левый клик по иконке в трее — настройки; правый — меню. Цвет иконки: зелёный — готов, красный — запись, оранжевый — распознавание, серый — выключено или ошибка.</li>" +
			"<li><b>config.json</b> — все настройки; правки вручную применяются кнопкой «Перечитать» в разделе «Система». Там же — «Открыть лог» и «Сбросить настройки»: сброс возвращает всё к заводскому виду, но модели, историю и промпты не трогает.</li>" +
			"<li><b>{log}</b> — журнал работы, автоматически ограничен ~2 МБ.</li>" +
			"<li><b>models/</b> — скачанные модели Whisper и LLM.</li>" +
			"</ul>",
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
		"S_TR_HINT":        "Translation is done by Whisper: to English — native mode, to other languages — experimentally, by forcing the output language (quality depends on the language pair).",
		"S_TR_TURBO":       "⚠ The active Turbo model is not trained for translation to English — pick another model in the Models section for translating.",
		"S_TR_DEFAULT":     "Always translate to the target language",
		"S_TR_TARGET":      "Target language",
		"S_TR_ASK":         "Language choice",
		"S_TR_ASK_NEVER":   "Don't ask",
		"S_TR_ASK_ALWAYS":  "Ask every time",
		"S_TR_ASK_TIMEOUT": "Ask with a timeout",
		"S_TR_SECONDS":     "Timeout, sec",
		"S_TR_LANGS":       "Languages in the dialog",
		"S_DICT_HINT":      "Comma-separated terms, names and abbreviations — a hint for the ear, not commands. It works for Whisper; Russian speech through GigaAM ignores it. The default set follows the recognition language until you write your own.",
		"S_LLM_HINT":       "Checked profiles apply one after another, top to bottom, on regular dictation. Nothing checked — text is inserted as is. A profile's hotkey applies just that profile once.",
		"S_UPDATED":        "Model last updated",
		"S_FIT_OK":         "fits",
		"S_FIT_WARN":       "tight",
		"S_FIT_BAD":        "not enough RAM",
		"S_RAM":            "Computer RAM:",
		"S_HF_PH":          "Model name — e.g. qwen2.5 instruct",
		"S_NO_LLM":         "No models installed yet — find and download one in the search field below.",
		"S_NO_LLM_PROF":    "Prompts become available once a model is installed — the Models section, the editing-model list at the bottom.",
		"S_PROF_EDIT":      "Edit",
		"S_PROF_CLOSE":     "Collapse",
		"S_CONFIRM_DEL":    "Delete the \"%s\" model? It can be downloaded again.",
		"S_DEL_ACTIVE":     "Delete the active model \"%s\"? Recognition stops until you pick another one — you can download it right here.",
		"S_WIZ_NEED_MODEL": "Download a model first — without one there is nothing to recognise with",
		"S_FREE":           "free",
		"S_SUB_PROMPTS":    "Prompts",
		"S_SUB_DICT":       "Dictionary",
		"S_PROF_ADD":       "Add profile",
		"S_PROF_NAME":      "Name",
		"S_PROF_PROMPT":    "Prompt",
		"S_PROF_HOTKEY":    "Hotkey",
		"S_PROF_SET":       "Set…",
		"S_PROF_CLEAR":     "Clear",
		"S_PROF_TEST":      "Test",
		"S_HOTKEY":         "Keyboard shortcut",
		"S_CHANGE":         "Change…",
		"S_UILANG":         "UI language",
		"S_AUTO":           "System default",
		"S_BEEP":           "Recording sound cues",
		"S_SEC_SOUND":      "Sound",
		"S_SEC_BEHAVIOR":   "Behavior",
		"S_SOUND":          "Cue sound",
		"S_SND_SPEECH":     "Windows voice",
		"S_SND_CHIME":      "Chime",
		"S_SND_SOFT":       "Soft",
		"S_SND_MARIMBA":    "Marimba",
		"S_SND_BLIP":       "Blip",
		"S_SND_POP":        "Pop",
		"S_AUTOENTER":      "Press Enter after paste (auto-submit)",
		"S_RESTORE":        "Restore clipboard after paste",
		"S_OVERLAY":        "On-screen indicator",
		"S_NAV_HISTORY":    "History", "S_HIST_ON": "Keep a history of dictations", "S_HIST_ON_SUB": "text only, on this computer; audio is never kept",
		"S_HIST_DAYS": "How many days to keep", "S_HIST_MAX": "How many entries to keep",
		"S_HIST_SKIP": "Never record from these programs", "S_HIST_SKIP_SUB": "comma-separated: keepass.exe, 1password.exe",
		"S_HIST_LIST": "Entries", "S_HIST_CLEAR": "Clear", "S_HIST_COPY": "Copy",
		"S_HIST_FIND": "Search the history…", "S_HIST_EMPTY": "No history yet", "S_HIST_ASK": "Delete the whole dictation history?",
		"S_SEC_CMD": "Voice commands", "S_CMD_HINT": "What you say turns into a line break, a symbol or a cancel instead of landing in the text. Matched as whole words, applied top to bottom, after the replacements.",
		"S_CMD_ADD": "Add a command", "S_CMD_PRESET": "Add the usual ones", "S_CMD_PH": "new line",
		"S_CMD_NEWLINE": "line break", "S_CMD_PARAGRAPH": "new paragraph", "S_CMD_TEXT": "insert text", "S_CMD_CANCEL": "cancel the dictation",
		"S_CMD_TEXT_PH": "what to insert", "S_CMD_EMPTY": "No commands yet", "S_CMD_DEL": "Delete the command",
		"S_CMD_P_NEWLINE": "new line", "S_CMD_P_PARAGRAPH": "new paragraph", "S_CMD_P_CANCEL": "cancel",
		"S_SEC_REPLACE": "Replacements after recognition", "S_REPLACE_HINT": "What was misheard becomes what you meant — right after recognition, before the prompts. Applied from top to bottom.",
		"S_REPL_ADD": "Add a replacement", "S_REPL_FROM_PH": "git hub", "S_REPL_TO_PH": "GitHub",
		"S_REPL_WHOLE": "whole words", "S_REPL_CASE": "case", "S_REPL_EMPTY": "No replacements yet",
		"S_REPL_DEL": "Delete the replacement", "S_REPL_TEST_PH": "type a phrase to try the replacements and commands",
		"S_SEC_RULES": "Rules per application", "S_RULES_HINT": "Insertion can work differently for particular programs. The first matching rule wins.",
		"S_RULE_ADD": "Add a rule", "S_RULE_PH": "chrome.exe, msedge.exe",
		"S_RULE_PASTE_INH": "insertion: as set", "S_RULE_ENTER_INH": "Enter: as set", "S_RULE_DELAY_NONE": "no delay", "S_RULE_PROMPT_INH": "prompts: as set",
		"S_RULE_CLIP": "clipboard", "S_RULE_TYPE": "character by character", "S_RULE_ENTER_ON": "with Enter", "S_RULE_ENTER_OFF": "without Enter",
		"S_RULE_NOPROMPT": "no prompts", "S_RULE_LAST": "last insertion: %s", "S_RULE_EMPTY": "No rules yet",
		"S_RULE_DEL": "Delete the rule", "S_RULE_PROMPTS": "Prompts",
		"S_PASTE_DELAY": "Delay before inserting", "S_PASTE_DELAY_SUB": "when the program is not ready for the text yet",
		"S_OVPOS": "Where to show the plate", "S_OVPOS_SUB": "at the cursor — next to where you type; if the app hides it, next to the mouse pointer",
		"S_OVPOS_BOTTOM": "Bottom of the screen", "S_OVPOS_TOP": "Top of the screen", "S_OVPOS_CARET": "At the cursor",
		"S_OVTEXT": "Show the recognised text", "S_OVTEXT_SUB": "on the plate after insertion, instead of the character count",
		"S_ANIM":           "Recording & transcribing animation",
		"S_TYPEMODE":       "Insert character-by-character instead of the clipboard",
		"S_RECLANG":        "Recognition language",
		"S_RECAUTO":        "Auto-detect",
		"S_DL":             "Download",
		"S_DEL":            "Delete",
		"S_M_BASE":         "fast, for weak PCs",
		"S_M_SMALL":        "balanced speed and accuracy",
		"S_M_MED":          "more accurate, recommended",
		"S_M_TURBO":        "best accuracy on CPU",
		"S_M_CUSTOM":       "custom (from config.json)",
		"S_M_GIGAAM":       "more accurate in Russian, punctuates by itself",
		"S_ADV_TITLE":      "Pick a model",
		"S_ADV_LANGQ":      "Which language do you dictate",
		"S_ADV_PRIOQ":      "What matters more",
		"S_ADV_TRQ":        "Translation needed",
		"S_ADV_GO":         "Recommend",
		"S_ADV_LANG":       "It is tuned for your language.",
		"S_ADV_ACC":        "Picked for accuracy.",
		"S_ADV_SPEED":      "Picked for speed.",
		"S_ADV_RAM":        "Larger ones would not fit the free memory.",
		"S_ADV_NONE":       "Not enough free memory for any model — close a few programs.",
		"S_F_ALL":          "all",
		"S_F_RU":           "Russian",
		"S_F_MULTI":        "many languages",
		"S_F_PUNCT":        "punctuates",
		"S_F_FIT":          "fits in memory",
		"S_PUNCT":          "Punctuation and capitals",
		"S_PUNCT_MODEL":    "from the model",
		"S_PUNCT_LLM":      "by the editor model",
		"S_PUNCT_OFF":      "strip",
		"S_SEARCH":         "Find a setting…",
		"S_LVL_SIMPLE":     "simple",
		"S_LVL_ALL":        "all",
		"S_MORE":           "%d more settings",
		"S_LESS":           "Collapse %d settings",
		"S_HOTMODE":        "Mode",
		"S_HOTMODE_HOLD":   "hold",
		"S_HOTMODE_TOGGLE": "toggle",
		"S_SUB_HOTMODE":    "hold the keys, or press once to start and once to stop",
		"S_GRP_GENERAL":    "General",
		"S_GRP_SPEECH":     "Speech",
		"S_GRP_TEXT":       "Text",
		"S_GRP_INFO":       "Info",
		"S_NAV_POST":       "Post-processing",
		"S_NAV_HELP":       "Help",
		"S_NAV_CONTACTS":   "Contacts",
		"S_NAV_STATE":      "Overview",
		"S_NAV_DICT":       "Dictation",
		"S_NAV_MIC":        "Microphone",
		"S_NAV_MODELS":     "Models",
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
		"S_SUB_TRTARGET":   "English is native for Whisper, other targets are experimental",
		"S_TR_EXP":         "for targets other than English the app forces the output language instead of translating — the text may stay in the spoken language",
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
		"S_STATE_RU":    "Russian speech",
		"S_STATE_OTHER": "Other languages",
		"S_MIC":         "Microphone",
		"S_MIC_DEFAULT": "System default",
		"S_MIC_CHECK":   "Check the microphone", "S_MIC_CHECK_SUB": "three seconds of recording, then a verdict: level, clipping, whether there is speech", "S_MIC_CHECKING": "Checking…",
		"S_PAUSE": "Pause the recording", "S_PAUSE_SUB": "works in toggle mode: press once and recording freezes, press again and it goes on",
		"S_MCHECK": "Check installed models", "S_MCHECK_SUB": "compares model files with reference hashes", "S_MCHECK_GO": "Check", "S_MCHECK_RUN": "Checking…",
		"S_HIST_INSERT": "Paste",
		"S_LISTS_HINT":  "Replacements and commands in one file — to move to another computer", "S_LISTS_EXPORT": "Save to a file", "S_LISTS_IMPORT": "Load from a file",
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
		"S_ADV_PRIMARY":     "main",
		"S_ADV_COMPANION":   "second",
		"S_ADV_HAVE":        "already here",
		"S_ADV_APPLY":       "Apply",
		"S_ADV_ASK":         "These will be downloaded: %s — %s in total. Start?",
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
		"S_ADV_ROLE_MAIN":   "for the language you picked",
		"S_ADV_ROLE_SECOND": "for other languages and translation",
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
			"<p class=\"wh\">Pausing the recording</p>" +
			"<p>In toggle mode (press once to start, press again to stop) a separate shortcut can pause: the Dictation section, the \"Pause the recording\" row. Press it and the recording freezes — the plate says Paused and nothing is recorded; press again and it goes on, keeping everything said before the pause. The recording length limit does not fire while paused.</p>" +
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
			"<div class=\"mock\"><div class=\"mock-row\"><span class=\"mock-radio on\"></span><b>Small</b><span class=\"mock-note\">speed/accuracy balance</span><span style=\"margin-left:auto\">466 MB</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-radio\"></span><b>Medium (q5)</b><span class=\"mock-note\">more accurate, recommended</span><span style=\"margin-left:auto\">539 MB ⭳</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-radio\"></span><b>Turbo (q5)</b><span class=\"mock-note\">best CPU accuracy</span><span style=\"margin-left:auto\">574 MB ⭳</span></div></div>" +
			"<ul>" +
			"<li><b>Two slots</b> — one model for Russian speech (GigaAM v3: more accurate on Russian and it punctuates by itself) and one for every other language and for translation. The Status screen and the routing table above the catalog show which model each case goes to; the picker answers with a model and the reason for it.</li>" +
			"<li><b>Models</b> — the Whisper catalog: Base (fast, weak PCs), Small (balanced), Medium and Turbo (more accurate, slower; \"q5\" means a quantized build — slightly smaller and faster with almost no quality loss). The radio selects the active one (applied at once; the recognizer restarts); clicking the radio of a missing model downloads it from the official whisper.cpp repository on Hugging Face. One model serves both recognition and translation.</li>" +
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
			"<li>A prompt can have its own hotkey: dictating with it applies only that prompt, once. The ✎ pencil opens the editor: name, prompt text, hotkey and a ▶ test field that runs a sample through the live model right from Settings.</li>" +
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
