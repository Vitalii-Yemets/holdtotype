package commands

import (
	"strings"

	"holdtotype/internal/replace"
)

const (
	ActionNewline   = "newline"
	ActionParagraph = "paragraph"
	ActionText      = "text"
	ActionCancel    = "cancel"
)

type Command struct {
	ID     string `json:"id"`
	Phrase string `json:"phrase"`
	Action string `json:"action"`
	Text   string `json:"text"`
}

type Result struct {
	Text      string
	Cancelled bool
	Applied   []string
}

func ValidAction(v string) bool {
	switch v {
	case ActionNewline, ActionParagraph, ActionText, ActionCancel:
		return true
	}
	return false
}

func Clean(cmds []Command) []Command {
	out := make([]Command, 0, len(cmds))
	for _, c := range cmds {
		c.Phrase = strings.TrimSpace(c.Phrase)
		if c.Phrase == "" {
			continue
		}
		if !ValidAction(c.Action) {
			c.Action = ActionNewline
		}
		if c.Action != ActionText {
			c.Text = ""
		}
		out = append(out, c)
	}
	return out
}

func output(c Command) string {
	switch c.Action {
	case ActionNewline:
		return "\n"
	case ActionParagraph:
		return "\n\n"
	case ActionText:
		return c.Text
	}
	return ""
}

func tidy(s string) string {
	for _, pair := range [][2]string{{" \n", "\n"}, {"\n ", "\n"}, {"\t\n", "\n"}, {"\n\t", "\n"}} {
		for strings.Contains(s, pair[0]) {
			s = strings.ReplaceAll(s, pair[0], pair[1])
		}
	}
	for strings.Contains(s, "\n\n\n") {
		s = strings.ReplaceAll(s, "\n\n\n", "\n\n")
	}
	return strings.Trim(s, " \t\n")
}

func Apply(cmds []Command, text string) Result {
	res := Result{Text: text}
	if text == "" || len(cmds) == 0 {
		return res
	}
	for _, c := range cmds {
		if strings.TrimSpace(c.Phrase) == "" {
			continue
		}
		rule := replace.Rule{From: c.Phrase, To: output(c), Whole: true}
		if c.Action == ActionCancel {
			if replace.Contains(rule, res.Text) {
				res.Cancelled = true
				res.Applied = append(res.Applied, c.Phrase)
				return res
			}
			continue
		}
		before := res.Text
		res.Text = replace.Apply([]replace.Rule{rule}, res.Text)
		if res.Text != before {
			res.Applied = append(res.Applied, c.Phrase)
		}
	}
	res.Text = tidy(res.Text)
	return res
}
