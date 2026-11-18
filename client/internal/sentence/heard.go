package sentence

import (
	"strings"
	"time"

	"holdtotype/internal/punct"
)

const (
	OpAdd = iota
	OpSet
)

type Op struct {
	Kind  int
	Index int
	Raw   string
	Word  string
	At    time.Time
}

type heardWord struct {
	raw   string
	index int
}

type Heard struct {
	cur   int
	done  int
	words []heardWord
	given int
	next  int
}

func NewHeard() *Heard { return &Heard{cur: -1, done: -1} }

func (h *Heard) Given() int { return h.given }

func (h *Heard) Waiting() int { return len(h.words) - h.given }

func (h *Heard) give(upto int, at time.Time, out []Op) []Op {
	for ; h.given < upto && h.given < len(h.words); h.given++ {
		w := &h.words[h.given]
		norm := punct.Norm(w.raw)
		if norm == "" {
			w.index = -1
			out = append(out, Op{Kind: OpAdd, Index: -1, Raw: w.raw, At: at})
			continue
		}
		w.index = h.next
		h.next++
		out = append(out, Op{Kind: OpAdd, Index: w.index, Raw: w.raw, Word: norm, At: at})
	}
	return out
}

func (h *Heard) observe(fields []string, at time.Time, out []Op) []Op {
	for i, f := range fields {
		if i >= len(h.words) {
			h.words = append(h.words, heardWord{raw: f, index: -1})
			continue
		}
		w := &h.words[i]
		if w.raw == f {
			continue
		}
		w.raw = f
		if i < h.given && w.index >= 0 {
			if norm := punct.Norm(f); norm != "" {
				out = append(out, Op{Kind: OpSet, Index: w.index, Raw: f, Word: norm, At: at})
			}
		}
	}
	keep := len(fields)
	if keep < h.given {
		keep = h.given
	}
	if keep < len(h.words) {
		h.words = h.words[:keep]
	}
	return out
}

func (h *Heard) closeSegment(at time.Time, out []Op) []Op {
	out = h.give(len(h.words), at, out)
	if h.cur > h.done {
		h.done = h.cur
	}
	h.cur = -1
	h.words = nil
	h.given = 0
	return out
}

func (h *Heard) Update(seg int, text string, final bool, at time.Time) []Op {
	if seg <= h.done {
		return nil
	}
	var out []Op
	if h.cur >= 0 && seg != h.cur {
		out = h.closeSegment(at, out)
	}
	h.cur = seg
	out = h.observe(strings.Fields(text), at, out)
	if final {
		return h.closeSegment(at, out)
	}
	return h.give(len(h.words)-1, at, out)
}

func (h *Heard) Whole(text string, at time.Time) []Op {
	var out []Op
	if h.cur >= 0 {
		out = h.closeSegment(at, out)
	}
	h.words = nil
	for _, f := range strings.Fields(text) {
		h.words = append(h.words, heardWord{raw: f, index: -1})
	}
	h.given = 0
	out = h.give(len(h.words), at, out)
	h.words = nil
	h.given = 0
	return out
}

func (h *Heard) Flush(at time.Time) []Op {
	out := h.give(len(h.words), at, nil)
	h.words = nil
	h.given = 0
	h.cur = -1
	h.done = -1
	return out
}
