package theme

import (
	"strings"
	"testing"
)

func TestSkinsAndColoursAreSeparate(t *testing.T) {
	skins := SkinIDs()
	if len(skins) != 5 {
		t.Fatalf("SkinIDs() = %v, want five skins", skins)
	}
	if skins[0] != DefaultSkin {
		t.Errorf("the first skin is %q, want the default %q", skins[0], DefaultSkin)
	}
	colours := ColourIDs("terminal")
	if len(colours) != 4 || colours[0] != DefaultPalette {
		t.Errorf("ColourIDs(terminal) = %v, want the four with green first", colours)
	}
	for _, id := range []string{"editor", "neon", "soft", "paper"} {
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
			"Text": p.Text, "Accent": p.Accent, "Dim": p.Dim, "Faint": p.Faint,
			"Warn": p.Warn, "Bad": p.Bad, "Rec": p.Rec,
			"Field": p.Field, "Soft": p.Soft, "NavOn": p.NavOn, "On": p.On,
			"BtnFg": p.BtnFg,
		} {
			if len(hex) != 7 || !strings.HasPrefix(hex, "#") {
				t.Errorf("%s.%s = %q, want #rrggbb", p.ID, name, hex)
			}
		}
		for name, v := range map[string]string{
			"TitleBg": p.TitleBg, "SideBg": p.SideBg, "KeyBg": p.KeyBg,
			"BtnBg": p.BtnBg, "BtnLine": p.BtnLine, "Scrim": p.Scrim,
		} {
			if v == "" {
				t.Errorf("%s.%s is empty", p.ID, name)
			}
		}
		if p.Text == p.Warn || p.Text == p.Bad {
			t.Errorf("%s: the text is not told apart from a warning or an error", p.ID)
		}
	}
}

func TestTheQuietSkinsSeparateTextFromAccent(t *testing.T) {
	for _, id := range []string{"editor", "neon", "paper"} {
		p := GetPalette(id)
		if p.Text == p.Accent {
			t.Errorf("%s: text and accent are both %q", id, p.Text)
		}
	}
	if p := GetPalette("green"); p.Text != p.Accent {
		t.Errorf("the terminal green is one colour: text %q, accent %q", p.Text, p.Accent)
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
		if s.PagePx < 11 || s.PagePx > 18 {
			t.Errorf("%s: page text is %d px", id, s.PagePx)
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

func TestOnlyTheTerminalShouts(t *testing.T) {
	if !GetSkin("terminal").Caps || !GetSkin("terminal").Flicker {
		t.Error("the terminal skin lost its capitals or its flicker")
	}
	for _, id := range []string{"editor", "neon", "soft", "paper"} {
		s := GetSkin(id)
		if s.Caps || s.Flicker {
			t.Errorf("%s: capitals %v, flicker %v — both should be off", id, s.Caps, s.Flicker)
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

func TestNothingIsDrawnInABoldFace(t *testing.T) {
	for _, id := range SkinIDs() {
		if w := GetSkin(id).Weight; w > 400 {
			t.Errorf("%s: the plate is drawn at weight %d — heavy type reads as shouting", id, w)
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

func TestCSSVarsCarryEverythingThePageNeeds(t *testing.T) {
	css := Current("editor", "").CSSVars()
	for _, want := range []string{
		"--bg:#1e1e1e", "--green:#d4d4d4", "--hi:#4fc1ff", "--rgb:212,212,212",
		"--glow:none", "--font:", "--fs:13px", "--r:3px", "--bw:1px", "--scan:0",
		"--shadow:0 10px 30px", "--titlebg:#323233", "--keybg:#3c3c3c",
		"--btnbg:#0e639c", "--btnfg:#ffffff", "--caps:none", "--ls:0", "--flicker:none",
	} {
		if !strings.Contains(css, want) {
			t.Errorf("CSSVars() misses %q: %s", want, css)
		}
	}
	term := Current("terminal", "green").CSSVars()
	for _, want := range []string{"--caps:uppercase", "--ls:1px", "--flicker:flicker 6s infinite", "--brandbg:none"} {
		if !strings.Contains(term, want) {
			t.Errorf("the terminal skin misses %q", want)
		}
	}
	if !strings.Contains(Current("terminal", "amber").CSSVars(), "--green:#ff9e2c") {
		t.Error("the terminal skin did not take the amber colour")
	}
	if !strings.Contains(term, "--glow:0 0 7px rgba(60,255,110,.55)") {
		t.Error("the terminal skin lost its halo")
	}
	neon := Current("neon", "").CSSVars()
	if !strings.Contains(neon, "--brandclip:text") || !strings.Contains(neon, "--brandfill:transparent") {
		t.Error("the neon title should be painted with its gradient")
	}
}

func TestSkinsNameTheFacesTheyShipWith(t *testing.T) {
	for _, c := range []struct{ skin, want string }{
		{"terminal", "IBM Plex Mono"},
		{"neon", "IBM Plex Sans"},
	} {
		look := Current(c.skin, "")
		if !strings.HasPrefix(look.FontCSS, `"`+c.want+`"`) {
			t.Errorf("the %s page asks for %q, want %q first", c.skin, look.FontCSS, c.want)
		}
		if look.FontGDI != c.want {
			t.Errorf("the %s windows are drawn with %q, want %q", c.skin, look.FontGDI, c.want)
		}
	}
	ed := Current("editor", "")
	if ed.FontGDI != "Cascadia Mono" {
		t.Errorf("the editor skin keeps its own face, got %q", ed.FontGDI)
	}
}

func TestTheLightSkinsAreLight(t *testing.T) {
	for _, id := range []string{"soft", "paper"} {
		p := Current(id, "").Palette
		if !p.Light() {
			t.Errorf("the %s skin stands on %q, which is not a light ground", id, p.Bg)
		}
		if luma(p.Text) > 110 {
			t.Errorf("the %s skin writes in %q, too pale for its ground", id, p.Text)
		}
		if luma(p.Bg)-luma(p.Text) < 120 {
			t.Errorf("the %s skin has too little between its ink and its paper", id)
		}
	}
	for _, id := range []string{"terminal", "editor", "neon"} {
		if Current(id, "").Palette.Light() {
			t.Errorf("the %s skin should stay dark", id)
		}
	}
}

func TestEveryPaletteFillsInTheDerivedColours(t *testing.T) {
	for _, id := range []string{"green", "amber", "blue", "pink", "editor", "neon", "soft", "paper"} {
		p := GetPalette(id)
		for name, v := range map[string]string{"Off": p.Off, "BadBg": p.BadBg, "BadLine": p.BadLine} {
			if len(v) != 7 || v[0] != '#' {
				t.Errorf("palette %s left %s as %q", id, name, v)
			}
		}
		if luma(p.BadBg) > luma(p.Bad) && !p.Light() {
			t.Errorf("palette %s got a hover ground brighter than the warning itself", id)
		}
	}
}

func TestALightSkinAsksForNoHalo(t *testing.T) {
	for _, id := range []string{"soft", "paper"} {
		look := Current(id, "")
		if look.Glow {
			t.Errorf("the %s skin should not glow", id)
		}
		if look.Scan != 0 {
			t.Errorf("the %s skin should have no scanlines", id)
		}
		css := look.CSSVars()
		for _, want := range []string{"--glow:none", "--higlow:none", "--amberglow:none", "--badglow:none", "--badfilter:none", "--scan:0"} {
			if !strings.Contains(css, want) {
				t.Errorf("the %s skin misses %q", id, want)
			}
		}
		if !strings.Contains(css, "--badbg:#") || !strings.Contains(css, "--badline:#") {
			t.Errorf("the %s skin did not fill the warning surfaces: %s", id, css)
		}
	}
	term := Current("terminal", "green").CSSVars()
	for _, want := range []string{"--amberglow:0 0 6px", "--badglow:0 0 7px", "--badfilter:drop-shadow"} {
		if !strings.Contains(term, want) {
			t.Errorf("the terminal skin lost %q", want)
		}
	}
}

func TestTheLightSkinsKeepTheirDrawnCharacter(t *testing.T) {
	soft := Current("soft", "")
	if !soft.Round || soft.Radius != 16 {
		t.Errorf("the soft skin should be round at 16, got round=%v r=%d", soft.Round, soft.Radius)
	}
	if soft.Level != "dots" {
		t.Errorf("the soft skin shows the level as %q, want bouncing dots", soft.Level)
	}
	if !strings.Contains(soft.CSSVars(), "--lvlw:10px") || !strings.Contains(soft.CSSVars(), "--lvlr:50%") {
		t.Error("the soft skin did not ask for round level marks")
	}
	if p := soft.Palette; p.Text != p.Accent {
		t.Errorf("the soft skin speaks in one voice, got text %q accent %q", p.Text, p.Accent)
	}

	paper := Current("paper", "")
	if !paper.Round || paper.Radius != 10 {
		t.Errorf("the paper skin should be rounded at 10, got round=%v r=%d", paper.Round, paper.Radius)
	}
	if !strings.Contains(paper.CSSVars(), "--barr:2px") {
		t.Error("the paper skin should keep its level marks nearly square")
	}
	if !strings.Contains(paper.CSSVars(), "--scheme:light") || !strings.Contains(soft.CSSVars(), "--scheme:light") {
		t.Error("a light skin must ask the system controls for a light scheme")
	}
	if !strings.Contains(Current("terminal", "green").CSSVars(), "--scheme:dark") {
		t.Error("a dark skin must keep the dark scheme")
	}
	if paper.Level != "bars" {
		t.Errorf("the paper skin shows the level as %q, want plain bars", paper.Level)
	}
	if p := paper.Palette; p.Ok == p.Accent {
		t.Error("the paper skin should light its lamps green, not in the link blue")
	}
	if !strings.Contains(paper.CSSVars(), "--ok:#1a7f37") {
		t.Error("the paper skin lost its green lamp")
	}
	for _, id := range []string{"terminal", "editor", "neon", "soft"} {
		p := GetPalette(GetSkin(id).Palette)
		if p.Ok != p.Accent {
			t.Errorf("the %s palette should light its lamps in its own accent, got %q", id, p.Ok)
		}
	}
}
