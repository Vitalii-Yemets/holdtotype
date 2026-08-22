package sherpaproto

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"
)

func wav(pcm []byte, rate int) []byte {
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

func TestHeader(t *testing.T) {
	h := Header(16000, 4096)
	if len(h) != 8 {
		t.Fatalf("длина заголовка %d, ожидалось 8", len(h))
	}
	if got := binary.LittleEndian.Uint32(h[0:4]); got != 16000 {
		t.Errorf("частота %d, ожидалось 16000", got)
	}
	if got := binary.LittleEndian.Uint32(h[4:8]); got != 4096 {
		t.Errorf("длина %d, ожидалось 4096", got)
	}
}

func TestFromPCM16(t *testing.T) {
	pcm := make([]byte, 8)
	binary.LittleEndian.PutUint16(pcm[0:], uint16(0))
	binary.LittleEndian.PutUint16(pcm[2:], uint16(int16(32767)))
	binary.LittleEndian.PutUint16(pcm[4:], uint16(0x8000))
	binary.LittleEndian.PutUint16(pcm[6:], uint16(int16(16384)))

	out := FromPCM16(pcm)
	if len(out) != 16 {
		t.Fatalf("длина %d, ожидалось 16", len(out))
	}
	want := []float32{0, 32767.0 / 32768, -1, 0.5}
	for i, w := range want {
		got := math.Float32frombits(binary.LittleEndian.Uint32(out[i*4:]))
		if math.Abs(float64(got-w)) > 1e-6 {
			t.Errorf("сэмпл %d: %v, ожидалось %v", i, got, w)
		}
	}
}

func TestFromPCM16OddTail(t *testing.T) {
	if out := FromPCM16([]byte{1, 2, 3}); len(out) != 4 {
		t.Fatalf("нечётный хвост должен отбрасываться, получено %d байт", len(out))
	}
	if out := FromPCM16(nil); len(out) != 0 {
		t.Fatal("пустой вход должен давать пустой выход")
	}
}

func TestPCMFromWAV(t *testing.T) {
	pcm := []byte{1, 0, 2, 0, 3, 0, 4, 0}
	got, rate, err := PCMFromWAV(wav(pcm, 16000))
	if err != nil {
		t.Fatalf("ошибка разбора: %v", err)
	}
	if rate != 16000 {
		t.Errorf("частота %d, ожидалось 16000", rate)
	}
	if !bytes.Equal(got, pcm) {
		t.Errorf("данные %v, ожидалось %v", got, pcm)
	}
}

func TestPCMFromWAVRejectsGarbage(t *testing.T) {
	if _, _, err := PCMFromWAV([]byte("not a wav at all")); err == nil {
		t.Fatal("мусор должен отвергаться")
	}
	if _, _, err := PCMFromWAV(nil); err == nil {
		t.Fatal("пустой вход должен отвергаться")
	}
}

func TestPCMFromWAVSkipsExtraChunks(t *testing.T) {
	pcm := []byte{7, 0, 8, 0}
	base := wav(pcm, 8000)
	var b bytes.Buffer
	b.Write(base[:12])
	b.WriteString("LIST")
	_ = binary.Write(&b, binary.LittleEndian, uint32(5))
	b.WriteString("junk\x00")
	b.WriteByte(0)
	b.Write(base[12:])
	got, rate, err := PCMFromWAV(b.Bytes())
	if err != nil {
		t.Fatalf("ошибка разбора: %v", err)
	}
	if rate != 8000 {
		t.Errorf("частота %d, ожидалось 8000", rate)
	}
	if !bytes.Equal(got, pcm) {
		t.Errorf("данные %v, ожидалось %v", got, pcm)
	}
}

func TestChunks(t *testing.T) {
	payload := make([]byte, 25)
	c := Chunks(payload, 10)
	if len(c) != 3 {
		t.Fatalf("кусков %d, ожидалось 3", len(c))
	}
	if len(c[0]) != 10 || len(c[2]) != 5 {
		t.Fatalf("размеры кусков %d и %d", len(c[0]), len(c[2]))
	}
	total := 0
	for _, part := range c {
		total += len(part)
	}
	if total != len(payload) {
		t.Fatalf("суммарно %d байт из %d", total, len(payload))
	}
	if Chunks(nil, 10) != nil {
		t.Fatal("пустой вход должен давать пустой список")
	}
}
