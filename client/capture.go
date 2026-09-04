package main

import (
	"fmt"
	"runtime"
	"sort"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"holdtotype/internal/appid"
)

const (
	wsOverlapped   = 0x00CF0000
	wsCaption      = 0x00C00000
	wsSysMenu      = 0x00080000
	wsVisible      = 0x10000000
	wsExTopmost    = 0x00000008
	wmDestroy      = 0x0002
	wmPaint        = 0x000F
	wmClose        = 0x0010
	wmUserRedraw   = 0x0400 + 1
	wmUserFinish   = 0x0400 + 2
	dtBlock        = 0x0001 | 0x0010
	dtCalcRect     = 0x0400
	capPadX        = 34
	capPadY        = 30
	capMinW        = 320
	capMaxW        = 560
	colorWindow    = 5
	defaultGUIFont = 17
	idcArrow       = 32512
)

var (
	gdi32                = windows.NewLazySystemDLL("gdi32.dll")
	procRegisterClassExW = user32.NewProc("RegisterClassExW")
	procCreateWindowExW  = user32.NewProc("CreateWindowExW")
	procDestroyWindow    = user32.NewProc("DestroyWindow")
	procDefWindowProcW   = user32.NewProc("DefWindowProcW")
	procPostQuitMessage  = user32.NewProc("PostQuitMessage")
	procPostMessageW     = user32.NewProc("PostMessageW")
	procTranslateMessage = user32.NewProc("TranslateMessage")
	procDispatchMessageW = user32.NewProc("DispatchMessageW")
	procInvalidateRect   = user32.NewProc("InvalidateRect")
	procBeginPaint       = user32.NewProc("BeginPaint")
	procEndPaint         = user32.NewProc("EndPaint")
	procDrawTextW        = user32.NewProc("DrawTextW")
	procGetClientRect    = user32.NewProc("GetClientRect")
	procGetSystemMetrics = user32.NewProc("GetSystemMetrics")
	procLoadCursorW      = user32.NewProc("LoadCursorW")
	procSetForegroundWnd = user32.NewProc("SetForegroundWindow")
	procGetStockObject   = gdi32.NewProc("GetStockObject")
	procSelectObject     = gdi32.NewProc("SelectObject")
	procSetBkMode        = gdi32.NewProc("SetBkMode")
)

type wndClassExW struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   uintptr
	Icon       uintptr
	Cursor     uintptr
	Background uintptr
	MenuName   *uint16
	ClassName  *uint16
	IconSm     uintptr
}

type rect struct{ Left, Top, Right, Bottom int32 }

type paintStruct struct {
	Hdc         uintptr
	Erase       int32
	RcPaint     rect
	Restore     int32
	IncUpdate   int32
	RgbReserved [32]byte
}

var (
	capMu     sync.Mutex
	capText   string
	capResult string
	capOK     bool
	capHwnd   uintptr

	capClassOnce sync.Once
	capWndProcCB uintptr
)

func captureWndProc(hwnd, msg, wParam, lParam uintptr) uintptr {
	switch msg {
	case wmPaint:
		var ps paintStruct
		hdc, _, _ := procBeginPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
		var rc rect
		procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
		border := rc
		brush, _, _ := procCreateSolidBrush.Call(colLine)
		procFrameRect.Call(hdc, uintptr(unsafe.Pointer(&border)), brush)
		procDeleteObject.Call(brush)
		procSelectObject.Call(hdc, uiFont(hwnd))
		procSetBkMode.Call(hdc, 1)
		procSetTextColor.Call(hdc, colGreen)
		capMu.Lock()
		text := capText
		capMu.Unlock()
		u, _ := windows.UTF16FromString(text)
		pad := scaleDPI(capPadX, dpiFor(hwnd))
		box := rect{Left: rc.Left + pad, Top: rc.Top, Right: rc.Right - pad, Bottom: rc.Bottom}
		measured := box
		procDrawTextW.Call(hdc, uintptr(unsafe.Pointer(&u[0])), uintptr(len(u)-1),
			uintptr(unsafe.Pointer(&measured)), dtBlock|dtCalcRect)
		if th := measured.Bottom - measured.Top; th > 0 && th < rc.Bottom-rc.Top {
			box.Top = rc.Top + (rc.Bottom-rc.Top-th)/2
			box.Bottom = box.Top + th
		}
		procDrawTextW.Call(hdc, uintptr(unsafe.Pointer(&u[0])), uintptr(len(u)-1), uintptr(unsafe.Pointer(&box)), dtBlock)
		procEndPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
		return 0
	case wmChromeRefresh:
		applyChrome(hwnd)
		return 0
	case wmUserRedraw:
		procInvalidateRect.Call(hwnd, 0, 1)
		return 0
	case wmUserFinish, wmClose:
		procDestroyWindow.Call(hwnd)
		return 0
	case wmDestroy:
		procPostQuitMessage.Call(0)
		return 0
	}
	r, _, _ := procDefWindowProcW.Call(hwnd, msg, wParam, lParam)
	return r
}

func captureHotkeyDialog(hook *hotkeyHook, current string) (string, bool) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	capClassOnce.Do(func() {
		capWndProcCB = syscall.NewCallback(captureWndProc)
		className, _ := windows.UTF16PtrFromString(appid.Class("CaptureWnd"))
		cursor, _, _ := procLoadCursorW.Call(0, idcArrow)
		bg, _, _ := procCreateSolidBrush.Call(colBg)
		wc := wndClassExW{
			Size:       uint32(unsafe.Sizeof(wndClassExW{})),
			WndProc:    capWndProcCB,
			Cursor:     cursor,
			Background: bg,
			ClassName:  className,
		}
		procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
	})

	capMu.Lock()
	capText = trf("cap.prompt", current)
	capResult, capOK = "", false
	capMu.Unlock()

	w, h, cx, cy := captureBox(capText)
	className, _ := windows.UTF16PtrFromString(appid.Class("CaptureWnd"))
	title, _ := windows.UTF16PtrFromString(tr("cap.title"))
	hwnd, _, _ := procCreateWindowExW.Call(
		wsExTopmost,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(title)),
		0x80000000|wsVisible,
		uintptr(cx), uintptr(cy), uintptr(w), uintptr(h),
		0, 0, 0, 0,
	)
	if hwnd == 0 {
		return "", false
	}
	applyDarkCaption(hwnd)
	capMu.Lock()
	capHwnd = hwnd
	capMu.Unlock()
	procSetForegroundWnd.Call(hwnd)

	hook.StartCapture(
		func(live string) {
			capMu.Lock()
			capText = live
			capMu.Unlock()
			procPostMessageW.Call(hwnd, wmUserRedraw, 0, 0)
		},
		func(combo string, ok bool) {
			capMu.Lock()
			capResult, capOK = combo, ok
			if ok {
				capText = trf("cap.selected", combo)
			} else {
				capText = tr("cap.cancelled")
			}
			capMu.Unlock()
			procPostMessageW.Call(hwnd, wmUserRedraw, 0, 0)
			time.Sleep(700 * time.Millisecond)
			procPostMessageW.Call(hwnd, wmUserFinish, 0, 0)
		},
	)

	var m msgStruct
	for {
		r, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(r) <= 0 {
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
	}

	hook.CancelCapture()

	capMu.Lock()
	defer capMu.Unlock()
	capHwnd = 0
	return capResult, capOK
}

func vkName(vk uint32) string {
	switch vk {
	case 0xA2, 0xA3, 0x11:
		return "ctrl"
	case 0xA0, 0xA1, 0x10:
		return "shift"
	case 0xA4, 0xA5, 0x12:
		return "alt"
	case 0x5B, 0x5C:
		return "win"
	case 0x20:
		return "space"
	case 0x09:
		return "tab"
	case 0x14:
		return "capslock"
	case 0x1B:
		return "esc"
	case 0x2D:
		return "insert"
	case 0x2E:
		return "delete"
	case 0x24:
		return "home"
	case 0x23:
		return "end"
	case 0x21:
		return "pageup"
	case 0x22:
		return "pagedown"
	case 0x13:
		return "pause"
	case 0x91:
		return "scrolllock"
	case 0x90:
		return "numlock"
	case 0x08:
		return "backspace"
	case 0x0D:
		return "enter"
	case 0x2C:
		return "printscreen"
	case 0x5D:
		return "menu"
	case 0x25:
		return "left"
	case 0x26:
		return "up"
	case 0x27:
		return "right"
	case 0x28:
		return "down"
	case 0x60:
		return "num0"
	case 0x61:
		return "num1"
	case 0x62:
		return "num2"
	case 0x63:
		return "num3"
	case 0x64:
		return "num4"
	case 0x65:
		return "num5"
	case 0x66:
		return "num6"
	case 0x67:
		return "num7"
	case 0x68:
		return "num8"
	case 0x69:
		return "num9"
	case 0x6A:
		return "multiply"
	case 0x6B:
		return "add"
	case 0x6D:
		return "subtract"
	case 0x6E:
		return "decimal"
	case 0x6F:
		return "divide"
	case 0xBA:
		return "semicolon"
	case 0xBB:
		return "equals"
	case 0xBC:
		return "comma"
	case 0xBD:
		return "minus"
	case 0xBE:
		return "period"
	case 0xBF:
		return "slash"
	case 0xC0:
		return "backquote"
	case 0xDB:
		return "lbracket"
	case 0xDC:
		return "backslash"
	case 0xDD:
		return "rbracket"
	case 0xDE:
		return "quote"
	case 0xAD:
		return "volumemute"
	case 0xAE:
		return "volumedown"
	case 0xAF:
		return "volumeup"
	case 0xB0:
		return "medianext"
	case 0xB1:
		return "mediaprev"
	case 0xB2:
		return "mediastop"
	case 0xB3:
		return "mediaplay"
	case 0xA6:
		return "browserback"
	case 0xA7:
		return "browserforward"
	case 0xA8:
		return "browserrefresh"
	case 0xA9:
		return "browserstop"
	case 0xAA:
		return "browsersearch"
	case 0xAB:
		return "browserfavorites"
	case 0xAC:
		return "browserhome"
	case 0xB4:
		return "launchmail"
	case 0xB5:
		return "launchmedia"
	case 0xB6:
		return "launchapp1"
	case 0xB7:
		return "launchapp2"
	}
	if vk >= 0x70 && vk <= 0x87 {
		return "f" + itoa(int(vk-0x70+1))
	}
	if vk >= 'A' && vk <= 'Z' {
		return string(rune('a' + vk - 'A'))
	}
	if vk >= '0' && vk <= '9' {
		return string(rune(vk))
	}
	if vk > 0 && vk < 256 {
		return fmt.Sprintf("vk0x%02X", vk)
	}
	return ""
}

func itoa(n int) string {
	if n < 10 {
		return string(rune('0' + n))
	}
	return string(rune('0'+n/10)) + string(rune('0'+n%10))
}

func comboString(vks map[uint32]bool) string {
	seen := map[string]bool{}
	var names []string
	for vk := range vks {
		n := vkName(vk)
		if n == "" || n == "esc" || seen[n] {
			continue
		}
		seen[n] = true
		names = append(names, n)
	}
	priority := map[string]int{"ctrl": 0, "shift": 1, "alt": 2, "win": 3}
	sort.Slice(names, func(i, j int) bool {
		pi, iok := priority[names[i]]
		pj, jok := priority[names[j]]
		if iok && jok {
			return pi < pj
		}
		if iok != jok {
			return iok
		}
		return names[i] < names[j]
	})
	out := ""
	for i, n := range names {
		if i > 0 {
			out += "+"
		}
		out += n
	}
	return out
}

func captureBox(text string) (w, h, x, y int32) {
	host := settingsHwnd.Load()
	if host == 0 {
		host, _, _ = procGetForegroundWindow.Call()
	}
	dpi := dpiFor(host)
	if dpi < 72 {
		dpi = 96
	}
	padX, padY := scaleDPI(capPadX, dpi), scaleDPI(capPadY, dpi)
	w, h = scaleDPI(420, dpi), scaleDPI(120, dpi)
	if hdc, _, _ := procGetDC.Call(0); hdc != 0 {
		defer procReleaseDC.Call(0, hdc)
		old, _, _ := procSelectObject.Call(hdc, uiFontDPI(dpi))
		if u, err := windows.UTF16FromString(text); err == nil {
			r := rect{Right: scaleDPI(capMaxW, dpi) - padX*2}
			procDrawTextW.Call(hdc, uintptr(unsafe.Pointer(&u[0])), uintptr(len(u)-1),
				uintptr(unsafe.Pointer(&r)), dtBlock|dtCalcRect)
			if r.Bottom > r.Top {
				w = r.Right - r.Left + padX*2
				h = r.Bottom - r.Top + padY*2
			}
		}
		if old != 0 {
			procSelectObject.Call(hdc, old)
		}
	}
	if min := scaleDPI(capMinW, dpi); w < min {
		w = min
	}
	if max := scaleDPI(capMaxW, dpi); w > max {
		w = max
	}
	if min := scaleDPI(96, dpi); h < min {
		h = min
	}
	wa := captureWorkArea(host)
	x = wa.Left + (wa.Right-wa.Left-w)/2
	y = wa.Top + (wa.Bottom-wa.Top-h)/2
	return w, h, x, y
}

func captureWorkArea(host uintptr) rect {
	if host != 0 && procMonitorFromWindow.Find() == nil && procGetMonitorInfoW.Find() == nil {
		if mon, _, _ := procMonitorFromWindow.Call(host, 2); mon != 0 {
			mi := monitorInfo{Size: uint32(unsafe.Sizeof(monitorInfo{}))}
			if r, _, _ := procGetMonitorInfoW.Call(mon, uintptr(unsafe.Pointer(&mi))); r != 0 {
				return mi.Work
			}
		}
	}
	var pt point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	return workAreaForPoint(pt.X, pt.Y)
}
