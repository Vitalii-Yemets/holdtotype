package ovplace

import "fmt"

type Rect struct {
	Left, Top, Right, Bottom int32
}

type Frac struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

const (
	PosBottom      = "bottom"
	PosTop         = "top"
	PosCaret       = "caret"
	PosTopLeft     = "top-left"
	PosTopRight    = "top-right"
	PosLeft        = "left"
	PosRight       = "right"
	PosBottomLeft  = "bottom-left"
	PosBottomRight = "bottom-right"
	PosCenter      = "center"
	PosCustom      = "custom"
)

var anchors = map[string]bool{
	PosBottom: true, PosTop: true, PosCaret: true,
	PosTopLeft: true, PosTopRight: true, PosLeft: true, PosRight: true,
	PosBottomLeft: true, PosBottomRight: true, PosCenter: true, PosCustom: true,
}

func Valid(mode string) bool { return anchors[mode] }

// ResKey names the monitor by its full resolution, so a spot dragged out on
// one screen never lands blindly on another.
func ResKey(screen Rect) string {
	return fmt.Sprintf("%dx%d", screen.Right-screen.Left, screen.Bottom-screen.Top)
}

func clamp(v, lo, hi int32) int32 {
	if hi < lo {
		return lo
	}
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// Anchor puts a w×h plate at one of the eight scheme positions inside the
// work area, a margin away from the edges.
func Anchor(mode string, wa Rect, w, h, margin int32) (int32, int32) {
	left := wa.Left + margin
	right := wa.Right - w - margin
	centerX := wa.Left + (wa.Right-wa.Left-w)/2
	top := wa.Top + margin
	bottom := wa.Bottom - h - margin
	centerY := wa.Top + (wa.Bottom-wa.Top-h)/2
	var x, y int32
	switch mode {
	case PosTopLeft:
		x, y = left, top
	case PosTop:
		x, y = centerX, top
	case PosTopRight:
		x, y = right, top
	case PosLeft:
		x, y = left, centerY
	case PosRight:
		x, y = right, centerY
	case PosBottomLeft:
		x, y = left, bottom
	case PosBottomRight:
		x, y = right, bottom
	case PosCenter:
		x, y = centerX, centerY
	default:
		x, y = centerX, bottom
	}
	return clamp(x, wa.Left, wa.Right-w), clamp(y, wa.Top, wa.Bottom-h)
}

// Custom places the plate's centre at the remembered fraction of the work
// area, clamped so the plate never leaves the screen.
func Custom(f Frac, wa Rect, w, h int32) (int32, int32) {
	cx := float64(wa.Left) + f.X*float64(wa.Right-wa.Left)
	cy := float64(wa.Top) + f.Y*float64(wa.Bottom-wa.Top)
	x := int32(cx) - w/2
	y := int32(cy) - h/2
	return clamp(x, wa.Left, wa.Right-w), clamp(y, wa.Top, wa.Bottom-h)
}

// FracOf answers where the plate's centre sits, as a fraction of the work
// area — the inverse of Custom, for saving a dragged spot.
func FracOf(x, y, w, h int32, wa Rect) Frac {
	ww := wa.Right - wa.Left
	wh := wa.Bottom - wa.Top
	if ww <= 0 || wh <= 0 {
		return Frac{X: 0.5, Y: 0.9}
	}
	fx := (float64(x) + float64(w)/2 - float64(wa.Left)) / float64(ww)
	fy := (float64(y) + float64(h)/2 - float64(wa.Top)) / float64(wh)
	if fx < 0 {
		fx = 0
	}
	if fx > 1 {
		fx = 1
	}
	if fy < 0 {
		fy = 0
	}
	if fy > 1 {
		fy = 1
	}
	return Frac{X: fx, Y: fy}
}

// CleanCustom keeps only sane entries: fractions inside the unit square
// under resolution-shaped keys.
func CleanCustom(in map[string]Frac) map[string]Frac {
	if len(in) == 0 {
		return nil
	}
	out := map[string]Frac{}
	for k, f := range in {
		var w, h int
		if n, err := fmt.Sscanf(k, "%dx%d", &w, &h); n != 2 || err != nil || w <= 0 || h <= 0 {
			continue
		}
		if f.X < 0 || f.X > 1 || f.Y < 0 || f.Y > 1 {
			continue
		}
		out[k] = f
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
