package main

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/gen2brain/malgo"
)

const sampleRate = 16000

const maxBufferBytes = sampleRate * 2 * 600

type Recorder struct {
	ctx    *malgo.AllocatedContext
	device *malgo.Device

	mu        sync.Mutex
	buf       []byte
	recording bool
	started   bool

	level atomic.Uint32
}

func (r *Recorder) Level() float64 {
	return float64(r.level.Load()) / 1000
}

func NewRecorder() (*Recorder, error) {
	r := &Recorder{}
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, fmt.Errorf("аудиоконтекст: %w", err)
	}
	r.ctx = ctx

	cfg := malgo.DefaultDeviceConfig(malgo.Capture)
	cfg.Capture.Format = malgo.FormatS16
	cfg.Capture.Channels = 1
	cfg.SampleRate = sampleRate

	callbacks := malgo.DeviceCallbacks{
		Data: func(_, in []byte, _ uint32) {
			var peak int32
			for i := 0; i+1 < len(in); i += 2 {
				v := int32(int16(uint16(in[i]) | uint16(in[i+1])<<8))
				if v < 0 {
					v = -v
				}
				if v > peak {
					peak = v
				}
			}
			r.level.Store(uint32(peak * 1000 / 32768))
			r.mu.Lock()
			if r.recording && len(r.buf) < maxBufferBytes {
				r.buf = append(r.buf, in...)
			}
			r.mu.Unlock()
		},
	}
	device, err := malgo.InitDevice(ctx.Context, cfg, callbacks)
	if err != nil {
		ctx.Uninit()
		ctx.Free()
		return nil, fmt.Errorf("микрофон: %w", err)
	}
	r.device = device
	return r, nil
}

func (r *Recorder) Start() error {
	r.mu.Lock()
	if r.started {
		r.mu.Unlock()
		return errors.New("запись уже идёт")
	}
	r.buf = nil
	r.recording = true
	r.started = true
	r.mu.Unlock()
	if err := r.device.Start(); err != nil {
		r.mu.Lock()
		r.recording = false
		r.started = false
		r.mu.Unlock()
		return err
	}
	return nil
}

func (r *Recorder) Stop() []byte {
	r.mu.Lock()
	wasStarted := r.started
	r.recording = false
	r.started = false
	r.mu.Unlock()
	if wasStarted {
		_ = r.device.Stop()
	}
	r.mu.Lock()
	pcm := r.buf
	r.buf = nil
	r.mu.Unlock()
	return pcm
}

func (r *Recorder) Close() {
	if r.device != nil {
		r.device.Uninit()
	}
	if r.ctx != nil {
		r.ctx.Uninit()
		r.ctx.Free()
	}
}
