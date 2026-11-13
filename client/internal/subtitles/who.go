package subtitles

const (
	ChoiceQuality = "quality"
	ChoiceServer  = "server"
	ChoiceFast    = "fast"

	KindQuality = "quality"
	KindServer  = "api"
	KindFast    = "fast"

	MissingCard   = "card"
	MissingModel  = "model"
	MissingStart  = "start"
	MissingServer = "server"
	MissingPack   = "pack"
	MissingPacks  = "packs"
	MissingLang   = "lang"
	MissingFast   = "faststart"
)

type Avail struct {
	ModelReady  bool
	CardFits    bool
	QualityDown bool
	Server      bool
	Fast        bool
	FastMissing string
}

type Choice struct {
	Kind    string
	Missing string
}

func (c Choice) OK() bool { return c.Kind != "" && c.Missing == "" }

type translatorEntry struct {
	kind    string
	choice  string
	missing func(Avail) string
}

var translators = []translatorEntry{
	{KindQuality, ChoiceQuality, func(a Avail) string {
		switch {
		case !a.CardFits:
			return MissingCard
		case !a.ModelReady:
			return MissingModel
		case a.QualityDown:
			return MissingStart
		}
		return ""
	}},
	{KindServer, ChoiceServer, func(a Avail) string {
		if !a.Server {
			return MissingServer
		}
		return ""
	}},
	{KindFast, ChoiceFast, func(a Avail) string {
		switch {
		case a.Fast:
			return ""
		case a.FastMissing != "":
			return a.FastMissing
		}
		return MissingPack
	}},
}

func ValidTranslatorChoice(c string) bool {
	for _, e := range translators {
		if e.choice == c {
			return true
		}
	}
	return false
}

func ChooseTranslator(choice string, a Avail) Choice {
	for _, e := range translators {
		if e.choice == choice {
			return Choice{Kind: e.kind, Missing: e.missing(a)}
		}
	}
	return Choice{}
}

type Have struct {
	CardFits   bool
	ModelReady bool
	Server     bool
	PairReady  bool
}

func SettleTranslator(h Have) string {
	switch {
	case h.CardFits && h.ModelReady:
		return ChoiceQuality
	case h.Server:
		return ChoiceServer
	case h.PairReady:
		return ChoiceFast
	case h.CardFits:
		return ChoiceQuality
	}
	return ChoiceFast
}

func TranslationSetting(lang string, on, saved bool, ui string) (string, bool) {
	if !saved {
		on = LanguageName(lang) != ""
	}
	if LanguageName(lang) == "" {
		lang = "en"
		if LanguageName(ui) != "" {
			lang = ui
		}
	}
	return lang, on
}
