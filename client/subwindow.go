package main

import (
	"fmt"
	"log"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"holdtotype/internal/appid"
	"holdtotype/internal/ovplace"
	"holdtotype/internal/subtitles"
)

const (
	subPhaseLoading = iota
	subPhaseListening
	subPhaseRestarting
)

const (
	wmSwSet         = 0x0400 + 20
	wmSwLayout      = 0x0400 + 21
	wmMouseActivate = 0x0021
	wmCaptureLost   = 0x0215
	wmPrintClient   = 0x0318
	maNoActivate    = 3

	swTimerID   = 1
	swHeaderDIP = 30
	swPadX      = 16
	swPadTop    = 4
	swPadBottom = 12
	swMinW      = 520
	swMaxW      = 1100
	swShare     = 56
	swBottomGap = 150
	swAlpha     = 235
	swTextRunes = 900
	swCloseDIP  = 40
)

var swSizes = map[string]int32{"s": 15, "m": 19, "l": 25}

type subWindowSource interface {
	windowPhase() int
	windowRows() []subtitles.Row
	windowTranslated() bool
}

var (
	swMu      sync.Mutex
	swHwnd    uintptr
	swOnce    sync.Once
	swReady   = make(chan struct{})
	swSrc     subWindowSource
	swLang    string
	swRows    int32 = defaultSubtitlesLines
	swSize          = "m"
	swCustom  map[string]ovplace.Frac
	swOnClose func()
	swOnMoved func(key string, f ovplace.Frac)

	swShown     bool
	swHover     bool
	swDragging  bool
	swDragMoved bool
	swDragDX    int32
	swDragDY    int32
	swX, swY    int32
	swW, swH    int32
	swLineH     int32
	swDPI       int32 = 96
	swLastRows  []subtitles.Row
	swLaidTr    bool
	swLastPhase = -1
	swLastPaint time.Time
	swTicks     int
	swBuf       backBuf

	swFontMu sync.Mutex
	swFonts  = map[string]uintptr{}
)

func subWindowConfigure(rows int, size string, custom map[string]ovplace.Frac) {
	if rows < minSubtitlesLines || rows > maxSubtitlesLines {
		rows = defaultSubtitlesLines
	}
	if !validSubtitlesSize(size) {
		size = "m"
	}
	places := map[string]ovplace.Frac{}
	for k, v := range custom {
		places[k] = v
	}
	swMu.Lock()
	swRows = int32(rows)
	swSize = size
	swCustom = places
	hwnd := swHwnd
	swMu.Unlock()
	if hwnd != 0 {
		procPostMessageW.Call(hwnd, wmSwLayout, 0, 0)
	}
}

func subWindowShow(src subWindowSource, langCode string, onClose func(), onMoved func(string, ovplace.Frac)) {
	swOnce.Do(startSubWindowThread)
	<-swReady
	swMu.Lock()
	swSrc = src
	swLang = strings.ToUpper(strings.TrimSpace(langCode))
	swOnClose = onClose
	swOnMoved = onMoved
	hwnd := swHwnd
	swMu.Unlock()
	if hwnd != 0 {
		procPostMessageW.Call(hwnd, wmSwSet, 0, 0)
	}
}

func subWindowSetLang(src subWindowSource, langCode string) {
	swMu.Lock()
	if swSrc == src {
		swLang = strings.ToUpper(strings.TrimSpace(langCode))
	}
	swMu.Unlock()
}

func subWindowHide(src subWindowSource) {
	swMu.Lock()
	if swSrc != src {
		swMu.Unlock()
		return
	}
	had := swSrc != nil
	swSrc = nil
	swOnClose = nil
	swOnMoved = nil
	hwnd := swHwnd
	swMu.Unlock()
	if had && hwnd != 0 {
		procPostMessageW.Call(hwnd, wmSwSet, 0, 0)
	}
}

func startSubWindowThread() {
	go func() {
		runtime.LockOSThread()
		className, _ := windows.UTF16PtrFromString(appid.Class("SubtitlesWnd"))
		wc := wndClassExW{
			Size:      uint32(unsafe.Sizeof(wndClassExW{})),
			WndProc:   syscall.NewCallback(subWndProc),
			ClassName: className,
		}
		procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
		hwnd, _, _ := procCreateWindowExW.Call(
			wsExLayered|wsExToolWindow|wsExNoActivate|0x00000008,
			uintptr(unsafe.Pointer(className)), 0,
			wsPopup,
			0, 0, 600, 140,
			0, 0, 0, 0,
		)
		if hwnd != 0 {
			applyDarkCaption(hwnd)
			procSetLayeredWindowAttributes.Call(hwnd, 0, swAlpha, lwaAlpha)
		} else {
			log.Printf("subtitles: the window was not created")
		}
		swMu.Lock()
		swHwnd = hwnd
		swMu.Unlock()
		close(swReady)
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

func swFont(face string, weight int32, px int32) uintptr {
	key := fmt.Sprintf("%s/%d/%d", face, weight, px)
	swFontMu.Lock()
	defer swFontMu.Unlock()
	if f, ok := swFonts[key]; ok {
		return f
	}
	name, _ := windows.UTF16PtrFromString(face)
	f, _, _ := procCreateFontW.Call(uintptr(^uintptr(px)+1), 0, 0, 0, uintptr(weight),
		0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(name)))
	swFonts[key] = f
	return f
}

func swTextFont(dpi int32) uintptr {
	swMu.Lock()
	size := swSize
	swMu.Unlock()
	skin := themeLook()
	return swFont(skin.FontGDI, skin.Weight, scaleDPI(swSizes[size], dpi))
}

func swMeasureLine(font uintptr) int32 {
	hdc, _, _ := procGetDC.Call(0)
	if hdc == 0 {
		return 0
	}
	defer procReleaseDC.Call(0, hdc)
	old, _, _ := procSelectObject.Call(hdc, font)
	defer procSelectObject.Call(hdc, old)
	probe, _ := windows.UTF16FromString("Ag")
	r := rect{}
	procDrawTextW.Call(hdc, uintptr(unsafe.Pointer(&probe[0])), uintptr(len(probe)-1),
		uintptr(unsafe.Pointer(&r)), 0x0400|0x0020)
	return r.Bottom - r.Top
}

func swPlace() (work, screen rect) {
	if swShown && swW > 0 {
		return monitorRectsForPoint(swX+swW/2, swY+swH/2)
	}
	var pt point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	return monitorRectsForPoint(pt.X, pt.Y)
}

func swMeasure(wa rect, dpi int32, translated bool) (w, h, lineH int32) {
	px := func(v int32) int32 { return scaleDPI(v, dpi) }
	waW := wa.Right - wa.Left
	w = waW * swShare / 100
	if w < px(swMinW) {
		w = px(swMinW)
	}
	if w > px(swMaxW) {
		w = px(swMaxW)
	}
	if w > waW {
		w = waW
	}
	lineH = swMeasureLine(swTextFont(dpi))
	if lineH <= 0 {
		lineH = px(24)
	}
	swMu.Lock()
	rows := swRows
	swMu.Unlock()
	if translated {
		rows *= 2
	}
	h = px(swHeaderDIP) + px(swPadTop) + rows*lineH + px(swPadBottom)
	return w, h, lineH
}

func swTranslatedNow() bool {
	swMu.Lock()
	src := swSrc
	swMu.Unlock()
	return src != nil && src.windowTranslated()
}

func swLayout(hwnd uintptr) {
	wa, screen := swPlace()
	dpi := dpiForPoint(wa.Left+(wa.Right-wa.Left)/2, wa.Top+(wa.Bottom-wa.Top)/2)
	translated := swTranslatedNow()
	w, h, lineH := swMeasure(wa, dpi, translated)
	swMu.Lock()
	custom := swCustom
	swMu.Unlock()
	var x, y int32
	if f, ok := custom[ovplace.ResKey(toPlace(screen))]; ok {
		x, y = ovplace.Custom(f, toPlace(wa), w, h)
	} else {
		x = wa.Left + (wa.Right-wa.Left-w)/2
		y = wa.Bottom - scaleDPI(swBottomGap, dpi) - h
		if y < wa.Top {
			y = wa.Top
		}
	}
	swX, swY, swW, swH, swLineH, swDPI = x, y, w, h, lineH, dpi
	swLaidTr = translated
	procSetWindowPos.Call(hwnd, ^uintptr(0), uintptr(x), uintptr(y), uintptr(w), uintptr(h), swpNoActivate)
	procSetWindowRgn.Call(hwnd, swRegion(w, h, dpi), 1)
}

func swRescale(hwnd uintptr) bool {
	var pt point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	cx, cy := pt.X, pt.Y
	if !swDragging {
		cx, cy = swX+swW/2, swY+swH/2
	}
	dpi := dpiForPoint(cx, cy)
	if dpi == swDPI || swW <= 0 || swH <= 0 {
		return false
	}
	wa, _ := monitorRectsForPoint(cx, cy)
	w, h, lineH := swMeasure(wa, dpi, swLaidTr)
	var x, y int32
	if swDragging {
		swDragDX = int32(int64(swDragDX) * int64(w) / int64(swW))
		swDragDY = int32(int64(swDragDY) * int64(h) / int64(swH))
		x, y = pt.X-swDragDX, pt.Y-swDragDY
	} else {
		x, y = cx-w/2, cy-h/2
	}
	log.Printf("subtitles: the window is on a screen at %d DPI now, %dx%d", dpi, w, h)
	swX, swY, swW, swH, swLineH, swDPI = x, y, w, h, lineH, dpi
	procSetWindowPos.Call(hwnd, ^uintptr(0), uintptr(x), uintptr(y), uintptr(w), uintptr(h), swpNoActivate)
	procSetWindowRgn.Call(hwnd, swRegion(w, h, dpi), 1)
	swRenderNow(hwnd)
	return true
}

func swRegion(w, h, dpi int32) uintptr {
	if !themeRoundCorners() || w <= 0 || h <= 0 {
		return 0
	}
	r := scaleDPI(themeLook().Radius, dpi)
	if r < 2 {
		r = 2
	}
	if max := h / 2; r > max {
		r = max
	}
	rgn, _, _ := procCreateRoundRectRgn.Call(0, 0, uintptr(w+1), uintptr(h+1), uintptr(2*r), uintptr(2*r))
	return rgn
}

func swInClose(x, y int32) bool {
	px := func(v int32) int32 { return scaleDPI(v, swDPI) }
	return x >= swW-px(swCloseDIP) && x < swW && y >= 0 && y < px(swHeaderDIP)
}

func swCursorClient(hwnd uintptr) (int32, int32) {
	var pt point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	procScreenToClient.Call(hwnd, uintptr(unsafe.Pointer(&pt)))
	return pt.X, pt.Y
}

func swContent() (int, []subtitles.Row, bool) {
	swMu.Lock()
	src := swSrc
	swMu.Unlock()
	if src == nil {
		return 0, nil, false
	}
	return src.windowPhase(), src.windowRows(), true
}

func swApplyState(hwnd uintptr) {
	_, _, on := swContent()
	if !on {
		if swDragging {
			procReleaseCapture.Call()
			swDragging = false
		}
		procKillTimer.Call(hwnd, swTimerID)
		procShowWindow.Call(hwnd, swHide)
		swShown = false
		swLastRows, swLastPhase = nil, -1
		return
	}
	if !swShown {
		swLayout(hwnd)
	}
	swShown = true
	procSetTimer.Call(hwnd, swTimerID, 100, 0)
	procShowWindow.Call(hwnd, swShowNoActivate)
	overlayRaise(hwnd)
	swRenderNow(hwnd)
}

func swRenderNow(hwnd uintptr) {
	hdc, _, _ := procGetDC.Call(hwnd)
	if hdc == 0 {
		return
	}
	swFrame(hwnd, hdc)
	procReleaseDC.Call(hwnd, hdc)
}

func swFrame(hwnd, hdc uintptr) {
	var rc rect
	procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
	mem := swBuf.begin(hdc, rc.Right-rc.Left, rc.Bottom-rc.Top)
	if mem == 0 {
		swRender(hwnd, hdc)
		return
	}
	swRender(hwnd, mem)
	swBuf.blit(hdc)
	swLastPaint = time.Now()
}

func swPhaseText(phase int) string {
	switch phase {
	case subPhaseLoading:
		return tr("sub.loading")
	case subPhaseRestarting:
		return tr("sub.restarting")
	}
	return tr("sub.listening")
}

func swRender(hwnd, hdc uintptr) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("subtitles window: panic: %v", r)
		}
	}()
	var rc rect
	procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
	dpi := swDPI
	px := func(v int32) int32 { return scaleDPI(v, dpi) }
	phase, rows, _ := swContent()
	swLastRows, swLastPhase = rows, phase

	fill := func(r rect, color uintptr) {
		br, _, _ := procCreateSolidBrush.Call(color)
		procFillRect.Call(hdc, uintptr(unsafe.Pointer(&r)), br)
		procDeleteObject.Call(br)
	}
	draw := func(s string, r rect, color uintptr, flags uintptr) rect {
		u, _ := windows.UTF16FromString(s)
		procSetTextColor.Call(hdc, color)
		procDrawTextW.Call(hdc, uintptr(unsafe.Pointer(&u[0])), uintptr(len(u)-1),
			uintptr(unsafe.Pointer(&r)), flags)
		return r
	}
	fill(rc, colBg)
	if themeScanlines() {
		for y := rc.Top + 3; y < rc.Bottom; y += 3 {
			fill(rect{Left: rc.Left, Top: y, Right: rc.Right, Bottom: y + 1}, colBgLine)
		}
	}
	brush, _, _ := procCreateSolidBrush.Call(colLine)
	if rgn := swRegion(rc.Right-rc.Left, rc.Bottom-rc.Top, dpi); rgn != 0 {
		procFrameRgn.Call(hdc, rgn, brush, 1, 1)
		procDeleteObject.Call(rgn)
	} else {
		border := rc
		procFrameRect.Call(hdc, uintptr(unsafe.Pointer(&border)), brush)
	}
	procDeleteObject.Call(brush)
	procSetBkMode.Call(hdc, 1)

	head := px(swHeaderDIP)
	cy := head / 2
	dot := colDot
	if phase != subPhaseListening {
		dot = colAmber
	}
	core := px(4)
	halo := core
	if themeGlow() {
		halo = core + px(4)
	}
	drawSmoothDot(hdc, px(18), cy, core, halo, dot, colBg, 0.34)

	procSelectObject.Call(hdc, uiFontDPI(dpi))
	title := tr("sub.title")
	measure := draw(title, rect{}, colGreen, 0x0400|0x0020)
	tr0 := rect{Left: px(32), Top: 0, Right: px(32) + (measure.Right - measure.Left), Bottom: head}
	draw(title, tr0, colGreen, 0x0020|0x0004)
	swMu.Lock()
	code := swLang
	swMu.Unlock()
	if code != "" {
		draw(code, rect{Left: tr0.Right + px(10), Top: 0, Right: rc.Right - px(swCloseDIP), Bottom: head}, colGreenDm, 0x0020|0x0004|0x8000)
	}
	xcol := colGreenDm
	if swHover {
		xcol = colGreen
	}
	draw("✕", rect{Left: rc.Right - px(swCloseDIP), Top: 0, Right: rc.Right - px(6), Bottom: head}, xcol, 0x0020|0x0004|0x0001)
	fill(rect{Left: px(10), Top: head, Right: rc.Right - px(10), Bottom: head + 1}, colLine)

	area := rect{Left: px(swPadX), Top: head + px(swPadTop), Right: rc.Right - px(swPadX), Bottom: rc.Bottom - px(swPadBottom)}
	font := swTextFont(dpi)
	procSelectObject.Call(hdc, font)
	if !swHasText(rows) {
		draw(swPhaseText(phase), rect{Left: area.Left, Top: area.Bottom - swLineH, Right: area.Right, Bottom: area.Bottom}, colGreenDm, 0x0020|0x0004|0x8000)
		return
	}
	saved, _, _ := procSaveDC.Call(hdc)
	procIntersectClipRect.Call(hdc, uintptr(area.Left), uintptr(area.Top), uintptr(area.Right), uintptr(area.Bottom))
	bottom := area.Bottom
	for i := len(rows) - 1; i >= 0 && bottom > area.Top; i-- {
		s := strings.TrimSpace(rows[i].Text)
		if s == "" {
			continue
		}
		color := colGreen
		if rows[i].Kind == subtitles.RowTranslation {
			color = colGreenDm
		}
		calc := draw(s, rect{Right: area.Right - area.Left}, color, 0x0400|0x0010)
		height := calc.Bottom - calc.Top
		draw(s, rect{Left: area.Left, Top: bottom - height, Right: area.Right, Bottom: bottom}, color, 0x0010)
		bottom -= height
	}
	procRestoreDC.Call(hdc, saved)
}

func swHasText(rows []subtitles.Row) bool {
	for _, r := range rows {
		if strings.TrimSpace(r.Text) != "" {
			return true
		}
	}
	return false
}

func swSavePlace() {
	wa, screen := monitorRectsForPoint(swX+swW/2, swY+swH/2)
	key := ovplace.ResKey(toPlace(screen))
	f := ovplace.FracOf(swX, swY, swW, swH, toPlace(wa))
	swMu.Lock()
	if swCustom == nil {
		swCustom = map[string]ovplace.Frac{}
	}
	swCustom[key] = f
	cb := swOnMoved
	swMu.Unlock()
	log.Printf("subtitles: the window was moved to %.3f,%.3f of the %s screen", f.X, f.Y, key)
	if cb != nil {
		go cb(key, f)
	}
}

func subWndProc(hwnd, msg, wParam, lParam uintptr) uintptr {
	switch msg {
	case wmSwSet:
		swApplyState(hwnd)
		return 0
	case wmDpiChanged:
		if swShown {
			swRescale(hwnd)
		}
		return 0
	case wmSwLayout:
		if swShown {
			swLayout(hwnd)
			swRenderNow(hwnd)
		}
		return 0
	case wmTimer:
		if !swShown {
			return 0
		}
		swTicks++
		if swTicks%30 == 0 {
			overlayRaise(hwnd)
		}
		if !swDragging && swTranslatedNow() != swLaidTr {
			swLayout(hwnd)
			swRenderNow(hwnd)
			return 0
		}
		x, y := swCursorClient(hwnd)
		hover := swInClose(x, y)
		phase, rows, _ := swContent()
		if hover != swHover || !slices.Equal(rows, swLastRows) || phase != swLastPhase || time.Since(swLastPaint) > time.Second {
			swHover = hover
			swRenderNow(hwnd)
		}
		return 0
	case wmPaint:
		var ps paintStruct
		hdc, _, _ := procBeginPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
		swFrame(hwnd, hdc)
		procEndPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
		return 0
	case wmPrintClient:
		swRender(hwnd, wParam)
		return 0
	case wmMouseActivate:
		return maNoActivate
	case wmLButtonDown:
		x := int32(int16(lParam & 0xFFFF))
		y := int32(int16(lParam >> 16 & 0xFFFF))
		if swInClose(x, y) {
			return 0
		}
		var pt point
		procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
		swDragging = true
		swDragMoved = false
		swDragDX, swDragDY = pt.X-swX, pt.Y-swY
		procSetCapture.Call(hwnd)
		return 0
	case wmMouseMove:
		if swDragging {
			var pt point
			procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
			if swRescale(hwnd) {
				return 0
			}
			nx, ny := pt.X-swDragDX, pt.Y-swDragDY
			if nx != swX || ny != swY {
				swX, swY = nx, ny
				swDragMoved = true
				procSetWindowPos.Call(hwnd, ^uintptr(0), uintptr(nx), uintptr(ny), 0, 0, swpNoSize|swpNoActivate)
			}
		}
		return 0
	case wmLBtnUp:
		x := int32(int16(lParam & 0xFFFF))
		y := int32(int16(lParam >> 16 & 0xFFFF))
		if swDragging {
			swDragging = false
			procReleaseCapture.Call()
			if swDragMoved {
				swSavePlace()
			}
			return 0
		}
		if swInClose(x, y) {
			swMu.Lock()
			cb := swOnClose
			swMu.Unlock()
			if cb != nil {
				log.Printf("subtitles: closed from the window")
				go cb()
			}
		}
		return 0
	case wmCaptureLost:
		if swDragging {
			swDragging = false
			if swDragMoved {
				swSavePlace()
			}
		}
		return 0
	case wmSetCursor:
		x, y := swCursorClient(hwnd)
		id := uintptr(32512)
		if swInClose(x, y) {
			id = 32649
		}
		cur, _, _ := procLoadCursorW.Call(0, id)
		procSetCursor.Call(cur)
		return 1
	case wmChromeRefresh:
		applyChrome(hwnd)
		return 0
	case wmDestroy:
		swBuf.release()
		procPostQuitMessage.Call(0)
		return 0
	}
	r, _, _ := procDefWindowProcW.Call(hwnd, msg, wParam, lParam)
	return r
}

type swDemo struct {
	mu        sync.Mutex
	phase     int
	board     *subtitles.Board
	translate bool
}

func (d *swDemo) windowPhase() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.phase
}

func (d *swDemo) windowRows() []subtitles.Row {
	return d.board.Tail(swTextRunes, d.translate).Rows()
}

func (d *swDemo) windowTranslated() bool { return d.translate }

var swDemoPairs = [][2]string{
	{"Good morning everyone.", "Доброе утро всем."},
	{"Thanks for joining the call today, I know it is early for some of you.", "Спасибо, что подключились, я знаю, что для некоторых из вас это рано."},
	{"Let us start with a quick look at last week: the new build went out on Tuesday, and the feedback so far has been better than we expected, especially from the people who use the export every day.", "Начнём с короткого обзора прошлой недели: новая сборка вышла во вторник, и отзывы пока лучше, чем мы ожидали, особенно от тех, кто каждый день пользуется экспортом."},
	{"Any questions so far?", "Пока есть вопросы?"},
	{"If not, Anna will walk us through the plan for the next two weeks.", "Если нет, Анна расскажет нам о плане на следующие две недели."},
}

const swDemoTrDelay = 900 * time.Millisecond

type swDemoTr struct {
	id   int
	text string
	due  time.Time
}

var swDemoText = []string{
	"Good morning everyone, thanks for joining the call today.",
	"Let us start with a quick look at what happened last week.",
	"The new build went out on Tuesday and the feedback so far is good.",
	"Two people asked about the export, so we will look at it next.",
	"Before we move on, does anyone have a question about the release?",
	"If not, Anna will walk us through the plan for the next two weeks.",
}

func runSubwindowDemo(args []string) {
	seconds := 30
	if len(args) > 0 {
		if n, err := strconv.Atoi(args[0]); err == nil && n > 0 {
			seconds = n
		}
	}
	langCode := "auto"
	if cfg, err := loadConfig("config.json"); err == nil && cfg != nil {
		registerBundledFonts()
		initLang(cfg.UILanguage)
		applyTheme(cfg.Skin, cfg.Theme)
		subWindowConfigure(cfg.SubtitlesLines, cfg.SubtitlesSize, cfg.SubtitlesXY)
		langCode = cfg.SubtitlesLanguage
	}
	translate := len(args) > 1 && strings.TrimSpace(args[1]) != ""
	texts := swDemoText
	if translate {
		texts = make([]string, 0, len(swDemoPairs))
		for _, p := range swDemoPairs {
			texts = append(texts, p[0])
		}
	}
	var pending []swDemoTr
	demo := &swDemo{phase: subPhaseLoading, board: subtitles.NewBoard(subBoardKeep), translate: translate}
	closed := make(chan struct{})
	var once sync.Once
	subWindowShow(demo, langCode, func() { once.Do(func() { close(closed) }) }, func(key string, f ovplace.Frac) {
		log.Printf("subwindow: the place %s %.3f,%.3f is not saved by the demo", key, f.X, f.Y)
	})
	if translate {
		log.Printf("subwindow: showing the subtitles window with translations for %d s", seconds)
	} else {
		log.Printf("subwindow: showing the subtitles window for %d s", seconds)
	}
	end := time.After(time.Duration(seconds) * time.Second)
	tick := time.NewTicker(220 * time.Millisecond)
	defer tick.Stop()
	started := time.Now()
	sentence, word := 0, 0
	var lineStart time.Time
loop:
	for {
		select {
		case <-closed:
			log.Printf("subwindow: closed with the cross")
			break loop
		case <-end:
			break loop
		case <-tick.C:
		}
		age := time.Since(started)
		demo.mu.Lock()
		switch {
		case age < 1500*time.Millisecond:
			demo.phase = subPhaseLoading
		default:
			demo.phase = subPhaseListening
		}
		demo.mu.Unlock()
		if age < 3*time.Second {
			continue
		}
		for len(pending) > 0 && !time.Now().Before(pending[0].due) {
			demo.board.SetTranslation(pending[0].id, pending[0].text)
			pending = pending[1:]
		}
		n := sentence % len(texts)
		words := strings.Fields(texts[n])
		if word == 0 {
			lineStart = time.Now()
		}
		word++
		if word > len(words) {
			ln := demo.board.Add(lineStart, strings.Join(words, " "))
			if translate {
				pending = append(pending, swDemoTr{id: ln.ID, text: swDemoPairs[n][1], due: time.Now().Add(swDemoTrDelay)})
			}
			demo.board.SetLive("")
			sentence++
			word = 0
			continue
		}
		demo.board.SetLive(strings.Join(words[:word], " "))
	}
	subWindowHide(demo)
	time.Sleep(300 * time.Millisecond)
	log.Printf("subwindow: done")
}
