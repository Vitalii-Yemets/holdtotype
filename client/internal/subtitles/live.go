package subtitles

import (
	"strings"
	"time"
)

type Finished struct {
	Start time.Time
	Text  string
}

type Live struct {
	cur   int
	done  int
	text  string
	start time.Time
}

func NewLive() *Live { return &Live{cur: -1, done: -1} }

func (l *Live) Update(seg int, text string, final bool, now time.Time) []Finished {
	var out []Finished
	if seg <= l.done {
		return nil
	}
	if l.cur >= 0 && seg > l.cur {
		out = l.close(out)
	}
	if l.cur < 0 {
		l.cur = seg
	}
	if t := strings.Join(strings.Fields(text), " "); t != "" {
		if l.text == "" {
			l.start = now
		}
		l.text = t
	}
	if final {
		out = l.close(out)
	}
	return out
}

func (l *Live) close(out []Finished) []Finished {
	if l.text != "" {
		out = append(out, Finished{Start: l.start, Text: l.text})
		if l.cur > l.done {
			l.done = l.cur
		}
	}
	l.cur = -1
	l.text = ""
	l.start = time.Time{}
	return out
}

func (l *Live) Current() (string, time.Time) { return l.text, l.start }

func (l *Live) Flush() []Finished {
	var out []Finished
	if l.text != "" {
		out = append(out, Finished{Start: l.start, Text: l.text})
	}
	l.cur = -1
	l.done = -1
	l.text = ""
	l.start = time.Time{}
	return out
}
