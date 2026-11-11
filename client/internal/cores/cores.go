package cores

import (
	"fmt"
	"math/bits"
	"sort"
)

type Core struct {
	Mask  uint64
	Class int
}

type Split struct {
	Recognition uint64
	Translation uint64
}

func Count(mask uint64) int { return bits.OnesCount64(mask) }

func Divide(list []Core, threads int) (Split, string) {
	if threads < 1 {
		return Split{}, "the editor model has no thread count"
	}
	if len(list) == 0 {
		return Split{}, "the processor layout is unknown"
	}
	ordered := append([]Core(nil), list...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].Class > ordered[j].Class })
	var s Split
	logical := 0
	at := len(ordered)
	for at > 0 && logical < threads {
		at--
		s.Translation |= ordered[at].Mask
		logical += Count(ordered[at].Mask)
	}
	if logical < threads {
		return Split{}, fmt.Sprintf("the editor model uses %d threads and the processor has %d logical processors", threads, logical)
	}
	if at < 2 {
		return Split{}, fmt.Sprintf("only %d core(s) would be left for recognition", at)
	}
	for _, c := range ordered[:at] {
		s.Recognition |= c.Mask
	}
	return s, ""
}
