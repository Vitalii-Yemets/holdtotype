package subtitles

import "time"

const (
	BacklogLimit = 60 * time.Second
	JoinLimit    = 25 * time.Second
)

type Backlog struct {
	rate  int
	max   int
	bytes int
	items []Phrase
}

func NewBacklog(rate int, limit time.Duration) *Backlog {
	if rate <= 0 {
		rate = 16000
	}
	return &Backlog{rate: rate, max: durationBytes(limit, rate)}
}

func durationBytes(d time.Duration, rate int) int {
	return int(int64(d)*int64(rate)/int64(time.Second)) * 2
}

func (b *Backlog) Push(p Phrase) []Phrase {
	b.items = append(b.items, p)
	b.bytes += len(p.PCM)
	var dropped []Phrase
	for b.bytes > b.max && len(b.items) > 1 {
		dropped = append(dropped, b.items[0])
		b.bytes -= len(b.items[0].PCM)
		b.items[0] = Phrase{}
		b.items = b.items[1:]
	}
	return dropped
}

func (b *Backlog) Take(limit time.Duration) (Phrase, int) {
	if len(b.items) == 0 {
		return Phrase{}, 0
	}
	room := durationBytes(limit, b.rate)
	n, size := 1, len(b.items[0].PCM)
	for n < len(b.items) && size+len(b.items[n].PCM) <= room {
		size += len(b.items[n].PCM)
		n++
	}
	first := b.items[0]
	out := Phrase{Start: first.Start, At: first.At, Cut: b.items[n-1].Cut}
	if n == 1 {
		out.PCM = first.PCM
		out.Speech = first.Speech
	} else {
		out.PCM = make([]byte, 0, size)
		for _, p := range b.items[:n] {
			out.PCM = append(out.PCM, p.PCM...)
			out.Speech += p.Speech
		}
	}
	for i := 0; i < n; i++ {
		b.items[i] = Phrase{}
	}
	b.items = b.items[n:]
	b.bytes -= size
	return out, n
}

func (b *Backlog) Len() int { return len(b.items) }

func (b *Backlog) Waiting() time.Duration {
	return SamplesDuration(int64(b.bytes/2), b.rate)
}

func (b *Backlog) Clear() {
	b.items = nil
	b.bytes = 0
}
