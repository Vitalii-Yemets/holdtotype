package appid

import (
	"strings"
	"testing"

	"holdtotype/internal/appver"
)

func TestVersionParses(t *testing.T) {
	if _, ok := appver.Parse(Version); !ok {
		t.Fatalf("Version %q не разбирается как x.y.z", Version)
	}
}

func TestSlugIsFileSafe(t *testing.T) {
	if Slug == "" {
		t.Fatal("пустой Slug")
	}
	for _, r := range Slug {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
			t.Fatalf("Slug содержит %q — недопустимо в именах файлов и URL", r)
		}
	}
}

func TestDerivedNames(t *testing.T) {
	cases := map[string]string{
		Exe:      ".exe",
		SetupExe: ".exe",
		LogFile:  ".log",
		Portable: ".zip",
	}
	for name, ext := range cases {
		if !strings.HasPrefix(name, Slug) {
			t.Errorf("%q не начинается со Slug", name)
		}
		if !strings.HasSuffix(name, ext) {
			t.Errorf("%q не оканчивается на %q", name, ext)
		}
	}
	if Exe == SetupExe {
		t.Error("имена приложения и установщика совпали")
	}
}

func TestRegistryAndURLsCarryName(t *testing.T) {
	if !strings.HasSuffix(UninstallKey, Name) {
		t.Errorf("ключ удаления %q не оканчивается именем", UninstallKey)
	}
	if !strings.Contains(MutexName, Name) {
		t.Errorf("имя мьютекса %q не содержит имени приложения", MutexName)
	}
	if !strings.Contains(LatestAPI, Repo) || !strings.Contains(RepoURL, Repo) {
		t.Error("ссылки не совпадают с Repo")
	}
	if !strings.HasSuffix(RepoURL, Slug) {
		t.Errorf("репозиторий %q не совпадает со Slug", RepoURL)
	}
}

func TestTempDirName(t *testing.T) {
	got := TempDirName("webview", 421)
	want := Slug + "-webview-421"
	if got != want {
		t.Fatalf("TempDirName = %q, ожидалось %q", got, want)
	}
	if !strings.HasPrefix(got, TempDirPrefix("webview")) {
		t.Fatal("префикс не совпадает с именем каталога")
	}
}

func TestClass(t *testing.T) {
	if Class("TrayWnd") != Name+"TrayWnd" {
		t.Fatal("Class собирает имя класса неверно")
	}
}
