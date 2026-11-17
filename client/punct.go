package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"syscall"
	"time"

	"holdtotype/internal/punct"
	"holdtotype/internal/sentence"
	"holdtotype/internal/subtitles"
)

const (
	punctExeName    = "punct-engine.exe"
	punctStartWait  = 20 * time.Second
	punctAnswerWait = 1500 * time.Millisecond
)

var punctModelDir = filepath.Join("models", "punct", "horizon")

type punctEngine struct {
	cmd    *exec.Cmd
	in     io.WriteCloser
	lines  chan []byte
	LoadMS float64
}

type punctAnswer struct {
	Ready  bool      `json:"ready"`
	LoadMS float64   `json:"load_ms"`
	End    []bool    `json:"end"`
	Prob   []float64 `json:"prob"`
	MS     float64   `json:"ms"`
	Tokens int       `json:"tokens"`
	Error  string    `json:"error"`
}

func startPunctEngine(threads int) (*punctEngine, error) {
	exe := filepath.Join(filepath.Dir(absPath(os.Args[0])), punctExeName)
	if _, err := os.Stat(exe); err != nil {
		return nil, fmt.Errorf("%s is missing next to the program", punctExeName)
	}
	cmd := exec.Command(exe, "-model", absPath(punctModelDir), "-threads", strconv.Itoa(threads))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	in, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("starting %s: %w", punctExeName, err)
	}
	e := &punctEngine{cmd: cmd, in: in, lines: make(chan []byte, 1)}
	go func() {
		r := bufio.NewReaderSize(outPipe, 1<<16)
		for {
			line, err := r.ReadBytes('\n')
			if len(line) > 0 {
				e.lines <- line
			}
			if err != nil {
				close(e.lines)
				return
			}
		}
	}()
	a, err := e.read(punctStartWait)
	if err != nil || !a.Ready {
		e.stop()
		if err == nil {
			err = errors.New("the engine did not say it is ready")
		}
		return nil, fmt.Errorf("%s did not start: %w", punctExeName, err)
	}
	e.LoadMS = a.LoadMS
	return e, nil
}

func (e *punctEngine) read(wait time.Duration) (punctAnswer, error) {
	var line []byte
	select {
	case l, ok := <-e.lines:
		if !ok {
			return punctAnswer{}, errors.New("the engine process ended")
		}
		line = l
	case <-time.After(wait):
		e.kill()
		return punctAnswer{}, fmt.Errorf("no answer within %.1f s, the engine is stopped", wait.Seconds())
	}
	var a punctAnswer
	if err := json.Unmarshal(line, &a); err != nil {
		return punctAnswer{}, fmt.Errorf("the engine answered %.200q", line)
	}
	if a.Error != "" {
		return a, errors.New(a.Error)
	}
	return a, nil
}

func (e *punctEngine) judge(words []string) (punctAnswer, error) {
	body, _ := json.Marshal(map[string][]string{"words": words})
	if _, err := e.in.Write(append(body, '\n')); err != nil {
		return punctAnswer{}, err
	}
	a, err := e.read(punctAnswerWait)
	if err != nil {
		return a, err
	}
	if len(a.End) != len(words) || len(a.Prob) != len(words) {
		return a, fmt.Errorf("the engine answered for %d word(s) instead of %d", len(a.End), len(words))
	}
	return a, nil
}

func (e *punctEngine) kill() {
	if e.cmd != nil && e.cmd.Process != nil {
		_ = e.cmd.Process.Kill()
	}
}

func (e *punctEngine) stop() {
	if e.in != nil {
		_ = e.in.Close()
	}
	if e.cmd != nil && e.cmd.Process != nil {
		done := make(chan struct{})
		go func() {
			_ = e.cmd.Wait()
			close(done)
		}()
		select {
		case <-done:
		case <-time.After(3 * time.Second):
			_ = e.cmd.Process.Kill()
		}
	}
}

type punctRun struct {
	Words int     `json:"words"`
	MS    float64 `json:"ms"`
	Round float64 `json:"round_ms"`
}

type punctSentence struct {
	First    int     `json:"first"`
	End      int     `json:"end"`
	Words    int     `json:"words"`
	How      string  `json:"how"`
	At       float64 `json:"at"`
	LastWord float64 `json:"last_word"`
}

func spread(xs []float64) map[string]float64 {
	if len(xs) == 0 {
		return nil
	}
	s := append([]float64(nil), xs...)
	sort.Float64s(s)
	q := func(f float64) float64 { return s[int(f*float64(len(s)-1)+0.5)] }
	return map[string]float64{"count": float64(len(s)), "median": q(0.5), "p90": q(0.9), "max": s[len(s)-1]}
}

func runPunctCheck(args []string) {
	if len(args) < 3 {
		log.Printf("punctcheck: needs <text file> <run start hh:mm:ss or all> <output file> [words ahead] [silence, s] [limit, words]")
		return
	}
	raw, err := os.ReadFile(args[0])
	if err != nil {
		log.Printf("punctcheck: %v", err)
		return
	}
	lines, _ := subtitles.ParseCall(string(raw), args[1])
	if len(lines) == 0 {
		log.Printf("punctcheck: no lines for the run %s", args[1])
		return
	}
	p := sentence.Defaults()
	if len(args) > 3 {
		if n, err := strconv.Atoi(args[3]); err == nil && n >= 0 {
			p.Ahead = n
		}
	}
	if len(args) > 4 {
		if f, err := strconv.ParseFloat(args[4], 64); err == nil && f >= 0 {
			p.Silence = time.Duration(f * float64(time.Second))
		}
	}
	if len(args) > 5 {
		if n, err := strconv.Atoi(args[5]); err == nil && n >= 0 {
			p.Limit = n
		}
	}
	threads := 1
	if n, err := strconv.Atoi(os.Getenv("HT_PUNCT_THREADS")); err == nil && n > 0 {
		threads = n
	}
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	var words []string
	var times []time.Time
	for _, l := range lines {
		for k, w := range l.Words {
			n := punct.Norm(w)
			if n == "" {
				continue
			}
			back := time.Duration(float64(len(l.Words)-1-k) / subtitles.CheckRate * float64(time.Second))
			words = append(words, n)
			times = append(times, base.Add(l.At-back))
		}
	}
	started := time.Now()
	eng, err := startPunctEngine(threads)
	if err != nil {
		log.Printf("punctcheck: %v", err)
		return
	}
	defer eng.stop()
	log.Printf("punctcheck: the engine is ready after %.1f s (model loaded in %.1f s), %d word(s), %d thread(s)", time.Since(started).Seconds(), eng.LoadMS/1000, len(words), threads)
	secs := func(t time.Time) float64 { return t.Sub(base).Seconds() }
	result := map[string]any{"run": args[1], "words": len(words), "threads": threads, "load_ms": eng.LoadMS,
		"params": map[string]any{"ahead": p.Ahead, "silence_s": p.Silence.Seconds(), "limit": p.Limit, "context": p.Context}}
	if os.Getenv("HT_PUNCT_BENCH") != "" {
		bench := map[string]any{}
		for _, size := range []int{20, 40, 70} {
			var model, round []float64
			for r := 0; r < 40 && r*7+size <= len(words); r++ {
				t0 := time.Now()
				a, err := eng.judge(words[r*7 : r*7+size])
				if err != nil {
					log.Printf("punctcheck: %v", err)
					return
				}
				if r >= 5 {
					model = append(model, a.MS)
					round = append(round, float64(time.Since(t0).Microseconds())/1000)
				}
			}
			bench[strconv.Itoa(size)] = map[string]any{"model_ms": spread(model), "round_ms": spread(round)}
		}
		result["bench"] = bench
	}
	var runs []punctRun
	judge := func(ws []string) ([]bool, []float64, error) {
		t0 := time.Now()
		a, err := eng.judge(ws)
		if err != nil {
			return nil, nil, err
		}
		runs = append(runs, punctRun{Words: len(ws), MS: a.MS, Round: float64(time.Since(t0).Microseconds()) / 1000})
		return a.End, a.Prob, nil
	}
	st := sentence.New(p, judge)
	var closes []sentence.Close
	t0 := time.Now()
	for i, w := range words {
		c, err := st.Add(w, times[i])
		if err != nil {
			log.Printf("punctcheck: %v", err)
			return
		}
		closes = append(closes, c...)
	}
	closes = append(closes, st.Flush()...)
	var out []punctSentence
	var lens, ms, rounds []float64
	how := map[string]int{}
	for _, c := range closes {
		out = append(out, punctSentence{First: c.First, End: c.End, Words: c.Words(), How: c.How, At: secs(c.At), LastWord: secs(times[c.End])})
		lens = append(lens, float64(c.Words()))
		how[c.How]++
	}
	for _, r := range runs {
		ms = append(ms, r.MS)
		rounds = append(rounds, r.Round)
	}
	result["sentences"] = out
	result["count"] = len(out)
	result["by"] = how
	result["lengths"] = spread(lens)
	result["model_ms"] = spread(ms)
	result["round_ms"] = spread(rounds)
	result["runs"] = runs
	result["seconds"] = time.Since(t0).Seconds()
	body, _ := json.MarshalIndent(result, "", " ")
	if err := os.WriteFile(args[2], body, 0o644); err != nil {
		log.Printf("punctcheck: %v", err)
		return
	}
	log.Printf("punctcheck: %d sentence(s) from %d word(s) in %d run(s) of the model, %.1f ms each (median), by %v", len(out), len(words), len(runs), spread(ms)["median"], how)
}
