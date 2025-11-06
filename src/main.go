package main

import (
	"bufio"
	"encoding/json"
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
		Messages: make([]Message, 0),
	}

	// Base system prompt
	systemPrompt := "You are an AI assistant. For multi-step tasks, chain commands with && (e.g., 'echo content > file.py && python3 file.py'). Use execute_command for shell tasks."

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

	agent.Messages = append(agent.Messages, Message{
		Role:    "system",
		Content: stringp(systemPrompt),
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
	fmt.Println(`
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
	defer rl.Close()

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
			output, err := executeCommand(userInput)
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
		agent.Messages = append(agent.Messages, Message{Role: "user", Content: stringp(userInput)})

		// Message history is now unlimited

		// Agentic loop
		for {
			// Auto-compress context if enabled and token count exceeds 75% of context length
			if config.AutoCompress && totalTokens > (config.ModelContextLength*3/4) {
				compressAndStartNewChat()
				fmt.Println("Context compressed due to token limit. New user input is required.")
				continue // Restart the outer loop to get fresh user input
			}

			resp, err := sendAPIRequest(agent, config)
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

			if assistantMsg.Content != nil {
				fmt.Printf("\033[34m%s\033[0m\n", *assistantMsg.Content)
			}

			// Update and display total tokens
			if resp.Usage.TotalTokens > 0 {
				totalTokens += resp.Usage.TotalTokens
				fmt.Printf("\033[32m[Tokens: %d]\033[0m\n", totalTokens)
			}

			if len(assistantMsg.ToolCalls) == 0 {
				break // No tool calls, so the agent's turn is over
			}

			for _, toolCall := range assistantMsg.ToolCalls {
				if toolCall.Function.Name == "execute_command" {
					var args CommandArgs
					if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
						fmt.Fprintf(os.Stderr, "Tool call argument error: %s\n", err)
						continue // Next tool call
					}

					output, err := executeCommand(args.Command)
					if err != nil {
						output = fmt.Sprintf("Command execution error: %s", err)
					}

					content := "Command output:\n" + output
					toolMsg := Message{
						Role:       "tool",
						ToolCallID: toolCall.ID,
						Content:    stringp(content),
					}
					agent.Messages = append(agent.Messages, toolMsg)
				}
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
			// File not found is not an error in this context.
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
		fmt.Println("  /quit              - Exit the application")
	case "/contextlength":
		if len(parts) > 1 {
			val, err := strconv.Atoi(parts[1])
			if err == nil && val > 0 {
				config.ModelContextLength = val
				saveConfig(config)
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
			saveConfig(config)
			fmt.Printf("Model set to: %s\n", config.Model)
		} else {
			fmt.Println("Usage: /model <model_name>")
		}
	case "/provider":
		if len(parts) > 1 {
			config.APIURL = parts[1]
			saveConfig(config)
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
	case "/rag":
		if len(parts) > 1 {
			switch parts[1] {
			case "on":
				config.RAGEnabled = true
				saveConfig(config)
				fmt.Println("RAG enabled.")
			case "off":
				config.RAGEnabled = false
				saveConfig(config)
				fmt.Println("RAG disabled.")
			case "path":
				if len(parts) > 2 {
					config.RAGPath = parts[2]
					saveConfig(config)
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
	default:
		fmt.Printf("Unknown command: %s\n", baseCommand)
	}
}

func stringp(s string) *string {
	return &s
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

	agent = &Agent{
		Messages: make([]Message, 0),
	}

	systemPrompt := fmt.Sprintf("Previous conversation context:\n\n%s\n\nYou are an AI assistant. For multi-step tasks, chain commands with && (e.g., 'echo content > file.py && python3 file.py'). Use execute_command for shell tasks.", compressedContent)

	agentInstructions, err := readAgentsFile("AGENTS.md")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not read AGENTS.md: %v\n", err)
	}

	if agentInstructions != "" {
		systemPrompt = agentInstructions + "\n\n" + systemPrompt
	}

	agent.Messages = append(agent.Messages, Message{
		Role:    "system",
		Content: stringp(systemPrompt),
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
		Messages: make([]Message, 0),
	}

	// Base system prompt
	systemPrompt := "You are an AI assistant. For multi-step tasks, chain commands with && (e.g., 'echo content > file.py && python3 file.py'). Use execute_command for shell tasks."

	// Check for AGENTS.md and prepend its content to the system prompt
	agentInstructions, err := readAgentsFile("AGENTS.md")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not read AGENTS.md: %v\n", err)
	}

	if agentInstructions != "" {
		systemPrompt = agentInstructions + "\n\n" + systemPrompt
	}

	agent.Messages = append(agent.Messages, Message{
		Role:    "system",
		Content: stringp(systemPrompt),
	})

	// Add the task as a user message
	agent.Messages = append(agent.Messages, Message{
		Role:    "user",
		Content: stringp(task),
	})

	// Execute the task using the agentic loop
	for {
		resp, err := sendAPIRequest(agent, config)
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

		if assistantMsg.Content != nil {
			fmt.Printf("%s\n", *assistantMsg.Content)
		}

		if len(assistantMsg.ToolCalls) == 0 {
			break // No tool calls, so the agent's turn is over
		}

		for _, toolCall := range assistantMsg.ToolCalls {
			if toolCall.Function.Name == "execute_command" {
				var args CommandArgs
				if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
					fmt.Fprintf(os.Stderr, "Tool call argument error: %s\n", err)
					continue // Next tool call
				}

				output, err := executeCommand(args.Command)
				if err != nil {
					output = fmt.Sprintf("Command execution error: %s", err)
				}

				content := "Command output:\n" + output
				toolMsg := Message{
					Role:       "tool",
					ToolCallID: toolCall.ID,
					Content:    stringp(content),
				}
				agent.Messages = append(agent.Messages, toolMsg)
			}
		}
		// Continue loop to send tool output back to API
	}
}
