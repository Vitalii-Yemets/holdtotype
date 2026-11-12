package subtitles

import "testing"

func TestChooseTranslatorAutomatically(t *testing.T) {
	for mask := 0; mask < 32; mask++ {
		a := Avail{
			ModelReady:  mask&1 != 0,
			CardFits:    mask&2 != 0,
			QualityDown: mask&4 != 0,
			Server:      mask&8 != 0,
			Editor:      mask&16 != 0,
		}
		want := Choice{Missing: MissingAll}
		switch {
		case a.ModelReady && a.CardFits && !a.QualityDown:
			want = Choice{Kind: KindQuality}
		case a.Server:
			want = Choice{Kind: KindServer}
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
	for mask := 0; mask < 32; mask++ {
		a := Avail{
			ModelReady:  mask&1 != 0,
			CardFits:    mask&2 != 0,
			QualityDown: mask&4 != 0,
			Server:      mask&8 != 0,
			Editor:      mask&16 != 0,
		}
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
		if q.Kind != KindQuality || s.Kind != KindServer {
			t.Errorf("%+v: an explicit choice was replaced: %+v %+v", a, q, s)
		}
	}
}

func TestTranslatorChoices(t *testing.T) {
	for _, c := range []string{"", "quality", "server"} {
		if !ValidTranslatorChoice(c) {
			t.Errorf("%q is not valid", c)
		}
	}
	for _, c := range []string{"local", "api", "fast", "Quality", " server"} {
		if ValidTranslatorChoice(c) {
			t.Errorf("%q is valid", c)
		}
	}
	if got := ChooseTranslator("unknown", Avail{Server: true}); got != (Choice{Kind: KindServer}) {
		t.Errorf("an unknown choice: %+v", got)
	}
}
