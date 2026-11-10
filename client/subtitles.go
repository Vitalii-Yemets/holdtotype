package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"holdtotype/internal/evqueue"
	"holdtotype/internal/livetail"
	"holdtotype/internal/replace"
	"holdtotype/internal/segments"
	"holdtotype/internal/subtitles"
)

const (
	modeLive    = "live"
	modePhrases = "phrases"

	subPrimaryWait   = 30 * time.Second
	subDictationWait = 20 * time.Second
	subPhraseTimeout = 2 * time.Minute
	subStopDrain     = 20 * time.Second
	subRestartDrain  = 5 * time.Second
	subExitDrain     = 3 * time.Second
	subFailDrain     = 3 * time.Second
	subLiveDrain     = 3 * time.Second
	subDeviceGrace   = 3 * time.Second
	subDeviceGiveUp  = 30 * time.Second
	subBoardKeep     = 50
	subRecoverWindow = 2 * time.Minute
	subRecoverMax    = 3
)

var errSubStopped = errors.New("subtitles are stopping")

type subTarget struct {
	lang   string
	engine string
	model  string
	name   string
	mode   string
	cfg    *Config
}

func (t *subTarget) same(o *subTarget) bool {
	return o != nil && t.lang == o.lang && t.engine == o.engine && t.model == o.model
}

func (t *subTarget) diff(o *subTarget) string {
	switch {
	case t.lang != o.lang:
		return fmt.Sprintf("the language changed from %s to %s", t.lang, o.lang)
	case t.engine != o.engine:
		return fmt.Sprintf("the engine changed from %s to %s", t.engine, o.engine)
	}
	return fmt.Sprintf("the model changed from %s to %s", t.name, o.name)
}

func resolveSubtitles(cfg *Config) (*subTarget, error) {
	c := *cfg
	lang := cfg.SubtitlesLanguage
	if !validPresetLang(lang) {
		lang = "auto"
	}
	c.Language = lang
	c.CanaryTarget = ""
	m := activeModel(&c)
	if m == nil {
		return nil, fmt.Errorf("no model is chosen for the language %s", lang)
	}
	external := m.Engine == engineWhisper && (c.SttSource == "remote" || !c.ServerAutostart)
	if !external && !m.installed() {
		return nil, fmt.Errorf("the model %s for the language %s is not installed", m.NameKey, lang)
	}
	mode := modePhrases
	if m.Engine == engineStream {
		mode = modeLive
	}
	return &subTarget{lang: lang, engine: m.Engine, model: m.modelPath(), name: m.NameKey, mode: mode, cfg: &c}, nil
}

func subCompatible(srv recognizer, t *subTarget) bool {
	if srv.engine() != t.engine {
		return false
	}
	if srv.external() {
		return t.engine == engineWhisper
	}
	if srv.model() != t.model {
		return false
	}
	if s, ok := srv.(*sherpaServer); ok && s.src != "" && s.src != canarySrc(t.lang) {
		return false
	}
	return true
}

func engineGone(r recognizer) bool {
	select {
	case <-r.done():
		return true
	default:
		return false
	}
}

type subOptions struct {
	logLive bool
	onEnd   func(*subRun)
}

type subRun struct {
	app     *App
	opts    subOptions
	started time.Time

	cap   *loopback
	pacer *subtitles.Pacer
	board *subtitles.Board

	stopCh    chan struct{}
	stopOnce  sync.Once
	doneCh    chan struct{}
	restartCh chan *subTarget
	drain     atomic.Int64

	mu        sync.Mutex
	want      *subTarget
	failMsg   string
	engName   string
	engPort   int
	engShared bool
	engModel  string

	fileMu   sync.Mutex
	file     *os.File
	filePath string
	fileErr  bool

	skipped atomic.Int64
	muted   atomic.Int64

	liveMu   sync.Mutex
	liveLast time.Time
	liveText string

	devBad    time.Time
	devTried  time.Time
	devLogged bool
}

func startSubtitlesRun(a *App, cfg *Config, opts subOptions) (*subRun, error) {
	t, err := resolveSubtitles(cfg)
	if err != nil {
		return nil, err
	}
	lb, err := openLoopback()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	r := &subRun{
		app:       a,
		opts:      opts,
		started:   now,
		cap:       lb,
		pacer:     subtitles.NewPacer(sampleRate, now),
		board:     subtitles.NewBoard(subBoardKeep),
		stopCh:    make(chan struct{}),
		doneCh:    make(chan struct{}),
		restartCh: make(chan *subTarget, 1),
		want:      t,
	}
	log.Printf("subtitles: started: language %s, model %s, mode %s", t.lang, t.name, t.mode)
	go r.loop(t)
	return r, nil
}

func (r *subRun) done() <-chan struct{} { return r.doneCh }

func (r *subRun) target() *subTarget {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.want
}

func (r *subRun) stop(why string, drain time.Duration) {
	r.stopOnce.Do(func() {
		r.drain.Store(int64(drain))
		log.Printf("subtitles: stopping: %s", why)
		close(r.stopCh)
	})
}

func (r *subRun) fail(msg string) {
	r.mu.Lock()
	first := r.failMsg == ""
	if first {
		r.failMsg = msg
	}
	r.mu.Unlock()
	if first {
		log.Printf("subtitles: error: %s", msg)
	}
	r.stop("an error", subFailDrain)
}

func (r *subRun) restart(t *subTarget) {
	r.mu.Lock()
	r.want = t
	r.mu.Unlock()
	for {
		select {
		case r.restartCh <- t:
			return
		default:
		}
		select {
		case <-r.restartCh:
		default:
		}
	}
}

func (r *subRun) noteEngine(rec recognizer, shared bool, t *subTarget) {
	r.mu.Lock()
	r.engName = rec.engine()
	r.engPort = recognizerPort(rec)
	r.engShared = shared
	r.engModel = t.name
	r.mu.Unlock()
}

func (r *subRun) loop(t *subTarget) {
	defer close(r.doneCh)
	pipe := newSubPipe(r, t)
	tick := time.NewTicker(subtitles.TickMs * time.Millisecond)
	defer tick.Stop()
	dev := time.NewTicker(time.Second)
	defer dev.Stop()
	muting := false
	for {
		select {
		case <-r.stopCh:
			pipe.stop(time.Duration(r.drain.Load()))
			r.finish()
			return
		case next := <-r.restartCh:
			log.Printf("subtitles: restarting in place: %s", pipe.target.diff(next))
			pipe.stop(subRestartDrain)
			pipe = newSubPipe(r, next)
		case <-tick.C:
			now := time.Now()
			chunk := r.cap.take()
			n := cueMutes.Apply(chunk, now, sampleRate)
			if n > 0 {
				r.muted.Add(int64(n))
				if !muting {
					log.Printf("subtitles: the app's own sound is playing, it is silenced")
				}
			}
			muting = n > 0
			out := r.pacer.Tick(now, chunk)
			if len(out) > 0 {
				pipe.feed(out)
			}
		case <-dev.C:
			r.checkDevice(time.Now())
		}
	}
}

func (r *subRun) checkDevice(now time.Time) {
	if r.cap.running() {
		if !r.devBad.IsZero() {
			log.Printf("subtitles: the playback device delivers again")
		}
		r.devBad = time.Time{}
		r.devLogged = false
		return
	}
	if r.devBad.IsZero() {
		r.devBad = now
		log.Printf("subtitles: listening to the playback device stopped, waiting for it to come back")
		return
	}
	if now.Sub(r.devBad) < subDeviceGrace || now.Sub(r.devTried) < subDeviceGrace {
		return
	}
	r.devTried = now
	err := r.cap.reopen()
	if err == nil {
		log.Printf("subtitles: the playback device was reopened after %.0f s", now.Sub(r.devBad).Seconds())
		r.devBad = time.Time{}
		r.devLogged = false
		return
	}
	if !r.devLogged {
		r.devLogged = true
		log.Printf("subtitles: reopening the playback device: %v", err)
	}
	if now.Sub(r.devBad) >= subDeviceGiveUp {
		r.fail(errNoPlayback.Error())
	}
}

func (r *subRun) setLive(text string) {
	r.board.SetLive(text)
	if !r.opts.logLive {
		return
	}
	r.liveMu.Lock()
	defer r.liveMu.Unlock()
	if text == "" || text == r.liveText || time.Since(r.liveLast) < time.Second {
		return
	}
	r.liveText = text
	r.liveLast = time.Now()
	log.Printf("subtitles: live %q", livetail.Tail(text, 160))
}

func (r *subRun) finishLine(at time.Time, text string, t *subTarget) {
	text = strings.TrimSpace(segments.Join(strings.TrimSpace(text)))
	if text == "" || onlyNoise.MatchString(text) {
		return
	}
	cfg := r.app.snapshot()
	text = strings.TrimSpace(replace.Apply(replace.ForLang(cfg.Replacements, t.lang), text))
	if text == "" {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	ln := r.board.Add(at, text)
	log.Printf("subtitles: line %q", text)
	if cfg.SubtitlesSave {
		r.write(ln, t)
	}
}

func (r *subRun) write(ln subtitles.Line, t *subTarget) {
	r.fileMu.Lock()
	defer r.fileMu.Unlock()
	if r.fileErr {
		return
	}
	if r.file == nil {
		if err := os.MkdirAll(subtitles.Folder, 0o755); err != nil {
			log.Printf("subtitles: the transcript folder was not created: %v", err)
			r.fileErr = true
			return
		}
		p := filepath.Join(subtitles.Folder, subtitles.FileName(r.started))
		f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			log.Printf("subtitles: the transcript was not created: %v", err)
			r.fileErr = true
			return
		}
		if _, err := f.WriteString(subtitles.Opening(r.started, t.lang, t.name)); err != nil {
			log.Printf("subtitles: writing the transcript: %v", err)
		}
		r.file = f
		r.filePath = p
		if abs, aerr := filepath.Abs(p); aerr == nil {
			r.filePath = abs
		}
		log.Printf("subtitles: transcript %s", r.filePath)
	}
	if _, err := r.file.WriteString(subtitles.Entry(ln.Start, ln.Text)); err != nil {
		log.Printf("subtitles: writing the transcript: %v", err)
		return
	}
	_ = r.file.Sync()
}

func (r *subRun) finish() {
	r.cap.close()
	r.fileMu.Lock()
	if r.file != nil {
		_ = r.file.Close()
	}
	path := r.filePath
	r.fileMu.Unlock()
	r.mu.Lock()
	name, port, shared, model := r.engName, r.engPort, r.engShared, r.engModel
	r.mu.Unlock()
	how := "dedicated"
	if shared {
		how = "shared with dictation"
	}
	if name == "" {
		name, how = "none", "not started"
	}
	if path == "" {
		path = "none"
	}
	log.Printf("subtitles: summary: %.1f s captured, %.1f s of silence added, topping up paused %d time(s) after %.0f s of silence, %d line(s), %.1f s skipped, %.1f s of the app's own sounds silenced, engine %s on port %d with %s (%s), %d device open(s), transcript %s",
		r.pacer.Captured().Seconds(), r.pacer.Added().Seconds(), r.pacer.Pauses(), subtitles.SilenceLimit.Seconds(),
		r.board.Count(), time.Duration(r.skipped.Load()).Seconds(), float64(r.muted.Load())/sampleRate, name, port, model, how, r.cap.opens.Load(), path)
	if r.opts.onEnd != nil {
		r.opts.onEnd(r)
	}
}

type subPipe struct {
	run    *subRun
	target *subTarget
	stopCh chan struct{}
	killCh chan struct{}
	exitCh chan struct{}
	wg     sync.WaitGroup

	mu      sync.Mutex
	rec     recognizer
	shared  bool
	gen     int
	readyCh chan struct{}
	closed  bool

	recoverMu sync.Mutex
	recovered []time.Time

	seg      *subtitles.Segmenter
	qmu      sync.Mutex
	q        *subtitles.Backlog
	qsig     chan struct{}
	draining atomic.Bool
	waitLog  bool

	smu  sync.Mutex
	sess *segSession
	lmu  sync.Mutex
	live *subtitles.Live
	gate *subtitles.Gate
}

func newSubPipe(r *subRun, t *subTarget) *subPipe {
	p := &subPipe{
		run:     r,
		target:  t,
		stopCh:  make(chan struct{}),
		killCh:  make(chan struct{}),
		exitCh:  make(chan struct{}),
		readyCh: make(chan struct{}),
		qsig:    make(chan struct{}, 1),
	}
	if t.mode == modeLive {
		p.live = subtitles.NewLive()
		p.gate = subtitles.NewGate(sampleRate)
	} else {
		p.seg = subtitles.NewSegmenter(sampleRate)
		p.q = subtitles.NewBacklog(sampleRate, subtitles.BacklogLimit)
	}
	p.wg.Add(2)
	go func() {
		defer p.wg.Done()
		p.firstEngine()
	}()
	go func() {
		defer p.wg.Done()
		defer close(p.exitCh)
		if t.mode == modeLive {
			p.liveLoop()
		} else {
			p.phraseLoop()
		}
	}()
	return p
}

func (p *subPipe) stopped() bool {
	select {
	case <-p.stopCh:
		return true
	default:
		return false
	}
}

func (p *subPipe) sleep(d time.Duration) bool {
	select {
	case <-time.After(d):
		return true
	case <-p.stopCh:
		return false
	}
}

func (p *subPipe) firstEngine() {
	rec, shared, err := p.acquire()
	if err != nil {
		if !p.stopped() {
			p.run.fail(fmt.Sprintf("the engine did not start: %s", logText(err)))
		}
		return
	}
	p.install(rec, shared)
}

func (p *subPipe) acquire() (recognizer, bool, error) {
	a := p.run.app
	deadline := time.Now().Add(subPrimaryWait)
	logged := false
	for {
		a.mu.Lock()
		srv, ready, parked, quitting := a.srv, a.ready, a.parked, a.quitting
		a.mu.Unlock()
		if quitting {
			return nil, false, errors.New("the app is closing")
		}
		if srv == nil || parked || !subCompatible(srv, p.target) {
			break
		}
		if ready && !srv.wasStopped() && !engineGone(srv) {
			return srv, true, nil
		}
		if time.Now().After(deadline) {
			log.Printf("subtitles: the dictation engine did not come back within %.0f s", subPrimaryWait.Seconds())
			break
		}
		if !logged {
			logged = true
			log.Printf("subtitles: waiting for the dictation engine, it has the same model")
		}
		if !p.sleep(500 * time.Millisecond) {
			return nil, false, errSubStopped
		}
	}
	if p.stopped() {
		return nil, false, errSubStopped
	}
	log.Printf("subtitles: starting an engine of its own: %s with %s", p.target.engine, p.target.name)
	started := time.Now()
	srv, err := startEngineInstance(p.target.cfg, p.target.engine, p.target.model, logFile)
	if err != nil {
		return nil, false, err
	}
	done := make(chan error, 1)
	go func() { done <- srv.waitReady(engineReadyTimeout(srv)) }()
	select {
	case err := <-done:
		if err != nil {
			srv.stop()
			return nil, false, err
		}
	case <-p.stopCh:
		srv.stop()
		return nil, false, errSubStopped
	}
	log.Printf("subtitles: its own engine is ready in %.1f s", time.Since(started).Seconds())
	return srv, false, nil
}

func (p *subPipe) install(rec recognizer, shared bool) {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		if !shared {
			rec.stop()
		}
		return
	}
	p.rec = rec
	p.shared = shared
	p.gen++
	gen := p.gen
	close(p.readyCh)
	p.mu.Unlock()
	p.run.noteEngine(rec, shared, p.target)
	how := "its own"
	if shared {
		how = "shared with dictation"
	}
	log.Printf("subtitles: engine %s on port %d with %s, %s", rec.engine(), recognizerPort(rec), rec.model(), how)
	go func() {
		select {
		case <-rec.done():
			p.lost(gen, "the engine process ended")
		case <-p.stopCh:
		}
	}()
}

func (p *subPipe) current() (recognizer, int, bool, bool) {
	for {
		p.mu.Lock()
		rec, gen, shared, ready := p.rec, p.gen, p.shared, p.readyCh
		p.mu.Unlock()
		if rec != nil {
			return rec, gen, shared, true
		}
		select {
		case <-ready:
		case <-p.stopCh:
			return nil, 0, false, false
		}
	}
}

func (p *subPipe) lost(gen int, why string) {
	p.mu.Lock()
	if p.closed || gen != p.gen || p.rec == nil {
		p.mu.Unlock()
		return
	}
	old, shared := p.rec, p.shared
	p.rec = nil
	p.readyCh = make(chan struct{})
	p.mu.Unlock()
	log.Printf("subtitles: the engine %s went away (%s), recovering", old.engine(), why)
	go p.recover(old, shared)
}

func (p *subPipe) recover(old recognizer, shared bool) {
	p.recoverMu.Lock()
	defer p.recoverMu.Unlock()
	if !shared {
		old.stop()
	}
	now := time.Now()
	recent := p.recovered[:0]
	for _, at := range p.recovered {
		if now.Sub(at) < subRecoverWindow {
			recent = append(recent, at)
		}
	}
	p.recovered = recent
	if len(p.recovered) >= subRecoverMax {
		p.run.fail(fmt.Sprintf("the engine went away %d times in %.0f minutes", len(p.recovered)+1, subRecoverWindow.Minutes()))
		return
	}
	p.recovered = append(p.recovered, now)
	if !p.sleep(time.Second) {
		return
	}
	rec, sh, err := p.acquire()
	if err != nil {
		if !p.stopped() {
			p.run.fail(fmt.Sprintf("the engine could not be brought back: %s", logText(err)))
		}
		return
	}
	p.install(rec, sh)
}

func (p *subPipe) feed(pcm []byte) {
	if p.target.mode == modeLive {
		if pcm = p.gate.Pass(pcm); len(pcm) == 0 {
			return
		}
		p.smu.Lock()
		s := p.sess
		p.smu.Unlock()
		if s != nil {
			_ = s.push(pcm)
		}
		return
	}
	for _, ph := range p.seg.Feed(pcm) {
		p.enqueue(ph)
	}
}

func (p *subPipe) enqueue(ph subtitles.Phrase) {
	ph.At = p.run.pacer.TimeAt(ph.Start)
	p.qmu.Lock()
	dropped := p.q.Push(ph)
	p.qmu.Unlock()
	if len(dropped) > 0 {
		var d time.Duration
		for _, x := range dropped {
			d += x.Duration(sampleRate)
		}
		p.run.skipped.Add(int64(d))
		log.Printf("subtitles: recognition falls behind, %d phrase(s) skipped (%.1f s)", len(dropped), d.Seconds())
	}
	select {
	case p.qsig <- struct{}{}:
	default:
	}
}

func (p *subPipe) nextPhrase() (subtitles.Phrase, bool) {
	for {
		p.qmu.Lock()
		ph, n := p.q.Take(subtitles.JoinLimit)
		p.qmu.Unlock()
		if n > 1 {
			log.Printf("subtitles: %d waiting phrases go as one request (%.1f s)", n, ph.Duration(sampleRate).Seconds())
		}
		if n > 0 {
			return ph, true
		}
		if p.draining.Load() {
			return subtitles.Phrase{}, false
		}
		select {
		case <-p.qsig:
		case <-p.stopCh:
			return subtitles.Phrase{}, false
		}
	}
}

func (p *subPipe) waitDictation() {
	a := p.run.app
	if a.state.Load() != stProcessing {
		return
	}
	if !p.waitLog {
		p.waitLog = true
		log.Printf("subtitles: a dictation is being recognized on the same engine, the phrase waits")
	}
	deadline := time.Now().Add(subDictationWait)
	for a.state.Load() == stProcessing && time.Now().Before(deadline) {
		if !p.sleep(100 * time.Millisecond) {
			return
		}
	}
}

func (p *subPipe) phraseLoop() {
	fails := 0
	for {
		ph, ok := p.nextPhrase()
		if !ok {
			return
		}
		rec, gen, shared, ok := p.current()
		if !ok {
			return
		}
		if shared {
			p.waitDictation()
		}
		ctx, cancel := context.WithTimeout(context.Background(), subPhraseTimeout)
		go func() {
			select {
			case <-p.killCh:
				cancel()
			case <-ctx.Done():
			}
		}()
		started := time.Now()
		text, err := rec.transcribe(ctx, wavFromPCM16(ph.PCM, sampleRate), p.target.lang, "", false)
		cancel()
		if err != nil {
			select {
			case <-p.killCh:
				return
			default:
			}
			fails++
			log.Printf("subtitles: a phrase of %.1f s was not recognized: %v", ph.Duration(sampleRate).Seconds(), err)
			if fails >= 2 || engineGone(rec) {
				fails = 0
				p.lost(gen, "recognition fails")
			}
			continue
		}
		fails = 0
		if p.run.opts.logLive {
			log.Printf("subtitles: phrase of %.1f s recognized in %.1f s", ph.Duration(sampleRate).Seconds(), time.Since(started).Seconds())
		}
		p.run.finishLine(ph.At, text, p.target)
	}
}

func (p *subPipe) onSegment(seg int, text string, final bool) {
	if final {
		p.gate.Final()
	}
	p.lmu.Lock()
	out := p.live.Update(seg, text, final, time.Now())
	cur, _ := p.live.Current()
	p.lmu.Unlock()
	for _, f := range out {
		p.run.finishLine(f.Start, f.Text, p.target)
	}
	p.run.setLive(cur)
}

func (p *subPipe) flushLive() {
	if p.live == nil {
		return
	}
	p.lmu.Lock()
	out := p.live.Flush()
	p.lmu.Unlock()
	for _, f := range out {
		p.run.finishLine(f.Start, f.Text, p.target)
	}
	p.run.setLive("")
}

func (p *subPipe) liveLoop() {
	fails := 0
	for {
		rec, gen, _, ok := p.current()
		if !ok {
			return
		}
		ss, isStream := rec.(*streamServer)
		if !isStream {
			p.run.fail("the engine for live subtitles is not a streaming engine")
			return
		}
		sess, err := ss.openSegmentSession(context.Background(), p.onSegment)
		if err != nil {
			fails++
			log.Printf("subtitles: the live stream did not open: %v", err)
			if fails >= 3 || engineGone(rec) {
				fails = 0
				p.lost(gen, "the live stream does not open")
			}
			if !p.sleep(time.Second) {
				return
			}
			continue
		}
		fails = 0
		p.smu.Lock()
		p.sess = sess
		p.smu.Unlock()
		log.Printf("subtitles: the live stream is open")
		var berr error
		select {
		case berr = <-sess.broken():
		case <-p.stopCh:
		}
		p.smu.Lock()
		if p.sess == sess {
			p.sess = nil
		}
		p.smu.Unlock()
		sess.close()
		p.flushLive()
		if berr == nil || p.draining.Load() {
			return
		}
		log.Printf("subtitles: the live stream broke off (%v), reconnecting", berr)
		if engineGone(rec) {
			p.lost(gen, "the engine process ended")
		}
		if !p.sleep(time.Second) {
			return
		}
	}
}

func (p *subPipe) stop(drain time.Duration) {
	if drain > 0 {
		if p.target.mode == modePhrases {
			for _, ph := range p.seg.Flush() {
				p.enqueue(ph)
			}
			p.draining.Store(true)
			select {
			case p.qsig <- struct{}{}:
			default:
			}
		} else {
			if drain > subLiveDrain {
				drain = subLiveDrain
			}
			p.draining.Store(true)
			p.smu.Lock()
			s := p.sess
			p.smu.Unlock()
			if s != nil {
				s.end()
			} else {
				drain = 0
			}
		}
		if drain > 0 {
			select {
			case <-p.exitCh:
			case <-time.After(drain):
			}
		}
	}
	close(p.killCh)
	p.mu.Lock()
	p.closed = true
	rec, shared := p.rec, p.shared
	p.rec = nil
	p.mu.Unlock()
	close(p.stopCh)
	p.wg.Wait()
	p.flushLive()
	if rec != nil && !shared {
		rec.stop()
		log.Printf("subtitles: its own engine %s stopped", rec.engine())
	}
}

func (a *App) startSubtitles() {
	a.mu.Lock()
	quitting := a.quitting
	a.mu.Unlock()
	if quitting {
		return
	}
	a.subMu.Lock()
	defer a.subMu.Unlock()
	if a.sub != nil {
		log.Printf("subtitles: already running")
		return
	}
	r, err := startSubtitlesRun(a, a.snapshot(), subOptions{onEnd: func(r *subRun) {
		a.subMu.Lock()
		if a.sub == r {
			a.sub = nil
		}
		a.subMu.Unlock()
	}})
	if err != nil {
		log.Printf("subtitles: did not start: %s", logText(err))
		return
	}
	a.sub = r
}

func (a *App) stopSubtitles(why string, drain time.Duration) bool {
	a.subMu.Lock()
	r := a.sub
	a.sub = nil
	a.subMu.Unlock()
	if r == nil {
		return false
	}
	r.stop(why, drain)
	select {
	case <-r.done():
	case <-time.After(drain + 15*time.Second):
		log.Printf("subtitles: the run did not finish in time")
	}
	return true
}

func (a *App) subtitlesCommand(action string) {
	switch action {
	case "start":
		a.startSubtitles()
	case "stop":
		if !a.stopSubtitles("stopped by the user", subStopDrain) {
			log.Printf("subtitles: not running")
		}
	case "toggle":
		if !a.stopSubtitles("stopped by the user", subStopDrain) {
			a.startSubtitles()
		}
	}
}

func (a *App) subtitlesConfigChanged() {
	a.subMu.Lock()
	r := a.sub
	a.subMu.Unlock()
	if r == nil {
		return
	}
	t, err := resolveSubtitles(a.snapshot())
	if err != nil {
		r.fail(fmt.Sprintf("the new settings leave subtitles without a model: %s", logText(err)))
		return
	}
	if cur := r.target(); cur != nil && !cur.same(t) {
		r.restart(t)
	}
}

func runSubcheck(args []string) {
	seconds := 20
	if len(args) > 0 {
		if n, err := strconv.Atoi(args[0]); err == nil && n > 0 {
			seconds = n
		}
	}
	cfg, err := loadConfig("config.json")
	if err != nil {
		log.Printf("config: %v", err)
		return
	}
	if len(args) > 1 {
		if !validPresetLang(args[1]) {
			log.Printf("subcheck: unknown language %q", args[1])
			return
		}
		cfg.SubtitlesLanguage = args[1]
	}
	app := &App{cfg: cfg, enabled: true, evq: evqueue.New[ptEvent](8)}
	log.Printf("subcheck: listening to the computer's sound for %d s", seconds)
	r, err := startSubtitlesRun(app, cfg, subOptions{logLive: true})
	if err != nil {
		log.Printf("subtitles: did not start: %s", logText(err))
		return
	}
	select {
	case <-time.After(time.Duration(seconds) * time.Second):
		r.stop("the check time is over", subStopDrain)
	case <-r.done():
	}
	<-r.done()
	log.Printf("subcheck: done")
}
