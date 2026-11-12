package subtitles

import (
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	TranslateWaiting  = 6
	TranslateMaxWait  = 10 * time.Second
	TranslateMerge    = 30
	TranslateHoldMax  = 40
	TranslateWarmWait = 90 * time.Second
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
	"So the first thing we noticed when we looked at the data",
	"was that most of the users who signed up in March",
	"never came back after the first week.",
}

var shotSourcesForEnglish = [3]string{
	"Así que lo primero que notamos cuando miramos los datos",
	"fue que la mayoría de los usuarios que se registraron en marzo",
	"nunca volvieron después de la primera semana.",
}

var shotAnswers = map[string][3]string{
	"en": shotSources,
	"ru": {"Итак, первое, что мы заметили, когда посмотрели на данные,", "это то, что большинство пользователей, зарегистрировавшихся в марте,", "так и не вернулись после первой недели."},
	"uk": {"Отже, перше, що ми помітили, коли подивилися на дані,", "це те, що більшість користувачів, які зареєструвалися в березні,", "так і не повернулися після першого тижня."},
	"de": {"Also, das Erste, was uns auffiel, als wir uns die Daten ansahen,", "war, dass die meisten Nutzer, die sich im März angemeldet hatten,", "nach der ersten Woche nie wiederkamen."},
	"fr": {"Donc la première chose que nous avons remarquée en regardant les données,", "c'est que la plupart des utilisateurs qui se sont inscrits en mars", "ne sont jamais revenus après la première semaine."},
	"es": {"Así que lo primero que notamos cuando miramos los datos", "fue que la mayoría de los usuarios que se registraron en marzo", "nunca volvieron después de la primera semana."},
	"it": {"Quindi la prima cosa che abbiamo notato guardando i dati", "è che la maggior parte degli utenti che si sono registrati a marzo", "non è mai tornata dopo la prima settimana."},
	"pl": {"A więc pierwszą rzeczą, którą zauważyliśmy, gdy spojrzeliśmy na dane,", "było to, że większość użytkowników, którzy zarejestrowali się w marcu,", "nigdy nie wróciła po pierwszym tygodniu."},
}

func TranslateSystem(target string) string {
	name := LanguageName(target)
	if name == "" {
		return ""
	}
	return "You translate live subtitles into " + name +
		". The messages are consecutive pieces of one speech, often only a part of a sentence. Reply with the " +
		name + " translation of the new piece only, worded so that it continues your previous replies. " +
		"Never repeat what you have already translated. Do not end the reply with a full stop unless the piece ends with one."
}

type Turn struct {
	Source      string
	Translation string
}

const WarmRounds = 3

func SendAgain(local, timedOut bool, attempt int) bool {
	return local && timedOut && attempt < 2
}

func WarmMessages(target string) [][]Message {
	answers, ok := shotAnswers[target]
	if !ok || TranslateSystem(target) == "" {
		return nil
	}
	sources := shotSources
	if target == "en" {
		sources = shotSourcesForEnglish
	}
	var turns []Turn
	for r := 0; r < WarmRounds; r++ {
		for i := range sources {
			turns = append(turns, Turn{Source: sources[i], Translation: answers[i]})
		}
	}
	return [][]Message{
		TranslateMessages(target, turns, "Thank you."),
		TranslateMessages(target, nil, "Hello."),
	}
}

func TranslateMessages(target string, history []Turn, text string) []Message {
	system := TranslateSystem(target)
	answers, ok := shotAnswers[target]
	if system == "" || !ok {
		return nil
	}
	sources := shotSources
	if target == "en" {
		sources = shotSourcesForEnglish
	}
	out := make([]Message, 0, 8+2*len(history))
	out = append(out, Message{Role: "system", Content: system})
	for i := range sources {
		out = append(out, Message{Role: "user", Content: sources[i]}, Message{Role: "assistant", Content: answers[i]})
	}
	for _, t := range history {
		out = append(out, Message{Role: "user", Content: oneLine(t.Source)}, Message{Role: "assistant", Content: oneLine(t.Translation)})
	}
	return append(out, Message{Role: "user", Content: oneLine(text)})
}

const (
	HistoryLimit      = 8
	HistoryKeepServer = 3
	HistoryKeepLocal  = 1
)

type History struct {
	limit int
	keep  int
	turns []Turn
}

func NewHistory(limit, keep int) *History {
	if limit < 1 {
		limit = 1
	}
	if keep < 0 {
		keep = 0
	}
	if keep >= limit {
		keep = limit - 1
	}
	return &History{limit: limit, keep: keep}
}

func (h *History) Turns() []Turn { return append([]Turn(nil), h.turns...) }

func (h *History) Len() int { return len(h.turns) }

func (h *History) Add(t Turn) {
	h.turns = append(h.turns, t)
	if len(h.turns) >= h.limit {
		h.turns = append([]Turn(nil), h.turns[len(h.turns)-h.keep:]...)
	}
}

func (h *History) Clear() { h.turns = nil }

func lastWord(text string) string {
	f := strings.Fields(text)
	if len(f) == 0 {
		return ""
	}
	return f[len(f)-1]
}

func DropAddedFullStop(source, answer string) string {
	if EndsSentenceBefore(lastWord(source), "", true) {
		return answer
	}
	core := strings.TrimRightFunc(answer, isCloser)
	if !strings.HasSuffix(core, ".") || strings.HasSuffix(core, "..") {
		return answer
	}
	return strings.TrimRightFunc(core[:len(core)-1], unicode.IsSpace) + answer[len(core):]
}

func isSentenceMark(r rune) bool { return r == '.' || r == '!' || r == '?' || r == '…' }

func KeepSentenceEnd(source, answer string) string {
	last := lastWord(source)
	if !EndsSentenceBefore(last, "", true) {
		return answer
	}
	src := strings.TrimRightFunc(last, isCloser)
	body := strings.TrimRightFunc(src, isSentenceMark)
	mark := src[len(body):]
	if mark == "" {
		mark = "."
	}
	trimmed := strings.TrimRightFunc(answer, unicode.IsSpace)
	core := strings.TrimRightFunc(trimmed, isCloser)
	tail := trimmed[len(core):]
	core = strings.TrimRightFunc(core, unicode.IsSpace)
	if core == "" {
		return answer
	}
	if r, _ := utf8.DecodeLastRuneInString(core); isSentenceMark(r) {
		return answer
	}
	cut := strings.TrimRightFunc(core, func(r rune) bool { return strings.ContainsRune(",;:-–—", r) })
	if cut != core {
		core = strings.TrimRightFunc(cut, unicode.IsSpace)
	}
	return core + mark + tail
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

	held  bool
	since time.Time
}

type Batch struct {
	Seqs   []int
	Text   string
	Words  int
	Target string
	Para   bool
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
	limit     int
	maxWait   time.Duration
	merge     int
	holdLimit int
	items     []Pending
	counts    QueueCounts

	holding   bool
	holdGen   int
	holdFor   string
	holdSince time.Time
	holdCount int
}

func NewTranslateQueue(limit int, maxWait time.Duration, merge int) *TranslateQueue {
	if limit < 1 {
		limit = 1
	}
	return &TranslateQueue{limit: limit, maxWait: maxWait, merge: merge, holdLimit: TranslateHoldMax}
}

func (q *TranslateQueue) dropAt(i int) Pending {
	p := q.items[i]
	copy(q.items[i:], q.items[i+1:])
	q.items[len(q.items)-1] = Pending{}
	q.items = q.items[:len(q.items)-1]
	return p
}

func (q *TranslateQueue) loose() int {
	n := 0
	for _, it := range q.items {
		if !it.held {
			n++
		}
	}
	return n
}

func (q *TranslateQueue) Push(p Pending) []Pending {
	if p.Words <= 0 {
		p.Words = len(strings.Fields(p.Text))
	}
	p.held, p.since = false, time.Time{}
	q.counts.Pushed++
	var dropped []Pending
	if q.holding {
		p.held = true
		q.holdCount++
		for len(q.items) >= q.holdLimit {
			dropped = append(dropped, q.dropAt(0))
		}
	} else {
		for q.loose() >= q.limit {
			for i := range q.items {
				if !q.items[i].held {
					dropped = append(dropped, q.dropAt(i))
					break
				}
			}
		}
	}
	q.counts.Overflow += len(dropped)
	q.items = append(q.items, p)
	return dropped
}

func (q *TranslateQueue) Hold(label string, now time.Time) (int, bool) {
	if q.holding && q.holdFor == label {
		return q.holdGen, false
	}
	if !q.holding {
		q.holdSince = now
		q.holdCount = 0
	}
	q.holding = true
	q.holdFor = label
	q.holdGen++
	for i := range q.items {
		q.items[i].held = true
	}
	return q.holdGen, true
}

func (q *TranslateQueue) Holding() (bool, string, time.Time) {
	return q.holding, q.holdFor, q.holdSince
}

func (q *TranslateQueue) Release(gen int, now time.Time) (int, bool) {
	if !q.holding || gen != q.holdGen {
		return 0, false
	}
	return q.release(now), true
}

func (q *TranslateQueue) ReleaseAll(now time.Time) (int, bool) {
	if !q.holding {
		return 0, false
	}
	return q.release(now), true
}

func (q *TranslateQueue) release(now time.Time) int {
	q.holding = false
	q.holdFor = ""
	n := 0
	for i := range q.items {
		if q.items[i].held {
			q.items[i].since = now
			n++
		}
	}
	return n
}

func (q *TranslateQueue) take() Pending {
	p := q.items[0]
	q.items[0] = Pending{}
	q.items = q.items[1:]
	return p
}

func (q *TranslateQueue) stale(p Pending, now time.Time) bool {
	from := p.Heard
	if !p.since.IsZero() {
		from = p.since
	}
	return q.maxWait > 0 && now.Sub(from) > q.maxWait
}

func (q *TranslateQueue) Pop(now time.Time) (Batch, bool, []Pending) {
	if q.holding {
		return Batch{}, false, nil
	}
	var dropped []Pending
	for len(q.items) > 0 && q.stale(q.items[0], now) {
		dropped = append(dropped, q.take())
	}
	q.counts.Stale += len(dropped)
	if len(q.items) == 0 {
		return Batch{}, false, dropped
	}
	first := q.take()
	b := Batch{Seqs: []int{first.Seq}, Text: first.Text, Words: first.Words, Target: first.Target, Para: first.Para, Start: first.Start, Heard: first.Heard}
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
