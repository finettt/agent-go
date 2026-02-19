package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// formatToolCallCompact returns a compact display string for a tool call, e.g. "Bash(ls -la)" or "create_todo(...)".
// Used for error display so the user sees what tool failed without verbose error messages.
func formatToolCallCompact(toolCall ToolCall) string {
	name := toolCall.Function.Name
	argsRaw := toolCall.Function.Arguments

	// For execute_command, show as Bash (command) — matches success log format in subagent
	if name == "execute_command" {
		var args CommandArgs
		if err := json.Unmarshal([]byte(argsRaw), &args); err == nil && args.Command != "" {
			cmd := args.Command
			if len(cmd) > 80 {
				cmd = cmd[:77] + "..."
			}
			return fmt.Sprintf("Bash (%s)", cmd)
		}
	}

	// For spawn_agent, show as SpawnAgent (task_summary)
	if name == "spawn_agent" {
		var args SubAgentTask
		if err := json.Unmarshal([]byte(argsRaw), &args); err == nil && args.Task != "" {
			task := args.Task
			if len(task) > 60 {
				task = task[:57] + "..."
			}
			label := "SpawnAgent"
			if args.Agent != "" {
				label = fmt.Sprintf("SpawnAgent[%s]", args.Agent)
			}
			return fmt.Sprintf("%s (%s)", label, task)
		}
	}

	// For use_mcp_tool, show as MCP:server.tool (...)
	if name == "use_mcp_tool" {
		var args UseMCPToolArgs
		if err := json.Unmarshal([]byte(argsRaw), &args); err == nil {
			return fmt.Sprintf("MCP:%s.%s (...)", args.ServerName, args.ToolName)
		}
	}

	// Generic: show ToolName (truncated args)
	summary := argsRaw
	if len(summary) > 60 {
		summary = summary[:57] + "..."
	}
	return fmt.Sprintf("%s (%s)", name, summary)
}

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
				modelName := strings.TrimSpace(args.Model)

				// Avoid dumping long task prompts into the user's console by default.
				// Only show the full task when sub-agent verbose mode is "Full" (2).
				if !pipelineMode {
					if config.SubAgentVerboseMode == 2 {
						if agentName != "" {
							fmt.Printf("%sSpawning sub-agent (%s, %s) for task: %s%s%s\n", ColorMeta, agentName, modelName, ColorHighlight, args.Task, ColorReset)
						} else {
							fmt.Printf("%sSpawning sub-agent (%s) for task: %s%s%s\n", ColorMeta, modelName, ColorHighlight, args.Task, ColorReset)
						}
					} else {
						if agentName != "" {
							fmt.Printf("%sSpawning sub-agent (%s, %s)%s\n", ColorMeta, agentName, modelName, ColorReset)
						} else {
							fmt.Printf("%sSpawning sub-agent (%s)%s\n", ColorMeta, modelName, ColorReset)
						}
					}
				}

				output, err = runSubAgentWithAgent(args.Task, agentName, modelName, config)
				logMessage = "Sub-agent finished task"
			}

			if err != nil {
				output = fmt.Sprintf("Tool execution error: %s", err)
				if !pipelineMode {
					fmt.Printf("%s==> %s%s\n", ColorRed, formatToolCallCompact(toolCall), ColorReset)
				}
			} else if logMessage != "" && !pipelineMode {
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

		// Auto-checkpoint for dangerous tools
		// We do this BEFORE the switch to ensure state is saved before any potential damage.
		// Dangerous tools: execute_command, spawn_agent, use_mcp_tool, kill_background_command
		dangerousTools := map[string]bool{
			"execute_command":         true,
			"spawn_agent":             true,
			"use_mcp_tool":            true,
			"kill_background_command": true,
		}

		if dangerousTools[toolCall.Function.Name] {
			// Create auto-checkpoint for dangerous tools
			// Note: We create checkpoints regardless of OperationMode since MCP tools
			// could potentially execute commands even in Plan mode
			if _, err := createCheckpoint(agent, config, fmt.Sprintf("Auto-checkpoint before %s", toolCall.Function.Name), true); err != nil {
				// Log error but proceed
				fmt.Printf("%sWarning: Failed to create auto-checkpoint: %v%s\n", ColorYellow, err, ColorReset)
			}
		}

		switch toolCall.Function.Name {
		case "execute_command":
			var args CommandArgs
			if unmarshalErr := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); unmarshalErr != nil {
				output = fmt.Sprintf("Failed to parse arguments: %s", unmarshalErr)
			} else {
				// Only show execution message if not in Plan mode and not in pipeline mode
				if config.OperationMode != Plan && !pipelineMode {
					logMessage = fmt.Sprintf("%sExecuting command: %s%s\n", ColorMeta, args.Command, ColorReset)
				}

				// Background execution is a user choice in Ask mode (not agent-controlled).
				if pipelineMode {
					// In pipeline mode, always execute silently in the foreground without prompts/logs.
					output, err = executeCommandSilent(args.Command)
				} else {
					output, err = confirmAndExecute(config, args.Command)
				}
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
				if !pipelineMode {
					fmt.Printf("\n%s%s%sSuggested Plan: %s%s\n", StyleItalic, StyleUnderline, ColorHighlight, args.Name, ColorReset)
					fmt.Printf("%s%s%s\n", ColorMain, args.Description, ColorReset)
					fmt.Printf("%sApprove this plan? [y/N]: %s", ColorCyan, ColorReset)
				}
				var response string
				if pipelineMode {
					// In pipeline mode we cannot interactively ask; default to rejection
					response = "n"
				} else {
					fmt.Scanln(&response)
				}
				if strings.ToLower(strings.TrimSpace(response)) == "y" {
					// User approved the plan: mark for deferred switch to build mode
					output = "Plan approved by user. Switching to build mode to implement the plan..."
					logMessage = "Plan approved - will switch to build agent"

					// Set flag for deferred agent switch in main loop
					shouldSwitchToBuild = true

					// Optionally keep deprecated OperationMode in sync
					config.OperationMode = Build
					if err := saveConfig(config); err != nil {
						fmt.Fprintf(os.Stderr, "Error saving config: %s\n", err)
					}

					// Save plan to file (same behavior as before)
					timestamp := time.Now().Format("20060102_150405")
					safeName := strings.ReplaceAll(strings.ToLower(args.Name), " ", "_")
					safeName = strings.ReplaceAll(safeName, "/", "-") // Basic sanitization
					filename := fmt.Sprintf("plan_%s_%s.md", timestamp, safeName)

					// Ensure plans directory exists
					cwd, _ := os.Getwd()
					agentGoDir := filepath.Join(cwd, ".agent-go")
					plansDir := filepath.Join(agentGoDir, "plans")
					if err := os.MkdirAll(plansDir, 0755); err != nil {
						if !pipelineMode {
							fmt.Printf("Error creating plans directory: %v\n", err)
						}
					}

					filePath := filepath.Join(plansDir, filename)
					content := fmt.Sprintf("# %s\n\n%s", args.Name, args.Description)
					if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
						if !pipelineMode {
							fmt.Printf("Error saving plan to file: %v\n", err)
						}
					} else if !pipelineMode {
						fmt.Printf("Plan saved to %s\n", filePath)
						// Also update '.agent-go/current_plan.md' for easy access / inclusion in prompts
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
				if !pipelineMode {
					fmt.Println(output)
				}
			}
		case "update_todo":
			output, err = updateTodo(agent.ID, toolCall.Function.Arguments)
			if err == nil {
				logMessage = "Updated todo item"
				if !pipelineMode {
					fmt.Println(output)
				}
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
		case "clear_todo":
			output, err = clearTodo(agent.ID)
			if err == nil {
				logMessage = "Cleared todo list"
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
		case "create_checkpoint":
			var args map[string]string
			if unmarshalErr := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); unmarshalErr != nil {
				output = fmt.Sprintf("Failed to parse arguments: %s", unmarshalErr)
			} else {
				name := args["name"]
				if name == "" {
					name = "Manual Checkpoint"
				}
				var id string
				id, err = createCheckpoint(agent, config, name, false)
				if err == nil {
					output = fmt.Sprintf("Checkpoint created with ID: %s", id)
					logMessage = fmt.Sprintf("Created checkpoint '%s'", name)
				}
			}
		case "list_checkpoints":
			checkpoints, errC := listCheckpoints(agent.ID)
			if errC != nil {
				err = errC
			} else {
				if len(checkpoints) == 0 {
					output = "No checkpoints found."
				} else {
					var sb strings.Builder
					sb.WriteString("Checkpoints:\n")
					for _, cp := range checkpoints {
						sb.WriteString(fmt.Sprintf("- %s (%s): %s\n", cp.ID, cp.CreatedAt.Format("2006-01-02 15:04:05"), cp.Name))
					}
					output = sb.String()
				}
				logMessage = "Listed checkpoints"
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
			if !pipelineMode {
				fmt.Printf("%s==> %s%s\n", ColorRed, formatToolCallCompact(toolCall), ColorReset)
			}
		} else if logMessage != "" && !pipelineMode {
			// Always log the action summary (formerly only in verbose)
			// In verbose mode, we might want even more details, but for now let's make the summary always visible
			// as requested: "what is currently output in verbose mode should be output always".
			// The previous code only printed logMessage if config.Verbose.
			// Now we print it always (except in pipeline mode).
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
