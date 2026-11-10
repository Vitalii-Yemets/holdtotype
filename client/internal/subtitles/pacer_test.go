package subtitles

import (
	"testing"
	"time"
)

func TestNothingCapturedIsToppedUpThenPaused(t *testing.T) {
	start := time.Date(2026, 11, 10, 14, 0, 0, 0, time.Local)
	p := NewPacer(rate, start)
	total := 0
	for i := 1; i <= 100; i++ {
		out := p.Tick(start.Add(time.Duration(i)*100*time.Millisecond), nil)
		if !Silent(out) {
			t.Fatal("top-up is not silence")
		}
		total += len(out) / 2
	}
	if !p.Suspended() || p.Pauses() != 1 {
		t.Fatalf("suspended=%v pauses=%d after ten seconds of nothing", p.Suspended(), p.Pauses())
	}
	near(t, "samples added", total, int(SilenceLimit/time.Millisecond)*rate/1000, rate*TickMs/1000)
	near(t, "added ms", ms(p.Added()), int(SilenceLimit/time.Millisecond), TickMs)
}

func TestRealTimeCaptureNeedsNoTopUp(t *testing.T) {
	start := time.Now()
	p := NewPacer(rate, start)
	chunk := tone(100, 0.3)
	for i := 1; i <= 50; i++ {
		out := p.Tick(start.Add(time.Duration(i)*100*time.Millisecond), chunk)
		if len(out) != len(chunk) {
			t.Fatalf("tick %d: %d bytes out for %d in", i, len(out), len(chunk))
		}
	}
	if p.Added() != 0 {
		t.Fatalf("%v of silence added to a continuous stream", p.Added())
	}
}

func TestJitterIsAbsorbed(t *testing.T) {
	start := time.Now()
	p := NewPacer(rate, start)
	for i := 1; i <= 40; i++ {
		c := tone(50, 0.3)
		if i%2 == 0 {
			c = tone(150, 0.3)
		}
		p.Tick(start.Add(time.Duration(i)*100*time.Millisecond), c)
	}
	if p.Added() != 0 {
		t.Fatalf("jitter of 50 ms added %v of silence", p.Added())
	}
}

func TestTopUpGoesBeforeTheSound(t *testing.T) {
	start := time.Now()
	p := NewPacer(rate, start)
	p.Tick(start.Add(100*time.Millisecond), nil)
	out := p.Tick(start.Add(500*time.Millisecond), tone(100, 0.3))
	if len(out) <= rate/10*2 {
		t.Fatalf("no top-up before a late chunk: %d bytes", len(out))
	}
	pad := len(out) - rate/10*2
	if !Silent(out[:pad]) || Silent(out[pad:]) {
		t.Fatal("the silence is not in front of the captured sound")
	}
	near(t, "fed ms", ms(SamplesDuration(p.Fed(), rate)), 500-SlackMs, 1)
}

func TestSoundResumesAfterAPause(t *testing.T) {
	start := time.Now()
	p := NewPacer(rate, start)
	var now time.Time
	for i := 1; i <= 60; i++ {
		now = start.Add(time.Duration(i) * 100 * time.Millisecond)
		p.Tick(now, nil)
	}
	fedAtPause := p.Fed()
	for i := 61; i <= 600; i++ {
		now = start.Add(time.Duration(i) * 100 * time.Millisecond)
		if out := p.Tick(now, quiet(100)); out != nil {
			t.Fatal("silence was passed on during the pause")
		}
	}
	if p.Fed() != fedAtPause {
		t.Fatal("the timeline grew during the pause")
	}
	chunk := tone(100, 0.3)
	out := p.Tick(now, chunk)
	if len(out) != len(chunk) || p.Suspended() {
		t.Fatalf("sound after the pause: %d bytes, suspended=%v", len(out), p.Suspended())
	}
	if got := p.TimeAt(p.Fed()); !got.Equal(now) {
		t.Fatalf("after resuming the timeline ends at %v, want %v", got.Sub(start), now.Sub(start))
	}
	next := p.Tick(now.Add(time.Second), nil)
	near(t, "top-up after resume ms", len(next)/2*1000/rate, 1000-SlackMs, 1)
	if p.Pauses() != 1 {
		t.Fatalf("pauses=%d", p.Pauses())
	}
}

func TestSilentCapturedStreamAlsoPauses(t *testing.T) {
	start := time.Now()
	p := NewPacer(rate, start)
	for i := 1; i <= 45; i++ {
		p.Tick(start.Add(time.Duration(i)*100*time.Millisecond), quiet(100))
	}
	if !p.Suspended() {
		t.Fatal("a stream of digital silence did not pause the feed")
	}
	if p.Added() != 0 {
		t.Fatalf("%v added although the capture covered the time", p.Added())
	}
}
