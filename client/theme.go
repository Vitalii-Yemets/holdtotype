package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"strings"
	"sync/atomic"

	"holdtotype/internal/theme"
)

var currentTheme atomic.Value

func themeID() string {
	if v, ok := currentTheme.Load().(string); ok && theme.Valid(v) {
		return v
	}
	return theme.Default
}

func themeSkin() theme.Skin { return theme.Get(themeID()) }

func themeRoundCorners() bool { return themeSkin().Round }

func themeGlow() bool { return themeSkin().Glow }

func themeScanlines() bool { return themeSkin().Scan > 0 }

func themeLevelStyle() string { return themeSkin().Level }

func themePulse() float64 { return themeSkin().Pulse }

func themeCSSVars() string { return themeSkin().CSSVars() }

// colorref turns a #rrggbb string into the BGR value GDI wants.
func colorref(hex string) uintptr {
	r, g, b := theme.RGB(hex)
	return uintptr(b)<<16 | uintptr(g)<<8 | uintptr(r)
}

func mixHex(hex string, t float64) uintptr {
	r, g, b := theme.RGB(hex)
	f := func(v uint8) uintptr { return uintptr(float64(v) * t) }
	return f(b)<<16 | f(g)<<8 | f(r)
}

func applyTheme(id string) {
	if !theme.Valid(id) {
		id = theme.Default
	}
	if themeID() == id && currentTheme.Load() != nil {
		return
	}
	currentTheme.Store(id)
	p := theme.Get(id)
	dropFontCache()

	colBg = colorref(p.Bg)
	colBgLine = mixHex(p.Bg, 0.75)
	colGreen = colorref(p.Accent)
	colGreenDm = colorref(p.Dim)
	colGreenLo = colorref(p.Faint)
	colAmber = colorref(p.Warn)
	colAmberDm = mixHex(p.Warn, 0.35)
	colBad = colorref(p.Bad)
	colBadDm = mixHex(p.Bad, 0.35)
	colRed = colBad
	colRedDm = colBadDm
	colAskBg = mixHex(p.Panel, 1.0)

	rebuildIcons(p)
	log.Printf("оформление: %s", id)
}

func rebuildIcons(p theme.Skin) {
	ar, ag, ab := theme.RGB(p.Accent)
	br, bg, bb := theme.RGB(p.Bad)
	wr, wg, wb := theme.RGB(p.Warn)
	accent := color.NRGBA{R: ar, G: ag, B: ab, A: 255}
	bad := color.NRGBA{R: br, G: bg, B: bb, A: 255}
	warn := color.NRGBA{R: wr, G: wg, B: wb, A: 255}
	off := color.NRGBA{R: 0x5A, G: 0x6E, B: 0x60, A: 255}
	iconIdle = iconPNG(accent)
	iconRecording = iconPNG(bad)
	iconProcessing = iconPNG(warn)
	iconDisabled = iconPNG(off)
	iconError = iconPNG(off, bad)
}

func themeListJSON() string {
	type entry struct {
		Bg      string `json:"bg"`
		Panel   string `json:"panel"`
		Line    string `json:"line"`
		Accent  string `json:"accent"`
		Dim     string `json:"dim"`
		Faint   string `json:"faint"`
		Warn    string `json:"warn"`
		Bad     string `json:"bad"`
		RGB     string `json:"rgb"`
		Glow    string `json:"glow"`
		Font    string `json:"font"`
		Radius  string `json:"r"`
		Border  string `json:"bw"`
		Scan    string `json:"scan"`
		Shadow  string `json:"shadow"`
		WBorder string `json:"wborder"`
		BarR    string `json:"barr"`
	}
	out := map[string]entry{}
	for _, id := range theme.IDs() {
		p := theme.Get(id)
		r, g, b := theme.RGB(p.Accent)
		out[id] = entry{
			Bg: p.Bg, Panel: p.Panel, Line: p.Line, Accent: p.Accent,
			Dim: p.Dim, Faint: p.Faint, Warn: p.Warn, Bad: p.Bad,
			RGB:     fmt.Sprintf("%d,%d,%d", r, g, b),
			Glow:    cssVar(p.CSSVars(), "--glow"),
			Font:    cssVar(p.CSSVars(), "--font"),
			Radius:  cssVar(p.CSSVars(), "--r"),
			Border:  cssVar(p.CSSVars(), "--bw"),
			Scan:    cssVar(p.CSSVars(), "--scan"),
			Shadow:  cssVar(p.CSSVars(), "--shadow"),
			WBorder: cssVar(p.CSSVars(), "--wborder"),
			BarR:    barRadius(p),
		}
	}
	data, err := json.Marshal(out)
	if err != nil {
		return "{}"
	}
	return string(data)
}

// refreshWindowChrome repaints the frame of every window that is open right now.
func refreshWindowChrome() {
	for _, h := range liveWindows() {
		if h != 0 {
			applyDarkCaption(h)
			procSetWindowPos.Call(h, 0, 0, 0, 0, 0, 0x0002|0x0001|0x0004|0x0020)
			procRedrawWindow.Call(h, 0, 0, 0x0001|0x0004|0x0100)
		}
	}
}

func liveWindows() []uintptr {
	out := []uintptr{settingsHwnd.Load(), overlayHwnd()}
	trayMu.Lock()
	out = append(out, trayHwnd)
	trayMu.Unlock()
	capMu.Lock()
	out = append(out, capHwnd)
	capMu.Unlock()
	return out
}

// cssVar pulls one value out of the string CSSVars renders.
func cssVar(vars, name string) string {
	for _, part := range strings.Split(vars, ";") {
		if strings.HasPrefix(part, name+":") {
			return strings.TrimPrefix(part, name+":")
		}
	}
	return ""
}

func barRadius(s theme.Skin) string {
	if s.Radius >= 10 {
		return "99px"
	}
	return "0"
}
