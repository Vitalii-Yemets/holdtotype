package subtitles

import (
	"strings"
	"time"
	"unicode"
)

const (
	TranslateWaiting = 6
	TranslateMaxWait = 10 * time.Second
	TranslateMerge   = 30
)

var TranslateTargets = []string{"de", "en", "es", "fr", "it", "pl", "uk", "ru"}

var languageNames = map[string]string{
	"de": "German", "en": "English", "es": "Spanish", "fr": "French",
	"it": "Italian", "pl": "Polish", "uk": "Ukrainian", "ru": "Russian",
}

func LanguageName(code string) string {
	return languageNames[code]
}

const (
	TranslateTimeout   = 8 * time.Second
	TranslateMinTokens = 48
	TranslateMaxTokens = 160
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

var shotSources = [3]string{
	"and then we went back to the office,",
	"which is why the second test took almost three minutes",
	"So if you look at the numbers here, Windows Eleven is",
}

var shotSourcesForEnglish = [3]string{
	"y luego volvimos a la oficina,",
	"weshalb der zweite Test fast drei Minuten gedauert hat",
	"Donc si vous regardez les chiffres ici, Windows Eleven est",
}

var shotAnswers = map[string][3]string{
	"en": shotSources,
	"ru": {"а потом мы вернулись в офис,", "поэтому второй тест занял почти три минуты", "Так что если посмотреть на цифры, Windows Eleven"},
	"uk": {"а потім ми повернулися до офісу,", "тому другий тест тривав майже три хвилини", "Тож якщо подивитися на цифри, Windows Eleven"},
	"de": {"und dann sind wir zurück ins Büro gegangen,", "weshalb der zweite Test fast drei Minuten gedauert hat", "Wenn man sich also die Zahlen hier ansieht, ist Windows Eleven"},
	"fr": {"et ensuite nous sommes retournés au bureau,", "c'est pourquoi le deuxième test a pris presque trois minutes", "Donc si vous regardez les chiffres ici, Windows Eleven est"},
	"es": {"y luego volvimos a la oficina,", "por eso la segunda prueba tardó casi tres minutos", "Así que si miran los números aquí, Windows Eleven es"},
	"it": {"e poi siamo tornati in ufficio,", "ed è per questo che il secondo test ha richiesto quasi tre minuti", "Quindi, se guardate i numeri qui, Windows Eleven è"},
	"pl": {"a potem wróciliśmy do biura,", "dlatego drugi test trwał prawie trzy minuty", "Więc jeśli spojrzycie na te liczby, Windows Eleven jest"},
}

func TranslateSystem(target string) string {
	name := LanguageName(target)
	if name == "" {
		return ""
	}
	return "You translate live subtitles into " + name +
		". Each message is a piece of speech, often only a part of a sentence. Reply with the " +
		name + " translation of that piece only."
}

func TranslateMessages(target, text string) []Message {
	system := TranslateSystem(target)
	answers, ok := shotAnswers[target]
	if system == "" || !ok {
		return nil
	}
	sources := shotSources
	if target == "en" {
		sources = shotSourcesForEnglish
	}
	out := make([]Message, 0, 8)
	out = append(out, Message{Role: "system", Content: system})
	for i := range sources {
		out = append(out, Message{Role: "user", Content: sources[i]}, Message{Role: "assistant", Content: answers[i]})
	}
	return append(out, Message{Role: "user", Content: oneLine(text)})
}

func TranslateTokens(words int) int {
	n := 3*words + 24
	if n < TranslateMinTokens {
		return TranslateMinTokens
	}
	if n > TranslateMaxTokens {
		return TranslateMaxTokens
	}
	return n
}

func StreamingText(sofar string) string {
	line := firstLine(sofar)
	for {
		r := []rune(line)
		if len(r) == 0 || !isQuote(r[0]) {
			break
		}
		line = strings.TrimSpace(string(r[1:]))
	}
	return oneLine(line)
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

func lowerWordsIn(answer string, script *unicode.RangeTable) int {
	n := 0
	for _, f := range strings.Fields(answer) {
		letters := 0
		ok := true
		for _, r := range f {
			if !unicode.IsLetter(r) {
				continue
			}
			letters++
			if !unicode.Is(script, r) || !unicode.IsLower(r) {
				ok = false
				break
			}
		}
		if ok && letters >= 3 {
			n++
		}
	}
	return n
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
		return 2*latin >= letters && lowerWordsIn(answer, unicode.Latin) >= 3
	case "de", "en", "es", "fr", "it", "pl":
		return 2*cyrillic >= letters && lowerWordsIn(answer, unicode.Cyrillic) >= 3
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
	Seq    int
	Text   string
	Words  int
	Target string
	Para   bool
	Start  time.Time
	Heard  time.Time
}

type Batch struct {
	Seqs   []int
	Text   string
	Words  int
	Target string
	Start  time.Time
	Heard  time.Time
}

func (b Batch) Seq() int {
	if len(b.Seqs) == 0 {
		return 0
	}
	return b.Seqs[0]
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
	merge   int
	items   []Pending
	counts  QueueCounts
}

func NewTranslateQueue(limit int, maxWait time.Duration, merge int) *TranslateQueue {
	if limit < 1 {
		limit = 1
	}
	return &TranslateQueue{limit: limit, maxWait: maxWait, merge: merge}
}

func (q *TranslateQueue) Push(p Pending) []Pending {
	if p.Words <= 0 {
		p.Words = len(strings.Fields(p.Text))
	}
	q.counts.Pushed++
	var dropped []Pending
	for len(q.items) >= q.limit {
		dropped = append(dropped, q.items[0])
		q.items[0] = Pending{}
		q.items = q.items[1:]
	}
	q.counts.Overflow += len(dropped)
	q.items = append(q.items, p)
	return dropped
}

func (q *TranslateQueue) take() Pending {
	p := q.items[0]
	q.items[0] = Pending{}
	q.items = q.items[1:]
	return p
}

func (q *TranslateQueue) stale(p Pending, now time.Time) bool {
	return q.maxWait > 0 && now.Sub(p.Heard) > q.maxWait
}

func (q *TranslateQueue) Pop(now time.Time) (Batch, bool, []Pending) {
	var dropped []Pending
	for len(q.items) > 0 && q.stale(q.items[0], now) {
		dropped = append(dropped, q.take())
	}
	q.counts.Stale += len(dropped)
	if len(q.items) == 0 {
		return Batch{}, false, dropped
	}
	first := q.take()
	b := Batch{Seqs: []int{first.Seq}, Text: first.Text, Words: first.Words, Target: first.Target, Start: first.Start, Heard: first.Heard}
	for len(q.items) > 0 {
		next := q.items[0]
		if next.Para || next.Target != b.Target || b.Words+next.Words > q.merge || q.stale(next, now) {
			break
		}
		q.take()
		b.Seqs = append(b.Seqs, next.Seq)
		b.Text += " " + next.Text
		b.Words += next.Words
	}
	return b, true, dropped
}

func (q *TranslateQueue) Clear() []Pending {
	out := q.items
	q.items = nil
	q.counts.Cleared += len(out)
	return out
}

func (q *TranslateQueue) Len() int { return len(q.items) }

func (q *TranslateQueue) Counts() QueueCounts { return q.counts }
