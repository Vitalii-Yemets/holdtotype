package subtitles

import (
	"strings"
	"time"
	"unicode"
)

const (
	TranslateWaiting = 3
	TranslateMaxWait = 20 * time.Second
)

var TranslateTargets = []string{"de", "en", "es", "fr", "it", "pl", "uk", "ru"}

var languageNames = map[string]string{
	"de": "German", "en": "English", "es": "Spanish", "fr": "French",
	"it": "Italian", "pl": "Polish", "uk": "Ukrainian", "ru": "Russian",
}

func LanguageName(code string) string {
	return languageNames[code]
}

func TranslatePrompt(target string) string {
	name := LanguageName(target)
	if name == "" {
		return ""
	}
	return "You translate subtitles. Translate the user's message into " + name +
		". Reply with the translation only: one line, no quotes, no notes, no explanations. Keep names and numbers as they are. If the message is already in " +
		name + ", repeat it unchanged."
}

type Verdict string

const (
	TrOK     Verdict = ""
	TrEmpty  Verdict = "empty"
	TrSame   Verdict = "same"
	TrRamble Verdict = "ramble"
	TrScript Verdict = "script"
)

const quoteRunes = "\"'`«»‹›„“”‚‘’「」『』〝〞＂"

func isQuote(r rune) bool {
	return strings.ContainsRune(quoteRunes, r)
}

func firstLine(s string) string {
	for _, l := range strings.FieldsFunc(s, func(r rune) bool { return r == '\n' || r == '\r' }) {
		if l = strings.TrimSpace(l); l != "" {
			return l
		}
	}
	return ""
}

func unquote(s string) string {
	for {
		r := []rune(s)
		if len(r) < 2 || !isQuote(r[0]) || !isQuote(r[len(r)-1]) {
			return s
		}
		s = strings.TrimSpace(string(r[1 : len(r)-1]))
	}
}

func plain(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

var foreignScripts = []*unicode.RangeTable{{R16: []unicode.Range16{
	{Lo: 0x2e80, Hi: 0x9fff, Stride: 1},
	{Lo: 0xac00, Hi: 0xd7af, Stride: 1},
	{Lo: 0xff00, Hi: 0xffef, Stride: 1},
}}}

const translateGrammar = "root ::= [^⺀-鿿가-힯぀-ヿ＀-￯]*"

func TranslateGrammar(target string) string {
	if LanguageName(target) == "" {
		return ""
	}
	return translateGrammar
}

func wrongScript(answer, target string) bool {
	latin, cyrillic, letters := 0, 0, 0
	for _, r := range answer {
		if unicode.In(r, foreignScripts...) {
			return true
		}
		if !unicode.IsLetter(r) {
			continue
		}
		letters++
		switch {
		case unicode.Is(unicode.Latin, r):
			latin++
		case unicode.Is(unicode.Cyrillic, r):
			cyrillic++
		}
	}
	if letters == 0 {
		return false
	}
	switch target {
	case "ru", "uk":
		return 2*latin >= letters
	case "de", "en", "es", "fr", "it", "pl":
		return 2*cyrillic >= letters
	}
	return false
}

func CleanTranslation(line, answer, target string) (string, Verdict) {
	out := oneLine(unquote(firstLine(answer)))
	p := plain(out)
	if p == "" {
		return "", TrEmpty
	}
	if p == plain(line) {
		return "", TrSame
	}
	if wrongScript(out, target) {
		return "", TrScript
	}
	if len([]rune(out)) > 4*len([]rune(line))+40 {
		return "", TrRamble
	}
	return out, TrOK
}

type Pending struct {
	ID     int
	At     time.Time
	Text   string
	Target string
	Queued time.Time
}

type QueueCounts struct {
	Pushed   int
	Overflow int
	Stale    int
	Cleared  int
}

func (c QueueCounts) Dropped() int { return c.Overflow + c.Stale + c.Cleared }

type TranslateQueue struct {
	limit   int
	maxWait time.Duration
	items   []Pending
	counts  QueueCounts
}

func NewTranslateQueue(limit int, maxWait time.Duration) *TranslateQueue {
	if limit < 1 {
		limit = 1
	}
	return &TranslateQueue{limit: limit, maxWait: maxWait}
}

func (q *TranslateQueue) Push(p Pending, now time.Time) int {
	p.Queued = now
	q.counts.Pushed++
	dropped := 0
	for len(q.items) >= q.limit {
		q.items[0] = Pending{}
		q.items = q.items[1:]
		dropped++
	}
	q.counts.Overflow += dropped
	q.items = append(q.items, p)
	return dropped
}

func (q *TranslateQueue) Pop(now time.Time) (Pending, bool, int) {
	stale := 0
	for len(q.items) > 0 {
		p := q.items[0]
		q.items[0] = Pending{}
		q.items = q.items[1:]
		if q.maxWait > 0 && now.Sub(p.Queued) > q.maxWait {
			stale++
			continue
		}
		q.counts.Stale += stale
		return p, true, stale
	}
	q.counts.Stale += stale
	return Pending{}, false, stale
}

func (q *TranslateQueue) Clear() int {
	n := len(q.items)
	q.items = nil
	q.counts.Cleared += n
	return n
}

func (q *TranslateQueue) Len() int { return len(q.items) }

func (q *TranslateQueue) Counts() QueueCounts { return q.counts }
