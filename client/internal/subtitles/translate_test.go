package subtitles

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestTranslationOnAKeptLine(t *testing.T) {
	b := NewBoard(5)
	t0 := time.Now()
	b.Add(t0, "first")
	ln := b.Add(t0, "second")
	if !b.SetTranslation(ln.ID, "  zweite \n Zeile ") {
		t.Fatal("a kept line did not take its translation")
	}
	lines := b.Lines()
	if lines[1].Tr != "zweite Zeile" || lines[0].Tr != "" {
		t.Fatalf("lines %+v", lines)
	}
}

func TestTranslationOnALineThatLeftTheBoard(t *testing.T) {
	b := NewBoard(2)
	first := b.Add(time.Now(), "first")
	b.Add(time.Now(), "second")
	b.Add(time.Now(), "third")
	if b.SetTranslation(first.ID, "erste") {
		t.Fatal("a line that left the board took a translation")
	}
	for _, ln := range b.Lines() {
		if ln.Tr != "" {
			t.Fatalf("another line took it: %+v", ln)
		}
	}
	if b.SetTranslation(99, "nothing") {
		t.Fatal("an unknown id took a translation")
	}
}

func TestRowsWithoutTranslations(t *testing.T) {
	b := NewBoard(10)
	b.Add(time.Now(), "one")
	b.Add(time.Now(), "two")
	want := []Row{{RowLine, "one"}, {RowLine, "two"}}
	if got := b.Tail(500).Rows(); !reflect.DeepEqual(got, want) {
		t.Fatalf("rows %+v", got)
	}
	b.SetLive("thr")
	want = append(want, Row{RowLive, "thr"})
	if got := b.Tail(500).Rows(); !reflect.DeepEqual(got, want) {
		t.Fatalf("rows with live text %+v", got)
	}
	if got := NewBoard(3).Tail(100).Rows(); len(got) != 0 {
		t.Fatalf("empty board %+v", got)
	}
}

func TestRowsPutEachTranslationUnderItsLine(t *testing.T) {
	b := NewBoard(10)
	one := b.Add(time.Now(), "one")
	b.Add(time.Now(), "two")
	three := b.Add(time.Now(), "three")
	b.SetTranslation(one.ID, "eins")
	b.SetTranslation(three.ID, "drei")
	b.SetLive("four")
	want := []Row{
		{RowLine, "one"}, {RowTranslation, "eins"},
		{RowLine, "two"},
		{RowLine, "three"}, {RowTranslation, "drei"},
		{RowLive, "four"},
	}
	v := b.Tail(500)
	if got := v.Rows(); !reflect.DeepEqual(got, want) {
		t.Fatalf("rows %+v", got)
	}
	if got := v.Text(); got != "one\ntwo\nthree\nfour" {
		t.Fatalf("text keeps the originals only, got %q", got)
	}
}

func TestTailCountsTranslationsInItsBudget(t *testing.T) {
	b := NewBoard(10)
	a := b.Add(time.Now(), "aaaaa")
	c := b.Add(time.Now(), "bbbbb")
	b.SetTranslation(a.ID, "AAAAA")
	b.SetTranslation(c.ID, "BBBBB")
	v := b.Tail(15)
	if len(v.Lines) != 1 || v.Lines[0].Text != "bbbbb" || v.Lines[0].Tr != "BBBBB" {
		t.Fatalf("tail %+v", v.Lines)
	}
	v = b.Tail(20)
	if len(v.Lines) != 2 {
		t.Fatalf("tail %+v", v.Lines)
	}
	b.SetLive("live!")
	if v = b.Tail(20); len(v.Lines) != 1 || v.Live != "live!" {
		t.Fatalf("tail with live text %+v", v)
	}
}

func TestTailCutsTheOriginalBeforeTheTranslation(t *testing.T) {
	b := NewBoard(10)
	ln := b.Add(time.Now(), strings.Repeat("word ", 10))
	b.SetTranslation(ln.ID, "Wort Wort")
	v := b.Tail(20)
	if len(v.Lines) != 1 || v.Lines[0].Tr != "Wort Wort" {
		t.Fatalf("tail %+v", v.Lines)
	}
	if n := len([]rune(v.Lines[0].Text)) + len([]rune(v.Lines[0].Tr)); n > 20 || !strings.HasPrefix(v.Lines[0].Text, "…") {
		t.Fatalf("cut %+v (%d runes)", v.Lines[0], n)
	}
	b.SetTranslation(ln.ID, strings.Repeat("Wort ", 10))
	v = b.Tail(20)
	if len(v.Lines) != 1 || v.Lines[0].Text != "" || len([]rune(v.Lines[0].Tr)) > 20 {
		t.Fatalf("an overlong translation %+v", v.Lines)
	}
	rows := v.Rows()
	if len(rows) != 1 || rows[0].Kind != RowTranslation {
		t.Fatalf("rows %+v", rows)
	}
	if v.Text() != "" {
		t.Fatalf("text %q", v.Text())
	}
}

func pending(id int) Pending {
	return Pending{ID: id, Text: "line", Target: "de"}
}

func TestQueueDropsTheOldestWaitingLine(t *testing.T) {
	q := NewTranslateQueue(TranslateWaiting, TranslateMaxWait)
	t0 := time.Now()
	for i := 1; i <= 3; i++ {
		if n := q.Push(pending(i), t0); n != 0 {
			t.Fatalf("push %d dropped %d", i, n)
		}
	}
	if n := q.Push(pending(4), t0); n != 1 {
		t.Fatalf("the fourth line dropped %d", n)
	}
	if q.Len() != 3 {
		t.Fatalf("len %d", q.Len())
	}
	var ids []int
	for {
		p, ok, _ := q.Pop(t0.Add(time.Second))
		if !ok {
			break
		}
		ids = append(ids, p.ID)
	}
	if !reflect.DeepEqual(ids, []int{2, 3, 4}) {
		t.Fatalf("popped %v", ids)
	}
	c := q.Counts()
	if c.Pushed != 4 || c.Overflow != 1 || c.Stale != 0 || c.Dropped() != 1 {
		t.Fatalf("counts %+v", c)
	}
}

func TestQueueDropsALineThatWaitedTooLong(t *testing.T) {
	q := NewTranslateQueue(TranslateWaiting, TranslateMaxWait)
	t0 := time.Now()
	q.Push(pending(1), t0)
	q.Push(pending(2), t0.Add(5*time.Second))
	q.Push(pending(3), t0.Add(10*time.Second))
	p, ok, stale := q.Pop(t0.Add(25*time.Second + time.Millisecond))
	if !ok || p.ID != 3 || stale != 2 {
		t.Fatalf("pop gave %+v %v, %d stale", p, ok, stale)
	}
	if !p.Queued.Equal(t0.Add(10 * time.Second)) {
		t.Fatal("the queue did not keep the time the line was pushed")
	}
	q.Push(pending(4), t0)
	if p, ok, _ := q.Pop(t0.Add(TranslateMaxWait)); !ok || p.ID != 4 {
		t.Fatal("a line that waited exactly the limit was dropped")
	}
	q.Push(pending(5), t0)
	if _, ok, stale := q.Pop(t0.Add(time.Minute)); ok || stale != 1 {
		t.Fatalf("a lone old line: ok %v, %d stale", ok, stale)
	}
	if _, ok, stale := q.Pop(t0.Add(time.Minute)); ok || stale != 0 {
		t.Fatal("an empty queue gave something")
	}
	q.Push(pending(6), t0)
	q.Push(pending(7), t0)
	if n := q.Clear(); n != 2 || q.Len() != 0 {
		t.Fatalf("clear %d, len %d", n, q.Len())
	}
	c := q.Counts()
	if c.Pushed != 7 || c.Stale != 3 || c.Cleared != 2 || c.Overflow != 0 || c.Dropped() != 5 {
		t.Fatalf("counts %+v", c)
	}
}

func TestCleanTranslation(t *testing.T) {
	cases := []struct {
		line, answer, want string
		verdict            Verdict
	}{
		{"Hello there", "Hallo da", "Hallo da", TrOK},
		{"Hello there", `"Hallo da"`, "Hallo da", TrOK},
		{"Hello there", "«Привет»", "Привет", TrOK},
		{"Hello there", "„Hallo da“", "Hallo da", TrOK},
		{"Hello there", "“Hallo da”", "Hallo da", TrOK},
		{"Hello there", "'Hallo da'", "Hallo da", TrOK},
		{"Hello there", "‘Ciao’", "Ciao", TrOK},
		{"Hello there", "「こんにちは」", "こんにちは", TrOK},
		{"Hello there", "\"« Hallo da »\"", "Hallo da", TrOK},
		{"Hello there", `"Hallo", sagte er.`, `"Hallo", sagte er.`, TrOK},
		{"Hello there", "  Hallo   \t da  ", "Hallo da", TrOK},
		{"Hello there", "\n\n  Hallo da\nNote: this is German.\n", "Hallo da", TrOK},
		{"Hello there", "\r\nHallo da\r\n(translated)", "Hallo da", TrOK},
		{"Hello there", "", "", TrEmpty},
		{"Hello there", " \n \"\" \n", "", TrEmpty},
		{"Hello there", "...", "", TrEmpty},
		{"Hello, there!", "hello there", "", TrSame},
		{"Hello there", "  HELLO,   THERE.  ", "", TrSame},
		{"Привет, мир", "«привет мир!»", "", TrSame},
		{"Hi", strings.Repeat("x", 49), "", TrRamble},
		{"Hi", strings.Repeat("x", 48), strings.Repeat("x", 48), TrOK},
		{"Hi", "\"" + strings.Repeat("x", 48) + "\"", strings.Repeat("x", 48), TrOK},
	}
	for _, c := range cases {
		got, v := CleanTranslation(c.line, c.answer)
		if got != c.want || v != c.verdict {
			t.Errorf("CleanTranslation(%q, %q) = %q, %q; want %q, %q", c.line, c.answer, got, v, c.want, c.verdict)
		}
	}
}

func TestCleanTranslationRambleLimit(t *testing.T) {
	line := strings.Repeat("a", 10)
	limit := 4*10 + 40
	if got, v := CleanTranslation(line, strings.Repeat("b", limit)); v != TrOK || len(got) != limit {
		t.Fatalf("an answer at the limit was thrown away: %q", v)
	}
	if _, v := CleanTranslation(line, strings.Repeat("b", limit+1)); v != TrRamble {
		t.Fatalf("an answer over the limit gave %q", v)
	}
}

func TestTranslatePromptNamesTheLanguage(t *testing.T) {
	names := map[string]string{
		"de": "German", "en": "English", "es": "Spanish", "fr": "French",
		"it": "Italian", "pl": "Polish", "uk": "Ukrainian", "ru": "Russian",
	}
	if len(TranslateTargets) != 8 {
		t.Fatalf("targets %v", TranslateTargets)
	}
	for _, code := range TranslateTargets {
		name := names[code]
		if LanguageName(code) != name {
			t.Errorf("name of %s is %q", code, LanguageName(code))
		}
		want := "You translate subtitles. Translate the user's message into " + name +
			". Reply with the translation only: one line, no quotes, no notes, no explanations. Keep names and numbers as they are. If the message is already in " +
			name + ", repeat it unchanged."
		if got := TranslatePrompt(code); got != want {
			t.Errorf("prompt for %s: %q", code, got)
		}
	}
	for _, code := range []string{"", "auto", "xx", "DE", "zh"} {
		if got := TranslatePrompt(code); got != "" {
			t.Errorf("prompt for an unknown code %q: %q", code, got)
		}
		if LanguageName(code) != "" {
			t.Errorf("name for an unknown code %q", code)
		}
	}
}

func TestTranscriptWithTranslation(t *testing.T) {
	at := time.Date(2026, 11, 10, 14, 5, 33, 0, time.Local)
	head := "HoldToType subtitles · 2026-11-10 14:05 · language: auto · model: Nemotron 3.5"
	if got := Opening(at, "auto", "Nemotron 3.5", ""); got != head+"\r\n\r\n" {
		t.Fatalf("opening without translation %q", got)
	}
	if got := Opening(at, "auto", "Nemotron 3.5", "de"); got != head+" · translation: de\r\n\r\n" {
		t.Fatalf("opening with translation %q", got)
	}
	line := at.Add(8 * time.Second)
	if got := TranslationText(line, " zwei\nZeilen  hier "); got != "[14:05:41] → zwei Zeilen hier" {
		t.Fatalf("translation text %q", got)
	}
	if got := TranslationEntry(line, "Hallo"); got != "[14:05:41] → Hallo\r\n" {
		t.Fatalf("translation entry %q", got)
	}
}
