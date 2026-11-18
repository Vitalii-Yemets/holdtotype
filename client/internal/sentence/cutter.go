package sentence

import (
	"strings"
	"time"
)

const ByFallback = "fallback"

type Piece struct {
	Text  string
	Words int
	Start time.Time
	Last  time.Time
	At    time.Time
	How   string
	Parts []string
	Times []time.Time
}

type token struct {
	pre  string
	word string
	post string
}

func (t token) text() string { return t.pre + t.word + t.post }

type Cutter struct {
	st    *Stream
	toks  []token
	times []time.Time
	lead  string
}

func NewCutter(p Params, judge Judge) *Cutter {
	return &Cutter{st: New(p, judge)}
}

func (c *Cutter) Open() int { return c.st.Open() }

func (c *Cutter) OpenText() (string, int, time.Time) {
	start := c.st.Start()
	if start >= len(c.toks) {
		return "", 0, time.Time{}
	}
	parts := make([]string, 0, len(c.toks)-start)
	for _, t := range c.toks[start:] {
		parts = append(parts, t.text())
	}
	return strings.Join(parts, " "), len(c.toks) - start, c.times[start]
}

func (c *Cutter) pieces(cl []Close) []Piece {
	out := make([]Piece, 0, len(cl))
	for _, x := range cl {
		parts := make([]string, 0, x.Words())
		for i := x.First; i <= x.End; i++ {
			parts = append(parts, c.toks[i].text())
		}
		text := strings.Join(parts, " ")
		times := append([]time.Time(nil), c.times[x.First:x.End+1]...)
		out = append(out, Piece{Text: text, Words: len(strings.Fields(text)), Start: c.times[x.First], Last: c.times[x.End], At: x.At, How: x.How, Parts: parts, Times: times})
	}
	return out
}

func (c *Cutter) Apply(op Op) ([]Piece, error) {
	switch {
	case op.Kind == OpSet:
		if c.st.Replace(op.Index, op.Word) {
			c.toks[op.Index].word = op.Raw
		}
		return nil, nil
	case op.Word == "":
		last := len(c.toks) - 1
		if last >= 0 && last >= c.st.Start() {
			c.toks[last].post += " " + op.Raw
		} else {
			c.lead += op.Raw + " "
		}
		return nil, nil
	}
	c.toks = append(c.toks, token{pre: c.lead, word: op.Raw})
	c.times = append(c.times, op.At)
	c.lead = ""
	cl, err := c.st.Add(op.Word, op.At)
	return c.pieces(cl), err
}

func (c *Cutter) Quiet(now time.Time) []Piece { return c.pieces(c.st.Quiet(now)) }

func (c *Cutter) Flush() []Piece {
	out := c.pieces(c.st.Flush())
	if strings.TrimSpace(c.lead) != "" {
		text := strings.TrimSpace(c.lead)
		at := time.Now()
		if n := len(c.times); n > 0 {
			at = c.times[n-1]
		}
		out = append(out, Piece{Text: text, Words: len(strings.Fields(text)), Start: at, Last: at, At: at, How: ByEnd, Parts: []string{text}, Times: []time.Time{at}})
	}
	c.lead = ""
	return out
}

func (c *Cutter) Abandon() []Piece {
	out := c.Flush()
	for i := range out {
		out[i].How = ByFallback
	}
	return out
}
