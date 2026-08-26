package livetail

import "strings"

// Tail keeps the end of a growing phrase so a live plate stays one quiet
// line: the freshest words fit, the beginning folds into an ellipsis.
func Tail(s string, max int) string {
	s = strings.Join(strings.Fields(s), " ")
	r := []rune(s)
	if max <= 1 || len(r) <= max {
		return s
	}
	cut := r[len(r)-max+1:]
	if at := indexRune(cut, ' '); at >= 0 && at+1 < len(cut) {
		cut = cut[at+1:]
	}
	return "…" + string(cut)
}

func indexRune(r []rune, want rune) int {
	for i, c := range r {
		if c == want {
			return i
		}
	}
	return -1
}
