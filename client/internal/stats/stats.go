package stats

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"holdtotype/internal/apprules"
)

const (
	OtherApp = "*"

	DefaultWPM = 40
	MinWPM     = 10
	MaxWPM     = 150
	StepWPM    = 5

	PeriodWeek    = "w"
	PeriodMonth   = "m"
	PeriodQuarter = "q"
	PeriodHalf    = "h"
	PeriodYear    = "y"
	PeriodAll     = "all"
	PeriodCustom  = "custom"
	DefaultPeriod = PeriodMonth

	UnitDay   = "day"
	UnitWeek  = "week"
	UnitMonth = "month"

	dayLayout   = "2006-01-02"
	topApps     = 5
	maxDayBars  = 31
	maxWeekBars = 186
)

var periodDays = map[string]int{PeriodWeek: 7, PeriodMonth: 30, PeriodQuarter: 90, PeriodHalf: 180, PeriodYear: 365}

func Program(app, skip string) string {
	app = strings.TrimSpace(app)
	if app == "" {
		return OtherApp
	}
	if skip = strings.TrimSpace(skip); skip != "" {
		if _, hit := apprules.Find([]apprules.Rule{{Match: skip}}, app); hit {
			return OtherApp
		}
	}
	return app
}

func TypingWPM(v int) int {
	if v <= 0 {
		return DefaultWPM
	}
	v = (v + StepWPM/2) / StepWPM * StepWPM
	return min(MaxWPM, max(MinWPM, v))
}

type Model struct {
	Dictations int   `json:"dictations"`
	WaitMs     int64 `json:"wait_ms"`
}

type Day struct {
	Dictations int              `json:"dictations,omitempty"`
	Words      int              `json:"words,omitempty"`
	SpeechMs   int64            `json:"speech_ms,omitempty"`
	Langs      map[string]int   `json:"langs,omitempty"`
	Apps       map[string]int   `json:"apps,omitempty"`
	Models     map[string]Model `json:"models,omitempty"`
	SubsMs     int64            `json:"subs_ms,omitempty"`
	Pairs      map[string]int64 `json:"pairs,omitempty"`
}

type stored struct {
	Version int             `json:"version"`
	Days    map[string]*Day `json:"days"`
}

type Dictation struct {
	Text     string
	SpeechMs int64
	WaitMs   int64
	Model    string
	Lang     string
	App      string
}

type Store struct {
	mu     sync.Mutex
	path   string
	logf   func(string, ...any)
	loaded bool
	days   map[string]*Day
}

func Open(path string, logf func(string, ...any)) *Store {
	if logf == nil {
		logf = func(string, ...any) {}
	}
	return &Store{path: path, logf: logf}
}

func (s *Store) loadLocked() {
	if s.loaded {
		return
	}
	s.loaded = true
	s.days = map[string]*Day{}
	if s.path == "" {
		return
	}
	data, err := os.ReadFile(s.path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.logf("statistics: %s cannot be read, starting empty: %v", s.path, err)
		}
		return
	}
	var f stored
	if err := json.Unmarshal(data, &f); err != nil {
		s.logf("statistics: %s is damaged, starting empty: %v", s.path, err)
		return
	}
	for key, d := range f.Days {
		if d == nil {
			continue
		}
		if _, err := time.Parse(dayLayout, key); err != nil {
			s.logf("statistics: skipping the day %q, it is not a date", key)
			continue
		}
		s.days[key] = d
	}
}

func (s *Store) dayLocked(at time.Time) *Day {
	key := at.Format(dayLayout)
	d := s.days[key]
	if d == nil {
		d = &Day{}
		s.days[key] = d
	}
	return d
}

func (s *Store) AddDictation(at time.Time, in Dictation) error {
	words := Words(in.Text)
	if words == 0 {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.loadLocked()
	d := s.dayLocked(at)
	d.Dictations++
	d.Words += words
	d.SpeechMs += max(0, in.SpeechMs)
	if in.Lang != "" {
		if d.Langs == nil {
			d.Langs = map[string]int{}
		}
		d.Langs[in.Lang] += words
	}
	app := in.App
	if app == "" {
		app = OtherApp
	}
	if d.Apps == nil {
		d.Apps = map[string]int{}
	}
	d.Apps[app] += words
	if in.Model != "" {
		if d.Models == nil {
			d.Models = map[string]Model{}
		}
		m := d.Models[in.Model]
		m.Dictations++
		m.WaitMs += max(0, in.WaitMs)
		d.Models[in.Model] = m
	}
	return s.saveLocked()
}

func (s *Store) AddSubtitles(at time.Time, listenMs int64, pairs map[string]int64) error {
	if listenMs <= 0 && len(pairs) == 0 {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.loadLocked()
	d := s.dayLocked(at)
	d.SubsMs += max(0, listenMs)
	for pair, ms := range pairs {
		if ms <= 0 || pair == "" {
			continue
		}
		if d.Pairs == nil {
			d.Pairs = map[string]int64{}
		}
		d.Pairs[pair] += ms
	}
	return s.saveLocked()
}

func (s *Store) Reset() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.loaded = true
	s.days = map[string]*Day{}
	return s.saveLocked()
}

func (s *Store) saveLocked() error {
	if s.path == "" {
		return nil
	}
	if len(s.days) == 0 {
		if err := os.Remove(s.path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
		return nil
	}
	out, err := json.MarshalIndent(stored{Version: 1, Days: s.days}, "", " ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, out, 0o600); err != nil {
		return err
	}
	if err := os.Rename(tmp, filepath.Clean(s.path)); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}

type Choice struct {
	Period string
	From   string
	To     string
}

func NormalizeChoice(period, from, to string) Choice {
	period = strings.ToLower(strings.TrimSpace(period))
	if _, ok := periodDays[period]; !ok && period != PeriodAll && period != PeriodCustom {
		period = DefaultPeriod
	}
	a, errA := time.Parse(dayLayout, strings.TrimSpace(from))
	b, errB := time.Parse(dayLayout, strings.TrimSpace(to))
	if errA != nil || errB != nil {
		if period == PeriodCustom {
			period = DefaultPeriod
		}
		return Choice{Period: period}
	}
	if a.After(b) {
		a, b = b, a
	}
	return Choice{Period: period, From: a.Format(dayLayout), To: b.Format(dayLayout)}
}

func UnitFor(period string, days int) string {
	switch {
	case period == PeriodAll:
		return UnitMonth
	case days <= maxDayBars:
		return UnitDay
	case days <= maxWeekBars:
		return UnitWeek
	}
	return UnitMonth
}

type Bar struct {
	From       string `json:"from"`
	To         string `json:"to"`
	Words      int    `json:"words"`
	Dictations int    `json:"dictations"`
}

type Lang struct {
	Code  string `json:"code"`
	Words int    `json:"words"`
}

type App struct {
	Name  string `json:"name"`
	Words int    `json:"words"`
	Other bool   `json:"other,omitempty"`
}

type ModelRow struct {
	Key        string `json:"key"`
	Name       string `json:"name"`
	Dictations int    `json:"dictations"`
	WaitMs     int64  `json:"wait_ms"`
}

type Pair struct {
	From string `json:"from"`
	To   string `json:"to"`
	Ms   int64  `json:"ms"`
}

type Period struct {
	Key        string     `json:"key"`
	Unit       string     `json:"unit"`
	From       string     `json:"from"`
	To         string     `json:"to"`
	Dictations int        `json:"dictations"`
	Words      int        `json:"words"`
	SpeechMs   int64      `json:"speech_ms"`
	TodayWords int        `json:"today_words"`
	Bars       []Bar      `json:"bars"`
	Langs      []Lang     `json:"langs"`
	Apps       []App      `json:"apps"`
	Models     []ModelRow `json:"models"`
	SubsMs     int64      `json:"subs_ms"`
	TrMs       int64      `json:"tr_ms"`
	Pairs      []Pair     `json:"pairs"`
}

type Card struct {
	Any             bool   `json:"any"`
	TodayWords      int    `json:"today_words"`
	TodayDictations int    `json:"today_dictations"`
	Today           string `json:"today"`
	Period          string `json:"period"`
	Unit            string `json:"unit"`
	From            string `json:"from"`
	To              string `json:"to"`
	Words           int    `json:"words"`
	Dictations      int    `json:"dictations"`
	SpeechMs        int64  `json:"speech_ms"`
	Bars            []Bar  `json:"bars"`
}

type Report struct {
	Empty  bool   `json:"empty"`
	Today  string `json:"today"`
	First  string `json:"first"`
	Period Period `json:"period"`
}

func dayStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func (s *Store) firstLocked(today time.Time) (time.Time, bool, bool) {
	first := today
	anyDict, anySubs := false, false
	for key, d := range s.days {
		t, err := time.ParseInLocation(dayLayout, key, today.Location())
		if err != nil || t.After(today) {
			continue
		}
		if d.Dictations > 0 {
			anyDict = true
		}
		if d.SubsMs > 0 {
			anySubs = true
		}
		if (d.Dictations > 0 || d.SubsMs > 0) && t.Before(first) {
			first = t
		}
	}
	return first, anyDict, anySubs
}

func rangeOf(c Choice, today, first time.Time) (string, time.Time, time.Time) {
	c = NormalizeChoice(c.Period, c.From, c.To)
	switch c.Period {
	case PeriodAll:
		return c.Period, first, today
	case PeriodCustom:
		from, _ := time.ParseInLocation(dayLayout, c.From, today.Location())
		to, _ := time.ParseInLocation(dayLayout, c.To, today.Location())
		if to.After(today) {
			to = today
		}
		if from.After(to) {
			from = to
		}
		return c.Period, from, to
	}
	return c.Period, today.AddDate(0, 0, 1-periodDays[c.Period]), today
}

func (s *Store) Report(now time.Time, c Choice, modelName func(string) string) Report {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.loadLocked()
	today := dayStart(now)
	first, anyDict, anySubs := s.firstLocked(today)
	key, from, to := rangeOf(c, today, first)
	return Report{
		Empty:  !anyDict && !anySubs,
		Today:  today.Format(dayLayout),
		First:  first.Format(dayLayout),
		Period: s.periodLocked(key, from, to, today, modelName),
	}
}

func (s *Store) Card(now time.Time, c Choice) Card {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.loadLocked()
	today := dayStart(now)
	first, anyDict, _ := s.firstLocked(today)
	key, from, to := rangeOf(c, today, first)
	p := s.periodLocked(key, from, to, today, nil)
	card := Card{Any: anyDict, Today: today.Format(dayLayout), Period: p.Key, Unit: p.Unit, From: p.From, To: p.To, Words: p.Words, Dictations: p.Dictations, SpeechMs: p.SpeechMs, Bars: p.Bars}
	if t := s.days[today.Format(dayLayout)]; t != nil {
		card.TodayWords = t.Words
		card.TodayDictations = t.Dictations
	}
	return card
}

func bucketEnd(t, to time.Time, unit string) time.Time {
	end := t
	switch unit {
	case UnitWeek:
		end = t.AddDate(0, 0, (7-int(t.Weekday()))%7)
	case UnitMonth:
		end = time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location())
	}
	if end.After(to) {
		end = to
	}
	return end
}

func (s *Store) periodLocked(key string, from, to, today time.Time, modelName func(string) string) Period {
	days := int(to.Sub(from).Hours()/24+0.5) + 1
	p := Period{Key: key, Unit: UnitFor(key, days), From: from.Format(dayLayout), To: to.Format(dayLayout)}
	if d := s.days[today.Format(dayLayout)]; d != nil {
		p.TodayWords = d.Words
	}
	langs := map[string]int{}
	apps := map[string]int{}
	models := map[string]Model{}
	pairs := map[string]int64{}
	end := from.AddDate(0, 0, -1)
	for t := from; !t.After(to); t = t.AddDate(0, 0, 1) {
		if t.After(end) {
			end = bucketEnd(t, to, p.Unit)
			p.Bars = append(p.Bars, Bar{From: t.Format(dayLayout), To: end.Format(dayLayout)})
		}
		d := s.days[t.Format(dayLayout)]
		if d == nil {
			continue
		}
		bar := &p.Bars[len(p.Bars)-1]
		bar.Words += d.Words
		bar.Dictations += d.Dictations
		p.Dictations += d.Dictations
		p.Words += d.Words
		p.SpeechMs += d.SpeechMs
		for k, v := range d.Langs {
			langs[k] += v
		}
		for k, v := range d.Apps {
			apps[k] += v
		}
		for k, v := range d.Models {
			m := models[k]
			m.Dictations += v.Dictations
			m.WaitMs += v.WaitMs
			models[k] = m
		}
		p.SubsMs += d.SubsMs
		for k, v := range d.Pairs {
			pairs[k] += v
			p.TrMs += v
		}
	}
	p.Langs = []Lang{}
	for code, n := range langs {
		if n > 0 {
			p.Langs = append(p.Langs, Lang{Code: code, Words: n})
		}
	}
	sort.Slice(p.Langs, func(i, j int) bool {
		if p.Langs[i].Words != p.Langs[j].Words {
			return p.Langs[i].Words > p.Langs[j].Words
		}
		return p.Langs[i].Code < p.Langs[j].Code
	})
	p.Apps = topAppsOf(apps)
	p.Models = []ModelRow{}
	for key, m := range models {
		if m.Dictations <= 0 {
			continue
		}
		name := key
		if modelName != nil {
			name = modelName(key)
		}
		p.Models = append(p.Models, ModelRow{Key: key, Name: name, Dictations: m.Dictations, WaitMs: m.WaitMs})
	}
	sort.Slice(p.Models, func(i, j int) bool {
		a, b := p.Models[i], p.Models[j]
		if l, r := a.WaitMs*int64(b.Dictations), b.WaitMs*int64(a.Dictations); l != r {
			return l < r
		}
		return a.Name < b.Name
	})
	p.Pairs = []Pair{}
	for key, ms := range pairs {
		from, to := splitPair(key)
		p.Pairs = append(p.Pairs, Pair{From: from, To: to, Ms: ms})
	}
	sort.Slice(p.Pairs, func(i, j int) bool {
		if p.Pairs[i].Ms != p.Pairs[j].Ms {
			return p.Pairs[i].Ms > p.Pairs[j].Ms
		}
		return p.Pairs[i].From+p.Pairs[i].To < p.Pairs[j].From+p.Pairs[j].To
	})
	return p
}

func PairKey(from, to string) string { return from + ">" + to }

func splitPair(key string) (string, string) {
	from, to, _ := strings.Cut(key, ">")
	return from, to
}

func topAppsOf(apps map[string]int) []App {
	named := []App{}
	other := 0
	for name, n := range apps {
		if n <= 0 {
			continue
		}
		if name == OtherApp {
			other += n
			continue
		}
		named = append(named, App{Name: name, Words: n})
	}
	sort.Slice(named, func(i, j int) bool {
		if named[i].Words != named[j].Words {
			return named[i].Words > named[j].Words
		}
		return named[i].Name < named[j].Name
	})
	if len(named) > topApps {
		for _, a := range named[topApps:] {
			other += a.Words
		}
		named = named[:topApps]
	}
	if other > 0 {
		named = append(named, App{Words: other, Other: true})
	}
	return named
}

func CheckReport(path string, args []string, now time.Time, modelName func(string) string) ([]byte, error) {
	if len(args) == 0 {
		return nil, errors.New("needs a period: w, m, q, h, y, all, or custom with two dates YYYY-MM-DD")
	}
	c := Choice{Period: args[0]}
	if len(args) >= 3 {
		c.From, c.To = args[1], args[2]
	}
	var problems []string
	s := Open(path, func(f string, a ...any) { problems = append(problems, fmt.Sprintf(f, a...)) })
	rep := s.Report(now, c, modelName)
	if len(problems) > 0 {
		return nil, errors.New(strings.Join(problems, "; "))
	}
	return json.Marshal(rep)
}
