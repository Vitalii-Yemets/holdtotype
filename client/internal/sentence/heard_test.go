package sentence

import (
	"strings"
	"testing"
	"time"
)

func run(t *testing.T, c *Cutter, ops []Op) []Piece {
	var out []Piece
	for _, op := range ops {
		p, err := c.Apply(op)
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, p...)
	}
	return out
}

func texts(ps []Piece) []string {
	var out []string
	for _, p := range ps {
		out = append(out, p.Text)
	}
	return out
}

func TestTheLastWordWaitsUntilTheRecognizerMovesOn(t *testing.T) {
	t0 := time.Date(2026, 11, 18, 10, 0, 0, 0, time.Local)
	h := NewHeard()
	ops := h.Update(1, "So we", false, t0)
	if len(ops) != 1 || ops[0].Raw != "So" || ops[0].Word != "so" || ops[0].Index != 0 {
		t.Fatalf("ops %+v", ops)
	}
	ops = h.Update(1, "So we looked", false, t0.Add(time.Second))
	if len(ops) != 1 || ops[0].Raw != "we" || ops[0].Index != 1 {
		t.Fatalf("ops %+v", ops)
	}
	ops = h.Update(1, "So we looked at the cohorts.", true, t0.Add(2*time.Second))
	if len(ops) != 4 || ops[3].Raw != "cohorts." || ops[3].Word != "cohorts" || ops[3].Index != 5 || h.Waiting() != 0 {
		t.Fatalf("ops %+v", ops)
	}
	if ops := h.Update(1, "late text", false, t0.Add(3*time.Second)); len(ops) != 0 {
		t.Fatalf("a finished segment came back: %+v", ops)
	}
	ops = h.Update(2, "Then —", false, t0.Add(4*time.Second))
	ops = append(ops, h.Update(3, "next", false, t0.Add(5*time.Second))...)
	if len(ops) != 2 || ops[0].Raw != "Then" || ops[0].Index != 6 || ops[1].Raw != "—" || ops[1].Word != "" || ops[1].Index != -1 {
		t.Fatalf("a new segment does not finish the old one: %+v", ops)
	}
}

func TestARewrittenWordInTheOpenTailIsCorrected(t *testing.T) {
	t0 := time.Date(2026, 11, 18, 10, 0, 0, 0, time.Local)
	h := NewHeard()
	c := NewCutter(Params{Ahead: 1, Limit: 40, Context: 30}, endsAfter("stop"))
	var got []Piece
	got = append(got, run(t, c, h.Update(1, "the approval reyt went", false, t0))...)
	ops := h.Update(1, "the approval rate went up", false, t0.Add(time.Second))
	if len(ops) != 2 || ops[0].Kind != OpSet || ops[0].Index != 2 || ops[0].Raw != "rate" || ops[1].Raw != "went" {
		t.Fatalf("ops %+v", ops)
	}
	got = append(got, run(t, c, ops)...)
	got = append(got, run(t, c, h.Flush(t0.Add(2*time.Second)))...)
	got = append(got, c.Flush()...)
	if strings.Join(texts(got), "|") != "the approval rate went up" {
		t.Fatalf("pieces %q", texts(got))
	}
}

func TestAClosedSentenceIsNotOpenedAgain(t *testing.T) {
	t0 := time.Date(2026, 11, 18, 10, 0, 0, 0, time.Local)
	h := NewHeard()
	c := NewCutter(Params{Ahead: 1, Limit: 40, Context: 30}, endsAfter("stop"))
	got := run(t, c, h.Update(1, "we will stop now and", false, t0))
	if strings.Join(texts(got), "|") != "we will stop" {
		t.Fatalf("pieces %q", texts(got))
	}
	got = run(t, c, h.Update(1, "we will stop here and then", false, t0.Add(time.Second)))
	got = append(got, run(t, c, h.Flush(t0.Add(2*time.Second)))...)
	got = append(got, c.Flush()...)
	if strings.Join(texts(got), "|") != "here and then" {
		t.Fatalf("pieces %q", texts(got))
	}
}

func TestEveryWordGoesOutOnceAsHeard(t *testing.T) {
	t0 := time.Date(2026, 11, 18, 10, 0, 0, 0, time.Local)
	h := NewHeard()
	c := NewCutter(Params{Ahead: 2, Silence: 2 * time.Second, Limit: 6, Context: 30, MinWords: 3}, endsAfter("one", "stop"))
	heard := []string{"Yes, one.", "— and the next stop. We", "carry on with eight more words here today please"}
	var got []Piece
	at := t0
	for i, line := range heard {
		fields := strings.Fields(line)
		for k := 1; k <= len(fields); k++ {
			got = append(got, run(t, c, h.Update(i+1, strings.Join(fields[:k], " "), k == len(fields), at))...)
			at = at.Add(400 * time.Millisecond)
		}
		got = append(got, c.Quiet(at)...)
	}
	got = append(got, run(t, c, h.Flush(at))...)
	got = append(got, c.Flush()...)
	if strings.Join(texts(got), " ") != strings.Join(heard, " ") {
		t.Fatalf("the words came out as %q", texts(got))
	}
	for _, p := range got {
		if p.Words < 3 && p.How == ByModel {
			t.Fatalf("a model end closed a piece of %d words: %+v", p.Words, p)
		}
	}
}

func TestAShortSentenceIsClosedOnlyBySilence(t *testing.T) {
	t0 := time.Date(2026, 11, 18, 10, 0, 0, 0, time.Local)
	s := New(Params{Ahead: 1, Silence: 2 * time.Second, Limit: 40, Context: 30, MinWords: 3}, endsAfter("no", "stop"))
	got, last := feed(t, s, t0, "no we will stop here", 400*time.Millisecond)
	if len(got) != 1 || got[0].First != 0 || got[0].End != 3 {
		t.Fatalf("closes %+v", got)
	}
	s2 := New(Params{Ahead: 1, Silence: 2 * time.Second, Limit: 40, Context: 30, MinWords: 3}, endsAfter("no"))
	_, last = feed(t, s2, t0, "no", 400*time.Millisecond)
	c := s2.Quiet(last.Add(2 * time.Second))
	if len(c) != 1 || c[0].Words() != 1 || c[0].How != BySilence {
		t.Fatalf("closes %+v", c)
	}
}

func TestAbandonGivesTheOpenWordsBack(t *testing.T) {
	t0 := time.Date(2026, 11, 18, 10, 0, 0, 0, time.Local)
	c := NewCutter(Params{Ahead: 3, Limit: 40, Context: 30}, endsAfter())
	h := NewHeard()
	run(t, c, h.Update(1, "words still open here", true, t0))
	got := c.Abandon()
	if len(got) != 1 || got[0].Text != "words still open here" || got[0].How != ByFallback {
		t.Fatalf("pieces %+v", got)
	}
}
