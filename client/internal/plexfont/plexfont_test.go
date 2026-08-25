package plexfont

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestBothFacesAreEmbedded(t *testing.T) {
	for name, data := range map[string][]byte{"Regular": Regular, "SemiBold": SemiBold} {
		if len(data) < 50_000 {
			t.Fatalf("%s is %d bytes — the font did not make it into the binary", name, len(data))
		}
		if got := data[:4]; got[0] != 0x00 || got[1] != 0x01 || got[2] != 0x00 || got[3] != 0x00 {
			t.Errorf("%s does not start with the TrueType mark: % x", name, got)
		}
	}
}

func TestFaceCSSCarriesTheFontsThemselves(t *testing.T) {
	css := FaceCSS()
	if n := strings.Count(css, "@font-face"); n != 2 {
		t.Fatalf("FaceCSS() holds %d faces, want two", n)
	}
	for _, weight := range []string{"font-weight:400", "font-weight:600"} {
		if !strings.Contains(css, weight) {
			t.Errorf("FaceCSS() misses %q", weight)
		}
	}
	if !strings.Contains(css, "font-family:'"+Family+"'") {
		t.Errorf("FaceCSS() names something other than %q", Family)
	}
	for _, part := range strings.Split(css, "base64,")[1:] {
		payload := part[:strings.Index(part, ")")]
		raw, err := base64.StdEncoding.DecodeString(payload)
		if err != nil {
			t.Fatalf("a face is not valid base64: %v", err)
		}
		if len(raw) < 50_000 || raw[1] != 0x01 {
			t.Errorf("a face decodes to %d bytes that do not look like a font", len(raw))
		}
	}
}

func TestFaceCSSIsBuiltOnce(t *testing.T) {
	if FaceCSS() != FaceCSS() {
		t.Error("FaceCSS() built two different strings")
	}
}
