package subtitles

import (
	"encoding/binary"
	"errors"
	"fmt"
)

func ParseWav(data []byte, rate int) ([]byte, error) {
	if len(data) < 12 || string(data[0:4]) != "RIFF" || string(data[8:12]) != "WAVE" {
		return nil, errors.New("not a wav file")
	}
	var format, channels, bitsPer uint16
	var sampleRate uint32
	haveFmt := false
	off := 12
	for off+8 <= len(data) {
		id := string(data[off : off+4])
		size := int(binary.LittleEndian.Uint32(data[off+4:]))
		body := off + 8
		if size < 0 || body+size > len(data) {
			if id == "data" && haveFmt {
				size = len(data) - body
			} else {
				return nil, fmt.Errorf("the %q chunk is cut short", id)
			}
		}
		switch id {
		case "fmt ":
			if size < 16 {
				return nil, errors.New("the format chunk is too short")
			}
			format = binary.LittleEndian.Uint16(data[body:])
			channels = binary.LittleEndian.Uint16(data[body+2:])
			sampleRate = binary.LittleEndian.Uint32(data[body+4:])
			bitsPer = binary.LittleEndian.Uint16(data[body+14:])
			if format == 0xfffe && size >= 40 {
				format = binary.LittleEndian.Uint16(data[body+24:])
			}
			haveFmt = true
		case "data":
			if !haveFmt {
				return nil, errors.New("the sound comes before its format")
			}
			if format != 1 || channels != 1 || bitsPer != 16 || int(sampleRate) != rate {
				return nil, fmt.Errorf("the sound is %d Hz, %d channel(s), %d bit, format %d; it must be %d Hz, mono, 16 bit PCM", sampleRate, channels, bitsPer, format, rate)
			}
			pcm := data[body : body+size]
			return pcm[:len(pcm)/2*2], nil
		}
		off = body + size + size%2
	}
	return nil, errors.New("the file has no sound")
}
