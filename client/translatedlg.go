package main

import (
	"runtime"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	"holdtotype/internal/appid"
)

const (
	tdBtnW   = 66
	tdBtnH   = 32
	tdGap    = 8
	tdH      = 104
	tdTimer  = 2
	wmTdInit = 0x0400 + 20
)

var (
	tdMu        sync.Mutex
	tdLangs     []string
	tdDefault   string
	tdCountdown int
	tdResult    chan string
	tdHwnd      uintptr

	tdClassOnce sync.Once
)

func tdWidth() int32 {
	tdMu.Lock()
	n := int32(len(tdLangs))
	tdMu.Unlock()
	w := 16 + n*(tdBtnW+tdGap) - tdGap + 16 + 34
	if w < 240 {
		w = 240
	}
	return w
}

func tdBtnRect(i int32) rect {
	x := int32(16) + i*(tdBtnW+tdGap)
	return rect{Left: x, Top: 56, Right: x + tdBtnW, Bottom: 56 + tdBtnH}
}

func tdHit(x, y int32) (string, bool) {
	tdMu.Lock()
	langs := append([]string(nil), tdLangs...)
	tdMu.Unlock()
	w := tdWidth()
	if x >= w-30 && x <= w-8 && y >= 8 && y <= 30 {
		return "", true
	}
	for i, l := range langs {
		r := tdBtnRect(int32(i))
		if x >= r.Left && x <= r.Right && y >= r.Top && y <= r.Bottom {
			return l, true
		}
	}
	return "", false
}

func tdFinish(res string) {
	tdMu.Lock()
	ch := tdResult
	tdResult = nil
	tdMu.Unlock()
	if ch != nil {
		ch <- res
	}
}

func tdAbort() {
	tdMu.Lock()
	active := tdResult != nil
	tdMu.Unlock()
	if active {
		tdFinish("")
	}
}

func tdWndProc(hwnd, msg, wParam, lParam uintptr) uintptr {
	switch msg {
	case wmPaint:
		tdPaint(hwnd)
		return 0
	case wmLBtnUp:
		x := int32(int16(lParam & 0xFFFF))
		y := int32(int16(lParam >> 16 & 0xFFFF))
		if res, hit := tdHit(x, y); hit {
			procKillTimer.Call(hwnd, tdTimer)
			procShowWindow.Call(hwnd, swHide)
			tdFinish(res)
		}
		return 0
	case wmSetCursor:
		var pt point
		procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
		procScreenToClient.Call(hwnd, uintptr(unsafe.Pointer(&pt)))
		if _, hit := tdHit(pt.X, pt.Y); hit {
			cur, _, _ := procLoadCursorW.Call(0, 32649)
			procSetCursor.Call(cur)
			return 1
		}
	case wmTimer:
		tdMu.Lock()
		if tdCountdown > 0 {
			tdCountdown--
		}
		done := tdCountdown == 0
		def := tdDefault
		tdMu.Unlock()
		if done {
			procKillTimer.Call(hwnd, tdTimer)
			procShowWindow.Call(hwnd, swHide)
			tdFinish(def)
		} else {
			procInvalidateRect.Call(hwnd, 0, 0)
			tdPaintDirect(hwnd)
		}
		return 0
	case wmDestroy:
		procPostQuitMessage.Call(0)
		return 0
	}
	r, _, _ := procDefWindowProcW.Call(hwnd, msg, wParam, lParam)
	return r
}

func tdPaintDirect(hwnd uintptr) {
	hdc, _, _ := procGetDC.Call(hwnd)
	if hdc == 0 {
		return
	}
	tdRender(hwnd, hdc)
	procReleaseDC.Call(hwnd, hdc)
}

func tdPaint(hwnd uintptr) {
	var ps paintStruct
	hdc, _, _ := procBeginPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
	tdRender(hwnd, hdc)
	procEndPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
}

func tdRender(hwnd, hdc uintptr) {
	var rc rect
	procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
	tdMu.Lock()
	langs := append([]string(nil), tdLangs...)
	def := tdDefault
	cd := tdCountdown
	tdMu.Unlock()

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

	title := tr("td.title")
	if cd > 0 {
		title += " (" + itoa(cd) + ")"
	}
	drawText := func(s string, r rect, color uintptr, flags uintptr) {
		u, _ := windows.UTF16FromString(s)
		procSetTextColor.Call(hdc, color)
		procDrawTextW.Call(hdc, uintptr(unsafe.Pointer(&u[0])), uintptr(len(u)-1),
			uintptr(unsafe.Pointer(&r)), flags)
	}
	drawText(title, rect{Left: 16, Top: 14, Right: rc.Right - 40, Bottom: 46}, colGreen, 0)
	drawText("✕", rect{Left: rc.Right - 30, Top: 8, Right: rc.Right - 8, Bottom: 30}, colGreenDm, 0x0020|0x0004|0x0001)

	for i, l := range langs {
		r := tdBtnRect(int32(i))
		var border uintptr = colGreenLo
		var txt uintptr = colGreenDm
		if l == def {
			border = colGreenDm
			txt = colGreen
		}
		fill(rect{Left: r.Left - 1, Top: r.Top - 1, Right: r.Right + 1, Bottom: r.Bottom + 1}, border)
		fill(r, 0x0B100D)
		drawText(langLabel(l), r, txt, 0x0020|0x0004|0x0001)
	}
}

func langLabel(code string) string {
	switch code {
	case "en":
		return "EN"
	case "de":
		return "DE"
	case "fr":
		return "FR"
	case "es":
		return "ES"
	case "it":
		return "IT"
	case "pl":
		return "PL"
	case "ru":
		return "RU"
	case "uk":
		return "UK"
	}
	return code
}

func askTranslateTarget(cfg *Config) string {
	ch := make(chan string, 1)
	tdMu.Lock()
	if tdResult != nil {
		tdMu.Unlock()
		return cfg.TranslateTarget
	}
	tdLangs = append([]string(nil), cfg.TranslateAskLangs...)
	tdDefault = cfg.TranslateTarget
	if cfg.TranslateAsk == "timeout" {
		tdCountdown = cfg.TranslateAskSeconds
	} else {
		tdCountdown = 0
	}
	tdResult = ch
	timeoutMode := cfg.TranslateAsk == "timeout"
	tdMu.Unlock()

	go func() {
		runtime.LockOSThread()
		tdClassOnce.Do(func() {
			className, _ := windows.UTF16PtrFromString(appid.Class("TranslateDlg"))
			cb := syscall.NewCallback(tdWndProc)
			wc := wndClassExW{
				Size:      uint32(unsafe.Sizeof(wndClassExW{})),
				WndProc:   cb,
				ClassName: className,
			}
			procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
		})
		className, _ := windows.UTF16PtrFromString(appid.Class("TranslateDlg"))
		w := tdWidth()
		var wa rect
		procSystemParametersInfoW.Call(0x30, 0, uintptr(unsafe.Pointer(&wa)), 0)
		x := int(wa.Left) + (int(wa.Right-wa.Left)-int(w))/2
		y := int(wa.Bottom) - tdH - ovH - 44

		hwnd, _, _ := procCreateWindowExW.Call(
			wsExLayered|wsExToolWindow|wsExNoActivate|0x00000008,
			uintptr(unsafe.Pointer(className)), 0,
			wsPopup,
			uintptr(x), uintptr(y), uintptr(w), tdH,
			0, 0, 0, 0,
		)
		if hwnd == 0 {
			tdFinish(cfg.TranslateTarget)
			return
		}
		applyDarkCaption(hwnd)
		procSetLayeredWindowAttributes.Call(hwnd, 0, 245, lwaAlpha)
		tdMu.Lock()
		tdHwnd = hwnd
		tdMu.Unlock()
		procShowWindow.Call(hwnd, swShowNoActivate)
		tdPaintDirect(hwnd)
		if timeoutMode {
			procSetTimer.Call(hwnd, tdTimer, 1000, 0)
		}
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
	tdMu.Lock()
	hwnd := tdHwnd
	tdHwnd = 0
	tdMu.Unlock()
	if hwnd != 0 {
		procPostMessageW.Call(hwnd, wmClose, 0, 0)
	}
	return res
}
