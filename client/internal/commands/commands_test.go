package commands

import "testing"

func TestNewlineReplacesPhraseAndTrimsSpaces(t *testing.T) {
	cmds := []Command{{Phrase: "новая строка", Action: ActionNewline}}
	got := Apply(cmds, "привет новая строка как дела")
	if got.Text != "привет\nкак дела" {
		t.Fatalf("перенос строки собран неверно: %q", got.Text)
	}
	if got.Cancelled {
		t.Fatal("команда переноса отменила диктовку")
	}
	if len(got.Applied) != 1 {
		t.Fatalf("сработавшая команда не отмечена: %+v", got.Applied)
	}
}

func TestParagraph(t *testing.T) {
	cmds := []Command{{Phrase: "новый абзац", Action: ActionParagraph}}
	got := Apply(cmds, "первый новый абзац второй")
	if got.Text != "первый\n\nвторой" {
		t.Fatalf("абзац собран неверно: %q", got.Text)
	}
}

func TestCancelStopsEverything(t *testing.T) {
	cmds := []Command{
		{Phrase: "новая строка", Action: ActionNewline},
		{Phrase: "отмена", Action: ActionCancel},
	}
	got := Apply(cmds, "что-то не то отмена")
	if !got.Cancelled {
		t.Fatalf("команда отмены не сработала: %+v", got)
	}
}

func TestCancelNeedsWholeWord(t *testing.T) {
	cmds := []Command{{Phrase: "отмена", Action: ActionCancel}}
	if got := Apply(cmds, "отменена подписка"); got.Cancelled {
		t.Fatal("отмена сработала внутри другого слова")
	}
}

func TestTextAction(t *testing.T) {
	cmds := []Command{{Phrase: "смайлик", Action: ActionText, Text: ":)"}}
	got := Apply(cmds, "спасибо смайлик")
	if got.Text != "спасибо :)" {
		t.Fatalf("подстановка текста не сработала: %q", got.Text)
	}
}

func TestOrderTopToBottom(t *testing.T) {
	cmds := []Command{
		{Phrase: "раз", Action: ActionText, Text: "два"},
		{Phrase: "два", Action: ActionText, Text: "три"},
	}
	if got := Apply(cmds, "раз"); got.Text != "три" {
		t.Fatalf("команды применились не по порядку: %q", got.Text)
	}
}

func TestNothingHappensWithoutCommands(t *testing.T) {
	if got := Apply(nil, "просто текст"); got.Text != "просто текст" || got.Cancelled {
		t.Fatalf("пустой список изменил текст: %+v", got)
	}
	if got := Apply([]Command{{Phrase: "  ", Action: ActionNewline}}, "текст"); got.Text != "текст" {
		t.Fatalf("пустая команда изменила текст: %+v", got)
	}
}

func TestTidyCollapsesRepeats(t *testing.T) {
	cmds := []Command{{Phrase: "нс", Action: ActionNewline}}
	got := Apply(cmds, "нс нс а нс")
	if got.Text != "а" {
		t.Fatalf("лишние переносы не убраны: %q", got.Text)
	}
}

func TestCleanFixesActionsAndDropsEmpty(t *testing.T) {
	out := Clean([]Command{
		{Phrase: "  "},
		{Phrase: " отмена ", Action: "странное"},
		{Phrase: "смайлик", Action: ActionText, Text: ":)"},
		{Phrase: "нс", Action: ActionNewline, Text: "мусор"},
	})
	if len(out) != 3 {
		t.Fatalf("пустая команда не выброшена: %d", len(out))
	}
	if out[0].Phrase != "отмена" || out[0].Action != ActionNewline {
		t.Fatalf("неизвестное действие не заменено: %+v", out[0])
	}
	if out[1].Text != ":)" {
		t.Fatal("текст подстановки потерялся")
	}
	if out[2].Text != "" {
		t.Fatal("лишний текст у команды переноса не очищен")
	}
}
