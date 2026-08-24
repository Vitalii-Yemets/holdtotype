package main

import (
	"log"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	ovAskH   = 46
	ovBtnH   = 26
	ovBtnGap = 8
	ovBtnPad = 13
	ovBtnMin = 46
	ovAskPad = 20
)

type ovChoice struct {
	id    string
	label string
	def   bool
}

var (
	askMu      sync.Mutex
	askOn      bool
	askPrompt  string
	askItems   []ovChoice
	askWidths  []int32
	askPromptW int32
	askEnd     time.Time
	askTotal   time.Duration
	askCh      chan string
	askForced  bool
)

func askActive() bool {
	askMu.Lock()
	defer askMu.Unlock()
	return askOn
}

func askDefault(choices []ovChoice) string {
	for _, c := range choices {
		if c.def {
			return c.id
		}
	}
	if len(choices) > 0 {
		return choices[0].id
	}
	return ""
}

func textWidthDIP(s string) int32 {
	if s == "" {
		return 0
	}
	hdc, _, _ := procGetDC.Call(0)
	if hdc == 0 {
		return int32(len([]rune(s))) * 8
	}
	defer procReleaseDC.Call(0, hdc)
	old, _, _ := procSelectObject.Call(hdc, overlayFont())
	u, err := windows.UTF16FromString(s)
	if err != nil {
		return int32(len([]rune(s))) * 8
	}
	r := rect{}
	procDrawTextW.Call(hdc, uintptr(unsafe.Pointer(&u[0])), uintptr(len(u)-1),
		uintptr(unsafe.Pointer(&r)), 0x0400|0x0020)
	if old != 0 {
		procSelectObject.Call(hdc, old)
	}
	dpi := dpiFor(overlayHwnd())
	if dpi < 72 {
		dpi = 96
	}
	return (r.Right - r.Left) * 96 / dpi
}

func askMaxWidthDIP() int32 {
	wa := overlayWorkArea()
	dpi := dpiFor(overlayHwnd())
	if dpi < 72 {
		dpi = 96
	}
	max := (wa.Right - wa.Left) * 4 / 5 * 96 / dpi
	if max < ovW {
		max = ovW
	}
	return max
}

func askLayoutLocked() (rows int32, width int32, pos []point) {
	if !askOn {
		return 0, 0, nil
	}
	left := int32(ovAskPad) + askPromptW + 16
	max := askMaxWidthDIP() - ovAskPad
	x, row, widest := left, int32(0), left
	pos = make([]point, 0, len(askWidths))
	for _, bw := range askWidths {
		if x > left && x+bw > max {
			row++
			x = left
		}
		pos = append(pos, point{X: x, Y: row})
		x += bw + ovBtnGap
		if x-ovBtnGap > widest {
			widest = x - ovBtnGap
		}
	}
	return row + 1, widest + ovAskPad, pos
}

func askRows() int32 {
	askMu.Lock()
	defer askMu.Unlock()
	rows, _, _ := askLayoutLocked()
	if rows < 1 {
		rows = 1
	}
	return rows
}

func askWidthDIP() int32 {
	askMu.Lock()
	defer askMu.Unlock()
	_, w, _ := askLayoutLocked()
	return w
}

func askButtonRects(dpi int32) []rect {
	askMu.Lock()
	defer askMu.Unlock()
	_, _, pos := askLayoutLocked()
	if pos == nil {
		return nil
	}
	out := make([]rect, 0, len(pos))
	for i, p := range pos {
		top := int32(ovH) + p.Y*ovAskH + (ovAskH-ovBtnH)/2
		out = append(out, rect{
			Left:   scaleDPI(p.X, dpi),
			Top:    scaleDPI(top, dpi),
			Right:  scaleDPI(p.X+askWidths[i], dpi),
			Bottom: scaleDPI(top+ovBtnH, dpi),
		})
	}
	return out
}

func askHit(x, y int32) (string, bool) {
	askMu.Lock()
	on := askOn
	items := append([]ovChoice(nil), askItems...)
	askMu.Unlock()
	if !on {
		return "", false
	}
	for i, r := range askButtonRects(dpiFor(overlayHwnd())) {
		if x >= r.Left && x <= r.Right && y >= r.Top && y <= r.Bottom && i < len(items) {
			return items[i].id, true
		}
	}
	return "", false
}

func askFinish(id string) {
	askMu.Lock()
	ch := askCh
	askCh = nil
	askMu.Unlock()
	if ch != nil {
		ch <- id
	}
}

func askKey(vk uint32) bool {
	askMu.Lock()
	on := askOn
	items := append([]ovChoice(nil), askItems...)
	askMu.Unlock()
	if !on {
		return false
	}
	switch {
	case vk == vkReturn:
		askFinish(askDefault(items))
		return true
	case vk >= '1' && vk <= '9':
		i := int(vk - '1')
		if i < len(items) {
			askFinish(items[i].id)
			return true
		}
	}
	return false
}

func overlayAsk(prompt string, choices []ovChoice, seconds int) string {
	if len(choices) == 0 {
		return ""
	}
	ovOnce.Do(startOverlayThread)
	<-ovReady

	ch := make(chan string, 1)
	widths := make([]int32, 0, len(choices))
	for _, c := range choices {
		w := textWidthDIP(c.label) + 2*ovBtnPad
		if w < ovBtnMin {
			w = ovBtnMin
		}
		widths = append(widths, w)
	}
	promptW := textWidthDIP(prompt)

	ovMu.Lock()
	hidden := ovState == ovHidden
	ovMu.Unlock()

	askMu.Lock()
	if askOn {
		askMu.Unlock()
		return askDefault(choices)
	}
	askOn = true
	askPrompt = prompt
	askItems = append([]ovChoice(nil), choices...)
	askWidths = widths
	askPromptW = promptW
	askCh = ch
	askForced = hidden
	askTotal = time.Duration(seconds) * time.Second
	if seconds > 0 {
		askEnd = time.Now().Add(askTotal)
	} else {
		askEnd = time.Time{}
	}
	askMu.Unlock()

	if hidden {
		overlaySet(ovProcessing, tr("ov.transcribing"))
	} else {
		overlayRefresh()
	}

	res := <-ch

	askMu.Lock()
	askOn = false
	askItems = nil
	askWidths = nil
	askPrompt = ""
	forced := askForced
	askMu.Unlock()

	if forced {
		overlayHide()
	} else {
		overlayRefresh()
	}
	log.Printf("вопрос %q → %q", prompt, res)
	return res
}

func askTick() {
	askMu.Lock()
	on := askOn
	end := askEnd
	items := append([]ovChoice(nil), askItems...)
	askMu.Unlock()
	if !on || end.IsZero() || time.Now().Before(end) {
		return
	}
	askFinish(askDefault(items))
}

func askLeft() (frac float64, live bool) {
	askMu.Lock()
	defer askMu.Unlock()
	if !askOn || askEnd.IsZero() || askTotal <= 0 {
		return 0, false
	}
	left := time.Until(askEnd)
	if left < 0 {
		left = 0
	}
	return float64(left) / float64(askTotal), true
}

func askRender(hwnd, hdc uintptr, rc rect, fill func(rect, uintptr), drawText func(string, rect, uintptr, uintptr)) {
	askMu.Lock()
	on := askOn
	prompt := askPrompt
	items := append([]ovChoice(nil), askItems...)
	askMu.Unlock()
	if !on {
		return
	}
	dpi := dpiFor(hwnd)
	px := func(v int32) int32 { return scaleDPI(v, dpi) }

	fill(rect{Left: 0, Top: px(ovH), Right: rc.Right, Bottom: px(ovH) + 1}, colGreenLo)
	drawText(prompt, rect{
		Left:   px(ovAskPad),
		Top:    px(ovH),
		Right:  px(ovAskPad) + px(askPromptW),
		Bottom: px(ovH + ovAskH),
	}, colGreenDm, 0x0020|0x0004)

	frac, live := askLeft()
	for i, r := range askButtonRects(dpi) {
		if i >= len(items) {
			break
		}
		border, txt := uintptr(colGreenDm), uintptr(colGreenDm)
		if items[i].def {
			border, txt = colGreen, colGreen
		}
		fill(rect{Left: r.Left - 1, Top: r.Top - 1, Right: r.Right + 1, Bottom: r.Bottom + 1}, border)
		fill(r, 0x0B100D)
		drawText(items[i].label, r, txt, 0x0020|0x0004|0x0001)
		if live && items[i].def {
			w := int32(float64(r.Right-r.Left) * frac)
			fill(rect{Left: r.Left, Top: r.Bottom + px(3), Right: r.Left + w, Bottom: r.Bottom + px(5)}, colGreen)
		}
	}
}

func askTranslateTarget(cfg *Config) string {
	choices := make([]ovChoice, 0, len(cfg.TranslateAskLangs)+1)
	for _, l := range cfg.TranslateAskLangs {
		choices = append(choices, ovChoice{id: l, label: langLabel(l), def: l == cfg.TranslateTarget})
	}
	choices = append(choices, ovChoice{id: "", label: tr("td.plain")})
	seconds := 0
	if cfg.TranslateAsk == "timeout" {
		seconds = cfg.TranslateAskSeconds
	}
	return overlayAsk(tr("td.title"), choices, seconds)
}

func askFocusMismatch() string {
	return overlayAsk(tr("fd.title"), []ovChoice{
		{id: "here", label: tr("fd.here"), def: true},
		{id: "copy", label: tr("fd.copy")},
	}, 30)
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

func askAbort() {
	if askActive() {
		askFinish("")
	}
}
