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

func TestCancelAloneOrAtEndOnly(t *testing.T) {
	cmds := []Command{{Phrase: "cancel", Action: ActionCancel}}
	for _, in := range []string{"Cancel.", "  cancel ", "that was wrong, cancel!"} {
		if got := Apply(cmds, in); !got.Cancelled || len(got.Applied) != 1 {
			t.Fatalf("%q: cancel did not fire: %+v", in, got)
		}
	}
	got := Apply(cmds, "please cancel my order")
	if got.Cancelled || got.Text != "please cancel my order" {
		t.Fatalf("cancel in the middle threw the dictation away: %+v", got)
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

func TestCleanKeepsParametersPerAction(t *testing.T) {
	keep := map[string]bool{
		ActionNewline: false, ActionParagraph: false, ActionText: true, ActionCancel: false,
		ActionUndo: false, ActionHeading: true, ActionList: true, ActionCapital: false,
		ActionQuotes: true, ActionNoStop: false,
	}
	for action, want := range keep {
		if !ValidAction(action) {
			t.Fatalf("%s is not a valid action", action)
		}
		out := Clean([]Command{{Phrase: "p", Action: action, Text: "«»"}})
		if got := out[0].Text != ""; got != want {
			t.Fatalf("%s: parameter kept=%v, want %v", action, got, want)
		}
		if out[0].Action != action {
			t.Fatalf("%s changed to %s", action, out[0].Action)
		}
	}
	if ValidAction("bold") {
		t.Fatal("unknown action accepted")
	}
}

func TestUndoOnlyWhenAlone(t *testing.T) {
	cmds := []Command{{Phrase: "delete that", Action: ActionUndo}}
	got := Apply(cmds, "Delete that.")
	if !got.Undo || got.Text != "" || got.Cancelled || len(got.Applied) != 1 {
		t.Fatalf("undo alone did not fire: %+v", got)
	}
	for _, in := range []string{"delete that please", "please delete that"} {
		got := Apply(cmds, in)
		if got.Undo || got.Text != in || len(got.Applied) != 0 {
			t.Fatalf("%q: undo fired inside a longer dictation: %+v", in, got)
		}
	}
}

func shapeCmds() []Command {
	return []Command{
		{Phrase: "cancel", Action: ActionCancel},
		{Phrase: "new line", Action: ActionNewline},
		{Phrase: "as a list", Action: ActionList},
		{Phrase: "as heading", Action: ActionHeading},
		{Phrase: "in quotes", Action: ActionQuotes, Text: "«»"},
		{Phrase: "no full stop", Action: ActionNoStop},
		{Phrase: "capital letter", Action: ActionCapital},
	}
}

func actions(cs []Command) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = c.Action
	}
	return out
}

func sameActions(cs []Command, want ...string) bool {
	got := actions(cs)
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

func TestShapeAtEndIsCutWithPunctuation(t *testing.T) {
	got := Apply(shapeCmds(), "Buy milk, eggs and bread as a list.")
	if got.Text != "Buy milk, eggs and bread" || !sameActions(got.Shapes, ActionList) {
		t.Fatalf("list at the end not recognised: %+v", got)
	}
	if len(got.Applied) != 1 || got.Applied[0] != "as a list" {
		t.Fatalf("shape not reported as applied: %+v", got.Applied)
	}
	got = Apply(shapeCmds(), "Buy bread. — As a list")
	if got.Text != "Buy bread." {
		t.Fatalf("full stop before the shape should stay, dashes should go: %q", got.Text)
	}
}

func TestShapeNotAtEndIsText(t *testing.T) {
	in := "put it as a list of things"
	got := Apply(shapeCmds(), in)
	if got.Text != in || len(got.Shapes) != 0 || len(got.Applied) != 0 {
		t.Fatalf("shape fired in the middle: %+v", got)
	}
}

func TestSeveralShapesInApplyOrder(t *testing.T) {
	got := Apply(shapeCmds(), "Hello world in quotes, no full stop.")
	if got.Text != "Hello world" || !sameActions(got.Shapes, ActionNoStop, ActionQuotes) {
		t.Fatalf("two shapes: %+v", got)
	}
	if got.Shapes[1].Text != "«»" {
		t.Fatalf("quote marks lost: %+v", got.Shapes[1])
	}
	got = Apply(shapeCmds(), "x as a list capital letter in quotes no full stop")
	if got.Text != "x as a list" || !sameActions(got.Shapes, ActionNoStop, ActionCapital, ActionQuotes) {
		t.Fatalf("more than three shapes taken: %+v", got)
	}
}

func TestShapeIgnoresCaseAndHyphens(t *testing.T) {
	got := Apply(shapeCmds(), "Done. No Full-Stop")
	if got.Text != "Done." || !sameActions(got.Shapes, ActionNoStop) {
		t.Fatalf("case or hyphen broke the match: %+v", got)
	}
}

func TestShapeAloneArmsNextDictation(t *testing.T) {
	got := Apply(shapeCmds(), "As a list.")
	if got.Text != "" || !sameActions(got.Shapes, ActionList) || len(got.Applied) != 1 {
		t.Fatalf("shape alone: %+v", got)
	}
}

func TestStructureNearerEndWins(t *testing.T) {
	got := Apply(shapeCmds(), "Plans as a list as heading")
	if got.Text != "Plans" || !sameActions(got.Shapes, ActionHeading) || len(got.Applied) != 2 {
		t.Fatalf("heading said last should win: %+v", got)
	}
}

func TestCancelBeatsShapes(t *testing.T) {
	got := Apply(shapeCmds(), "Buy milk as a list, cancel.")
	if !got.Cancelled || len(got.Shapes) != 0 {
		t.Fatalf("cancel lost to a shape: %+v", got)
	}
}

func TestInlineCommandsStillWorkWithShapes(t *testing.T) {
	got := Apply(shapeCmds(), "first new line second as a list")
	if got.Text != "first\nsecond" || !sameActions(got.Shapes, ActionList) {
		t.Fatalf("newline with a shape: %+v", got)
	}
}

func TestShapeHeading(t *testing.T) {
	h := []Command{{Action: ActionHeading}}
	if got := Shape("Weekly report.", h, ShapeOptions{}); got != "# Weekly report\n\n" {
		t.Fatalf("default heading: %q", got)
	}
	h[0].Text = "## "
	if got := Shape("Part one\nnotes.", h, ShapeOptions{}); got != "## Part one notes\n\n" {
		t.Fatalf("heading with a marker: %q", got)
	}
	if got := Shape("Is it done?", h, ShapeOptions{}); got != "## Is it done?\n\n" {
		t.Fatalf("question mark lost: %q", got)
	}
}

func TestShapeList(t *testing.T) {
	l := []Command{{Action: ActionList}}
	en := ShapeOptions{Lang: "en"}
	cases := []struct {
		in, want string
		opts     ShapeOptions
	}{
		{"Shopping: milk, eggs and bread.", "Shopping:\n- milk\n- eggs\n- bread", en},
		{"tape, rubber band", "- tape\n- rubber band", en},
		{"tea and coffee", "- tea\n- coffee", en},
		{"salt and pepper, bread and butter", "- salt and pepper\n- bread\n- butter", en},
		{"milk, eggs, and bread.", "- milk\n- eggs\n- bread", en},
		{"tea and coffee", "- tea and coffee", ShapeOptions{Lang: "ja"}},
		{"чай и кофе", "- чай\n- кофе", ShapeOptions{Lang: "ru-RU"}},
		{"хліб та сіль", "- хліб\n- сіль", ShapeOptions{Lang: "uk"}},
		{"Milk, Eggs! and bread.", "- milk\n- eggs\n- bread", ShapeOptions{Lang: "en", StripPunct: true}},
		{"meet at 10:30, lunch", "- meet at 10:30\n- lunch", en},
	}
	for _, c := range cases {
		if got := Shape(c.in, l, c.opts); got != c.want {
			t.Fatalf("%q: got %q, want %q", c.in, got, c.want)
		}
	}
	custom := []Command{{Action: ActionList, Text: "plus, as well as"}}
	if got := Shape("a, b as well as c", custom, en); got != "- a\n- b\n- c" {
		t.Fatalf("own joining words: %q", got)
	}
	if got := Shape("a and b plus c", custom, en); got != "- a and b\n- c" {
		t.Fatalf("own joining words replace the table: %q", got)
	}
}

func TestShapeCapital(t *testing.T) {
	c := []Command{{Action: ActionCapital}}
	if got := Shape("«привет», сказал он", c, ShapeOptions{}); got != "«Привет», сказал он" {
		t.Fatalf("capital after a quote: %q", got)
	}
	if got := Shape("Hello, World.", c, ShapeOptions{StripPunct: true}); got != "Hello world" {
		t.Fatalf("capital after stripping: %q", got)
	}
	both := []Command{{Action: ActionCapital}, {Action: ActionList}}
	if got := Shape("milk, eggs and bread", both, ShapeOptions{Lang: "en"}); got != "- Milk\n- Eggs\n- Bread" {
		t.Fatalf("capital in a list: %q", got)
	}
}

func TestShapeQuotes(t *testing.T) {
	cases := []struct{ marks, in, want string }{
		{"", "Hello world.", "“Hello world”"},
		{"„“", "Hallo.", "„Hallo“"},
		{"«  »", "Bonjour.", "« Bonjour »"},
		{"", "Really?", "“Really?”"},
		{"« »", "oui", "«oui»"},
	}
	for _, c := range cases {
		q := []Command{{Action: ActionQuotes, Text: c.marks}}
		if got := Shape(c.in, q, ShapeOptions{}); got != c.want {
			t.Fatalf("%q with %q: got %q, want %q", c.in, c.marks, got, c.want)
		}
	}
}

func TestShapeNoStop(t *testing.T) {
	n := []Command{{Action: ActionNoStop}}
	cases := map[string]string{"Done.": "Done", "Done...": "Done", "Done…": "Done", "Really?": "Really?", "Wow!": "Wow!"}
	for in, want := range cases {
		if got := Shape(in, n, ShapeOptions{}); got != want {
			t.Fatalf("%q: got %q, want %q", in, got, want)
		}
	}
}

func TestShapeOrder(t *testing.T) {
	all := []Command{{Action: ActionNoStop}, {Action: ActionCapital}, {Action: ActionQuotes}, {Action: ActionHeading}}
	if got := Shape("hello.", all, ShapeOptions{}); got != "# “Hello”\n\n" {
		t.Fatalf("all shapes together: %q", got)
	}
	if got := Shape("", all, ShapeOptions{}); got != "" {
		t.Fatalf("empty text shaped: %q", got)
	}
}

func TestApplyThenShape(t *testing.T) {
	res := Apply(shapeCmds(), "Buy milk, eggs and bread as a list.")
	if got := Shape(res.Text, res.Shapes, ShapeOptions{Lang: "en-US"}); got != "- Buy milk\n- eggs\n- bread" {
		t.Fatalf("dictated list: %q", got)
	}
}
