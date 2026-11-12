package subtitles

import "testing"

func availOf(mask int) Avail {
	a := Avail{
		ModelReady:  mask&1 != 0,
		CardFits:    mask&2 != 0,
		QualityDown: mask&4 != 0,
		Server:      mask&8 != 0,
		Editor:      mask&16 != 0,
		Fast:        mask&32 != 0,
	}
	if !a.Fast {
		a.FastMissing = []string{"", MissingPack, MissingLang, MissingFast}[mask%4]
	}
	return a
}

func TestChooseTranslatorAutomatically(t *testing.T) {
	for mask := 0; mask < 64; mask++ {
		a := availOf(mask)
		want := Choice{Missing: MissingAll}
		switch {
		case a.ModelReady && a.CardFits && !a.QualityDown:
			want = Choice{Kind: KindQuality}
		case a.Server:
			want = Choice{Kind: KindServer}
		case a.Fast:
			want = Choice{Kind: KindFast}
		case a.Editor:
			want = Choice{Kind: KindEditor}
		}
		got := ChooseTranslator(ChoiceAuto, a)
		if got != want {
			t.Errorf("%+v: %+v, want %+v", a, got, want)
		}
		if got.OK() != (want.Kind != "") {
			t.Errorf("%+v: OK %v", a, got.OK())
		}
	}
}

func TestChooseTranslatorExplicitly(t *testing.T) {
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
			t.Errorf("%+v: an explicit choice was replaced: %+v %+v %+v", a, q, s, f)
		}
	}
}

func TestTranslatorChoices(t *testing.T) {
	for _, c := range []string{"", "quality", "server", "fast"} {
		if !ValidTranslatorChoice(c) {
			t.Errorf("%q is not valid", c)
		}
	}
	for _, c := range []string{"local", "api", "Fast", "Quality", " server"} {
		if ValidTranslatorChoice(c) {
			t.Errorf("%q is valid", c)
		}
	}
	if got := ChooseTranslator("unknown", Avail{Server: true}); got != (Choice{Kind: KindServer}) {
		t.Errorf("an unknown choice: %+v", got)
	}
	if got := ChooseTranslator(ChoiceAuto, Avail{Fast: true, Editor: true}); got != (Choice{Kind: KindFast}) {
		t.Errorf("fast before the editor model: %+v", got)
	}
	if got := ChooseTranslator(ChoiceAuto, Avail{Fast: true, Server: true}); got != (Choice{Kind: KindServer}) {
		t.Errorf("the server before fast: %+v", got)
	}
}
