package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func loadConfig() *Config {
	config := &Config{
		Temp:                  DefaultTemp,
		MaxTokens:             DefaultMaxTokens,
		APIURL:                DefaultAPIURL,
		Model:                 DefaultModel,
		RAGEnabled:            false,
		RAGSnippets:           DefaultRAGSnippets,
		AutoCompress:          true,
		AutoCompressThreshold: DefaultAutoCompressThreshold,
		ModelContextLength:    DefaultModelContextLength,
		Stream:                false,
		SubagentsEnabled:      true,
		ExecutionMode:         Ask,
		OperationMode:         Build,
	}
	config.MCPs = make(map[string]MCPServer)

	home, err := os.UserHomeDir()
	if err == nil {
		configPath := filepath.Join(home, ".config", "agent-go", "config.json")
		if _, err := os.Stat(configPath); err == nil {
			file, err := os.ReadFile(configPath)
			if err == nil {
				if err := json.Unmarshal(file, &config); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to parse config file: %v\n", err)
				}
			}
		}
	}

	// Add default context7 MCP if no other MCPs are configured after loading
	if config.MCPs == nil || len(config.MCPs) == 0 {
		if config.MCPs == nil {
			config.MCPs = make(map[string]MCPServer)
		}
		config.MCPs["context7"] = MCPServer{
			Name:    "context7",
			Command: "npx -y @upstash/context7-mcp",
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
	if autoCompress := os.Getenv("AUTO_COMPRESS"); autoCompress == "1" {
		config.AutoCompress = true
	}
	if autoCompressThreshold := os.Getenv("AUTO_COMPRESS_THRESHOLD"); autoCompressThreshold != "" {
		if val, err := strconv.Atoi(autoCompressThreshold); err == nil && val > 0 {
			config.AutoCompressThreshold = val
		}
	}
	if modelContextLength := os.Getenv("MODEL_CONTEXT_LENGTH"); modelContextLength != "" {
		if val, err := strconv.Atoi(modelContextLength); err == nil && val > 0 {
			config.ModelContextLength = val
		}
	}
	if stream := os.Getenv("STREAM_ENABLED"); stream == "1" || stream == "true" {
		config.Stream = true
	}
	if subagents := os.Getenv("SUBAGENTS_ENABLED"); subagents == "0" || subagents == "false" {
		config.SubagentsEnabled = false
	}
	if executionMode := os.Getenv("EXECUTION_MODE"); executionMode != "" {
		if executionMode == "yolo" {
			config.ExecutionMode = YOLO
		} else {
			config.ExecutionMode = Ask
		}
	}
	if operationMode := os.Getenv("OPERATION_MODE"); operationMode != "" {
		if operationMode == "plan" {
			config.OperationMode = Plan
		} else {
			config.OperationMode = Build
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
