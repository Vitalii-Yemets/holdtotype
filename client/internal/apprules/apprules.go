package apprules

import "strings"

const (
	PasteInherit   = ""
	PasteClipboard = "clipboard"
	PasteType      = "type"

	EnterInherit = ""
	EnterOn      = "on"
	EnterOff     = "off"
)

type Rule struct {
	ID          string   `json:"id"`
	Match       string   `json:"match"`
	Paste       string   `json:"paste_mode"`
	Enter       string   `json:"auto_enter"`
	DelayMs     int      `json:"delay_ms"`
	UseProfiles bool     `json:"use_profiles"`
	Profiles    []string `json:"profiles"`
}

func ValidPaste(v string) bool {
	return v == PasteInherit || v == PasteClipboard || v == PasteType
}

func ValidEnter(v string) bool {
	return v == EnterInherit || v == EnterOn || v == EnterOff
}

func exeName(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.Trim(s, `"`)
	if i := strings.LastIndexAny(s, `\/`); i >= 0 {
		s = s[i+1:]
	}
	return s
}

func matches(pattern, target string) bool {
	p := exeName(pattern)
	if p == "" || target == "" {
		return false
	}
	if strings.HasSuffix(p, "*") {
		return strings.HasPrefix(target, strings.TrimSuffix(p, "*"))
	}
	if !strings.Contains(p, ".") {
		return p == strings.TrimSuffix(target, ".exe")
	}
	return p == target
}

func Find(rules []Rule, exe string) (Rule, bool) {
	target := exeName(exe)
	if target == "" {
		return Rule{}, false
	}
	for _, r := range rules {
		for _, part := range strings.Split(r.Match, ",") {
			if matches(part, target) {
				return r, true
			}
		}
	}
	return Rule{}, false
}

func Clean(rules []Rule) []Rule {
	out := make([]Rule, 0, len(rules))
	for _, r := range rules {
		r.Match = strings.TrimSpace(r.Match)
		if r.Match == "" {
			continue
		}
		if !ValidPaste(r.Paste) {
			r.Paste = PasteInherit
		}
		if !ValidEnter(r.Enter) {
			r.Enter = EnterInherit
		}
		if r.DelayMs < 0 || r.DelayMs > 5000 {
			r.DelayMs = 0
		}
		if !r.UseProfiles {
			r.Profiles = nil
		}
		out = append(out, r)
	}
	return out
}
