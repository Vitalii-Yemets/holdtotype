package subwin

import (
	"reflect"
	"testing"
)

func runes(s string) int32 { return int32(len([]rune(s))) }

func TestWrap(t *testing.T) {
	cases := []struct {
		text  string
		width int32
		want  []string
	}{
		{"one two three four", 9, []string{"one two", "three", "four"}},
		{"one two three four", 100, []string{"one two three four"}},
		{"  spaced   out  ", 100, []string{"spaced out"}},
		{"", 10, nil},
		{"a supercalifragilistic b", 8, []string{"a", "supercal", "ifragili", "stic b"}},
		{"Доброе утро всем", 11, []string{"Доброе утро", "всем"}},
		{"one two", 0, []string{"one two"}},
	}
	for _, c := range cases {
		if got := Wrap(c.text, c.width, runes); !reflect.DeepEqual(got, c.want) {
			t.Errorf("Wrap(%q, %d) = %q, want %q", c.text, c.width, got, c.want)
		}
	}
}

func refs(keys ...int) []LineRef {
	var out []LineRef
	line := map[int]int{}
	for _, k := range keys {
		out = append(out, LineRef{Key: k, Line: line[k]})
		line[k]++
	}
	return out
}

func TestScrollFollowsTheLatestText(t *testing.T) {
	var s Scroll
	lines := refs(1, 1, 2, 3, 3, 3, 4)
	if s.Back() || s.Top(lines, 3) != 4 {
		t.Fatalf("top %d", s.Top(lines, 3))
	}
	if s.Top(lines[:2], 3) != 0 {
		t.Fatal("fewer lines than rows")
	}
}

func TestWheelHoldsTheTextStillWhileNewTextArrives(t *testing.T) {
	var s Scroll
	lines := refs(1, 1, 2, 3, 3, 3, 4)
	s.Wheel(lines, 3, -3)
	if !s.Back() || s.Top(lines, 3) != 1 {
		t.Fatalf("after one notch back: top %d", s.Top(lines, 3))
	}
	more := refs(1, 1, 2, 3, 3, 3, 4, 4, 5, 6)
	if s.Top(more, 3) != 1 {
		t.Fatalf("new text moved the rows being read: top %d", s.Top(more, 3))
	}
	dropped := refs(2, 3, 3, 3, 4, 4, 5, 6)
	if got := dropped[s.Top(dropped, 3)]; got != (LineRef{Key: 2, Line: 0}) {
		t.Fatalf("after the oldest row left: %+v", got)
	}
	s.Wheel(more, 3, -3)
	if s.Top(more, 3) != 0 {
		t.Fatalf("not past the beginning: %d", s.Top(more, 3))
	}
	s.Wheel(more, 3, 3)
	if !s.Back() || s.Top(more, 3) != 3 {
		t.Fatalf("one notch forward: %d", s.Top(more, 3))
	}
	s.Wheel(more, 3, 3)
	s.Wheel(more, 3, 3)
	if s.Back() || s.Top(more, 3) != 7 {
		t.Fatalf("scrolling to the end follows again: back %v top %d", s.Back(), s.Top(more, 3))
	}
}

func TestWheelWithNothingToScroll(t *testing.T) {
	var s Scroll
	s.Wheel(refs(1, 2), 3, -3)
	if s.Back() {
		t.Fatal("a window with all text visible was scrolled back")
	}
	s.Wheel(nil, 3, -3)
	if s.Back() {
		t.Fatal("an empty window was scrolled back")
	}
}

func TestScrolledRowThatWrapsDifferently(t *testing.T) {
	var s Scroll
	lines := refs(1, 2, 2, 2, 2, 3, 4, 5)
	s.Wheel(lines, 2, -3)
	if lines[s.Top(lines, 2)] != (LineRef{Key: 2, Line: 2}) {
		t.Fatalf("top %+v", lines[s.Top(lines, 2)])
	}
	narrower := refs(1, 2, 2, 3, 4, 5)
	if got := narrower[s.Top(narrower, 2)]; got != (LineRef{Key: 2, Line: 1}) {
		t.Fatalf("after a re-wrap %+v", got)
	}
	s.Follow()
	if s.Back() || s.Top(lines, 2) != 6 {
		t.Fatal("follow did not return to the end")
	}
}
