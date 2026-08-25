package audiolevel

// SilenceThreshold is the peak amplitude (0..1) below which a recording is
// treated as silence and never sent to recognition.
const SilenceThreshold = 0.012

func Peak(pcm []byte) float64 {
	var peak int32
	for i := 0; i+1 < len(pcm); i += 2 {
		v := int32(int16(uint16(pcm[i]) | uint16(pcm[i+1])<<8))
		if v < 0 {
			v = -v
		}
		if v > peak {
			peak = v
		}
	}
	return float64(peak) / 32768
}

func IsSilent(pcm []byte) bool {
	return Peak(pcm) < SilenceThreshold
}

const (
	HeardFloor = 0.008
	HeardTop   = 0.35
)

func Heard(level float64) float64 {
	if level <= HeardFloor {
		return 0
	}
	if level > 1 {
		level = 1
	}
	span := DBFS(HeardTop) - DBFS(HeardFloor)
	part := (DBFS(level) - DBFS(HeardFloor)) / span
	if part < 0 {
		return 0
	}
	if part > 1 {
		return 1
	}
	return part
}
