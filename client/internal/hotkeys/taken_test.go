package hotkeys

import "testing"

func TestNormalizePutsModifiersFirst(t *testing.T) {
	cases := map[string]string{
		"win+ctrl+d":  "ctrl+win+d",
		"D+WIN+CTRL":  "ctrl+win+d",
		" alt + tab ": "alt+tab",
		"win+win+l":   "win+l",
	}
	for in, want := range cases {
		if got := Normalize(in); got != want {
			t.Errorf("Normalize(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestTakenBySystemKnowsWindowsShortcuts(t *testing.T) {
	taken := map[string]string{
		"win+l":          "lock",
		"WIN+L":          "lock",
		"l+win":          "lock",
		"alt+tab":        "switch",
		"ctrl+shift+esc": "taskmgr",
		"shift+win+s":    "snip",
		"win+ctrl+d":     "newdesktop",
	}
	for combo, want := range taken {
		if got := TakenBySystem(combo); got != want {
			t.Errorf("TakenBySystem(%q) = %q, want %q", combo, got, want)
		}
	}
}

func TestTakenBySystemLeavesOrdinaryCombosAlone(t *testing.T) {
	free := []string{"ctrl+win", "ctrl+alt+j", "f9", "ctrl+shift+space", "", "win+j"}
	for _, combo := range free {
		if got := TakenBySystem(combo); got != "" {
			t.Errorf("TakenBySystem(%q) = %q, want empty", combo, got)
		}
	}
}

func TestLoneKeyWarnsOnlyForAnOrdinaryKeyByItself(t *testing.T) {
	for combo, want := range map[string]bool{
		"backspace": true, "f8": true, "a": true, "space": true, "num5": true,
		"ctrl": false, "shift": false, "win": false,
		"shift+backspace": false, "ctrl+win": false, "a+s": false, "": false,
	} {
		if got := LoneKey(combo); got != want {
			t.Errorf("LoneKey(%q) = %v, want %v", combo, got, want)
		}
	}
}
