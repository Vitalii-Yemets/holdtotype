package sentence

type Numbers struct {
	last  int
	draft int
}

func NewNumbers(last int) *Numbers { return &Numbers{last: last} }

func (n *Numbers) Last() int { return n.last }

func (n *Numbers) Reserve() (int, bool) {
	if n.draft != 0 {
		return n.draft, false
	}
	n.last++
	n.draft = n.last
	return n.draft, true
}

func (n *Numbers) Take() (int, bool) {
	if n.draft != 0 {
		seq := n.draft
		n.draft = 0
		return seq, true
	}
	n.last++
	return n.last, false
}
