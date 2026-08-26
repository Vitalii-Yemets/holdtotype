package preset

import "testing"

var testModels = map[string]*Model{
	"medium":  {ID: "medium", Engine: "whisper", Langs: []string{"*"}, Auto: true, Translate: true},
	"gigaam":  {ID: "gigaam", Engine: "sherpa", Langs: []string{"ru"}},
	"parakit": {ID: "parakit", Engine: "sherpa", Langs: []string{"en", "de", "ru", "uk"}, Auto: true},
	"local":   {ID: "local", Engine: "whisper", Langs: []string{"*"}},
}

func find(id string) *Model { return testModels[id] }

func known(id string) bool { return testModels[id] != nil }

func TestResolvePrefersExplicitAssignment(t *testing.T) {
	assign := map[string]string{"ru": "gigaam", "auto": "parakit"}
	if got := Resolve(assign, "ru", "medium", known); got != "gigaam" {
		t.Fatalf("получено %q, язык с назначением должен брать свою модель", got)
	}
}

func TestResolveFallsBackToUniversal(t *testing.T) {
	assign := map[string]string{"auto": "parakit"}
	if got := Resolve(assign, "de", "medium", known); got != "parakit" {
		t.Fatalf("получено %q, язык без назначения должен брать универсальную", got)
	}
}

func TestResolveFallsBackToDefault(t *testing.T) {
	if got := Resolve(nil, "fr", "medium", known); got != "medium" {
		t.Fatalf("получено %q, без назначений работает встроенный по умолчанию", got)
	}
	assign := map[string]string{"fr": "deleted-model"}
	if got := Resolve(assign, "fr", "medium", known); got != "medium" {
		t.Fatalf("получено %q, назначение на несуществующую модель игнорируется", got)
	}
}

func TestResolveNormalisesLanguage(t *testing.T) {
	assign := map[string]string{"ru": "gigaam"}
	if got := Resolve(assign, " RU ", "medium", known); got != "gigaam" {
		t.Fatalf("получено %q, регистр и пробелы не должны мешать", got)
	}
	if got := Resolve(map[string]string{"auto": "parakit"}, "", "medium", known); got != "parakit" {
		t.Fatalf("получено %q, пустой язык — это «определять сам»", got)
	}
}

func TestCanServeChecksCoverage(t *testing.T) {
	if CanServe(find("gigaam"), "en") {
		t.Fatal("русская модель не должна назначаться английскому")
	}
	if !CanServe(find("gigaam"), "ru") {
		t.Fatal("русская модель обязана назначаться русскому")
	}
	if !CanServe(find("medium"), "pl") {
		t.Fatal("универсальная модель обязана покрывать любой язык")
	}
}

func TestAutoNeedsAutoCapableModel(t *testing.T) {
	if CanServe(find("gigaam"), "auto") {
		t.Fatal("одноязычная модель не умеет определять язык")
	}
	if !CanServe(find("parakit"), "auto") {
		t.Fatal("модель с автоопределением обязана подходить для «определять сам»")
	}
	if CanServe(find("local"), "auto") {
		t.Fatal("принесённая руками модель с неизвестными свойствами не подходит для «определять сам»")
	}
}

func TestCleanDropsBrokenEntries(t *testing.T) {
	assign := map[string]string{"ru": "gigaam", "en": "gigaam", "de": "gone", "auto": "medium"}
	out, dropped := Clean(assign, find)
	if out["ru"] != "gigaam" || out["auto"] != "medium" {
		t.Fatalf("правильные назначения потеряны: %v", out)
	}
	if _, ok := out["en"]; ok {
		t.Fatal("назначение мимо языков модели должно отбрасываться")
	}
	if _, ok := out["de"]; ok {
		t.Fatal("назначение на несуществующую модель должно отбрасываться")
	}
	if len(dropped) != 2 {
		t.Fatalf("ожидались две отброшенные записи, получено %v", dropped)
	}
}
