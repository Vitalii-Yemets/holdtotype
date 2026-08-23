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
	got, err := File(path)
	if err != nil {
		return err
	}
	if got != want {
		return fmt.Errorf("%s: %s ≠ %s", path, got[:12], want[:12])
	}
	return nil
}
