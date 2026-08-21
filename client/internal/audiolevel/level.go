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
