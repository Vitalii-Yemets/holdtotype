package main

import (
	"time"

	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"github.com/gen2brain/malgo"

	"holdtotype/internal/audiolevel"
)

const sampleRate = 16000

const defaultMaxSeconds = 600

const monitorLinger = 2 * time.Second
const micSilenceLimit = 1500 * time.Millisecond

type Recorder struct {
	ctx    *malgo.AllocatedContext
	device *malgo.Device

	devMu sync.Mutex

	mu        sync.Mutex
	buf       []byte
	maxBytes  int
	overflow  bool
	recording bool
	paused    bool
	started   bool
	monitor   bool
	monUntil  time.Time
	deviceID  string

	done chan struct{}
	once sync.Once

	level    atomic.Uint32
	peak     atomic.Uint32
	lastData atomic.Int64
	lost     atomic.Bool
	onLost   func()
}

type micDevice struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Default bool   `json:"default"`
	System  bool   `json:"system"`
}

func (r *Recorder) Level() float64 {
	if r.lost.Load() {
		return 0
	}
	return float64(r.level.Load()) / 1000
}

func (r *Recorder) Lost() bool { return r.lost.Load() }

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
	for i := range infos {
		id := hex.EncodeToString(infos[i].ID[:])
		out = append(out, micDevice{
			ID:      id,
			Name:    infos[i].Name(),
			Default: id == cur,
			System:  infos[i].IsDefault != 0,
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
			r.lastData.Store(time.Now().UnixNano())
			r.level.Store(lvl)
			if lvl > r.peak.Load() {
				r.peak.Store(lvl)
			}
			r.mu.Lock()
			if r.recording && !r.paused {
				if len(r.buf) < r.maxBytes {
					r.buf = append(r.buf, in...)
				} else {
					r.overflow = true
				}
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
	r := &Recorder{
		maxBytes: sampleRate * 2 * defaultMaxSeconds,
		done:     make(chan struct{}),
	}
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, fmt.Errorf("аудиоконтекст: %w", err)
	}
	r.ctx = ctx

	if err := r.openDevice(deviceID); err != nil {
		opened := false
		if deviceID != "" {
			log.Printf("микрофон %s недоступен (%v) — беру системный по умолчанию", deviceID, err)
			opened = r.openDevice("") == nil
		}
		if !opened {
			ctx.Uninit()
			ctx.Free()
			return nil, fmt.Errorf("микрофон: %w", err)
		}
	}
	go r.watchMonitor()
	return r, nil
}

func (r *Recorder) SetDevice(deviceID string) error {
	r.mu.Lock()
	same := r.deviceID == deviceID
	busy := r.recording
	wasMonitoring := r.monitor
	if !same && !busy {
		r.monitor = false
	}
	r.mu.Unlock()
	if same {
		return nil
	}
	if busy {
		return errRecorderBusy
	}
	r.stopDevice()
	err := r.openDevice(deviceID)
	if wasMonitoring {
		r.mu.Lock()
		r.monitor = true
		r.monUntil = time.Now().Add(monitorLinger)
		r.mu.Unlock()
		if err == nil {
			err = r.startDevice()
		}
	}
	return err
}

func (r *Recorder) startDevice() error {
	r.devMu.Lock()
	defer r.devMu.Unlock()
	r.mu.Lock()
	if r.started {
		r.mu.Unlock()
		return nil
	}
	dev := r.device
	devID := r.deviceID
	r.mu.Unlock()
	if dev == nil {
		return errors.New("устройство записи не открыто")
	}

	err := dev.Start()
	if err != nil {
		log.Printf("микрофон не запустился (%v) — переоткрываю устройство", err)
		rerr := r.openDevice(devID)
		if rerr != nil && devID != "" {
			rerr = r.openDevice("")
		}
		if rerr != nil {
			return err
		}
		r.mu.Lock()
		dev = r.device
		r.mu.Unlock()
		if err2 := dev.Start(); err2 != nil {
			return err2
		}
	}
	r.mu.Lock()
	r.started = true
	r.mu.Unlock()
	return nil
}

func (r *Recorder) stopDevice() {
	r.devMu.Lock()
	defer r.devMu.Unlock()
	r.mu.Lock()
	if !r.started || r.recording || r.monitor {
		r.mu.Unlock()
		return
	}
	dev := r.device
	r.started = false
	r.mu.Unlock()
	if dev != nil {
		_ = dev.Stop()
	}
	r.level.Store(0)
}

func (r *Recorder) Start(maxSeconds int) error {
	r.mu.Lock()
	if r.recording {
		r.mu.Unlock()
		return errors.New("запись уже идёт")
	}
	if maxSeconds <= 0 {
		maxSeconds = defaultMaxSeconds
	}
	r.buf = nil
	r.overflow = false
	r.maxBytes = sampleRate * 2 * maxSeconds
	r.recording = true
	r.paused = false
	r.lost.Store(false)
	r.lastData.Store(time.Now().UnixNano())
	r.mu.Unlock()
	r.resetPeak()

	if err := r.startDevice(); err != nil {
		r.mu.Lock()
		r.recording = false
		r.mu.Unlock()
		r.stopDevice()
		return err
	}
	return nil
}

func (r *Recorder) SetPaused(v bool) {
	r.mu.Lock()
	r.paused = v
	r.mu.Unlock()
}

func (r *Recorder) Stop() []byte {
	r.mu.Lock()
	r.recording = false
	r.paused = false
	pcm := r.buf
	over := r.overflow
	max := r.maxBytes
	r.buf = nil
	r.overflow = false
	r.mu.Unlock()
	r.stopDevice()
	if over {
		log.Printf("запись упёрлась в предел буфера (%d с), хвост отброшен", max/(sampleRate*2))
	}
	return pcm
}

func (r *Recorder) MonitorPing() {
	r.mu.Lock()
	r.monitor = true
	r.monUntil = time.Now().Add(monitorLinger)
	need := !r.started
	r.mu.Unlock()
	if !need {
		return
	}
	if err := r.startDevice(); err != nil {
		log.Printf("монитор микрофона: %v", err)
		r.mu.Lock()
		r.monitor = false
		r.mu.Unlock()
	}
}

func (r *Recorder) watchMonitor() {
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-r.done:
			return
		case <-t.C:
		}
		r.mu.Lock()
		expired := r.monitor && time.Now().After(r.monUntil)
		if expired {
			r.monitor = false
		}
		recording, paused := r.recording, r.paused
		r.mu.Unlock()
		if expired {
			r.stopDevice()
		}
		if recording && !paused && !r.lost.Load() {
			last := r.lastData.Load()
			if last > 0 && time.Since(time.Unix(0, last)) > micSilenceLimit {
				r.lost.Store(true)
				r.level.Store(0)
				log.Printf("микрофон перестал отдавать звук — считаю его отключённым")
				if r.onLost != nil {
					go r.onLost()
				}
			}
		}
	}
}

func (r *Recorder) Close() {
	r.once.Do(func() { close(r.done) })
	r.mu.Lock()
	dev := r.device
	r.device = nil
	r.recording = false
	r.monitor = false
	r.started = false
	r.mu.Unlock()
	if dev != nil {
		_ = dev.Stop()
		dev.Uninit()
	}
	if r.ctx != nil {
		r.ctx.Uninit()
		r.ctx.Free()
	}
}

func pcmPeak(pcm []byte) float64 { return audiolevel.Peak(pcm) }

func pcmIsSilent(pcm []byte) bool { return audiolevel.IsSilent(pcm) }

var errRecorderBusy = errors.New("микрофон занят записью")
