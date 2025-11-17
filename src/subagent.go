package main

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// runSubAgent executes a task in a separate sub-agent context
func runSubAgent(task string, config *Config) (string, error) {
	systemPrompt := getSystemInfo() + "\n\nYou are a sub-agent tasked with completing a specific goal. You have access to the 'execute_command' and todo list management tools. Plan your steps and execute them sequentially. When you have finished the task, output the final result as a single response."

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

			switch toolCall.Function.Name {
			case "execute_command":
				var args CommandArgs
				if unmarshalErr := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); unmarshalErr != nil {
					output = fmt.Sprintf("Failed to parse arguments: %s", unmarshalErr)
				} else {
					output, err = confirmAndExecute(config, args.Command)
				}
			case "create_todo":
				output, err = createTodo(subAgent.ID, toolCall.Function.Arguments)
			case "update_todo":
				output, err = updateTodo(subAgent.ID, toolCall.Function.Arguments)
			case "get_todo_list":
				output, err = getTodoList(subAgent.ID)
			default:
				output = fmt.Sprintf("Unknown tool: %s", toolCall.Function.Name)
			}

			if err != nil {
				output = fmt.Sprintf("Tool execution error: %s", err)
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
