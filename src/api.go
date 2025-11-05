package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

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
