package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"golang.org/x/sys/windows"

	"holdtotype/internal/ort"
	"holdtotype/internal/punct"
)

type request struct {
	Words []string `json:"words"`
}

type answer struct {
	Ready  bool      `json:"ready,omitempty"`
	LoadMS float64   `json:"load_ms,omitempty"`
	End    []bool    `json:"end,omitempty"`
	Prob   []float64 `json:"prob,omitempty"`
	MS     float64   `json:"ms,omitempty"`
	Tokens int       `json:"tokens,omitempty"`
	Error  string    `json:"error,omitempty"`
}

func main() {
	dir := flag.String("model", filepath.Join("models", "punct", "horizon"), "")
	threads := flag.Int("threads", 1, "")
	dll := flag.String("ort", "", "")
	arena := flag.Bool("arena", false, "")
	pattern := flag.Bool("pattern", false, "")
	flag.Parse()
	out := json.NewEncoder(os.Stdout)
	if *dll == "" {
		exe, _ := os.Executable()
		*dll = filepath.Join(filepath.Dir(exe), "onnxruntime.dll")
	}
	start := time.Now()
	if err := punct.CheckRuntime(loadSystem); err != nil {
		_ = out.Encode(answer{Error: err.Error()})
		os.Exit(1)
	}
	tok, err := punct.LoadTokenizer(filepath.Join(*dir, "tokenizer.json"))
	if err != nil {
		_ = out.Encode(answer{Error: err.Error()})
		os.Exit(1)
	}
	s, err := ort.Open(*dll, filepath.Join(*dir, "model_quantized.onnx"), ort.Options{Threads: *threads, NoArena: !*arena, NoPattern: !*pattern})
	if err != nil {
		_ = out.Encode(answer{Error: err.Error()})
		os.Exit(1)
	}
	defer s.Close()
	debug.FreeOSMemory()
	_ = out.Encode(answer{Ready: true, LoadMS: ms(start)})
	in := bufio.NewScanner(os.Stdin)
	in.Buffer(make([]byte, 1<<20), 1<<20)
	for in.Scan() {
		var req request
		if err := json.Unmarshal(in.Bytes(), &req); err != nil {
			_ = out.Encode(answer{Error: fmt.Sprintf("the request cannot be read: %v", err)})
			continue
		}
		if len(req.Words) == 0 {
			_ = out.Encode(answer{End: []bool{}, Prob: []float64{}})
			continue
		}
		t0 := time.Now()
		ids, last := tok.Words(req.Words)
		wide := make([]int64, len(ids))
		for i, id := range ids {
			wide[i] = int64(id)
		}
		logits, dims, err := s.Logits(wide)
		if err == nil && (len(dims) != 3 || dims[1] != int64(len(ids)) || dims[2] != punct.Labels) {
			err = fmt.Errorf("the answer has the shape %v", dims)
		}
		if err != nil {
			_ = out.Encode(answer{Error: err.Error()})
			continue
		}
		end, prob := punct.Decide(logits, last)
		_ = out.Encode(answer{End: end, Prob: prob, MS: ms(t0), Tokens: len(ids)})
	}
}

func loadSystem(name string) error {
	h, err := windows.LoadLibraryEx(name, 0, windows.LOAD_LIBRARY_SEARCH_SYSTEM32)
	if err != nil {
		return err
	}
	return windows.FreeLibrary(h)
}

func ms(t time.Time) float64 { return float64(time.Since(t).Microseconds()) / 1000 }
