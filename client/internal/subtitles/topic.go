package subtitles

import (
	"regexp"
	"strings"
	"time"
)

const (
	TopicFirstWords  = 150
	TopicEvery       = 3 * time.Minute
	TopicRecentWords = 400
	TopicTokens      = 70
	TermsMax         = 30
	TermsPerCall     = 6
	TermsPairs       = 25
	TermsTokens      = 220
	TopicTimeout     = 30 * time.Second
)

type Step int

const (
	StepNone Step = iota
	StepTopic
	StepTerms
)

type Term struct {
	Source string
	Target string
}

type heardAt struct {
	word string
	at   time.Time
}

type pairAt struct {
	source string
	target string
	at     time.Time
}

type Topic struct {
	target  string
	words   []heardAt
	total   int
	heard   strings.Builder
	pairs   []pairAt
	text    string
	terms   []Term
	pending Step
	last    time.Time
	asked   time.Time
	updates int
	Counts  TopicCounts
}

type TopicCounts struct {
	TopicBad    int
	TermsBad    int
	Known       int
	Unheard     int
	Plain       int
	Example     int
	Full        int
	TermsCalls  int
	TopicCalls  int
	TermsAdded  int
	TopicFailed int
	TermsFailed int
}

func NewTopic(target string) *Topic { return &Topic{target: target} }

func (t *Topic) Target() string { return t.target }

func (t *Topic) Text() string { return t.text }

func (t *Topic) Terms() []Term { return append([]Term(nil), t.terms...) }

func (t *Topic) Updates() int { return t.updates }

func (t *Topic) Heard(text string, at time.Time) {
	for _, w := range strings.Fields(text) {
		t.words = append(t.words, heardAt{w, at})
		t.total++
		t.heard.WriteString(" ")
		t.heard.WriteString(strings.ToLower(w))
	}
	if extra := len(t.words) - 4*TopicRecentWords; extra > 0 {
		t.words = append([]heardAt(nil), t.words[extra:]...)
	}
}

func (t *Topic) Translated(source, target string, at time.Time) {
	if strings.TrimSpace(source) == "" || strings.TrimSpace(target) == "" {
		return
	}
	t.pairs = append(t.pairs, pairAt{source, target, at})
	if extra := len(t.pairs) - 4*TermsPairs; extra > 0 {
		t.pairs = append([]pairAt(nil), t.pairs[extra:]...)
	}
}

func (t *Topic) Due(now time.Time) Step {
	if t.pending != StepNone {
		return t.pending
	}
	if t.last.IsZero() {
		if t.total >= TopicFirstWords {
			return StepTopic
		}
		return StepNone
	}
	if now.Sub(t.last) >= TopicEvery {
		return StepTopic
	}
	return StepNone
}

func (t *Topic) recentWords(now time.Time) []string {
	var out []string
	if t.last.IsZero() {
		from := len(t.words) - TopicFirstWords
		if from < 0 {
			from = 0
		}
		for _, w := range t.words[from:] {
			out = append(out, w.word)
		}
		return out
	}
	for _, w := range t.words {
		if w.at.After(now.Add(-TopicEvery)) {
			out = append(out, w.word)
		}
	}
	if len(out) > TopicRecentWords {
		out = out[len(out)-TopicRecentWords:]
	}
	return out
}

const topicSystem = "You read part of a live transcript of a conversation. The speech recognizer adds little punctuation and mishears some words. " +
	"Write one short line in English, at most 25 words: what the conversation is about and which field or industry it belongs to. " +
	"If an earlier description is given, keep what is still true and make it more precise."

const TopicGrammar = "root ::= \"Topic: \" [^\\n\\r*#]+\n"

const (
	letters  = `a-zA-Z\u00C0-\u024F\u0400-\u04FF`
	cyrillic = `\u0400-\u04FF`
	latin    = `a-zA-Z\u00C0-\u024F`
)

func TermsGrammar(target string) string {
	first := latin
	if target == "ru" || target == "uk" {
		first = cyrillic
	}
	return "root ::= \"none\" | entry (\"\\n\" entry){0,5}\n" +
		"entry ::= src \" = \" dst\n" +
		"src ::= [" + letters + "] [" + letters + "0-9 '/&-]{0,40}\n" +
		"dst ::= [" + first + "] [" + letters + "0-9 ,'()/&-]{0,60}\n"
}

func (t *Topic) TopicMessages(now time.Time) []Message {
	user := ""
	if t.text != "" {
		user = "Earlier description: " + t.text + "\n"
	}
	user += "Transcript: " + strings.Join(t.recentWords(now), " ")
	t.asked = now
	return []Message{{Role: "system", Content: topicSystem}, {Role: "user", Content: user}}
}

func (t *Topic) TakeTopic(raw string, now time.Time) bool {
	t.Counts.TopicCalls++
	line := strings.TrimSpace(raw)
	if !strings.HasPrefix(line, "Topic: ") || strings.ContainsAny(line, "\r\n") || len(line) < 15 {
		t.Counts.TopicBad++
		t.finish(now)
		return false
	}
	t.text = strings.TrimSpace(strings.TrimPrefix(line, "Topic: "))
	t.updates++
	t.pending = StepTerms
	return true
}

func (t *Topic) Failed(step Step, now time.Time) {
	if step == StepTerms {
		t.Counts.TermsFailed++
	} else {
		t.Counts.TopicFailed++
	}
	t.finish(now)
}

func (t *Topic) finish(now time.Time) {
	t.pending = StepNone
	if t.asked.IsZero() {
		t.asked = now
	}
	t.last = t.asked
}

var termExamples = map[string][2]string{
	"ru": {"blood pressure", "артериальное давление"},
	"uk": {"blood pressure", "артеріальний тиск"},
	"de": {"blood pressure", "Blutdruck"},
	"fr": {"blood pressure", "tension artérielle"},
	"es": {"blood pressure", "presión arterial"},
	"it": {"blood pressure", "pressione sanguigna"},
	"pl": {"blood pressure", "ciśnienie krwi"},
	"en": {"presión arterial", "blood pressure"},
}

func termExample(target string) [2]string {
	if e, ok := termExamples[target]; ok {
		return e
	}
	return termExamples["ru"]
}

func termsSystem(target string) string {
	name := LanguageName(target)
	ex := termExample(target)
	return "You keep a glossary for the live translation of a conversation into " + name + ". From the topic and the recent lines with their " +
		"machine translation, pick the key terms of this field that the conversation uses: words and short phrases that a specialist " +
		"would translate in a particular way, not everyday words. Skip terms that are already in the glossary. For each term give the " +
		name + " term a specialist in this field would use. Do not copy the machine translation: it is often a literal guess, so think how " +
		"the term is really called in " + name + " in this field. " +
		"Write one entry per line: the term, an equals sign and its translation, for example: " + ex[0] + " = " + ex[1] + ". " +
		"Write at most six entries. If there is nothing to add, write none."
}

func (t *Topic) TermsMessages(now time.Time) []Message {
	var recent []pairAt
	for _, p := range t.pairs {
		if !p.at.After(t.asked) && p.at.After(t.asked.Add(-TopicEvery)) {
			recent = append(recent, p)
		}
	}
	if len(recent) == 0 {
		for _, p := range t.pairs {
			if !p.at.After(t.asked) {
				recent = append(recent, p)
			}
		}
	}
	if len(recent) > TermsPairs {
		recent = recent[len(recent)-TermsPairs:]
	}
	gl := "empty"
	if len(t.terms) > 0 {
		parts := make([]string, len(t.terms))
		for i, e := range t.terms {
			parts[i] = e.Source + " = " + e.Target
		}
		gl = strings.Join(parts, "; ")
	}
	lines := make([]string, len(recent))
	for i, p := range recent {
		lines[i] = "- " + oneLine(p.source) + " → " + oneLine(p.target)
	}
	user := "Topic: " + t.text + "\nGlossary: " + gl + "\nRecent lines:\n" + strings.Join(lines, "\n")
	return []Message{{Role: "system", Content: termsSystem(t.target)}, {Role: "user", Content: user}}
}

var stopTerms = func() map[string]bool {
	m := map[string]bool{}
	for _, w := range strings.Fields("the a an of to and but so because which that in on at for with from by if or is are was were be been do does did have has had get got make " +
		"go it this these those we you i they he she our my your their there here what how when where who very just like can will would " +
		"one some all any thing things way kind sort lot") {
		m[w] = true
	}
	return m
}()

var spaceOrDash = regexp.MustCompile(`[\s-]+`)

func heardLoosely(term, heard string) bool {
	parts := spaceOrDash.Split(strings.TrimSuffix(strings.ToLower(strings.TrimSpace(term)), "s"), -1)
	for i, p := range parts {
		parts[i] = regexp.QuoteMeta(p)
	}
	re, err := regexp.Compile(`(^|[^\pL\pN])` + strings.Join(parts, `[\s-]+`) + `s?([^\pL\pN]|$)`)
	if err != nil {
		return false
	}
	return re.MatchString(heard)
}

func termKey(term string) string {
	return strings.TrimSuffix(strings.Join(spaceOrDash.Split(strings.ToLower(strings.TrimSpace(term)), -1), " "), "s")
}

func plainTerm(term string) bool {
	if len([]rune(term)) <= 2 {
		return true
	}
	for _, w := range strings.Fields(strings.ToLower(term)) {
		if !stopTerms[w] {
			return false
		}
	}
	return true
}

func (t *Topic) TakeTerms(raw string, now time.Time) int {
	t.Counts.TermsCalls++
	defer t.finish(now)
	text := strings.TrimSpace(raw)
	if text == "none" || text == "" {
		return 0
	}
	known := map[string]bool{}
	for _, e := range t.terms {
		known[termKey(e.Source)] = true
	}
	ex := termExample(t.target)
	heard := t.heard.String()
	added := 0
	for _, line := range strings.Split(text, "\n") {
		src, dst, ok := strings.Cut(strings.TrimSpace(line), " = ")
		if i := strings.Index(dst, " ("); i > 0 {
			dst = dst[:i]
		}
		src, dst = strings.TrimSpace(src), strings.TrimRight(strings.TrimSpace(dst), ",;.")
		if !ok || src == "" || dst == "" {
			t.Counts.TermsBad++
			continue
		}
		key := termKey(src)
		switch {
		case strings.EqualFold(src, ex[0]) || strings.Contains(strings.ToLower(dst), strings.ToLower(ex[1])):
			t.Counts.Example++
		case plainTerm(src):
			t.Counts.Plain++
		case known[key]:
			t.Counts.Known++
		case !heardLoosely(src, heard):
			t.Counts.Unheard++
		case len(t.terms) >= TermsMax:
			t.Counts.Full++
		default:
			known[key] = true
			t.terms = append(t.terms, Term{Source: src, Target: dst})
			added++
		}
	}
	t.Counts.TermsAdded += added
	return added
}

func (t *Topic) Context() string {
	var b strings.Builder
	if t.text != "" {
		b.WriteString(" The conversation: ")
		b.WriteString(t.text)
	}
	if len(t.terms) > 0 {
		b.WriteString(" Translate these terms this way: ")
		for i, e := range t.terms {
			if i > 0 {
				b.WriteString("; ")
			}
			b.WriteString(e.Source)
			b.WriteString(" = ")
			b.WriteString(e.Target)
		}
		b.WriteString(".")
	}
	return b.String()
}
