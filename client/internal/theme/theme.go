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
	DotR    string
	BadgeR  string
	PanelR  string
	SwitchR string
	CtlH    string

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
		Label: "#f2fff5", Brand: "", Scrim: "rgba(3,7,4,.78)"},

	{ID: "amber", Bg: "#100c0a", Panel: "#17110d", Line: "#4a3018",
		Text: "#ff9e2c", Accent: "#ff9e2c", Dim: "#b56a12", Faint: "#8a4f0d",
		Warn: "#ffd24a", Bad: "#ff6b5b", Rec: "#ff5b4d",
		Field: "#120c07", Soft: "#2a1a0d", NavOn: "#22160c", On: "#402611",
		TitleBg: "transparent", SideBg: "transparent", KeyBg: "transparent",
		BtnBg: "#22160c", BtnFg: "#ff9e2c", BtnLine: "#b56a12", SelBg: "#ff9e2c", SelFg: "#100c0a",
		Label: "#ffe9c9", Brand: "", Scrim: "rgba(8,5,3,.78)"},

	{ID: "blue", Bg: "#0b0e10", Panel: "#0e1317", Line: "#1d3a4a",
		Text: "#4cc3ff", Accent: "#4cc3ff", Dim: "#1c7fb8", Faint: "#14608f",
		Warn: "#ffb347", Bad: "#ff7b6b", Rec: "#ff5b4d",
		Field: "#070f14", Soft: "#12222c", NavOn: "#101c24", On: "#123a52",
		TitleBg: "transparent", SideBg: "transparent", KeyBg: "transparent",
		BtnBg: "#101c24", BtnFg: "#4cc3ff", BtnLine: "#1c7fb8", SelBg: "#4cc3ff", SelFg: "#0b0e10",
		Label: "#e4f6ff", Brand: "", Scrim: "rgba(3,6,8,.78)"},

	{ID: "pink", Bg: "#100b0e", Panel: "#170e14", Line: "#4a1d3a",
		Text: "#ff6ec7", Accent: "#ff6ec7", Dim: "#b82f86", Faint: "#8f2467",
		Warn: "#ffb347", Bad: "#ff6b6b", Rec: "#ff5b4d",
		Field: "#120810", Soft: "#2a1222", NavOn: "#22101c", On: "#40183a",
		TitleBg: "transparent", SideBg: "transparent", KeyBg: "transparent",
		BtnBg: "#22101c", BtnFg: "#ff6ec7", BtnLine: "#b82f86", SelBg: "#ff6ec7", SelFg: "#100b0e",
		Label: "#ffe6f4", Brand: "", Scrim: "rgba(8,3,6,.78)"},

	{ID: "editor", Bg: "#1e1e1e", Panel: "#252526", Line: "#3c3c3c",
		Text: "#d4d4d4", Accent: "#4fc1ff", Dim: "#9d9d9d", Faint: "#6e6e6e",
		Warn: "#cca700", Bad: "#f14c4c", Rec: "#f14c4c",
		Field: "#3c3c3c", Card: "#252526", Soft: "#2d2d2d", NavOn: "#37373d", On: "#094771",
		TitleBg: "#323233", SideBg: "#252526", KeyBg: "#3c3c3c",
		BtnBg: "#0e639c", BtnFg: "#ffffff", BtnLine: "#0e639c", SelBg: "#0e639c", SelFg: "#ffffff",
		BtnBgH: "#1177bb", Btn2Bg: "#0e639c", Btn2Fg: "#ffffff", Btn2Line: "transparent", Btn2BgH: "#1177bb",
		DangerBg: "#a1260d", DangerFg: "#ffffff", DangerBgH: "#c42b1c", Dot: "#89d185",
		Focus: "#007fd4",
		Label: "#e0e0e0", Brand: "", Scrim: "rgba(0,0,0,.6)"},

	{ID: "neon", Bg: "#150a22", Panel: "#1d0e30", Line: "#4a2472",
		Text: "#f3b6e4", Accent: "#46e0ff", Dim: "#b06ee0", Faint: "#7d4fae",
		Warn: "#ffd24a", Bad: "#ff4d7d", Rec: "#ff4d7d",
		Field: "#1e0f33", Card: "#231039", Soft: "#2a1442", NavOn: "#2b1240", On: "#4a2472",
		TitleBg: "linear-gradient(90deg,#26103f,#1a0b2b)", SideBg: "#190c29", KeyBg: "linear-gradient(90deg,rgba(255,95,200,.18),rgba(70,224,255,.14))",
		BtnBg: "linear-gradient(90deg,rgba(255,95,200,.20),rgba(70,224,255,.16))", BtnFg: "#ffffff", BtnLine: "#4a2472",
		SelBg: "linear-gradient(90deg,#ff5fc8,#46e0ff)", SelFg: "#150a22",
		BtnBgH: "linear-gradient(90deg,rgba(255,95,200,.34),rgba(70,224,255,.28))",
		Btn2Bg: "linear-gradient(90deg,rgba(255,95,200,.20),rgba(70,224,255,.16))", Btn2Fg: "#ffffff", Btn2Line: "#4a2472",
		Btn2BgH: "linear-gradient(90deg,rgba(255,95,200,.34),rgba(70,224,255,.28))",
		DangerBg: "linear-gradient(90deg,rgba(255,77,125,.30),rgba(255,77,125,.14))", DangerFg: "#ffd7e3",
		DangerBgH: "linear-gradient(90deg,rgba(255,77,125,.46),rgba(255,77,125,.24))",
		Dot: "#46e0ff", Focus: "#46e0ff",
		Label: "#ffffff", Brand: "linear-gradient(90deg,#ff5fc8,#46e0ff)", Scrim: "rgba(10,4,18,.72)"},

	{ID: "soft", Bg: "#ddc9d5", Panel: "#e7d5df", Line: "#c9a5b9",
		Text: "#8f2f60", Accent: "#b8407c", Dim: "#9d6885", Faint: "#ae8b9e",
		Warn: "#7d4a00", Bad: "#ab2445", Rec: "#c74a6c",
		Field: "#f4e8ee", Card: "#eadbe4", Soft: "#d2b9c7", NavOn: "#d8c0ce", On: "#e6c6d8",
		TitleBg: "#e2cedb", SideBg: "#e4d2dd", KeyBg: "#e2cedb",
		BtnBg: "#c2467f", BtnFg: "#ffffff", BtnLine: "transparent", SelBg: "#c2467f", SelFg: "#ffffff",
		BtnBgH: "#ad3b71", Btn2Bg: "#d5b6c8", Btn2Fg: "#6d2349", Btn2Line: "#bf9bb0", Btn2BgH: "#cba9be",
		DangerBg: "#ab2445", DangerFg: "#ffffff", DangerBgH: "#8f1c39", Dot: "#3f8a58", Focus: "#b8407c",
		Label: "#6d2349", Brand: "", Scrim: "rgba(120,45,85,.35)"},

	{ID: "paper", Bg: "#e7eaee", Panel: "#f3f5f7", Line: "#ced5dc",
		Text: "#1f2328", Accent: "#0969da", Dim: "#59636e", Faint: "#818b98",
		Warn: "#9a6700", Bad: "#cf222e", Rec: "#cf222e",
		Field: "#ffffff", Card: "#f7f9fa", Soft: "#dde2e8", NavOn: "#dfe5eb", On: "#ddf4ff",
		TitleBg: "#eff2f5", SideBg: "#eef1f4", KeyBg: "#e4e9ee",
		BtnBg: "#0969da", BtnFg: "#ffffff", BtnLine: "#0969da", SelBg: "#0969da", SelFg: "#ffffff",
		BtnBgH: "#0860ca", Btn2Bg: "#d8dee5", Btn2Fg: "#24292f", Btn2Line: "#bcc5ce", Btn2BgH: "#c9d1da",
		DangerBg: "#cf222e", DangerFg: "#ffffff", DangerBgH: "#a40e26", Dot: "#1a7f37", Focus: "#0969da",
		Ok: "#1a7f37", Label: "#111418", Brand: "", Scrim: "rgba(31,35,40,.35)"},
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
		FontCSS: `"Segoe UI Variable Text","Segoe UI",system-ui,sans-serif`, FontGDI: "Segoe UI",
		PagePx: 13, FontPx: 14, Weight: 400, BrandLS: ".02em",
		FieldPad: "6px 11px", CtlFS: "12.5px", WeightB: 600,
		Radius: 8, BarR: "0", Border: 1, Round: false,
		DotR: "50%", BadgeR: "10px", PanelR: "6px", SwitchR: "999px",
		Glow: false, Scan: 0, Shadow: "0 10px 30px rgba(0,0,0,.45)",
		Level: "flat", Pulse: 1.5, Mark: "mic", Flash: "none"},

	{ID: "neon", Palette: "neon",
		FontCSS: `"IBM Plex Sans","Segoe UI",system-ui,sans-serif`, FontGDI: "IBM Plex Sans",
		PagePx: 15, FontPx: 16, Weight: 400, BrandLS: ".04em",
		FieldPad: "8px 14px", CtlFS: "13.5px", WeightB: 600,
		Radius: 14, BarR: "99px", Border: 1, Round: true,
		DotR: "50%", BadgeR: "999px", PanelR: "20px", SwitchR: "999px", CtlH: "36px",
		Glow: true, Scan: 0.18, Shadow: "0 18px 46px rgba(150,40,220,.35)",
		Level: "bars", Pulse: 0.8, Mark: "mic", Flash: "glow"},

	{ID: "soft", Palette: "soft",
		FontCSS: `"Comic Sans MS","Segoe UI Variable Display","Segoe UI",sans-serif`, FontGDI: "Comic Sans MS",
		PagePx: 15, FontPx: 16, Weight: 400, BrandLS: ".02em",
		FieldPad: "8px 14px", CtlFS: "13.5px", WeightB: 600,
		Radius: 16, BarR: "99px", Border: 1, Round: true,
		DotR: "50%", BadgeR: "999px", PanelR: "22px", SwitchR: "999px", CtlH: "36px",
		Glow: false, Scan: 0, Shadow: "0 14px 34px rgba(255,140,190,.28)",
		Level: "dots", Pulse: 1.2, Mark: "face", Flash: "bounce"},

	{ID: "paper", Palette: "paper",
		FontCSS: `"Segoe UI",system-ui,sans-serif`, FontGDI: "Segoe UI",
		PagePx: 14, FontPx: 15, Weight: 400, BrandLS: "-.01em",
		FieldPad: "7px 11px", CtlFS: "12.5px", WeightB: 600,
		Radius: 10, BarR: "2px", Border: 1, Round: true,
		DotR: "50%", BadgeR: "999px", PanelR: "12px", SwitchR: "999px",
		Glow: false, Scan: 0, Shadow: "0 8px 24px rgba(31,35,40,.12)",
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
