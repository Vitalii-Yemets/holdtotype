package subtitles

import (
	"strings"
	"testing"
	"time"
)

func words(n int, w string) string {
	return strings.TrimSpace(strings.Repeat(w+" ", n))
}

func TestTheTopicComesAfterTheFirstWordsAndThenEveryThreeMinutes(t *testing.T) {
	t0 := time.Date(2026, 11, 16, 19, 27, 0, 0, time.Local)
	tp := NewTopic("ru")
	tp.Heard(words(149, "credit"), t0)
	if s := tp.Due(t0); s != StepNone {
		t.Fatalf("due before 150 words: %v", s)
	}
	tp.Heard("cohort", t0.Add(time.Second))
	if s := tp.Due(t0.Add(time.Second)); s != StepTopic {
		t.Fatalf("not due at 150 words: %v", s)
	}
	msgs := tp.TopicMessages(t0.Add(time.Second))
	if len(msgs) != 2 || msgs[0].Role != "system" || !strings.HasPrefix(msgs[1].Content, "Transcript: ") || strings.Count(msgs[1].Content, " ") != TopicFirstWords {
		t.Fatalf("topic request: %+v", msgs)
	}
	if !tp.TakeTopic("Topic: Credit risk in a lending company.", t0.Add(2*time.Second)) {
		t.Fatal("a good topic line is refused")
	}
	if s := tp.Due(t0.Add(3 * time.Second)); s != StepTerms {
		t.Fatalf("terms are not asked right after the topic: %v", s)
	}
	tp.TakeTerms("none", t0.Add(4*time.Second))
	if s := tp.Due(t0.Add(2 * time.Minute)); s != StepNone {
		t.Fatalf("due again too early: %v", s)
	}
	if s := tp.Due(t0.Add(time.Second + TopicEvery)); s != StepTopic {
		t.Fatalf("not due after three minutes: %v", s)
	}
	tp.Heard("new words here", t0.Add(TopicEvery))
	msgs = tp.TopicMessages(t0.Add(time.Second + TopicEvery))
	if !strings.HasPrefix(msgs[1].Content, "Earlier description: Credit risk in a lending company.\nTranscript: ") || !strings.HasSuffix(msgs[1].Content, "new words here") || strings.Contains(msgs[1].Content, "credit credit") {
		t.Fatalf("update request: %q", msgs[1].Content)
	}
}

func TestAModelWithoutTermMemoryAsksOnlyTheTopic(t *testing.T) {
	t0 := time.Now()
	tp := NewTopic("ru")
	tp.Heard("the approval rate of the cohort "+words(200, "x"), t0)
	tp.TopicMessages(t0)
	tp.TakeTopic("Topic: Lending.", t0)
	tp.TakeTerms("approval rate = коэффициент одобрения", t0)
	if len(tp.Terms()) != 1 {
		t.Fatalf("the term is not kept: %+v", tp.Terms())
	}
	tp.UseTerms(false, t0)
	if len(tp.Terms()) != 0 || tp.Context() != " The conversation: Lending." {
		t.Fatalf("the terms stay after the memory is switched off: %+v %q", tp.Terms(), tp.Context())
	}
	at := t0.Add(TopicEvery)
	if s := tp.Due(at); s != StepTopic {
		t.Fatalf("the topic is not due: %v", s)
	}
	tp.TopicMessages(at)
	if !tp.TakeTopic("Topic: Lending and credit risk.", at) {
		t.Fatal("a good topic line is refused")
	}
	if s := tp.Due(at.Add(time.Second)); s != StepNone {
		t.Fatalf("the terms are asked after the topic: %v", s)
	}
	if tp.Context() != " The conversation: Lending and credit risk." || tp.Counts.TermsCalls != 1 {
		t.Fatalf("context %q, counts %+v", tp.Context(), tp.Counts)
	}
	if s := tp.Due(at.Add(TopicEvery)); s != StepTopic {
		t.Fatalf("the next topic round never comes: %v", s)
	}
	tp.pending = StepTerms
	tp.UseTerms(false, at)
	if s := tp.Due(at.Add(time.Second)); s != StepNone {
		t.Fatalf("a pending terms request survives: %v", s)
	}
}

func TestABrokenTopicAnswerSkipsTheRound(t *testing.T) {
	t0 := time.Now()
	tp := NewTopic("ru")
	tp.Heard(words(200, "x"), t0)
	tp.TopicMessages(t0)
	if tp.TakeTopic("nonsense", t0) {
		t.Fatal("a line without the prefix is taken")
	}
	if tp.Text() != "" || tp.Due(t0.Add(time.Minute)) != StepNone || tp.Counts.TopicBad != 1 {
		t.Fatalf("after a broken answer: %q %v %+v", tp.Text(), tp.Due(t0.Add(time.Minute)), tp.Counts)
	}
	tp.Failed(StepTopic, t0)
	if tp.Due(t0.Add(TopicEvery)) != StepTopic {
		t.Fatal("the next round never comes")
	}
}

func TestTheTermsAreFilteredAndKept(t *testing.T) {
	t0 := time.Now()
	tp := NewTopic("ru")
	tp.Heard("the user goes into charge off and the approval rate is low, cohorts of borrowings, an X credit", t0)
	tp.TopicMessages(t0)
	tp.TakeTopic("Topic: Lending.", t0)
	added := tp.TakeTerms("charge-off = списание\napproval rate = коэффициент одобрения (or процент одобрения)\nis = артериальное давление, for example\nblood pressure = давление\nthe = это\nX = икс\ncohort = когорта (или группа)\nloan origination = выдача кредита\nbroken line", t0)
	if added != 3 {
		t.Fatalf("added %d: %+v", added, tp.Terms())
	}
	got := tp.Terms()
	if got[0] != (Term{"charge-off", "списание"}) || got[1] != (Term{"approval rate", "коэффициент одобрения"}) || got[2] != (Term{"cohort", "когорта"}) {
		t.Fatalf("terms: %+v", got)
	}
	c := tp.Counts
	if c.Example != 2 || c.Plain != 2 || c.Unheard != 1 || c.TermsBad != 1 {
		t.Fatalf("counts: %+v", c)
	}
	tp.Due(t0)
	tp.pending = StepTerms
	if tp.TakeTerms("charge off = просрочка\nCohort = группа", t0) != 0 || tp.Terms()[0].Target != "списание" || tp.Counts.Known != 2 {
		t.Fatalf("a written term changed: %+v %+v", tp.Terms(), tp.Counts)
	}
}

func TestTheMemoryStopsAtThirtyTerms(t *testing.T) {
	t0 := time.Now()
	tp := NewTopic("de")
	var heard, answer []string
	for i := 0; i < 40; i++ {
		w := "term" + string(rune('a'+i%26)) + string(rune('a'+i/26))
		heard = append(heard, w)
		answer = append(answer, w+" = Begriff")
	}
	tp.Heard(strings.Join(heard, " "), t0)
	for i := 0; i < 40; i += 6 {
		end := i + 6
		if end > 40 {
			end = 40
		}
		tp.TakeTerms(strings.Join(answer[i:end], "\n"), t0)
	}
	if len(tp.Terms()) != TermsMax || tp.Counts.Full != 40-TermsMax {
		t.Fatalf("terms %d, full %d", len(tp.Terms()), tp.Counts.Full)
	}
}

func TestTheTermsRequestUsesTheRecentPairs(t *testing.T) {
	t0 := time.Date(2026, 11, 16, 19, 27, 0, 0, time.Local)
	tp := NewTopic("ru")
	tp.Translated("old line", "старая строка", t0)
	tp.Translated("charge off", "просрочка", t0.Add(4*time.Minute))
	tp.Heard(words(150, "y"), t0.Add(4*time.Minute))
	tp.TopicMessages(t0.Add(5 * time.Minute))
	tp.TakeTopic("Topic: Lending in a fintech company.", t0.Add(5*time.Minute))
	m := tp.TermsMessages(t0.Add(5 * time.Minute))
	if !strings.Contains(m[0].Content, "into Russian") || !strings.Contains(m[0].Content, "blood pressure = артериальное давление") {
		t.Fatalf("system: %q", m[0].Content)
	}
	want := "Topic: Lending in a fintech company.\nGlossary: empty\nRecent lines:\n- charge off → просрочка"
	if m[1].Content != want {
		t.Fatalf("user:\n%q\nwant\n%q", m[1].Content, want)
	}
	tp.Translated("later", "позже", t0.Add(6*time.Minute))
	if strings.Contains(tp.TermsMessages(t0.Add(6 * time.Minute))[1].Content, "позже") {
		t.Fatal("a line heard after the topic request is used")
	}
}

func TestTheStartOfTheRequestOnlyChangesWithAnUpdate(t *testing.T) {
	t0 := time.Now()
	tp := NewTopic("ru")
	tp.Heard("charge off cohort "+words(200, "z"), t0)
	if tp.Context() != "" {
		t.Fatalf("context before the topic: %q", tp.Context())
	}
	tp.TopicMessages(t0)
	tp.TakeTopic("Topic: Lending.", t0)
	tp.TakeTerms("charge off = списание\ncohort = когорта", t0)
	ctx := tp.Context()
	if ctx != " The conversation: Lending. Translate these terms this way: charge off = списание; cohort = когорта." {
		t.Fatalf("context: %q", ctx)
	}
	h := []Turn{{Source: "a", Translation: "б"}}
	one := TranslateMessagesWith("ru", ctx, h, "first piece")
	tp.Heard("more words", t0.Add(time.Minute))
	tp.Translated("first piece", "первый кусок", t0.Add(time.Minute))
	two := TranslateMessagesWith("ru", tp.Context(), append(h, Turn{Source: "first piece", Translation: "первый кусок"}), "second piece")
	for i := 0; i < len(one)-1; i++ {
		if one[i] != two[i] {
			t.Fatalf("message %d changed between updates:\n%+v\n%+v", i, one[i], two[i])
		}
	}
	if !strings.HasSuffix(one[0].Content, ctx) || !strings.HasPrefix(one[0].Content, TranslateSystem("ru")) {
		t.Fatalf("system: %q", one[0].Content)
	}
}

func TestTheGrammarsAreTheOnesOfTheExperiment(t *testing.T) {
	if TopicGrammar != "root ::= \"Topic: \" [^\\n\\r*#]+\n" {
		t.Fatalf("topic grammar: %q", TopicGrammar)
	}
	ru := TermsGrammar("ru")
	if !strings.HasPrefix(ru, "root ::= \"none\" | entry (\"\\n\" entry){0,5}\n") || !strings.Contains(ru, "entry ::= src \" = \" dst\n") {
		t.Fatalf("terms grammar: %q", ru)
	}
	if !strings.Contains(ru, "dst ::= [\\u0400-\\u04FF] [") || !strings.Contains(TermsGrammar("de"), "dst ::= [a-zA-Z\\u00C0-\\u024F] [") {
		t.Fatalf("the translation starts in the script of the target: %q", ru)
	}
	for _, g := range []string{ru, TermsGrammar("de")} {
		if strings.Contains(g, "→") || strings.Contains(g, "^") || strings.Count(g, ":") != 8 {
			t.Fatalf("the grammar lets anything in: %q", g)
		}
	}
}

func TestAServiceRequestHoldsPiecesNoLongerThanATranslation(t *testing.T) {
	if TopicTimeout > TranslateTimeout {
		t.Fatalf("timeout %s, translation %s", TopicTimeout, TranslateTimeout)
	}
}

func termsDue(t *testing.T) (*Topic, time.Time) {
	t0 := time.Now()
	tp := NewTopic("ru")
	tp.Heard("the approval rate of the cohort and the charge off of the borrowing "+words(200, "x"), t0)
	tp.TopicMessages(t0)
	if !tp.TakeTopic("Topic: Lending.", t0) {
		t.Fatal("a good topic line is refused")
	}
	if s := tp.Due(t0); s != StepTerms {
		t.Fatalf("the terms are not due: %v", s)
	}
	return tp, t0
}

func TestATermsRequestThatGivesWayKeepsItsCompleteLines(t *testing.T) {
	tp, t0 := termsDue(t)
	out := tp.Answer(StepTerms, "approval rate = коэффициент одобрения\ncohort = когорта\ncharge off = спис", EndYield, t0)
	if out.Lines != 2 || out.Added != 2 || out.Again || out.GaveUp {
		t.Fatalf("outcome %+v", out)
	}
	if got := tp.Terms(); len(got) != 2 || got[1] != (Term{"cohort", "когорта"}) {
		t.Fatalf("terms %+v", got)
	}
	if s := tp.Due(t0.Add(time.Second)); s != StepNone {
		t.Fatalf("the round is not closed: %v", s)
	}
	if tp.Counts.Yielded != 1 {
		t.Fatalf("counts %+v", tp.Counts)
	}
}

func TestATermsRequestThatGivesWayWithoutALineIsAskedAgainThreeTimesAtMost(t *testing.T) {
	tp, t0 := termsDue(t)
	for i := 1; i < ServiceTries; i++ {
		out := tp.Answer(StepTerms, "approval rate = коэфф", EndYield, t0)
		if !out.Again || out.Lines != 0 || len(tp.Terms()) != 0 {
			t.Fatalf("try %d: %+v %+v", i, out, tp.Terms())
		}
		if s := tp.Due(t0); s != StepTerms {
			t.Fatalf("try %d: the terms are not due again: %v", i, s)
		}
	}
	out := tp.Answer(StepTerms, "", EndYield, t0)
	if out.Again || !out.GaveUp || tp.Counts.GaveUp != 1 || tp.Counts.Yielded != ServiceTries {
		t.Fatalf("after the last try: %+v %+v", out, tp.Counts)
	}
	if s := tp.Due(t0.Add(time.Second)); s != StepNone {
		t.Fatalf("the round is not closed: %v", s)
	}
	if s := tp.Due(t0.Add(TopicEvery)); s != StepTopic {
		t.Fatalf("the next round never comes: %v", s)
	}
}

func TestACutLastLineIsNotTakenAndAWholeAnswerIs(t *testing.T) {
	tp, t0 := termsDue(t)
	out := tp.Answer(StepTerms, "approval rate = коэффициент одобрения\ncohort = ког", EndLimit, t0)
	if out.Added != 1 || len(tp.Terms()) != 1 || tp.Terms()[0].Source != "approval rate" {
		t.Fatalf("a cut answer: %+v %+v", out, tp.Terms())
	}
	tp2, t1 := termsDue(t)
	out = tp2.Answer(StepTerms, "approval rate = коэффициент одобрения\ncohort = когорта", EndWhole, t1)
	if out.Added != 2 || out.Lines != 2 || len(tp2.Terms()) != 2 {
		t.Fatalf("a whole answer: %+v %+v", out, tp2.Terms())
	}
	tp3, t2 := termsDue(t)
	out = tp3.Answer(StepTerms, "approval rate = коэффициент одобрения\ncoh", EndTimeout, t2)
	if out.Added != 1 || out.Failed || tp3.Counts.TermsFailed != 0 {
		t.Fatalf("a timed out answer with a complete line: %+v %+v", out, tp3.Counts)
	}
	tp4, t3 := termsDue(t)
	out = tp4.Answer(StepTerms, "approval", EndTimeout, t3)
	if !out.Failed || tp4.Counts.TermsFailed != 1 || tp4.Due(t3.Add(time.Second)) != StepNone {
		t.Fatalf("a timed out answer without a line: %+v %+v", out, tp4.Counts)
	}
}

func TestATopicRequestThatGivesWayIsAskedAgain(t *testing.T) {
	t0 := time.Now()
	tp := NewTopic("ru")
	tp.Heard(words(200, "credit"), t0)
	tp.TopicMessages(t0)
	out := tp.Answer(StepTopic, "Topic: Credit risk in a lend", EndYield, t0)
	if !out.Again || tp.Text() != "" {
		t.Fatalf("outcome %+v, text %q", out, tp.Text())
	}
	if s := tp.Due(t0.Add(time.Second)); s != StepTopic {
		t.Fatalf("the topic is not asked again: %v", s)
	}
	tp.TopicMessages(t0.Add(time.Second))
	out = tp.Answer(StepTopic, "Topic: Credit risk in a lending company.", EndWhole, t0.Add(2*time.Second))
	if out.Failed || tp.Text() != "Credit risk in a lending company." || tp.Due(t0.Add(3*time.Second)) != StepTerms {
		t.Fatalf("outcome %+v, text %q", out, tp.Text())
	}
	tp.Answer(StepTerms, "", EndYield, t0)
	tp.Answer(StepTerms, "", EndYield, t0)
	if out := tp.Answer(StepTerms, "none", EndWhole, t0); out.Lines != 1 || out.Again {
		t.Fatalf("the tries of the topic step are counted for the terms: %+v", out)
	}
	if tp.Counts.GaveUp != 0 {
		t.Fatalf("counts %+v", tp.Counts)
	}
}

func TestATimedOutTopicIsAFailure(t *testing.T) {
	t0 := time.Now()
	tp := NewTopic("ru")
	tp.Heard(words(200, "credit"), t0)
	tp.TopicMessages(t0)
	out := tp.Answer(StepTopic, "Topic: Credit", EndTimeout, t0)
	if !out.Failed || tp.Counts.TopicFailed != 1 || tp.Text() != "" || tp.Due(t0.Add(time.Second)) != StepNone {
		t.Fatalf("outcome %+v, counts %+v", out, tp.Counts)
	}
}
