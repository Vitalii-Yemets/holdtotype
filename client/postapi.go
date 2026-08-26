package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

func postAPIOn(cfg *Config) bool { return strings.TrimSpace(cfg.PostAPIURL) != "" }

// postReady says whether the prompt chain has anything to run on: the local
// model, or the external server the user set up on purpose.
func postReady(cfg *Config) bool { return postAPIOn(cfg) || llmInstalled(cfg) }

const dpapiNoUI = 1

func protectKey(plain string) (string, error) {
	if plain == "" {
		return "", nil
	}
	in := []byte(plain)
	var out windows.DataBlob
	blob := windows.DataBlob{Size: uint32(len(in)), Data: &in[0]}
	if err := windows.CryptProtectData(&blob, nil, nil, 0, nil, dpapiNoUI, &out); err != nil {
		return "", err
	}
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data)))
	enc := unsafe.Slice(out.Data, out.Size)
	return base64.StdEncoding.EncodeToString(enc), nil
}

func unprotectKey(stored string) (string, error) {
	if stored == "" {
		return "", nil
	}
	raw, err := base64.StdEncoding.DecodeString(stored)
	if err != nil || len(raw) == 0 {
		return "", fmt.Errorf("ключ повреждён")
	}
	var out windows.DataBlob
	blob := windows.DataBlob{Size: uint32(len(raw)), Data: &raw[0]}
	if err := windows.CryptUnprotectData(&blob, nil, nil, 0, nil, dpapiNoUI, &out); err != nil {
		return "", err
	}
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data)))
	return string(unsafe.Slice(out.Data, out.Size)), nil
}

func validPostAPIURL(raw string) bool {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return true
	}
	u, err := url.Parse(raw)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func externalChat(ctx context.Context, cfg *Config, prompt, text string) (string, error) {
	base := strings.TrimRight(strings.TrimSpace(cfg.PostAPIURL), "/")
	timeout := time.Duration(cfg.PostAPITimeout) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	tctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	body, err := json.Marshal(chatRequest{
		Model: strings.TrimSpace(cfg.PostAPIModel),
		Messages: []chatMessage{
			{Role: "system", Content: prompt},
			{Role: "user", Content: text},
		},
		Temperature: 0.2,
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(tctx, http.MethodPost, base+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if key, kerr := unprotectKey(cfg.PostAPIKey); kerr != nil {
		return "", fmt.Errorf("ключ API: %w", kerr)
	} else if key != "" {
		req.Header.Set("Authorization", "Bearer "+key)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	if err != nil {
		return "", err
	}
	var parsed chatResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("неожиданный ответ сервера (%d): %.200s", resp.StatusCode, string(raw))
	}
	if parsed.Error.Message != "" {
		return "", fmt.Errorf("сервер постобработки: %s", parsed.Error.Message)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("сервер постобработки: HTTP %d: %.200s", resp.StatusCode, string(raw))
	}
	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("пустой ответ сервера постобработки")
	}
	log.Printf("постобработка через внешний сервер %s: %d → %d символов", base, len([]rune(text)), len([]rune(parsed.Choices[0].Message.Content)))
	return parsed.Choices[0].Message.Content, nil
}
