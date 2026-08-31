package history

import (
	"os"
	"path/filepath"
	"testing"
)

const day = int64(24 * 60 * 60 * 1000)

const week = 7 * 24 * 60

func TestPruneDropsOldEntries(t *testing.T) {
	now := int64(10 * day)
	items := []Item{
		{At: now - day, Text: "вчера"},
		{At: now - 8*day, Text: "неделю с лишним назад"},
	}
	out := Prune(items, now, week, 100)
	if len(out) != 1 || out[0].Text != "вчера" {
		t.Fatalf("старая запись не выброшена: %+v", out)
	}
}

func TestPruneKeepsNewestFirstAndCapsCount(t *testing.T) {
	now := int64(10 * day)
	items := []Item{
		{At: now - 3, Text: "третья"},
		{At: now - 1, Text: "первая"},
		{At: now - 2, Text: "вторая"},
	}
	out := Prune(items, now, week, 2)
	if len(out) != 2 {
		t.Fatalf("лимит записей не сработал: %d", len(out))
	}
	if out[0].Text != "первая" || out[1].Text != "вторая" {
		t.Fatalf("порядок не от новых к старым: %+v", out)
	}
}

func TestPruneDropsEmptyText(t *testing.T) {
	now := int64(10 * day)
	out := Prune([]Item{{At: now, Text: "   "}, {At: now, Text: "есть"}}, now, week, 100)
	if len(out) != 1 {
		t.Fatalf("пустая запись сохранилась: %+v", out)
	}
}

func TestMatchIsCaseInsensitiveAndCoversApp(t *testing.T) {
	it := Item{Text: "Выложи на GitHub", App: "chrome.exe"}
	for _, q := range []string{"", "github", "ВЫЛОЖИ", "chrome"} {
		if !Match(it, q) {
			t.Fatalf("запрос %q не нашёл запись", q)
		}
	}
	if Match(it, "телеграм") {
		t.Fatal("посторонний запрос нашёл запись")
	}
}

func TestStoreRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")
	now := int64(10 * day)

	s := Open(path)
	if s.Count() != 0 {
		t.Fatal("новое хранилище не пустое")
	}
	if err := s.Add(Item{At: now, Text: "первая", App: "code.exe"}, now, week, 100); err != nil {
		t.Fatalf("запись не добавилась: %v", err)
	}
	if err := s.Add(Item{At: now + 1, Text: "вторая", App: "chrome.exe"}, now, week, 100); err != nil {
		t.Fatalf("вторая запись не добавилась: %v", err)
	}

	again := Open(path)
	if again.Count() != 2 {
		t.Fatalf("после перезапуска записей %d вместо 2", again.Count())
	}
	if again.Items()[0].Text != "вторая" {
		t.Fatalf("порядок после чтения неверный: %+v", again.Items())
	}
	found := again.Search("chrome", 10)
	if len(found) != 1 || found[0].Text != "вторая" {
		t.Fatalf("поиск по программе не сработал: %+v", found)
	}
}

func TestClearRemovesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")
	now := int64(10 * day)
	s := Open(path)
	if err := s.Add(Item{At: now, Text: "что-то"}, now, week, 100); err != nil {
		t.Fatalf("запись не добавилась: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("файл не создан: %v", err)
	}
	if err := s.Clear(); err != nil {
		t.Fatalf("очистка не удалась: %v", err)
	}
	if s.Count() != 0 {
		t.Fatal("после очистки остались записи")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("файл истории остался на диске после очистки")
	}
}

func TestSearchLimit(t *testing.T) {
	s := &Store{items: []Item{{At: 3, Text: "a"}, {At: 2, Text: "a"}, {At: 1, Text: "a"}}}
	if got := len(s.Search("a", 2)); got != 2 {
		t.Fatalf("предел выдачи не соблюдён: %d", got)
	}
}

func TestEnforceDropsWhatRetentionForbids(t *testing.T) {
	dir := t.TempDir()
	s := Open(filepath.Join(dir, "h.json"))
	now := int64(1_700_000_000_000)
	if err := s.Add(Item{At: now - 30*day, Text: "старое"}, now-30*day, 365*24*60, 100); err != nil {
		t.Fatalf("не добавилось: %v", err)
	}
	if err := s.Add(Item{At: now, Text: "свежее"}, now, 365*24*60, 100); err != nil {
		t.Fatalf("не добавилось: %v", err)
	}
	dropped, err := s.Enforce(now, week, 100)
	if err != nil {
		t.Fatalf("применение срока: %v", err)
	}
	if dropped != 1 {
		t.Fatalf("выброшено %d записей вместо одной", dropped)
	}
	if got := s.Count(); got != 1 {
		t.Fatalf("осталось %d записей вместо одной", got)
	}
	if again, _ := s.Enforce(now, week, 100); again != 0 {
		t.Fatalf("повторное применение выбросило ещё %d", again)
	}
}

func TestDeleteDropsOneEntry(t *testing.T) {
	dir := t.TempDir()
	s := Open(dir + "/hist.json")
	s.Add(Item{At: 1, Text: "first"}, 10, 7, 100)
	s.Add(Item{At: 2, Text: "second"}, 10, 7, 100)
	if err := s.Delete(1); err != nil {
		t.Fatalf("delete: %v", err)
	}
	items := s.Items()
	if len(items) != 1 || items[0].At != 2 {
		t.Fatalf("want only the second entry, got %+v", items)
	}
}
