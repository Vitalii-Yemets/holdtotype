package subtitles

import (
	"bytes"
	"encoding/binary"
	"strings"
	"testing"
)

func wavBytes(format, channels uint16, rate uint32, bitsPer uint16, extra []byte, pcm []byte) []byte {
	var b bytes.Buffer
	b.WriteString("RIFF")
	binary.Write(&b, binary.LittleEndian, uint32(0))
	b.WriteString("WAVE")
	b.WriteString("fmt ")
	binary.Write(&b, binary.LittleEndian, uint32(16))
	binary.Write(&b, binary.LittleEndian, format)
	binary.Write(&b, binary.LittleEndian, channels)
	binary.Write(&b, binary.LittleEndian, rate)
	binary.Write(&b, binary.LittleEndian, rate*uint32(channels)*uint32(bitsPer/8))
	binary.Write(&b, binary.LittleEndian, channels*bitsPer/8)
	binary.Write(&b, binary.LittleEndian, bitsPer)
	b.Write(extra)
	b.WriteString("data")
	binary.Write(&b, binary.LittleEndian, uint32(len(pcm)))
	b.Write(pcm)
	return b.Bytes()
}

func TestParseWavTakesTheSound(t *testing.T) {
	pcm := []byte{1, 2, 3, 4, 5, 6}
	got, err := ParseWav(wavBytes(1, 1, 16000, 16, nil, pcm), 16000)
	if err != nil || !bytes.Equal(got, pcm) {
		t.Fatalf("%v %v", got, err)
	}
	list := []byte("LIST")
	list = binary.LittleEndian.AppendUint32(list, 3)
	list = append(list, 'a', 'b', 'c', 0)
	got, err = ParseWav(wavBytes(1, 1, 16000, 16, list, pcm), 16000)
	if err != nil || !bytes.Equal(got, pcm) {
		t.Fatalf("with a list chunk: %v %v", got, err)
	}
	raw := wavBytes(1, 1, 16000, 16, nil, pcm)
	binary.LittleEndian.PutUint32(raw[40:], 0xffffffff)
	got, err = ParseWav(raw, 16000)
	if err != nil || !bytes.Equal(got, pcm) {
		t.Fatalf("a streamed data size: %v %v", got, err)
	}
}

func TestParseWavRefusesOtherSound(t *testing.T) {
	pcm := make([]byte, 8)
	cases := map[string][]byte{
		"48000 Hz":   wavBytes(1, 1, 48000, 16, nil, pcm),
		"2 channel":  wavBytes(1, 2, 16000, 16, nil, pcm),
		"8 bit":      wavBytes(1, 1, 16000, 8, nil, pcm),
		"format 3":   wavBytes(3, 1, 16000, 32, nil, pcm),
		"not a wav":  []byte("OggS and more bytes"),
		"has no":     []byte("RIFF\x04\x00\x00\x00WAVE"),
		"cut short":  append([]byte("RIFF\x00\x00\x00\x00WAVEfmt "), 0x40, 0, 0, 0, 1),
		"before its": append([]byte("RIFF\x00\x00\x00\x00WAVEdata"), 2, 0, 0, 0, 1, 2),
	}
	for want, data := range cases {
		_, err := ParseWav(data, 16000)
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: %v", want, err)
		}
	}
}
