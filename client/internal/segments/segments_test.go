package segments

import "testing"

func TestJoin(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"plain text is left alone", "привет как дела", "привет как дела"},
		{"one break becomes a space", "первый кусок\nвторой кусок", "первый кусок второй кусок"},
		{"windows breaks count too", "первый\r\nвторой", "первый второй"},
		{"a lone carriage return counts", "первый\rвторой", "первый второй"},
		{"empty lines do not leave double spaces", "первый\n\n\nвторой", "первый второй"},
		{"spaces around the break are eaten", "первый  \n  второй", "первый второй"},
		{"nothing but breaks comes back empty", "\n\n", ""},
		{"an empty string stays empty", "", ""},
		{"inner spacing is untouched", "два  пробела\nздесь", "два  пробела здесь"},
	}
	for _, c := range cases {
		if got := Join(c.in); got != c.want {
			t.Errorf("%s: Join(%q) = %q, want %q", c.name, c.in, got, c.want)
		}
	}
}
