package subtitles

import (
	"reflect"
	"testing"
	"time"
)

func drain(q *TranslateQueue, now time.Time) ([]Batch, []int) {
	var out []Batch
	var stale []int
	for {
		b, ok, st := q.Pop(now)
		stale = append(stale, seqsOf(st)...)
		if !ok {
			return out, stale
		}
		out = append(out, b)
	}
}

func batchSeqs(bs []Batch) []int {
	var out []int
	for _, b := range bs {
		out = append(out, b.Seqs...)
	}
	return out
}

func TestPiecesHeardDuringTheStartAllGoAfterwardsInOrder(t *testing.T) {
	q := newQueue()
	t0 := time.Date(2026, 11, 12, 20, 0, 0, 0, time.UTC)
	gen, fresh := q.Hold("the editor model", t0)
	if !fresh || gen == 0 {
		t.Fatalf("hold %d %v", gen, fresh)
	}
	var want []int
	for i := 1; i <= 25; i++ {
		if d := q.Push(phrase(i, 4, t0.Add(time.Duration(i)*1200*time.Millisecond), i%5 == 1)); len(d) != 0 {
			t.Fatalf("piece %d dropped %v while the engine starts", i, seqsOf(d))
		}
		want = append(want, i)
		if b, ok, st := q.Pop(t0.Add(time.Duration(i) * 1300 * time.Millisecond)); ok || len(st) != 0 {
			t.Fatalf("a piece left while the engine starts: %+v %v", b, seqsOf(st))
		}
	}
	if held, label, since := q.Holding(); !held || label != "the editor model" || !since.Equal(t0) {
		t.Fatalf("holding %v %q %v", held, label, since)
	}
	warm := t0.Add(45 * time.Second)
	n, ok := q.Release(gen, warm)
	if !ok || n != 25 {
		t.Fatalf("release %d %v", n, ok)
	}
	bs, stale := drain(q, warm.Add(9*time.Second))
	if len(stale) != 0 {
		t.Fatalf("stale after the release %v", stale)
	}
	if got := batchSeqs(bs); !reflect.DeepEqual(got, want) {
		t.Fatalf("order %v", got)
	}
	if len(bs) >= 25 {
		t.Fatalf("not merged: %d batches", len(bs))
	}
	for _, b := range bs {
		if b.Words > TranslateMerge {
			t.Fatalf("a batch of %d words", b.Words)
		}
	}
	if c := q.Counts(); c.Dropped() != 0 || c.Pushed != 25 {
		t.Fatalf("counts %+v", c)
	}
}

func TestTheWaitAfterTheStartCountsFromTheWarmMoment(t *testing.T) {
	q := newQueue()
	t0 := time.Date(2026, 11, 12, 20, 0, 0, 0, time.UTC)
	gen, _ := q.Hold("x", t0)
	for i := 1; i <= 3; i++ {
		q.Push(phrase(i, 5, t0, true))
	}
	warm := t0.Add(time.Minute)
	q.Release(gen, warm)
	if _, ok, st := q.Pop(warm.Add(TranslateMaxWait)); !ok || len(st) != 0 {
		t.Fatal("a held piece was stale at the limit after the warm moment")
	}
	_, stale := drain(q, warm.Add(TranslateMaxWait+time.Second))
	if !reflect.DeepEqual(stale, []int{2, 3}) {
		t.Fatalf("stale %v", stale)
	}
}

func TestAtMostFortyPiecesWait(t *testing.T) {
	q := newQueue()
	t0 := time.Date(2026, 11, 12, 20, 0, 0, 0, time.UTC)
	gen, _ := q.Hold("x", t0)
	var dropped []int
	for i := 1; i <= TranslateHoldMax+5; i++ {
		dropped = append(dropped, seqsOf(q.Push(phrase(i, 2, t0.Add(time.Duration(i)*time.Second), true)))...)
	}
	if !reflect.DeepEqual(dropped, []int{1, 2, 3, 4, 5}) || q.Len() != TranslateHoldMax {
		t.Fatalf("dropped %v, %d wait", dropped, q.Len())
	}
	q.Release(gen, t0.Add(time.Minute))
	bs, _ := drain(q, t0.Add(time.Minute))
	if got := batchSeqs(bs); len(got) != TranslateHoldMax || got[0] != 6 || got[len(got)-1] != TranslateHoldMax+5 {
		t.Fatalf("popped %v", got)
	}
	if c := q.Counts(); c.Overflow != 5 {
		t.Fatalf("counts %+v", c)
	}
}

func TestHeldPiecesAreOutsideTheQueueLimitAfterTheRelease(t *testing.T) {
	q := newQueue()
	t0 := time.Date(2026, 11, 12, 20, 0, 0, 0, time.UTC)
	gen, _ := q.Hold("x", t0)
	for i := 1; i <= 20; i++ {
		q.Push(phrase(i, 2, t0, true))
	}
	q.Release(gen, t0.Add(30*time.Second))
	var dropped []int
	for i := 21; i <= 30; i++ {
		dropped = append(dropped, seqsOf(q.Push(phrase(i, 2, t0.Add(31*time.Second), true)))...)
	}
	if !reflect.DeepEqual(dropped, []int{21, 22, 23, 24}) {
		t.Fatalf("dropped %v", dropped)
	}
	bs, stale := drain(q, t0.Add(32*time.Second))
	want := []int{}
	for i := 1; i <= 20; i++ {
		want = append(want, i)
	}
	want = append(want, 25, 26, 27, 28, 29, 30)
	if got := batchSeqs(bs); !reflect.DeepEqual(got, want) || len(stale) != 0 {
		t.Fatalf("popped %v, stale %v", got, stale)
	}
}

func TestAHoldIsReleasedOnlyByItsOwnStart(t *testing.T) {
	q := newQueue()
	t0 := time.Date(2026, 11, 12, 20, 0, 0, 0, time.UTC)
	g1, fresh := q.Hold("the translation model on the graphics card", t0)
	if _, again := q.Hold("the translation model on the graphics card", t0.Add(time.Second)); again {
		t.Fatal("the same start was begun twice")
	}
	q.Push(phrase(1, 3, t0.Add(2*time.Second), true))
	g2, fresh2 := q.Hold("the editor model", t0.Add(20*time.Second))
	if !fresh || !fresh2 || g2 == g1 {
		t.Fatal("a hold for the next translator is not a new one")
	}
	if _, ok := q.Release(g1, t0.Add(21*time.Second)); ok {
		t.Fatal("the failed start released the hold of the next translator")
	}
	if held, label, since := q.Holding(); !held || label != "the editor model" || !since.Equal(t0) {
		t.Fatalf("holding %v %q %v", held, label, since)
	}
	q.Push(phrase(2, 3, t0.Add(25*time.Second), true))
	if n, ok := q.Release(g2, t0.Add(40*time.Second)); !ok || n != 2 {
		t.Fatalf("release %d %v", n, ok)
	}
	bs, _ := drain(q, t0.Add(41*time.Second))
	if got := batchSeqs(bs); !reflect.DeepEqual(got, []int{1, 2}) {
		t.Fatalf("popped %v", got)
	}
	if _, ok := q.ReleaseAll(t0); ok {
		t.Fatal("released twice")
	}
}

func TestAFailedStartHandsTheWaitingPiecesOn(t *testing.T) {
	for _, c := range []struct {
		choice string
		want   Choice
	}{
		{ChoiceAuto, Choice{Kind: KindServer}},
		{ChoiceQuality, Choice{Kind: KindQuality, Missing: MissingStart}},
	} {
		q := newQueue()
		t0 := time.Date(2026, 11, 12, 20, 0, 0, 0, time.UTC)
		gen, _ := q.Hold("the translation model on the graphics card", t0)
		for i := 1; i <= 8; i++ {
			q.Push(phrase(i, 3, t0.Add(time.Duration(i)*time.Second), false))
		}
		after := ChooseTranslator(c.choice, Avail{ModelReady: true, CardFits: true, QualityDown: true, Server: true, Editor: true})
		if after != c.want {
			t.Fatalf("%q: after the failed start %+v", c.choice, after)
		}
		q.Release(gen, t0.Add(30*time.Second))
		bs, stale := drain(q, t0.Add(31*time.Second))
		if got := batchSeqs(bs); len(stale) != 0 || !reflect.DeepEqual(got, []int{1, 2, 3, 4, 5, 6, 7, 8}) {
			t.Fatalf("%q: popped %v, stale %v", c.choice, got, stale)
		}
	}
}

func TestWithoutAStartTheQueueIsAsBefore(t *testing.T) {
	q := newQueue()
	t0 := time.Date(2026, 11, 12, 20, 0, 0, 0, time.UTC)
	var dropped []int
	for i := 1; i <= 9; i++ {
		dropped = append(dropped, seqsOf(q.Push(phrase(i, 5, t0, true)))...)
	}
	if !reflect.DeepEqual(dropped, []int{1, 2, 3}) {
		t.Fatalf("dropped %v", dropped)
	}
	if held, _, _ := q.Holding(); held {
		t.Fatal("holding without a start")
	}
	_, stale := drain(q, t0.Add(TranslateMaxWait+time.Second))
	if !reflect.DeepEqual(stale, []int{4, 5, 6, 7, 8, 9}) {
		t.Fatalf("stale %v", stale)
	}
	if TranslateHoldMax != 40 || TranslateWarmWait != 90*time.Second {
		t.Fatal("constants")
	}
}
