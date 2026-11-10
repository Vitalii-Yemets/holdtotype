package subtitles

import (
	"encoding/binary"
	"testing"
	"time"
)

func wavOf(pcm []byte, rate int, extra bool) []byte {
	le := binary.LittleEndian
	var b []byte
	w32 := func(v uint32) { b = le.AppendUint32(b, v) }
	w16 := func(v uint16) { b = le.AppendUint16(b, v) }
	b = append(b, "RIFF"...)
	w32(0)
	b = append(b, "WAVE"...)
	b = append(b, "fmt "...)
	w32(16)
	w16(1)
	w16(1)
	w32(uint32(rate))
	w32(uint32(rate * 2))
	w16(2)
	w16(16)
	if extra {
		b = append(b, "LIST"...)
		w32(5)
		b = append(b, "abcde"...)
		b = append(b, 0)
	}
	b = append(b, "data"...)
	w32(uint32(len(pcm)))
	b = append(b, pcm...)
	return b
}

func TestWavDuration(t *testing.T) {
	d, ok := WavDuration(wavOf(make([]byte, 44100*2*3/4), 44100, false))
	if !ok || d != 750*time.Millisecond {
		t.Fatalf("plain wav: %v %v", d, ok)
	}
	d, ok = WavDuration(wavOf(make([]byte, 16000), 16000, true))
	if !ok || d != 500*time.Millisecond {
		t.Fatalf("wav with a list chunk: %v %v", d, ok)
	}
	if _, ok := WavDuration([]byte("not a wav at all")); ok {
		t.Fatal("garbage has a duration")
	}
}

func TestMutesSilenceTheCueWindowOnly(t *testing.T) {
	m := NewMutes(MuteBefore, MuteAfter)
	t0 := time.Date(2026, 11, 10, 14, 0, 10, 0, time.Local)
	m.Add(t0, 300*time.Millisecond)
	pcm := tone(2000, 0.3)
	end := t0.Add(time.Second)
	muted := m.Apply(pcm, end, rate)
	from := 1000 - ms(MuteBefore)
	to := 1000 + 300 + ms(MuteAfter)
	near(t, "muted samples", muted, (to-from)*rate/1000, rate/100)
	at := func(msec int) []byte { i := msec * rate / 1000 * 2; return pcm[i : i+2*rate/100] }
	if !Silent(at(from+20)) || !Silent(at(to-30)) {
		t.Fatal("the cue window was not silenced")
	}
	if Silent(at(from-30)) || Silent(at(to+20)) {
		t.Fatal("sound outside the cue window was silenced")
	}
}

func TestMutesLeaveOtherChunksAlone(t *testing.T) {
	m := NewMutes(MuteBefore, MuteAfter)
	t0 := time.Now()
	m.Add(t0, 200*time.Millisecond)
	pcm := tone(500, 0.3)
	if n := m.Apply(pcm, t0.Add(3*time.Second), rate); n != 0 || Silent(pcm) {
		t.Fatalf("a later chunk lost %d samples", n)
	}
	if n := m.Apply(pcm, t0.Add(-time.Second), rate); n != 0 {
		t.Fatalf("an earlier chunk lost %d samples", n)
	}
}

func TestMutesWithUnknownLengthUseTheGuess(t *testing.T) {
	m := NewMutes(0, 0)
	t0 := time.Now()
	m.Add(t0, 0)
	pcm := tone(3000, 0.3)
	muted := m.Apply(pcm, t0.Add(3*time.Second), rate)
	near(t, "muted ms", muted*1000/rate, ms(CueGuess), 10)
}
