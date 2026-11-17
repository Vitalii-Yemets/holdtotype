package punct

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

const space = "▁"

type Tokenizer struct {
	vocab  map[string]int32
	merges map[uint64]merge
	Bos    int32
	Eos    int32
	mu     sync.Mutex
	cache  map[string][]int32
}

type merge struct {
	rank int32
	id   int32
}

type tokenizerFile struct {
	AddedTokens []struct {
		ID      int32  `json:"id"`
		Content string `json:"content"`
	} `json:"added_tokens"`
	Model struct {
		Type   string            `json:"type"`
		Vocab  map[string]int32  `json:"vocab"`
		Merges []json.RawMessage `json:"merges"`
	} `json:"model"`
}

func LoadTokenizer(path string) (*Tokenizer, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseTokenizer(raw)
}

func ParseTokenizer(raw []byte) (*Tokenizer, error) {
	var f tokenizerFile
	if err := json.Unmarshal(raw, &f); err != nil {
		return nil, fmt.Errorf("the tokenizer file cannot be read: %w", err)
	}
	if f.Model.Type != "BPE" {
		return nil, fmt.Errorf("the tokenizer is %q, not BPE", f.Model.Type)
	}
	t := &Tokenizer{vocab: f.Model.Vocab, merges: make(map[uint64]merge, len(f.Model.Merges)), cache: map[string][]int32{}}
	for _, a := range f.AddedTokens {
		if _, ok := t.vocab[a.Content]; !ok {
			t.vocab[a.Content] = a.ID
		}
	}
	for rank, m := range f.Model.Merges {
		a, b, err := mergePair(m)
		if err != nil {
			return nil, err
		}
		ia, oka := t.vocab[a]
		ib, okb := t.vocab[b]
		id, okc := t.vocab[a+b]
		if !oka || !okb || !okc {
			continue
		}
		key := pairKey(ia, ib)
		if _, seen := t.merges[key]; !seen {
			t.merges[key] = merge{rank: int32(rank), id: id}
		}
	}
	var ok1, ok2 bool
	t.Bos, ok1 = t.vocab["<bos>"]
	t.Eos, ok2 = t.vocab["<eos>"]
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("the tokenizer has no <bos> or <eos>")
	}
	return t, nil
}

func mergePair(m json.RawMessage) (string, string, error) {
	var pair []string
	if err := json.Unmarshal(m, &pair); err == nil && len(pair) == 2 {
		return pair[0], pair[1], nil
	}
	var line string
	if err := json.Unmarshal(m, &line); err == nil {
		if a, b, ok := strings.Cut(line, " "); ok {
			return a, b, nil
		}
	}
	return "", "", fmt.Errorf("a merge of the tokenizer cannot be read: %s", m)
}

func pairKey(a, b int32) uint64 { return uint64(uint32(a))<<32 | uint64(uint32(b)) }

func (t *Tokenizer) Word(w string) []int32 {
	t.mu.Lock()
	ids, ok := t.cache[w]
	t.mu.Unlock()
	if ok {
		return ids
	}
	ids = t.bpe(space + strings.ReplaceAll(w, " ", space))
	t.mu.Lock()
	t.cache[w] = ids
	t.mu.Unlock()
	return ids
}

func (t *Tokenizer) bpe(s string) []int32 {
	var sym []int32
	for _, r := range s {
		ch := string(r)
		if id, ok := t.vocab[ch]; ok {
			sym = append(sym, id)
			continue
		}
		for _, b := range []byte(ch) {
			if id, ok := t.vocab[fmt.Sprintf("<0x%02X>", b)]; ok {
				sym = append(sym, id)
			}
		}
	}
	for len(sym) > 1 {
		best, at := int32(-1), -1
		var to int32
		for i := 0; i+1 < len(sym); i++ {
			m, ok := t.merges[pairKey(sym[i], sym[i+1])]
			if ok && (best < 0 || m.rank < best) {
				best, at, to = m.rank, i, m.id
			}
		}
		if at < 0 {
			break
		}
		sym[at] = to
		sym = append(sym[:at+1], sym[at+2:]...)
	}
	return sym
}

func (t *Tokenizer) Words(words []string) ([]int32, []int) {
	ids := []int32{t.Bos}
	last := make([]int, len(words))
	for i, w := range words {
		ids = append(ids, t.Word(w)...)
		last[i] = len(ids) - 1
	}
	ids = append(ids, t.Eos)
	return ids, last
}
