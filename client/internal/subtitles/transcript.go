package subtitles

import (
	"strings"
	"time"

	"holdtotype/internal/appid"
)

const Folder = "transcripts"

func fileStamp(start time.Time) string {
	return start.Format("2006-01-02_15-04-05")
}

func langCode(lang string) string {
	l := strings.ToUpper(strings.TrimSpace(lang))
	if l == "" || l == "AUTO" || strings.ContainsFunc(l, func(r rune) bool { return r < 'A' || r > 'Z' }) {
		return ""
	}
	return l
}

func OriginalFileName(start time.Time, lang string, translating bool) string {
	name := fileStamp(start)
	if translating {
		name += "_original"
		if code := langCode(lang); code != "" {
			name += "_" + code
		}
	}
	return name + ".txt"
}

func TranslationFileName(start time.Time, to string) string {
	name := fileStamp(start) + "_translation"
	if code := langCode(to); code != "" {
		name += "_" + code
	}
	return name + ".txt"
}

func oneLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func headStart(start time.Time) string {
	return appid.Name + " subtitles · " + start.Format("2006-01-02 15:04")
}

func Header(start time.Time, lang, model string) string {
	return headStart(start) + " · language: " + oneLine(lang) + " · model: " + oneLine(model)
}

func OriginalHeader(start time.Time, lang, model string, translating bool) string {
	if !translating {
		return Header(start, lang, model)
	}
	return headStart(start) + " · original · language: " + oneLine(lang) + " · model: " + oneLine(model)
}

func TranslationHeader(start time.Time, from, to, by string) string {
	src := strings.ToLower(oneLine(from))
	if langCode(src) == "" {
		src = "auto"
	}
	return headStart(start) + " · translation: " + src + " → " + strings.ToLower(oneLine(to)) + " · by: " + oneLine(by)
}

func Opening(header string) string {
	return header + "\r\n\r\n"
}

func LineText(at time.Time, text string) string {
	return "[" + at.Format("15:04:05") + "] " + oneLine(text)
}

func Entry(at time.Time, text string) string {
	return LineText(at, text) + "\r\n"
}
