package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var config *Config
var agent *Agent

func main() {
	config = loadConfig()
	if config.APIKey == "" {
		fmt.Fprintln(os.Stderr, "OPENAI_KEY required")
		os.Exit(1)
	}

	agent = &Agent{
		Messages: make([]Message, 0, 20),
	}
	agent.Messages = append(agent.Messages, Message{
		Role:    "system",
		Content: "You are an AI assistant. For multi-step tasks, chain commands with && (e.g., 'echo content > file.py && python3 file.py'). Use execute_command for shell tasks.",
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
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Agent-Go is ready. Type your requests.")
	for {
		fmt.Print("> ")
		userInput, _ := reader.ReadString('\n')
		if userInput == "\n" {
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
		agent.Messages = append(agent.Messages, Message{Role: "user", Content: userInput})

		// Manage message history to stay within limits
		if len(agent.Messages) > 20 {
			agent.Messages = append(agent.Messages[:1], agent.Messages[len(agent.Messages)-19:]...)
		}

		resp, err := sendAPIRequest(agent, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			continue
		}

		if len(resp.Choices) > 0 {
			msg := resp.Choices[0].Message
			if msg.Content != nil {
				fmt.Printf("\033[34m%s\033[0m\n", *msg.Content)
				agent.Messages = append(agent.Messages, Message{Role: "assistant", Content: *msg.Content})
			}

			if len(msg.ToolCalls) > 0 {
				// For simplicity, handle one tool call at a time
				toolCall := msg.ToolCalls[0]
				if toolCall.Function.Name == "execute_command" {
					var args CommandArgs
					if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err == nil {
						output, err := executeCommand(args.Command)
						if err != nil {
							fmt.Fprintf(os.Stderr, "Command execution error: %s\n", err)
						}
						// Add tool output to message history and resend
						agent.Messages = append(agent.Messages, Message{Role: "tool", Content: "Command output:\n" + output})
						// This will require another call to sendAPIRequest, creating a loop.
						// For this iteration, we will just print the output.
						// A more robust implementation would loop here.
						fmt.Println("Tool output:", output)
					}
				}
			}
		}
	}
}