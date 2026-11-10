package subtitles

import "sync/atomic"

type Gate struct {
	hold    atomic.Bool
	lead    []byte
	leadMax int
	held    int64
}

func NewGate(rate int) *Gate {
	if rate <= 0 {
		rate = 16000
	}
	return &Gate{leadMax: rate * 2 * 200 / 1000}
}

func (g *Gate) Final() { g.hold.Store(true) }

func (g *Gate) Holding() bool { return g.hold.Load() }

func (g *Gate) Held() int64 { return g.held }

func (g *Gate) Pass(pcm []byte) []byte {
	if !g.hold.Load() {
		return pcm
	}
	if Silent(pcm) {
		g.held += int64(len(pcm) / 2)
		g.lead = append(g.lead, pcm...)
		if over := len(g.lead) - g.leadMax; over > 0 {
			n := copy(g.lead, g.lead[over:])
			g.lead = g.lead[:n]
		}
		return nil
	}
	g.hold.Store(false)
	out := make([]byte, 0, len(g.lead)+len(pcm))
	out = append(out, g.lead...)
	out = append(out, pcm...)
	g.held -= int64(len(g.lead) / 2)
	g.lead = g.lead[:0]
	return out
}
