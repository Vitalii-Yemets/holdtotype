package theme

import (
	"strings"
	"testing"
)

func TestSixSkinsAndTheDefault(t *testing.T) {
	ids := IDs()
	if len(ids) != 6 {
		t.Fatalf("IDs() = %v, want six skins", ids)
	}
	if ids[0] != Default {
		t.Errorf("the first skin is %q, want the default %q", ids[0], Default)
	}
	for _, id := range ids {
		if !Valid(id) {
			t.Errorf("Valid(%q) = false", id)
		}
	}
	if Valid("chartreuse") {
		t.Error("Valid accepted a skin that does not exist")
	}
	if Get("chartreuse").ID != Default {
		t.Error("Get falls back to something other than the default")
	}
}

func TestEverySkinIsComplete(t *testing.T) {
	for _, id := range IDs() {
		s := Get(id)
		for name, hex := range map[string]string{
			"Bg": s.Bg, "Panel": s.Panel, "Line": s.Line,
			"Accent": s.Accent, "Dim": s.Dim, "Faint": s.Faint,
			"Warn": s.Warn, "Bad": s.Bad,
		} {
			if len(hex) != 7 || !strings.HasPrefix(hex, "#") {
				t.Errorf("%s.%s = %q, want #rrggbb", id, name, hex)
			}
		}
		if s.Accent == s.Warn || s.Accent == s.Bad {
			t.Errorf("%s: the accent is not told apart from a warning or an error", id)
		}
		if s.FontCSS == "" || s.FontGDI == "" {
			t.Errorf("%s: no font for the page or for the plate", id)
		}
		if s.FontPx < 11 || s.FontPx > 22 {
			t.Errorf("%s: plate text is %d px, want something readable", id, s.FontPx)
		}
		if s.Weight != 400 && s.Weight != 600 {
			t.Errorf("%s: weight %d, want 400 or 600", id, s.Weight)
		}
		switch s.Level {
		case "bars", "flat", "dots":
		default:
			t.Errorf("%s: level meter %q is not one we can draw", id, s.Level)
		}
		if s.Pulse <= 0 {
			t.Errorf("%s: pulse speed %v", id, s.Pulse)
		}
		if s.Scan < 0 || s.Scan > 1 {
			t.Errorf("%s: scanlines %v, want 0…1", id, s.Scan)
		}
	}
}

func TestTheDarkSkinsDifferInMoreThanColour(t *testing.T) {
	term, editor, neon := Get("green"), Get("editor"), Get("neon")
	if term.FontGDI == neon.FontGDI {
		t.Error("Terminal and Neon are drawn with the same font")
	}
	if editor.Glow || editor.Scan != 0 {
		t.Error("Editor should be flat: no halo, no scanlines")
	}
	if term.Radius != 0 || neon.Radius == 0 {
		t.Error("Terminal should be square and Neon rounded")
	}
	if term.Round || !neon.Round {
		t.Error("window corners: Terminal square, Neon rounded")
	}
	if editor.Level != "flat" {
		t.Error("Editor should get the quiet level meter")
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

func TestCSSVarsCarryEverythingThePageNeeds(t *testing.T) {
	css := Get("editor").CSSVars()
	for _, want := range []string{
		"--bg:#1e1e1e", "--green:#4fc1ff", "--rgb:79,193,255",
		"--glow:none", "--font:", "--r:3px", "--bw:1px", "--scan:0", "--shadow:0 10px 30px",
	} {
		if !strings.Contains(css, want) {
			t.Errorf("CSSVars() misses %q: %s", want, css)
		}
	}
	if !strings.Contains(Get("green").CSSVars(), "--glow:0 0 7px rgba(60,255,110,.55)") {
		t.Error("the green skin lost its halo")
	}
	if !strings.Contains(Get("neon").CSSVars(), "--scan:.35") {
		t.Error("the neon skin lost its scanlines")
	}
}
