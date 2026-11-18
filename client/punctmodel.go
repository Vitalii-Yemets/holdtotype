package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"holdtotype/internal/punct"
	"holdtotype/internal/trmodel"
)

var punctModel = punctCatalog()

func punctCatalog() modelInfo {
	m := modelInfo{ID: "punct-horizon", SizeMB: punct.SizeMB(), NameKey: punct.ModelName, Dir: punct.ModelDir,
		BaseURL: punct.BaseURL(), Remote: map[string]string{}, Hashes: map[string]string{}}
	for _, f := range punct.Files {
		m.Files = append(m.Files, f.Name)
		m.Remote[f.Name] = f.Remote
		m.Hashes[f.Name] = f.SHA256
	}
	return m
}

var (
	punctLive     atomic.Int32
	punctDeleting atomic.Bool
	punctTurnMu   sync.Mutex
	punctTurnOn   bool
)

type punctState struct {
	Name    string `json:"name"`
	State   string `json:"state"`
	Pct     int    `json:"pct"`
	Err     string `json:"err"`
	Size    int    `json:"size"`
	RAM     int    `json:"ram"`
	Langs   int    `json:"langs"`
	On      bool   `json:"on"`
	Running bool   `json:"running"`
}

func punctInstalled() bool {
	return !punctDeleting.Load() && punctModel.installed()
}

func punctStateNow(cfg *Config) punctState {
	st := punctState{Name: punctModel.NameKey, Size: punctModel.SizeMB, RAM: punct.RAMMB, Langs: punct.LangCount,
		On: cfg.SubtitlesSentences, Running: punctLive.Load() > 0}
	dlMu.Lock()
	d := dl[punctModel.ID]
	dlMu.Unlock()
	switch {
	case d != nil && d.active:
		st.State, st.Pct = "downloading", d.pct
	case punctInstalled():
		st.State = "ready"
	default:
		st.State = "absent"
		if d != nil {
			st.Err = d.err
		}
	}
	return st
}

func (a *App) punctAction(action string) string {
	switch action {
	case "download":
		a.punctDownload(false)
	case "download:on":
		a.punctDownload(true)
	case "cancel":
		punctTurnMu.Lock()
		punctTurnOn = false
		punctTurnMu.Unlock()
		cancelDownload(punctModel.ID)
	case "delete":
		return a.punctDelete()
	case "terms":
		shellOpenURL(trmodel.ApacheTermsURL)
	}
	return ""
}

func (a *App) punctDownload(turnOn bool) {
	if punctInstalled() {
		if turnOn {
			a.punctSwitchOn()
			a.subtitlesConfigChanged()
		}
		return
	}
	punctTurnMu.Lock()
	punctTurnOn = punctTurnOn || turnOn
	punctTurnMu.Unlock()
	log.Printf("punctuation model %s: download started", punctModel.NameKey)
	a.startMultiDownloadThen(punctModel.ID, &punctModel, func() {
		punctTurnMu.Lock()
		on := punctTurnOn
		punctTurnOn = false
		punctTurnMu.Unlock()
		if on {
			a.punctSwitchOn()
		}
		a.subtitlesConfigChanged()
	})
}

func (a *App) punctSwitchOn() {
	a.mu.Lock()
	if a.cfg.SubtitlesSentences {
		a.mu.Unlock()
		return
	}
	c := *a.cfg
	c.SubtitlesSentences = true
	if err := saveConfig("config.json", &c); err != nil {
		a.mu.Unlock()
		log.Printf("turning on whole sentences after the download: %v", err)
		return
	}
	a.cfg = &c
	a.mu.Unlock()
	log.Printf("whole sentences are turned on: the punctuation model is downloaded")
}

func (a *App) punctDelete() string {
	dlMu.Lock()
	if d := dl[punctModel.ID]; d != nil && d.active {
		dlMu.Unlock()
		return ""
	}
	dlMu.Unlock()
	punctDeleting.Store(true)
	defer punctDeleting.Store(false)
	a.dropSentences("the punctuation model is being deleted")
	dir := filepath.Join("models", punctModel.Dir)
	var err error
	for i := 0; i < 30; i++ {
		if err = os.RemoveAll(dir); err == nil {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	if err != nil {
		log.Printf("deleting the punctuation model: %v", err)
		return humanError(err)
	}
	_ = os.Remove(filepath.Dir(dir))
	dlMu.Lock()
	delete(dl, punctModel.ID)
	dlMu.Unlock()
	log.Printf("punctuation model %s deleted", punctModel.NameKey)
	return tr("model.del.ok")
}

func (a *App) dropSentences(why string) {
	a.subMu.Lock()
	r := a.sub
	a.subMu.Unlock()
	if r == nil {
		return
	}
	r.phMu.Lock()
	defer r.phMu.Unlock()
	s := r.sent
	if s == nil {
		return
	}
	if s.active() {
		s.halt(why)
	}
	if given, seq, ok := s.takeOver(); ok {
		r.phraser.Resume(given, seq)
	}
	r.sent = nil
}

func runPunctDownload(args []string) {
	cfg, err := loadConfig("config.json")
	if err != nil {
		log.Printf("punctdownload: config: %v", err)
		return
	}
	turnOn := len(args) > 0 && args[0] == "on"
	app := &App{cfg: cfg, enabled: true}
	started := time.Now()
	app.punctDownload(turnOn)
	last := -1
	for {
		time.Sleep(200 * time.Millisecond)
		dlMu.Lock()
		d := dl[punctModel.ID]
		active := d != nil && d.active
		pct := 0
		if d != nil {
			pct = d.pct
		}
		dlMu.Unlock()
		if pct/10 != last/10 {
			log.Printf("punctdownload: %d%%", pct)
			last = pct
		}
		if !active {
			break
		}
	}
	time.Sleep(300 * time.Millisecond)
	st := punctStateNow(app.snapshot())
	saved, _ := loadConfig("config.json")
	log.Printf("punctdownload: state %s, error %q, %.0f s, whole sentences in config.json: %v", st.State, st.Err, time.Since(started).Seconds(), saved != nil && saved.SubtitlesSentences)
	if st.State == "ready" {
		log.Printf("punctdownload: model check: %s", app.verifyModels())
	}
}
