package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func sendAPIRequest(agent *Agent, config *Config, includeSpawn bool) (*APIResponse, error) {
	apiURL := strings.TrimSuffix(config.APIURL, "/") + "/v1/chat/completions"

	requestBody := APIRequest{
		Model:       config.Model,
		Messages:    agent.Messages,
		Temperature: config.Temp,
		MaxTokens:   config.MaxTokens,
		Stream:      config.Stream,
		ToolChoice:  "auto",
		Tools:       getAvailableTools(includeSpawn),
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var apiResponse APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &apiResponse, nil
}

func getAvailableTools(includeSpawn bool) []Tool {
	tools := []Tool{
		{
			Type: "function",
			Function: FunctionDefinition{
				Name:        "execute_command",
				Description: "Execute shell command",
				Parameters: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{"command": map[string]string{"type": "string"}},
					"required":   []string{"command"},
				},
			},
		},
	}

	if includeSpawn {
		tools = append(tools, Tool{
			Type: "function",
			Function: FunctionDefinition{
				Name:        "spawn_agent",
				Description: "Spawn a sub-agent to perform a specific task and return the result.",
				Parameters: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{"task": map[string]string{"type": "string"}},
					"required":   []string{"task"},
				},
			},
		})
	}

	return tools
}

func compressContext(agent *Agent, config *Config) (string, error) {
	if len(agent.Messages) <= 1 {
		return "", fmt.Errorf("no messages to compress")
	}

	var messagesToCompress []Message
	for _, msg := range agent.Messages {
		if msg.Role != "system" {
			messagesToCompress = append(messagesToCompress, msg)
		}
	}

	// Формируем промпт для сжатия
	compressionPrompt := "Сжати следующую беседу в краткое резюме (1-3 предложения), сохраняя ключевые детали и контекст:\n\n"
	for _, msg := range messagesToCompress {
		if msg.Role == "user" {
			compressionPrompt += fmt.Sprintf("Пользователь: %s\n", *msg.Content)
		} else if msg.Role == "assistant" {
			compressionPrompt += fmt.Sprintf("Ассистент: %s\n", *msg.Content)
		}
	}
	compressionPrompt += "\nКраткое резюме:"

	// Создаем запрос для сжатия
	requestBody := APIRequest{
		Model:       config.Model,
		Messages:    []Message{{Role: "user", Content: &compressionPrompt}},
		Temperature: 0.3, // Низкая температура для более точного сжатия
		MaxTokens:   500, // Ограничение на длину сжатого текста
		Stream:      false,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal compression request: %w", err)
	}

	apiURL := strings.TrimSuffix(config.APIURL, "/") + "/v1/chat/completions"
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create compression request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send compression request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("compression API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var apiResponse APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return "", fmt.Errorf("failed to decode compression response: %w", err)
	}

	if len(apiResponse.Choices) == 0 {
		return "", fmt.Errorf("received empty response from compression API")
	}

	compressedContent := apiResponse.Choices[0].Message.Content
	if compressedContent == nil {
		return "", fmt.Errorf("received empty content from compression API")
	}

	return *compressedContent, nil
}

// sendAPIRequestStreaming handles streaming API responses
func sendAPIRequestStreaming(agent *Agent, config *Config, includeSpawn bool) (*APIResponse, error) {
	apiURL := strings.TrimSuffix(config.APIURL, "/") + "/v1/chat/completions"

	requestBody := APIRequest{
		Model:       config.Model,
		Messages:    agent.Messages,
		Temperature: config.Temp,
		MaxTokens:   config.MaxTokens,
		Stream:      true,
		ToolChoice:  "auto",
		Tools:       getAvailableTools(includeSpawn),
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("Accept", "text/event-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Process the stream
	scanner := bufio.NewScanner(resp.Body)
	var fullContent strings.Builder
	var toolCalls []ToolCall
	var role string = "assistant"
	var promptTokens, completionTokens, totalTokens int

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check for data lines
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")

		// Check for stream end
		if data == "[DONE]" {
			break
		}

		// Parse the JSON chunk
		var chunk StreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			// Skip invalid JSON chunks
			continue
		}

		// Process each choice
		for _, choice := range chunk.Choices {
			// Handle content delta
			if choice.Delta.Content != nil {
				content := *choice.Delta.Content
				fullContent.WriteString(content)
				// Print the content as it arrives with blue color
				fmt.Printf("\033[34m%s\033[0m", content)
			}

			// Handle tool calls
			if len(choice.Delta.ToolCalls) > 0 {
				// Accumulate tool calls
				for _, tc := range choice.Delta.ToolCalls {
					toolCalls = append(toolCalls, tc)
				}
			}

			// Handle role if present
			if choice.Delta.Role != nil {
				role = *choice.Delta.Role
			}
		}
	}

	if scanner.Err() != nil {
		return nil, fmt.Errorf("error reading stream: %w", scanner.Err())
	}

	// Print a newline after streaming content
	if fullContent.Len() > 0 {
		fmt.Println()
	}

	// Construct the final response
	finalContent := fullContent.String()
	response := &APIResponse{
		Choices: []Choice{
			{
				Message: Message{
					Role:      role,
					Content:   &finalContent,
					ToolCalls: toolCalls,
				},
			},
		},
		Usage: Usage{
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      totalTokens,
		},
	}

	return response, nil
}
