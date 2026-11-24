package main

import (
	"log"
	"path/filepath"
	"strings"
	"time"

	"holdtotype/internal/appid"
	"holdtotype/internal/appname"
	"holdtotype/internal/stats"
	"holdtotype/internal/subtitles"
)

var statsStore = stats.Open(appid.StatsFile, log.Printf)

const (
	statsFilePrefix   = "file:"
	subStatsEvery     = 2 * time.Minute
	subStatsMaxStep   = 5 * time.Second
	statsLangFallback = "auto"
)

type dictMeta struct {
	stopAt   time.Time
	speechMs int64
	model    string
	lang     string
}

func statsModelKey(srv recognizer, active *modelInfo) string {
	p := ""
	if srv != nil && !srv.external() {
		p = strings.TrimSpace(srv.model())
	}
	if i := strings.IndexByte(p, '#'); i >= 0 {
		p = p[:i]
	}
	if p != "" {
		base := filepath.Base(filepath.Clean(p))
		for i := range modelCatalog {
			m := &modelCatalog[i]
			if (m.File != "" && filepath.Base(m.File) == base) || (m.Dir != "" && filepath.Base(filepath.Clean(m.Dir)) == base) {
				return m.ID
			}
		}
		return statsFilePrefix + base
	}
	if active != nil {
		return active.ID
	}
	if srv != nil {
		return srv.engine()
	}
	return ""
}

func statsModelName(key string) string {
	if name, ok := strings.CutPrefix(key, statsFilePrefix); ok {
		return name
	}
	if m := findModel(key); m != nil {
		return m.NameKey
	}
	return key
}

func statsLang(lang, text string) string {
	lang = strings.ToLower(strings.TrimSpace(lang))
	if lang != "" && lang != "auto" {
		return lang
	}
	if g := subtitles.GuessLanguage(strings.Fields(text)); validTranslateLang(g) {
		return g
	}
	return statsLangFallback
}

var statsAppNames = appname.NewCache(appname.ReadVersion)

func statsProgram(cfg *Config, app, path string) string {
	prog := stats.Program(app, cfg.HistorySkip)
	if prog == stats.OtherApp || strings.TrimSpace(path) == "" {
		return prog
	}
	if name := statsAppNames.Name(path); name != "" {
		return name
	}
	return prog
}

func runAppNameCheck(paths []string) {
	if len(paths) == 0 {
		log.Printf("appnamecheck: needs one or more paths to .exe files")
		return
	}
	for _, p := range paths {
		desc, product := appname.ReadVersion(p)
		log.Printf("appnamecheck: %s -> %q (description %q, product %q)", p, statsAppNames.Name(p), desc, product)
	}
	log.Printf("appnamecheck: done")
}

func strOr(p *string, def string) string {
	if p == nil {
		return def
	}
	return *p
}

func runStatsCheck(args []string) {
	out, err := stats.CheckReport(appid.StatsFile, args, time.Now(), statsModelName)
	if err != nil {
		log.Printf("statscheck: %v", err)
		return
	}
	log.Printf("statscheck: %s", out)
}

func (a *App) countDictation(cfg *Config, text, app, path string, meta dictMeta, insertedAt time.Time) {
	d := stats.Dictation{
		Text:     text,
		SpeechMs: meta.speechMs,
		Model:    meta.model,
		Lang:     statsLang(meta.lang, text),
		App:      statsProgram(cfg, app, path),
	}
	if !meta.stopAt.IsZero() && insertedAt.After(meta.stopAt) {
		d.WaitMs = insertedAt.Sub(meta.stopAt).Milliseconds()
	} else {
		d.Model = ""
	}
	if err := statsStore.AddDictation(insertedAt, d); err != nil {
		log.Printf("statistics: %v", err)
	}
}

type subTally struct {
	last    time.Time
	flushed time.Time
	listen  time.Duration
	pairs   map[string]time.Duration
}

func (t *subTally) tick(now time.Time, listening bool, pair string) {
	if t.last.IsZero() {
		t.last, t.flushed = now, now
		return
	}
	step := now.Sub(t.last)
	t.last = now
	if step <= 0 || !listening {
		return
	}
	step = min(step, subStatsMaxStep)
	t.listen += step
	if pair != "" {
		if t.pairs == nil {
			t.pairs = map[string]time.Duration{}
		}
		t.pairs[pair] += step
	}
}

func (t *subTally) flush(now time.Time, force bool) {
	if !force && now.Sub(t.flushed) < subStatsEvery {
		return
	}
	t.flushed = now
	if t.listen <= 0 && len(t.pairs) == 0 {
		return
	}
	pairs := map[string]int64{}
	for k, v := range t.pairs {
		pairs[k] = v.Milliseconds()
	}
	listen := t.listen.Milliseconds()
	t.listen, t.pairs = 0, nil
	if err := statsStore.AddSubtitles(now, listen, pairs); err != nil {
		log.Printf("statistics: subtitles: %v", err)
	}
}

func (r *subRun) tallyPair() string {
	if !r.trShown.Load() {
		return ""
	}
	to, _ := r.trTo.Load().(string)
	t := r.target()
	if to == "" || t == nil {
		return ""
	}
	from := t.lang
	if from == "" || from == "auto" {
		if g, _, _ := fastRun.snapshot(); g != "" {
			from = g
		} else {
			from = "auto"
		}
	}
	if from == to {
		return ""
	}
	return stats.PairKey(from, to)
}

func (r *subRun) tally(now time.Time, final bool) {
	r.stat.tick(now, r.windowPhase() == subPhaseListening, r.tallyPair())
	r.stat.flush(now, final)
}
