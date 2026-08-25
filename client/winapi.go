package main

import (
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"holdtotype/internal/appid"
)

const injectedMarker uintptr = 0x56325454

var (
	user32   = windows.NewLazySystemDLL("user32.dll")
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	shell32  = windows.NewLazySystemDLL("shell32.dll")
	shcore   = windows.NewLazySystemDLL("shcore.dll")

	procSetWindowsHookExW    = user32.NewProc("SetWindowsHookExW")
	procCallNextHookEx       = user32.NewProc("CallNextHookEx")
	procGetMessageW          = user32.NewProc("GetMessageW")
	procSendInput            = user32.NewProc("SendInput")
	procGetAsyncKeyState     = user32.NewProc("GetAsyncKeyState")
	procOpenClipboard        = user32.NewProc("OpenClipboard")
	procEnumClipboardFormats = user32.NewProc("EnumClipboardFormats")
	procGetForegroundWindow  = user32.NewProc("GetForegroundWindow")
	procCloseClipboard       = user32.NewProc("CloseClipboard")
	procEmptyClipboard       = user32.NewProc("EmptyClipboard")
	procGetClipboardData     = user32.NewProc("GetClipboardData")
	procSetClipboardData     = user32.NewProc("SetClipboardData")
	procMessageBoxW          = user32.NewProc("MessageBoxW")

	procGlobalAlloc  = kernel32.NewProc("GlobalAlloc")
	procGlobalFree   = kernel32.NewProc("GlobalFree")
	procGlobalLock   = kernel32.NewProc("GlobalLock")
	procGlobalUnlock = kernel32.NewProc("GlobalUnlock")
	procGlobalSize   = kernel32.NewProc("GlobalSize")
	procCreateMutexW = kernel32.NewProc("CreateMutexW")

	procShellExecuteW = shell32.NewProc("ShellExecuteW")

	dwmapi                    = windows.NewLazySystemDLL("dwmapi.dll")
	procDwmSetWindowAttribute = dwmapi.NewProc("DwmSetWindowAttribute")

	procGetWindowLongPtrW = user32.NewProc("GetWindowLongPtrW")
	procSetWindowLongPtrW = user32.NewProc("SetWindowLongPtrW")
	procSetWindowPos      = user32.NewProc("SetWindowPos")
	procRedrawWindow      = user32.NewProc("RedrawWindow")
	procSendMessageW      = user32.NewProc("SendMessageW")
	procReleaseCapture    = user32.NewProc("ReleaseCapture")
)

func applyDarkCaption(hwnd uintptr) {
	if hwnd == 0 {
		return
	}
	set := func(attr uintptr, value uint32) {
		procDwmSetWindowAttribute.Call(hwnd, attr, uintptr(unsafe.Pointer(&value)), 4)
	}
	set(20, 1)
	set(35, uint32(colBg))
	set(36, uint32(colGreen))
	set(34, uint32(colGreenLo))
	set(33, 2)
}

var (
	procSetClassLongPtrW = user32.NewProc("SetClassLongPtrW")
	procFindWindowW      = user32.NewProc("FindWindowW")
	procIsWindowVisible  = user32.NewProc("IsWindowVisible")
)

var (
	darkBrushOnce sync.Once
	darkBrush     uintptr
	darkBgMu      sync.Mutex
	darkBgDone    = map[uintptr]bool{}
)

func setDarkClientBackground(hwnd uintptr) {
	darkBgMu.Lock()
	if darkBgDone[hwnd] {
		darkBgMu.Unlock()
		return
	}
	darkBgDone[hwnd] = true
	darkBgMu.Unlock()
	darkBrushOnce.Do(func() {
		darkBrush, _, _ = procCreateSolidBrush.Call(0x0C0F0B)
	})
	procSetClassLongPtrW.Call(hwnd, ^uintptr(9), darkBrush)
}

const offscreenXY int32 = -32000

func offscreenPos() uintptr {
	v := offscreenXY
	return uintptr(uint32(v))
}

func hideWebViewWindowEarly(title string) func() {
	done := make(chan struct{})
	go func() {
		cls, _ := windows.UTF16PtrFromString("webview")
		t, _ := windows.UTF16PtrFromString(title)
		var styled uintptr
		for {
			select {
			case <-done:
				return
			default:
			}
			h, _, _ := procFindWindowW.Call(uintptr(unsafe.Pointer(cls)), uintptr(unsafe.Pointer(t)))
			if h != 0 {
				var rc rect
				procGetWindowRect.Call(h, uintptr(unsafe.Pointer(&rc)))
				if rc.Left > offscreenXY+1000 {
					procSetWindowPos.Call(h, 0, offscreenPos(), offscreenPos(), 0, 0, 0x0001|0x0004|0x0010)
				}
				if h != styled {
					setDarkClientBackground(h)
					applyDarkCaption(h)
					styled = h
				}
			}
			time.Sleep(2 * time.Millisecond)
		}
	}()
	return func() { close(done) }
}

func revealWindowCentered(hwnd uintptr, w, h int) {
	var wa rect
	procSystemParametersInfoW.Call(0x30, 0, uintptr(unsafe.Pointer(&wa)), 0)
	x := wa.Left + (wa.Right-wa.Left-int32(w))/2
	y := wa.Top + (wa.Bottom-wa.Top-int32(h))/2
	procSetWindowPos.Call(hwnd, 0, uintptr(uint32(x)), uintptr(uint32(y)), uintptr(w), uintptr(h), 0x0004|0x0010)
	procShowWindow.Call(hwnd, 5)
	procSetForegroundWnd.Call(hwnd)
}

func makeBorderless(hwnd uintptr) {
	if hwnd == 0 {
		return
	}
	const gwlStyle = ^uintptr(15)
	style, _, _ := procGetWindowLongPtrW.Call(hwnd, gwlStyle)
	style &^= 0x00C00000 | 0x00080000 | 0x00020000 | 0x00010000
	style |= 0x00040000
	procSetWindowLongPtrW.Call(hwnd, gwlStyle, style)
	procSetWindowPos.Call(hwnd, 0, 0, 0, 0, 0, 0x0027)
}

func beginWindowDrag(hwnd uintptr) {
	procReleaseCapture.Call()
	procSendMessageW.Call(hwnd, 0x00A1, 2, 0)
}

type pointL struct{ X, Y int32 }

type minMaxInfo struct {
	Reserved     pointL
	MaxSize      pointL
	MaxPosition  pointL
	MinTrackSize pointL
	MaxTrackSize pointL
}

var (
	procCallWindowProcW = user32.NewProc("CallWindowProcW")

	minSizeOld               uintptr
	minSizeW, minSizeH       int32
	minSizeDIPW, minSizeDIPH int32
	minSizeOnce              sync.Once
	minSizeCB                uintptr
)

var (
	procGetWindowRect  = user32.NewProc("GetWindowRect")
	lastWndW, lastWndH int32
)

func minSizeProc(hwnd, msg, wp, lp uintptr) uintptr {
	if msg == 0x02E0 && lp != 0 {
		suggested := (*rect)(unsafe.Pointer(lp))
		procSetWindowPos.Call(hwnd, 0, uintptr(suggested.Left), uintptr(suggested.Top),
			uintptr(suggested.Right-suggested.Left), uintptr(suggested.Bottom-suggested.Top), 0x0004|0x0010)
		dpi := int32(wp & 0xFFFF)
		if dpi >= 72 {
			minSizeW = scaleDPI(minSizeDIPW, dpi)
			minSizeH = scaleDPI(minSizeDIPH, dpi)
		}
	}
	r, _, _ := procCallWindowProcW.Call(minSizeOld, hwnd, msg, wp, lp)
	if msg == 0x0024 && lp != 0 {
		mmi := (*minMaxInfo)(unsafe.Pointer(lp))
		mmi.MinTrackSize.X = minSizeW
		mmi.MinTrackSize.Y = minSizeH
	}
	if msg == 0x0005 && wp != 1 {
		var rc rect
		procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
		lastWndW, lastWndH = rc.Right-rc.Left, rc.Bottom-rc.Top
	}
	return r
}

func applyMinSize(hwnd uintptr, w, h int32) {
	minSizeOnce.Do(func() { minSizeCB = syscall.NewCallback(minSizeProc) })
	minSizeDIPW, minSizeDIPH = w, h
	dpi := dpiFor(hwnd)
	minSizeW, minSizeH = scaleDPI(w, dpi), scaleDPI(h, dpi)
	const gwlpWndproc = ^uintptr(3)
	old, _, _ := procSetWindowLongPtrW.Call(hwnd, gwlpWndproc, minSizeCB)
	minSizeOld = old
}

func msgBox(title, text string) {
	t, _ := windows.UTF16PtrFromString(title)
	m, _ := windows.UTF16PtrFromString(text)
	procMessageBoxW.Call(0, uintptr(unsafe.Pointer(m)), uintptr(unsafe.Pointer(t)), 0x10)
}

func shellOpen(path string) {
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	verb, _ := windows.UTF16PtrFromString("open")
	p, _ := windows.UTF16PtrFromString(abs)
	r, _, _ := procShellExecuteW.Call(0, uintptr(unsafe.Pointer(verb)), uintptr(unsafe.Pointer(p)), 0, 0, 1)
	if r <= 32 {
		notepad, _ := windows.UTF16PtrFromString("notepad.exe")
		procShellExecuteW.Call(0, uintptr(unsafe.Pointer(verb)), uintptr(unsafe.Pointer(notepad)), uintptr(unsafe.Pointer(p)), 0, 1)
	}
}

func msgBoxYesNo(title, text string) bool {
	t, _ := windows.UTF16PtrFromString(title)
	m, _ := windows.UTF16PtrFromString(text)
	r, _, _ := procMessageBoxW.Call(0, uintptr(unsafe.Pointer(m)), uintptr(unsafe.Pointer(t)), 0x04|0x20)
	return r == 6
}

func shellOpenURL(url string) {
	verb, _ := windows.UTF16PtrFromString("open")
	p, _ := windows.UTF16PtrFromString(url)
	procShellExecuteW.Call(0, uintptr(unsafe.Pointer(verb)), uintptr(unsafe.Pointer(p)), 0, 0, 1)
}

func acquireSingleInstance() bool {
	name, _ := windows.UTF16PtrFromString(appid.MutexName)
	_, _, err := procCreateMutexW.Call(0, 0, uintptr(unsafe.Pointer(name)))
	return err != windows.ERROR_ALREADY_EXISTS
}

func waitModifiersReleased(timeout time.Duration) {
	mods := []uintptr{0x10, 0x11, 0x12, 0x5B, 0x5C}
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		down := false
		for _, vk := range mods {
			s, _, _ := procGetAsyncKeyState.Call(vk)
			if s&0x8000 != 0 {
				down = true
				break
			}
		}
		if !down {
			time.Sleep(30 * time.Millisecond)
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
}

var procGetWindowTextW = user32.NewProc("GetWindowTextW")

func windowTitle(hwnd uintptr) string {
	if hwnd == 0 || hwnd == 1 {
		return ""
	}
	buf := make([]uint16, 256)
	n, _, _ := procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if n == 0 {
		return ""
	}
	return windows.UTF16ToString(buf[:n])
}

var (
	procGetDpiForWindow = user32.NewProc("GetDpiForWindow")
	procGetDpiForSystem = user32.NewProc("GetDpiForSystem")
)

func dpiFor(hwnd uintptr) int32 {
	if hwnd != 0 && procGetDpiForWindow.Find() == nil {
		if d, _, _ := procGetDpiForWindow.Call(hwnd); d >= 72 {
			return int32(d)
		}
	}
	if procGetDpiForSystem.Find() == nil {
		if d, _, _ := procGetDpiForSystem.Call(); d >= 72 {
			return int32(d)
		}
	}
	return 96
}

func scaleDPI(v, dpi int32) int32 { return v * dpi / 96 }

var (
	procMonitorFromPoint = user32.NewProc("MonitorFromPoint")
	procGetDpiForMonitor = shcore.NewProc("GetDpiForMonitor")
	procGetCursorPosDPI  = user32.NewProc("GetCursorPos")
)

func dpiForPoint(x, y int32) int32 {
	if procMonitorFromPoint.Find() == nil && procGetDpiForMonitor.Find() == nil {
		mon, _, _ := procMonitorFromPoint.Call(uintptr(uint32(x))|uintptr(uint32(y))<<32, 2)
		if mon != 0 {
			var dx, dy uint32
			r, _, _ := procGetDpiForMonitor.Call(mon, 0, uintptr(unsafe.Pointer(&dx)), uintptr(unsafe.Pointer(&dy)))
			if r == 0 && dx >= 72 {
				return int32(dx)
			}
		}
	}
	return dpiFor(0)
}

func dpiForCursor() int32 {
	var pt point
	procGetCursorPosDPI.Call(uintptr(unsafe.Pointer(&pt)))
	return dpiForPoint(pt.X, pt.Y)
}

type monitorInfo struct {
	Size    uint32
	Monitor rect
	Work    rect
	Flags   uint32
}

var procGetMonitorInfoW = user32.NewProc("GetMonitorInfoW")

func workAreaForPoint(x, y int32) rect {
	var wa rect
	if procMonitorFromPoint.Find() == nil && procGetMonitorInfoW.Find() == nil {
		mon, _, _ := procMonitorFromPoint.Call(uintptr(uint32(x))|uintptr(uint32(y))<<32, 2)
		if mon != 0 {
			mi := monitorInfo{Size: uint32(unsafe.Sizeof(monitorInfo{}))}
			if r, _, _ := procGetMonitorInfoW.Call(mon, uintptr(unsafe.Pointer(&mi))); r != 0 {
				return mi.Work
			}
		}
	}
	procSystemParametersInfoW.Call(0x30, 0, uintptr(unsafe.Pointer(&wa)), 0)
	return wa
}

var (
	procOpenProcess                = kernel32.NewProc("OpenProcess")
	procQueryFullProcessImageNameW = kernel32.NewProc("QueryFullProcessImageNameW")
	procCloseHandleW               = kernel32.NewProc("CloseHandle")
	procGetWindowThreadPID         = user32.NewProc("GetWindowThreadProcessId")
)

func processNameOf(hwnd uintptr) string {
	if hwnd == 0 || hwnd == 1 {
		return ""
	}
	var pid uint32
	procGetWindowThreadPID.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	if pid == 0 {
		return ""
	}
	h, _, _ := procOpenProcess.Call(0x1000, 0, uintptr(pid))
	if h == 0 {
		return ""
	}
	defer procCloseHandleW.Call(h)
	buf := make([]uint16, 520)
	size := uint32(len(buf))
	ok, _, _ := procQueryFullProcessImageNameW.Call(h, 0,
		uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&size)))
	if ok == 0 {
		return ""
	}
	full := windows.UTF16ToString(buf[:size])
	if i := strings.LastIndexAny(full, `\/`); i >= 0 {
		return full[i+1:]
	}
	return full
}

func ownWindow(hwnd uintptr) bool {
	if hwnd == 0 {
		return false
	}
	var pid uint32
	procGetWindowThreadPID.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	return pid != 0 && pid == uint32(windows.GetCurrentProcessId())
}

func openSettingsInRunningInstance() bool {
	cls, err := windows.UTF16PtrFromString(appid.Class("TrayWnd"))
	if err != nil {
		return false
	}
	hwnd, _, _ := procFindWindowW.Call(uintptr(unsafe.Pointer(cls)), 0)
	if hwnd == 0 {
		return false
	}
	r, _, _ := procPostMessageW.Call(hwnd, wmTrayCallback, 0, wmLButtonUp)
	return r != 0
}
