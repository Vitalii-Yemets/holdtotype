package main

import (
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"holdtotype/internal/theme"
)

var (
	user32   = windows.NewLazySystemDLL("user32.dll")
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	dwmapi   = windows.NewLazySystemDLL("dwmapi.dll")
	shell32  = windows.NewLazySystemDLL("shell32.dll")

	procMessageBoxW           = user32.NewProc("MessageBoxW")
	procGetWindowLongPtrW     = user32.NewProc("GetWindowLongPtrW")
	procSetWindowLongPtrW     = user32.NewProc("SetWindowLongPtrW")
	procSetWindowPos          = user32.NewProc("SetWindowPos")
	procReleaseCapture        = user32.NewProc("ReleaseCapture")
	procSendMessageW          = user32.NewProc("SendMessageW")
	procShowWindow            = user32.NewProc("ShowWindow")
	procSetForegroundWnd      = user32.NewProc("SetForegroundWindow")
	procDwmSetWindowAttribute = dwmapi.NewProc("DwmSetWindowAttribute")
	procShellExecuteW         = shell32.NewProc("ShellExecuteW")
	procGetUserDefaultUILang  = kernel32.NewProc("GetUserDefaultUILanguage")

	gdi32                     = windows.NewLazySystemDLL("gdi32.dll")
	procCreateSolidBrush      = gdi32.NewProc("CreateSolidBrush")
	procSetClassLongPtrW      = user32.NewProc("SetClassLongPtrW")
	procFindWindowW           = user32.NewProc("FindWindowW")
	procGetWindowRect         = user32.NewProc("GetWindowRect")
	procSystemParametersInfoW = user32.NewProc("SystemParametersInfoW")

	ole32                    = windows.NewLazySystemDLL("ole32.dll")
	procSHBrowseForFolderW   = shell32.NewProc("SHBrowseForFolderW")
	procSHGetPathFromIDListW = shell32.NewProc("SHGetPathFromIDListW")
	procCoTaskMemFree        = ole32.NewProc("CoTaskMemFree")

	shcore                = windows.NewLazySystemDLL("shcore.dll")
	procGetDpiForMonitor  = shcore.NewProc("GetDpiForMonitor")
	procGetDpiForWindow   = user32.NewProc("GetDpiForWindow")
	procCallWindowProcW   = user32.NewProc("CallWindowProcW")
	procGetCursorPos      = user32.NewProc("GetCursorPos")
	procMonitorFromPoint  = user32.NewProc("MonitorFromPoint")
	procMonitorFromWindow = user32.NewProc("MonitorFromWindow")
	procMonitorFromRect   = user32.NewProc("MonitorFromRect")
	procGetMonitorInfoW   = user32.NewProc("GetMonitorInfoW")
)

type browseInfoW struct {
	Owner       uintptr
	Root        uintptr
	DisplayName *uint16
	Title       *uint16
	Flags       uint32
	Callback    uintptr
	LParam      uintptr
	Image       int32
}

func browseFolder(owner uintptr, title string) string {
	display := make([]uint16, 260)
	t, _ := windows.UTF16PtrFromString(title)
	bi := browseInfoW{
		Owner:       owner,
		DisplayName: &display[0],
		Title:       t,
		Flags:       0x00000041,
	}
	pidl, _, _ := procSHBrowseForFolderW.Call(uintptr(unsafe.Pointer(&bi)))
	if pidl == 0 {
		return ""
	}
	defer procCoTaskMemFree.Call(pidl)
	buf := make([]uint16, 300)
	if r, _, _ := procSHGetPathFromIDListW.Call(pidl, uintptr(unsafe.Pointer(&buf[0]))); r == 0 {
		return ""
	}
	return windows.UTF16ToString(buf)
}

func msgBox(title, text string) {
	t, _ := windows.UTF16PtrFromString(title)
	m, _ := windows.UTF16PtrFromString(text)
	procMessageBoxW.Call(0, uintptr(unsafe.Pointer(m)), uintptr(unsafe.Pointer(t)), 0x10)
}

func shellOpenURL(url string) {
	verb, _ := windows.UTF16PtrFromString("open")
	p, _ := windows.UTF16PtrFromString(url)
	procShellExecuteW.Call(0, uintptr(unsafe.Pointer(verb)), uintptr(unsafe.Pointer(p)), 0, 0, 1)
}

func colorref(hex string) int32 {
	r, g, b := theme.RGB(hex)
	return int32(b)<<16 | int32(g)<<8 | int32(r)
}

func applyCaption(hwnd uintptr, l theme.Look) {
	p := l.Palette
	dark := int32(1)
	if p.Light() {
		dark = 0
	}
	procDwmSetWindowAttribute.Call(hwnd, 20, uintptr(unsafe.Pointer(&dark)), 4)
	capColor := colorref(p.Bg)
	procDwmSetWindowAttribute.Call(hwnd, 35, uintptr(unsafe.Pointer(&capColor)), 4)
	txtColor := colorref(p.Text)
	procDwmSetWindowAttribute.Call(hwnd, 36, uintptr(unsafe.Pointer(&txtColor)), 4)
	border := int32(-2)
	corner := int32(1)
	if l.Round {
		border = colorref(p.Line)
		corner = 2
		if l.SmallR {
			corner = 3
		}
	}
	procDwmSetWindowAttribute.Call(hwnd, 34, uintptr(unsafe.Pointer(&border)), 4)
	procDwmSetWindowAttribute.Call(hwnd, 33, uintptr(unsafe.Pointer(&corner)), 4)
}

func setClientBackground(hwnd uintptr, p theme.Palette) {
	br, _, _ := procCreateSolidBrush.Call(uintptr(uint32(colorref(p.Bg))))
	procSetClassLongPtrW.Call(hwnd, ^uintptr(9), br)
}

const offscreenXY int32 = -32000

func offscreenPos() uintptr {
	v := offscreenXY
	return uintptr(uint32(v))
}

type winRect struct{ Left, Top, Right, Bottom int32 }

func hideWebViewWindowEarly(title string, l theme.Look) func() {
	done := make(chan struct{})
	go func() {
		cls, _ := windows.UTF16PtrFromString("webview")
		t, _ := windows.UTF16PtrFromString(title)
		for {
			select {
			case <-done:
				return
			default:
			}
			h, _, _ := procFindWindowW.Call(uintptr(unsafe.Pointer(cls)), uintptr(unsafe.Pointer(t)))
			if h != 0 {
				var rc winRect
				procGetWindowRect.Call(h, uintptr(unsafe.Pointer(&rc)))
				if rc.Left > offscreenXY+1000 {
					procSetWindowPos.Call(h, 0, offscreenPos(), offscreenPos(), 0, 0, 0x0001|0x0004|0x0010)
				}
				setClientBackground(h, l.Palette)
				applyCaption(h, l)
			}
			time.Sleep(2 * time.Millisecond)
		}
	}()
	return func() { close(done) }
}

const winDIPW int32 = 560

var (
	winDIPH    int32 = 560
	fitOldProc uintptr
	fitProcCB  uintptr
)

type monInfo struct {
	Size    uint32
	Monitor winRect
	Work    winRect
	Flags   uint32
}

func scaleDIP(v, dpi int32) int32 { return (v*dpi + 95) / 96 }

func windowDPI(hwnd uintptr) int32 {
	if procGetDpiForWindow.Find() == nil {
		if d, _, _ := procGetDpiForWindow.Call(hwnd); d >= 72 {
			return int32(d)
		}
	}
	return 96
}

func monitorDPI(mon uintptr) int32 {
	if mon != 0 && procGetDpiForMonitor.Find() == nil {
		var dx, dy uint32
		if r, _, _ := procGetDpiForMonitor.Call(mon, 0, uintptr(unsafe.Pointer(&dx)), uintptr(unsafe.Pointer(&dy))); r == 0 && dx >= 72 {
			return int32(dx)
		}
	}
	return 96
}

func workArea(mon uintptr) winRect {
	if mon != 0 {
		mi := monInfo{Size: uint32(unsafe.Sizeof(monInfo{}))}
		if r, _, _ := procGetMonitorInfoW.Call(mon, uintptr(unsafe.Pointer(&mi))); r != 0 {
			return mi.Work
		}
	}
	var wa winRect
	procSystemParametersInfoW.Call(0x30, 0, uintptr(unsafe.Pointer(&wa)), 0)
	return wa
}

func cursorMonitor() uintptr {
	var pt struct{ X, Y int32 }
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	mon, _, _ := procMonitorFromPoint.Call(uintptr(uint32(pt.X))|uintptr(uint32(pt.Y))<<32, 2)
	return mon
}

func windowSize(dpi int32, work winRect) (int32, int32) {
	w, h := scaleDIP(winDIPW, dpi), scaleDIP(winDIPH, dpi)
	if mw := work.Right - work.Left; w > mw {
		w = mw
	}
	if mh := work.Bottom - work.Top; h > mh {
		h = mh
	}
	return w, h
}

func setWindowRect(hwnd uintptr, x, y, w, h int32) {
	procSetWindowPos.Call(hwnd, 0, uintptr(uint32(x)), uintptr(uint32(y)), uintptr(w), uintptr(h), 0x0004|0x0010)
}

func sizeForDPI(hwnd uintptr) {
	dpi := windowDPI(hwnd)
	procSetWindowPos.Call(hwnd, 0, 0, 0, uintptr(scaleDIP(winDIPW, dpi)), uintptr(scaleDIP(winDIPH, dpi)), 0x0002|0x0004|0x0010)
}

func followDPI(hwnd uintptr) {
	fitProcCB = windows.NewCallback(fitProc)
	const gwlpWndproc = ^uintptr(3)
	fitOldProc, _, _ = procSetWindowLongPtrW.Call(hwnd, gwlpWndproc, fitProcCB)
}

func fitProc(hwnd, msg, wp, lp uintptr) uintptr {
	if msg == 0x02E0 && lp != 0 {
		s := (*winRect)(unsafe.Pointer(lp))
		mon, _, _ := procMonitorFromRect.Call(lp, 2)
		w, h := windowSize(int32(wp&0xFFFF), workArea(mon))
		setWindowRect(hwnd, s.Left, s.Top, w, h)
		return 0
	}
	r, _, _ := procCallWindowProcW.Call(fitOldProc, hwnd, msg, wp, lp)
	return r
}

func revealWindowCentered(hwnd uintptr) {
	mon := cursorMonitor()
	work := workArea(mon)
	for _, dpi := range []int32{monitorDPI(mon), 0} {
		if dpi == 0 {
			dpi = windowDPI(hwnd)
		}
		w, h := windowSize(dpi, work)
		setWindowRect(hwnd, work.Left+(work.Right-work.Left-w)/2, work.Top+(work.Bottom-work.Top-h)/2, w, h)
	}
	procShowWindow.Call(hwnd, 5)
	procSetForegroundWnd.Call(hwnd)
}

func fitWindow(hwnd uintptr) {
	var rc winRect
	procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
	mon, _, _ := procMonitorFromWindow.Call(hwnd, 2)
	work := workArea(mon)
	w, h := windowSize(windowDPI(hwnd), work)
	y := rc.Top
	if y+h > work.Bottom {
		y = work.Bottom - h
	}
	if y < work.Top {
		y = work.Top
	}
	setWindowRect(hwnd, rc.Left, y, w, h)
}

func makeBorderless(hwnd uintptr) {
	const gwlStyle = ^uintptr(15)
	style, _, _ := procGetWindowLongPtrW.Call(hwnd, gwlStyle)
	style &^= 0x00C00000 | 0x00040000
	procSetWindowLongPtrW.Call(hwnd, gwlStyle, style)
	procSetWindowPos.Call(hwnd, 0, 0, 0, 0, 0, 0x0001|0x0002|0x0004|0x0020)
}

func beginWindowDrag(hwnd uintptr) {
	procReleaseCapture.Call()
	procSendMessageW.Call(hwnd, 0x00A1, 2, 0)
}

func uiLang() string {
	id, _, _ := procGetUserDefaultUILang.Call()
	switch id & 0x3FF {
	case 0x19:
		return "ru"
	case 0x22:
		return "uk"
	case 0x07:
		return "de"
	case 0x0C:
		return "fr"
	case 0x0A:
		return "es"
	case 0x10:
		return "it"
	case 0x15:
		return "pl"
	}
	return "en"
}

func runHidden(name string, args ...string) error {
	return makeCmd(name, args...).Run()
}
