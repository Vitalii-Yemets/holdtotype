package appver

import "testing"

func TestParse(t *testing.T) {
	cases := []struct {
		in   string
		want [3]int
		ok   bool
	}{
		{"0.7.3", [3]int{0, 7, 3}, true},
		{"v0.7.3", [3]int{0, 7, 3}, true},
		{" v1.10.0 ", [3]int{1, 10, 0}, true},
		{"0.7", [3]int{}, false},
		{"v0.7.x", [3]int{}, false},
		{"", [3]int{}, false},
	}
	for _, c := range cases {
		got, ok := Parse(c.in)
		if ok != c.ok || (ok && got != c.want) {
			t.Errorf("Parse(%q) = %v, %v; want %v, %v", c.in, got, ok, c.want, c.ok)
		}
	}
}

func TestIsNewer(t *testing.T) {
	cases := []struct {
		latest, current string
		want            bool
	}{
		{"v0.8.0", "0.7.3", true},
		{"v0.7.4", "0.7.3", true},
		{"v1.0.0", "0.9.9", true},
		{"v0.7.3", "0.7.3", false},
		{"v0.7.2", "0.7.3", false},
		{"v0.10.0", "0.9.0", true},
		{"garbage", "0.7.3", false},
		{"v0.8.0", "not-a-version", false},
	}
	for _, c := range cases {
		if got := IsNewer(c.latest, c.current); got != c.want {
			t.Errorf("IsNewer(%q, %q) = %v; want %v", c.latest, c.current, got, c.want)
		}
	}
}
