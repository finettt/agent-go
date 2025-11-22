package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/chzyer/readline"
)

// ModelListResponse represents the expected JSON structure from an OpenAI-compatible /models endpoint.
type ModelListResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

// fetchAvailableModels retrieves the list of available models from the provider's API.
func fetchAvailableModels(config *Config) ([]string, error) {
	if config.APIURL == "" || config.APIKey == "" {
		return nil, fmt.Errorf("API URL or API Key not configured")
	}

	// Construct the request for an OpenAI-compatible models endpoint.
	url := strings.TrimSuffix(config.APIURL, "/") + "/v1/models"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+config.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send request to %s: %w", url, err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch models, status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}

	var modelList ModelListResponse
	if err := json.Unmarshal(body, &modelList); err != nil {
		return nil, fmt.Errorf("could not unmarshal models response: %w", err)
	}

	var models []string
	for _, model := range modelList.Data {
		models = append(models, model.ID)
	}
	return models, nil
}

// buildCompleter creates the completer for readline.
// It fetches model names once on startup for autocompletion.
func buildCompleter(config *Config) readline.AutoCompleter {
	// Prepare model completions
	modelCompleters := make([]readline.PrefixCompleterInterface, 0)
	models, err := fetchAvailableModels(config)
	if err == nil {
		for _, modelName := range models {
			modelCompleters = append(modelCompleters, readline.PcItem(modelName))
		}
	} else {
		// Log error but don't crash, model completion will just be unavailable
		fmt.Printf("Warning: could not fetch models for autocompletion: %v\n", err)
	}

	// Create the main completer with static and semi-static items.
	// Dynamic file path completion is not implemented to maintain compatibility
	// with readline v1.5.1 without a custom completer implementation.

	// Prepare MCP server completions for the /mcp remove command
	mcpServerCompleters := make([]readline.PrefixCompleterInterface, 0)
	if config.MCPs != nil {
		for name := range config.MCPs {
			mcpServerCompleters = append(mcpServerCompleters, readline.PcItem(name))
		}
	}

	var completer = readline.NewPrefixCompleter(
		readline.PcItem("/help"),
		readline.PcItem("/model", modelCompleters...),
		readline.PcItem("/provider",
			readline.PcItem("https://"),
			readline.PcItem("http://"),
		),
		readline.PcItem("/config"),
		readline.PcItem("/rag",
			readline.PcItem("on"),
			readline.PcItem("off"),
			readline.PcItem("path"), // Only suggests the "path" keyword.
		),
		readline.PcItem("/mcp",
			readline.PcItem("add"),
			readline.PcItem("remove", mcpServerCompleters...),
			readline.PcItem("list"),
		),
		readline.PcItem("/compress"),
		readline.PcItem("/clear"),
		readline.PcItem("/contextlength"),
		readline.PcItem("/stream",
			readline.PcItem("on"),
			readline.PcItem("off"),
		),
		readline.PcItem("/subagents",
			readline.PcItem("on"),
			readline.PcItem("off"),
		),
		readline.PcItem("/shell"),
		readline.PcItem("/quit"),
	)

	return completer
}
