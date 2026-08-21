package main

import (
	"os"
	"sync"
)

type rotatingWriter struct {
	mu      sync.Mutex
	path    string
	maxSize int64
	f       *os.File
	size    int64
}

func newRotatingWriter(path string, maxSize int64) *rotatingWriter {
	w := &rotatingWriter{path: path, maxSize: maxSize}
	if st, err := os.Stat(path); err == nil && st.Size() > maxSize {
		w.rotate()
	}
	w.open()
	return w
}

func (w *rotatingWriter) open() {
	f, err := os.OpenFile(w.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	w.f = f
	if st, err := f.Stat(); err == nil {
		w.size = st.Size()
	}
}

func (w *rotatingWriter) rotate() {
	if w.f != nil {
		_ = w.f.Close()
		w.f = nil
	}
	old := w.path + ".old"
	_ = os.Remove(old)
	_ = os.Rename(w.path, old)
	w.size = 0
}

func (w *rotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.f == nil {
		w.open()
		if w.f == nil {
			return len(p), nil
		}
	}
	n, err := w.f.Write(p)
	w.size += int64(n)
	if w.size >= w.maxSize {
		w.rotate()
		w.open()
	}
	return n, err
}

func (w *rotatingWriter) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.f != nil {
		_ = w.f.Close()
		w.f = nil
	}
}
