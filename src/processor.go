package main

import (
	"encoding/json"
	"fmt"
)

// processToolCalls handles the logic for executing tool calls from the API response
func processToolCalls(agent *Agent, toolCalls []ToolCall, config *Config) {
	for _, toolCall := range toolCalls {
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
			output, err = createTodo(agent.ID, toolCall.Function.Arguments)
		case "update_todo":
			output, err = updateTodo(agent.ID, toolCall.Function.Arguments)
		case "get_todo_list":
			output, err = getTodoList(agent.ID)
		case "spawn_agent":
			var args SubAgentTask
			if unmarshalErr := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); unmarshalErr != nil {
				output = fmt.Sprintf("Failed to parse arguments: %s", unmarshalErr)
			} else {
				fmt.Printf("%sSpawning sub-agent for task: %s%s\n", ColorYellow, args.Task, ColorReset)
				output, err = runSubAgent(args.Task, config)
				fmt.Printf("%sSub-agent finished with result: %s%s\n", ColorYellow, output, ColorReset)
			}
		case "use_mcp_tool":
			var args UseMCPToolArgs
			if unmarshalErr := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); unmarshalErr != nil {
				output = fmt.Sprintf("Failed to parse arguments for use_mcp_tool: %s", unmarshalErr)
			} else {
				output, err = useMCPTool(args.ServerName, args.ToolName, args.Arguments)
			}
		default:
			output = fmt.Sprintf("Unknown tool: %s", toolCall.Function.Name)
		}

		if err != nil {
			output = fmt.Sprintf("Tool execution error: %s", err)
		}

		// Print tool output for visibility
		fmt.Printf("%s%s%s\n", ColorPurple, output, ColorReset)

		toolMsg := Message{
			Role:       "tool",
			ToolCallID: toolCall.ID,
			Content:    &output,
		}
		agent.Messages = append(agent.Messages, toolMsg)
	}
}