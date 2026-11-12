package editormodel

import (
	"os"
	"path/filepath"
	"strings"

	"holdtotype/internal/trmodel"
)

type File struct {
	Name   string
	SizeMB int
}

func List(dir string) []File {
	entries, _ := os.ReadDir(dir)
	var out []File
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".gguf") || strings.HasPrefix(name, "mmproj-") || trmodel.Is(name) {
			continue
		}
		info, err := e.Info()
		if err != nil || !info.Mode().IsRegular() {
			continue
		}
		out = append(out, File{Name: name, SizeMB: int(info.Size() / (1024 * 1024))})
	}
	return out
}

func Chosen(path string) (File, bool) {
	p := strings.TrimSpace(path)
	if p == "" {
		return File{}, false
	}
	info, err := os.Stat(p)
	if err != nil || !info.Mode().IsRegular() {
		return File{}, false
	}
	return File{Name: filepath.Base(p), SizeMB: int(info.Size() / (1024 * 1024))}, true
}
