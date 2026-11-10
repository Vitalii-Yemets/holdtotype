package subtitles

type Backlog struct {
	max   int
	items []Phrase
}

func NewBacklog(max int) *Backlog {
	if max < 1 {
		max = 1
	}
	return &Backlog{max: max}
}

func (b *Backlog) Push(p Phrase) []Phrase {
	b.items = append(b.items, p)
	over := len(b.items) - b.max
	if over <= 0 {
		return nil
	}
	dropped := append([]Phrase(nil), b.items[:over]...)
	n := copy(b.items, b.items[over:])
	for i := n; i < len(b.items); i++ {
		b.items[i] = Phrase{}
	}
	b.items = b.items[:n]
	return dropped
}

func (b *Backlog) PushFront(p Phrase) []Phrase {
	b.items = append([]Phrase{p}, b.items...)
	over := len(b.items) - b.max
	if over <= 0 {
		return nil
	}
	dropped := append([]Phrase(nil), b.items[1:1+over]...)
	b.items = append(b.items[:1], b.items[1+over:]...)
	return dropped
}

func (b *Backlog) Pop() (Phrase, bool) {
	if len(b.items) == 0 {
		return Phrase{}, false
	}
	p := b.items[0]
	b.items[0] = Phrase{}
	b.items = b.items[1:]
	return p, true
}

func (b *Backlog) Len() int { return len(b.items) }

func (b *Backlog) Clear() { b.items = nil }
