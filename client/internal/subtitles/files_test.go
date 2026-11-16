package subtitles

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func readOrEmpty(t *testing.T, p string) string {
	t.Helper()
	b, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return ""
	}
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestFilesFollowTheSaveSwitchPhraseByPhrase(t *testing.T) {
	dir := filepath.Join(t.TempDir(), Folder)
	start := time.Date(2026, 11, 16, 20, 1, 2, 0, time.Local)
	at := start.Add(5 * time.Second)
	f := NewFiles(dir, start, nil)
	orig := filepath.Join(dir, OriginalFileName(start, "en", true))
	trans := filepath.Join(dir, TranslationFileName(start, "ru"))

	f.Original(false, "en", "Nemotron 3.5", true, Entry(at, "one"))
	f.Translation(false, "en", "ru", "Fast", Entry(at, "один"))
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatalf("with saving off nothing is created, the folder: %v", err)
	}

	f.Original(true, "en", "Nemotron 3.5", true, Entry(at, "two"))
	f.Translation(true, "en", "ru", "Fast", Entry(at, "два"))
	o, tr := readOrEmpty(t, orig), readOrEmpty(t, trans)
	if !strings.HasPrefix(o, Opening(OriginalHeader(start, "en", "Nemotron 3.5", true))) || !strings.HasPrefix(tr, Opening(TranslationHeader(start, "en", "ru", "Fast"))) {
		t.Fatalf("turned on in the middle, each file starts with its header:\n%q\n%q", o, tr)
	}
	if strings.Contains(o, "one") || !strings.Contains(o, "two") || strings.Contains(tr, "один") || !strings.Contains(tr, "два") {
		t.Fatalf("the files begin with the first phrase after saving was turned on:\n%q\n%q", o, tr)
	}

	f.Original(false, "en", "Nemotron 3.5", true, Entry(at, "three"))
	f.Translation(false, "en", "ru", "Fast", Entry(at, "три"))
	f.Original(true, "en", "Nemotron 3.5", true, Entry(at, "four"))
	f.Translation(true, "en", "ru", "Fast", Entry(at, "четыре"))
	o, tr = readOrEmpty(t, orig), readOrEmpty(t, trans)
	if strings.Contains(o, "three") || strings.Contains(tr, "три") {
		t.Fatalf("a phrase heard while saving was off is not written:\n%q\n%q", o, tr)
	}
	if !strings.HasSuffix(o, Entry(at, "two")+Entry(at, "four")) || !strings.HasSuffix(tr, Entry(at, "два")+Entry(at, "четыре")) {
		t.Fatalf("turned on again, the same files go on:\n%q\n%q", o, tr)
	}
	if strings.Count(o, "subtitles ·") != 1 || strings.Count(tr, "subtitles ·") != 1 {
		t.Fatalf("one header per file:\n%q\n%q", o, tr)
	}

	op, tp := f.Close()
	if filepath.Base(op) != filepath.Base(orig) || filepath.Base(tp) != filepath.Base(trans) {
		t.Fatalf("closing names the files: %q %q", op, tp)
	}
	f.Original(true, "en", "Nemotron 3.5", true, Entry(at, "five"))
	f.Translation(true, "en", "ru", "Fast", Entry(at, "пять"))
	if strings.Contains(readOrEmpty(t, orig), "five") || strings.Contains(readOrEmpty(t, trans), "пять") {
		t.Fatal("nothing is written after the run is closed")
	}
}

func TestFilesNeverOnNameNothing(t *testing.T) {
	f := NewFiles(filepath.Join(t.TempDir(), Folder), time.Now(), nil)
	f.Original(false, "en", "m", false, "x\r\n")
	f.Translation(false, "en", "ru", "Fast", "x\r\n")
	if op, tp := f.Close(); op != "" || tp != "" {
		t.Fatalf("with saving off the whole run there are no files: %q %q", op, tp)
	}
}
