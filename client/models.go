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
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sys/windows"

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
	TrLangs   string
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

var modelCatalog = []modelInfo{
	{ID: "tiny", File: "ggml-tiny.bin", SizeMB: 74, NameKey: "Whisper Tiny", DescKey: "S_M_TINY",
		Hashes: map[string]string{"ggml-tiny.bin": "be07e048e1e599ad46341c8d2a135645097a538221678b7acdd1b1919c6e1b21"},
		Engine: engineWhisper, Langs: "*", Auto: true, Translate: true, Speed: 5, Accuracy: 1},
	{ID: "base", File: "ggml-base.bin", SizeMB: 142, NameKey: "Whisper Base", DescKey: "S_M_BASE",
		Hashes: map[string]string{"ggml-base.bin": "60ed5bc3dd14eea856493d334349b405782ddcaf0028d4b5df4088345fba2efe"},
		Engine: engineWhisper, Langs: "*", Auto: true, Translate: true, Speed: 5, Accuracy: 2},
	{ID: "small", File: "ggml-small.bin", SizeMB: 466, NameKey: "Whisper Small", DescKey: "S_M_SMALL",
		Hashes: map[string]string{"ggml-small.bin": "1be3a9b2063867b937e64e2ec7483364a79917e157fa98c5d94b5c1fffea987b"},
		Engine: engineWhisper, Langs: "*", Auto: true, Translate: true, Speed: 3, Accuracy: 3},
	{ID: "medium-q5_0", File: "ggml-medium-q5_0.bin", SizeMB: 539, NameKey: "Whisper Medium (q5)", DescKey: "S_M_MED",
		Hashes: map[string]string{"ggml-medium-q5_0.bin": "19fea4b380c3a618ec4723c3eef2eb785ffba0d0538cf43f8f235e7b3b34220f"},
		Engine: engineWhisper, Langs: "*", Auto: true, Translate: true, Speed: 2, Accuracy: 4},
	{ID: "large-v3-turbo-q5_0", File: "ggml-large-v3-turbo-q5_0.bin", SizeMB: 574, NameKey: "Whisper Large v3 Turbo (q5)", DescKey: "S_M_TURBO",
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
	{ID: "parakeet-v3", SizeMB: 670, NameKey: "Parakeet TDT 0.6B v3", DescKey: "S_M_PARAKEET",
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
	{ID: "canary-180m", SizeMB: 198, NameKey: "Canary 180M Flash", DescKey: "S_M_CANARY",
		Engine: engineSherpa, Dir: "canary-180m", Langs: "en,de,es,fr", Punct: true,
		Translate: true, TrLangs: "en,de,es,fr", Speed: 5, Accuracy: 3,
		Files: []string{"encoder.int8.onnx", "decoder.int8.onnx", "tokens.txt"},
		Hashes: map[string]string{
			"encoder.int8.onnx": "7a75b4e2a5857a6dcc0819503bbe3fad66943db4a3ccf21d3f27c633667d303f",
			"decoder.int8.onnx": "e41a2ab9c0c2fe81a1e8ade5a45fb02a74bc4db7d1f91b89a54a25e2cf79cba2",
			"tokens.txt":        "2dae6fc7815f9640645e0c765522b278ee0cef49b482d91f6913e334628d3e77",
		},
		BaseURL: "https://huggingface.co/csukuangfj/sherpa-onnx-nemo-canary-180m-flash-en-es-de-fr-int8/resolve/main/"},
	{ID: "qwen3-asr", SizeMB: 937, NameKey: "Qwen3 ASR 0.6B", DescKey: "S_M_QWEN3",
		Engine: engineSherpa, Dir: "qwen3-asr", Langs: "*", Auto: true, Punct: true,
		Speed: 2, Accuracy: 5,
		Files: []string{"conv_frontend.onnx", "decoder.int8.onnx", "encoder.int8.onnx",
			"tokenizer/merges.txt", "tokenizer/tokenizer_config.json", "tokenizer/vocab.json"},
		Hashes: map[string]string{
			"conv_frontend.onnx":              "d22dc4423e0940e49884e903d2ea2f7e5567c14fc1aed97e4e26d6b8f208ef9e",
			"decoder.int8.onnx":               "4f6885be5959ae26af3089d38ee7972c5fafbeeb1cf8d5e76eab6d8b61ca5771",
			"encoder.int8.onnx":               "60748d3e6744a57c9c91e1b17424a6c2990567e8adceb0783940c03ed98fa9d9",
			"tokenizer/merges.txt":            "8831e4f1a044471340f7c0a83d7bd71306a5b867e95fd870f74d0c5308a904d5",
			"tokenizer/tokenizer_config.json": "4942d005604266809309cabc9f4e9cb89ce855d59b14681fdc0e1cc62ea26c4c",
			"tokenizer/vocab.json":            "ca10d7e9fb3ed18575dd1e277a2579c16d108e32f27439684afa0e10b1440910",
		},
		BaseURL: "https://huggingface.co/csukuangfj2/sherpa-onnx-qwen3-asr-0.6B-int8-2026-03-25/resolve/main/"},
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
				log.Printf("models folder: found your own model %s (%d MB)", name, info.Size()/(1024*1024))
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
			log.Printf("models folder: found your own model %s (%d MB)", name, size)
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
	TrL    string   `json:"trlangs"`
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
			Trans: m.Translate, TrL: m.TrLangs, RAM: m.ramEstimateMB(), Fit: ramFit(m.ramEstimateMB(), freeRAM),
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
			log.Printf("download of %s cancelled", m.ID)
			delete(dl, key)
		case err != nil:
			log.Printf("download of %s: %v", m.ID, err)
			dl[key] = &dlState{err: humanError(err)}
		default:
			log.Printf("model %s downloaded in full", m.ID)
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
		return uiErr(fmt.Sprintf("not enough disk space: %d MB free, %d MB needed", free, m.SizeMB), trf("err.disk.space", free, m.SizeMB))
	}
	total := int64(m.SizeMB) * 1024 * 1024
	var doneBytes int64
	for _, f := range m.Files {
		if sub := filepath.Dir(filepath.Join(dir, f)); sub != dir {
			if err := os.MkdirAll(sub, 0o755); err != nil {
				return err
			}
		}
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
				log.Printf("hash mismatch: %v", verr)
				return uiErr("the downloaded file does not match its checksum", tr("err.hash"))
			}
		}
		doneBytes += written
		log.Printf("model %s: file %s ready (%d MB)", m.ID, f, written/(1024*1024))
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
		log.Printf("%s: resuming from %d MB", filepath.Base(final), have/(1024*1024))
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
		return 0, fmt.Errorf("the file was not received in full: %d of %d bytes", written, expected)
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
			log.Printf("download of %s cancelled", file)
			delete(dl, key)
		case err != nil:
			log.Printf("download of %s: %v", file, err)
			dl[key] = &dlState{err: humanError(err)}
		default:
			log.Printf("model %s downloaded", file)
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
				log.Printf("unfinished download %s (%d MB) kept — it can be resumed", p, info.Size()/(1024*1024))
				continue
			}
			log.Printf("removing the unfinished download %s (%d MB)", p, info.Size()/(1024*1024))
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
			log.Printf("%s: %d MB downloaded, ~%d MB to go", file, int(fi.Size()/(1024*1024)), need)
		}
		if free := freeDiskMB("models"); free >= 0 && free < need+512 {
			return uiErr(fmt.Sprintf("not enough disk space: %d MB free, %d MB needed", free, need), trf("err.disk.space", free, need))
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
				log.Printf("hash mismatch: %v", verr)
				return uiErr("the downloaded file does not match its checksum", tr("err.hash"))
			}
			log.Printf("%s: hash matches", file)
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
		log.Printf("deleting the active model %s as confirmed — recognition stops until another one is picked", m.ID)
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
		log.Printf("model %s deleted", m.ID)
		return tr("model.del.ok")
	}
	if err := os.Remove(filepath.Join("models", m.File)); err != nil {
		return humanError(err)
	}
	log.Printf("model %s deleted", m.File)
	return tr("model.del.ok")
}

type stateModelRow struct {
	Model   string `json:"model"`
	Langs   string `json:"langs"`
	State   string `json:"state"`
	Current bool   `json:"current"`
}

type widthResult struct {
	Prev  int `json:"prev"`
	Width int `json:"width"`
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
	Loaded      string          `json:"loaded_now"`
	WeekLine    string          `json:"week_line"`
	LLMOK       bool            `json:"llm_ok"`
	MicOK       bool            `json:"mic_ok"`
	StatusLine  string          `json:"status_line"`
	Remote      bool            `json:"remote"`
	SttSource   string          `json:"stt_source"`
	SttBroken   bool            `json:"stt_broken"`
	WhisperNow  bool            `json:"whisper_now"`
	BackendErr  string          `json:"backend_err"`
	PostErr     string          `json:"post_err"`
	Enabled     bool            `json:"enabled"`
	Autostart   bool            `json:"autostart"`
	DiskMB      int             `json:"disk_mb"`
	WeekChars   int             `json:"week_chars"`
	WeekCount   int             `json:"week_count"`
	TodayCount  int             `json:"today_count"`
	WeekApps    []appShare      `json:"week_apps"`
	RAMFreeMB   int             `json:"ram_free_mb"`
	RAMTotalMB  int             `json:"ram_mb"`
	Idle        []idleModel     `json:"idle_models"`
	PostModel   string          `json:"post_model"`
	PostSizeMB  int             `json:"post_size"`
	PostPrompts int             `json:"post_prompts"`
	PostRemote  bool            `json:"post_remote"`
	UpdChecked  int64           `json:"upd_checked"`
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
	if cfg.SttSource == "remote" {
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
	if m.Engine == engineWhisper && cfg.SttSource == "remote" {
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

type idleModel struct {
	Name   string `json:"name"`
	SizeMB int    `json:"size"`
}

// idleModels names what is downloaded but serves no language: the summary
// shows them so the disk space has a face.
func idleModels(cfg *Config) []idleModel {
	busy := map[string]bool{}
	for _, row := range assignedModelRows(cfg) {
		busy[row.Model] = true
	}
	var out []idleModel
	for i := range modelCatalog {
		m := &modelCatalog[i]
		if !m.installed() || busy[m.NameKey] {
			continue
		}
		out = append(out, idleModel{Name: m.NameKey, SizeMB: m.SizeMB})
	}
	return out
}

func postModelName(cfg *Config) string {
	name := filepath.Base(strings.TrimSpace(cfg.LLMModel))
	if name == "." || name == string(filepath.Separator) {
		return ""
	}
	return strings.TrimSuffix(name, ".gguf")
}

func postModelSizeMB(cfg *Config) int {
	name := filepath.Base(strings.TrimSpace(cfg.LLMModel))
	if name == "" {
		return 0
	}
	info, err := os.Stat(filepath.Join("models", name))
	if err != nil {
		return 0
	}
	return int(info.Size() / (1024 * 1024))
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
	verdict := a.lastVerdict
	enabled := a.enabled
	postErr := a.postErr
	postErrProf := a.postErrProf
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
	if !cfg.PostEnabled {
		llm = strS("S_POST_OFF")
	} else if postAPIOn(cfg) {
		llm = strS("S_POSTAPI_BADGE")
	} else if cfg.PostSource == "local" && llmInstalled(cfg) {
		llm = filepath.Base(cfg.LLMModel)
	}
	total, free := ramMB()
	sttDown := sttRemoteBroken(cfg) && primaryEngine(cfg) == engineWhisper
	if sttDown {
		ready = false
	}
	status := tr("status.loading")
	if sttDown {
		status = strS("S_SRV_DOWN")
	} else if ready {
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
		if verdict != "" {
			parts = append(parts, tr("snd."+verdict))
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
		RAM:         trf("state.ram.free", free),
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
		Loaded:      a.loadedModelsLine(),
		WeekLine:    weekLine(),
		LLMOK:       postReady(cfg),
		MicOK:       rec != nil,
		StatusLine:  statusLine(cfg, ready, free),
		Remote:      cfg.SttSource == "remote",
		SttSource:   cfg.SttSource,
		SttBroken:   sttRemoteBroken(cfg),
		WhisperNow:  primaryEngine(cfg) == engineWhisper,
		BackendErr:  backendErr,
		PostErr:     postErrLine(postErrProf, postErr),
		Enabled:     enabled,
		Autostart:   autorunEnabled(),
		DiskMB:      installedDiskMB(),
		RAMFreeMB:   free,
		RAMTotalMB:  total,
		Idle:        idleModels(cfg),
		PostModel:   postModelName(cfg),
		PostSizeMB:  postModelSizeMB(cfg),
		PostPrompts: len(cfg.ActiveProfiles),
		PostRemote:  postAPIOn(cfg),
	}
	a.mu.Lock()
	st.UpdVersion = a.updVer
	st.UpdChecked = a.updChecked.UnixMilli()
	a.mu.Unlock()
	st.WeekCount, st.WeekChars = histStore.Stats(time.Now().Add(-7 * 24 * time.Hour).UnixMilli())
	st.TodayCount, _ = histStore.Stats(startOfDayMs())
	st.WeekApps = weekByApp()
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

func postErrLine(profile, msg string) string {
	if msg == "" {
		return ""
	}
	if profile == "" {
		return msg
	}
	return profile + ": " + msg
}

func itoaSafe(n int) string {
	return fmt.Sprintf("%d", n)
}

// loadedModelsLine names what actually sits in memory right now, so the
// free-memory number next to it means something.
func (a *App) loadedModelsLine() string {
	var names []string
	for p := range a.loadedModelPaths() {
		names = append(names, modelNameForPath(p))
	}
	sort.Strings(names)
	if len(names) == 0 {
		return tr("state.loaded.none")
	}
	return strings.Join(names, " + ")
}

type appShare struct {
	App   string `json:"app"`
	Count int    `json:"count"`
}

func startOfDayMs() int64 {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).UnixMilli()
}

// weekByApp counts where the week went: which programs took the dictations,
// biggest first, so the summary can draw them.
func weekByApp() []appShare {
	since := time.Now().Add(-7 * 24 * time.Hour).UnixMilli()
	counts := map[string]int{}
	for _, it := range histStore.Items() {
		if it.At < since {
			continue
		}
		app := strings.TrimSpace(it.App)
		if app == "" {
			app = "—"
		}
		counts[app]++
	}
	out := make([]appShare, 0, len(counts))
	for app, n := range counts {
		out = append(out, appShare{App: app, Count: n})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].App < out[j].App
	})
	if len(out) > 6 {
		rest := 0
		for _, s := range out[5:] {
			rest += s.Count
		}
		out = append(out[:5], appShare{App: strS("S_WEEK_OTHER"), Count: rest})
	}
	return out
}

// installedDiskMB adds up what the downloaded models weigh.
func installedDiskMB() int {
	total := 0
	for _, dir := range []string{"models"} {
		entries, _ := os.ReadDir(dir)
		for _, e := range entries {
			if e.IsDir() {
				total += dirSizeMB(filepath.Join(dir, e.Name()))
				continue
			}
			if info, err := e.Info(); err == nil {
				total += int(info.Size() / (1024 * 1024))
			}
		}
	}
	return total
}

func dirSizeMB(dir string) int {
	total := int64(0)
	filepath.Walk(dir, func(_ string, info os.FileInfo, err error) error {
		if err == nil && info != nil && !info.IsDir() {
			total += info.Size()
		}
		return nil
	})
	return int(total / (1024 * 1024))
}

func weekLine() string {

	since := time.Now().Add(-7 * 24 * time.Hour).UnixMilli()
	n, chars := histStore.Stats(since)
	if n == 0 {
		return ""
	}
	return trf("state.week", n, chars)
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
	if sttRemoteBroken(cfg) && primaryEngine(cfg) == engineWhisper {
		return strS("S_SRV_DOWN")
	}
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
			rel := filepath.ToSlash(p)
			if m.Dir != "" {
				rel = strings.TrimPrefix(rel, "models/"+m.Dir+"/")
			}
			want := m.Hashes[rel]
			if want == "" {
				want = m.Hashes[filepath.Base(p)]
			}
			if want == "" {
				continue
			}
			if err := checksum.Verify(p, want); err != nil {
				log.Printf("model check: %v", err)
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
	log.Printf("model check: %d checked, %d broken", checked, broken)
	out, _ := json.Marshal(map[string]any{"rows": rows, "text": text, "ok": broken == 0})
	return string(out)
}
