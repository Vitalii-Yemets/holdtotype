package subwin

import (
	"testing"

	"holdtotype/internal/ovplace"
)

var wa = ovplace.Rect{Left: 0, Top: 0, Right: 1920, Bottom: 1040}

func TestPlaceAnchors(t *testing.T) {
	cases := map[string][2]int32{
		"top-left":     {24, 24},
		"top":          {660, 24},
		"top-right":    {1296, 24},
		"left":         {24, 470},
		"center":       {660, 470},
		"right":        {1296, 470},
		"bottom-left":  {24, 790},
		"bottom":       {660, 790},
		"bottom-right": {1296, 790},
		"custom":       {660, 790},
		"":             {660, 790},
	}
	for mode, want := range cases {
		x, y := Place(mode, wa, 600, 100, 24, 150)
		if x != want[0] || y != want[1] {
			t.Errorf("%q: %d,%d, want %d,%d", mode, x, y, want[0], want[1])
		}
	}
}

func TestPlaceOnAnotherMonitorAndTooBig(t *testing.T) {
	second := ovplace.Rect{Left: 1920, Top: -200, Right: 4480, Bottom: 1240}
	if x, y := Place("top-right", second, 600, 100, 24, 150); x != 4480-624 || y != -176 {
		t.Fatalf("second monitor %d,%d", x, y)
	}
	if x, y := Place("bottom", wa, 1900, 1000, 24, 150); x != 10 || y != 0 {
		t.Fatalf("a window that barely fits %d,%d", x, y)
	}
}

func TestValidPosition(t *testing.T) {
	for _, m := range []string{"top-left", "top", "top-right", "left", "center", "right", "bottom-left", "bottom", "bottom-right", "custom"} {
		if !ValidPosition(m) {
			t.Errorf("%q refused", m)
		}
	}
	for _, m := range []string{"caret", "", "middle", "Bottom"} {
		if ValidPosition(m) {
			t.Errorf("%q accepted", m)
		}
	}
}

func TestAlign(t *testing.T) {
	cases := map[string][2]int{
		"top-left": {AlignStart, AlignStart}, "bottom": {AlignCenter, AlignEnd}, "right": {AlignEnd, AlignCenter},
		"center": {AlignCenter, AlignCenter}, "custom": {AlignFree, AlignFree}, "bottom-right": {AlignEnd, AlignEnd},
	}
	for mode, want := range cases {
		if h, v := Align(mode); h != want[0] || v != want[1] {
			t.Errorf("%q: %d %d", mode, h, v)
		}
	}
}

func TestPositionOfAnOlderConfig(t *testing.T) {
	cases := []struct {
		mode      string
		hasCustom bool
		want      string
		moved     bool
	}{
		{"", true, "custom", true},
		{"", false, "bottom", true},
		{"caret", true, "bottom", true},
		{"nonsense", false, "bottom", true},
		{"top-left", true, "top-left", false},
		{"custom", false, "custom", false},
	}
	for _, c := range cases {
		if got, moved := Position(c.mode, c.hasCustom); got != c.want || moved != c.moved {
			t.Errorf("Position(%q, %v) = %q, %v", c.mode, c.hasCustom, got, moved)
		}
	}
}
