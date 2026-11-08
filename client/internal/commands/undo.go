package commands

import (
	"strings"
	"time"
	"unicode/utf8"
)

const (
	UndoWindow   = 5 * time.Minute
	UndoMaxRunes = 2000
)

type UndoRecord struct {
	Text      string
	Wnd       uintptr
	Focus     uintptr
	At        time.Time
	Typed     uint64
	Clicked   uint64
	EnterSent bool
	Mode      string
}

func Backspaces(s string) int {
	return utf8.RuneCountInString(s) - strings.Count(s, "\r\n")
}

func CanUndo(r *UndoRecord, now time.Time, wnd, focus uintptr, typed, clicked uint64) (bool, string) {
	switch {
	case r == nil || r.Text == "":
		return false, "none"
	case typed != r.Typed || clicked != r.Clicked:
		return false, "typed"
	case wnd != r.Wnd || (focus != 0 && r.Focus != 0 && focus != r.Focus):
		return false, "window"
	case r.EnterSent:
		return false, "sent"
	case now.Sub(r.At) > UndoWindow:
		return false, "old"
	case strings.ContainsAny(r.Text, "\r\n"):
		return false, "breaks"
	case utf8.RuneCountInString(r.Text) > UndoMaxRunes:
		return false, "long"
	}
	return true, ""
}

func isModifier(vk uint32) bool {
	switch {
	case vk >= 0x10 && vk <= 0x12, vk >= 0xA0 && vk <= 0xA5, vk == 0x5B, vk == 0x5C:
		return true
	}
	return false
}

func CountsAsTyping(vk uint32, down bool, comboMatched bool) bool {
	return down && !comboMatched && vk != 0 && !isModifier(vk)
}
