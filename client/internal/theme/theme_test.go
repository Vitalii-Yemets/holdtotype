package theme

import (
	"strings"
	"testing"
)

func TestSkinsAndColoursAreSeparate(t *testing.T) {
	skins := SkinIDs()
	if len(skins) != 3 {
		t.Fatalf("SkinIDs() = %v, want three skins", skins)
	}
	if skins[0] != DefaultSkin {
		t.Errorf("the first skin is %q, want the default %q", skins[0], DefaultSkin)
	}
	colours := ColourIDs("terminal")
	if len(colours) != 4 || colours[0] != DefaultPalette {
		t.Errorf("ColourIDs(terminal) = %v, want the four with green first", colours)
	}
	for _, id := range []string{"editor", "neon"} {
		if got := ColourIDs(id); len(got) != 0 {
			t.Errorf("ColourIDs(%q) = %v, want none — the skin carries its own", id, got)
		}
	}
}

func TestASkinWithOwnColoursIgnoresTheChoice(t *testing.T) {
	if got := Current("editor", "pink").Palette.ID; got != "editor" {
		t.Errorf("Current(editor, pink) drawn with %q, want editor's own", got)
	}
	if got := Current("terminal", "pink").Palette.ID; got != "pink" {
		t.Errorf("Current(terminal, pink) drawn with %q, want pink", got)
	}
	if got := Current("terminal", "nonsense").Palette.ID; got != DefaultPalette {
		t.Errorf("an unknown colour gave %q, want the default", got)
	}
	if got := Current("nonsense", "amber"); got.ID != DefaultSkin || got.Palette.ID != "amber" {
		t.Errorf("an unknown skin gave %q/%q", got.ID, got.Palette.ID)
	}
}

func TestMigrateSplitsTheOldSingleValue(t *testing.T) {
	cases := map[string][2]string{
		"green":  {"terminal", "green"},
		"amber":  {"terminal", "amber"},
		"blue":   {"terminal", "blue"},
		"pink":   {"terminal", "pink"},
		"editor": {"editor", "green"},
		"neon":   {"neon", "green"},
		"":       {"terminal", "green"},
		"what":   {"terminal", "green"},
	}
	for old, want := range cases {
		skin, colour := Migrate(old)
		if skin != want[0] || colour != want[1] {
			t.Errorf("Migrate(%q) = %q/%q, want %q/%q", old, skin, colour, want[0], want[1])
		}
	}
}

func TestEveryPaletteIsComplete(t *testing.T) {
	for _, p := range palettes {
		for name, hex := range map[string]string{
			"Bg": p.Bg, "Panel": p.Panel, "Line": p.Line,
			"Accent": p.Accent, "Dim": p.Dim, "Faint": p.Faint,
			"Warn": p.Warn, "Bad": p.Bad,
			"Field": p.Field, "Soft": p.Soft, "NavOn": p.NavOn, "On": p.On,
		} {
			if len(hex) != 7 || !strings.HasPrefix(hex, "#") {
				t.Errorf("%s.%s = %q, want #rrggbb", p.ID, name, hex)
			}
		}
		if p.Accent == p.Warn || p.Accent == p.Bad {
			t.Errorf("%s: the accent is not told apart from a warning or an error", p.ID)
		}
	}
}

func TestEverySkinIsComplete(t *testing.T) {
	for _, id := range SkinIDs() {
		s := GetSkin(id)
		if s.FontCSS == "" || s.FontGDI == "" {
			t.Errorf("%s: no font for the page or for the plate", id)
		}
		if s.FontPx < 11 || s.FontPx > 22 {
			t.Errorf("%s: plate text is %d px", id, s.FontPx)
		}
		if s.Weight != 400 && s.Weight != 600 {
			t.Errorf("%s: weight %d, want 400 or 600", id, s.Weight)
		}
		switch s.Level {
		case "bars", "flat", "dots":
		default:
			t.Errorf("%s: level meter %q is not one we can draw", id, s.Level)
		}
		if s.Pulse <= 0 || s.Scan < 0 || s.Scan > 1 {
			t.Errorf("%s: pulse %v, scanlines %v", id, s.Pulse, s.Scan)
		}
		if !ValidColour(s.Palette) && s.Palette != id {
			t.Errorf("%s: its own palette is %q, which is neither a colour nor its namesake", id, s.Palette)
		}
	}
}

func TestTheThreeSkinsDifferInMoreThanColour(t *testing.T) {
	term, editor, neon := GetSkin("terminal"), GetSkin("editor"), GetSkin("neon")
	if term.FontGDI == neon.FontGDI {
		t.Error("Terminal and Neon are drawn with the same font")
	}
	if editor.Glow || editor.Scan != 0 {
		t.Error("Editor should be flat: no halo, no scanlines")
	}
	if term.Radius != 0 || neon.Radius == 0 || term.Round || !neon.Round {
		t.Error("Terminal should be square and Neon rounded")
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
	css := Current("editor", "").CSSVars()
	for _, want := range []string{
		"--bg:#1e1e1e", "--green:#4fc1ff", "--rgb:79,193,255",
		"--glow:none", "--font:", "--r:3px", "--bw:1px", "--scan:0", "--shadow:0 10px 30px",
	} {
		if !strings.Contains(css, want) {
			t.Errorf("CSSVars() misses %q: %s", want, css)
		}
	}
	if !strings.Contains(Current("terminal", "amber").CSSVars(), "--green:#ff9e2c") {
		t.Error("the terminal skin did not take the amber colour")
	}
	if !strings.Contains(Current("terminal", "green").CSSVars(), "--glow:0 0 7px rgba(60,255,110,.55)") {
		t.Error("the terminal skin lost its halo")
	}
}
