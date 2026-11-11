package subwin

import (
	"strconv"
	"strings"

	"holdtotype/internal/ovplace"
)

const (
	EdgeLeft = 1 << iota
	EdgeRight
	EdgeTop
	EdgeBottom
)

func Hit(x, y, w, h, grip int32) int {
	if x < 0 || y < 0 || x >= w || y >= h {
		return 0
	}
	e := 0
	if x < grip {
		e |= EdgeLeft
	} else if x >= w-grip {
		e |= EdgeRight
	}
	if y < grip {
		e |= EdgeTop
	} else if y >= h-grip {
		e |= EdgeBottom
	}
	return e
}

type Limits struct {
	MinW    int32
	MaxW    int32
	Chrome  int32
	LineH   int32
	MinRows int32
	MaxRows int32
}

func NewLimits(minW, waW, waH, chrome, lineH, minRows int32) Limits {
	if lineH <= 0 {
		lineH = 1
	}
	maxW := waW
	if minW > maxW {
		minW = maxW
	}
	maxRows := (waH*4/5 - chrome) / lineH
	if maxRows < minRows {
		maxRows = minRows
	}
	return Limits{MinW: minW, MaxW: maxW, Chrome: chrome, LineH: lineH, MinRows: minRows, MaxRows: maxRows}
}

func (l Limits) Rows(rows int32) int32 { return clamp(rows, l.MinRows, l.MaxRows) }

func (l Limits) Width(w int32) int32 { return clamp(w, l.MinW, l.MaxW) }

func (l Limits) Height(rows int32) int32 { return l.Chrome + l.Rows(rows)*l.LineH }

func (l Limits) RowsFor(h int32) int32 {
	if l.LineH <= 0 {
		return l.MinRows
	}
	return l.Rows((h - l.Chrome + l.LineH/2) / l.LineH)
}

func resizeAxis(lo, hi int32, grabStart, grabEnd bool, delta int32, align int, size func(int32) int32) (int32, int32) {
	if !grabStart && !grabEnd {
		return lo, hi
	}
	out := delta
	if grabStart {
		out = -delta
	}
	cur := hi - lo
	want := cur + out
	if align == AlignCenter {
		want = cur + 2*out
	}
	n := size(want)
	switch {
	case align == AlignCenter:
		mid := lo + cur/2
		return mid - n/2, mid - n/2 + n
	case align == AlignStart:
		return lo, lo + n
	case align == AlignEnd:
		return hi - n, hi
	case grabStart:
		return hi - n, hi
	}
	return lo, lo + n
}

func Resize(start ovplace.Rect, edges int, dx, dy int32, mode string, l Limits) (ovplace.Rect, int32) {
	ha, va := Align(mode)
	r := start
	r.Left, r.Right = resizeAxis(start.Left, start.Right, edges&EdgeLeft != 0, edges&EdgeRight != 0, dx, ha, l.Width)
	rows := l.RowsFor(start.Bottom - start.Top)
	if edges&(EdgeTop|EdgeBottom) != 0 {
		var raw int32
		r.Top, r.Bottom = resizeAxis(start.Top, start.Bottom, edges&EdgeTop != 0, edges&EdgeBottom != 0, dy, va, func(h int32) int32 {
			raw = h
			return h
		})
		rows = l.RowsFor(raw)
		r.Top, r.Bottom = resizeAxis(start.Top, start.Bottom, edges&EdgeTop != 0, edges&EdgeBottom != 0, dy, va, func(int32) int32 { return l.Height(rows) })
	}
	return r, rows
}

type Box struct {
	W    int32 `json:"w"`
	Rows int32 `json:"rows"`
}

const (
	BoxMinW    = 360
	BoxMaxW    = 20000
	BoxMaxRows = 100
)

func validKey(k string) bool {
	a, b, ok := strings.Cut(k, "x")
	if !ok {
		return false
	}
	w, err1 := strconv.Atoi(a)
	h, err2 := strconv.Atoi(b)
	return err1 == nil && err2 == nil && w > 0 && h > 0 && strconv.Itoa(w) == a && strconv.Itoa(h) == b
}

func CleanBoxes(in map[string]Box) map[string]Box {
	out := map[string]Box{}
	for k, b := range in {
		if !validKey(k) || b.W < BoxMinW || b.W > BoxMaxW || b.Rows < 1 || b.Rows > BoxMaxRows {
			continue
		}
		out[k] = b
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
