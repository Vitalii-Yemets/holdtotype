package hotkeys

import (
	"sort"
	"strings"
)

func CleanPerLang(keys map[string]string, knownLang func(string) bool, parses func(string) bool) (map[string]string, []string) {
	out := map[string]string{}
	var dropped []string
	for code, combo := range keys {
		code = strings.ToLower(strings.TrimSpace(code))
		combo = strings.TrimSpace(combo)
		if code == "" || combo == "" {
			continue
		}
		if !knownLang(code) || !parses(combo) {
			dropped = append(dropped, code)
			continue
		}
		out[code] = combo
	}
	sort.Strings(dropped)
	return out, dropped
}

func DropClashing(keys map[string]string, taken ...string) (map[string]string, []string) {
	seen := map[string]bool{}
	for _, t := range taken {
		if n := Normalize(t); n != "" {
			seen[n] = true
		}
	}
	codes := make([]string, 0, len(keys))
	for code := range keys {
		codes = append(codes, code)
	}
	sort.Strings(codes)
	out := map[string]string{}
	var dropped []string
	for _, code := range codes {
		n := Normalize(keys[code])
		if n == "" || seen[n] {
			dropped = append(dropped, code)
			continue
		}
		seen[n] = true
		out[code] = keys[code]
	}
	return out, dropped
}

func WithPerLang(combos []string, keys map[string]string) []string {
	codes := make([]string, 0, len(keys))
	for code := range keys {
		codes = append(codes, code)
	}
	sort.Strings(codes)
	out := append([]string(nil), combos...)
	for _, code := range codes {
		out = append(out, keys[code])
	}
	return out
}

func Codes(keys map[string]string) []string {
	codes := make([]string, 0, len(keys))
	for code := range keys {
		codes = append(codes, code)
	}
	sort.Strings(codes)
	return codes
}
