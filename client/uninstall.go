package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"golang.org/x/sys/windows/registry"
)

const (
	uninstRegKey = `Software\Microsoft\Windows\CurrentVersion\Uninstall\VoxTerminal`
	runRegKey    = `Software\Microsoft\Windows\CurrentVersion\Run`
	runRegValue  = "VoxTerminal"
)

func hiddenCmd(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	return cmd
}

func runUninstall(silent bool) {
	if !silent && !msgBoxYesNo(tr("un.title"), tr("un.confirm")) {
		return
	}
	exe, err := os.Executable()
	if err != nil {
		return
	}
	dir := filepath.Dir(exe)

	_ = hiddenCmd("taskkill", "/F", "/IM", "voxterminal.exe", "/FI", fmt.Sprintf("PID ne %d", os.Getpid())).Run()
	_ = hiddenCmd("taskkill", "/F", "/IM", "whisper-server.exe").Run()
	_ = hiddenCmd("taskkill", "/F", "/IM", "llama-server.exe").Run()
	time.Sleep(500 * time.Millisecond)

	delData := false
	if !silent {
		delData = msgBoxYesNo(tr("un.title"), tr("un.data"))
	}
	files := []string{"whisper-server.exe", "llama-server.exe", "config.default.json", "README.md"}
	if delData {
		files = append(files, "config.json", "voxterminal.log", "voxterminal.log.old")
		_ = os.RemoveAll(filepath.Join(dir, "models"))
	}
	for _, f := range files {
		_ = os.Remove(filepath.Join(dir, f))
	}
	_ = os.Remove(filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", "Vox Terminal.lnk"))

	_ = registry.DeleteKey(registry.CURRENT_USER, uninstRegKey)
	if k, kerr := registry.OpenKey(registry.CURRENT_USER, runRegKey, registry.SET_VALUE); kerr == nil {
		_ = k.DeleteValue(runRegValue)
		k.Close()
	}

	if !silent {
		msgBox(tr("un.title"), tr("un.done"))
	}
	_ = hiddenCmd("cmd", "/C",
		`ping -n 3 127.0.0.1 >nul & del /f /q "`+exe+`" & rd "`+dir+`"`).Start()
}
