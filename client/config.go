package main

import (
	"holdtotype/internal/commands"
	"holdtotype/internal/history"
	"holdtotype/internal/mojibake"
	"holdtotype/internal/ovplace"
	"holdtotype/internal/preset"
	"holdtotype/internal/profiles"
	"holdtotype/internal/replace"
	"holdtotype/internal/theme"

	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
)

type Profile = profiles.Profile

const configVersion = 5

const (
	hotkeyHold   = "hold"
	hotkeyToggle = "toggle"
	levelSimple  = "simple"
	levelAll     = "all"
)

func validHotkeyMode(v string) bool { return v == hotkeyHold || v == hotkeyToggle }

func validUILevel(v string) bool { return v == levelSimple || v == levelAll }

const (
	ovPosBottom = ovplace.PosBottom
	ovPosTop    = ovplace.PosTop
	ovPosCaret  = ovplace.PosCaret
	ovPosCustom = ovplace.PosCustom
)

func validOverlayPos(v string) bool { return ovplace.Valid(v) }

func validOverlayMonitor(v string) bool {
	if v == "" || v == "cursor" {
		return true
	}
	if len(v) > 2 {
		return false
	}
	for _, r := range v {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

const (
	punctFromModel = "model"
	punctByLLM     = "llm"
	punctOff       = "off"
)

func validPunctuation(v string) bool {
	return v == punctFromModel || v == punctByLLM || v == punctOff
}

type Config struct {
	Hotkey           string            `json:"hotkey"`
	Language         string            `json:"language"`
	Model            string            `json:"model"`
	LangModels       map[string]string `json:"lang_models"`
	Threads          int               `json:"threads"`
	ServerPort       int               `json:"server_port"`
	ServerAutostart  bool              `json:"server_autostart"`
	ServerExe        string            `json:"server_exe"`
	ServerURL        string            `json:"server_url"`
	SherpaExe        string            `json:"sherpa_exe"`
	SherpaPort       int    `json:"sherpa_port"`
	SherpaModel      string `json:"sherpa_model"`
	StreamExe        string `json:"stream_exe"`
	StreamPort       int    `json:"stream_port"`
	StreamModel      string `json:"stream_model"`
	ConfigVersion    int    `json:"config_version"`
	SherpaThreads    int    `json:"sherpa_threads"`
	EngineIdleMin    int    `json:"engine_idle_minutes"`
	Punctuation      string `json:"punctuation"`
	HotkeyMode       string `json:"hotkey_mode"`
	UILevel          string `json:"ui_level"`
	Skin             string `json:"skin"`
	Theme            string `json:"theme"`
	PasteMode        string `json:"paste_mode"`
	RestoreClipboard bool   `json:"restore_clipboard"`
	Beep             bool   `json:"beep"`
	SoundTheme       string `json:"sound_theme"`
	AutoEnter        bool   `json:"auto_enter"`
	Overlay          bool   `json:"overlay"`
	OverlayPos       string `json:"overlay_position"`
	OverlayMonitor   string `json:"overlay_monitor"`
	OverlayXY        map[string]ovplace.Frac `json:"overlay_custom"`
	OverlayText      bool   `json:"overlay_text"`
	PasteDelayMs     int    `json:"paste_delay_ms"`
	UILanguage       string `json:"ui_language"`
	MaxRecordSeconds int    `json:"max_record_seconds"`
	MinRecordMs      int    `json:"min_record_ms"`

	WhisperPrompt       string             `json:"whisper_prompt"`
	TranslateHotkey     string             `json:"translate_hotkey"`
	TranslateTarget     string             `json:"translate_target"`
	TranslateAsk        string             `json:"translate_ask"`
	TranslateAskSeconds int                `json:"translate_ask_seconds"`
	TranslateAskLangs   []string           `json:"translate_ask_langs"`
	DefaultProfile      string             `json:"default_profile"`
	TranslateDefault    bool               `json:"translate_default"`
	ActiveProfiles      []string           `json:"active_profiles"`
	Profiles            []Profile          `json:"profiles"`
	LLMPort             int                `json:"llm_port"`
	LLMExe              string             `json:"llm_exe"`
	LLMModel            string             `json:"llm_model"`
	SettingsW           int                `json:"settings_width"`
	SettingsH           int                `json:"settings_height"`
	CheckUpdates        bool               `json:"check_updates"`
	MicDevice           string             `json:"mic_device"`
	MicDeviceName       string             `json:"mic_device_name"`
	Replacements        []replace.Rule     `json:"replacements"`
	HistoryOn           bool               `json:"history"`
	HistoryDays         int                `json:"history_days"`
	HistoryMax          int                `json:"history_max"`
	HistorySkip         string             `json:"history_skip"`
	PostEnabled         bool               `json:"post_enabled"`
	PostSource          string             `json:"post_source"`
	PostAPIURL          string             `json:"post_api_url"`
	PostAPIModel        string             `json:"post_api_model"`
	PostAPIKey          string             `json:"post_api_key"`
	PostAPITimeout      int                `json:"post_api_timeout_s"`
	Commands            []commands.Command `json:"commands"`
	WizardDone          bool               `json:"wizard_done"`

	CanaryTarget string `json:"-"`
}

func presetProfiles() []Profile {
	cleanPrompt := "You clean up dictated text: remove filler words, repetitions and false starts; fix punctuation and capitalization. Always answer in the same language as the input, never translate.\nExample input: nu eee koroche ya khotel nu skazat privet vsem privet\nExample output: Ya khotel skazat: privet vsem!\nReturn only the cleaned text, nothing else."
	formalPrompt := "You rewrite dictated text in a polite, formal business style. Always answer in the same language as the input, never translate.\nExample input: надо переделать отчет до завтра\nExample output: Пожалуйста, переработайте отчёт к завтрашнему дню.\nReturn only the rewritten text, nothing else."
	translatePrompt := "Translate the following text to English. Return only the translation, nothing else."
	return []Profile{
		{ID: "clean", Name: "Cleanup", Prompt: cleanPrompt},
		{ID: "formal", Name: "Business style", Prompt: formalPrompt},
		{ID: "translate-en", Name: "Translate → English (quality)", Prompt: translatePrompt},
	}
}

var legacyProfileNames = map[string][]string{
	"clean":        {"Чистка речи", "Очищення мовлення"},
	"formal":       {"Деловой стиль", "Діловий стиль"},
	"translate-en": {"Перевод → English (качественный)", "Переклад → English (якісний)"},
}

func renameLegacyProfiles(cfg *Config) bool {
	changed := false
	presets := presetProfiles()
	for i := range cfg.Profiles {
		old, ok := legacyProfileNames[cfg.Profiles[i].ID]
		if !ok {
			continue
		}
		for _, name := range old {
			if cfg.Profiles[i].Name != name {
				continue
			}
			for _, p := range presets {
				if p.ID == cfg.Profiles[i].ID {
					log.Printf("built-in prompt %s renamed to English", p.ID)
					cfg.Profiles[i].Name = p.Name
					changed = true
				}
			}
		}
	}
	return changed
}

func defaultConfig() *Config {
	return &Config{
		Hotkey:           "ctrl+win",
		Language:         "ru",
		Model:            "models/ggml-medium-q5_0.bin",
		LangModels:       map[string]string{"auto": defaultPresetModel},
		Threads:          4,
		ServerPort:       8910,
		ServerAutostart:  true,
		ServerExe:        "whisper-server.exe",
		SherpaExe:        "sherpa-server.exe",
		SherpaPort:       8912,
		StreamExe:        "sherpa-online-server.exe",
		StreamPort:       8913,
		StreamModel:      "models/nemotron-3.5",
		ConfigVersion:    configVersion,
		SherpaThreads:    4,
		EngineIdleMin:    10,
		Punctuation:      punctFromModel,
		HotkeyMode:       hotkeyHold,
		UILevel:          levelAll,
		Skin:             theme.DefaultSkin,
		Theme:            theme.DefaultPalette,
		SherpaModel:      "models/gigaam-v3",
		PasteMode:        "clipboard",
		RestoreClipboard: true,
		Beep:             true,
		SoundTheme:       "speech",
		Overlay:          true,
		OverlayPos:       ovPosBottom,
		OverlayText:      true,
		UILanguage:       "auto",
		MaxRecordSeconds: 120,
		MinRecordMs:      300,
		WhisperPrompt:    builtinDictionary(lang()),
		LLMPort:          8911,
		LLMExe:           "llama-server.exe",
		LLMModel:         "models/" + llmFile,
		PostEnabled:      false,
		PostSource:       "local",
		WizardDone:       true,
		HistoryDays:      history.DefaultKeepDays,
		HistoryMax:       history.DefaultMax,
	}
}

var saveMu sync.Mutex

func withBOM(out []byte) []byte {
	return append([]byte{0xEF, 0xBB, 0xBF}, append(out, '\n')...)
}

func saveConfig(path string, cfg *Config) error {
	saveMu.Lock()
	defer saveMu.Unlock()
	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, withBOM(out), 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}

func fixConfigText(cfg *Config) bool {
	changed := false
	fix := func(p *string) {
		if v := mojibake.Fix(*p); v != *p {
			*p = v
			changed = true
		}
	}
	fix(&cfg.WhisperPrompt)
	for i := range cfg.Profiles {
		fix(&cfg.Profiles[i].Name)
		fix(&cfg.Profiles[i].Prompt)
	}
	return changed
}

// fileNamesSkin says whether the file on disk already speaks of a design of
// its own — a file written before the split names only the colour.
func fileNamesSkin(data []byte) bool {
	var probe struct {
		Skin *string `json:"skin"`
	}
	if json.Unmarshal(data, &probe) != nil {
		return false
	}
	return probe.Skin != nil
}

func loadConfig(path string) (*Config, error) {
	scanLocalModels()
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
	cfg.LangModels = nil
	cfg.PostSource = ""
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})
	logUnknownConfigKeys(data, cfg)
	dec := json.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(cfg); err != nil {
		broken := path + ".broken"
		_ = os.Rename(path, broken)
		cfg = defaultConfig()
		if def, defErr := os.ReadFile("config.default.json"); defErr == nil {
			d2 := json.NewDecoder(bytes.NewReader(bytes.TrimPrefix(def, []byte{0xEF, 0xBB, 0xBF})))
			_ = d2.Decode(cfg)
		}
		_ = saveConfig(path, cfg)
		go msgBox(tr("cfg.err.title"), trf("cfg.err.recovered", humanError(err), broken))
		return cfg, nil
	}
	fileVer := fileConfigVersion(data)
	migrated := fileVer != configVersion
	if !firstRun && fileVer != configVersion {
		backupConfig(path, data, fileVer)
	}
	if fileVer > configVersion {
		log.Printf("config comes from a newer version (v%d, this build understands v%d) — unknown fields will be lost on the first save", fileVer, configVersion)
	}
	if syncDictionary(cfg) {
		migrated = true
	}
	if fixConfigText(cfg) {
		log.Printf("config: the text was mangled by another editor — encoding restored")
		migrated = true
	}
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
	if fileConfigVersion(data) < 4 {
		migrateToPresets(cfg, data)
		migrated = true
	}
	if !firstRun && fileConfigVersion(data) < 5 && (cfg.TranslateDefault || cfg.TranslateAsk == "always" || cfg.TranslateAsk == "timeout") {
		cfg.TranslateDefault = false
		cfg.TranslateAsk = "never"
		log.Printf("translation is off by default in the new scheme — turn it on in the settings if you want it")
		migrated = true
	}
	if cleaned, dropped := preset.Clean(cfg.LangModels, func(id string) *preset.Model {
		return presetView(findModel(id))
	}); len(dropped) > 0 {
		log.Printf("presets: dropped assignments %s", strings.Join(dropped, ", "))
		cfg.LangModels = cleaned
		migrated = true
	} else {
		cfg.LangModels = cleaned
	}
	if cfg.LangModels == nil {
		cfg.LangModels = map[string]string{}
	}
	if applyPreset(cfg) {
		migrated = true
	}
	cfg.ConfigVersion = configVersion
	if !validPunctuation(cfg.Punctuation) {
		cfg.Punctuation = punctFromModel
	}
	if !validHotkeyMode(cfg.HotkeyMode) {
		cfg.HotkeyMode = hotkeyHold
	}
	if !fileNamesSkin(data) && cfg.Theme != "" && !theme.ValidColour(cfg.Theme) {
		// older versions kept one value for both; split it
		skin, colour := theme.Migrate(cfg.Theme)
		cfg.Skin, cfg.Theme = skin, colour
		log.Printf("appearance split into skin %s and colour %s", skin, colour)
		migrated = true
	}
	if !theme.ValidSkin(cfg.Skin) {
		cfg.Skin = theme.DefaultSkin
	}
	if !theme.ValidColour(cfg.Theme) {
		cfg.Theme = theme.DefaultPalette
	}
	if !validUILevel(cfg.UILevel) {
		cfg.UILevel = levelAll
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
	if cfg.StreamExe == "" {
		cfg.StreamExe = "sherpa-online-server.exe"
	}
	if cfg.StreamPort <= 0 {
		cfg.StreamPort = 8913
	}
	if cfg.StreamModel == "" {
		cfg.StreamModel = "models/nemotron-3.5"
	}
	if !validOverlayPos(cfg.OverlayPos) {
		cfg.OverlayPos = ovPosBottom
	}
	if !validOverlayMonitor(cfg.OverlayMonitor) {
		cfg.OverlayMonitor = ""
	}
	cfg.OverlayXY = ovplace.CleanCustom(cfg.OverlayXY)
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
	if !validPostAPIURL(cfg.PostAPIURL) {
		log.Printf("post-processing server address %q could not be parsed — reset", cfg.PostAPIURL)
		cfg.PostAPIURL = ""
	}
	if cfg.PostAPITimeout < 5 || cfg.PostAPITimeout > 120 {
		cfg.PostAPITimeout = 30
	}
	if cfg.PostSource != "local" && cfg.PostSource != "api" {
		if strings.TrimSpace(cfg.PostAPIURL) != "" {
			cfg.PostSource = "api"
		} else {
			cfg.PostSource = "local"
		}
		migrated = true
	}
	if cfg.Profiles == nil {
		cfg.Profiles = presetProfiles()
	}
	if renameLegacyProfiles(cfg) {
		migrated = true
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
			log.Printf("first run: recognition language taken from the system — %s", sys)
			cfg.Language = sys
			migrated = true
		}
	}
	if migrated {
		if err := saveConfig(path, cfg); err != nil {
			log.Printf("saving the upgraded config: %v", err)
		} else if fileVer != configVersion {
			log.Printf("config upgraded: v%d → v%d", fileVer, configVersion)
		}
	}
	return cfg, nil
}

// backupConfig keeps the untouched file of the old version next to the
// config, once per version — so any migration can be undone by hand.
func backupConfig(path string, data []byte, ver int) {
	bak := fmt.Sprintf("%s.v%d.bak", path, ver)
	if _, err := os.Stat(bak); err == nil {
		return
	}
	if err := os.WriteFile(bak, data, 0o644); err != nil {
		log.Printf("backup of the previous config was not written: %v", err)
		return
	}
	log.Printf("backup of the previous config (v%d) saved: %s", ver, bak)
}

var translateLangNames = map[string]string{
	"en": "English", "de": "German", "fr": "French", "es": "Spanish",
	"it": "Italian", "pl": "Polish", "ru": "Russian", "uk": "Ukrainian",
}

var translateLangOrder = []string{"de", "en", "es", "fr", "it", "pl", "uk", "ru"}

func translateLangCodes() []string { return translateLangOrder }

// migrateToPresets rebuilds the two-slot world as language presets: the old
// whisper slot becomes the universal model, and if a sherpa model stood in
// its slot and was installed, its languages keep going to it — exactly what
// the old routing did.
func migrateToPresets(cfg *Config, data []byte) {
	if cfg.LangModels == nil {
		cfg.LangModels = map[string]string{}
	}
	wID := ""
	file := filepath.Base(cfg.Model)
	for i := range modelCatalog {
		if modelCatalog[i].Engine != engineSherpa && modelCatalog[i].File == file {
			wID = modelCatalog[i].ID
		}
	}
	if wID == "" {
		wID = defaultPresetModel
	}
	if cfg.LangModels["auto"] == "" {
		cfg.LangModels["auto"] = wID
	}
	if fileSTTEngine(data) == "whisper" {
		log.Printf("config from an older version: language → model, everything on %s", wID)
		return
	}
	dir := filepath.Base(filepath.Clean(cfg.SherpaModel))
	for i := range modelCatalog {
		m := &modelCatalog[i]
		if m.Engine != engineSherpa || m.Dir != dir || !m.installed() {
			continue
		}
		for _, l := range strings.Split(m.Langs, ",") {
			l = strings.TrimSpace(l)
			if validTranslateLang(l) && cfg.LangModels[l] == "" {
				cfg.LangModels[l] = m.ID
			}
		}
	}
	log.Printf("config from an older version: language → model, universal %s", wID)
}

func fileSTTEngine(data []byte) string {
	var probe struct {
		STTEngine string `json:"stt_engine"`
	}
	if json.Unmarshal(data, &probe) != nil {
		return ""
	}
	return probe.STTEngine
}

func validTranslateLang(l string) bool {
	_, ok := translateLangNames[l]
	return ok
}

func profileByID(cfg *Config, id string) *Profile {
	return profiles.ByID(cfg.Profiles, id)
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

func logUnknownConfigKeys(data []byte, cfg *Config) {
	var raw map[string]json.RawMessage
	if json.Unmarshal(data, &raw) != nil {
		return
	}
	known := map[string]bool{}
	t := reflect.TypeOf(*cfg)
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		known[strings.Split(tag, ",")[0]] = true
	}
	var extra []string
	for k := range raw {
		if !known[k] {
			extra = append(extra, k)
		}
	}
	if len(extra) > 0 {
		sort.Strings(extra)
		log.Printf("config: unknown fields skipped — %s", strings.Join(extra, ", "))
	}
}

func configUILanguage(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "auto"
	}
	var probe struct {
		UILanguage string `json:"ui_language"`
	}
	if json.Unmarshal(bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF}), &probe) != nil {
		return "auto"
	}
	if probe.UILanguage == "" {
		return "auto"
	}
	return probe.UILanguage
}
