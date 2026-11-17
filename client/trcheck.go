package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"holdtotype/internal/subtitles"
	"holdtotype/internal/trmodel"
)

type trCheckRow struct {
	Kind     string  `json:"kind"`
	Seqs     []int   `json:"seqs,omitempty"`
	Source   string  `json:"source,omitempty"`
	Result   string  `json:"result,omitempty"`
	Verdict  string  `json:"verdict,omitempty"`
	Ready    float64 `json:"ready"`
	Start    float64 `json:"start"`
	End      float64 `json:"end"`
	Seconds  float64 `json:"seconds"`
	LastWord float64 `json:"last_word,omitempty"`
	Delay    float64 `json:"delay,omitempty"`
	Context  string  `json:"context,omitempty"`
	Terms    int     `json:"terms,omitempty"`
	Err      string  `json:"err,omitempty"`
	Yielded  bool    `json:"yielded,omitempty"`
	Lines    int     `json:"lines,omitempty"`
	Hold     float64 `json:"hold,omitempty"`
}

func percentile(xs []float64, q float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	s := append([]float64(nil), xs...)
	sort.Float64s(s)
	i := int(q*float64(len(s)-1) + 0.5)
	return s[i]
}

func runTrCheck(args []string) {
	if len(args) < 3 {
		log.Printf("trcheck: needs <text file> <run start hh:mm:ss or all> <output file> [model id or own:file] [target language] [off, topic or terms]")
		return
	}
	raw, err := os.ReadFile(args[0])
	if err != nil {
		log.Printf("trcheck: %v", err)
		return
	}
	lines, lang := subtitles.ParseCall(string(raw), args[1])
	if len(lines) == 0 {
		log.Printf("trcheck: no lines for the run %s", args[1])
		return
	}
	cfg, err := loadConfig("config.json")
	if err != nil {
		cfg = defaultConfig()
	}
	if len(args) > 3 && args[3] != "" {
		cfg.SubtitlesModel = args[3]
	}
	for _, m := range trmodel.Catalog {
		if st, err := trmodel.Check("models", m.Spec()); err == nil && st != trmodel.Absent {
			log.Printf("trcheck: %s is %s", m.Name, st)
		}
	}
	forgetTrModelCheck()
	target := cfg.SubtitlesTranslate
	if len(args) > 4 {
		target = args[4]
	}
	if subtitles.LanguageName(target) == "" {
		target = "ru"
	}
	m, ok := currentTrModel(cfg)
	if !ok || (cfg.SubtitlesModel != "" && m.Key() != cfg.SubtitlesModel) {
		log.Printf("trcheck: the model %q is not ready", cfg.SubtitlesModel)
		return
	}
	app := &App{cfg: cfg, enabled: true}
	defer app.qualityShutdown()
	started := time.Now()
	f := app.warmEngine(&app.gpuEng, gpuEngine)
	<-f.done
	srv, err := asLlama(f.srv, f.err)
	if err != nil {
		fail := asQualityFail(err)
		log.Printf("trcheck: the engine did not start (%s): %s", fail.code, fail.detail)
		return
	}
	loaded := time.Since(started)
	if err := srv.prime(target); err != nil {
		log.Printf("trcheck: the priming request failed: %v", err)
		return
	}
	log.Printf("trcheck: %s loaded in %.1f s, warm after %.1f s", m.Name, loaded.Seconds(), time.Since(started).Seconds())
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.Local)
	pieces := subtitles.PhraseCall(lines, lang, base)
	from := speechLang(lang)
	src := editorSource{kind: editorQuality, model: m}
	useTopic, useTerms := cfg.SubtitlesTopic && m.Topic, m.Terms
	mode := ""
	if len(args) > 5 {
		mode = args[5]
	}
	switch mode {
	case "off":
		useTopic, useTerms = false, false
	case "topic":
		useTopic, useTerms = true, false
	case "terms":
		useTopic, useTerms = true, true
	}
	var topic *subtitles.Topic
	if useTopic {
		topic = subtitles.NewTopic(target)
		topic.UseTerms(useTerms, base)
	}
	calls := subtitles.CheckCalls{
		Translate: func(b subtitles.Batch, prior []subtitles.Turn, context string) (string, subtitles.Verdict, time.Duration, error) {
			req := trRequestFor(src, from, target, prior, b.Text, context)
			ctx, cancel := ctxTimeout(subtitles.TranslateTimeout * 4)
			defer cancel()
			t0 := time.Now()
			var out string
			var err error
			if req.prompt != "" {
				out, err = srv.completeStream(ctx, req.prompt, subtitles.TranslateTokens(b.Words), nil)
			} else {
				out, err = srv.chatStream(ctx, req.msgs, subtitles.TranslateTokens(b.Words), "", nil)
			}
			d := time.Since(t0)
			if err != nil {
				return b.Text, subtitles.TrEmpty, d, err
			}
			text, verdict := subtitles.CleanTranslation(b.Text, out, target)
			switch verdict {
			case subtitles.TrOK:
				text = subtitles.KeepSentenceEnd(b.Text, subtitles.DropAddedFullStop(b.Text, text))
			case subtitles.TrSame:
				text = b.Text
			default:
				text = b.Text
			}
			return text, verdict, d, nil
		},
	}
	if topic != nil {
		calls.Service = func(step subtitles.Step, msgs []subtitles.Message, limit time.Duration) (string, subtitles.ServiceEnd, time.Duration, error) {
			grammar, tokens := subtitles.TopicGrammar, subtitles.TopicTokens
			if step == subtitles.StepTerms {
				grammar, tokens = subtitles.TermsGrammar(target), subtitles.TermsTokens
			}
			ctx, cancel := ctxTimeout(subtitles.TopicTimeout)
			defer cancel()
			yctx := ctx
			if limit > 0 && limit < subtitles.TopicTimeout {
				var stop context.CancelFunc
				yctx, stop = context.WithTimeout(ctx, limit)
				defer stop()
			}
			var sofar string
			t0 := time.Now()
			out, finish, err := srv.chatStreamEnd(yctx, msgs, tokens, grammar, func(s string) { sofar = s })
			d := time.Since(t0)
			if out == "" {
				out = sofar
			}
			switch {
			case err == nil && finish == "length":
				return out, subtitles.EndLimit, d, nil
			case err == nil:
				return out, subtitles.EndWhole, d, nil
			case ctx.Err() == context.DeadlineExceeded:
				return out, subtitles.EndTimeout, d, err
			case yctx.Err() == context.DeadlineExceeded:
				return out, subtitles.EndYield, d, err
			}
			return out, subtitles.EndWhole, d, err
		}
	}
	events := subtitles.RunCheck(pieces, target, topic, calls)
	last := map[int]time.Time{}
	for _, p := range pieces {
		last[p.Seq] = p.LastWord
	}
	secs := func(t time.Time) float64 { return t.Sub(base).Seconds() }
	var rows []trCheckRow
	var delays []float64
	gpu := 0.0
	count := map[string]int{}
	yields, hold := 0, 0.0
	for _, e := range events {
		row := trCheckRow{Kind: e.Kind, Seqs: e.Seqs, Source: e.Text, Result: e.Result, Verdict: string(e.Verdict), Ready: secs(e.Ready), Start: secs(e.Start), End: secs(e.End), Seconds: e.Duration.Seconds(), Context: e.Context, Terms: e.Terms, Err: e.Err, Yielded: e.Yielded, Lines: e.Lines, Hold: e.Hold.Seconds()}
		count[e.Kind]++
		if e.Yielded {
			yields++
		}
		if e.Hold.Seconds() > hold {
			hold = e.Hold.Seconds()
		}
		gpu += e.Duration.Seconds()
		if e.Kind == "piece" {
			for _, s := range e.Seqs {
				d := e.End.Sub(last[s]).Seconds()
				delays = append(delays, d)
				if row.LastWord == 0 || secs(last[s]) > row.LastWord {
					row.LastWord = secs(last[s])
				}
			}
			row.Delay = e.End.Sub(last[e.Seqs[len(e.Seqs)-1]]).Seconds()
		}
		rows = append(rows, row)
	}
	speech := pieces[len(pieces)-1].LastWord.Sub(pieces[0].LastWord).Minutes()
	if first := lines[0]; len(first.Words) > 0 {
		speech = pieces[len(pieces)-1].LastWord.Sub(base.Add(first.At - time.Duration(float64(len(first.Words)-1)/subtitles.CheckRate*float64(time.Second)))).Minutes()
	}
	summary := map[string]any{
		"model": m.Name, "model_id": m.Key(), "target": target, "run": args[1], "speech_language": lang,
		"loaded_s": loaded.Seconds(), "pieces": len(pieces), "requests": count["piece"], "stale": count["stale"],
		"topic_requests": count["topic"], "terms_requests": count["terms"],
		"delay_median_s": percentile(delays, 0.5), "delay_p90_s": percentile(delays, 0.9), "delay_worst_s": percentile(delays, 1),
		"speech_minutes": speech, "gpu_s_per_min": gpu / speech, "topic_on": topic != nil, "terms_on": topic != nil && useTerms, "mode": mode,
		"service_yielded": yields, "service_hold_max_s": hold,
	}
	if topic != nil {
		summary["topic"] = topic.Text()
		summary["terms"] = topic.Terms()
		summary["topic_counts"] = topic.Counts
	}
	out, _ := json.MarshalIndent(map[string]any{"summary": summary, "events": rows}, "", " ")
	if err := os.WriteFile(args[2], out, 0o644); err != nil {
		log.Printf("trcheck: %v", err)
		return
	}
	log.Printf("trcheck: %d piece(s) in %d request(s), delay median %.1f s, worst %.1f s, %s", len(pieces), count["piece"], percentile(delays, 0.5), percentile(delays, 1), strings.TrimSpace(fmt.Sprintf("%d topic and %d terms request(s)", count["topic"], count["terms"])))
}

func ctxTimeout(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d)
}
