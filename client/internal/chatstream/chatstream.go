package chatstream

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type chunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Content string `json:"content"`
	Error   *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (c chunk) text() string {
	if len(c.Choices) == 0 {
		return c.Content
	}
	if t := c.Choices[0].Delta.Content; t != "" {
		return t
	}
	return c.Choices[0].Message.Content
}

func (c chunk) finish() string {
	if len(c.Choices) == 0 {
		return ""
	}
	return c.Choices[0].FinishReason
}

func (c chunk) failure() error {
	if c.Error == nil {
		return nil
	}
	msg := strings.TrimSpace(c.Error.Message)
	if msg == "" {
		msg = "the server reported an error"
	}
	return errors.New(msg)
}

func whole(raw string, onText func(string)) (string, string, error) {
	var c chunk
	if err := json.Unmarshal([]byte(raw), &c); err != nil {
		return "", "", fmt.Errorf("unexpected answer: %.200s", strings.TrimSpace(raw))
	}
	if err := c.failure(); err != nil {
		return "", "", err
	}
	if len(c.Choices) == 0 && c.Content == "" {
		return "", "", fmt.Errorf("empty answer: %.200s", strings.TrimSpace(raw))
	}
	text := c.text()
	if text != "" && onText != nil {
		onText(text)
	}
	return text, c.finish(), nil
}

func Read(r io.Reader, onText func(string)) (string, error) {
	text, _, err := ReadEnd(r, onText)
	return text, err
}

func ReadEnd(r io.Reader, onText func(string)) (string, string, error) {
	br := bufio.NewReader(r)
	var b strings.Builder
	first := true
	end := ""
	for {
		line, err := br.ReadString('\n')
		if first && strings.TrimSpace(line) != "" {
			first = false
			if t := strings.TrimSpace(line); strings.HasPrefix(t, "{") || strings.HasPrefix(t, "[") {
				rest, rerr := io.ReadAll(br)
				if rerr != nil {
					return "", "", rerr
				}
				return whole(line+string(rest), onText)
			}
		}
		if t := strings.TrimSpace(line); strings.HasPrefix(t, "data:") {
			data := strings.TrimSpace(strings.TrimPrefix(t, "data:"))
			if data == "[DONE]" {
				return b.String(), end, nil
			}
			var c chunk
			if json.Unmarshal([]byte(data), &c) == nil {
				if ferr := c.failure(); ferr != nil {
					return b.String(), end, ferr
				}
				if f := c.finish(); f != "" {
					end = f
				}
				if piece := c.text(); piece != "" {
					b.WriteString(piece)
					if onText != nil {
						onText(b.String())
					}
				}
			}
		}
		if err == io.EOF {
			return b.String(), end, nil
		}
		if err != nil {
			return b.String(), end, err
		}
	}
}
