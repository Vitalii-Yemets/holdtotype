package sentence

import "time"

const (
	ByModel   = "model"
	BySilence = "silence"
	ByLimit   = "limit"
	ByEnd     = "end"
)

type Params struct {
	Ahead   int
	Silence time.Duration
	Limit   int
	Context int
}

func Defaults() Params {
	return Params{Ahead: 3, Silence: 2 * time.Second, Limit: 40, Context: 30}
}

type Judge func(words []string) ([]bool, []float64, error)

type Close struct {
	First int
	End   int
	At    time.Time
	How   string
}

func (c Close) Words() int { return c.End - c.First + 1 }

type Stream struct {
	p      Params
	judge  Judge
	words  []string
	times  []time.Time
	start  int
	lastLo int
	lastJ  int
	end    []bool
	prob   []float64
}

func New(p Params, judge Judge) *Stream {
	if p.Ahead < 0 {
		p.Ahead = 0
	}
	if p.Context < 0 {
		p.Context = 0
	}
	return &Stream{p: p, judge: judge, lastLo: -1, lastJ: -1}
}

func (s *Stream) Open() int { return len(s.words) - s.start }

func (s *Stream) Quiet(now time.Time) []Close {
	j := len(s.words) - 1
	if s.p.Silence <= 0 || j < s.start || now.Sub(s.times[j]) < s.p.Silence {
		return nil
	}
	return []Close{s.close(j, s.times[j].Add(s.p.Silence), BySilence)}
}

func (s *Stream) Add(word string, at time.Time) ([]Close, error) {
	out := s.Quiet(at)
	s.words = append(s.words, word)
	s.times = append(s.times, at)
	j := len(s.words) - 1
	for {
		lo := s.start - s.p.Context
		if lo < 0 {
			lo = 0
		}
		if lo != s.lastLo || j != s.lastJ {
			end, prob, err := s.judge(s.words[lo : j+1])
			if err != nil {
				return out, err
			}
			s.end, s.prob, s.lastLo, s.lastJ = end, prob, lo, j
		}
		hit := -1
		for i := s.start; i <= j-s.p.Ahead; i++ {
			if i-lo < len(s.end) && s.end[i-lo] {
				hit = i
				break
			}
		}
		if hit >= 0 {
			out = append(out, s.close(hit, at, ByModel))
			continue
		}
		if s.p.Limit > 0 && j-s.start+1 > s.p.Limit {
			out = append(out, s.close(s.likeliest(lo, j), at, ByLimit))
			continue
		}
		return out, nil
	}
}

func (s *Stream) likeliest(lo, j int) int {
	best := -1
	for i := s.start + 3; i <= j-3; i++ {
		if i-lo >= len(s.prob) {
			break
		}
		if best < 0 || s.prob[i-lo] > s.prob[best-lo] {
			best = i
		}
	}
	if best < 0 {
		return j
	}
	return best
}

func (s *Stream) Flush() []Close {
	j := len(s.words) - 1
	if j < s.start {
		return nil
	}
	return []Close{s.close(j, s.times[j].Add(s.p.Silence), ByEnd)}
}

func (s *Stream) close(end int, at time.Time, how string) Close {
	c := Close{First: s.start, End: end, At: at, How: how}
	s.start = end + 1
	return c
}
