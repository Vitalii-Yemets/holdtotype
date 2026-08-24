package theme

import "strings"

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
}

const Default = "green"

var palettes = []Palette{
	{ID: "green", Bg: "#0b0f0c", Panel: "#0e1410", Line: "#1d4a2b",
		Accent: "#3cff6e", Dim: "#20a34a", Faint: "#14803a", Warn: "#ffb347", Bad: "#ff7b6b"},
	{ID: "amber", Bg: "#100c0a", Panel: "#17110d", Line: "#4a3018",
		Accent: "#ff9e2c", Dim: "#b56a12", Faint: "#8a4f0d", Warn: "#ffd24a", Bad: "#ff6b5b"},
	{ID: "blue", Bg: "#0b0e10", Panel: "#0e1317", Line: "#1d3a4a",
		Accent: "#4cc3ff", Dim: "#1c7fb8", Faint: "#14608f", Warn: "#ffb347", Bad: "#ff7b6b"},
	{ID: "pink", Bg: "#100b0e", Panel: "#170e14", Line: "#4a1d3a",
		Accent: "#ff6ec7", Dim: "#b82f86", Faint: "#8f2467", Warn: "#ffb347", Bad: "#ff6b6b"},
}

func IDs() []string {
	out := make([]string, 0, len(palettes))
	for _, p := range palettes {
		out = append(out, p.ID)
	}
	return out
}

func Valid(id string) bool {
	for _, p := range palettes {
		if p.ID == id {
			return true
		}
	}
	return false
}

func Get(id string) Palette {
	for _, p := range palettes {
		if p.ID == id {
			return p
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

// CSSVars renders the palette as the contents of a :root rule.
func (p Palette) CSSVars() string {
	r, g, b := RGB(p.Accent)
	glow := "0 0 7px rgba(" + itoa(int(r)) + "," + itoa(int(g)) + "," + itoa(int(b)) + ",.55)"
	rgb := itoa(int(r)) + "," + itoa(int(g)) + "," + itoa(int(b))
	return "--bg:" + p.Bg + ";--panel:" + p.Panel + ";--line:" + p.Line +
		";--green:" + p.Accent + ";--dim:" + p.Dim + ";--faint:" + p.Faint +
		";--amber:" + p.Warn + ";--bad:" + p.Bad + ";--rgb:" + rgb + ";--glow:" + glow
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	digits := ""
	for n > 0 {
		digits = string(rune('0'+n%10)) + digits
		n /= 10
	}
	return digits
}
