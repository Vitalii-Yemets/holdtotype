package commands

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"holdtotype/internal/textshape"
)

type ShapeOptions struct {
	Lang       string
	StripPunct bool
}

const (
	defaultHeading = "# "
	defaultQuotes  = "“”"
)

var defaultAnd = map[string]string{
	"en": "and",
	"ru": "и",
	"uk": "і, й, та",
	"de": "und",
	"fr": "et",
	"es": "y, e",
	"it": "e, ed",
	"pl": "i, oraz",
}

func normLang(lang string) string {
	lang = strings.ToLower(strings.TrimSpace(lang))
	if i := strings.IndexAny(lang, "-_"); i >= 0 {
		lang = lang[:i]
	}
	return lang
}

func DefaultAnd(lang string) string {
	return defaultAnd[normLang(lang)]
}

func andWords(param, lang string) [][]string {
	src := param
	if strings.TrimSpace(src) == "" {
		src = DefaultAnd(lang)
	}
	var out [][]string
	for _, w := range strings.Split(src, ",") {
		if p := words(w); len(p) > 0 {
			out = append(out, p)
		}
	}
	return out
}

func quoteMarks(m string) (string, string) {
	if strings.TrimSpace(m) == "" {
		m = defaultQuotes
	}
	r := []rune(m)
	if len(r) == 1 {
		return m, m
	}
	h := len(r) / 2
	if len(r)%2 == 1 {
		return string(r[:h]), string(r[h+1:])
	}
	return string(r[:h]), string(r[h:])
}

func trimPeriod(s string) string {
	s = strings.TrimRightFunc(s, unicode.IsSpace)
	if strings.HasSuffix(s, ".") && !strings.HasSuffix(s, "..") {
		s = strings.TrimRightFunc(s[:len(s)-1], unicode.IsSpace)
	}
	return s
}

func dropStops(s string) string {
	return strings.TrimRightFunc(s, func(r rune) bool {
		return r == '.' || r == '…' || unicode.IsSpace(r)
	})
}

func capitalize(s string) string {
	for i, r := range s {
		if unicode.IsSpace(r) || unicode.IsPunct(r) {
			continue
		}
		if !unicode.IsLetter(r) {
			return s
		}
		return s[:i] + string(unicode.ToTitle(r)) + s[i+utf8.RuneLen(r):]
	}
	return s
}

func joinLines(s string) string {
	var parts []string
	for _, l := range strings.Split(s, "\n") {
		if l = strings.TrimSpace(l); l != "" {
			parts = append(parts, l)
		}
	}
	return strings.Join(parts, " ")
}

func leadColon(s string) int {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case ',', ';', '\n':
			return -1
		case ':':
			if i+1 == len(s) || unicode.IsSpace(rune(s[i+1])) {
				return i
			}
		}
	}
	return -1
}

func splitLastAnd(s string, and [][]string) (string, string, bool) {
	toks := tokenize(s)
	for i := len(toks) - 1; i >= 0; i-- {
		for _, p := range and {
			if i+len(p) > len(toks) || !endsWith(toks[:i+len(p)], p) {
				continue
			}
			return s[:toks[i].start], s[toks[i+len(p)-1].end:], true
		}
	}
	return s, "", false
}

func splitList(text string, and [][]string) (string, []string) {
	text = strings.TrimSpace(text)
	lead := ""
	if i := leadColon(text); i >= 0 {
		if strings.TrimSpace(text[:i]) != "" {
			lead = strings.TrimSpace(text[:i+1])
		}
		text = text[i+1:]
	}
	parts := strings.FieldsFunc(text, func(r rune) bool {
		return r == ',' || r == ';' || r == '\n'
	})
	if n := len(parts); n > 0 {
		if before, after, ok := splitLastAnd(parts[n-1], and); ok {
			parts = append(parts[:n-1], before, after)
		}
	}
	var items []string
	for _, p := range parts {
		if p = trimPeriod(strings.TrimSpace(p)); p != "" {
			items = append(items, p)
		}
	}
	return lead, items
}

func Shape(text string, shapes []Command, opts ShapeOptions) string {
	if strings.TrimSpace(text) == "" || len(shapes) == 0 {
		return text
	}
	var nostop, capital bool
	var quotes, structure *Command
	for i := range shapes {
		switch shapes[i].Action {
		case ActionNoStop:
			nostop = true
		case ActionCapital:
			capital = true
		case ActionQuotes:
			quotes = &shapes[i]
		case ActionList, ActionHeading:
			structure = &shapes[i]
		}
	}
	form := func(s string, item bool) string {
		if opts.StripPunct {
			s = textshape.StripPunctuation(s)
		}
		if nostop {
			s = dropStops(s)
		}
		if capital {
			s = capitalize(s)
		}
		if quotes != nil && item && strings.TrimSpace(s) != "" {
			open, end := quoteMarks(quotes.Text)
			s = open + trimPeriod(s) + end
		}
		return s
	}
	switch {
	case structure != nil && structure.Action == ActionList:
		lead, items := splitList(text, andWords(structure.Text, opts.Lang))
		var lines []string
		if lead != "" {
			if l := form(lead, false); strings.TrimSpace(l) != "" {
				lines = append(lines, l)
			}
		}
		for _, it := range items {
			if it = trimPeriod(strings.TrimSpace(form(it, true))); it != "" {
				lines = append(lines, "- "+it)
			}
		}
		return strings.Join(lines, "\n")
	case structure != nil:
		marker := structure.Text
		if marker == "" {
			marker = defaultHeading
		}
		body := trimPeriod(strings.TrimSpace(form(joinLines(text), true)))
		if body == "" {
			return ""
		}
		return marker + body + "\n\n"
	}
	return form(strings.TrimSpace(text), true)
}
