package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode"

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

// AgentCompleter handles both slash commands and @filename completion
type AgentCompleter struct {
	slashCompleter *readline.PrefixCompleter
}

func (c *AgentCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	// Check if we are in a slash command context (line starts with /)
	// But wait, we might want to autocomplete files even inside a slash command (e.g. /somecmd @file)
	// However, readline's PrefixCompleter expects to handle the whole line if it matches.
	
	// Let's check the word being typed.
	word, _ := getLastWord(line, pos)
	
	if strings.HasPrefix(word, "@") {
		// File completion logic
		prefix := word[1:] // remove @
		matches, _ := filepath.Glob(prefix + "*")
		
		// If no matches, maybe it's a directory inside?
		// Simple implementation: list current directory and filter
		if matches == nil {
			// Try listing all files and filtering
			files, err := os.ReadDir(".")
			if err == nil {
				for _, f := range files {
					name := f.Name()
					if strings.HasPrefix(name, prefix) {
						matches = append(matches, name)
					}
				}
			}
		}

		// We want to return candidates without the @ if we are appending, 
		// but readline expects the replacement.
		// Actually, readline's Do returns valid suffixes usually?
		// No, it returns the full string to replace the input "length" characters back?
		// Let's look at standard behavior.
		// If I type "@fil", I expect suggestions like "filename.txt".
		// If I select it, it becomes "@filename.txt".
		
		var candidates [][]rune
		for _, m := range matches {
			// We suggest the full filename
			// The 'length' return value tells readline how many chars before cursor to replace.
			// We want to replace 'prefix' (the part after @).
			// So length = len(prefix).
			// And the candidate should be the rest of the filename?
			// Or the full filename?
			// If length is len(prefix), then we replace 'fil' with 'filename.txt'.
			// So the result is '@filename.txt'.
			candidates = append(candidates, []rune(m))
		}
		return candidates, len(prefix)
	}

	// Delegate to slash completer if line starts with /
	if len(line) > 0 && line[0] == '/' {
		return c.slashCompleter.Do(line, pos)
	}

	return nil, 0
}

func getLastWord(line []rune, pos int) (string, int) {
	if pos == 0 {
		return "", 0
	}
	// Scan backwards for whitespace
	start := pos
	for start > 0 && !unicode.IsSpace(line[start-1]) {
		start--
	}
	return string(line[start:pos]), start
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

	// Prepare MCP server completions for the /mcp remove command
	mcpServerCompleters := make([]readline.PrefixCompleterInterface, 0)
	if config.MCPs != nil {
		for name := range config.MCPs {
			mcpServerCompleters = append(mcpServerCompleters, readline.PcItem(name))
		}
	}

	var slashCompleter = readline.NewPrefixCompleter(
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

	return &AgentCompleter{slashCompleter: slashCompleter}
}