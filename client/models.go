package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"holdtotype/internal/checksum"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sys/windows"

	"holdtotype/internal/advisor"
)

type modelInfo struct {
	ID        string
	File      string
	SizeMB    int
	NameKey   string
	DescKey   string
	Engine    string
	Dir       string
	Files     []string
	BaseURL   string
	Langs     string
	Punct     bool
	Translate bool
	Speed     int
	Hashes    map[string]string
	Accuracy  int
	Auto      bool
	Manual    bool
	LinkURL   string
	Custom    bool
}

func (m *modelInfo) ramEstimateMB() int {
	if m.Engine == engineSherpa || m.Engine == engineStream {
		return m.SizeMB * 12 / 10
	}
	return m.SizeMB*15/10 + 60
}

func advisorCatalog() []advisor.Model {
	out := make([]advisor.Model, 0, len(modelCatalog))
	for i := range modelCatalog {
		m := &modelCatalog[i]
		if m.Custom || (m.Manual && !m.installed()) {
			continue
		}
		langs := strings.Split(m.Langs, ",")
		out = append(out, advisor.Model{
			ID: m.ID, Engine: m.Engine, Langs: langs, SizeMB: m.SizeMB,
			RAMMB: m.ramEstimateMB(), Punct: m.Punct, Translate: m.Translate,
			Speed: m.Speed, Accuracy: m.Accuracy,
		})
	}
	return out
}

var modelCatalog = []modelInfo{
	{ID: "base", File: "ggml-base.bin", SizeMB: 142, NameKey: "Base", DescKey: "S_M_BASE",
		Hashes: map[string]string{"ggml-base.bin": "60ed5bc3dd14eea856493d334349b405782ddcaf0028d4b5df4088345fba2efe"},
		Engine: engineWhisper, Langs: "*", Auto: true, Translate: true, Speed: 5, Accuracy: 2},
	{ID: "small", File: "ggml-small.bin", SizeMB: 466, NameKey: "Small", DescKey: "S_M_SMALL",
		Hashes: map[string]string{"ggml-small.bin": "1be3a9b2063867b937e64e2ec7483364a79917e157fa98c5d94b5c1fffea987b"},
		Engine: engineWhisper, Langs: "*", Auto: true, Translate: true, Speed: 3, Accuracy: 3},
	{ID: "medium-q5_0", File: "ggml-medium-q5_0.bin", SizeMB: 539, NameKey: "Medium (q5)", DescKey: "S_M_MED",
		Hashes: map[string]string{"ggml-medium-q5_0.bin": "19fea4b380c3a618ec4723c3eef2eb785ffba0d0538cf43f8f235e7b3b34220f"},
		Engine: engineWhisper, Langs: "*", Auto: true, Translate: true, Speed: 2, Accuracy: 4},
	{ID: "large-v3-turbo-q5_0", File: "ggml-large-v3-turbo-q5_0.bin", SizeMB: 574, NameKey: "Turbo (q5)", DescKey: "S_M_TURBO",
		Hashes: map[string]string{"ggml-large-v3-turbo-q5_0.bin": "394221709cd5ad1f40c46e6031ca61bce88931e6e088c188294c6d5a55ffa7e2"},
		Engine: engineWhisper, Langs: "*", Auto: true, Translate: false, Speed: 4, Accuracy: 4},
	{ID: "gigaam-v3", SizeMB: 232, NameKey: "GigaAM v3", DescKey: "S_M_GIGAAM",
		Engine: engineSherpa, Dir: "gigaam-v3", Langs: "ru", Punct: true, Speed: 5, Accuracy: 5,
		Files: []string{"encoder.int8.onnx", "decoder.onnx", "joiner.onnx", "tokens.txt"},
		Hashes: map[string]string{
			"encoder.int8.onnx": "369f35a71bf288d3b8e0391fabd8dba5f2314088d440bca474056b7b4b6e66bf",
			"decoder.onnx":      "38fc7475443ea2a26f63211ca350f73ac50fff824ab7a3876ee2bd610c53bbc4",
			"joiner.onnx":       "602ff7017a93311aad34df1437c8d7f49911353c13d6eae7a6ee7b041339465c",
			"tokens.txt":        "39abae20e692998290c574e606f11a9edef2902a1995463fcff63d1490cf22b7",
		},
		BaseURL: "https://huggingface.co/csukuangfj/sherpa-onnx-nemo-transducer-punct-giga-am-v3-russian-2025-12-16/resolve/main/"},
	{ID: "gigaam-v2", SizeMB: 231, NameKey: "GigaAM v2", DescKey: "S_M_GIGAAM2",
		Engine: engineSherpa, Dir: "gigaam-v2", Langs: "ru", Speed: 5, Accuracy: 4,
		Files: []string{"encoder.int8.onnx", "decoder.onnx", "joiner.onnx", "tokens.txt"},
		Hashes: map[string]string{
			"encoder.int8.onnx": "b51efc61e3c0037ad1cb804079975468de3d175324fe8323aef5be4f5c6a38a1",
			"decoder.onnx":      "208e24cc150fb0ebca3fab169502796daa12e0255dcf7b4acf65015c436e9f76",
			"joiner.onnx":       "4b02eced18e033fc5173e6c47b6ab166b5efea8d35c3f33a6755ff0d622fb5b0",
			"tokens.txt":        "17cc514451bcceac9c280068c71502f8448f99e9fb1456b8d0761651fd0392f2",
		},
		BaseURL: "https://huggingface.co/csukuangfj/sherpa-onnx-nemo-transducer-giga-am-v2-russian-2025-04-19/resolve/main/"},
	{ID: "parakeet-v3", SizeMB: 670, NameKey: "Parakeet v3", DescKey: "S_M_PARAKEET",
		Engine: engineSherpa, Dir: "parakeet-v3", Auto: true,
		Langs: "en,es,fr,de,bg,hr,cs,da,nl,et,fi,el,hu,it,lv,lt,mt,pl,pt,ro,sk,sl,sv,ru,uk",
		Punct: true, Speed: 5, Accuracy: 4,
		Files: []string{"encoder.int8.onnx", "decoder.int8.onnx", "joiner.int8.onnx", "tokens.txt"},
		Hashes: map[string]string{
			"encoder.int8.onnx": "acfc2b4456377e15d04f0243af540b7fe7c992f8d898d751cf134c3a55fd2247",
			"decoder.int8.onnx": "179e50c43d1a9de79c8a24149a2f9bac6eb5981823f2a2ed88d655b24248db4e",
			"joiner.int8.onnx":  "3164c13fc2821009440d20fcb5fdc78bff28b4db2f8d0f0b329101719c0948b3",
			"tokens.txt":        "d58544679ea4bc6ac563d1f545eb7d474bd6cfa467f0a6e2c1dc1c7d37e3c35d",
		},
		BaseURL: "https://huggingface.co/csukuangfj/sherpa-onnx-nemo-parakeet-tdt-0.6b-v3-int8/resolve/main/"},
	{ID: "nemotron-3.5", SizeMB: 651, NameKey: "Nemotron 3.5", DescKey: "S_M_NEMOTRON",
		Engine: engineStream, Dir: "nemotron-3.5", Langs: "*", Auto: true, Punct: true,
		Speed: 3, Accuracy: 5,
		Files: []string{"encoder.int8.onnx", "decoder.int8.onnx", "joiner.int8.onnx", "tokens.txt"},
		Hashes: map[string]string{
			"encoder.int8.onnx": "874275f509c86e331eb1c1f1e2f7fa48c39144de94f98ccbe6ae3fcdc18df38a",
			"decoder.int8.onnx": "19f9c98fc6d0a2c33a65a43b36fdb2e914c26c0aa9764be3aebc502a1e982fb0",
			"joiner.int8.onnx":  "4101c7c679a0bc30483794b27a059e34e79232aa2068d78d51231a22c8b0d7ce",
			"tokens.txt":        "32be3ebfabfff475d64d7829b435f1c7856a1c497907def5c41d54ca9f1eccfd",
		},
		BaseURL: "https://huggingface.co/Masterx/sherpa-onnx-nemotron-3.5-asr-streaming-0.6b-560ms-2026-06-11/resolve/main/"},
	{ID: "moonshine-uk", SizeMB: 135, NameKey: "Moonshine Base uk", DescKey: "S_M_MOONUK",
		Engine: engineSherpa, Dir: "moonshine-uk", Langs: "uk", Speed: 5, Accuracy: 3,
		Manual:  true,
		LinkURL: "https://github.com/k2-fsa/sherpa-onnx/releases/download/asr-models/sherpa-onnx-moonshine-base-uk-quantized-2026-02-27.tar.bz2"},
}

const defaultPresetModel = "medium-q5_0"

const modelBaseURL = "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/"

var localScanOnce sync.Once

func scanLocalModels() {
	localScanOnce.Do(func() {
		entries, err := os.ReadDir("models")
		if err != nil {
			return
		}
		knownFile := map[string]bool{}
		knownDir := map[string]bool{}
		for i := range modelCatalog {
			if modelCatalog[i].File != "" {
				knownFile[modelCatalog[i].File] = true
			}
			if modelCatalog[i].Dir != "" {
				knownDir[modelCatalog[i].Dir] = true
			}
		}
		for _, e := range entries {
			name := e.Name()
			if !e.IsDir() {
				if !strings.HasSuffix(name, ".bin") || knownFile[name] {
					continue
				}
				info, ierr := e.Info()
				if ierr != nil {
					continue
				}
				modelCatalog = append(modelCatalog, modelInfo{
					ID: "local:" + name, File: name, SizeMB: int(info.Size() / (1024 * 1024)),
					NameKey: strings.TrimSuffix(strings.TrimPrefix(name, "ggml-"), ".bin"),
					DescKey: "S_M_LOCAL", Engine: engineWhisper, Langs: "*", Custom: true,
				})
				log.Printf("папка models: найдена своя модель %s (%d МБ)", name, info.Size()/(1024*1024))
				continue
			}
			if knownDir[name] {
				continue
			}
			dir := filepath.Join("models", name)
			if _, aerr := sherpaModelArgs(dir); aerr != nil {
				continue
			}
			size := 0
			if items, derr := os.ReadDir(dir); derr == nil {
				for _, it := range items {
					if fi, ferr := it.Info(); ferr == nil {
						size += int(fi.Size() / (1024 * 1024))
					}
				}
			}
			modelCatalog = append(modelCatalog, modelInfo{
				ID: "local:" + name, Dir: name, SizeMB: size, NameKey: name,
				DescKey: "S_M_LOCAL", Engine: engineSherpa, Langs: "*", Custom: true,
			})
			log.Printf("папка models: найдена своя модель %s (%d МБ)", name, size)
		}
	})
}

func (m *modelInfo) paths() []string {
	if m.Dir != "" {
		out := make([]string, 0, len(m.Files))
		for _, f := range m.Files {
			out = append(out, filepath.Join("models", m.Dir, f))
		}
		return out
	}
	return []string{filepath.Join("models", m.File)}
}

func (m *modelInfo) installed() bool {
	if m.Dir != "" && len(m.Files) == 0 {
		_, err := sherpaModelArgs(filepath.Join("models", m.Dir))
		return err == nil
	}
	for _, p := range m.paths() {
		if _, err := os.Stat(p); err != nil {
			return false
		}
	}
	return true
}

func (m *modelInfo) isActive(cfg *Config) bool {
	return m.ID == presetModelID(cfg, cfg.Language)
}

func servesLangs(cfg *Config, id string) []string {
	var out []string
	if cfg.LangModels["auto"] == id {
		out = append(out, "auto")
	}
	for _, l := range translateLangCodes() {
		if cfg.LangModels[l] == id {
			out = append(out, l)
		}
	}
	return out
}

type dlState struct {
	active bool
	pct    int
	err    string
	cancel context.CancelFunc
}

const partKeep = 7 * 24 * time.Hour

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
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Desc   string   `json:"desc"`
	Size   int      `json:"size"`
	State  string   `json:"state"`
	Pct    int      `json:"pct"`
	Err    string   `json:"err"`
	Engine string   `json:"engine"`
	Langs  string   `json:"langs"`
	Punct  bool     `json:"punct"`
	Trans  bool     `json:"translate"`
	RAM    int      `json:"ram"`
	Fit    string   `json:"fit"`
	Serves []string `json:"serves"`
	Speed  int      `json:"speed"`
	Acc    int      `json:"accuracy"`
	Auto   bool     `json:"auto"`
	Manual bool     `json:"manual"`
	Custom bool     `json:"custom"`
	Loaded bool     `json:"loaded"`
	Link   string   `json:"link"`
}

func ramFit(needMB, freeMB int) string {
	if freeMB <= 0 || needMB <= 0 {
		return ""
	}
	switch {
	case needMB*10 <= freeMB*6:
		return "ok"
	case needMB*10 <= freeMB*9:
		return "warn"
	default:
		return "bad"
	}
}

func (a *App) modelRows() string {
	cfg := a.snapshot()
	_, freeRAM := ramMB()
	loaded := a.loadedModelPaths()
	var rows []modelRow
	for i := range modelCatalog {
		m := &modelCatalog[i]
		row := modelRow{
			ID: m.ID, Name: m.NameKey, Desc: strS(m.DescKey), Size: m.SizeMB,
			Engine: m.Engine, Langs: m.Langs, Punct: m.Punct,
			Trans: m.Translate, RAM: m.ramEstimateMB(), Fit: ramFit(m.ramEstimateMB(), freeRAM),
			Speed: m.Speed, Acc: m.Accuracy, Auto: m.Auto, Manual: m.Manual,
			Custom: m.Custom, Link: m.LinkURL,
			Serves: servesLangs(cfg, m.ID),
		}
		have := m.installed()
		active := m.isActive(cfg)
		row.Loaded = have && loaded[m.modelPath()]
		dlMu.Lock()
		st := dl[m.ID]
		dlMu.Unlock()
		switch {
		case st != nil && st.active:
			row.State = "downloading"
			row.Pct = st.pct
		case active && have:
			row.State = "active"
		case have:
			row.State = "installed"
		default:
			row.State = "absent"
			if st != nil && st.err != "" {
				row.Err = st.err
			}
		}
		rows = append(rows, row)
	}
	out, _ := json.Marshal(rows)
	return string(out)
}

func (m *modelInfo) modelPath() string {
	if m.Dir != "" {
		return "models/" + m.Dir
	}
	return "models/" + m.File
}

func (a *App) loadedModelPaths() map[string]bool {
	out := map[string]bool{}
	a.mu.Lock()
	srv := a.srv
	ready := a.ready
	a.mu.Unlock()
	if srv != nil && ready && !srv.external() {
		out[filepath.ToSlash(srv.model())] = true
	}
	a.altMu.Lock()
	if a.alt != nil {
		out[filepath.ToSlash(a.alt.model())] = true
	}
	a.altMu.Unlock()
	return out
}

func (a *App) downloadModel(id string) {
	m := findModel(id)
	if m == nil {
		return
	}
	if m.Dir != "" {
		a.startMultiDownload(id, m)
		return
	}
	a.startDownload(id, m.File, modelBaseURL+m.File, m.SizeMB)
}

func (a *App) startMultiDownload(key string, m *modelInfo) {
	dlMu.Lock()
	if st := dl[key]; st != nil && st.active {
		dlMu.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	dl[key] = &dlState{active: true, cancel: cancel}
	dlMu.Unlock()

	go func() {
		defer cancel()
		err := a.doMultiDownload(ctx, key, m)
		dlMu.Lock()
		switch {
		case errors.Is(err, context.Canceled):
			log.Printf("скачивание %s отменено", m.ID)
			delete(dl, key)
		case err != nil:
			log.Printf("скачивание %s: %v", m.ID, err)
			dl[key] = &dlState{err: humanError(err)}
		default:
			log.Printf("модель %s скачана целиком", m.ID)
			dl[key] = &dlState{pct: 100}
		}
		dlMu.Unlock()
	}()
}

func (a *App) doMultiDownload(ctx context.Context, key string, m *modelInfo) error {
	dir := filepath.Join("models", m.Dir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	if free := freeDiskMB("models"); free >= 0 && free < m.SizeMB+256 {
		return fmt.Errorf("%s", trf("err.disk.space", free, m.SizeMB))
	}
	total := int64(m.SizeMB) * 1024 * 1024
	var doneBytes int64
	for _, f := range m.Files {
		written, err := downloadFile(ctx, m.BaseURL+f, filepath.Join(dir, f), func(n int64) {
			dlMu.Lock()
			if st := dl[key]; st != nil && total > 0 {
				pct := int((doneBytes + n) * 100 / total)
				if pct > 99 {
					pct = 99
				}
				st.pct = pct
			}
			dlMu.Unlock()
		})
		if err != nil {
			return fmt.Errorf("%s: %w", f, err)
		}
		if want := m.Hashes[f]; want != "" {
			if verr := checksum.Verify(filepath.Join(dir, f), want); verr != nil {
				_ = os.Remove(filepath.Join(dir, f))
				log.Printf("хеш не сошёлся: %v", verr)
				return fmt.Errorf("%s", tr("err.hash"))
			}
		}
		doneBytes += written
		log.Printf("модель %s: файл %s готов (%d МБ)", m.ID, f, written/(1024*1024))
	}
	return nil
}

func downloadFile(ctx context.Context, url, final string, progress func(written int64)) (int64, error) {
	tmp := final + ".part"
	have := int64(0)
	if fi, err := os.Stat(tmp); err == nil {
		have = fi.Size()
	}
	client := &http.Client{Timeout: 0}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	if have > 0 {
		req.Header.Set("Range", "bytes="+strconv.FormatInt(have, 10)+"-")
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	resumed := resp.StatusCode == http.StatusPartialContent
	if resp.StatusCode != http.StatusOK && !resumed {
		return 0, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	if !resumed {
		have = 0
	} else {
		log.Printf("%s: продолжаю с %d МБ", filepath.Base(final), have/(1024*1024))
	}
	var f *os.File
	if resumed {
		f, err = os.OpenFile(tmp, os.O_WRONLY|os.O_APPEND, 0o644)
	} else {
		f, err = os.Create(tmp)
	}
	if err != nil {
		return 0, err
	}
	written := have
	expected := resp.ContentLength
	if expected > 0 {
		expected += have
	}
	buf := make([]byte, 256*1024)
	lastUpd := time.Now()
	for {
		n, rerr := resp.Body.Read(buf)
		if n > 0 {
			if _, werr := f.Write(buf[:n]); werr != nil {
				f.Close()
				os.Remove(tmp)
				return 0, werr
			}
			written += int64(n)
			if progress != nil && time.Since(lastUpd) > 300*time.Millisecond {
				lastUpd = time.Now()
				progress(written)
			}
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			f.Close()
			if ctx.Err() != nil {
				return written, ctx.Err()
			}
			os.Remove(tmp)
			return 0, rerr
		}
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return 0, err
	}
	if expected > 0 && written != expected {
		os.Remove(tmp)
		return 0, fmt.Errorf("файл получен не целиком: %d из %d байт", written, expected)
	}
	os.Remove(final)
	if err := os.Rename(tmp, final); err != nil {
		return 0, err
	}
	if progress != nil {
		progress(written)
	}
	return written, nil
}

func (a *App) startDownload(key, file, url string, sizeMB int) {
	dlMu.Lock()
	if st := dl[key]; st != nil && st.active {
		dlMu.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	dl[key] = &dlState{active: true, cancel: cancel}
	dlMu.Unlock()

	go func() {
		defer cancel()
		err := a.doDownload(ctx, key, file, url, sizeMB)
		dlMu.Lock()
		switch {
		case errors.Is(err, context.Canceled):
			log.Printf("скачивание %s отменено", file)
			delete(dl, key)
		case err != nil:
			log.Printf("скачивание %s: %v", file, err)
			dl[key] = &dlState{err: humanError(err)}
		default:
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
	dirs := []string{"models"}
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, filepath.Join("models", e.Name()))
		}
	}
	for _, dir := range dirs {
		items, derr := os.ReadDir(dir)
		if derr != nil {
			continue
		}
		for _, e := range items {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".part") {
				continue
			}
			p := filepath.Join(dir, e.Name())
			info, ierr := e.Info()
			if ierr != nil {
				_ = os.Remove(p)
				continue
			}
			if time.Since(info.ModTime()) < partKeep {
				log.Printf("незавершённая загрузка %s (%d МБ) сохранена — можно продолжить", p, info.Size()/(1024*1024))
				continue
			}
			log.Printf("удаляю незавершённую загрузку %s (%d МБ)", p, info.Size()/(1024*1024))
			_ = os.Remove(p)
		}
	}
}

func (a *App) doDownload(ctx context.Context, key, file, url string, sizeMB int) error {
	if err := os.MkdirAll("models", 0o755); err != nil {
		return err
	}
	final := filepath.Join("models", file)
	if m := findModel(key); m != nil && m.SizeMB > 0 {
		sizeMB = m.SizeMB
	}
	if sizeMB > 0 {
		need := sizeMB
		if fi, serr := os.Stat(final + ".part"); serr == nil {
			need -= int(fi.Size() / (1024 * 1024))
			if need < 0 {
				need = 0
			}
			log.Printf("%s: уже скачано %d МБ, осталось ~%d МБ", file, int(fi.Size()/(1024*1024)), need)
		}
		if free := freeDiskMB("models"); free >= 0 && free < need+512 {
			return fmt.Errorf("%s", trf("err.disk.space", free, need))
		}
	}
	total := int64(sizeMB) * 1024 * 1024
	_, err := downloadFile(ctx, url, final, func(written int64) {
		if total <= 0 {
			return
		}
		pct := int(written * 100 / total)
		if pct > 99 {
			pct = 99
		}
		dlMu.Lock()
		if st := dl[key]; st != nil {
			st.pct = pct
		}
		dlMu.Unlock()
	})
	if err != nil {
		return err
	}
	if m := findModel(key); m != nil {
		if want := m.Hashes[file]; want != "" {
			if verr := checksum.Verify(final, want); verr != nil {
				_ = os.Remove(final)
				log.Printf("хеш не сошёлся: %v", verr)
				return fmt.Errorf("%s", tr("err.hash"))
			}
			log.Printf("%s: хеш совпал", file)
		}
	}
	return nil
}

func (a *App) deleteModel(id string, force bool) string {
	m := findModel(id)
	if m == nil {
		return ""
	}
	if m.isActive(a.snapshot()) {
		if !force {
			return tr("model.del.active")
		}
		log.Printf("удаляю активную модель %s по подтверждению — распознавание остановится до выбора другой", m.ID)
		a.requestServerRestart()
		time.Sleep(800 * time.Millisecond)
	}
	dlMu.Lock()
	if st := dl[id]; st != nil && st.active {
		dlMu.Unlock()
		return ""
	}
	dlMu.Unlock()
	if m.Dir != "" {
		dir := filepath.Join("models", m.Dir)
		if err := os.RemoveAll(dir); err != nil {
			return humanError(err)
		}
		log.Printf("модель %s удалена", m.ID)
		return tr("model.del.ok")
	}
	if err := os.Remove(filepath.Join("models", m.File)); err != nil {
		return humanError(err)
	}
	log.Printf("модель %s удалена", m.File)
	return tr("model.del.ok")
}

type advicePart struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	SizeMB    int    `json:"size"`
	Installed bool   `json:"installed"`
}

type adviceOut struct {
	Primary   string       `json:"primary"`
	Companion string       `json:"companion"`
	Text      string       `json:"text"`
	RAM       string       `json:"ram"`
	Plan      []advicePart `json:"plan"`
	NeedMB    int          `json:"need"`
}

func advicePartOf(id string) (advicePart, bool) {
	m := findModel(id)
	if m == nil {
		return advicePart{}, false
	}
	return advicePart{ID: m.ID, Name: m.NameKey, SizeMB: m.SizeMB, Installed: m.installed()}, true
}

func modelDisplayName(id string) string {
	if m := findModel(id); m != nil {
		return m.NameKey
	}
	return id
}

func modelNameForPath(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return ""
	}
	base := filepath.Base(filepath.Clean(p))
	for i := range modelCatalog {
		m := &modelCatalog[i]
		if m.File != "" && filepath.Base(m.File) == base {
			return m.NameKey
		}
		if m.Dir != "" && filepath.Base(filepath.Clean(m.Dir)) == base {
			return m.NameKey
		}
	}
	return base
}

func adviseModel(lang, priority string, needTranslate bool) string {
	_, free := ramMB()
	res := advisor.Recommend(advisor.Input{
		Lang: lang, Priority: priority, RAMFreeMB: free, Translate: needTranslate,
	}, advisorCatalog())

	var parts []string
	if res.Primary == "" {
		parts = append(parts, strS("S_ADV_NONE"))
	} else {
		parts = append(parts, trf("adv.pick", modelDisplayName(res.Primary)))
		for _, why := range res.Why {
			switch why {
			case advisor.WhyLanguage:
				parts = append(parts, strS("S_ADV_LANG"))
			case advisor.WhyAccuracy:
				parts = append(parts, strS("S_ADV_ACC"))
			case advisor.WhySpeed:
				parts = append(parts, strS("S_ADV_SPEED"))
			case advisor.WhyRAM:
				parts = append(parts, strS("S_ADV_RAM"))
			}
		}
		if res.Companion != "" {
			parts = append(parts, trf("adv.companion", modelDisplayName(res.Companion)))
		}
	}
	var plan []advicePart
	need := 0
	for _, id := range []string{res.Primary, res.Companion} {
		if id == "" {
			continue
		}
		if p, ok := advicePartOf(id); ok {
			plan = append(plan, p)
			if !p.Installed {
				need += p.SizeMB
			}
		}
	}
	out, _ := json.Marshal(adviceOut{
		Primary:   res.Primary,
		Companion: res.Companion,
		Text:      strings.Join(parts, " "),
		RAM:       trf("adv.ram", free),
		Plan:      plan,
		NeedMB:    need,
	})
	return string(out)
}

type stateModelRow struct {
	Model   string `json:"model"`
	Langs   string `json:"langs"`
	State   string `json:"state"`
	Current bool   `json:"current"`
}

type stateOut struct {
	Hotkey   string `json:"hotkey"`
	Mic      string `json:"mic"`
	Engine   string `json:"engine"`
	LLM      string `json:"llm"`
	RAM      string `json:"ram"`
	Last     string `json:"last"`
	LastMeta string `json:"last_meta"`
	LastAt   int64  `json:"last_at"`
	LastApp  string `json:"last_app"`
	Ready    bool   `json:"ready"`
	Status   string `json:"status"`

	ActiveModel string          `json:"active_model"`
	ActiveState string          `json:"active_state"`
	ActiveLang  string          `json:"active_lang"`
	Assigned    []stateModelRow `json:"assigned"`
	InstalledMs []string        `json:"installed_models"`
	LLMOK       bool            `json:"llm_ok"`
	MicOK       bool            `json:"mic_ok"`
	StatusLine  string          `json:"status_line"`
	Remote      bool            `json:"remote"`
	BackendErr  string          `json:"backend_err"`
	UpdVersion  string          `json:"upd_version"`
	Badges      struct {
		Mic     string `json:"mic"`
		Models  string `json:"models"`
		System  string `json:"system"`
		History string `json:"history"`
	} `json:"badges"`
}

func agoLabel(d time.Duration) string {
	switch {
	case d < time.Minute:
		return tr("ago.now")
	case d < time.Hour:
		return trf("ago.min", int(d.Minutes()))
	default:
		return trf("ago.hour", int(d.Hours()))
	}
}

func installedModelCount() int {
	n := 0
	for i := range modelCatalog {
		if modelCatalog[i].installed() {
			n++
		}
	}
	return n
}

func systemWarnings(cfg *Config) int {
	n := 0
	if cfg.ServerURL != "" {
		n++
	}
	if !cfg.ServerAutostart {
		n++
	}
	return n
}

func stateForModel(cfg *Config, m *modelInfo) string {
	if m == nil {
		return "missing"
	}
	if m.Engine == engineWhisper && strings.TrimSpace(cfg.ServerURL) != "" {
		return "remote"
	}
	if m.installed() {
		return "ready"
	}
	dlMu.Lock()
	st := dl[m.ID]
	dlMu.Unlock()
	if st != nil && st.active {
		return "downloading"
	}
	return "missing"
}

func assignedModelRows(cfg *Config) []stateModelRow {
	langs := append([]string{"auto"}, translateLangCodes()...)
	byModel := map[string][]string{}
	var order []string
	for _, l := range langs {
		id := presetModelID(cfg, l)
		if _, ok := byModel[id]; !ok {
			order = append(order, id)
		}
		byModel[id] = append(byModel[id], l)
	}
	currentID := presetModelID(cfg, cfg.Language)
	rows := make([]stateModelRow, 0, len(order))
	for _, id := range order {
		m := findModel(id)
		if m == nil {
			continue
		}
		label := ""
		if len(byModel[id]) == len(langs) {
			label = strS("S_ALL_LANGS")
		} else {
			var parts []string
			for _, l := range byModel[id] {
				if l == "auto" {
					parts = append(parts, strS("S_RECAUTO"))
				} else {
					parts = append(parts, langLabel(l))
				}
			}
			label = strings.Join(parts, ", ")
		}
		rows = append(rows, stateModelRow{
			Model: m.NameKey, Langs: label,
			State: stateForModel(cfg, m), Current: id == currentID,
		})
	}
	return rows
}

func installedModelNames() []string {
	var out []string
	for i := range modelCatalog {
		if modelCatalog[i].installed() {
			out = append(out, modelCatalog[i].NameKey)
		}
	}
	return out
}

func (a *App) stateSnapshot() string {
	cfg := a.snapshot()
	a.mu.Lock()
	ready := a.ready
	backendErr := a.backendErr
	stamp := a.lastResultAt
	lastApp := a.lastProcess
	target := a.lastTarget
	rec := a.rec
	last := a.lastResult
	a.mu.Unlock()

	mic := strS("S_MIC_DEFAULT")
	if cfg.MicDeviceName != "" {
		mic = cfg.MicDeviceName
	} else if rec != nil {
		for _, d := range rec.devices() {
			if d.System || d.Default {
				mic = d.Name
				if d.System {
					break
				}
			}
		}
	}
	llm := strS("S_NO_LLM")
	if postAPIOn(cfg) {
		llm = strS("S_POSTAPI_BADGE")
	} else if llmInstalled(cfg) {
		llm = filepath.Base(cfg.LLMModel)
	}
	_, free := ramMB()
	status := tr("status.loading")
	if ready {
		status = trf("status.ready", cfg.Hotkey)
	} else if backendErr != "" {
		status = backendErr
	}
	if last == "" {
		last = "—"
	}
	lastMeta := ""
	var lastAt int64
	if stamp.IsZero() {
		last = "—"
	} else {
		lastAt = stamp.UnixMilli()
		parts := []string{agoLabel(time.Since(stamp)), trf("chars", len([]rune(last)))}
		if target != "" {
			parts = append(parts, trf("inserted.into", target))
		}
		lastMeta = strings.Join(parts, " · ")
	}

	active := activeModel(cfg)
	activeName := strS("S_NOT_INSTALLED")
	if active != nil {
		activeName = active.NameKey
	}
	activeLang := strS("S_RECAUTO")
	if l := strings.ToLower(strings.TrimSpace(cfg.Language)); l != "" && l != "auto" {
		activeLang = langLabel(l)
	}
	st := stateOut{
		Hotkey:      cfg.Hotkey,
		Mic:         mic,
		Engine:      primaryEngine(cfg) + " · " + modelNameForPath(activeModelPath(cfg)),
		LLM:         llm,
		RAM:         trf("adv.ram", free),
		Last:        last,
		LastMeta:    lastMeta,
		LastAt:      lastAt,
		LastApp:     lastApp,
		Ready:       ready,
		Status:      status,
		ActiveModel: activeName,
		ActiveState: stateForModel(cfg, active),
		ActiveLang:  activeLang,
		Assigned:    assignedModelRows(cfg),
		InstalledMs: installedModelNames(),
		LLMOK:       postReady(cfg),
		MicOK:       rec != nil,
		StatusLine:  statusLine(cfg, ready, free),
		Remote:      strings.TrimSpace(cfg.ServerURL) != "",
		BackendErr:  backendErr,
	}
	a.mu.Lock()
	st.UpdVersion = a.updVer
	a.mu.Unlock()
	st.Badges.Mic = micBadge(mic)
	st.Badges.Models = itoaSafe(installedModelCount())
	st.Badges.History = histBadge()
	warn := systemWarnings(cfg)
	if st.UpdVersion != "" {
		warn++
	}
	if warn > 0 {
		st.Badges.System = itoaSafe(warn)
	}
	out, _ := json.Marshal(st)
	return string(out)
}

func itoaSafe(n int) string {
	return fmt.Sprintf("%d", n)
}

func shortLabel(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max-1]) + "…"
}

var micNoise = map[string]bool{
	"microphone": true, "mic": true, "audio": true, "input": true, "device": true,
	"микрофон": true, "звук": true, "устройство": true, "system": true, "системный": true,
}

func micBadge(name string) string {
	clean := strings.Map(func(r rune) rune {
		if r == '(' || r == ')' || r == '-' || r == ',' {
			return ' '
		}
		return r
	}, name)
	for _, w := range strings.Fields(clean) {
		if len([]rune(w)) < 3 {
			continue
		}
		if micNoise[strings.ToLower(w)] {
			continue
		}
		return shortLabel(w, 10)
	}
	return shortLabel(strings.TrimSpace(name), 10)
}

func statusLine(cfg *Config, ready bool, freeMB int) string {
	if !ready {
		return tr("status.loading")
	}
	models := modelNameForPath(activeModelPath(cfg))
	if primaryEngine(cfg) == engineSherpa {
		models += " + " + modelNameForPath(cfg.Model)
	}
	return trf("status.line", models, float64(freeMB)/1024)
}

func cancelDownload(id string) bool {
	dlMu.Lock()
	st := dl[id]
	dlMu.Unlock()
	if st == nil || !st.active || st.cancel == nil {
		return false
	}
	st.cancel()
	return true
}

func histBadge() string {
	n := histStore.Count()
	if n == 0 {
		return ""
	}
	return itoaSafe(n)
}

func (a *App) verifyModels() string {
	type row struct {
		Name string `json:"name"`
		OK   bool   `json:"ok"`
		Note string `json:"note"`
	}
	var rows []row
	for i := range modelCatalog {
		m := &modelCatalog[i]
		if !m.installed() || len(m.Hashes) == 0 {
			continue
		}
		bad := ""
		for _, p := range m.paths() {
			want := m.Hashes[filepath.Base(p)]
			if want == "" {
				continue
			}
			if err := checksum.Verify(p, want); err != nil {
				log.Printf("проверка моделей: %v", err)
				bad = filepath.Base(p)
				break
			}
		}
		rows = append(rows, row{Name: m.NameKey, OK: bad == "", Note: bad})
	}
	checked, broken := 0, 0
	for _, r := range rows {
		checked++
		if !r.OK {
			broken++
		}
	}
	text := trf("models.check.ok", checked)
	if checked == 0 {
		text = tr("models.check.none")
	} else if broken > 0 {
		var names []string
		for _, r := range rows {
			if !r.OK {
				names = append(names, r.Name+" ("+r.Note+")")
			}
		}
		text = trf("models.check.bad", strings.Join(names, ", "))
	}
	log.Printf("проверка моделей: проверено %d, повреждено %d", checked, broken)
	out, _ := json.Marshal(map[string]any{"rows": rows, "text": text, "ok": broken == 0})
	return string(out)
}
