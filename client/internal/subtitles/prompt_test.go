package subtitles

import (
	"strings"
	"testing"
)

func TestTheTranslateGemmaPromptIsItsOwnTemplate(t *testing.T) {
	got := TranslateGemmaPrompt("en", "ru", []Turn{{Source: " charge  off ", Translation: "списание"}}, "So we go on")
	head := "<start_of_turn>user\nYou are a professional English (en) to Russian (ru) translator. Your goal is to accurately convey the meaning and " +
		"nuances of the original English text while adhering to Russian grammar, vocabulary, and cultural sensitivities.\n" +
		"Produce only the Russian translation, without any additional explanations or commentary. Please translate the following English text into Russian:\n\n\n"
	want := head + "So the first thing we noticed when we looked at the data<end_of_turn>\n<start_of_turn>model\nИтак, первое, что мы заметили, когда посмотрели на данные,<end_of_turn>\n" +
		head + "was that most of the users who signed up in March<end_of_turn>\n<start_of_turn>model\nэто то, что большинство пользователей, зарегистрировавшихся в марте,<end_of_turn>\n" +
		head + "never came back after the first week.<end_of_turn>\n<start_of_turn>model\nтак и не вернулись после первой недели.<end_of_turn>\n" +
		head + "charge off<end_of_turn>\n<start_of_turn>model\nсписание<end_of_turn>\n" +
		head + "So we go on<end_of_turn>\n<start_of_turn>model\n"
	if got != want {
		t.Fatalf("prompt:\n%q\nwant\n%q", got, want)
	}
}

func TestTheTranslateGemmaPromptNamesEachTurnsOwnLanguage(t *testing.T) {
	got := TranslateGemmaPrompt("de", "en", nil, "Guten Morgen")
	if !strings.HasPrefix(got, "<start_of_turn>user\nYou are a professional Spanish (es) to English (en) translator.") {
		t.Fatalf("the examples for English are Spanish: %q", got[:120])
	}
	if !strings.Contains(got, "You are a professional German (de) to English (en) translator.") || !strings.HasSuffix(got, "Guten Morgen<end_of_turn>\n<start_of_turn>model\n") {
		t.Fatalf("the piece: %q", got[len(got)-200:])
	}
	if p := TranslateGemmaPrompt("auto", "ru", nil, "x"); strings.Count(p, "English (en) to Russian (ru)") != 4 {
		t.Fatalf("an unknown speech language is taken as English: %q", p)
	}
	if TranslateGemmaPrompt("en", "xx", nil, "x") != "" {
		t.Fatal("an unknown target gives a prompt")
	}
}

func TestTheChatRequestCarriesTheContextAtTheEndOfTheSystemMessage(t *testing.T) {
	plain := TranslateMessages("ru", nil, "hello")
	with := TranslateMessagesWith("ru", " The conversation: Lending.", nil, "hello")
	if len(plain) != len(with) || with[0].Content != plain[0].Content+" The conversation: Lending." {
		t.Fatalf("system: %q", with[0].Content)
	}
	for i := 1; i < len(plain); i++ {
		if plain[i] != with[i] {
			t.Fatalf("message %d differs", i)
		}
	}
	same := TranslateMessagesWith("ru", "", nil, "hello")
	if same[0] != plain[0] {
		t.Fatal("an empty context changes the request")
	}
}
