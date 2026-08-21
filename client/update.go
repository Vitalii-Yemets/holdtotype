package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const updateLatestURL = "https://api.github.com/repos/Vitalii-Yemets/vox-terminal/releases/latest"

func parseVer(s string) ([3]int, bool) {
	var v [3]int
	s = strings.TrimPrefix(strings.TrimSpace(s), "v")
	parts := strings.SplitN(s, ".", 3)
	if len(parts) != 3 {
		return v, false
	}
	for i, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return v, false
		}
		v[i] = n
	}
	return v, true
}

func verNewer(latest, current string) bool {
	l, ok1 := parseVer(latest)
	c, ok2 := parseVer(current)
	if !ok1 || !ok2 {
		return false
	}
	for i := 0; i < 3; i++ {
		if l[i] != c[i] {
			return l[i] > c[i]
		}
	}
	return false
}

func fetchLatestRelease() (tag, setupURL string, err error) {
	client := &http.Client{Timeout: 20 * time.Second}
	req, _ := http.NewRequest(http.MethodGet, updateLatestURL, nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	var rel struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name string `json:"name"`
			URL  string `json:"browser_download_url"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return "", "", err
	}
	for _, a := range rel.Assets {
		if a.Name == "voxterminal-setup.exe" {
			return rel.TagName, a.URL, nil
		}
	}
	return rel.TagName, "", nil
}

func (a *App) startupUpdateCheck() {
	if !a.snapshot().CheckUpdates {
		return
	}
	tag, url, err := fetchLatestRelease()
	if err != nil {
		log.Printf("проверка обновлений: %v", err)
		return
	}
	if verNewer(tag, appVersion) && url != "" {
		log.Printf("доступно обновление %s (текущая %s)", tag, appVersion)
		a.mu.Lock()
		a.updVer, a.updURL = tag, url
		a.mu.Unlock()
	}
}

func downloadSetup(url string, progress func(pct int)) (string, error) {
	client := &http.Client{Timeout: 30 * time.Minute}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	dst := filepath.Join(os.TempDir(), "voxterminal-setup-update.exe")
	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	total := resp.ContentLength
	buf := make([]byte, 256*1024)
	var done int64
	for {
		n, rerr := resp.Body.Read(buf)
		if n > 0 {
			if _, werr := out.Write(buf[:n]); werr != nil {
				out.Close()
				os.Remove(dst)
				return "", werr
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
			os.Remove(dst)
			return "", rerr
		}
	}
	if err := out.Close(); err != nil {
		os.Remove(dst)
		return "", err
	}
	return dst, nil
}

func launchUpdater(setupPath string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	cmd := exec.Command(setupPath, "-update", "-dir", filepath.Dir(exe))
	return cmd.Start()
}

func (a *App) quitForUpdate() {
	log.Printf("выхожу для обновления")
	a.onExit()
	trayMu.Lock()
	h := trayHwnd
	trayMu.Unlock()
	procPostMessageW.Call(h, wmClose, 0, 0)
}
