package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
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
					fmt.Printf("%sSpawning sub-agent for task: %s%s%s\n", ColorMeta, ColorHighlight, args.Task, ColorReset)
					output, err = runSubAgent(args.Task, config)
					logMessage = "Sub-agent finished task"
				}

				if err != nil {
					output = fmt.Sprintf("Tool execution error: %s", err)
					fmt.Printf("%s%s%s\n", ColorRed, output, ColorReset)
				} else if logMessage != "" {
					fmt.Printf("%s%s%s\n", ColorMeta, logMessage, ColorReset)
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
				// Only show execution message if not in Plan mode
				if config.OperationMode != Plan {
					logMessage = fmt.Sprintf("%sExecuting command: %s (Background: %t)%s\n", ColorMeta, args.Command, args.Background, ColorReset)
				}

				output, err = confirmAndExecute(config, args.Command, args.Background)
			}
		case "kill_background_command":
			var args KillBackgroundCommandArgs
			if unmarshalErr := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); unmarshalErr != nil {
				output = fmt.Sprintf("Failed to parse arguments: %s", unmarshalErr)
			} else {
				output, err = killBackgroundCommand(args.PID)
				if err == nil {
					logMessage = fmt.Sprintf("Killed background process %d", args.PID)
				}
			}
		case "get_background_logs":
			var args GetBackgroundLogsArgs
			if unmarshalErr := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); unmarshalErr != nil {
				output = fmt.Sprintf("Failed to parse arguments: %s", unmarshalErr)
			} else {
				output, err = getBackgroundLogs(args.PID)
				if err == nil {
					logMessage = fmt.Sprintf("Retrieved logs for PID %d", args.PID)
				}
			}
		case "list_background_commands":
			output = listBackgroundCommands()
			logMessage = "Listed background commands"
		case "suggest_plan":
			var args SuggestPlanArgs
			if unmarshalErr := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); unmarshalErr != nil {
				output = fmt.Sprintf("Failed to parse arguments: %s", unmarshalErr)
			} else {
				fmt.Printf("\n%sSuggested Plan:%s\n%s\n", ColorHighlight, ColorReset, args.Plan)
				fmt.Printf("%sApprove this plan? [y/N]: %s", ColorRed, ColorReset)
				var response string
				fmt.Scanln(&response)
				if strings.ToLower(strings.TrimSpace(response)) == "y" {
					config.OperationMode = Build
					output = "Plan approved by user. Switched to Build mode. You may now execute commands."
					logMessage = "Plan approved - Switched to Build mode"
					if err := saveConfig(config); err != nil {
						fmt.Fprintf(os.Stderr, "Error saving config: %s\n", err)
					}
				} else {
					output = "Plan rejected by user."
					logMessage = "Plan rejected"
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
			// Always log the action summary (formerly only in verbose)
			// In verbose mode, we might want even more details, but for now let's make the summary always visible
			// as requested: "what is currently output in verbose mode should be output always".
			// The previous code only printed logMessage if config.Verbose.
			// Now we print it always.
			fmt.Printf("%s%s%s\n", ColorMeta, logMessage, ColorReset)

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
