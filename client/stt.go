package main

import (
	"context"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"holdtotype/internal/preset"
)

const portFreeWait = 2 * time.Second

const (
	engineWhisper = "whisper"
	engineSherpa  = "sherpa"
	engineStream  = "stream"
)

type recognizer interface {
	transcribe(ctx context.Context, wav []byte, language, prompt string, translate bool) (string, error)
	waitReady(timeout time.Duration) error
	stop()
	done() <-chan struct{}
	external() bool
	wasStopped() bool
	engine() string
	model() string
}

func portFree(port int) bool {
	ln, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(port)))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

func waitPortFree(port int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for {
		if portFree(port) {
			return true
		}
		if time.Now().After(deadline) {
			return false
		}
		time.Sleep(150 * time.Millisecond)
	}
}

func dialOK(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func startEngine(cfg *Config, engine string, logw io.Writer) (recognizer, error) {
	switch engine {
	case engineSherpa:
		return startSherpaServer(cfg, logw)
	case engineStream:
		return startStreamServer(cfg, logw)
	}
	return startWhisperServer(cfg, logw)
}

func startRecognizer(cfg *Config, logw io.Writer) (recognizer, error) {
	return startEngine(cfg, primaryEngine(cfg), logw)
}

func sherpaInstalled(cfg *Config) bool {
	_, err := sherpaModelArgs(cfg.SherpaModel)
	return err == nil
}

func presetView(m *modelInfo) *preset.Model {
	if m == nil {
		return nil
	}
	return &preset.Model{
		ID: m.ID, Engine: m.Engine, Langs: strings.Split(m.Langs, ","),
		Auto: m.Auto, Translate: m.Translate,
	}
}

func presetModelID(cfg *Config, lang string) string {
	return preset.Resolve(cfg.LangModels, lang, defaultPresetModel,
		func(id string) bool { return findModel(id) != nil })
}

func modelForLang(cfg *Config, lang string) *modelInfo {
	return findModel(presetModelID(cfg, lang))
}

func activeModel(cfg *Config) *modelInfo {
	return modelForLang(cfg, cfg.Language)
}

func primaryEngine(cfg *Config) string {
	if m := activeModel(cfg); m != nil {
		return m.Engine
	}
	return engineWhisper
}

func modelTranslates(m *modelInfo, srcLang, target string) bool {
	if m == nil || !m.Translate {
		return false
	}
	if m.TrLangs == "" {
		return true
	}
	return preset.Covers(strings.Split(m.TrLangs, ","), target) &&
		preset.Covers(strings.Split(m.Langs, ","), srcLang)
}

func bestInstalledWhisper() *modelInfo {
	var best *modelInfo
	for i := range modelCatalog {
		m := &modelCatalog[i]
		if m.Engine != engineWhisper || !m.Translate || m.Custom || !m.installed() {
			continue
		}
		if best == nil || m.Accuracy > best.Accuracy {
			best = m
		}
	}
	return best
}

// applyPreset writes the paths the engines actually load, derived from the
// language presets: the active model's own slot, and — when the active model
// is not a whisper — the best installed whisper standing by for translation.
func applyPreset(cfg *Config) bool {
	m := activeModel(cfg)
	if m == nil {
		return false
	}
	changed := false
	switch m.Engine {
	case engineSherpa:
		if nd := m.modelPath(); cfg.SherpaModel != nd {
			cfg.SherpaModel = nd
			changed = true
		}
	case engineStream:
		if nd := m.modelPath(); cfg.StreamModel != nd {
			cfg.StreamModel = nd
			changed = true
		}
	default:
		if nm := m.modelPath(); cfg.Model != nm {
			cfg.Model = nm
			changed = true
		}
	}
	if m.Engine != engineWhisper {
		if w := bestInstalledWhisper(); w != nil {
			if nm := w.modelPath(); cfg.Model != nm {
				cfg.Model = nm
				changed = true
			}
		}
	}
	return changed
}

func engineTranslates(name string) bool { return name == engineWhisper }

func activeModelPath(cfg *Config) string {
	switch primaryEngine(cfg) {
	case engineSherpa:
		return cfg.SherpaModel
	case engineStream:
		return cfg.StreamModel
	}
	return cfg.Model
}

func missingModelPath(cfg *Config) string {
	switch primaryEngine(cfg) {
	case engineSherpa:
		if _, err := sherpaModelArgs(cfg.SherpaModel); err != nil {
			return cfg.SherpaModel
		}
		return ""
	case engineStream:
		if _, err := streamModelArgs(cfg.StreamModel); err != nil {
			return cfg.StreamModel
		}
		return ""
	}
	if cfg.SttSource == "remote" {
		return ""
	}
	if _, err := os.Stat(cfg.Model); err != nil {
		return cfg.Model
	}
	return ""
}

func engineReadyTimeout(r recognizer) time.Duration {
	if r.external() {
		return 20 * time.Second
	}
	return 3 * time.Minute
}
