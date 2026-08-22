package main

import (
	"context"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"holdtotype/internal/routing"
)

const (
	engineWhisper = "whisper"
	engineSherpa  = "sherpa"
)

type recognizer interface {
	transcribe(ctx context.Context, wav []byte, language, prompt string, translate bool) (string, error)
	waitReady(timeout time.Duration) error
	stop()
	done() <-chan struct{}
	external() bool
	wasStopped() bool
	engine() string
}

func dialOK(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func validEngine(name string) bool { return routing.ValidMode(name) }

func startEngine(cfg *Config, engine string, logw io.Writer) (recognizer, error) {
	if engine == engineSherpa {
		return startSherpaServer(cfg, logw)
	}
	return startWhisperServer(cfg, logw)
}

func startRecognizer(cfg *Config, logw io.Writer) (recognizer, error) {
	return startEngine(cfg, primaryEngine(cfg), logw)
}

func sherpaInstalled(cfg *Config) bool {
	_, _, _, _, err := sherpaModelFiles(cfg.SherpaModel)
	return err == nil
}

func sherpaLangs(cfg *Config) []string {
	dir := filepath.Base(filepath.Clean(cfg.SherpaModel))
	for i := range modelCatalog {
		m := &modelCatalog[i]
		if m.Engine == engineSherpa && m.Dir == dir && m.Langs != "" {
			return strings.Split(m.Langs, ",")
		}
	}
	return []string{"ru"}
}

func routingInput(cfg *Config, translate bool) routing.Input {
	return routing.Input{
		Mode:        cfg.STTEngine,
		Language:    cfg.Language,
		Translate:   translate,
		SherpaReady: sherpaInstalled(cfg),
		SherpaLangs: sherpaLangs(cfg),
	}
}

func pickEngine(cfg *Config, translate bool) routing.Decision {
	return routing.Pick(routingInput(cfg, translate))
}

func primaryEngine(cfg *Config) string {
	return pickEngine(cfg, false).Engine
}

func engineModelName(cfg *Config, engine string) string {
	if engine == engineSherpa {
		return filepath.Base(filepath.Clean(cfg.SherpaModel))
	}
	return filepath.Base(cfg.Model)
}

type routeRow struct {
	Cond   string `json:"cond"`
	Engine string `json:"engine"`
	Why    string `json:"why"`
}

func routeRows(cfg *Config) []routeRow {
	d := pickEngine(cfg, false)
	rows := []routeRow{{
		Cond:   trf("route.speech", routeLangLabel(cfg.Language)),
		Engine: engineModelName(cfg, d.Engine),
		Why:    tr("route.why." + d.Reason),
	}}
	if d.Engine == engineSherpa {
		rows = append(rows, routeRow{
			Cond:   tr("route.other"),
			Engine: engineModelName(cfg, engineWhisper),
			Why:    tr("route.why.otherlang"),
		})
	}
	rows = append(rows, routeRow{
		Cond:   tr("route.translate"),
		Engine: engineModelName(cfg, engineWhisper),
		Why:    tr("route.why.translate"),
	})
	return rows
}

func routeLangLabel(code string) string {
	code = strings.ToLower(strings.TrimSpace(code))
	if code == "" || code == "auto" {
		return tr("route.lang.auto")
	}
	return langLabel(code)
}

func engineTranslates(name string) bool { return name == engineWhisper }

func activeModelPath(cfg *Config) string {
	if primaryEngine(cfg) == engineSherpa {
		return cfg.SherpaModel
	}
	return cfg.Model
}

func missingModelPath(cfg *Config) string {
	if primaryEngine(cfg) == engineSherpa {
		if _, _, _, _, err := sherpaModelFiles(cfg.SherpaModel); err != nil {
			return cfg.SherpaModel
		}
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
