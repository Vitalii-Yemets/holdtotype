package main

import (
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	aboutOpen      atomic.Bool
	aboutClassOnce sync.Once
	aboutMu        sync.Mutex
	aboutText      string
)

func aboutWndProc(hwnd, msg, wParam, lParam uintptr) uintptr {
	switch msg {
	case wmPaint:
		var ps paintStruct
		hdc, _, _ := procBeginPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
		var rc rect
		procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
		rc.Left += 20
		rc.Top += 16
		rc.Right -= 20
		font, _, _ := procGetStockObject.Call(defaultGUIFont)
		procSelectObject.Call(hdc, font)
		procSetBkMode.Call(hdc, 1)
		procSetTextColor.Call(hdc, colGreen)
		aboutMu.Lock()
		text := aboutText
		aboutMu.Unlock()
		u, _ := windows.UTF16FromString(text)
		procDrawTextW.Call(hdc, uintptr(unsafe.Pointer(&u[0])), uintptr(len(u)-1),
			uintptr(unsafe.Pointer(&rc)), 0x0010 )
		procEndPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
		return 0
	case wmClose:
		procDestroyWindow.Call(hwnd)
		return 0
	case wmDestroy:
		procPostQuitMessage.Call(0)
		return 0
	}
	r, _, _ := procDefWindowProcW.Call(hwnd, msg, wParam, lParam)
	return r
}

func showAbout(cfg *Config) {
	if !aboutOpen.CompareAndSwap(false, true) {
		return
	}
	defer aboutOpen.Store(false)

	aboutMu.Lock()
	aboutText = trf("about.text", appVersion, cfg.Hotkey, cfg.Model, cfg.Language)
	aboutMu.Unlock()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	aboutClassOnce.Do(func() {
		className, _ := windows.UTF16PtrFromString("V2TAboutWnd")
		cursor, _, _ := procLoadCursorW.Call(0, idcArrow)
		bg, _, _ := procCreateSolidBrush.Call(colBg)
		wc := wndClassExW{
			Size:       uint32(unsafe.Sizeof(wndClassExW{})),
			WndProc:    syscall.NewCallback(aboutWndProc),
			Cursor:     cursor,
			Background: bg,
			ClassName:  className,
		}
		procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
	})

	const w, h = 480, 330
	sw, _, _ := procGetSystemMetrics.Call(0)
	sh, _, _ := procGetSystemMetrics.Call(1)
	className, _ := windows.UTF16PtrFromString("V2TAboutWnd")
	title, _ := windows.UTF16PtrFromString(tr("menu.about") + " — Vox Terminal")
	hwnd, _, _ := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(title)),
		wsCaption|wsSysMenu|wsVisible,
		(sw-w)/2, (sh-h)/2, w, h,
		0, 0, 0, 0,
	)
	if hwnd == 0 {
		return
	}
	applyDarkCaption(hwnd)
	procSetForegroundWnd.Call(hwnd)

	var m msgStruct
	for {
		r, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(r) <= 0 {
			return
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
	}
}
