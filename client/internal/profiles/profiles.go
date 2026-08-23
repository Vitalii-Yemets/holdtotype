package profiles

const Translate = "wtranslate"

type Profile struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Prompt string `json:"prompt"`
	Hotkey string `json:"hotkey"`
}

func ByID(all []Profile, id string) *Profile {
	for i := range all {
		if all[i].ID == id {
			return &all[i]
		}
	}
	return nil
}

func Chain(all []Profile, active []string, forced string) []Profile {
	if forced == Translate {
		return nil
	}
	var ids []string
	if forced != "" {
		ids = []string{forced}
	} else {
		want := make(map[string]bool, len(active))
		for _, id := range active {
			want[id] = true
		}
		for i := range all {
			if want[all[i].ID] {
				ids = append(ids, all[i].ID)
			}
		}
	}
	var out []Profile
	seen := make(map[string]bool, len(ids))
	for _, id := range ids {
		if seen[id] {
			continue
		}
		seen[id] = true
		if p := ByID(all, id); p != nil && p.Prompt != "" {
			out = append(out, *p)
		}
	}
	return out
}
