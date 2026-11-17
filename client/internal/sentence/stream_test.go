package sentence

import (
	"strings"
	"testing"
	"time"
)

func endsAfter(marks ...string) Judge {
	set := map[string]bool{}
	for _, m := range marks {
		set[m] = true
	}
	return func(words []string) ([]bool, []float64, error) {
		end := make([]bool, len(words))
		prob := make([]float64, len(words))
		for i, w := range words {
			end[i] = set[w]
			if strings.HasPrefix(w, "p") {
				prob[i] = 0.9
			}
		}
		return end, prob, nil
	}
}

func feed(t *testing.T, s *Stream, t0 time.Time, words string, gap time.Duration) ([]Close, time.Time) {
	var out []Close
	at := t0
	for _, w := range strings.Fields(words) {
		c, err := s.Add(w, at)
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, c...)
		at = at.Add(gap)
	}
	return out, at.Add(-gap)
}

func TestAnEndIsTakenOnlyWhenTheModelStillSeesItKWordsLater(t *testing.T) {
	t0 := time.Date(2026, 11, 18, 10, 0, 0, 0, time.Local)
	p := Params{Ahead: 3, Limit: 40, Context: 30}
	s := New(p, endsAfter("stop"))
	got, _ := feed(t, s, t0, "one two stop four five", 400*time.Millisecond)
	if len(got) != 0 {
		t.Fatalf("taken with two words after it: %+v", got)
	}
	got, _ = feed(t, s, t0.Add(2*time.Second), "six", 0)
	if len(got) != 1 || got[0].First != 0 || got[0].End != 2 || got[0].How != ByModel || !got[0].At.Equal(t0.Add(2*time.Second)) {
		t.Fatalf("closes %+v", got)
	}
	if s.Open() != 3 {
		t.Fatalf("open words %d", s.Open())
	}
}

func TestSilenceClosesTheOpenSentence(t *testing.T) {
	t0 := time.Date(2026, 11, 18, 10, 0, 0, 0, time.Local)
	s := New(Params{Ahead: 3, Silence: 2 * time.Second, Limit: 40, Context: 30}, endsAfter())
	_, last := feed(t, s, t0, "one two three", 400*time.Millisecond)
	if c := s.Quiet(last.Add(1900 * time.Millisecond)); len(c) != 0 {
		t.Fatalf("closed before the silence: %+v", c)
	}
	c := s.Quiet(last.Add(2 * time.Second))
	if len(c) != 1 || c[0].End != 2 || c[0].How != BySilence || !c[0].At.Equal(last.Add(2*time.Second)) {
		t.Fatalf("closes %+v", c)
	}
	got, _ := feed(t, s, last.Add(5*time.Second), "four", 0)
	if len(got) != 0 || s.Open() != 1 {
		t.Fatalf("after the silence: %+v, open %d", got, s.Open())
	}
	s2 := New(Params{Ahead: 3, Silence: 2 * time.Second, Limit: 40, Context: 30}, endsAfter())
	feed(t, s2, t0, "one two", 400*time.Millisecond)
	got, _ = feed(t, s2, t0.Add(10*time.Second), "three", 0)
	if len(got) != 1 || got[0].End != 1 || got[0].How != BySilence {
		t.Fatalf("a late word does not close the sentence before it: %+v", got)
	}
	s3 := New(Params{Ahead: 3, Limit: 40, Context: 30}, endsAfter())
	feed(t, s3, t0, "one two", 400*time.Millisecond)
	if c := s3.Quiet(t0.Add(time.Hour)); len(c) != 0 {
		t.Fatalf("closed by silence while silence is off: %+v", c)
	}
}

func TestASentenceOverTheLimitClosesAtItsLikeliestEnd(t *testing.T) {
	t0 := time.Date(2026, 11, 18, 10, 0, 0, 0, time.Local)
	s := New(Params{Ahead: 3, Limit: 10, Context: 30}, endsAfter())
	got, _ := feed(t, s, t0, "a b c d e f pause h i j", 400*time.Millisecond)
	if len(got) != 0 {
		t.Fatalf("closed at the limit: %+v", got)
	}
	got, _ = feed(t, s, t0.Add(4*time.Second), "k", 0)
	if len(got) != 1 || got[0].End != 6 || got[0].How != ByLimit || s.Open() != 4 {
		t.Fatalf("closes %+v, open %d", got, s.Open())
	}
}

func TestTheModelSeesTheContextBeforeTheOpenSentence(t *testing.T) {
	t0 := time.Date(2026, 11, 18, 10, 0, 0, 0, time.Local)
	var seen [][]string
	judge := func(words []string) ([]bool, []float64, error) {
		seen = append(seen, append([]string(nil), words...))
		return endsAfter("stop")(words)
	}
	s := New(Params{Ahead: 1, Limit: 40, Context: 2}, judge)
	feed(t, s, t0, "a b c stop e f", 400*time.Millisecond)
	last := seen[len(seen)-1]
	if strings.Join(last, " ") != "c stop e f" {
		t.Fatalf("the model saw %q", last)
	}
	if c := s.Flush(); len(c) != 1 || c[0].First != 4 || c[0].End != 5 || c[0].How != ByEnd {
		t.Fatalf("flush %+v", c)
	}
	if c := s.Flush(); len(c) != 0 {
		t.Fatalf("a second flush %+v", c)
	}
}

func TestSeveralEndsInOneTailAreTakenInOrder(t *testing.T) {
	t0 := time.Date(2026, 11, 18, 10, 0, 0, 0, time.Local)
	calls := 0
	judge := func(words []string) ([]bool, []float64, error) {
		calls++
		return endsAfter("x", "y")(words)
	}
	s := New(Params{Ahead: 0, Limit: 40, Context: 30}, judge)
	got, _ := feed(t, s, t0, "a x", 400*time.Millisecond)
	if len(got) != 1 || got[0].End != 1 {
		t.Fatalf("closes %+v", got)
	}
	s2 := New(Params{Ahead: 2, Limit: 40, Context: 30}, endsAfter("x", "y"))
	got, _ = feed(t, s2, t0, "a x y b c", 400*time.Millisecond)
	if len(got) != 2 || got[0].End != 1 || got[1].End != 2 {
		t.Fatalf("closes %+v", got)
	}
}
