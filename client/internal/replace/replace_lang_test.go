package replace

import "testing"

func TestForLangKeepsSharedAndOwn(t *testing.T) {
	rules := []Rule{
		{ID: "a", From: "x", To: "y"},
		{ID: "b", From: "x", To: "y", Lang: "ru"},
		{ID: "c", From: "x", To: "y", Lang: "en"},
	}
	got := ForLang(rules, "ru")
	if len(got) != 2 || got[0].ID != "a" || got[1].ID != "b" {
		t.Fatalf("получено %v — русская диктовка берёт общие и русские правила", got)
	}
}

func TestForLangAutoTakesOnlyShared(t *testing.T) {
	rules := []Rule{
		{ID: "a", From: "x", To: "y"},
		{ID: "b", From: "x", To: "y", Lang: "ru"},
	}
	got := ForLang(rules, "auto")
	if len(got) != 1 || got[0].ID != "a" {
		t.Fatalf("получено %v — при автоязыке применимы только общие правила", got)
	}
}

func TestForLangIgnoresCase(t *testing.T) {
	rules := []Rule{{ID: "b", From: "x", To: "y", Lang: "RU"}}
	if got := ForLang(rules, " ru "); len(got) != 1 {
		t.Fatalf("получено %v — регистр и пробелы не должны мешать", got)
	}
}
