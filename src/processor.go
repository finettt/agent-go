package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

// processToolCalls handles the logic for executing tool calls from the API response
func processToolCalls(agent *Agent, toolCalls []ToolCall, config *Config) {
	var wg sync.WaitGroup
	// A channel to collect results from parallel sub-agent execution if needed.
	// But for now, we append messages to the agent sequentially after wait, or we need a thread-safe way.
	// Since the order of tool outputs in the messages list matters for the LLM to match calls with outputs,
	// we should maintain order or just append them. However, `agent.Messages` modification must be thread-safe.
	// Let's use a mutex for updating agent.Messages.
	var agentMutex sync.Mutex

	// Channel to collect tool outputs to ensure deterministic order if possible,
	// or just use a slice with index. For simplicity in this "parallel sub-agent" task,
	// we will just lock when appending results.

	for _, toolCall := range toolCalls {
		// If it's a spawn_agent call, we run it in a goroutine
		if toolCall.Function.Name == "spawn_agent" {
			wg.Add(1)
			go func(tc ToolCall) {
				defer wg.Done()
				var output string
				var err error
				var logMessage string

				var args SubAgentTask
				if unmarshalErr := json.Unmarshal([]byte(tc.Function.Arguments), &args); unmarshalErr != nil {
					output = fmt.Sprintf("Failed to parse arguments: %s", unmarshalErr)
				} else {
					if config.Verbose {
						fmt.Printf("%sSpawning sub-agent for task: %s%s%s\n", ColorMeta, ColorHighlight, args.Task, ColorReset)
					}
					output, err = runSubAgent(args.Task, config)
					logMessage = "Sub-agent finished task"
				}

				if err != nil {
					output = fmt.Sprintf("Tool execution error: %s", err)
					fmt.Printf("%s%s%s\n", ColorRed, output, ColorReset)
				} else if logMessage != "" {
					if config.Verbose {
						fmt.Printf("%s%s%s\n", ColorMeta, logMessage, ColorReset)
					}
				}

				toolMsg := Message{
					Role:       "tool",
					ToolCallID: tc.ID,
					Content:    &output,
				}

				agentMutex.Lock()
				agent.Messages = append(agent.Messages, toolMsg)
				agentMutex.Unlock()
			}(toolCall)
			continue
		}

		// For other tools, run sequentially (or we could parallelize everything, but let's stick to sub-agents first)
		// Actually, if we run others sequentially but sub-agents in parallel, we might mix up the order of completion.
		// But since we are inside a loop, "continue" above means we skip the sequential block for sub-agents.
		// The sequential block below handles non-sub-agent tools.

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
			// Only log if verbose or if it's a significant action
			if config.Verbose {
				fmt.Printf("%s%s%s\n", ColorMeta, logMessage, ColorReset)
			}
		}

		toolMsg := Message{
			Role:       "tool",
			ToolCallID: toolCall.ID,
			Content:    &output,
		}
		agentMutex.Lock()
		agent.Messages = append(agent.Messages, toolMsg)
		agentMutex.Unlock()
	}

	// Wait for all spawned sub-agents to complete
	wg.Wait()
}
