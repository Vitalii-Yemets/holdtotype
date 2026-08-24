package main

import (
	"sync/atomic"
	"runtime"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	"holdtotype/internal/appid"
)

const (
	tmWidth          = 260
	tmItemH          = 30
	tmSepH           = 9
	tmPad            = 6
	tmTextPad        = 18
	wmLBtnDown       = 0x0201
	wmRBtnDown       = 0x0204
	wmCaptureChanged = 0x0215
)

type tmItem struct {
	text   string
	id     uintptr
	grayed bool
	sep    bool
}

var (
	tmMu    sync.Mutex
	tmItems []tmItem
	tmHover int = -1
	tmHwnd  uintptr
	tmDone  chan uintptr

	tmClassOnce sync.Once

	procSetCapture = user32.NewProc("SetCapture")
)

var tmDpiVal atomic.Int32

func tmDPI() int32 {
	if v := tmDpiVal.Load(); v >= 72 {
		return v
	}
	return 96
}

var tmWidthDIP atomic.Int32

func tmMeasure(items []tmItem) int32 {
	widest := int32(tmWidth)
	for _, it := range items {
		if it.sep || it.text == "" {
			continue
		}
		if w := textWidthDIP(it.text) + 2*tmTextPad; w > widest {
			widest = w
		}
	}
	max := askMaxWidthDIP()
	if widest > max {
		widest = max
	}
	tmWidthDIP.Store(widest)
	return widest
}

func tmW() int32 {
	w := tmWidthDIP.Load()
	if w <= 0 {
		w = tmWidth
	}
	return scaleDPI(w, tmDPI())
}

func tmItemHeight() int32 { return scaleDPI(tmItemH, tmDPI()) }

func tmSepHeight() int32 { return scaleDPI(tmSepH, tmDPI()) }

func tmPadding() int32 { return scaleDPI(tmPad, tmDPI()) }

func tmHeight() int32 {
	h := tmPadding() * 2
	itemH, sepH := tmItemHeight(), tmSepHeight()
	tmMu.Lock()
	for _, it := range tmItems {
		if it.sep {
			h += sepH
		} else {
			h += itemH
		}
	}
	tmMu.Unlock()
	return h
}

func tmItemAt(y int32) int {
	cur := tmPadding()
	itemH, sepH := tmItemHeight(), tmSepHeight()
	tmMu.Lock()
	defer tmMu.Unlock()
	for i, it := range tmItems {
		hh := itemH
		if it.sep {
			hh = sepH
		}
		if y >= cur && y < cur+hh {
			if it.sep || it.grayed {
				return -1
			}
			return i
		}
		cur += hh
	}
	return -1
}

func tmFinish(id uintptr) {
	tmMu.Lock()
	ch := tmDone
	tmDone = nil
	tmMu.Unlock()
	if ch != nil {
		ch <- id
	}
}

func tmWndProc(hwnd, msg, wParam, lParam uintptr) uintptr {
	switch msg {
	case wmPaint:
		var ps paintStruct
		hdc, _, _ := procBeginPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
		tmRender(hwnd, hdc)
		procEndPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
		return 0
	case wmMouseMove:
		y := int32(int16(lParam >> 16 & 0xFFFF))
		x := int32(int16(lParam & 0xFFFF))
		idx := -1
		if x >= 0 && x < tmW() {
			idx = tmItemAt(y)
		}
		tmMu.Lock()
		changed := idx != tmHover
		tmHover = idx
		tmMu.Unlock()
		if changed {
			hdc, _, _ := procGetDC.Call(hwnd)
			if hdc != 0 {
				tmRender(hwnd, hdc)
				procReleaseDC.Call(hwnd, hdc)
			}
		}
		return 0
	case wmSetCursor:
		tmMu.Lock()
		hov := tmHover
		tmMu.Unlock()
		if hov >= 0 {
			cur, _, _ := procLoadCursorW.Call(0, 32649)
			procSetCursor.Call(cur)
			return 1
		}
	case wmLBtnDown, wmRBtnDown:
		x := int32(int16(lParam & 0xFFFF))
		y := int32(int16(lParam >> 16 & 0xFFFF))
		if x < 0 || x >= tmW() || y < 0 || y >= tmHeight() {
			procShowWindow.Call(hwnd, swHide)
			tmFinish(0)
		}
		return 0
	case wmLBtnUp:
		y := int32(int16(lParam >> 16 & 0xFFFF))
		if idx := tmItemAt(y); idx >= 0 {
			tmMu.Lock()
			id := tmItems[idx].id
			tmMu.Unlock()
			procShowWindow.Call(hwnd, swHide)
			tmFinish(id)
		}
		return 0
	case wmCaptureChanged:
		tmMu.Lock()
		open := tmDone != nil
		tmMu.Unlock()
		if open {
			procShowWindow.Call(hwnd, swHide)
			tmFinish(0)
		}
		return 0
	case wmDestroy:
		procPostQuitMessage.Call(0)
		return 0
	}
	r, _, _ := procDefWindowProcW.Call(hwnd, msg, wParam, lParam)
	return r
}

func tmRender(hwnd, hdc uintptr) {
	var rc rect
	procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
	fill := func(r rect, color uintptr) {
		br, _, _ := procCreateSolidBrush.Call(color)
		procFillRect.Call(hdc, uintptr(unsafe.Pointer(&r)), br)
		procDeleteObject.Call(br)
	}
	fill(rc, colBg)
	for y := rc.Top + 3; y < rc.Bottom; y += 3 {
		fill(rect{Left: rc.Left, Top: y, Right: rc.Right, Bottom: y + 1}, colBgLine)
	}

	procSelectObject.Call(hdc, uiFont(hwnd))
	procSetBkMode.Call(hdc, 1)

	tmMu.Lock()
	items := append([]tmItem(nil), tmItems...)
	hov := tmHover
	tmMu.Unlock()

	dpi := dpiFor(hwnd)
	px := func(v int32) int32 { return scaleDPI(v, dpi) }
	itemH, sepH := px(tmItemH), px(tmSepH)
	cur := px(tmPad)
	for i, it := range items {
		if it.sep {
			mid := cur + sepH/2
			fill(rect{Left: px(10), Top: mid, Right: rc.Right - px(10), Bottom: mid + 1}, colGreenLo)
			cur += sepH
			continue
		}
		row := rect{Left: px(4), Top: cur, Right: rc.Right - px(4), Bottom: cur + itemH}
		if i == hov {
			fill(row, 0x223F12)
		}
		var color uintptr = colGreen
		if it.grayed {
			color = colGreenDm
		}
		procSetTextColor.Call(hdc, color)
		txt := rect{Left: px(16), Top: cur, Right: rc.Right - px(10), Bottom: cur + itemH}
		u, _ := windows.UTF16FromString(it.text)
		procDrawTextW.Call(hdc, uintptr(unsafe.Pointer(&u[0])), uintptr(len(u)-1),
			uintptr(unsafe.Pointer(&txt)), 0x0020|0x0004|0x8000)
		cur += itemH
	}
}

func showTrayMenu(items []tmItem) uintptr {
	ch := make(chan uintptr, 1)
	tmMeasure(items)
	tmMu.Lock()
	if tmDone != nil {
		tmMu.Unlock()
		return 0
	}
	tmItems = items
	tmHover = -1
	tmDone = ch
	tmMu.Unlock()

	go func() {
		runtime.LockOSThread()
		tmClassOnce.Do(func() {
			className, _ := windows.UTF16PtrFromString(appid.Class("TrayMenu"))
			cb := syscall.NewCallback(tmWndProc)
			wc := wndClassExW{
				Size:      uint32(unsafe.Sizeof(wndClassExW{})),
				WndProc:   cb,
				ClassName: className,
			}
			procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
		})
		className, _ := windows.UTF16PtrFromString(appid.Class("TrayMenu"))
		var pt point
		procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
		tmDpiVal.Store(dpiForPoint(pt.X, pt.Y))
		h := tmHeight()
		wa := workAreaForPoint(pt.X, pt.Y)
		x := pt.X
		mw := tmW()
		if x+mw > wa.Right {
			x = wa.Right - mw - 4
		}
		y := pt.Y - h
		if y < wa.Top {
			y = wa.Top + 4
		}
		hwnd, _, _ := procCreateWindowExW.Call(
			wsExLayered|wsExToolWindow|wsExNoActivate|0x00000008,
			uintptr(unsafe.Pointer(className)), 0,
			wsPopup,
			uintptr(x), uintptr(y), uintptr(mw), uintptr(h),
			0, 0, 0, 0,
		)
		if hwnd == 0 {
			tmFinish(0)
			return
		}
		applyDarkCaption(hwnd)
		procSetLayeredWindowAttributes.Call(hwnd, 0, 247, lwaAlpha)
		tmMu.Lock()
		tmHwnd = hwnd
		tmMu.Unlock()
		procShowWindow.Call(hwnd, swShowNoActivate)
		hdc, _, _ := procGetDC.Call(hwnd)
		if hdc != 0 {
			tmRender(hwnd, hdc)
			procReleaseDC.Call(hwnd, hdc)
		}
		procSetCapture.Call(hwnd)
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

	id := <-ch
	tmMu.Lock()
	hwnd := tmHwnd
	tmHwnd = 0
	tmMu.Unlock()
	if hwnd != 0 {
		procPostMessageW.Call(hwnd, wmClose, 0, 0)
	}
	return id
}
