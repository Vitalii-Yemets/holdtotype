package theme

import "strings"

// Skin is everything that makes the program look the way it does: colour,
// typography, shape, effects and the character of its motion.
type Skin struct {
	ID string

	// colour
	Bg     string
	Panel  string
	Line   string
	Accent string
	Dim    string
	Faint  string
	Warn   string
	Bad    string

	// typography
	FontCSS string // font stack for the settings page and the installer
	FontGDI string // family the plate and the shortcut window are drawn with
	FontPx  int32  // plate text size in device-independent pixels
	Weight  int32  // GDI weight: 400 regular, 600 semibold
	BrandLS string // letter-spacing of the title

	// shape
	Radius int32 // corner radius on the page, px
	Border int32 // border width, px
	Round  bool  // ask Windows for rounded window corners

	// effects
	Glow   bool    // green-terminal style halo around text and dots
	Scan   float64 // scanline overlay, 0…1
	Shadow string  // CSS shadow of the window and the plate

	// motion
	Level string  // level meter: "bars", "flat", "dots"
	Pulse float64 // recording dot: 1 is the usual speed, more is slower
}

const Default = "green"

var skins = []Skin{
	{ID: "green",
		Bg: "#0b0f0c", Panel: "#0e1410", Line: "#1d4a2b",
		Accent: "#3cff6e", Dim: "#20a34a", Faint: "#14803a", Warn: "#ffb347", Bad: "#ff7b6b",
		FontCSS: `Consolas,"Cascadia Mono",monospace`, FontGDI: "Consolas", FontPx: 15, Weight: 400, BrandLS: ".18em",
		Radius: 0, Border: 1, Round: false,
		Glow: true, Scan: 1, Shadow: "none",
		Level: "bars", Pulse: 1},

	{ID: "amber",
		Bg: "#100c0a", Panel: "#17110d", Line: "#4a3018",
		Accent: "#ff9e2c", Dim: "#b56a12", Faint: "#8a4f0d", Warn: "#ffd24a", Bad: "#ff6b5b",
		FontCSS: `Consolas,"Cascadia Mono",monospace`, FontGDI: "Consolas", FontPx: 15, Weight: 400, BrandLS: ".18em",
		Radius: 0, Border: 1, Round: false,
		Glow: true, Scan: 1, Shadow: "none",
		Level: "bars", Pulse: 1},

	{ID: "blue",
		Bg: "#0b0e10", Panel: "#0e1317", Line: "#1d3a4a",
		Accent: "#4cc3ff", Dim: "#1c7fb8", Faint: "#14608f", Warn: "#ffb347", Bad: "#ff7b6b",
		FontCSS: `Consolas,"Cascadia Mono",monospace`, FontGDI: "Consolas", FontPx: 15, Weight: 400, BrandLS: ".18em",
		Radius: 0, Border: 1, Round: false,
		Glow: true, Scan: 1, Shadow: "none",
		Level: "bars", Pulse: 1},

	{ID: "pink",
		Bg: "#100b0e", Panel: "#170e14", Line: "#4a1d3a",
		Accent: "#ff6ec7", Dim: "#b82f86", Faint: "#8f2467", Warn: "#ffb347", Bad: "#ff6b6b",
		FontCSS: `Consolas,"Cascadia Mono",monospace`, FontGDI: "Consolas", FontPx: 15, Weight: 400, BrandLS: ".18em",
		Radius: 0, Border: 1, Round: false,
		Glow: true, Scan: 1, Shadow: "none",
		Level: "bars", Pulse: 1},

	{ID: "editor",
		Bg: "#1e1e1e", Panel: "#252526", Line: "#3c3c3c",
		Accent: "#4fc1ff", Dim: "#9d9d9d", Faint: "#6e6e6e", Warn: "#cca700", Bad: "#f14c4c",
		FontCSS: `"Cascadia Mono",Consolas,"Segoe UI",sans-serif`, FontGDI: "Cascadia Mono", FontPx: 15, Weight: 400, BrandLS: ".02em",
		Radius: 3, Border: 1, Round: false,
		Glow: false, Scan: 0, Shadow: "0 10px 30px rgba(0,0,0,.45)",
		Level: "flat", Pulse: 1.5},

	{ID: "neon",
		Bg: "#150a22", Panel: "#1d0e30", Line: "#4a2472",
		Accent: "#ff5fc8", Dim: "#b06ee0", Faint: "#7d4fae", Warn: "#ffd24a", Bad: "#ff4d7d",
		FontCSS: `"Segoe UI Variable Display","Segoe UI",system-ui,sans-serif`, FontGDI: "Segoe UI Variable Display", FontPx: 16, Weight: 600, BrandLS: ".08em",
		Radius: 14, Border: 1, Round: true,
		Glow: true, Scan: 0.35, Shadow: "0 18px 46px rgba(150,40,220,.35)",
		Level: "bars", Pulse: 0.8},
}

func IDs() []string {
	out := make([]string, 0, len(skins))
	for _, s := range skins {
		out = append(out, s.ID)
	}
	return out
}

func Valid(id string) bool {
	for _, s := range skins {
		if s.ID == id {
			return true
		}
	}
	return false
}

func Get(id string) Skin {
	for _, s := range skins {
		if s.ID == id {
			return s
		}
	}
	return Get(Default)
}

// RGB returns the three channels of a #rrggbb string.
func RGB(hex string) (r, g, b uint8) {
	hex = strings.TrimPrefix(strings.TrimSpace(hex), "#")
	if len(hex) != 6 {
		return 0, 0, 0
	}
	v := make([]uint8, 3)
	for i := 0; i < 3; i++ {
		v[i] = uint8(hexPair(hex[i*2], hex[i*2+1]))
	}
	return v[0], v[1], v[2]
}

func hexPair(a, b byte) int {
	return hexDigit(a)*16 + hexDigit(b)
}

func hexDigit(c byte) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c-'a') + 10
	case c >= 'A' && c <= 'F':
		return int(c-'A') + 10
	}
	return 0
}

// CSSVars renders the skin as the contents of a :root rule.
func (s Skin) CSSVars() string {
	r, g, b := RGB(s.Accent)
	rgb := itoa(int(r)) + "," + itoa(int(g)) + "," + itoa(int(b))
	glow := "none"
	if s.Glow {
		glow = "0 0 7px rgba(" + rgb + ",.55)"
	}
	wborder := "none"
	if !s.Round {
		wborder = itoa(int(s.Border)) + "px solid " + s.Line
	}
	return "--wborder:" + wborder +
		";--bg:" + s.Bg + ";--panel:" + s.Panel + ";--line:" + s.Line +
		";--green:" + s.Accent + ";--dim:" + s.Dim + ";--faint:" + s.Faint +
		";--amber:" + s.Warn + ";--bad:" + s.Bad + ";--rgb:" + rgb +
		";--glow:" + glow +
		";--font:" + s.FontCSS +
		";--r:" + itoa(int(s.Radius)) + "px" +
		";--bw:" + itoa(int(s.Border)) + "px" +
		";--scan:" + dec(s.Scan) +
		";--shadow:" + s.Shadow +
		";--brandls:" + s.BrandLS
}

func dec(v float64) string {
	if v <= 0 {
		return "0"
	}
	if v >= 1 {
		return "1"
	}
	hundredths := int(v*100 + 0.5)
	return "." + pad2(hundredths)
}

func pad2(n int) string {
	if n < 10 {
		return "0" + itoa(n)
	}
	return itoa(n)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	digits := ""
	for n > 0 {
		digits = string(rune('0'+n%10)) + digits
		n /= 10
	}
	if neg {
		return "-" + digits
	}
	return digits
}
