package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// processToolCalls handles the logic for executing tool calls from the API response
func processToolCalls(agent *Agent, toolCalls []ToolCall, config *Config) {

	for _, toolCall := range toolCalls {
		// If it's a spawn_agent call, we run it sequentially
		if toolCall.Function.Name == "spawn_agent" {
			var output string
			var err error
			var logMessage string

			var args SubAgentTask
			if unmarshalErr := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); unmarshalErr != nil {
				output = fmt.Sprintf("Failed to parse arguments: %s", unmarshalErr)
			} else {
				agentName := strings.TrimSpace(args.Agent)

				// Avoid dumping long task prompts into the user's console by default.
				// Only show the full task when sub-agent verbose mode is "Full" (2).
				if config.SubAgentVerboseMode == 2 {
					if agentName != "" {
						fmt.Printf("%sSpawning sub-agent (%s) for task: %s%s%s\n", ColorMeta, agentName, ColorHighlight, args.Task, ColorReset)
					} else {
						fmt.Printf("%sSpawning sub-agent for task: %s%s%s\n", ColorMeta, ColorHighlight, args.Task, ColorReset)
					}
				} else {
					if agentName != "" {
						fmt.Printf("%sSpawning sub-agent (%s)%s\n", ColorMeta, agentName, ColorReset)
					} else {
						fmt.Printf("%sSpawning sub-agent%s\n", ColorMeta, ColorReset)
					}
				}

				output, err = runSubAgentWithAgent(args.Task, agentName, config)
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
				ToolCallID: toolCall.ID,
				Content:    &output,
			}

			agent.Messages = append(agent.Messages, toolMsg)
			continue
		}

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
					logMessage = fmt.Sprintf("%sExecuting command: %s%s\n", ColorMeta, args.Command, ColorReset)
				}

				// Background execution is a user choice in Ask mode (not agent-controlled).
				output, err = confirmAndExecute(config, args.Command)
				if output == "Command not executed by user." {
					logMessage = fmt.Sprintf("%sCommand not executed by user.%s\n", ColorMeta, ColorReset)
				}
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
				fmt.Printf("\n%sSuggested Plan: %s%s\n", ColorHighlight, args.Name, ColorReset)
				fmt.Printf("%s%s%s\n", ColorMain, args.Description, ColorReset)
				fmt.Printf("%sApprove this plan? [y/N]: %s", ColorCyan, ColorReset)
				var response string
				fmt.Scanln(&response)
				if strings.ToLower(strings.TrimSpace(response)) == "y" {
					config.OperationMode = Build
					output = "Plan approved by user. Switched to Build mode. You may now execute commands."
					logMessage = "Plan approved - Switched to Build mode"
					if err := saveConfig(config); err != nil {
						fmt.Fprintf(os.Stderr, "Error saving config: %s\n", err)
					}

					// Save plan to file
					timestamp := time.Now().Format("20060102_150405")
					safeName := strings.ReplaceAll(strings.ToLower(args.Name), " ", "_")
					safeName = strings.ReplaceAll(safeName, "/", "-") // Basic sanitization
					filename := fmt.Sprintf("plan_%s_%s.md", timestamp, safeName)

					// Ensure plans directory exists
					cwd, _ := os.Getwd()
					agentGoDir := filepath.Join(cwd, ".agent-go")
					plansDir := filepath.Join(agentGoDir, "plans")
					if err := os.MkdirAll(plansDir, 0755); err != nil {
						fmt.Printf("Error creating plans directory: %v\n", err)
					}

					filePath := filepath.Join(plansDir, filename)
					content := fmt.Sprintf("# %s\n\n%s", args.Name, args.Description)
					if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
						fmt.Printf("Error saving plan to file: %v\n", err)
					} else {
						fmt.Printf("Plan saved to %s\n", filePath)
						// Update config to track current plan? Or just symlink?
						// Let's create a symlink or just copy to '.agent-go/current_plan.md' for easy access
						currentPlanPath := filepath.Join(agentGoDir, "current_plan.md")
						if err := os.WriteFile(currentPlanPath, []byte(content), 0644); err != nil {
							fmt.Printf("Error saving current_plan.md: %v\n", err)
						}
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
		case "get_current_task":
			output, err = getCurrentTask(agent.ID)
			if err == nil {
				logMessage = "Retrieved current task"
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
		case "create_agent_definition":
			output, err = createAgentDefinition(toolCall.Function.Arguments)
			if err == nil {
				logMessage = "Created agent definition"
			}
		default:
			// Check if it's a custom skill
			var skillExecuted bool
			for _, skill := range config.Skills {
				if skill.Name == toolCall.Function.Name {
					// Prepare command with arguments
					var argsMap map[string]interface{}
					if unmarshalErr := json.Unmarshal([]byte(toolCall.Function.Arguments), &argsMap); unmarshalErr != nil {
						output = fmt.Sprintf("Failed to parse arguments for skill %s: %s", skill.Name, unmarshalErr)
					} else {
						// Pass arguments as JSON environment variable
						argsJSON, _ := json.Marshal(argsMap)

						// Execute skill command
						if config.OperationMode != Plan {
							logMessage = fmt.Sprintf("%sExecuting skill: %s%s\n", ColorMeta, skill.Name, ColorReset)
						}

						// Use executeSkill to handle .sh files directly or fallback to shell
						output, err = executeSkill(skill.Command, argsJSON)
						if err == nil {
							logMessage = fmt.Sprintf("Executed skill: %s", skill.Name)
						}
					}
					skillExecuted = true
					break
				}
			}

			if !skillExecuted {
				output = fmt.Sprintf("Unknown tool: %s", toolCall.Function.Name)
				logMessage = fmt.Sprintf("Executed unknown tool: %s", toolCall.Function.Name)
			}
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
		agent.Messages = append(agent.Messages, toolMsg)
	}
}
