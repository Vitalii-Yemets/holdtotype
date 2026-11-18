package subtitles

import (
	"holdtotype/internal/sentence"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const CheckRate = 2.5

type CallLine struct {
	At    time.Duration
	Words []string
}

var (
	callHead = regexp.MustCompile(`^===.*?(\d\d):(\d\d):(\d\d)(.*)$`)
	callLang = regexp.MustCompile(`·\s*[^·:]*:\s*([a-z]{2,4})\s*·`)
	callLine = regexp.MustCompile(`^\[(\d\d):(\d\d):(\d\d)\]\s*(.*)$`)
)

func clock(h, m, s string) time.Duration {
	hh, _ := strconv.Atoi(h)
	mm, _ := strconv.Atoi(m)
	ss, _ := strconv.Atoi(s)
	return time.Duration(hh)*time.Hour + time.Duration(mm)*time.Minute + time.Duration(ss)*time.Second
}

func ParseCall(raw, run string) ([]CallLine, string) {
	var out []CallLine
	lang := ""
	on := run == "" || run == "all"
	for _, line := range strings.Split(strings.ReplaceAll(raw, "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(line)
		if m := callHead.FindStringSubmatch(line); m != nil {
			start := m[1] + ":" + m[2] + ":" + m[3]
			if run != "" && run != "all" {
				on = start == run
			}
			if on {
				if l := callLang.FindStringSubmatch(line); l != nil && lang == "" {
					lang = l[1]
				}
			}
			continue
		}
		m := callLine.FindStringSubmatch(line)
		if m == nil || !on {
			continue
		}
		words := strings.Fields(m[4])
		if len(words) == 0 {
			continue
		}
		out = append(out, CallLine{At: clock(m[1], m[2], m[3]), Words: words})
	}
	return out, lang
}

type CheckPiece struct {
	Seq      int
	Text     string
	Words    int
	Para     bool
	Release  time.Time
	LastWord time.Time
	Parts    []string
	Times    []time.Time
}

func wordTime(base time.Time, l CallLine, k int) time.Time {
	back := time.Duration(float64(len(l.Words)-1-k) / CheckRate * float64(time.Second))
	return base.Add(l.At - back)
}

func PhraseCall(lines []CallLine, lang string, base time.Time) []CheckPiece {
	p := NewPhraser()
	p.SetLanguage(lang)
	var times []time.Time
	var out []CheckPiece
	used := 0
	take := func(ps []TextPhrase) {
		for _, ph := range ps {
			last := used + ph.Words - 1
			lw := ph.Heard
			if last < len(times) {
				lw = times[last]
			}
			used += ph.Words
			out = append(out, CheckPiece{Seq: ph.Seq, Text: ph.Text, Words: ph.Words, Para: ph.Para, Release: ph.Heard, LastWord: lw})
		}
	}
	for seg, l := range lines {
		for k := range l.Words {
			at := wordTime(base, l, k)
			times = append(times, at)
			take(p.Update(seg, strings.Join(l.Words[:k+1], " "), false, at))
		}
		end := wordTime(base, l, len(l.Words)-1).Add(PhraseQuiet)
		take(p.Update(seg, strings.Join(l.Words, " "), true, end))
	}
	if len(lines) > 0 {
		take(p.Flush(wordTime(base, lines[len(lines)-1], len(lines[len(lines)-1].Words)-1).Add(PhraseQuiet)))
	}
	return out
}

type CheckEvent struct {
	Kind     string
	Seqs     []int
	Text     string
	Result   string
	Verdict  Verdict
	Ready    time.Time
	Start    time.Time
	End      time.Time
	Duration time.Duration
	Context  string
	Terms    int
	Err      string
	Yielded  bool
	Lines    int
	Hold     time.Duration
}

type CheckCalls struct {
	Translate func(b Batch, prior []Turn, context string) (string, Verdict, time.Duration, error)
	Service   func(step Step, msgs []Message, limit time.Duration) (string, ServiceEnd, time.Duration, error)
	Draft     func(text string, prior []Turn, context string, limit time.Duration) (string, Verdict, ServiceEnd, time.Duration, error)
}

func heardBy(times []time.Time, at time.Time) int {
	n := 0
	for _, t := range times {
		if !t.After(at) {
			n++
		}
	}
	return n
}

func RunCheck(pieces []CheckPiece, target string, topic *Topic, calls CheckCalls) []CheckEvent {
	q := NewTranslateQueue(TranslateWaiting, TranslateMaxWait, TranslateMerge)
	hist := NewHistory(HistoryLimit, HistoryKeep)
	var out []CheckEvent
	if len(pieces) == 0 {
		return out
	}
	free := pieces[0].Release
	i := 0
	draftFor := -1
	var last sentence.Draft
	draftDue := func() (time.Time, bool) {
		if calls.Draft == nil || i >= len(pieces) || len(pieces[i].Times) == 0 || len(pieces[i].Parts) != len(pieces[i].Times) {
			return time.Time{}, false
		}
		if draftFor != i {
			draftFor, last = i, sentence.Draft{}
		}
		at, ok := sentence.DraftAt(pieces[i].Times, last)
		if !ok || !at.Before(pieces[i].Release) {
			return time.Time{}, false
		}
		return at, true
	}
	for i < len(pieces) || q.Len() > 0 {
		for i < len(pieces) && !pieces[i].Release.After(free) {
			p := pieces[i]
			q.Push(Pending{Seq: p.Seq, Text: p.Text, Words: p.Words, Target: target, Para: p.Para, Start: p.Release, Heard: p.Release})
			if topic != nil {
				topic.Heard(p.Text, p.Release)
			}
			i++
		}
		b, ok, stale := q.Pop(free)
		for _, s := range stale {
			out = append(out, CheckEvent{Kind: "stale", Seqs: []int{s.Seq}, Text: s.Text, Ready: s.Heard, Start: free, End: free})
		}
		if ok {
			ctx := ""
			if topic != nil {
				ctx = topic.Context()
			}
			text, verdict, d, err := calls.Translate(b, hist.Turns(), ctx)
			ev := CheckEvent{Kind: "piece", Seqs: b.Seqs, Text: b.Text, Result: text, Verdict: verdict, Ready: b.Heard, Start: free, End: free.Add(d), Duration: d, Context: ctx}
			if err != nil {
				ev.Err = err.Error()
			} else if verdict == TrOK || verdict == TrSame {
				hist.Add(Turn{Source: b.Text, Translation: text})
				if topic != nil {
					topic.Translated(b.Text, text, ev.End)
				}
			}
			out = append(out, ev)
			free = ev.End
			continue
		}
		if at, due := draftDue(); due && q.Len() == 0 && !at.After(free) {
			p := pieces[i]
			k := heardBy(p.Times, free)
			text := strings.Join(p.Parts[:k], " ")
			ctx := ""
			if topic != nil {
				ctx = topic.Context()
			}
			res, verdict, how, d, err := calls.Draft(text, hist.Turns(), ctx, p.Release.Sub(free))
			ev := CheckEvent{Kind: "draft", Seqs: []int{p.Seq}, Text: text, Result: res, Verdict: verdict, Ready: free, Start: free, End: free.Add(d), Duration: d, Yielded: how == EndYield, Lines: k}
			if err != nil {
				ev.Err = err.Error()
			}
			out = append(out, ev)
			last = sentence.Draft{Words: k, At: free}
			free = ev.End
			continue
		}
		if q.Len() == 0 && topic != nil && calls.Service != nil {
			if step := topic.Due(free); step != StepNone {
				var msgs []Message
				kind := "topic"
				if step == StepTopic {
					msgs = topic.TopicMessages(free)
				} else {
					msgs = topic.TermsMessages(free)
					kind = "terms"
				}
				var limit time.Duration
				if i < len(pieces) {
					limit = pieces[i].Release.Sub(free)
				}
				res, how, d, err := calls.Service(step, msgs, limit)
				end := free.Add(d)
				ev := CheckEvent{Kind: kind, Result: res, Ready: free, Start: free, End: end, Duration: d, Yielded: how == EndYield}
				if i < len(pieces) && end.After(pieces[i].Release) {
					ev.Hold = end.Sub(pieces[i].Release)
				}
				if err != nil && how != EndYield && how != EndTimeout {
					ev.Err = err.Error()
					topic.Failed(step, end)
				} else {
					if err != nil {
						ev.Err = err.Error()
					}
					ev.Lines = topic.Answer(step, res, how, free).Lines
				}
				ev.Context = topic.Context()
				ev.Terms = len(topic.Terms())
				out = append(out, ev)
				free = end
				continue
			}
		}
		if i < len(pieces) {
			next := pieces[i].Release
			if at, due := draftDue(); due && at.After(free) && at.Before(next) {
				next = at
			}
			if next.After(free) {
				free = next
			}
			continue
		}
		break
	}
	return out
}
