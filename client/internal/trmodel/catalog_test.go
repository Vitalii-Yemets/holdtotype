package trmodel

import "testing"

func readyOf(ids ...string) func(Model) bool {
	return func(m Model) bool {
		for _, id := range ids {
			if m.ID == id {
				return true
			}
		}
		return false
	}
}

func TestTheCatalogKnowsFourCheckedModels(t *testing.T) {
	if len(Catalog) != 4 {
		t.Fatalf("models: %d", len(Catalog))
	}
	rec := 0
	for _, m := range Catalog {
		if m.Recommended {
			rec++
		}
		if m.File == "" || m.URL == "" || m.Size <= 0 || len(m.SHA256) != 64 || m.TermsURL == "" || m.License == "" {
			t.Fatalf("%s lacks a field: %+v", m.ID, m)
		}
		if m.Accuracy < 1 || m.Accuracy > 5 || m.Speed < 1 || m.Speed > 5 {
			t.Fatalf("%s bars: %d %d", m.ID, m.Accuracy, m.Speed)
		}
		if m.Prompt != PromptChat && m.Prompt != PromptTranslate {
			t.Fatalf("%s prompt: %q", m.ID, m.Prompt)
		}
		if m.FreeMB < m.VramMB || m.CardMB < m.FreeMB {
			t.Fatalf("%s memory: card %d free %d vram %d", m.ID, m.CardMB, m.FreeMB, m.VramMB)
		}
	}
	if rec != 1 || Recommended().ID != "gemma-4-e4b" {
		t.Fatalf("recommended: %d %s", rec, Recommended().ID)
	}
	g3, ok := ByID("gemma-3-4b")
	if !ok || g3.File != File || g3.SHA256 != SHA256 || g3.URL != URL {
		t.Fatalf("Gemma 3 4B changed: %+v", g3)
	}
	tg, _ := ByID("translategemma-4b")
	if tg.Service || tg.Prompt != PromptTranslate || len(tg.ChatArgs) != 2 || tg.ChatArgs[1] != "gemma" {
		t.Fatalf("TranslateGemma: %+v", tg)
	}
	eu, _ := ByID("eurollm-9b")
	if eu.CardMB != 8192 || !eu.Service {
		t.Fatalf("EuroLLM: %+v", eu)
	}
}

func TestEveryCatalogFileIsATranslationModel(t *testing.T) {
	for _, m := range Catalog {
		if !Is(m.File) || !Is("models/"+m.File) {
			t.Fatalf("%s is not recognised", m.File)
		}
	}
	if Is("qwen2.5-1.5b-instruct-q4_k_m.gguf") {
		t.Fatal("an editor model is taken for a translation model")
	}
}

func TestTheChosenModelFollowsTheRules(t *testing.T) {
	own := []OwnFile{{Name: "my-model.gguf", SizeMB: 4000}}
	cases := []struct {
		choice string
		ready  []string
		want   string
		ok     bool
	}{
		{"", []string{"gemma-4-e4b", "gemma-3-4b"}, "gemma-4-e4b", true},
		{"", []string{"gemma-3-4b"}, "gemma-3-4b", true},
		{"", []string{"eurollm-9b", "translategemma-4b"}, "translategemma-4b", true},
		{"eurollm-9b", []string{"gemma-4-e4b", "eurollm-9b"}, "eurollm-9b", true},
		{"eurollm-9b", []string{"gemma-3-4b"}, "gemma-3-4b", true},
		{"nonsense", []string{"gemma-3-4b"}, "gemma-3-4b", true},
		{"own:my-model.gguf", nil, "own:my-model.gguf", true},
		{"own:MY-MODEL.gguf", nil, "own:my-model.gguf", true},
		{"own:gone.gguf", []string{"gemma-4-e4b"}, "gemma-4-e4b", true},
		{"own:../x.gguf", nil, "", false},
		{"", nil, "", false},
	}
	for _, c := range cases {
		m, ok := Resolve(c.choice, readyOf(c.ready...), own)
		if ok != c.ok || (ok && m.Key() != c.want) {
			t.Fatalf("%q with %v: got %q %v, want %q %v", c.choice, c.ready, m.Key(), ok, c.want, c.ok)
		}
	}
	m, _ := Resolve("own:my-model.gguf", readyOf(), own)
	if !m.Own || m.Prompt != PromptChat || m.Service || m.FreeMB != 5000 || m.CardMB < m.FreeMB {
		t.Fatalf("own model: %+v", m)
	}
}

func TestOnlyKnownChoicesAreSaved(t *testing.T) {
	for _, c := range []string{"", "gemma-4-e4b", "eurollm-9b", "own:my-model.gguf"} {
		if !ValidChoice(c) {
			t.Fatalf("%q is refused", c)
		}
	}
	for _, c := range []string{"gemma-5", "own:", "own:../x.gguf", `own:a\b.gguf`, "own:notes.txt", "own:" + File} {
		if ValidChoice(c) {
			t.Fatalf("%q is taken", c)
		}
	}
}

func TestTheShownModelIsTheChosenOneEvenWhenItIsMissing(t *testing.T) {
	if m := Shown("eurollm-9b", readyOf(), nil); m.ID != "eurollm-9b" {
		t.Fatalf("got %s", m.ID)
	}
	if m := Shown("", readyOf(), nil); m.ID != "gemma-4-e4b" {
		t.Fatalf("got %s", m.ID)
	}
	if m := Shown("", readyOf("gemma-3-4b"), nil); m.ID != "gemma-3-4b" {
		t.Fatalf("got %s", m.ID)
	}
}
