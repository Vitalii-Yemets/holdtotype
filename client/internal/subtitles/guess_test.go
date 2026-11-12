package subtitles

import (
	"strings"
	"testing"
	"time"
)

const guessEnglish = "So the first thing we noticed when we looked at the data was that most of the users who signed up in March " +
	"never came back after the first week, and that is why we started to look at what happened during the first session. " +
	"It turned out that the setup took too long for people who came from a phone, because the form asked for things " +
	"they did not have at hand, so we cut it down to two fields and moved the rest into the settings, where it can wait " +
	"until somebody actually needs it. After that change the share of people who stayed for a second week went up by about a third, " +
	"which is more than we expected from such a small change."

const guessRussian = "Итак, первое, что мы заметили, когда посмотрели на данные, это то, что большинство пользователей, которые " +
	"зарегистрировались в марте, так и не вернулись после первой недели, и поэтому мы стали смотреть, что происходило во время " +
	"первого входа. Оказалось, что настройка занимала слишком много времени у тех, кто приходил с телефона, потому что форма " +
	"спрашивала то, чего у них не было под рукой, и мы сократили её до двух полей, а остальное перенесли в настройки, где оно " +
	"может подождать, пока кому-то не понадобится. После этого изменения доля тех, кто остался на вторую неделю, выросла " +
	"примерно на треть, что больше, чем мы ждали от такой небольшой правки."

const guessUkrainian = "Отже, перше, що ми помітили, коли подивилися на дані, це те, що більшість користувачів, які " +
	"зареєструвалися в березні, так і не повернулися після першого тижня, і тому ми почали дивитися, що відбувалося під час " +
	"першого входу. Виявилося, що налаштування займало забагато часу в тих, хто приходив із телефона, бо форма питала те, " +
	"чого в них не було під рукою, і ми скоротили її до двох полів, а решту перенесли до налаштувань, де воно може почекати, " +
	"поки комусь не знадобиться. Після цієї зміни частка тих, хто залишився на другий тиждень, зросла приблизно на третину, " +
	"що більше, ніж ми чекали від такої невеликої правки."

func TestGuessTheLanguageFromTheWords(t *testing.T) {
	cases := []struct {
		name, text, want string
	}{
		{"english", guessEnglish, "en"},
		{"russian", guessRussian, "ru"},
		{"ukrainian", guessUkrainian, "uk"},
		{"a short beginning", "So the first thing we noticed", ""},
		{"nothing", "", ""},
		{"numbers only", "one two three four five six seven eight nine ten eleven twelve", ""},
	}
	for _, c := range cases {
		if got := GuessLanguage(strings.Fields(c.text)); got != c.want {
			t.Errorf("%s: %q, want %q", c.name, got, c.want)
		}
	}
}

func TestGuessFollowsASwitchOfLanguage(t *testing.T) {
	en := strings.Fields(guessEnglish)
	ru := strings.Fields(guessRussian)
	if got := GuessLanguage(append(append([]string{}, en...), ru[:15]...)); got != "en" {
		t.Fatalf("a few Russian words after English: %q", got)
	}
	mixed := append(append([]string{}, en[len(en)-50:]...), ru[:50]...)
	if got := GuessLanguage(mixed); got != "" {
		t.Fatalf("half and half: %q", got)
	}
	if got := GuessLanguage(append(append([]string{}, en...), ru...)); got != "ru" {
		t.Fatalf("after the switch: %q", got)
	}
	uk := strings.Fields(guessUkrainian)
	if got := GuessLanguage(append(append([]string{}, ru...), uk...)); got != "uk" {
		t.Fatalf("russian then ukrainian: %q", got)
	}
}

func TestAutoUsesTheGuessedLanguageAtTheHardLimit(t *testing.T) {
	p := NewPhraser()
	p.SetLanguage("auto")
	t0 := time.Date(2026, 11, 12, 10, 0, 0, 0, time.UTC)
	got := p.Whole(guessEnglish, t0, t0, t0)
	if len(got) == 0 {
		t.Fatal("no phrases")
	}
	tail := strings.Fields("you will see that the build using Opus five point five on a Claude Max subscription are absolutely phenomenal and")
	words := append([]string{}, tail...)
	out := p.Update(0, strings.Join(words, " "), false, t0.Add(time.Minute))
	if len(out) != 1 || out[0].Text != "you will see that the build using Opus five point five" {
		t.Fatalf("phrases %v", textsOf(out))
	}
}
