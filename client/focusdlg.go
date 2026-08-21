package main

import (
	"runtime"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	fdBtnH    = 32
	fdBtn1W   = 160
	fdBtn2W   = 130
	fdGapX    = 8
	fdHeight  = 104
	fdTimerID = 3
)

var (
	fdMu        sync.Mutex
	fdCountdown int
	fdResult    chan string
	fdHwnd      uintptr

	fdClassOnce sync.Once
)

func fdWidth() int32 {
	w := int32(16 + fdBtn1W + fdGapX + fdBtn2W + 16 + 34)
	if w < 300 {
		w = 300
	}
	return w
}

func fdBtnRect(i int32) rect {
	x := int32(16)
	w := int32(fdBtn1W)
	if i == 1 {
		x += fdBtn1W + fdGapX
		w = fdBtn2W
	}
	return rect{Left: x, Top: 56, Right: x + w, Bottom: 56 + fdBtnH}
}

func fdHit(x, y int32) (string, bool) {
	w := fdWidth()
	if x >= w-30 && x <= w-8 && y >= 8 && y <= 30 {
		return "", true
	}
	for i, id := range []string{"here", "copy"} {
		r := fdBtnRect(int32(i))
		if x >= r.Left && x <= r.Right && y >= r.Top && y <= r.Bottom {
			return id, true
		}
	}
	return "", false
}

func fdFinish(res string) {
	fdMu.Lock()
	ch := fdResult
	fdResult = nil
	fdMu.Unlock()
	if ch != nil {
		ch <- res
	}
}

func fdAbort() {
	fdMu.Lock()
	active := fdResult != nil
	fdMu.Unlock()
	if active {
		fdFinish("")
	}
}

func fdWndProc(hwnd, msg, wParam, lParam uintptr) uintptr {
	switch msg {
	case wmPaint:
		fdPaint(hwnd)
		return 0
	case wmLBtnUp:
		x := int32(int16(lParam & 0xFFFF))
		y := int32(int16(lParam >> 16 & 0xFFFF))
		if res, hit := fdHit(x, y); hit {
			procKillTimer.Call(hwnd, fdTimerID)
			procShowWindow.Call(hwnd, swHide)
			fdFinish(res)
		}
		return 0
	case wmSetCursor:
		var pt point
		procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
		procScreenToClient.Call(hwnd, uintptr(unsafe.Pointer(&pt)))
		if _, hit := fdHit(pt.X, pt.Y); hit {
			cur, _, _ := procLoadCursorW.Call(0, 32649)
			procSetCursor.Call(cur)
			return 1
		}
	case wmTimer:
		fdMu.Lock()
		if fdCountdown > 0 {
			fdCountdown--
		}
		done := fdCountdown == 0
		fdMu.Unlock()
		if done {
			procKillTimer.Call(hwnd, fdTimerID)
			procShowWindow.Call(hwnd, swHide)
			fdFinish("")
		} else {
			procInvalidateRect.Call(hwnd, 0, 0)
			fdPaintDirect(hwnd)
		}
		return 0
	case wmDestroy:
		procPostQuitMessage.Call(0)
		return 0
	}
	r, _, _ := procDefWindowProcW.Call(hwnd, msg, wParam, lParam)
	return r
}

func fdPaintDirect(hwnd uintptr) {
	hdc, _, _ := procGetDC.Call(hwnd)
	if hdc == 0 {
		return
	}
	fdRender(hwnd, hdc)
	procReleaseDC.Call(hwnd, hdc)
}

func fdPaint(hwnd uintptr) {
	var ps paintStruct
	hdc, _, _ := procBeginPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
	fdRender(hwnd, hdc)
	procEndPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
}

func fdRender(hwnd, hdc uintptr) {
	var rc rect
	procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
	fdMu.Lock()
	cd := fdCountdown
	fdMu.Unlock()

	fill := func(r rect, color uintptr) {
		br, _, _ := procCreateSolidBrush.Call(color)
		procFillRect.Call(hdc, uintptr(unsafe.Pointer(&r)), br)
		procDeleteObject.Call(br)
	}
	fill(rc, colBg)
	for y := rc.Top + 3; y < rc.Bottom; y += 3 {
		fill(rect{Left: rc.Left, Top: y, Right: rc.Right, Bottom: y + 1}, colBgLine)
	}

	if ovFont == 0 {
		face, _ := windows.UTF16PtrFromString("Consolas")
		ovFont, _, _ = procCreateFontW.Call(uintptr(^uintptr(15)+1), 0, 0, 0, 400,
			0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(face)))
	}
	procSelectObject.Call(hdc, ovFont)
	procSetBkMode.Call(hdc, 1)

	title := tr("fd.title")
	if cd > 0 {
		title += " (" + itoa(cd) + ")"
	}
	drawText := func(s string, r rect, color uintptr, flags uintptr) {
		u, _ := windows.UTF16FromString(s)
		procSetTextColor.Call(hdc, color)
		procDrawTextW.Call(hdc, uintptr(unsafe.Pointer(&u[0])), uintptr(len(u)-1),
			uintptr(unsafe.Pointer(&r)), flags)
	}
	drawText(title, rect{Left: 16, Top: 14, Right: rc.Right - 40, Bottom: 46}, colAmber, 0)
	drawText("✕", rect{Left: rc.Right - 30, Top: 8, Right: rc.Right - 8, Bottom: 30}, colGreenDm, 0x0020|0x0004|0x0001)

	labels := []string{tr("fd.here"), tr("fd.copy")}
	for i, l := range labels {
		r := fdBtnRect(int32(i))
		border := uintptr(colGreenLo)
		txt := uintptr(colGreenDm)
		if i == 0 {
			border = colGreenDm
			txt = colGreen
		}
		fill(rect{Left: r.Left - 1, Top: r.Top - 1, Right: r.Right + 1, Bottom: r.Bottom + 1}, border)
		fill(r, 0x0B100D)
		drawText(l, r, txt, 0x0020|0x0004|0x0001)
	}
}

func askFocusMismatch() string {
	ch := make(chan string, 1)
	fdMu.Lock()
	if fdResult != nil {
		fdMu.Unlock()
		return ""
	}
	fdCountdown = 30
	fdResult = ch
	fdMu.Unlock()

	go func() {
		runtime.LockOSThread()
		fdClassOnce.Do(func() {
			className, _ := windows.UTF16PtrFromString("V2TFocusDlg")
			cb := syscall.NewCallback(fdWndProc)
			wc := wndClassExW{
				Size:      uint32(unsafe.Sizeof(wndClassExW{})),
				WndProc:   cb,
				ClassName: className,
			}
			procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
		})
		className, _ := windows.UTF16PtrFromString("V2TFocusDlg")
		w := fdWidth()
		var wa rect
		procSystemParametersInfoW.Call(0x30, 0, uintptr(unsafe.Pointer(&wa)), 0)
		x := int(wa.Left) + (int(wa.Right-wa.Left)-int(w))/2
		y := int(wa.Bottom) - fdHeight - ovH - 44

		hwnd, _, _ := procCreateWindowExW.Call(
			wsExLayered|wsExToolWindow|wsExNoActivate|0x00000008,
			uintptr(unsafe.Pointer(className)), 0,
			wsPopup,
			uintptr(x), uintptr(y), uintptr(w), fdHeight,
			0, 0, 0, 0,
		)
		if hwnd == 0 {
			fdFinish("")
			return
		}
		applyDarkCaption(hwnd)
		procSetLayeredWindowAttributes.Call(hwnd, 0, 245, lwaAlpha)
		fdMu.Lock()
		fdHwnd = hwnd
		fdMu.Unlock()
		procShowWindow.Call(hwnd, swShowNoActivate)
		fdPaintDirect(hwnd)
		procSetTimer.Call(hwnd, fdTimerID, 1000, 0)
		var m msgStruct
		for {
			r, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
			if int32(r) <= 0 {
				break
			}
			procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
			procDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
		}
	}()

	res := <-ch
	fdMu.Lock()
	hwnd := fdHwnd
	fdHwnd = 0
	fdMu.Unlock()
	if hwnd != 0 {
		procPostMessageW.Call(hwnd, wmClose, 0, 0)
	}
	return res
}
