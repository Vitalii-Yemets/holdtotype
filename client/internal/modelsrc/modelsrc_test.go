package modelsrc

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var client = filepath.Join("..", "..")

func TestEveryAddressIsPinnedToARevision(t *testing.T) {
	pinned := regexp.MustCompile(`^https://huggingface\.co/[\w.-]+/[\w.-]+/resolve/[0-9a-f]{40}/$`)
	want := map[string]string{
		"Whisper":  "5359861c739e955e79d9a303bcbc70fb988958b1",
		"GigaAMv3": "a6039be7cee829a9044a69ac0ebaf1c191217c97",
		"GigaAMv2": "d7002de64f758983ac1146695e5de5a381c36785",
		"Parakeet": "2bda32ec70b097a55adaa07d9a7173915b43cc78",
		"Nemotron": "3f1f6020667b151d0b6820cc4b8ae4f42b6d35fc",
		"Canary":   "9077164e0d3dd1d5353743e89ceaa1d3a770838c",
		"Qwen3ASR": "68818b2313fe77bd06f6a7c5068ff3ef59d02b8a",
	}
	all := All()
	if len(all) != len(want) {
		t.Fatalf("%d addresses, want %d", len(all), len(want))
	}
	for name, u := range all {
		if !pinned.MatchString(u) {
			t.Fatalf("%s is not pinned to a revision: %s", name, u)
		}
		if !strings.Contains(u, "/resolve/"+want[name]+"/") {
			t.Fatalf("%s is pinned to another revision: %s", name, u)
		}
	}
}

func TestNoDownloadInTheProgramFollowsABranch(t *testing.T) {
	lit := regexp.MustCompile(`"[^"\n]*/resolve/([^/"\n]+)/[^"\n]*"`)
	hex := regexp.MustCompile(`^[0-9a-f]{40}$`)
	seen := 0
	err := filepath.WalkDir(client, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") {
			return nil
		}
		raw, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		for _, m := range lit.FindAllStringSubmatch(string(raw), -1) {
			seen++
			if hex.MatchString(m[1]) {
				continue
			}
			if filepath.Base(p) == "llm.go" && m[0] == `"/resolve/main/"` {
				continue
			}
			t.Errorf("%s downloads from a branch instead of a revision: %s", p, m[0])
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if seen == 0 {
		t.Fatal("no download address was found, the sources were not read")
	}
}

func catalogBases(t *testing.T, file string) map[string]string {
	raw, err := os.ReadFile(filepath.Join(client, file))
	if err != nil {
		t.Fatal(err)
	}
	src := string(raw)
	if regexp.MustCompile(`BaseURL:\s*"`).MatchString(src) {
		t.Fatalf("%s writes a download address by hand instead of taking it from modelsrc", file)
	}
	dir := regexp.MustCompile(`Dir:\s*"([^"]+)"`)
	base := regexp.MustCompile(`BaseURL:\s*modelsrc\.(\w+)`)
	out := map[string]string{}
	locs := dir.FindAllStringSubmatchIndex(src, -1)
	for i, l := range locs {
		end := len(src)
		if i+1 < len(locs) {
			end = locs[i+1][0]
		}
		if m := base.FindStringSubmatch(src[l[1]:end]); m != nil {
			out[src[l[2]:l[3]]] = m[1]
		}
	}
	return out
}

func TestTheProgramAndTheInstallerDownloadAModelFromOneAddress(t *testing.T) {
	app := catalogBases(t, "models.go")
	setup := catalogBases(t, filepath.Join("setup", "install.go"))
	if len(setup) == 0 || len(app) == 0 {
		t.Fatalf("catalogs not found: program %v, installer %v", app, setup)
	}
	for dir, name := range setup {
		if _, ok := All()[name]; !ok {
			t.Fatalf("the installer takes %s from an unknown address %s", dir, name)
		}
		if app[dir] != name {
			t.Fatalf("%s: the installer downloads from %s, the program from %s", dir, name, app[dir])
		}
	}
	for _, f := range []struct{ file, line string }{{"models.go", "const modelBaseURL = modelsrc.Whisper"}, {filepath.Join("setup", "install.go"), "const whisperBaseURL = modelsrc.Whisper"}} {
		raw, err := os.ReadFile(filepath.Join(client, f.file))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(raw), f.line) {
			t.Fatalf("%s does not take the Whisper models from modelsrc", f.file)
		}
	}
}
