package replace

import "testing"

func TestApplyPlainPhrase(t *testing.T) {
	rules := []Rule{{From: "гит хаб", To: "GitHub", Whole: true}}
	got := Apply(rules, "выложи это на гит хаб сегодня")
	if got != "выложи это на GitHub сегодня" {
		t.Fatalf("фраза не заменилась: %q", got)
	}
}

func TestApplyIgnoresCaseByDefault(t *testing.T) {
	rules := []Rule{{From: "гит хаб", To: "GitHub", Whole: true}}
	got := Apply(rules, "Гит Хаб — это удобно")
	if got != "GitHub — это удобно" {
		t.Fatalf("регистр помешал замене: %q", got)
	}
}

func TestApplyMatchCase(t *testing.T) {
	rules := []Rule{{From: "ит", To: "IT", Whole: true, MatchCase: true}}
	got := Apply(rules, "ит и Ит")
	if got != "IT и Ит" {
		t.Fatalf("с учётом регистра заменилось лишнее: %q", got)
	}
}

func TestApplyWholeWordsOnly(t *testing.T) {
	rules := []Rule{{From: "код", To: "code", Whole: true}}
	got := Apply(rules, "код в кодировке")
	if got != "code в кодировке" {
		t.Fatalf("замена залезла внутрь слова: %q", got)
	}
}

func TestApplyWithoutWholeWords(t *testing.T) {
	rules := []Rule{{From: "ё", To: "е"}}
	got := Apply(rules, "ёлки-палки")
	if got != "елки-палки" {
		t.Fatalf("замена внутри слова не сработала: %q", got)
	}
}

func TestApplyLatinBoundaries(t *testing.T) {
	rules := []Rule{{From: "js", To: "JavaScript", Whole: true}}
	got := Apply(rules, "js в jsx и js.")
	if got != "JavaScript в jsx и JavaScript." {
		t.Fatalf("границы латиницы посчитаны неверно: %q", got)
	}
}

func TestApplyOrderIsListOrder(t *testing.T) {
	rules := []Rule{
		{From: "а", To: "б"},
		{From: "б", To: "в"},
	}
	if got := Apply(rules, "а"); got != "в" {
		t.Fatalf("правила применились не по порядку: %q", got)
	}
}

func TestApplyNoLoopWhenResultContainsSource(t *testing.T) {
	rules := []Rule{{From: "код", To: "код ревью", Whole: true}}
	got := Apply(rules, "код")
	if got != "код ревью" {
		t.Fatalf("замена зациклилась или не сработала: %q", got)
	}
}

func TestApplyEmptyInputs(t *testing.T) {
	if got := Apply(nil, "текст"); got != "текст" {
		t.Fatalf("пустой список изменил текст: %q", got)
	}
	if got := Apply([]Rule{{From: "  ", To: "x"}}, "текст"); got != "текст" {
		t.Fatalf("пустой шаблон изменил текст: %q", got)
	}
	if got := Apply([]Rule{{From: "а", To: ""}}, ""); got != "" {
		t.Fatalf("пустой текст перестал быть пустым: %q", got)
	}
}

func TestCleanDropsEmptyAndTrims(t *testing.T) {
	out := Clean([]Rule{{From: "  ", To: "x"}, {From: " гит хаб ", To: " GitHub "}})
	if len(out) != 1 {
		t.Fatalf("пустая замена не выброшена: %d", len(out))
	}
	if out[0].From != "гит хаб" || out[0].To != "GitHub" {
		t.Fatalf("пробелы не убраны: %+v", out[0])
	}
}
