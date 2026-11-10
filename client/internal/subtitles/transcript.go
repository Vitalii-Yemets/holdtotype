package subtitles

import (
	"strings"
	"time"

	"holdtotype/internal/appid"
)

const Folder = "transcripts"

func FileName(start time.Time) string {
	return start.Format("2006-01-02 15-04-05") + ".txt"
}

func oneLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func Header(start time.Time, lang, model string) string {
	return appid.Name + " subtitles · " + start.Format("2006-01-02 15:04") +
		" · language: " + oneLine(lang) + " · model: " + oneLine(model)
}

func LineText(at time.Time, text string) string {
	return "[" + at.Format("15:04:05") + "] " + oneLine(text)
}

func Opening(start time.Time, lang, model, translation string) string {
	h := Header(start, lang, model)
	if t := oneLine(translation); t != "" {
		h += " · translation: " + t
	}
	return h + "\r\n\r\n"
}

func Entry(at time.Time, text string) string {
	return LineText(at, text) + "\r\n"
}

func TranslationText(at time.Time, text string) string {
	return "[" + at.Format("15:04:05") + "] → " + oneLine(text)
}

func TranslationEntry(at time.Time, text string) string {
	return TranslationText(at, text) + "\r\n"
}
