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

const (
	RowLine = iota
	RowLive
	RowPiece
	RowHeard
)

type Row struct {
	Kind int
	Text string
}

const (
	PieceWaiting = iota
	PieceRunning
	PieceDone
	PieceFailed
)

type piece struct {
	seq    int
	para   bool
	source string
	text   string
	state  int
}

type Board struct {
	mu      sync.Mutex
	keep    int
	next    int
	lines   []Line
	live    string
	pieces  []piece
	hearing bool
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

func (v View) Text() string {
	parts := make([]string, 0, len(v.Lines)+1)
	for _, ln := range v.Lines {
		if ln.Text != "" {
			parts = append(parts, ln.Text)
		}
	}
	if v.Live != "" {
		parts = append(parts, v.Live)
	}
	return strings.Join(parts, "\n")
}

func (v View) Rows() []Row {
	rows := make([]Row, 0, len(v.Lines)+1)
	for _, ln := range v.Lines {
		if ln.Text != "" {
			rows = append(rows, Row{Kind: RowLine, Text: ln.Text})
		}
	}
	if v.Live != "" {
		rows = append(rows, Row{Kind: RowLive, Text: v.Live})
	}
	return rows
}

func (b *Board) findPiece(seq int) int {
	for i := len(b.pieces) - 1; i >= 0; i-- {
		if b.pieces[i].seq == seq {
			return i
		}
	}
	return -1
}

func (b *Board) AddPiece(seq int, para bool, source string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.pieces = append(b.pieces, piece{seq: seq, para: para, source: oneLine(source)})
	if over := len(b.pieces) - b.keep; over > 0 {
		n := copy(b.pieces, b.pieces[over:])
		b.pieces = b.pieces[:n]
	}
}

func (b *Board) JoinPieces(seqs []int) {
	if len(seqs) < 2 {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	at := b.findPiece(seqs[0])
	if at < 0 {
		return
	}
	for _, s := range seqs[1:] {
		i := b.findPiece(s)
		if i <= at {
			continue
		}
		b.pieces[at].source = oneLine(b.pieces[at].source + " " + b.pieces[i].source)
		b.pieces = append(b.pieces[:i], b.pieces[i+1:]...)
	}
}

func (b *Board) SetPiece(seq int, text string, state int) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	i := b.findPiece(seq)
	if i < 0 {
		return false
	}
	b.pieces[i].text = oneLine(text)
	b.pieces[i].state = state
	return true
}

func (b *Board) SetHearing(on bool) {
	b.mu.Lock()
	b.hearing = on
	b.mu.Unlock()
}

func (b *Board) TranslationRows(maxRunes int) []Row {
	b.mu.Lock()
	defer b.mu.Unlock()
	if maxRunes <= 0 {
		return nil
	}
	var rows []Row
	pending := b.hearing
	breakNext := true
	for _, pc := range b.pieces {
		if pc.para {
			breakNext = true
		}
		switch {
		case pc.state == PieceWaiting, pc.state == PieceRunning && pc.text == "":
			pending = true
			continue
		case pc.state == PieceFailed:
			if pc.source != "" {
				rows = append(rows, Row{Kind: RowHeard, Text: pc.source})
			}
			breakNext = true
			continue
		case pc.text == "":
			continue
		}
		if breakNext || len(rows) == 0 || rows[len(rows)-1].Kind != RowPiece {
			rows = append(rows, Row{Kind: RowPiece, Text: pc.text})
		} else {
			rows[len(rows)-1].Text += " " + pc.text
		}
		breakNext = false
	}
	if pending {
		if len(rows) == 0 {
			rows = append(rows, Row{Kind: RowPiece, Text: "…"})
		} else {
			rows[len(rows)-1].Text += " …"
		}
	}
	budget := maxRunes
	start := len(rows)
	for start > 0 && budget > 0 {
		n := len([]rune(rows[start-1].Text))
		if n > budget {
			rows[start-1].Text = livetail.Tail(rows[start-1].Text, budget)
			start--
			break
		}
		budget -= n
		start--
	}
	return rows[start:]
}
