package subwin

import "testing"

func TestTitleFollowsTheTranslation(t *testing.T) {
	if TitleKey(true) != "sub.title.tr" || TitleKey(false) != "sub.title" {
		t.Fatalf("%q %q", TitleKey(true), TitleKey(false))
	}
}

func TestCodes(t *testing.T) {
	cases := []struct {
		lang, to   string
		translated bool
		want       string
	}{
		{"en", "ru", true, "EN → RU"},
		{"auto", "ru", true, "RU"},
		{"", "de", true, "DE"},
		{"en", "ru", false, "EN"},
		{"auto", "ru", false, "AUTO"},
		{" uk ", "", true, "UK"},
		{"", "", false, ""},
	}
	for _, c := range cases {
		if got := Codes(c.lang, c.to, c.translated); got != c.want {
			t.Errorf("Codes(%q, %q, %v) = %q, want %q", c.lang, c.to, c.translated, got, c.want)
		}
	}
}
