package subtitles

import (
	"context"
	"strings"
	"testing"
	"time"
)

const sampleCall = "=== Запуск субтитров 16.11.2026 в 14:09:41 · язык речи: es · модель: Nemotron 3.5 · строк: 1 · слов: 2 ===\n\n" +
	"[14:11:14] Segura, por\n\n" +
	"=== Запуск субтитров 16.11.2026 в 19:26:55 · язык речи: en · модель: Nemotron 3.5 · строк: 2 · слов: 13 ===\r\n\r\n" +
	"[19:27:00] Repaid the previous credit and goes into the next one\r\n" +
	"[19:27:04] Which are quite important.\r\n"

func TestTheCallFileIsReadByRun(t *testing.T) {
	lines, lang := ParseCall(sampleCall, "19:26:55")
	if lang != "en" || len(lines) != 2 || lines[0].At != 19*time.Hour+27*time.Minute || len(lines[0].Words) != 10 || lines[1].Words[3] != "important." {
		t.Fatalf("got %q %+v", lang, lines)
	}
	all, lang := ParseCall(sampleCall, "all")
	if lang != "es" || len(all) != 3 {
		t.Fatalf("all runs: %q %d", lang, len(all))
	}
	if none, _ := ParseCall(sampleCall, "10:00:00"); len(none) != 0 {
		t.Fatalf("an unknown run gives lines: %+v", none)
	}
}

func TestTheCallIsCutLikeTheProgramCutsLiveSpeech(t *testing.T) {
	lines, lang := ParseCall(sampleCall, "19:26:55")
	base := time.Date(2026, 11, 16, 0, 0, 0, 0, time.Local)
	pieces := PhraseCall(lines, lang, base)
	var texts []string
	words := 0
	for _, p := range pieces {
		texts = append(texts, p.Text)
		words += p.Words
		if p.LastWord.After(p.Release) {
			t.Fatalf("a piece is released before its last word: %+v", p)
		}
	}
	if words != 14 || strings.Join(texts, " ") != "Repaid the previous credit and goes into the next one Which are quite important." {
		t.Fatalf("pieces: %q", texts)
	}
	first := pieces[0]
	if !first.LastWord.Equal(base.Add(19*time.Hour+27*time.Minute)) || first.Release.Sub(first.LastWord) != PhraseQuiet {
		t.Fatalf("the first line ends at its time and goes at the pause: %+v", first)
	}
}

func checkPieces(t0 time.Time, n int, every time.Duration, words int) []CheckPiece {
	var out []CheckPiece
	for i := 0; i < n; i++ {
		at := t0.Add(time.Duration(i) * every)
		out = append(out, CheckPiece{Seq: i + 1, Text: strings.TrimSpace(strings.Repeat("word ", words)), Words: words, Release: at, LastWord: at})
	}
	return out
}

func TestServiceRequestsOnlyRunWhenNoPieceWaits(t *testing.T) {
	t0 := time.Date(2026, 11, 16, 19, 27, 0, 0, time.Local)
	pieces := checkPieces(t0, 60, 3*time.Second, 10)
	tp := NewTopic("ru")
	calls := CheckCalls{
		Translate: func(b Batch, prior []Turn, context string) (string, Verdict, time.Duration, error) {
			return "перевод", TrOK, 500 * time.Millisecond, nil
		},
		Service: func(step Step, msgs []Message, limit time.Duration) (string, ServiceEnd, time.Duration, error) {
			if step == StepTopic {
				return "Topic: Something.", EndWhole, 2500 * time.Millisecond, nil
			}
			return "word = слово", EndWhole, 2500 * time.Millisecond, nil
		},
	}
	ev := RunCheck(pieces, "ru", tp, calls)
	services := 0
	for k, e := range ev {
		if e.Kind == "piece" || e.Kind == "stale" {
			continue
		}
		services++
		for _, p := range pieces {
			if p.Release.After(e.Start) {
				continue
			}
			done := false
			for _, f := range ev[:k] {
				if f.Kind == "piece" && containsSeq(f.Seqs, p.Seq) && !f.End.After(e.Start) {
					done = true
				}
			}
			if !done {
				t.Fatalf("%s request at %s started while piece %d waited", e.Kind, e.Start.Format("15:04:05.000"), p.Seq)
			}
		}
		if k+1 < len(ev) && ev[k+1].Kind == "piece" {
			next := ev[k+1]
			if next.Start.After(e.End) || next.Start.Sub(next.Ready) > e.Duration {
				t.Fatalf("a piece waited longer than the service request: %+v after %+v", next, e)
			}
		}
	}
	if services != 2 || tp.Updates() != 1 || len(tp.Terms()) != 1 {
		t.Fatalf("services %d, updates %d, terms %+v", services, tp.Updates(), tp.Terms())
	}
	last := ev[len(ev)-1]
	if last.Kind != "piece" || last.Context != " The conversation: Something. Translate these terms this way: word = слово." {
		t.Fatalf("the last piece goes without the context: %+v", last)
	}
}

func TestPiecesQueuedBehindASlowRequestAreJoined(t *testing.T) {
	t0 := time.Now()
	pieces := checkPieces(t0, 4, 100*time.Millisecond, 5)
	var batches [][]int
	calls := CheckCalls{Translate: func(b Batch, prior []Turn, context string) (string, Verdict, time.Duration, error) {
		batches = append(batches, b.Seqs)
		return "x", TrOK, time.Second, nil
	}}
	RunCheck(pieces, "ru", nil, calls)
	if len(batches) != 2 || len(batches[0]) != 1 || len(batches[1]) != 3 {
		t.Fatalf("batches: %v", batches)
	}
}

func containsSeq(seqs []int, s int) bool {
	for _, x := range seqs {
		if x == s {
			return true
		}
	}
	return false
}

func TestAServiceRequestGivesWayToTheNextPieceAndKeepsItsLines(t *testing.T) {
	t0 := time.Date(2026, 11, 16, 19, 27, 0, 0, time.Local)
	pieces := checkPieces(t0, 60, 3*time.Second, 10)
	tp := NewTopic("ru")
	calls := CheckCalls{
		Translate: func(b Batch, prior []Turn, context string) (string, Verdict, time.Duration, error) {
			return "перевод", TrOK, 500 * time.Millisecond, nil
		},
		Service: func(step Step, msgs []Message, limit time.Duration) (string, ServiceEnd, time.Duration, error) {
			if step == StepTopic {
				return "Topic: Something.", EndWhole, time.Second, nil
			}
			if limit > 0 && limit < 4*time.Second {
				return "word = слово\nwor", EndYield, limit, context.Canceled
			}
			return "word = слово", EndWhole, 4 * time.Second, nil
		},
	}
	ev := RunCheck(pieces, "ru", tp, calls)
	yielded := 0
	for _, e := range ev {
		if e.Kind != "terms" {
			continue
		}
		if e.Yielded {
			yielded++
		}
		if e.Hold > 0 {
			t.Fatalf("a piece waited for a request that gave way: %+v", e)
		}
		if e.Yielded && e.Lines != 1 {
			t.Fatalf("the complete line is not counted: %+v", e)
		}
	}
	if yielded == 0 || len(tp.Terms()) != 1 || tp.Counts.Yielded != yielded {
		t.Fatalf("yielded %d, terms %+v, counts %+v", yielded, tp.Terms(), tp.Counts)
	}
}

func TestADraftOnlyFillsTheTimeBeforeTheSentenceCloses(t *testing.T) {
	t0 := time.Date(2026, 11, 18, 10, 0, 0, 0, time.Local)
	var pieces []CheckPiece
	at := t0
	for s := 1; s <= 3; s++ {
		p := CheckPiece{Seq: s}
		n := 6
		if s == 2 {
			n = 30
		}
		for k := 0; k < n; k++ {
			p.Parts = append(p.Parts, "word")
			p.Times = append(p.Times, at)
			at = at.Add(400 * time.Millisecond)
		}
		p.Text = strings.Join(p.Parts, " ")
		p.Words = n
		p.LastWord = p.Times[n-1]
		p.Release = p.LastWord.Add(500 * time.Millisecond)
		pieces = append(pieces, p)
	}
	drafts := 0
	calls := CheckCalls{
		Translate: func(b Batch, prior []Turn, context string) (string, Verdict, time.Duration, error) {
			return "перевод", TrOK, 300 * time.Millisecond, nil
		},
		Draft: func(text string, prior []Turn, about string, limit time.Duration) (string, Verdict, ServiceEnd, time.Duration, error) {
			drafts++
			if limit < time.Second {
				return "", TrEmpty, EndYield, limit, context.DeadlineExceeded
			}
			return "черновик", TrOK, EndWhole, 700 * time.Millisecond, nil
		},
	}
	ev := RunCheck(pieces, "ru", nil, calls)
	var draftEv []CheckEvent
	for _, e := range ev {
		if e.Kind == "draft" {
			draftEv = append(draftEv, e)
			if e.Seqs[0] != 2 {
				t.Fatalf("a draft for a short sentence: %+v", e)
			}
			if e.End.After(pieces[1].Release) {
				t.Fatalf("a draft ran past the end of its sentence: %+v", e)
			}
		}
		if e.Kind == "piece" && e.Start.After(pieces[e.Seqs[0]-1].Release.Add(time.Millisecond)) && e.Seqs[0] != 3 {
			t.Fatalf("a whole sentence waited: %+v", e)
		}
	}
	if len(draftEv) < 2 || draftEv[0].Lines < 6 || !draftEv[0].Start.Equal(pieces[1].Times[0].Add(4*time.Second)) {
		t.Fatalf("drafts %+v", draftEv)
	}
	if draftEv[1].Lines < draftEv[0].Lines+6 || draftEv[1].Start.Sub(draftEv[0].Start) < 3*time.Second {
		t.Fatalf("the update came too soon: %+v", draftEv)
	}
}
