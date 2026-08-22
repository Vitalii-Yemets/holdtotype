package logfilter

import (
	"bytes"
	"strings"
	"testing"
)

func TestDropsMatchingLines(t *testing.T) {
	var out bytes.Buffer
	w := New(&out, "handle_read_frame error")
	in := "первая строка\n[error] handle_read_frame error: asio 10058\nвторая строка\n"
	if _, err := w.Write([]byte(in)); err != nil {
		t.Fatalf("запись: %v", err)
	}
	got := out.String()
	if strings.Contains(got, "handle_read_frame") {
		t.Errorf("шумная строка не отфильтрована: %q", got)
	}
	if !strings.Contains(got, "первая строка") || !strings.Contains(got, "вторая строка") {
		t.Errorf("полезные строки потеряны: %q", got)
	}
}

func TestSplitWrites(t *testing.T) {
	var out bytes.Buffer
	w := New(&out, "шум")
	for _, part := range []string{"нача", "ло стро", "ки\nшум тут\nхво", "ст\n"} {
		if _, err := w.Write([]byte(part)); err != nil {
			t.Fatalf("запись: %v", err)
		}
	}
	got := out.String()
	if got != "начало строки\nхвост\n" {
		t.Fatalf("склейка сломана: %q", got)
	}
}

func TestFlushWritesTailWithoutNewline(t *testing.T) {
	var out bytes.Buffer
	w := New(&out, "шум")
	_, _ = w.Write([]byte("без перевода строки"))
	if out.Len() != 0 {
		t.Fatal("незавершённая строка не должна уходить до Flush")
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}
	if out.String() != "без перевода строки" {
		t.Fatalf("хвост потерян: %q", out.String())
	}
}

func TestFlushDropsNoisyTail(t *testing.T) {
	var out bytes.Buffer
	w := New(&out, "шум")
	_, _ = w.Write([]byte("тут шум и нет перевода строки"))
	_ = w.Flush()
	if out.Len() != 0 {
		t.Fatalf("шумный хвост должен отбрасываться, получено %q", out.String())
	}
}

func TestReturnsFullLength(t *testing.T) {
	var out bytes.Buffer
	w := New(&out, "шум")
	p := []byte("шум\nданные\n")
	n, err := w.Write(p)
	if err != nil {
		t.Fatalf("запись: %v", err)
	}
	if n != len(p) {
		t.Fatalf("вернулось %d байт из %d — вызывающий сочтёт это ошибкой", n, len(p))
	}
}

func TestNoSkipPassesEverything(t *testing.T) {
	var out bytes.Buffer
	w := New(&out)
	_, _ = w.Write([]byte("строка\n"))
	if out.String() != "строка\n" {
		t.Fatalf("без фильтров ничего не должно теряться: %q", out.String())
	}
}
