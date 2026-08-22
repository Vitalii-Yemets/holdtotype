package main

import (
	"path/filepath"
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
	set(35, 0x0C0F0B)
	set(36, 0x6EFF3C)
	set(34, 0x2B4A1D)
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

	minSizeOld         uintptr
	minSizeW, minSizeH int32
	minSizeOnce        sync.Once
	minSizeCB          uintptr
)

var (
	procGetWindowRect  = user32.NewProc("GetWindowRect")
	lastWndW, lastWndH int32
)

func minSizeProc(hwnd, msg, wp, lp uintptr) uintptr {
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
	minSizeW, minSizeH = w, h
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
