package uninst

import "holdtotype/internal/appid"

var programFiles = []string{"whisper-server.exe", "llama-server.exe", "llama-server-gpu.exe", "sherpa-server.exe", "sherpa-online-server.exe", "config.default.json", "README.md"}

func Plan(delData bool) (files, dirs []string) {
	files = append(files, programFiles...)
	if !delData {
		return files, nil
	}
	files = append(files, "config.json", appid.LogFile, appid.LogFile+".old", appid.HistoryFile)
	dirs = []string{"models", "transcripts"}
	return files, dirs
}
