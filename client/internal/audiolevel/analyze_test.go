package audiolevel

import (
	"math"
	"testing"
)

func square(samples int, amplitude float64) []byte {
	out := make([]byte, samples*2)
	for i := 0; i < samples; i++ {
		v := int16(amplitude * 32767)
		if i%2 == 0 {
			v = -v
		}
		out[i*2] = byte(uint16(v) & 0xFF)
		out[i*2+1] = byte(uint16(v) >> 8)
	}
	return out
}

func TestAnalyzeEmpty(t *testing.T) {
	rep := Analyze(nil)
	if rep.Samples != 0 || Verdict(rep) != VerdictSilent {
		t.Fatalf("пустая запись не признана тишиной: %+v", rep)
	}
}

func TestVerdictSilenceForDigitalZero(t *testing.T) {
	rep := Analyze(make([]byte, 16000))
	if Verdict(rep) != VerdictSilent {
		t.Fatalf("цифровая тишина не признана тишиной: %+v", rep)
	}
}

func TestVerdictNormalSpeechLevel(t *testing.T) {
	rep := Analyze(tone(16000, 0.4))
	if v := Verdict(rep); v != VerdictOK {
		t.Fatalf("нормальный уровень принят за %s: %+v", v, rep)
	}
	if rep.VoiceRatio < 0.9 {
		t.Fatalf("речь не распознана как речь: %+v", rep)
	}
}

func TestVerdictQuiet(t *testing.T) {
	rep := Analyze(tone(16000, 0.05))
	if v := Verdict(rep); v != VerdictQuiet {
		t.Fatalf("тихая запись принята за %s: %+v", v, rep)
	}
}

func TestVerdictClipped(t *testing.T) {
	rep := Analyze(square(16000, 1.0))
	if v := Verdict(rep); v != VerdictClipped {
		t.Fatalf("перегруз не пойман: %s %+v", v, rep)
	}
	if rep.ClipRatio < 0.9 {
		t.Fatalf("доля обрезанных отсчётов посчитана неверно: %+v", rep)
	}
}

func TestClippingWinsOverQuietOnlyWhenReal(t *testing.T) {
	loud := tone(16000, 0.9)
	rep := Analyze(loud)
	if v := Verdict(rep); v != VerdictOK {
		t.Fatalf("громкая, но не обрезанная запись принята за %s: %+v", v, rep)
	}
}

func TestVoiceRatioCountsSilentFrames(t *testing.T) {
	speech := tone(8000, 0.4)
	quiet := make([]byte, 8000*2)
	rep := Analyze(append(speech, quiet...))
	if rep.VoiceRatio < 0.4 || rep.VoiceRatio > 0.6 {
		t.Fatalf("половина тишины даёт %v вместо ~0.5", rep.VoiceRatio)
	}
}

func TestDBFS(t *testing.T) {
	if db := DBFS(1); math.Abs(db) > 0.01 {
		t.Fatalf("полная шкала должна быть 0 дБ, получено %v", db)
	}
	if db := DBFS(0.5); math.Abs(db+6.02) > 0.1 {
		t.Fatalf("половина шкалы должна быть около -6 дБ, получено %v", db)
	}
	if db := DBFS(0); db != -80 {
		t.Fatalf("тишина должна упираться в -80 дБ, получено %v", db)
	}
}
