package trmodel

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"holdtotype/internal/checksum"
)

const (
	File     = "gemma-3-4b-it-Q4_K_M.gguf"
	Size     = int64(2489757856)
	SHA256   = "be49949e48422e4547b00af14179a193d3777eea7fbbd7d6e1b0861304628a01"
	URL      = "https://huggingface.co/lmstudio-community/gemma-3-4b-it-GGUF/resolve/d650fa07be1a9252c9f7c6597fadc729a377254b/gemma-3-4b-it-Q4_K_M.gguf"
	TermsURL = "https://ai.google.dev/gemma/terms"
	Mark     = ".checked"
)

func SizeMB() int { return int(Size / (1024 * 1024)) }

type State string

const (
	Absent    State = "absent"
	WrongSize State = "size"
	WrongHash State = "hash"
	Ready     State = "ready"
)

type Spec struct {
	Name string
	Size int64
	Hash string
}

func Fixed() Spec { return Spec{Name: File, Size: Size, Hash: SHA256} }

func Is(name string) bool { return strings.EqualFold(filepath.Base(name), File) }

func stamp(spec Spec, fi os.FileInfo) string {
	return spec.Hash + " " + strconv.FormatInt(fi.Size(), 10) + " " + strconv.FormatInt(fi.ModTime().UnixNano(), 10)
}

func Known(dir string, spec Spec) (State, bool) {
	path := filepath.Join(dir, spec.Name)
	fi, err := os.Stat(path)
	if err != nil || !fi.Mode().IsRegular() {
		return Absent, true
	}
	if fi.Size() != spec.Size {
		return WrongSize, true
	}
	raw, err := os.ReadFile(path + Mark)
	if err == nil && strings.TrimSpace(string(raw)) == stamp(spec, fi) {
		return Ready, true
	}
	return "", false
}

func Check(dir string, spec Spec) (State, error) {
	if st, ok := Known(dir, spec); ok {
		return st, nil
	}
	path := filepath.Join(dir, spec.Name)
	fi, err := os.Stat(path)
	if err != nil {
		return Absent, nil
	}
	got, err := checksum.File(path)
	if err != nil {
		return "", err
	}
	if got != checksum.Normalize(spec.Hash) {
		_ = os.Remove(path + Mark)
		return WrongHash, nil
	}
	if after, err := os.Stat(path); err != nil || after.Size() != fi.Size() || !after.ModTime().Equal(fi.ModTime()) {
		return "", fmt.Errorf("%s changed while it was being checked", spec.Name)
	}
	_ = os.WriteFile(path+Mark, []byte(stamp(spec, fi)+"\r\n"), 0o644)
	return Ready, nil
}

func Remove(dir string, spec Spec) error {
	path := filepath.Join(dir, spec.Name)
	var err error
	for i := 0; i < 10; i++ {
		err = os.Remove(path)
		if err == nil || os.IsNotExist(err) {
			err = nil
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	_ = os.Remove(path + Mark)
	_ = os.Remove(path + ".part")
	return err
}
