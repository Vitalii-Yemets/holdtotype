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
	"holdtotype/internal/subwin"
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
	swSideGap   = 24
	swGripDIP   = 6
	swAlpha     = 235
	swTextRunes = 900
	swCloseDIP  = 40
)

var swSizes = map[string]int32{"s": 15, "m": 19, "l": 25}

type subWindowSource interface {
	windowPhase() int
	windowRows() []subtitles.Row
	windowHeader() (bool, string, string)
}

type subPart struct {
	rows     func() []subtitles.Row
	area     rect
	dim      bool
	lastRows []subtitles.Row
}

func (p *subPart) fetch() []subtitles.Row {
	if p.rows == nil {
		return nil
	}
	return p.rows()
}

type swHeader struct {
	title string
	codes string
}

func (h swHeader) key() string { return h.title + "\x00" + h.codes }

type subWindow struct {
	class string
	once  sync.Once
	ready chan struct{}

	mu      sync.Mutex
	hwnd    uintptr
	src     subWindowSource
	rows    int32
	size    string
	pos     string
	monitor string
	custom  map[string]ovplace.Frac
	box     map[string]subwin.Box
	onClose func()
	onMoved func(key string, f ovplace.Frac)
	onSized func(key string, b subwin.Box)

	shown     bool
	hover     bool
	dragging  bool
	dragMoved bool
	dragDX    int32
	dragDY    int32
	resizing  bool
	resEdges  int
	resStart  ovplace.Rect
	resPt     point
	resRows   int32
	x, y      int32
	w, h      int32
	lineH     int32
	dpi       int32
	parts     []*subPart
	lastPhase int
	lastHead  string
	lastPaint time.Time
	ticks     int
	buf       backBuf
}

func newSubWindow(class string) *subWindow {
	return &subWindow{
		class:     class,
		ready:     make(chan struct{}),
		rows:      defaultSubtitlesLines,
		size:      "m",
		pos:       ovPosBottom,
		dpi:       96,
		lastPhase: -1,
	}
}

var (
	swMain     = newSubWindow("SubtitlesWnd")
	swByHwnd   sync.Map
	swProcOnce sync.Once
	swProc     uintptr

	swFontMu sync.Mutex
	swFonts  = map[string]uintptr{}
)

func subWindowConfigure(cfg *Config) {
	swMain.configure(cfg)
}

func subWindowShow(src subWindowSource, onClose func(), onMoved func(string, ovplace.Frac), onSized func(string, subwin.Box)) {
	swMain.show(src, onClose, onMoved, onSized)
}

func subWindowHide(src subWindowSource) {
	swMain.hide(src)
}

func (sw *subWindow) configure(cfg *Config) {
	rows := cfg.SubtitlesLines
	if rows < minSubtitlesLines || rows > maxSubtitlesLines {
		rows = defaultSubtitlesLines
	}
	size := cfg.SubtitlesSize
	if !validSubtitlesSize(size) {
		size = "m"
	}
	pos := cfg.SubtitlesPos
	if !subwin.ValidPosition(pos) {
		pos = ovPosBottom
	}
	mon := cfg.SubtitlesMonitor
	if !validOverlayMonitor(mon) {
		mon = ""
	}
	places := map[string]ovplace.Frac{}
	for k, v := range cfg.SubtitlesXY {
		places[k] = v
	}
	boxes := map[string]subwin.Box{}
	for k, v := range cfg.SubtitlesBox {
		boxes[k] = v
	}
	sw.mu.Lock()
	sw.rows = int32(rows)
	sw.size = size
	sw.pos = pos
	sw.monitor = mon
	sw.custom = places
	sw.box = boxes
	hwnd := sw.hwnd
	sw.mu.Unlock()
	if hwnd != 0 {
		procPostMessageW.Call(hwnd, wmSwLayout, 0, 0)
	}
}

func (sw *subWindow) show(src subWindowSource, onClose func(), onMoved func(string, ovplace.Frac), onSized func(string, subwin.Box)) {
	sw.once.Do(sw.startThread)
	<-sw.ready
	sw.mu.Lock()
	sw.src = src
	sw.onClose = onClose
	sw.onMoved = onMoved
	sw.onSized = onSized
	hwnd := sw.hwnd
	sw.mu.Unlock()
	if hwnd != 0 {
		procPostMessageW.Call(hwnd, wmSwSet, 0, 0)
	}
}

func (sw *subWindow) hide(src subWindowSource) {
	sw.mu.Lock()
	if sw.src != src {
		sw.mu.Unlock()
		return
	}
	had := sw.src != nil
	sw.src = nil
	sw.onClose = nil
	sw.onMoved = nil
	sw.onSized = nil
	hwnd := sw.hwnd
	sw.mu.Unlock()
	if had && hwnd != 0 {
		procPostMessageW.Call(hwnd, wmSwSet, 0, 0)
	}
}

func subWndProc(hwnd, msg, wParam, lParam uintptr) uintptr {
	if v, ok := swByHwnd.Load(hwnd); ok {
		return v.(*subWindow).proc(hwnd, msg, wParam, lParam)
	}
	r, _, _ := procDefWindowProcW.Call(hwnd, msg, wParam, lParam)
	return r
}

func (sw *subWindow) startThread() {
	swProcOnce.Do(func() { swProc = syscall.NewCallback(subWndProc) })
	go func() {
		runtime.LockOSThread()
		className, _ := windows.UTF16PtrFromString(appid.Class(sw.class))
		wc := wndClassExW{
			Size:      uint32(unsafe.Sizeof(wndClassExW{})),
			WndProc:   swProc,
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
			swByHwnd.Store(hwnd, sw)
			applyDarkCaption(hwnd)
			procSetLayeredWindowAttributes.Call(hwnd, 0, swAlpha, lwaAlpha)
		} else {
			log.Printf("subtitles: the window was not created")
		}
		sw.mu.Lock()
		sw.hwnd = hwnd
		sw.mu.Unlock()
		close(sw.ready)
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

func (sw *subWindow) textFont(dpi int32) uintptr {
	sw.mu.Lock()
	size := sw.size
	sw.mu.Unlock()
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

func (sw *subWindow) place() (work, screen rect) {
	sw.mu.Lock()
	pick := sw.monitor
	sw.mu.Unlock()
	if pick != "" && pick != "cursor" {
		idx, err := strconv.Atoi(pick)
		mons := listMonitors()
		if err == nil && idx >= 0 && idx < len(mons) {
			return mons[idx].Work, mons[idx].Screen
		}
	}
	if sw.shown && sw.w > 0 {
		return monitorRectsForPoint(sw.x+sw.w/2, sw.y+sw.h/2)
	}
	var pt point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	return monitorRectsForPoint(pt.X, pt.Y)
}

func (sw *subWindow) limits(wa rect, dpi, lineH int32) subwin.Limits {
	px := func(v int32) int32 { return scaleDPI(v, dpi) }
	chrome := px(swHeaderDIP) + px(swPadTop) + px(swPadBottom)
	return subwin.NewLimits(px(subwin.BoxMinW), wa.Right-wa.Left, wa.Bottom-wa.Top, chrome, lineH, 1)
}

func (sw *subWindow) measure(wa, screen rect, dpi int32) (w, h, lineH int32) {
	px := func(v int32) int32 { return scaleDPI(v, dpi) }
	waW := wa.Right - wa.Left
	w = waW * swShare / 100
	if w < px(swMinW) {
		w = px(swMinW)
	}
	if w > px(swMaxW) {
		w = px(swMaxW)
	}
	lineH = swMeasureLine(sw.textFont(dpi))
	if lineH <= 0 {
		lineH = px(24)
	}
	sw.mu.Lock()
	rows := sw.rows
	b, sized := sw.box[ovplace.ResKey(toPlace(screen))]
	sw.mu.Unlock()
	if sized {
		w, rows = px(b.W), b.Rows
	}
	l := sw.limits(wa, dpi, lineH)
	return l.Width(w), l.Height(rows), lineH
}

func (sw *subWindow) setParts() {
	sw.mu.Lock()
	src := sw.src
	sw.mu.Unlock()
	if src == nil {
		sw.parts = nil
		return
	}
	sw.parts = []*subPart{{rows: src.windowRows}}
	sw.layoutParts()
}

func (sw *subWindow) layoutParts() {
	px := func(v int32) int32 { return scaleDPI(v, sw.dpi) }
	area := rect{Left: px(swPadX), Top: px(swHeaderDIP) + px(swPadTop), Right: sw.w - px(swPadX), Bottom: sw.h - px(swPadBottom)}
	for _, p := range sw.parts {
		p.area = area
	}
}

func (sw *subWindow) setRect(hwnd uintptr, x, y, w, h, lineH, dpi int32) {
	sw.x, sw.y, sw.w, sw.h, sw.lineH, sw.dpi = x, y, w, h, lineH, dpi
	sw.layoutParts()
	procSetWindowPos.Call(hwnd, ^uintptr(0), uintptr(x), uintptr(y), uintptr(w), uintptr(h), swpNoActivate)
	procSetWindowRgn.Call(hwnd, swRegion(w, h, dpi), 1)
}

func (sw *subWindow) layout(hwnd uintptr) {
	wa, screen := sw.place()
	dpi := dpiForPoint(wa.Left+(wa.Right-wa.Left)/2, wa.Top+(wa.Bottom-wa.Top)/2)
	w, h, lineH := sw.measure(wa, screen, dpi)
	sw.mu.Lock()
	custom, pos := sw.custom, sw.pos
	sw.mu.Unlock()
	var x, y int32
	if f, ok := custom[ovplace.ResKey(toPlace(screen))]; ok && pos == ovPosCustom {
		x, y = ovplace.Custom(f, toPlace(wa), w, h)
	} else {
		x, y = subwin.Place(pos, toPlace(wa), w, h, scaleDPI(swSideGap, dpi), scaleDPI(swBottomGap, dpi))
	}
	sw.setRect(hwnd, x, y, w, h, lineH, dpi)
}

func (sw *subWindow) rescale(hwnd uintptr) bool {
	var pt point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	cx, cy := pt.X, pt.Y
	if !sw.dragging {
		cx, cy = sw.x+sw.w/2, sw.y+sw.h/2
	}
	dpi := dpiForPoint(cx, cy)
	if dpi == sw.dpi || sw.w <= 0 || sw.h <= 0 {
		return false
	}
	wa, screen := monitorRectsForPoint(cx, cy)
	w, h, lineH := sw.measure(wa, screen, dpi)
	var x, y int32
	if sw.dragging {
		sw.dragDX = int32(int64(sw.dragDX) * int64(w) / int64(sw.w))
		sw.dragDY = int32(int64(sw.dragDY) * int64(h) / int64(sw.h))
		x, y = pt.X-sw.dragDX, pt.Y-sw.dragDY
	} else {
		x, y = cx-w/2, cy-h/2
	}
	log.Printf("subtitles: the window is on a screen at %d DPI now, %dx%d", dpi, w, h)
	sw.setRect(hwnd, x, y, w, h, lineH, dpi)
	sw.renderNow(hwnd)
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

func (sw *subWindow) inClose(x, y int32) bool {
	px := func(v int32) int32 { return scaleDPI(v, sw.dpi) }
	return x >= sw.w-px(swCloseDIP) && x < sw.w && y >= 0 && y < px(swHeaderDIP)
}

func swCursorClient(hwnd uintptr) (int32, int32) {
	var pt point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	procScreenToClient.Call(hwnd, uintptr(unsafe.Pointer(&pt)))
	return pt.X, pt.Y
}

func (sw *subWindow) content() (int, swHeader, bool) {
	sw.mu.Lock()
	src := sw.src
	sw.mu.Unlock()
	if src == nil {
		return 0, swHeader{}, false
	}
	translated, lang, to := src.windowHeader()
	head := swHeader{title: tr(subwin.TitleKey(translated)), codes: subwin.Codes(lang, to, translated)}
	return src.windowPhase(), head, true
}

func (sw *subWindow) applyState(hwnd uintptr) {
	_, _, on := sw.content()
	if !on {
		if sw.dragging {
			procReleaseCapture.Call()
			sw.dragging = false
		}
		procKillTimer.Call(hwnd, swTimerID)
		procShowWindow.Call(hwnd, swHide)
		sw.shown = false
		sw.parts = nil
		sw.lastPhase, sw.lastHead = -1, ""
		return
	}
	sw.setParts()
	if !sw.shown {
		sw.layout(hwnd)
	}
	sw.shown = true
	procSetTimer.Call(hwnd, swTimerID, 100, 0)
	procShowWindow.Call(hwnd, swShowNoActivate)
	overlayRaise(hwnd)
	sw.renderNow(hwnd)
}

func (sw *subWindow) renderNow(hwnd uintptr) {
	hdc, _, _ := procGetDC.Call(hwnd)
	if hdc == 0 {
		return
	}
	sw.frame(hwnd, hdc)
	procReleaseDC.Call(hwnd, hdc)
}

func (sw *subWindow) frame(hwnd, hdc uintptr) {
	var rc rect
	procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
	mem := sw.buf.begin(hdc, rc.Right-rc.Left, rc.Bottom-rc.Top)
	if mem == 0 {
		sw.render(hwnd, hdc)
		return
	}
	sw.render(hwnd, mem)
	sw.buf.blit(hdc)
	sw.lastPaint = time.Now()
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

func (sw *subWindow) render(hwnd, hdc uintptr) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("subtitles window: panic: %v", r)
		}
	}()
	var rc rect
	procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
	dpi := sw.dpi
	px := func(v int32) int32 { return scaleDPI(v, dpi) }
	phase, hdr, _ := sw.content()
	sw.lastPhase, sw.lastHead = phase, hdr.key()
	partRows := make([][]subtitles.Row, len(sw.parts))
	hasText := false
	for i, p := range sw.parts {
		partRows[i] = p.fetch()
		p.lastRows = partRows[i]
		if swHasText(partRows[i]) {
			hasText = true
		}
	}

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
	title := hdr.title
	measure := draw(title, rect{}, colGreen, 0x0400|0x0020)
	tr0 := rect{Left: px(32), Top: 0, Right: px(32) + (measure.Right - measure.Left), Bottom: head}
	draw(title, tr0, colGreen, 0x0020|0x0004)
	if code := hdr.codes; code != "" {
		draw(code, rect{Left: tr0.Right + px(10), Top: 0, Right: rc.Right - px(swCloseDIP), Bottom: head}, colGreenDm, 0x0020|0x0004|0x8000)
	}
	xcol := colGreenDm
	if sw.hover {
		xcol = colGreen
	}
	draw("✕", rect{Left: rc.Right - px(swCloseDIP), Top: 0, Right: rc.Right - px(6), Bottom: head}, xcol, 0x0020|0x0004|0x0001)
	fill(rect{Left: px(10), Top: head, Right: rc.Right - px(10), Bottom: head + 1}, colLine)

	font := sw.textFont(dpi)
	procSelectObject.Call(hdc, font)
	if !hasText {
		area := rect{Left: px(swPadX), Top: head + px(swPadTop), Right: rc.Right - px(swPadX), Bottom: rc.Bottom - px(swPadBottom)}
		draw(swPhaseText(phase), rect{Left: area.Left, Top: area.Bottom - sw.lineH, Right: area.Right, Bottom: area.Bottom}, colGreenDm, 0x0020|0x0004|0x8000)
		return
	}
	for i, p := range sw.parts {
		sw.drawPart(hdc, p, partRows[i], draw)
	}
}

func (sw *subWindow) drawPart(hdc uintptr, p *subPart, rows []subtitles.Row, draw func(string, rect, uintptr, uintptr) rect) {
	area := p.area
	saved, _, _ := procSaveDC.Call(hdc)
	procIntersectClipRect.Call(hdc, uintptr(area.Left), uintptr(area.Top), uintptr(area.Right), uintptr(area.Bottom))
	bottom := area.Bottom
	for i := len(rows) - 1; i >= 0 && bottom > area.Top; i-- {
		s := strings.TrimSpace(rows[i].Text)
		if s == "" {
			continue
		}
		color := colGreen
		if p.dim || rows[i].Kind == subtitles.RowHeard {
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

func (sw *subWindow) savePlace() {
	wa, screen := monitorRectsForPoint(sw.x+sw.w/2, sw.y+sw.h/2)
	key := ovplace.ResKey(toPlace(screen))
	f := ovplace.FracOf(sw.x, sw.y, sw.w, sw.h, toPlace(wa))
	sw.mu.Lock()
	if sw.custom == nil {
		sw.custom = map[string]ovplace.Frac{}
	}
	sw.custom[key] = f
	sw.pos = ovPosCustom
	cb := sw.onMoved
	sw.mu.Unlock()
	log.Printf("subtitles: the window was moved to %.3f,%.3f of the %s screen", f.X, f.Y, key)
	if cb != nil {
		go cb(key, f)
	}
}

func (sw *subWindow) resizeTo(hwnd uintptr) {
	var pt point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	wa, _ := monitorRectsForPoint(sw.x+sw.w/2, sw.y+sw.h/2)
	sw.mu.Lock()
	pos := sw.pos
	sw.mu.Unlock()
	l := sw.limits(wa, sw.dpi, sw.lineH)
	r, rows := subwin.Resize(sw.resStart, sw.resEdges, pt.X-sw.resPt.X, pt.Y-sw.resPt.Y, pos, l)
	w, h := r.Right-r.Left, r.Bottom-r.Top
	x, y := r.Left, r.Top
	if pos != ovPosCustom {
		x, y = subwin.Place(pos, toPlace(wa), w, h, scaleDPI(swSideGap, sw.dpi), scaleDPI(swBottomGap, sw.dpi))
	}
	sw.resRows = rows
	if x == sw.x && y == sw.y && w == sw.w && h == sw.h {
		return
	}
	sw.setRect(hwnd, x, y, w, h, sw.lineH, sw.dpi)
	sw.renderNow(hwnd)
}

func (sw *subWindow) endResize() {
	if !sw.resizing {
		return
	}
	sw.resizing = false
	_, screen := monitorRectsForPoint(sw.x+sw.w/2, sw.y+sw.h/2)
	key := ovplace.ResKey(toPlace(screen))
	b := subwin.Box{W: sw.w * 96 / sw.dpi, Rows: sw.resRows}
	if b.W < subwin.BoxMinW {
		b.W = subwin.BoxMinW
	}
	sw.mu.Lock()
	if sw.box == nil {
		sw.box = map[string]subwin.Box{}
	}
	sw.box[key] = b
	cb := sw.onSized
	pos := sw.pos
	sw.mu.Unlock()
	log.Printf("subtitles: the window was resized to %d dip and %d rows on the %s screen", b.W, b.Rows, key)
	if cb != nil {
		go cb(key, b)
	}
	if pos == ovPosCustom {
		sw.savePlace()
	}
}

func swEdgeCursor(edges int) uintptr {
	switch {
	case edges == subwin.EdgeLeft|subwin.EdgeTop, edges == subwin.EdgeRight|subwin.EdgeBottom:
		return 32642
	case edges == subwin.EdgeRight|subwin.EdgeTop, edges == subwin.EdgeLeft|subwin.EdgeBottom:
		return 32643
	case edges&(subwin.EdgeLeft|subwin.EdgeRight) != 0:
		return 32644
	case edges != 0:
		return 32645
	}
	return 0
}

func (sw *subWindow) edgesAt(x, y int32) int {
	return subwin.Hit(x, y, sw.w, sw.h, scaleDPI(swGripDIP, sw.dpi))
}

func (sw *subWindow) changed() bool {
	phase, hdr, _ := sw.content()
	if phase != sw.lastPhase || hdr.key() != sw.lastHead {
		return true
	}
	for _, p := range sw.parts {
		if !slices.Equal(p.fetch(), p.lastRows) {
			return true
		}
	}
	return false
}

func (sw *subWindow) proc(hwnd, msg, wParam, lParam uintptr) uintptr {
	switch msg {
	case wmSwSet:
		sw.applyState(hwnd)
		return 0
	case wmDpiChanged:
		if sw.shown {
			sw.rescale(hwnd)
		}
		return 0
	case wmSwLayout:
		if sw.shown {
			sw.layout(hwnd)
			sw.renderNow(hwnd)
		}
		return 0
	case wmTimer:
		if !sw.shown {
			return 0
		}
		sw.ticks++
		if sw.ticks%30 == 0 {
			overlayRaise(hwnd)
		}
		x, y := swCursorClient(hwnd)
		hover := sw.inClose(x, y)
		if hover != sw.hover || sw.changed() || time.Since(sw.lastPaint) > time.Second {
			sw.hover = hover
			sw.renderNow(hwnd)
		}
		return 0
	case wmPaint:
		var ps paintStruct
		hdc, _, _ := procBeginPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
		sw.frame(hwnd, hdc)
		procEndPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
		return 0
	case wmPrintClient:
		sw.render(hwnd, wParam)
		return 0
	case wmMouseActivate:
		return maNoActivate
	case wmLButtonDown:
		x := int32(int16(lParam & 0xFFFF))
		y := int32(int16(lParam >> 16 & 0xFFFF))
		var pt point
		procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
		if edges := sw.edgesAt(x, y); edges != 0 {
			sw.resizing = true
			sw.resEdges = edges
			sw.resStart = ovplace.Rect{Left: sw.x, Top: sw.y, Right: sw.x + sw.w, Bottom: sw.y + sw.h}
			sw.resPt = pt
			wa, _ := monitorRectsForPoint(sw.x+sw.w/2, sw.y+sw.h/2)
			sw.resRows = sw.limits(wa, sw.dpi, sw.lineH).RowsFor(sw.h)
			procSetCapture.Call(hwnd)
			return 0
		}
		if sw.inClose(x, y) {
			return 0
		}
		sw.dragging = true
		sw.dragMoved = false
		sw.dragDX, sw.dragDY = pt.X-sw.x, pt.Y-sw.y
		procSetCapture.Call(hwnd)
		return 0
	case wmMouseMove:
		if sw.resizing {
			sw.resizeTo(hwnd)
			return 0
		}
		if sw.dragging {
			var pt point
			procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
			if sw.rescale(hwnd) {
				return 0
			}
			nx, ny := pt.X-sw.dragDX, pt.Y-sw.dragDY
			if nx != sw.x || ny != sw.y {
				sw.x, sw.y = nx, ny
				sw.dragMoved = true
				procSetWindowPos.Call(hwnd, ^uintptr(0), uintptr(nx), uintptr(ny), 0, 0, swpNoSize|swpNoActivate)
			}
		}
		return 0
	case wmLBtnUp:
		x := int32(int16(lParam & 0xFFFF))
		y := int32(int16(lParam >> 16 & 0xFFFF))
		if sw.resizing {
			procReleaseCapture.Call()
			sw.endResize()
			return 0
		}
		if sw.dragging {
			sw.dragging = false
			procReleaseCapture.Call()
			if sw.dragMoved {
				sw.savePlace()
			}
			return 0
		}
		if sw.inClose(x, y) {
			sw.mu.Lock()
			cb := sw.onClose
			sw.mu.Unlock()
			if cb != nil {
				log.Printf("subtitles: closed from the window")
				go cb()
			}
		}
		return 0
	case wmCaptureLost:
		sw.endResize()
		if sw.dragging {
			sw.dragging = false
			if sw.dragMoved {
				sw.savePlace()
			}
		}
		return 0
	case wmSetCursor:
		x, y := swCursorClient(hwnd)
		id := uintptr(32512)
		if c := swEdgeCursor(sw.edgesAt(x, y)); c != 0 {
			id = c
		} else if sw.inClose(x, y) {
			id = 32649
		}
		cur, _, _ := procLoadCursorW.Call(0, id)
		procSetCursor.Call(cur)
		return 1
	case wmChromeRefresh:
		applyChrome(hwnd)
		return 0
	case wmDestroy:
		sw.buf.release()
		swByHwnd.Delete(hwnd)
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
	lang      string
	to        string
}

func (d *swDemo) windowPhase() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.phase
}

func (d *swDemo) windowHeader() (bool, string, string) {
	return d.translate, d.lang, d.to
}

func (d *swDemo) windowRows() []subtitles.Row {
	if d.translate {
		return d.board.TranslationRows(swTextRunes)
	}
	return d.board.Tail(swTextRunes).Rows()
}

type swDemoPiece struct {
	source string
	text   string
	para   bool
	fail   bool
}

var swDemoPieces = []swDemoPiece{
	{"Good morning everyone,", "Доброе утро всем,", true, false},
	{"thanks for joining the call today, I know it is early for some of you.", "спасибо, что подключились, я знаю, что для некоторых из вас это рано.", false, false},
	{"Let us start with a quick look at last week:", "Начнём с короткого обзора прошлой недели:", true, false},
	{"the new build went out on Tuesday,", "новая сборка вышла во вторник,", false, false},
	{"and the feedback so far has been better than we expected,", "и отзывы пока лучше, чем мы ожидали,", false, false},
	{"especially from the people who use the export every day.", "особенно от тех, кто каждый день пользуется экспортом.", false, false},
	{"Any questions so far?", "", true, true},
	{"If not, Anna will walk us through the plan for the next two weeks.", "Если нет, Анна расскажет нам о плане на следующие две недели.", true, false},
}

func swDemoScript(d *swDemo) []func() {
	var steps []func()
	idle := func(n int) {
		for i := 0; i < n; i++ {
			steps = append(steps, func() {})
		}
	}
	hearing := func(on bool) { steps = append(steps, func() { d.board.SetHearing(on) }) }
	for round := 0; round < 50; round++ {
		for i, pc := range swDemoPieces {
			seq := round*len(swDemoPieces) + i + 1
			if pc.para {
				hearing(false)
				idle(12)
			}
			hearing(true)
			idle(len(strings.Fields(pc.source)) * 2)
			steps = append(steps, func() { d.board.AddPiece(seq, pc.para, pc.source) })
			idle(1)
			steps = append(steps, func() { d.board.SetPiece(seq, "", subtitles.PieceRunning) })
			idle(2)
			if pc.fail {
				steps = append(steps, func() { d.board.SetPiece(seq, "", subtitles.PieceFailed) })
				continue
			}
			words := strings.Fields(pc.text)
			for n := 2; n < len(words); n += 2 {
				shown := strings.Join(words[:n], " ")
				steps = append(steps, func() { d.board.SetPiece(seq, shown, subtitles.PieceRunning) })
			}
			steps = append(steps, func() { d.board.SetPiece(seq, pc.text, subtitles.PieceDone) })
		}
	}
	return steps
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
	langCode, to := "auto", "ru"
	if cfg, err := loadConfig("config.json"); err == nil && cfg != nil {
		registerBundledFonts()
		initLang(cfg.UILanguage)
		applyTheme(cfg.Skin, cfg.Theme)
		subWindowConfigure(cfg)
		langCode = cfg.SubtitlesLanguage
		if cfg.SubtitlesTranslate != "" {
			to = cfg.SubtitlesTranslate
		}
	}
	translate := len(args) > 1 && strings.TrimSpace(args[1]) != ""
	demo := &swDemo{phase: subPhaseLoading, board: subtitles.NewBoard(subBoardKeep), translate: translate, lang: langCode, to: to}
	var script []func()
	if translate {
		script = swDemoScript(demo)
	}
	closed := make(chan struct{})
	var once sync.Once
	subWindowShow(demo, func() { once.Do(func() { close(closed) }) }, func(key string, f ovplace.Frac) {
		log.Printf("subwindow: the place %s %.3f,%.3f is not saved by the demo", key, f.X, f.Y)
	}, func(key string, b subwin.Box) {
		log.Printf("subwindow: the size %s %d dip, %d rows is not saved by the demo", key, b.W, b.Rows)
	})
	if translate {
		log.Printf("subwindow: showing the subtitles window with translation for %d s", seconds)
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
		if translate {
			if len(script) > 0 {
				script[0]()
				script = script[1:]
			}
			continue
		}
		words := strings.Fields(swDemoText[sentence%len(swDemoText)])
		if word == 0 {
			lineStart = time.Now()
		}
		word++
		if word > len(words) {
			demo.board.Add(lineStart, strings.Join(words, " "))
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
