package uninst

import (
	"slices"
	"testing"
)

func TestKeepingDataLeavesTranscriptsAlone(t *testing.T) {
	files, dirs := Plan(false)
	if len(dirs) != 0 {
		t.Fatalf("folders removed without asking: %v", dirs)
	}
	if slices.Contains(files, "config.json") {
		t.Fatal("the settings are removed without asking")
	}
	if !slices.Contains(files, "sherpa-online-server.exe") {
		t.Fatal("the program files are not removed")
	}
}

func TestDeletingDataTakesTranscriptsWithIt(t *testing.T) {
	files, dirs := Plan(true)
	for _, d := range []string{"models", "transcripts"} {
		if !slices.Contains(dirs, d) {
			t.Fatalf("%s is not removed with the data: %v", d, dirs)
		}
	}
	for _, f := range []string{"config.json", "holdtotype.log", "holdtotype-history.json"} {
		if !slices.Contains(files, f) {
			t.Fatalf("%s is not removed with the data", f)
		}
	}
}
