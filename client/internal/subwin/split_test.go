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
