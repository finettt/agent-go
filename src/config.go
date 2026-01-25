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
		MiniModel:             DefaultMiniModel,
		RAGEnabled:            false,
		RAGSnippets:           DefaultRAGSnippets,
		AutoCompress:          true,
		AutoCompressThreshold: DefaultAutoCompressThreshold,
		ModelContextLength:    DefaultModelContextLength,
		SubagentsEnabled:      true,
		ExecutionMode:         Ask,
		OperationMode:         Build,
		UsageVerboseMode:      UsageSilent,
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

	// MIGRATION: If OperationMode is set, warn user
	if config.OperationMode == Plan {
		fmt.Println("Notice: operation_mode='plan' is deprecated. The 'plan' agent will be activated instead.")
		fmt.Println("  Update your config to remove 'operation_mode'. Use /plan to switch between plan/build agents.")
	} else if config.OperationMode == Build {
		// Build is default, only warn if explicitly set in env var
		if _, exists := os.LookupEnv("OPERATION_MODE"); exists {
			fmt.Println("Notice: operation_mode config is deprecated. Use /plan to switch between plan/build agents.")
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

	// Load skills
	skills, err := loadSkills()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load skills: %v\n", err)
	} else {
		config.Skills = skills
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
	if miniModel := os.Getenv("OPENAI_MINI_MODEL"); miniModel != "" {
		config.MiniModel = miniModel
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
