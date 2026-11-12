package subtitles

import (
	"strings"
	"testing"
	"time"
)

func TestTheSplitOfWayC(t *testing.T) {
	cases := []struct {
		name, answer, want string
		ok                 bool
	}{
		{"clean", "<b>и я думаю, что это проделало феноменальную работу</b>, если мы посмотрим на тест AGI,", ", если мы посмотрим на тест AGI,", true},
		{"clean with a space", "  <b>Взгляните на разницу здесь.</b> GPT Six Sol сейчас безумно дешев", "GPT Six Sol сейчас безумно дешев", true},
		{"entities", "<b>a</b> b &amp; c &lt;d&gt;", "b & c <d>", true},
		{"no tag", "просто перевод без меток", "", false},
		{"two tagged parts", "<b>это также</b> и мы только <b>что сказали</b>, что это", "", false},
		{"empty inside", "<b> </b> всё остальное", "", false},
		{"empty tail", "<b>всё внутри</b>  ", "", false},
		{"tag in the middle", "Итак, <b>это невероятно</b> дешевле", "", false},
		{"only an opening tag", "<b>начало без конца", "", false},
		{"a closing tag first", "</b>что-то<b>", "", false},
	}
	for _, c := range cases {
		got, ok := SplitTagged(c.answer)
		if ok != c.ok || got != c.want {
			t.Errorf("%s: %q %v, want %q %v", c.name, got, ok, c.want, c.ok)
		}
	}
	if got := TagPrevious("a < b & c", "d > e"); got != "<b>a &lt; b &amp; c</b> d &gt; e" {
		t.Fatalf("tagged %q", got)
	}
}

func TestWhenAPieceIsTranslatedAlone(t *testing.T) {
	cases := []struct {
		prev string
		para bool
		text string
		want bool
	}{
		{"", false, "So the first thing", true},
		{"we looked at the data.", false, "Most users left", true},
		{"we looked at the data", false, "and most users left", false},
		{"we looked at the data", true, "and most users left", true},
		{"you can really get into it.", false, "But then", true},
		{"you can really get into it.", false, "but then", false},
		{"is it cheaper?", false, "It is.", true},
		{"compared with Opus Five,", false, "which cost less", false},
	}
	for _, c := range cases {
		if got := BeginsSentence(c.prev, c.para, c.text); got != c.want {
			t.Errorf("BeginsSentence(%q, %v, %q) = %v", c.prev, c.para, c.text, got)
		}
	}
}

func TestThePunctuationRuleForTheWindowAndTheFile(t *testing.T) {
	for _, c := range []struct {
		text, file string
		glue       bool
	}{
		{", если мы посмотрим", "если мы посмотрим", true},
		{". Это всё", "Это всё", true},
		{"» и дальше", "и дальше", true},
		{"и дальше", "и дальше", false},
		{"", "", false},
	} {
		if GluesToPrevious(c.text) != c.glue || WithoutLeadingMark(c.text) != c.file {
			t.Errorf("%q: glue %v, file %q", c.text, GluesToPrevious(c.text), WithoutLeadingMark(c.text))
		}
	}
	rows := pieceRows([]PieceView{
		{Seq: 1, Para: true, Text: "и я думаю, что это проделало работу", State: PieceDone},
		{Seq: 2, Text: ", если мы посмотрим на тест,", State: PieceDone},
		{Seq: 3, Text: "который прошёл", State: PieceDone},
	}, func(i int) string {
		return []string{"и я думаю, что это проделало работу", ", если мы посмотрим на тест,", "который прошёл"}[i]
	}, false)
	if len(rows) != 1 || rows[0].Text != "и я думаю, что это проделало работу, если мы посмотрим на тест, который прошёл" {
		t.Fatalf("rows %+v", rows)
	}
}

func TestTheFastCleaning(t *testing.T) {
	cases := []struct{ src, out, want string }{
		{"and then we went back to the", "а потом мы вернулись в.", "а потом мы вернулись в"},
		{"we went back to the office.", "мы вернулись в офис,", "мы вернулись в офис."},
		{"than Deep Seek V four point one flash by seventy five percent,", ", чем Deep Seek V четыре точки одна вспышка на семьдесят пять процентов.", ", чем Deep Seek V четыре точки одна вспышка на семьдесят пять процентов"},
		{"GPT Six Sol", "GPT Six Sol", ""},
	}
	for _, c := range cases {
		got, v := CleanFastTranslation(c.src, c.out)
		if c.want == "" {
			if v != TrSame {
				t.Errorf("%q: verdict %q", c.out, v)
			}
			continue
		}
		if v != TrOK || got != c.want {
			t.Errorf("%q: %q %q", c.out, got, v)
		}
	}
	if _, v := CleanFastTranslation("Opus five point five is cheaper", "Opus 5.5 is cheaper now"); v != TrOK {
		t.Fatalf("an answer made mostly of Latin names was rejected: %q", v)
	}
	if _, v := CleanFastTranslation("x", " "); v != TrEmpty {
		t.Fatal("empty")
	}
}

func TestTheBatchKeepsTheParagraphMark(t *testing.T) {
	q := newQueue()
	t0 := time.Date(2026, 11, 12, 20, 0, 0, 0, time.UTC)
	q.Push(Pending{Seq: 1, Text: "first piece", Target: "ru", Para: true, Heard: t0})
	b, ok, _ := q.Pop(t0)
	if !ok || !b.Para {
		t.Fatalf("%+v", b)
	}
	if !strings.Contains(b.Text, "first") {
		t.Fatal(b.Text)
	}
}

func TestTheSpeechLanguage(t *testing.T) {
	limit := 30 * time.Second
	cases := []struct {
		setting, guess string
		waited         time.Duration
		lang           string
		state          SpeechState
	}{
		{"en", "", 0, "en", SpeechKnown},
		{"DE", "ru", time.Minute, "de", SpeechKnown},
		{"auto", "ru", 0, "ru", SpeechKnown},
		{"auto", "", 10 * time.Second, "", SpeechWaiting},
		{"", "", 29 * time.Second, "", SpeechWaiting},
		{"auto", "", 30 * time.Second, "", SpeechLost},
		{"auto", "uk", 40 * time.Second, "uk", SpeechKnown},
	}
	for _, c := range cases {
		lang, st := SpeechLanguage(c.setting, c.guess, c.waited, limit)
		if lang != c.lang || st != c.state {
			t.Errorf("%+v: %q %q", c, lang, st)
		}
	}
	words := strings.Fields(guessEnglish)
	before, _ := SpeechLanguage("auto", GuessLanguage(words), 0, limit)
	words = append(words, strings.Fields(guessRussian)...)
	after, _ := SpeechLanguage("auto", GuessLanguage(words), 0, limit)
	if before != "en" || after != "ru" {
		t.Fatalf("a change during a run: %q then %q", before, after)
	}
}
