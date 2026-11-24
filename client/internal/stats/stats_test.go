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
	again := Open(path, t.Logf).Card(at("2026-11-22", 23), Choice{})
	if again.TodayWords != 6 || again.TodayDictations != 2 {
		t.Fatalf("the reopened store says %+v", again)
	}
}

func TestNormalizeChoice(t *testing.T) {
	cases := []struct {
		in   Choice
		want Choice
	}{
		{Choice{}, Choice{Period: "m"}},
		{Choice{Period: "fortnight"}, Choice{Period: "m"}},
		{Choice{Period: " W "}, Choice{Period: "w"}},
		{Choice{Period: "all"}, Choice{Period: "all"}},
		{Choice{Period: "custom"}, Choice{Period: "m"}},
		{Choice{Period: "custom", From: "2026-11-01", To: "tomorrow"}, Choice{Period: "m"}},
		{Choice{Period: "custom", From: "2026-02-30", To: "2026-03-01"}, Choice{Period: "m"}},
		{Choice{Period: "custom", From: "2026-11-10", To: "2026-11-01"}, Choice{Period: "custom", From: "2026-11-01", To: "2026-11-10"}},
		{Choice{Period: "q", From: "2026-11-01", To: "2026-11-10"}, Choice{Period: "q", From: "2026-11-01", To: "2026-11-10"}},
		{Choice{Period: "y", From: "x", To: "2026-11-10"}, Choice{Period: "y"}},
	}
	for _, c := range cases {
		if got := NormalizeChoice(c.in.Period, c.in.From, c.in.To); got != c.want {
			t.Errorf("NormalizeChoice(%+v) = %+v, want %+v", c.in, got, c.want)
		}
	}
}

func TestUnitFor(t *testing.T) {
	cases := []struct {
		period string
		days   int
		want   string
	}{
		{"w", 7, UnitDay}, {"m", 30, UnitDay}, {"custom", 31, UnitDay}, {"custom", 32, UnitWeek},
		{"q", 90, UnitWeek}, {"h", 180, UnitWeek}, {"custom", 186, UnitWeek}, {"custom", 187, UnitMonth},
		{"y", 365, UnitMonth}, {"all", 1, UnitMonth},
	}
	for _, c := range cases {
		if got := UnitFor(c.period, c.days); got != c.want {
			t.Errorf("UnitFor(%s, %d) = %s, want %s", c.period, c.days, got, c.want)
		}
	}
}

func addWords(t *testing.T, s *Store, day string, words int) {
	t.Helper()
	text := strings.TrimSpace(strings.Repeat("слово ", words))
	if err := s.AddDictation(at(day, 12), dict(text, "ru", "Code.exe", "gigaam-v3", int64(words)*600, 400)); err != nil {
		t.Fatal(err)
	}
}

func barSum(p Period) (int, int) {
	w, d := 0, 0
	for _, b := range p.Bars {
		w += b.Words
		d += b.Dictations
	}
	return w, d
}

func spans(p Period) string {
	out := []string{}
	for _, b := range p.Bars {
		out = append(out, b.From+".."+b.To)
	}
	return strings.Join(out, " ")
}

func TestPeriods(t *testing.T) {
	s, _ := tempStore(t)
	addWords(t, s, "2025-09-30", 9)
	addWords(t, s, "2026-06-09", 5)
	addWords(t, s, "2026-09-30", 7)
	addWords(t, s, "2026-10-23", 100)
	addWords(t, s, "2026-10-24", 11)
	addWords(t, s, "2026-11-15", 50)
	addWords(t, s, "2026-11-16", 2)
	addWords(t, s, "2026-11-22", 3)
	addWords(t, s, "2026-11-22", 4)
	addWords(t, s, "2026-11-23", 1000)
	now := at("2026-11-22", 20)
	want := []struct {
		period, from, unit string
		bars, words        int
	}{
		{"w", "2026-11-16", UnitDay, 7, 9},
		{"m", "2026-10-24", UnitDay, 30, 9 + 50 + 11},
		{"q", "2026-08-25", UnitWeek, 13, 9 + 50 + 11 + 100 + 7},
		{"h", "2026-05-27", UnitWeek, 26, 9 + 50 + 11 + 100 + 7 + 5},
		{"y", "2025-11-23", UnitMonth, 13, 9 + 50 + 11 + 100 + 7 + 5},
		{"all", "2025-09-30", UnitMonth, 15, 9 + 50 + 11 + 100 + 7 + 5 + 9},
	}
	prev := 0
	for _, w := range want {
		r := s.Report(now, Choice{Period: w.period}, nil)
		p := r.Period
		if r.Empty || p.Key != w.period || p.From != w.from || p.To != "2026-11-22" || p.Unit != w.unit || len(p.Bars) != w.bars || p.Words != w.words {
			t.Fatalf("%s: key %s, %s..%s, unit %s, %d bars, %d words", w.period, p.Key, p.From, p.To, p.Unit, len(p.Bars), p.Words)
		}
		if words, dicts := barSum(p); words != p.Words || dicts != p.Dictations {
			t.Fatalf("%s: the bars hold %d words and %d dictations, the period %d and %d", w.period, words, dicts, p.Words, p.Dictations)
		}
		if p.TodayWords != 7 || p.SpeechMs != int64(p.Words)*600 {
			t.Fatalf("%s: today %d, speech %d", w.period, p.TodayWords, p.SpeechMs)
		}
		if p.Words < prev {
			t.Fatalf("%s holds fewer words than the shorter period before it", w.period)
		}
		prev = p.Words
		if r.First != "2025-09-30" {
			t.Fatalf("the history starts %s", r.First)
		}
	}
	week := s.Report(now, Choice{Period: "w"}, nil).Period
	if week.Bars[0] != (Bar{From: "2026-11-16", To: "2026-11-16", Words: 2, Dictations: 1}) || week.Bars[6].Words != 7 || week.Bars[3].Words != 0 {
		t.Fatalf("day bars %+v", week.Bars)
	}
	all := s.Report(now, Choice{Period: "all"}, nil).Period
	if all.Bars[0] != (Bar{From: "2025-09-30", To: "2025-09-30", Words: 9, Dictations: 1}) || all.Bars[1].From != "2025-10-01" || all.Bars[1].To != "2025-10-31" || all.Bars[14].From != "2026-11-01" || all.Bars[14].To != "2026-11-22" {
		t.Fatalf("month bars %s", spans(all))
	}
}

func TestWeekBuckets(t *testing.T) {
	s, _ := tempStore(t)
	addWords(t, s, "2026-08-25", 3)
	addWords(t, s, "2026-08-30", 4)
	addWords(t, s, "2026-08-31", 5)
	addWords(t, s, "2026-11-25", 6)
	p := s.Report(at("2026-11-25", 9), Choice{Period: "q"}, nil).Period
	if p.From != "2026-08-28" || p.Unit != UnitWeek {
		t.Fatalf("a quarter ending on a Wednesday starts %s, unit %s", p.From, p.Unit)
	}
	first, last := p.Bars[0], p.Bars[len(p.Bars)-1]
	if first.From != "2026-08-28" || first.To != "2026-08-30" || first.Words != 4 {
		t.Fatalf("the partial first week is %+v", first)
	}
	if p.Bars[1].From != "2026-08-31" || p.Bars[1].To != "2026-09-06" || p.Bars[1].Words != 5 {
		t.Fatalf("the first full week is %+v", p.Bars[1])
	}
	if last.From != "2026-11-23" || last.To != "2026-11-25" || last.Words != 6 {
		t.Fatalf("the partial last week is %+v", last)
	}
	for i := 1; i < len(p.Bars)-1; i++ {
		a, _ := time.Parse(dayLayout, p.Bars[i].From)
		b, _ := time.Parse(dayLayout, p.Bars[i].To)
		if a.Weekday() != time.Monday || b.Weekday() != time.Sunday || b.Sub(a) != 6*24*time.Hour {
			t.Fatalf("week %d is %s..%s", i, p.Bars[i].From, p.Bars[i].To)
		}
	}
	if len(p.Bars) != 14 {
		t.Fatalf("%d weeks: %s", len(p.Bars), spans(p))
	}
}

func TestCustomRanges(t *testing.T) {
	s, _ := tempStore(t)
	addWords(t, s, "2025-06-15", 1)
	addWords(t, s, "2025-12-31", 2)
	addWords(t, s, "2026-01-01", 3)
	addWords(t, s, "2026-11-05", 4)
	addWords(t, s, "2026-11-20", 5)
	now := at("2026-11-22", 9)
	in := s.Report(now, Choice{Period: "custom", From: "2026-11-03", To: "2026-11-10"}, nil).Period
	if in.Unit != UnitDay || len(in.Bars) != 8 || in.Words != 4 || in.Bars[2].Words != 4 || in.From != "2026-11-03" || in.To != "2026-11-10" {
		t.Fatalf("a range inside one month: %s, %d bars, %d words", in.Unit, len(in.Bars), in.Words)
	}
	across := s.Report(now, Choice{Period: "custom", From: "2025-12-20", To: "2026-01-10"}, nil).Period
	if across.Unit != UnitDay || len(across.Bars) != 22 || across.Words != 5 {
		t.Fatalf("a range across the new year: %s, %d bars, %d words", across.Unit, len(across.Bars), across.Words)
	}
	weeks := s.Report(now, Choice{Period: "custom", From: "2025-11-20", To: "2026-02-10"}, nil).Period
	if weeks.Unit != UnitWeek || weeks.Bars[0].From != "2025-11-20" || weeks.Bars[0].To != "2025-11-23" || weeks.Words != 5 {
		t.Fatalf("weeks across the year: %s, %s, %d words", weeks.Unit, spans(weeks), weeks.Words)
	}
	months := s.Report(now, Choice{Period: "custom", From: "2025-06-10", To: "2026-02-10"}, nil).Period
	if months.Unit != UnitMonth || len(months.Bars) != 9 || months.Bars[0].From != "2025-06-10" || months.Bars[0].To != "2025-06-30" || months.Bars[8].From != "2026-02-01" || months.Bars[8].To != "2026-02-10" || months.Words != 6 {
		t.Fatalf("months across the year: %s", spans(months))
	}
	future := s.Report(now, Choice{Period: "custom", From: "2026-11-20", To: "2026-12-31"}, nil).Period
	if future.To != "2026-11-22" || len(future.Bars) != 3 || future.Words != 5 {
		t.Fatalf("a range into the future: %s..%s, %d bars", future.From, future.To, len(future.Bars))
	}
	ahead := s.Report(now, Choice{Period: "custom", From: "2026-12-01", To: "2026-12-31"}, nil).Period
	if ahead.From != "2026-11-22" || ahead.To != "2026-11-22" || len(ahead.Bars) != 1 {
		t.Fatalf("a range wholly in the future: %s..%s", ahead.From, ahead.To)
	}
	broken := s.Report(now, Choice{Period: "custom", From: "soon"}, nil).Period
	if broken.Key != "m" || broken.From != "2026-10-24" {
		t.Fatalf("a broken range gives %s from %s", broken.Key, broken.From)
	}
}

func TestShortHistory(t *testing.T) {
	s, _ := tempStore(t)
	addWords(t, s, "2026-11-20", 3)
	addWords(t, s, "2026-11-21", 4)
	now := at("2026-11-22", 9)
	y := s.Report(now, Choice{Period: "y"}, nil).Period
	if len(y.Bars) != 13 || y.Words != 7 || y.Bars[12].Words != 7 || y.Bars[0].Words != 0 {
		t.Fatalf("a year over two days of history: %d bars, %d words", len(y.Bars), y.Words)
	}
	all := s.Report(now, Choice{Period: "all"}, nil).Period
	if all.From != "2026-11-20" || len(all.Bars) != 1 || all.Words != 7 {
		t.Fatalf("all time over two days: %s, %d bars", all.From, len(all.Bars))
	}
}

func TestCardFollowsThePeriod(t *testing.T) {
	s, _ := tempStore(t)
	addWords(t, s, "2026-08-01", 40)
	addWords(t, s, "2026-11-01", 30)
	addWords(t, s, "2026-11-22", 3)
	addWords(t, s, "2026-11-22", 4)
	now := at("2026-11-22", 20)
	for _, c := range []Choice{{Period: "w"}, {Period: "m"}, {Period: "q"}, {Period: "all"}, {Period: "custom", From: "2026-10-30", To: "2026-11-02"}} {
		p := s.Report(now, c, nil).Period
		card := s.Card(now, c)
		if !card.Any || card.Today != "2026-11-22" || card.TodayWords != 7 || card.TodayDictations != 2 || card.Period != p.Key || card.Unit != p.Unit || card.From != p.From || card.To != p.To || card.Words != p.Words || card.Dictations != p.Dictations || card.SpeechMs != p.SpeechMs || len(card.Bars) != len(p.Bars) {
			t.Fatalf("%+v: the card %+v does not match the period", c, card)
		}
	}
	if c := s.Card(now, Choice{Period: "custom", From: "2026-10-30", To: "2026-11-02"}); c.Words != 30 || c.Period != "custom" {
		t.Fatalf("the card over a range: %+v", c)
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
	apps := s.Report(day, Choice{Period: "w"}, nil).Period.Apps
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
	p := s.Report(day, Choice{Period: "m"}, func(k string) string { return "name of " + k }).Period
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
	now := at("2026-11-22", 9)
	r := s.Report(now, Choice{Period: "all"}, nil)
	if r.Empty || s.Card(now, Choice{}).Any {
		t.Fatalf("empty %v, the card says dictated", r.Empty)
	}
	if r.Period.From != "2026-11-20" || r.First != "2026-11-20" {
		t.Fatalf("all time starts %s, the history %s", r.Period.From, r.First)
	}
}

func TestEmptyStore(t *testing.T) {
	s, path := tempStore(t)
	now := at("2026-11-22", 9)
	r := s.Report(now, Choice{Period: "w"}, nil)
	if !r.Empty || s.Card(now, Choice{}).Any || r.First != "2026-11-22" {
		t.Fatal("a new store is not empty")
	}
	if len(r.Period.Bars) != 7 || r.Period.Words != 0 {
		t.Fatalf("an empty week %+v", r.Period)
	}
	all := s.Report(now, Choice{Period: "all"}, nil).Period
	if len(all.Bars) != 1 || all.From != "2026-11-22" || all.Unit != UnitMonth {
		t.Fatalf("empty all time %+v", all)
	}
	if card := s.Card(now, Choice{}); len(card.Bars) != 30 || card.Words != 0 || card.Period != "m" {
		t.Fatalf("an empty card %+v", card)
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
	if r := s.Report(at("2026-11-22", 9), Choice{}, nil); !r.Empty {
		t.Fatal("the store still has numbers after a reset")
	}
	if r := Open(path, t.Logf).Report(at("2026-11-22", 9), Choice{}, nil); !r.Empty {
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
			if r := s.Report(at("2026-11-22", 9), Choice{}, nil); !r.Empty {
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
	r := Open(path, t.Logf).Report(at("2026-11-22", 9), Choice{Period: "w"}, nil)
	if p := r.Period; p.Dictations != 1 || p.Words != 2 {
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
			_ = s.Report(day, Choice{}, nil)
			_ = s.Card(day, Choice{})
		}()
	}
	wg.Wait()
	p := Open(path, t.Logf).Report(day, Choice{Period: "w"}, nil).Period
	if p.Dictations != 40 || p.Words != 120 || p.SpeechMs != 4000 || p.SubsMs != 40000 || p.TrMs != 20000 {
		t.Fatalf("after concurrent adds: %d dictations, %d words, %d ms speech, %d ms subtitles, %d ms translated", p.Dictations, p.Words, p.SpeechMs, p.SubsMs, p.TrMs)
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Fatal("a temporary file is left behind")
	}
}

func TestCheckReport(t *testing.T) {
	s, path := tempStore(t)
	addWords(t, s, "2026-11-20", 3)
	addWords(t, s, "2026-11-22", 4)
	now := at("2026-11-22", 9)
	if _, err := CheckReport(path, nil, now, nil); err == nil {
		t.Fatal("no period gave a report")
	}
	cases := []struct {
		args       []string
		key, from  string
		bars, word int
	}{
		{[]string{"w"}, "w", "2026-11-16", 7, 7},
		{[]string{"all"}, "all", "2026-11-20", 1, 7},
		{[]string{"custom", "2026-11-21", "2026-11-20"}, "custom", "2026-11-20", 2, 3},
		{[]string{"custom"}, "m", "2026-10-24", 30, 7},
		{[]string{"decade"}, "m", "2026-10-24", 30, 7},
	}
	for _, c := range cases {
		out, err := CheckReport(path, c.args, now, nil)
		if err != nil {
			t.Fatalf("%v: %v", c.args, err)
		}
		var r Report
		if err := json.Unmarshal(out, &r); err != nil {
			t.Fatalf("%v: %v in %s", c.args, err, out)
		}
		if r.Period.Key != c.key || r.Period.From != c.from || len(r.Period.Bars) != c.bars || r.Period.Words != c.word || r.Empty {
			t.Fatalf("%v: %s", c.args, out)
		}
	}
	bad := filepath.Join(t.TempDir(), "holdtotype-stats.json")
	_ = os.WriteFile(bad, []byte("{broken"), 0o600)
	if _, err := CheckReport(bad, []string{"w"}, now, nil); err == nil || !strings.Contains(err.Error(), "damaged") {
		t.Fatalf("a damaged file gave %v", err)
	}
	if out, err := CheckReport(filepath.Join(t.TempDir(), "none.json"), []string{"w"}, now, nil); err != nil || !strings.Contains(string(out), `"empty":true`) {
		t.Fatalf("a missing file gave %s, %v", out, err)
	}
}
