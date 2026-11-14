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

func Parts(translated, translationOnly bool) (bool, bool) {
	if !translated {
		return true, false
	}
	return !translationOnly, true
}

func Reuse(had, want []int) []int {
	out := make([]int, len(want))
	taken := make([]bool, len(had))
	for i, k := range want {
		out[i] = -1
		for j, h := range had {
			if h == k && h >= 0 && !taken[j] {
				out[i] = j
				taken[j] = true
				break
			}
		}
	}
	return out
}

func DefaultRows(lines, parts int32) (int32, int32) {
	if parts < 2 {
		return lines, 0
	}
	return lines * parts, lines
}
