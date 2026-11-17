package trmodel

import (
	"path/filepath"
	"sort"
	"strings"
)

const (
	PromptChat      = "chat"
	PromptTranslate = "translategemma"

	GemmaTermsURL  = "https://ai.google.dev/gemma/terms"
	ApacheTermsURL = "https://www.apache.org/licenses/LICENSE-2.0"
)

type Model struct {
	ID          string
	Name        string
	File        string
	URL         string
	Size        int64
	SHA256      string
	License     string
	TermsURL    string
	CardMB      int
	FreeMB      int
	VramMB      int
	Prompt      string
	ChatArgs    []string
	Topic       bool
	Terms       bool
	Recommended bool
	Accuracy    int
	Speed       int
	Own         bool
}

var Catalog = []Model{
	{
		ID: "gemma-4-e4b", Name: "Gemma 4 E4B", File: "gemma-4-E4B_q4_0-it.gguf",
		URL:  "https://huggingface.co/google/gemma-4-E4B-it-qat-q4_0-gguf/resolve/4b4a2c1d584be7264f87aac328a1bc739ce81b6c/gemma-4-E4B_q4_0-it.gguf",
		Size: 5154941280, SHA256: "676c35070db6dbe52f93e9c864ee0fba4eddea94b9c875d9cb10daff453fbaee",
		License: "Apache 2.0", TermsURL: ApacheTermsURL,
		CardMB: 4096, FreeMB: 3300, VramMB: 3015, Prompt: PromptChat, Topic: true, Terms: true, Recommended: true, Accuracy: 5, Speed: 5,
	},
	{
		ID: "gemma-3-4b", Name: Name, File: File, URL: URL, Size: Size, SHA256: SHA256,
		License: "Gemma", TermsURL: TermsURL,
		CardMB: 4096, FreeMB: 3300, VramMB: 3022, Prompt: PromptChat, Topic: true, Terms: true, Accuracy: 4, Speed: 5,
	},
	{
		ID: "translategemma-4b", Name: "TranslateGemma 4B", File: "translategemma-4b-it-Q4_K_M.gguf",
		URL:  "https://huggingface.co/bullerwins/translategemma-4b-it-GGUF/resolve/7c938465a870d8624bcfa98a8e4a3510053c19a8/translategemma-4b-it-Q4_K_M.gguf",
		Size: 2489909312, SHA256: "7f7357c14abd9da4eb200b38b05da502cd6e10d7e1d403fbc9f78c19f3209b72",
		License: "Gemma", TermsURL: GemmaTermsURL,
		CardMB: 4096, FreeMB: 3300, VramMB: 3022, Prompt: PromptTranslate, ChatArgs: []string{"--chat-template", "gemma"}, Accuracy: 2, Speed: 5,
	},
	{
		ID: "eurollm-9b", Name: "EuroLLM-9B", File: "EuroLLM-9B-Instruct-2512.Q4_K_M.gguf",
		URL:  "https://huggingface.co/mradermacher/EuroLLM-9B-Instruct-2512-GGUF/resolve/5138603fccb9dc1f35c25cfd1a4981a78617e9ff/EuroLLM-9B-Instruct-2512.Q4_K_M.gguf",
		Size: 5582838912, SHA256: "fe9b50d4ba67eb0131f1ceeaab2e56ea8fb6b2bbc37391d3219c791c3b140a28",
		License: "Apache 2.0", TermsURL: ApacheTermsURL,
		CardMB: 8192, FreeMB: 6100, VramMB: 5828, Prompt: PromptChat, Topic: true, Terms: true, Accuracy: 3, Speed: 4,
	},
}

func (m Model) Spec() Spec { return Spec{Name: m.File, Size: m.Size, Hash: m.SHA256} }

func (m Model) SizeMB() int { return int(m.Size / (1024 * 1024)) }

func (m Model) Key() string {
	if m.Own {
		return "own:" + m.File
	}
	return m.ID
}

func ByID(id string) (Model, bool) {
	for _, m := range Catalog {
		if m.ID == id {
			return m, true
		}
	}
	return Model{}, false
}

func ByFile(name string) (Model, bool) {
	base := filepath.Base(name)
	for _, m := range Catalog {
		if strings.EqualFold(m.File, base) {
			return m, true
		}
	}
	return Model{}, false
}

func Recommended() Model {
	for _, m := range Catalog {
		if m.Recommended {
			return m
		}
	}
	return Catalog[0]
}

func OwnMB(sizeMB int) (card, free int) {
	free = sizeMB + 1000
	return free + 256, free
}

func Own(file string, sizeMB int) Model {
	card, free := OwnMB(sizeMB)
	return Model{
		ID: "own:" + file, Name: strings.TrimSuffix(file, filepath.Ext(file)), File: file, Size: int64(sizeMB) * 1024 * 1024,
		CardMB: card, FreeMB: free, VramMB: free, Prompt: PromptChat, Own: true,
	}
}

func OwnChoice(choice string) (string, bool) {
	if strings.HasPrefix(choice, "own:") {
		f := strings.TrimPrefix(choice, "own:")
		return f, f != "" && !strings.ContainsAny(f, `/\`)
	}
	return "", false
}

func ValidChoice(choice string) bool {
	if choice == "" {
		return true
	}
	if f, ok := OwnChoice(choice); ok {
		return strings.HasSuffix(strings.ToLower(f), ".gguf") && !Is(f)
	}
	_, ok := ByID(choice)
	return ok
}

type OwnFile struct {
	Name   string
	SizeMB int
}

func Resolve(choice string, ready func(Model) bool, own []OwnFile) (Model, bool) {
	if f, ok := OwnChoice(choice); ok {
		for _, o := range own {
			if strings.EqualFold(o.Name, f) {
				return Own(o.Name, o.SizeMB), true
			}
		}
	} else if m, ok := ByID(choice); ok && ready(m) {
		return m, true
	}
	if r := Recommended(); ready(r) {
		return r, true
	}
	for _, m := range Catalog {
		if ready(m) {
			return m, true
		}
	}
	return Model{}, false
}

func Shown(choice string, ready func(Model) bool, own []OwnFile) Model {
	if m, ok := Resolve(choice, ready, own); ok {
		return m
	}
	if m, ok := ByID(choice); ok {
		return m
	}
	return Recommended()
}

func SortOwn(list []OwnFile) []OwnFile {
	out := append([]OwnFile(nil), list...)
	sort.Slice(out, func(i, j int) bool { return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name) })
	return out
}
