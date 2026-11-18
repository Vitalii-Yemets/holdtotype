package sentence

import "time"

const (
	DraftAfter     = 4 * time.Second
	DraftMinWords  = 6
	DraftMoreWords = 6
	DraftEvery     = 3 * time.Second
)

type Draft struct {
	Words int
	At    time.Time
}

func DraftDue(words int, first time.Time, last Draft, now time.Time) bool {
	if words < DraftMinWords || now.Sub(first) < DraftAfter {
		return false
	}
	if last.Words == 0 {
		return true
	}
	return words-last.Words >= DraftMoreWords && now.Sub(last.At) >= DraftEvery
}

func DraftAt(times []time.Time, last Draft) (time.Time, bool) {
	need := DraftMinWords
	if last.Words > 0 {
		need = last.Words + DraftMoreWords
	}
	if need > len(times) || len(times) == 0 {
		return time.Time{}, false
	}
	at := times[0].Add(DraftAfter)
	if times[need-1].After(at) {
		at = times[need-1]
	}
	if last.Words > 0 && last.At.Add(DraftEvery).After(at) {
		at = last.At.Add(DraftEvery)
	}
	return at, true
}
