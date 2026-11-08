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
	ActionUndo      = "undo"
	ActionHeading   = "heading"
	ActionList      = "list"
	ActionCapital   = "capital"
	ActionQuotes    = "quotes"
	ActionNoStop    = "nostop"
)

const maxShapes = 3

type Command struct {
	ID     string `json:"id"`
	Phrase string `json:"phrase"`
	Action string `json:"action"`
	Text   string `json:"text"`
}

type Result struct {
	Text      string
	Cancelled bool
	Undo      bool
	Shapes    []Command
	Applied   []string
}

func ValidAction(v string) bool {
	switch v {
	case ActionNewline, ActionParagraph, ActionText, ActionCancel,
		ActionUndo, ActionHeading, ActionList, ActionCapital, ActionQuotes, ActionNoStop:
		return true
	}
	return false
}

func IsShape(v string) bool {
	switch v {
	case ActionHeading, ActionList, ActionCapital, ActionQuotes, ActionNoStop:
		return true
	}
	return false
}

func hasParam(v string) bool {
	switch v {
	case ActionText, ActionHeading, ActionList, ActionQuotes:
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
		if !hasParam(c.Action) {
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
	toks := tokenize(text)
	phrases := make([][]string, len(cmds))
	for i, c := range cmds {
		phrases[i] = words(c.Phrase)
	}
	for i, c := range cmds {
		if c.Action != ActionUndo && c.Action != ActionCancel {
			continue
		}
		if len(phrases[i]) == 0 || len(phrases[i]) != len(toks) || !endsWith(toks, phrases[i]) {
			continue
		}
		res.Applied = append(res.Applied, c.Phrase)
		if c.Action == ActionUndo {
			res.Undo = true
			res.Text = ""
		} else {
			res.Cancelled = true
		}
		return res
	}
	for i, c := range cmds {
		if c.Action == ActionCancel && endsWith(toks, phrases[i]) {
			res.Cancelled = true
			res.Applied = append(res.Applied, c.Phrase)
			return res
		}
	}
	body, shapes, used := cutShapes(cmds, phrases, text, toks)
	res.Text = body
	res.Shapes = shapes
	res.Applied = append(res.Applied, used...)
	for _, c := range cmds {
		if strings.TrimSpace(c.Phrase) == "" || c.Action == ActionUndo || c.Action == ActionCancel || IsShape(c.Action) {
			continue
		}
		rule := replace.Rule{From: c.Phrase, To: output(c), Whole: true}
		before := res.Text
		res.Text = replace.Apply([]replace.Rule{rule}, res.Text)
		if res.Text != before {
			res.Applied = append(res.Applied, c.Phrase)
		}
	}
	res.Text = tidy(res.Text)
	return res
}
