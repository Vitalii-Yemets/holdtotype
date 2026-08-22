package logfilter

import (
	"bytes"
	"io"
	"strings"
	"sync"
)

type Writer struct {
	dst  io.Writer
	skip []string

	mu  sync.Mutex
	buf []byte
}

func New(dst io.Writer, skip ...string) *Writer {
	return &Writer{dst: dst, skip: skip}
}

func (w *Writer) drop(line string) bool {
	for _, s := range w.skip {
		if s != "" && strings.Contains(line, s) {
			return true
		}
	}
	return false
}

func (w *Writer) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	n := len(p)
	w.buf = append(w.buf, p...)
	for {
		i := bytes.IndexByte(w.buf, '\n')
		if i < 0 {
			break
		}
		line := w.buf[:i+1]
		w.buf = w.buf[i+1:]
		if w.drop(string(line)) {
			continue
		}
		if _, err := w.dst.Write(line); err != nil {
			return n, err
		}
	}
	if len(w.buf) > 64<<10 {
		tail := w.buf
		w.buf = nil
		if !w.drop(string(tail)) {
			if _, err := w.dst.Write(tail); err != nil {
				return n, err
			}
		}
	}
	return n, nil
}

func (w *Writer) Flush() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.buf) == 0 {
		return nil
	}
	tail := w.buf
	w.buf = nil
	if w.drop(string(tail)) {
		return nil
	}
	_, err := w.dst.Write(tail)
	return err
}
