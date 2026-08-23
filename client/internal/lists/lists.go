package lists

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"holdtotype/internal/commands"
	"holdtotype/internal/replace"
)

const (
	Kind    = "holdtotype-lists"
	Version = 1
)

type File struct {
	Kind         string             `json:"kind"`
	Version      int                `json:"version"`
	Replacements []replace.Rule     `json:"replacements"`
	Commands     []commands.Command `json:"commands"`
}

var ErrNotOurs = errors.New("файл не похож на списки")

func Encode(rules []replace.Rule, cmds []commands.Command) ([]byte, error) {
	f := File{Kind: Kind, Version: Version, Replacements: replace.Clean(rules), Commands: commands.Clean(cmds)}
	if f.Replacements == nil {
		f.Replacements = []replace.Rule{}
	}
	if f.Commands == nil {
		f.Commands = []commands.Command{}
	}
	return json.MarshalIndent(f, "", "  ")
}

func Parse(data []byte) (File, error) {
	var f File
	if err := json.Unmarshal(data, &f); err != nil {
		return File{}, ErrNotOurs
	}
	if f.Kind != Kind {
		return File{}, ErrNotOurs
	}
	f.Replacements = replace.Clean(f.Replacements)
	f.Commands = commands.Clean(f.Commands)
	return f, nil
}

func key(s string) string { return strings.ToLower(strings.TrimSpace(s)) }

func freeID(want string, taken map[string]bool) string {
	if want == "" {
		want = "imp"
	}
	id := want
	for i := 2; taken[id]; i++ {
		id = want + "-" + strconv.Itoa(i)
	}
	taken[id] = true
	return id
}

func MergeRules(cur, add []replace.Rule) ([]replace.Rule, int, int) {
	seen := map[string]bool{}
	ids := map[string]bool{}
	out := append([]replace.Rule{}, cur...)
	for _, r := range out {
		seen[key(r.From)] = true
		ids[r.ID] = true
	}
	added, skipped := 0, 0
	for _, r := range add {
		if seen[key(r.From)] {
			skipped++
			continue
		}
		seen[key(r.From)] = true
		r.ID = freeID(r.ID, ids)
		out = append(out, r)
		added++
	}
	return out, added, skipped
}

func MergeCommands(cur, add []commands.Command) ([]commands.Command, int, int) {
	seen := map[string]bool{}
	ids := map[string]bool{}
	out := append([]commands.Command{}, cur...)
	for _, c := range out {
		seen[key(c.Phrase)] = true
		ids[c.ID] = true
	}
	added, skipped := 0, 0
	for _, c := range add {
		if seen[key(c.Phrase)] {
			skipped++
			continue
		}
		seen[key(c.Phrase)] = true
		c.ID = freeID(c.ID, ids)
		out = append(out, c)
		added++
	}
	return out, added, skipped
}
