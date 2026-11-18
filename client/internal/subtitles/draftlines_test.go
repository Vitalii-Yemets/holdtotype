package subtitles

import "testing"

func TestALateDraftIsDropped(t *testing.T) {
	var d DraftLines
	if d.Show(3) {
		t.Fatal("a draft shown without a line")
	}
	d.Offer(3)
	if !d.Show(3) {
		t.Fatal("a draft for an open line is not shown")
	}
	d.Final(3)
	if d.Show(3) {
		t.Fatal("a draft shown after the whole sentence came")
	}
	if !d.TakeShown(3) || d.TakeShown(3) {
		t.Fatal("the final translation does not know the line had a draft")
	}
	d.Offer(4)
	if !d.Drop(4) || d.Open(4) || d.Drop(4) {
		t.Fatal("a dropped line stays open")
	}
}
