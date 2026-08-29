package sherpaproto

import (
	"encoding/binary"
	"errors"
	"math"
)

const ChunkBytes = 10240

var ErrShortWAV = errors.New("the file is shorter than the WAV header")

func Header(sampleRate, payloadBytes int) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint32(b[0:4], uint32(sampleRate))
	binary.LittleEndian.PutUint32(b[4:8], uint32(payloadBytes))
	return b
}

func FromPCM16(pcm []byte) []byte {
	n := len(pcm) / 2
	out := make([]byte, n*4)
	for i := 0; i < n; i++ {
		s := int16(binary.LittleEndian.Uint16(pcm[i*2:]))
		binary.LittleEndian.PutUint32(out[i*4:], math.Float32bits(float32(s)/32768))
	}
	return out
}

func PCMFromWAV(wav []byte) ([]byte, int, error) {
	if len(wav) < 12 || string(wav[0:4]) != "RIFF" || string(wav[8:12]) != "WAVE" {
		return nil, 0, ErrShortWAV
	}
	rate := 0
	pos := 12
	for pos+8 <= len(wav) {
		id := string(wav[pos : pos+4])
		size := int(binary.LittleEndian.Uint32(wav[pos+4 : pos+8]))
		body := pos + 8
		if size < 0 || body+size > len(wav) {
			size = len(wav) - body
		}
		switch id {
		case "fmt ":
			if size >= 8 {
				rate = int(binary.LittleEndian.Uint32(wav[body+4 : body+8]))
			}
		case "data":
			return wav[body : body+size], rate, nil
		}
		pos = body + size
		if size%2 == 1 {
			pos++
		}
	}
	return nil, rate, ErrShortWAV
}

func Chunks(payload []byte, size int) [][]byte {
	if size <= 0 {
		size = ChunkBytes
	}
	var out [][]byte
	for i := 0; i < len(payload); i += size {
		end := i + size
		if end > len(payload) {
			end = len(payload)
		}
		out = append(out, payload[i:end])
	}
	return out
}
