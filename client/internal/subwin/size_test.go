package subwin

import (
	"reflect"
	"testing"

	"holdtotype/internal/ovplace"
)

var lim = Limits{MinW: 360, MaxW: 1920, Chrome: 46, LineH: 24, MinRows: 1, MaxRows: 32}

func TestHitFindsEdgesAndCorners(t *testing.T) {
	cases := []struct {
		x, y int32
		want int
	}{
		{0, 0, EdgeLeft | EdgeTop},
		{599, 100, EdgeRight},
		{300, 199, EdgeBottom},
		{5, 199, EdgeLeft | EdgeBottom},
		{599, 0, EdgeRight | EdgeTop},
		{300, 100, 0},
		{6, 6, 0},
		{593, 193, 0},
		{594, 194, EdgeRight | EdgeBottom},
		{-1, 5, 0},
		{600, 5, 0},
	}
	for _, c := range cases {
		if got := Hit(c.x, c.y, 600, 200, 6); got != c.want {
			t.Errorf("Hit(%d, %d) = %d, want %d", c.x, c.y, got, c.want)
		}
	}
}

func TestLimits(t *testing.T) {
	l := NewLimits(360, 1920, 1040, 46, 24, 1)
	if l != lim {
		t.Fatalf("limits %+v", l)
	}
	if got := NewLimits(360, 300, 100, 46, 24, 2); got.MinW != 300 || got.MaxW != 300 || got.MaxRows != 2 {
		t.Fatalf("a tiny screen %+v", got)
	}
	if lim.Height(3) != 118 || lim.Height(0) != 70 || lim.Height(99) != 46+32*24 {
		t.Fatalf("heights %d %d %d", lim.Height(3), lim.Height(0), lim.Height(99))
	}
	for h, want := range map[int32]int32{118: 3, 129: 3, 130: 4, 70: 1, 10: 1, 5000: 32} {
		if got := lim.RowsFor(h); got != want {
			t.Errorf("RowsFor(%d) = %d, want %d", h, got, want)
		}
	}
}

func rect(l, t, r, b int32) ovplace.Rect { return ovplace.Rect{Left: l, Top: t, Right: r, Bottom: b} }

func TestResizeKeepsTheAnchor(t *testing.T) {
	bottom := rect(660, 772, 1260, 890)
	cases := []struct {
		name   string
		start  ovplace.Rect
		edges  int
		dx, dy int32
		mode   string
		want   ovplace.Rect
		rows   int32
	}{
		{"bottom anchor, top edge up", bottom, EdgeTop, 0, -50, "bottom", rect(660, 724, 1260, 890), 5},
		{"bottom anchor, right edge out grows both sides", bottom, EdgeRight, 40, 0, "bottom", rect(620, 772, 1300, 890), 3},
		{"bottom anchor, as wide as the screen at most", bottom, EdgeLeft, -1000, 0, "bottom", rect(0, 772, 1920, 890), 3},
		{"bottom anchor, as narrow as the limit at least", bottom, EdgeRight, -500, 0, "bottom", rect(780, 772, 1140, 890), 3},
		{"bottom anchor, the bottom edge keeps its place", bottom, EdgeBottom, 0, 48, "bottom", rect(660, 724, 1260, 890), 5},
		{"custom, left edge moves alone", bottom, EdgeLeft, -100, 0, "custom", rect(560, 772, 1260, 890), 3},
		{"custom, bottom edge moves alone", bottom, EdgeBottom, 0, 48, "custom", rect(660, 772, 1260, 938), 5},
		{"top-right anchor, right edge keeps its place", rect(1296, 24, 1896, 142), EdgeRight, 50, 0, "top-right", rect(1246, 24, 1896, 142), 3},
		{"top-right anchor, bottom edge down", rect(1296, 24, 1896, 142), EdgeBottom, 0, 24, "top-right", rect(1296, 24, 1896, 166), 4},
		{"bottom-left anchor, top-right corner", rect(24, 772, 624, 890), EdgeTop | EdgeRight, 100, -24, "bottom-left", rect(24, 748, 724, 890), 4},
		{"centre anchor, bottom edge grows both ways", rect(660, 461, 1260, 579), EdgeBottom, 0, 24, "center", rect(660, 437, 1260, 603), 5},
		{"rows snap", bottom, EdgeTop, 0, -11, "bottom", rect(660, 772, 1260, 890), 3},
		{"rows snap up from half a row", bottom, EdgeTop, 0, -12, "bottom", rect(660, 748, 1260, 890), 4},
		{"no more rows than the limit", bottom, EdgeTop, 0, -2000, "bottom", rect(660, 890-46-32*24, 1260, 890), 32},
		{"at least one row", bottom, EdgeTop, 0, 500, "bottom", rect(660, 820, 1260, 890), 1},
		{"no edge changes nothing", bottom, 0, 300, 300, "bottom", bottom, 3},
	}
	for _, c := range cases {
		got, rows := Resize(c.start, c.edges, c.dx, c.dy, c.mode, lim)
		if got != c.want || rows != c.rows {
			t.Errorf("%s: %+v rows %d, want %+v rows %d", c.name, got, rows, c.want, c.rows)
		}
	}
}

func TestCleanBoxes(t *testing.T) {
	in := map[string]Box{
		"1920x1080": {W: 900, Rows: 4},
		"bad":       {W: 900, Rows: 4},
		"1920x0":    {W: 900, Rows: 4},
		"x1080":     {W: 900, Rows: 4},
		"2560x1440": {W: 100, Rows: 3},
		"3840x2160": {W: 1000, Rows: 0},
		"1280x720":  {W: 360, Rows: 100},
		"01x2":      {W: 900, Rows: 4},
	}
	want := map[string]Box{"1920x1080": {W: 900, Rows: 4}, "1280x720": {W: 360, Rows: 100}}
	if got := CleanBoxes(in); !reflect.DeepEqual(got, want) {
		t.Fatalf("%+v", got)
	}
	if CleanBoxes(map[string]Box{"bad": {}}) != nil || CleanBoxes(nil) != nil {
		t.Fatal("an empty result is not nil")
	}
}
