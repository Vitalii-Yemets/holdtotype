package idleunload

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMinutesKeepsOnlyTheListedChoices(t *testing.T) {
	for in, want := range map[int]int{0: 0, 5: 5, 10: 10, 30: 30, 60: 60, -5: 0, 1: 0, 7: 0, 15: 0, 45: 0, 61: 0, 120: 0, 1 << 30: 0} {
		if got := Minutes(in); got != want {
			t.Errorf("Minutes(%d) = %d, want %d", in, got, want)
		}
	}
}

func TestTheConfigValueReadsAsMinutes(t *testing.T) {
	cases := map[string]int{
		`{}`:                             0,
		`{"models_idle_minutes":10}`:     10,
		`{"models_idle_minutes":60}`:     60,
		`{"models_idle_minutes":0}`:      0,
		`{"models_idle_minutes":20}`:     0,
		`{"models_idle_minutes":-30}`:    0,
		`{"models_idle_minutes":"10"}`:   0,
		`{"models_idle_minutes":null}`:   0,
		`{"models_idle_minutes":5.0}`:    0,
		`{"models_idle_minutes":true}`:   0,
		`{"models_idle_minutes":300000}`: 0,
	}
	for in, want := range cases {
		var c struct {
			M int `json:"models_idle_minutes"`
		}
		_ = json.Unmarshal([]byte(in), &c)
		if got := Minutes(c.M); got != want {
			t.Errorf("%s reads as %d, want %d", in, got, want)
		}
	}
}

func TestDueAfterTheChosenTime(t *testing.T) {
	now := time.Date(2026, 11, 15, 12, 0, 0, 0, time.UTC)
	for _, m := range []int{5, 10, 30, 60} {
		d := time.Duration(m) * time.Minute
		loaded := State{Loaded: true}
		loaded.Last = now.Add(-d + time.Second)
		if Due(m, loaded, now) {
			t.Errorf("%d minutes: due a second early", m)
		}
		loaded.Last = now.Add(-d)
		if !Due(m, loaded, now) {
			t.Errorf("%d minutes: not due on time", m)
		}
		loaded.Last = now.Add(-3 * d)
		if !Due(m, loaded, now) {
			t.Errorf("%d minutes: not due long after", m)
		}
	}
}

func TestNeverDueWhenOffBusyOrEmpty(t *testing.T) {
	now := time.Date(2026, 11, 15, 12, 0, 0, 0, time.UTC)
	long := now.Add(-48 * time.Hour)
	cases := []struct {
		name    string
		minutes int
		s       State
	}{
		{"switched off", 0, State{Loaded: true, Last: long}},
		{"a value not in the list", 15, State{Loaded: true, Last: long}},
		{"a negative value", -5, State{Loaded: true, Last: long}},
		{"recording, recognizing, editing or subtitles", 5, State{Busy: true, Loaded: true, Last: long}},
		{"nothing loaded", 5, State{Last: long}},
		{"no use counted yet", 5, State{Loaded: true}},
		{"used just now", 5, State{Loaded: true, Last: now}},
		{"used in the future after a clock change", 5, State{Loaded: true, Last: now.Add(time.Hour)}},
	}
	for _, c := range cases {
		if Due(c.minutes, c.s, now) {
			t.Errorf("%s: due", c.name)
		}
	}
}
