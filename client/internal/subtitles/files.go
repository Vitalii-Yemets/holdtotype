package subtitles

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Files struct {
	Dir  string
	Logf func(format string, args ...any)

	started  time.Time
	mu       sync.Mutex
	orig     *os.File
	origPath string
	origErr  bool
	tr       *os.File
	trTo     string
	trPaths  []string
	trErr    bool
	done     bool
}

func NewFiles(dir string, started time.Time, logf func(format string, args ...any)) *Files {
	return &Files{Dir: dir, Logf: logf, started: started}
}

func (f *Files) logf(format string, args ...any) {
	if f.Logf != nil {
		f.Logf(format, args...)
	}
}

func (f *Files) open(name, header string) (*os.File, string, error) {
	if err := os.MkdirAll(f.Dir, 0o755); err != nil {
		return nil, "", fmt.Errorf("the transcript folder was not created: %w", err)
	}
	p := filepath.Join(f.Dir, name)
	fh, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, "", fmt.Errorf("the transcript was not created: %w", err)
	}
	if _, err := fh.WriteString(Opening(header)); err != nil {
		f.logf("subtitles: writing the transcript: %v", err)
	}
	if abs, aerr := filepath.Abs(p); aerr == nil {
		p = abs
	}
	return fh, p, nil
}

func (f *Files) write(fh *os.File, entry string) {
	if _, err := fh.WriteString(entry); err != nil {
		f.logf("subtitles: writing the transcript: %v", err)
		return
	}
	_ = fh.Sync()
}

func (f *Files) Original(save bool, lang, model string, translating bool, entry string) {
	if !save {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.origErr || f.done {
		return
	}
	if f.orig == nil {
		fh, p, err := f.open(OriginalFileName(f.started, lang, translating), OriginalHeader(f.started, lang, model, translating))
		if err != nil {
			f.logf("subtitles: %v", err)
			f.origErr = true
			return
		}
		f.orig, f.origPath = fh, p
		f.logf("subtitles: transcript %s", p)
	}
	f.write(f.orig, entry)
}

func (f *Files) Translation(save bool, from, to, by, entry string) {
	if !save {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.trErr || f.done {
		return
	}
	if f.tr != nil && f.trTo != to {
		_ = f.tr.Close()
		f.tr = nil
	}
	if f.tr == nil {
		fh, p, err := f.open(TranslationFileName(f.started, to), TranslationHeader(f.started, from, to, by))
		if err != nil {
			f.logf("subtitles: %v", err)
			f.trErr = true
			return
		}
		f.tr, f.trTo = fh, to
		f.trPaths = append(f.trPaths, p)
		f.logf("subtitles: translation transcript %s", p)
	}
	f.write(f.tr, entry)
}

func (f *Files) Close() (string, string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.orig != nil {
		_ = f.orig.Close()
	}
	if f.tr != nil {
		_ = f.tr.Close()
	}
	f.done = true
	return f.origPath, strings.Join(f.trPaths, ", ")
}
