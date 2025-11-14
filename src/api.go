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
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}()

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
		return "", fmt.Errorf("not enough messages to compress")
	}

	var messagesToCompress []Message
	for _, msg := range agent.Messages {
		if msg.Role != "system" {
			messagesToCompress = append(messagesToCompress, msg)
		}
	}

	// Build the prompt for compression
	var compressionBuilder strings.Builder
	compressionBuilder.WriteString("Compress the following conversation into a brief summary (1-3 sentences), preserving key details and context:\n\n")
	for _, msg := range messagesToCompress {
		switch msg.Role {
		case "user":
			compressionBuilder.WriteString(fmt.Sprintf("User: %s\n", *msg.Content))
		case "assistant":
			compressionBuilder.WriteString(fmt.Sprintf("Assistant: %s\n", *msg.Content))
		}
	}
	compressionBuilder.WriteString("\nBrief summary:")
	compressionPrompt := compressionBuilder.String()

	// Create the compression request
	requestBody := APIRequest{
		Model:       config.Model,
		Messages:    []Message{{Role: "user", Content: &compressionPrompt}},
		Temperature: CompressionTemp,      // Low temperature for more precise compression
		MaxTokens:   CompressionMaxTokens, // Limit the length of the compressed text
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
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}()

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
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Process the stream
	scanner := bufio.NewScanner(resp.Body)
	var fullContent strings.Builder
	var toolCalls []ToolCall
	finalRole := "assistant"
	var finalUsage Usage

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue // Skip empty lines
		}
		if !strings.HasPrefix(line, StreamDataPrefix) {
			continue // Skip non-data lines
		}

		data := strings.TrimPrefix(line, StreamDataPrefix)
		if data == StreamDoneMarker {
			break // End of stream
		}

		var chunk StreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			// Silently skip invalid JSON chunks
			continue
		}
		
		if len(chunk.Choices) == 0 {
			continue
		}

		delta := chunk.Choices[0].Delta
		if delta.Role != "" {
			finalRole = delta.Role
		}
		if delta.Content != "" {
			fullContent.WriteString(delta.Content)
			fmt.Printf("%s%s%s", ColorBlue, delta.Content, ColorReset)
		}
		if len(delta.ToolCalls) > 0 {
			// In streaming, tool calls can be sent incrementally.
			// We need to merge them based on their index.
			for _, toolCallChunk := range delta.ToolCalls {
				if toolCallChunk.Index >= len(toolCalls) {
					toolCalls = append(toolCalls, ToolCall{
						ID:       toolCallChunk.ID,
						Type:     toolCallChunk.Type,
						Function: FunctionCall{},
					})
				}
				if toolCallChunk.Function.Name != "" {
					toolCalls[toolCallChunk.Index].Function.Name = toolCallChunk.Function.Name
				}
				if toolCallChunk.Function.Arguments != "" {
					toolCalls[toolCallChunk.Index].Function.Arguments += toolCallChunk.Function.Arguments
				}
			}
		}

		// Check for usage data in the stream (some providers send it)
		if chunk.Usage != nil {
			finalUsage = *chunk.Usage
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
	finalContentStr := fullContent.String()
	response := &APIResponse{
		Choices: []Choice{
			{
				Message: Message{
					Role:      finalRole,
					Content:   &finalContentStr,
					ToolCalls: toolCalls,
				},
			},
		},
		Usage: finalUsage,
	}

	return response, nil
}
