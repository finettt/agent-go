package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type TodoItem struct {
	ID     int    `json:"id"`
	Task   string `json:"task"`
	Status string `json:"status"` // "pending", "in-progress", "completed"
}

type TodoList struct {
	AgentID string     `json:"agent_id"`
	Todos   []TodoItem `json:"todos"`
	NextID  int        `json:"next_id"`
}

func getTodoListPath(agentID string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "agent-go", "todos", fmt.Sprintf("%s.json", agentID)), nil
}

func loadTodoList(agentID string) (*TodoList, error) {
	path, err := getTodoListPath(agentID)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &TodoList{AgentID: agentID, Todos: []TodoItem{}, NextID: 1}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var todoList TodoList
	if err := json.Unmarshal(data, &todoList); err != nil {
		return nil, err
	}

	return &todoList, nil
}

func saveTodoList(todoList *TodoList) error {
	path, err := getTodoListPath(todoList.AgentID)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(todoList, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func getCurrentTask(agentID string) (string, error) {
	todoList, err := loadTodoList(agentID)
	if err != nil {
		return "", err
	}

	for _, todo := range todoList.Todos {
		if todo.Status == "in-progress" {
			return fmt.Sprintf("%s", todo.Task), nil
		}
	}

	return "No task in progress.", nil
}
