package subwin

import "holdtotype/internal/ovplace"

const (
	AlignStart = iota
	AlignCenter
	AlignEnd
	AlignFree
)

func Align(mode string) (h, v int) {
	switch mode {
	case ovplace.PosTopLeft:
		return AlignStart, AlignStart
	case ovplace.PosTop:
		return AlignCenter, AlignStart
	case ovplace.PosTopRight:
		return AlignEnd, AlignStart
	case ovplace.PosLeft:
		return AlignStart, AlignCenter
	case ovplace.PosCenter:
		return AlignCenter, AlignCenter
	case ovplace.PosRight:
		return AlignEnd, AlignCenter
	case ovplace.PosBottomLeft:
		return AlignStart, AlignEnd
	case ovplace.PosBottomRight:
		return AlignEnd, AlignEnd
	case ovplace.PosCustom:
		return AlignFree, AlignFree
	}
	return AlignCenter, AlignEnd
}

func ValidPosition(mode string) bool {
	return mode != ovplace.PosCaret && ovplace.Valid(mode)
}

func clamp(v, lo, hi int32) int32 {
	if hi < lo {
		return lo
	}
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func alignOn(a int, lo, hi, size, startGap, endGap int32) int32 {
	switch a {
	case AlignStart:
		return lo + startGap
	case AlignEnd:
		return hi - size - endGap
	}
	return lo + (hi-lo-size)/2
}

func Place(mode string, wa ovplace.Rect, w, h, side, bottom int32) (int32, int32) {
	ha, va := Align(mode)
	if ha == AlignFree {
		ha, va = AlignCenter, AlignEnd
	}
	endY := side
	if va == AlignEnd {
		endY = bottom
	}
	x := alignOn(ha, wa.Left, wa.Right, w, side, side)
	y := alignOn(va, wa.Top, wa.Bottom, h, side, endY)
	return clamp(x, wa.Left, wa.Right-w), clamp(y, wa.Top, wa.Bottom-h)
}

func Position(mode string, hasCustom bool) (string, bool) {
	if ValidPosition(mode) {
		return mode, false
	}
	if mode == "" && hasCustom {
		return ovplace.PosCustom, true
	}
	return ovplace.PosBottom, true
}
