package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// sendAPIRequest отправляет запрос к API и возвращает ответ.
func sendAPIRequest(agent *Agent, config *Config) (*APIResponse, error) {
	apiURL := strings.TrimSuffix(config.APIURL, "/") + "/v1/chat/completions"

	requestBody := APIRequest{
		Model:       config.Model,
		Messages:    agent.Messages,
		Temperature: config.Temp,
		MaxTokens:   config.MaxTokens,
		Stream:      false,
		ToolChoice:  "auto",
		Tools: []Tool{
			{
				Type: "function",
				Function: FunctionDefinition{
					Name:        "execute_command",
					Description: "Execute shell command",
					Parameters: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"command": map[string]string{"type": "string"},
						},
						"required": []string{"command"},
					},
				},
			},
		},
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