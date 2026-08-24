package main

import "holdtotype/internal/hotkeys"

func hotkeyWarning(combo string) string {
	what := hotkeys.TakenBySystem(combo)
	if what == "" {
		return ""
	}
	return trf("hk.taken", combo, tr("hk."+what))
}
