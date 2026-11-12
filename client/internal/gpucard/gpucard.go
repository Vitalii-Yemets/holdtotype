package gpucard

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	NeedMB  = 4096
	SlackMB = 256
)

type Adapter struct {
	Name        string
	DedicatedMB int
	Integrated  bool
	Hardware    bool
}

func Fits(a Adapter) bool {
	return a.Hardware && !a.Integrated && a.DedicatedMB >= NeedMB-SlackMB
}

func Pick(list []Adapter) (Adapter, string) {
	var fit []Adapter
	for _, a := range list {
		if Fits(a) {
			fit = append(fit, a)
		}
	}
	if len(fit) > 0 {
		sort.SliceStable(fit, func(i, j int) bool { return fit[i].DedicatedMB > fit[j].DedicatedMB })
		return fit[0], ""
	}
	return Adapter{}, Why(list)
}

const (
	ReasonNone       = "none"
	ReasonIntegrated = "integrated"
	ReasonSmall      = "small"
	ReasonUnknown    = "unknown"
)

func Reason(list []Adapter) (string, int) {
	if pick, why := Pick(list); why == "" {
		return "", pick.DedicatedMB
	}
	best := -1
	integrated := false
	for _, a := range list {
		if !a.Hardware {
			continue
		}
		if a.Integrated {
			integrated = true
			continue
		}
		if a.DedicatedMB > best {
			best = a.DedicatedMB
		}
	}
	switch {
	case best >= 0:
		return ReasonSmall, best
	case integrated:
		return ReasonIntegrated, 0
	}
	return ReasonNone, 0
}

func Why(list []Adapter) string {
	best := -1
	integrated := false
	for _, a := range list {
		if !a.Hardware {
			continue
		}
		if a.Integrated {
			integrated = true
			continue
		}
		if a.DedicatedMB > best {
			best = a.DedicatedMB
		}
	}
	switch {
	case best >= 0:
		return fmt.Sprintf("the graphics card has %d MB of its own memory, %d MB are needed", best, NeedMB)
	case integrated:
		return "only integrated graphics, which share the computer's memory"
	}
	return "no graphics card is reported"
}

type Device struct {
	ID      string
	Name    string
	TotalMB int
	FreeMB  int
}

var deviceLine = regexp.MustCompile(`^\s*(Vulkan\d+):\s*(.*?)\s*\((\d+)\s*MiB,\s*(\d+)\s*MiB free\)\s*$`)

func ParseDevices(out string) []Device {
	var list []Device
	for _, line := range strings.Split(strings.ReplaceAll(out, "\r", ""), "\n") {
		m := deviceLine.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		total, _ := strconv.Atoi(m[3])
		free, _ := strconv.Atoi(m[4])
		list = append(list, Device{ID: m[1], Name: m[2], TotalMB: total, FreeMB: free})
	}
	return list
}

func plainName(s string) string {
	s = strings.ToLower(s)
	for _, junk := range []string{"(r)", "(tm)", "®", "™"} {
		s = strings.ReplaceAll(s, junk, "")
	}
	return strings.Join(strings.Fields(s), " ")
}

func Match(list []Device, name string) (Device, bool) {
	want := plainName(name)
	if want == "" {
		return Device{}, false
	}
	for _, d := range list {
		if plainName(d.Name) == want {
			return d, true
		}
	}
	var hit []Device
	for _, d := range list {
		got := plainName(d.Name)
		if got != "" && (strings.Contains(got, want) || strings.Contains(want, got)) {
			hit = append(hit, d)
		}
	}
	if len(hit) == 1 {
		return hit[0], true
	}
	return Device{}, false
}
