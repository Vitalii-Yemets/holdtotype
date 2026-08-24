package hotkeys

import (
	"sort"
	"strings"
)

var systemCombos = map[string]string{
	"win+l":           "lock",
	"win+d":           "desktop",
	"win+e":           "explorer",
	"win+r":           "run",
	"win+i":           "settings",
	"win+s":           "search",
	"win+a":           "center",
	"win+x":           "menu",
	"win+v":           "clipboard",
	"win+g":           "gamebar",
	"win+h":           "voice",
	"win+p":           "project",
	"win+tab":         "tasks",
	"win+space":       "layout",
	"ctrl+win+d":      "newdesktop",
	"ctrl+win+f4":     "closedesktop",
	"shift+win+s":     "snip",
	"alt+tab":         "switch",
	"alt+f4":          "close",
	"alt+esc":         "cycle",
	"ctrl+esc":        "start",
	"ctrl+shift+esc":  "taskmgr",
	"ctrl+alt+delete": "secure",
}

func Normalize(combo string) string {
	parts := strings.Split(strings.ToLower(strings.TrimSpace(combo)), "+")
	out := make([]string, 0, len(parts))
	seen := map[string]bool{}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" || seen[p] {
			continue
		}
		seen[p] = true
		out = append(out, p)
	}
	order := map[string]int{"ctrl": 0, "shift": 1, "alt": 2, "win": 3}
	sort.SliceStable(out, func(i, j int) bool {
		pi, iok := order[out[i]]
		pj, jok := order[out[j]]
		if iok && jok {
			return pi < pj
		}
		if iok != jok {
			return iok
		}
		return out[i] < out[j]
	})
	return strings.Join(out, "+")
}

func TakenBySystem(combo string) string {
	return systemCombos[Normalize(combo)]
}
