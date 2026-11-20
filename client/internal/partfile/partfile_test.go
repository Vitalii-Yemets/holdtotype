package partfile

import (
	"os"
	"path/filepath"
	"testing"
)

const (
	oldURL = "https://huggingface.co/csukuangfj/sherpa-onnx-nemo-parakeet-tdt-0.6b-v3-int8/resolve/main/encoder.int8.onnx"
	newURL = "https://huggingface.co/csukuangfj/sherpa-onnx-nemo-parakeet-tdt-0.6b-v3-int8/resolve/2bda32ec70b097a55adaa07d9a7173915b43cc78/encoder.int8.onnx"
)

func exists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func TestAPartWithoutItsAddressStartsOver(t *testing.T) {
	part := filepath.Join(t.TempDir(), "encoder.int8.onnx.part")
	if err := os.WriteFile(part, []byte("bytes of an older version"), 0o644); err != nil {
		t.Fatal(err)
	}
	have, dropped := Resumable(part, newURL)
	if have != 0 || !dropped || exists(part) {
		t.Fatalf("a part left by an older program is resumed: have %d, dropped %v, still there %v", have, dropped, exists(part))
	}
}

func TestAPartFromAnotherAddressStartsOver(t *testing.T) {
	part := filepath.Join(t.TempDir(), "encoder.int8.onnx.part")
	if err := os.WriteFile(part, []byte("bytes of the main branch"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Begin(part, oldURL); err != nil {
		t.Fatal(err)
	}
	if Kept(part, newURL) != 0 {
		t.Fatal("a part of another address is counted as downloaded")
	}
	have, dropped := Resumable(part, newURL)
	if have != 0 || !dropped || exists(part) || exists(Tag(part)) {
		t.Fatalf("a part of another address is glued to the pinned one: have %d, dropped %v", have, dropped)
	}
}

func TestAPartOfTheSameAddressIsResumed(t *testing.T) {
	part := filepath.Join(t.TempDir(), "encoder.int8.onnx.part")
	if err := os.WriteFile(part, []byte("0123456789"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Begin(part, newURL); err != nil {
		t.Fatal(err)
	}
	if Kept(part, newURL) != 10 {
		t.Fatalf("kept %d bytes", Kept(part, newURL))
	}
	have, dropped := Resumable(part, newURL)
	if have != 10 || dropped || !exists(part) {
		t.Fatalf("the same address is not resumed: have %d, dropped %v", have, dropped)
	}
	Finish(part)
	if exists(Tag(part)) {
		t.Fatal("the record of the address stays after the download")
	}
}

func TestNoPartLeavesNoRecord(t *testing.T) {
	part := filepath.Join(t.TempDir(), "tokens.txt.part")
	if err := os.WriteFile(Tag(part), []byte(newURL), 0o644); err != nil {
		t.Fatal(err)
	}
	have, dropped := Resumable(part, newURL)
	if have != 0 || dropped || exists(Tag(part)) {
		t.Fatalf("an orphan record: have %d, dropped %v, record there %v", have, dropped, exists(Tag(part)))
	}
	if !IsTag("tokens.txt.part.src") || IsTag("tokens.txt.part") || IsTag("tokens.txt") {
		t.Fatal("the record is not told from the part")
	}
}
