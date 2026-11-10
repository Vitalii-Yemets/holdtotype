package subtitles

import (
	"math"
	"testing"
	"time"
)

const rate = 16000

func tone(ms int, amp float64) []byte {
	n := rate * ms / 1000
	out := make([]byte, n*2)
	for i := 0; i < n; i++ {
		v := int16(amp * 32767 * math.Sin(2*math.Pi*440*float64(i)/rate))
		out[i*2] = byte(uint16(v))
		out[i*2+1] = byte(uint16(v) >> 8)
	}
	return out
}

func quiet(ms int) []byte { return make([]byte, rate*ms/1000*2) }

func join(parts ...[]byte) []byte {
	var out []byte
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

func ms(d time.Duration) int { return int(d / time.Millisecond) }

func feedAll(s *Segmenter, pcm []byte, chunk int) []Phrase {
	var out []Phrase
	for len(pcm) > 0 {
		n := chunk
		if n > len(pcm) {
			n = len(pcm)
		}
		out = append(out, s.Feed(pcm[:n])...)
		pcm = pcm[n:]
	}
	return out
}

func near(t *testing.T, what string, got, want, tol int) {
	t.Helper()
	if got < want-tol || got > want+tol {
		t.Fatalf("%s: got %d, want %d ± %d", what, got, want, tol)
	}
}

func TestSilenceGivesNothingAndHoldsNoMemory(t *testing.T) {
	s := NewSegmenter(rate)
	for i := 0; i < 600; i++ {
		if got := s.Feed(quiet(100)); len(got) != 0 {
			t.Fatalf("silence produced %d phrases", len(got))
		}
	}
	if b := s.Buffered(); b > s.preMax+s.frameBytes {
		t.Fatalf("a minute of silence keeps %d bytes", b)
	}
}

func TestOnePhraseKeepsLeadAndTail(t *testing.T) {
	s := NewSegmenter(rate)
	got := feedAll(s, join(quiet(1000), tone(1000, 0.3), quiet(1000)), 3200)
	if len(got) != 1 {
		t.Fatalf("got %d phrases, want 1", len(got))
	}
	p := got[0]
	near(t, "start sample", int(p.Start), rate, rate*FrameMs/1000)
	near(t, "length ms", ms(p.Duration(rate)), PrerollMs+1000+TailMs, 2*FrameMs)
	near(t, "speech ms", ms(p.Speech), 1000, 2*FrameMs)
	if p.Cut {
		t.Fatal("a phrase that ended on a pause is marked as cut")
	}
	lead := p.PCM[:rate*PrerollMs/1000*2-2*FrameMs*rate/1000*2]
	if !Silent(lead) {
		t.Fatal("the lead before the speech is not the quiet audio that preceded it")
	}
}

func TestShortPauseKeepsOnePhrase(t *testing.T) {
	s := NewSegmenter(rate)
	got := feedAll(s, join(quiet(500), tone(700, 0.3), quiet(400), tone(700, 0.3), quiet(1000)), 1600)
	if len(got) != 1 {
		t.Fatalf("a 400 ms pause split the phrase: %d phrases", len(got))
	}
}

func TestLongPauseSplits(t *testing.T) {
	s := NewSegmenter(rate)
	got := feedAll(s, join(quiet(500), tone(700, 0.3), quiet(800), tone(700, 0.3), quiet(1000)), 1600)
	if len(got) != 2 {
		t.Fatalf("an 800 ms pause gave %d phrases, want 2", len(got))
	}
	if got[1].Start <= got[0].Start {
		t.Fatal("phrases are out of order")
	}
	near(t, "second start", int(got[1].Start), rate*(500+700+800)/1000, rate*FrameMs/1000)
}

func TestClickIsDropped(t *testing.T) {
	s := NewSegmenter(rate)
	got := feedAll(s, join(quiet(500), tone(200, 0.5), quiet(1500)), 3200)
	if len(got) != 0 {
		t.Fatalf("a 200 ms click became %d phrases", len(got))
	}
	if s.InPhrase() {
		t.Fatal("the segmenter stayed inside a phrase after the click")
	}
}

func TestLongSoundIsCutAtTheQuietestSpot(t *testing.T) {
	s := NewSegmenter(rate)
	music := join(tone(10500, 0.3), tone(100, 0.12), tone(9400, 0.3))
	got := feedAll(s, join(quiet(300), music), 3200)
	if len(got) != 1 {
		t.Fatalf("twenty seconds of sound gave %d phrases before the end, want 1", len(got))
	}
	p := got[0]
	if !p.Cut {
		t.Fatal("the forced cut is not marked")
	}
	near(t, "cut point ms", ms(p.Duration(rate)), 300+10500-PrerollMs+60, 60)
	if d := ms(p.Duration(rate)); d > MaxPhraseMs {
		t.Fatalf("a phrase of %d ms is longer than the limit", d)
	}
	rest := s.Flush()
	if len(rest) != 1 {
		t.Fatalf("flush gave %d phrases, want the remainder", len(rest))
	}
	cutAt := p.Start - int64(rate*PrerollMs/1000) + int64(len(p.PCM)/2)
	if rest[0].Start < cutAt || rest[0].Start > cutAt+int64(rate*FrameMs/1000) {
		t.Fatal("the remainder starts before the cut")
	}
}

func TestCutWithoutDipStaysWithinLimit(t *testing.T) {
	s := NewSegmenter(rate)
	got := feedAll(s, join(quiet(300), tone(30000, 0.3)), 3200)
	if len(got) < 2 {
		t.Fatalf("thirty seconds of sound gave %d phrases, want at least 2", len(got))
	}
	for i, p := range got {
		if d := ms(p.Duration(rate)); d > MaxPhraseMs || d < MaxPhraseMs-CutWindowMs-FrameMs {
			t.Fatalf("phrase %d is %d ms", i, d)
		}
	}
}

func TestChunkSizeDoesNotMatter(t *testing.T) {
	pcm := join(quiet(600), tone(900, 0.3), quiet(900), tone(500, 0.3), quiet(1200))
	a := feedAll(NewSegmenter(rate), pcm, 3200)
	b := feedAll(NewSegmenter(rate), pcm, 77)
	if len(a) != len(b) || len(a) != 2 {
		t.Fatalf("chunks of 3200 gave %d phrases, of 77 gave %d", len(a), len(b))
	}
	for i := range a {
		if a[i].Start != b[i].Start || len(a[i].PCM) != len(b[i].PCM) {
			t.Fatalf("phrase %d differs: %d/%d vs %d/%d", i, a[i].Start, len(a[i].PCM), b[i].Start, len(b[i].PCM))
		}
	}
}

func TestFlushEmitsTheOpenPhrase(t *testing.T) {
	s := NewSegmenter(rate)
	if got := feedAll(s, join(quiet(300), tone(800, 0.3)), 3200); len(got) != 0 {
		t.Fatalf("the phrase ended before its pause: %d", len(got))
	}
	got := s.Flush()
	if len(got) != 1 {
		t.Fatalf("flush gave %d phrases", len(got))
	}
	if s.InPhrase() {
		t.Fatal("still inside a phrase after flush")
	}
	if more := s.Flush(); len(more) != 0 {
		t.Fatal("a second flush repeated the phrase")
	}
}

func TestQuietBackgroundRaisesNothing(t *testing.T) {
	s := NewSegmenter(rate)
	got := feedAll(s, join(tone(3000, 0.002), tone(800, 0.3), tone(1000, 0.002)), 3200)
	if len(got) != 1 {
		t.Fatalf("speech over a quiet hum gave %d phrases, want 1", len(got))
	}
}

func TestASegmenterStartedLaterCountsFromItsStart(t *testing.T) {
	s := NewSegmenter(rate)
	s.StartAt(int64(rate) * 60)
	got := feedAll(s, join(quiet(1000), tone(1000, 0.3), quiet(1000)), 3200)
	if len(got) != 1 {
		t.Fatalf("got %d phrases", len(got))
	}
	near(t, "start sample", int(got[0].Start), rate*61, rate*FrameMs/1000)
}
