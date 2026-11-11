package subtitles

import (
	"strings"
	"time"
)

const (
	RevealRate     = 30
	RevealFastRate = 60
	RevealFastOver = 60
	RevealAllOver  = 240
)

type Revealer struct {
	shown  map[int]int
	failed map[int]bool
	credit float64
	last   time.Time
}

func NewRevealer() *Revealer {
	return &Revealer{shown: map[int]int{}, failed: map[int]bool{}}
}

func backToWord(r []rune, at int) int {
	if at >= len(r) {
		return len(r)
	}
	for at > 0 && r[at] != ' ' {
		at--
	}
	return at
}

func completeWords(r []rune) int {
	i := len(r)
	for i > 0 && r[i-1] != ' ' {
		i--
	}
	if i > 0 {
		i--
	}
	return i
}

func nextWordEnd(r []rune, from, limit int) int {
	i := from
	for i < limit && r[i] == ' ' {
		i++
	}
	for i < limit && r[i] != ' ' {
		i++
	}
	return i
}

func (v *Revealer) Rows(views []PieceView, hearing bool, now time.Time) ([]Row, bool) {
	dt := 0.0
	if !v.last.IsZero() && now.After(v.last) {
		dt = now.Sub(v.last).Seconds()
	}
	v.last = now
	texts := make([][]rune, len(views))
	avail := make([]int, len(views))
	seen := map[int]bool{}
	backlog := 0
	pending := hearing
	for i, p := range views {
		seen[p.Seq] = true
		r := []rune(p.Text)
		texts[i] = r
		if waitingView(p) {
			pending = true
			continue
		}
		if p.State == PieceFailed && !v.failed[p.Seq] {
			v.failed[p.Seq] = true
			v.shown[p.Seq] = backToWord(r, v.shown[p.Seq])
		}
		if v.shown[p.Seq] > len(r) {
			v.shown[p.Seq] = backToWord(r, len(r))
		}
		a := len(r)
		if p.State == PieceRunning {
			a = completeWords(r)
		}
		if a < v.shown[p.Seq] {
			a = v.shown[p.Seq]
		}
		avail[i] = a
		backlog += a - v.shown[p.Seq]
	}
	for seq := range v.shown {
		if !seen[seq] {
			delete(v.shown, seq)
			delete(v.failed, seq)
		}
	}
	switch {
	case backlog > RevealAllOver:
		for i, p := range views {
			if v.shown[p.Seq] < avail[i] {
				v.shown[p.Seq] = avail[i]
			}
		}
		v.credit = 0
	case backlog > 0:
		rate := float64(RevealRate)
		if backlog > RevealFastOver {
			rate = RevealFastRate
		}
		v.credit += rate * dt
	reveal:
		for i, p := range views {
			for v.shown[p.Seq] < avail[i] {
				end := nextWordEnd(texts[i], v.shown[p.Seq], avail[i])
				need := float64(end - v.shown[p.Seq])
				if need > v.credit {
					break reveal
				}
				v.credit -= need
				v.shown[p.Seq] = end
			}
		}
	default:
		v.credit = 0
	}
	hidden := false
	left := 0
	for i, p := range views {
		if waitingView(p) {
			continue
		}
		if v.shown[p.Seq] < len(texts[i]) {
			hidden = true
		}
		left += avail[i] - v.shown[p.Seq]
	}
	if left == 0 {
		v.credit = 0
	}
	rows := pieceRows(views, func(i int) string {
		if waitingView(views[i]) {
			return ""
		}
		return strings.TrimSpace(string(texts[i][:v.shown[views[i].Seq]]))
	}, pending || hidden)
	return rows, hidden
}
