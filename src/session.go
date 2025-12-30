package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Session represents a saved agent session
type Session struct {
	ID               string    `json:"id"`
	Messages         []Message `json:"messages"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	TotalTokens      int       `json:"total_tokens"`
	PromptTokens     int       `json:"prompt_tokens"`
	CompletionTokens int       `json:"completion_tokens"`
	ToolCalls        int       `json:"tool_calls"`
}

// getSessionsDir returns the path to the sessions directory
func getSessionsDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "agent-go", "sessions")
}

// ensureSessionsDir creates the sessions directory if it doesn't exist
func ensureSessionsDir() error {
	sessionsDir := getSessionsDir()
	return os.MkdirAll(sessionsDir, 0755)
}

// getSessionPath returns the path to a specific session file
func getSessionPath(id string) string {
	// Sanitize the id to be a valid filename
	safeID := strings.ReplaceAll(id, "/", "_")
	safeID = strings.ReplaceAll(safeID, "\\", "_")
	safeID = strings.ReplaceAll(safeID, "..", "_")
	return filepath.Join(getSessionsDir(), safeID+".json")
}

// saveSession saves the current agent state as a session
func saveSession(agent *Agent) error {
	if err := ensureSessionsDir(); err != nil {
		return err
	}

	// If the agent ID is "main" (default), we might want to generate a timestamp-based ID
	// or keep it as "main" if the user hasn't renamed it.
	// However, the requirement is "give each session a name (by llm tool name_session)".
	// If the user uses the tool, agent.ID will be updated.

	session := Session{
		ID:               agent.ID,
		Messages:         agent.Messages,
		UpdatedAt:        time.Now(),
		TotalTokens:      totalTokens,
		PromptTokens:     totalPromptTokens,
		CompletionTokens: totalCompletionTokens,
		ToolCalls:        totalToolCalls,
	}
	// Try to load existing session to preserve CreatedAt
	existing, err := loadSession(agent.ID)
	if err == nil {
		session.CreatedAt = existing.CreatedAt
	} else {
		session.CreatedAt = time.Now()
	}

	path := getSessionPath(agent.ID)
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// loadSession loads a session from disk
func loadSession(id string) (*Session, error) {
	path := getSessionPath(id)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// listSessions returns a list of available sessions
func listSessions() ([]Session, error) {
	sessionsDir := getSessionsDir()
	if _, err := os.Stat(sessionsDir); os.IsNotExist(err) {
		return []Session{}, nil
	}

	entries, err := os.ReadDir(sessionsDir)
	if err != nil {
		return nil, err
	}

	var sessions []Session
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		id := strings.TrimSuffix(entry.Name(), ".json")
		session, err := loadSession(id)
		if err != nil {
			continue
		}
		sessions = append(sessions, *session)
	}

	// Sort by UpdatedAt descending
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].UpdatedAt.After(sessions[j].UpdatedAt)
	})

	return sessions, nil
}

// deleteSession deletes a session from disk
func deleteSession(id string) error {
	path := getSessionPath(id)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("session '%s' not found", id)
	}
	return os.Remove(path)
}

// formatSessionsList returns a formatted string of available sessions
func formatSessionsList() string {
	sessions, err := listSessions()
	if err != nil {
		return fmt.Sprintf("Error listing sessions: %s", err)
	}

	if len(sessions) == 0 {
		return "No saved sessions found."
	}

	var sb strings.Builder
	sb.WriteString("Available Sessions:\n")
	for _, s := range sessions {
		sb.WriteString(fmt.Sprintf("- %s (Last updated: %s)\n", s.ID, s.UpdatedAt.Format("2006-01-02 15:04:05")))
	}
	return sb.String()
}

// renameSession renames a session file and updates the internal ID if it's the current one
func renameSession(oldID, newID string) error {
	oldPath := getSessionPath(oldID)
	newPath := getSessionPath(newID)

	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("session '%s' already exists", newID)
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	// We also need to update the ID inside the JSON file
	session, err := loadSession(newID)
	if err == nil {
		session.ID = newID
		data, _ := json.MarshalIndent(session, "", "  ")
		os.WriteFile(newPath, data, 0644)
	}

	return nil
}

// NameSessionArgs represents arguments for naming a session
type NameSessionArgs struct {
	Name string `json:"name"`
}

// nameSession updates the agent's ID to the new name
func nameSession(agent *Agent, argsJSON string) (string, error) {
	var args NameSessionArgs
	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}

	name := strings.TrimSpace(args.Name)
	if name == "" {
		return "", fmt.Errorf("session name cannot be empty")
	}

	// Sanitize name
	safeName := strings.ReplaceAll(name, " ", "-")
	safeName = strings.ReplaceAll(safeName, "/", "-")
	safeName = strings.ReplaceAll(safeName, "\\", "-")

	// Rename existing session file if it exists
	if err := renameSession(agent.ID, safeName); err != nil {
		// If we can't rename (maybe old session wasn't saved yet), just proceed
		// But if new name already exists, we should warn or error?
		// renameSession returns error if newID exists.
		// Let's return the error if it's "already exists"
		if strings.Contains(err.Error(), "already exists") {
			return "", err
		}
		// Other errors might be "file not found" which is fine if we haven't saved yet
	}

	oldID := agent.ID
	agent.ID = safeName

	// Also rename todo list if it exists
	oldTodoPath, _ := getTodoListPath(oldID)
	newTodoPath, _ := getTodoListPath(safeName)
	if _, err := os.Stat(oldTodoPath); err == nil {
		os.Rename(oldTodoPath, newTodoPath)
		// Update the agent_id inside the todo json
		if todoList, err := loadTodoList(safeName); err == nil {
			todoList.AgentID = safeName
			saveTodoList(todoList)
		}
	}

	return fmt.Sprintf("Session renamed from '%s' to '%s'.", oldID, safeName), nil
}
