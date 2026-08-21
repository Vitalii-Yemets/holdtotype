package main

import (
	"bytes"
	"encoding/binary"
)

func wavFromPCM16(pcm []byte, rate int) []byte {
	var b bytes.Buffer
	le := binary.LittleEndian
	w32 := func(v uint32) { _ = binary.Write(&b, le, v) }
	w16 := func(v uint16) { _ = binary.Write(&b, le, v) }

	b.WriteString("RIFF")
	w32(uint32(36 + len(pcm)))
	b.WriteString("WAVE")
	b.WriteString("fmt ")
	w32(16)
	w16(1)
	w16(1)
	w32(uint32(rate))
	w32(uint32(rate * 2))
	w16(2)
	w16(16)
	b.WriteString("data")
	w32(uint32(len(pcm)))
	b.Write(pcm)
	return b.Bytes()
}
