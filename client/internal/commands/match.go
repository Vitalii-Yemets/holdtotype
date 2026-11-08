package commands

import (
	"strings"
	"unicode"
)

type token struct {
	start, end int
	word       string
}

func isTokenRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsMark(r)
}

func tokenize(s string) []token {
	var out []token
	start := -1
	for i, r := range s {
		if isTokenRune(r) {
			if start < 0 {
				start = i
			}
			continue
		}
		if start >= 0 {
			out = append(out, token{start: start, end: i, word: strings.ToLower(s[start:i])})
			start = -1
		}
	}
	if start >= 0 {
		out = append(out, token{start: start, end: len(s), word: strings.ToLower(s[start:])})
	}
	return out
}

func words(s string) []string {
	toks := tokenize(s)
	out := make([]string, len(toks))
	for i, t := range toks {
		out[i] = t.word
	}
	return out
}

func endsWith(toks []token, p []string) bool {
	if len(p) == 0 || len(p) > len(toks) {
		return false
	}
	off := len(toks) - len(p)
	for i, w := range p {
		if toks[off+i].word != w {
			return false
		}
	}
	return true
}

func trimTail(s string) string {
	return strings.TrimRightFunc(s, func(r rune) bool {
		return unicode.IsSpace(r) || strings.ContainsRune(",;:—–-", r)
	})
}

func cutShapes(cmds []Command, phrases [][]string, text string, toks []token) (string, []Command, []string) {
	cut := len(toks)
	var hits []Command
	var used []string
	for n := 0; n < maxShapes; n++ {
		best, size := -1, 0
		for i, c := range cmds {
			if !IsShape(c.Action) || len(phrases[i]) <= size {
				continue
			}
			if endsWith(toks[:cut], phrases[i]) {
				best, size = i, len(phrases[i])
			}
		}
		if best < 0 {
			break
		}
		hits = append(hits, cmds[best])
		used = append(used, cmds[best].Phrase)
		cut -= size
	}
	if len(hits) == 0 {
		return text, nil, nil
	}
	return trimTail(text[:toks[cut].start]), ordered(hits), used
}

func ordered(hits []Command) []Command {
	const structure = "structure"
	pick := map[string]Command{}
	for _, h := range hits {
		key := h.Action
		if key == ActionList || key == ActionHeading {
			key = structure
		}
		if _, ok := pick[key]; !ok {
			pick[key] = h
		}
	}
	var out []Command
	for _, key := range []string{ActionNoStop, ActionCapital, ActionQuotes, structure} {
		if c, ok := pick[key]; ok {
			out = append(out, c)
		}
	}
	return out
}
