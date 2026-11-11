package subtitles

import (
	"strings"
	"time"
	"unicode"
)

const (
	PhraseMinSentence = 4
	PhraseSoft        = 10
	PhraseHard        = 18
	PhraseMinClause   = 5
	PhraseKeep        = 3
	PhraseCarry       = 1500 * time.Millisecond
	PhraseQuiet       = 1000 * time.Millisecond
	PhrasePause       = 2 * time.Second
)

type TextPhrase struct {
	Seq   int
	Text  string
	Words int
	Start time.Time
	Heard time.Time
	Para  bool
}

type heardWord struct {
	text string
	at   time.Time
}

type Phraser struct {
	cur     int
	done    int
	words   []heardWord
	given   int
	lastNew time.Time
	carry   []heardWord
	carryAt time.Time
	prevAt  time.Time
	seq     int
}

func NewPhraser() *Phraser { return &Phraser{cur: -1, done: -1} }

func isCloser(r rune) bool {
	return strings.ContainsRune("\"')]}»”’›", r)
}

func isDash(w string) bool {
	return w == "-" || w == "–" || w == "—"
}

func EndsSentence(w string) bool {
	core := strings.TrimRightFunc(w, isCloser)
	if core == "" {
		return false
	}
	body := strings.TrimRightFunc(core, func(r rune) bool { return r == '.' || r == '!' || r == '?' || r == '…' })
	if body == core {
		return false
	}
	var kept []rune
	for _, r := range body {
		if r == '.' || unicode.IsLetter(r) || unicode.IsDigit(r) {
			kept = append(kept, r)
		}
	}
	if strings.ContainsRune(string(kept), '.') {
		return false
	}
	if !strings.HasSuffix(core, ".") || strings.HasSuffix(core, "..") {
		return true
	}
	if len(kept) > 2 {
		return true
	}
	if len(kept) == 0 {
		return false
	}
	for _, r := range kept {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func EndsClause(w string) bool {
	if isDash(w) {
		return true
	}
	core := strings.TrimRightFunc(w, isCloser)
	return strings.HasSuffix(core, ",") || strings.HasSuffix(core, ";") || strings.HasSuffix(core, ":")
}

func phraseCut(ws []heardWord) int {
	n := len(ws)
	for i := n - 1; i >= PhraseMinSentence-1; i-- {
		if EndsSentence(ws[i].text) {
			return i + 1
		}
	}
	if n >= PhraseSoft {
		for i := n - 1; i >= PhraseMinClause-1; i-- {
			if EndsClause(ws[i].text) {
				return i + 1
			}
		}
	}
	if n >= PhraseHard {
		return n - PhraseKeep
	}
	return 0
}

func (p *Phraser) emit(ws []heardWord, now time.Time, out []TextPhrase) []TextPhrase {
	if len(ws) == 0 {
		return out
	}
	parts := make([]string, len(ws))
	for i, w := range ws {
		parts[i] = w.text
	}
	para := p.prevAt.IsZero() || ws[0].at.Sub(p.prevAt) >= PhrasePause
	p.prevAt = ws[len(ws)-1].at
	p.seq++
	return append(out, TextPhrase{
		Seq:   p.seq,
		Text:  strings.Join(parts, " "),
		Words: len(ws),
		Start: ws[0].at,
		Heard: now,
		Para:  para,
	})
}

func (p *Phraser) cutAll(pend []heardWord, now time.Time, out []TextPhrase) ([]heardWord, []TextPhrase) {
	n := 0
	for n < len(pend) {
		n++
		for {
			c := phraseCut(pend[:n])
			if c == 0 {
				break
			}
			out = p.emit(pend[:c], now, out)
			pend = pend[c:]
			n -= c
		}
	}
	return pend, out
}

func (p *Phraser) usable(final bool) int {
	n := len(p.words)
	if !final {
		n--
	}
	if n < p.given {
		return p.given
	}
	return n
}

func (p *Phraser) advance(final bool, now time.Time, out []TextPhrase) []TextPhrase {
	use := p.usable(final)
	pend := append(append([]heardWord(nil), p.carry...), p.words[p.given:use]...)
	total := len(pend)
	left, out := p.cutAll(pend, now, out)
	took := total - len(left)
	fromCarry := took
	if fromCarry > len(p.carry) {
		fromCarry = len(p.carry)
	}
	p.carry = append([]heardWord(nil), p.carry[fromCarry:]...)
	p.given += took - fromCarry
	return out
}

func (p *Phraser) observe(fields []string, now time.Time) {
	for i, f := range fields {
		if i < len(p.words) {
			p.words[i].text = f
			continue
		}
		p.words = append(p.words, heardWord{text: f, at: now})
		p.lastNew = now
	}
	if len(fields) < len(p.words) {
		p.words = p.words[:len(fields)]
	}
	if p.given > len(p.words) {
		p.given = len(p.words)
	}
}

func (p *Phraser) closeSegment(now time.Time, out []TextPhrase) []TextPhrase {
	out = p.advance(true, now, out)
	left := append(append([]heardWord(nil), p.carry...), p.words[p.given:]...)
	p.carry = nil
	if len(left) > 0 {
		if EndsSentence(left[len(left)-1].text) || now.Sub(p.lastNew) >= PhraseQuiet {
			out = p.emit(left, now, out)
		} else {
			p.carry = left
			p.carryAt = now
		}
	}
	if p.cur > p.done {
		p.done = p.cur
	}
	p.cur = -1
	p.words = nil
	p.given = 0
	return out
}

func (p *Phraser) Update(seg int, text string, final bool, now time.Time) []TextPhrase {
	if seg <= p.done {
		return nil
	}
	var out []TextPhrase
	if p.cur >= 0 && seg != p.cur {
		out = p.closeSegment(now, out)
	}
	p.cur = seg
	if fields := strings.Fields(text); len(fields) > 0 {
		p.observe(fields, now)
	}
	if final {
		return p.closeSegment(now, out)
	}
	return p.advance(false, now, out)
}

func (p *Phraser) Tick(now time.Time) []TextPhrase {
	if len(p.carry) == 0 || len(p.words) > 0 || now.Sub(p.carryAt) < PhraseCarry {
		return nil
	}
	left := p.carry
	p.carry = nil
	return p.emit(left, now, nil)
}

func (p *Phraser) Whole(text string, from, to, now time.Time) []TextPhrase {
	fields := strings.Fields(text)
	if len(fields) == 0 {
		return nil
	}
	ws := make([]heardWord, len(fields))
	for i, f := range fields {
		ws[i] = heardWord{text: f, at: from}
	}
	left, out := p.cutAll(ws, now, nil)
	out = p.emit(left, now, out)
	if !to.Before(from) {
		p.prevAt = to
	}
	return out
}

func (p *Phraser) Flush(now time.Time) []TextPhrase {
	pend := append(append([]heardWord(nil), p.carry...), p.words[p.given:]...)
	left, out := p.cutAll(pend, now, nil)
	out = p.emit(left, now, out)
	p.carry = nil
	p.words = nil
	p.given = 0
	p.cur = -1
	p.done = -1
	return out
}

func (p *Phraser) Waiting() int {
	return len(p.carry) + len(p.words) - p.given
}
