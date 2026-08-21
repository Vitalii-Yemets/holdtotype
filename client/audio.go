package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"github.com/gen2brain/malgo"

	"voxterminal/internal/audiolevel"
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
	deviceID  string

	level atomic.Uint32
	peak  atomic.Uint32
}

type micDevice struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Default bool   `json:"default"`
}

func (r *Recorder) Level() float64 {
	return float64(r.level.Load()) / 1000
}

func (r *Recorder) resetPeak() {
	r.peak.Store(0)
}

func (r *Recorder) PeakLevel() float64 {
	return float64(r.peak.Load()) / 1000
}

func (r *Recorder) devices() []micDevice {
	r.mu.Lock()
	ctx := r.ctx
	cur := r.deviceID
	r.mu.Unlock()
	if ctx == nil {
		return nil
	}
	infos, err := ctx.Devices(malgo.Capture)
	if err != nil {
		log.Printf("список микрофонов: %v", err)
		return nil
	}
	out := make([]micDevice, 0, len(infos))
	for _, info := range infos {
		id := hex.EncodeToString(info.ID[:])
		out = append(out, micDevice{
			ID:      id,
			Name:    info.Name(),
			Default: id == cur,
		})
	}
	return out
}

func decodeDeviceID(id string) *malgo.DeviceID {
	raw, err := hex.DecodeString(id)
	if err != nil || len(raw) == 0 {
		return nil
	}
	var did malgo.DeviceID
	if len(raw) > len(did) {
		raw = raw[:len(did)]
	}
	copy(did[:], raw)
	return &did
}

func (r *Recorder) openDevice(deviceID string) error {
	cfg := malgo.DefaultDeviceConfig(malgo.Capture)
	cfg.Capture.Format = malgo.FormatS16
	cfg.Capture.Channels = 1
	cfg.SampleRate = sampleRate
	if did := decodeDeviceID(deviceID); did != nil {
		cfg.Capture.DeviceID = did.Pointer()
	}

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
			lvl := uint32(peak * 1000 / 32768)
			r.level.Store(lvl)
			if lvl > r.peak.Load() {
				r.peak.Store(lvl)
			}
			r.mu.Lock()
			if r.recording && len(r.buf) < maxBufferBytes {
				r.buf = append(r.buf, in...)
			}
			r.mu.Unlock()
		},
	}
	device, err := malgo.InitDevice(r.ctx.Context, cfg, callbacks)
	if err != nil {
		return err
	}
	r.mu.Lock()
	old := r.device
	r.device = device
	r.deviceID = deviceID
	r.started = false
	r.mu.Unlock()
	if old != nil {
		old.Uninit()
	}
	return nil
}

func NewRecorder(deviceID string) (*Recorder, error) {
	r := &Recorder{}
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, fmt.Errorf("аудиоконтекст: %w", err)
	}
	r.ctx = ctx

	if err := r.openDevice(deviceID); err != nil {
		if deviceID != "" {
			log.Printf("микрофон %s недоступен (%v) — беру системный по умолчанию", deviceID, err)
			if err2 := r.openDevice(""); err2 == nil {
				return r, nil
			}
		}
		ctx.Uninit()
		ctx.Free()
		return nil, fmt.Errorf("микрофон: %w", err)
	}
	return r, nil
}

func (r *Recorder) SetDevice(deviceID string) error {
	r.mu.Lock()
	same := r.deviceID == deviceID
	busy := r.started
	r.mu.Unlock()
	if same || busy {
		return nil
	}
	return r.openDevice(deviceID)
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
	dev := r.device
	devID := r.deviceID
	r.mu.Unlock()
	r.resetPeak()

	err := dev.Start()
	if err != nil {
		log.Printf("микрофон не запустился (%v) — переоткрываю устройство", err)
		if rerr := r.openDevice(devID); rerr != nil && devID != "" {
			rerr = r.openDevice("")
			if rerr != nil {
				r.failStart()
				return err
			}
		} else if rerr != nil {
			r.failStart()
			return err
		}
		r.mu.Lock()
		dev = r.device
		r.recording = true
		r.started = true
		r.mu.Unlock()
		if err2 := dev.Start(); err2 != nil {
			r.failStart()
			return err2
		}
	}
	return nil
}

func (r *Recorder) failStart() {
	r.mu.Lock()
	r.recording = false
	r.started = false
	r.mu.Unlock()
}

func (r *Recorder) Stop() []byte {
	r.mu.Lock()
	wasStarted := r.started
	dev := r.device
	r.recording = false
	r.started = false
	r.mu.Unlock()
	if wasStarted && dev != nil {
		_ = dev.Stop()
	}
	r.mu.Lock()
	pcm := r.buf
	r.buf = nil
	r.mu.Unlock()
	return pcm
}

func (r *Recorder) Close() {
	r.mu.Lock()
	dev := r.device
	r.device = nil
	r.mu.Unlock()
	if dev != nil {
		dev.Uninit()
	}
	if r.ctx != nil {
		r.ctx.Uninit()
		r.ctx.Free()
	}
}

func pcmPeak(pcm []byte) float64 { return audiolevel.Peak(pcm) }

func pcmIsSilent(pcm []byte) bool { return audiolevel.IsSilent(pcm) }
