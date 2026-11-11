package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"holdtotype/internal/evqueue"
	"holdtotype/internal/livetail"
	"holdtotype/internal/ovplace"
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
	subPendingBytes  = 10 * sampleRate * 2
	subDeviceGrace   = 3 * time.Second
	subDeviceGiveUp  = 30 * time.Second
	subBoardKeep     = 50
	subRecoverWindow = 2 * time.Minute
	subRecoverMax    = 3

	subTranslateStop   = 3 * time.Second
	subTranslateLogGap = 10 * time.Second
)

var errSubStopped = errors.New("subtitles are stopping")

const subtitlesComboID = "subtitles"

const (
	subErrNoModel = "nomodel"
	subErrEngine  = "engine"
	subErrDevice  = "device"
	subErrLost    = "lost"
)

type subError struct {
	kind string
	msg  string
}

func (e *subError) Error() string { return e.msg }

func subErrorKind(err error) string {
	var se *subError
	if errors.As(err, &se) {
		return se.kind
	}
	return subErrEngine
}

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
		return nil, &subError{subErrNoModel, fmt.Sprintf("no model is chosen for the language %s", lang)}
	}
	external := m.Engine == engineWhisper && (c.SttSource == "remote" || !c.ServerAutostart)
	if !external && !m.installed() {
		return nil, &subError{subErrNoModel, fmt.Sprintf("the model %s for the language %s is not installed", m.NameKey, lang)}
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
	failKind  string
	engName   string
	engPort   int
	engShared bool
	engModel  string
	engPid    int

	cores coreDivider

	fileMu   sync.Mutex
	file     *os.File
	filePath string
	fileErr  bool
	fileDone bool

	trMu     sync.Mutex
	tr       *subTranslator
	trClosed bool
	trSeen   string
	trLang   string
	trShown  atomic.Bool

	phMu     sync.Mutex
	phraser  *subtitles.Phraser
	outMu    sync.Mutex
	outOrder []int
	outDone  map[int]subPieceOut

	skipped atomic.Int64
	muted   atomic.Int64

	liveMu   sync.Mutex
	liveLast time.Time
	liveText string

	devBad    time.Time
	devTried  time.Time
	devLogged bool

	phase atomic.Int32
}

func startSubtitlesRun(a *App, cfg *Config, opts subOptions) (*subRun, error) {
	t, err := resolveSubtitles(cfg)
	if err != nil {
		return nil, err
	}
	lb, err := openLoopback()
	if err != nil {
		return nil, &subError{subErrDevice, err.Error()}
	}
	now := time.Now()
	r := &subRun{
		app:       a,
		opts:      opts,
		started:   now,
		cap:       lb,
		pacer:     subtitles.NewPacer(sampleRate, now),
		board:     subtitles.NewBoard(subBoardKeep),
		phraser:   subtitles.NewPhraser(),
		outDone:   map[int]subPieceOut{},
		stopCh:    make(chan struct{}),
		doneCh:    make(chan struct{}),
		restartCh: make(chan *subTarget, 1),
		want:      t,
	}
	log.Printf("subtitles: started: language %s, model %s, mode %s", t.lang, t.name, t.mode)
	r.syncTranslation(cfg, t)
	r.refreshShown()
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

func (r *subRun) failure() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.failKind
}

func (r *subRun) setPhase(p int) { r.phase.Store(int32(p)) }

func (r *subRun) windowPhase() int { return int(r.phase.Load()) }

func (r *subRun) windowRows() []subtitles.Row {
	if r.trShown.Load() {
		return r.board.TranslationRows(swTextRunes)
	}
	return r.board.Tail(swTextRunes).Rows()
}

func (r *subRun) refreshShown() {
	t := r.target()
	r.trShown.Store(t != nil && translationTarget(r.app.snapshot(), t) != "")
}

func (r *subRun) fail(kind, msg string) {
	r.mu.Lock()
	first := r.failMsg == ""
	if first {
		r.failMsg = msg
		r.failKind = kind
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
	r.engPid = recognizerPid(rec)
	r.mu.Unlock()
	r.syncCores()
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
			r.hear(pipe.target, func(ph *subtitles.Phraser, at time.Time) []subtitles.TextPhrase { return ph.Tick(at) })
		case <-dev.C:
			r.checkDevice(time.Now())
			r.refreshShown()
			r.syncCores()
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
		r.fail(subErrDevice, errNoPlayback.Error())
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

func (r *subRun) cleanText(cfg *Config, t *subTarget, text string) string {
	text = strings.TrimSpace(segments.Join(strings.TrimSpace(text)))
	if text == "" || onlyNoise.MatchString(text) {
		return ""
	}
	return strings.TrimSpace(replace.Apply(replace.ForLang(cfg.Replacements, t.lang), text))
}

func (r *subRun) finishLine(at time.Time, text string, t *subTarget) {
	cfg := r.app.snapshot()
	text = r.cleanText(cfg, t, text)
	if text == "" {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	ln := r.board.Add(at, text)
	log.Printf("subtitles: line %q", text)
	if cfg.SubtitlesSave && translationTarget(cfg, t) == "" {
		r.appendTranscript(t, "", subtitles.Entry(ln.Start, ln.Text))
	}
}

func (r *subRun) hear(t *subTarget, cut func(*subtitles.Phraser, time.Time) []subtitles.TextPhrase) {
	r.phMu.Lock()
	defer r.phMu.Unlock()
	out := cut(r.phraser, time.Now())
	r.board.SetHearing(r.phraser.Waiting() > 0)
	for _, ph := range out {
		r.phraseCut(ph, t)
	}
}

func (r *subRun) phraseCut(ph subtitles.TextPhrase, t *subTarget) {
	cfg := r.app.snapshot()
	text := r.cleanText(cfg, t, ph.Text)
	if text == "" {
		return
	}
	if r.opts.logLive {
		log.Printf("subtitles: phrase %d, %d words, cut %.1f s after its first word: %q", ph.Seq, ph.Words, ph.Heard.Sub(ph.Start).Seconds(), text)
	}
	to := translationTarget(cfg, t)
	if to == "" {
		return
	}
	if w := r.translator(); w != nil {
		w.push(ph, text, to)
	}
}

func (r *subRun) appendTranscript(t *subTarget, to string, entries ...string) {
	r.fileMu.Lock()
	defer r.fileMu.Unlock()
	if r.fileErr || r.fileDone {
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
		if _, err := f.WriteString(subtitles.Opening(r.started, t.lang, t.name, to)); err != nil {
			log.Printf("subtitles: writing the transcript: %v", err)
		}
		r.file = f
		r.filePath = p
		if abs, aerr := filepath.Abs(p); aerr == nil {
			r.filePath = abs
		}
		log.Printf("subtitles: transcript %s", r.filePath)
	}
	for _, e := range entries {
		if _, err := r.file.WriteString(e); err != nil {
			log.Printf("subtitles: writing the transcript: %v", err)
			return
		}
	}
	_ = r.file.Sync()
}

type subPieceOut struct {
	at          time.Time
	source      string
	translation string
	target      string
}

func (r *subRun) pieceQueued(seq int) {
	r.outMu.Lock()
	r.outOrder = append(r.outOrder, seq)
	r.outMu.Unlock()
}

func (r *subRun) piecesJoined(seqs []int) {
	r.outMu.Lock()
	defer r.outMu.Unlock()
	gone := map[int]bool{}
	for _, s := range seqs[1:] {
		gone[s] = true
	}
	kept := r.outOrder[:0]
	for _, s := range r.outOrder {
		if !gone[s] {
			kept = append(kept, s)
		}
	}
	r.outOrder = kept
}

func (r *subRun) pieceFinished(seq int, out subPieceOut) {
	r.outMu.Lock()
	defer r.outMu.Unlock()
	r.outDone[seq] = out
	cfg := r.app.snapshot()
	t := r.target()
	for len(r.outOrder) > 0 {
		e, ok := r.outDone[r.outOrder[0]]
		if !ok {
			break
		}
		delete(r.outDone, r.outOrder[0])
		r.outOrder = r.outOrder[1:]
		if !cfg.SubtitlesSave || t == nil {
			continue
		}
		entries := []string{subtitles.Entry(e.at, e.source)}
		if e.translation != "" {
			entries = append(entries, subtitles.TranslationEntry(e.at, e.translation))
		}
		r.appendTranscript(t, e.target, entries...)
	}
}

func (r *subRun) finish() {
	r.cap.close()
	r.trMu.Lock()
	w := r.tr
	r.trClosed = true
	r.trMu.Unlock()
	trSummary := ""
	if w != nil {
		w.stop()
		trSummary = w.summary()
	}
	r.cores.close()
	r.board.SetHearing(false)
	r.fileMu.Lock()
	if r.file != nil {
		_ = r.file.Close()
	}
	r.fileDone = true
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
	log.Printf("subtitles: summary: %.1f s captured, %.1f s of silence added, topping up paused %d time(s) after %.0f s of silence, %d line(s), %.1f s skipped, %.1f s of the app's own sounds silenced, engine %s on port %d with %s (%s), %d device open(s), transcript %s%s",
		r.pacer.Captured().Seconds(), r.pacer.Added().Seconds(), r.pacer.Pauses(), subtitles.SilenceLimit.Seconds(),
		r.board.Count(), time.Duration(r.skipped.Load()).Seconds(), float64(r.muted.Load())/sampleRate, name, port, model, how, r.cap.opens.Load(), path, trSummary)
	if r.opts.onEnd != nil {
		r.opts.onEnd(r)
	}
}

func translationTarget(cfg *Config, t *subTarget) string {
	to := cfg.SubtitlesTranslate
	if to == "" || to == t.lang || !translationSource(cfg).ok() {
		return ""
	}
	return to
}

type subTrState struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

func subtitlesTranslator(cfg *Config) subTrState {
	src := translationSource(cfg)
	return subTrState{Kind: src.kind, Name: src.name}
}

func (r *subRun) syncTranslation(cfg *Config, t *subTarget) {
	to := cfg.SubtitlesTranslate
	r.trMu.Lock()
	prev, prevLang := r.trSeen, r.trLang
	r.trSeen, r.trLang = to, t.lang
	r.trMu.Unlock()
	if to == prev && (to == "" || t.lang == prevLang) {
		return
	}
	if to == "" {
		log.Printf("subtitles: translation is off from now on")
		return
	}
	name := subtitles.LanguageName(to)
	if to == t.lang {
		log.Printf("subtitles: translation into %s: the speech language is the same, nothing is translated", name)
		return
	}
	src := translationSource(cfg)
	if !src.ok() {
		log.Printf("subtitles: translation into %s: no editor model is installed and no post-processing server is set, the text stays untranslated", name)
		return
	}
	log.Printf("subtitles: translation into %s by %s", name, src.label())
	if r.translator() == nil || src.kind != editorLocal {
		return
	}
	f := r.app.llmWarmUp()
	go func() {
		<-f.done
		if f.err != nil {
			log.Printf("subtitles: the editor model did not start: %s", logText(f.err))
			return
		}
		r.syncCores()
		r.primeTranslation(f.srv, to)
	}()
}

func (r *subRun) primeTranslation(llm *llamaServer, to string) {
	if llm == nil {
		return
	}
	started := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), subtitles.TranslateTimeout)
	defer cancel()
	if _, err := llm.chatStream(ctx, subtitles.TranslateMessages(to, "Hello."), 1, subtitles.TranslateGrammar(to), nil); err != nil {
		log.Printf("subtitles: the translation request did not load: %s", logText(err))
		return
	}
	log.Printf("subtitles: the translation request is loaded into the editor model in %.1f s", time.Since(started).Seconds())
}

func (r *subRun) translator() *subTranslator {
	r.trMu.Lock()
	defer r.trMu.Unlock()
	if r.tr == nil && !r.trClosed {
		r.tr = newSubTranslator(r)
	}
	return r.tr
}

type logGap struct {
	last  time.Time
	quiet int
}

func (g *logGap) pass(now time.Time) (bool, int) {
	if !g.last.IsZero() && now.Sub(g.last) < subTranslateLogGap {
		g.quiet++
		return false, 0
	}
	n := g.quiet
	g.quiet = 0
	g.last = now
	return true, n
}

type subTranslator struct {
	run    *subRun
	ctx    context.Context
	cancel context.CancelFunc
	stopCh chan struct{}
	exitCh chan struct{}
	sig    chan struct{}

	mu      sync.Mutex
	q       *subtitles.TranslateQueue
	closed  bool
	target  string
	used    string
	failLog logGap
	dropLog logGap
	waitLog bool
	delays  []time.Duration

	pieces     atomic.Int64
	translated atomic.Int64
	heard      atomic.Int64
	dropped    atomic.Int64
}

func newSubTranslator(r *subRun) *subTranslator {
	ctx, cancel := context.WithCancel(context.Background())
	w := &subTranslator{
		run:    r,
		ctx:    ctx,
		cancel: cancel,
		stopCh: make(chan struct{}),
		exitCh: make(chan struct{}),
		sig:    make(chan struct{}, 1),
		q:      subtitles.NewTranslateQueue(subtitles.TranslateWaiting, subtitles.TranslateMaxWait, subtitles.TranslateMerge),
	}
	go w.loop()
	return w
}

func (w *subTranslator) stopped() bool {
	select {
	case <-w.stopCh:
		return true
	default:
		return false
	}
}

func (w *subTranslator) push(ph subtitles.TextPhrase, text, to string) {
	r := w.run
	w.mu.Lock()
	if w.closed {
		w.mu.Unlock()
		return
	}
	w.target = to
	r.board.AddPiece(ph.Seq, ph.Para, text)
	r.pieceQueued(ph.Seq)
	dropped := w.q.Push(subtitles.Pending{Seq: ph.Seq, Text: text, Words: len(strings.Fields(text)), Target: to, Para: ph.Para, Start: ph.Start, Heard: ph.Heard})
	w.mu.Unlock()
	if len(dropped) > 0 {
		w.drop(dropped, "the queue is full")
		w.noteDrop(fmt.Sprintf("translation falls behind, %d waiting phrase(s) are shown as heard", len(dropped)))
	}
	select {
	case w.sig <- struct{}{}:
	default:
	}
}

func (w *subTranslator) drop(ps []subtitles.Pending, why string) {
	r := w.run
	now := time.Now()
	for _, p := range ps {
		w.pieces.Add(1)
		w.dropped.Add(1)
		r.board.SetPiece(p.Seq, "", subtitles.PieceFailed)
		r.pieceFinished(p.Seq, subPieceOut{at: p.Start, source: p.Text, target: p.Target})
		if r.opts.logLive {
			log.Printf("subtitles: piece %d dropped after %.1f s: %s", p.Seq, now.Sub(p.Heard).Seconds(), why)
		}
	}
}

func (w *subTranslator) noteDrop(msg string) {
	w.mu.Lock()
	ok, quiet := w.dropLog.pass(time.Now())
	w.mu.Unlock()
	if !ok {
		return
	}
	if quiet > 0 {
		msg += fmt.Sprintf(" (%d more such case(s) not logged)", quiet)
	}
	log.Printf("subtitles: %s", msg)
}

func (w *subTranslator) noteFail(err error) {
	w.mu.Lock()
	ok, quiet := w.failLog.pass(time.Now())
	w.mu.Unlock()
	if !ok {
		return
	}
	more := ""
	if quiet > 0 {
		more = fmt.Sprintf(" (%d more failure(s) not logged)", quiet)
	}
	log.Printf("subtitles: a piece was not translated: %s%s", logText(err), more)
}

func (w *subTranslator) waitDictation() bool {
	a := w.run.app
	if a.state.Load() != stProcessing || w.stopped() {
		return w.ctx.Err() == nil
	}
	if !w.waitLog {
		w.waitLog = true
		log.Printf("subtitles: a dictation is being processed, the translation waits")
	}
	deadline := time.Now().Add(subDictationWait)
	for a.state.Load() == stProcessing && time.Now().Before(deadline) {
		select {
		case <-time.After(100 * time.Millisecond):
		case <-w.stopCh:
			return w.ctx.Err() == nil
		}
	}
	return w.ctx.Err() == nil
}

func (w *subTranslator) loop() {
	defer close(w.exitCh)
	for {
		select {
		case <-w.sig:
		case <-w.stopCh:
		}
		for {
			if !w.waitDictation() {
				return
			}
			w.mu.Lock()
			b, ok, stale := w.q.Pop(time.Now())
			w.mu.Unlock()
			if len(stale) > 0 {
				w.drop(stale, "it waited too long")
				w.noteDrop(fmt.Sprintf("%d phrase(s) waited more than %.0f s for translation and are shown as heard", len(stale), subtitles.TranslateMaxWait.Seconds()))
			}
			if !ok {
				break
			}
			w.handle(b)
			if w.ctx.Err() != nil {
				return
			}
		}
		if w.stopped() {
			return
		}
	}
}

func (w *subTranslator) fail(b subtitles.Batch, err error) {
	r := w.run
	seq := b.Seq()
	r.board.SetPiece(seq, "", subtitles.PieceFailed)
	r.pieceFinished(seq, subPieceOut{at: b.Start, source: b.Text, target: b.Target})
	if w.ctx.Err() != nil {
		w.dropped.Add(1)
		if r.opts.logLive {
			log.Printf("subtitles: piece %d dropped: the run ended before its translation", seq)
		}
		return
	}
	w.heard.Add(1)
	w.noteFail(err)
	if r.opts.logLive {
		log.Printf("subtitles: piece %d shown as heard after %.1f s: %s", seq, time.Since(b.Heard).Seconds(), logText(err))
	}
}

func (w *subTranslator) handle(b subtitles.Batch) {
	r := w.run
	seq := b.Seq()
	if len(b.Seqs) > 1 {
		r.board.JoinPieces(b.Seqs)
		r.piecesJoined(b.Seqs)
	}
	w.pieces.Add(1)
	cfg := r.app.snapshot()
	src := translationSource(cfg)
	if !src.ok() {
		w.fail(b, errNoEditor)
		return
	}
	w.mu.Lock()
	w.used = src.label()
	w.mu.Unlock()
	r.board.SetPiece(seq, "", subtitles.PieceRunning)
	ctx, cancel := context.WithTimeout(w.ctx, subtitles.TranslateTimeout)
	var first time.Time
	out, err := r.app.editorStream(ctx, cfg, src, subtitles.TranslateMessages(b.Target, b.Text), subtitles.TranslateTokens(b.Words), subtitles.TranslateGrammar(b.Target), func(sofar string) {
		if first.IsZero() {
			first = time.Now()
		}
		r.board.SetPiece(seq, subtitles.StreamingText(sofar), subtitles.PieceRunning)
	})
	timedOut := ctx.Err() == context.DeadlineExceeded
	cancel()
	done := time.Now()
	if err != nil {
		if timedOut {
			err = fmt.Errorf("no answer within %.0f s (%s)", subtitles.TranslateTimeout.Seconds(), logText(err))
		}
		w.fail(b, err)
		return
	}
	text, verdict := subtitles.CleanTranslation(b.Text, out, b.Target)
	switch verdict {
	case subtitles.TrSame:
		text = b.Text
	case subtitles.TrEmpty:
		w.fail(b, errors.New("the answer is empty"))
		return
	case subtitles.TrRamble:
		w.fail(b, fmt.Errorf("an answer of %d characters for a piece of %d is thrown away", len([]rune(out)), len([]rune(b.Text))))
		return
	case subtitles.TrScript:
		w.fail(b, fmt.Errorf("the answer is not written in %s", subtitles.LanguageName(b.Target)))
		return
	}
	r.board.SetPiece(seq, text, subtitles.PieceDone)
	r.pieceFinished(seq, subPieceOut{at: b.Start, source: b.Text, translation: text, target: b.Target})
	w.translated.Add(1)
	w.mu.Lock()
	w.delays = append(w.delays, done.Sub(b.Heard))
	waiting := w.q.Len()
	w.mu.Unlock()
	if r.opts.logLive {
		if first.IsZero() {
			first = done
		}
		same := ""
		if verdict == subtitles.TrSame {
			same = ", the same as heard"
		}
		log.Printf("subtitles: piece %d translated: first word after %.1f s, complete after %.1f s, %d waiting%s: %q", seq, first.Sub(b.Heard).Seconds(), done.Sub(b.Heard).Seconds(), waiting, same, text)
	}
}

func (w *subTranslator) stop() {
	w.mu.Lock()
	w.closed = true
	w.mu.Unlock()
	close(w.stopCh)
	select {
	case <-w.exitCh:
	case <-time.After(subTranslateStop):
		log.Printf("subtitles: the translation did not finish within %.0f s, it is cancelled", subTranslateStop.Seconds())
		w.cancel()
		select {
		case <-w.exitCh:
		case <-time.After(time.Second):
		}
	}
	w.cancel()
	w.mu.Lock()
	rest := w.q.Clear()
	w.mu.Unlock()
	if len(rest) > 0 {
		w.drop(rest, "the run ended")
	}
}

func (w *subTranslator) summary() string {
	w.mu.Lock()
	target, used := w.target, w.used
	delays := append([]time.Duration(nil), w.delays...)
	w.mu.Unlock()
	if used == "" {
		used = "none (no request sent)"
	}
	if target == "" {
		target = "none"
	}
	timing := "none complete"
	if len(delays) > 0 {
		slices.Sort(delays)
		timing = fmt.Sprintf("complete after %.1f s (median), %.1f s (worst)", delays[len(delays)/2].Seconds(), delays[len(delays)-1].Seconds())
	}
	return fmt.Sprintf(", translation into %s: %d pieces, %d translated, %d shown as heard, %d dropped, %s, by %s",
		target, w.pieces.Load(), w.translated.Load(), w.heard.Load(), w.dropped.Load(), timing, used)
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

	smu     sync.Mutex
	sess    *segSession
	pending []byte
	lmu     sync.Mutex
	live    *subtitles.Live
	gate    *subtitles.Gate
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
		p.seg.StartAt(r.pacer.Fed())
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
	p.run.setPhase(subPhaseLoading)
	rec, shared, err := p.acquire()
	if err != nil {
		if !p.stopped() {
			p.run.fail(subErrorKind(err), fmt.Sprintf("the engine did not start: %s", logText(err)))
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
	p.run.setPhase(subPhaseListening)
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
	p.run.setPhase(subPhaseRestarting)
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
		p.run.fail(subErrLost, fmt.Sprintf("the engine went away %d times in %.0f minutes", len(p.recovered)+1, subRecoverWindow.Minutes()))
		return
	}
	p.recovered = append(p.recovered, now)
	if !p.sleep(time.Second) {
		return
	}
	rec, sh, err := p.acquire()
	if err != nil {
		if !p.stopped() {
			p.run.fail(subErrEngine, fmt.Sprintf("the engine could not be brought back: %s", logText(err)))
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
		if s == nil {
			p.pending = append(p.pending, pcm...)
			if over := len(p.pending) - subPendingBytes; over > 0 {
				n := copy(p.pending, p.pending[over:])
				p.pending = p.pending[:n]
			}
		}
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
		from, to := ph.At, ph.At.Add(ph.Duration(sampleRate))
		p.run.hear(p.target, func(x *subtitles.Phraser, at time.Time) []subtitles.TextPhrase {
			return x.Whole(segments.Join(text), from, to, at)
		})
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
	p.run.hear(p.target, func(ph *subtitles.Phraser, at time.Time) []subtitles.TextPhrase {
		return ph.Update(seg, text, final, at)
	})
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
	p.run.hear(p.target, func(ph *subtitles.Phraser, at time.Time) []subtitles.TextPhrase { return ph.Flush(at) })
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
			p.run.fail(subErrEngine, "the engine for live subtitles is not a streaming engine")
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
		kept := len(p.pending)
		if kept > 0 {
			_ = sess.push(p.pending)
			p.pending = p.pending[:0]
		}
		p.smu.Unlock()
		log.Printf("subtitles: the live stream is open, %.1f s of sound kept while it opened", float64(kept)/(2*sampleRate))
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

func (a *App) subFlash(kind string) {
	key := "sub.err." + kind
	if _, ok := msgs["en"][key]; !ok {
		key = "sub.err." + subErrEngine
	}
	if a.snapshot().Overlay {
		overlaySet(ovFlashErr, tr(key))
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
	cfg := a.snapshot()
	r, err := startSubtitlesRun(a, cfg, subOptions{onEnd: func(r *subRun) {
		a.subMu.Lock()
		if a.sub == r {
			a.sub = nil
		}
		a.subMu.Unlock()
		subWindowHide(r)
		if kind := r.failure(); kind != "" {
			a.subFlash(kind)
		}
	}})
	if err != nil {
		log.Printf("subtitles: did not start: %s", logText(err))
		a.subFlash(subErrorKind(err))
		return
	}
	a.sub = r
	subWindowConfigure(cfg.SubtitlesLines, cfg.SubtitlesSize, cfg.SubtitlesXY)
	subWindowShow(r, r.target().lang, func() {
		a.stopSubtitles("closed from its window", subStopDrain)
	}, a.saveSubtitlesPlace)
}

func (a *App) saveSubtitlesPlace(key string, f ovplace.Frac) {
	a.mu.Lock()
	c := *a.cfg
	places := map[string]ovplace.Frac{}
	for k, v := range c.SubtitlesXY {
		places[k] = v
	}
	places[key] = f
	c.SubtitlesXY = places
	a.mu.Unlock()
	if err := saveConfig("config.json", &c); err != nil {
		log.Printf("subtitles: the window place was not saved: %v", err)
		return
	}
	a.mu.Lock()
	cur := *a.cfg
	cur.SubtitlesXY = places
	a.cfg = &cur
	a.mu.Unlock()
}

func (a *App) subtitlesKey() {
	if !a.subBusy.CompareAndSwap(false, true) {
		log.Printf("subtitles: the shortcut came while subtitles are starting or stopping, ignored")
		return
	}
	go func() {
		defer a.subBusy.Store(false)
		a.subtitlesCommand("toggle")
	}()
}

func (a *App) subtitlesRunning() bool {
	a.subMu.Lock()
	defer a.subMu.Unlock()
	return a.sub != nil
}

func (a *App) stopSubtitles(why string, drain time.Duration) bool {
	a.subMu.Lock()
	r := a.sub
	a.sub = nil
	a.subMu.Unlock()
	if r == nil {
		return false
	}
	subWindowHide(r)
	a.subMu.Lock()
	a.subEnding = append(a.subEnding, r)
	a.subMu.Unlock()
	r.stop(why, drain)
	select {
	case <-r.done():
	case <-time.After(drain + 15*time.Second):
		log.Printf("subtitles: the run did not finish in time")
	}
	a.subMu.Lock()
	for i, x := range a.subEnding {
		if x == r {
			a.subEnding = append(a.subEnding[:i], a.subEnding[i+1:]...)
			break
		}
	}
	a.subMu.Unlock()
	return true
}

func (a *App) closeSubtitles() {
	a.stopSubtitles("the app is closing", subExitDrain)
	a.subMu.Lock()
	ending := append([]*subRun(nil), a.subEnding...)
	a.subMu.Unlock()
	deadline := time.After(subExitDrain + 2*time.Second)
	for _, r := range ending {
		r.stop("the app is closing", subExitDrain)
		select {
		case <-r.done():
		case <-deadline:
			log.Printf("subtitles: a stopping run did not finish before the exit")
			return
		}
	}
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
	cfg := a.snapshot()
	subWindowConfigure(cfg.SubtitlesLines, cfg.SubtitlesSize, cfg.SubtitlesXY)
	t, err := resolveSubtitles(cfg)
	if err != nil {
		r.fail(subErrorKind(err), fmt.Sprintf("the new settings leave subtitles without a model: %s", logText(err)))
		return
	}
	if cur := r.target(); cur != nil && !cur.same(t) {
		r.restart(t)
		subWindowSetLang(r, t.lang)
	}
	r.syncTranslation(cfg, t)
	r.refreshShown()
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
	if len(args) > 2 {
		if !validTranslateLang(args[2]) {
			log.Printf("subcheck: unknown translation language %q", args[2])
			return
		}
		cfg.SubtitlesTranslate = args[2]
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
	app.llmShutdown()
	log.Printf("subcheck: done")
}

type subModelInfo struct {
	Lang      string `json:"lang"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Engine    string `json:"engine"`
	Mode      string `json:"mode"`
	Installed bool   `json:"installed"`
	Whisper   bool   `json:"whisper"`
	Remote    bool   `json:"remote"`
}

func (a *App) subtitlesModelInfo(lang string) string {
	cfg := a.snapshot()
	if !validPresetLang(lang) {
		lang = cfg.SubtitlesLanguage
	}
	c := *cfg
	c.Language = lang
	info := subModelInfo{Lang: lang}
	if m := activeModel(&c); m != nil {
		remote := m.Engine == engineWhisper && c.SttSource == "remote"
		external := m.Engine == engineWhisper && (remote || !c.ServerAutostart)
		info.ID = m.ID
		info.Name = m.NameKey
		info.Engine = m.Engine
		info.Mode = modePhrases
		if m.Engine == engineStream {
			info.Mode = modeLive
		}
		info.Installed = external || m.installed()
		info.Whisper = m.Engine == engineWhisper
		info.Remote = remote
	}
	out, _ := json.Marshal(info)
	return string(out)
}

func openTranscriptsFolder() {
	if err := os.MkdirAll(subtitles.Folder, 0o755); err != nil {
		log.Printf("subtitles: the transcript folder: %v", err)
		return
	}
	if abs, err := filepath.Abs(subtitles.Folder); err == nil {
		shellOpenURL(abs)
	}
}

type subState struct {
	Phase string `json:"phase"`
	Mode  string `json:"mode"`
	Model string `json:"model"`
	Lang  string `json:"lang"`
}

func (a *App) subtitlesState() subState {
	a.subMu.Lock()
	r := a.sub
	a.subMu.Unlock()
	if r == nil {
		return subState{Phase: "off"}
	}
	t := r.target()
	st := subState{Mode: t.mode, Model: t.name, Lang: t.lang}
	switch r.windowPhase() {
	case subPhaseLoading:
		st.Phase = "loading"
	case subPhaseRestarting:
		st.Phase = "restarting"
	default:
		st.Phase = "listening"
	}
	return st
}
