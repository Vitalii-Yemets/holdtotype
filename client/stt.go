package main

import (
	"context"
	"io"
	"net"
	"os"
	"time"
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

func validEngine(name string) bool {
	return name == engineWhisper || name == engineSherpa
}

func startRecognizer(cfg *Config, logw io.Writer) (recognizer, error) {
	if cfg.STTEngine == engineSherpa {
		return startSherpaServer(cfg, logw)
	}
	return startWhisperServer(cfg, logw)
}

func engineTranslates(name string) bool { return name == engineWhisper }

func activeModelPath(cfg *Config) string {
	if cfg.STTEngine == engineSherpa {
		return cfg.SherpaModel
	}
	return cfg.Model
}

func missingModelPath(cfg *Config) string {
	if cfg.STTEngine == engineSherpa {
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
