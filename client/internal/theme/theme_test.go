package theme

import (
	"strings"
	"testing"
)

func TestFourThemesAndTheDefault(t *testing.T) {
	ids := IDs()
	if len(ids) != 4 {
		t.Fatalf("IDs() = %v, want four themes", ids)
	}
	if ids[0] != Default {
		t.Errorf("the first theme is %q, want the default %q", ids[0], Default)
	}
	for _, id := range ids {
		if !Valid(id) {
			t.Errorf("Valid(%q) = false", id)
		}
	}
	if Valid("chartreuse") {
		t.Error("Valid accepted a theme that does not exist")
	}
	if Get("chartreuse").ID != Default {
		t.Error("Get falls back to something other than the default")
	}
}

func TestEveryPaletteIsComplete(t *testing.T) {
	for _, id := range IDs() {
		p := Get(id)
		for name, hex := range map[string]string{
			"Bg": p.Bg, "Panel": p.Panel, "Line": p.Line,
			"Accent": p.Accent, "Dim": p.Dim, "Faint": p.Faint,
			"Warn": p.Warn, "Bad": p.Bad,
		} {
			if len(hex) != 7 || !strings.HasPrefix(hex, "#") {
				t.Errorf("%s.%s = %q, want #rrggbb", id, name, hex)
			}
		}
		if p.Accent == p.Warn || p.Accent == p.Bad {
			t.Errorf("%s: the accent is not told apart from a warning or an error", id)
		}
	}
}

func TestRGBReadsTheChannels(t *testing.T) {
	r, g, b := RGB("#3cff6e")
	if r != 0x3c || g != 0xff || b != 0x6e {
		t.Errorf("RGB(#3cff6e) = %d,%d,%d", r, g, b)
	}
	if r, g, b := RGB("nonsense"); r != 0 || g != 0 || b != 0 {
		t.Errorf("RGB(nonsense) = %d,%d,%d, want zeros", r, g, b)
	}
}

func TestCSSVarsCarryEveryColour(t *testing.T) {
	css := Get("amber").CSSVars()
	for _, want := range []string{"--bg:", "--panel:", "--line:", "--green:#ff9e2c", "--dim:", "--faint:", "--amber:", "--bad:", "--glow:0 0 7px rgba(255,158,44,.55)"} {
		if !strings.Contains(css, want) {
			t.Errorf("CSSVars() misses %q: %s", want, css)
		}
	}
}
