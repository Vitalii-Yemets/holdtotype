package main

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync/atomic"

	"holdtotype/internal/subtitles"
)

type draftReq struct {
	seq  int
	text string
	to   string
}

type draftCounts struct {
	asked  atomic.Int64
	shown  atomic.Int64
	yield  atomic.Int64
	failed atomic.Int64
	late   atomic.Int64
}

func (r *subRun) offerDraft(seq int, para bool, raw string) {
	cfg := r.app.snapshot()
	if !cfg.SubtitlesSentenceDraft {
		return
	}
	t := r.target()
	if t == nil {
		return
	}
	text := r.cleanText(cfg, t, raw)
	to := translationTarget(cfg, t)
	if text == "" || to == "" || translationSource(cfg).kind == editorFast {
		return
	}
	w := r.translator()
	if w == nil {
		return
	}
	w.mu.Lock()
	if w.closed {
		w.mu.Unlock()
		return
	}
	r.board.AddPiece(seq, para, text)
	w.draftLines.Offer(seq)
	w.draft = &draftReq{seq: seq, text: text, to: to}
	w.mu.Unlock()
	w.drafts.asked.Add(1)
	if r.opts.logLive {
		log.Printf("subtitles: a draft is asked for line %d, %d words", seq, len(strings.Fields(text)))
	}
	w.wake()
}

func (r *subRun) dropDraft(seq int, para bool, text string) {
	r.trMu.Lock()
	w := r.tr
	r.trMu.Unlock()
	if w == nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.draftLines.Drop(seq) {
		return
	}
	if w.draft != nil && w.draft.seq == seq {
		w.draft = nil
	}
	r.board.AddPiece(seq, para, text)
	r.board.SetPiece(seq, "", subtitles.PieceFailed)
	log.Printf("subtitles: line %d had a draft and is shown as heard: translation is off", seq)
}

func (w *subTranslator) showDraft(seq int, text string) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.draftLines.Show(seq) {
		return false
	}
	w.run.board.SetPiece(seq, text, subtitles.PieceRunning)
	return true
}

func (w *subTranslator) draftStep() bool {
	r := w.run
	cfg := r.app.snapshot()
	src := translationSource(cfg)
	w.mu.Lock()
	d := w.draft
	if d == nil || w.q.Len() > 0 {
		w.mu.Unlock()
		return false
	}
	if !src.ok() || src.kind == editorFast || needsStart(src.kind) && !r.app.translatorWarm(src.kind) {
		w.mu.Unlock()
		return false
	}
	w.draft = nil
	if !w.draftLines.Open(d.seq) {
		w.mu.Unlock()
		return true
	}
	yctx, giveWay := context.WithCancelCause(w.ctx)
	w.giveWay = giveWay
	w.mu.Unlock()
	hist := w.history(src.kind, d.to)
	from := ""
	if t := r.target(); t != nil {
		from = speechLang(t.lang)
	}
	req := trRequestFor(src, from, d.to, hist.Turns(), d.text, w.topicContext(cfg, src, d.to))
	ctx, cancel := context.WithTimeout(yctx, subtitles.TranslateTimeout)
	out, err := r.app.translateStream(ctx, cfg, src, req, subtitles.TranslateTokens(len(strings.Fields(d.text))), func(sofar string) {
		w.showDraft(d.seq, subtitles.StreamingText(sofar))
	})
	cancel()
	yielded := errors.Is(context.Cause(yctx), errGiveWay)
	giveWay(nil)
	w.mu.Lock()
	w.giveWay = nil
	w.mu.Unlock()
	if w.ctx.Err() != nil {
		return true
	}
	if err != nil {
		if yielded {
			w.drafts.yield.Add(1)
		} else {
			w.drafts.failed.Add(1)
		}
		return true
	}
	text, verdict := subtitles.CleanTranslation(d.text, out, d.to)
	if verdict != subtitles.TrOK {
		w.drafts.failed.Add(1)
		return true
	}
	if w.showDraft(d.seq, text) {
		w.drafts.shown.Add(1)
	} else {
		w.drafts.late.Add(1)
	}
	return true
}

func (w *subTranslator) logDrafts() {
	c := &w.drafts
	if c.asked.Load() == 0 {
		return
	}
	log.Printf("subtitles: drafts: %d asked, %d shown, %d gave way to a piece, %d failed, %d came after the whole sentence", c.asked.Load(), c.shown.Load(), c.yield.Load(), c.failed.Load(), c.late.Load())
}
