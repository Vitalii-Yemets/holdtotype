package audiolevel

import (
	"math"
	"testing"
	"time"
)

const epRate = 16000

func epSamples(ms int) int { return epRate * ms / 1000 }

func epBytes(ms int) int { return epSamples(ms) * 2 }

func noiseAt(ms int, db float64, seed uint32) []byte {
	n := epSamples(ms)
	out := make([]byte, n*2)
	amp := math.Pow(10, db/20) * math.Sqrt(3)
	x := seed | 1
	for i := 0; i < n; i++ {
		x ^= x << 13
		x ^= x >> 17
		x ^= x << 5
		u := float64(x)/float64(math.MaxUint32)*2 - 1
		v := int16(math.Round(u * amp * 32767))
		out[i*2] = byte(uint16(v) & 0xFF)
		out[i*2+1] = byte(uint16(v) >> 8)
	}
	return out
}

func speechAt(ms int, db float64) []byte {
	return tone(epSamples(ms), math.Pow(10, db/20)*math.Sqrt2)
}

func zerosFor(ms int) []byte { return make([]byte, epBytes(ms)) }

func joined(parts ...[]byte) []byte {
	var out []byte
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

func newTestEndpoint(pauseMs, guardMs int) *Endpoint {
	return NewEndpoint(epRate, time.Duration(pauseMs)*time.Millisecond, time.Duration(guardMs)*time.Millisecond)
}

type feedResult struct {
	fires int
	start int
	end   int
}

func feedBy(e *Endpoint, pcm []byte, chunk int, base int) feedResult {
	r := feedResult{start: -1, end: -1}
	for off := 0; off < len(pcm); off += chunk {
		end := min(off+chunk, len(pcm))
		if e.Feed(pcm[off:end]) {
			r.fires++
			if r.end < 0 {
				r.start, r.end = base+off, base+end
			}
		}
	}
	return r
}

func TestEndpointNeverFiresWithoutSpeech(t *testing.T) {
	cases := []struct {
		name string
		pcm  []byte
	}{
		{"шум -50 dB", noiseAt(30000, -50, 1)},
		{"шум -65 dB", noiseAt(30000, -65, 2)},
		{"цифровые нули", zerosFor(30000)},
		{"громкая комната -30 dB", noiseAt(30000, -30, 3)},
		{"громкая комната -34 dB", noiseAt(30000, -34, 4)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			e := newTestEndpoint(2000, 300)
			if r := feedBy(e, c.pcm, 3200, 0); r.fires != 0 {
				t.Fatalf("сработал без речи: %+v", r)
			}
			if e.Heard() || e.Fired() {
				t.Fatalf("речь услышана в тишине: floor %.1f, threshold %.1f, speech %d ms", e.Floor(), e.Threshold(), e.SpeechMs())
			}
		})
	}
}

func TestEndpointFiresAfterPause(t *testing.T) {
	e := newTestEndpoint(2000, 300)
	head := joined(noiseAt(500, -50, 5), speechAt(1000, -20), noiseAt(1900, -50, 6))
	if r := feedBy(e, head, 640, 0); r.fires != 0 {
		t.Fatalf("сработал раньше паузы: %+v", r)
	}
	if !e.Heard() {
		t.Fatalf("речь не услышана: floor %.1f, threshold %.1f", e.Floor(), e.Threshold())
	}
	if got, want := e.SpeechEnd(), epBytes(1500); got != want {
		t.Fatalf("конец речи %d, ожидался %d", got, want)
	}
	if got := e.SpeechMs(); got != 1000 {
		t.Fatalf("речи насчитано %d ms вместо 1000", got)
	}
	r := feedBy(e, noiseAt(200, -50, 7), 640, len(head))
	if r.fires != 1 {
		t.Fatalf("не сработал после паузы: %+v", r)
	}
	want := e.SpeechEnd() + epBytes(2000)
	if d := r.end - want; d < -640 || d > 640 {
		t.Fatalf("сработал на %d байте, ожидался около %d", r.end, want)
	}
	if !e.Fired() {
		t.Fatalf("Fired() не выставлен")
	}
}

func TestEndpointGapsInSpeechDoNotFire(t *testing.T) {
	pcm := joined(
		noiseAt(500, -50, 8),
		speechAt(1000, -20), noiseAt(300, -50, 9),
		speechAt(1000, -20), noiseAt(1200, -50, 10),
		speechAt(1000, -20), noiseAt(1500, -50, 11),
	)
	e := newTestEndpoint(2000, 300)
	if r := feedBy(e, pcm, 3200, 0); r.fires != 0 {
		t.Fatalf("сработал на паузе внутри речи: %+v", r)
	}
	if !e.Heard() {
		t.Fatalf("речь не услышана")
	}
	if got, want := e.SpeechEnd(), epBytes(500+1000+300+1000+1200+1000); got != want {
		t.Fatalf("конец речи %d, ожидался %d", got, want)
	}
}

func TestEndpointGuardIgnoresCue(t *testing.T) {
	cases := []struct {
		name    string
		guardMs int
		cueMs   int
	}{
		{"со звуком 800 ms", 800, 780},
		{"без звука 300 ms", 300, 280},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			pcm := joined(speechAt(c.cueMs, -6), noiseAt(5000, -50, 12))
			e := newTestEndpoint(2000, c.guardMs)
			if r := feedBy(e, pcm, 3200, 0); r.fires != 0 {
				t.Fatalf("звук старта принят за речь: %+v", r)
			}
			if e.Heard() {
				t.Fatalf("звук старта взвёл детектор")
			}
			if f := e.Floor(); f > -48 {
				t.Fatalf("звук старта поднял уровень шума до %.1f", f)
			}
		})
	}
}

func TestEndpointClickDoesNotArm(t *testing.T) {
	cases := []struct {
		name   string
		leadMs int
	}{
		{"по границе кадра", 1000},
		{"со сдвигом", 1010},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			pcm := joined(noiseAt(c.leadMs, -50, 13), speechAt(60, -10), noiseAt(5000, -50, 14))
			e := newTestEndpoint(2000, 300)
			if r := feedBy(e, pcm, 3200, 0); r.fires != 0 {
				t.Fatalf("щелчок остановил запись: %+v", r)
			}
			if e.Heard() {
				t.Fatalf("щелчок взвёл детектор: %d ms речи", e.SpeechMs())
			}
		})
	}
}

func TestEndpointNoiseFloor(t *testing.T) {
	cases := []struct {
		name     string
		floorDB  float64
		speechDB float64
		fires    bool
	}{
		{"шум -40, речь -20", -40, -20, true},
		{"шум -30 без речи", -30, 0, false},
		{"шум -30, речь -20", -30, -20, true},
		{"шум -60, речь -45", -60, -45, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			pcm := noiseAt(30000, c.floorDB, 15)
			if c.speechDB != 0 {
				pcm = joined(noiseAt(500, c.floorDB, 16), speechAt(1000, c.speechDB), noiseAt(3000, c.floorDB, 17))
			}
			e := newTestEndpoint(2000, 300)
			r := feedBy(e, pcm, 3200, 0)
			if r.fires > 1 || (r.fires == 1) != c.fires {
				t.Fatalf("срабатываний %d, ожидалось %v: floor %.1f, threshold %.1f, speech %d ms", r.fires, c.fires, e.Floor(), e.Threshold(), e.SpeechMs())
			}
			if c.speechDB == 0 && e.Heard() {
				t.Fatalf("шум принят за речь: floor %.1f, threshold %.1f", e.Floor(), e.Threshold())
			}
			if f := e.Floor(); math.Abs(f-c.floorDB) > 1.5 {
				t.Fatalf("уровень шума %.1f вместо %.0f", f, c.floorDB)
			}
		})
	}
}

func TestEndpointThresholdLimits(t *testing.T) {
	cases := []struct {
		floorDB float64
		lo, hi  float64
	}{
		{-75, -50, -50},
		{-60, -49, -47.5},
		{-40, -30, -30},
		{-30, -23, -21.5},
	}
	for _, c := range cases {
		e := newTestEndpoint(2000, 300)
		feedBy(e, noiseAt(1000, c.floorDB, 18), 640, 0)
		if th := e.Threshold(); th < c.lo || th > c.hi {
			t.Fatalf("при шуме %.0f порог %.1f, ожидался в %.1f..%.1f", c.floorDB, th, c.lo, c.hi)
		}
	}
}

func TestEndpointZeroFramesDoNotStickFloor(t *testing.T) {
	cases := []struct {
		name string
		pcm  []byte
	}{
		{"нули перед шумом и речью", joined(zerosFor(1000), noiseAt(300, -45, 19), speechAt(1000, -20), noiseAt(3000, -45, 20))},
		{"нули вперемешку с шумом", joined(zerosFor(200), noiseAt(100, -45, 21), zerosFor(400), noiseAt(300, -45, 22), speechAt(1000, -20), noiseAt(3000, -45, 23))},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			e := newTestEndpoint(2000, 300)
			if r := feedBy(e, c.pcm, 3200, 0); r.fires != 1 {
				t.Fatalf("срабатываний %d вместо 1: floor %.1f, threshold %.1f, speech %d ms", r.fires, e.Floor(), e.Threshold(), e.SpeechMs())
			}
			if f := e.Floor(); f < -47 || f > -43 {
				t.Fatalf("нули утянули уровень шума до %.1f", f)
			}
		})
	}
}

func TestEndpointZerosAfterSpeechCountAsSilence(t *testing.T) {
	pcm := joined(noiseAt(500, -50, 24), speechAt(1000, -20), zerosFor(2100))
	e := newTestEndpoint(2000, 300)
	if r := feedBy(e, pcm, 3200, 0); r.fires != 1 {
		t.Fatalf("выключенный микрофон после речи не дал остановки: %+v", r)
	}
}

func TestEndpointSoftTailKeepsAlive(t *testing.T) {
	e := newTestEndpoint(2000, 300)
	head := joined(noiseAt(500, -50, 25), speechAt(1000, -20), noiseAt(3000, -40, 26), noiseAt(1900, -50, 27))
	if r := feedBy(e, head, 640, 0); r.fires != 0 {
		t.Fatalf("тихий хвост речи оборван: %+v, floor %.1f, threshold %.1f", r, e.Floor(), e.Threshold())
	}
	if got, want := e.SpeechEnd(), epBytes(4500); got != want {
		t.Fatalf("конец речи %d, ожидался %d", got, want)
	}
	if r := feedBy(e, noiseAt(200, -50, 28), 640, len(head)); r.fires != 1 {
		t.Fatalf("после хвоста не сработал: %+v", r)
	}
}

func TestEndpointChunkSizesAgree(t *testing.T) {
	pcm := joined(noiseAt(500, -50, 29), speechAt(1000, -20), noiseAt(300, -50, 30), speechAt(700, -22), noiseAt(3000, -50, 31))
	ref := newTestEndpoint(2000, 300)
	want := feedBy(ref, pcm, 640, 0)
	if want.fires != 1 {
		t.Fatalf("эталон не сработал: %+v", want)
	}
	for _, chunk := range []int{160, 3200, 1, 7, 333, 1001, 4799} {
		e := newTestEndpoint(2000, 300)
		got := feedBy(e, pcm, chunk, 0)
		if got.fires != 1 {
			t.Fatalf("кусок %d: срабатываний %d", chunk, got.fires)
		}
		if got.start >= want.end || got.end < want.end {
			t.Fatalf("кусок %d: сработал в %d..%d, эталон на %d", chunk, got.start, got.end, want.end)
		}
		if e.SpeechEnd() != ref.SpeechEnd() || e.SpeechMs() != ref.SpeechMs() || e.Floor() != ref.Floor() {
			t.Fatalf("кусок %d: end %d/%d, speech %d/%d, floor %v/%v", chunk, e.SpeechEnd(), ref.SpeechEnd(), e.SpeechMs(), ref.SpeechMs(), e.Floor(), ref.Floor())
		}
	}
}

func TestEndpointFiresOnce(t *testing.T) {
	e := newTestEndpoint(2000, 300)
	first := joined(noiseAt(500, -50, 32), speechAt(1000, -20), noiseAt(2500, -50, 33))
	if r := feedBy(e, first, 3200, 0); r.fires != 1 {
		t.Fatalf("первое срабатывание: %+v", r)
	}
	end := e.SpeechEnd()
	again := joined(speechAt(1000, -20), noiseAt(3000, -50, 34))
	if r := feedBy(e, again, 3200, len(first)); r.fires != 0 {
		t.Fatalf("сработал повторно: %+v", r)
	}
	if e.Feed(nil) || !e.Fired() || e.SpeechEnd() != end {
		t.Fatalf("состояние после срабатывания изменилось: fired %v, end %d/%d", e.Fired(), e.SpeechEnd(), end)
	}
}

func TestEndpointPauseLength(t *testing.T) {
	for _, pauseMs := range []int{1000, 3500, 10000} {
		e := newTestEndpoint(pauseMs, 300)
		pcm := joined(noiseAt(500, -50, 35), speechAt(1000, -20), noiseAt(pauseMs+500, -50, 36))
		r := feedBy(e, pcm, 640, 0)
		if r.fires != 1 {
			t.Fatalf("пауза %d ms: срабатываний %d", pauseMs, r.fires)
		}
		if want := epBytes(1500 + pauseMs); r.end != want {
			t.Fatalf("пауза %d ms: сработал на %d байте вместо %d", pauseMs, r.end, want)
		}
	}
}

func TestValidPauseMs(t *testing.T) {
	cases := []struct {
		ms int
		ok bool
	}{
		{-1, false}, {0, false}, {999, false}, {1000, true}, {DefaultPauseMs, true}, {10000, true}, {10001, false},
	}
	for _, c := range cases {
		if got := ValidPauseMs(c.ms); got != c.ok {
			t.Errorf("ValidPauseMs(%d) = %v; want %v", c.ms, got, c.ok)
		}
	}
}
