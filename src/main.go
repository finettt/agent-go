package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
var totalTokens = 0

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
		fmt.Println("\nBye!")
		os.Exit(0)
	}()

	runCLI()
}
func printLogo() {
	fmt.Print(`
  /$$$$$$                                  /$$            /$$$$$$
 /$$__  $$                                | $$           /$$__  $$
| $$  \ $$  /$$$$$$   /$$$$$$  /$$$$$$$  /$$$$$$        | $$  \__/  /$$$$$$
| $$$$$$$$ /$$__  $$ /$$__  $$| $$__  $$|_  $$_/        | $$ /$$$$ /$$__  $$
| $$__  $$| $$  \ $$| $$$$$$$$| $$  \ $$  | $$          | $$|_  $$| $$  \ $$
| $$  | $$| $$  | $$| $$_____/| $$  | $$  | $$ /$$      | $$  \ $$| $$  | $$
| $$  | $$|  $$$$$$$|  $$$$$$$| $$  | $$  |  $$$$/      |  $$$$$$/|  $$$$$$/
|__/  |__/ \____  $$ \_______/|__/  |__/   \___/         \______/  \______/
           /$$  \ $$
          |  $$$$$$/
           \______/
`)
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

	fmt.Println("Agent-Go is ready. Type your requests, or /help for a list of commands.")

	for {
		if shellMode {
			rl.SetPrompt("shell> ")
		} else {
			rl.SetPrompt("> ")
		}

		userInput, err := rl.Readline()
		if err == readline.ErrInterrupt {
			continue
		} else if err == io.EOF {
			break
		}

		userInput = strings.TrimSpace(userInput)

		if shellMode {
			if userInput == "exit" {
				shellMode = false
				fmt.Println("Exited shell mode.")
				continue
			}
			if userInput == "" {
				continue
			}
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
			if !config.Stream && assistantMsg.Content != nil {
				fmt.Printf("%s%s%s\n", ColorBlue, *assistantMsg.Content, ColorReset)
			}

			// Update and display total tokens
			if resp.Usage.TotalTokens > 0 {
				totalTokens += resp.Usage.TotalTokens
				fmt.Printf("Used %s%d%s tokens on %s\n", ColorGreen, totalTokens, ColorReset, config.Model)
			}

			if len(assistantMsg.ToolCalls) > 0 {
				processToolCalls(agent, assistantMsg.ToolCalls, config)
			} else {
				break // No more tools to call, end agent turn
			}
			// Continue loop to send tool output back to API
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
	case "/help":
		fmt.Println("Available commands:")
		fmt.Println("  /help              - Show this help message")
		fmt.Println("  /model <name>      - Set the AI model (e.g., gpt-4)")
		fmt.Println("  /provider <url>    - Set the API provider URL")
		fmt.Println("  /config            - Display current configuration")
		fmt.Println("  /rag on|off        - Toggle RAG feature")
		fmt.Println("  /rag path <path>   - Set the RAG documents path")
		fmt.Println("  /shell             - Enter shell mode for direct command execution")
		fmt.Println("  /compress          - Compress context and start new chat thread")
		fmt.Println("  /contextlength <value> - Set the model context length (e.g., 131072)")
		fmt.Println("  /stream on|off     - Toggle streaming mode")
		fmt.Println("  /subagents on|off  - Toggle sub-agent spawning")
		fmt.Println("  /todo              - Display the current todo list")
		fmt.Println("  /mcp add <name> <command> - Add an MCP server")
		fmt.Println("  /mcp remove <name> - Remove an MCP server")
		fmt.Println("  /mcp list          - List MCP servers")
		fmt.Println("  /mode              - Switch between ASK and YOLO mode")
		fmt.Println("  /quit              - Exit the application")
	case "/mode":
		if config.ExecutionMode == Ask {
			config.ExecutionMode = YOLO
			fmt.Println("Switched to YOLO mode.")
		} else {
			config.ExecutionMode = Ask
			fmt.Println("Switched to ASK mode.")
		}
		if err := saveConfig(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %s\n", err)
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
	case "/shell":
		shellMode = true
		fmt.Println("Entered shell mode. Type 'exit' to return.")
	case "/quit":
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
			default:
				fmt.Println("Usage: /subagents [on|off]")
			}
		} else {
			if config.SubagentsEnabled {
				fmt.Println("Sub-agent spawning is currently enabled.")
			} else {
				fmt.Println("Sub-agent spawning is currently disabled.")
			}
		}
	default:
		fmt.Printf("Unknown command: %s\n", baseCommand)
	}
}

func buildSystemPrompt(contextSummary string) string {
	// Base system prompt
	basePrompt := "You are an AI assistant. You can manage a todo list by using the `create_todo`, `update_todo`, and `get_todo_list` tools. For multi-step tasks, chain commands with && (e.g., 'echo content > file.py && python3 file.py'). Use execute_command for shell tasks."

	// Add compressed context if available
	if contextSummary != "" {
		basePrompt = fmt.Sprintf("Previous conversation context:\n\n%s\n\n%s", contextSummary, basePrompt)
	}

	// Add detailed MCP server and tool info to the prompt
	basePrompt += getMCPToolInfo()

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
		if !config.Stream && assistantMsg.Content != nil {
			fmt.Printf("%s\n", *assistantMsg.Content)
		}

		if len(assistantMsg.ToolCalls) > 0 {
			processToolCalls(agent, assistantMsg.ToolCalls, config)
		} else {
			break // No more tools to call, end agent turn
		}
		// Continue loop to send tool output back to API
	}
}
