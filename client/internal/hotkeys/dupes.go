package hotkeys

func FindDuplicate(combos []string) string {
	seen := map[string]bool{}
	for _, c := range combos {
		c = Normalize(c)
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

func Same(a, b string) bool {
	n := Normalize(a)
	return n != "" && n == Normalize(b)
}
