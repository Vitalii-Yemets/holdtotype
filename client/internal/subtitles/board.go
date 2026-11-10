package subtitles

import (
	"strings"
	"sync"
	"time"

	"holdtotype/internal/livetail"
)

type Line struct {
	ID    int
	Start time.Time
	Text  string
}

type View struct {
	Lines []Line
	Live  string
}

type Board struct {
	mu    sync.Mutex
	keep  int
	next  int
	lines []Line
	live  string
}

func NewBoard(keep int) *Board {
	if keep < 1 {
		keep = 1
	}
	return &Board{keep: keep, next: 1}
}

func (b *Board) Add(start time.Time, text string) Line {
	b.mu.Lock()
	defer b.mu.Unlock()
	ln := Line{ID: b.next, Start: start, Text: text}
	b.next++
	b.lines = append(b.lines, ln)
	if over := len(b.lines) - b.keep; over > 0 {
		n := copy(b.lines, b.lines[over:])
		b.lines = b.lines[:n]
	}
	return ln
}

func (b *Board) SetLive(text string) {
	b.mu.Lock()
	b.live = strings.Join(strings.Fields(text), " ")
	b.mu.Unlock()
}

func (b *Board) Live() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.live
}

func (b *Board) Count() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.next - 1
}

func (b *Board) Lines() []Line {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]Line(nil), b.lines...)
}

func (b *Board) Tail(maxRunes int) View {
	b.mu.Lock()
	defer b.mu.Unlock()
	var v View
	if maxRunes <= 0 {
		return v
	}
	budget := maxRunes
	if b.live != "" {
		v.Live = livetail.Tail(b.live, maxRunes)
		budget -= len([]rune(v.Live))
	}
	var picked []Line
	for i := len(b.lines) - 1; i >= 0 && budget > 0; i-- {
		ln := b.lines[i]
		n := len([]rune(ln.Text))
		if n > budget {
			if len(picked) == 0 {
				ln.Text = livetail.Tail(ln.Text, budget)
				picked = append(picked, ln)
			}
			break
		}
		picked = append(picked, ln)
		budget -= n
	}
	for i := len(picked) - 1; i >= 0; i-- {
		v.Lines = append(v.Lines, picked[i])
	}
	return v
}
