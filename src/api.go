package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func sendAPIRequest(agent *Agent, config *Config, includeSpawn bool, agentDef *AgentDefinition) (*APIResponse, error) {
	apiURL := strings.TrimSuffix(config.APIURL, "/") + "/v1/chat/completions"

	// Build base tools (now includes all tools)
	baseTools := getAvailableTools(config, includeSpawn, config.OperationMode)

	// Apply operation mode filtering and agent-specific policy
	tools := filterToolsByPolicy(baseTools, agentDef, config.OperationMode)

	// Create a copy of messages with time context injected
	messagesWithTime := make([]Message, len(agent.Messages))
	copy(messagesWithTime, agent.Messages)

	// Inject current time as a system message at the end
	timeContext := getCurrentTimeContext()
	timeMsg := Message{
		Role:    "system",
		Content: &timeContext,
	}
	messagesWithTime = append(messagesWithTime, timeMsg)

	requestBody := APIRequest{
		Model:       config.Model,
		Messages:    messagesWithTime,
		Temperature: config.Temp,
		MaxTokens:   config.MaxTokens,
		ToolChoice:  "auto",
		Tools:       tools,
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

// sendMiniLLMRequest sends a request to the configured mini model (or fallback to main model)
func sendMiniLLMRequest(config *Config, messages []Message) (string, error) {
	model := config.MiniModel
	if model == "" {
		model = config.Model
	}

	// Create a copy of messages with time context injected
	messagesWithTime := make([]Message, len(messages))
	copy(messagesWithTime, messages)

	// Inject current time as a system message at the end
	timeContext := getCurrentTimeContext()
	timeMsg := Message{
		Role:    "system",
		Content: &timeContext,
	}
	messagesWithTime = append(messagesWithTime, timeMsg)

	requestBody := APIRequest{
		Model:       model,
		Messages:    messagesWithTime,
		Temperature: CompressionTemp,      // Reuse default temp for utility tasks
		MaxTokens:   CompressionMaxTokens, // Reuse default max tokens for utility tasks
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	apiURL := strings.TrimSuffix(config.APIURL, "/") + "/v1/chat/completions"
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var apiResponse APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(apiResponse.Choices) == 0 {
		return "", fmt.Errorf("received empty response from API")
	}

	content := apiResponse.Choices[0].Message.Content
	if content == nil {
		return "", fmt.Errorf("received empty content from API")
	}

	return *content, nil
}

// separateMessages separates system messages from conversation messages (user/assistant/tool)
func separateMessages(messages []Message) ([]Message, []Message) {
	systemMsgs := []Message{}
	conversationMsgs := []Message{}

	for _, msg := range messages {
		if msg.Role == "system" {
			systemMsgs = append(systemMsgs, msg)
		} else {
			conversationMsgs = append(conversationMsgs, msg)
		}
	}

	return systemMsgs, conversationMsgs
}

func compressContext(agent *Agent, config *Config) (string, error) {
	if len(agent.Messages) <= 1 {
		return "", fmt.Errorf("not enough messages to compress")
	}

	// Separate system messages from conversation messages
	systemMsgs, conversationMsgs := separateMessages(agent.Messages)

	// Log message breakdown
	fmt.Printf("%sCompressing context: %d total messages (%d system, %d conversation)%s\n",
		ColorMeta, len(agent.Messages), len(systemMsgs), len(conversationMsgs), ColorReset)

	if len(conversationMsgs) == 0 {
		return "", fmt.Errorf("no conversation messages to compress (only system messages)")
	}

	messagesToCompress := conversationMsgs

	// Build the prompt for compression
	var compressionBuilder strings.Builder
	compressionBuilder.WriteString("Compress the following conversation into a brief summary (1-3 sentences), preserving key details and context:\n\n")
	for _, msg := range messagesToCompress {
		if msg.Content == nil {
			continue // Skip messages with no content (e.g. tool calls only)
		}
		switch msg.Role {
		case "user":
			compressionBuilder.WriteString(fmt.Sprintf("User: %s\n", *msg.Content))
		case "assistant":
			compressionBuilder.WriteString(fmt.Sprintf("Assistant: %s\n", *msg.Content))
		}
	}
	compressionBuilder.WriteString("\nBrief summary:")
	compressionPrompt := compressionBuilder.String()

	// Use the MAIN model for compression to ensure accuracy and handle larger contexts if needed.
	// Mini models might have smaller context windows or lower reasoning capabilities for complex summaries.
	requestBody := APIRequest{
		Model:       config.Model,
		Messages:    []Message{{Role: "user", Content: &compressionPrompt}},
		Temperature: CompressionTemp,      // Low temperature for more precise compression
		MaxTokens:   CompressionMaxTokens, // Limit the length of the compressed text
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
