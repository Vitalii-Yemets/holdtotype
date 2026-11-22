package appname

import "testing"

func TestChoose(t *testing.T) {
	cases := []struct {
		desc, product, path, want string
	}{
		{"Google Chrome", "Google Chrome", `C:\Program Files\Google\Chrome\Application\chrome.exe`, "Google Chrome"},
		{"", "Microsoft Office", `C:\Office\WINWORD.EXE`, "Microsoft Office"},
		{"  ", "\x00", `C:\Tools\keepass.exe`, "keepass"},
		{"Telegram  Desktop\x00", "Telegram", `C:\Telegram\Telegram.exe`, "Telegram Desktop"},
		{"", "", `D:/apps/Code.EXE`, "Code"},
		{"", "", `C:\bin\tool`, "tool"},
		{"", "", "", ""},
		{"", "", `C:\x\.exe`, ".exe"},
	}
	for _, c := range cases {
		if got := Choose(c.desc, c.product, c.path); got != c.want {
			t.Errorf("Choose(%q, %q, %q) = %q, want %q", c.desc, c.product, c.path, got, c.want)
		}
	}
}

func TestCacheReadsEachProgramOnce(t *testing.T) {
	reads := map[string]int{}
	c := NewCache(func(path string) (string, string) {
		reads[path]++
		if path == `C:\a\chrome.exe` {
			return "Google Chrome", ""
		}
		return "", ""
	})
	for i := 0; i < 3; i++ {
		if got := c.Name(`C:\a\chrome.exe`); got != "Google Chrome" {
			t.Fatalf("got %q", got)
		}
		if got := c.Name(`C:\A\CHROME.EXE`); got != "Google Chrome" {
			t.Fatalf("the same program by another case got %q", got)
		}
		if got := c.Name(`C:\b\tool.exe`); got != "tool" {
			t.Fatalf("got %q", got)
		}
	}
	if reads[`C:\a\chrome.exe`] != 1 || reads[`C:\b\tool.exe`] != 1 || len(reads) != 2 {
		t.Fatalf("version resources were read %v", reads)
	}
	if c.Name("  ") != "" || len(reads) != 2 {
		t.Fatal("an empty path is not a program")
	}
}
