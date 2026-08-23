package audiolevel

import "math"

const (
	VerdictOK      = "ok"
	VerdictSilent  = "silent"
	VerdictQuiet   = "quiet"
	VerdictClipped = "clipped"

	clipLevel   = 0.985
	clipShare   = 0.004
	quietPeak   = 0.10
	voiceFrame  = 320
	voiceEnergy = 0.02
)

type Report struct {
	Peak       float64 `json:"peak"`
	RMS        float64 `json:"rms"`
	ClipRatio  float64 `json:"clip"`
	VoiceRatio float64 `json:"voice"`
	Samples    int     `json:"samples"`
}

func sampleAt(pcm []byte, i int) float64 {
	v := int16(uint16(pcm[i]) | uint16(pcm[i+1])<<8)
	return float64(v) / 32768
}

func Analyze(pcm []byte) Report {
	var rep Report
	n := len(pcm) / 2
	if n == 0 {
		return rep
	}
	rep.Samples = n

	var sum float64
	var clipped int
	frameSum := 0.0
	frameCount := 0
	frames, voiced := 0, 0
	for i := 0; i+1 < len(pcm); i += 2 {
		v := sampleAt(pcm, i)
		a := math.Abs(v)
		if a > rep.Peak {
			rep.Peak = a
		}
		if a >= clipLevel {
			clipped++
		}
		sum += v * v
		frameSum += v * v
		frameCount++
		if frameCount == voiceFrame {
			frames++
			if math.Sqrt(frameSum/float64(frameCount)) >= voiceEnergy {
				voiced++
			}
			frameSum, frameCount = 0, 0
		}
	}
	if frameCount > 0 {
		frames++
		if math.Sqrt(frameSum/float64(frameCount)) >= voiceEnergy {
			voiced++
		}
	}
	rep.RMS = math.Sqrt(sum / float64(n))
	rep.ClipRatio = float64(clipped) / float64(n)
	if frames > 0 {
		rep.VoiceRatio = float64(voiced) / float64(frames)
	}
	return rep
}

func Verdict(r Report) string {
	switch {
	case r.Samples == 0 || r.Peak < SilenceThreshold || r.VoiceRatio == 0:
		return VerdictSilent
	case r.ClipRatio >= clipShare:
		return VerdictClipped
	case r.Peak < quietPeak:
		return VerdictQuiet
	default:
		return VerdictOK
	}
}

func DBFS(v float64) float64 {
	if v <= 0.0001 {
		return -80
	}
	db := 20 * math.Log10(v)
	if db < -80 {
		return -80
	}
	if db > 0 {
		return 0
	}
	return db
}
