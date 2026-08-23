package appid

import "testing"

func TestIdentityIsOneName(t *testing.T) {
	if Name != "HoldToType" || Slug != "holdtotype" {
		t.Fatalf("identity drifted: %s / %s", Name, Slug)
	}
	for _, s := range []string{Exe, SetupExe, LogFile, Portable, LLMAlias, Repo, RepoURL, LatestAPI} {
		if !contains(s, Slug) {
			t.Errorf("%q does not carry the slug", s)
		}
	}
	for _, s := range []string{InstallDirName, ShortcutName, RunValue, UninstallKey, MutexName} {
		if !contains(s, Name) {
			t.Errorf("%q does not carry the name", s)
		}
	}
}

func TestPreviousNameIsOnlyUsedForCleanup(t *testing.T) {
	if PrevSlug == Slug {
		t.Fatal("the previous name must differ from the current one")
	}
	if got := PrevTempPrefix("webview"); got != "voxterminal-webview" {
		t.Fatalf("PrevTempPrefix = %q", got)
	}
	for _, s := range []string{Exe, SetupExe, LogFile, Portable, InstallDirName, ShortcutName, RunValue, UninstallKey, MutexName, Repo, LLMAlias} {
		if contains(s, PrevSlug) || contains(s, "Vox") {
			t.Errorf("%q still carries the old name", s)
		}
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
