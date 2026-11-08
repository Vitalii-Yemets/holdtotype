package audiolevel

import (
	"math"
	"time"
)

const (
	DefaultPauseMs = 2000
	MinPauseMs     = 1000
	MaxPauseMs     = 10000

	endpointFrameMs = 20
	zeroDB          = -80
	floorLearnMs    = 500
	floorRise       = 0.02
	onAboveFloor    = 12
	onLowest        = -50
	onHighest       = -30
	onMinMargin     = 8
	offBelowOn      = 4
	voiceRunFrames  = 3
	armSpeechMs     = 400
)

func ValidPauseMs(ms int) bool { return ms >= MinPauseMs && ms <= MaxPauseMs }

type endpointState int

const (
	epGuard endpointState = iota
	epWaiting
	epArmed
	epFired
)

type Endpoint struct {
	rate         int
	frame        int
	guardSamples int
	pauseSamples int
	learnSamples int
	armSamples   int

	state     endpointState
	pending   []byte
	frames    int
	learned   int
	floor     float64
	floorSet  bool
	run       int
	speech    int
	quiet     int
	speechEnd int
}

func NewEndpoint(rate int, pause time.Duration, guard time.Duration) *Endpoint {
	if rate <= 0 {
		rate = 16000
	}
	if pause <= 0 {
		pause = DefaultPauseMs * time.Millisecond
	}
	if guard < 0 {
		guard = 0
	}
	frame := rate * endpointFrameMs / 1000
	if frame < 1 {
		frame = 1
	}
	return &Endpoint{
		rate:         rate,
		frame:        frame,
		guardSamples: durationSamples(guard, rate),
		pauseSamples: durationSamples(pause, rate),
		learnSamples: rate * floorLearnMs / 1000,
		armSamples:   rate * armSpeechMs / 1000,
		pending:      make([]byte, 0, frame*2),
	}
}

func durationSamples(d time.Duration, rate int) int {
	return int(int64(d) * int64(rate) / int64(time.Second))
}

func (e *Endpoint) Feed(pcm []byte) bool {
	frameBytes := e.frame * 2
	for len(pcm) > 0 && e.state != epFired {
		if len(e.pending) == 0 && len(pcm) >= frameBytes {
			if e.step(pcm[:frameBytes]) {
				return true
			}
			pcm = pcm[frameBytes:]
			continue
		}
		need := frameBytes - len(e.pending)
		if need > len(pcm) {
			need = len(pcm)
		}
		e.pending = append(e.pending, pcm[:need]...)
		pcm = pcm[need:]
		if len(e.pending) == frameBytes {
			fired := e.step(e.pending)
			e.pending = e.pending[:0]
			if fired {
				return true
			}
		}
	}
	return false
}

func (e *Endpoint) step(frame []byte) bool {
	if e.state == epGuard && e.frames*e.frame >= e.guardSamples {
		e.state = epWaiting
	}
	var sum float64
	for i := 0; i+1 < len(frame); i += 2 {
		v := sampleAt(frame, i)
		sum += v * v
	}
	db := DBFS(math.Sqrt(sum / float64(len(frame)/2)))
	e.frames++
	end := e.frames * len(frame)
	zero := db <= zeroDB
	learning := !zero && e.learned*e.frame < e.learnSamples
	if learning {
		if !e.floorSet || db < e.floor {
			e.floor = db
			e.floorSet = true
		}
		e.learned++
	}
	on := e.Threshold()
	voice := !zero && db >= on
	if e.state == epArmed {
		voice = !zero && db >= on-offBelowOn
	}
	if !zero && !voice && !learning {
		e.floor = math.Min(db, e.floor+floorRise)
	}
	switch e.state {
	case epWaiting:
		if !voice {
			e.run = 0
			return false
		}
		e.run++
		switch {
		case e.run == voiceRunFrames:
			e.speech += voiceRunFrames
		case e.run > voiceRunFrames:
			e.speech++
		default:
			return false
		}
		e.speechEnd = end
		if e.speech*e.frame >= e.armSamples {
			e.state = epArmed
			e.quiet = 0
		}
	case epArmed:
		if voice {
			e.quiet = 0
			e.speech++
			e.speechEnd = end
			return false
		}
		e.quiet++
		if e.quiet*e.frame >= e.pauseSamples {
			e.state = epFired
			return true
		}
	}
	return false
}

func (e *Endpoint) Heard() bool { return e.state == epArmed || e.state == epFired }

func (e *Endpoint) Fired() bool { return e.state == epFired }

func (e *Endpoint) SpeechMs() int { return e.speech * e.frame * 1000 / e.rate }

func (e *Endpoint) SpeechEnd() int { return e.speechEnd }

func (e *Endpoint) Floor() float64 {
	if !e.floorSet {
		return zeroDB
	}
	return e.floor
}

func (e *Endpoint) Threshold() float64 {
	floor := e.Floor()
	on := math.Max(onLowest, math.Min(onHighest, floor+onAboveFloor))
	return math.Max(on, floor+onMinMargin)
}
