package hotkeys

import "testing"

func TestFindDuplicate(t *testing.T) {
	cases := []struct {
		in   []string
		want string
	}{
		{[]string{"ctrl+win", "", "alt+f9"}, ""},
		{[]string{"ctrl+win", "ctrl+win"}, "ctrl+win"},
		{[]string{"ctrl+win", "CTRL+WIN"}, "ctrl+win"},
		{[]string{"ctrl+win", " ctrl+win "}, "ctrl+win"},
		{[]string{"", "", ""}, ""},
		{[]string{"a", "b", "c", "b"}, "b"},
		{nil, ""},
		{[]string{"ctrl+win+z", "win+ctrl+z"}, "ctrl+win+z"},
		{[]string{"ctrl+win", "ctrl+win+z"}, ""},
		{[]string{"alt+shift+f9", " Shift + Alt + F9 "}, "shift+alt+f9"},
	}
	for _, c := range cases {
		if got := FindDuplicate(c.in); got != c.want {
			t.Errorf("FindDuplicate(%v) = %q; want %q", c.in, got, c.want)
		}
	}
}

func TestSame(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"ctrl+win+z", "win+ctrl+z", true},
		{"CTRL+WIN", " ctrl + win ", true},
		{"ctrl+win", "ctrl+win+z", false},
		{"", "", false},
		{"", "ctrl+win", false},
	}
	for _, c := range cases {
		if got := Same(c.a, c.b); got != c.want {
			t.Errorf("Same(%q, %q) = %v; want %v", c.a, c.b, got, c.want)
		}
	}
}
