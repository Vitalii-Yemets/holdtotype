package subwin

import "testing"

func TestSplitRows(t *testing.T) {
	cases := []struct{ total, upper, wantU, wantL int32 }{
		{5, 0, 2, 3},
		{3, 0, 2, 1},
		{2, 0, 1, 1},
		{8, 6, 6, 2},
		{8, 8, 7, 1},
		{8, -3, 2, 6},
		{4, 1, 1, 3},
		{1, 0, 0, 1},
	}
	for _, c := range cases {
		u, l := SplitRows(c.total, c.upper)
		if u != c.wantU || l != c.wantL {
			t.Errorf("SplitRows(%d, %d) = %d, %d, want %d, %d", c.total, c.upper, u, l, c.wantU, c.wantL)
		}
	}
}

func TestSplitAtSnapsToRows(t *testing.T) {
	cases := []struct{ y, want int32 }{
		{100 + 24*2, 2},
		{100 + 24*2 + 11, 2},
		{100 + 24*2 + 12, 3},
		{100 + 24*2 - 12, 2},
		{100 + 24*2 - 13, 1},
		{100, 1},
		{50, 1},
		{100 + 24*10, 5},
	}
	for _, c := range cases {
		if got := SplitAt(c.y, 100, 24, 6); got != c.want {
			t.Errorf("SplitAt(%d) = %d, want %d", c.y, got, c.want)
		}
	}
	if SplitAt(150, 100, 24, 1) != 0 || SplitAt(150, 100, 0, 6) != 0 {
		t.Fatal("a split without room")
	}
}

func TestHitSplit(t *testing.T) {
	for y, want := range map[int32]bool{94: true, 100: true, 106: true, 93: false, 107: false} {
		if HitSplit(y, 100, 6) != want {
			t.Errorf("HitSplit(%d) = %v", y, !want)
		}
	}
}

func TestLinesPerPart(t *testing.T) {
	for lines := int32(2); lines <= 8; lines++ {
		rows, top := DefaultRows(lines, 1)
		if rows != lines || top != 0 {
			t.Errorf("one part of %d lines: %d rows, %d above the line", lines, rows, top)
		}
		rows, top = DefaultRows(lines, 2)
		if rows != 2*lines || top != lines {
			t.Errorf("two parts of %d lines: %d rows, %d above the line", lines, rows, top)
		}
		upper, lower := SplitRows(rows, top)
		if upper != lines || lower != lines {
			t.Errorf("two parts of %d lines are split %d and %d", lines, upper, lower)
		}
	}
}

func TestThePartsOfTheWindow(t *testing.T) {
	cases := []struct {
		translated, only, heard, translation bool
	}{
		{false, false, true, false},
		{false, true, true, false},
		{true, false, true, true},
		{true, true, false, true},
	}
	for _, c := range cases {
		h, tr := Parts(c.translated, c.only)
		if h != c.heard || tr != c.translation {
			t.Errorf("translated %v, only %v: heard %v, translation %v", c.translated, c.only, h, tr)
		}
		if !h && !tr {
			t.Errorf("translated %v, only %v: the window has no part", c.translated, c.only)
		}
	}
	if rows, top := DefaultRows(3, 1); rows != 3 || top != 0 {
		t.Errorf("the translation alone: %d rows, %d above the line", rows, top)
	}
}

func TestThePartsThatStayAreKept(t *testing.T) {
	cases := []struct {
		had, want, got []int
	}{
		{[]int{0, 1}, []int{1}, []int{1}},
		{[]int{1}, []int{0, 1}, []int{-1, 0}},
		{[]int{0}, []int{0, 1}, []int{0, -1}},
		{[]int{0, 1}, []int{0}, []int{0}},
		{[]int{0, 1}, []int{0, 1}, []int{0, 1}},
		{nil, []int{0, 1}, []int{-1, -1}},
		{[]int{-1, -1}, []int{0, 1}, []int{-1, -1}},
		{[]int{-1, 1}, []int{0, 1}, []int{-1, 1}},
		{[]int{1, 1}, []int{1}, []int{0}},
	}
	for _, c := range cases {
		got := Reuse(c.had, c.want)
		if len(got) != len(c.got) {
			t.Fatalf("%v → %v: %v", c.had, c.want, got)
		}
		for i := range got {
			if got[i] != c.got[i] {
				t.Errorf("%v → %v: %v, want %v", c.had, c.want, got, c.got)
				break
			}
		}
	}
}
