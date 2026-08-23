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
