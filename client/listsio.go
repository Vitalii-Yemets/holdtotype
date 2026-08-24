package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"holdtotype/internal/appid"
	"holdtotype/internal/commands"
	"holdtotype/internal/lists"
	"holdtotype/internal/replace"
)

type listsPayload struct {
	Replacements []replace.Rule     `json:"replacements"`
	Commands     []commands.Command `json:"commands"`
}

type listsReply struct {
	OK           bool               `json:"ok"`
	Cancelled    bool               `json:"cancelled"`
	Text         string             `json:"text"`
	Replacements []replace.Rule     `json:"replacements,omitempty"`
	Commands     []commands.Command `json:"commands,omitempty"`
}

func listsAnswer(r listsReply) string {
	out, _ := json.Marshal(r)
	return string(out)
}

func exportLists(payload string) string {
	var in listsPayload
	if err := json.Unmarshal([]byte(payload), &in); err != nil {
		log.Printf("списки: разбор страницы: %v", err)
		return listsAnswer(listsReply{Text: tr("lists.bad")})
	}
	data, err := lists.Encode(in.Replacements, in.Commands)
	if err != nil {
		log.Printf("списки: сборка файла: %v", err)
		return listsAnswer(listsReply{Text: tr("lists.bad")})
	}
	path := askFilePath(true, tr("lists.save.title"), appid.Slug+"-lists.json")
	if path == "" {
		return listsAnswer(listsReply{Cancelled: true})
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		log.Printf("списки: запись %s: %v", path, err)
		return listsAnswer(listsReply{Text: err.Error()})
	}
	log.Printf("списки сохранены: %s (замен %d, команд %d)", path, len(in.Replacements), len(in.Commands))
	return listsAnswer(listsReply{OK: true, Text: trf("lists.saved", filepath.Base(path))})
}

func importLists(payload string) string {
	var in listsPayload
	if err := json.Unmarshal([]byte(payload), &in); err != nil {
		log.Printf("списки: разбор страницы: %v", err)
		return listsAnswer(listsReply{Text: tr("lists.bad")})
	}
	path := askFilePath(false, tr("lists.open.title"), "")
	if path == "" {
		return listsAnswer(listsReply{Cancelled: true})
	}
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("списки: чтение %s: %v", path, err)
		return listsAnswer(listsReply{Text: err.Error()})
	}
	f, err := lists.Parse(data)
	if err != nil {
		log.Printf("списки: %s не наш файл", path)
		return listsAnswer(listsReply{Text: tr("lists.bad")})
	}
	rules, addedR, skipR := lists.MergeRules(in.Replacements, f.Replacements)
	cmds, addedC, skipC := lists.MergeCommands(in.Commands, f.Commands)
	log.Printf("списки загружены из %s: добавлено замен %d, команд %d, пропущено %d",
		path, addedR, addedC, skipR+skipC)
	return listsAnswer(listsReply{
		OK:           true,
		Text:         trf("lists.added", addedR+addedC, skipR+skipC),
		Replacements: rules,
		Commands:     cmds,
	})
}

func (a *App) insertTarget() uintptr {
	own := settingsHwnd.Load()
	a.mu.Lock()
	prev, last := a.settingsPrev, a.lastWnd
	a.mu.Unlock()
	for _, wnd := range []uintptr{prev, last} {
		if wnd == 0 || wnd == own {
			continue
		}
		if vis, _, _ := procIsWindowVisible.Call(wnd); vis != 0 {
			return wnd
		}
	}
	return 0
}

func (a *App) insertFromHistory(text string) string {
	if strings.TrimSpace(text) == "" {
		return listsAnswer(listsReply{Text: tr("hist.insert.gone")})
	}
	wnd := a.insertTarget()
	if wnd == 0 {
		if err := setClipboardText(text); err != nil {
			log.Printf("история: копирование: %v", err)
			return listsAnswer(listsReply{Text: err.Error()})
		}
		log.Printf("история: некуда вставлять, текст скопирован")
		return listsAnswer(listsReply{Text: tr("hist.insert.nowin")})
	}
	procSetForegroundWnd.Call(wnd)
	time.Sleep(250 * time.Millisecond)
	if cur, _, _ := procGetForegroundWindow.Call(); cur != wnd {
		_ = setClipboardText(text)
		log.Printf("история: окно не вышло вперёд, текст скопирован")
		return listsAnswer(listsReply{Text: tr("hist.insert.nowin")})
	}
	cfg := a.snapshot()
	if err := pasteText(cfg, text, wnd); err != nil {
		log.Printf("история: вставка: %v", err)
		return listsAnswer(listsReply{Text: err.Error()})
	}
	title := windowTitle(wnd)
	log.Printf("история: вставлено %d символов в «%s»", len([]rune(text)), title)
	return listsAnswer(listsReply{OK: true, Text: trf("hist.insert.ok", title)})
}

func (a *App) copyLastResult() (bool, string) {
	a.mu.Lock()
	text := a.lastResult
	a.mu.Unlock()
	if strings.TrimSpace(text) == "" {
		return false, tr("copy.none")
	}
	if err := setClipboardText(text); err != nil {
		log.Printf("копирование последнего результата: %v", err)
		return false, trf("copy.fail", err.Error())
	}
	log.Printf("последний результат скопирован: %d символов", len([]rune(text)))
	return true, tr("copy.ok")
}
