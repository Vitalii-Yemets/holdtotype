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
	}
	for _, c := range cases {
		if got := FindDuplicate(c.in); got != c.want {
			t.Errorf("FindDuplicate(%v) = %q; want %q", c.in, got, c.want)
		}
	}
}
