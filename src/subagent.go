package main

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

func runSubAgent(task string, config *Config) (string, error) {
	subAgent := &Agent{
		ID: uuid.New().String(),
		Messages: []Message{
			{
				Role:    "system",
				Content: stringp(getSystemInfo() + "\n\nYou are a sub-agent tasked with completing a specific goal. You have access to the 'execute_command' and todo list management tools. Plan your steps and execute them sequentially. When you have finished the task, output the final result as a single response."),
			},
			{
				Role:    "user",
				Content: stringp(task),
			},
		},
	}

	for {
		resp, err := sendAPIRequest(subAgent, config, true) // true: include all tools for sub-agents
		if err != nil {
			return "", fmt.Errorf("sub-agent API request failed: %w", err)
		}

		if len(resp.Choices) == 0 {
			return "", fmt.Errorf("sub-agent received an empty response from the API")
		}

		assistantMsg := resp.Choices[0].Message
		subAgent.Messages = append(subAgent.Messages, assistantMsg)

		if len(assistantMsg.ToolCalls) == 0 {
			// If there are no more tool calls, the sub-agent's work is done.
			// The last message's content is considered the final result.
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
				if err = json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err == nil {
					output, err = executeCommand(args.Command)
				}
			case "create_todo":
				output, err = createTodo(subAgent.ID, toolCall.Function.Arguments)
			case "update_todo":
				output, err = updateTodo(subAgent.ID, toolCall.Function.Arguments)
			case "get_todo_list":
				output, err = getTodoList(subAgent.ID)
			}

			if err != nil {
				output = fmt.Sprintf("Tool execution error: %s", err)
			}

			toolMsg := Message{
				Role:       "tool",
				ToolCallID: toolCall.ID,
				Content:    stringp(output),
			}
			subAgent.Messages = append(subAgent.Messages, toolMsg)
		}
	}
}
