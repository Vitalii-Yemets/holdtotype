package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/sys/windows"

	"holdtotype/internal/fastpack"
	"holdtotype/internal/subtitles"
	"holdtotype/internal/trmodel"
)

const (
	fastExeName  = "fast-translator.exe"
	fastLangWait = 30 * time.Second
)

type fastServer struct {
	label  string
	cmd    *exec.Cmd
	job    uintptr
	in     io.WriteCloser
	lines  chan string
	readyC chan struct{}

	mu     sync.Mutex
	exited bool

	reqMu  sync.Mutex
	nextID int
	loaded map[string]bool
}

func (s *fastServer) alive() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return !s.exited
}

func (s *fastServer) pid() int { return cmdPid(s.cmd) }

func (s *fastServer) stop() {
	if s.cmd != nil && s.cmd.Process != nil {
		_ = s.cmd.Process.Kill()
	}
	if s.job != 0 {
		_ = windows.CloseHandle(windows.Handle(s.job))
		s.job = 0
	}
}

func (s *fastServer) waitReady(timeout time.Duration) error {
	deadline := time.After(timeout)
	for {
		select {
		case <-s.readyC:
			return nil
		case <-deadline:
			return fmt.Errorf("%s %w: nothing within %s", s.label, errEngineTimeout, timeout)
		case <-time.After(200 * time.Millisecond):
			if !s.alive() {
				return fmt.Errorf("%s %w", s.label, errEngineExited)
			}
		}
	}
}

func fastExePath(cfg *Config) string {
	return filepath.Join(filepath.Dir(absPath(cfg.LLMExe)), fastExeName)
}

func startFastServer(spec engineSpec, logw io.Writer) (engineProc, error) {
	s := &fastServer{label: spec.label, lines: make(chan string, 16), readyC: make(chan struct{}), loaded: map[string]bool{}}
	cmd := exec.Command(spec.exe)
	cmd.Stderr = logw
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	in, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("starting %s: %w", filepath.Base(spec.exe), err)
	}
	s.cmd, s.in = cmd, in
	s.job = uintptr(attachProcessToJob(cmd.Process.Pid))
	go func() {
		r := bufio.NewReaderSize(out, 1<<16)
		once := false
		for {
			line, err := r.ReadString('\n')
			if line = strings.TrimSpace(line); line != "" {
				if !once && strings.Contains(line, `"ready"`) {
					once = true
					close(s.readyC)
				} else {
					s.lines <- line
				}
			}
			if err != nil {
				break
			}
		}
		close(s.lines)
	}()
	go func() {
		_ = cmd.Wait()
		s.mu.Lock()
		s.exited = true
		s.mu.Unlock()
	}()
	return s, nil
}

type fastReply struct {
	ID    int     `json:"id"`
	OK    bool    `json:"ok"`
	Text  string  `json:"text"`
	Error string  `json:"error"`
	MS    float64 `json:"ms"`
}

func (s *fastServer) call(ctx context.Context, req map[string]any) (fastReply, error) {
	s.reqMu.Lock()
	defer s.reqMu.Unlock()
	return s.callLocked(ctx, req)
}

func (s *fastServer) callLocked(ctx context.Context, req map[string]any) (fastReply, error) {
	s.nextID++
	id := s.nextID
	req["id"] = id
	raw, _ := json.Marshal(req)
	if _, err := s.in.Write(append(raw, '\n')); err != nil {
		return fastReply{}, fmt.Errorf("%s: %w", s.label, err)
	}
	for {
		select {
		case line, ok := <-s.lines:
			if !ok {
				return fastReply{}, fmt.Errorf("%s %w", s.label, errEngineExited)
			}
			var r fastReply
			if err := json.Unmarshal([]byte(line), &r); err != nil {
				continue
			}
			if r.ID < id {
				continue
			}
			if !r.OK {
				return r, fmt.Errorf("%s: %s", s.label, r.Error)
			}
			return r, nil
		case <-ctx.Done():
			return fastReply{}, ctx.Err()
		}
	}
}

func (s *fastServer) load(ctx context.Context, packs []fastpack.Pack) error {
	s.reqMu.Lock()
	defer s.reqMu.Unlock()
	for _, p := range packs {
		if s.loaded[p.Key()] {
			continue
		}
		started := time.Now()
		_, err := s.callLocked(ctx, map[string]any{
			"op": "load", "from": p.From, "to": p.To,
			"model": p.Path("models", "model"), "vocab": p.Path("models", "vocab"), "lex": p.Path("models", "lex"),
		})
		if err != nil {
			return err
		}
		s.loaded[p.Key()] = true
		log.Printf("fast-translator: language pack %s loaded in %.1f s", p.Key(), time.Since(started).Seconds())
	}
	return nil
}

func (s *fastServer) has(packs []fastpack.Pack) bool {
	s.reqMu.Lock()
	defer s.reqMu.Unlock()
	for _, p := range packs {
		if !s.loaded[p.Key()] {
			return false
		}
	}
	return len(packs) > 0
}

func (s *fastServer) translate(ctx context.Context, from, to, text string, html bool) (string, error) {
	r, err := s.call(ctx, map[string]any{"op": "translate", "from": from, "to": to, "text": text, "html": html})
	return r.Text, err
}

func fastEngine(cfg *Config) (engineSpec, error) {
	exe := fastExePath(cfg)
	if _, err := os.Stat(exe); err != nil {
		return engineSpec{}, fmt.Errorf("%s is missing next to the program", fastExeName)
	}
	return engineSpec{label: "fast-translator", exe: exe, wait: 30 * time.Second, launch: startFastServer}, nil
}

func asFast(p engineProc, err error) (*fastServer, error) {
	if err != nil {
		return nil, err
	}
	s, ok := p.(*fastServer)
	if !ok || s == nil {
		return nil, errors.New("the engine is not the light translation engine")
	}
	return s, nil
}

func (a *App) fastWait(ctx context.Context) (*fastServer, error) {
	return asFast(a.waitEngine(ctx, &a.fastEng, fastEngine, "the light translation engine is still starting"))
}

func (a *App) fastShutdown() {
	if a.fastEng.shutdown() {
		log.Printf("fast-translator stopped")
	}
}

type fastRunState struct {
	mu       sync.Mutex
	active   bool
	lang     string
	lost     bool
	waitFrom time.Time
	down     string
}

var fastRun fastRunState

func (f *fastRunState) begin() {
	f.mu.Lock()
	f.active, f.lang, f.lost, f.waitFrom, f.down = true, "", false, time.Time{}, ""
	f.mu.Unlock()
}

func (f *fastRunState) end() {
	f.mu.Lock()
	f.active = false
	f.lang = ""
	f.waitFrom = time.Time{}
	f.mu.Unlock()
}

func (f *fastRunState) setLang(lang string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	changed := f.lang != lang
	f.lang = lang
	if lang != "" {
		f.lost = false
	}
	return changed
}

func (f *fastRunState) snapshot() (lang string, lost bool, down string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.lang, f.lost, f.down
}

func (f *fastRunState) setDown(why string) {
	f.mu.Lock()
	f.down = why
	f.mu.Unlock()
}

func fastSpeechLang(cfg *Config) string {
	if l := strings.ToLower(strings.TrimSpace(cfg.SubtitlesLanguage)); l != "" && l != "auto" {
		return l
	}
	lang, _, _ := fastRun.snapshot()
	return lang
}

var fastChecks sync.Map

func packState(p fastpack.Pack) fastpack.State {
	if st, ok := fastpack.Known("models", p); ok {
		fastChecks.Delete(p.Key())
		return st
	}
	if v, ok := fastChecks.Load(p.Key()); ok {
		if st, done := v.(fastpack.State); done && st != "" {
			return st
		}
		return fastpack.Checking
	}
	fastChecks.Store(p.Key(), fastpack.State(""))
	go func() {
		st, err := fastpack.Check("models", p)
		if err != nil {
			log.Printf("language pack %s cannot be checked: %v", p.Key(), err)
			st = fastpack.Bad
		}
		fastChecks.Store(p.Key(), st)
	}()
	return fastpack.Checking
}

func routeReady(packs []fastpack.Pack) (bool, string) {
	for _, p := range packs {
		if packState(p) != fastpack.Ready {
			return false, p.Key()
		}
	}
	return len(packs) > 0, ""
}

var fastWaitForGuess atomic.Bool

func installedInto(to string) []string {
	var out []string
	for _, from := range fastpack.Languages {
		if from == to {
			continue
		}
		if packs, ok := fastpack.Route(from, to); ok {
			if ready, _ := routeReady(packs); ready {
				out = append(out, from)
			}
		}
	}
	return out
}

type fastStatus struct {
	ok      bool
	missing string
	detail  string
	lang    string
	packs   []fastpack.Pack
}

func fastAvail(cfg *Config) fastStatus {
	to := cfg.SubtitlesTranslate
	if to == "" {
		return fastStatus{missing: subtitles.MissingPack}
	}
	if _, err := os.Stat(fastExePath(cfg)); err != nil {
		return fastStatus{missing: subtitles.MissingFast, detail: "exe"}
	}
	lang, lost, down := fastRun.snapshot()
	if down != "" {
		return fastStatus{missing: subtitles.MissingFast, detail: down}
	}
	if l := strings.ToLower(strings.TrimSpace(cfg.SubtitlesLanguage)); l != "" && l != "auto" {
		lang = l
	}
	if lang == "" {
		installed := installedInto(to)
		switch {
		case len(installed) == 1 && !fastWaitForGuess.Load():
			lang = installed[0]
		case lost:
			return fastStatus{missing: subtitles.MissingLang}
		case len(installed) > 1 || (len(installed) == 1 && fastWaitForGuess.Load()):
			return fastStatus{ok: true}
		default:
			return fastStatus{missing: subtitles.MissingPack, detail: "en-" + to}
		}
	}
	if lang == to {
		return fastStatus{ok: true, lang: lang}
	}
	packs, ok := fastpack.Route(lang, to)
	if !ok {
		return fastStatus{missing: subtitles.MissingPack, detail: lang + "-" + to, lang: lang}
	}
	if ready, missing := routeReady(packs); !ready {
		return fastStatus{missing: subtitles.MissingPack, detail: missing, lang: lang, packs: packs}
	}
	return fastStatus{ok: true, lang: lang, packs: packs}
}

func routeName(lang, to string) string {
	if lang == "" {
		return "speech language not known yet → " + to
	}
	if lang != to && lang != fastpack.Pivot && to != fastpack.Pivot {
		if _, direct := fastpack.Find(lang, to); !direct {
			return lang + " → en → " + to
		}
	}
	return lang + " → " + to
}

func (a *App) fastWarm() bool {
	s, ok := a.fastEng.live().(*fastServer)
	if !ok || s == nil {
		return false
	}
	st := fastAvail(a.snapshot())
	if !st.ok || st.lang == "" {
		return false
	}
	if len(st.packs) == 0 {
		return true
	}
	return s.has(st.packs)
}

func (a *App) warmFast(ctx context.Context) error {
	for {
		st := fastAvail(a.snapshot())
		if !st.ok {
			return &fastFail{st.missing, st.detail}
		}
		if st.lang != "" {
			break
		}
		fastRun.mu.Lock()
		if fastRun.waitFrom.IsZero() {
			fastRun.waitFrom = time.Now()
		}
		_, state := subtitles.SpeechLanguage(a.snapshot().SubtitlesLanguage, fastRun.lang, time.Since(fastRun.waitFrom), fastLangWait)
		if state == subtitles.SpeechLost {
			fastRun.lost = true
			fastRun.mu.Unlock()
			log.Printf("subtitles: the speech language was not recognized within %.0f s, the light engine cannot translate in this run", fastLangWait.Seconds())
			return &fastFail{subtitles.MissingLang, ""}
		}
		fastRun.mu.Unlock()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(250 * time.Millisecond):
		}
	}
	f := a.warmEngine(&a.fastEng, fastEngine)
	<-f.done
	s, err := asFast(f.srv, f.err)
	if err != nil {
		return &fastFail{subtitles.MissingFast, logText(err)}
	}
	st := fastAvail(a.snapshot())
	if !st.ok {
		return &fastFail{st.missing, st.detail}
	}
	if len(st.packs) == 0 {
		return nil
	}
	if err := s.load(ctx, st.packs); err != nil {
		return &fastFail{subtitles.MissingFast, logText(err)}
	}
	return nil
}

type fastFail struct {
	missing string
	detail  string
}

func (f *fastFail) Error() string {
	switch f.missing {
	case subtitles.MissingLang:
		return "the speech language is not known"
	case subtitles.MissingPack:
		return "the language pack " + f.detail + " is not downloaded"
	}
	return "the light translation engine did not start: " + f.detail
}

type fastCounts struct {
	alone, split, fallback, trimmed atomic.Int64
}

const (
	fastAlone    = "alone"
	fastSplit    = "split"
	fastTrimmed  = "trimmed"
	fastFallback = "fallback"
)

func fastPiece(ctx context.Context, srv *fastServer, from, to, prev, shown, text string, para bool) (string, string, string, error) {
	if subtitles.BeginsSentence(prev, para, text) {
		out, err := srv.translate(ctx, from, to, text, false)
		return out, fastAlone, "", err
	}
	out, err := srv.translate(ctx, from, to, subtitles.TagPrevious(prev, text), true)
	if err != nil {
		return "", "", "", err
	}
	why := "the split did not work: " + out
	if tail, ok := subtitles.SplitTagged(out); ok {
		trimmed, left := subtitles.TrimRepeated(shown, tail)
		if left {
			if trimmed != tail {
				return trimmed, fastTrimmed, tail, nil
			}
			return trimmed, fastSplit, "", nil
		}
		why = "it only repeated the end of the previous piece: " + tail
	}
	alone, err := srv.translate(ctx, from, to, text, false)
	return alone, fastFallback, why, err
}

func (r *subRun) handleFast(w *subTranslator, b subtitles.Batch, src editorSource) (string, error) {
	cfg := r.app.snapshot()
	st := fastAvail(cfg)
	if !st.ok || st.lang == "" {
		return "", &fastFail{st.missing, st.detail}
	}
	w.mu.Lock()
	prev, shown := w.prevText, w.prevShown
	w.mu.Unlock()
	if st.lang == b.Target {
		return b.Text, nil
	}
	ctx, cancel := context.WithTimeout(w.ctx, subtitles.TranslateTimeout)
	defer cancel()
	srv, err := r.app.fastWait(ctx)
	if err != nil {
		return "", err
	}
	if !srv.has(st.packs) {
		if err := srv.load(ctx, st.packs); err != nil {
			return "", err
		}
	}
	out, way, note, err := fastPiece(ctx, srv, st.lang, b.Target, prev, shown, b.Text, b.Para)
	if err != nil {
		return "", err
	}
	switch way {
	case fastAlone:
		w.fast.alone.Add(1)
	case fastSplit:
		w.fast.split.Add(1)
	case fastTrimmed:
		w.fast.split.Add(1)
		w.fast.trimmed.Add(1)
		if r.opts.logLive {
			log.Printf("subtitles: piece %d repeated the end of the previous piece, %q is shown as %q", b.Seq(), note, out)
		}
	case fastFallback:
		w.fast.fallback.Add(1)
		if r.opts.logLive {
			log.Printf("subtitles: piece %d is translated alone, %s", b.Seq(), note)
		}
	}
	return out, nil
}

func (r *subRun) noteSpeechLangLocked() {
	cfg := r.app.snapshot()
	if l := strings.ToLower(strings.TrimSpace(cfg.SubtitlesLanguage)); l != "" && l != "auto" {
		return
	}
	g := r.phraser.Guess()
	if g == "" {
		return
	}
	if fastRun.setLang(g) {
		log.Printf("subtitles: the speech language is taken as %s from the words heard, %.1f s after the start", g, time.Since(r.started).Seconds())
	}
}

func (w *subTranslator) handleFastBatch(b subtitles.Batch, src editorSource) {
	r := w.run
	seq := b.Seq()
	shown := ""
	defer func() {
		w.mu.Lock()
		w.prevText, w.prevShown = b.Text, shown
		w.mu.Unlock()
	}()
	var out string
	var err error
	for attempt := 1; ; attempt++ {
		out, err = r.handleFast(w, b, src)
		timedOut := errors.Is(err, context.DeadlineExceeded)
		if err == nil || !subtitles.SendAgain(true, timedOut, attempt) || w.ctx.Err() != nil {
			break
		}
		w.retried.Add(1)
		log.Printf("subtitles: piece %d got no answer within %.0f s from %s, it is sent once more", seq, subtitles.TranslateTimeout.Seconds(), src.label())
	}
	done := time.Now()
	if err != nil {
		w.fail(b, err)
		return
	}
	text, verdict := subtitles.CleanFastTranslation(b.Text, out)
	switch verdict {
	case subtitles.TrSame:
		text = b.Text
	case subtitles.TrEmpty:
		w.fail(b, errors.New("the answer is empty"))
		return
	case subtitles.TrRamble:
		w.fail(b, fmt.Errorf("an answer of %d characters for a piece of %d is thrown away", len([]rune(out)), len([]rune(b.Text))))
		return
	}
	shown = text
	r.board.SetPiece(seq, text, subtitles.PieceDone)
	r.pieceFinished(seq, subPieceOut{at: b.Start, source: b.Text, translation: subtitles.WithoutLeadingMark(text), target: b.Target, by: src.label()})
	w.translated.Add(1)
	w.mu.Lock()
	w.delays = append(w.delays, done.Sub(b.Heard))
	waiting := w.q.Len()
	w.mu.Unlock()
	if r.opts.logLive {
		log.Printf("subtitles: piece %d translated by the light engine: complete after %.2f s, %d waiting: %q", seq, done.Sub(b.Heard).Seconds(), waiting, text)
	}
}

func packKey(p fastpack.Pack) string { return "fastpack-" + p.Key() }

func (a *App) packDownload(p fastpack.Pack, base string) {
	key := packKey(p)
	dlMu.Lock()
	if st := dl[key]; st != nil && st.active {
		dlMu.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	dl[key] = &dlState{active: true, cancel: cancel}
	dlMu.Unlock()
	go func() {
		defer cancel()
		err := doPackDownload(ctx, key, p, base)
		fastChecks.Delete(p.Key())
		dlMu.Lock()
		switch {
		case errors.Is(err, context.Canceled):
			log.Printf("download of the language pack %s cancelled", p.Key())
			delete(dl, key)
		case err != nil:
			log.Printf("download of the language pack %s: %v", p.Key(), err)
			dl[key] = &dlState{err: humanError(err)}
		default:
			log.Printf("language pack %s downloaded and checked", p.Key())
			delete(dl, key)
		}
		dlMu.Unlock()
	}()
}

func doPackDownload(ctx context.Context, key string, p fastpack.Pack, base string) error {
	dir := p.Dir("models")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return uiErr("mkdir "+dir+": "+err.Error(), trf("err.model.dir", dir))
	}
	need := int64(0)
	for _, f := range p.Files {
		if st, ok := trmodelKnown(dir, f); ok && st {
			continue
		}
		need += f.Size
	}
	if free := freeDiskMB("models"); free >= 0 && int64(free) < need/(1024*1024)+256 {
		return uiErr(fmt.Sprintf("not enough disk space: %d MB free, %d MB needed", free, need/(1024*1024)), trf("err.disk.space", free, int(need/(1024*1024))))
	}
	total := p.Size()
	var done int64
	for _, f := range p.Files {
		if ok, _ := trmodelKnown(dir, f); ok {
			done += f.Size
			continue
		}
		url := f.URL
		if base != "" {
			url = strings.TrimRight(base, "/") + "/" + f.Name
		}
		written, err := downloadFile(ctx, url, filepath.Join(dir, f.Name), func(n int64) {
			dlMu.Lock()
			if st := dl[key]; st != nil && total > 0 {
				st.pct = min(99, int((done+n)*100/total))
			}
			dlMu.Unlock()
		})
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return err
			}
			return uiErr(f.Name+": "+logText(err), f.Name+" — "+humanError(err))
		}
		st, err := trmodelCheck(dir, f)
		if err != nil {
			return err
		}
		if !st {
			log.Printf("%s: the downloaded file does not match its checksum", f.Name)
			_ = os.Remove(filepath.Join(dir, f.Name))
			return uiErr("the downloaded file does not match its checksum", tr("err.hash"))
		}
		done += written
	}
	return nil
}

func trmodelKnown(dir string, f fastpack.File) (bool, bool) {
	st, ok := trmodel.Known(dir, f.Spec())
	return ok && st == trmodel.Ready, ok
}

func trmodelCheck(dir string, f fastpack.File) (bool, error) {
	st, err := trmodel.Check(dir, f.Spec())
	return st == trmodel.Ready, err
}

func (a *App) packDelete(p fastpack.Pack) error {
	if s, ok := a.fastEng.live().(*fastServer); ok && s != nil && s.has([]fastpack.Pack{p}) {
		a.fastShutdown()
	}
	cancelDownload(packKey(p))
	err := fastpack.Remove("models", p)
	fastChecks.Delete(p.Key())
	dlMu.Lock()
	delete(dl, packKey(p))
	dlMu.Unlock()
	log.Printf("language pack %s deleted", p.Key())
	return err
}

type packRoute struct {
	Lang   string   `json:"lang"`
	Packs  []string `json:"packs"`
	SizeMB int      `json:"size"`
	State  string   `json:"state"`
	Pct    int      `json:"pct"`
	Err    string   `json:"err"`
}

type fastUI struct {
	Lang   string      `json:"lang"`
	Auto   bool        `json:"auto"`
	Exe    bool        `json:"exe"`
	Down   string      `json:"down"`
	Routes []packRoute `json:"routes"`
}

func routeView(lang, to string) packRoute {
	v := packRoute{Lang: lang}
	packs, ok := fastpack.Route(lang, to)
	if !ok {
		v.State = "none"
		return v
	}
	var size, doneSize int64
	states := map[string]bool{}
	for _, p := range packs {
		v.Packs = append(v.Packs, p.Key())
		size += p.Size()
		dlMu.Lock()
		st := dl[packKey(p)]
		switch {
		case st != nil && st.active:
			states["downloading"] = true
			doneSize += p.Size() * int64(st.pct) / 100
		case st != nil && st.err != "":
			states["error"] = true
			v.Err = st.err
		default:
			s := packState(p)
			states[string(s)] = true
			if s == fastpack.Ready {
				doneSize += p.Size()
			}
		}
		dlMu.Unlock()
	}
	v.SizeMB = int((size + 512*1024) / (1024 * 1024))
	for _, s := range []string{"downloading", "error", "bad", "checking", "absent"} {
		if states[s] {
			v.State = s
			break
		}
	}
	if v.State == "" {
		v.State = "ready"
	}
	if v.State == "downloading" && size > 0 {
		v.Pct = int(doneSize * 100 / size)
	}
	return v
}

func (a *App) fastView(cfg *Config) fastUI {
	to := cfg.SubtitlesTranslate
	v := fastUI{Auto: strings.ToLower(strings.TrimSpace(cfg.SubtitlesLanguage)) == "auto" || strings.TrimSpace(cfg.SubtitlesLanguage) == ""}
	v.Lang = fastSpeechLang(cfg)
	_, err := os.Stat(fastExePath(cfg))
	v.Exe = err == nil
	_, _, v.Down = fastRun.snapshot()
	if to == "" {
		return v
	}
	for _, lang := range fastpack.Languages {
		if lang != to {
			v.Routes = append(v.Routes, routeView(lang, to))
		}
	}
	return v
}

func (a *App) fastPackAction(action, lang string) string {
	cfg := a.snapshot()
	packs, ok := fastpack.Route(lang, cfg.SubtitlesTranslate)
	if !ok {
		return ""
	}
	switch action {
	case "download":
		for _, p := range packs {
			if packState(p) != fastpack.Ready {
				log.Printf("language pack %s: download started", p.Key())
				a.packDownload(p, "")
			}
		}
	case "cancel":
		for _, p := range packs {
			cancelDownload(packKey(p))
		}
	case "delete":
		for _, p := range packs {
			if err := a.packDelete(p); err != nil {
				return humanError(err)
			}
		}
		return tr("model.del.ok")
	}
	return ""
}

func runPackDownload(args []string) {
	if len(args) < 2 {
		log.Printf("packdownload: needs <from> <to> [base url]")
		return
	}
	base := ""
	if len(args) > 2 {
		base = args[2]
	}
	p, ok := fastpack.Find(args[0], args[1])
	if !ok {
		log.Printf("packdownload: no pack %s-%s in the catalog", args[0], args[1])
		return
	}
	app := &App{cfg: defaultConfig(), enabled: true}
	started := time.Now()
	app.packDownload(p, base)
	for {
		time.Sleep(200 * time.Millisecond)
		dlMu.Lock()
		st := dl[packKey(p)]
		active := st != nil && st.active
		dlMu.Unlock()
		if !active {
			break
		}
	}
	v := routeView(p.From, p.To)
	log.Printf("packdownload: %s is %s %s after %.1f s", p.Key(), v.State, v.Err, time.Since(started).Seconds())
}

func runFastCheck(args []string) {
	if len(args) < 4 {
		log.Printf("fastcheck: needs <from> <to> <pieces file> <output file>")
		return
	}
	from, to := args[0], args[1]
	packs, ok := fastpack.Route(from, to)
	if !ok {
		log.Printf("fastcheck: no route %s → %s", from, to)
		return
	}
	raw, err := os.ReadFile(args[2])
	if err != nil {
		log.Printf("fastcheck: %v", err)
		return
	}
	cfg, err := loadConfig("config.json")
	if err != nil {
		cfg = defaultConfig()
	}
	app := &App{cfg: cfg, enabled: true}
	defer app.fastShutdown()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	f := app.warmEngine(&app.fastEng, fastEngine)
	<-f.done
	srv, err := asFast(f.srv, f.err)
	if err != nil {
		log.Printf("fastcheck: the engine did not start: %v", err)
		return
	}
	if err := srv.load(ctx, packs); err != nil {
		log.Printf("fastcheck: %v", err)
		return
	}
	var out strings.Builder
	prev, shown := "", ""
	para := true
	counts := map[string]int{}
	for _, line := range strings.Split(strings.ReplaceAll(string(raw), "\r\n", "\n"), "\n") {
		text := strings.TrimSpace(line)
		if text == "" {
			para = true
			continue
		}
		started := time.Now()
		got, way, note, err := fastPiece(ctx, srv, from, to, prev, shown, text, para)
		ms := time.Since(started).Seconds() * 1000
		if err != nil {
			log.Printf("fastcheck: %q: %v", text, err)
			return
		}
		cleaned, verdict := subtitles.CleanFastTranslation(text, got)
		if verdict == subtitles.TrSame {
			cleaned = text
		}
		counts[way]++
		fmt.Fprintf(&out, "%s\t%s\t%.0f ms\t%s\t%s\n", way, text, ms, cleaned, note)
		prev, shown, para = text, cleaned, false
	}
	if err := os.WriteFile(args[3], []byte(out.String()), 0o644); err != nil {
		log.Printf("fastcheck: %v", err)
		return
	}
	log.Printf("fastcheck: %s → %s through %d pack(s): %d alone, %d split, %d with a repeated word removed, %d fell back to alone", from, to, len(packs), counts[fastAlone], counts[fastSplit], counts[fastTrimmed], counts[fastFallback])
}
