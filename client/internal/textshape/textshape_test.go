package textshape

import (
	"strings"
	"testing"
)

func oldStrip(text string) string {
	for _, m := range punctMarks {
		text = strings.ReplaceAll(text, m, "")
	}
	text = strings.Join(strings.Fields(text), " ")
	return strings.ToLower(strings.TrimSpace(text))
}

func TestSingleLineMatchesOldBehaviour(t *testing.T) {
	for _, in := range []string{
		"Привет, Мир!",
		"  Hello …  world — again; ok: done.  ",
		"Что? Нет – да.",
		"tabs\tand  spaces",
		"",
		"...",
		"- Milk, «eggs» and #tag",
	} {
		if got, want := StripPunctuation(in), oldStrip(in); got != want {
			t.Fatalf("%q: got %q, want %q", in, got, want)
		}
	}
}

func TestKeepsLineBreaks(t *testing.T) {
	cases := map[string]string{
		"Первая строка.\nВторая, строка!": "первая строка\nвторая строка",
		"А.\n\nБ.":                "а\n\nб",
		"  x  \n\n\n y ":          "x\n\ny",
		"a\r\nb":                  "a\nb",
		"\nОдна строка.\n":        "одна строка",
		"- Milk.\n- Eggs":         "- milk\n- eggs",
		"# Title.\n\nBody, text.": "# title\n\nbody text",
		"a\n—\nb":                 "a\n\nb",
	}
	for in, want := range cases {
		if got := StripPunctuation(in); got != want {
			t.Fatalf("%q: got %q, want %q", in, got, want)
		}
	}
}
