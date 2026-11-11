package subtitles

import (
	"strings"
	"testing"
	"time"
)

func TestBoardNumbersLinesAndKeepsAFewInMemory(t *testing.T) {
	b := NewBoard(3)
	t0 := time.Date(2026, 11, 10, 14, 5, 41, 0, time.Local)
	for i := 0; i < 5; i++ {
		ln := b.Add(t0.Add(time.Duration(i)*time.Second), "line")
		if ln.ID != i+1 {
			t.Fatalf("line %d got id %d", i, ln.ID)
		}
	}
	lines := b.Lines()
	if len(lines) != 3 || lines[0].ID != 3 || lines[2].ID != 5 {
		t.Fatalf("kept %+v", lines)
	}
	if b.Count() != 5 {
		t.Fatalf("count %d", b.Count())
	}
	if !lines[0].Start.Equal(t0.Add(2 * time.Second)) {
		t.Fatal("the start time of a line was lost")
	}
}

func TestTailTakesTheNewestThatFit(t *testing.T) {
	b := NewBoard(10)
	now := time.Now()
	b.Add(now, "first line here")
	b.Add(now, "second line")
	b.Add(now, "third")
	b.SetLive("  live   words ")
	v := b.Tail(26)
	if v.Live != "live words" {
		t.Fatalf("live %q", v.Live)
	}
	if len(v.Lines) != 2 || v.Lines[0].Text != "second line" || v.Lines[1].Text != "third" {
		t.Fatalf("tail %+v", v.Lines)
	}
}

func TestTailCutsAnOverlongNewestLine(t *testing.T) {
	b := NewBoard(10)
	b.Add(time.Now(), "older")
	b.Add(time.Now(), strings.Repeat("word ", 20))
	v := b.Tail(30)
	if len(v.Lines) != 1 {
		t.Fatalf("tail %+v", v.Lines)
	}
	if n := len([]rune(v.Lines[0].Text)); n > 30 || !strings.HasPrefix(v.Lines[0].Text, "…") {
		t.Fatalf("cut line %q (%d runes)", v.Lines[0].Text, n)
	}
	if v.Lines[0].ID != 2 {
		t.Fatal("the cut line lost its id")
	}
}

func TestTailOfAnEmptyBoard(t *testing.T) {
	v := NewBoard(3).Tail(100)
	if len(v.Lines) != 0 || v.Live != "" {
		t.Fatalf("%+v", v)
	}
}

func TestLiveSegmentsBecomeLines(t *testing.T) {
	l := NewLive()
	t0 := time.Now()
	if out := l.Update(0, "", false, t0); len(out) != 0 {
		t.Fatal("an empty message finished something")
	}
	l.Update(0, "hello", false, t0.Add(time.Second))
	l.Update(0, "hello there", false, t0.Add(2*time.Second))
	text, start := l.Current()
	if text != "hello there" || !start.Equal(t0.Add(time.Second)) {
		t.Fatalf("current %q from %v", text, start.Sub(t0))
	}
	out := l.Update(0, "hello there friend", true, t0.Add(3*time.Second))
	if len(out) != 1 || out[0].Text != "hello there friend" || !out[0].Start.Equal(t0.Add(time.Second)) {
		t.Fatalf("final gave %+v", out)
	}
	if out := l.Update(0, "hello there friend", true, t0.Add(4*time.Second)); len(out) != 0 {
		t.Fatal("a repeated final made a second line")
	}
	l.Update(1, "next", false, t0.Add(5*time.Second))
	out = l.Update(2, "", false, t0.Add(6*time.Second))
	if len(out) != 1 || out[0].Text != "next" {
		t.Fatalf("a higher segment did not close the previous one: %+v", out)
	}
	if out := l.Update(1, "late", false, t0.Add(7*time.Second)); len(out) != 0 {
		t.Fatal("a stale segment came back")
	}
	if text, _ := l.Current(); text != "" {
		t.Fatalf("a stale segment became live text %q", text)
	}
}

func TestAnEmptyFinalDoesNotCloseTheSegmentNumber(t *testing.T) {
	l := NewLive()
	t0 := time.Now()
	var lines []Finished
	step := func(seg int, text string, final bool, at time.Duration) {
		lines = append(lines, l.Update(seg, text, final, t0.Add(at))...)
	}
	step(0, "first words", false, time.Second)
	step(0, "first words here", true, 2*time.Second)
	step(1, "", false, 3*time.Second)
	step(1, "", true, 5*time.Second)
	step(1, "second", false, 8*time.Second)
	step(1, "second utterance", true, 9*time.Second)
	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2: %+v", len(lines), lines)
	}
	if lines[1].Text != "second utterance" || !lines[1].Start.Equal(t0.Add(8*time.Second)) {
		t.Fatalf("second line %+v", lines[1])
	}
}

func TestEmptyFinalsBeforeSpeechDoNotBlockIt(t *testing.T) {
	l := NewLive()
	t0 := time.Now()
	var lines []Finished
	for i := 0; i < 3; i++ {
		lines = append(lines, l.Update(0, "", true, t0.Add(time.Duration(i)*time.Second))...)
	}
	lines = append(lines, l.Update(0, "hello", false, t0.Add(4*time.Second))...)
	if text, _ := l.Current(); text != "hello" {
		t.Fatalf("live text %q after empty finals", text)
	}
	lines = append(lines, l.Update(0, "hello world", true, t0.Add(5*time.Second))...)
	if len(lines) != 1 || lines[0].Text != "hello world" {
		t.Fatalf("lines %+v", lines)
	}
}

func TestLiveFlushAfterABrokenConnection(t *testing.T) {
	l := NewLive()
	t0 := time.Now()
	l.Update(4, "half a sen", false, t0)
	out := l.Flush()
	if len(out) != 1 || out[0].Text != "half a sen" {
		t.Fatalf("flush %+v", out)
	}
	l.Update(0, "new session", false, t0.Add(time.Second))
	if text, _ := l.Current(); text != "new session" {
		t.Fatalf("numbering did not restart: %q", text)
	}
}

func TestTranscriptFormatting(t *testing.T) {
	at := time.Date(2026, 11, 10, 14, 5, 33, 0, time.Local)
	if got := FileName(at); got != "2026-11-10 14-05-33.txt" {
		t.Fatalf("file name %q", got)
	}
	want := "HoldToType subtitles · 2026-11-10 14:05 · language: auto · model: Nemotron 3.5"
	if got := Header(at, "auto", "Nemotron 3.5"); got != want {
		t.Fatalf("header %q", got)
	}
	if got := Opening(at, "auto", "Nemotron 3.5", ""); got != want+"\r\n\r\n" {
		t.Fatalf("opening %q", got)
	}
	line := at.Add(8 * time.Second)
	if got := LineText(line, " two\nlines  here "); got != "[14:05:41] two lines here" {
		t.Fatalf("line %q", got)
	}
	if got := Entry(line, "x"); got != "[14:05:41] x\r\n" {
		t.Fatalf("entry %q", got)
	}
}

func TestViewTextPutsEachLineOnItsOwnRow(t *testing.T) {
	b := NewBoard(10)
	b.Add(time.Now(), "next to the cafeteria")
	b.Add(time.Now(), "Tomorrow we will visit the old castle")
	b.SetLive("Could you send me")
	if got := b.Tail(500).Text(); got != "next to the cafeteria\nTomorrow we will visit the old castle\nCould you send me" {
		t.Fatalf("text %q", got)
	}
	b.SetLive("")
	if got := b.Tail(500).Text(); got != "next to the cafeteria\nTomorrow we will visit the old castle" {
		t.Fatalf("text without live %q", got)
	}
	if got := NewBoard(3).Tail(100).Text(); got != "" {
		t.Fatalf("empty board %q", got)
	}
}
