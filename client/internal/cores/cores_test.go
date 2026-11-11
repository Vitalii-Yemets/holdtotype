package cores

import (
	"strings"
	"testing"
)

func layout(groups ...[3]int) []Core {
	var out []Core
	bit := 0
	for _, g := range groups {
		count, threads, class := g[0], g[1], g[2]
		for i := 0; i < count; i++ {
			var m uint64
			for t := 0; t < threads; t++ {
				m |= 1 << bit
				bit++
			}
			out = append(out, Core{Mask: m, Class: class})
		}
	}
	return out
}

func TestHybridProcessor(t *testing.T) {
	s, why := Divide(layout([3]int{6, 2, 1}, [3]int{8, 1, 0}), 8)
	if why != "" {
		t.Fatal(why)
	}
	if s.Recognition != 0xfff || s.Translation != 0xff000 {
		t.Fatalf("split %x %x", s.Recognition, s.Translation)
	}
	if Count(s.Recognition) != 12 || Count(s.Translation) != 8 {
		t.Fatalf("counts %d %d", Count(s.Recognition), Count(s.Translation))
	}
}

func TestEfficiencyCoresListedFirst(t *testing.T) {
	s, why := Divide(layout([3]int{8, 1, 0}, [3]int{6, 2, 1}), 8)
	if why != "" {
		t.Fatal(why)
	}
	if s.Translation != 0xff || s.Recognition != 0xfff00 {
		t.Fatalf("split %x %x", s.Recognition, s.Translation)
	}
}

func TestSameCores(t *testing.T) {
	cases := []struct {
		cores, threads, want int
		rec, tr              int
	}{
		{8, 2, 8, 8, 8},
		{6, 2, 8, 4, 8},
		{4, 2, 4, 4, 4},
		{8, 1, 4, 4, 4},
		{16, 2, 6, 26, 6},
		{5, 2, 5, 4, 6},
	}
	for _, c := range cases {
		s, why := Divide(layout([3]int{c.cores, c.threads, 0}), c.want)
		if why != "" {
			t.Errorf("%d cores x %d, %d threads: %s", c.cores, c.threads, c.want, why)
			continue
		}
		if Count(s.Recognition) != c.rec || Count(s.Translation) != c.tr || s.Recognition&s.Translation != 0 {
			t.Errorf("%d cores x %d, %d threads: %x %x", c.cores, c.threads, c.want, s.Recognition, s.Translation)
		}
	}
	s, _ := Divide(layout([3]int{6, 2, 0}), 8)
	if s.Recognition != 0xf || s.Translation != 0xff0 {
		t.Fatalf("whole cores from the end: %x %x", s.Recognition, s.Translation)
	}
}

func TestNotDivided(t *testing.T) {
	cases := []struct {
		list    []Core
		threads int
		reason  string
	}{
		{layout([3]int{4, 2, 0}), 8, "only 0 core"},
		{layout([3]int{4, 1, 0}), 4, "only 0 core"},
		{layout([3]int{3, 2, 0}), 4, "only 1 core"},
		{layout([3]int{2, 2, 0}), 8, "8 threads"},
		{layout([3]int{8, 2, 0}), 0, "thread count"},
		{nil, 8, "unknown"},
	}
	for _, c := range cases {
		s, why := Divide(c.list, c.threads)
		if why == "" || !strings.Contains(why, c.reason) || s != (Split{}) {
			t.Errorf("%d cores, %d threads: %q %+v", len(c.list), c.threads, why, s)
		}
	}
}
