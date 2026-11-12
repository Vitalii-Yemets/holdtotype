package subtitles

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

type feed struct {
	p   *Phraser
	t0  time.Time
	at  time.Duration
	out []TextPhrase
}

func newFeed() *feed {
	return &feed{p: NewPhraser(), t0: time.Date(2026, 11, 11, 10, 0, 0, 0, time.UTC)}
}

func (f *feed) now() time.Time { return f.t0.Add(f.at) }

func (f *feed) wait(d time.Duration) { f.at += d }

func (f *feed) send(seg int, text string, final bool) []TextPhrase {
	got := f.p.Update(seg, text, final, f.now())
	f.out = append(f.out, got...)
	return got
}

func (f *feed) tick() []TextPhrase {
	got := f.p.Tick(f.now())
	f.out = append(f.out, got...)
	return got
}

func partial(w string, step int) string {
	trimmed := strings.TrimRight(w, ".,;:!?…")
	if trimmed != w && trimmed != "" {
		return trimmed
	}
	r := []rune(w)
	if step%2 == 0 && len(r) > 4 {
		return string(r[:len(r)/2])
	}
	return w
}

func (f *feed) speak(seg int, words []string) {
	steps := []int{2, 1, 3, 4, 2, 1, 3}
	shown := 0
	for i := 0; shown < len(words); i++ {
		shown += steps[i%len(steps)]
		if shown > len(words) {
			shown = len(words)
		}
		part := append([]string(nil), words[:shown]...)
		part[len(part)-1] = partial(part[len(part)-1], i)
		f.wait(600 * time.Millisecond)
		f.send(seg, strings.Join(part, " "), false)
		if part[len(part)-1] != words[shown-1] {
			f.wait(600 * time.Millisecond)
			f.send(seg, strings.Join(words[:shown], " "), false)
		}
	}
}

func textsOf(ps []TextPhrase) []string {
	var out []string
	for _, p := range ps {
		out = append(out, p.Text)
	}
	return out
}

func TestSentenceEndIsCutOnlyAfterTheNextWord(t *testing.T) {
	f := newFeed()
	if got := f.send(0, "We went back to the office.", false); len(got) != 0 {
		t.Fatalf("the unstable last word was used: %v", textsOf(got))
	}
	f.wait(600 * time.Millisecond)
	got := f.send(0, "We went back to the office. Then", false)
	if !reflect.DeepEqual(textsOf(got), []string{"We went back to the office."}) {
		t.Fatalf("phrases %v", textsOf(got))
	}
	if got[0].Words != 6 || got[0].Seq != 1 || !got[0].Para {
		t.Fatalf("phrase %+v", got[0])
	}
	if !got[0].Start.Equal(f.t0) || !got[0].Heard.Equal(f.t0.Add(600*time.Millisecond)) {
		t.Fatalf("times %+v", got[0])
	}
	if f.p.Waiting() != 1 {
		t.Fatalf("waiting %d", f.p.Waiting())
	}
}

func TestEverySentenceEndClosesAPiece(t *testing.T) {
	f := newFeed()
	f.speak(0, strings.Fields("Yes it is. We will see you there tomorrow. Fine. And"))
	want := []string{"Yes it is.", "We will see you there tomorrow.", "Fine."}
	if !reflect.DeepEqual(textsOf(f.out), want) {
		t.Fatalf("phrases %v", textsOf(f.out))
	}
	p := NewPhraser()
	t0 := time.Now()
	got := p.Whole("Right. So the long sentence that goes on and on finally ends here. Yes.", t0, t0, t0)
	if !reflect.DeepEqual(textsOf(got), []string{"Right.", "So the long sentence that goes on and on finally ends here.", "Yes."}) {
		t.Fatalf("whole %v", textsOf(got))
	}
}

func TestShortWordBeforeAFullStop(t *testing.T) {
	cases := []struct {
		w, next string
		last    bool
		want    bool
	}{
		{"it.", "But", false, true},
		{"it.", "but", false, false},
		{"it.", "", true, true},
		{"it.", "", false, false},
		{"It.", "But", false, false},
		{"It.", "", true, false},
		{"Mr.", "Smith", false, false},
		{"A.", "Smith", false, false},
		{"a.", "Smith", false, true},
		{"it.", "«But", false, true},
		{"it.", "2026", false, false},
		{"42.", "next", false, true},
		{"e.g.", "The", false, false},
		{"office.", "and", false, true},
		{"so", "But", false, false},
	}
	for _, c := range cases {
		if got := EndsSentenceBefore(c.w, c.next, c.last); got != c.want {
			t.Errorf("EndsSentenceBefore(%q, %q, %v) = %v", c.w, c.next, c.last, got)
		}
	}
	f := newFeed()
	f.send(0, "you can really get into it.", false)
	if len(f.out) != 0 {
		t.Fatalf("the unstable last word was used: %v", textsOf(f.out))
	}
	f.wait(600 * time.Millisecond)
	f.send(0, "you can really get into it. Bu", false)
	if !reflect.DeepEqual(textsOf(f.out), []string{"you can really get into it."}) {
		t.Fatalf("the unstable next word did not decide: %v", textsOf(f.out))
	}
	g := newFeed()
	g.send(0, "you can really get into it. but then we", false)
	if len(g.out) != 0 {
		t.Fatalf("a lower-case next word ended the sentence: %v", textsOf(g.out))
	}
	h := newFeed()
	h.send(0, "and that is all there is to it.", false)
	h.wait(300 * time.Millisecond)
	if got := h.send(0, "and that is all there is to it.", true); !reflect.DeepEqual(textsOf(got), []string{"and that is all there is to it."}) {
		t.Fatalf("the end of a finished segment %v", textsOf(got))
	}
	if h.p.Waiting() != 0 {
		t.Fatalf("waiting %d", h.p.Waiting())
	}
}

func numbered(n int) []string {
	names := strings.Fields("one two three four five six seven eight nine ten eleven twelve thirteen fourteen fifteen sixteen seventeen eighteen nineteen twenty twentyone twentytwo twentythree twentyfour")
	return append([]string(nil), names[:n]...)
}

func hardCut(lang string, words []string) []string {
	p := NewPhraser()
	if lang != "-" {
		p.SetLanguage(lang)
	}
	t0 := time.Date(2026, 11, 12, 10, 0, 0, 0, time.UTC)
	return textsOf(p.Update(0, strings.Join(words, " "), false, t0))
}

func TestHardLimitCutsBeforeALinkingWord(t *testing.T) {
	words := numbered(20)
	words[11] = "because"
	got := hardCut("en", words)
	if len(got) != 1 || got[0] != strings.Join(words[:11], " ") {
		t.Fatalf("phrases %v", got)
	}
	words = numbered(20)
	words[7], words[13] = "and", "with"
	got = hardCut("en", words)
	if len(got) != 1 || got[0] != strings.Join(words[:13], " ") {
		t.Fatalf("the last linking word was not used: %v", got)
	}
	words = numbered(20)
	words[17] = "and"
	got = hardCut("en", words)
	if len(got) != 1 || got[0] != strings.Join(words[:17], " ") {
		t.Fatalf("a linking word as the newest word: %v", got)
	}
}

func TestHardLimitFallsBackToKeepingThreeWords(t *testing.T) {
	words := numbered(20)
	got := hardCut("en", words)
	if len(got) != 1 || got[0] != strings.Join(words[:PhraseHard-PhraseKeep], " ") {
		t.Fatalf("no linking word: %v", got)
	}
	words = numbered(20)
	words[4] = "and"
	got = hardCut("en", words)
	if len(got) != 1 || got[0] != strings.Join(words[:PhraseHard-PhraseKeep], " ") {
		t.Fatalf("a linking word that leaves too short a piece: %v", got)
	}
	words = numbered(20)
	words[5] = "and"
	got = hardCut("en", words)
	if len(got) != 1 || got[0] != strings.Join(words[:5], " ") {
		t.Fatalf("a linking word that leaves exactly the shortest piece: %v", got)
	}
}

func TestHardLimitUsesTheSpeechLanguage(t *testing.T) {
	words := numbered(20)
	words[10] = "und"
	cut := strings.Join(words[:10], " ")
	keep := strings.Join(words[:PhraseHard-PhraseKeep], " ")
	cases := map[string]string{"de": cut, "en": keep, "auto": cut, "-": cut, "xx": cut, "ru": keep}
	for lang, want := range cases {
		if got := hardCut(lang, words); len(got) != 1 || got[0] != want {
			t.Errorf("language %q: %v", lang, got)
		}
	}
	words = numbered(20)
	words[9] = "который"
	if got := hardCut("ru", words); len(got) != 1 || got[0] != strings.Join(words[:9], " ") {
		t.Errorf("russian: %v", got)
	}
}

func TestHardLimitDoesNotApplyAtTheSoftLimit(t *testing.T) {
	words := numbered(17)
	words[8] = "because"
	if got := hardCut("en", words); len(got) != 0 {
		t.Fatalf("a cut before the hard limit: %v", got)
	}
}

func TestLinkingWords(t *testing.T) {
	cases := []struct {
		lang, w string
		want    bool
	}{
		{"en", "And,", true}, {"en", "because", true}, {"en", "than", true}, {"en", "of", false}, {"en", "to", false}, {"en", "the", false},
		{"ru", "который", true}, {"ru", "В", true}, {"ru", "of", false},
		{"uk", "якщо", true}, {"de", "dass", true}, {"de", "zu", false}, {"de", "der", false},
		{"fr", "qu'il", true}, {"fr", "lorsqu’on", true}, {"fr", "parce", true}, {"fr", "de", false}, {"fr", "des", false}, {"fr", "l'utilisateur", false},
		{"es", "porque", true}, {"es", "de", false}, {"es", "del", false},
		{"it", "nell'anno", true}, {"it", "perché", true}, {"it", "di", false}, {"it", "della", false},
		{"pl", "ponieważ", true}, {"pl", "żeby", true},
		{"auto", "und", true}, {"auto", "который", true}, {"xx", "because", true}, {"", "the", false},
	}
	for _, c := range cases {
		if got := IsLinkingWord(c.lang, c.w); got != c.want {
			t.Errorf("IsLinkingWord(%q, %q) = %v", c.lang, c.w, got)
		}
	}
	for _, lang := range TranslateTargets {
		if len(linkingSets[lang]) == 0 {
			t.Errorf("no linking words for %s", lang)
		}
		for _, w := range []string{"of", "to", "the", "de", "du", "di", "del", "della", "zu", "der", "die", "das", "le", "la", "el", "il"} {
			if linkingSets[lang][w] {
				t.Errorf("%s: %q must stay with the word before it", lang, w)
			}
		}
	}
}

func TestLongClauseIsCutAtTheComma(t *testing.T) {
	f := newFeed()
	f.speak(0, strings.Fields("When the server started this morning at nine, everyone in the room looked at the"))
	if len(f.out) != 1 || f.out[0].Text != "When the server started this morning at nine," || f.out[0].Words != 8 {
		t.Fatalf("phrases %v", textsOf(f.out))
	}
	f2 := newFeed()
	f2.speak(0, strings.Fields("Well, so the thing is that we never really looked at the logs"))
	if len(f2.out) != 0 {
		t.Fatalf("a clause of two words was cut: %v", textsOf(f2.out))
	}
	f3 := newFeed()
	f3.speak(0, strings.Fields("the first build was fine and the second one — it was slower than anything"))
	if len(f3.out) != 1 || f3.out[0].Text != "the first build was fine and the second one —" {
		t.Fatalf("a lone dash did not end the clause: %v", textsOf(f3.out))
	}
}

func TestNoPunctuationCutsAtTheHardLimit(t *testing.T) {
	var words []string
	for i := 0; i < 40; i++ {
		words = append(words, "word")
	}
	f := newFeed()
	f.speak(0, words)
	if len(f.out) != 2 {
		t.Fatalf("phrases %v", textsOf(f.out))
	}
	for _, p := range f.out {
		if p.Words != PhraseHard-PhraseKeep {
			t.Fatalf("phrase of %d words", p.Words)
		}
	}
	if f.p.Waiting() != 10 {
		t.Fatalf("waiting %d", f.p.Waiting())
	}
}

func TestFinalAfterSilenceGivesOutTheRest(t *testing.T) {
	f := newFeed()
	f.speak(0, strings.Fields("and then we simply went home"))
	f.wait(1200 * time.Millisecond)
	got := f.send(0, "and then we simply went home", true)
	if !reflect.DeepEqual(textsOf(got), []string{"and then we simply went home"}) {
		t.Fatalf("phrases %v", textsOf(got))
	}
	if f.p.Waiting() != 0 {
		t.Fatalf("waiting %d", f.p.Waiting())
	}
	f.wait(400 * time.Millisecond)
	if got := f.send(1, "", false); len(got) != 0 {
		t.Fatalf("an empty segment gave %v", textsOf(got))
	}
	f.wait(5 * time.Second)
	if got := f.tick(); len(got) != 0 {
		t.Fatalf("tick without a carry gave %v", textsOf(got))
	}
}

func TestFinalOfAFinishedSentenceGivesItOutAtOnce(t *testing.T) {
	f := newFeed()
	f.speak(0, strings.Fields("Ok so"))
	f.wait(300 * time.Millisecond)
	got := f.send(0, "Ok so that is the new one.", true)
	if !reflect.DeepEqual(textsOf(got), []string{"Ok so that is the new one."}) {
		t.Fatalf("phrases %v", textsOf(got))
	}
}

func TestFinalAtTheLimitKeepsTheCarry(t *testing.T) {
	f := newFeed()
	f.speak(0, strings.Fields("The second test took almost three minutes. And the third one was"))
	if !reflect.DeepEqual(textsOf(f.out), []string{"The second test took almost three minutes."}) {
		t.Fatalf("phrases before the cut %v", textsOf(f.out))
	}
	f.wait(500 * time.Millisecond)
	if got := f.send(0, "The second test took almost three minutes. And the third one was even", true); len(got) != 0 {
		t.Fatalf("the cut in the middle of speech gave %v", textsOf(got))
	}
	if f.p.Waiting() != 6 {
		t.Fatalf("waiting %d", f.p.Waiting())
	}
	f.wait(500 * time.Millisecond)
	f.send(1, "", false)
	f.wait(600 * time.Millisecond)
	f.send(1, "Slower than", false)
	f.speak(1, strings.Fields("Slower than the first. So we"))
	want := []string{"The second test took almost three minutes.", "And the third one was even Slower than the first."}
	if !reflect.DeepEqual(textsOf(f.out), want) {
		t.Fatalf("phrases %v", textsOf(f.out))
	}
	if f.out[1].Para {
		t.Fatal("a phrase continued after the limit started a paragraph")
	}
	f.wait(2 * time.Second)
	if got := f.tick(); len(got) != 0 {
		t.Fatalf("tick gave the words of an open segment: %v", textsOf(got))
	}
}

func TestCarryLeavesByTick(t *testing.T) {
	f := newFeed()
	f.speak(0, strings.Fields("and this is where it gets"))
	f.wait(400 * time.Millisecond)
	f.send(0, "and this is where it gets", true)
	if len(f.out) != 0 || f.p.Waiting() != 6 {
		t.Fatalf("phrases %v, waiting %d", textsOf(f.out), f.p.Waiting())
	}
	f.wait(500 * time.Millisecond)
	f.send(1, "", false)
	f.wait(PhraseCarry - 600*time.Millisecond)
	if got := f.tick(); len(got) != 0 {
		t.Fatalf("the carry left too early: %v", textsOf(got))
	}
	f.wait(100 * time.Millisecond)
	got := f.tick()
	if !reflect.DeepEqual(textsOf(got), []string{"and this is where it gets"}) {
		t.Fatalf("tick gave %v", textsOf(got))
	}
	if f.p.Waiting() != 0 {
		t.Fatalf("waiting %d", f.p.Waiting())
	}
	if got := f.tick(); len(got) != 0 {
		t.Fatalf("a second tick gave %v", textsOf(got))
	}
}

func TestAbbreviationsAndNumbersDoNotEndASentence(t *testing.T) {
	ends := map[string]bool{
		"office.": true, "Really?": true, "Stop!": true, "wait…": true, "so...": true,
		`there."`: true, "(there.)": true, "»Ende.«": false, "Ende.»": true, "42.": true, "1999.": true,
		"Mr.": false, "Dr.": false, "e.g.": false, "U.S.": false, "i.e.,": false,
		"3.5.": false, "v1.2.": false, "it.": false, "A.": false, "word": false, "word,": false,
		"": false, ".": false, "Tera.": true, "сегодня.": true, "им.": false,
	}
	for w, want := range ends {
		if got := EndsSentence(w); got != want {
			t.Errorf("EndsSentence(%q) = %v", w, got)
		}
	}
	clauses := map[string]bool{
		"office,": true, "first;": true, "note:": true, `said,"`: true, "—": true, "-": true, "–": true,
		"10:30": false, "word": false, "word.": false, "well-known": false,
	}
	for w, want := range clauses {
		if got := EndsClause(w); got != want {
			t.Errorf("EndsClause(%q) = %v", w, got)
		}
	}
	text := "Then Mr. Smith from the U.S. office showed version 3.5. of the app and"
	f := newFeed()
	f.speak(0, strings.Fields(text))
	if len(f.out) != 0 {
		t.Fatalf("a phrase ended at an abbreviation: %v", textsOf(f.out))
	}
	f.wait(1200 * time.Millisecond)
	f.send(0, text, true)
	if !reflect.DeepEqual(textsOf(f.out), []string{text}) {
		t.Fatalf("phrases %v", textsOf(f.out))
	}
}

func TestShrinkingTextNeverTakesWordsBack(t *testing.T) {
	f := newFeed()
	f.send(0, "This is the first sentence here. And then", false)
	if !reflect.DeepEqual(textsOf(f.out), []string{"This is the first sentence here."}) {
		t.Fatalf("phrases %v", textsOf(f.out))
	}
	f.wait(600 * time.Millisecond)
	f.send(0, "This is the", false)
	if f.p.Waiting() != 0 {
		t.Fatalf("waiting after the shrink %d", f.p.Waiting())
	}
	f.wait(600 * time.Millisecond)
	f.send(0, "This is the one we built for the new office. It", false)
	want := []string{"This is the first sentence here.", "one we built for the new office."}
	if !reflect.DeepEqual(textsOf(f.out), want) {
		t.Fatalf("phrases %v", textsOf(f.out))
	}
}

func TestParagraphAfterAPause(t *testing.T) {
	f := newFeed()
	f.send(0, "Good morning to all of you.", false)
	f.wait(1200 * time.Millisecond)
	f.send(0, "Good morning to all of you.", true)
	f.wait(2500 * time.Millisecond)
	f.send(1, "Let us look at the numbers now.", false)
	f.wait(1200 * time.Millisecond)
	f.send(1, "Let us look at the numbers now.", true)
	f.wait(500 * time.Millisecond)
	f.send(2, "They are better than last week.", false)
	f.wait(1200 * time.Millisecond)
	f.send(2, "They are better than last week.", true)
	if len(f.out) != 3 {
		t.Fatalf("phrases %v", textsOf(f.out))
	}
	if !f.out[0].Para || !f.out[1].Para || f.out[2].Para {
		t.Fatalf("paragraphs %v %v %v", f.out[0].Para, f.out[1].Para, f.out[2].Para)
	}
}

func TestSegmentsThatAreClosedOrSkipped(t *testing.T) {
	f := newFeed()
	f.send(0, "one two three four", false)
	f.wait(1500 * time.Millisecond)
	got := f.send(1, "five", false)
	if !reflect.DeepEqual(textsOf(got), []string{"one two three four"}) {
		t.Fatalf("a new segment did not close the old one: %v", textsOf(got))
	}
	f.wait(1500 * time.Millisecond)
	f.send(1, "five six.", true)
	if got := f.send(0, "late words", false); len(got) != 0 || f.p.Waiting() != 0 {
		t.Fatalf("a closed segment came back: %v, waiting %d", textsOf(got), f.p.Waiting())
	}
}

func TestWhole(t *testing.T) {
	p := NewPhraser()
	t0 := time.Date(2026, 11, 11, 10, 0, 0, 0, time.UTC)
	text := "Let us start with a quick look at last week, the new build went out on Tuesday. The feedback so far has been good and"
	got := p.Whole(text, t0, t0.Add(8*time.Second), t0.Add(9*time.Second))
	want := []string{"Let us start with a quick look at last week,", "the new build went out on Tuesday.", "The feedback so far has been good and"}
	if !reflect.DeepEqual(textsOf(got), want) {
		t.Fatalf("phrases %v", textsOf(got))
	}
	if !got[0].Para || got[1].Para || got[2].Para {
		t.Fatalf("paragraphs %+v", got)
	}
	for _, ph := range got {
		if !ph.Start.Equal(t0) || !ph.Heard.Equal(t0.Add(9*time.Second)) {
			t.Fatalf("times %+v", ph)
		}
	}
	next := p.Whole("Okay.", t0.Add(8500*time.Millisecond), t0.Add(9*time.Second), t0.Add(10*time.Second))
	if len(next) != 1 || next[0].Para || next[0].Text != "Okay." || next[0].Seq != 4 {
		t.Fatalf("the next piece %+v", next)
	}
	after := p.Whole("Any questions?", t0.Add(12*time.Second), t0.Add(13*time.Second), t0.Add(14*time.Second))
	if len(after) != 1 || !after[0].Para {
		t.Fatalf("a piece after a pause %+v", after)
	}
	if got := p.Whole("   ", t0, t0, t0); len(got) != 0 {
		t.Fatalf("empty text %v", got)
	}
	if p.Waiting() != 0 {
		t.Fatalf("waiting %d", p.Waiting())
	}
}

func TestFlush(t *testing.T) {
	f := newFeed()
	f.send(0, "and the last words were", false)
	if f.p.Waiting() != 5 {
		t.Fatalf("waiting %d", f.p.Waiting())
	}
	got := f.p.Flush(f.now())
	if !reflect.DeepEqual(textsOf(got), []string{"and the last words were"}) {
		t.Fatalf("flush %v", textsOf(got))
	}
	if f.p.Waiting() != 0 {
		t.Fatalf("waiting %d", f.p.Waiting())
	}
	if got := f.send(0, "A new stream starts again here.", true); len(got) != 1 {
		t.Fatalf("a new stream after the flush: %v", textsOf(got))
	}
	f.send(1, "carried words", true)
	if got := f.p.Flush(f.now()); !reflect.DeepEqual(textsOf(got), []string{"carried words"}) {
		t.Fatalf("flush of a carry %v", textsOf(got))
	}
	if got := f.p.Flush(f.now()); len(got) != 0 {
		t.Fatalf("an empty flush %v", textsOf(got))
	}
}
