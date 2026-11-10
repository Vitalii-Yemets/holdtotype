package subtitles

import (
	"testing"
	"time"
)

func phraseOf(start int64, secs float64) Phrase {
	n := int(secs * rate)
	return Phrase{PCM: make([]byte, n*2), Start: start, Speech: time.Duration(secs * float64(time.Second)), At: time.Unix(start, 0)}
}

func TestOneWaitingPhraseGoesAlone(t *testing.T) {
	b := NewBacklog(rate, BacklogLimit)
	b.Push(phraseOf(1, 3.5))
	p, n := b.Take(JoinLimit)
	if n != 1 || p.Start != 1 || len(p.PCM) != int(3.5*rate)*2 {
		t.Fatalf("took %d phrases, start %d, %d bytes", n, p.Start, len(p.PCM))
	}
	if _, n := b.Take(JoinLimit); n != 0 || b.Len() != 0 || b.Waiting() != 0 {
		t.Fatal("the backlog is not empty after taking its only phrase")
	}
}

func TestWaitingPhrasesAreJoinedUpToTheLimit(t *testing.T) {
	b := NewBacklog(rate, BacklogLimit)
	for i := int64(1); i <= 8; i++ {
		b.Push(phraseOf(i, 4))
	}
	p, n := b.Take(JoinLimit)
	if n != 6 {
		t.Fatalf("joined %d phrases of 4 s under a 25 s limit, want 6", n)
	}
	if p.Start != 1 || !p.At.Equal(time.Unix(1, 0)) {
		t.Fatalf("the joined phrase starts at %d, not at the first", p.Start)
	}
	if got := p.Duration(rate); got != 24*time.Second || p.Speech != 24*time.Second {
		t.Fatalf("joined %v of audio, %v of speech", got, p.Speech)
	}
	if b.Len() != 2 || b.Waiting() != 8*time.Second {
		t.Fatalf("left %d phrases, %v", b.Len(), b.Waiting())
	}
	p, n = b.Take(JoinLimit)
	if n != 2 || p.Start != 7 {
		t.Fatalf("second take: %d phrases from %d", n, p.Start)
	}
}

func TestALongPhraseStillGoesAlone(t *testing.T) {
	b := NewBacklog(rate, BacklogLimit)
	b.Push(phraseOf(1, 12))
	b.Push(phraseOf(2, 12))
	b.Push(phraseOf(3, 2))
	_, n := b.Take(20 * time.Second)
	if n != 1 {
		t.Fatalf("took %d phrases, the second does not fit", n)
	}
}

func TestTheBacklogIsLimitedByAudioLength(t *testing.T) {
	b := NewBacklog(rate, BacklogLimit)
	var dropped []Phrase
	for i := int64(1); i <= 20; i++ {
		dropped = append(dropped, b.Push(phraseOf(i, 4))...)
	}
	if b.Waiting() != BacklogLimit || b.Len() != 15 {
		t.Fatalf("waiting %v in %d phrases", b.Waiting(), b.Len())
	}
	if len(dropped) != 5 || dropped[0].Start != 1 || dropped[4].Start != 5 {
		t.Fatalf("dropped %d, oldest first expected", len(dropped))
	}
	if p, _ := b.Take(JoinLimit); p.Start != 6 {
		t.Fatalf("the oldest kept phrase is %d", p.Start)
	}
}

func TestShortPhrasesAreNotDroppedByCount(t *testing.T) {
	b := NewBacklog(rate, BacklogLimit)
	for i := int64(1); i <= 10; i++ {
		if d := b.Push(phraseOf(i, 1)); len(d) != 0 {
			t.Fatalf("phrase %d dropped %d with only %v waiting", i, len(d), b.Waiting())
		}
	}
}
