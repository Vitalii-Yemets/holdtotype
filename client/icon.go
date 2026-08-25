package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"

	"holdtotype/internal/theme"
)

func micDist(fx, fy float64) float64 {
	d := capsuleDist(fx, fy, 15.5, 9.5, 15.5, 15.5) - 3.6
	if v := math.Abs(math.Hypot(fx-15.5, fy-14.5)-6.3) - 1.0; fy >= 14.5 && v < d {
		d = v
	}
	if v := capsuleDist(fx, fy, 15.5, 19.8, 15.5, 22.0) - 1.1; v < d {
		d = v
	}
	if v := capsuleDist(fx, fy, 11.2, 23.6, 19.8, 23.6) - 1.0; v < d {
		d = v
	}
	return d
}

func lerpC(a, b color.NRGBA, t float64) color.NRGBA {
	return color.NRGBA{
		R: uint8(float64(a.R)*(1-t) + float64(b.R)*t),
		G: uint8(float64(a.G)*(1-t) + float64(b.G)*t),
		B: uint8(float64(a.B)*(1-t) + float64(b.B)*t),
		A: 255,
	}
}

type iconTile struct {
	top  color.NRGBA
	bot  color.NRGBA
	core color.NRGBA
}

func iconPNG(tile iconTile, glow color.NRGBA, badge ...color.NRGBA) []byte {
	const sz = 32
	img := image.NewNRGBA(image.Rect(0, 0, sz, sz))
	bgTop := tile.top
	bgBot := tile.bot
	const radius = 7.5

	for y := 0; y < sz; y++ {
		t := float64(y) / (sz - 1)
		row := lerpC(bgTop, bgBot, t)
		for x := 0; x < sz; x++ {
			a := roundRectAlpha(float64(x), float64(y), 0.5, 0.5, sz-1.5, sz-1.5, radius)
			if a <= 0 {
				continue
			}
			fx, fy := float64(x), float64(y)
			c := row
			d := micDist(fx, fy)
			if d <= 0 {
				c = lerpC(glow, tile.core, 0.25)
			} else if d < 3.0 {
				c = lerpC(row, glow, (1-d/3.0)*0.6)
			}
			c.A = uint8(a * 255)
			img.SetNRGBA(x, y, c)
		}
	}

	if len(badge) > 0 {
		b := badge[0]
		for y := sz - 11; y < sz-2; y++ {
			for x := sz - 11; x < sz-2; x++ {
				img.SetNRGBA(x, y, b)
			}
		}
	}

	if len(badge) > 0 {
		b := badge[0]
		for y := sz - 12; y < sz-3; y++ {
			for x := sz - 12; x < sz-3; x++ {
				img.SetNRGBA(x, y, b)
			}
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil
	}
	return buf.Bytes()
}

func roundRectAlpha(px, py, x0, y0, x1, y1, r float64) float64 {
	cx := math.Max(x0+r, math.Min(px, x1-r))
	cy := math.Max(y0+r, math.Min(py, y1-r))
	d := math.Hypot(px-cx, py-cy)
	if d <= r-0.8 {
		return 1
	}
	if d >= r+0.2 {
		return 0
	}
	return (r + 0.2 - d) / 1.0
}

func capsuleDist(px, py, x1, y1, x2, y2 float64) float64 {
	dx, dy := x2-x1, y2-y1
	l2 := dx*dx + dy*dy
	if l2 == 0 {
		return math.Hypot(px-x1, py-y1)
	}
	t := ((px-x1)*dx + (py-y1)*dy) / l2
	t = math.Max(0, math.Min(1, t))
	return math.Hypot(px-(x1+t*dx), py-(y1+t*dy))
}

var (
	startupPalette = theme.GetPalette(theme.DefaultPalette)
	startupTile    = paletteTile(startupPalette)

	iconIdle       = iconPNG(startupTile, nrgba(startupPalette.Accent))
	iconRecording  = iconPNG(startupTile, nrgba(startupPalette.Bad))
	iconProcessing = iconPNG(startupTile, nrgba(startupPalette.Warn))
	iconDisabled   = iconPNG(startupTile, nrgba(startupPalette.Off))
	iconError      = iconPNG(startupTile, nrgba(startupPalette.Off), nrgba(startupPalette.Bad))
)

func loadTrayIcons() map[int]uintptr {
	return map[int]uintptr{
		trayIdle:       hIconFromPNG(iconIdle),
		trayRecording:  hIconFromPNG(iconRecording),
		trayProcessing: hIconFromPNG(iconProcessing),
		trayOff:        hIconFromPNG(iconDisabled),
		trayError:      hIconFromPNG(iconError),
	}
}
