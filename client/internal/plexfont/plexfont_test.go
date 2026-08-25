package plexfont

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestEveryFaceIsEmbedded(t *testing.T) {
	faces := Faces()
	if len(faces) != 4 {
		t.Fatalf("Faces() carries %d faces, want four", len(faces))
	}
	for _, f := range faces {
		if len(f.Data) < 50_000 {
			t.Fatalf("%s %s is %d bytes — the font did not make it into the binary", f.Family, f.Weight, len(f.Data))
		}
		if g := f.Data[:4]; g[0] != 0x00 || g[1] != 0x01 || g[2] != 0x00 || g[3] != 0x00 {
			t.Errorf("%s %s does not start with the TrueType mark: % x", f.Family, f.Weight, g)
		}
	}
}

func TestBothFamiliesAreThere(t *testing.T) {
	seen := map[string][]string{}
	for _, f := range Faces() {
		seen[f.Family] = append(seen[f.Family], f.Weight)
	}
	for _, family := range []string{Sans, Mono} {
		weights := seen[family]
		if len(weights) != 2 || weights[0] != "400" || weights[1] != "600" {
			t.Errorf("%s carries %v, want the regular and the semibold", family, weights)
		}
	}
}

func TestFaceCSSCarriesTheFontsThemselves(t *testing.T) {
	css := FaceCSS()
	if n := strings.Count(css, "@font-face"); n != 4 {
		t.Fatalf("FaceCSS() holds %d faces, want four", n)
	}
	for _, family := range []string{Sans, Mono} {
		if n := strings.Count(css, "font-family:'"+family+"'"); n != 2 {
			t.Errorf("FaceCSS() names %q %d times, want twice", family, n)
		}
	}
	for _, weight := range []string{"font-weight:400", "font-weight:600"} {
		if n := strings.Count(css, weight); n != 2 {
			t.Errorf("FaceCSS() carries %q %d times, want twice", weight, n)
		}
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
