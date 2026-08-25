package theme

import "strings"

type Palette struct {
	ID     string
	Bg     string
	Panel  string
	Line   string
	Text   string
	Accent string
	Dim    string
	Faint  string
	Warn   string
	Bad    string
	Rec    string

	Field   string
	Soft    string
	NavOn   string
	On      string
	TitleBg string
	SideBg  string
	KeyBg   string
	BtnBg   string
	BtnFg   string
	BtnLine string
	SelBg   string
	SelFg   string
	Brand   string
	Scrim   string
	Halo    string
}

type Skin struct {
	ID      string
	Palette string
	Colours bool

	FontCSS string
	FontGDI string
	PagePx  int32
	FontPx  int32
	Weight  int32
	BrandLS string
	CtlPad  string
	FieldPad string
	Caps    bool
	Flicker bool

	Radius int32
	Border int32
	Round  bool

	Glow   bool
	Scan   float64
	Shadow string

	Level string
	Pulse float64
}

const (
	DefaultSkin    = "terminal"
	DefaultPalette = "green"
)

var palettes = []Palette{
	{ID: "green", Bg: "#0b0f0c", Panel: "#0e1410", Line: "#1d4a2b",
		Text: "#3cff6e", Accent: "#3cff6e", Dim: "#20a34a", Faint: "#14803a",
		Warn: "#ffb347", Bad: "#ff7b6b", Rec: "#ff5b4d",
		Field: "#08100b", Soft: "#12241a", NavOn: "#101d14", On: "#123f22",
		TitleBg: "transparent", SideBg: "transparent", KeyBg: "transparent",
		BtnBg: "#101d14", BtnFg: "#3cff6e", BtnLine: "#20a34a", SelBg: "#3cff6e", SelFg: "#0b0f0c",
		Brand: "", Scrim: "rgba(3,7,4,.78)"},

	{ID: "amber", Bg: "#100c0a", Panel: "#17110d", Line: "#4a3018",
		Text: "#ff9e2c", Accent: "#ff9e2c", Dim: "#b56a12", Faint: "#8a4f0d",
		Warn: "#ffd24a", Bad: "#ff6b5b", Rec: "#ff5b4d",
		Field: "#120c07", Soft: "#2a1a0d", NavOn: "#22160c", On: "#402611",
		TitleBg: "transparent", SideBg: "transparent", KeyBg: "transparent",
		BtnBg: "#22160c", BtnFg: "#ff9e2c", BtnLine: "#b56a12", SelBg: "#ff9e2c", SelFg: "#100c0a",
		Brand: "", Scrim: "rgba(8,5,3,.78)"},

	{ID: "blue", Bg: "#0b0e10", Panel: "#0e1317", Line: "#1d3a4a",
		Text: "#4cc3ff", Accent: "#4cc3ff", Dim: "#1c7fb8", Faint: "#14608f",
		Warn: "#ffb347", Bad: "#ff7b6b", Rec: "#ff5b4d",
		Field: "#070f14", Soft: "#12222c", NavOn: "#101c24", On: "#123a52",
		TitleBg: "transparent", SideBg: "transparent", KeyBg: "transparent",
		BtnBg: "#101c24", BtnFg: "#4cc3ff", BtnLine: "#1c7fb8", SelBg: "#4cc3ff", SelFg: "#0b0e10",
		Brand: "", Scrim: "rgba(3,6,8,.78)"},

	{ID: "pink", Bg: "#100b0e", Panel: "#170e14", Line: "#4a1d3a",
		Text: "#ff6ec7", Accent: "#ff6ec7", Dim: "#b82f86", Faint: "#8f2467",
		Warn: "#ffb347", Bad: "#ff6b6b", Rec: "#ff5b4d",
		Field: "#120810", Soft: "#2a1222", NavOn: "#22101c", On: "#40183a",
		TitleBg: "transparent", SideBg: "transparent", KeyBg: "transparent",
		BtnBg: "#22101c", BtnFg: "#ff6ec7", BtnLine: "#b82f86", SelBg: "#ff6ec7", SelFg: "#100b0e",
		Brand: "", Scrim: "rgba(8,3,6,.78)"},

	{ID: "editor", Bg: "#1e1e1e", Panel: "#252526", Line: "#3c3c3c",
		Text: "#d4d4d4", Accent: "#4fc1ff", Dim: "#9d9d9d", Faint: "#6e6e6e",
		Warn: "#cca700", Bad: "#f14c4c", Rec: "#f14c4c",
		Field: "#3c3c3c", Soft: "#2d2d2d", NavOn: "#37373d", On: "#094771",
		TitleBg: "#323233", SideBg: "#252526", KeyBg: "#3c3c3c",
		BtnBg: "#0e639c", BtnFg: "#ffffff", BtnLine: "#0e639c", SelBg: "#0e639c", SelFg: "#ffffff",
		Brand: "", Scrim: "rgba(0,0,0,.6)"},

	{ID: "neon", Bg: "#150a22", Panel: "#1d0e30", Line: "#4a2472",
		Text: "#f3b6e4", Accent: "#46e0ff", Dim: "#b06ee0", Faint: "#7d4fae",
		Warn: "#ffd24a", Bad: "#ff4d7d", Rec: "#ff4d7d",
		Field: "#1e0f33", Soft: "#2a1442", NavOn: "#2b1240", On: "#4a2472",
		TitleBg: "linear-gradient(90deg,#26103f,#1a0b2b)", SideBg: "#190c29", KeyBg: "linear-gradient(90deg,rgba(255,95,200,.18),rgba(70,224,255,.14))",
		BtnBg: "transparent", BtnFg: "#f3b6e4", BtnLine: "#4a2472", SelBg: "linear-gradient(90deg,#ff5fc8,#46e0ff)", SelFg: "#150a22",
		Brand: "linear-gradient(90deg,#ff5fc8,#46e0ff)", Scrim: "rgba(10,4,18,.72)", Halo: "#a03ce0"},
}

var colourChoice = []string{"green", "amber", "blue", "pink"}

var skins = []Skin{
	{ID: "terminal", Palette: "green", Colours: true,
		FontCSS: `Consolas,"Cascadia Mono",monospace`, FontGDI: "Consolas",
		PagePx: 14, FontPx: 15, Weight: 400, BrandLS: ".18em", Caps: true, Flicker: true,
		CtlPad: "4px 8px", FieldPad: "7px 10px",
		Radius: 0, Border: 1, Round: false,
		Glow: true, Scan: 1, Shadow: "none",
		Level: "bars", Pulse: 1},

	{ID: "editor", Palette: "editor",
		FontCSS: `"Cascadia Mono",Consolas,"Segoe UI",sans-serif`, FontGDI: "Cascadia Mono",
		PagePx: 13, FontPx: 14, Weight: 400, BrandLS: ".02em",
		CtlPad: "5px 9px", FieldPad: "7px 11px",
		Radius: 3, Border: 1, Round: false,
		Glow: false, Scan: 0, Shadow: "0 10px 30px rgba(0,0,0,.45)",
		Level: "flat", Pulse: 1.5},

	{ID: "neon", Palette: "neon",
		FontCSS: `"Segoe UI Variable Display","Segoe UI",system-ui,sans-serif`, FontGDI: "Segoe UI Variable Display",
		PagePx: 13, FontPx: 15, Weight: 400, BrandLS: ".08em",
		CtlPad: "7px 13px", FieldPad: "10px 14px",
		Radius: 14, Border: 1, Round: true,
		Glow: true, Scan: 0.18, Shadow: "0 18px 46px rgba(150,40,220,.35)",
		Level: "bars", Pulse: 0.8},
}

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

func Current(skinID, colourID string) Look {
	s := GetSkin(skinID)
	id := s.Palette
	if s.Colours && ValidColour(colourID) {
		id = colourID
	}
	return Look{Skin: s, Palette: GetPalette(id)}
}

func Migrate(old string) (skin, colour string) {
	switch {
	case ValidColour(old):
		return DefaultSkin, old
	case ValidSkin(old):
		return old, DefaultPalette
	}
	return DefaultSkin, DefaultPalette
}

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

func (l Look) CSSVars() string {
	p := l.Palette
	r, g, b := RGB(p.Text)
	rgb := itoa(int(r)) + "," + itoa(int(g)) + "," + itoa(int(b))
	hr, hg, hb := RGB(p.Accent)
	hirgb := itoa(int(hr)) + "," + itoa(int(hg)) + "," + itoa(int(hb))
	glow := "none"
	if l.Glow {
		glow = "0 0 7px rgba(" + rgb + ",.55)"
	}
	hiGlow := "none"
	iconGlow := "none"
	if l.Glow {
		hiGlow = "0 0 8px rgba(" + hirgb + ",.6)"
		iconGlow = "drop-shadow(0 0 6px rgba(" + hirgb + ",.7))"
	}
	wborder := "none"
	if !l.Round {
		wborder = itoa(int(l.Border)) + "px solid " + p.Line
	}
	caps, ls := "none", "0"
	if l.Caps {
		caps, ls = "uppercase", "1px"
	}
	flicker := "none"
	if l.Flicker {
		flicker = "flicker 6s infinite"
	}
	barr := "0"
	if l.Radius >= 10 {
		barr = "99px"
	}
	brandBg, brandClip, brandFill := "none", "border-box", "currentColor"
	if p.Brand != "" {
		brandBg, brandClip, brandFill = p.Brand, "text", "transparent"
	}
	return "--wborder:" + wborder +
		";--bg:" + p.Bg + ";--panel:" + p.Panel + ";--line:" + p.Line +
		";--green:" + p.Text + ";--hi:" + p.Accent +
		";--dim:" + p.Dim + ";--faint:" + p.Faint +
		";--amber:" + p.Warn + ";--bad:" + p.Bad + ";--rec:" + p.Rec +
		";--rgb:" + rgb +
		";--field:" + p.Field + ";--soft:" + p.Soft +
		";--navon:" + p.NavOn + ";--on:" + p.On +
		";--titlebg:" + p.TitleBg + ";--sidebg:" + p.SideBg + ";--keybg:" + p.KeyBg +
		";--btnbg:" + p.BtnBg + ";--btnfg:" + p.BtnFg + ";--btnline:" + p.BtnLine +
		";--selbg:" + p.SelBg + ";--selfg:" + p.SelFg +
		";--brandbg:" + brandBg + ";--brandclip:" + brandClip + ";--brandfill:" + brandFill +
		";--scrim:" + p.Scrim +
		";--glow:" + glow + ";--higlow:" + hiGlow + ";--iconglow:" + iconGlow +
		";--font:" + l.FontCSS +
		";--fs:" + itoa(int(l.PagePx)) + "px" +
		";--caps:" + caps + ";--ls:" + ls + ";--flicker:" + flicker +
		";--ctlpad:" + l.CtlPad + ";--fieldpad:" + l.FieldPad +
		";--r:" + itoa(int(l.Radius)) + "px" +
		";--barr:" + barr +
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
