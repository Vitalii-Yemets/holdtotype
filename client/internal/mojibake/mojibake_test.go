package mojibake

import "testing"

func mangle(s string) string {
	inv := map[byte]rune{}
	for r, b := range cp1252 {
		inv[b] = r
	}
	out := make([]rune, 0, len(s))
	for _, b := range []byte(s) {
		if r, ok := inv[b]; ok {
			out = append(out, r)
			continue
		}
		out = append(out, rune(b))
	}
	return string(out)
}

func TestFixUndoesAnEditorsMistake(t *testing.T) {
	for _, want := range []string{
		"Чистка речи",
		"Словарь, термины и сокращения",
		"für Größe",
		"тире — и «кавычки»",
	} {
		t.Run(want, func(t *testing.T) {
			in := mangle(want)
			if in == want {
				t.Fatalf("test data is not mangled: %q", in)
			}
			if got := Fix(in); got != want {
				t.Fatalf("Fix(%q) = %q, want %q", in, got, want)
			}
		})
	}
}

func TestFixLeavesHealthyTextAlone(t *testing.T) {
	for _, s := range []string{
		"",
		"Whisper, Docker, Go",
		"Чистка речи",
		"тире — и «кавычки»",
		"Чистка Ð§Ð¸",
		"Café",
		"für Größe",
	} {
		t.Run(s, func(t *testing.T) {
			if got := Fix(s); got != s {
				t.Fatalf("Fix(%q) = %q, want it unchanged", s, got)
			}
		})
	}
}

func TestFixIsStableOnRepairedText(t *testing.T) {
	once := Fix(mangle("Словарь"))
	if once != "Словарь" {
		t.Fatalf("first pass = %q", once)
	}
	if twice := Fix(once); twice != once {
		t.Fatalf("second pass changed the text: %q", twice)
	}
}
