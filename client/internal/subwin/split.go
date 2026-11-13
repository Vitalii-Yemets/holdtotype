package subwin

const DefaultUpperRows = 2

func SplitRows(total, upper int32) (int32, int32) {
	if total < 2 {
		return 0, total
	}
	if upper <= 0 {
		upper = DefaultUpperRows
	}
	upper = clamp(upper, 1, total-1)
	return upper, total - upper
}

func SplitAt(y, areaTop, lineH, total int32) int32 {
	if lineH <= 0 || total < 2 {
		return 0
	}
	rows := (y - areaTop + lineH/2) / lineH
	if y < areaTop {
		rows = 0
	}
	return clamp(rows, 1, total-1)
}

func HitSplit(y, lineY, grip int32) bool {
	return y >= lineY-grip && y <= lineY+grip
}

func DefaultRows(lines, parts int32) (int32, int32) {
	if parts < 2 {
		return lines, 0
	}
	return lines * parts, lines
}
