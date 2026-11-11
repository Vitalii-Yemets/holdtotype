package subwin

import "strings"

func Wrap(text string, width int32, measure func(string) int32) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}
	if width <= 0 {
		return []string{strings.Join(words, " ")}
	}
	space := measure(" ")
	var lines []string
	var cur []string
	var curW int32
	flush := func() {
		if len(cur) > 0 {
			lines = append(lines, strings.Join(cur, " "))
			cur, curW = nil, 0
		}
	}
	for _, w := range words {
		ww := measure(w)
		if ww > width {
			flush()
			r := []rune(w)
			start := 0
			for start < len(r) {
				end := start + 1
				for end < len(r) && measure(string(r[start:end+1])) <= width {
					end++
				}
				piece := string(r[start:end])
				if end == len(r) {
					cur, curW = []string{piece}, measure(piece)
				} else {
					lines = append(lines, piece)
				}
				start = end
			}
			continue
		}
		if len(cur) > 0 && curW+space+ww <= width {
			cur = append(cur, w)
			curW += space + ww
			continue
		}
		flush()
		cur, curW = []string{w}, ww
	}
	flush()
	return lines
}

type LineRef struct {
	Key  int
	Line int
}

type Scroll struct {
	back bool
	key  int
	line int
}

func (s *Scroll) Back() bool { return s.back }

func (s *Scroll) Follow() { *s = Scroll{} }

func lastStart(n, visible int) int {
	if n > visible {
		return n - visible
	}
	return 0
}

func (s *Scroll) Top(lines []LineRef, visible int) int {
	end := lastStart(len(lines), visible)
	if !s.back {
		return end
	}
	at := -1
	for i, l := range lines {
		if l.Key == s.key && l.Line <= s.line {
			at = i
		}
		if l.Key > s.key {
			if at < 0 {
				at = i
			}
			break
		}
	}
	if at < 0 {
		at = end
	}
	if at > end {
		at = end
	}
	return at
}

func (s *Scroll) Wheel(lines []LineRef, visible, delta int) {
	end := lastStart(len(lines), visible)
	top := s.Top(lines, visible) + delta
	if end == 0 || top >= end {
		s.Follow()
		return
	}
	if top < 0 {
		top = 0
	}
	s.back = true
	s.key = lines[top].Key
	s.line = lines[top].Line
}
