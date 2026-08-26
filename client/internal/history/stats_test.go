package history

import "testing"

func TestStatsCountsSince(t *testing.T) {
	s := &Store{items: []Item{
		{At: 100, Text: "старое"},
		{At: 900, Text: "раз"},
		{At: 950, Text: "два три"},
	}}
	n, chars := s.Stats(500)
	if n != 2 {
		t.Fatalf("получено %d диктовок, старые не должны считаться", n)
	}
	if chars != len([]rune("раз"))+len([]rune("два три")) {
		t.Fatalf("получено %d знаков", chars)
	}
}
