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

// getAvailableTools returns the list of tools available to the agent
func getAvailableTools(includeSpawn bool) []Tool {
	tools := []Tool{
		{
			Type: "function",
			Function: FunctionDefinition{
				Name:        "execute_command",
				Description: "Execute shell command",
				Parameters: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{"command": map[string]string{"type": "string"}},
					"required":   []string{"command"},
				},
			},
		},
		{
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
		},
		{
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
		},
		{
			Type: "function",
			Function: FunctionDefinition{
				Name:        "get_todo_list",
				Description: "Get the current list of todo items.",
				Parameters:  map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			},
		},
	}

	if includeSpawn {
		tools = append(tools, Tool{
			Type: "function",
			Function: FunctionDefinition{
				Name:        "spawn_agent",
				Description: "Spawn a sub-agent to perform a specific task and return the result.",
				Parameters: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{"task": map[string]string{"type": "string"}},
					"required":   []string{"task"},
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
		builder.WriteString(fmt.Sprintf("- [ID: %d] %s (%s)\n", todo.ID, todo.Task, todo.Status))
	}

	return builder.String(), nil
}
