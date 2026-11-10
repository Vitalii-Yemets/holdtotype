package subtitles

import "time"

const (
	TickMs       = 100
	SlackMs      = 100
	SilenceLimit = 4 * time.Second
	silentPeak   = 66
)

type Pacer struct {
	rate  int
	slack int64
	limit int64

	base      time.Time
	fed       int64
	captured  int64
	added     int64
	silentRun int64
	suspended bool
	pauses    int
}

func NewPacer(rate int, start time.Time) *Pacer {
	if rate <= 0 {
		rate = 16000
	}
	return &Pacer{
		rate:  rate,
		slack: int64(rate) * SlackMs / 1000,
		limit: int64(SilenceLimit) * int64(rate) / int64(time.Second),
		base:  start,
	}
}

func Silent(pcm []byte) bool {
	for i := 0; i+1 < len(pcm); i += 2 {
		v := int32(int16(uint16(pcm[i]) | uint16(pcm[i+1])<<8))
		if v < 0 {
			v = -v
		}
		if v >= silentPeak {
			return false
		}
	}
	return true
}

func (p *Pacer) samplesSince(now time.Time) int64 {
	d := now.Sub(p.base)
	if d <= 0 {
		return 0
	}
	return int64(d) * int64(p.rate) / int64(time.Second)
}

func (p *Pacer) Tick(now time.Time, captured []byte) []byte {
	n := int64(len(captured) / 2)
	captured = captured[:n*2]
	p.captured += n
	silent := Silent(captured)
	if p.suspended {
		if silent {
			return nil
		}
		p.suspended = false
		p.silentRun = 0
		p.base = now.Add(-SamplesDuration(p.fed+n, p.rate))
		p.fed += n
		return captured
	}
	var pad int64
	if target := p.samplesSince(now) - p.slack; p.fed+n < target {
		pad = target - (p.fed + n)
	}
	var out []byte
	if pad == 0 {
		out = captured
	} else {
		out = make([]byte, (pad+n)*2)
		copy(out[pad*2:], captured)
	}
	p.fed += pad + n
	p.added += pad
	if silent {
		p.silentRun += pad + n
	} else {
		p.silentRun = 0
	}
	if p.silentRun >= p.limit {
		p.suspended = true
		p.pauses++
	}
	return out
}

func (p *Pacer) TimeAt(sample int64) time.Time {
	return p.base.Add(SamplesDuration(sample, p.rate))
}

func (p *Pacer) Fed() int64 { return p.fed }

func (p *Pacer) Captured() time.Duration { return SamplesDuration(p.captured, p.rate) }

func (p *Pacer) Added() time.Duration { return SamplesDuration(p.added, p.rate) }

func (p *Pacer) Suspended() bool { return p.suspended }

func (p *Pacer) Pauses() int { return p.pauses }
