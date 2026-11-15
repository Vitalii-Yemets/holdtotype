package subwin

const (
	LayoutRows    = "rows"
	LayoutColumns = "columns"

	ColWhole   = 1000
	ColDefault = ColWhole / 2
)

func Layout(v string) string {
	if v == LayoutColumns {
		return LayoutColumns
	}
	return LayoutRows
}

func InColumns(layout string, parts int) bool {
	return parts == 2 && Layout(layout) == LayoutColumns
}

func MinRows(parts int, columns bool) int32 {
	if parts > 1 && !columns {
		return 2
	}
	return 1
}

func ValidCol(c int32) bool { return c > 0 && c < ColWhole }

func colRange(w, pad, minCol int32) (int32, int32, bool) {
	lo := 2*pad + minCol
	hi := w - 2*pad - minCol
	return lo, hi, lo <= hi
}

func ColX(w, pad, minCol, col int32) int32 {
	if !ValidCol(col) {
		col = ColDefault
	}
	lo, hi, ok := colRange(w, pad, minCol)
	if !ok {
		return w / 2
	}
	return clamp(int32((int64(w)*int64(col)+ColWhole/2)/ColWhole), lo, hi)
}

func ColAt(x, w, pad, minCol int32) int32 {
	if w <= 0 {
		return ColDefault
	}
	lo, hi, ok := colRange(w, pad, minCol)
	if !ok {
		return ColDefault
	}
	x = clamp(x, lo, hi)
	return clamp(int32((int64(x)*ColWhole+int64(w)/2)/int64(w)), 1, ColWhole-1)
}

type Area struct {
	Left, Top, Right, Bottom int32
}

func ColAreas(w, top, bottom, pad, minCol, col int32) (Area, Area, int32) {
	x := ColX(w, pad, minCol, col)
	return Area{Left: pad, Top: top, Right: x - pad, Bottom: bottom}, Area{Left: x + pad, Top: top, Right: w - pad, Bottom: bottom}, x
}

func HitColumn(x, y, lineX, top, bottom, grip int32) bool {
	return y >= top && y < bottom && x >= lineX-grip && x <= lineX+grip
}
