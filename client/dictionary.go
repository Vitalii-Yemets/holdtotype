package main

import (
	"log"
	"strings"
)

var builtinDictionaries = map[string]string{
	"en": "GitHub, GitLab, Docker, Kubernetes, Python, TypeScript, JavaScript, React, Node.js, Postgres, Redis, nginx, API, JSON, YAML, SQL, CLI, UI, UX, backend, frontend, DevOps, CI/CD, pull request, merge, commit, rebase, deploy, rollback, repository, endpoint, webhook, latency, cache, ChatGPT, Claude, Gemini, Copilot, Hugging Face, Cursor, prompt, token, embedding, stream, subscriber, donation, emote, clip, shorts, thumbnail",
	"ru": "GitHub, Docker, Kubernetes, Python, TypeScript, JavaScript, React, Postgres, Redis, nginx, API, JSON, YAML, SQL, бэкенд, фронтенд, деплой, коммит, ребейз, пул-реквест, мёрдж, откат, репозиторий, эндпоинт, вебхук, кэш, ChatGPT, Claude, Gemini, Copilot, Hugging Face, Cursor, промпт, токен, эмбеддинг, нейросеть, стрим, донат, подписчик, эмоут, клип, шортс, превью, хоткей, буфер обмена, распознавание, диктовка",
	"uk": "GitHub, Docker, Kubernetes, Python, TypeScript, JavaScript, React, Postgres, Redis, nginx, API, JSON, YAML, SQL, бекенд, фронтенд, деплой, коміт, ребейз, пул-реквест, мердж, відкат, репозиторій, ендпоінт, вебхук, кеш, ChatGPT, Claude, Gemini, Copilot, Hugging Face, Cursor, промпт, токен, ембединг, нейромережа, стрім, донат, підписник, емоут, кліп, шортс, прев'ю, гарячі клавіші, буфер обміну, розпізнавання, диктування",
	"de": "GitHub, Docker, Kubernetes, Python, TypeScript, JavaScript, React, Postgres, Redis, nginx, API, JSON, YAML, SQL, Backend, Frontend, Deployment, Commit, Rebase, Pull Request, Merge, Rollback, Repository, Endpunkt, Webhook, Cache, ChatGPT, Claude, Gemini, Copilot, Hugging Face, Cursor, Prompt, Token, Embedding, KI, Stream, Spende, Abonnent, Emote, Clip, Shorts, Tastenkürzel, Zwischenablage, Spracherkennung, Diktat",
	"fr": "GitHub, Docker, Kubernetes, Python, TypeScript, JavaScript, React, Postgres, Redis, nginx, API, JSON, YAML, SQL, backend, frontend, déploiement, commit, rebase, pull request, merge, rollback, dépôt, point de terminaison, webhook, cache, ChatGPT, Claude, Gemini, Copilot, Hugging Face, Cursor, prompt, token, embedding, IA, stream, don, abonné, emote, clip, shorts, raccourci clavier, presse-papiers, reconnaissance vocale, dictée",
	"es": "GitHub, Docker, Kubernetes, Python, TypeScript, JavaScript, React, Postgres, Redis, nginx, API, JSON, YAML, SQL, backend, frontend, despliegue, commit, rebase, pull request, merge, rollback, repositorio, endpoint, webhook, caché, ChatGPT, Claude, Gemini, Copilot, Hugging Face, Cursor, prompt, token, embedding, IA, directo, donación, suscriptor, emote, clip, shorts, atajo de teclado, portapapeles, reconocimiento de voz, dictado",
	"it": "GitHub, Docker, Kubernetes, Python, TypeScript, JavaScript, React, Postgres, Redis, nginx, API, JSON, YAML, SQL, backend, frontend, deploy, commit, rebase, pull request, merge, rollback, repository, endpoint, webhook, cache, ChatGPT, Claude, Gemini, Copilot, Hugging Face, Cursor, prompt, token, embedding, IA, diretta, donazione, iscritto, emote, clip, shorts, scorciatoia, appunti, riconoscimento vocale, dettatura",
	"pl": "GitHub, Docker, Kubernetes, Python, TypeScript, JavaScript, React, Postgres, Redis, nginx, API, JSON, YAML, SQL, backend, frontend, wdrożenie, commit, rebase, pull request, merge, rollback, repozytorium, endpoint, webhook, cache, ChatGPT, Claude, Gemini, Copilot, Hugging Face, Cursor, prompt, token, embedding, SI, stream, donacja, subskrybent, emotka, klip, shorts, skrót klawiszowy, schowek, rozpoznawanie mowy, dyktowanie",
}

func builtinDictionary(lang string) string {
	if d, ok := builtinDictionaries[lang]; ok {
		return d
	}
	return builtinDictionaries["en"]
}

func isBuiltinDictionary(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return true
	}
	for _, d := range builtinDictionaries {
		if s == d {
			return true
		}
	}
	return s == legacyDictionary
}

func syncDictionary(cfg *Config) bool {
	if !isBuiltinDictionary(cfg.WhisperPrompt) {
		return false
	}
	want := builtinDictionary(cfg.Language)
	if cfg.WhisperPrompt == want {
		return false
	}
	cfg.WhisperPrompt = want
	log.Printf("recognition dictionary replaced with the set for language %s", cfg.Language)
	return true
}

var legacyDictionary = "Whisper, whisper.cpp, Docker, Go, LLM, llama.cpp, Qwen, UI, UX, API, HTTP, JSON, GitHub, Windows, exe, ggml, промпт, хоткей, чекбокс, радиокнопка, таймаут, конфиг, вкладка, секция, диктовка, распознавание, постобработка, перевод, интерфейс, локализация, сочетание клавиш, буфер обмена, курсор, микрофон, модель, сервер, трей, оверлей, плашка, скролл, ползунок, диалог, кнопка, профиль, hotkey, checkbox, timeout, dictation, transcription, translation, clipboard, cursor, microphone, overlay, slider, settings, диктування, розпізнавання, налаштування, буфер обміну, Einstellungen, Tastenkürzel, Zwischenablage, Übersetzung, paramètres, raccourci clavier, presse-papiers, traduction, ajustes, atajo de teclado, portapapeles, traducción, impostazioni, scorciatoia, appunti, traduzione, ustawienia, skrót klawiszowy, schowek, tłumaczenie"
