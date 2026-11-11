package subtitles

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestRowsOfLinesAndLiveText(t *testing.T) {
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

func TestPiecesOfOneParagraphFlowAsOneRow(t *testing.T) {
	b := NewBoard(50)
	b.AddPiece(1, true, "Good morning everyone,")
	b.AddPiece(2, false, "thanks for joining.")
	b.SetPiece(1, "Доброе утро всем,", PieceDone)
	b.SetPiece(2, "спасибо, что", PieceRunning)
	want := []Row{{RowPiece, "Доброе утро всем, спасибо, что"}}
	if got := b.TranslationRows(500); !reflect.DeepEqual(got, want) {
		t.Fatalf("rows %+v", got)
	}
	b.SetPiece(2, "спасибо, что пришли.", PieceDone)
	b.AddPiece(3, true, "Let us start.")
	b.SetPiece(3, "Начнём.", PieceDone)
	want = []Row{{RowPiece, "Доброе утро всем, спасибо, что пришли."}, {RowPiece, "Начнём."}}
	if got := b.TranslationRows(500); !reflect.DeepEqual(got, want) {
		t.Fatalf("a paragraph did not start a row: %+v", got)
	}
}

func TestPendingSpeechEndsTheTextWithAnEllipsis(t *testing.T) {
	b := NewBoard(50)
	if got := b.TranslationRows(500); len(got) != 0 {
		t.Fatalf("an empty board %+v", got)
	}
	b.SetHearing(true)
	if got := b.TranslationRows(500); !reflect.DeepEqual(got, []Row{{RowPiece, "…"}}) {
		t.Fatalf("hearing with nothing shown %+v", got)
	}
	b.SetHearing(false)
	b.AddPiece(1, true, "first piece")
	if got := b.TranslationRows(500); !reflect.DeepEqual(got, []Row{{RowPiece, "…"}}) {
		t.Fatalf("a waiting piece %+v", got)
	}
	b.SetPiece(1, "", PieceRunning)
	if got := b.TranslationRows(500); !reflect.DeepEqual(got, []Row{{RowPiece, "…"}}) {
		t.Fatalf("a running piece without text %+v", got)
	}
	b.SetPiece(1, "erstes", PieceRunning)
	if got := b.TranslationRows(500); !reflect.DeepEqual(got, []Row{{RowPiece, "erstes"}}) {
		t.Fatalf("a running piece with text %+v", got)
	}
	b.SetPiece(1, "erstes Stück", PieceDone)
	b.AddPiece(2, false, "second piece")
	if got := b.TranslationRows(500); !reflect.DeepEqual(got, []Row{{RowPiece, "erstes Stück …"}}) {
		t.Fatalf("a done piece and a waiting one %+v", got)
	}
	b.SetPiece(2, "zweites", PieceDone)
	b.SetHearing(true)
	if got := b.TranslationRows(500); !reflect.DeepEqual(got, []Row{{RowPiece, "erstes Stück zweites …"}}) {
		t.Fatalf("hearing after done pieces %+v", got)
	}
}

func TestFailedPieceIsShownAsHeardOnItsOwnRow(t *testing.T) {
	b := NewBoard(50)
	b.AddPiece(1, true, "one two")
	b.AddPiece(2, false, "three  four\nfive")
	b.AddPiece(3, false, "six seven")
	b.SetPiece(1, "eins zwei", PieceDone)
	b.SetPiece(2, "", PieceFailed)
	b.SetPiece(3, "sechs sieben", PieceDone)
	want := []Row{{RowPiece, "eins zwei"}, {RowHeard, "three four five"}, {RowPiece, "sechs sieben"}}
	if got := b.TranslationRows(500); !reflect.DeepEqual(got, want) {
		t.Fatalf("rows %+v", got)
	}
	b.AddPiece(4, false, "eight")
	b.SetPiece(4, "", PieceFailed)
	b.SetHearing(true)
	want = append(want, Row{RowHeard, "eight …"})
	if got := b.TranslationRows(500); !reflect.DeepEqual(got, want) {
		t.Fatalf("rows ending with a failed piece %+v", got)
	}
}

func TestJoinedPiecesBecomeOne(t *testing.T) {
	b := NewBoard(50)
	b.AddPiece(1, true, "alpha")
	b.AddPiece(2, false, "beta")
	b.AddPiece(3, false, "gamma")
	b.AddPiece(4, false, "delta")
	b.JoinPieces([]int{2, 3})
	b.SetPiece(1, "A", PieceDone)
	b.SetPiece(3, "lost", PieceDone)
	b.SetPiece(2, "", PieceFailed)
	b.SetPiece(4, "D", PieceDone)
	want := []Row{{RowPiece, "A"}, {RowHeard, "beta gamma"}, {RowPiece, "D"}}
	if got := b.TranslationRows(500); !reflect.DeepEqual(got, want) {
		t.Fatalf("rows %+v", got)
	}
	if b.SetPiece(3, "x", PieceDone) || b.SetPiece(99, "x", PieceDone) {
		t.Fatal("a joined or unknown piece took a text")
	}
}

func TestTranslationRowsKeepTheTailThatFits(t *testing.T) {
	b := NewBoard(50)
	b.AddPiece(1, true, "a")
	b.SetPiece(1, strings.Repeat("alt ", 10), PieceDone)
	b.AddPiece(2, true, "b")
	b.SetPiece(2, "neu eins zwei", PieceDone)
	rows := b.TranslationRows(20)
	if len(rows) != 2 || rows[1].Text != "neu eins zwei" || !strings.HasPrefix(rows[0].Text, "…") {
		t.Fatalf("rows %+v", rows)
	}
	n := 0
	for _, r := range rows {
		n += len([]rune(r.Text))
	}
	if n > 20 {
		t.Fatalf("%d runes in %+v", n, rows)
	}
	if got := b.TranslationRows(13); !reflect.DeepEqual(got, []Row{{RowPiece, "neu eins zwei"}}) {
		t.Fatalf("an exact fit %+v", got)
	}
	if got := b.TranslationRows(0); len(got) != 0 {
		t.Fatalf("no room %+v", got)
	}
}

func TestBoardKeepsAFewPieces(t *testing.T) {
	b := NewBoard(3)
	for i := 1; i <= 5; i++ {
		b.AddPiece(i, true, "x")
		b.SetPiece(i, string(rune('a'+i)), PieceDone)
	}
	want := []Row{{RowPiece, "d"}, {RowPiece, "e"}, {RowPiece, "f"}}
	if got := b.TranslationRows(500); !reflect.DeepEqual(got, want) {
		t.Fatalf("rows %+v", got)
	}
	if b.SetPiece(1, "y", PieceDone) {
		t.Fatal("a piece that left the board took a text")
	}
}

func phrase(seq, words int, heard time.Time, para bool) Pending {
	ws := make([]string, words)
	for i := range ws {
		ws[i] = fmt.Sprintf("w%d.%d", seq, i)
	}
	return Pending{Seq: seq, Text: strings.Join(ws, " "), Words: words, Target: "de", Para: para, Start: heard.Add(-3 * time.Second), Heard: heard}
}

func seqsOf(ps []Pending) []int {
	var out []int
	for _, p := range ps {
		out = append(out, p.Seq)
	}
	return out
}

func newQueue() *TranslateQueue {
	return NewTranslateQueue(TranslateWaiting, TranslateMaxWait, TranslateMerge)
}

func TestQueueDropsTheOldestWaitingPhrase(t *testing.T) {
	q := newQueue()
	t0 := time.Now()
	for i := 1; i <= 6; i++ {
		if d := q.Push(phrase(i, 5, t0, true)); len(d) != 0 {
			t.Fatalf("push %d dropped %v", i, seqsOf(d))
		}
	}
	d := q.Push(phrase(7, 5, t0, true))
	if !reflect.DeepEqual(seqsOf(d), []int{1}) {
		t.Fatalf("the seventh phrase dropped %v", seqsOf(d))
	}
	if q.Len() != 6 {
		t.Fatalf("len %d", q.Len())
	}
	var got []int
	for {
		b, ok, _ := q.Pop(t0.Add(time.Second))
		if !ok {
			break
		}
		got = append(got, b.Seqs...)
	}
	if !reflect.DeepEqual(got, []int{2, 3, 4, 5, 6, 7}) {
		t.Fatalf("popped %v", got)
	}
	c := q.Counts()
	if c.Pushed != 7 || c.Overflow != 1 || c.Stale != 0 || c.Dropped() != 1 {
		t.Fatalf("counts %+v", c)
	}
}

func TestQueueDropsAPhraseThatWaitedTooLong(t *testing.T) {
	q := newQueue()
	t0 := time.Now()
	q.Push(phrase(1, 5, t0, true))
	q.Push(phrase(2, 5, t0.Add(5*time.Second), true))
	q.Push(phrase(3, 5, t0.Add(10*time.Second), true))
	b, ok, stale := q.Pop(t0.Add(15*time.Second + time.Millisecond))
	if !ok || b.Seq() != 3 || !reflect.DeepEqual(seqsOf(stale), []int{1, 2}) {
		t.Fatalf("pop gave %+v %v, stale %v", b, ok, seqsOf(stale))
	}
	if !b.Heard.Equal(t0.Add(10*time.Second)) || !b.Start.Equal(t0.Add(7*time.Second)) {
		t.Fatalf("the batch did not keep the times of its phrase: %+v", b)
	}
	q.Push(phrase(4, 5, t0, true))
	if b, ok, _ := q.Pop(t0.Add(TranslateMaxWait)); !ok || b.Seq() != 4 {
		t.Fatal("a phrase that waited exactly the limit was dropped")
	}
	q.Push(phrase(5, 5, t0, true))
	if _, ok, stale := q.Pop(t0.Add(time.Minute)); ok || len(stale) != 1 {
		t.Fatalf("a lone old phrase: ok %v, stale %v", ok, seqsOf(stale))
	}
	if _, ok, stale := q.Pop(t0.Add(time.Minute)); ok || len(stale) != 0 {
		t.Fatal("an empty queue gave something")
	}
	q.Push(phrase(6, 5, t0, true))
	q.Push(phrase(7, 5, t0, true))
	if d := q.Clear(); !reflect.DeepEqual(seqsOf(d), []int{6, 7}) || q.Len() != 0 {
		t.Fatalf("clear %v, len %d", seqsOf(d), q.Len())
	}
	c := q.Counts()
	if c.Pushed != 7 || c.Stale != 3 || c.Cleared != 2 || c.Overflow != 0 || c.Dropped() != 5 {
		t.Fatalf("counts %+v", c)
	}
}

func TestQueueMergesWaitingPhrasesOfOneParagraph(t *testing.T) {
	q := newQueue()
	t0 := time.Now()
	q.Push(phrase(1, 12, t0, true))
	q.Push(phrase(2, 10, t0, false))
	q.Push(phrase(3, 8, t0, false))
	q.Push(phrase(4, 2, t0, false))
	b, ok, _ := q.Pop(t0)
	if !ok || !reflect.DeepEqual(b.Seqs, []int{1, 2, 3}) || b.Words != 30 {
		t.Fatalf("batch %+v", b)
	}
	if want := phrase(1, 12, t0, true).Text + " " + phrase(2, 10, t0, false).Text + " " + phrase(3, 8, t0, false).Text; b.Text != want {
		t.Fatalf("text %q", b.Text)
	}
	if b, _, _ := q.Pop(t0); !reflect.DeepEqual(b.Seqs, []int{4}) {
		t.Fatalf("the rest %+v", b)
	}

	q.Push(phrase(5, 6, t0, false))
	q.Push(phrase(6, 6, t0, true))
	q.Push(phrase(7, 6, t0, false))
	if b, _, _ := q.Pop(t0); !reflect.DeepEqual(b.Seqs, []int{5}) {
		t.Fatalf("a paragraph start joined the batch: %+v", b)
	}
	if b, _, _ := q.Pop(t0); !reflect.DeepEqual(b.Seqs, []int{6, 7}) {
		t.Fatalf("a paragraph did not start its own batch: %+v", b)
	}

	q.Push(phrase(8, 6, t0, false))
	other := phrase(9, 6, t0, false)
	other.Target = "fr"
	q.Push(other)
	b, _, _ = q.Pop(t0)
	if !reflect.DeepEqual(b.Seqs, []int{8}) || b.Target != "de" {
		t.Fatalf("another target joined the batch: %+v", b)
	}
	if b, _, _ := q.Pop(t0); !reflect.DeepEqual(b.Seqs, []int{9}) || b.Target != "fr" {
		t.Fatalf("the other target %+v", b)
	}

	q.Push(phrase(10, 18, t0, false))
	q.Push(phrase(11, 13, t0, false))
	if b, _, _ := q.Pop(t0); !reflect.DeepEqual(b.Seqs, []int{10}) {
		t.Fatalf("a batch went over the limit: %+v", b)
	}
	q.Clear()

	q.Push(Pending{Seq: 12, Text: "three words here", Target: "de", Heard: t0})
	if b, _, _ := q.Pop(t0); b.Words != 3 {
		t.Fatalf("words were not counted: %+v", b)
	}
}

func TestCleanTranslation(t *testing.T) {
	cases := []struct {
		line, answer, target, want string
		verdict                    Verdict
	}{
		{"Hello there", "Hallo da", "de", "Hallo da", TrOK},
		{"Hello there", `"Hallo da"`, "de", "Hallo da", TrOK},
		{"Hello there", "«Привет»", "ru", "Привет", TrOK},
		{"Hello there", "„Hallo da“", "de", "Hallo da", TrOK},
		{"Hello there", "“Hallo da”", "de", "Hallo da", TrOK},
		{"Hello there", "'Hallo da'", "de", "Hallo da", TrOK},
		{"Hello there", "‘Ciao’", "it", "Ciao", TrOK},
		{"Hello there", "「Hallo da」", "de", "Hallo da", TrOK},
		{"Hello there", "\"« Hallo da »\"", "de", "Hallo da", TrOK},
		{"Hello there", `"Hallo", sagte er.`, "de", `"Hallo", sagte er.`, TrOK},
		{"Hello there", "  Hallo   \t da  ", "de", "Hallo da", TrOK},
		{"Hello there", "\n\n  Hallo da\nNote: this is German.\n", "de", "Hallo da", TrOK},
		{"Hello there", "\r\nHallo da\r\n(translated)", "de", "Hallo da", TrOK},
		{"Hello there", "", "de", "", TrEmpty},
		{"Hello there", " \n \"\" \n", "de", "", TrEmpty},
		{"Hello there", "...", "de", "", TrEmpty},
		{"Hello, there!", "hello there", "de", "", TrSame},
		{"Hello there", "  HELLO,   THERE.  ", "de", "", TrSame},
		{"Привет, мир", "«привет мир!»", "ru", "", TrSame},
		{"Hi", strings.Repeat("x", 49), "de", "", TrRamble},
		{"Hi", strings.Repeat("x", 48), "de", strings.Repeat("x", 48), TrOK},
		{"Hi", "\"" + strings.Repeat("x", 48) + "\"", "de", strings.Repeat("x", 48), TrOK},
	}
	for _, c := range cases {
		got, v := CleanTranslation(c.line, c.answer, c.target)
		if got != c.want || v != c.verdict {
			t.Errorf("CleanTranslation(%q, %q, %q) = %q, %q; want %q, %q", c.line, c.answer, c.target, got, v, c.want, c.verdict)
		}
	}
}

func TestCleanTranslationRejectsTheWrongScript(t *testing.T) {
	cases := []struct {
		line, answer, target string
		verdict              Verdict
	}{
		{"Press the keyboard shortcut to start.", "Нажмите клавиатурным快捷方式, чтобы начать.", "ru", TrScript},
		{"Look at your screen.", "Посмотрите на ваш屏幕.", "ru", TrScript},
		{"So basically what happened is the server ran out of memory.", "Так что basically что happened is the server ran out of memory", "ru", TrScript},
		{"The server restarted at night.", "Сервер перезапустился сегодня ночью.", "en", TrScript},
		{"The server restarted at night.", "Сервер перезапустился ночью.", "en", TrOK},
		{"The server restarted at night.", "Сервер перезапустився сьогодні вночі.", "de", TrScript},
		{"Anna will put it on GitHub tomorrow.", "Анна выложит это на GitHub завтра.", "ru", TrOK},
		{"Good morning everyone", "Доброе утро everyone", "ru", TrOK},
		{"Good day, Kyiv", "Добрий день, Kyiv", "uk", TrOK},
		{"Good morning everyone", "Guten Morgen zusammen", "de", TrOK},
		{"Good morning everyone", "Dzień dobry wszystkim", "pl", TrOK},
		{"Fine, 1234567890 — 42%!", "Ок, 1234567890 — 42%!", "ru", TrOK},
		{"Fine, 1234567890 — 42%!", "OK, 1234567890 — 42%!", "ru", TrOK},
		{"So GPT Six Sol is like GPT Tera.", "Так GPT Six Sol — это GPT Tera.", "ru", TrOK},
		{"Claude Opus Five Point Five on Bridgebench,", "Claude Opus Five Point Five на Bridgebench,", "ru", TrOK},
		{"and this one is actually really interesting.", "and this one is actually really interesting, yes", "ru", TrScript},
		{"Here is the AGI test.", "Вот тест AGI, сделанный в OpenAI и NVIDIA.", "uk", TrOK},
		{"it was cheap", "it was дёшево", "ru", TrOK},
		{"it was really very cheap", "it was really very дёшево", "ru", TrScript},
		{"Москва и Киев", "Москва и Киев", "de", TrSame},
		{"the city of Moscow", "Москва, Київ, Львів, Одеса", "de", TrOK},
		{"the cities", "мы видели эти города", "de", TrScript},
		{"Ok", "Хорошо, мистер Smith", "ru", TrOK},
		{"Hello", "Hallo！", "de", TrScript},
		{"Hello", "Hello 안녕", "fr", TrScript},
		{"Hello", "Hallo ハロー", "de", TrScript},
		{"Hello", "Привет", "", TrOK},
		{"Hello", "Hallo 中", "", TrScript},
		{"Quote", "«Цитата» — №1…", "ru", TrOK},
	}
	for _, c := range cases {
		_, v := CleanTranslation(c.line, c.answer, c.target)
		if v != c.verdict {
			t.Errorf("CleanTranslation(%q, %q, %q) gave %q, want %q", c.line, c.answer, c.target, v, c.verdict)
		}
	}
}

func TestTranslateGrammarForbidsTheForeignScripts(t *testing.T) {
	want := "root ::= [^⺀-鿿가-힯぀-ヿ＀-￯]*"
	for _, code := range TranslateTargets {
		if got := TranslateGrammar(code); got != want {
			t.Errorf("grammar for %s: %q", code, got)
		}
	}
	for _, r := range []rune{0x2e80, 0x9fff, 0xac00, 0xd7af, 0x3040, 0x30ff, 0xff00, 0xffef} {
		if !strings.ContainsRune(want, r) {
			t.Errorf("the grammar lacks the real character %U", r)
		}
	}
	for _, code := range []string{"", "auto", "xx", "DE", "zh"} {
		if got := TranslateGrammar(code); got != "" {
			t.Errorf("grammar for an unknown code %q: %q", code, got)
		}
	}
	for _, r := range []rune{0x2e80, 0x4e2d, 0x9fff, 0xac00, 0xd7af, 0x3042, 0x30a2, 0xff01, 0xffef} {
		if _, v := CleanTranslation("Hello", "Hallo "+string(r), "de"); v != TrScript {
			t.Errorf("an answer with %U was not rejected: %q", r, v)
		}
	}
}

func TestCleanTranslationRambleLimit(t *testing.T) {
	line := strings.Repeat("a", 10)
	limit := 4*10 + 40
	if got, v := CleanTranslation(line, strings.Repeat("b", limit), "de"); v != TrOK || len(got) != limit {
		t.Fatalf("an answer at the limit was thrown away: %q", v)
	}
	if _, v := CleanTranslation(line, strings.Repeat("b", limit+1), "de"); v != TrRamble {
		t.Fatalf("an answer over the limit gave %q", v)
	}
}

func TestLanguageNames(t *testing.T) {
	names := map[string]string{
		"de": "German", "en": "English", "es": "Spanish", "fr": "French",
		"it": "Italian", "pl": "Polish", "uk": "Ukrainian", "ru": "Russian",
	}
	if len(TranslateTargets) != 8 {
		t.Fatalf("targets %v", TranslateTargets)
	}
	for _, code := range TranslateTargets {
		if LanguageName(code) != names[code] {
			t.Errorf("name of %s is %q", code, LanguageName(code))
		}
	}
	for _, code := range []string{"", "auto", "xx", "DE", "zh"} {
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

func TestTranslateMessages(t *testing.T) {
	sys := map[string]string{
		"ru": "You translate live subtitles into Russian. Each message is a piece of speech, often only a part of a sentence. Reply with the Russian translation of that piece only.",
		"en": "You translate live subtitles into English. Each message is a piece of speech, often only a part of a sentence. Reply with the English translation of that piece only.",
	}
	for code, want := range sys {
		if got := TranslateSystem(code); got != want {
			t.Errorf("system for %s: %q", code, got)
		}
	}
	for _, code := range TranslateTargets {
		m := TranslateMessages(code, "  the build went out\n on Tuesday ")
		if len(m) != 8 {
			t.Fatalf("%s: %d messages", code, len(m))
		}
		if m[0].Role != "system" || m[0].Content != TranslateSystem(code) || !strings.Contains(m[0].Content, LanguageName(code)) {
			t.Errorf("%s: system %+v", code, m[0])
		}
		for i := 1; i < 7; i += 2 {
			if m[i].Role != "user" || m[i+1].Role != "assistant" || m[i].Content == "" || m[i+1].Content == "" {
				t.Errorf("%s: example turn %d: %+v %+v", code, i, m[i], m[i+1])
			}
		}
		if m[7].Role != "user" || m[7].Content != "the build went out on Tuesday" {
			t.Errorf("%s: the piece %+v", code, m[7])
		}
		if code != "en" && m[1].Content != "and then we went back to the office," {
			t.Errorf("%s: the first example is not the English piece: %q", code, m[1].Content)
		}
	}
	ru := TranslateMessages("ru", "x")
	if ru[2].Content != "а потом мы вернулись в офис," || ru[4].Content != "поэтому второй тест занял почти три минуты" || ru[6].Content != "Так что если посмотреть на цифры, Windows Eleven" {
		t.Errorf("russian answers %+v", ru)
	}
	en := TranslateMessages("en", "x")
	wantEn := []Message{
		{"user", "y luego volvimos a la oficina,"}, {"assistant", "and then we went back to the office,"},
		{"user", "weshalb der zweite Test fast drei Minuten gedauert hat"}, {"assistant", "which is why the second test took almost three minutes"},
		{"user", "Donc si vous regardez les chiffres ici, Windows Eleven est"}, {"assistant", "So if you look at the numbers here, Windows Eleven is"},
	}
	if !reflect.DeepEqual(en[1:7], wantEn) {
		t.Errorf("english examples %+v", en[1:7])
	}
	for _, code := range []string{"", "auto", "xx", "DE", "zh"} {
		if m := TranslateMessages(code, "hello"); m != nil {
			t.Errorf("messages for an unknown code %q: %+v", code, m)
		}
		if TranslateSystem(code) != "" {
			t.Errorf("system for an unknown code %q", code)
		}
	}
}

func TestTranslateTokens(t *testing.T) {
	cases := map[int]int{0: 48, 1: 48, 8: 48, 9: 51, 15: 69, 18: 78, 30: 114, 45: 159, 46: 160, 200: 160}
	for words, want := range cases {
		if got := TranslateTokens(words); got != want {
			t.Errorf("TranslateTokens(%d) = %d, want %d", words, got, want)
		}
	}
}

func TestStreamingText(t *testing.T) {
	cases := map[string]string{
		"":                         "",
		"«Привет":                  "Привет",
		"\"Hallo da":               "Hallo da",
		"\" « Hallo":               "Hallo",
		"\n\n  Так что   если":     "Так что если",
		"Bonjour\nNote: French":    "Bonjour",
		"„Guten Morgen“, sagte er": "Guten Morgen“, sagte er",
		"\"":                       "",
		"Ciao \"amici\"":           "Ciao \"amici\"",
	}
	for in, want := range cases {
		if got := StreamingText(in); got != want {
			t.Errorf("StreamingText(%q) = %q, want %q", in, got, want)
		}
	}
}
