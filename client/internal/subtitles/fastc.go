package subtitles

import (
	"html"
	"regexp"
	"strings"
	"time"
	"unicode"
)

var tagMark = regexp.MustCompile(`</?b>`)

var markupEscaper = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")

func firstWord(text string) string {
	f := strings.Fields(text)
	if len(f) == 0 {
		return ""
	}
	return f[0]
}

func BeginsSentence(prev string, para bool, text string) bool {
	if para || strings.TrimSpace(prev) == "" {
		return true
	}
	return EndsSentenceBefore(lastWord(prev), firstWord(text), false)
}

func TagPrevious(prev, text string) string {
	return "<b>" + markupEscaper.Replace(oneLine(prev)) + "</b> " + markupEscaper.Replace(oneLine(text))
}

func SplitTagged(answer string) (string, bool) {
	tags := tagMark.FindAllString(answer, -1)
	if len(tags) != 2 || tags[0] != "<b>" || tags[1] != "</b>" {
		return "", false
	}
	rest := strings.TrimLeftFunc(answer, unicode.IsSpace)
	if !strings.HasPrefix(rest, "<b>") {
		return "", false
	}
	end := strings.Index(rest, "</b>")
	inside := strings.TrimSpace(rest[len("<b>"):end])
	tail := oneLine(html.UnescapeString(rest[end+len("</b>"):]))
	if inside == "" || tail == "" {
		return "", false
	}
	return tail, true
}

func Untagged(answer string) string {
	return oneLine(html.UnescapeString(tagMark.ReplaceAllString(answer, "")))
}

func CleanFastTranslation(line, answer string) (string, Verdict) {
	out := oneLine(answer)
	p := plain(out)
	if p == "" {
		return "", TrEmpty
	}
	if p == plain(line) {
		return "", TrSame
	}
	if len([]rune(out)) > 4*len([]rune(line))+40 {
		return "", TrRamble
	}
	return KeepSentenceEnd(line, DropAddedFullStop(line, out)), TrOK
}

func isLeadingMark(r rune) bool {
	return strings.ContainsRune(",.;:!?…)]}»”’›", r)
}

func GluesToPrevious(text string) bool {
	for _, r := range text {
		return isLeadingMark(r)
	}
	return false
}

func WithoutLeadingMark(text string) string {
	return strings.TrimLeftFunc(strings.TrimLeftFunc(text, isLeadingMark), unicode.IsSpace)
}

type SpeechState string

const (
	SpeechKnown   SpeechState = "known"
	SpeechWaiting SpeechState = "waiting"
	SpeechLost    SpeechState = "lost"
)

func SpeechLanguage(setting, guess string, waited, limit time.Duration) (string, SpeechState) {
	if s := strings.ToLower(strings.TrimSpace(setting)); s != "" && s != "auto" {
		return s, SpeechKnown
	}
	if guess != "" {
		return guess, SpeechKnown
	}
	if waited >= limit {
		return "", SpeechLost
	}
	return "", SpeechWaiting
}

func plainWord(w string) string {
	return strings.ToLower(strings.TrimFunc(w, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) }))
}

func TrimRepeated(shown, part string) (string, bool) {
	var tail []string
	for _, w := range strings.Fields(shown) {
		if p := plainWord(w); p != "" {
			tail = append(tail, p)
		}
	}
	fields := strings.Fields(part)
	var at []int
	for i, w := range fields {
		if plainWord(w) != "" {
			at = append(at, i)
		}
	}
	for k := 3; k >= 1; k-- {
		if k > len(tail) || k > len(at) {
			continue
		}
		same := true
		for j := 0; j < k; j++ {
			if tail[len(tail)-k+j] != plainWord(fields[at[j]]) {
				same = false
				break
			}
		}
		if !same || (k == 1 && len([]rune(tail[len(tail)-1])) < 4) {
			continue
		}
		drop := map[int]bool{}
		for j := 0; j < k; j++ {
			drop[at[j]] = true
		}
		var kept []string
		left := false
		for i, w := range fields {
			if drop[i] {
				continue
			}
			kept = append(kept, w)
			if plainWord(w) != "" {
				left = true
			}
		}
		if !left {
			return "", false
		}
		return strings.Join(kept, " "), true
	}
	return part, true
}
