package main

import "holdtotype/internal/hotkeys"

func hotkeyWarning(combo string) string {
	if hotkeys.LoneKey(combo) {
		return trf("hk.lone", combo)
	}
	what := hotkeys.TakenBySystem(combo)
	if what == "" {
		return ""
	}
	return trf("hk.taken", combo, tr("hk."+what))
}

const langProfilePrefix = "lang:"

func langCombos(cfg *Config) []comboDef {
	var combos []comboDef
	for _, code := range hotkeys.Codes(cfg.LangHotkeys) {
		if groups, err := parseHotkey(cfg.LangHotkeys[code]); err == nil {
			combos = append(combos, comboDef{id: langProfilePrefix + code, groups: groups})
		}
	}
	return combos
}
