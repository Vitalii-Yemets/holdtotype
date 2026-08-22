package routing

import "testing"

func TestAutoRoutesRussianToSherpa(t *testing.T) {
	d := Pick(Input{Mode: ModeAuto, Language: "ru", SherpaReady: true, SherpaLangs: []string{"ru"}})
	if d.Engine != Sherpa || d.Reason != ReasonLanguage {
		t.Fatalf("получено %+v, ожидался sherpa по языку", d)
	}
}

func TestAutoRoutesOtherLanguagesToWhisper(t *testing.T) {
	for _, lang := range []string{"en", "de", "uk"} {
		d := Pick(Input{Mode: ModeAuto, Language: lang, SherpaReady: true, SherpaLangs: []string{"ru"}})
		if d.Engine != Whisper || d.Reason != ReasonOtherLang {
			t.Errorf("язык %s: получено %+v, ожидался whisper", lang, d)
		}
	}
}

func TestTranslationAlwaysGoesToWhisper(t *testing.T) {
	d := Pick(Input{Mode: ModeAuto, Language: "ru", Translate: true, SherpaReady: true, SherpaLangs: []string{"ru"}})
	if d.Engine != Whisper || d.Reason != ReasonTranslate {
		t.Fatalf("получено %+v, перевод обязан идти в whisper", d)
	}
	forced := Pick(Input{Mode: ModeSherpa, Language: "ru", Translate: true, SherpaReady: true, SherpaLangs: []string{"ru"}})
	if forced.Engine != Whisper || forced.Reason != ReasonTranslate {
		t.Fatalf("получено %+v, перевод обязан идти в whisper даже при принудительном режиме", forced)
	}
}

func TestMissingSherpaFallsBack(t *testing.T) {
	d := Pick(Input{Mode: ModeAuto, Language: "ru", SherpaReady: false, SherpaLangs: []string{"ru"}})
	if d.Engine != Whisper || d.Reason != ReasonNotInstalled {
		t.Fatalf("получено %+v, без модели ожидался whisper", d)
	}
	forced := Pick(Input{Mode: ModeSherpa, Language: "ru", SherpaReady: false})
	if forced.Engine != Whisper || forced.Reason != ReasonNotInstalled {
		t.Fatalf("получено %+v, принудительный режим без модели должен откатываться", forced)
	}
}

func TestForcedModes(t *testing.T) {
	d := Pick(Input{Mode: ModeWhisper, Language: "ru", SherpaReady: true, SherpaLangs: []string{"ru"}})
	if d.Engine != Whisper || d.Reason != ReasonForced {
		t.Errorf("получено %+v, ожидался принудительный whisper", d)
	}
	s := Pick(Input{Mode: ModeSherpa, Language: "en", SherpaReady: true, SherpaLangs: []string{"ru"}})
	if s.Engine != Sherpa || s.Reason != ReasonForced {
		t.Errorf("получено %+v: принудительный sherpa должен работать и на чужом языке", s)
	}
}

func TestUnknownLanguageGoesToWhisper(t *testing.T) {
	for _, lang := range []string{"", "auto", "   "} {
		d := Pick(Input{Mode: ModeAuto, Language: lang, SherpaReady: true, SherpaLangs: []string{"ru"}})
		if d.Engine != Whisper || d.Reason != ReasonUnknownLang {
			t.Errorf("язык %q: получено %+v, ожидался whisper", lang, d)
		}
	}
}

func TestLanguageMatchIgnoresCaseAndSpaces(t *testing.T) {
	d := Pick(Input{Mode: ModeAuto, Language: " RU ", SherpaReady: true, SherpaLangs: []string{"Ru"}})
	if d.Engine != Sherpa {
		t.Fatalf("получено %+v, регистр и пробелы не должны мешать", d)
	}
}

func TestUnknownModeBehavesAsAuto(t *testing.T) {
	d := Pick(Input{Mode: "чепуха", Language: "ru", SherpaReady: true, SherpaLangs: []string{"ru"}})
	if d.Engine != Sherpa {
		t.Fatalf("получено %+v, неизвестный режим должен вести себя как авто", d)
	}
}

func TestValidMode(t *testing.T) {
	for _, m := range []string{ModeAuto, ModeWhisper, ModeSherpa} {
		if !ValidMode(m) {
			t.Errorf("режим %q должен быть допустимым", m)
		}
	}
	if ValidMode("чепуха") {
		t.Error("мусорный режим не должен считаться допустимым")
	}
}
