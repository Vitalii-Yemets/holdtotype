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
	Card    string
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
	Label   string
	Brand   string
	Scrim   string

	Off     string
	Ok      string
	BadBg   string
	BadLine string

	BtnBgH  string
	DangerBg string
	DangerFg string
	DangerBgH string
	Dot      string
	Btn2Bg  string
	Btn2Fg  string
	Btn2Line string
	Btn2BgH string
	Focus   string
	SwBg     string
	SwOnBg   string
	SwKnob   string
	SwOnLine string
	KeyFg    string
	KeyLine  string
	TagFg    string
	TagLine  string
	TagBg    string
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
	FieldPad string
	CtlFS    string
	WeightB  int32
	Caps    bool
	Flicker bool

	Radius int32
	BarR   string
	Border int32
	Round  bool
	SmallR bool
	DotR    string
	BadgeR  string
	PanelR  string
	SwitchR string
	CtlH    string
	BtnR    string
	FieldR  string
	CardR   string
	KeyR    string

	Glow   bool
	Scan   float64
	Shadow string

	Level string
	Pulse float64

	Mark  string
	Brackets bool
	Flash string
}

const (
	DefaultSkin    = "editor"
	DefaultPalette = "green"
)

var palettes = []Palette{
	{ID: "green", Bg: "#0b0f0c", Panel: "#0e1410", Line: "#1d4a2b",
		Text: "#3cff6e", Accent: "#3cff6e", Dim: "#20a34a", Faint: "#14803a",
		Warn: "#ffb347", Bad: "#ff7b6b", Rec: "#ff5b4d",
		Field: "#08100b", Soft: "#12241a", NavOn: "#101d14", On: "#123f22",
		TitleBg: "transparent", SideBg: "transparent", KeyBg: "transparent",
		BtnBg: "#101d14", BtnFg: "#3cff6e", BtnLine: "#20a34a", SelBg: "#3cff6e", SelFg: "#0b0f0c",
		Label: "#f2fff5", Brand: "", Scrim: "rgba(3,7,4,.78)", KeyFg: "#f2fff5"},

	{ID: "amber", Bg: "#100c0a", Panel: "#17110d", Line: "#4a3018",
		Text: "#ff9e2c", Accent: "#ff9e2c", Dim: "#b56a12", Faint: "#8a4f0d",
		Warn: "#ffd24a", Bad: "#ff6b5b", Rec: "#ff5b4d",
		Field: "#120c07", Soft: "#2a1a0d", NavOn: "#22160c", On: "#402611",
		TitleBg: "transparent", SideBg: "transparent", KeyBg: "transparent",
		BtnBg: "#22160c", BtnFg: "#ff9e2c", BtnLine: "#b56a12", SelBg: "#ff9e2c", SelFg: "#100c0a",
		Label: "#ffe9c9", Brand: "", Scrim: "rgba(8,5,3,.78)", KeyFg: "#ffe9c9"},

	{ID: "blue", Bg: "#0b0e10", Panel: "#0e1317", Line: "#1d3a4a",
		Text: "#4cc3ff", Accent: "#4cc3ff", Dim: "#1c7fb8", Faint: "#14608f",
		Warn: "#ffb347", Bad: "#ff7b6b", Rec: "#ff5b4d",
		Field: "#070f14", Soft: "#12222c", NavOn: "#101c24", On: "#123a52",
		TitleBg: "transparent", SideBg: "transparent", KeyBg: "transparent",
		BtnBg: "#101c24", BtnFg: "#4cc3ff", BtnLine: "#1c7fb8", SelBg: "#4cc3ff", SelFg: "#0b0e10",
		Label: "#e4f6ff", Brand: "", Scrim: "rgba(3,6,8,.78)", KeyFg: "#e4f6ff"},

	{ID: "pink", Bg: "#100b0e", Panel: "#170e14", Line: "#4a1d3a",
		Text: "#ff6ec7", Accent: "#ff6ec7", Dim: "#b82f86", Faint: "#8f2467",
		Warn: "#ffb347", Bad: "#ff6b6b", Rec: "#ff5b4d",
		Field: "#120810", Soft: "#2a1222", NavOn: "#22101c", On: "#40183a",
		TitleBg: "transparent", SideBg: "transparent", KeyBg: "transparent",
		BtnBg: "#22101c", BtnFg: "#ff6ec7", BtnLine: "#b82f86", SelBg: "#ff6ec7", SelFg: "#100b0e",
		Label: "#ffe6f4", Brand: "", Scrim: "rgba(8,3,6,.78)", KeyFg: "#ffe6f4"},

	{ID: "editor", Bg: "#1b1f24", Panel: "#20252b", Line: "#2c333c",
		Text: "#d7dce2", Accent: "#4fc1ff", Dim: "#9aa4b1", Faint: "#6b7480",
		Warn: "#e3b341", Bad: "#e5534b", Rec: "#e5534b",
		Field: "#161a1f", Card: "#20252b", Soft: "#262c34", NavOn: "#2a3139", On: "#1e3a55",
		TitleBg: "#1f242a", SideBg: "#1a1e23", KeyBg: "#161a1f",
		BtnBg: "#2f81c8", BtnFg: "#ffffff", BtnLine: "#2f81c8", SelBg: "#2f81c8", SelFg: "#ffffff",
		BtnBgH: "#3a8fd8", Btn2Bg: "#262c34", Btn2Fg: "#d7dce2", Btn2Line: "#39414b", Btn2BgH: "#2f3640",
		DangerBg: "#b23a3a", DangerFg: "#ffffff", DangerBgH: "#c44545", Dot: "#7ed49b", Ok: "#7ed49b",
		Focus: "#4fc1ff", SwBg: "#161a1f", SwOnBg: "#2f81c8", SwKnob: "#ffffff", SwOnLine: "#4fc1ff",
		KeyFg: "#4fc1ff", KeyLine: "#39414b",
		Label: "#e6eaef", Brand: "", Scrim: "rgba(0,0,0,.6)"},

	{ID: "neon", Bg: "#120a1e", Panel: "#190f2b", Line: "#3a2360",
		Text: "#f1e4ff", Accent: "#46e0ff", Dim: "#b79bd6", Faint: "#7e63a3",
		Warn: "#ffd24a", Bad: "#ff4d7d", Rec: "#ff4d7d",
		Field: "#150b25", Card: "#190f2b", Soft: "#24163d", NavOn: "#221540", On: "#2b1a4a",
		TitleBg: "#160c26", SideBg: "#140b22", KeyBg: "#1a0f2e",
		BtnBg: "linear-gradient(90deg,#ff5fc8,#46e0ff)", BtnFg: "#120a1e", BtnLine: "transparent",
		SelBg: "linear-gradient(90deg,#ff5fc8,#46e0ff)", SelFg: "#120a1e",
		BtnBgH: "linear-gradient(90deg,#ff7ad2,#6be7ff)",
		Btn2Bg: "#1f1236", Btn2Fg: "#f1e4ff", Btn2Line: "#4a2f78", Btn2BgH: "#291a45",
		DangerBg: "#ff4d7d", DangerFg: "#1b0710", DangerBgH: "#ff6690",
		Dot: "#5cf2c4", Ok: "#5cf2c4", Focus: "#46e0ff",
		SwBg: "#150b25", SwOnBg: "#2b1a4a", SwKnob: "#46e0ff", SwOnLine: "#46e0ff",
		KeyFg: "#46e0ff", KeyLine: "#46e0ff", TagFg: "#ff8fd9", TagLine: "#ff5fc8",
		Label: "#ffffff", Brand: "linear-gradient(90deg,#ff5fc8,#46e0ff)", Scrim: "rgba(10,4,18,.72)"},

	{ID: "soft", Bg: "#f2eef8", Panel: "#fbf9fe", Line: "#e0d8ee",
		Text: "#2b2438", Accent: "#7c5cff", Dim: "#6f6684", Faint: "#9a92ab",
		Warn: "#b26a00", Bad: "#c93d64", Rec: "#e2557a",
		Field: "#ffffff", Card: "#fbf9fe", Soft: "#e6dff3", NavOn: "#e8e0f6", On: "#e8e0f6",
		TitleBg: "#f6f2fb", SideBg: "#f4f0fa", KeyBg: "#ffffff",
		BtnBg: "#7c5cff", BtnFg: "#ffffff", BtnLine: "transparent", SelBg: "#7c5cff", SelFg: "#ffffff",
		BtnBgH: "#6c4cf0", Btn2Bg: "#ebe5f7", Btn2Fg: "#4a3f66", Btn2Line: "transparent", Btn2BgH: "#e1d9f2",
		DangerBg: "#e2557a", DangerFg: "#ffffff", DangerBgH: "#d1436a", Dot: "#2e9e6a", Ok: "#2e9e6a", Focus: "#7c5cff",
		SwBg: "#e6dff3", SwOnBg: "#7c5cff", SwKnob: "#ffffff", SwOnLine: "#7c5cff",
		KeyFg: "#7c5cff", KeyLine: "#d8cdf2", TagFg: "#6f6684", TagLine: "transparent", TagBg: "#efe9fb",
		Label: "#241d33", Brand: "", Scrim: "rgba(80,60,120,.35)"},

	{ID: "paper", Bg: "#f4f1ea", Panel: "#fbfaf6", Line: "#d9d4c8",
		Text: "#1d1d1b", Accent: "#1f4fbf", Dim: "#5d6067", Faint: "#8b8e94",
		Warn: "#8a5a00", Bad: "#b3261e", Rec: "#b3261e",
		Field: "#ffffff", Card: "#fbfaf6", Soft: "#e9e5db", NavOn: "#e6e1d5", On: "#dfe6f7",
		TitleBg: "#efebe2", SideBg: "#f1ede4", KeyBg: "#ffffff",
		BtnBg: "#1f4fbf", BtnFg: "#ffffff", BtnLine: "#1f4fbf", SelBg: "#1f4fbf", SelFg: "#ffffff",
		BtnBgH: "#1a45a8", Btn2Bg: "#e9e5db", Btn2Fg: "#1d1d1b", Btn2Line: "#cfc9bb", Btn2BgH: "#dfdacd",
		DangerBg: "#b3261e", DangerFg: "#ffffff", DangerBgH: "#9a1f18", Dot: "#1f7a3f", Focus: "#1f4fbf",
		Ok: "#1f7a3f", SwBg: "#e9e5db", SwOnBg: "#1f4fbf", SwKnob: "#ffffff", SwOnLine: "#1f4fbf",
		KeyFg: "#1f4fbf", KeyLine: "#c7c1b3", TagFg: "#5d6067", TagLine: "#cfc9bb", TagBg: "#ffffff",
		Label: "#111110", Brand: "", Scrim: "rgba(40,35,20,.35)"},


	{ID: "fluent", Bg: "#202020", Panel: "#2b2b2b", Line: "#3a3a3a",
		Text: "#ffffff", Accent: "#60cdff", Dim: "#c5c5c5", Faint: "#8a8a8a",
		Warn: "#fce100", Bad: "#ff99a4", Rec: "#ff6b6b",
		Field: "#323232", Card: "#2b2b2b", Soft: "#323232", NavOn: "#2d2d2d", On: "#24404f",
		TitleBg: "#202020", SideBg: "#202020", KeyBg: "#323232",
		BtnBg: "#60cdff", BtnFg: "#000000", BtnLine: "#60cdff", SelBg: "#60cdff", SelFg: "#000000",
		BtnBgH: "#56b8e6", Btn2Bg: "#323232", Btn2Fg: "#ffffff", Btn2Line: "#3f3f3f", Btn2BgH: "#3a3a3a",
		DangerBg: "#c42b1c", DangerFg: "#ffffff", DangerBgH: "#d33a2b", Dot: "#6ccb5f", Ok: "#6ccb5f", Focus: "#60cdff",
		SwBg: "transparent", SwOnBg: "#60cdff", SwKnob: "#000000", SwOnLine: "#60cdff",
		KeyFg: "#60cdff", KeyLine: "#454545", TagFg: "#c5c5c5", TagLine: "transparent", TagBg: "#323232",
		Label: "#ffffff", Brand: "", Scrim: "rgba(0,0,0,.6)"},

	{ID: "studio", Bg: "#1c1d1f", Panel: "#242527", Line: "#343538",
		Text: "#e6e3dc", Accent: "#ff9f43", Dim: "#a19d94", Faint: "#6f6c66",
		Warn: "#ffd166", Bad: "#e06060", Rec: "#d64545",
		Field: "#161718", Card: "#242527", Soft: "#2c2d30", NavOn: "#2c2d30", On: "#3a2a16",
		TitleBg: "#191a1c", SideBg: "#1a1b1d", KeyBg: "#161718",
		BtnBg: "#ff9f43", BtnFg: "#1c1d1f", BtnLine: "#ff9f43", SelBg: "#ff9f43", SelFg: "#1c1d1f",
		BtnBgH: "#ffb066", Btn2Bg: "#2c2d30", Btn2Fg: "#e6e3dc", Btn2Line: "#45464a", Btn2BgH: "#35363a",
		DangerBg: "#d64545", DangerFg: "#ffffff", DangerBgH: "#e05555", Dot: "#7bd88f", Ok: "#7bd88f", Focus: "#ff9f43",
		SwBg: "#161718", SwOnBg: "#3a2a16", SwKnob: "#ff9f43", SwOnLine: "#ff9f43",
		KeyFg: "#ff9f43", KeyLine: "#45464a", TagFg: "#ff9f43", TagLine: "#5a4a33", TagBg: "#161718",
		Label: "#f3f1ea", Brand: "", Scrim: "rgba(0,0,0,.65)"},

}

var colourChoice = []string{"green", "amber", "blue", "pink"}

var skins = []Skin{
	{ID: "terminal", Palette: "green", Colours: true,
		FontCSS: `"IBM Plex Mono",Consolas,monospace`, FontGDI: "IBM Plex Mono",
		PagePx: 14, FontPx: 15, Weight: 400, BrandLS: ".18em", Caps: true, Flicker: true,
		FieldPad: "6px 10px", CtlFS: "12.5px", WeightB: 700,
		Radius: 0, BarR: "0", Border: 1, Round: false,
		Glow: true, Scan: 1, Shadow: "none",
		Level: "bars", Pulse: 1, Mark: "mic", Flash: "blink", Brackets: true},

	{ID: "editor", Palette: "editor",
		FontCSS: `"IBM Plex Mono",Consolas,monospace`, FontGDI: "IBM Plex Mono",
		PagePx: 13, FontPx: 14, Weight: 400, BrandLS: ".02em",
		FieldPad: "6px 11px", CtlFS: "12.5px", WeightB: 600,
		Radius: 8, BarR: "1px", Border: 1, Round: true,
		DotR: "50%", BadgeR: "6px", PanelR: "8px", SwitchR: "999px",
		BtnR: "6px", FieldR: "6px", CardR: "8px", KeyR: "6px",
		Glow: false, Scan: 0, Shadow: "0 10px 30px rgba(0,0,0,.45)",
		Level: "flat", Pulse: 1.5, Mark: "mic", Flash: "none"},

	{ID: "neon", Palette: "neon",
		FontCSS: `"IBM Plex Mono",Consolas,monospace`, FontGDI: "IBM Plex Mono",
		PagePx: 15, FontPx: 16, Weight: 400, BrandLS: ".04em",
		FieldPad: "8px 14px", CtlFS: "13.5px", WeightB: 600,
		Radius: 10, BarR: "3px", Border: 1, Round: true,
		DotR: "50%", BadgeR: "999px", PanelR: "14px", SwitchR: "999px", CtlH: "36px",
		BtnR: "10px", FieldR: "10px", CardR: "14px", KeyR: "10px",
		Glow: true, Scan: 0.18, Shadow: "0 18px 46px rgba(120,40,220,.35)",
		Level: "bars", Pulse: 0.8, Mark: "mic", Flash: "glow"},


	{ID: "fluent", Palette: "fluent",
		FontCSS: `"IBM Plex Mono",Consolas,monospace`, FontGDI: "IBM Plex Mono",
		PagePx: 13, FontPx: 14, Weight: 400, BrandLS: ".02em",
		FieldPad: "6px 11px", CtlFS: "12.5px", WeightB: 600,
		Radius: 8, BarR: "2px", Border: 1, Round: true,
		DotR: "50%", BadgeR: "4px", PanelR: "8px", SwitchR: "999px",
		BtnR: "4px", FieldR: "4px", CardR: "8px", KeyR: "4px",
		Glow: false, Scan: 0, Shadow: "0 32px 64px rgba(0,0,0,.5)",
		Level: "bars", Pulse: 1.4, Mark: "mic", Flash: "none"},

	{ID: "studio", Palette: "studio",
		FontCSS: `"IBM Plex Mono",Consolas,monospace`, FontGDI: "IBM Plex Mono",
		PagePx: 13, FontPx: 14, Weight: 400, BrandLS: ".02em",
		FieldPad: "6px 11px", CtlFS: "12.5px", WeightB: 600,
		Radius: 4, BarR: "1px", Border: 1, Round: true, SmallR: true,
		DotR: "50%", BadgeR: "3px", PanelR: "4px", SwitchR: "999px",
		BtnR: "3px", FieldR: "3px", CardR: "4px", KeyR: "3px",
		Glow: false, Scan: 0, Shadow: "0 24px 60px rgba(0,0,0,.6)",
		Level: "bars", Pulse: 1.4, Mark: "mic", Flash: "none"},

	{ID: "soft", Palette: "soft",
		FontCSS: `"IBM Plex Mono",Consolas,monospace`, FontGDI: "IBM Plex Mono",
		PagePx: 15, FontPx: 16, Weight: 400, BrandLS: ".02em",
		FieldPad: "8px 14px", CtlFS: "13.5px", WeightB: 600,
		Radius: 20, BarR: "99px", Border: 1, Round: true,
		DotR: "50%", BadgeR: "999px", PanelR: "20px", SwitchR: "999px", CtlH: "36px",
		BtnR: "14px", FieldR: "14px", CardR: "20px", KeyR: "14px",
		Glow: false, Scan: 0, Shadow: "0 14px 34px rgba(124,92,255,.18)",
		Level: "dots", Pulse: 1.2, Mark: "mic", Flash: "bounce"},

	{ID: "paper", Palette: "paper",
		FontCSS: `"IBM Plex Mono",Consolas,monospace`, FontGDI: "IBM Plex Mono",
		PagePx: 14, FontPx: 15, Weight: 400, BrandLS: "-.01em",
		FieldPad: "7px 11px", CtlFS: "12.5px", WeightB: 600,
		Radius: 8, BarR: "1px", Border: 1, Round: true,
		DotR: "50%", BadgeR: "4px", PanelR: "8px", SwitchR: "999px",
		BtnR: "6px", FieldR: "6px", CardR: "8px", KeyR: "6px",
		Glow: false, Scan: 0, Shadow: "0 12px 30px rgba(40,35,20,.14)",
		Level: "bars", Pulse: 1.4, Mark: "mic", Flash: "none"},

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
			return p.filled()
		}
	}
	return GetPalette(DefaultPalette)
}

func (p Palette) filled() Palette {
	if p.Off == "" {
		p.Off = grey(0.45*luma(p.Text) + 0.55*luma(p.Bg))
	}
	if p.Ok == "" {
		p.Ok = p.Accent
	}
	if p.BadBg == "" {
		p.BadBg = blend(p.Bad, p.Bg, 0.80)
	}
	if p.BadLine == "" {
		p.BadLine = blend(p.Bad, p.Bg, 0.52)
	}
	if p.BtnBgH == "" {
		p.BtnBgH = p.BtnBg
	}
	if p.Btn2Bg == "" {
		p.Btn2Bg = "transparent"
	}
	if p.Btn2Fg == "" {
		p.Btn2Fg = p.Dim
	}
	if p.Btn2Line == "" {
		p.Btn2Line = p.BtnLine
	}
	if p.Btn2BgH == "" {
		p.Btn2BgH = p.Btn2Bg
	}
	if p.Focus == "" {
		p.Focus = p.Dim
	}
	if p.Card == "" {
		p.Card = p.Field
	}
	if p.DangerBg == "" {
		p.DangerBg = "transparent"
	}
	if p.DangerFg == "" {
		p.DangerFg = p.Bad
	}
	if p.DangerBgH == "" {
		p.DangerBgH = p.BadBg
	}
	if p.Dot == "" {
		p.Dot = p.Text
	}
	return p
}

func (p Palette) Light() bool { return luma(p.Bg) > 140 }

func luma(hex string) float64 {
	r, g, b := RGB(hex)
	return 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)
}

func blend(fg, bg string, t float64) string {
	fr, fg2, fb := RGB(fg)
	br, bg2, bb := RGB(bg)
	mix := func(a, b uint8) uint8 { return uint8(float64(a)*(1-t) + float64(b)*t + 0.5) }
	return "#" + hex2(mix(fr, br)) + hex2(mix(fg2, bg2)) + hex2(mix(fb, bb))
}

func grey(v float64) string {
	if v < 0 {
		v = 0
	}
	if v > 255 {
		v = 255
	}
	c := hex2(uint8(v + 0.5))
	return "#" + c + c + c
}

func hex2(v uint8) string {
	const digits = "0123456789abcdef"
	return string([]byte{digits[v>>4], digits[v&0x0f]})
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
		return "terminal", old
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
	p := l.Palette.filled()
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
	warnGlow := "none"
	badGlow := "none"
	badFilter := "none"
	if l.Glow {
		hiGlow = "0 0 8px rgba(" + hirgb + ",.6)"
		iconGlow = "drop-shadow(0 0 6px rgba(" + hirgb + ",.7))"
		wr, wg, wb := RGB(p.Warn)
		warnrgb := itoa(int(wr)) + "," + itoa(int(wg)) + "," + itoa(int(wb))
		br2, bg2, bb2 := RGB(p.Bad)
		badrgb := itoa(int(br2)) + "," + itoa(int(bg2)) + "," + itoa(int(bb2))
		warnGlow = "0 0 6px rgba(" + warnrgb + ",.5)"
		badGlow = "0 0 7px rgba(" + badrgb + ",.5)"
		badFilter = "drop-shadow(0 0 4px rgba(" + badrgb + ",.5))"
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
	barr := l.BarR
	if barr == "" {
		barr = "0"
	}
	markMic, markFace := "block", "none"
	if l.Mark == "face" {
		markMic, markFace = "none", "block"
	}
	lvlw, lvlr := "4px", barr
	switch l.Level {
	case "dots":
		lvlw, lvlr = "10px", "50%"
	case "flat":
		lvlw, lvlr = "3px", "0"
	}
	scheme := "dark"
	if p.Light() {
		scheme = "light"
	}
	brandBg, brandClip, brandFill := "none", "border-box", "currentColor"
	if p.Brand != "" {
		brandBg, brandClip, brandFill = p.Brand, "text", "transparent"
	}
	label := p.Label
	if label == "" {
		label = p.Text
	}
	labelGlow := "none"
	if l.Glow {
		lr, lg, lb := RGB(label)
		labelGlow = "0 0 7px rgba(" + itoa(int(lr)) + "," + itoa(int(lg)) + "," + itoa(int(lb)) + ",.4)"
	}
	dotR := l.DotR
	if dotR == "" {
		dotR = "0"
	}
	badgeR := l.BadgeR
	if badgeR == "" {
		badgeR = "calc(" + itoa(int(l.Radius)) + "px * .4)"
	}
	panelR := l.PanelR
	if panelR == "" {
		panelR = itoa(int(l.Radius)) + "px"
	}
	switchR := l.SwitchR
	if switchR == "" {
		switchR = "calc(" + itoa(int(l.Radius)) + "px * .8)"
	}
	ctlH := l.CtlH
	if ctlH == "" {
		ctlH = "30px"
	}
	btnR := l.BtnR
	if btnR == "" {
		btnR = "calc(" + itoa(int(l.Radius)) + "px * .5)"
	}
	fieldR := l.FieldR
	if fieldR == "" {
		fieldR = "calc(" + itoa(int(l.Radius)) + "px * .55)"
	}
	cardR := l.CardR
	if cardR == "" {
		cardR = "calc(" + itoa(int(l.Radius)) + "px * .6)"
	}
	keyR := l.KeyR
	if keyR == "" {
		keyR = "calc(" + itoa(int(l.Radius)) + "px * .6)"
	}
	keyFg := p.KeyFg
	if keyFg == "" {
		keyFg = p.Text
	}
	keyLine := p.KeyLine
	if keyLine == "" {
		keyLine = p.Line
	}
	tagFg := p.TagFg
	if tagFg == "" {
		tagFg = label
	}
	tagLine := p.TagLine
	if tagLine == "" {
		tagLine = p.Line
	}
	tagBg := p.TagBg
	if tagBg == "" {
		tagBg = "transparent"
	}
	swBg := p.SwBg
	if swBg == "" {
		swBg = "transparent"
	}
	swOnBg := p.SwOnBg
	if swOnBg == "" {
		swOnBg = "transparent"
	}
	swKnob := p.SwKnob
	if swKnob == "" {
		swKnob = p.Accent
	}
	swOnLine := p.SwOnLine
	if swOnLine == "" {
		swOnLine = p.Dim
	}
	btnBo, btnBc := `""`, `""`
	if l.Brackets {
		btnBo, btnBc = "\"[ \"", "\" ]\""
	}
	return "--wborder:" + wborder +
		";--bg:" + p.Bg + ";--panel:" + p.Panel + ";--line:" + p.Line +
		";--green:" + p.Text + ";--hi:" + p.Accent +
		";--dim:" + p.Dim + ";--faint:" + p.Faint +
		";--amber:" + p.Warn + ";--bad:" + p.Bad + ";--rec:" + p.Rec +
		";--rgb:" + rgb +
		";--field:" + p.Field + ";--card:" + p.Card + ";--soft:" + p.Soft +
		";--navon:" + p.NavOn + ";--on:" + p.On +
		";--titlebg:" + p.TitleBg + ";--sidebg:" + p.SideBg + ";--keybg:" + p.KeyBg +
		";--btnbg:" + p.BtnBg + ";--btnfg:" + p.BtnFg + ";--btnline:" + p.BtnLine +
		";--btnbgh:" + p.BtnBgH + ";--btn2bg:" + p.Btn2Bg + ";--btn2fg:" + p.Btn2Fg +
		";--dangerbg:" + p.DangerBg + ";--dangerfg:" + p.DangerFg + ";--dangerbgh:" + p.DangerBgH +
		";--switchr:" + switchR + ";--ctlh:" + ctlH +
		";--btnr:" + btnR + ";--fieldr:" + fieldR + ";--cardr:" + cardR + ";--keyr:" + keyR + ";--keyfg:" + keyFg + ";--keyline:" + keyLine +
		";--tagfg:" + tagFg + ";--tagline:" + tagLine + ";--tagbg:" + tagBg +
		";--swbg:" + swBg + ";--swonbg:" + swOnBg + ";--swknob:" + swKnob + ";--swonline:" + swOnLine +
		";--btn2line:" + p.Btn2Line + ";--btn2bgh:" + p.Btn2BgH + ";--focus:" + p.Focus +
		";--dotr:" + dotR + ";--badger:" + badgeR + ";--panelr:" + panelR +
		";--btnbo:" + btnBo + ";--btnbc:" + btnBc +
		";--lbl:" + label + ";--lblglow:" + labelGlow +
		";--selbg:" + p.SelBg + ";--selfg:" + p.SelFg +
		";--brandbg:" + brandBg + ";--brandclip:" + brandClip + ";--brandfill:" + brandFill +
		";--scrim:" + p.Scrim +
		";--ok:" + p.Ok + ";--scheme:" + scheme +
		";--badbg:" + p.BadBg + ";--badline:" + p.BadLine +
		";--lvlw:" + lvlw + ";--lvlr:" + lvlr +
		";--markmic:" + markMic + ";--markface:" + markFace +
		";--glow:" + glow + ";--higlow:" + hiGlow + ";--iconglow:" + iconGlow +
		";--amberglow:" + warnGlow + ";--badglow:" + badGlow + ";--badfilter:" + badFilter +
		";--font:" + l.FontCSS +
		";--fs:" + itoa(int(l.PagePx)) + "px" +
		";--caps:" + caps + ";--ls:" + ls + ";--flicker:" + flicker +
		";--fieldpad:" + l.FieldPad + ";--ctlfs:" + l.CtlFS +
		";--wb:" + itoa(int(l.WeightB)) +
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
