package fastpack

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

func keysOf(packs []Pack) []string {
	var out []string
	for _, p := range packs {
		out = append(out, p.Key())
	}
	return out
}

func TestPairsAreCheckedAndKeptInOrder(t *testing.T) {
	if from, to, ok := ParsePair("fr-ru"); !ok || from != "fr" || to != "ru" {
		t.Fatalf("fr-ru: %q %q %v", from, to, ok)
	}
	for _, bad := range []string{"", "fr", "fr-fr", "xx-ru", "fr-", "-ru", "fr-ru-en", "FR-RU"} {
		if _, _, ok := ParsePair(bad); ok {
			t.Fatalf("%q is taken as a pair", bad)
		}
	}
	got := CleanPairs([]string{"en-ru", "bad", "fr-ru", "en-ru", "de-en"})
	if strings.Join(got, ",") != "en-ru,fr-ru,de-en" {
		t.Fatalf("cleaned: %v", got)
	}
	if got := AddPair([]string{"en-ru"}, "fr-ru"); strings.Join(got, ",") != "en-ru,fr-ru" {
		t.Fatalf("added: %v", got)
	}
	if got := AddPair([]string{"en-ru"}, "en-ru"); strings.Join(got, ",") != "en-ru" {
		t.Fatalf("added twice: %v", got)
	}
	if got := AddPair(nil, "ru-ru"); len(got) != 0 {
		t.Fatalf("a pair into its own language: %v", got)
	}
}

func TestRemovingAPairKeepsThePacksAnotherPairNeeds(t *testing.T) {
	pairs := []string{"en-ru", "fr-ru", "fr-de"}
	left, drop := DropPair(pairs, "fr-ru")
	if strings.Join(left, ",") != "en-ru,fr-de" {
		t.Fatalf("left: %v", left)
	}
	if len(drop) != 0 {
		t.Fatalf("fr-en is needed by fr-de and en-ru by en-ru, yet %v is deleted", keysOf(drop))
	}
	left, drop = DropPair(left, "en-ru")
	if strings.Join(left, ",") != "fr-de" || strings.Join(keysOf(drop), ",") != "en-ru" {
		t.Fatalf("left %v, deleted %v", left, keysOf(drop))
	}
	left, drop = DropPair(left, "fr-de")
	if len(left) != 0 || strings.Join(keysOf(drop), ",") != "fr-en,en-de" {
		t.Fatalf("left %v, deleted %v", left, keysOf(drop))
	}
	left, drop = DropPair([]string{"de-ru", "es-ru"}, "de-ru")
	if strings.Join(left, ",") != "es-ru" || strings.Join(keysOf(drop), ",") != "de-en" {
		t.Fatalf("the shared English half: left %v, deleted %v", left, keysOf(drop))
	}
	left, drop = DropPair([]string{"en-ru"}, "de-ru")
	if strings.Join(left, ",") != "en-ru" || len(drop) != 0 {
		t.Fatalf("a pair that was not added: left %v, deleted %v", left, keysOf(drop))
	}
	need := Needed([]string{"de-ru", "en-uk"})
	if len(need) != 3 || !need["de-en"] || !need["en-ru"] || !need["en-uk"] {
		t.Fatalf("needed: %v", need)
	}
}

func TestTheFirstListIsOnePairPerPackOnDisk(t *testing.T) {
	on := map[string]bool{"en-ru": true, "de-en": true, "uk-en": true}
	got := PairsOnDisk(func(p Pack) bool { return on[p.Key()] })
	if len(got) != 3 {
		t.Fatalf("pairs: %v", got)
	}
	for _, k := range got {
		if !on[k] {
			t.Fatalf("pair %s has no pack", k)
		}
		if _, _, ok := ParsePair(k); !ok {
			t.Fatalf("pair %s is not a pair", k)
		}
	}
	if got := PairsOnDisk(func(Pack) bool { return false }); len(got) != 0 {
		t.Fatalf("no packs: %v", got)
	}
}
