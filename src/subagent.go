package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// runSubAgent executes a task in a separate sub-agent context (backward-compatible default behavior).
func runSubAgent(task string, config *Config) (string, error) {
	return runSubAgentWithAgent(task, "", config)
}

// runSubAgentWithAgent executes a task in a separate sub-agent context, optionally using a saved task-specific agent
// definition (including the built-in "default") as additional system instructions.
func runSubAgentWithAgent(task string, agentName string, config *Config) (string, error) {
	sysInfo := getSystemInfo()

	basePrompt := "You are a sub-agent tasked with completing a specific goal. You have access to the 'execute_command' and todo list management tools. Plan your steps and execute them sequentially. When you have finished the task, output the final result as a single response."

	systemPrompt := sysInfo + "\n\n" + basePrompt

	if strings.TrimSpace(agentName) != "" {
		def, err := loadAgentDefinition(agentName)
		if err != nil {
			return "", fmt.Errorf("failed to load agent '%s' for sub-agent: %w", agentName, err)
		}
		systemPrompt = sysInfo + "\n\n" + fmt.Sprintf("=== Task-Specific Agent: %s ===\n%s\n\n%s", def.Name, def.SystemPrompt, basePrompt)
	}

	subAgent := &Agent{
		ID: uuid.New().String(),
		Messages: []Message{
			{
				Role:    "system",
				Content: &systemPrompt,
			},
			{
				Role:    "user",
				Content: &task,
			},
		},
	}

	// Limit iterations to prevent infinite loops
	for iteration := 0; iteration < MaxSubAgentIterations; iteration++ {
		// Use false for includeSpawn to prevent sub-agents from creating more sub-agents
		resp, err := sendAPIRequest(subAgent, config, false)
		if err != nil {
			return "", fmt.Errorf("sub-agent API request failed: %w", err)
		}

		if len(resp.Choices) == 0 {
			return "", fmt.Errorf("sub-agent received an empty response from the API")
		}

		assistantMsg := resp.Choices[0].Message
		subAgent.Messages = append(subAgent.Messages, assistantMsg)

		if len(assistantMsg.ToolCalls) == 0 {
			// If there are no more tool calls, the sub-agent's work is done
			if assistantMsg.Content != nil {
				return *assistantMsg.Content, nil
			}
			return "", fmt.Errorf("sub-agent finished without providing a result")
		}

		for _, toolCall := range assistantMsg.ToolCalls {
			var output string
			var err error
			var logMessage string

			switch toolCall.Function.Name {
			case "execute_command":
				var args CommandArgs
				if unmarshalErr := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); unmarshalErr != nil {
					output = fmt.Sprintf("Failed to parse arguments: %s", unmarshalErr)
				} else {
					output, err = confirmAndExecute(config, args.Command)
					if err == nil {
						logMessage = fmt.Sprintf("Bash %s(%s)%s", ColorMeta, args.Command, ColorReset)
					}
				}
			case "create_todo":
				output, err = createTodo(subAgent.ID, toolCall.Function.Arguments)
				if err == nil {
					logMessage = "Created todo item"
				}
			case "update_todo":
				output, err = updateTodo(subAgent.ID, toolCall.Function.Arguments)
				if err == nil {
					logMessage = "Updated todo item"
				}
			case "get_todo_list":
				output, err = getTodoList(subAgent.ID)
				if err == nil {
					logMessage = "Retrieved todo list"
				}
			default:
				output = fmt.Sprintf("Unknown tool: %s", toolCall.Function.Name)
			}

			if err != nil {
				output = fmt.Sprintf("Tool execution error: %s", err)
				fmt.Printf("%s%s%s\n", ColorRed, output, ColorReset)
			} else if logMessage != "" {
				fmt.Printf("%s%s==> %s%s%s\n", StyleBold, ColorHighlight, ColorReset, logMessage, ColorReset)
			}

			toolMsg := Message{
				Role:       "tool",
				ToolCallID: toolCall.ID,
				Content:    &output,
			}
			subAgent.Messages = append(subAgent.Messages, toolMsg)
		}
	}

	return "", fmt.Errorf("sub-agent exceeded maximum iterations (%d)", MaxSubAgentIterations)
}
