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
