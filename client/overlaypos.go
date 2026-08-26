package main

import (
	"log"
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"

	"holdtotype/internal/ovplace"
)

var (
	procGetGUIThreadInfo     = user32.NewProc("GetGUIThreadInfo")
	procClientToScreen       = user32.NewProc("ClientToScreen")
	procEnumDisplayMonitors  = user32.NewProc("EnumDisplayMonitors")
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

var (
	ovPosMode    atomic.Value
	ovMonitorSel atomic.Value
	ovCustomXY   atomic.Value
)

func setOverlayPos(mode string) {
	if !validOverlayPos(mode) {
		mode = ovPosBottom
	}
	ovPosMode.Store(mode)
}

func setOverlayPlacement(cfg *Config) {
	setOverlayPos(cfg.OverlayPos)
	mon := cfg.OverlayMonitor
	if !validOverlayMonitor(mon) {
		mon = ""
	}
	ovMonitorSel.Store(mon)
	custom := map[string]ovplace.Frac{}
	for k, v := range cfg.OverlayXY {
		custom[k] = v
	}
	ovCustomXY.Store(custom)
}

func overlayPosMode() string {
	if v, ok := ovPosMode.Load().(string); ok && validOverlayPos(v) {
		return v
	}
	return ovPosBottom
}

func overlayMonitorPick() string {
	if v, ok := ovMonitorSel.Load().(string); ok {
		return v
	}
	return ""
}

func overlayCustomMap() map[string]ovplace.Frac {
	if v, ok := ovCustomXY.Load().(map[string]ovplace.Frac); ok {
		return v
	}
	return nil
}

type monitorEntry struct {
	Index   int   `json:"index"`
	Work    rect  `json:"-"`
	Screen  rect  `json:"-"`
	W       int32 `json:"w"`
	H       int32 `json:"h"`
	Primary bool  `json:"primary"`
}

var (
	monCbOnce sync.Once
	monCb     uintptr
	monMu     sync.Mutex
	monAcc    []monitorEntry
)

func monEnumProc(hmon, hdc uintptr, rc *rect, lp uintptr) uintptr {
	mi := monitorInfo{Size: uint32(unsafe.Sizeof(monitorInfo{}))}
	if r, _, _ := procGetMonitorInfoW.Call(hmon, uintptr(unsafe.Pointer(&mi))); r != 0 {
		monAcc = append(monAcc, monitorEntry{
			Index: len(monAcc), Work: mi.Work, Screen: mi.Monitor,
			W: mi.Monitor.Right - mi.Monitor.Left, H: mi.Monitor.Bottom - mi.Monitor.Top,
			Primary: mi.Flags&1 != 0,
		})
	}
	return 1
}

func listMonitors() []monitorEntry {
	if procEnumDisplayMonitors.Find() != nil || procGetMonitorInfoW.Find() != nil {
		return nil
	}
	monCbOnce.Do(func() { monCb = syscall.NewCallback(monEnumProc) })
	monMu.Lock()
	defer monMu.Unlock()
	monAcc = nil
	procEnumDisplayMonitors.Call(0, 0, monCb, 0)
	return append([]monitorEntry(nil), monAcc...)
}

func monitorRectsForPoint(x, y int32) (work, screen rect) {
	if procMonitorFromPoint.Find() == nil && procGetMonitorInfoW.Find() == nil {
		mon, _, _ := procMonitorFromPoint.Call(uintptr(uint32(x))|uintptr(uint32(y))<<32, 2)
		if mon != 0 {
			mi := monitorInfo{Size: uint32(unsafe.Sizeof(monitorInfo{}))}
			if r, _, _ := procGetMonitorInfoW.Call(mon, uintptr(unsafe.Pointer(&mi))); r != 0 {
				return mi.Work, mi.Monitor
			}
		}
	}
	var wa rect
	procSystemParametersInfoW.Call(0x30, 0, uintptr(unsafe.Pointer(&wa)), 0)
	return wa, wa
}

// overlayArea answers where the plate lives: the chosen monitor's work area,
// or the one under the anchor when the choice is "the screen with the cursor".
func overlayArea() (work, screen rect) {
	pick := overlayMonitorPick()
	if overlayPosMode() != ovPosCaret && pick != "" && pick != "cursor" {
		idx := 0
		for _, r := range pick {
			idx = idx*10 + int(r-'0')
		}
		mons := listMonitors()
		if idx >= 0 && idx < len(mons) {
			return mons[idx].Work, mons[idx].Screen
		}
	}
	anchor := anchorRect()
	return monitorRectsForPoint(anchor.Left, anchor.Top)
}

func toPlace(r rect) ovplace.Rect {
	return ovplace.Rect{Left: r.Left, Top: r.Top, Right: r.Right, Bottom: r.Bottom}
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
	wa, screen := overlayArea()
	margin := scaleDPI(28, dpiForPoint(wa.Left+(wa.Right-wa.Left)/2, wa.Top+(wa.Bottom-wa.Top)/2))
	if mode == ovPosCustom {
		if f, ok := overlayCustomMap()[ovplace.ResKey(toPlace(screen))]; ok {
			return ovplace.Custom(f, toPlace(wa), w, h)
		}
		mode = ovPosBottom
	}
	return ovplace.Anchor(mode, toPlace(wa), w, h, margin)
}
