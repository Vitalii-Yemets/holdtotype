package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type Item struct {
	At   int64  `json:"at"`
	Text string `json:"text"`
	App  string `json:"app"`
}

const (
	DefaultKeepDays = 7
	DefaultMax      = 200
	MaxAllowed      = 2000
)

func Prune(items []Item, nowMs int64, keepDays, max int) []Item {
	if keepDays <= 0 {
		keepDays = DefaultKeepDays
	}
	if max <= 0 || max > MaxAllowed {
		max = DefaultMax
	}
	cutoff := nowMs - int64(keepDays)*24*60*60*1000
	kept := make([]Item, 0, len(items))
	for _, it := range items {
		if it.At >= cutoff && strings.TrimSpace(it.Text) != "" {
			kept = append(kept, it)
		}
	}
	sort.SliceStable(kept, func(i, j int) bool { return kept[i].At > kept[j].At })
	if len(kept) > max {
		kept = kept[:max]
	}
	return kept
}

func Match(it Item, query string) bool {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return true
	}
	return strings.Contains(strings.ToLower(it.Text), q) || strings.Contains(strings.ToLower(it.App), q)
}

type Store struct {
	mu    sync.Mutex
	path  string
	items []Item
}

func Open(path string) *Store {
	s := &Store{path: path}
	data, err := os.ReadFile(path)
	if err != nil {
		return s
	}
	var items []Item
	if json.Unmarshal(data, &items) == nil {
		s.items = items
	}
	return s
}

func (s *Store) Items() []Item {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Item, len(s.items))
	copy(out, s.items)
	return out
}

func (s *Store) Search(query string, limit int) []Item {
	if limit <= 0 {
		limit = DefaultMax
	}
	out := make([]Item, 0, limit)
	for _, it := range s.Items() {
		if Match(it, query) {
			out = append(out, it)
			if len(out) >= limit {
				break
			}
		}
	}
	return out
}

func (s *Store) Add(it Item, nowMs int64, keepDays, max int) error {
	s.mu.Lock()
	s.items = Prune(append([]Item{it}, s.items...), nowMs, keepDays, max)
	s.mu.Unlock()
	return s.save()
}

func (s *Store) Clear() error {
	s.mu.Lock()
	s.items = nil
	s.mu.Unlock()
	return s.save()
}

func (s *Store) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.items)
}

// Stats counts the dictations and their characters since the given moment.
func (s *Store) Stats(sinceMs int64) (count, chars int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, it := range s.items {
		if it.At >= sinceMs {
			count++
			chars += len([]rune(it.Text))
		}
	}
	return count, chars
}

func (s *Store) save() error {
	items := s.Items()
	s.mu.Lock()
	path := s.path
	s.mu.Unlock()
	if path == "" {
		return nil
	}
	if len(items) == 0 {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	}
	out, err := json.MarshalIndent(items, "", " ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, out, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Clean(path))
}

func (s *Store) Enforce(nowMs int64, keepDays, max int) (int, error) {
	s.mu.Lock()
	before := len(s.items)
	s.items = Prune(s.items, nowMs, keepDays, max)
	dropped := before - len(s.items)
	s.mu.Unlock()
	if dropped == 0 {
		return 0, nil
	}
	return dropped, s.save()
}
