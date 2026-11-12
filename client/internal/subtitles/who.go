package subtitles

const (
	ChoiceAuto    = ""
	ChoiceQuality = "quality"
	ChoiceServer  = "server"
	ChoiceFast    = "fast"

	KindQuality = "quality"
	KindServer  = "api"
	KindEditor  = "local"
	KindFast    = "fast"

	MissingCard   = "card"
	MissingModel  = "model"
	MissingStart  = "start"
	MissingServer = "server"
	MissingPack   = "pack"
	MissingLang   = "lang"
	MissingFast   = "faststart"
	MissingAll    = "none"
)

type Avail struct {
	ModelReady  bool
	CardFits    bool
	QualityDown bool
	Server      bool
	Editor      bool
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

var translatorOrder = []translatorEntry{
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
	{KindEditor, "", func(a Avail) string {
		if !a.Editor {
			return MissingAll
		}
		return ""
	}},
}

func ValidTranslatorChoice(c string) bool {
	if c == ChoiceAuto {
		return true
	}
	for _, e := range translatorOrder {
		if e.choice != "" && e.choice == c {
			return true
		}
	}
	return false
}

func ChooseTranslator(choice string, a Avail) Choice {
	if choice != ChoiceAuto {
		for _, e := range translatorOrder {
			if e.choice == choice {
				return Choice{Kind: e.kind, Missing: e.missing(a)}
			}
		}
	}
	for _, e := range translatorOrder {
		if e.missing(a) == "" {
			return Choice{Kind: e.kind}
		}
	}
	return Choice{Missing: MissingAll}
}
