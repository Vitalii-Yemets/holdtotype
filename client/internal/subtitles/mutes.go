package subtitles

import (
	"encoding/binary"
	"sync"
	"time"
)

const (
	MuteBefore = 150 * time.Millisecond
	MuteAfter  = 500 * time.Millisecond
	CueGuess   = 1500 * time.Millisecond
)

func WavDuration(wav []byte) (time.Duration, bool) {
	if len(wav) < 12 || string(wav[0:4]) != "RIFF" || string(wav[8:12]) != "WAVE" {
		return 0, false
	}
	var byteRate uint32
	for off := 12; off+8 <= len(wav); {
		id := string(wav[off : off+4])
		size := int(binary.LittleEndian.Uint32(wav[off+4 : off+8]))
		body := off + 8
		switch id {
		case "fmt ":
			if body+12 <= len(wav) {
				byteRate = binary.LittleEndian.Uint32(wav[body+8 : body+12])
			}
		case "data":
			if byteRate == 0 {
				return 0, false
			}
			if body+size > len(wav) || size < 0 {
				size = len(wav) - body
			}
			return time.Duration(int64(size) * int64(time.Second) / int64(byteRate)), true
		}
		if size < 0 {
			return 0, false
		}
		off = body + size + size%2
	}
	return 0, false
}

type span struct {
	from, to time.Time
}

type Mutes struct {
	mu     sync.Mutex
	before time.Duration
	after  time.Duration
	spans  []span
}

func NewMutes(before, after time.Duration) *Mutes {
	return &Mutes{before: before, after: after}
}

func (m *Mutes) Add(at time.Time, length time.Duration) {
	if length <= 0 {
		length = CueGuess
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	s := span{from: at.Add(-m.before), to: at.Add(length + m.after)}
	m.spans = append(m.spans, s)
	if len(m.spans) > 32 {
		m.spans = append([]span(nil), m.spans[len(m.spans)-32:]...)
	}
}

func (m *Mutes) Apply(pcm []byte, end time.Time, rate int) int {
	n := len(pcm) / 2
	if n == 0 || rate <= 0 {
		return 0
	}
	start := end.Add(-SamplesDuration(int64(n), rate))
	m.mu.Lock()
	kept := m.spans[:0]
	var hit []span
	for _, s := range m.spans {
		if !s.to.After(start.Add(-time.Minute)) {
			continue
		}
		kept = append(kept, s)
		if s.to.After(start) && s.from.Before(end) {
			hit = append(hit, s)
		}
	}
	m.spans = kept
	m.mu.Unlock()
	muted := 0
	for _, s := range hit {
		a := int64(s.from.Sub(start)) * int64(rate) / int64(time.Second)
		b := int64(s.to.Sub(start)) * int64(rate) / int64(time.Second)
		if a < 0 {
			a = 0
		}
		if b > int64(n) {
			b = int64(n)
		}
		for i := a; i < b; i++ {
			if pcm[i*2] != 0 || pcm[i*2+1] != 0 {
				pcm[i*2], pcm[i*2+1] = 0, 0
				muted++
			}
		}
	}
	return muted
}
