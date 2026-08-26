package replace

import (
	"strings"
	"unicode"
)

type Rule struct {
	ID        string `json:"id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Whole     bool   `json:"whole"`
	MatchCase bool   `json:"match_case"`
	Lang      string `json:"lang,omitempty"`
}

// ForLang keeps the rules that apply to the given recognition language: the
// ones written for every language, and the ones pinned to this one.
func ForLang(rules []Rule, lang string) []Rule {
	lang = strings.ToLower(strings.TrimSpace(lang))
	out := make([]Rule, 0, len(rules))
	for _, r := range rules {
		rl := strings.ToLower(strings.TrimSpace(r.Lang))
		if rl == "" || (rl == lang && lang != "" && lang != "auto") {
			out = append(out, r)
		}
	}
	return out
}

func isWordRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

func equalRune(a, b rune, matchCase bool) bool {
	if a == b {
		return true
	}
	if matchCase {
		return false
	}
	return unicode.ToLower(a) == unicode.ToLower(b)
}

func matchAt(src []rune, at int, pat []rune, matchCase bool) bool {
	if at+len(pat) > len(src) {
		return false
	}
	for i, r := range pat {
		if !equalRune(src[at+i], r, matchCase) {
			return false
		}
	}
	return true
}

func boundaryOK(src []rune, at, length int) bool {
	if at > 0 && isWordRune(src[at-1]) && isWordRune(src[at]) {
		return false
	}
	end := at + length
	if end < len(src) && isWordRune(src[end-1]) && isWordRune(src[end]) {
		return false
	}
	return true
}

func applyOne(r Rule, text string) string {
	pat := []rune(r.From)
	if len(pat) == 0 {
		return text
	}
	src := []rune(text)
	to := []rune(r.To)
	out := make([]rune, 0, len(src))
	for i := 0; i < len(src); {
		if matchAt(src, i, pat, r.MatchCase) && (!r.Whole || boundaryOK(src, i, len(pat))) {
			out = append(out, to...)
			i += len(pat)
			continue
		}
		out = append(out, src[i])
		i++
	}
	return string(out)
}

func Apply(rules []Rule, text string) string {
	if text == "" {
		return text
	}
	for _, r := range rules {
		if strings.TrimSpace(r.From) == "" {
			continue
		}
		text = applyOne(r, text)
	}
	return text
}

func Clean(rules []Rule) []Rule {
	out := make([]Rule, 0, len(rules))
	for _, r := range rules {
		r.From = strings.TrimSpace(r.From)
		r.To = strings.TrimSpace(r.To)
		if r.From == "" {
			continue
		}
		out = append(out, r)
	}
	return out
}

func Contains(r Rule, text string) bool {
	pat := []rune(r.From)
	if len(pat) == 0 || text == "" {
		return false
	}
	src := []rune(text)
	for i := 0; i+len(pat) <= len(src); i++ {
		if matchAt(src, i, pat, r.MatchCase) && (!r.Whole || boundaryOK(src, i, len(pat))) {
			return true
		}
	}
	return false
}
