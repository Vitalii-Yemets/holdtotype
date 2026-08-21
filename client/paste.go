package main

import (
	"errors"
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
	return errors.New("буфер обмена занят другим приложением")
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
		return errors.New("GlobalAlloc не выделил память")
	}
	p, _, _ := procGlobalLock.Call(h)
	if p == 0 {
		procGlobalFree.Call(h)
		return errors.New("GlobalLock не удался")
	}
	dst := unsafe.Slice((*byte)(unsafe.Pointer(p)), size)
	src := unsafe.Slice((*byte)(unsafe.Pointer(&u[0])), size)
	copy(dst, src)
	procGlobalUnlock.Call(h)

	if r, _, _ := procSetClipboardData.Call(cfUnicodeText, h); r == 0 {
		procGlobalFree.Call(h)
		return errors.New("SetClipboardData не удался")
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

func pasteText(cfg *Config, text string) error {
	waitModifiersReleased(3 * time.Second)

	if cfg.PasteMode == "type" {
		return typeUnicode(text)
	}

	old, hadOld := "", false
	if cfg.RestoreClipboard {
		old, hadOld = getClipboardText()
	}
	if err := setClipboardText(text); err != nil {
		return err
	}
	time.Sleep(50 * time.Millisecond)
	if err := sendCtrlV(); err != nil {
		return err
	}
	if cfg.RestoreClipboard && hadOld {
		time.Sleep(400 * time.Millisecond)
		if cur, ok := getClipboardText(); ok && cur == text {
			_ = setClipboardText(old)
		}
	}
	return nil
}
