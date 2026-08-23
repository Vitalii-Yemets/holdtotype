package main

import (
	"holdtotype/internal/routing"

	"bytes"
	"encoding/json"
	"log"
	"os"
	"sync"
)

type Profile struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Prompt string `json:"prompt"`
	Hotkey string `json:"hotkey"`
}

const configVersion = 3

const (
	hotkeyHold   = "hold"
	hotkeyToggle = "toggle"
	levelSimple  = "simple"
	levelAll     = "all"
)

func validHotkeyMode(v string) bool { return v == hotkeyHold || v == hotkeyToggle }

func validUILevel(v string) bool { return v == levelSimple || v == levelAll }

const (
	punctFromModel = "model"
	punctByLLM     = "llm"
	punctOff       = "off"
)

func validPunctuation(v string) bool {
	return v == punctFromModel || v == punctByLLM || v == punctOff
}

type Config struct {
	Hotkey           string `json:"hotkey"`
	Language         string `json:"language"`
	Model            string `json:"model"`
	Threads          int    `json:"threads"`
	ServerPort       int    `json:"server_port"`
	ServerAutostart  bool   `json:"server_autostart"`
	ServerExe        string `json:"server_exe"`
	ServerURL        string `json:"server_url"`
	STTEngine        string `json:"stt_engine"`
	SherpaExe        string `json:"sherpa_exe"`
	SherpaPort       int    `json:"sherpa_port"`
	SherpaModel      string `json:"sherpa_model"`
	ConfigVersion    int    `json:"config_version"`
	SherpaThreads    int    `json:"sherpa_threads"`
	EngineIdleMin    int    `json:"engine_idle_minutes"`
	Punctuation      string `json:"punctuation"`
	HotkeyMode       string `json:"hotkey_mode"`
	UILevel          string `json:"ui_level"`
	PasteMode        string `json:"paste_mode"`
	RestoreClipboard bool   `json:"restore_clipboard"`
	Beep             bool   `json:"beep"`
	SoundTheme       string `json:"sound_theme"`
	AutoEnter        bool   `json:"auto_enter"`
	Overlay          bool   `json:"overlay"`
	Animation        bool   `json:"animation"`
	UILanguage       string `json:"ui_language"`
	MaxRecordSeconds int    `json:"max_record_seconds"`
	MinRecordMs      int    `json:"min_record_ms"`

	WhisperPrompt       string    `json:"whisper_prompt"`
	TranslateHotkey     string    `json:"translate_hotkey"`
	TranslateTarget     string    `json:"translate_target"`
	TranslateAsk        string    `json:"translate_ask"`
	TranslateAskSeconds int       `json:"translate_ask_seconds"`
	TranslateAskLangs   []string  `json:"translate_ask_langs"`
	DefaultProfile      string    `json:"default_profile"`
	TranslateDefault    bool      `json:"translate_default"`
	ActiveProfiles      []string  `json:"active_profiles"`
	Profiles            []Profile `json:"profiles"`
	LLMPort             int       `json:"llm_port"`
	LLMExe              string    `json:"llm_exe"`
	LLMModel            string    `json:"llm_model"`
	SettingsW           int       `json:"settings_width"`
	SettingsH           int       `json:"settings_height"`
	CheckUpdates        bool      `json:"check_updates"`
	MicDevice           string    `json:"mic_device"`
	MicDeviceName       string    `json:"mic_device_name"`
}

func presetProfiles() []Profile {
	cleanPrompt := "You clean up dictated text: remove filler words, repetitions and false starts; fix punctuation and capitalization. Always answer in the same language as the input, never translate.\nExample input: nu eee koroche ya khotel nu skazat privet vsem privet\nExample output: Ya khotel skazat: privet vsem!\nReturn only the cleaned text, nothing else."
	formalPrompt := "You rewrite dictated text in a polite, formal business style. Always answer in the same language as the input, never translate.\nExample input: надо переделать отчет до завтра\nExample output: Пожалуйста, переработайте отчёт к завтрашнему дню.\nReturn only the rewritten text, nothing else."
	translatePrompt := "Translate the following text to English. Return only the translation, nothing else."
	if lang() == "ru" {
		return []Profile{
			{ID: "clean", Name: "Чистка речи", Prompt: cleanPrompt},
			{ID: "formal", Name: "Деловой стиль", Prompt: formalPrompt},
			{ID: "translate-en", Name: "Перевод → English (качественный)", Prompt: translatePrompt},
		}
	}
	return []Profile{
		{ID: "clean", Name: "Cleanup", Prompt: cleanPrompt},
		{ID: "formal", Name: "Business style", Prompt: formalPrompt},
		{ID: "translate-en", Name: "Translate → English (quality)", Prompt: translatePrompt},
	}
}

func defaultConfig() *Config {
	return &Config{
		Hotkey:           "ctrl+win",
		Language:         "ru",
		Model:            "models/ggml-small.bin",
		Threads:          4,
		ServerPort:       8910,
		ServerAutostart:  true,
		ServerExe:        "whisper-server.exe",
		STTEngine:        routing.ModeAuto,
		SherpaExe:        "sherpa-server.exe",
		SherpaPort:       8912,
		ConfigVersion:    configVersion,
		SherpaThreads:    4,
		EngineIdleMin:    10,
		Punctuation:      punctFromModel,
		HotkeyMode:       hotkeyHold,
		UILevel:          levelSimple,
		SherpaModel:      "models/gigaam-v3",
		PasteMode:        "clipboard",
		RestoreClipboard: true,
		Beep:             true,
		SoundTheme:       "speech",
		Overlay:          true,
		Animation:        true,
		UILanguage:       "auto",
		MaxRecordSeconds: 120,
		MinRecordMs:      300,
		WhisperPrompt:    "Whisper, whisper.cpp, Docker, Go, LLM, llama.cpp, Qwen, UI, UX, API, HTTP, JSON, GitHub, Windows, exe, ggml, промпт, хоткей, чекбокс, радиокнопка, таймаут, конфиг, вкладка, секция, диктовка, распознавание, постобработка, перевод, интерфейс, локализация, сочетание клавиш, буфер обмена, курсор, микрофон, модель, сервер, трей, оверлей, плашка, скролл, ползунок, диалог, кнопка, профиль, hotkey, checkbox, timeout, dictation, transcription, translation, clipboard, cursor, microphone, overlay, slider, settings, диктування, розпізнавання, налаштування, буфер обміну, Einstellungen, Tastenkürzel, Zwischenablage, Übersetzung, paramètres, raccourci clavier, presse-papiers, traduction, ajustes, atajo de teclado, portapapeles, traducción, impostazioni, scorciatoia, appunti, traduzione, ustawienia, skrót klawiszowy, schowek, tłumaczenie",
		LLMPort:          8911,
		LLMExe:           "llama-server.exe",
		LLMModel:         "models/" + llmFile,
	}
}

var saveMu sync.Mutex

func saveConfig(path string, cfg *Config) error {
	saveMu.Lock()
	defer saveMu.Unlock()
	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(out, '\n'), 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}

func loadConfig(path string) (*Config, error) {
	cfg := defaultConfig()
	data, err := os.ReadFile(path)
	firstRun := os.IsNotExist(err)
	if firstRun {
		if def, defErr := os.ReadFile("config.default.json"); defErr == nil {
			data = def
			_ = os.WriteFile(path, def, 0o644)
			err = nil
		} else {
			out, _ := json.MarshalIndent(cfg, "", "  ")
			_ = os.WriteFile(path, out, 0o644)
			return cfg, nil
		}
	}
	if err != nil {
		return nil, err
	}
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(cfg); err != nil {
		broken := path + ".broken"
		_ = os.Rename(path, broken)
		cfg = defaultConfig()
		if def, defErr := os.ReadFile("config.default.json"); defErr == nil {
			d2 := json.NewDecoder(bytes.NewReader(bytes.TrimPrefix(def, []byte{0xEF, 0xBB, 0xBF})))
			_ = d2.Decode(cfg)
		}
		_ = saveConfig(path, cfg)
		go msgBox(tr("cfg.err.title"), trf("cfg.err.recovered", err.Error(), broken))
		return cfg, nil
	}
	migrated := fileConfigVersion(data) != configVersion
	if cfg.PasteMode != "clipboard" && cfg.PasteMode != "type" {
		cfg.PasteMode = "clipboard"
	}
	if cfg.Threads <= 0 {
		cfg.Threads = 4
	}
	if cfg.MaxRecordSeconds <= 0 {
		cfg.MaxRecordSeconds = 120
	}
	if cfg.ServerPort <= 0 {
		cfg.ServerPort = 8910
	}
	if !validEngine(cfg.STTEngine) {
		cfg.STTEngine = routing.ModeAuto
	}
	if fileConfigVersion(data) < 2 && cfg.STTEngine != routing.ModeAuto {
		log.Printf("конфигурация из прошлой версии: движок %q превращён в автоматический выбор по языку", cfg.STTEngine)
		cfg.STTEngine = routing.ModeAuto
		migrated = true
	}
	cfg.ConfigVersion = configVersion
	if !validPunctuation(cfg.Punctuation) {
		cfg.Punctuation = punctFromModel
	}
	if !validHotkeyMode(cfg.HotkeyMode) {
		cfg.HotkeyMode = hotkeyHold
	}
	if !validUILevel(cfg.UILevel) {
		cfg.UILevel = levelSimple
	}
	if fileConfigVersion(data) < 3 {
		cfg.UILevel = levelAll
		migrated = true
	}
	if cfg.SherpaThreads <= 0 {
		cfg.SherpaThreads = 4
	}
	if cfg.EngineIdleMin <= 0 {
		cfg.EngineIdleMin = 10
	}
	if cfg.SherpaExe == "" {
		cfg.SherpaExe = "sherpa-server.exe"
	}
	if cfg.SherpaPort <= 0 {
		cfg.SherpaPort = 8912
	}
	if cfg.SherpaModel == "" {
		cfg.SherpaModel = "models/gigaam-v3"
	}
	if !validUILang(cfg.UILanguage) {
		cfg.UILanguage = "auto"
	}
	if !validSoundTheme(cfg.SoundTheme) {
		cfg.SoundTheme = "speech"
	}
	if cfg.LLMPort <= 0 {
		cfg.LLMPort = 8911
	}
	if cfg.LLMExe == "" {
		cfg.LLMExe = "llama-server.exe"
	}
	if cfg.LLMModel == "" {
		cfg.LLMModel = "models/" + llmFile
	}
	if cfg.Profiles == nil {
		cfg.Profiles = presetProfiles()
	}
	if cfg.DefaultProfile == "wtranslate" {
		cfg.TranslateDefault = true
		cfg.DefaultProfile = ""
	} else if cfg.DefaultProfile != "" {
		if profileByID(cfg, cfg.DefaultProfile) != nil {
			cfg.ActiveProfiles = append(cfg.ActiveProfiles, cfg.DefaultProfile)
		}
		cfg.DefaultProfile = ""
	}
	var aps []string
	seen := map[string]bool{}
	for _, id := range cfg.ActiveProfiles {
		if !seen[id] && profileByID(cfg, id) != nil {
			seen[id] = true
			aps = append(aps, id)
		}
	}
	cfg.ActiveProfiles = aps
	if !validTranslateLang(cfg.TranslateTarget) {
		cfg.TranslateTarget = "en"
	}
	switch cfg.TranslateAsk {
	case "always", "timeout", "never":
	default:
		cfg.TranslateAsk = "never"
	}
	if cfg.TranslateAskSeconds < 1 || cfg.TranslateAskSeconds > 10 {
		cfg.TranslateAskSeconds = 3
	}
	var langs []string
	for _, l := range cfg.TranslateAskLangs {
		if validTranslateLang(l) {
			langs = append(langs, l)
		}
	}
	if len(langs) == 0 {
		langs = []string{"en", "de"}
	}
	cfg.TranslateAskLangs = langs
	if firstRun {
		if sys := systemLang(); sys != "" && sys != cfg.Language {
			log.Printf("первый запуск: язык распознавания по системе — %s", sys)
			cfg.Language = sys
			migrated = true
		}
	}
	if migrated {
		_ = saveConfig(path, cfg)
	}
	return cfg, nil
}

var translateLangs = []string{"en", "de", "fr", "es", "it", "pl", "ru"}

var translateLangNames = map[string]string{
	"en": "English", "de": "German", "fr": "French", "es": "Spanish",
	"it": "Italian", "pl": "Polish", "ru": "Russian", "uk": "Ukrainian",
}

func validTranslateLang(l string) bool {
	_, ok := translateLangNames[l]
	return ok
}

func profileByID(cfg *Config, id string) *Profile {
	for i := range cfg.Profiles {
		if cfg.Profiles[i].ID == id {
			return &cfg.Profiles[i]
		}
	}
	return nil
}

func fileConfigVersion(data []byte) int {
	var probe struct {
		ConfigVersion int `json:"config_version"`
	}
	if json.Unmarshal(data, &probe) != nil {
		return 0
	}
	return probe.ConfigVersion
}
