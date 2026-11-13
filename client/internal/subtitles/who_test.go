package subtitles

import "testing"

func availOf(mask int) Avail {
	a := Avail{
		ModelReady:  mask&1 != 0,
		CardFits:    mask&2 != 0,
		QualityDown: mask&4 != 0,
		Server:      mask&8 != 0,
		Fast:        mask&16 != 0,
	}
	if !a.Fast {
		a.FastMissing = []string{"", MissingPack, MissingPacks, MissingLang, MissingFast}[mask%5]
	}
	return a
}

func TestChooseTranslatorWithoutAnAutomaticOrder(t *testing.T) {
	for mask := 0; mask < 64; mask++ {
		a := availOf(mask)
		q := ChooseTranslator(ChoiceQuality, a)
		wantQ := Choice{Kind: KindQuality}
		switch {
		case !a.CardFits:
			wantQ.Missing = MissingCard
		case !a.ModelReady:
			wantQ.Missing = MissingModel
		case a.QualityDown:
			wantQ.Missing = MissingStart
		}
		if q != wantQ {
			t.Errorf("quality %+v: %+v, want %+v", a, q, wantQ)
		}
		s := ChooseTranslator(ChoiceServer, a)
		wantS := Choice{Kind: KindServer}
		if !a.Server {
			wantS.Missing = MissingServer
		}
		if s != wantS {
			t.Errorf("server %+v: %+v, want %+v", a, s, wantS)
		}
		f := ChooseTranslator(ChoiceFast, a)
		wantF := Choice{Kind: KindFast}
		if !a.Fast {
			wantF.Missing = a.FastMissing
			if wantF.Missing == "" {
				wantF.Missing = MissingPack
			}
		}
		if f != wantF {
			t.Errorf("fast %+v: %+v, want %+v", a, f, wantF)
		}
		if q.Kind != KindQuality || s.Kind != KindServer || f.Kind != KindFast {
			t.Errorf("%+v: a choice was replaced by another translator: %+v %+v %+v", a, q, s, f)
		}
	}
}

func TestTranslatorChoices(t *testing.T) {
	for _, c := range []string{"quality", "server", "fast"} {
		if !ValidTranslatorChoice(c) {
			t.Errorf("%q is not valid", c)
		}
	}
	for _, c := range []string{"", "local", "api", "Fast", "Quality", " server"} {
		if ValidTranslatorChoice(c) {
			t.Errorf("%q is valid", c)
		}
	}
	all := Avail{ModelReady: true, CardFits: true, Server: true, Fast: true}
	for _, c := range []string{"", "unknown", "local"} {
		if got := ChooseTranslator(c, all); got.OK() || got.Kind != "" {
			t.Errorf("%q chooses %+v", c, got)
		}
	}
}

func TestSettleTranslatorFromAnOldSettingsFile(t *testing.T) {
	for mask := 0; mask < 16; mask++ {
		h := Have{CardFits: mask&1 != 0, ModelReady: mask&2 != 0, Server: mask&4 != 0, PairReady: mask&8 != 0}
		var want string
		switch {
		case h.CardFits && h.ModelReady:
			want = ChoiceQuality
		case h.Server:
			want = ChoiceServer
		case h.PairReady:
			want = ChoiceFast
		case h.CardFits:
			want = ChoiceQuality
		default:
			want = ChoiceFast
		}
		got := SettleTranslator(h)
		if got != want {
			t.Errorf("%+v: %q, want %q", h, got, want)
		}
		if !ValidTranslatorChoice(got) {
			t.Errorf("%+v: %q is not a choice", h, got)
		}
	}
	if got := SettleTranslator(Have{}); got != ChoiceFast {
		t.Errorf("a fresh install without a fitting card: %q", got)
	}
	if got := SettleTranslator(Have{CardFits: true}); got != ChoiceQuality {
		t.Errorf("a fresh install with a fitting card: %q", got)
	}
	if got := SettleTranslator(Have{ModelReady: true, Server: true}); got != ChoiceServer {
		t.Errorf("a model on disk without a fitting card does not win over the server: %q", got)
	}
}

func TestTheTranslationSettingOfAnOldSettingsFile(t *testing.T) {
	cases := []struct {
		lang      string
		on, saved bool
		ui        string
		wantLang  string
		wantOn    bool
	}{
		{"ru", false, false, "en", "ru", true},
		{"de", false, false, "ru", "de", true},
		{"", false, false, "ru", "ru", false},
		{"", false, false, "pl", "pl", false},
		{"", false, false, "zz", "en", false},
		{"xx", false, false, "uk", "uk", false},
		{"fr", false, true, "ru", "fr", false},
		{"fr", true, true, "ru", "fr", true},
		{"", true, true, "de", "de", true},
	}
	for _, c := range cases {
		lang, on := TranslationSetting(c.lang, c.on, c.saved, c.ui)
		if lang != c.wantLang || on != c.wantOn {
			t.Errorf("%+v: %q %v", c, lang, on)
		}
	}
	for _, ui := range TranslateTargets {
		if lang, on := TranslationSetting("", false, false, ui); lang != ui || on {
			t.Errorf("nothing chosen, interface %s: %q %v", ui, lang, on)
		}
	}
}
