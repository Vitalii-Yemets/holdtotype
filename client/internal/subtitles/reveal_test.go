package subtitles

import (
	"strings"
	"testing"
	"time"
)

func visibleText(rows []Row) string {
	parts := make([]string, 0, len(rows))
	for _, r := range rows {
		parts = append(parts, r.Text)
	}
	return strings.Join(parts, "\n")
}

func done(seq int, para bool, text string) PieceView {
	return PieceView{Seq: seq, Para: para, State: PieceDone, Text: text}
}

func TestShortPieceAppearsOverItsLengthByThirty(t *testing.T) {
	text := "ab cd ef gh ij kl mn op qr st"
	if len([]rune(text)) != 29 {
		t.Fatalf("length %d", len([]rune(text)))
	}
	v := NewRevealer()
	t0 := time.Date(2026, 11, 11, 10, 0, 0, 0, time.UTC)
	views := []PieceView{done(1, true, text)}
	rows, busy := v.Rows(views, false, t0)
	if !busy || visibleText(rows) != "…" {
		t.Fatalf("at once: %q %v", visibleText(rows), busy)
	}
	rows, _ = v.Rows(views, false, t0.Add(500*time.Millisecond))
	if got := visibleText(rows); got != "ab cd ef gh ij …" {
		t.Fatalf("after half a second %q", got)
	}
	rows, _ = v.Rows(views, false, t0.Add(900*time.Millisecond))
	if got := visibleText(rows); got != "ab cd ef gh ij kl mn op qr …" {
		t.Fatalf("after 0.9 s %q", got)
	}
	rows, busy = v.Rows(views, false, t0.Add(1000*time.Millisecond))
	if busy || visibleText(rows) != text {
		t.Fatalf("after 29/30 s: %q %v", visibleText(rows), busy)
	}
	rows, busy = v.Rows(views, false, t0.Add(5*time.Second))
	if busy || visibleText(rows) != text {
		t.Fatalf("later %q %v", visibleText(rows), busy)
	}
}

func revealedAfter(v *Revealer, views []PieceView, t0 time.Time, d time.Duration) int {
	v.Rows(views, false, t0)
	rows, _ := v.Rows(views, false, t0.Add(d))
	return len([]rune(strings.TrimSuffix(visibleText(rows), " …")))
}

func TestBacklogOverSixtyDoublesThePace(t *testing.T) {
	t0 := time.Date(2026, 11, 11, 10, 0, 0, 0, time.UTC)
	short := []PieceView{done(1, true, strings.Repeat("ab ", 20)[:59])}
	long := []PieceView{done(1, true, strings.Repeat("ab ", 30)[:89])}
	slow := revealedAfter(NewRevealer(), short, t0, time.Second)
	fast := revealedAfter(NewRevealer(), long, t0, time.Second)
	if slow < 26 || slow > 30 {
		t.Fatalf("59 characters in a second: %d shown", slow)
	}
	if fast < 56 || fast > 60 {
		t.Fatalf("89 characters in a second: %d shown", fast)
	}
}

func TestBacklogOverTwoFortyShowsEverything(t *testing.T) {
	t0 := time.Date(2026, 11, 11, 10, 0, 0, 0, time.UTC)
	views := []PieceView{done(1, true, strings.Repeat("word ", 30)), done(2, false, strings.Repeat("next ", 20))}
	v := NewRevealer()
	rows, busy := v.Rows(views, false, t0)
	want := strings.TrimSpace(strings.Repeat("word ", 30)) + " " + strings.TrimSpace(strings.Repeat("next ", 20))
	if busy || visibleText(rows) != want {
		t.Fatalf("%q %v", visibleText(rows), busy)
	}
	views = append(views, done(3, false, "one two three"))
	rows, busy = v.Rows(views, false, t0.Add(10*time.Millisecond))
	if !busy || !strings.HasSuffix(visibleText(rows), "next …") {
		t.Fatalf("a short piece after that is revealed at the normal pace: %q", visibleText(rows))
	}
}

func TestWordsAreNeverCut(t *testing.T) {
	text := "Здравствуйте коллеги сегодня мы посмотрим на результаты прошлой недели"
	views := []PieceView{done(1, true, text)}
	words := map[string]bool{}
	for _, w := range strings.Fields(text) {
		words[w] = true
	}
	v := NewRevealer()
	t0 := time.Date(2026, 11, 11, 10, 0, 0, 0, time.UTC)
	for ms := 0; ms <= 3000; ms += 50 {
		rows, _ := v.Rows(views, false, t0.Add(time.Duration(ms)*time.Millisecond))
		got := strings.TrimSuffix(visibleText(rows), " …")
		if got == "…" {
			continue
		}
		if !strings.HasPrefix(text, got) {
			t.Fatalf("%d ms: %q is not a prefix", ms, got)
		}
		for _, w := range strings.Fields(got) {
			if !words[w] {
				t.Fatalf("%d ms: a cut word in %q", ms, got)
			}
		}
	}
}

func TestAStreamedPieceKeepsItsVisiblePart(t *testing.T) {
	t0 := time.Date(2026, 11, 11, 10, 0, 0, 0, time.UTC)
	v := NewRevealer()
	run := func(text string) []PieceView {
		return []PieceView{{Seq: 1, Para: true, State: PieceRunning, Text: text}}
	}
	v.Rows(run("Добрый день, кол"), false, t0)
	rows, busy := v.Rows(run("Добрый день, кол"), false, t0.Add(time.Second))
	if !busy || visibleText(rows) != "Добрый день, …" {
		t.Fatalf("the unfinished last word waits: %q", visibleText(rows))
	}
	rows, _ = v.Rows(run("Добрый день, коллеги, сегодня"), false, t0.Add(1100*time.Millisecond))
	if got := visibleText(rows); !strings.HasPrefix(got, "Добрый день,") {
		t.Fatalf("the visible part was lost: %q", got)
	}
	rows, busy = v.Rows([]PieceView{done(1, true, "Добрый день, коллеги, сегодня")}, false, t0.Add(3*time.Second))
	if busy || visibleText(rows) != "Добрый день, коллеги, сегодня" {
		t.Fatalf("complete %q %v", visibleText(rows), busy)
	}
}

func TestAPieceShownAsHeardReplacesItsVisiblePart(t *testing.T) {
	t0 := time.Date(2026, 11, 11, 10, 0, 0, 0, time.UTC)
	v := NewRevealer()
	run := []PieceView{{Seq: 1, Para: true, State: PieceRunning, Text: "Это перевод который пошёл не туда и"}}
	v.Rows(run, false, t0)
	rows, _ := v.Rows(run, false, t0.Add(500*time.Millisecond))
	before := strings.TrimSuffix(visibleText(rows), " …")
	if before == "" || before == "…" {
		t.Fatalf("nothing shown before the failure: %q", before)
	}
	failed := []PieceView{{Seq: 1, Para: true, State: PieceFailed, Text: "This is what was heard on the call"}}
	rows, _ = v.Rows(failed, false, t0.Add(510*time.Millisecond))
	if len(rows) != 1 || rows[0].Kind != RowHeard {
		t.Fatalf("rows %+v", rows)
	}
	got := strings.TrimSuffix(rows[0].Text, " …")
	if !strings.HasPrefix("This is what was heard on the call", got) || len([]rune(got)) > len([]rune(before)) || strings.Contains(got, "Это") {
		t.Fatalf("the heard text did not replace the visible part: %q after %q", got, before)
	}
	rows, busy := v.Rows(failed, false, t0.Add(3*time.Second))
	if busy || len(rows) != 1 || rows[0].Text != "This is what was heard on the call" || rows[0].Kind != RowHeard {
		t.Fatalf("the heard text in the end %+v", rows)
	}
}

func TestEllipsisWhileSomethingIsHiddenWaitingOrHeard(t *testing.T) {
	t0 := time.Date(2026, 11, 11, 10, 0, 0, 0, time.UTC)
	v := NewRevealer()
	views := []PieceView{done(1, true, "one two")}
	v.Rows(views, false, t0)
	rows, busy := v.Rows(views, false, t0.Add(time.Second))
	if busy || visibleText(rows) != "one two" {
		t.Fatalf("all shown %q", visibleText(rows))
	}
	rows, _ = v.Rows(views, true, t0.Add(1100*time.Millisecond))
	if visibleText(rows) != "one two …" {
		t.Fatalf("hearing %q", visibleText(rows))
	}
	views = append(views, PieceView{Seq: 2, State: PieceWaiting, Text: ""})
	rows, busy = v.Rows(views, false, t0.Add(1200*time.Millisecond))
	if busy || visibleText(rows) != "one two …" {
		t.Fatalf("waiting %q %v", visibleText(rows), busy)
	}
}

func TestRevealedRowsKeepParagraphsAndKeys(t *testing.T) {
	t0 := time.Date(2026, 11, 11, 10, 0, 0, 0, time.UTC)
	v := NewRevealer()
	views := []PieceView{done(4, true, "a b"), done(5, false, "c d"), done(6, true, "e f"), {Seq: 7, State: PieceFailed, Text: "g h"}}
	v.Rows(views, false, t0)
	rows, _ := v.Rows(views, false, t0.Add(10*time.Second))
	want := []Row{{Kind: RowPiece, Text: "a b c d", Key: 4}, {Kind: RowPiece, Text: "e f", Key: 6}, {Kind: RowHeard, Text: "g h", Key: 7}}
	if len(rows) != len(want) {
		t.Fatalf("rows %+v", rows)
	}
	for i := range want {
		if rows[i] != want[i] {
			t.Fatalf("rows %+v", rows)
		}
	}
	v2 := NewRevealer()
	v2.Rows(views[:1], false, t0)
	v2.Rows(nil, false, t0.Add(time.Second))
	if len(v2.shown) != 0 {
		t.Fatalf("pieces that left the board are still remembered: %v", v2.shown)
	}
}

func TestRowsCarryKeysForScrolling(t *testing.T) {
	b := NewBoard(10)
	b.Add(time.Now(), "one")
	two := b.Add(time.Now(), "two")
	b.SetLive("thr")
	rows := b.Tail(500).Rows()
	if len(rows) != 3 || rows[0].Key != 1 || rows[1].Key != two.ID || rows[2].Key != TailKey {
		t.Fatalf("line rows %+v", rows)
	}
	b.AddPiece(8, true, "x y")
	b.SetPiece(8, "x y", PieceDone)
	b.AddPiece(9, false, "z")
	b.SetPiece(9, "z", PieceDone)
	b.AddPiece(10, true, "w")
	b.SetPiece(10, "", PieceFailed)
	tr := b.TranslationRows(500)
	if len(tr) != 2 || tr[0].Key != 8 || tr[0].Text != "x y z" || tr[1].Key != 10 || tr[1].Kind != RowHeard {
		t.Fatalf("piece rows %+v", tr)
	}
	views, hearing := b.Pieces()
	if hearing || len(views) != 3 || views[2].Text != "w" || views[2].State != PieceFailed || !views[0].Para || views[1].Para {
		t.Fatalf("views %+v", views)
	}
}
