package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/sys/windows"
)

type modelInfo struct {
	ID      string
	File    string
	SizeMB  int
	NameKey string
	DescKey string
}

var modelCatalog = []modelInfo{
	{"base", "ggml-base.bin", 142, "Base", "S_M_BASE"},
	{"small", "ggml-small.bin", 466, "Small", "S_M_SMALL"},
	{"medium-q5_0", "ggml-medium-q5_0.bin", 539, "Medium (q5)", "S_M_MED"},
	{"large-v3-turbo-q5_0", "ggml-large-v3-turbo-q5_0.bin", 574, "Turbo (q5)", "S_M_TURBO"},
}

const modelBaseURL = "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/"

type dlState struct {
	active bool
	pct    int
	err    string
}

var (
	dlMu sync.Mutex
	dl   = map[string]*dlState{}
)

func findModel(id string) *modelInfo {
	for i := range modelCatalog {
		if modelCatalog[i].ID == id {
			return &modelCatalog[i]
		}
	}
	return nil
}

type modelRow struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Size  int    `json:"size"`
	State string `json:"state"`
	Pct   int    `json:"pct"`
	Err   string `json:"err"`
}

func (a *App) modelRows() string {
	activeFile := filepath.Base(a.snapshot().Model)
	var rows []modelRow
	known := false
	for _, m := range modelCatalog {
		row := modelRow{ID: m.ID, Name: m.NameKey, Desc: strS(m.DescKey), Size: m.SizeMB}
		path := filepath.Join("models", m.File)
		_, statErr := os.Stat(path)
		dlMu.Lock()
		st := dl[m.ID]
		dlMu.Unlock()
		switch {
		case st != nil && st.active:
			row.State = "downloading"
			row.Pct = st.pct
		case m.File == activeFile && statErr == nil:
			row.State = "active"
			known = true
		case statErr == nil:
			row.State = "installed"
		default:
			if m.File == activeFile {
				known = true
			}
			row.State = "absent"
			if st != nil && st.err != "" {
				row.Err = st.err
			}
		}
		rows = append(rows, row)
	}
	if !known {
		rows = append(rows, modelRow{ID: "custom", Name: activeFile, Desc: strS("S_M_CUSTOM"), State: "active"})
	}
	out, _ := json.Marshal(rows)
	return string(out)
}

func (a *App) downloadModel(id string) {
	m := findModel(id)
	if m == nil {
		return
	}
	a.startDownload(id, m.File, modelBaseURL+m.File)
}

func (a *App) startDownload(key, file, url string) {
	dlMu.Lock()
	if st := dl[key]; st != nil && st.active {
		dlMu.Unlock()
		return
	}
	dl[key] = &dlState{active: true}
	dlMu.Unlock()

	go func() {
		err := a.doDownload(key, file, url)
		dlMu.Lock()
		if err != nil {
			log.Printf("скачивание %s: %v", file, err)
			dl[key] = &dlState{err: err.Error()}
		} else {
			log.Printf("модель %s скачана", file)
			dl[key] = &dlState{pct: 100}
		}
		dlMu.Unlock()
	}()
}

func freeDiskMB(dir string) int {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return -1
	}
	var free uint64
	p, _ := windows.UTF16PtrFromString(abs)
	if err := windows.GetDiskFreeSpaceEx(p, &free, nil, nil); err != nil {
		return -1
	}
	return int(free / (1024 * 1024))
}

func cleanupStaleParts() {
	entries, err := os.ReadDir("models")
	if err != nil {
		return
	}
	dlMu.Lock()
	activeAny := false
	for _, st := range dl {
		if st.active {
			activeAny = true
		}
	}
	dlMu.Unlock()
	if activeAny {
		return
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".part") {
			continue
		}
		p := filepath.Join("models", e.Name())
		if info, ierr := e.Info(); ierr == nil {
			log.Printf("удаляю незавершённую загрузку %s (%d МБ)", e.Name(), info.Size()/(1024*1024))
		}
		_ = os.Remove(p)
	}
}

func (a *App) doDownload(key, file, url string) error {
	if err := os.MkdirAll("models", 0o755); err != nil {
		return err
	}
	tmp := filepath.Join("models", file+".part")
	final := filepath.Join("models", file)

	client := &http.Client{Timeout: 0}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	total := resp.ContentLength
	if total > 0 {
		if free := freeDiskMB("models"); free >= 0 && int64(free) < total/(1024*1024)+512 {
			return fmt.Errorf("%s", trf("err.disk.space", free, total/(1024*1024)))
		}
	}

	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	var written int64
	buf := make([]byte, 256*1024)
	lastUpd := time.Now()
	for {
		n, rerr := resp.Body.Read(buf)
		if n > 0 {
			if _, werr := f.Write(buf[:n]); werr != nil {
				f.Close()
				os.Remove(tmp)
				return werr
			}
			written += int64(n)
			if total > 0 && time.Since(lastUpd) > 300*time.Millisecond {
				lastUpd = time.Now()
				dlMu.Lock()
				if st := dl[key]; st != nil {
					st.pct = int(written * 100 / total)
				}
				dlMu.Unlock()
			}
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			f.Close()
			os.Remove(tmp)
			return rerr
		}
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	os.Remove(final)
	return os.Rename(tmp, final)
}

func (a *App) deleteModel(id string) string {
	m := findModel(id)
	if m == nil {
		return ""
	}
	if filepath.Base(a.snapshot().Model) == m.File {
		return tr("model.del.active")
	}
	dlMu.Lock()
	if st := dl[id]; st != nil && st.active {
		dlMu.Unlock()
		return ""
	}
	dlMu.Unlock()
	path := filepath.Join("models", m.File)
	if err := os.Remove(path); err != nil {
		return err.Error()
	}
	log.Printf("модель %s удалена", m.File)
	return tr("model.del.ok")
}
