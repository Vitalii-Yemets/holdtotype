package sentence

import "testing"

func TestATailSplitInTwoGivesTheDraftLineToTheFirst(t *testing.T) {
	n := NewNumbers(4)
	seq, fresh := n.Reserve()
	if seq != 5 || !fresh {
		t.Fatalf("reserved %d %v", seq, fresh)
	}
	if again, fresh := n.Reserve(); again != 5 || fresh {
		t.Fatalf("an update took a new line: %d %v", again, fresh)
	}
	first, had := n.Take()
	second, had2 := n.Take()
	if first != 5 || !had || second != 6 || had2 {
		t.Fatalf("the two sentences got %d %v and %d %v", first, had, second, had2)
	}
	if next, _ := n.Reserve(); next != 7 || n.Last() != 7 {
		t.Fatalf("the next draft line %d, last %d", next, n.Last())
	}
}
