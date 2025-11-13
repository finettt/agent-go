package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type CreateTodoArgs struct {
	Task string `json:"task"`
}

type UpdateTodoArgs struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

func createTodo(agentID, argsJSON string) (string, error) {
	var args CreateTodoArgs
	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
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

func updateTodo(agentID, argsJSON string) (string, error) {
	var args UpdateTodoArgs
	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
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
