package stats

import (
	"strings"
	"unicode"
)

func spaceless(r rune) bool {
	return unicode.In(r, unicode.Han, unicode.Hiragana, unicode.Katakana, unicode.Thai)
}

func Words(text string) int {
	n := 0
	for _, tok := range strings.Fields(text) {
		run := false
		for _, r := range tok {
			switch {
			case spaceless(r):
				if run {
					n++
					run = false
				}
				if !unicode.In(r, unicode.Mn, unicode.Mc) {
					n++
				}
			case unicode.IsLetter(r) || unicode.IsDigit(r):
				run = true
			}
		}
		if run {
			n++
		}
	}
	return n
}
