package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/chzyer/readline"
)

var config *Config
var agent *Agent
var shellMode = false

// Agent Studio + task-specific agents
type AgentConfigSnapshot struct {
	Model     string
	Temp      float32
	MaxTokens int
}

var prevAgentConfigSnapshot *AgentConfigSnapshot

var totalTokens = 0
var totalPromptTokens = 0
var totalCompletionTokens = 0
var totalToolCalls = 0

func formatTokenCount(count int) string {
	if count >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(count)/1000000)
	}
	if count >= 1000 {
		return fmt.Sprintf("%.1fK", float64(count)/1000)
	}
	return fmt.Sprintf("%d", count)
}

func main() {
	// Check for command line task argument
	if len(os.Args) > 1 {
		task := strings.Join(os.Args[1:], " ")
		runTask(task)
		return
	}

	// Display ASCII logo
	printLogo()

	config = loadConfig()
	if config.APIKey == "" {
		runSetup()
	}

	agent = &Agent{
		ID:       "main",
		Messages: make([]Message, 0),
	}

	systemPrompt := buildSystemPrompt("")
	agent.Messages = append(agent.Messages, Message{
		Role:    "system",
		Content: &systemPrompt,
	})

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if agent != nil && len(agent.Messages) > 1 {
			if err := saveSession(agent); err != nil {
				fmt.Fprintf(os.Stderr, "\nFailed to save session: %v\n", err)
			} else {
				fmt.Printf("\nSession '%s' saved.\n", agent.ID)
			}
		}
		// Check for running background processes
		if hasRunningBackgroundProcesses() {
			fmt.Println("\nWarning: There are running background processes. They will be terminated.")
			// We could ask for confirmation or list them, but for now let's just warn and exit.
		}

		fmt.Println("\nBye!")
		os.Exit(0)
	}()

	runCLI()
}
func printLogo() {
	fmt.Print(ColorHighlight + `
 █████╗  ██████╗ ███████╗███╗   ██╗████████╗    ██████╗  ██████╗
██╔══██╗██╔════╝ ██╔════╝████╗  ██║╚══██╔══╝   ██╔════╝ ██╔═══██╗
███████║██║  ███╗█████╗  ██╔██╗ ██║   ██║█████╗██║  ███╗██║   ██║
██╔══██║██║   ██║██╔══╝  ██║╚██╗██║   ██║╚════╝██║   ██║██║   ██║
██║  ██║╚██████╔╝███████╗██║ ╚████║   ██║      ╚██████╔╝╚██████╔╝
╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚═╝  ╚═══╝   ╚═╝       ╚═════╝  ╚═════╝

` + ColorReset)
}

// showHelp displays all available slash commands grouped by family
func showHelp() {
	// Helper for printing main commands
	printCmd := func(cmd, desc string) {
		fmt.Printf("%s%s%s - %s%s%s\n", ColorHighlight, cmd, ColorReset, ColorMeta, desc, ColorReset)
	}

	// Helper for printing subcommands (indented)
	printSubCmd := func(subCmd, desc string) {
		fmt.Printf("    %s%s%s - %s%s%s\n", ColorMain, subCmd, ColorReset, ColorMeta, desc, ColorReset)
	}

	fmt.Println("Available commands:")

	printCmd("/help, /?", "Show this help message")
	printCmd("/config", "Display current configuration")
	printCmd("/shell", "Enter shell mode for direct command execution")
	printCmd("/bg", "Background process management")
	printSubCmd("list", "List background processes")
	printSubCmd("view <pid>", "View logs (stdout/stderr) for a background process")
	printSubCmd("kill <pid>", "Kill a background process")
	printCmd("/clear", "Clear context without compressing")
	printCmd("/compress", "Compress context and start new chat thread")
	printCmd("/edit", "Edit prompt in nano editor")
	printCmd("/quit", "Exit the application")
	printCmd("/sandbox", "Relaunch agent-go in a Docker sandbox")

	printCmd("/model <name>", "Set the AI model (e.g., gpt-4)")
	printCmd("/provider <url>", "Set the API provider URL")

	printCmd("/contextlength <val>", "Set the model context length (e.g., 131072)")
	printCmd("/stream on|off", "Toggle streaming mode")

	printCmd("/rag", "Retrieval-Augmented Generation controls")
	printSubCmd("on|off", "Toggle RAG feature")
	printSubCmd("path <path>", "Set the RAG documents path")

	printCmd("/usage <1|2|3>", "Set usage verbosity (1: Silent, 2: Basic, 3: Detailed)")
	printCmd("/cost", "Show current usage statistics")

	printCmd("/todo", "Display the current todo list")

	printCmd("/notes", "Notes management")
	printSubCmd("list", "List all notes")
	printSubCmd("view <name>", "View a specific note")

	printCmd("/session", "Manage chat sessions")
	printSubCmd("list", "List saved sessions")
	printSubCmd("restore <name>", "Restore a session")
	printSubCmd("new", "Create a new session with fresh context")
	printSubCmd("rm <name>", "Delete a saved session")

	printCmd("/mcp", "Model Context Protocol server management")
	printSubCmd("add <name> <cmd>", "Add an MCP server")
	printSubCmd("remove <name>", "Remove an MCP server")
	printSubCmd("list", "List MCP servers")

	printCmd("/agent", "Autonomous agent management")
	printSubCmd("studio [spec]", "Start Agent Studio to create a task-specific agent")
	printSubCmd("list", "List saved task-specific agents")
	printSubCmd("view <name>", "View a saved agent definition")
	printSubCmd("use <name>", "Activate a saved agent for the current chat")
	printSubCmd("clear", "Clear active agent and restore previous model settings")
	printSubCmd("rm <name>", "Delete a saved agent definition")

	printCmd("/subagents", "Configure autonomous sub-agents")
	printSubCmd("[on|off]", "Enable or disable sub-agents")
	printSubCmd("verbose <1|2>", "Set verbosity level")
	printCmd("/security", "Spawn a subagent to review current changes")

	printCmd("/mode", "Toggle between Plan and Build operation modes")
	printCmd("/plan", "Alias for toggling operation mode")
	printCmd("/ask on|off", "Enable/Disable confirmation for commands (Ask vs YOLO)")
}

func runCLI() {
	home, _ := os.UserHomeDir()
	historyFile := filepath.Join(home, ".config", "agent-go", "history.txt")
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "> ",
		HistoryFile:  historyFile,
		AutoComplete: buildCompleter(config),
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		err := rl.Close()
		if err != nil {
			fmt.Printf("failed to close readline: %v\n", err)
		}
	}()

	cwd, _ := os.Getwd()
	sandboxStatus := "Sandbox: Off"
	if _, err := os.Stat("/.dockerenv"); err == nil {
		sandboxStatus = "Sandbox: On"
	}
	fmt.Printf("Welcome to Agent-Go!\n%s%s • %s • %s%s\n", ColorMeta, config.Model, cwd, sandboxStatus, ColorReset)

	for {
		if agentStudioMode {
			rl.SetPrompt(StyleBold + ColorHighlight + ">>> ")
		} else if shellMode {
			rl.SetPrompt(StyleBold + ColorCyan + "! ")
		} else {
			rl.SetPrompt(StyleBold + "> ")
		}

		userInput, err := rl.Readline()
		fmt.Print(ColorReset) // Reset after input is complete

		if err == readline.ErrInterrupt {
			continue
		} else if err == io.EOF {
			break
		}

		userInput = strings.TrimSpace(userInput)

		// Process @filename mentions
		userInput = processFileMentions(userInput)

		if agentStudioMode {
			if userInput == "exit" {
				agentStudioMode = false
				studioAgent = nil
				fmt.Println("Exited Agent Studio.")
				continue
			}
			if userInput == "" {
				continue
			}

			// Add user message to studio history
			studioAgent.Messages = append(studioAgent.Messages, Message{Role: "user", Content: &userInput})

			// Run one (possibly multi-tool) studio turn
			if err := runAgentStudioTurn(config); err != nil {
				fmt.Fprintf(os.Stderr, "Agent Studio error: %v\n", err)
			}

			// If studio finished (agent created), clear studio state.
			if !agentStudioMode {
				studioAgent = nil
				fmt.Println("Agent Studio finished.")
			}
			continue
		}

		if shellMode {
			if userInput == "exit" {
				shellMode = false
				fmt.Println("Exited shell mode.")
				continue
			}
			if userInput == "" {
				continue
			}
			// Shell mode commands can be run in foreground or background (Ask mode prompts the user).
			output, err := confirmAndExecute(config, userInput)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			}
			fmt.Println(output)
			continue
		}

		if userInput == "" {
			continue
		}

		if strings.HasPrefix(userInput, "/") {
			handleSlashCommand(userInput)
			continue
		}

		// RAG processing
		if config.RAGEnabled && config.RAGPath != "" {
			snippets, err := searchRAGFiles(config.RAGPath, userInput, config.RAGSnippets)
			if err == nil && snippets != "" {
				userInput = fmt.Sprintf("User asked: %s\n\nRelevant snippets from local documents:\n%s\n\nPlease answer based on the user's request and the provided context.", userInput, snippets)
			}
		}

		// Add user message to agent history
		agent.Messages = append(agent.Messages, Message{Role: "user", Content: &userInput})

		// Message history is now unlimited

		// Agentic loop
		for {
			// Auto-compress context if enabled and token count exceeds 75% of context length
			if config.AutoCompress && totalTokens > (config.ModelContextLength*3/4) {
				compressAndStartNewChat()
				fmt.Println("Context compressed due to token limit. New user input is required.")
				continue // Restart the outer loop to get fresh user input
			}

			var resp *APIResponse
			var err error

			// Use streaming or regular API based on config
			if config.Stream {
				resp, err = sendAPIRequestStreaming(agent, config, config.SubagentsEnabled)
			} else {
				resp, err = sendAPIRequest(agent, config, config.SubagentsEnabled)
			}

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				break // Break from the agentic loop
			}

			if len(resp.Choices) == 0 {
				fmt.Fprintf(os.Stderr, "Error: received an empty response from the API\n")
				break
			}

			assistantMsg := resp.Choices[0].Message
			agent.Messages = append(agent.Messages, assistantMsg)

			// Only print content if not streaming (streaming already printed it)
			if !config.Stream {
				if assistantMsg.ReasoningContent != nil && *assistantMsg.ReasoningContent != "" {
					fmt.Printf("%sThink...\n%s", ColorMeta, ColorReset)
				}
				if assistantMsg.Content != nil && *assistantMsg.Content != "" {
					fmt.Printf("%s%s%s\n", ColorMain, *assistantMsg.Content, ColorReset)
				}
			}

			// Update and display total tokens
			if resp.Usage.TotalTokens > 0 {
				totalTokens += resp.Usage.TotalTokens
				totalPromptTokens += resp.Usage.PromptTokens
				totalCompletionTokens += resp.Usage.CompletionTokens
			}

			// Update tool call count
			if len(assistantMsg.ToolCalls) > 0 {
				totalToolCalls += len(assistantMsg.ToolCalls)
			}

			// Display usage based on verbose mode
			switch config.UsageVerboseMode {
			case UsageSilent:
				// Do nothing
			case UsageDetailed:
				if resp.Usage.TotalTokens > 0 {
					fmt.Printf("%sUsage: %d prompt + %d completion = %d total tokens\n", ColorMeta, resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
					fmt.Printf("Total: %s tokens (%d prompt, %d completion), %d tool calls%s\n", formatTokenCount(totalTokens), totalPromptTokens, totalCompletionTokens, totalToolCalls, ColorReset)
				}
			case UsageBasic:
				fallthrough
			default:
				// Default behavior (Basic)
				if resp.Usage.TotalTokens > 0 {
					fmt.Printf("%sUsed %s%s%s tokens on %s%s\n", ColorMeta, ColorHighlight, formatTokenCount(totalTokens), ColorMeta, config.Model, ColorReset)
				}
			}

			if len(assistantMsg.ToolCalls) > 0 {
				processToolCalls(agent, assistantMsg.ToolCalls, config)
			} else {
				// Check if we got an empty response
				if assistantMsg.Content == nil || *assistantMsg.Content == "" {
					fmt.Printf("%sWarning: Received empty response from model%s\n", ColorMeta, ColorReset)
				}
				break // No more tools to call, end agent turn
			}
			// Continue loop to send tool output back to API

			// Inject reminder if background processes are running
			if hasRunningBackgroundProcesses() {
				reminder := fmt.Sprintf("REMINDER: You have running background processes:\n%s", listBackgroundCommands())
				agent.Messages = append(agent.Messages, Message{
					Role:    "system",
					Content: &reminder,
				})
				// We don't print this to the user, it's just for the agent's context
			}
		}
	}
}

// readAgentsFile checks for AGENTS.md and returns its content.
func readAgentsFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// File not found is not a fatal error.
			return "", nil
		}
		// Any other error while reading the file is a problem.
		return "", err
	}
	return string(content), nil
}

func handleSlashCommand(command string) {
	parts := strings.Fields(command)
	baseCommand := parts[0]

	switch baseCommand {
	case "/help", "/?":
		showHelp()
	case "/edit":
		editCommand()
	case "/sandbox":
		fmt.Println("Building Docker image...")
		buildCmd := exec.Command("docker", "pull", "ghcr.io/finettt/agent-go")
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr
		if err := buildCmd.Run(); err != nil {
			fmt.Printf("Failed to build Docker image: %v\n", err)
			return
		}

		cwd, _ := os.Getwd()
		home, _ := os.UserHomeDir()
		configDir := filepath.Join(home, ".config", "agent-go")

		fmt.Println("Starting sandbox environment...")
		runCmd := exec.Command("docker", "run", "-it", "--rm",
			"-v", fmt.Sprintf("%s:/workspace", cwd),
			"-v", fmt.Sprintf("%s:/home/appuser/.config/agent-go", configDir),
			"agent-go-sandbox")

		runCmd.Stdin = os.Stdin
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr

		if err := runCmd.Run(); err != nil {
			fmt.Printf("Sandbox exited with error: %v\n", err)
		} else {
			fmt.Println("Sandbox session ended.")
		}

	case "/security":
		if !config.SubagentsEnabled {
			fmt.Println("Subagents are disabled. Enable them with /subagents on to use this command.")
			return
		}
		task := "Review the current changes in the branch/working directory for security issues, bugs, and best practices. Provide a summary of findings."
		fmt.Println("Spawning subagent to review changes...")
		result, err := runSubAgent(task, config)
		if err != nil {
			fmt.Printf("Security review failed: %v\n", err)
		} else {
			fmt.Printf("\n=== Security Review ===\n%s\n", result)
		}

	case "/usage":
		if len(parts) < 2 {
			fmt.Printf("Current usage verbose mode: %d\n", config.UsageVerboseMode)
			fmt.Println("Usage: /usage <1|2|3> (1: Silent, 2: Basic, 3: Detailed)")
			return
		}
		mode, err := strconv.Atoi(parts[1])
		if err != nil || mode < 1 || mode > 3 {
			fmt.Println("Invalid mode. Please use 1, 2, or 3.")
			return
		}
		config.UsageVerboseMode = mode
		if err := saveConfig(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %s\n", err)
		}
		fmt.Printf("Usage verbose mode set to %d\n", mode)

	case "/cost":
		cwd, _ := os.Getwd()
		fmt.Printf("\n=== Current Usage Stats ===\n")
		fmt.Printf("Model:             %s\n", config.Model)
		fmt.Printf("Working Directory: %s\n", cwd)
		fmt.Printf("Total Tokens:      %d\n", totalTokens)
		fmt.Printf("  - Prompt:        %d\n", totalPromptTokens)
		fmt.Printf("  - Completion:    %d\n", totalCompletionTokens)
		fmt.Printf("Total Tool Calls:  %d\n", totalToolCalls)
		fmt.Println("===========================")

	case "/session":
		if len(parts) < 2 {
			fmt.Println("Usage: /session [list|restore <name>|rm <name>]")
			return
		}
		switch parts[1] {
		case "list":
			fmt.Println(formatSessionsList())
		case "restore":
			if len(parts) < 3 {
				fmt.Println("Usage: /session restore <name>")
				return
			}
			name := parts[2]
			// Save current session first if it has content
			if len(agent.Messages) > 1 {
				if err := saveSession(agent); err != nil {
					fmt.Fprintf(os.Stderr, "Error saving current session: %v\n", err)
				} else {
					fmt.Printf("Current session '%s' saved.\n", agent.ID)
				}
			}

			loadedSession, err := loadSession(name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading session '%s': %v\n", name, err)
				return
			}

			// Reconstruct agent from session
			agent = &Agent{
				ID:       loadedSession.ID,
				Messages: loadedSession.Messages,
			}
			// Recalculate tokens? We don't store them, so they'll reset to 0
			// Ideally we should probably estimate them or store them in session.
			// For now, resetting to 0 is acceptable behavior, though auto-compress might trigger late.
			totalTokens = 0
			fmt.Printf("Session '%s' restored.\n", name)
		case "new":
			// Save current session first if it has content
			if len(agent.Messages) > 1 {
				if err := saveSession(agent); err != nil {
					fmt.Fprintf(os.Stderr, "Error saving current session: %v\n", err)
				} else {
					fmt.Printf("Current session '%s' saved.\n", agent.ID)
				}
			}

			// Create new session with fresh context
			agent = &Agent{
				ID:       "main",
				Messages: make([]Message, 0),
			}

			// Reset token counters
			totalTokens = 0
			totalPromptTokens = 0
			totalCompletionTokens = 0

			// Add new system prompt
			systemPrompt := buildSystemPrompt("")
			agent.Messages = append(agent.Messages, Message{
				Role:    "system",
				Content: &systemPrompt,
			})

			fmt.Println("New session created with fresh context.")
		case "rm":
			if len(parts) < 3 {
				fmt.Println("Usage: /session rm <name>")
				return
			}
			name := parts[2]
			if agent.ID == name {
				fmt.Println("Cannot delete the active session.")
				return
			}
			if err := deleteSession(name); err != nil {
				fmt.Fprintf(os.Stderr, "Error deleting session: %v\n", err)
			} else {
				fmt.Printf("Session '%s' deleted.\n", name)
			}
		default:
			fmt.Println("Usage: /session [list|restore <name>|new|rm <name>]")
		}

	case "/mode":
		if config.OperationMode == Build {
			config.OperationMode = Plan
			fmt.Println("Switched to Plan mode.")
		} else {
			config.OperationMode = Build
			fmt.Println("Switched to Build mode.")
		}
		if err := saveConfig(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %s\n", err)
		}
	case "/plan":
		if config.OperationMode == Build {
			config.OperationMode = Plan
			fmt.Println("Switched to Plan mode.")
		} else {
			config.OperationMode = Build
			fmt.Println("Switched to Build mode.")
		}
		if err := saveConfig(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %s\n", err)
		}
	case "/ask":
		if len(parts) > 1 {
			switch parts[1] {
			case "on":
				config.ExecutionMode = Ask
				fmt.Println("Switched to Ask mode.")
			case "off":
				config.ExecutionMode = YOLO
				fmt.Println("Switched to YOLO mode.")
			default:
				fmt.Println("Usage: /ask [on|off]")
			}
			if err := saveConfig(config); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving config: %s\n", err)
			}
		} else {
			if config.ExecutionMode == Ask {
				fmt.Println("Ask mode is currently enabled.")
			} else {
				fmt.Println("Ask mode is currently disabled (YOLO mode).")
			}
		}
	case "/mcp":
		if len(parts) < 2 {
			fmt.Println("Usage: /mcp [add|remove|list]")
			return
		}
		switch parts[1] {
		case "add":
			if len(parts) < 4 {
				fmt.Println("Usage: /mcp add <name> <command>")
				return
			}
			name := parts[2]
			command := strings.Join(parts[3:], " ")
			if config.MCPs == nil {
				config.MCPs = make(map[string]MCPServer)
			}
			config.MCPs[name] = MCPServer{Name: name, Command: command}
			if err := saveConfig(config); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving config: %s\n", err)
			}
			fmt.Printf("MCP server '%s' added.\n", name)
		case "remove":
			if len(parts) < 3 {
				fmt.Println("Usage: /mcp remove <name>")
				return
			}
			name := parts[2]
			if _, ok := config.MCPs[name]; ok {
				delete(config.MCPs, name)
				if err := saveConfig(config); err != nil {
					fmt.Fprintf(os.Stderr, "Error saving config: %s\n", err)
				}
				fmt.Printf("MCP server '%s' removed.\n", name)
			} else {
				fmt.Printf("MCP server '%s' not found.\n", name)
			}
		case "list":
			if len(config.MCPs) == 0 {
				fmt.Println("No MCP servers configured.")
				return
			}
			fmt.Println("Configured MCP servers:")
			for name, server := range config.MCPs {
				fmt.Printf("- %s: %s\n", name, server.Command)
			}
		default:
			fmt.Println("Usage: /mcp [add|remove|list]")
		}
	case "/todo":
		list, err := getTodoList(agent.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting todo list: %s\n", err)
		} else {
			fmt.Println(list)
		}
	case "/notes":
		if len(parts) < 2 {
			fmt.Println("Usage: /notes [list|view <name>]")
			return
		}
		switch parts[1] {
		case "list":
			fmt.Println(formatNotesList())
		case "view":
			if len(parts) < 3 {
				fmt.Println("Usage: /notes view <name>")
				return
			}
			noteName := strings.Join(parts[2:], " ")
			fmt.Println(formatNoteView(noteName))
		default:
			fmt.Println("Usage: /notes [list|view <name>]")
		}
	case "/contextlength":
		if len(parts) > 1 {
			val, err := strconv.Atoi(parts[1])
			if err == nil && val > 0 {
				config.ModelContextLength = val
				if err := saveConfig(config); err != nil {
					fmt.Fprintf(os.Stderr, "Error saving config: %s\n", err)
				}
				fmt.Printf("Model context length set to: %d\n", config.ModelContextLength)
			} else {
				fmt.Println("Usage: /contextlength <positive_integer_value>")
			}
		} else {
			fmt.Println("Usage: /contextlength <value>")
		}
	case "/bg":
		if len(parts) < 2 {
			fmt.Println("Usage: /bg [list|view <pid>|kill <pid>]")
			return
		}
		switch parts[1] {
		case "list":
			fmt.Println(listBackgroundCommands())
		case "view", "logs":
			if len(parts) < 3 {
				fmt.Println("Usage: /bg view <pid>")
				return
			}
			pid, err := strconv.Atoi(parts[2])
			if err != nil {
				fmt.Println("Invalid pid. Usage: /bg view <pid>")
				return
			}
			logs, err := getBackgroundLogs(pid)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting logs: %v\n", err)
				return
			}
			if strings.TrimSpace(logs) == "" {
				fmt.Println("(no logs)")
				return
			}
			fmt.Println(logs)
		case "kill", "stop":
			if len(parts) < 3 {
				fmt.Println("Usage: /bg kill <pid>")
				return
			}
			pid, err := strconv.Atoi(parts[2])
			if err != nil {
				fmt.Println("Invalid pid. Usage: /bg kill <pid>")
				return
			}
			result, err := killBackgroundCommand(pid)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error killing process: %v\n", err)
				return
			}
			fmt.Println(result)
		default:
			fmt.Println("Usage: /bg [list|view <pid>|kill <pid>]")
		}
	case "/shell":
		shellMode = true
		fmt.Println("Entered shell mode. Type 'exit' to return.")
	case "/quit":
		if hasRunningBackgroundProcesses() {
			fmt.Println("Warning: You have running background processes.")
			fmt.Println(listBackgroundCommands())
			fmt.Print("Are you sure you want to quit? [y/N]: ")
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(strings.TrimSpace(response)) != "y" {
				return
			}
		}

		if agent != nil && len(agent.Messages) > 1 {
			if err := saveSession(agent); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to save session: %v\n", err)
			} else {
				fmt.Printf("Session '%s' saved.\n", agent.ID)
			}
		}
		fmt.Println("Bye!")
		os.Exit(0)
	case "/model":
		if len(parts) > 1 {
			config.Model = parts[1]
			if err := saveConfig(config); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving config: %s\n", err)
			}
			fmt.Printf("Model set to: %s\n", config.Model)
		} else {
			fmt.Println("Usage: /model <model_name>")
		}
	case "/provider":
		if len(parts) > 1 {
			config.APIURL = parts[1]
			if err := saveConfig(config); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving config: %s\n", err)
			}
			fmt.Printf("Provider URL set to: %s\n", config.APIURL)
		} else {
			fmt.Println("Usage: /provider <api_url>")
		}
	case "/config":
		fmt.Printf("Model: %s\n", config.Model)
		fmt.Printf("Provider: %s\n", config.APIURL)
		fmt.Printf("RAG Enabled: %t\n", config.RAGEnabled)
		fmt.Printf("RAG Path: %s\n", config.RAGPath)
		fmt.Printf("Operation Mode: %s\n", config.OperationMode)
		fmt.Printf("Execution Mode: %s\n", config.ExecutionMode)
		fmt.Printf("Auto Compress Enabled: %t\n", config.AutoCompress)
		fmt.Printf("Auto Compress Threshold: %d\n", config.AutoCompressThreshold)
		fmt.Printf("Model Context Length: %d\n", config.ModelContextLength)
		fmt.Printf("Stream Enabled: %t\n", config.Stream)
		fmt.Printf("Subagents Enabled: %t\n", config.SubagentsEnabled)
		if len(config.MCPs) > 0 {
			fmt.Println("MCP Servers:")
			for name, server := range config.MCPs {
				fmt.Printf("  - %s: %s\n", name, server.Command)
			}
		}
	case "/rag":
		if len(parts) > 1 {
			switch parts[1] {
			case "on":
				config.RAGEnabled = true
				if err := saveConfig(config); err != nil {
					fmt.Fprintf(os.Stderr, "Error saving config: %s\n", err)
				}
				fmt.Println("RAG enabled.")
			case "off":
				config.RAGEnabled = false
				if err := saveConfig(config); err != nil {
					fmt.Fprintf(os.Stderr, "Error saving config: %s\n", err)
				}
				fmt.Println("RAG disabled.")
			case "path":
				if len(parts) > 2 {
					config.RAGPath = parts[2]
					if err := saveConfig(config); err != nil {
						fmt.Fprintf(os.Stderr, "Error saving config: %s\n", err)
					}
					fmt.Printf("RAG path set to: %s\n", config.RAGPath)
				} else {
					fmt.Println("Usage: /rag path <path>")
				}
			default:
				fmt.Println("Usage: /rag [on|off|path <path>]")
			}
		} else {
			fmt.Println("Usage: /rag [on|off|path <path>]")
		}
	case "/compress":
		compressAndStartNewChat()
	case "/clear":
		fmt.Println("Clearing context (messages). This does NOT delete the saved session from disk.")
		if len(agent.Messages) > 1 {
			if err := saveSession(agent); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving session before clear: %v\n", err)
			} else {
				fmt.Printf("Session '%s' saved before clearing.\n", agent.ID)
			}
		}

		// Maintain the same ID so we are "in the same session" but cleared
		currentID := agent.ID
		agent = &Agent{
			ID:       currentID,
			Messages: make([]Message, 0),
		}
		systemPrompt := buildSystemPrompt("")
		agent.Messages = append(agent.Messages, Message{
			Role:    "system",
			Content: &systemPrompt,
		})
		totalTokens = 0
		fmt.Println("Context cleared.")
	case "/stream":
		if len(parts) > 1 {
			switch parts[1] {
			case "on":
				config.Stream = true
				if err := saveConfig(config); err != nil {
					fmt.Fprintf(os.Stderr, "Error saving config: %s\n", err)
				}
				fmt.Println("Streaming enabled.")
			case "off":
				config.Stream = false
				if err := saveConfig(config); err != nil {
					fmt.Fprintf(os.Stderr, "Error saving config: %s\n", err)
				}
				fmt.Println("Streaming disabled.")
			default:
				fmt.Println("Usage: /stream [on|off]")
			}
		} else {
			if config.Stream {
				fmt.Println("Streaming is currently enabled.")
			} else {
				fmt.Println("Streaming is currently disabled.")
			}
		}
	case "/subagents":
		if len(parts) > 1 {
			switch parts[1] {
			case "on":
				config.SubagentsEnabled = true
				if err := saveConfig(config); err != nil {
					fmt.Fprintf(os.Stderr, "Error saving config: %s\n", err)
				}
				fmt.Println("Sub-agent spawning enabled.")
			case "off":
				config.SubagentsEnabled = false
				if err := saveConfig(config); err != nil {
					fmt.Fprintf(os.Stderr, "Error saving config: %s\n", err)
				}
				fmt.Println("Sub-agent spawning disabled.")
			case "verbose":
				if len(parts) > 2 {
					mode, err := strconv.Atoi(parts[2])
					if err != nil || mode < 1 || mode > 2 {
						fmt.Println("Usage: /subagents verbose <1|2> (1: Default, 2: Full)")
					} else {
						config.SubAgentVerboseMode = mode
						if err := saveConfig(config); err != nil {
							fmt.Fprintf(os.Stderr, "Error saving config: %s\n", err)
						}
						fmt.Printf("Sub-agent verbose mode set to %d\n", mode)
					}
				} else {
					fmt.Printf("Current sub-agent verbose mode: %d\n", config.SubAgentVerboseMode)
					fmt.Println("Usage: /subagents verbose <1|2> (1: Default, 2: Full)")
				}
			default:
				fmt.Println("Usage: /subagents [on|off|verbose]")
			}
		} else {
			if config.SubagentsEnabled {
				fmt.Println("Sub-agent spawning is currently enabled.")
			} else {
				fmt.Println("Sub-agent spawning is currently disabled.")
			}
			fmt.Printf("Sub-agent verbose mode: %d\n", config.SubAgentVerboseMode)
		}
	case "/agent":
		if len(parts) < 2 {
			fmt.Println("Usage: /agent [studio|list|view <name>|use <name>|clear|rm <name>]")
			return
		}
		switch parts[1] {
		case "studio":
			// Optional: allow seeding studio with a one-line spec.
			spec := ""
			if len(parts) > 2 {
				spec = strings.Join(parts[2:], " ")
			}
			startAgentStudio(spec)
		case "list":
			fmt.Println(formatAgentsList())
		case "view":
			if len(parts) < 3 {
				fmt.Println("Usage: /agent view <name>")
				return
			}
			name := strings.Join(parts[2:], " ")
			fmt.Println(formatAgentView(name))
		case "use":
			if len(parts) < 3 {
				fmt.Println("Usage: /agent use <name>")
				return
			}
			name := strings.Join(parts[2:], " ")
			def, err := loadAgentDefinition(name)
			if err != nil {
				fmt.Printf("Error loading agent '%s': %v\n", name, err)
				return
			}

			// Snapshot current settings only when activating from "no active agent".
			if activeAgentDef == nil {
				prevAgentConfigSnapshot = &AgentConfigSnapshot{
					Model:     config.Model,
					Temp:      config.Temp,
					MaxTokens: config.MaxTokens,
				}
			}

			activeAgentDef = def

			// Apply optional overrides (in-memory only).
			if strings.TrimSpace(def.Model) != "" {
				config.Model = strings.TrimSpace(def.Model)
			}
			if def.Temperature != nil {
				config.Temp = *def.Temperature
			}
			if def.MaxTokens != nil {
				config.MaxTokens = *def.MaxTokens
			}

			// Clear context and rebuild system prompt so the new agent prompt takes effect.
			if len(agent.Messages) > 1 {
				if err := saveSession(agent); err != nil {
					fmt.Fprintf(os.Stderr, "Error saving session before agent switch: %v\n", err)
				} else {
					fmt.Printf("Session '%s' saved before switching agent.\n", agent.ID)
				}
			}
			currentID := agent.ID
			agent = &Agent{ID: currentID, Messages: make([]Message, 0)}
			systemPrompt := buildSystemPrompt("")
			agent.Messages = append(agent.Messages, Message{Role: "system", Content: &systemPrompt})
			totalTokens = 0
			fmt.Printf("Active agent set to '%s'. Context cleared.\n", activeAgentDef.Name)

		case "clear":
			if activeAgentDef == nil {
				fmt.Println("No active agent.")
				return
			}
			activeAgentDef = nil
			if prevAgentConfigSnapshot != nil {
				config.Model = prevAgentConfigSnapshot.Model
				config.Temp = prevAgentConfigSnapshot.Temp
				config.MaxTokens = prevAgentConfigSnapshot.MaxTokens
				prevAgentConfigSnapshot = nil
			}

			// Clear context and rebuild base system prompt.
			if len(agent.Messages) > 1 {
				if err := saveSession(agent); err != nil {
					fmt.Fprintf(os.Stderr, "Error saving session before clearing agent: %v\n", err)
				} else {
					fmt.Printf("Session '%s' saved before clearing agent.\n", agent.ID)
				}
			}
			currentID := agent.ID
			agent = &Agent{ID: currentID, Messages: make([]Message, 0)}
			systemPrompt := buildSystemPrompt("")
			agent.Messages = append(agent.Messages, Message{Role: "system", Content: &systemPrompt})
			totalTokens = 0
			fmt.Println("Active agent cleared. Context cleared.")

		case "rm":
			if len(parts) < 3 {
				fmt.Println("Usage: /agent rm <name>")
				return
			}
			name := strings.Join(parts[2:], " ")
			if err := deleteAgentDefinition(name); err != nil {
				fmt.Printf("Error deleting agent '%s': %v\n", name, err)
				return
			}
			// If the deleted agent is active, clear it.
			if activeAgentDef != nil && activeAgentDef.Name == name {
				activeAgentDef = nil
			}
			fmt.Printf("Agent '%s' deleted.\n", name)
		default:
			fmt.Println("Usage: /agent [studio|list|view <name>|use <name>|clear|rm <name>]")
		}
	default:
		fmt.Printf("Unknown command: %s\n", baseCommand)
	}
}

func buildSystemPrompt(contextSummary string) string {
	var basePrompt string
	if config.OperationMode == Plan {
		basePrompt = "You are an AI assistant in PLAN mode. Your goals are to:\n1. Analyze the user's request.\n2. Create a detailed implementation plan.\n3. Generate a comprehensive TODO list using the `create_todo` tool.\n4. Present the plan to the user using the `suggest_plan` tool for approval.\n\nIMPORTANT: You CANNOT execute shell commands in this mode. Focus purely on planning. Use the `suggest_plan` tool to show your plan and ask for confirmation. If the user approves (answers 'y' to the prompt), the system will automatically switch to 'build' mode for you to start implementation."
	} else {
		basePrompt = "You are an AI assistant in BUILD mode. You can execute commands, write code, and implement solutions. You can manage a todo list by using the `create_todo`, `update_todo`, and `get_todo_list` tools. You can also create notes using `create_note`, `update_note`, and `delete_note` tools. Notes persist across sessions. For multi-step tasks, chain commands with && (e.g., 'echo content > file.py && python3 file.py'). Use execute_command for shell tasks."
	}

	// If a task-specific agent is active, prepend its system prompt.
	if activeAgentDef != nil && strings.TrimSpace(activeAgentDef.SystemPrompt) != "" {
		basePrompt = fmt.Sprintf("=== Active Task-Specific Agent: %s ===\n%s\n\n%s", activeAgentDef.Name, activeAgentDef.SystemPrompt, basePrompt)
	}

	// Add compressed context if available
	if contextSummary != "" {
		basePrompt = fmt.Sprintf("Previous conversation context:\n\n%s\n\n%s", contextSummary, basePrompt)
	}

	// Add a list of available task-specific agents (helps the model choose agent names).
	agentsContent := getAgentsForSystemPrompt()
	if agentsContent != "" {
		basePrompt += agentsContent
	}

	// Add detailed MCP server and tool info to the prompt
	basePrompt += getMCPToolInfo()

	// Add notes to the prompt
	notesContent := getNotesForSystemPrompt()
	if notesContent != "" {
		basePrompt += notesContent
	}

	systemPrompt := getSystemInfo() + "\n\n" + basePrompt

	// Check for AGENTS.md and prepend its content to the system prompt
	agentInstructions, err := readAgentsFile("AGENTS.md")
	if err != nil {
		// Log the error but continue, as it's not a fatal issue.
		fmt.Fprintf(os.Stderr, "Warning: could not read AGENTS.md: %v\n", err)
	}

	if agentInstructions != "" {
		systemPrompt = agentInstructions + "\n\n" + systemPrompt
		fmt.Println("Found and loaded instructions from AGENTS.md")
	}

	return systemPrompt
}

func compressAndStartNewChat() {
	if len(agent.Messages) <= 1 {
		fmt.Println("No messages to compress. Start a conversation first.")
		return
	}

	fmt.Println("Compressing context...")
	compressedContent, err := compressContext(agent, config)
	if err != nil {
		fmt.Printf("Error compressing context: %s\n", err)
		return
	}

	// Create a new agent with the compressed context
	agent = &Agent{
		ID:       "main",
		Messages: make([]Message, 0),
	}

	systemPrompt := buildSystemPrompt(compressedContent)
	agent.Messages = append(agent.Messages, Message{
		Role:    "system",
		Content: &systemPrompt,
	})

	fmt.Println("Context compressed. Starting new chat with compressed summary as system message.")

	// Reset total tokens after compression
	totalTokens = 0
}

func runSetup() {
	fmt.Println("OpenAI API key is not set. Let's set it up.")
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your OpenAI API key: ")
	apiKey, _ := reader.ReadString('\n')
	config.APIKey = strings.TrimSpace(apiKey)

	if err := saveConfig(config); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to save config: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Configuration saved successfully.")
}

func runTask(task string) {
	// Load configuration
	config = loadConfig()
	if config.APIKey == "" {
		fmt.Fprintln(os.Stderr, "Error: API key not set. Please run the interactive setup first.")
		os.Exit(1)
	}

	// Create agent instance
	agent = &Agent{
		ID:       "main",
		Messages: make([]Message, 0),
	}

	systemPrompt := buildSystemPrompt("")
	agent.Messages = append(agent.Messages, Message{
		Role:    "system",
		Content: &systemPrompt,
	})

	// Add the task as a user message
	agent.Messages = append(agent.Messages, Message{
		Role:    "user",
		Content: &task,
	})

	// Execute the task using the agentic loop
	for {
		var resp *APIResponse
		var err error

		// Use streaming or regular API based on config
		if config.Stream {
			resp, err = sendAPIRequestStreaming(agent, config, config.SubagentsEnabled)
		} else {
			resp, err = sendAPIRequest(agent, config, config.SubagentsEnabled)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			break
		}

		if len(resp.Choices) == 0 {
			fmt.Fprintf(os.Stderr, "Error: received an empty response from the API\n")
			break
		}

		assistantMsg := resp.Choices[0].Message
		agent.Messages = append(agent.Messages, assistantMsg)

		// Only print content if not streaming (streaming already printed it)
		if !config.Stream && assistantMsg.Content != nil && *assistantMsg.Content != "" {
			fmt.Printf("%s\n", *assistantMsg.Content)
		}

		if len(assistantMsg.ToolCalls) > 0 {
			processToolCalls(agent, assistantMsg.ToolCalls, config)
		} else {
			// Check if we got an empty response
			if assistantMsg.Content == nil || *assistantMsg.Content == "" {
				fmt.Printf("%sWarning: Received empty response from model%s\n", ColorYellow, ColorReset)
			}
			break // No more tools to call, end agent turn
		}
		// Continue loop to send tool output back to API
	}
}

// editCommand creates a temporary file in /tmp, opens it in nano,
// and then adds its content to the prompt when the file is saved.
func editCommand() {
	// Create a temporary file in /tmp
	tmpFile, err := os.CreateTemp("/tmp", "agent-go-edit-*.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating temporary file: %s\n", err)
		return
	}
	tmpFileName := tmpFile.Name()
	tmpFile.Close()

	// Open the file in nano editor
	fmt.Printf("Opening %s in nano editor...\n", tmpFileName)
	fmt.Println("Make your changes and save the file (Ctrl+O, Enter, then Ctrl+X to exit).")
	cmd := exec.Command("nano", tmpFileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running nano: %s\n", err)
		// Clean up the temp file on error
		os.Remove(tmpFileName)
		return
	}

	// Read the content of the edited file
	content, err := os.ReadFile(tmpFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading temporary file: %s\n", err)
		os.Remove(tmpFileName)
		return
	}

	// Clean up the temporary file
	os.Remove(tmpFileName)

	// Only proceed if the file has content
	if len(content) == 0 {
		fmt.Println("No content in the file. No prompt added.")
		return
	}

	// Add the content to the agent's messages
	promptContent := string(content)
	agent.Messages = append(agent.Messages, Message{
		Role:    "user",
		Content: &promptContent,
	})

	// Print a success message
	fmt.Printf("Added content to prompt from %s\n", tmpFileName)

	// Now trigger the model response (agentic loop)
	// We'll call the same loop that processes user input
	for {
		// Auto-compress context if enabled and token count exceeds 75% of context length
		if config.AutoCompress && totalTokens > (config.ModelContextLength*3/4) {
			compressAndStartNewChat()
			fmt.Println("Context compressed due to token limit. New user input is required.")
			return // Exit after compression
		}

		var resp *APIResponse
		var err error

		// Use streaming or regular API based on config
		if config.Stream {
			resp, err = sendAPIRequestStreaming(agent, config, config.SubagentsEnabled)
		} else {
			resp, err = sendAPIRequest(agent, config, config.SubagentsEnabled)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			return // Break from the agentic loop
		}

		if len(resp.Choices) == 0 {
			fmt.Fprintf(os.Stderr, "Error: received an empty response from the API\n")
			return // Break from the agentic loop
		}

		assistantMsg := resp.Choices[0].Message
		agent.Messages = append(agent.Messages, assistantMsg)

		// Only print content if not streaming (streaming already printed it)
		if !config.Stream {
			if assistantMsg.ReasoningContent != nil && *assistantMsg.ReasoningContent != "" {
				fmt.Printf("%sThink...\n%s", ColorMeta, ColorReset)
			}
			if assistantMsg.Content != nil && *assistantMsg.Content != "" {
				fmt.Printf("%s%s%s\n", ColorMain, *assistantMsg.Content, ColorReset)
			}
		}

		// Update and display total tokens
		if resp.Usage.TotalTokens > 0 {
			totalTokens += resp.Usage.TotalTokens
			totalPromptTokens += resp.Usage.PromptTokens
			totalCompletionTokens += resp.Usage.CompletionTokens
		}

		// Update tool call count
		if len(assistantMsg.ToolCalls) > 0 {
			totalToolCalls += len(assistantMsg.ToolCalls)
		}

		// Display usage based on verbose mode
		switch config.UsageVerboseMode {
		case UsageSilent:
			// Do nothing
		case UsageDetailed:
			if resp.Usage.TotalTokens > 0 {
				fmt.Printf("%sUsage: %d prompt + %d completion = %d total tokens\n", ColorMeta, resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
				fmt.Printf("Total: %s tokens (%d prompt, %d completion), %d tool calls%s\n", formatTokenCount(totalTokens), totalPromptTokens, totalCompletionTokens, totalToolCalls, ColorReset)
			}
		case UsageBasic:
			fallthrough
		default:
			// Default behavior (Basic)
			if resp.Usage.TotalTokens > 0 {
				fmt.Printf("%sUsed %s%s%s tokens on %s%s\n", ColorMeta, ColorHighlight, formatTokenCount(totalTokens), ColorMeta, config.Model, ColorReset)
			}
		}

		if len(assistantMsg.ToolCalls) > 0 {
			processToolCalls(agent, assistantMsg.ToolCalls, config)
		} else {
			// Check if we got an empty response
			if assistantMsg.Content == nil || *assistantMsg.Content == "" {
				fmt.Printf("%sWarning: Received empty response from model%s\n", ColorMeta, ColorReset)
			}
			break // No more tools to call, end agent turn
		}
		// Continue loop to send tool output back to API

		// Inject reminder if background processes are running
		if hasRunningBackgroundProcesses() {
			reminder := fmt.Sprintf("REMINDER: You have running background processes:\n%s", listBackgroundCommands())
			agent.Messages = append(agent.Messages, Message{
				Role:    "system",
				Content: &reminder,
			})
			// We don't print this to the user, it's just for the agent's context
		}
	}
}
