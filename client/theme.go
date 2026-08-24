package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
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

func themePalette() theme.Palette { return theme.Get(themeID()) }

func themeCSSVars() string { return themePalette().CSSVars() }

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

func themeListJSON() string {
	type entry struct {
		Bg     string `json:"bg"`
		Panel  string `json:"panel"`
		Line   string `json:"line"`
		Accent string `json:"accent"`
		Dim    string `json:"dim"`
		Faint  string `json:"faint"`
		Warn   string `json:"warn"`
		Bad    string `json:"bad"`
		RGB    string `json:"rgb"`
		Glow   string `json:"glow"`
	}
	out := map[string]entry{}
	for _, id := range theme.IDs() {
		p := theme.Get(id)
		r, g, b := theme.RGB(p.Accent)
		out[id] = entry{
			Bg: p.Bg, Panel: p.Panel, Line: p.Line, Accent: p.Accent,
			Dim: p.Dim, Faint: p.Faint, Warn: p.Warn, Bad: p.Bad,
			RGB:  fmt.Sprintf("%d,%d,%d", r, g, b),
			Glow: fmt.Sprintf("0 0 7px rgba(%d,%d,%d,.55)", r, g, b),
		}
	}
	data, err := json.Marshal(out)
	if err != nil {
		return "{}"
	}
	return string(data)
}
