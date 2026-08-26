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
	"holdtotype/internal/audiolevel"
)

const (
	ovHidden = iota
	ovRecording
	ovProcessing
	ovFlashOK
	ovFlashErr
	ovPaused
)

const (
	wmOvSet      = 0x0400 + 10
	wmTimer      = 0x0113
	wmMouseMove  = 0x0200
	wmLBtnUp     = 0x0202
	wmSetCursor  = 0x0020
	wmDpiChanged = 0x02E0
	ovW          = 390
	ovH          = 52
	ovMaxLines   = 6
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
	procFrameRect                  = user32.NewProc("FrameRect")
	procMonitorFromWindow          = user32.NewProc("MonitorFromWindow")
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

var (
	colBg      uintptr = 0x0C0F0B
	colBgLine  uintptr = 0x0A0C09
	colGreen   uintptr = 0x6EFF3C
	colGreenDm uintptr = 0x4AA320
	colGreenLo uintptr = 0x2B4A1D
	colRed     uintptr = 0x6B6BFF
	colRedDm   uintptr = 0x26265C
	colAmber   uintptr = 0x47B3FF
	colAmberDm uintptr = 0x20455C
	colBad     uintptr = 0x6B6BFF
	colBadDm   uintptr = 0x26265C
	colAskBg   uintptr = 0x0D100B
	colHi      uintptr = 0x6EFF3C
	colHiLo    uintptr = 0x2B4A1D
	colLine    uintptr = 0x2B4A1D
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
	ovRows    atomic.Int32
	ovLineH   atomic.Int32
	ovFlashAt int
	ovTick    int

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
	return st == ovRecording || st == ovProcessing || st == ovPaused
}

func ovShowsClose(st int) bool {
	return st == ovRecording || st == ovProcessing || st == ovPaused || st == ovFlashOK || st == ovFlashErr
}

func ovFlashState() bool {
	ovMu.Lock()
	st := ovState
	ovMu.Unlock()
	return st == ovFlashOK || st == ovFlashErr
}

var (
	ovWidth   = int32(ovW)
	ovHeight  = int32(ovH)
	ovX       int32
	ovY       int32
	ovWidthMu sync.Mutex
)

func overlayWidth() int32 {
	ovWidthMu.Lock()
	defer ovWidthMu.Unlock()
	return ovWidth
}

func overlayFont() uintptr { return uiFont(overlayHwnd()) }

func overlayMeasureFont() uintptr { return uiFontDPI(overlayDPI()) }

var (
	fontMu    sync.Mutex
	fontCache = map[int32]uintptr{}
)

func uiFont(hwnd uintptr) uintptr { return uiFontDPI(dpiFor(hwnd)) }

func uiFontDPI(dpi int32) uintptr {
	skin := themeLook()
	fontMu.Lock()
	defer fontMu.Unlock()
	if f, ok := fontCache[dpi]; ok {
		return f
	}
	face, _ := windows.UTF16PtrFromString(skin.FontGDI)
	h := scaleDPI(skin.FontPx, dpi)
	f, _, _ := procCreateFontW.Call(uintptr(^uintptr(h)+1), 0, 0, 0, uintptr(skin.Weight),
		0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(face)))
	fontCache[dpi] = f
	log.Printf("шрифт интерфейса создан: %s %d px для %d DPI", skin.FontGDI, skin.FontPx, dpi)
	return f
}

// dropFontCache is called when the skin changes: the next paint asks for a new face.
func dropFontCache() {
	fontMu.Lock()
	old := fontCache
	fontCache = map[int32]uintptr{}
	fontMu.Unlock()
	for _, f := range old {
		if f != 0 {
			procDeleteObject.Call(f)
		}
	}
}

func overlayHwnd() uintptr {
	ovMu.Lock()
	defer ovMu.Unlock()
	return ovHwnd
}

func measureOverlayWidth(state int, text string) int32 {
	need := measureStatusWidth(state, text)
	dpi := overlayDPI()
	if aw := askWidthDIP(); aw > 0 {
		if w := scaleDPI(aw, dpi); w > need {
			need = w
		}
	}
	if state == ovFlashOK || state == ovFlashErr {
		if lim := scaleDPI(640, dpi); need > lim {
			need = lim
		}
	}
	wa := overlayWorkArea()
	if max := (wa.Right - wa.Left) * 4 / 5; max > 0 && need > max {
		need = max
	}
	return need
}

func measureStatusWidth(state int, text string) int32 {
	dpi := overlayDPI()
	base := scaleDPI(ovW, dpi)
	if text == "" || state == ovRecording {
		return base
	}
	hdc, _, _ := procGetDC.Call(0)
	if hdc == 0 {
		return base
	}
	defer procReleaseDC.Call(0, hdc)
	old, _, _ := procSelectObject.Call(hdc, overlayMeasureFont())
	if state == ovProcessing {
		text += "..."
	}
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
	wa := overlayWorkArea()
	max := (wa.Right - wa.Left) * 4 / 5
	if max > 0 && need > max {
		need = max
	}
	return need
}

func overlayHeightPx(dpi int32) int32 {
	base := scaleDPI(ovH, dpi)
	if askActive() {
		return base + scaleDPI(ovAskH*askRows(), dpi)
	}
	rows, lineH := ovRows.Load(), ovLineH.Load()
	if rows > 1 && lineH > 0 {
		return base + (rows-1)*lineH
	}
	return base
}

func ovTextWraps(state int) bool { return state == ovFlashOK || state == ovFlashErr }

func ovTextRoom(state int, width int32, dpi int32) int32 {
	px := func(v int32) int32 { return scaleDPI(v, dpi) }
	right := width - px(12)
	if ovShowsClose(state) {
		right = width - px(44)
	}
	return right - px(44)
}

func measureTextRows(state int, text string, width int32) (rows, lineH int32) {
	dpi := overlayDPI()
	if !ovTextWraps(state) || askActive() || strings.TrimSpace(text) == "" {
		return 1, 0
	}
	room := ovTextRoom(state, width, dpi)
	if room <= 0 {
		return 1, 0
	}
	hdc, _, _ := procGetDC.Call(0)
	if hdc == 0 {
		return 1, 0
	}
	defer procReleaseDC.Call(0, hdc)
	old, _, _ := procSelectObject.Call(hdc, uiFontDPI(dpi))
	defer func() {
		if old != 0 {
			procSelectObject.Call(hdc, old)
		}
	}()
	probe, err := windows.UTF16FromString("Ag")
	if err != nil {
		return 1, 0
	}
	one := rect{}
	procDrawTextW.Call(hdc, uintptr(unsafe.Pointer(&probe[0])), uintptr(len(probe)-1),
		uintptr(unsafe.Pointer(&one)), 0x0400|0x0020)
	lineH = one.Bottom - one.Top
	if lineH <= 0 {
		return 1, 0
	}
	u, uerr := windows.UTF16FromString(text)
	if uerr != nil {
		return 1, lineH
	}
	r := rect{Right: room}
	procDrawTextW.Call(hdc, uintptr(unsafe.Pointer(&u[0])), uintptr(len(u)-1),
		uintptr(unsafe.Pointer(&r)), 0x0400|0x0010)
	rows = (r.Bottom - r.Top + lineH - 1) / lineH
	if rows < 1 {
		rows = 1
	}
	if rows > ovMaxLines {
		rows = ovMaxLines
	}
	return rows, lineH
}

func resizeOverlay(hwnd uintptr, width int32) {
	if hwnd == 0 {
		return
	}
	h := overlayHeightPx(overlayDPI())
	x, y := overlayOrigin(width, h)
	ovWidthMu.Lock()
	same := ovWidth == width && ovHeight == h && ovX == x && ovY == y
	ovWidth, ovHeight, ovX, ovY = width, h, x, y
	ovWidthMu.Unlock()
	if same {
		return
	}
	procSetWindowPos.Call(hwnd, 0, uintptr(x), uintptr(y), uintptr(width), uintptr(h), 0x0004|0x0010)
}

func overlayRefresh() {
	hwnd := overlayHwnd()
	if hwnd == 0 {
		return
	}
	ovMu.Lock()
	state, text := ovState, ovText
	ovMu.Unlock()
	width := measureOverlayWidth(state, text)
	rows, lineH := measureTextRows(state, text, width)
	ovRows.Store(rows)
	ovLineH.Store(lineH)
	resizeOverlay(hwnd, width)
	procPostMessageW.Call(hwnd, wmOvSet, 0, 0)
}
func overlaySet(state int, text string) {
	ovOnce.Do(startOverlayThread)
	<-ovReady
	ovMu.Lock()
	ovState = state
	ovText = text
	if state == ovFlashOK {
		ovFlashEnd = time.Now().Add(flashLife(text, 1500))
	}
	if state == ovFlashErr {
		ovFlashEnd = time.Now().Add(flashLife(text, 3000))
	}
	hwnd := ovHwnd
	ovMu.Unlock()
	width := measureOverlayWidth(state, text)
	rows, lineH := measureTextRows(state, text, width)
	ovRows.Store(rows)
	ovLineH.Store(lineH)
	resizeOverlay(hwnd, width)
	if hwnd != 0 {
		procPostMessageW.Call(hwnd, wmOvSet, 0, 0)
	}
}

func flashLife(text string, base int) time.Duration {
	d := time.Duration(base)*time.Millisecond + time.Duration(len([]rune(text)))*30*time.Millisecond
	if d > 5*time.Second {
		d = 5 * time.Second
	}
	return d
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

		sysDPI := int(dpiForCursor())
		startW := ovW * sysDPI / 96
		startH := ovH * sysDPI / 96
		sx, sy := overlayOrigin(int32(startW), int32(startH))
		x, y := int(sx), int(sy)

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
			if st == ovFlashOK {
				ovFlashAt = ovTick
			}
			procSetTimer.Call(hwnd, ovTimerID, 33, 0)
			procShowWindow.Call(hwnd, swShowNoActivate)
			procInvalidateRect.Call(hwnd, 0, 0)
		}
		return 0
	case wmDpiChanged:
		log.Printf("оверлей: масштаб экрана сменился, пересчитываю плашку")
		go overlayRefresh()
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
		askTick()
		_, counting := askLeft()
		if !counting && st != ovRecording && (!ovAnim.Load() || st != ovProcessing) {
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
		if id, hit := askHit(x, y); hit {
			askFinish(id)
			return 0
		}
		if ovInCloseZone(x, y) && ovFlashState() {
			log.Printf("оверлей: вспышка закрыта нажатием")
			overlayHide()
			return 0
		}
		if ovInCloseZone(x, y) && (ovCancelActive() || askActive()) {
			if askActive() {
				askFinish("")
			}
			if ovOnCancel != nil {
				ovOnCancel()
			}
		}
		return 0
	case wmSetCursor:
		var pt point
		procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
		procScreenToClient.Call(hwnd, uintptr(unsafe.Pointer(&pt)))
		if _, hit := askHit(pt.X, pt.Y); hit {
			cur, _, _ := procLoadCursorW.Call(0, 32649)
			procSetCursor.Call(cur)
			return 1
		}
		if (ovCancelActive() || askActive() || ovFlashState()) && ovInCloseZone(pt.X, pt.Y) {
			cur, _, _ := procLoadCursorW.Call(0, 32649)
			procSetCursor.Call(cur)
			return 1
		}
	case wmDestroy:
		ovBuf.release()
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
	case ovPaused:
		return colAmber, colAmberDm
	case ovFlashErr:
		return colBad, colBadDm
	default:
		return colGreen, colGreenDm
	}
}

var ovBuf backBuf

func overlayFrame(hwnd, hdc uintptr) {
	var rc rect
	procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
	mem := ovBuf.begin(hdc, rc.Right-rc.Left, rc.Bottom-rc.Top)
	if mem == 0 {
		overlayRender(hwnd, hdc)
		return
	}
	overlayRender(hwnd, mem)
	ovBuf.blit(hdc)
}

func overlayPaint(hwnd uintptr) {
	var ps paintStruct
	hdc, _, _ := procBeginPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
	overlayFrame(hwnd, hdc)
	procEndPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
}

func overlayRenderDirect(hwnd uintptr) {
	hdc, _, _ := procGetDC.Call(hwnd)
	if hdc == 0 {
		return
	}
	overlayFrame(hwnd, hdc)
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
	left := overlayCountdown()
	bright, _ := stateColors(st)
	if st == ovRecording && left >= 0 {
		bright = colAmber
	}

	fill := func(r rect, color uintptr) {
		br, _, _ := procCreateSolidBrush.Call(color)
		procFillRect.Call(hdc, uintptr(unsafe.Pointer(&r)), br)
		procDeleteObject.Call(br)
	}

	fill(rc, colBg)
	if themeScanlines() {
		for y := rc.Top + 3; y < rc.Bottom; y += 3 {
			fill(rect{Left: rc.Left, Top: y, Right: rc.Right, Bottom: y + 1}, colBgLine)
		}
	}
	if !themeRoundCorners() {
		border := rc
		brush, _, _ := procCreateSolidBrush.Call(colLine)
		procFrameRect.Call(hdc, uintptr(unsafe.Pointer(&border)), brush)
		procDeleteObject.Call(brush)
	}
	pulse := 0.0
	if anim && !askActive() && (st == ovRecording || st == ovProcessing) {
		pulse = math.Abs(math.Sin(float64(ovTick) * 0.18 / themePulse()))
	}
	cy := px(ovH) / 2
	dotX := px(25)
	core := int32(float64(px(5)) * (0.88 + 0.2*pulse))
	halo := core
	if themeGlow() {
		halo = core + px(5)
	}
	hideDot := false
	if anim && st == ovFlashOK && !askActive() {
		age := float64(ovTick - ovFlashAt)
		switch themeFlash() {
		case "blink":
			hideDot = age < 12 && int(age/3)%2 == 1
		case "glow":
			if k := 1 - age/18; k > 0 {
				halo = core + px(5) + int32(float64(px(10))*k)
			}
		case "bounce":
			if k := math.Exp(-age / 7); k > 0.02 {
				core = int32(float64(core) * (1 + 0.55*k*math.Abs(math.Sin(age*0.45))))
			}
		}
	}
	switch st {
	case ovPaused:
		w, h, gap := px(2), px(10), px(3)
		fill(rect{Left: dotX - gap - w, Top: cy - h/2, Right: dotX - gap, Bottom: cy + h/2}, bright)
		fill(rect{Left: dotX + gap, Top: cy - h/2, Right: dotX + gap + w, Bottom: cy + h/2}, bright)
	case ovFlashErr:
		s := px(5)
		fill(rect{Left: dotX - s, Top: cy - s, Right: dotX + s, Bottom: cy + s}, bright)
	default:
		if !hideDot {
			drawSmoothDot(hdc, dotX, cy, core, halo, bright, colBg, 0.34+0.12*pulse)
		}
	}

	overlayFont()
	if st == ovProcessing && !askActive() {
		if text == "" {
			text = tr("ov.transcribing")
		}
		if anim {
			text += strings.Repeat(".", 1+(ovTick/10)%3)
		} else {
			text += "…"
		}
	}
	procSelectObject.Call(hdc, overlayFont())
	procSetBkMode.Call(hdc, 1)
	txtRc := rect{Left: px(44), Top: 0, Right: rc.Right - px(12), Bottom: px(ovH)}
	if ovShowsClose(st) {
		txtRc.Right = rc.Right - px(44)
	}
	if st == ovRecording {
		txtRc.Right = rc.Right - px(196)
	}
	textFlags := uintptr(0x0020 | 0x0004 | 0x8000)
	if rows, lineH := ovRows.Load(), ovLineH.Load(); rows > 1 && lineH > 0 && ovTextWraps(st) && !askActive() {
		txtRc.Top = cy - lineH/2
		txtRc.Bottom = txtRc.Top + lineH*rows
		textFlags = 0x0010 | 0x8000
	}
	u, _ := windows.UTF16FromString(text)
	drawText := func(r rect, color uintptr) {
		procSetTextColor.Call(hdc, color)
		procDrawTextW.Call(hdc, uintptr(unsafe.Pointer(&u[0])), uintptr(len(u)-1),
			uintptr(unsafe.Pointer(&r)), textFlags)
	}
	// the words keep the skin's own colour; what state the program is in is
	// said by the dot, not by tinting the text
	textCol, haloCol := colGreen, colGreenLo
	if st == ovFlashErr {
		textCol, haloCol = colBad, colBadDm
	}
	if themeGlow() {
		for _, off := range [][2]int32{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			drawText(rect{Left: txtRc.Left + off[0], Top: txtRc.Top + off[1], Right: txtRc.Right + off[0], Bottom: txtRc.Bottom + off[1]}, haloCol)
		}
	}
	drawText(txtRc, textCol)

	if st == ovRecording && left >= 0 {
		lr := rect{Left: rc.Right - px(190), Top: 0, Right: rc.Right - px(36), Bottom: px(ovH)}
		ls, _ := windows.UTF16FromString(trf("ov.left", left))
		procSetTextColor.Call(hdc, colAmber)
		procDrawTextW.Call(hdc, uintptr(unsafe.Pointer(&ls[0])), uintptr(len(ls)-1),
			uintptr(unsafe.Pointer(&lr)), 0x0020|0x0004|0x0002)
	} else if st == ovRecording && anim {
		style := themeLevelStyle()
		for i, v := range ovHistory {
			x := rc.Right - px(186) + int32(i)*px(7)
			lv := audiolevel.Heard(v)
			switch style {
			case "dots":
				r := px(3)
				drawSmoothDot(hdc, x+px(2), cy-int32(float64(px(8))*lv), r, r, colHi, colBg, 0)
			case "flat":
				half := (px(3) + int32(float64(px(13))*lv)) / 2
				fill(rect{Left: x, Top: cy - half, Right: x + px(4), Bottom: cy + half}, colHi)
			default:
				half := (px(4) + int32(float64(px(18))*lv)) / 2
				if themeGlow() && lv > 0 {
					fill(rect{Left: x - px(1), Top: cy - half - px(1), Right: x + px(5), Bottom: cy + half + px(1)}, colHiLo)
				}
				fill(rect{Left: x, Top: cy - half, Right: x + px(4), Bottom: cy + half}, colHi)
			}
		}
	}

	if ovShowsClose(st) {
		xr := rect{Left: rc.Right - px(34), Top: 0, Right: rc.Right - px(8), Bottom: px(ovH)}
		xs, _ := windows.UTF16FromString("✕")
		procSetTextColor.Call(hdc, colGreenDm)
		procDrawTextW.Call(hdc, uintptr(unsafe.Pointer(&xs[0])), uintptr(len(xs)-1),
			uintptr(unsafe.Pointer(&xr)), 0x0020|0x0004|0x0001)
	}

	askRender(hwnd, hdc, rc, fill, func(s string, r rect, color uintptr, flags uintptr) {
		t, _ := windows.UTF16FromString(s)
		procSetTextColor.Call(hdc, color)
		procDrawTextW.Call(hdc, uintptr(unsafe.Pointer(&t[0])), uintptr(len(t)-1),
			uintptr(unsafe.Pointer(&r)), flags)
	})
}

func overlayNote(text string) {
	ovMu.Lock()
	st := ovState
	ovMu.Unlock()
	if st != ovProcessing && st != ovRecording {
		overlaySet(ovFlashErr, text)
		return
	}
	log.Printf("оверлей: предупреждение поверх работы — %s", text)
	overlaySet(st, text)
}

func overlayWorkArea() rect {
	wa, _ := overlayArea()
	return wa
}

func overlayDPI() int32 {
	wa, _ := overlayArea()
	if d := dpiForPoint(wa.Left+(wa.Right-wa.Left)/2, wa.Top+(wa.Bottom-wa.Top)/2); d >= 72 {
		return d
	}
	return dpiFor(overlayHwnd())
}

var ovDeadline atomic.Int64

func overlaySetDeadline(t time.Time) { ovDeadline.Store(t.UnixMilli()) }

func overlayClearDeadline() { ovDeadline.Store(0) }

func overlayCountdown() int {
	ms := ovDeadline.Load()
	if ms == 0 {
		return -1
	}
	left := time.Until(time.UnixMilli(ms))
	if left > 10*time.Second || left < 0 {
		return -1
	}
	sec := int(left/time.Second) + 1
	if sec > 10 {
		sec = 10
	}
	return sec
}

