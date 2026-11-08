package textshape

import "strings"

var punctMarks = []string{".", ",", "!", "?", ";", ":", "…", "—", "–"}

func stripLine(line string) string {
	for _, m := range punctMarks {
		line = strings.ReplaceAll(line, m, "")
	}
	return strings.ToLower(strings.Join(strings.Fields(line), " "))
}

func StripPunctuation(s string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = stripLine(l)
	}
	out := strings.Join(lines, "\n")
	for strings.Contains(out, "\n\n\n") {
		out = strings.ReplaceAll(out, "\n\n\n", "\n\n")
	}
	return strings.Trim(out, "\n")
}
