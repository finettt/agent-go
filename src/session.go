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
	AgentDefName     string    `json:"agent_def_name,omitempty"`
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
		AgentDefName:     agent.AgentDefName,
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

// generateSessionName uses the mini model to generate a descriptive name for the session
func generateSessionName(agent *Agent, config *Config) (string, error) {
	// Construct a prompt for the mini model
	var promptBuilder strings.Builder
	promptBuilder.WriteString("Analyze the following conversation and generate a short, descriptive, filesystem-friendly name for this session.\n")
	promptBuilder.WriteString("- The name should be lowercased, using dashes for spaces (e.g., 'implement-auth-feature').\n")
	promptBuilder.WriteString("- Maximum length: 50 characters.\n")
	promptBuilder.WriteString("- Do not include file extensions or special characters other than dashes.\n")
	promptBuilder.WriteString("- Return ONLY the name, nothing else.\n\n")
	promptBuilder.WriteString("Conversation Context:\n")

	// Include a few recent messages for context, up to a limit
	start := len(agent.Messages) - 10
	if start < 0 {
		start = 0
	}
	for i := start; i < len(agent.Messages); i++ {
		msg := agent.Messages[i]
		if msg.Content != nil {
			role := msg.Role
			content := *msg.Content
			// Truncate long messages
			if len(content) > 200 {
				content = content[:200] + "..."
			}
			promptBuilder.WriteString(fmt.Sprintf("%s: %s\n", role, content))
		}
	}

	name, err := sendMiniLLMRequest(config, []Message{{Role: "user", Content: genericStringPointer(promptBuilder.String())}})
	if err != nil {
		return "", err
	}

	// Clean up the response
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, "\n", "")
	name = strings.ReplaceAll(name, "\r", "")
	name = strings.ReplaceAll(name, "\"", "")
	name = strings.ReplaceAll(name, "'", "")
	name = strings.ToLower(name)

	return name, nil
}

// Helper to get string pointer
func genericStringPointer(s string) *string {
	return &s
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

// formatSessionView returns a formatted string of a specific session
func formatSessionView(id string) string {
	session, err := loadSession(id)
	if err != nil {
		return fmt.Sprintf("Error loading session '%s': %s", id, err)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Session: %s\n", session.ID))
	sb.WriteString(fmt.Sprintf("Created: %s\n", session.CreatedAt.Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("Updated: %s\n", session.UpdatedAt.Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("Messages: %d\n", len(session.Messages)))
	sb.WriteString(fmt.Sprintf("Tokens: %d (Prompt: %d, Completion: %d)\n", session.TotalTokens, session.PromptTokens, session.CompletionTokens))
	sb.WriteString(fmt.Sprintf("Tool Calls: %d\n", session.ToolCalls))
	sb.WriteString("\nRecent Messages:\n")

	// Show last 5 messages or fewer if there aren't that many
	start := len(session.Messages) - 5
	if start < 0 {
		start = 0
	}

	for i := start; i < len(session.Messages); i++ {
		msg := session.Messages[i]
		role := msg.Role
		content := ""
		if msg.Content != nil {
			content = *msg.Content
		} else if len(msg.ToolCalls) > 0 {
			content = fmt.Sprintf("[Tool Call: %s]", msg.ToolCalls[0].Function.Name)
		}

		// Truncate content if it's too long
		if len(content) > 100 {
			content = content[:97] + "..."
		}
		// Replace newlines with spaces for preview
		content = strings.ReplaceAll(content, "\n", " ")

		sb.WriteString(fmt.Sprintf("- [%s]: %s\n", role, content))
	}

	return sb.String()
}
