package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

func main() {
	out := flag.String("o", "app.ico", "output .ico path")
	flag.Parse()

	sizes := []int{16, 24, 32, 48, 256}
	var images [][]byte
	for _, sz := range sizes {
		images = append(images, drawPNG(sz))
	}

	var ico bytes.Buffer
	le := binary.LittleEndian
	w16 := func(v uint16) { _ = binary.Write(&ico, le, v) }
	w32 := func(v uint32) { _ = binary.Write(&ico, le, v) }
	w16(0)
	w16(1)
	w16(uint16(len(sizes)))
	offset := 6 + 16*len(sizes)
	for i, sz := range sizes {
		b := byte(sz)
		if sz >= 256 {
			b = 0
		}
		ico.WriteByte(b)
		ico.WriteByte(b)
		ico.WriteByte(0)
		ico.WriteByte(0)
		w16(1)
		w16(32)
		w32(uint32(len(images[i])))
		w32(uint32(offset))
		offset += len(images[i])
	}
	for _, img := range images {
		ico.Write(img)
	}
	if err := os.WriteFile(*out, ico.Bytes(), 0o644); err != nil {
		panic(err)
	}
}

func drawPNG(sz int) []byte {
	img := image.NewNRGBA(image.Rect(0, 0, sz, sz))
	s := float64(sz) / 32.0
	bgTop := color.NRGBA{R: 0x14, G: 0x1A, B: 0x16, A: 255}
	bgBot := color.NRGBA{R: 0x0A, G: 0x0F, B: 0x0C, A: 255}
	glow := color.NRGBA{R: 0x3C, G: 0xFF, B: 0x6E, A: 255}
	white := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	radius := 7.5 * s

	for y := 0; y < sz; y++ {
		t := float64(y) / float64(sz-1)
		row := lerpC(bgTop, bgBot, t)
		for x := 0; x < sz; x++ {
			a := rrAlpha(float64(x), float64(y), 0.5*s, 0.5*s, float64(sz)-1.5*s, float64(sz)-1.5*s, radius, s)
			if a <= 0 {
				continue
			}
			fx, fy := float64(x)/s, float64(y)/s
			c := row
			d := micDist(fx, fy)
			if d <= 0 {
				c = lerpC(glow, white, 0.25)
			} else if d < 3.0 {
				c = lerpC(row, glow, (1-d/3.0)*0.6)
			}
			c.A = uint8(a * 255)
			img.SetNRGBA(x, y, c)
		}
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

func micDist(fx, fy float64) float64 {
	d := capDist(fx, fy, 15.5, 9.5, 15.5, 15.5) - 3.6
	if v := math.Abs(math.Hypot(fx-15.5, fy-14.5)-6.3) - 1.0; fy >= 14.5 && v < d {
		d = v
	}
	if v := capDist(fx, fy, 15.5, 19.8, 15.5, 22.0) - 1.1; v < d {
		d = v
	}
	if v := capDist(fx, fy, 11.2, 23.6, 19.8, 23.6) - 1.0; v < d {
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

func rrAlpha(px, py, x0, y0, x1, y1, r, s float64) float64 {
	cx := math.Max(x0+r, math.Min(px, x1-r))
	cy := math.Max(y0+r, math.Min(py, y1-r))
	d := math.Hypot(px-cx, py-cy)
	edge := 0.8 * s
	if d <= r-edge {
		return 1
	}
	if d >= r+0.2*s {
		return 0
	}
	return (r + 0.2*s - d) / (edge + 0.2*s)
}

func capDist(px, py, x1, y1, x2, y2 float64) float64 {
	dx, dy := x2-x1, y2-y1
	l2 := dx*dx + dy*dy
	if l2 == 0 {
		return math.Hypot(px-x1, py-y1)
	}
	t := ((px-x1)*dx + (py-y1)*dy) / l2
	t = math.Max(0, math.Min(1, t))
	return math.Hypot(px-(x1+t*dx), py-(y1+t*dy))
}
