package editormodel

import (
	"os"
	"path/filepath"
	"testing"
)

func put(t *testing.T, path string, size int) {
	t.Helper()
	if err := os.WriteFile(path, make([]byte, size), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestListShowsOnlyEditorModels(t *testing.T) {
	dir := t.TempDir()
	put(t, filepath.Join(dir, "a.gguf"), 3<<20)
	put(t, filepath.Join(dir, "mmproj-a.gguf"), 10)
	put(t, filepath.Join(dir, "gemma-3-4b-it-Q4_K_M.gguf"), 10)
	put(t, filepath.Join(dir, "notes.txt"), 10)
	if err := os.Mkdir(filepath.Join(dir, "b.gguf"), 0o755); err != nil {
		t.Fatal(err)
	}
	got := List(dir)
	if len(got) != 1 || got[0].Name != "a.gguf" || got[0].SizeMB != 3 {
		t.Fatalf("listed %+v", got)
	}
}

func TestListOfAMissingFolderIsEmpty(t *testing.T) {
	if got := List(filepath.Join(t.TempDir(), "nope")); len(got) != 0 {
		t.Fatalf("listed %+v", got)
	}
}

func TestChosenIsNamedOnlyWhileItsFileIsThere(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "Qwen2.5-3B-Instruct-IQ2_M.gguf")
	put(t, p, 2<<20)
	f, ok := Chosen(p)
	if !ok || f.Name != "Qwen2.5-3B-Instruct-IQ2_M.gguf" || f.SizeMB != 2 {
		t.Fatalf("chosen %+v %v", f, ok)
	}
	if err := os.Remove(p); err != nil {
		t.Fatal(err)
	}
	if f, ok := Chosen(p); ok || f.Name != "" {
		t.Fatalf("a deleted model is still named: %+v", f)
	}
}

func TestNoChoiceAndAFolderNameNothing(t *testing.T) {
	for _, p := range []string{"", "   ", t.TempDir()} {
		if f, ok := Chosen(p); ok || f.Name != "" || f.SizeMB != 0 {
			t.Fatalf("%q named %+v", p, f)
		}
	}
}

func TestDeletingAnotherModelLeavesTheChosenOne(t *testing.T) {
	dir := t.TempDir()
	keep := filepath.Join(dir, "keep.gguf")
	gone := filepath.Join(dir, "gone.gguf")
	put(t, keep, 1<<20)
	put(t, gone, 1<<20)
	if err := os.Remove(gone); err != nil {
		t.Fatal(err)
	}
	if f, ok := Chosen(keep); !ok || f.Name != "keep.gguf" {
		t.Fatalf("the chosen model is lost: %+v", f)
	}
	if got := List(dir); len(got) != 1 || got[0].Name != "keep.gguf" {
		t.Fatalf("listed %+v", got)
	}
}
