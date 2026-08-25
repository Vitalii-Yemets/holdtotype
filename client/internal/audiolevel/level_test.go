package audiolevel

import (
	"math"
	"testing"
)

func tone(samples int, amplitude float64) []byte {
	out := make([]byte, samples*2)
	for i := 0; i < samples; i++ {
		v := int16(amplitude * 32767 * math.Sin(float64(i)/8))
		out[i*2] = byte(uint16(v) & 0xFF)
		out[i*2+1] = byte(uint16(v) >> 8)
	}
	return out
}

func TestPeak(t *testing.T) {
	if p := Peak(nil); p != 0 {
		t.Errorf("Peak(nil) = %v; want 0", p)
	}
	if p := Peak(tone(500, 0.5)); p < 0.45 || p > 0.55 {
		t.Errorf("Peak(half amplitude) = %v; want ~0.5", p)
	}
	if p := Peak(tone(500, 1.0)); p < 0.9 {
		t.Errorf("Peak(full amplitude) = %v; want ~1.0", p)
	}
}

func TestIsSilent(t *testing.T) {
	cases := []struct {
		name string
		pcm  []byte
		want bool
	}{
		{"empty", nil, true},
		{"digital silence", make([]byte, 4000), true},
		{"faint noise below threshold", tone(4000, 0.005), true},
		{"quiet speech above threshold", tone(4000, 0.05), false},
		{"normal speech", tone(4000, 0.4), false},
	}
	for _, c := range cases {
		if got := IsSilent(c.pcm); got != c.want {
			t.Errorf("IsSilent(%s) = %v; want %v (peak %.4f)", c.name, got, c.want, Peak(c.pcm))
		}
	}
}

func TestHeardOpensUpAQuietMicrophone(t *testing.T) {
	if got := Heard(HeardFloor); got != 0 {
		t.Errorf("Heard(%v) = %v, want a quiet room to read zero", HeardFloor, got)
	}
	if got := Heard(0.004); got != 0 {
		t.Errorf("Heard(0.004) = %v, want zero below the floor", got)
	}
	if got := Heard(0.06); got < 0.5 || got > 0.7 {
		t.Errorf("Heard(0.06) = %v, want a little over half for ordinary speech", got)
	}
	if got := Heard(0.15); got < 0.75 {
		t.Errorf("Heard(0.15) = %v, want a loud phrase near the top", got)
	}
	if got := Heard(HeardTop); got != 1 {
		t.Errorf("Heard(%v) = %v, want a full meter at the top", HeardTop, got)
	}
	if got := Heard(0.9); got != 1 {
		t.Errorf("Heard(0.9) = %v, want the meter to stay full above the top", got)
	}
	prev := 0.0
	for v := HeardFloor + 0.001; v < HeardTop; v += 0.005 {
		got := Heard(v)
		if got < prev {
			t.Fatalf("Heard fell from %v to %v at %v", prev, got, v)
		}
		prev = got
	}
}
