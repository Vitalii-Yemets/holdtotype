package main

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/sys/windows/registry"

	"holdtotype/internal/appid"
	"holdtotype/internal/theme"
)

//go:embed payload.zip
var payload []byte

const (
	appName   = appid.Name
	exeName   = appid.Exe
	uninstKey = appid.UninstallKey
	runKey    = `Software\Microsoft\Windows\CurrentVersion\Run`
	runValue  = appid.RunValue
)

var appVersion = appid.Version

func existingInstall() string {
	k, err := registry.OpenKey(registry.CURRENT_USER, uninstKey, registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer k.Close()
	dir, _, err := k.GetStringValue("InstallLocation")
	if err != nil || dir == "" {
		return ""
	}
	if _, err := os.Stat(filepath.Join(dir, exeName)); err != nil {
		return ""
	}
	return dir
}

type modelOpt struct {
	ID      string
	Name    string
	File    string
	SizeMB  int
	Dir      string
	Files    []string
	BaseURL  string
	Lang     string
	PresetID string
}

const whisperBaseURL = "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/"

var modelOpts = []modelOpt{
	{ID: "gigaam-v3", Name: "GigaAM v3", SizeMB: 232, Dir: "gigaam-v3", Lang: "ru", PresetID: "gigaam-v3",
		Files:   []string{"encoder.int8.onnx", "decoder.onnx", "joiner.onnx", "tokens.txt"},
		BaseURL: "https://huggingface.co/csukuangfj/sherpa-onnx-nemo-transducer-punct-giga-am-v3-russian-2025-12-16/resolve/main/"},
	{ID: "base", Name: "Base", File: "ggml-base.bin", SizeMB: 142, BaseURL: whisperBaseURL, PresetID: "base"},
	{ID: "small", Name: "Small", File: "ggml-small.bin", SizeMB: 466, BaseURL: whisperBaseURL, PresetID: "small"},
	{ID: "medium", Name: "Medium (q5)", File: "ggml-medium-q5_0.bin", SizeMB: 539, BaseURL: whisperBaseURL, PresetID: "medium-q5_0"},
	{ID: "turbo", Name: "Turbo (q5)", File: "ggml-large-v3-turbo-q5_0.bin", SizeMB: 574, BaseURL: whisperBaseURL, PresetID: "large-v3-turbo-q5_0"},
}

func modelByID(id string) *modelOpt {
	for i := range modelOpts {
		if modelOpts[i].ID == id {
			return &modelOpts[i]
		}
	}
	return nil
}

func makeCmd(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	return cmd
}

func defaultInstallDir() string {
	base := os.Getenv("LOCALAPPDATA")
	if base == "" {
		base = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
	}
	return filepath.Join(base, "Programs", appid.InstallDirName)
}

func startMenuLnk() string {
	return filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", appid.ShortcutName)
}

func extractFile(f *zip.File, dir string) error {
	dst := filepath.Join(dir, f.Name)
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	var out *os.File
	for i := 0; ; i++ {
		out, err = os.Create(dst)
		if err == nil {
			break
		}
		if i >= 9 {
			return err
		}
		time.Sleep(300 * time.Millisecond)
	}
	defer out.Close()
	_, err = io.Copy(out, rc)
	return err
}

func makeShortcut(lnk, target, workDir string) error {
	ps := fmt.Sprintf(
		`$s=(New-Object -ComObject WScript.Shell).CreateShortcut('%s');$s.TargetPath='%s';$s.WorkingDirectory='%s';$s.Save()`,
		strings.ReplaceAll(lnk, "'", "''"),
		strings.ReplaceAll(target, "'", "''"),
		strings.ReplaceAll(workDir, "'", "''"))
	return runHidden("powershell", "-NoProfile", "-NonInteractive", "-Command", ps)
}

func dirSizeKB(dir string) uint32 {
	var total int64
	_ = filepath.Walk(dir, func(_ string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			total += info.Size()
		}
		return nil
	})
	return uint32(total / 1024)
}

func writeUninstallEntry(dir string) error {
	k, _, err := registry.CreateKey(registry.CURRENT_USER, uninstKey, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer k.Close()
	exe := filepath.Join(dir, exeName)
	_ = k.SetStringValue("DisplayName", appName)
	_ = k.SetStringValue("DisplayVersion", appVersion)
	_ = k.SetStringValue("DisplayIcon", exe)
	_ = k.SetStringValue("InstallLocation", dir)
	_ = k.SetStringValue("UninstallString", `"`+exe+`" -uninstall`)
	_ = k.SetDWordValue("NoModify", 1)
	_ = k.SetDWordValue("NoRepair", 1)
	_ = k.SetDWordValue("EstimatedSize", dirSizeKB(dir))
	return nil
}

func setAutorun(dir string, enable bool) {
	k, _, err := registry.CreateKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		return
	}
	defer k.Close()
	if enable {
		_ = k.SetStringValue(runValue, `"`+filepath.Join(dir, exeName)+`"`)
	} else {
		_ = k.DeleteValue(runValue)
	}
}

func patchDefaultConfig(dir string, m *modelOpt, updates bool) error {
	path := filepath.Join(dir, "config.default.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var cfg map[string]any
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	cfg["check_updates"] = updates
	if m != nil {
		presets := map[string]any{}
		if prev, ok := cfg["lang_models"].(map[string]any); ok {
			presets = prev
		}
		if m.Dir != "" {
			cfg["sherpa_model"] = "models/" + m.Dir
			if m.Lang != "" {
				cfg["language"] = m.Lang
				presets[m.Lang] = m.PresetID
			}
		} else {
			cfg["model"] = "models/" + m.File
			presets["auto"] = m.PresetID
		}
		cfg["lang_models"] = presets
	}
	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0o644)
}

var dlCancel atomic.Bool

func cancelDownload() { dlCancel.Store(true) }

var errDownloadCancelled = errors.New("cancelled")

func downloadModel(dir string, m *modelOpt, progress func(pct int)) error {
	if m.Dir != "" {
		return downloadModelDir(dir, m, progress)
	}
	return downloadFile(filepath.Join(dir, "models"), m.BaseURL+m.File, m.File, progress)
}

func downloadModelDir(dir string, m *modelOpt, progress func(pct int)) error {
	target := filepath.Join(dir, "models", m.Dir)
	if err := os.MkdirAll(target, 0o755); err != nil {
		return err
	}
	for i, f := range m.Files {
		base, span := i*100/len(m.Files), 100/len(m.Files)
		if err := downloadFile(target, m.BaseURL+f, f, func(pct int) {
			progress(base + span*pct/100)
		}); err != nil {
			return err
		}
	}
	progress(100)
	return nil
}

func downloadFile(dir, url, name string, progress func(pct int)) error {
	client := &http.Client{Timeout: 60 * time.Minute}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	total := resp.ContentLength
	part := filepath.Join(dir, name+".part")
	out, err := os.Create(part)
	if err != nil {
		return err
	}
	buf := make([]byte, 256*1024)
	var done int64
	for {
		if dlCancel.Load() {
			out.Close()
			os.Remove(part)
			return errDownloadCancelled
		}
		n, rerr := resp.Body.Read(buf)
		if n > 0 {
			if _, werr := out.Write(buf[:n]); werr != nil {
				out.Close()
				os.Remove(part)
				return werr
			}
			done += int64(n)
			if total > 0 {
				progress(int(done * 100 / total))
			}
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			out.Close()
			os.Remove(part)
			return rerr
		}
	}
	if err := out.Close(); err != nil {
		os.Remove(part)
		return err
	}
	return os.Rename(part, filepath.Join(dir, name))
}

func install(dir string, shortcut, autorun, touchAutorun, updates bool, modelID string, progress func(pct int, name string)) (string, error) {
	dlCancel.Store(false)
	if strings.TrimSpace(dir) == "" {
		return "", fmt.Errorf("empty path")
	}
	model := modelByID(modelID)
	filesSpan := 100
	if model != nil {
		filesSpan = 25
	}
	_ = runHidden("taskkill", "/IM", exeName, "/F")
	time.Sleep(400 * time.Millisecond)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Join(dir, "models"), 0o755); err != nil {
		return "", err
	}
	zr, err := zip.NewReader(bytes.NewReader(payload), int64(len(payload)))
	if err != nil {
		return "", err
	}
	for i, f := range zr.File {
		progress(i*filesSpan/len(zr.File), f.Name)
		if err := extractFile(f, dir); err != nil {
			return "", fmt.Errorf("%s: %w", f.Name, err)
		}
	}
	warn := ""
	if err := patchDefaultConfig(dir, model, updates); err != nil {
		warn = err.Error()
	} else if model != nil {
		name := model.Name
		progress(filesSpan, name)
		if derr := downloadModel(dir, model, func(pct int) {
			progress(filesSpan+(100-filesSpan-4)*pct/100, name)
		}); derr != nil {
			if errors.Is(derr, errDownloadCancelled) {
				warn = "cancelled"
			} else {
				warn = derr.Error()
			}
		}
	}
	progress(96, "")
	if shortcut {
		if err := makeShortcut(startMenuLnk(), filepath.Join(dir, exeName), dir); err != nil {
			return warn, fmt.Errorf("shortcut: %w", err)
		}
	}
	if touchAutorun {
		setAutorun(dir, autorun)
	}
	if err := writeUninstallEntry(dir); err != nil {
		return warn, fmt.Errorf("registry: %w", err)
	}
	progress(100, "")
	return warn, nil
}

// installedLook reads the skin and colour of the copy being updated, so the
// installer wears the same clothes.
func installedLook(dir string) (skin, colour string) {
	skin, colour = theme.DefaultSkin, theme.DefaultPalette
	if dir == "" {
		return skin, colour
	}
	data, err := os.ReadFile(filepath.Join(dir, "config.json"))
	if err != nil {
		return skin, colour
	}
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})
	var cfg struct {
		Skin  string `json:"skin"`
		Theme string `json:"theme"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return skin, colour
	}
	if cfg.Skin == "" && cfg.Theme != "" && !theme.ValidColour(cfg.Theme) {
		return theme.Migrate(cfg.Theme)
	}
	if theme.ValidSkin(cfg.Skin) {
		skin = cfg.Skin
	}
	if theme.ValidColour(cfg.Theme) {
		colour = cfg.Theme
	}
	return skin, colour
}
