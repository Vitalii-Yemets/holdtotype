package subtitles

import (
	"math"
	"time"
)

const (
	FrameMs     = 20
	PrerollMs   = 150
	PauseMs     = 600
	TailMs      = 200
	MaxPhraseMs = 12000
	CutWindowMs = 3000
	MinSpeechMs = 300

	onsetFrames  = 3
	floorLearnMs = 500
	floorRise    = 0.02
	onAboveFloor = 12.0
	onLowest     = -50.0
	onHighest    = -30.0
	onMinMargin  = 8.0
	offBelowOn   = 4.0
	zeroDB       = -80.0
)

type Phrase struct {
	PCM    []byte
	Start  int64
	Speech time.Duration
	Cut    bool
	At     time.Time
}

func (p Phrase) Duration(rate int) time.Duration {
	return SamplesDuration(int64(len(p.PCM)/2), rate)
}

func SamplesDuration(samples int64, rate int) time.Duration {
	if rate <= 0 {
		return 0
	}
	return time.Duration(samples * int64(time.Second) / int64(rate))
}

type Segmenter struct {
	rate       int
	frameBytes int
	preMax     int
	maxBytes   int
	pauseFr    int
	tailBytes  int
	cutFr      int
	minSpeech  int

	pending []byte
	pos     int64

	pre []byte
	run int

	inPhrase bool
	cur      []byte
	curStart int64
	dbs      []float64
	voice    []bool
	quiet    int

	floor    float64
	floorSet bool
	learned  int
}

func NewSegmenter(rate int) *Segmenter {
	if rate <= 0 {
		rate = 16000
	}
	frame := rate * FrameMs / 1000
	ms := func(v int) int { return rate * v / 1000 * 2 }
	return &Segmenter{
		rate:       rate,
		frameBytes: frame * 2,
		preMax:     ms(PrerollMs) + onsetFrames*frame*2,
		maxBytes:   ms(MaxPhraseMs),
		pauseFr:    PauseMs / FrameMs,
		tailBytes:  ms(TailMs),
		cutFr:      CutWindowMs / FrameMs,
		minSpeech:  MinSpeechMs / FrameMs,
		pending:    make([]byte, 0, frame*2),
	}
}

func (s *Segmenter) InPhrase() bool { return s.inPhrase }

func (s *Segmenter) Buffered() int { return len(s.pre) + len(s.cur) + len(s.pending) }

func (s *Segmenter) Floor() float64 {
	if !s.floorSet {
		return zeroDB
	}
	return s.floor
}

func (s *Segmenter) Threshold() float64 {
	floor := s.Floor()
	on := math.Max(onLowest, math.Min(onHighest, floor+onAboveFloor))
	return math.Max(on, floor+onMinMargin)
}

func (s *Segmenter) Feed(pcm []byte) []Phrase {
	var out []Phrase
	for len(pcm) > 0 {
		if len(s.pending) == 0 && len(pcm) >= s.frameBytes {
			out = s.step(pcm[:s.frameBytes], out)
			pcm = pcm[s.frameBytes:]
			continue
		}
		need := s.frameBytes - len(s.pending)
		if need > len(pcm) {
			need = len(pcm)
		}
		s.pending = append(s.pending, pcm[:need]...)
		pcm = pcm[need:]
		if len(s.pending) == s.frameBytes {
			out = s.step(s.pending, out)
			s.pending = s.pending[:0]
		}
	}
	return out
}

func (s *Segmenter) Flush() []Phrase {
	if !s.inPhrase {
		return nil
	}
	var out []Phrase
	out = s.finish(out)
	return out
}

func frameDB(frame []byte) float64 {
	var sum float64
	n := len(frame) / 2
	for i := 0; i+1 < len(frame); i += 2 {
		v := float64(int16(uint16(frame[i])|uint16(frame[i+1])<<8)) / 32768
		sum += v * v
	}
	if n == 0 {
		return zeroDB
	}
	rms := math.Sqrt(sum / float64(n))
	if rms <= 0.0001 {
		return zeroDB
	}
	db := 20 * math.Log10(rms)
	if db < zeroDB {
		return zeroDB
	}
	return db
}

func (s *Segmenter) step(frame []byte, out []Phrase) []Phrase {
	db := frameDB(frame)
	s.pos += int64(len(frame) / 2)
	zero := db <= zeroDB
	learning := !zero && db < onHighest && s.learned*FrameMs < floorLearnMs
	if learning {
		if !s.floorSet || db < s.floor {
			s.floor = db
			s.floorSet = true
		}
		s.learned++
	}
	on := s.Threshold()
	voiced := !zero && db >= on
	if s.inPhrase {
		voiced = !zero && db >= on-offBelowOn
	}
	if !zero && !voiced && !learning {
		s.floor = math.Min(db, s.floor+floorRise)
	}
	if !s.inPhrase {
		s.pushPre(frame)
		if !voiced {
			s.run = 0
			return out
		}
		s.run++
		if s.run < onsetFrames {
			return out
		}
		s.begin()
		return out
	}
	s.cur = append(s.cur, frame...)
	s.dbs = append(s.dbs, db)
	s.voice = append(s.voice, voiced)
	if voiced {
		s.quiet = 0
	} else {
		s.quiet++
	}
	if s.quiet >= s.pauseFr {
		return s.finish(out)
	}
	if len(s.cur) >= s.maxBytes {
		return s.cut(out)
	}
	return out
}

func (s *Segmenter) pushPre(frame []byte) {
	s.pre = append(s.pre, frame...)
	if over := len(s.pre) - s.preMax; over > 0 {
		n := copy(s.pre, s.pre[over:])
		s.pre = s.pre[:n]
	}
}

func (s *Segmenter) begin() {
	frames := len(s.pre) / s.frameBytes
	s.inPhrase = true
	s.cur = append(s.cur[:0], s.pre...)
	s.curStart = s.pos - int64(len(s.pre)/2)
	s.dbs = s.dbs[:0]
	s.voice = s.voice[:0]
	for i := 0; i < frames; i++ {
		s.dbs = append(s.dbs, frameDB(s.pre[i*s.frameBytes:(i+1)*s.frameBytes]))
		s.voice = append(s.voice, i >= frames-onsetFrames)
	}
	s.quiet = 0
	s.run = 0
	s.pre = s.pre[:0]
}

func (s *Segmenter) lastVoiceEnd(frames int) int {
	for i := frames - 1; i >= 0; i-- {
		if s.voice[i] {
			return (i + 1) * s.frameBytes
		}
	}
	return 0
}

func (s *Segmenter) emit(out []Phrase, bytes int, cut bool) []Phrase {
	frames := bytes / s.frameBytes
	speech := 0
	first := -1
	for i := 0; i < frames && i < len(s.voice); i++ {
		if s.voice[i] {
			speech++
			if first < 0 {
				first = i
			}
		}
	}
	if speech < s.minSpeech || first < 0 {
		return out
	}
	pcm := make([]byte, bytes)
	copy(pcm, s.cur[:bytes])
	return append(out, Phrase{
		PCM:    pcm,
		Start:  s.curStart + int64(first*s.frameBytes/2),
		Speech: time.Duration(speech*FrameMs) * time.Millisecond,
		Cut:    cut,
	})
}

func (s *Segmenter) finish(out []Phrase) []Phrase {
	frames := len(s.cur) / s.frameBytes
	end := s.lastVoiceEnd(frames) + s.tailBytes
	if end > frames*s.frameBytes {
		end = frames * s.frameBytes
	}
	out = s.emit(out, end, false)
	s.pre = s.pre[:0]
	tail := len(s.cur) - s.preMax
	if tail < 0 {
		tail = 0
	}
	s.pre = append(s.pre, s.cur[tail:]...)
	s.inPhrase = false
	s.cur = s.cur[:0]
	s.dbs = s.dbs[:0]
	s.voice = s.voice[:0]
	s.quiet = 0
	s.run = 0
	return out
}

func (s *Segmenter) cut(out []Phrase) []Phrase {
	frames := len(s.dbs)
	from := frames - s.cutFr
	if from < 1 {
		from = 1
	}
	best := frames - 1
	for i := from; i < frames; i++ {
		if s.dbs[i] < s.dbs[best] {
			best = i
		}
	}
	at := (best + 1) * s.frameBytes
	out = s.emit(out, at, true)
	restFrames := frames - (best + 1)
	n := copy(s.cur, s.cur[at:])
	s.cur = s.cur[:n]
	copy(s.dbs, s.dbs[best+1:])
	s.dbs = s.dbs[:restFrames]
	copy(s.voice, s.voice[best+1:])
	s.voice = s.voice[:restFrames]
	s.curStart += int64(at / 2)
	s.quiet = 0
	for i := restFrames - 1; i >= 0 && !s.voice[i]; i-- {
		s.quiet++
	}
	if s.quiet >= s.pauseFr {
		return s.finish(out)
	}
	return out
}

func (s *Segmenter) StartAt(pos int64) {
	s.pos = pos
}
