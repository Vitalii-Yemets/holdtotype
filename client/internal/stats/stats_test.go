package stats

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestWords(t *testing.T) {
	cases := []struct {
		text string
		want int
	}{
		{"", 0},
		{"   \n\t ", 0},
		{"Привет, мир!", 2},
		{"Hello world", 2},
		{"it's a well-known fact", 4},
		{"— … !!! ,", 0},
		{"в 2026 году 1,5 раза", 5},
		{"GitHub.com", 1},
		{"Що ж, добре", 3},
		{"zażółć gęślą jaźń", 3},
		{"你好世界", 4},
		{"これはテストです", 8},
		{"Windows用户", 3},
		{"test。这", 2},
		{"สวัสดี", 4},
		{"«Да», — сказал он.", 3},
	}
	for _, c := range cases {
		if got := Words(c.text); got != c.want {
			t.Errorf("Words(%q) = %d, want %d", c.text, got, c.want)
		}
	}
}

func TestTypingWPM(t *testing.T) {
	cases := map[int]int{0: 40, -5: 40, 1: 10, 10: 10, 12: 10, 13: 15, 40: 40, 42: 40, 43: 45, 150: 150, 151: 150, 999: 150}
	for in, want := range cases {
		if got := TypingWPM(in); got != want {
			t.Errorf("TypingWPM(%d) = %d, want %d", in, got, want)
		}
	}
}

func TestProgram(t *testing.T) {
	cases := []struct {
		app, skip, want string
	}{
		{"Telegram.exe", "", "Telegram.exe"},
		{"", "", OtherApp},
		{"   ", "keepass.exe", OtherApp},
		{"KeePass.exe", "keepass.exe,1password.exe", OtherApp},
		{"Teams.exe", "teams*", OtherApp},
		{"chrome.exe", "keepass.exe", "chrome.exe"},
	}
	for _, c := range cases {
		if got := Program(c.app, c.skip); got != c.want {
			t.Errorf("Program(%q, %q) = %q, want %q", c.app, c.skip, got, c.want)
		}
	}
}

func at(day string, hour int) time.Time {
	t, err := time.ParseInLocation("2006-01-02", day, time.Local)
	if err != nil {
		panic(err)
	}
	return t.Add(time.Duration(hour) * time.Hour)
}

func tempStore(t *testing.T) (*Store, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "holdtotype-stats.json")
	return Open(path, t.Logf), path
}

func dict(text, lang, app, model string, speech, wait int64) Dictation {
	return Dictation{Text: text, Lang: lang, App: app, Model: model, SpeechMs: speech, WaitMs: wait}
}

func TestAddKeepsNumbersByDayWithoutText(t *testing.T) {
	s, path := tempStore(t)
	secret := "пароль от банка девятнадцать"
	if err := s.AddDictation(at("2026-11-22", 9), dict(secret, "ru", "Telegram.exe", "gigaam-v3", 4000, 300)); err != nil {
		t.Fatal(err)
	}
	if err := s.AddDictation(at("2026-11-22", 10), dict("hello there", "en", "chrome.exe", "nemotron-3.5", 1500, 500)); err != nil {
		t.Fatal(err)
	}
	if err := s.AddDictation(at("2026-11-22", 11), dict("  ", "en", "chrome.exe", "nemotron-3.5", 900, 100)); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range strings.Fields(secret + " hello there") {
		if strings.Contains(string(raw), w) {
			t.Fatalf("the file holds the dictated word %q:\n%s", w, raw)
		}
	}
	var f stored
	if err := json.Unmarshal(raw, &f); err != nil {
		t.Fatal(err)
	}
	d := f.Days["2026-11-22"]
	if d == nil || d.Dictations != 2 || d.Words != 6 || d.SpeechMs != 5500 {
		t.Fatalf("the day is %+v", d)
	}
	if d.Langs["ru"] != 4 || d.Langs["en"] != 2 {
		t.Fatalf("languages %v", d.Langs)
	}
	if d.Apps["Telegram.exe"] != 4 || d.Apps["chrome.exe"] != 2 {
		t.Fatalf("programs %v", d.Apps)
	}
	if d.Models["gigaam-v3"] != (Model{1, 300}) || d.Models["nemotron-3.5"] != (Model{1, 500}) {
		t.Fatalf("models %v", d.Models)
	}
	again := Open(path, t.Logf).Report(at("2026-11-22", 23), nil)
	if again.Card.TodayWords != 6 || again.Card.TodayDictations != 2 {
		t.Fatalf("the reopened store says %+v", again.Card)
	}
}

func TestPeriods(t *testing.T) {
	s, _ := tempStore(t)
	add := func(day string, words int) {
		text := strings.TrimSpace(strings.Repeat("слово ", words))
		if err := s.AddDictation(at(day, 12), dict(text, "ru", "Code.exe", "gigaam-v3", int64(words)*600, 400)); err != nil {
			t.Fatal(err)
		}
	}
	add("2026-06-09", 5)
	add("2026-09-30", 7)
	add("2026-10-24", 11)
	add("2026-10-23", 100)
	add("2026-11-16", 2)
	add("2026-11-15", 50)
	add("2026-11-22", 3)
	add("2026-11-22", 4)
	add("2026-11-23", 1000)
	r := s.Report(at("2026-11-22", 20), nil)
	if r.Empty || !r.Card.Any {
		t.Fatal("a store with dictations reads as empty")
	}
	p7, p30, all := r.Periods["7"], r.Periods["30"], r.Periods["all"]
	if p7.From != "2026-11-16" || p7.To != "2026-11-22" || len(p7.Bars) != 7 {
		t.Fatalf("7 days: %s…%s, %d bars", p7.From, p7.To, len(p7.Bars))
	}
	if p7.Words != 9 || p7.Dictations != 3 || p7.SpeechMs != 9*600 || p7.TodayWords != 7 {
		t.Fatalf("7 days sum to %d words, %d dictations, %d ms, today %d", p7.Words, p7.Dictations, p7.SpeechMs, p7.TodayWords)
	}
	if p7.Bars[0].Key != "2026-11-16" || p7.Bars[0].Words != 2 || p7.Bars[6].Words != 7 || p7.Bars[3].Words != 0 {
		t.Fatalf("7 day bars %+v", p7.Bars)
	}
	if p30.From != "2026-10-24" || len(p30.Bars) != 30 || p30.Words != 9+50+11 {
		t.Fatalf("30 days: from %s, %d bars, %d words", p30.From, len(p30.Bars), p30.Words)
	}
	if all.From != "2026-06-09" || all.Words != 5+7+11+100+2+50+3+4 {
		t.Fatalf("all time: from %s, %d words", all.From, all.Words)
	}
	keys := []string{}
	sum := 0
	for _, b := range all.Bars {
		keys = append(keys, b.Key)
		sum += b.Words
	}
	if strings.Join(keys, " ") != "2026-06 2026-07 2026-08 2026-09 2026-10 2026-11" {
		t.Fatalf("months %v", keys)
	}
	if sum != all.Words || all.Bars[1].Words != 0 || all.Bars[4].Words != 111 {
		t.Fatalf("month bars %+v", all.Bars)
	}
	if r.Card.WeekWords != p7.Words || r.Card.WeekDictations != p7.Dictations || r.Card.WeekSpeechMs != p7.SpeechMs || r.Card.TodayWords != 7 || r.Card.TodayDictations != 2 {
		t.Fatalf("the card %+v does not match the week", r.Card)
	}
	if !(p7.Words <= p30.Words && p30.Words <= all.Words) {
		t.Fatal("periods are not nested")
	}
}

func TestOtherPrograms(t *testing.T) {
	s, path := tempStore(t)
	day := at("2026-11-22", 9)
	words := map[string]int{"a.exe": 9, "b.exe": 8, "c.exe": 7, "d.exe": 6, "e.exe": 5, "f.exe": 4, "g.exe": 3}
	for app, n := range words {
		if err := s.AddDictation(day, dict(strings.Repeat("x ", n), "en", app, "", 0, 0)); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.AddDictation(day, dict("secret bank words", "en", Program("KeePass.exe", "keepass.exe"), "", 0, 0)); err != nil {
		t.Fatal(err)
	}
	if err := s.AddDictation(day, dict("nowhere", "en", "", "", 0, 0)); err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(path)
	if strings.Contains(strings.ToLower(string(raw)), "keepass") {
		t.Fatal("an excluded program is written down by name")
	}
	apps := s.Report(day, nil).Periods["7"].Apps
	if len(apps) != 6 {
		t.Fatalf("want five programs and the others, got %+v", apps)
	}
	for i, want := range []string{"a.exe", "b.exe", "c.exe", "d.exe", "e.exe"} {
		if apps[i].Name != want || apps[i].Other {
			t.Fatalf("row %d is %+v, want %s", i, apps[i], want)
		}
	}
	last := apps[5]
	if !last.Other || last.Name != "" || last.Words != 4+3+3+1 {
		t.Fatalf("the others row is %+v", last)
	}
}

func TestModelsAndSubtitles(t *testing.T) {
	s, _ := tempStore(t)
	day := at("2026-11-22", 9)
	for i := 0; i < 3; i++ {
		_ = s.AddDictation(day, dict("one two", "ru", "x.exe", "whisper", 0, 1800))
	}
	_ = s.AddDictation(day, dict("one two", "ru", "x.exe", "gigaam-v3", 0, 400))
	_ = s.AddDictation(day, dict("one two", "ru", "x.exe", "gigaam-v3", 0, 200))
	if err := s.AddSubtitles(day, 60000, map[string]int64{PairKey("en", "ru"): 45000}); err != nil {
		t.Fatal(err)
	}
	if err := s.AddSubtitles(day.Add(time.Hour), 30000, map[string]int64{PairKey("en", "ru"): 5000, PairKey("de", "ru"): 20000}); err != nil {
		t.Fatal(err)
	}
	p := s.Report(day, func(k string) string { return "name of " + k }).Periods["30"]
	if len(p.Models) != 2 || p.Models[0].Key != "gigaam-v3" || p.Models[0].Dictations != 2 || p.Models[0].WaitMs != 600 || p.Models[0].Name != "name of gigaam-v3" {
		t.Fatalf("models %+v", p.Models)
	}
	if p.SubsMs != 90000 || p.TrMs != 70000 {
		t.Fatalf("subtitles %d listened, %d translated", p.SubsMs, p.TrMs)
	}
	if len(p.Pairs) != 2 || p.Pairs[0] != (Pair{"en", "ru", 50000}) || p.Pairs[1] != (Pair{"de", "ru", 20000}) {
		t.Fatalf("pairs %+v", p.Pairs)
	}
}

func TestSubtitlesAloneAreNotDictations(t *testing.T) {
	s, _ := tempStore(t)
	_ = s.AddSubtitles(at("2026-11-20", 9), 1000, nil)
	r := s.Report(at("2026-11-22", 9), nil)
	if r.Empty || r.Card.Any {
		t.Fatalf("empty %v, card says dictated %v", r.Empty, r.Card.Any)
	}
	if r.Periods["all"].From != "2026-11-20" {
		t.Fatalf("all time starts %s", r.Periods["all"].From)
	}
}

func TestEmptyStore(t *testing.T) {
	s, path := tempStore(t)
	r := s.Report(at("2026-11-22", 9), nil)
	if !r.Empty || r.Card.Any {
		t.Fatal("a new store is not empty")
	}
	if len(r.Periods["7"].Bars) != 7 || len(r.Periods["all"].Bars) != 1 || r.Periods["all"].From != "2026-11-22" {
		t.Fatalf("empty periods %+v", r.Periods)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("reading wrote a file")
	}
}

func TestReset(t *testing.T) {
	s, path := tempStore(t)
	_ = s.AddDictation(at("2026-11-22", 9), dict("one", "ru", "x.exe", "m", 1, 1))
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
	if err := s.Reset(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("the file is still there after a reset")
	}
	if r := s.Report(at("2026-11-22", 9), nil); !r.Empty {
		t.Fatal("the store still has numbers after a reset")
	}
	if r := Open(path, t.Logf).Report(at("2026-11-22", 9), nil); !r.Empty {
		t.Fatal("a reopened store still has numbers after a reset")
	}
}

func TestDamagedFile(t *testing.T) {
	for name, body := range map[string]string{
		"garbage":   "{not json",
		"truncated": `{"version":1,"days":{"2026-11-22":{"dictations":3,`,
		"wrong":     `[1,2,3]`,
	} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "holdtotype-stats.json")
			if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
				t.Fatal(err)
			}
			var logged []string
			s := Open(path, func(f string, a ...any) { logged = append(logged, fmt.Sprintf(f, a...)) })
			if r := s.Report(at("2026-11-22", 9), nil); !r.Empty {
				t.Fatal("a damaged file gave numbers")
			}
			if len(logged) != 1 || !strings.Contains(logged[0], "damaged") {
				t.Fatalf("logged %q", logged)
			}
			if err := s.AddDictation(at("2026-11-22", 9), dict("one two", "ru", "x.exe", "m", 1, 1)); err != nil {
				t.Fatal(err)
			}
			raw, _ := os.ReadFile(path)
			var f stored
			if err := json.Unmarshal(raw, &f); err != nil || f.Days["2026-11-22"].Words != 2 {
				t.Fatalf("the replaced file is %s (%v)", raw, err)
			}
		})
	}
}

func TestBadDayKeysAreSkipped(t *testing.T) {
	path := filepath.Join(t.TempDir(), "holdtotype-stats.json")
	_ = os.WriteFile(path, []byte(`{"version":1,"days":{"yesterday":{"dictations":5,"words":9},"2026-11-21":{"dictations":1,"words":2},"2026-11-20":null}}`), 0o600)
	r := Open(path, t.Logf).Report(at("2026-11-22", 9), nil)
	if p := r.Periods["7"]; p.Dictations != 1 || p.Words != 2 {
		t.Fatalf("got %d dictations, %d words", p.Dictations, p.Words)
	}
}

func TestConcurrentAdds(t *testing.T) {
	s, path := tempStore(t)
	day := at("2026-11-22", 9)
	var wg sync.WaitGroup
	for i := 0; i < 40; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			if err := s.AddDictation(day, dict("one two three", "ru", "x.exe", "m", 100, 10)); err != nil {
				t.Error(err)
			}
		}()
		go func() {
			defer wg.Done()
			if err := s.AddSubtitles(day, 1000, map[string]int64{PairKey("en", "ru"): 500}); err != nil {
				t.Error(err)
			}
			_ = s.Report(day, nil)
		}()
	}
	wg.Wait()
	p := Open(path, t.Logf).Report(day, nil).Periods["7"]
	if p.Dictations != 40 || p.Words != 120 || p.SpeechMs != 4000 || p.SubsMs != 40000 || p.TrMs != 20000 {
		t.Fatalf("after concurrent adds: %d dictations, %d words, %d ms speech, %d ms subtitles, %d ms translated", p.Dictations, p.Words, p.SpeechMs, p.SubsMs, p.TrMs)
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Fatal("a temporary file is left behind")
	}
}
