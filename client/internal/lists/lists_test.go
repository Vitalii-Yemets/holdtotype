package lists

import (
	"encoding/json"
	"testing"

	"holdtotype/internal/commands"
	"holdtotype/internal/replace"
)

func TestEncodeThenParseKeepsLists(t *testing.T) {
	data, err := Encode(
		[]replace.Rule{{ID: "a", From: "чат жпт", To: "ChatGPT", Whole: true}},
		[]commands.Command{{ID: "b", Phrase: "новая строка", Action: commands.ActionNewline}},
	)
	if err != nil {
		t.Fatalf("файл не собран: %v", err)
	}
	f, err := Parse(data)
	if err != nil {
		t.Fatalf("свой же файл не прочитан: %v", err)
	}
	if len(f.Replacements) != 1 || f.Replacements[0].To != "ChatGPT" {
		t.Fatalf("замены потерялись: %+v", f.Replacements)
	}
	if len(f.Commands) != 1 || f.Commands[0].Phrase != "новая строка" {
		t.Fatalf("команды потерялись: %+v", f.Commands)
	}
}

func TestEncodeDropsEmptyEntries(t *testing.T) {
	data, _ := Encode(
		[]replace.Rule{{ID: "a", From: "  ", To: "x"}},
		[]commands.Command{{ID: "b", Phrase: "", Action: commands.ActionNewline}},
	)
	f, err := Parse(data)
	if err != nil {
		t.Fatalf("файл не прочитан: %v", err)
	}
	if len(f.Replacements) != 0 || len(f.Commands) != 0 {
		t.Fatalf("пустые строки попали в файл: %+v %+v", f.Replacements, f.Commands)
	}
}

func TestParseRejectsForeignFile(t *testing.T) {
	if _, err := Parse([]byte(`{"kind":"что-то другое"}`)); err == nil {
		t.Fatal("чужой файл принят")
	}
	if _, err := Parse([]byte("не json")); err == nil {
		t.Fatal("мусор принят")
	}
}

func TestEncodeWritesEmptyArraysNotNull(t *testing.T) {
	data, err := Encode(nil, nil)
	if err != nil {
		t.Fatalf("файл не собран: %v", err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("файл не разобран: %v", err)
	}
	for _, k := range []string{"replacements", "commands"} {
		if string(raw[k]) != "[]" {
			t.Fatalf("%s записан как %s", k, raw[k])
		}
	}
}

func TestMergeRulesAddsOnlyNew(t *testing.T) {
	cur := []replace.Rule{{ID: "a", From: "чат жпт", To: "ChatGPT"}}
	add := []replace.Rule{{ID: "z", From: "ЧАТ ЖПТ", To: "другое"}, {ID: "y", From: "гит хаб", To: "GitHub"}}
	out, added, skipped := MergeRules(cur, add)
	if added != 1 || skipped != 1 {
		t.Fatalf("посчитано неверно: добавлено %d, пропущено %d", added, skipped)
	}
	if len(out) != 2 || out[0].To != "ChatGPT" || out[1].To != "GitHub" {
		t.Fatalf("список собран неверно: %+v", out)
	}
}

func TestMergeRulesKeepsIDsUnique(t *testing.T) {
	cur := []replace.Rule{{ID: "a", From: "один"}}
	add := []replace.Rule{{ID: "a", From: "два"}, {ID: "a", From: "три"}}
	out, _, _ := MergeRules(cur, add)
	seen := map[string]bool{}
	for _, r := range out {
		if seen[r.ID] {
			t.Fatalf("id повторился: %+v", out)
		}
		seen[r.ID] = true
	}
}

func TestMergeCommandsAddsOnlyNew(t *testing.T) {
	cur := []commands.Command{{ID: "a", Phrase: "точка", Action: commands.ActionText, Text: "."}}
	add := []commands.Command{{ID: "b", Phrase: "Точка", Action: commands.ActionNewline}, {ID: "c", Phrase: "отмена", Action: commands.ActionCancel}}
	out, added, skipped := MergeCommands(cur, add)
	if added != 1 || skipped != 1 {
		t.Fatalf("посчитано неверно: добавлено %d, пропущено %d", added, skipped)
	}
	if len(out) != 2 || out[1].Phrase != "отмена" {
		t.Fatalf("список собран неверно: %+v", out)
	}
}

func TestMergeDoesNotTouchOriginal(t *testing.T) {
	cur := []replace.Rule{{ID: "a", From: "один"}}
	MergeRules(cur, []replace.Rule{{ID: "b", From: "два"}})
	if len(cur) != 1 {
		t.Fatalf("исходный список изменился: %+v", cur)
	}
}
