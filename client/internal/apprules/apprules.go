package apprules

import "strings"

type Rule struct {
	ID    string `json:"id"`
	Match string `json:"match"`
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
