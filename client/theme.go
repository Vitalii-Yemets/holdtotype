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

var (
	currentSkin   atomic.Value
	currentColour atomic.Value
)

func skinID() string {
	if v, ok := currentSkin.Load().(string); ok && theme.ValidSkin(v) {
		return v
	}
	return theme.DefaultSkin
}

func colourID() string {
	if v, ok := currentColour.Load().(string); ok && theme.ValidColour(v) {
		return v
	}
	return theme.DefaultPalette
}

func themeLook() theme.Look { return theme.Current(skinID(), colourID()) }

func themeRoundCorners() bool { return themeLook().Round }

func themeGlow() bool { return themeLook().Glow }

func themeScanlines() bool { return themeLook().Scan > 0 }

func themeLevelStyle() string { return themeLook().Level }

func themePulse() float64 { return themeLook().Pulse }

func themeCSSVars() string { return themeLook().CSSVars() }

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

func applyTheme(skin, colour string) {
	if !theme.ValidSkin(skin) {
		skin = theme.DefaultSkin
	}
	if !theme.ValidColour(colour) {
		colour = theme.DefaultPalette
	}
	if skinID() == skin && colourID() == colour && currentSkin.Load() != nil {
		return
	}
	currentSkin.Store(skin)
	currentColour.Store(colour)
	dropFontCache()

	p := theme.Current(skin, colour).Palette
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
	log.Printf("оформление: %s, цвет %s", skin, p.ID)
}

func rebuildIcons(p theme.Palette) {
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

// skinListJSON hands the page every skin and every colour, so it can repaint
// itself the moment one is picked.
func skinListJSON() string {
	type entry struct {
		Skin    string `json:"skin"`
		Colour  string `json:"colour"`
		Bg      string `json:"bg"`
		Panel   string `json:"panel"`
		Line    string `json:"line"`
		Accent  string `json:"accent"`
		Dim     string `json:"dim"`
		Faint   string `json:"faint"`
		Warn    string `json:"warn"`
		Bad     string `json:"bad"`
		Field   string `json:"field"`
		Soft    string `json:"soft"`
		NavOn   string `json:"navon"`
		On      string `json:"on"`
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
	add := func(skin, colour string) {
		look := theme.Current(skin, colour)
		p := look.Palette
		r, g, b := theme.RGB(p.Accent)
		css := look.CSSVars()
		key := skin
		if look.Colours {
			key = skin + ":" + p.ID
		}
		out[key] = entry{
			Skin: skin, Colour: p.ID,
			Bg: p.Bg, Panel: p.Panel, Line: p.Line, Accent: p.Accent,
			Dim: p.Dim, Faint: p.Faint, Warn: p.Warn, Bad: p.Bad,
			Field: p.Field, Soft: p.Soft, NavOn: p.NavOn, On: p.On,
			RGB:     fmt.Sprintf("%d,%d,%d", r, g, b),
			Glow:    cssVar(css, "--glow"),
			Font:    cssVar(css, "--font"),
			Radius:  cssVar(css, "--r"),
			Border:  cssVar(css, "--bw"),
			Scan:    cssVar(css, "--scan"),
			Shadow:  cssVar(css, "--shadow"),
			WBorder: cssVar(css, "--wborder"),
			BarR:    barRadius(look),
		}
	}
	for _, skin := range theme.SkinIDs() {
		colours := theme.ColourIDs(skin)
		if len(colours) == 0 {
			add(skin, "")
			continue
		}
		for _, c := range colours {
			add(skin, c)
		}
	}
	data, err := json.Marshal(out)
	if err != nil {
		return "{}"
	}
	return string(data)
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

func barRadius(l theme.Look) string {
	if l.Radius >= 10 {
		return "99px"
	}
	return "0"
}
