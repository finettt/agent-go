package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ExportSessionArgs represents arguments for exporting a session
type ExportSessionArgs struct {
	SessionID       string `json:"session_id,omitempty"`       // Specific session to export (defaults to current)
	Format          string `json:"format"`                     // "markdown", "json", or "txt"
	Filename        string `json:"filename,omitempty"`         // Custom filename (auto-generated if not provided)
	IncludeMetadata bool   `json:"include_metadata,omitempty"` // Include session metadata (default: true)
}

// getExportDir returns the path to the export directory
func getExportDir() string {
	return filepath.Join(".agent-go", "exports")
}

// ensureExportDir creates the export directory if it doesn't exist
func ensureExportDir() error {
	exportDir := getExportDir()
	return os.MkdirAll(exportDir, 0755)
}

// generateExportFilename creates a filename for the exported session
func generateExportFilename(session *Session, format string) string {
	timestamp := time.Now().Format("20060102-150405")
	safeName := strings.ReplaceAll(session.ID, " ", "-")
	safeName = strings.ReplaceAll(safeName, "/", "_")
	safeName = strings.ReplaceAll(safeName, "\\", "_")
	safeName = strings.ReplaceAll(safeName, "..", "_")

	// Limit filename length
	if len(safeName) > 50 {
		safeName = safeName[:50]
	}

	return fmt.Sprintf("session-%s-%s.%s", safeName, timestamp, format)
}

// formatSessionMarkdown formats a session as markdown
func formatSessionMarkdown(session *Session, includeMetadata bool) string {
	var sb strings.Builder

	if includeMetadata {
		sb.WriteString(fmt.Sprintf("# Session: %s\n\n", session.ID))
		sb.WriteString("## Metadata\n\n")
		sb.WriteString(fmt.Sprintf("- **Created**: %s\n", session.CreatedAt.Format("2006-01-02 15:04:05")))
		sb.WriteString(fmt.Sprintf("- **Updated**: %s\n", session.UpdatedAt.Format("2006-01-02 15:04:05")))
		if session.AgentDefName != "" {
			sb.WriteString(fmt.Sprintf("- **Agent Definition**: %s\n", session.AgentDefName))
		}
		sb.WriteString(fmt.Sprintf("- **Total Messages**: %d\n", len(session.Messages)))
		sb.WriteString(fmt.Sprintf("- **Total Tokens**: %d\n", session.TotalTokens))
		sb.WriteString(fmt.Sprintf("- **Tool Calls**: %d\n\n", session.ToolCalls))
	}

	sb.WriteString("## Conversation History\n\n")

	// Group messages by role for better readability
	for i, msg := range session.Messages {
		role := strings.ToUpper(msg.Role)
		sb.WriteString(fmt.Sprintf("### %s\n\n", role))

		if msg.Content != nil && *msg.Content != "" {
			content := *msg.Content
			// Format as code block if it looks like code
			if strings.Contains(content, "```") || strings.Contains(content, "function") ||
				strings.Contains(content, "package") || strings.Contains(content, "import") ||
				strings.Contains(content, "def ") || strings.Contains(content, "class ") {
				sb.WriteString(fmt.Sprintf("```\n%s\n```\n\n", content))
			} else {
				sb.WriteString(fmt.Sprintf("%s\n\n", content))
			}
		}

		if msg.ReasoningContent != nil && *msg.ReasoningContent != "" {
			sb.WriteString("**Reasoning:**\n\n")
			sb.WriteString(fmt.Sprintf("%s\n\n", *msg.ReasoningContent))
		}

		if len(msg.ToolCalls) > 0 {
			sb.WriteString("**Tool Calls:**\n\n")
			for _, toolCall := range msg.ToolCalls {
				sb.WriteString(fmt.Sprintf("- `%s`(%s)\n", toolCall.Function.Name, toolCall.ID))
				if toolCall.Function.Arguments != "" {
					sb.WriteString(fmt.Sprintf("  ```json\n%s\n  ```\n", toolCall.Function.Arguments))
				}
			}
			sb.WriteString("\n")
		}

		if msg.ToolCallID != "" {
			sb.WriteString(fmt.Sprintf("**Tool Response (ID: %s):**\n\n", msg.ToolCallID))
			if msg.Content != nil {
				sb.WriteString(fmt.Sprintf("%s\n\n", *msg.Content))
			}
		}

		// Add separator between messages (except last one)
		if i < len(session.Messages)-1 {
			sb.WriteString("---\n\n")
		}
	}

	if includeMetadata {
		sb.WriteString("## Token Usage Summary\n\n")
		sb.WriteString(fmt.Sprintf("- **Prompt Tokens**: %d\n", session.PromptTokens))
		sb.WriteString(fmt.Sprintf("- **Completion Tokens**: %d\n", session.CompletionTokens))
		sb.WriteString(fmt.Sprintf("- **Total Tokens**: %d\n", session.TotalTokens))
		sb.WriteString(fmt.Sprintf("- **Current Context**: %d tokens\n", session.CurrentContextTokens))
	}

	return sb.String()
}

// formatSessionText formats a session as plain text
func formatSessionText(session *Session, includeMetadata bool) string {
	var sb strings.Builder

	if includeMetadata {
		sb.WriteString(fmt.Sprintf("Session: %s\n", session.ID))
		sb.WriteString(fmt.Sprintf("Created: %s\n", session.CreatedAt.Format("2006-01-02 15:04:05")))
		sb.WriteString(fmt.Sprintf("Updated: %s\n", session.UpdatedAt.Format("2006-01-02 15:04:05")))
		if session.AgentDefName != "" {
			sb.WriteString(fmt.Sprintf("Agent Definition: %s\n", session.AgentDefName))
		}
		sb.WriteString(fmt.Sprintf("Total Messages: %d\n", len(session.Messages)))
		sb.WriteString(fmt.Sprintf("Total Tokens: %d\n", session.TotalTokens))
		sb.WriteString(fmt.Sprintf("Tool Calls: %d\n\n", session.ToolCalls))
	}

	sb.WriteString("Conversation History\n")
	sb.WriteString(strings.Repeat("=", 50) + "\n\n")

	for i, msg := range session.Messages {
		role := strings.ToUpper(msg.Role)
		sb.WriteString(fmt.Sprintf("[%s]\n", role))

		if msg.Content != nil && *msg.Content != "" {
			sb.WriteString(fmt.Sprintf("%s\n", *msg.Content))
		}

		if msg.ReasoningContent != nil && *msg.ReasoningContent != "" {
			sb.WriteString(fmt.Sprintf("Reasoning: %s\n", *msg.ReasoningContent))
		}

		if len(msg.ToolCalls) > 0 {
			sb.WriteString("Tool Calls:\n")
			for _, toolCall := range msg.ToolCalls {
				sb.WriteString(fmt.Sprintf("  - %s(%s): %s\n", toolCall.Function.Name, toolCall.ID, toolCall.Function.Arguments))
			}
		}

		if msg.ToolCallID != "" {
			sb.WriteString(fmt.Sprintf("Tool Response (ID: %s):\n", msg.ToolCallID))
			if msg.Content != nil {
				sb.WriteString(fmt.Sprintf("%s\n", *msg.Content))
			}
		}

		// Add separator between messages (except last one)
		if i < len(session.Messages)-1 {
			sb.WriteString("\n" + strings.Repeat("-", 30) + "\n\n")
		}
	}

	if includeMetadata {
		sb.WriteString("\n" + strings.Repeat("=", 50) + "\n")
		sb.WriteString("Token Usage Summary\n")
		sb.WriteString(fmt.Sprintf("Prompt Tokens: %d\n", session.PromptTokens))
		sb.WriteString(fmt.Sprintf("Completion Tokens: %d\n", session.CompletionTokens))
		sb.WriteString(fmt.Sprintf("Total Tokens: %d\n", session.TotalTokens))
		sb.WriteString(fmt.Sprintf("Current Context: %d tokens\n", session.CurrentContextTokens))
	}

	return sb.String()
}

// formatSessionJSON formats a session as JSON with additional export metadata
func formatSessionJSON(session *Session, includeMetadata bool) (string, error) {
	// Create export structure
	type ExportSession struct {
		Session
		ExportedAt   time.Time `json:"exported_at"`
		ExportFormat string    `json:"export_format"`
		MessageCount int       `json:"message_count"`
	}

	export := ExportSession{
		Session:      *session,
		ExportedAt:   time.Now(),
		ExportFormat: "json",
		MessageCount: len(session.Messages),
	}

	// If not including metadata, clear the metadata fields
	if !includeMetadata {
		export.CurrentContextTokens = 0
		export.CurrentPromptTokens = 0
		export.CurrentCompletionTokens = 0
		export.TotalTokens = 0
		export.PromptTokens = 0
		export.CompletionTokens = 0
		export.ToolCalls = 0
	}

	data, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal session as JSON: %w", err)
	}

	return string(data), nil
}

// exportSession is the main export function
func exportSession(agent *Agent, argsJSON string) (string, error) {
	var args ExportSessionArgs
	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}

	// Validate format
	format := strings.ToLower(args.Format)
	if format != "markdown" && format != "json" && format != "txt" {
		return "", fmt.Errorf("format must be 'markdown', 'json', or 'txt'")
	}

	// Default metadata inclusion to true
	includeMetadata := args.IncludeMetadata
	if args.IncludeMetadata == false {
		// This handles the explicit false case
	} else {
		includeMetadata = true
	}

	// Get session to export
	var session *Session
	var err error

	if args.SessionID != "" {
		// Export specific session
		session, err = loadSession(args.SessionID)
		if err != nil {
			return "", fmt.Errorf("failed to load session '%s': %w", args.SessionID, err)
		}
	} else {
		// Export current session
		session = &Session{
			ID:           agent.ID,
			Messages:     agent.Messages,
			AgentDefName: agent.AgentDefName,
			CreatedAt:    time.Now(), // We don't have creation time for current session
			UpdatedAt:    time.Now(),
			// Use current token tracking variables
			CurrentContextTokens:    currentContextTokens,
			CurrentPromptTokens:     currentPromptTokens,
			CurrentCompletionTokens: currentCompletionTokens,
			TotalTokens:             totalTokens,
			PromptTokens:            totalPromptTokens,
			CompletionTokens:        totalCompletionTokens,
			ToolCalls:               totalToolCalls,
		}
	}

	// Generate filename
	filename := args.Filename
	if filename == "" {
		filename = generateExportFilename(session, format)
	} else {
		// Ensure filename has correct extension
		ext := filepath.Ext(filename)
		if ext == "" {
			filename += "." + format
		} else if ext != "."+format {
			filename = strings.TrimSuffix(filename, ext) + "." + format
		}
	}

	// Ensure export directory exists
	if err := ensureExportDir(); err != nil {
		return "", fmt.Errorf("failed to create export directory: %w", err)
	}

	// Format content based on format
	var content string
	switch format {
	case "markdown":
		content = formatSessionMarkdown(session, includeMetadata)
	case "txt":
		content = formatSessionText(session, includeMetadata)
	case "json":
		content, err = formatSessionJSON(session, includeMetadata)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}

	// Write to file
	exportPath := filepath.Join(getExportDir(), filename)
	if err := os.WriteFile(exportPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write export file: %w", err)
	}

	// Return success message
	return fmt.Sprintf("Session exported successfully to: %s\nFormat: %s\nMessages: %d\nTokens: %d",
		exportPath, format, len(session.Messages), session.TotalTokens), nil
}

// listExports returns a list of available export files
func listExports() ([]string, error) {
	exportDir := getExportDir()
	if _, err := os.Stat(exportDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	entries, err := os.ReadDir(exportDir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		files = append(files, entry.Name())
	}

	// Sort by modification time (newest first)
	type fileInfo struct {
		name    string
		modTime time.Time
	}

	var fileInfos []fileInfo
	for _, filename := range files {
		path := filepath.Join(exportDir, filename)
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		fileInfos = append(fileInfos, fileInfo{
			name:    filename,
			modTime: info.ModTime(),
		})
	}

	// Sort by mod time descending
	for i := 0; i < len(fileInfos)-1; i++ {
		for j := i + 1; j < len(fileInfos); j++ {
			if fileInfos[i].modTime.Before(fileInfos[j].modTime) {
				fileInfos[i], fileInfos[j] = fileInfos[j], fileInfos[i]
			}
		}
	}

	// Extract just filenames
	var sortedFiles []string
	for _, fi := range fileInfos {
		sortedFiles = append(sortedFiles, fi.name)
	}

	return sortedFiles, nil
}

// formatExportsList returns a formatted string of available exports
func formatExportsList() string {
	files, err := listExports()
	if err != nil {
		return fmt.Sprintf("Error listing exports: %s", err)
	}

	if len(files) == 0 {
		return "No exported sessions found."
	}

	var sb strings.Builder
	sb.WriteString("Available Exports:\n")
	exportDir := getExportDir()

	for _, filename := range files {
		path := filepath.Join(exportDir, filename)
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		// Get file size in a readable format
		size := info.Size()
		sizeStr := "B"
		if size > 1024 {
			size = size / 1024
			sizeStr = "KB"
		}
		if size > 1024 {
			size = size / 1024
			sizeStr = "MB"
		}

		sb.WriteString(fmt.Sprintf("- %s (%d %s, Modified: %s)\n",
			filename, size, sizeStr, info.ModTime().Format("2006-01-02 15:04:05")))
	}

	return sb.String()
}
