package sentence

import (
	"testing"
	"time"
)

func TestADraftWaitsForTimeAndWords(t *testing.T) {
	t0 := time.Date(2026, 11, 18, 10, 0, 0, 0, time.Local)
	if DraftDue(5, t0, Draft{}, t0.Add(10*time.Second)) {
		t.Fatal("a draft for five words")
	}
	if DraftDue(8, t0, Draft{}, t0.Add(3900*time.Millisecond)) {
		t.Fatal("a draft before four seconds")
	}
	if !DraftDue(6, t0, Draft{}, t0.Add(4*time.Second)) {
		t.Fatal("no draft at four seconds and six words")
	}
	last := Draft{Words: 6, At: t0.Add(4 * time.Second)}
	if DraftDue(11, t0, last, t0.Add(10*time.Second)) {
		t.Fatal("an update for five new words")
	}
	if DraftDue(12, t0, last, t0.Add(6900*time.Millisecond)) {
		t.Fatal("an update sooner than three seconds")
	}
	if !DraftDue(12, t0, last, t0.Add(7*time.Second)) {
		t.Fatal("no update after six words and three seconds")
	}
}

func TestTheTimeADraftBecomesDue(t *testing.T) {
	t0 := time.Date(2026, 11, 18, 10, 0, 0, 0, time.Local)
	var times []time.Time
	for i := 0; i < 14; i++ {
		times = append(times, t0.Add(time.Duration(i)*400*time.Millisecond))
	}
	at, ok := DraftAt(times, Draft{})
	if !ok || !at.Equal(t0.Add(4*time.Second)) {
		t.Fatalf("first draft at %v", at.Sub(t0))
	}
	at, ok = DraftAt(times, Draft{Words: 6, At: t0.Add(4 * time.Second)})
	if !ok || !at.Equal(t0.Add(7*time.Second)) {
		t.Fatalf("update at %v", at.Sub(t0))
	}
	slow := []time.Time{t0, t0.Add(time.Second), t0.Add(2 * time.Second), t0.Add(3 * time.Second), t0.Add(5 * time.Second), t0.Add(9 * time.Second)}
	at, ok = DraftAt(slow, Draft{})
	if !ok || !at.Equal(t0.Add(9*time.Second)) {
		t.Fatalf("first draft waits for the sixth word: %v", at.Sub(t0))
	}
	if _, ok := DraftAt(times, Draft{Words: 9, At: t0.Add(5 * time.Second)}); ok {
		t.Fatal("an update without six new words")
	}
}

func TestTheOpenTailIsWhatADraftTranslates(t *testing.T) {
	t0 := time.Date(2026, 11, 18, 10, 0, 0, 0, time.Local)
	h := NewHeard()
	c := NewCutter(Params{Ahead: 1, Limit: 40, Context: 30}, endsAfter("stop"))
	run(t, c, h.Update(1, "we will stop now and then the rest", true, t0))
	text, n, first := c.OpenText()
	if text != "now and then the rest" || n != 5 || !first.Equal(t0) {
		t.Fatalf("open tail %q %d %v", text, n, first)
	}
	ps := c.Flush()
	if len(ps) != 1 || len(ps[0].Parts) != 5 || len(ps[0].Times) != 5 {
		t.Fatalf("pieces %+v", ps)
	}
}
