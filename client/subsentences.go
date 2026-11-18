package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"holdtotype/internal/sentence"
	"holdtotype/internal/subtitles"
)

const (
	sentQuietEvery = 200 * time.Millisecond
	sentFlushWait  = 5 * time.Second
)

type heardKind int

const (
	heardTick heardKind = iota
	heardUpdate
	heardWhole
	heardFlush
)

type heardEvent struct {
	kind  heardKind
	seg   int
	text  string
	final bool
	from  time.Time
	to    time.Time
}

func (ev heardEvent) cut(ph *subtitles.Phraser, at time.Time) []subtitles.TextPhrase {
	switch ev.kind {
	case heardTick:
		return ph.Tick(at)
	case heardUpdate:
		return ph.Update(ev.seg, ev.text, ev.final, at)
	case heardWhole:
		return ph.Whole(ev.text, ev.from, ev.to, at)
	default:
		return ph.Flush(at)
	}
}

type sentItem struct {
	op    sentence.Op
	flush chan struct{}
}

type sentRun struct {
	r      *subRun
	params sentence.Params
	heard  *sentence.Heard

	mu      sync.Mutex
	queue   []sentItem
	failed  bool
	resumed bool
	why     string
	seq     int
	closed  bool
	wake    chan struct{}
	done    chan struct{}
	counts  map[string]int
	lastEnd time.Time
}

func punctModelMissing() string {
	for _, f := range []string{"model_quantized.onnx", "tokenizer.json"} {
		if _, err := os.Stat(filepath.Join(punctModelDir, f)); err != nil {
			return fmt.Sprintf("the punctuation model is missing: %s", filepath.Join(punctModelDir, f))
		}
	}
	return ""
}

func sentencesWanted(cfg *Config, t *subTarget) bool {
	return cfg.SubtitlesSentences && translationTarget(cfg, t) != ""
}

func newSentRun(r *subRun, cfg *Config, t *subTarget, from subtitles.PhraserState, seq int) *sentRun {
	if !cfg.SubtitlesSentences {
		return nil
	}
	if translationTarget(cfg, t) == "" {
		log.Printf("subtitles: sentences are not cut by the punctuation model: translation is off")
		return nil
	}
	if why := punctModelMissing(); why != "" {
		log.Printf("subtitles: sentences are not cut by the punctuation model: %s", why)
		return nil
	}
	s := &sentRun{r: r, params: sentence.Defaults(), heard: sentence.NewHeard(), wake: make(chan struct{}, 1), done: make(chan struct{}), counts: map[string]int{}, seq: seq}
	for _, op := range s.heard.Seed(from.Seg, from.Done, from.Words, from.Given, from.Carry, time.Now()) {
		s.queue = append(s.queue, sentItem{op: op})
	}
	p := s.params
	log.Printf("subtitles: sentence ends come from the punctuation model: %d word(s) ahead, silence %.1f s, limit %d words, shortest %d words", p.Ahead, p.Silence.Seconds(), p.Limit, p.MinWords)
	go s.loop()
	return s
}

func (s *sentRun) feed(ev heardEvent, at time.Time) bool {
	s.mu.Lock()
	if s.failed || s.closed {
		s.mu.Unlock()
		return false
	}
	var ops []sentence.Op
	switch ev.kind {
	case heardUpdate:
		ops = s.heard.Update(ev.seg, ev.text, ev.final, at)
	case heardWhole:
		ops = s.heard.Whole(ev.text, at)
	case heardFlush:
		ops = s.heard.Flush(at)
	}
	for _, op := range ops {
		s.queue = append(s.queue, sentItem{op: op})
	}
	var ack chan struct{}
	if ev.kind == heardFlush {
		ack = make(chan struct{})
		s.queue = append(s.queue, sentItem{flush: ack})
	}
	s.mu.Unlock()
	s.signal()
	if ack != nil {
		select {
		case <-ack:
		case <-s.done:
		case <-time.After(sentFlushWait):
			log.Printf("subtitles: the sentence ends were not all found within %.0f s of the end", sentFlushWait.Seconds())
		}
	}
	return true
}

func (s *sentRun) signal() {
	select {
	case s.wake <- struct{}{}:
	default:
	}
}

func (s *sentRun) active() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return !s.failed && !s.closed
}

func (s *sentRun) waiting() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.heard.Waiting() + len(s.queue)
}

func (s *sentRun) takeOver() (int, int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.failed || s.resumed {
		return 0, 0, false
	}
	s.resumed = true
	return s.heard.Given(), s.seq, true
}

func (s *sentRun) emit(ps []sentence.Piece) {
	for _, p := range ps {
		s.mu.Lock()
		s.seq++
		seq := s.seq
		para := s.lastEnd.IsZero() || p.Start.Sub(s.lastEnd) >= subtitles.PhrasePause
		s.lastEnd = p.Last
		s.counts[p.How]++
		s.mu.Unlock()
		if s.r.opts.logLive {
			log.Printf("subtitles: sentence %d, %d words, closed by %s %.1f s after its last word", seq, p.Words, p.How, p.At.Sub(p.Last).Seconds())
		}
		s.r.phraseCut(subtitles.TextPhrase{Seq: seq, Text: p.Text, Words: p.Words, Start: p.Start, Heard: time.Now(), Para: para}, s.r.target())
	}
}

func (s *sentRun) fail(c *sentence.Cutter, why string) {
	s.mu.Lock()
	s.failed = true
	s.why = why
	items := s.queue
	s.queue = nil
	var acks []chan struct{}
	if c != nil {
		for _, it := range items {
			if it.flush != nil {
				acks = append(acks, it.flush)
				continue
			}
			_, _ = c.Apply(it.op)
		}
	}
	log.Printf("subtitles: sentence ends are cut by phrases from now on: %s", why)
	if c != nil {
		s.mu.Unlock()
		s.emit(c.Abandon())
		s.mu.Lock()
	}
	s.mu.Unlock()
	for _, a := range acks {
		close(a)
	}
}

func (s *sentRun) loop() {
	defer close(s.done)
	started := time.Now()
	eng, err := startPunctEngine(1)
	if err != nil {
		s.fail(nil, err.Error())
		return
	}
	defer eng.stop()
	log.Printf("subtitles: the punctuation engine is ready after %.1f s", time.Since(started).Seconds())
	broken := false
	judge := func(ws []string) ([]bool, []float64, error) {
		if broken {
			return nil, nil, fmt.Errorf("the punctuation engine stopped")
		}
		a, err := eng.judge(ws)
		if err != nil {
			broken = true
			return nil, nil, err
		}
		return a.End, a.Prob, nil
	}
	c := sentence.NewCutter(s.params, judge)
	tick := time.NewTicker(sentQuietEvery)
	defer tick.Stop()
	for {
		select {
		case <-s.wake:
		case <-tick.C:
			s.emit(c.Quiet(time.Now()))
			continue
		}
		for {
			s.mu.Lock()
			if len(s.queue) == 0 {
				closed := s.closed
				s.mu.Unlock()
				if closed {
					s.emit(c.Flush())
					return
				}
				break
			}
			it := s.queue[0]
			s.queue = s.queue[1:]
			s.mu.Unlock()
			if it.flush != nil {
				s.emit(c.Flush())
				close(it.flush)
				continue
			}
			ps, err := c.Apply(it.op)
			s.emit(ps)
			if err != nil {
				s.fail(c, err.Error())
				return
			}
		}
	}
}

func (s *sentRun) halt(why string) {
	s.mu.Lock()
	s.failed = true
	s.why = why
	s.mu.Unlock()
	log.Printf("subtitles: the punctuation engine stops: %s", why)
	s.close()
}

func (s *sentRun) close() {
	s.mu.Lock()
	s.closed = true
	s.mu.Unlock()
	s.signal()
	select {
	case <-s.done:
	case <-time.After(sentFlushWait):
	}
	s.mu.Lock()
	n := s.counts
	s.mu.Unlock()
	log.Printf("subtitles: sentences: %d closed by the model, %d by silence, %d by the limit, %d at the end, %d handed over when the engine stopped",
		n[sentence.ByModel], n[sentence.BySilence], n[sentence.ByLimit], n[sentence.ByEnd], n[sentence.ByFallback])
}
