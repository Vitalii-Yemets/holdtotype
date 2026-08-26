package livetail

import (
	"strings"
	"testing"
)

func TestShortTextStaysWhole(t *testing.T) {
	if got := Tail("привет всем", 90); got != "привет всем" {
		t.Fatalf("получено %q — короткий текст не должен трогаться", got)
	}
}

func TestLongTextKeepsTheEnd(t *testing.T) {
	long := strings.Repeat("слово ", 40) + "конец"
	got := Tail(long, 30)
	if !strings.HasPrefix(got, "…") {
		t.Fatalf("получено %q — начало должно сворачиваться в многоточие", got)
	}
	if !strings.HasSuffix(got, "конец") {
		t.Fatalf("получено %q — конец фразы обязан остаться", got)
	}
	if n := len([]rune(got)); n > 30 {
		t.Fatalf("длина %d — хвост не должен превышать предел", n)
	}
}

func TestCutsAtWordBoundary(t *testing.T) {
	orig := "one two three four five six seven eight nine ten"
	got := Tail(orig, 20)
	body := strings.TrimPrefix(got, "…")
	if strings.HasPrefix(body, " ") {
		t.Fatalf("получено %q — после многоточия не должно быть пробела", got)
	}
	if !strings.HasSuffix(orig, body) {
		t.Fatalf("получено %q — хвост обязан быть концом исходной фразы", got)
	}
	at := len(orig) - len(body)
	if at > 0 && orig[at-1] != ' ' {
		t.Fatalf("получено %q — разрез пришёлся на середину слова", got)
	}
}

func TestCollapsesWhitespace(t *testing.T) {
	if got := Tail("a  b\n c", 90); got != "a b c" {
		t.Fatalf("получено %q — переносы и двойные пробелы должны схлопываться", got)
	}
}
