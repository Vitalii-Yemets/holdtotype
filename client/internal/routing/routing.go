package routing

import "strings"

const (
	Whisper = "whisper"
	Sherpa  = "sherpa"
)

const (
	ModeAuto    = "auto"
	ModeWhisper = "whisper"
	ModeSherpa  = "sherpa"
)

const (
	ReasonForced       = "forced"
	ReasonTranslate    = "translate"
	ReasonNotInstalled = "notinstalled"
	ReasonLanguage     = "language"
	ReasonOtherLang    = "otherlang"
	ReasonUnknownLang  = "unknownlang"
)

type Input struct {
	Mode        string
	Language    string
	Translate   bool
	SherpaReady bool
	SherpaLangs []string
}

type Decision struct {
	Engine string
	Reason string
}

func ValidMode(m string) bool {
	return m == ModeAuto || m == ModeWhisper || m == ModeSherpa
}

func covers(langs []string, lang string) bool {
	lang = strings.ToLower(strings.TrimSpace(lang))
	if lang == "" {
		return false
	}
	for _, l := range langs {
		if strings.ToLower(strings.TrimSpace(l)) == lang {
			return true
		}
	}
	return false
}

func Pick(in Input) Decision {
	switch in.Mode {
	case ModeWhisper:
		return Decision{Whisper, ReasonForced}
	case ModeSherpa:
		if !in.SherpaReady {
			return Decision{Whisper, ReasonNotInstalled}
		}
		if in.Translate {
			return Decision{Whisper, ReasonTranslate}
		}
		return Decision{Sherpa, ReasonForced}
	}

	if in.Translate {
		return Decision{Whisper, ReasonTranslate}
	}
	if !in.SherpaReady {
		return Decision{Whisper, ReasonNotInstalled}
	}
	lang := strings.ToLower(strings.TrimSpace(in.Language))
	if lang == "" || lang == "auto" {
		return Decision{Whisper, ReasonUnknownLang}
	}
	if covers(in.SherpaLangs, lang) {
		return Decision{Sherpa, ReasonLanguage}
	}
	return Decision{Whisper, ReasonOtherLang}
}
