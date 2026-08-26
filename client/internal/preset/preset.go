package preset

import "strings"

type Model struct {
	ID        string
	Engine    string
	Langs     []string
	Auto      bool
	Translate bool
}

func Norm(lang string) string {
	lang = strings.ToLower(strings.TrimSpace(lang))
	if lang == "" {
		return "auto"
	}
	return lang
}

func Covers(langs []string, lang string) bool {
	lang = Norm(lang)
	for _, l := range langs {
		l = strings.ToLower(strings.TrimSpace(l))
		if l == "*" || l == lang {
			return true
		}
	}
	return false
}

func CanServe(m *Model, lang string) bool {
	if m == nil {
		return false
	}
	if Norm(lang) == "auto" {
		return m.Auto
	}
	return Covers(m.Langs, lang)
}

func Resolve(assign map[string]string, lang, def string, known func(string) bool) string {
	lang = Norm(lang)
	if id := assign[lang]; id != "" && known(id) {
		return id
	}
	if lang != "auto" {
		if id := assign["auto"]; id != "" && known(id) {
			return id
		}
	}
	return def
}

func Clean(assign map[string]string, find func(string) *Model) (map[string]string, []string) {
	if len(assign) == 0 {
		return assign, nil
	}
	out := make(map[string]string, len(assign))
	var dropped []string
	for lang, id := range assign {
		lang = Norm(lang)
		m := find(id)
		if m == nil || !CanServe(m, lang) {
			dropped = append(dropped, lang+"="+id)
			continue
		}
		out[lang] = id
	}
	return out, dropped
}
