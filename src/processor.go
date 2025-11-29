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
		var logMessage string

		switch toolCall.Function.Name {
		case "execute_command":
			var args CommandArgs
			if unmarshalErr := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); unmarshalErr != nil {
				output = fmt.Sprintf("Failed to parse arguments: %s", unmarshalErr)
			} else {
				output, err = confirmAndExecute(config, args.Command)
				if err == nil {
					logMessage = "Executed bash command"
				}
			}
		case "create_todo":
			output, err = createTodo(agent.ID, toolCall.Function.Arguments)
			if err == nil {
				logMessage = "Created todo item"
			}
		case "update_todo":
			output, err = updateTodo(agent.ID, toolCall.Function.Arguments)
			if err == nil {
				logMessage = "Updated todo item"
			}
		case "get_todo_list":
			output, err = getTodoList(agent.ID)
			if err == nil {
				logMessage = "Retrieved todo list"
			}
		case "spawn_agent":
			var args SubAgentTask
			if unmarshalErr := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); unmarshalErr != nil {
				output = fmt.Sprintf("Failed to parse arguments: %s", unmarshalErr)
			} else {
				fmt.Printf("%sSpawning sub-agent for task: %s%s%s\n", ColorMeta, ColorHighlight, args.Task, ColorReset)
				output, err = runSubAgent(args.Task, config)
				// Sub-agent logs are handled within runSubAgent, but we can add a summary here if needed
				logMessage = "Sub-agent finished task"
			}
		case "use_mcp_tool":
			var args UseMCPToolArgs
			if unmarshalErr := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); unmarshalErr != nil {
				output = fmt.Sprintf("Failed to parse arguments for use_mcp_tool: %s", unmarshalErr)
			} else {
				output, err = useMCPTool(args.ServerName, args.ToolName, args.Arguments)
				if err == nil {
					logMessage = fmt.Sprintf("Called MCP server: %s (%s)", args.ServerName, args.ToolName)
				}
			}
		case "create_note":
			output, err = createNote(toolCall.Function.Arguments)
			if err == nil {
				logMessage = "Created note"
			}
		case "update_note":
			output, err = updateNote(toolCall.Function.Arguments)
			if err == nil {
				logMessage = "Updated note"
			}
		case "delete_note":
			output, err = deleteNote(toolCall.Function.Arguments)
			if err == nil {
				logMessage = "Deleted note"
			}
		case "name_session":
			output, err = nameSession(agent, toolCall.Function.Arguments)
			if err == nil {
				logMessage = "Named session"
			}
		default:
			output = fmt.Sprintf("Unknown tool: %s", toolCall.Function.Name)
			logMessage = fmt.Sprintf("Executed unknown tool: %s", toolCall.Function.Name)
		}

		if err != nil {
			output = fmt.Sprintf("Tool execution error: %s", err)
			// Print error in meta color or red?
			fmt.Printf("%s%s%s\n", ColorRed, output, ColorReset)
		} else if logMessage != "" {
			// Simplified output style
			fmt.Printf("%s%s%s\n", ColorMeta, logMessage, ColorReset)
		}

		toolMsg := Message{
			Role:       "tool",
			ToolCallID: toolCall.ID,
			Content:    &output,
		}
		agent.Messages = append(agent.Messages, toolMsg)
	}
}
