package punct

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestModelFilesArePinned(t *testing.T) {
	if !regexp.MustCompile(`^[0-9a-f]{40}$`).MatchString(Revision) {
		t.Fatalf("the revision %q is not a full commit hash", Revision)
	}
	base := BaseURL()
	if base != "https://huggingface.co/Horizon-Labs/punctuation-restoration-small/resolve/a1818fae73a1d09b5496e67747dda7bd667b79a8/" {
		t.Fatalf("base address %q", base)
	}
	if strings.Contains(base, "/main/") {
		t.Fatalf("the address follows a branch: %q", base)
	}
	want := map[string]File{
		"model_quantized.onnx": {Remote: "onnx/model_quantized.onnx", Bytes: 268414381, SHA256: "5679870be92c7fa47c9be9b2872fac389ef28d43d2c33f35e1315b93a82f4ace"},
		"tokenizer.json":       {Remote: "tokenizer.json", Bytes: 34363286, SHA256: "f34b44a6dc8a0eef651bfa10e8325057530152606a69d9c3b6ae02efb1558f0a"},
	}
	if len(Files) != len(want) {
		t.Fatalf("%d files, want %d", len(Files), len(want))
	}
	hex := regexp.MustCompile(`^[0-9a-f]{64}$`)
	for _, f := range Files {
		w, ok := want[f.Name]
		if !ok {
			t.Fatalf("unexpected file %q", f.Name)
		}
		if f.Remote != w.Remote || f.Bytes != w.Bytes || f.SHA256 != w.SHA256 {
			t.Fatalf("%s: got %+v", f.Name, f)
		}
		if !hex.MatchString(f.SHA256) {
			t.Fatalf("%s: the sum %q is not sha256", f.Name, f.SHA256)
		}
		if strings.ContainsAny(f.Name, `/\`) {
			t.Fatalf("%s: the local name must sit right in the model folder", f.Name)
		}
	}
	if SizeMB() != 289 {
		t.Fatalf("size %d MB, want 289", SizeMB())
	}
	if ModelDir != "punct/horizon" || ModelName != "Horizon Punctuation Small" || LangCount != 89 || RAMMB != 340 {
		t.Fatalf("catalog record changed: %s %s %d %d", ModelDir, ModelName, LangCount, RAMMB)
	}
}

func TestMissingNamesTheFirstAbsentFile(t *testing.T) {
	dir := t.TempDir()
	if why := Missing(dir); !strings.Contains(why, "model_quantized.onnx") {
		t.Fatalf("empty folder: %q", why)
	}
	if err := os.WriteFile(filepath.Join(dir, "model_quantized.onnx"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if why := Missing(dir); !strings.Contains(why, "tokenizer.json") {
		t.Fatalf("without the tokenizer: %q", why)
	}
	if err := os.WriteFile(filepath.Join(dir, "tokenizer.json.part"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if why := Missing(dir); !strings.Contains(why, "tokenizer.json") {
		t.Fatalf("an unfinished download counts as missing: %q", why)
	}
	if err := os.WriteFile(filepath.Join(dir, "tokenizer.json"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if why := Missing(dir); why != "" {
		t.Fatalf("both files: %q", why)
	}
}

func TestOffReason(t *testing.T) {
	cases := []struct {
		on, translating bool
		missing, want   string
	}{
		{false, true, "", "whole sentences are turned off"},
		{false, false, "gone", "whole sentences are turned off"},
		{true, false, "", "translation is off"},
		{true, true, "the punctuation model is missing: x", "the punctuation model is missing: x"},
		{true, true, "", ""},
	}
	for _, c := range cases {
		if got := OffReason(c.on, c.translating, c.missing); got != c.want {
			t.Fatalf("OffReason(%v, %v, %q) = %q, want %q", c.on, c.translating, c.missing, got, c.want)
		}
	}
}
