package hotkeys

import (
	"reflect"
	"testing"
)

func TestCleanPerLangKeepsOnlyKnownLanguagesWithReadableCombos(t *testing.T) {
	known := func(code string) bool { return code == "en" || code == "ru" || code == "auto" }
	parses := func(combo string) bool { return combo != "nonsense" }
	got, dropped := CleanPerLang(map[string]string{
		" EN ": " ctrl+alt+e ",
		"ru":   "nonsense",
		"xx":   "ctrl+alt+x",
		"auto": "",
		"":     "ctrl+alt+z",
		"Auto": "ctrl+alt+a",
	}, known, parses)
	want := map[string]string{"en": "ctrl+alt+e", "auto": "ctrl+alt+a"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("CleanPerLang kept %v, want %v", got, want)
	}
	if !reflect.DeepEqual(dropped, []string{"ru", "xx"}) {
		t.Errorf("dropped %v, want [ru xx]", dropped)
	}
}

func TestCleanPerLangNilMapGivesEmptyMap(t *testing.T) {
	got, dropped := CleanPerLang(nil, func(string) bool { return true }, func(string) bool { return true })
	if got == nil || len(got) != 0 || len(dropped) != 0 {
		t.Errorf("CleanPerLang(nil) = %v, %v; want an empty map and nothing dropped", got, dropped)
	}
}

func TestDropClashingRemovesCombosTakenElsewhere(t *testing.T) {
	got, dropped := DropClashing(map[string]string{
		"en": "ctrl+alt+e",
		"ru": "win+ctrl",
		"de": "ctrl+alt+e",
		"pl": "alt+ctrl+e",
		"uk": "ctrl+alt+u",
	}, "ctrl+win", "")
	want := map[string]string{"de": "ctrl+alt+e", "uk": "ctrl+alt+u"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("DropClashing kept %v, want %v", got, want)
	}
	if !reflect.DeepEqual(dropped, []string{"en", "pl", "ru"}) {
		t.Errorf("dropped %v, want [en pl ru]", dropped)
	}
}

func TestWithPerLangListsLanguageCombosInCodeOrder(t *testing.T) {
	got := WithPerLang([]string{"ctrl+win", ""}, map[string]string{"ru": "ctrl+alt+r", "en": "ctrl+alt+e"})
	want := []string{"ctrl+win", "", "ctrl+alt+e", "ctrl+alt+r"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WithPerLang = %v, want %v", got, want)
	}
	if dup := FindDuplicate(WithPerLang([]string{"ctrl+win"}, map[string]string{"en": "CTRL+WIN"})); dup != "ctrl+win" {
		t.Errorf("a language combo equal to the main one must count as a duplicate, got %q", dup)
	}
}
