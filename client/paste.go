package main

import (
	"errors"
	"log"
	"strings"
	"time"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	inputKeyboard    = 1
	keyeventfKeyUp   = 0x0002
	keyeventfUnicode = 0x0004
	cfUnicodeText    = 13
	gmemMoveable     = 0x0002
)

type keybdInput struct {
	Vk        uint16
	Scan      uint16
	Flags     uint32
	Time      uint32
	ExtraInfo uintptr
}

type inputEvent struct {
	Type uint32
	_    uint32
	Ki   keybdInput
	_    [8]byte
}

func keyEvent(vk uint16, up bool) inputEvent {
	var flags uint32
	if up {
		flags = keyeventfKeyUp
	}
	return inputEvent{
		Type: inputKeyboard,
		Ki:   keybdInput{Vk: vk, Flags: flags, ExtraInfo: injectedMarker},
	}
}

func sendInputs(events []inputEvent) error {
	if len(events) == 0 {
		return nil
	}
	n, _, callErr := procSendInput.Call(
		uintptr(len(events)),
		uintptr(unsafe.Pointer(&events[0])),
		unsafe.Sizeof(events[0]),
	)
	if int(n) != len(events) {
		return errors.New("SendInput: " + callErr.Error())
	}
	return nil
}

func pressEnter() error {
	return sendInputs([]inputEvent{
		keyEvent(0x0D, false),
		keyEvent(0x0D, true),
	})
}

func sendCtrlV() error {
	return sendInputs([]inputEvent{
		keyEvent(0x11, false),
		keyEvent(0x56, false),
		keyEvent(0x56, true),
		keyEvent(0x11, true),
	})
}

func typeUnicode(s string) error {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	var events []inputEvent
	for _, r := range s {
		if r == '\n' {
			events = append(events, keyEvent(0x0D, false), keyEvent(0x0D, true))
			continue
		}
		for _, u := range utf16.Encode([]rune{r}) {
			down := inputEvent{Type: inputKeyboard, Ki: keybdInput{Scan: u, Flags: keyeventfUnicode, ExtraInfo: injectedMarker}}
			up := down
			up.Ki.Flags |= keyeventfKeyUp
			events = append(events, down, up)
		}
	}
	const chunk = 64
	for i := 0; i < len(events); i += chunk {
		end := i + chunk
		if end > len(events) {
			end = len(events)
		}
		if err := sendInputs(events[i:end]); err != nil {
			return err
		}
		time.Sleep(5 * time.Millisecond)
	}
	return nil
}

func openClipboardRetry() error {
	for i := 0; i < 20; i++ {
		r, _, _ := procOpenClipboard.Call(0)
		if r != 0 {
			return nil
		}
		time.Sleep(15 * time.Millisecond)
	}
	return errors.New("the clipboard is held by another application")
}

func setClipboardText(s string) error {
	s = strings.ReplaceAll(s, "\x00", "")
	u, err := windows.UTF16FromString(s)
	if err != nil {
		return err
	}
	if err := openClipboardRetry(); err != nil {
		return err
	}
	defer procCloseClipboard.Call()
	procEmptyClipboard.Call()

	size := len(u) * 2
	h, _, _ := procGlobalAlloc.Call(gmemMoveable, uintptr(size))
	if h == 0 {
		return errors.New("GlobalAlloc did not allocate memory")
	}
	p, _, _ := procGlobalLock.Call(h)
	if p == 0 {
		procGlobalFree.Call(h)
		return errors.New("GlobalLock failed")
	}
	dst := unsafe.Slice((*byte)(unsafe.Pointer(p)), size)
	src := unsafe.Slice((*byte)(unsafe.Pointer(&u[0])), size)
	copy(dst, src)
	procGlobalUnlock.Call(h)

	if r, _, _ := procSetClipboardData.Call(cfUnicodeText, h); r == 0 {
		procGlobalFree.Call(h)
		return errors.New("SetClipboardData failed")
	}
	return nil
}

func getClipboardText() (string, bool) {
	if err := openClipboardRetry(); err != nil {
		return "", false
	}
	defer procCloseClipboard.Call()
	h, _, _ := procGetClipboardData.Call(cfUnicodeText)
	if h == 0 {
		return "", false
	}
	p, _, _ := procGlobalLock.Call(h)
	if p == 0 {
		return "", false
	}
	defer procGlobalUnlock.Call(h)
	size, _, _ := procGlobalSize.Call(h)
	if size < 2 {
		return "", false
	}
	u := unsafe.Slice((*uint16)(unsafe.Pointer(p)), size/2)
	return windows.UTF16ToString(u), true
}

type clipFormat struct {
	id   uint32
	data []byte
}

var clipSkipFormats = map[uint32]bool{
	2: true, 3: true, 9: true, 14: true,
	0x0080: true, 0x0082: true, 0x0083: true, 0x008E: true,
}

func snapshotClipboard() ([]clipFormat, bool) {
	if err := openClipboardRetry(); err != nil {
		return nil, false
	}
	defer procCloseClipboard.Call()
	var out []clipFormat
	total := 0
	complete := true
	fmtID := uintptr(0)
	for {
		fmtID, _, _ = procEnumClipboardFormats.Call(fmtID)
		if fmtID == 0 {
			break
		}
		id := uint32(fmtID)
		if clipSkipFormats[id] {
			complete = false
			continue
		}
		h, _, _ := procGetClipboardData.Call(fmtID)
		if h == 0 {
			complete = false
			continue
		}
		size, _, _ := procGlobalSize.Call(h)
		if size == 0 {
			continue
		}
		total += int(size)
		if total > 64<<20 {
			complete = false
			break
		}
		p, _, _ := procGlobalLock.Call(h)
		if p == 0 {
			complete = false
			continue
		}
		buf := make([]byte, size)
		copy(buf, unsafe.Slice((*byte)(unsafe.Pointer(p)), size))
		procGlobalUnlock.Call(h)
		out = append(out, clipFormat{id: id, data: buf})
	}
	return out, complete
}

func restoreClipboard(fmts []clipFormat) error {
	if err := openClipboardRetry(); err != nil {
		return err
	}
	defer procCloseClipboard.Call()
	procEmptyClipboard.Call()
	for _, f := range fmts {
		h, _, _ := procGlobalAlloc.Call(gmemMoveable, uintptr(len(f.data)))
		if h == 0 {
			continue
		}
		p, _, _ := procGlobalLock.Call(h)
		if p == 0 {
			procGlobalFree.Call(h)
			continue
		}
		copy(unsafe.Slice((*byte)(unsafe.Pointer(p)), len(f.data)), f.data)
		procGlobalUnlock.Call(h)
		if r, _, _ := procSetClipboardData.Call(uintptr(f.id), h); r == 0 {
			procGlobalFree.Call(h)
		}
	}
	return nil
}

var errFocusMoved = errors.New("the input window changed before pasting")

func focusStillOn(expect uintptr) bool {
	if expect == 0 {
		return true
	}
	cur, _, _ := procGetForegroundWindow.Call()
	if cur == expect {
		return true
	}
	log.Printf("the input window changed: expected [%s], now [%s]", windowTitle(expect), windowTitle(cur))
	return false
}

func pasteText(cfg *Config, text string, expect uintptr) error {
	waitModifiersReleased(3 * time.Second)
	if cfg.PasteDelayMs > 0 {
		time.Sleep(time.Duration(cfg.PasteDelayMs) * time.Millisecond)
	}
	if !focusStillOn(expect) {
		return errFocusMoved
	}

	if cfg.PasteMode == "type" {
		return typeUnicode(text)
	}

	var snap []clipFormat
	if cfg.RestoreClipboard {
		var complete bool
		snap, complete = snapshotClipboard()
		if !complete {
			log.Printf("the clipboard cannot be preserved in full — typing character by character, leaving the clipboard alone")
			return typeUnicode(text)
		}
	}
	if err := setClipboardText(text); err != nil {
		return err
	}
	time.Sleep(50 * time.Millisecond)
	if err := sendCtrlV(); err != nil {
		return err
	}
	if cfg.RestoreClipboard && len(snap) > 0 {
		time.Sleep(400 * time.Millisecond)
		if cur, ok := getClipboardText(); ok && cur == text {
			_ = restoreClipboard(snap)
		}
	}
	return nil
}
