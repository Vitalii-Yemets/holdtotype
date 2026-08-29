package main

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

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
	Index   int    `json:"index"`
	Work    rect   `json:"-"`
	Screen  rect   `json:"-"`
	W       int32  `json:"w"`
	H       int32  `json:"h"`
	Primary bool   `json:"primary"`
	Name    string `json:"name"`
	Device  string `json:"-"`
}

type monitorInfoEx struct {
	Size    uint32
	Monitor rect
	Work    rect
	Flags   uint32
	Device  [32]uint16
}

type luid struct {
	Low  uint32
	High int32
}

type dcPathSourceInfo struct {
	AdapterID   luid
	ID          uint32
	ModeInfoIdx uint32
	StatusFlags uint32
}

type dcRational struct {
	Numerator   uint32
	Denominator uint32
}

type dcPathTargetInfo struct {
	AdapterID        luid
	ID               uint32
	ModeInfoIdx      uint32
	OutputTechnology uint32
	Rotation         uint32
	Scaling          uint32
	RefreshRate      dcRational
	ScanLineOrdering uint32
	TargetAvailable  int32
	StatusFlags      uint32
}

type dcPathInfo struct {
	SourceInfo dcPathSourceInfo
	TargetInfo dcPathTargetInfo
	Flags      uint32
}

type dcModeInfo struct {
	InfoType  uint32
	ID        uint32
	AdapterID luid
	Data      [48]byte
}

type dcDeviceInfoHeader struct {
	Type      uint32
	Size      uint32
	AdapterID luid
	ID        uint32
}

type dcTargetDeviceName struct {
	Header            dcDeviceInfoHeader
	Flags             uint32
	OutputTechnology  uint32
	EdidManufactureID uint16
	EdidProductCodeID uint16
	ConnectorInstance uint32
	FriendlyName      [64]uint16
	DevicePath        [128]uint16
}

type dcSourceDeviceName struct {
	Header      dcDeviceInfoHeader
	GdiDeviceName [32]uint16
}

const (
	qdcOnlyActivePaths = 2
	dcInfoGetSourceName = 1
	dcInfoGetTargetName = 2
)

var (
	procGetDisplayConfigBufferSizes = user32.NewProc("GetDisplayConfigBufferSizes")
	procQueryDisplayConfig          = user32.NewProc("QueryDisplayConfig")
	procDisplayConfigGetDeviceInfo  = user32.NewProc("DisplayConfigGetDeviceInfo")
)

// monitorNames asks Windows what the screens are actually called, keyed by the
// device name the enumeration hands out (\\.\DISPLAY1 and friends). Screens
// whose firmware keeps quiet simply stay out of the map.
func monitorNames() map[string]string {
	out := map[string]string{}
	if procGetDisplayConfigBufferSizes.Find() != nil || procQueryDisplayConfig.Find() != nil || procDisplayConfigGetDeviceInfo.Find() != nil {
		return out
	}
	var pathCount, modeCount uint32
	if r, _, _ := procGetDisplayConfigBufferSizes.Call(qdcOnlyActivePaths,
		uintptr(unsafe.Pointer(&pathCount)), uintptr(unsafe.Pointer(&modeCount))); r != 0 || pathCount == 0 {
		return out
	}
	paths := make([]dcPathInfo, pathCount)
	modes := make([]dcModeInfo, modeCount)
	if modeCount == 0 {
		modes = make([]dcModeInfo, 1)
	}
	if r, _, _ := procQueryDisplayConfig.Call(qdcOnlyActivePaths,
		uintptr(unsafe.Pointer(&pathCount)), uintptr(unsafe.Pointer(&paths[0])),
		uintptr(unsafe.Pointer(&modeCount)), uintptr(unsafe.Pointer(&modes[0])), 0); r != 0 {
		return out
	}
	for i := uint32(0); i < pathCount; i++ {
		src := dcSourceDeviceName{}
		src.Header.Type = dcInfoGetSourceName
		src.Header.Size = uint32(unsafe.Sizeof(src))
		src.Header.AdapterID = paths[i].SourceInfo.AdapterID
		src.Header.ID = paths[i].SourceInfo.ID
		if r, _, _ := procDisplayConfigGetDeviceInfo.Call(uintptr(unsafe.Pointer(&src))); r != 0 {
			continue
		}
		tgt := dcTargetDeviceName{}
		tgt.Header.Type = dcInfoGetTargetName
		tgt.Header.Size = uint32(unsafe.Sizeof(tgt))
		tgt.Header.AdapterID = paths[i].TargetInfo.AdapterID
		tgt.Header.ID = paths[i].TargetInfo.ID
		if r, _, _ := procDisplayConfigGetDeviceInfo.Call(uintptr(unsafe.Pointer(&tgt))); r != 0 {
			continue
		}
		name := strings.TrimSpace(windows.UTF16ToString(tgt.FriendlyName[:]))
		device := windows.UTF16ToString(src.GdiDeviceName[:])
		if name == "" || device == "" {
			continue
		}
		out[device] = name
	}
	return out
}

var (
	monCbOnce sync.Once
	monCb     uintptr
	monMu     sync.Mutex
	monAcc    []monitorEntry
)

func monEnumProc(hmon, hdc uintptr, rc *rect, lp uintptr) uintptr {
	mi := monitorInfoEx{Size: uint32(unsafe.Sizeof(monitorInfoEx{}))}
	if r, _, _ := procGetMonitorInfoW.Call(hmon, uintptr(unsafe.Pointer(&mi))); r != 0 {
		monAcc = append(monAcc, monitorEntry{
			Index: len(monAcc), Work: mi.Work, Screen: mi.Monitor,
			W: mi.Monitor.Right - mi.Monitor.Left, H: mi.Monitor.Bottom - mi.Monitor.Top,
			Primary: mi.Flags&1 != 0,
			Device:  windows.UTF16ToString(mi.Device[:]),
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
	out := append([]monitorEntry(nil), monAcc...)
	names := monitorNames()
	seen := map[string]int{}
	for _, n := range names {
		seen[n]++
	}
	for i := range out {
		n := names[out[i].Device]
		if n == "" {
			continue
		}
		if seen[n] > 1 {
			n = n + " #" + strconv.Itoa(out[i].Index+1)
		}
		out[i].Name = n
	}
	return out
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
		log.Printf("plate at the cursor: anchor=%d,%d-%d,%d position=%d,%d",
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
