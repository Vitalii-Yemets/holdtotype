package relaunch

import (
	"slices"
	"sync"
)

const (
	FallbackTab = "state"
	UpdatedFlag = "-updated"
	TabFlag     = "-tab"
)

var Tabs = []string{"state", "system", "mic", "history", "dictation", "models", "text", "post", "subs", "help", "about", "contacts"}

func ValidTab(tab string) bool {
	return slices.Contains(Tabs, tab)
}

func SetupArgs(dir, tab string) []string {
	args := []string{"-update", "-dir", dir}
	if ValidTab(tab) {
		args = append(args, TabFlag, tab)
	}
	return args
}

func TabArg(args []string) string {
	for i, a := range args {
		if a == TabFlag && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

func AppArgs(tab string) []string {
	if !ValidTab(tab) {
		tab = FallbackTab
	}
	return []string{"-settings", tab, UpdatedFlag}
}

func Updated(args []string) bool {
	return slices.Contains(args, UpdatedFlag)
}

type Notice struct {
	mu sync.Mutex
	v  string
}

func (n *Notice) Set(v string) {
	n.mu.Lock()
	n.v = v
	n.mu.Unlock()
}

func (n *Notice) Take() string {
	n.mu.Lock()
	defer n.mu.Unlock()
	v := n.v
	n.v = ""
	return v
}
