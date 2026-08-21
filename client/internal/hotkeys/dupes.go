package hotkeys

import "strings"

func FindDuplicate(combos []string) string {
	seen := map[string]bool{}
	for _, c := range combos {
		c = strings.ToLower(strings.TrimSpace(c))
		if c == "" {
			continue
		}
		if seen[c] {
			return c
		}
		seen[c] = true
	}
	return ""
}
