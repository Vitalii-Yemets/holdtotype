package commands

import (
	"strings"
	"testing"
	"time"
)

func TestBackspaces(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want int
	}{
		{"empty", "", 0},
		{"crlf pair is one press", "a\r\nb", 3},
		{"lone lf is one press", "a\nb", 3},
		{"cyrillic counts runes not bytes", "привет", 6},
		{"emoji with a skin tone counts per rune", "👍🏽", 2},
		{"two crlf pairs", "a\r\n\r\nb", 4},
	}
	for _, c := range cases {
		if got := Backspaces(c.in); got != c.want {
			t.Fatalf("%s: %q gives %d, want %d", c.name, c.in, got, c.want)
		}
	}
}

func TestCanUndoReasons(t *testing.T) {
	t0 := time.Date(2026, 11, 9, 12, 0, 0, 0, time.UTC)
	base := UndoRecord{Text: "hello", Wnd: 1, Focus: 2, At: t0, Typed: 5, Clicked: 7, Mode: "paste"}
	with := func(f func(r *UndoRecord)) *UndoRecord {
		r := base
		f(&r)
		return &r
	}
	cases := []struct {
		name   string
		rec    *UndoRecord
		now    time.Time
		wnd    uintptr
		focus  uintptr
		typed  uint64
		click  uint64
		ok     bool
		reason string
	}{
		{"fresh record", &base, t0.Add(time.Minute), 1, 2, 5, 7, true, ""},
		{"no record", nil, t0, 1, 2, 5, 7, false, "none"},
		{"empty text", with(func(r *UndoRecord) { r.Text = "" }), t0, 1, 2, 5, 7, false, "none"},
		{"key typed since", &base, t0, 1, 2, 6, 7, false, "typed"},
		{"mouse clicked since", &base, t0, 1, 2, 5, 8, false, "typed"},
		{"other window", &base, t0, 3, 2, 5, 7, false, "window"},
		{"other focus", &base, t0, 1, 4, 5, 7, false, "window"},
		{"focus unknown now", &base, t0, 1, 0, 5, 7, true, ""},
		{"focus unknown then", with(func(r *UndoRecord) { r.Focus = 0 }), t0, 1, 4, 5, 7, true, ""},
		{"enter was sent", with(func(r *UndoRecord) { r.EnterSent = true }), t0, 1, 2, 5, 7, false, "sent"},
		{"exactly five minutes", &base, t0.Add(5 * time.Minute), 1, 2, 5, 7, true, ""},
		{"older than five minutes", &base, t0.Add(5*time.Minute + time.Second), 1, 2, 5, 7, false, "old"},
		{"line break", with(func(r *UndoRecord) { r.Text = "a\nb" }), t0, 1, 2, 5, 7, false, "breaks"},
		{"carriage return", with(func(r *UndoRecord) { r.Text = "a\r\nb" }), t0, 1, 2, 5, 7, false, "breaks"},
		{"2000 runes", with(func(r *UndoRecord) { r.Text = strings.Repeat("я", 2000) }), t0, 1, 2, 5, 7, true, ""},
		{"2001 runes", with(func(r *UndoRecord) { r.Text = strings.Repeat("я", 2001) }), t0, 1, 2, 5, 7, false, "long"},
	}
	for _, c := range cases {
		ok, reason := CanUndo(c.rec, c.now, c.wnd, c.focus, c.typed, c.click)
		if ok != c.ok || reason != c.reason {
			t.Fatalf("%s: got (%v, %q), want (%v, %q)", c.name, ok, reason, c.ok, c.reason)
		}
	}
}

func TestCountsAsTyping(t *testing.T) {
	cases := []struct {
		name  string
		vk    uint32
		down  bool
		combo bool
		want  bool
	}{
		{"letter down", 0x41, true, false, true},
		{"backspace down", 0x08, true, false, true},
		{"letter up", 0x41, false, false, false},
		{"letter inside a shortcut", 0x41, true, true, false},
		{"shift", 0x10, true, false, false},
		{"ctrl", 0x11, true, false, false},
		{"alt", 0x12, true, false, false},
		{"left shift", 0xA0, true, false, false},
		{"right alt", 0xA5, true, false, false},
		{"left win", 0x5B, true, false, false},
		{"right win", 0x5C, true, false, false},
		{"apps key", 0x5D, true, false, true},
		{"enter", 0x0D, true, false, true},
	}
	for _, c := range cases {
		if got := CountsAsTyping(c.vk, c.down, c.combo); got != c.want {
			t.Fatalf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}
