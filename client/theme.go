package main

import (
	"encoding/json"
	"image/color"
	"log"
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
	colGreen = colorref(p.Text)
	colGreenDm = colorref(p.Dim)
	colGreenLo = colorref(p.Faint)
	colHi = colorref(p.Accent)
	colHiLo = mixHex(p.Accent, 0.45)
	colAmber = colorref(p.Warn)
	colAmberDm = mixHex(p.Warn, 0.35)
	colBad = colorref(p.Bad)
	colBadDm = mixHex(p.Bad, 0.35)
	colRed = colorref(p.Rec)
	colRedDm = mixHex(p.Rec, 0.35)
	colAskBg = mixHex(p.Panel, 1.0)
	colLine = colorref(p.Line)

	rebuildIcons(p)
	log.Printf("оформление: %s, цвет %s", skin, p.ID)
}

func nrgba(hex string) color.NRGBA {
	r, g, b := theme.RGB(hex)
	return color.NRGBA{R: r, G: g, B: b, A: 255}
}

func lift(c color.NRGBA, by int) color.NRGBA {
	clamp := func(v uint8) uint8 {
		n := int(v) + by
		if n < 0 {
			n = 0
		}
		if n > 255 {
			n = 255
		}
		return uint8(n)
	}
	return color.NRGBA{R: clamp(c.R), G: clamp(c.G), B: clamp(c.B), A: 255}
}

func paletteTile(p theme.Palette) iconTile {
	bot := nrgba(p.Bg)
	core := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	if p.Light() {
		core = color.NRGBA{R: 0x1a, G: 0x1a, B: 0x1a, A: 255}
	}
	return iconTile{top: lift(bot, 10), bot: bot, core: core}
}

func rebuildIcons(p theme.Palette) {
	tile := paletteTile(p)
	accent := nrgba(p.Text)
	bad := nrgba(p.Bad)
	warn := nrgba(p.Warn)
	off := nrgba(p.Off)
	iconIdle = iconPNG(tile, accent)
	iconRecording = iconPNG(tile, bad)
	iconProcessing = iconPNG(tile, warn)
	iconDisabled = iconPNG(tile, off)
	iconError = iconPNG(tile, off, bad)
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

func skinListJSON() string {
	type entry struct {
		Skin   string `json:"skin"`
		Colour string `json:"colour"`
		Accent string `json:"accent"`
		Vars   string `json:"vars"`
	}
	out := map[string]entry{}
	add := func(skin, colour string) {
		look := theme.Current(skin, colour)
		p := look.Palette
		key := skin
		if look.Colours {
			key = skin + ":" + p.ID
		}
		out[key] = entry{Skin: skin, Colour: p.ID, Accent: p.Text, Vars: look.CSSVars()}
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
