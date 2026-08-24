package main

import (
	"log"
	"sync/atomic"
	"unsafe"
)

var (
	procGetGUIThreadInfo = user32.NewProc("GetGUIThreadInfo")
	procClientToScreen   = user32.NewProc("ClientToScreen")
)

type guiThreadInfo struct {
	Size      uint32
	Flags     uint32
	Active    uintptr
	Focus     uintptr
	Capture   uintptr
	MenuOwner uintptr
	MoveSize  uintptr
	Caret     uintptr
	RcCaret   rect
}

var ovPosMode atomic.Value

func setOverlayPos(mode string) {
	if !validOverlayPos(mode) {
		mode = ovPosBottom
	}
	ovPosMode.Store(mode)
}

func overlayPosMode() string {
	if v, ok := ovPosMode.Load().(string); ok && validOverlayPos(v) {
		return v
	}
	return ovPosBottom
}

func caretRect() (rect, bool) {
	var gui guiThreadInfo
	gui.Size = uint32(unsafe.Sizeof(gui))
	ok, _, _ := procGetGUIThreadInfo.Call(0, uintptr(unsafe.Pointer(&gui)))
	if ok == 0 || gui.Caret == 0 {
		return rect{}, false
	}
	if gui.RcCaret.Bottom <= gui.RcCaret.Top {
		return rect{}, false
	}
	top := point{X: gui.RcCaret.Left, Y: gui.RcCaret.Top}
	bottom := point{X: gui.RcCaret.Right, Y: gui.RcCaret.Bottom}
	procClientToScreen.Call(gui.Caret, uintptr(unsafe.Pointer(&top)))
	procClientToScreen.Call(gui.Caret, uintptr(unsafe.Pointer(&bottom)))
	return rect{Left: top.X, Top: top.Y, Right: bottom.X, Bottom: bottom.Y}, true
}

func anchorRect() rect {
	if r, ok := caretRect(); ok {
		return r
	}
	var pt point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	return rect{Left: pt.X, Top: pt.Y, Right: pt.X, Bottom: pt.Y}
}

func overlayOrigin(w, h int32) (int32, int32) {
	mode := overlayPosMode()
	if mode == ovPosCaret {
		anchor := anchorRect()
		wa := workAreaForPoint(anchor.Left, anchor.Top)
		gap := scaleDPI(18, dpiForPoint(anchor.Left, anchor.Top))
		x := anchor.Left - w/2
		if anchor.Right > anchor.Left {
			x = (anchor.Left+anchor.Right)/2 - w/2
		}
		y := anchor.Bottom + gap
		if y+h > wa.Bottom {
			y = anchor.Top - gap - h
		}
		if y < wa.Top {
			y = wa.Top + gap
		}
		if x < wa.Left+gap {
			x = wa.Left + gap
		}
		if x+w > wa.Right-gap {
			x = wa.Right - gap - w
		}
		log.Printf("плашка у курсора: якорь=%d,%d-%d,%d позиция=%d,%d",
			anchor.Left, anchor.Top, anchor.Right, anchor.Bottom, x, y)
		return x, y
	}
	anchor := anchorRect()
	wa := workAreaForPoint(anchor.Left, anchor.Top)
	margin := scaleDPI(28, dpiForPoint(anchor.Left, anchor.Top))
	x := wa.Left + (wa.Right-wa.Left-w)/2
	if mode == ovPosTop {
		return x, wa.Top + margin
	}
	return x, wa.Bottom - h - margin
}
