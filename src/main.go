package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/chzyer/readline"
)

var config *Config
var agent *Agent

func main() {
	config = loadConfig()
	if config.APIKey == "" {
		runSetup()
	}

	agent = &Agent{
		Messages: make([]Message, 0, 20),
	}
	agent.Messages = append(agent.Messages, Message{
		Role:    "system",
		Content: stringp("You are an AI assistant. For multi-step tasks, chain commands with && (e.g., 'echo content > file.py && python3 file.py'). Use execute_command for shell tasks."),
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

func runCLI() {
	home, _ := os.UserHomeDir()
	historyFile := filepath.Join(home, ".config", "agent-go", "history.txt")

	rl, err := readline.NewEx(&readline.Config{
		Prompt:      "> ",
		HistoryFile: historyFile,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	fmt.Println("Agent-Go is ready. Type your requests, or /help for a list of commands.")

	for {
		userInput, err := rl.Readline()
		if err == readline.ErrInterrupt {
			continue
		} else if err == io.EOF {
			break
		}

		userInput = strings.TrimSpace(userInput)
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

		// Manage message history to stay within limits
		if len(agent.Messages) > 20 {
			agent.Messages = append(agent.Messages[:1], agent.Messages[len(agent.Messages)-19:]...)
		}

		// Agentic loop
		for {
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
		fmt.Println("  /quit              - Exit the application")
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
	default:
		fmt.Printf("Unknown command: %s\n", baseCommand)
	}
}

func stringp(s string) *string {
	return &s
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