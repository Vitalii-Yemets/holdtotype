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

const guessGerman = "Also das Erste, was uns aufgefallen ist, als wir uns die Zahlen angesehen haben, war, dass die meisten Nutzer, die sich im März angemeldet haben, " +
	"nach der ersten Woche nicht mehr gekommen sind. Wir haben dann geschaut, was beim ersten Besuch passiert ist, und es hat sich gezeigt, dass die Einrichtung " +
	"am Telefon viel zu lange gedauert hat, weil das Formular nach Dingen gefragt hat, die man nicht zur Hand hatte. Jetzt sind es nur noch zwei Felder, " +
	"und der Rest wartet in den Einstellungen, bis ihn wirklich jemand braucht. Danach ist der Anteil der Leute, die auch in der zweiten Woche noch da waren, " +
	"um etwa ein Drittel gestiegen, und das ist mehr, als wir von so einer kleinen Änderung erwartet hätten."

const guessPolish = "A więc pierwsza rzecz, którą zauważyliśmy, kiedy spojrzeliśmy na dane, była taka, że większość użytkowników, którzy zarejestrowali się w marcu, " +
	"nie wróciła po pierwszym tygodniu. Potem sprawdziliśmy, co się działo przy pierwszej wizycie, i okazało się, że konfiguracja na telefonie trwała " +
	"zdecydowanie za długo, bo formularz pytał o rzeczy, których ludzie nie mieli pod ręką. Teraz są tylko dwa pola, a reszta czeka w ustawieniach, " +
	"aż ktoś naprawdę będzie jej potrzebował. Po tej zmianie odsetek osób, które zostały na drugi tydzień, wzrósł mniej więcej o jedną trzecią, " +
	"i to jest więcej, niż się spodziewaliśmy po tak małej zmianie."

func settleAt(words []string, lang string, common bool) int {
	at := 0
	for n := 1; n <= len(words); n++ {
		if guessWith(words[:n], common) != lang {
			at = 0
			continue
		}
		if at == 0 {
			at = n
		}
	}
	return at
}

func TestTheGuessSettlesSooner(t *testing.T) {
	cases := []struct {
		lang, text string
		most       int
	}{
		{"en", guessEnglish, 20},
		{"ru", guessRussian, 25},
		{"uk", guessUkrainian, 25},
		{"de", guessGerman, 20},
		{"pl", guessPolish, 25},
	}
	for _, c := range cases {
		words := strings.Fields(c.text)
		before := settleAt(words, c.lang, false)
		after := settleAt(words, c.lang, true)
		t.Logf("%s: settles after %d words with the frequent words, %d with the linking words alone, of %d", c.lang, after, before, len(words))
		if after == 0 || after > c.most {
			t.Errorf("%s settles after %d words", c.lang, after)
		}
		if before != 0 && after > before {
			t.Errorf("%s: the frequent words made it slower: %d against %d", c.lang, after, before)
		}
		if got := GuessLanguage(words); got != c.lang {
			t.Errorf("%s: the whole text gives %q", c.lang, got)
		}
	}
}

func TestTheFrequentWordListsBarelyOverlap(t *testing.T) {
	seen := map[string][]string{}
	for _, lang := range guessOrder {
		n := 0
		for w := range commonSets[lang] {
			seen[w] = append(seen[w], lang)
			n++
		}
		if n < 25 || n > 40 {
			t.Errorf("%s: %d frequent words", lang, n)
		}
	}
	for w, langs := range seen {
		if len(langs) > 1 {
			t.Errorf("%q is frequent in %v", w, langs)
		}
	}
}
