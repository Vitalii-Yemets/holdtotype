package subtitles

type DraftLines struct {
	open  map[int]bool
	shown map[int]bool
}

func (d *DraftLines) Offer(seq int) {
	if d.open == nil {
		d.open, d.shown = map[int]bool{}, map[int]bool{}
	}
	d.open[seq] = true
}

func (d *DraftLines) Open(seq int) bool { return d.open[seq] }

func (d *DraftLines) Show(seq int) bool {
	if !d.open[seq] {
		return false
	}
	d.shown[seq] = true
	return true
}

func (d *DraftLines) Final(seq int) { delete(d.open, seq) }

func (d *DraftLines) TakeShown(seq int) bool {
	had := d.shown[seq]
	delete(d.shown, seq)
	return had
}

func (d *DraftLines) Drop(seq int) bool {
	if !d.open[seq] {
		return false
	}
	delete(d.open, seq)
	delete(d.shown, seq)
	return true
}
