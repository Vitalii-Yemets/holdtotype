package subwin

import "strings"

const (
	TitlePlain       = "sub.title"
	TitleTranslation = "sub.title.tr"
)

func TitleKey(translated bool) string {
	if translated {
		return TitleTranslation
	}
	return TitlePlain
}

func Codes(lang, to string, translated bool) string {
	lang = strings.ToUpper(strings.TrimSpace(lang))
	to = strings.ToUpper(strings.TrimSpace(to))
	if !translated || to == "" {
		return lang
	}
	if lang == "" || lang == "AUTO" {
		return to
	}
	return lang + " → " + to
}
