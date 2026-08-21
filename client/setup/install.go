package main

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/windows/registry"
)

//go:embed payload.zip
var payload []byte

const (
	appName   = "Vox Terminal"
	exeName   = "voxterminal.exe"
	uninstKey = `Software\Microsoft\Windows\CurrentVersion\Uninstall\VoxTerminal`
	runKey    = `Software\Microsoft\Windows\CurrentVersion\Run`
	runValue  = "VoxTerminal"
)

var appVersion = "0.7.3"

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
	ID     string
	Name   string
	File   string
	SizeMB int
}

var modelOpts = []modelOpt{
	{"base", "Base", "ggml-base.bin", 142},
	{"small", "Small", "ggml-small.bin", 466},
	{"medium", "Medium (q5)", "ggml-medium-q5_0.bin", 539},
	{"turbo", "Turbo (q5)", "ggml-large-v3-turbo-q5_0.bin", 574},
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
	return filepath.Join(base, "Programs", "VoxTerminal")
}

func startMenuLnk() string {
	return filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", appName+".lnk")
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

func patchDefaultModel(dir, modelFile string) error {
	path := filepath.Join(dir, "config.default.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var cfg map[string]any
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	cfg["model"] = "models/" + modelFile
	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0o644)
}

func downloadModel(dir string, m *modelOpt, progress func(pct int)) error {
	client := &http.Client{Timeout: 60 * time.Minute}
	resp, err := client.Get("https://huggingface.co/ggerganov/whisper.cpp/resolve/main/" + m.File)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	total := resp.ContentLength
	part := filepath.Join(dir, "models", m.File+".part")
	out, err := os.Create(part)
	if err != nil {
		return err
	}
	buf := make([]byte, 256*1024)
	var done int64
	for {
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
	return os.Rename(part, filepath.Join(dir, "models", m.File))
}

func install(dir string, shortcut, autorun, touchAutorun bool, modelID string, progress func(pct int, name string)) (string, error) {
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
	if model != nil {
		if err := patchDefaultModel(dir, model.File); err == nil {
			progress(filesSpan, model.File)
			if derr := downloadModel(dir, model, func(pct int) {
				progress(filesSpan+(100-filesSpan-4)*pct/100, model.File)
			}); derr != nil {
				warn = derr.Error()
			}
		} else {
			warn = err.Error()
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
