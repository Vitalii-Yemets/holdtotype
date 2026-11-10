package main

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gen2brain/malgo"
)

const loopbackKeepSeconds = 30

type loopback struct {
	ctx *malgo.AllocatedContext

	devMu sync.Mutex
	mu    sync.Mutex
	dev   *malgo.Device
	buf   []byte
	max   int

	dropped  atomic.Int64
	lastData atomic.Int64
	opens    atomic.Int32
}

func openLoopback() (*loopback, error) {
	ctx, err := malgo.InitContext([]malgo.Backend{malgo.BackendWasapi}, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, fmt.Errorf("audio context: %w", err)
	}
	l := &loopback{ctx: ctx, max: sampleRate * 2 * loopbackKeepSeconds}
	if err := l.open(); err != nil {
		ctx.Uninit()
		ctx.Free()
		return nil, err
	}
	return l, nil
}

func (l *loopback) onData(_, in []byte, _ uint32) {
	l.lastData.Store(time.Now().UnixNano())
	l.mu.Lock()
	l.buf = append(l.buf, in...)
	if over := len(l.buf) - l.max; over > 0 {
		n := copy(l.buf, l.buf[over:])
		l.buf = l.buf[:n]
		l.dropped.Add(int64(over))
	}
	l.mu.Unlock()
}

func (l *loopback) open() error {
	l.devMu.Lock()
	defer l.devMu.Unlock()
	cfg := malgo.DefaultDeviceConfig(malgo.Loopback)
	cfg.Capture.Format = malgo.FormatS16
	cfg.Capture.Channels = 1
	cfg.SampleRate = sampleRate
	dev, err := malgo.InitDevice(l.ctx.Context, cfg, malgo.DeviceCallbacks{Data: l.onData})
	if err != nil {
		return fmt.Errorf("the playback device cannot be listened to: %w", err)
	}
	if err := dev.Start(); err != nil {
		dev.Uninit()
		return fmt.Errorf("listening to the playback device did not start: %w", err)
	}
	l.mu.Lock()
	old := l.dev
	l.dev = dev
	l.mu.Unlock()
	if old != nil {
		old.Uninit()
	}
	l.opens.Add(1)
	return nil
}

func (l *loopback) reopen() error { return l.open() }

func (l *loopback) running() bool {
	l.mu.Lock()
	dev := l.dev
	l.mu.Unlock()
	return dev != nil && dev.IsStarted()
}

func (l *loopback) take() []byte {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.buf) == 0 {
		return nil
	}
	out := l.buf
	l.buf = make([]byte, 0, len(out))
	return out
}

func (l *loopback) close() {
	l.devMu.Lock()
	defer l.devMu.Unlock()
	l.mu.Lock()
	dev := l.dev
	l.dev = nil
	l.buf = nil
	l.mu.Unlock()
	if dev != nil {
		_ = dev.Stop()
		dev.Uninit()
	}
	if l.ctx != nil {
		l.ctx.Uninit()
		l.ctx.Free()
		l.ctx = nil
	}
}

var errNoPlayback = errors.New("no playback device to listen to")
