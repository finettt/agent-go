package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
)

func loadConfig() *Config {

	config := &Config{
		Temp:        0.1,
		MaxTokens:   1000,
		APIURL:      "https://api.openai.com",
		Model:       "gpt-3.5-turbo",
		RAGEnabled:  false,
		RAGSnippets: 5,
	}

	home, err := os.UserHomeDir()
	if err == nil {
		configPath := filepath.Join(home, ".config", "agent-go", "config.json")
		if _, err := os.Stat(configPath); err == nil {
			file, err := os.ReadFile(configPath)
			if err == nil {
				json.Unmarshal(file, config)
			}
		}
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

func saveConfig(config *Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configDir := filepath.Join(home, ".config", "agent-go")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	configPath := filepath.Join(configDir, "config.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}
