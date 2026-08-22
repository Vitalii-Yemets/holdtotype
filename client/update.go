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
	"time"

	"holdtotype/internal/appver"

	"holdtotype/internal/appid"
)

const updateLatestURL = appid.LatestAPI

func verNewer(latest, current string) bool {
	return appver.IsNewer(latest, current)
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
		if a.Name == appid.SetupExe {
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
	dst := filepath.Join(os.TempDir(), appid.Slug+"-setup-update.exe")
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
