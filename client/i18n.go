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
		"status.restart.needed": "Модель/язык/сервер изменятся после перезапуска приложения",
		"menu.settings":         "Настройки…",
		"menu.enable":           "Включить",
		"menu.disable":          "Выключить",
		"menu.reload":           "Перечитать config.json",
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
		"ov.speak":              "Говорите…",
		"ov.transcribing":       "Распознаю",
		"ov.inserted":           "Вставлено: %d симв.",
		"ov.err.mic":            "Микрофон недоступен — проверьте устройство в настройках",
		"ov.err.recognize":      "Ошибка распознавания (см. лог)",
		"ov.err.paste":          "Ошибка вставки (см. лог)",
		"ov.silence":            "Тишина — текст не распознан",
		"ov.notranslate":        "Активная модель не переводит — вставлен исходный текст",
		"ov.engine.fallback":    "Второй движок не поднялся — распознаю текущим",
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
		"err.sherpa.notfound":   "Распознаватель sherpa не найден: %s",
		"err.sherpa.start":      "sherpa-server завершился при старте (см. лог)",
		"err.sherpa.translate":  "эта модель не умеет переводить",
		"err.sherpa.model":      "Файл модели не найден: %s — скачайте её в настройках или исправьте sherpa_model в config.json",
		"ov.server.loading":     "Сервер ещё загружается",
		"ov.cancelled":          "Отменено",
		"ov.editing":            "Редактирую: %s",
		"ov.translating":        "Перевожу",
		"ov.llm.needed":         "Для этого языка нужен LLM-модуль",
		"td.title":              "Переводить на:",
		"cap.title":             "{app} — сочетание клавиш",
		"cap.prompt":            "Нажмите новое сочетание клавиш…\n\n(сейчас: %s)\n\nEsc — отмена",
		"cap.selected":          "Выбрано: %s",
		"cap.cancelled":         "Отменено",
		"err.hotkey.dup":        "Сочетание «%s» назначено дважды — хоткеи не должны совпадать",
		"cfg.err.recovered":     "config.json повреждён (%s).\nФайл сохранён как %s, настройки сброшены к значениям по умолчанию.",
		"err.disk.space":        "мало места на диске: свободно %d МБ, нужно ~%d МБ",
		"un.title":              "{app} — удаление",
		"un.confirm":            "Удалить {app} с этого компьютера?",
		"un.data":               "Удалить также настройки и скачанные модели?",
		"un.done":               "{app} удалён.",
		"model.switching":       "Переключаю модель — распознаватель перезапускается…",
		"model.del.active":      "Нельзя удалить активную модель",
		"model.del.ok":          "Модель удалена",
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
		"status.restart.needed": "Model/language/server changes apply after app restart",
		"menu.settings":         "Settings…",
		"menu.enable":           "Enable",
		"menu.disable":          "Disable",
		"menu.reload":           "Reload config.json",
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
		"ov.speak":              "Speak…",
		"ov.transcribing":       "Transcribing",
		"ov.inserted":           "Inserted: %d chars",
		"ov.err.mic":            "Microphone unavailable — check the device in Settings",
		"ov.err.recognize":      "Recognition error (see log)",
		"ov.err.paste":          "Paste error (see log)",
		"ov.silence":            "Silence — nothing recognized",
		"ov.notranslate":        "The active model cannot translate — inserted as recognized",
		"ov.engine.fallback":    "The other engine did not start — using the current one",
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
		"err.sherpa.notfound":   "sherpa recognizer not found: %s",
		"err.sherpa.start":      "sherpa-server exited during startup (see log)",
		"err.sherpa.translate":  "this model cannot translate",
		"err.sherpa.model":      "Model file not found: %s — download it in Settings or fix sherpa_model in config.json",
		"ov.server.loading":     "Server is still loading",
		"ov.cancelled":          "Cancelled",
		"ov.editing":            "Editing: %s",
		"ov.translating":        "Translating",
		"ov.llm.needed":         "This language requires the LLM module",
		"td.title":              "Translate to:",
		"cap.title":             "{app} — shortcut",
		"cap.prompt":            "Press a new key combination…\n\n(current: %s)\n\nEsc — cancel",
		"cap.selected":          "Selected: %s",
		"cap.cancelled":         "Cancelled",
		"err.hotkey.dup":        "The \"%s\" shortcut is assigned twice — hotkeys must be unique",
		"cfg.err.recovered":     "config.json is corrupted (%s).\nThe file was saved as %s and settings were reset to defaults.",
		"err.disk.space":        "low disk space: %d MB free, ~%d MB needed",
		"un.title":              "{app} — Uninstall",
		"un.confirm":            "Remove {app} from this computer?",
		"un.data":               "Also delete settings and downloaded models?",
		"un.done":               "{app} has been removed.",
		"model.switching":       "Switching model — the recognizer is restarting…",
		"model.del.active":      "Cannot delete the active model",
		"model.del.ok":          "Model deleted",
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
		"S_TAB_GENERAL":    "Основные",
		"S_TAB_REC":        "Распознавание",
		"S_TAB_PROC":       "Постобработка",
		"S_TR":             "Перевод",
		"S_TR_HINT":        "Перевод выполняет Whisper: на английский — штатным режимом, на остальные языки — экспериментально, принудительным языком вывода (качество зависит от пары языков).",
		"S_TR_TURBO":       "⚠ Активная модель Turbo не обучена переводу на английский — для перевода выберите другую модель на вкладке «Модели».",
		"S_TR_DEFAULT":     "Всегда переводить на целевой язык",
		"S_TR_TARGET":      "Целевой язык",
		"S_TR_ASK":         "Выбор языка",
		"S_TR_ASK_NEVER":   "Не спрашивать (язык по умолчанию)",
		"S_TR_ASK_ALWAYS":  "Спрашивать каждый раз",
		"S_TR_ASK_TIMEOUT": "Спрашивать с таймаутом",
		"S_TR_SECONDS":     "Таймаут, сек",
		"S_TR_LANGS":       "Языки в диалоге",
		"S_TAB_SERVER":     "Сервер",
		"S_PIPE":           "голос ▸ распознавание ▸ редактура ▸ вставка",
		"S_DICT":           "Словарь распознавания",
		"S_DICT_HINT":      "Термины, имена и аббревиатуры через запятую — подсказка слуху, не команды.",
		"S_LLM":            "Профили обработки",
		"S_LLM_HINT":       "Отмеченные профили применяются по очереди, сверху вниз, при обычной диктовке. Ничего не отмечено — текст вставляется как есть. Хоткей профиля применяет разово только его.",
		"S_LLM_MODEL":      "Модель-редактор (LLM)",
		"S_UPDATED":        "Дата последнего обновления модели",
		"S_FIT_OK":         "поместится",
		"S_FIT_WARN":       "впритык",
		"S_FIT_BAD":        "не хватит памяти",
		"S_RAM":            "Память компьютера:",
		"S_HF_PH":          "Название модели — например, qwen2.5 instruct",
		"S_NO_LLM":         "Пока не установлено ни одной модели — найдите и скачайте на вкладке «Поиск».",
		"S_NO_LLM_PROF":    "Промпты станут доступны после установки модели (вкладки «Модели» и «Поиск»).",
		"S_PROF_EDIT":      "Редактировать",
		"S_PROF_CLOSE":     "Свернуть",
		"S_CONFIRM_DEL":    "Удалить модель «%s»? Её можно будет скачать заново.",
		"S_FREE":           "свободно",
		"S_SUB_MODELS":     "Модели",
		"S_SUB_SEARCH":     "Поиск",
		"S_SUB_PROMPTS":    "Промпты",
		"S_SUB_PARAMS":     "Параметры",
		"S_SUB_DICT":       "Словарь",
		"S_SUB_TR":         "Перевод",
		"S_PROF_ASIS":      "Как есть",
		"S_PROF_WT":        "Перевод → English (быстрый)",
		"S_PROF_ADD":       "Добавить профиль",
		"S_PROF_NAME":      "Имя",
		"S_PROF_PROMPT":    "Промпт",
		"S_PROF_HOTKEY":    "Хоткей",
		"S_PROF_SET":       "Задать…",
		"S_PROF_CLEAR":     "Сброс",
		"S_PROF_TEST":      "Проверка",
		"S_TAB_ABOUT":      "О программе",
		"S_HOTKEY":         "Сочетание клавиш",
		"S_CHANGE":         "Изменить…",
		"S_UILANG":         "Язык интерфейса",
		"S_AUTO":           "Как в системе",
		"S_BEEP":           "Звуковые сигналы записи",
		"S_SEC_SOUND":      "Звук",
		"S_SEC_BEHAVIOR":   "Поведение",
		"S_SOUND":          "Звук сигналов",
		"S_SND_SPEECH":     "Системный (речь)",
		"S_SND_CHIME":      "Колокольчик",
		"S_SND_SOFT":       "Мягкий",
		"S_SND_MARIMBA":    "Маримба",
		"S_SND_BLIP":       "Блип",
		"S_SND_POP":        "Поп",
		"S_AUTOENTER":      "Нажимать Enter после вставки (автоотправка)",
		"S_RESTORE":        "Восстанавливать буфер обмена после вставки",
		"S_OVERLAY":        "Индикатор внизу экрана",
		"S_ANIM":           "Анимация записи и распознавания",
		"S_TYPEMODE":       "Посимвольный ввод (для полей, где запрещена вставка)",
		"S_RECLANG":        "Язык распознавания",
		"S_RECAUTO":        "Автоопределение",
		"S_MODELS":         "Модели распознавания",
		"S_DL":             "Скачать",
		"S_USE":            "Включить",
		"S_DEL":            "Удалить",
		"S_ACTIVE":         "Активна",
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
		"S_HIDDEN":         "Простой режим · скрыто настроек: %d",
		"S_ALLSHOWN":       "Показаны все настройки",
		"S_SHOWALL":        "Показать всё",
		"S_SHOWSIMPLE":     "Вернуть простой режим",
		"S_HOTMODE":        "Режим",
		"S_HOTMODE_HOLD":   "удержание",
		"S_HOTMODE_TOGGLE": "фиксация",
		"S_SUB_HOTMODE":    "удерживать клавиши или включать и выключать нажатием",
		"S_GRP_WORK":       "Работа",
		"S_GRP_REC":        "Распознавание",
		"S_GRP_OTHER":      "Прочее",
		"S_NAV_STATE":      "Состояние",
		"S_NAV_DICT":       "Диктовка",
		"S_NAV_MIC":        "Микрофон",
		"S_NAV_MODELS":     "Модели",
		"S_NAV_TEXT":       "Текст",
		"S_NAV_TR":         "Перевод",
		"S_NAV_SYSTEM":     "Система",
		"S_NAV_ABOUT":      "О программе",
		"S_STATE_HINT":     "зажмите и говорите — текст появится там, где курсор",
		"S_STATE_ENGINE":   "Распознавание",
		"S_STATE_PROC":     "Обработка",
		"S_STATE_MEM":      "Память",
		"S_STATE_MEM_SUB":  "модели держатся загруженными, первая фраза без задержки",
		"S_SUB_MINMS":      "отсекает случайные нажатия",
		"S_SUB_ENTER":      "отправляет сообщение сразу",
		"S_SUB_CLIP":       "картинки и файлы возвращаются как были",
		"S_SUB_TYPE":       "для полей, где вставка из буфера запрещена",
		"S_SUB_THREADS":    "больше потоков — не всегда быстрее, проверьте на своей машине",
		"S_SUB_PUNCT":      "откуда берутся знаки препинания и заглавные",
		"S_SUB_TRTARGET":   "на английский переводит Whisper, остальные — экспериментально",
		"S_SUB_AUTOSTART":  "выключите, если сервер поднимаете сами",
		"S_SUB_PORT":       "применится после перезапуска",
		"S_SUB_UPD":        "единственный сетевой запрос, кроме загрузки моделей",
		"S_STATE_LAST":     "Последняя диктовка",
		"S_STATE_COPY":     "Копировать",
		"S_SEC_OVERLAY":    "Плашка на экране",
		"S_SEC_SERVICE":    "Служебное",
		"S_SEC_LLM":        "Модель-редактор",
		"S_NOT_INSTALLED":  "не установлена",
		"S_CHANGE_MODEL":   "Сменить",
		"S_PICK_MODEL":     "Подобрать",
		"S_STATE_RU":       "Русская речь",
		"S_STATE_OTHER":    "Прочие языки",
		"S_ENG_WHISPER":    "whisper.cpp · 99 языков · перевод",
		"S_ENG_SHERPA":     "sherpa-onnx · только русский",
		"S_MODEL":          "Файл модели",
		"S_MIC":            "Микрофон",
		"S_MIC_DEFAULT":    "Системный по умолчанию",
		"S_MIC_REFRESH":    "Обновить список",
		"S_MIC_LEVEL":      "Уровень сигнала",
		"S_MIC_QUIET":      "тихо",
		"S_THREADS":        "Потоки CPU",
		"S_MINMS":          "Минимальная запись, мс",
		"S_MAXSEC":         "Максимальная запись, сек",
		"S_AUTOSTART":      "Запускать whisper-server автоматически",
		"S_PORT":           "Порт",
		"S_SERVEREXE":      "Путь к whisper-server",
		"S_SERVERURL":      "Внешний сервер (URL)",
		"S_URLHINT":        "Если задан — свой сервер не запускается",
		"S_SAVE":           "Сохранить",
		"S_SAVED":          "Сохранено",
		"S_RESTART":        "Сохранено. Модель/язык/сервер применятся после перезапуска",
		"S_SUB_INFO":       "Информация",
		"S_SUB_HELP":       "Справка",
		"S_UPD":            "Обновления",
		"S_UPD_CHECK":      "Проверить обновления",
		"S_UPD_AUTO":       "Проверять при запуске",
		"S_UPD_NONE":       "Установлена последняя версия",
		"S_UPD_AVAIL":      "Доступна версия %s.",
		"S_UPD_GO":         "Обновить",
		"S_UPD_ERR":        "Не удалось проверить обновления",
		"S_UPD_DL":         "Скачиваю обновление…",
		"S_SUB_AUTHOR":     "Об авторе",
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
			"<p>Полный конвейер: <b>запись → распознавание (Whisper) → перевод (если включён) → промпты (LLM) → вставка</b>. Каждая стадия видна на оверлее.</p>" +
			"<p class=\"wh\">Оверлей</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\"></span><span>Говорите…</span><span class=\"mock-bars\"><i style=\"height:6px\"></i><i style=\"height:13px\"></i><i style=\"height:9px\"></i><i style=\"height:16px\"></i><i style=\"height:7px\"></i><i style=\"height:12px\"></i><i style=\"height:5px\"></i></span><span class=\"mock-x\">✕</span></div></div>" +
			"<ul>" +
			"<li><b>Говорите…</b> — идёт запись: красная точка и живые столбики уровня голоса.</li>" +
			"<li><b>Распознаю…</b> — Whisper обрабатывает звук; при переводе — «Перевожу», при промптах — «Редактирую: имя (1/2)».</li>" +
			"<li><b>Вставлено: N симв.</b> — готово; при ошибке или тишине — короткое сообщение о причине.</li>" +
			"<li>Крестик ✕ справа отменяет операцию на любой стадии; фокус ввода оверлей не забирает. Показывать плашку и её анимацию можно отключить на «Основных».</li>" +
			"</ul>" +
			"<p class=\"wh\">Диалог выбора языка</p>" +
			"<div class=\"mock\"><div style=\"display:flex;justify-content:space-between;margin-bottom:8px\"><b>Переводить на: (3)</b><span class=\"mock-x\">✕</span></div><div style=\"display:flex;gap:8px\"><span class=\"mock-btn on\">EN</span><span class=\"mock-btn\">DE</span><span class=\"mock-btn\">ES</span></div></div>" +
			"<p>Появляется над оверлеем после отпускания хоткея в режимах «Спрашивать всегда/с таймаутом». Набор кнопок — «Языки в диалоге», подсвечен целевой язык. В режиме с таймаутом в заголовке идёт отсчёт секунд, по истечении применяется целевой язык. Крестик диалога — вставить без перевода; крестик оверлея отменяет операцию целиком.</p>" +
			"<p class=\"wh\">Безопасная вставка</p>" +
			"<div class=\"mock\"><div style=\"display:flex;justify-content:space-between;margin-bottom:8px\"><b style=\"color:var(--amber)\">Фокус сменился — вставить? (30)</b><span class=\"mock-x\">✕</span></div><div style=\"display:flex;gap:8px\"><span class=\"mock-btn on\">Вставить сюда</span><span class=\"mock-btn\">Копировать</span></div></div>" +
			"<ul>" +
			"<li>Окно-цель запоминается в момент нажатия хоткея. Если за время обработки фокус сменился, текст не вставляется — появляется диалог: <b>Вставить сюда</b> (в текущее окно), <b>Копировать</b> (в буфер обмена) или ✕. По истечении отсчёта вставка отменяется, текст остаётся в «Последнем результате».</li>" +
			"<li>Enter после вставки нажимается только если окно-цель не менялось.</li>" +
			"<li><b>Последний результат</b> — финальный текст каждой диктовки хранится в памяти до следующей; в меню трея есть пункт «Копировать последний результат». Ошибка вставки или смена фокуса не теряют надиктованное.</li>" +
			"</ul>" +
			"<p class=\"wh\">Основные</p>" +
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
			"<p class=\"wh\">Распознавание (Whisper)</p>" +
			"<div class=\"mock\"><div class=\"mock-row\"><span class=\"mock-radio on\"></span><b>Small</b><span class=\"mock-note\">баланс скорости и точности</span><span style=\"margin-left:auto\">466 MB</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-radio\"></span><b>Medium (q5)</b><span class=\"mock-note\">точнее, рекомендуем</span><span style=\"margin-left:auto\">539 MB ⭳</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-radio\"></span><b>Turbo (q5)</b><span class=\"mock-note\">максимум точности на CPU</span><span style=\"margin-left:auto\">574 MB ⭳</span></div></div>" +
			"<ul>" +
			"<li><b>Модели</b> — каталог моделей Whisper: Base (быстрая, для слабых ПК), Small (баланс), Medium и Turbo (точнее и медленнее; «q5» — квантованная версия: чуть меньше и быстрее почти без потери качества). Радиокнопка выбирает активную (применяется по «Сохранить», распознаватель перезапустится); клик по радио нескачанной модели сразу загружает её с официального репозитория whisper.cpp на Hugging Face. Одна модель обслуживает и распознавание, и перевод.</li>" +
			"<li>Модель держится в памяти сервером whisper-server между фразами — поэтому первая диктовка после запуска чуть дольше (загрузка), дальше распознавание занимает 1–3 секунды.</li>" +
			"<li><b>Словарь</b> — термины, имена и аббревиатуры через запятую. Это подсказка «слуху» Whisper, чтобы редкие слова распознавались правильно; это не команды.</li>" +
			"<li><b>Параметры</b> — выбор микрофона со шкалой уровня (говорите — полоса двигается, значит устройство слышит), потоки CPU (больше — быстрее распознавание), минимальная длительность записи (отсекает случайные нажатия), максимальная (автостоп записи). Если выбранное устройство отключить, приложение само переключится на системное; запись без речи не отправляется на распознавание — покажет «Тишина».</li>" +
			"<li><b>Сервер</b> — whisper-server запускается автоматически и работает локально. Можно сменить порт или указать URL внешнего сервера — тогда локальный не используется.</li>" +
			"<li><b>Перевод</b> — весь перевод выполняет Whisper: на английский — встроенным режимом перевода, на остальные языки — <b>экспериментально</b>, принудительным языком вывода (качество зависит от пары языков; на крупные языки лучше). Модель Turbo переводу не обучена — при ней настройки покажут предупреждение. «Всегда переводить на целевой язык» — каждая диктовка основным сочетанием переводится на выбранную цель без вопросов. Со снятым чекбоксом работает режим «Спрашивать»: всегда или с таймаутом — перед распознаванием появляется диалог выбора языка, по истечении секунд берётся целевой. Отдельный хоткей перевода переводит разово, не трогая обычную диктовку. Неприменимые при текущем режиме настройки автоматически гаснут серым.</li>" +
			"</ul>" +
			"<p class=\"wh\">Постобработка (LLM)</p>" +
			"<p>Второй, необязательный слой: локальная языковая модель (llama.cpp) редактирует уже распознанный текст по вашим промптам — чистит речь от слов-паразитов, меняет стиль, форматирует. Работает полностью офлайн на CPU.</p>" +
			"<div class=\"mock\"><div class=\"mock-row\" style=\"color:var(--dim)\">▸ Qwen/Qwen2.5-3B-Instruct-GGUF<span style=\"margin-left:auto;color:var(--faint)\">2024-09-20</span><span>↓303k</span><span>↗</span></div>" +
			"<div class=\"mock-row\" style=\"padding-left:18px\">q4_k_m.gguf<span style=\"margin-left:auto\">1.9 GB</span><span style=\"color:var(--green)\">● ≈3.7 GB</span><span style=\"color:var(--dim)\">⭳</span></div>" +
			"<div class=\"mock-row\" style=\"padding-left:18px\">q8_0.gguf<span style=\"margin-left:auto\">3.6 GB</span><span style=\"color:var(--amber)\">● ≈8.2 GB</span><span style=\"color:var(--dim)\">⭳</span></div></div>" +
			"<ul>" +
			"<li><b>Модели</b> — установленные LLM-модели редактора; радиокнопка выбирает активную (по «Сохранить»), крестик удаляет (можно и активную — тогда постобработка отключится). Здесь же виден прогресс загрузок.</li>" +
			"<li><b>Поиск</b> — GGUF-модели на Hugging Face по имени (например, «qwen2.5 instruct»). У репозитория: дата последнего обновления, число загрузок и ссылка ↗ на страницу модели; клик по строке раскрывает файлы-кванты. Индикатор ● ≈N GB сравнивается со <b>свободной</b> оперативной памятью (она показана над списком).</li>" +
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
			"<li><b>Портативный вариант</b> — просто скопируйте папку с exe целиком (на флешку, другой ПК): настройки, модели и лог живут рядом с exe и переезжают вместе с ним. В реестр ничего не пишется.</li>" +
			"<li>При первом запуске без модели распознавания приложение само откроет каталог моделей и дождётся скачивания.</li>" +
			"<li>Требования: Windows 10/11 x64, CPU с AVX2 (~2013+), WebView2 Runtime для окна настроек (в Windows 11 есть).</li>" +
			"</ul>" +
			"<p class=\"wh\">Трей и файлы</p>" +
			"<div class=\"mock\" style=\"max-width:290px\"><div class=\"mock-mi dim\">{app} — Готов…</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Настройки…</div><div class=\"mock-mi\">Выключить</div><div class=\"mock-mi\">Копировать последний результат</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Перечитать config.json</div><div class=\"mock-mi\">Открыть config.json</div><div class=\"mock-mi\">Открыть лог</div><div class=\"mock-mi\">О приложении</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Выход</div></div>" +
			"<ul>" +
			"<li>Левый клик по иконке в трее — настройки; правый — меню. Цвет иконки: зелёный — готов, красный — запись, оранжевый — распознавание, серый — выключено или ошибка.</li>" +
			"<li><b>config.json</b> — все настройки; правки вручную применяются через «Перечитать config.json» в трее.</li>" +
			"<li><b>{log}</b> — журнал работы, автоматически ограничен ~2 МБ.</li>" +
			"<li><b>models/</b> — скачанные модели Whisper и LLM.</li>" +
			"</ul>",
		"S_UNSAVED":    "Есть несохранённые изменения. Сохранить?",
		"S_SAVE_YES":   "Сохранить",
		"S_SAVE_NO":    "Не сохранять",
		"S_APPLY_HINT": "Модель переключится после нажатия «Сохранить»",
		"S_VERSION":    "Версия",
	},
	"en": {
		"S_TITLE":          "{app} — Settings",
		"S_TAB_GENERAL":    "General",
		"S_TAB_REC":        "Recognition",
		"S_TAB_PROC":       "Post-processing",
		"S_TR":             "Translation",
		"S_TR_HINT":        "Translation is done by Whisper: to English — native mode, to other languages — experimentally, by forcing the output language (quality depends on the language pair).",
		"S_TR_TURBO":       "⚠ The active Turbo model is not trained for translation to English — pick another model on the Models tab for translating.",
		"S_TR_DEFAULT":     "Always translate to the target language",
		"S_TR_TARGET":      "Target language",
		"S_TR_ASK":         "Language choice",
		"S_TR_ASK_NEVER":   "Don't ask (use default)",
		"S_TR_ASK_ALWAYS":  "Ask every time",
		"S_TR_ASK_TIMEOUT": "Ask with a timeout",
		"S_TR_SECONDS":     "Timeout, sec",
		"S_TR_LANGS":       "Languages in the dialog",
		"S_TAB_SERVER":     "Server",
		"S_PIPE":           "voice ▸ recognition ▸ editing ▸ paste",
		"S_DICT":           "Recognition dictionary",
		"S_DICT_HINT":      "Comma-separated terms, names and abbreviations — a hint for the ear, not commands.",
		"S_LLM":            "Processing profiles",
		"S_LLM_HINT":       "Checked profiles apply one after another, top to bottom, on regular dictation. Nothing checked — text is inserted as is. A profile's hotkey applies just that profile once.",
		"S_LLM_MODEL":      "Editing model (LLM)",
		"S_UPDATED":        "Model last updated",
		"S_FIT_OK":         "fits",
		"S_FIT_WARN":       "tight",
		"S_FIT_BAD":        "not enough RAM",
		"S_RAM":            "Computer RAM:",
		"S_HF_PH":          "Model name — e.g. qwen2.5 instruct",
		"S_NO_LLM":         "No models installed yet — find and download one on the Search tab.",
		"S_NO_LLM_PROF":    "Prompts become available once a model is installed (see Models and Search tabs).",
		"S_PROF_EDIT":      "Edit",
		"S_PROF_CLOSE":     "Collapse",
		"S_CONFIRM_DEL":    "Delete the \"%s\" model? It can be downloaded again.",
		"S_FREE":           "free",
		"S_SUB_MODELS":     "Models",
		"S_SUB_SEARCH":     "Search",
		"S_SUB_PROMPTS":    "Prompts",
		"S_SUB_PARAMS":     "Parameters",
		"S_SUB_DICT":       "Dictionary",
		"S_SUB_TR":         "Translation",
		"S_PROF_ASIS":      "As is",
		"S_PROF_WT":        "Translate → English (fast)",
		"S_PROF_ADD":       "Add profile",
		"S_PROF_NAME":      "Name",
		"S_PROF_PROMPT":    "Prompt",
		"S_PROF_HOTKEY":    "Hotkey",
		"S_PROF_SET":       "Set…",
		"S_PROF_CLEAR":     "Clear",
		"S_PROF_TEST":      "Test",
		"S_TAB_ABOUT":      "About",
		"S_HOTKEY":         "Keyboard shortcut",
		"S_CHANGE":         "Change…",
		"S_UILANG":         "UI language",
		"S_AUTO":           "System default",
		"S_BEEP":           "Recording sound cues",
		"S_SEC_SOUND":      "Sound",
		"S_SEC_BEHAVIOR":   "Behavior",
		"S_SOUND":          "Cue sound",
		"S_SND_SPEECH":     "System (speech)",
		"S_SND_CHIME":      "Chime",
		"S_SND_SOFT":       "Soft",
		"S_SND_MARIMBA":    "Marimba",
		"S_SND_BLIP":       "Blip",
		"S_SND_POP":        "Pop",
		"S_AUTOENTER":      "Press Enter after paste (auto-submit)",
		"S_RESTORE":        "Restore clipboard after paste",
		"S_OVERLAY":        "On-screen indicator",
		"S_ANIM":           "Recording & transcribing animation",
		"S_TYPEMODE":       "Type character-by-character (for paste-blocked fields)",
		"S_RECLANG":        "Recognition language",
		"S_RECAUTO":        "Auto-detect",
		"S_MODELS":         "Recognition models",
		"S_DL":             "Download",
		"S_USE":            "Use",
		"S_DEL":            "Delete",
		"S_ACTIVE":         "Active",
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
		"S_HIDDEN":         "Simple mode · settings hidden: %d",
		"S_ALLSHOWN":       "Everything is shown",
		"S_SHOWALL":        "Show everything",
		"S_SHOWSIMPLE":     "Back to simple",
		"S_HOTMODE":        "Mode",
		"S_HOTMODE_HOLD":   "hold",
		"S_HOTMODE_TOGGLE": "toggle",
		"S_SUB_HOTMODE":    "hold the keys, or press once to start and once to stop",
		"S_GRP_WORK":       "Work",
		"S_GRP_REC":        "Recognition",
		"S_GRP_OTHER":      "Other",
		"S_NAV_STATE":      "Status",
		"S_NAV_DICT":       "Dictation",
		"S_NAV_MIC":        "Microphone",
		"S_NAV_MODELS":     "Models",
		"S_NAV_TEXT":       "Text",
		"S_NAV_TR":         "Translation",
		"S_NAV_SYSTEM":     "System",
		"S_NAV_ABOUT":      "About",
		"S_STATE_HINT":     "hold it and speak — the text lands where the cursor is",
		"S_STATE_ENGINE":   "Recognition",
		"S_STATE_PROC":     "Post-processing",
		"S_STATE_MEM":      "Memory",
		"S_STATE_MEM_SUB":  "models stay loaded, the first phrase has no delay",
		"S_SUB_MINMS":      "ignores accidental key presses",
		"S_SUB_ENTER":      "sends the message right away",
		"S_SUB_CLIP":       "images and files come back as they were",
		"S_SUB_TYPE":       "for fields where pasting is blocked",
		"S_SUB_THREADS":    "more threads is not always faster — measure on your machine",
		"S_SUB_PUNCT":      "where punctuation and capitals come from",
		"S_SUB_TRTARGET":   "English is native for Whisper, other targets are experimental",
		"S_SUB_AUTOSTART":  "turn off if you start the server yourself",
		"S_SUB_PORT":       "applies after a restart",
		"S_SUB_UPD":        "the only network request besides model downloads",
		"S_STATE_LAST":     "Last dictation",
		"S_STATE_COPY":     "Copy",
		"S_SEC_OVERLAY":    "On-screen overlay",
		"S_SEC_SERVICE":    "Service",
		"S_SEC_LLM":        "Editor model",
		"S_NOT_INSTALLED":  "not installed",
		"S_CHANGE_MODEL":   "Change",
		"S_PICK_MODEL":     "Pick",
		"S_STATE_RU":       "Russian speech",
		"S_STATE_OTHER":    "Other languages",
		"S_ENG_WHISPER":    "whisper.cpp · 99 languages · translation",
		"S_ENG_SHERPA":     "sherpa-onnx · Russian only",
		"S_MODEL":          "Model file",
		"S_MIC":            "Microphone",
		"S_MIC_DEFAULT":    "System default",
		"S_MIC_REFRESH":    "Refresh list",
		"S_MIC_LEVEL":      "Input level",
		"S_MIC_QUIET":      "quiet",
		"S_THREADS":        "CPU threads",
		"S_MINMS":          "Min recording, ms",
		"S_MAXSEC":         "Max recording, s",
		"S_AUTOSTART":      "Start whisper-server automatically",
		"S_PORT":           "Port",
		"S_SERVEREXE":      "whisper-server path",
		"S_SERVERURL":      "External server (URL)",
		"S_URLHINT":        "If set, the local server is not started",
		"S_SAVE":           "Save",
		"S_SAVED":          "Saved",
		"S_RESTART":        "Saved. Model/language/server apply after restart",
		"S_SUB_INFO":       "Info",
		"S_SUB_HELP":       "Guide",
		"S_UPD":            "Updates",
		"S_UPD_CHECK":      "Check for updates",
		"S_UPD_AUTO":       "Check on startup",
		"S_UPD_NONE":       "You are on the latest version",
		"S_UPD_AVAIL":      "Version %s is available.",
		"S_UPD_GO":         "Update",
		"S_UPD_ERR":        "Update check failed",
		"S_UPD_DL":         "Downloading the update…",
		"S_SUB_AUTHOR":     "Author",
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
			"<p>Full pipeline: <b>record → transcribe (Whisper) → translate (if enabled) → prompts (LLM) → paste</b>. Every stage is visible on the overlay.</p>" +
			"<p class=\"wh\">Overlay</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\"></span><span>Speak…</span><span class=\"mock-bars\"><i style=\"height:6px\"></i><i style=\"height:13px\"></i><i style=\"height:9px\"></i><i style=\"height:16px\"></i><i style=\"height:7px\"></i><i style=\"height:12px\"></i><i style=\"height:5px\"></i></span><span class=\"mock-x\">✕</span></div></div>" +
			"<ul>" +
			"<li><b>Speak…</b> — recording: a red dot and live voice-level bars.</li>" +
			"<li><b>Transcribing…</b> — Whisper is processing; while translating — \"Translating\", while running prompts — \"Editing: name (1/2)\".</li>" +
			"<li><b>Inserted: N chars</b> — done; on errors or silence a short reason is shown.</li>" +
			"<li>The ✕ on the right cancels at any stage; the overlay never steals input focus. The overlay and its animation can be turned off on the General tab.</li>" +
			"</ul>" +
			"<p class=\"wh\">Language dialog</p>" +
			"<div class=\"mock\"><div style=\"display:flex;justify-content:space-between;margin-bottom:8px\"><b>Translate to: (3)</b><span class=\"mock-x\">✕</span></div><div style=\"display:flex;gap:8px\"><span class=\"mock-btn on\">EN</span><span class=\"mock-btn\">DE</span><span class=\"mock-btn\">ES</span></div></div>" +
			"<p>Appears above the overlay after releasing the hotkey in the \"always ask\" / \"ask with timeout\" modes. The button set comes from \"Languages in the dialog\"; the target language is highlighted. In timeout mode the title counts down and the target is applied when time runs out. The dialog ✕ inserts without translation; the overlay ✕ cancels the whole operation.</p>" +
			"<p class=\"wh\">Safe insertion</p>" +
			"<div class=\"mock\"><div style=\"display:flex;justify-content:space-between;margin-bottom:8px\"><b style=\"color:var(--amber)\">Focus changed — insert? (30)</b><span class=\"mock-x\">✕</span></div><div style=\"display:flex;gap:8px\"><span class=\"mock-btn on\">Insert here</span><span class=\"mock-btn\">Copy</span></div></div>" +
			"<ul>" +
			"<li>The target window is captured the moment you press the hotkey. If focus changed while the speech was being processed, nothing is pasted — a dialog offers <b>Insert here</b> (into the current window), <b>Copy</b> (to the clipboard) or ✕. When the countdown expires the insertion is cancelled and the text stays in Last Result.</li>" +
			"<li>Enter after paste is pressed only when the target window has not changed.</li>" +
			"<li><b>Last Result</b> — the final text of every dictation is kept in memory until the next one; the tray menu has \"Copy last result\". A failed paste or focus change never loses a dictation.</li>" +
			"</ul>" +
			"<p class=\"wh\">General</p>" +
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
			"<p class=\"wh\">Recognition (Whisper)</p>" +
			"<div class=\"mock\"><div class=\"mock-row\"><span class=\"mock-radio on\"></span><b>Small</b><span class=\"mock-note\">speed/accuracy balance</span><span style=\"margin-left:auto\">466 MB</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-radio\"></span><b>Medium (q5)</b><span class=\"mock-note\">more accurate, recommended</span><span style=\"margin-left:auto\">539 MB ⭳</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-radio\"></span><b>Turbo (q5)</b><span class=\"mock-note\">best CPU accuracy</span><span style=\"margin-left:auto\">574 MB ⭳</span></div></div>" +
			"<ul>" +
			"<li><b>Models</b> — the Whisper catalog: Base (fast, weak PCs), Small (balanced), Medium and Turbo (more accurate, slower; \"q5\" means a quantized build — slightly smaller and faster with almost no quality loss). The radio selects the active one (applied on Save; the recognizer restarts); clicking the radio of a missing model downloads it from the official whisper.cpp repository on Hugging Face. One model serves both recognition and translation.</li>" +
			"<li>whisper-server keeps the model in memory between phrases — the first dictation after startup is slower (loading), afterwards recognition takes 1–3 seconds.</li>" +
			"<li><b>Dictionary</b> — comma-separated terms, names and abbreviations. A hint for Whisper's \"ear\" so rare words are recognized correctly; not commands.</li>" +
			"<li><b>Parameters</b> — microphone selection with a live level meter (speak and the bar moves, so you know the device is heard), CPU threads (more = faster transcription), minimum recording length (filters accidental presses), maximum length (auto-stop). If the chosen device is unplugged the app falls back to the system default; a recording with no speech is never sent for recognition — it reports \"Silence\" instead.</li>" +
			"<li><b>Server</b> — whisper-server starts automatically and runs locally. You can change the port or point to an external server URL — then the local one is not used.</li>" +
			"<li><b>Translation</b> — all translation is done by Whisper: to English via its native translate mode, to other languages <b>experimentally</b>, by forcing the output language (quality depends on the language pair; major languages work best). The Turbo model is not trained for translation — the settings show a warning when it is active. \"Always translate to the target language\" makes every main-hotkey dictation translate to the chosen target with no questions. With the checkbox off, the ask mode applies: always or with a timeout — a language dialog appears before transcription and the target is used when time runs out. The separate translation hotkey translates once without affecting normal dictation. Settings that do not apply in the current mode are greyed out automatically.</li>" +
			"</ul>" +
			"<p class=\"wh\">Post-processing (LLM)</p>" +
			"<p>An optional second layer: a local language model (llama.cpp) edits the transcribed text according to your prompts — removes filler words, changes style, formats. Fully offline, CPU only.</p>" +
			"<div class=\"mock\"><div class=\"mock-row\" style=\"color:var(--dim)\">▸ Qwen/Qwen2.5-3B-Instruct-GGUF<span style=\"margin-left:auto;color:var(--faint)\">2024-09-20</span><span>↓303k</span><span>↗</span></div>" +
			"<div class=\"mock-row\" style=\"padding-left:18px\">q4_k_m.gguf<span style=\"margin-left:auto\">1.9 GB</span><span style=\"color:var(--green)\">● ≈3.7 GB</span><span style=\"color:var(--dim)\">⭳</span></div>" +
			"<div class=\"mock-row\" style=\"padding-left:18px\">q8_0.gguf<span style=\"margin-left:auto\">3.6 GB</span><span style=\"color:var(--amber)\">● ≈8.2 GB</span><span style=\"color:var(--dim)\">⭳</span></div></div>" +
			"<ul>" +
			"<li><b>Models</b> — installed LLM editing models; the radio selects the active one (on Save), ✕ deletes (the active one too — post-processing then turns off). Download progress shows here as well.</li>" +
			"<li><b>Search</b> — GGUF models on Hugging Face by name (e.g. \"qwen2.5 instruct\"). Each repository shows its last update date, download count and a ↗ link to the model page; clicking a row expands its quant files. The ● ≈N GB indicator is compared against the <b>free</b> RAM (shown above the list).</li>" +
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
			"<li><b>Portable</b> — just copy the whole folder with the exe (to a USB stick, another PC): settings, models and the log live next to the exe and travel with it. Nothing is written to the registry.</li>" +
			"<li>On first run without a recognition model the app opens the model catalog itself and waits for the download.</li>" +
			"<li>Requirements: Windows 10/11 x64, a CPU with AVX2 (~2013+), WebView2 Runtime for the settings window (included in Windows 11).</li>" +
			"</ul>" +
			"<p class=\"wh\">Tray and files</p>" +
			"<div class=\"mock\" style=\"max-width:290px\"><div class=\"mock-mi dim\">{app} — Ready…</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Settings…</div><div class=\"mock-mi\">Disable</div><div class=\"mock-mi\">Copy last result</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Reload config.json</div><div class=\"mock-mi\">Open config.json</div><div class=\"mock-mi\">Open log</div><div class=\"mock-mi\">About</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Quit</div></div>" +
			"<ul>" +
			"<li>Left-click the tray icon — settings; right-click — the menu. Icon colors: green — ready, red — recording, orange — transcribing, grey — disabled or error.</li>" +
			"<li><b>config.json</b> — all settings; manual edits apply via \"Reload config.json\" in the tray menu.</li>" +
			"<li><b>{log}</b> — the log, automatically capped at ~2 MB.</li>" +
			"<li><b>models/</b> — downloaded Whisper and LLM models.</li>" +
			"</ul>",
		"S_UNSAVED":    "You have unsaved changes. Save them?",
		"S_SAVE_YES":   "Save",
		"S_SAVE_NO":    "Don't save",
		"S_APPLY_HINT": "The model will switch after you press Save",
		"S_VERSION":    "Version",
	},
}
