package checksum

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

func File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func Normalize(want string) string {
	want = strings.TrimSpace(strings.ToLower(want))
	return strings.TrimPrefix(want, "sha256:")
}

func Verify(path, want string) error {
	want = Normalize(want)
	if want == "" {
		return nil
	}
	return compare(path, want)
}

func Require(path, want string) error {
	want = Normalize(want)
	if want == "" {
		return fmt.Errorf("%s: no reference hash to check against", path)
	}
	return compare(path, want)
}

func compare(path, want string) error {
	got, err := File(path)
	if err != nil {
		return err
	}
	if got != want {
		return fmt.Errorf("%s: %s ≠ %s", path, head(got, 12), head(want, 12))
	}
	return nil
}

func head(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
