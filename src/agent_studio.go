package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// agentStudioMode toggles a dedicated interactive flow to create task-specific agents.
var agentStudioMode = false

// studioAgent holds the in-progress “studio chat” conversation state.
var studioAgent *Agent

// activeAgentDef is the currently selected task-specific agent definition used to build the system prompt.
var activeAgentDef *AgentDefinition

func buildAgentStudioSystemPrompt() string {
	// We deliberately do NOT include getMCPToolInfo() here to avoid costly MCP connects
	// and to keep the studio focused. The studio tool processing also rejects tools other than
	// create_agent_definition.
	return strings.TrimSpace(`
You are Agent Studio.

Goal:
- Help the user create a task-specific agent (a reusable assistant configuration).

Rules:
- Ask concise clarifying questions until you have enough info.
- When ready, call the tool create_agent_definition EXACTLY ONCE with valid JSON arguments.
- Do NOT call execute_command, spawn_agent, use_mcp_tool, create_todo, update_todo, get_todo_list, create_note, update_note, delete_note, name_session, or suggest_plan.
- After calling create_agent_definition, provide a short confirmation summary (name + what it does).

Tool: create_agent_definition
Arguments JSON schema:
{
  "name": "string (required, short, filesystem-friendly)",
  "description": "string (optional)",
  "system_prompt": "string (required, the full system prompt for the new agent)",
  "model": "string (optional, e.g. gpt-4.1)",
  "temperature": 0.0-2.0 (optional),
  "max_tokens": integer (optional)
}

System prompt guidance:
- Be explicit about scope, inputs/outputs, constraints, safety, style, and what tools it should/shouldn't use.
- Keep it practical and runnable for this CLI environment.
`)
}

// startAgentStudio enables studio mode and optionally seeds it with an initial spec.
func startAgentStudio(initialSpec string) {
	agentStudioMode = true

	systemPrompt := buildAgentStudioSystemPrompt()
	studioAgent = &Agent{
		ID:       "studio",
		Messages: make([]Message, 0),
	}
	studioAgent.Messages = append(studioAgent.Messages, Message{
		Role:    "system",
		Content: &systemPrompt,
	})

	if strings.TrimSpace(initialSpec) != "" {
		spec := strings.TrimSpace(initialSpec)
		studioAgent.Messages = append(studioAgent.Messages, Message{
			Role:    "user",
			Content: &spec,
		})
	}
	fmt.Println("Agent Studio started. Type your requirements. Type 'exit' to leave without creating an agent.")
}

// runAgentStudioTurn sends the current studio conversation to the LLM and processes tool calls.
// It is safe by design: it only executes create_agent_definition; all other tools are rejected.
func runAgentStudioTurn(cfg *Config) error {
	if studioAgent == nil {
		return fmt.Errorf("studio is not initialized")
	}

	// Use a config copy to ensure execute_command is not even offered as an available tool.
	// Plan mode tools may still be offered; we explicitly reject them in tool processing.
	cfgCopy := *cfg
	cfgCopy.OperationMode = Plan
	cfgCopy.Stream = false

	for {
		resp, err := sendAPIRequest(studioAgent, &cfgCopy, false)
		if err != nil {
			return err
		}
		if len(resp.Choices) == 0 {
			return fmt.Errorf("received an empty response from the API")
		}

		assistantMsg := resp.Choices[0].Message
		studioAgent.Messages = append(studioAgent.Messages, assistantMsg)

		if assistantMsg.Content != nil && strings.TrimSpace(*assistantMsg.Content) != "" {
			fmt.Printf("%s%s%s\n", ColorMain, *assistantMsg.Content, ColorReset)
		}

		if len(assistantMsg.ToolCalls) == 0 {
			break
		}

		if err := processAgentStudioToolCalls(studioAgent, assistantMsg.ToolCalls); err != nil {
			return err
		}
		// Continue loop to send tool outputs back to the API.
	}

	return nil
}

func processAgentStudioToolCalls(a *Agent, toolCalls []ToolCall) error {
	for _, tc := range toolCalls {
		var output string
		var err error

		switch tc.Function.Name {
		case "create_agent_definition":
			// Validate JSON early to give clear feedback to the model if it emits junk.
			var tmp map[string]interface{}
			if jerr := json.Unmarshal([]byte(tc.Function.Arguments), &tmp); jerr != nil {
				output = fmt.Sprintf("Failed to parse arguments for create_agent_definition: %v", jerr)
				err = nil
				break
			}

			output, err = createAgentDefinition(tc.Function.Arguments)
			if err == nil {
				// Studio goal achieved; exit studio mode.
				agentStudioMode = false
			}

		default:
			// Hard reject all other tools inside studio to prevent side effects.
			output = fmt.Sprintf("Tool '%s' is not allowed in Agent Studio. Only 'create_agent_definition' is permitted.", tc.Function.Name)
			err = nil
		}

		if err != nil {
			output = fmt.Sprintf("Tool execution error: %v", err)
			fmt.Printf("%s%s%s\n", ColorRed, output, ColorReset)
		}

		toolMsg := Message{
			Role:       "tool",
			ToolCallID: tc.ID,
			Content:    &output,
		}
		a.Messages = append(a.Messages, toolMsg)

		// Print tool output for the user (keeps studio transparent).
		fmt.Printf("%s%s%s\n", ColorMeta, output, ColorReset)
	}

	return nil
}
