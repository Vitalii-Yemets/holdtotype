package fastpack

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"
)

func TestTheCatalogHasFourteenDirections(t *testing.T) {
	all := All()
	if len(all) != 14 {
		t.Fatalf("%d directions", len(all))
	}
	seen := map[string]bool{}
	hex64 := regexp.MustCompile(`^[0-9a-f]{64}$`)
	for _, p := range all {
		if seen[p.Key()] {
			t.Fatalf("%s twice", p.Key())
		}
		seen[p.Key()] = true
		if (p.From == Pivot) == (p.To == Pivot) {
			t.Fatalf("%s does not go to or from English", p.Key())
		}
		if p.Version == "" || len(p.Files) != 3 {
			t.Fatalf("%s: %+v", p.Key(), p)
		}
		for _, kind := range []string{"model", "vocab", "lex"} {
			f, ok := p.File(kind)
			if !ok || f.Size <= 0 || !hex64.MatchString(f.SHA256) || f.Name == "" {
				t.Fatalf("%s %s: %+v", p.Key(), kind, f)
			}
			if !regexp.MustCompile(`^https://firefox-settings-attachments\.cdn\.mozilla\.net/`).MatchString(f.URL) {
				t.Fatalf("%s %s: %s", p.Key(), kind, f.URL)
			}
		}
		if p.SizeMB() < 15 || p.SizeMB() > 60 {
			t.Fatalf("%s: %d MB", p.Key(), p.SizeMB())
		}
	}
	for _, l := range Languages {
		if l == Pivot {
			continue
		}
		if !seen[Pivot+"-"+l] || !seen[l+"-"+Pivot] {
			t.Fatalf("%s is missing a direction", l)
		}
	}
}

func TestTheEnglishToRussianPackIsTheDesktopRelease(t *testing.T) {
	p, ok := Find("en", "ru")
	if !ok || p.Version != "2.0" {
		t.Fatalf("%+v", p)
	}
	want := map[string]string{
		"model": "0ef9a209c5edc46692750e7505b3695655b1c7c3ec73058b641201ef18c481ce",
		"vocab": "56ee63e14e8cb926c394242adc3ed7cc602644c3d33058cff2ce2959d52a6258",
		"lex":   "3587f93c10c1d457f874c2ba7ff4d5e2686fdaf7f3fd179cfefeacab3c6990b9",
	}
	for kind, h := range want {
		if f, _ := p.File(kind); f.SHA256 != h {
			t.Fatalf("%s: %s", kind, f.SHA256)
		}
	}
}

func TestRoutes(t *testing.T) {
	cases := []struct {
		from, to string
		want     []string
	}{
		{"en", "ru", []string{"en-ru"}},
		{"ru", "en", []string{"ru-en"}},
		{"de", "ru", []string{"de-en", "en-ru"}},
		{"uk", "pl", []string{"uk-en", "en-pl"}},
		{"ru", "ru", nil},
		{"", "ru", nil},
		{"en", "zh", nil},
		{"zh", "ru", nil},
	}
	for _, c := range cases {
		packs, ok := Route(c.from, c.to)
		var got []string
		for _, p := range packs {
			got = append(got, p.Key())
		}
		if ok != (c.want != nil) || len(got) != len(c.want) {
			t.Fatalf("%s→%s: %v %v", c.from, c.to, got, ok)
		}
		for i := range got {
			if got[i] != c.want[i] {
				t.Fatalf("%s→%s: %v", c.from, c.to, got)
			}
		}
	}
}

func put(t *testing.T, dir string, f File, body []byte) File {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, f.Name), body, 0o644); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(body)
	f.Size = int64(len(body))
	f.SHA256 = hex.EncodeToString(sum[:])
	return f
}

func TestTheStateOfAPack(t *testing.T) {
	root := t.TempDir()
	p := Pack{From: "en", To: "xx", Version: "1", Files: []File{{Kind: "model", Name: "m.bin"}, {Kind: "vocab", Name: "v.spm"}, {Kind: "lex", Name: "l.bin"}}}
	if st, ok := Known(root, p); !ok || st != Absent {
		t.Fatalf("empty: %s %v", st, ok)
	}
	dir := p.Dir(root)
	for i := range p.Files {
		p.Files[i] = put(t, dir, p.Files[i], []byte("file "+p.Files[i].Name))
	}
	if st, ok := Known(root, p); ok || st != Checking {
		t.Fatalf("before the check: %s %v", st, ok)
	}
	if st, err := Check(root, p); err != nil || st != Ready {
		t.Fatalf("check: %s %v", st, err)
	}
	if st, ok := Known(root, p); !ok || st != Ready {
		t.Fatalf("after the check: %s %v", st, ok)
	}
	if err := os.WriteFile(filepath.Join(dir, "v.spm"), []byte("file v.spX"), 0o644); err != nil {
		t.Fatal(err)
	}
	later := time.Now().Add(time.Minute)
	_ = os.Chtimes(filepath.Join(dir, "v.spm"), later, later)
	if st, _ := Check(root, p); st != Bad {
		t.Fatalf("a changed file: %s", st)
	}
	if err := os.Remove(filepath.Join(dir, "l.bin")); err != nil {
		t.Fatal(err)
	}
	if st, ok := Known(root, p); !ok || st != Absent {
		t.Fatalf("a missing file: %s %v", st, ok)
	}
	if err := Remove(root, p); err != nil {
		t.Fatal(err)
	}
	for _, f := range p.Files {
		if _, err := os.Stat(filepath.Join(dir, f.Name)); !os.IsNotExist(err) {
			t.Fatalf("%s stays", f.Name)
		}
	}
	if filepath.Base(filepath.Dir(dir)) != Folder || filepath.Base(dir) != "en-xx" {
		t.Fatalf("dir %s", dir)
	}
}
