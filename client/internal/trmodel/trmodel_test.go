package trmodel

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func small(t *testing.T, dir string, body string) Spec {
	t.Helper()
	sum := sha256.Sum256([]byte(body))
	spec := Spec{Name: "tr.gguf", Size: int64(len(body)), Hash: hex.EncodeToString(sum[:])}
	if err := os.WriteFile(filepath.Join(dir, spec.Name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return spec
}

func TestTheFileCheck(t *testing.T) {
	dir := t.TempDir()
	spec := small(t, dir, "the translation model")
	if st, ok := Known(dir, spec); ok {
		t.Fatalf("known before the check: %s", st)
	}
	if st, err := Check(dir, spec); err != nil || st != Ready {
		t.Fatalf("right file: %s %v", st, err)
	}
	if st, ok := Known(dir, spec); !ok || st != Ready {
		t.Fatalf("not remembered: %s %v", st, ok)
	}

	wrong := spec
	wrong.Hash = "00" + spec.Hash[2:]
	if st, ok := Known(dir, wrong); ok {
		t.Fatalf("a mark for another hash is trusted: %s", st)
	}
	if st, err := Check(dir, wrong); err != nil || st != WrongHash {
		t.Fatalf("wrong hash: %s %v", st, err)
	}
	if _, err := os.Stat(filepath.Join(dir, spec.Name+Mark)); !os.IsNotExist(err) {
		t.Fatal("the mark stays after a wrong hash")
	}

	size := spec
	size.Size++
	if st, err := Check(dir, size); err != nil || st != WrongSize {
		t.Fatalf("wrong size: %s %v", st, err)
	}

	gone := spec
	gone.Name = "absent.gguf"
	if st, err := Check(dir, gone); err != nil || st != Absent {
		t.Fatalf("absent: %s %v", st, err)
	}
}

func TestAChangedFileIsCheckedAgain(t *testing.T) {
	dir := t.TempDir()
	spec := small(t, dir, "abcdef")
	if st, _ := Check(dir, spec); st != Ready {
		t.Fatalf("first check %s", st)
	}
	path := filepath.Join(dir, spec.Name)
	if err := os.WriteFile(path, []byte("abcdeg"), 0o644); err != nil {
		t.Fatal(err)
	}
	later := time.Now().Add(time.Minute)
	_ = os.Chtimes(path, later, later)
	if st, ok := Known(dir, spec); ok {
		t.Fatalf("a changed file is still known: %s", st)
	}
	if st, _ := Check(dir, spec); st != WrongHash {
		t.Fatalf("changed file %s", st)
	}
}

func TestRemoveTakesTheMarkAndThePart(t *testing.T) {
	dir := t.TempDir()
	spec := small(t, dir, "x")
	_, _ = Check(dir, spec)
	_ = os.WriteFile(filepath.Join(dir, spec.Name+".part"), []byte("y"), 0o644)
	if err := Remove(dir, spec); err != nil {
		t.Fatal(err)
	}
	for _, n := range []string{spec.Name, spec.Name + Mark, spec.Name + ".part"} {
		if _, err := os.Stat(filepath.Join(dir, n)); !os.IsNotExist(err) {
			t.Errorf("%s stays", n)
		}
	}
	if err := Remove(dir, spec); err != nil {
		t.Fatalf("removing nothing: %v", err)
	}
}

func TestTheFixedModel(t *testing.T) {
	f := Fixed()
	if f.Name != "gemma-3-4b-it-Q4_K_M.gguf" || f.Size != 2489757856 || f.Hash != "be49949e48422e4547b00af14179a193d3777eea7fbbd7d6e1b0861304628a01" {
		t.Fatalf("%+v", f)
	}
	if URL != "https://huggingface.co/lmstudio-community/gemma-3-4b-it-GGUF/resolve/d650fa07be1a9252c9f7c6597fadc729a377254b/gemma-3-4b-it-Q4_K_M.gguf" {
		t.Fatal(URL)
	}
	if SizeMB() != 2374 {
		t.Fatalf("size %d MB", SizeMB())
	}
	if !Is("models/GEMMA-3-4b-it-Q4_K_M.gguf") || Is("qwen.gguf") {
		t.Fatal("Is")
	}
}
