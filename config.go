package main

import (
	"os"
	"strconv"
)

// loadConfig загружает конфигурацию из переменных окружения.
func loadConfig() *Config {
	config := &Config{
		Temp:        0.1,
		MaxTokens:   1000,
		APIURL:      "https://api.openai.com",
		Model:       "gpt-3.5-turbo",
		RAGEnabled:  false,
		RAGSnippets: 5,
	}

	if key := os.Getenv("OPENAI_KEY"); key != "" {
		config.APIKey = key
	}
	if url := os.Getenv("OPENAI_BASE"); url != "" {
		config.APIURL = url
	}
	if model := os.Getenv("OPENAI_MODEL"); model != "" {
		config.Model = model
	}
	if ragPath := os.Getenv("RAG_PATH"); ragPath != "" {
		config.RAGPath = ragPath
	}
	if ragEnabled := os.Getenv("RAG_ENABLED"); ragEnabled == "1" {
		config.RAGEnabled = true
	}
	if ragSnippets := os.Getenv("RAG_SNIPPETS"); ragSnippets != "" {
		if val, err := strconv.Atoi(ragSnippets); err == nil && val > 0 {
			config.RAGSnippets = val
		}
	}

	return config
}