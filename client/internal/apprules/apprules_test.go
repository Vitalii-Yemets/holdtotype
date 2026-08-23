package apprules

import "testing"

func TestFindMatchesExeName(t *testing.T) {
	rules := []Rule{{ID: "a", Match: "chrome.exe", Paste: PasteType}}
	got, ok := Find(rules, `C:\Program Files\Google\Chrome\Application\chrome.exe`)
	if !ok || got.ID != "a" {
		t.Fatalf("полный путь не совпал с правилом: %v %v", got, ok)
	}
}

func TestFindIgnoresCase(t *testing.T) {
	rules := []Rule{{ID: "a", Match: "Telegram.EXE"}}
	if _, ok := Find(rules, "telegram.exe"); !ok {
		t.Fatal("регистр помешал совпадению")
	}
}

func TestFindCommaList(t *testing.T) {
	rules := []Rule{{ID: "browsers", Match: "chrome.exe, msedge.exe , firefox.exe"}}
	for _, exe := range []string{"chrome.exe", "msedge.exe", "firefox.exe"} {
		if _, ok := Find(rules, exe); !ok {
			t.Fatalf("%s не совпал со списком", exe)
		}
	}
	if _, ok := Find(rules, "notepad.exe"); ok {
		t.Fatal("посторонняя программа совпала со списком")
	}
}

func TestFindWildcard(t *testing.T) {
	rules := []Rule{{ID: "t", Match: "teams*"}}
	if _, ok := Find(rules, "ms-teams.exe"); ok {
		t.Fatal("звёздочка совпала не с началом имени")
	}
	if _, ok := Find(rules, "teamspeak.exe"); !ok {
		t.Fatal("звёздочка не поймала имя с тем же началом")
	}
}

func TestFindWithoutExtension(t *testing.T) {
	rules := []Rule{{ID: "n", Match: "notepad"}}
	if _, ok := Find(rules, "notepad.exe"); !ok {
		t.Fatal("имя без .exe не совпало")
	}
}

func TestFindFirstRuleWins(t *testing.T) {
	rules := []Rule{
		{ID: "first", Match: "code.exe"},
		{ID: "second", Match: "code.exe"},
	}
	got, _ := Find(rules, "code.exe")
	if got.ID != "first" {
		t.Fatalf("выиграло не первое правило: %s", got.ID)
	}
}

func TestFindEmpty(t *testing.T) {
	if _, ok := Find(nil, "chrome.exe"); ok {
		t.Fatal("пустой список дал совпадение")
	}
	if _, ok := Find([]Rule{{Match: "  "}}, "chrome.exe"); ok {
		t.Fatal("пустой шаблон дал совпадение")
	}
	if _, ok := Find([]Rule{{Match: "chrome.exe"}}, ""); ok {
		t.Fatal("пустое имя процесса дало совпадение")
	}
}

func TestCleanDropsEmptyAndFixesFields(t *testing.T) {
	in := []Rule{
		{Match: "  "},
		{Match: " chrome.exe ", Paste: "whatever", Enter: "maybe", DelayMs: -5, Profiles: []string{"clean"}},
		{Match: "code.exe", DelayMs: 99999, UseProfiles: true, Profiles: []string{"formal"}},
	}
	out := Clean(in)
	if len(out) != 2 {
		t.Fatalf("правило без программы не выброшено: %d", len(out))
	}
	if out[0].Match != "chrome.exe" {
		t.Fatalf("пробелы не убраны: %q", out[0].Match)
	}
	if out[0].Paste != PasteInherit || out[0].Enter != EnterInherit || out[0].DelayMs != 0 {
		t.Fatalf("мусорные значения не сброшены: %+v", out[0])
	}
	if out[0].Profiles != nil {
		t.Fatal("промпты остались, хотя правило их не переопределяет")
	}
	if out[1].DelayMs != 0 {
		t.Fatalf("слишком большая задержка не сброшена: %d", out[1].DelayMs)
	}
	if len(out[1].Profiles) != 1 {
		t.Fatal("промпты потерялись у правила, которое их переопределяет")
	}
}
