package main

import (
	"log"
	"math"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"holdtotype/internal/appid"
)

const (
	ovHidden = iota
	ovRecording
	ovProcessing
	ovFlashOK
	ovFlashErr
)

const (
	wmOvSet      = 0x0400 + 10
	wmTimer      = 0x0113
	wmMouseMove  = 0x0200
	wmLBtnUp     = 0x0202
	wmSetCursor  = 0x0020
	ovW          = 390
	ovH          = 52
	ovTimerID    = 1

	wsPopup          = 0x80000000
	wsExLayered      = 0x00080000
	wsExToolWindow   = 0x00000080
	wsExNoActivate   = 0x08000000
	swShowNoActivate = 4
	swHide           = 0
	nullPen          = 8
	lwaAlpha         = 2
)

var (
	procSetTimer                   = user32.NewProc("SetTimer")
	procKillTimer                  = user32.NewProc("KillTimer")
	procShowWindow                 = user32.NewProc("ShowWindow")
	procFillRect                   = user32.NewProc("FillRect")
	procSystemParametersInfoW      = user32.NewProc("SystemParametersInfoW")
	procSetWindowRgn               = user32.NewProc("SetWindowRgn")
	procSetLayeredWindowAttributes = user32.NewProc("SetLayeredWindowAttributes")
	procScreenToClient             = user32.NewProc("ScreenToClient")
	procSetCursor                  = user32.NewProc("SetCursor")
	procCreateSolidBrush           = gdi32.NewProc("CreateSolidBrush")
	procDeleteObject               = gdi32.NewProc("DeleteObject")
	procEllipse                    = gdi32.NewProc("Ellipse")
	procCreateRoundRectRgn         = gdi32.NewProc("CreateRoundRectRgn")
	procSetTextColor               = gdi32.NewProc("SetTextColor")
	procCreateFontW                = gdi32.NewProc("CreateFontW")
	procFrameRgn                   = gdi32.NewProc("FrameRgn")
	procGetDC                      = user32.NewProc("GetDC")
	procReleaseDC                  = user32.NewProc("ReleaseDC")
)

func blendCol(a, b uintptr, t float64) uintptr {
	ch := func(v uintptr, sh uint) float64 { return float64((v >> sh) & 0xFF) }
	mix := func(sh uint) uintptr {
		v := uintptr(ch(a, sh)*(1-t) + ch(b, sh)*t)
		return (v & 0xFF) << sh
	}
	return mix(0) | mix(8) | mix(16)
}

const (
	colBg      = 0x0C0F0B
	colBgLine  = 0x0A0C09
	colGreen   = 0x6EFF3C
	colGreenDm = 0x4AA320
	colGreenLo = 0x2B4A1D
	colRed     = 0x6B6BFF
	colRedDm   = 0x26265C
	colAmber   = 0x47B3FF
	colAmberDm = 0x20455C
)

var (
	ovMu       sync.Mutex
	ovState    int
	ovText     string
	ovFlashEnd time.Time
	ovRecorder *Recorder
	ovHwnd     uintptr

	ovAnim atomic.Bool

	ovOnce  sync.Once
	ovReady = make(chan struct{})

	ovHistory [22]float64
	ovTick    int
	ovFont    uintptr

	ovOnCancel func()
)

func ovInCloseZone(x, y int32) bool {
	w := overlayWidth()
	dpi := dpiFor(overlayHwnd())
	h := scaleDPI(ovH, dpi)
	return x >= w-scaleDPI(40, dpi) && x <= w-scaleDPI(4, dpi) && y >= 0 && y <= h
}

func ovCancelActive() bool {
	ovMu.Lock()
	st := ovState
	ovMu.Unlock()
	return st == ovRecording || st == ovProcessing
}


var (
	ovWidth   = scaleDPI(ovW, dpiFor(0))
	ovWidthMu sync.Mutex
	ovFontDPI int32
)

func overlayWidth() int32 {
	ovWidthMu.Lock()
	defer ovWidthMu.Unlock()
	return ovWidth
}

func overlayFont() uintptr { return uiFont(overlayHwnd()) }

func uiFont(hwnd uintptr) uintptr {
	dpi := dpiFor(hwnd)
	if ovFont != 0 && ovFontDPI == dpi {
		return ovFont
	}
	if ovFont != 0 {
		procDeleteObject.Call(ovFont)
	}
	face, _ := windows.UTF16PtrFromString("Consolas")
	h := scaleDPI(15, dpi)
	ovFont, _, _ = procCreateFontW.Call(uintptr(^uintptr(h)+1), 0, 0, 0, 400,
		0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(face)))
	ovFontDPI = dpi
	return ovFont
}

func overlayHwnd() uintptr {
	ovMu.Lock()
	defer ovMu.Unlock()
	return ovHwnd
}

func measureOverlayWidth(state int, text string) int32 {
	dpi := dpiFor(overlayHwnd())
	base := scaleDPI(ovW, dpi)
	if text == "" || state == ovRecording {
		return base
	}
	hdc, _, _ := procGetDC.Call(0)
	if hdc == 0 {
		return base
	}
	defer procReleaseDC.Call(0, hdc)
	old, _, _ := procSelectObject.Call(hdc, overlayFont())
	u, err := windows.UTF16FromString(text)
	if err != nil {
		return base
	}
	r := rect{}
	procDrawTextW.Call(hdc, uintptr(unsafe.Pointer(&u[0])), uintptr(len(u)-1),
		uintptr(unsafe.Pointer(&r)), 0x0400|0x0020)
	if old != 0 {
		procSelectObject.Call(hdc, old)
	}
	need := r.Right - r.Left + scaleDPI(44+36+12, dpi)
	if need < base {
		return base
	}
	var wa rect
	procSystemParametersInfoW.Call(0x30, 0, uintptr(unsafe.Pointer(&wa)), 0)
	max := (wa.Right - wa.Left) * 4 / 5
	if max > 0 && need > max {
		need = max
	}
	return need
}

func resizeOverlay(hwnd uintptr, width int32) {
	ovWidthMu.Lock()
	same := ovWidth == width
	ovWidth = width
	ovWidthMu.Unlock()
	if same || hwnd == 0 {
		return
	}
	var wa rect
	procSystemParametersInfoW.Call(0x30, 0, uintptr(unsafe.Pointer(&wa)), 0)
	x := wa.Left + (wa.Right-wa.Left-width)/2
	dpi := dpiFor(hwnd)
	h := scaleDPI(ovH, dpi)
	y := wa.Bottom - h - scaleDPI(28, dpi)
	procSetWindowPos.Call(hwnd, 0, uintptr(x), uintptr(y), uintptr(width), uintptr(h), 0x0004|0x0010)
}
func overlaySet(state int, text string) {
	ovOnce.Do(startOverlayThread)
	<-ovReady
	ovMu.Lock()
	ovState = state
	ovText = text
	if state == ovFlashOK {
		ovFlashEnd = time.Now().Add(1500 * time.Millisecond)
	}
	if state == ovFlashErr {
		ovFlashEnd = time.Now().Add(5500 * time.Millisecond)
	}
	hwnd := ovHwnd
	ovMu.Unlock()
	resizeOverlay(hwnd, measureOverlayWidth(state, text))
	if hwnd != 0 {
		procPostMessageW.Call(hwnd, wmOvSet, 0, 0)
	}
}

func overlayHide() { overlaySet(ovHidden, "") }

func startOverlayThread() {
	ovAnim.Store(true)
	go func() {
		runtime.LockOSThread()
		className, _ := windows.UTF16PtrFromString(appid.Class("OverlayWnd"))
		cb := syscall.NewCallback(overlayWndProc)
		wc := wndClassExW{
			Size:      uint32(unsafe.Sizeof(wndClassExW{})),
			WndProc:   cb,
			ClassName: className,
		}
		procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))

		var wa rect
		procSystemParametersInfoW.Call(0x30, 0, uintptr(unsafe.Pointer(&wa)), 0)
		sysDPI := int(dpiFor(0))
		startW := ovW * sysDPI / 96
		startH := ovH * sysDPI / 96
		x := int(wa.Left) + (int(wa.Right-wa.Left)-startW)/2
		y := int(wa.Bottom) - startH - 28*sysDPI/96

		hwnd, _, _ := procCreateWindowExW.Call(
			wsExLayered|wsExToolWindow|wsExNoActivate|0x00000008,
			uintptr(unsafe.Pointer(className)), 0,
			wsPopup,
			uintptr(x), uintptr(y), uintptr(startW), uintptr(startH),
			0, 0, 0, 0,
		)
		if hwnd != 0 {
			applyDarkCaption(hwnd)
			procSetLayeredWindowAttributes.Call(hwnd, 0, 240, lwaAlpha)
		}
		ovMu.Lock()
		ovHwnd = hwnd
		ovMu.Unlock()
		close(ovReady)
		if hwnd == 0 {
			return
		}
		var m msgStruct
		for {
			r, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
			if int32(r) <= 0 {
				return
			}
			procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
			procDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
		}
	}()
}

func overlayWndProc(hwnd, msg, wParam, lParam uintptr) uintptr {
	switch msg {
	case wmOvSet:
		ovMu.Lock()
		st := ovState
		ovMu.Unlock()
		if st == ovHidden {
			procKillTimer.Call(hwnd, ovTimerID)
			procShowWindow.Call(hwnd, swHide)
		} else {
			if st == ovRecording {
				ovHistory = [22]float64{}
			}
			procSetTimer.Call(hwnd, ovTimerID, 33, 0)
			procShowWindow.Call(hwnd, swShowNoActivate)
			procInvalidateRect.Call(hwnd, 0, 0)
		}
		return 0
	case wmTimer:
		ovTick++
		ovMu.Lock()
		st := ovState
		flashEnd := ovFlashEnd
		rec := ovRecorder
		ovMu.Unlock()
		if st == ovRecording && rec != nil {
			copy(ovHistory[:], ovHistory[1:])
			ovHistory[len(ovHistory)-1] = rec.Level()
		}
		if (st == ovFlashOK || st == ovFlashErr) && time.Now().After(flashEnd) {
			ovMu.Lock()
			ovState = ovHidden
			ovMu.Unlock()
			procKillTimer.Call(hwnd, ovTimerID)
			procShowWindow.Call(hwnd, swHide)
			return 0
		}
		overlayRenderDirect(hwnd)
		return 0
	case wmPaint:
		overlayPaint(hwnd)
		return 0
	case wmLBtnUp:
		x := int32(int16(lParam & 0xFFFF))
		y := int32(int16(lParam >> 16 & 0xFFFF))
		if ovCancelActive() && ovInCloseZone(x, y) && ovOnCancel != nil {
			ovOnCancel()
		}
		return 0
	case wmSetCursor:
		var pt point
		procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
		procScreenToClient.Call(hwnd, uintptr(unsafe.Pointer(&pt)))
		if ovCancelActive() && ovInCloseZone(pt.X, pt.Y) {
			cur, _, _ := procLoadCursorW.Call(0, 32649)
			procSetCursor.Call(cur)
			return 1
		}
	case wmDestroy:
		procPostQuitMessage.Call(0)
		return 0
	}
	r, _, _ := procDefWindowProcW.Call(hwnd, msg, wParam, lParam)
	return r
}

func stateColors(st int) (bright, dim uintptr) {
	switch st {
	case ovRecording:
		return colRed, colRedDm
	case ovProcessing:
		return colAmber, colAmberDm
	case ovFlashErr:
		return colAmber, colAmberDm
	default:
		return colGreen, colGreenDm
	}
}

func overlayPaint(hwnd uintptr) {
	var ps paintStruct
	hdc, _, _ := procBeginPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
	overlayRender(hwnd, hdc)
	procEndPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
}

func overlayRenderDirect(hwnd uintptr) {
	hdc, _, _ := procGetDC.Call(hwnd)
	if hdc == 0 {
		return
	}
	overlayRender(hwnd, hdc)
	procReleaseDC.Call(hwnd, hdc)
}

func overlayRender(hwnd, hdc uintptr) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("overlayRender: panic: %v", r)
		}
	}()
	var rc rect
	procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
	dpi := dpiFor(hwnd)
	px := func(v int32) int32 { return scaleDPI(v, dpi) }

	ovMu.Lock()
	st := ovState
	text := ovText
	ovMu.Unlock()
	anim := ovAnim.Load()
	bright, _ := stateColors(st)

	fill := func(r rect, color uintptr) {
		br, _, _ := procCreateSolidBrush.Call(color)
		procFillRect.Call(hdc, uintptr(unsafe.Pointer(&r)), br)
		procDeleteObject.Call(br)
	}

	fill(rc, colBg)
	for y := rc.Top + 3; y < rc.Bottom; y += 3 {
		fill(rect{Left: rc.Left, Top: y, Right: rc.Right, Bottom: y + 1}, colBgLine)
	}
	pulse := 0.0
	if anim && (st == ovRecording || st == ovProcessing) {
		pulse = math.Abs(math.Sin(float64(ovTick) * 0.18))
	}
	drawDot := func(cx, cy, r int32, color uintptr) {
		br, _, _ := procCreateSolidBrush.Call(color)
		oldBr, _, _ := procSelectObject.Call(hdc, br)
		pen, _, _ := procGetStockObject.Call(nullPen)
		oldPen, _, _ := procSelectObject.Call(hdc, pen)
		procEllipse.Call(hdc, uintptr(cx-r), uintptr(cy-r), uintptr(cx+r+1), uintptr(cy+r+1))
		procSelectObject.Call(hdc, oldPen)
		procSelectObject.Call(hdc, oldBr)
		procDeleteObject.Call(br)
	}
	cy := rc.Bottom / 2
	glowR := px(11) + int32(pulse*3)
	steps := int32(5)
	for i := int32(0); i < steps; i++ {
		r := glowR - i*(glowR-px(5))/(steps-1)
		t := 0.10 + 0.16*float64(i)
		drawDot(px(25), cy, r, blendCol(colBg, bright, t))
	}
	drawDot(px(25), cy, px(5), bright)
	drawDot(px(25), cy, px(2), blendCol(bright, 0xFFFFFF, 0.45))

	overlayFont()
	if st == ovProcessing {
		if text == "" {
			text = tr("ov.transcribing")
		}
		if anim {
			text += strings.Repeat(".", 1+(ovTick/10)%3)
		} else {
			text += "…"
		}
	}
	procSelectObject.Call(hdc, ovFont)
	procSetBkMode.Call(hdc, 1)
	txtRc := rect{Left: px(44), Top: 0, Right: rc.Right - px(12), Bottom: rc.Bottom}
	if st == ovRecording || st == ovProcessing {
		txtRc.Right = rc.Right - px(36)
	}
	if st == ovRecording {
		txtRc.Right = rc.Right - px(190)
	}
	u, _ := windows.UTF16FromString(text)
	drawText := func(r rect, color uintptr) {
		procSetTextColor.Call(hdc, color)
		procDrawTextW.Call(hdc, uintptr(unsafe.Pointer(&u[0])), uintptr(len(u)-1),
			uintptr(unsafe.Pointer(&r)), 0x0020|0x0004|0x8000)
	}
	for _, off := range [][2]int32{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
		drawText(rect{Left: txtRc.Left + off[0], Top: txtRc.Top + off[1], Right: txtRc.Right + off[0], Bottom: txtRc.Bottom + off[1]}, colGreenLo)
	}
	drawText(txtRc, bright)

	if st == ovRecording && anim {
		for i, v := range ovHistory {
			h := px(int32(3 + v*36))
			if h > px(42) {
				h = px(42)
			}
			x := rc.Right - px(186) + int32(i)*px(7)
			fill(rect{Left: x - px(1), Top: cy - h/2 - px(1), Right: x + px(5), Bottom: cy + h/2 + px(1)}, colGreenLo)
			fill(rect{Left: x, Top: cy - h/2, Right: x + px(4), Bottom: cy + h/2}, bright)
		}
	}

	if st == ovRecording || st == ovProcessing {
		xr := rect{Left: rc.Right - px(34), Top: 0, Right: rc.Right - px(8), Bottom: rc.Bottom}
		xs, _ := windows.UTF16FromString("✕")
		procSetTextColor.Call(hdc, colGreenDm)
		procDrawTextW.Call(hdc, uintptr(unsafe.Pointer(&xs[0])), uintptr(len(xs)-1),
			uintptr(unsafe.Pointer(&xr)), 0x0020|0x0004|0x0001)
	}
}
