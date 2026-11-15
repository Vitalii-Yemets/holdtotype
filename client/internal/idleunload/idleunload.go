package idleunload

import "time"

var Choices = []int{0, 5, 10, 30, 60}

func Minutes(v int) int {
	for _, c := range Choices {
		if v == c {
			return v
		}
	}
	return 0
}

type State struct {
	Busy   bool
	Loaded bool
	Last   time.Time
}

func Due(minutes int, s State, now time.Time) bool {
	m := Minutes(minutes)
	if m == 0 || s.Busy || !s.Loaded || s.Last.IsZero() {
		return false
	}
	return now.Sub(s.Last) >= time.Duration(m)*time.Minute
}
