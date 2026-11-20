package partfile

import (
	"os"
	"strings"
)

const tagSuffix = ".src"

func Tag(part string) string {
	return part + tagSuffix
}

func IsTag(name string) bool {
	return strings.HasSuffix(name, ".part"+tagSuffix)
}

func Kept(part, url string) int64 {
	fi, err := os.Stat(part)
	if err != nil {
		return 0
	}
	src, err := os.ReadFile(Tag(part))
	if err != nil || string(src) != url {
		return 0
	}
	return fi.Size()
}

func Resumable(part, url string) (int64, bool) {
	if n := Kept(part, url); n > 0 {
		return n, false
	}
	_, err := os.Stat(part)
	Remove(part)
	return 0, err == nil
}

func Begin(part, url string) error {
	return os.WriteFile(Tag(part), []byte(url), 0o644)
}

func Finish(part string) {
	_ = os.Remove(Tag(part))
}

func Remove(part string) {
	_ = os.Remove(part)
	_ = os.Remove(Tag(part))
}
