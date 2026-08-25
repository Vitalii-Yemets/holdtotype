package theme

import "strings"

// Palette is colour and nothing else: the four the terminal look ships with,
// plus the one that belongs to each of the other skins.
type Palette struct {
	ID     string
	Bg     string
	Panel  string
	Line   string
	Accent string
	Dim    string
	Faint  string
	Warn   string
	Bad    string

	// the surfaces the page is built from: a field is what you type into,
	// soft is a quiet chip or a hairline, navon marks the section you are in
	// and on marks a choice that is made
	Field string
	Soft  string
	NavOn string
	On    string
}

// Skin is the design: typography, shape, effects and the character of motion.
// Colour comes from a palette — the terminal skin lets you pick one, the
// others carry their own.
type Skin struct {
	ID      string
	Palette string // the palette it is drawn with by default
	Colours bool   // whether the colour choice is offered for this skin

	FontCSS string // font stack for the settings page and the installer
	FontGDI string // family the plate and the shortcut window are drawn with
	FontPx  int32  // plate text size in device-independent pixels
	Weight  int32  // GDI weight: 400 regular, 600 semibold
	BrandLS string // letter-spacing of the title

	Radius int32 // corner radius on the page, px
	Border int32 // border width, px
	Round  bool  // ask Windows for rounded window corners

	Glow   bool    // halo around text and dots
	Scan   float64 // scanline overlay, 0…1
	Shadow string  // CSS shadow of the window and the plate

	Level string  // level meter: "bars", "flat", "dots"
	Pulse float64 // recording dot: 1 is the usual speed, more is slower
}

const (
	DefaultSkin    = "terminal"
	DefaultPalette = "green"
)

var palettes = []Palette{
	{ID: "green", Bg: "#0b0f0c", Panel: "#0e1410", Line: "#1d4a2b",
		Accent: "#3cff6e", Dim: "#20a34a", Faint: "#14803a", Warn: "#ffb347", Bad: "#ff7b6b",
		Field: "#08100b", Soft: "#12241a", NavOn: "#101d14", On: "#123f22"},
	{ID: "amber", Bg: "#100c0a", Panel: "#17110d", Line: "#4a3018",
		Accent: "#ff9e2c", Dim: "#b56a12", Faint: "#8a4f0d", Warn: "#ffd24a", Bad: "#ff6b5b",
		Field: "#120c07", Soft: "#2a1a0d", NavOn: "#22160c", On: "#402611"},
	{ID: "blue", Bg: "#0b0e10", Panel: "#0e1317", Line: "#1d3a4a",
		Accent: "#4cc3ff", Dim: "#1c7fb8", Faint: "#14608f", Warn: "#ffb347", Bad: "#ff7b6b",
		Field: "#070f14", Soft: "#12222c", NavOn: "#101c24", On: "#123a52"},
	{ID: "pink", Bg: "#100b0e", Panel: "#170e14", Line: "#4a1d3a",
		Accent: "#ff6ec7", Dim: "#b82f86", Faint: "#8f2467", Warn: "#ffb347", Bad: "#ff6b6b",
		Field: "#120810", Soft: "#2a1222", NavOn: "#22101c", On: "#40183a"},

	// palettes that belong to a skin and are not offered as a choice
	{ID: "editor", Bg: "#1e1e1e", Panel: "#252526", Line: "#3c3c3c",
		Accent: "#4fc1ff", Dim: "#9d9d9d", Faint: "#6e6e6e", Warn: "#cca700", Bad: "#f14c4c",
		Field: "#3c3c3c", Soft: "#2d2d30", NavOn: "#37373d", On: "#094771"},
	{ID: "neon", Bg: "#150a22", Panel: "#1d0e30", Line: "#4a2472",
		Accent: "#ff5fc8", Dim: "#b06ee0", Faint: "#7d4fae", Warn: "#ffd24a", Bad: "#ff4d7d",
		Field: "#1e0f33", Soft: "#2a1349", NavOn: "#26123f", On: "#4a2472"},
}

// colourChoice lists the palettes the terminal skin offers, in order.
var colourChoice = []string{"green", "amber", "blue", "pink"}

var skins = []Skin{
	{ID: "terminal", Palette: "green", Colours: true,
		FontCSS: `Consolas,"Cascadia Mono",monospace`, FontGDI: "Consolas", FontPx: 15, Weight: 400, BrandLS: ".18em",
		Radius: 0, Border: 1, Round: false,
		Glow: true, Scan: 1, Shadow: "none",
		Level: "bars", Pulse: 1},

	{ID: "editor", Palette: "editor",
		FontCSS: `"Cascadia Mono",Consolas,"Segoe UI",sans-serif`, FontGDI: "Cascadia Mono", FontPx: 15, Weight: 400, BrandLS: ".02em",
		Radius: 3, Border: 1, Round: false,
		Glow: false, Scan: 0, Shadow: "0 10px 30px rgba(0,0,0,.45)",
		Level: "flat", Pulse: 1.5},

	{ID: "neon", Palette: "neon",
		FontCSS: `"Segoe UI Variable Display","Segoe UI",system-ui,sans-serif`, FontGDI: "Segoe UI Variable Display", FontPx: 16, Weight: 600, BrandLS: ".08em",
		Radius: 14, Border: 1, Round: true,
		Glow: true, Scan: 0.35, Shadow: "0 18px 46px rgba(150,40,220,.35)",
		Level: "bars", Pulse: 0.8},
}

// Look is a skin with the colours it is being drawn in.
type Look struct {
	Skin
	Palette Palette
}

func SkinIDs() []string {
	out := make([]string, 0, len(skins))
	for _, s := range skins {
		out = append(out, s.ID)
	}
	return out
}

// ColourIDs lists the colours a skin offers; empty when it carries its own.
func ColourIDs(skinID string) []string {
	if !GetSkin(skinID).Colours {
		return nil
	}
	return append([]string(nil), colourChoice...)
}

func ValidSkin(id string) bool {
	for _, s := range skins {
		if s.ID == id {
			return true
		}
	}
	return false
}

func ValidColour(id string) bool {
	for _, c := range colourChoice {
		if c == id {
			return true
		}
	}
	return false
}

func GetSkin(id string) Skin {
	for _, s := range skins {
		if s.ID == id {
			return s
		}
	}
	return GetSkin(DefaultSkin)
}

func GetPalette(id string) Palette {
	for _, p := range palettes {
		if p.ID == id {
			return p
		}
	}
	return GetPalette(DefaultPalette)
}

// Current puts a skin and a colour together. A skin that carries its own
// colours ignores the choice.
func Current(skinID, colourID string) Look {
	s := GetSkin(skinID)
	id := s.Palette
	if s.Colours && ValidColour(colourID) {
		id = colourID
	}
	return Look{Skin: s, Palette: GetPalette(id)}
}

// Migrate turns the single value older versions stored into a skin and a colour.
func Migrate(old string) (skin, colour string) {
	switch {
	case ValidColour(old):
		return DefaultSkin, old
	case ValidSkin(old):
		return old, DefaultPalette
	}
	return DefaultSkin, DefaultPalette
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

func hexPair(a, b byte) int { return hexDigit(a)*16 + hexDigit(b) }

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

// CSSVars renders the look as the contents of a :root rule.
func (l Look) CSSVars() string {
	p := l.Palette
	r, g, b := RGB(p.Accent)
	rgb := itoa(int(r)) + "," + itoa(int(g)) + "," + itoa(int(b))
	glow := "none"
	if l.Glow {
		glow = "0 0 7px rgba(" + rgb + ",.55)"
	}
	wborder := "none"
	if !l.Round {
		wborder = itoa(int(l.Border)) + "px solid " + p.Line
	}
	return "--wborder:" + wborder +
		";--bg:" + p.Bg + ";--panel:" + p.Panel + ";--line:" + p.Line +
		";--green:" + p.Accent + ";--dim:" + p.Dim + ";--faint:" + p.Faint +
		";--amber:" + p.Warn + ";--bad:" + p.Bad + ";--rgb:" + rgb +
		";--field:" + p.Field + ";--soft:" + p.Soft +
		";--navon:" + p.NavOn + ";--on:" + p.On +
		";--glow:" + glow +
		";--font:" + l.FontCSS +
		";--r:" + itoa(int(l.Radius)) + "px" +
		";--bw:" + itoa(int(l.Border)) + "px" +
		";--scan:" + dec(l.Scan) +
		";--shadow:" + l.Shadow +
		";--brandls:" + l.BrandLS
}

func dec(v float64) string {
	if v <= 0 {
		return "0"
	}
	if v >= 1 {
		return "1"
	}
	return "." + pad2(int(v*100+0.5))
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
