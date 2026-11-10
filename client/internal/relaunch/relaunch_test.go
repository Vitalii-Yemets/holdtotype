package relaunch

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"testing"
)

func TestTheAppPassesTheTabItWasOn(t *testing.T) {
	got := SetupArgs(`C:\Apps\HoldToType`, "post")
	want := []string{"-update", "-dir", `C:\Apps\HoldToType`, "-tab", "post"}
	if !slices.Equal(got, want) {
		t.Fatalf("got %q", got)
	}
}

func TestAnUnknownTabIsNotPassedToTheInstaller(t *testing.T) {
	for _, tab := range []string{"", "wizard", "general", "post -silent", `..\x`} {
		got := SetupArgs(`C:\x`, tab)
		if !slices.Equal(got, []string{"-update", "-dir", `C:\x`}) {
			t.Fatalf("%q gave %q", tab, got)
		}
	}
}

func TestTheInstallerReadsTheTab(t *testing.T) {
	if got := TabArg([]string{"-update", "-dir", `C:\x`, "-tab", "subs"}); got != "subs" {
		t.Fatalf("got %q", got)
	}
	if got := TabArg([]string{"-update", "-dir", `C:\x`}); got != "" {
		t.Fatalf("an older app passes no tab, got %q", got)
	}
	if got := TabArg([]string{"-update", "-tab"}); got != "" {
		t.Fatalf("a dangling flag gave %q", got)
	}
}

func TestTheRestartedAppOpensTheSameTab(t *testing.T) {
	got := AppArgs(TabArg(SetupArgs(`C:\x`, "about")))
	if !slices.Equal(got, []string{"-settings", "about", "-updated"}) {
		t.Fatalf("got %q", got)
	}
}

func TestWithoutATabTheStatusTabOpens(t *testing.T) {
	for _, tab := range []string{"", "wizard", "nonsense"} {
		got := AppArgs(tab)
		if !slices.Equal(got, []string{"-settings", "state", "-updated"}) {
			t.Fatalf("%q gave %q", tab, got)
		}
	}
}

func TestTheNoticeIsAskedForOnlyByTheFlag(t *testing.T) {
	if !Updated([]string{"-settings", "state", "-updated"}) {
		t.Fatal("the flag is missed")
	}
	if Updated(nil) || Updated([]string{"-settings", "state"}) {
		t.Fatal("an ordinary start asks for the notice")
	}
}

func TestTheNoticeIsShownOnce(t *testing.T) {
	var n Notice
	if n.Take() != "" {
		t.Fatal("a notice without an update")
	}
	n.Set("0.81.2")
	if got := n.Take(); got != "0.81.2" {
		t.Fatalf("got %q", got)
	}
	if got := n.Take(); got != "" {
		t.Fatalf("shown again: %q", got)
	}
}

func TestTheTabsAreThePagesTabs(t *testing.T) {
	src, err := os.ReadFile(filepath.Join("..", "..", "settings.go"))
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, m := range regexp.MustCompile(`<button class="nav"[^>]*data-p="([a-z]+)"`).FindAllStringSubmatch(string(src), -1) {
		seen[m[1]] = true
	}
	var page []string
	for k := range seen {
		page = append(page, k)
	}
	sort.Strings(page)
	ours := slices.Clone(Tabs)
	sort.Strings(ours)
	if !slices.Equal(page, ours) {
		t.Fatalf("the page has %q, the list has %q", page, ours)
	}
}
