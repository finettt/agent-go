package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// CreateTodoArgs represents arguments for creating a todo item
type CreateTodoArgs struct {
	Task string `json:"task"`
}

// UpdateTodoArgs represents arguments for updating a todo item
type UpdateTodoArgs struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type SuggestPlanArgs struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// getAvailableTools returns the list of tools available to the agent
func getAvailableTools(config *Config, includeSpawn bool, operationMode OperationMode) []Tool {
	tools := []Tool{}

	// Add custom skills
	for _, skill := range config.Skills {
		tools = append(tools, Tool{
			Type: "function",
			Function: FunctionDefinition{
				Name:        skill.Name,
				Description: skill.Description,
				Parameters:  skill.Parameters,
			},
		})
	}

	if operationMode == Build {
		tools = append(tools, Tool{
			Type: "function",
			Function: FunctionDefinition{
				Name:        "execute_command",
				Description: "Execute shell command (foreground). In Ask mode the user may choose to run it in the background at approval time.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"command": map[string]string{"type": "string"},
					},
					"required": []string{"command"},
				},
			},
		})
	}

	tools = append(tools, Tool{
		Type: "function",
		Function: FunctionDefinition{
			Name:        "create_todo",
			Description: "Create a new todo item.",
			Parameters: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{"task": map[string]string{"type": "string"}},
				"required":   []string{"task"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: FunctionDefinition{
			Name:        "update_todo",
			Description: "Update a todo item's status.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id":     map[string]interface{}{"type": "integer"},
					"status": map[string]interface{}{"type": "string", "enum": []string{"pending", "in-progress", "completed"}},
				},
				"required": []string{"id", "status"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: FunctionDefinition{
			Name:        "get_todo_list",
			Description: "Get the current list of todo items.",
			Parameters:  map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: FunctionDefinition{
			Name:        "get_current_task",
			Description: "Get the current in-progress task.",
			Parameters:  map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: FunctionDefinition{
			Name:        "clear_todo",
			Description: "Clear all todo items.",
			Parameters:  map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: FunctionDefinition{
			Name:        "create_note",
			Description: "Create a note. Notes persist across sessions and are injected into the system prompt.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name":    map[string]string{"type": "string", "description": "Unique name for the note"},
					"content": map[string]string{"type": "string", "description": "Content of the note"},
				},
				"required": []string{"name", "content"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: FunctionDefinition{
			Name:        "update_note",
			Description: "Update an existing note's content.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name":    map[string]string{"type": "string", "description": "Name of the note to update"},
					"content": map[string]string{"type": "string", "description": "New content for the note"},
				},
				"required": []string{"name", "content"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: FunctionDefinition{
			Name:        "delete_note",
			Description: "Delete a note.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]string{"type": "string", "description": "Name of the note to delete"},
				},
				"required": []string{"name"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: FunctionDefinition{
			Name:        "name_session",
			Description: "Give the current session a name. This helps organize and restore sessions later.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]string{"type": "string", "description": "New name for the session (use dashes for spaces, e.g., 'implement-auth-feature')"},
				},
				"required": []string{"name"},
			},
		},
	})

	// Add background command tools only in Build mode
	if operationMode == Build {
		tools = append(tools, Tool{
			Type: "function",
			Function: FunctionDefinition{
				Name:        "kill_background_command",
				Description: "Kill a running background command by PID.",
				Parameters: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{"pid": map[string]interface{}{"type": "integer"}},
					"required":   []string{"pid"},
				},
			},
		})
		tools = append(tools, Tool{
			Type: "function",
			Function: FunctionDefinition{
				Name:        "get_background_logs",
				Description: "Get the logs (stdout/stderr) of a background command by PID.",
				Parameters: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{"pid": map[string]interface{}{"type": "integer"}},
					"required":   []string{"pid"},
				},
			},
		})
		tools = append(tools, Tool{
			Type: "function",
			Function: FunctionDefinition{
				Name:        "list_background_commands",
				Description: "List all running background commands.",
				Parameters:  map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			},
		})
	}

	// Add suggest_plan tool only in Plan mode
	if operationMode == Plan {
		tools = append(tools, Tool{
			Type: "function",
			Function: FunctionDefinition{
				Name:        "suggest_plan",
				Description: "Suggest a plan to the user and ask for approval.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name":        map[string]string{"type": "string", "description": "A short, descriptive name for the plan"},
						"description": map[string]string{"type": "string", "description": "A high-level description of what the plan entails"},
					},
					"required": []string{"name", "description"},
				},
			},
		})

		// Agent Studio: allow creating task-specific agents (persisted on disk).
		tools = append(tools, Tool{
			Type: "function",
			Function: FunctionDefinition{
				Name:        "create_agent_definition",
				Description: "Create a new task-specific agent definition (name + system prompt) and save it for later use via /agent commands.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name":          map[string]string{"type": "string", "description": "Unique agent name (short, filesystem-friendly; spaces become dashes)."},
						"description":   map[string]string{"type": "string", "description": "Optional short description of what the agent does."},
						"system_prompt": map[string]string{"type": "string", "description": "The full system prompt for the agent."},
						"model":         map[string]string{"type": "string", "description": "Optional model override for this agent."},
						"temperature":   map[string]string{"type": "number", "description": "Optional temperature override (0.0-2.0)."},
						"max_tokens":    map[string]string{"type": "integer", "description": "Optional max tokens override."},
					},
					"required": []string{"name", "system_prompt"},
				},
			},
		})
	}

	// Add generic MCP tool
	tools = append(tools, Tool{
		Type: "function",
		Function: FunctionDefinition{
			Name:        "use_mcp_tool",
			Description: "Call a tool on a connected MCP server.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"server_name": map[string]string{"type": "string"},
					"tool_name":   map[string]string{"type": "string"},
					"arguments":   map[string]interface{}{"type": "object"},
				},
				"required": []string{"server_name", "tool_name", "arguments"},
			},
		},
	})

	if includeSpawn {
		tools = append(tools, Tool{
			Type: "function",
			Function: FunctionDefinition{
				Name:        "spawn_agent",
				Description: "Spawn a sub-agent to perform a specific task and return the result. Optionally choose a task-specific agent definition (including built-in 'default').",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"task":  map[string]string{"type": "string"},
						"agent": map[string]string{"type": "string", "description": "Optional agent name (e.g. 'default' or a saved agent) to use as the sub-agent's system prompt."},
					},
					"required": []string{"task"},
				},
			},
		})
	}

	return tools
}

// validateTodoStatus validates if a status is valid
func validateTodoStatus(status string) error {
	if !ValidTodoStatuses[status] {
		return fmt.Errorf("invalid status: %s (must be: pending, in-progress, completed)", status)
	}
	return nil
}

// createTodo creates a new todo item for the given agent
func createTodo(agentID, argsJSON string) (string, error) {
	var args CreateTodoArgs
	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}

	// Validate task is not empty
	if strings.TrimSpace(args.Task) == "" {
		return "", fmt.Errorf("task cannot be empty")
	}

	todoList, err := loadTodoList(agentID)
	if err != nil {
		return "", err
	}

	newTodo := TodoItem{
		ID:     todoList.NextID,
		Task:   args.Task,
		Status: "pending",
	}
	todoList.NextID++
	todoList.Todos = append(todoList.Todos, newTodo)

	if err := saveTodoList(todoList); err != nil {
		return "", err
	}

	return getTodoList(agentID)
}

// updateTodo updates the status of an existing todo item
func updateTodo(agentID, argsJSON string) (string, error) {
	var args UpdateTodoArgs
	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}

	// Validate status
	if err := validateTodoStatus(args.Status); err != nil {
		return "", err
	}

	todoList, err := loadTodoList(agentID)
	if err != nil {
		return "", err
	}

	var found bool
	for i, todo := range todoList.Todos {
		if todo.ID == args.ID {
			todoList.Todos[i].Status = args.Status
			found = true
			break
		}
	}

	if !found {
		return "", fmt.Errorf("todo item with ID %d not found", args.ID)
	}

	if err := saveTodoList(todoList); err != nil {
		return "", err
	}

	return getTodoList(agentID)
}

// getTodoList returns the current todo list for the given agent
func getTodoList(agentID string) (string, error) {
	todoList, err := loadTodoList(agentID)
	if err != nil {
		return "", err
	}

	if len(todoList.Todos) == 0 {
		return "Todo list is empty.", nil
	}

	var builder strings.Builder
	builder.WriteString("Current Todo List:\n")
	for _, todo := range todoList.Todos {
		checkbox := " "
		switch todo.Status {
		case "completed":
			checkbox = "x"
		case "in-progress":
			checkbox = "-"
		}

		builder.WriteString(fmt.Sprintf("%d. [%s] %s \n", todo.ID, checkbox, todo.Task))
	}

	return builder.String(), nil
}

// clearTodo clears the todo list for the given agent
func clearTodo(agentID string) (string, error) {
	if err := clearTodoList(agentID); err != nil {
		return "", err
	}
	return "Todo list cleared.", nil
}
