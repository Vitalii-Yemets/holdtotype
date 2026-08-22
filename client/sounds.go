package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	winmm          = windows.NewLazySystemDLL("winmm.dll")
	procPlaySoundW = winmm.NewProc("PlaySoundW")
)

const (
	sndAsync     = 0x0001
	sndNoDefault = 0x0002
	sndMemory    = 0x0004
	sndFilename  = 0x00020000
)

const (
	cueStart = iota
	cueStop
	cueError
)

type note struct {
	freq  float64
	durMs int
	amp   float64
}

type timbre struct {
	harmonic float64
	decay    float64
	attackMs int
	noise    float64
}

var timbres = map[string]timbre{
	"chime":   {harmonic: 0.35, decay: 4.5, attackMs: 6, noise: 0},
	"soft":    {harmonic: 0.05, decay: 3.0, attackMs: 25, noise: 0},
	"marimba": {harmonic: 0.6, decay: 8.0, attackMs: 3, noise: 0},
	"blip":    {harmonic: 0.9, decay: 14.0, attackMs: 1, noise: 0},
	"pop":     {harmonic: 0.2, decay: 18.0, attackMs: 1, noise: 0.15},
}

func synthWavT(notes []note, tb timbre) []byte {
	const rate = 44100
	var pcm []int16
	seed := uint32(12345)
	gap := make([]int16, rate*25/1000)
	for ni, n := range notes {
		if ni > 0 {
			pcm = append(pcm, gap...)
		}
		total := rate * n.durMs / 1000
		attack := rate * tb.attackMs / 1000
		if attack < 1 {
			attack = 1
		}
		release := rate * 12 / 1000
		for i := 0; i < total; i++ {
			t := float64(i) / rate
			env := math.Exp(-tb.decay * float64(i) / float64(total))
			if i < attack {
				env *= float64(i) / float64(attack)
			}
			if left := total - i; left < release {
				env *= float64(left) / float64(release)
			}
			s := math.Sin(2*math.Pi*n.freq*t) + tb.harmonic*math.Sin(4*math.Pi*n.freq*t)
			if tb.noise > 0 {
				seed = seed*1664525 + 1013904223
				s += tb.noise * (float64(seed>>8&0xFFFF)/32768 - 1)
			}
			pcm = append(pcm, int16(s*env*n.amp*20000))
		}
	}
	var b bytes.Buffer
	le := binary.LittleEndian
	w32 := func(v uint32) { _ = binary.Write(&b, le, v) }
	w16 := func(v uint16) { _ = binary.Write(&b, le, v) }
	dataLen := len(pcm) * 2
	b.WriteString("RIFF")
	w32(uint32(36 + dataLen))
	b.WriteString("WAVE")
	b.WriteString("fmt ")
	w32(16)
	w16(1)
	w16(1)
	w32(rate)
	w32(rate * 2)
	w16(2)
	w16(16)
	b.WriteString("data")
	w32(uint32(dataLen))
	_ = binary.Write(&b, le, pcm)
	return b.Bytes()
}

type cueSet struct {
	start, stop, fail []note
	tb                timbre
}

var soundThemes = map[string]cueSet{
	"chime": {
		start: []note{{523.25, 70, 0.30}, {783.99, 110, 0.32}},
		stop:  []note{{659.25, 70, 0.28}, {523.25, 120, 0.28}},
		fail:  []note{{233.08, 180, 0.32}, {207.65, 200, 0.30}},
		tb:    timbres["chime"],
	},
	"soft": {
		start: []note{{440.00, 150, 0.32}},
		stop:  []note{{329.63, 190, 0.30}},
		fail:  []note{{196.00, 260, 0.32}},
		tb:    timbres["soft"],
	},
	"marimba": {
		start: []note{{587.33, 90, 0.34}, {880.00, 130, 0.30}},
		stop:  []note{{880.00, 80, 0.30}, {587.33, 140, 0.32}},
		fail:  []note{{261.63, 150, 0.32}, {196.00, 200, 0.30}},
		tb:    timbres["marimba"],
	},
	"blip": {
		start: []note{{1046.50, 45, 0.26}},
		stop:  []note{{698.46, 55, 0.26}},
		fail:  []note{{349.23, 90, 0.28}, {261.63, 110, 0.26}},
		tb:    timbres["blip"],
	},
	"pop": {
		start: []note{{392.00, 60, 0.34}},
		stop:  []note{{294.00, 70, 0.32}},
		fail:  []note{{174.61, 140, 0.34}},
		tb:    timbres["pop"],
	},
}

var fileThemes = map[string][3]string{
	"speech": {"Speech On.wav", "Speech Off.wav", "Speech Misrecognition.wav"},
}

func validSoundTheme(name string) bool {
	if _, ok := soundThemes[name]; ok {
		return true
	}
	_, ok := fileThemes[name]
	return ok
}

var (
	wavCache  sync.Map
	pathCache sync.Map
)

func cueFilePath(theme string, kind int) *uint16 {
	names, ok := fileThemes[theme]
	if !ok {
		return nil
	}
	key := fmt.Sprintf("%s/%d", theme, kind)
	if v, ok := pathCache.Load(key); ok {
		if v == nil {
			return nil
		}
		return v.(*uint16)
	}
	idx := kind
	if idx < 0 || idx > 2 {
		idx = 2
	}
	full := filepath.Join(os.Getenv("WINDIR"), "Media", names[idx])
	if _, err := os.Stat(full); err != nil {
		pathCache.Store(key, nil)
		return nil
	}
	p, _ := windows.UTF16PtrFromString(full)
	pathCache.Store(key, p)
	return p
}

func cueWav(theme string, kind int) []byte {
	set, ok := soundThemes[theme]
	if !ok {
		set = soundThemes["chime"]
		theme = "chime"
	}
	key := fmt.Sprintf("%s/%d", theme, kind)
	if v, ok := wavCache.Load(key); ok {
		return v.([]byte)
	}
	var notes []note
	switch kind {
	case cueStart:
		notes = set.start
	case cueStop:
		notes = set.stop
	default:
		notes = set.fail
	}
	wav := synthWavT(notes, set.tb)
	wavCache.Store(key, wav)
	return wav
}

func playWav(wav []byte) {
	if len(wav) == 0 {
		return
	}
	procPlaySoundW.Call(uintptr(unsafe.Pointer(&wav[0])), 0, sndMemory|sndAsync|sndNoDefault)
}

func playTheme(theme string, kind int) {
	if p := cueFilePath(theme, kind); p != nil {
		procPlaySoundW.Call(uintptr(unsafe.Pointer(p)), 0, sndFilename|sndAsync|sndNoDefault)
		return
	}
	playWav(cueWav(theme, kind))
}

func playCue(enabled bool, theme string, kind int) {
	if !enabled {
		return
	}
	playTheme(theme, kind)
}

func previewTheme(theme string) {
	go func() {
		playTheme(theme, cueStart)
		time.Sleep(900 * time.Millisecond)
		playTheme(theme, cueStop)
	}()
}
