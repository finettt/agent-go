package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Note represents a note created by the agent
type Note struct {
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateNoteArgs represents arguments for creating a note
type CreateNoteArgs struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// UpdateNoteArgs represents arguments for updating a note
type UpdateNoteArgs struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// DeleteNoteArgs represents arguments for deleting a note
type DeleteNoteArgs struct {
	Name string `json:"name"`
}

// getNotesDir returns the path to the notes directory
func getNotesDir() string {
	return filepath.Join(".agent-go", "notes")
}

// ensureNotesDir creates the notes directory if it doesn't exist
func ensureNotesDir() error {
	notesDir := getNotesDir()
	return os.MkdirAll(notesDir, 0755)
}

// getNotePath returns the path to a specific note file
func getNotePath(name string) string {
	// Sanitize the name to be a valid filename
	safeName := strings.ReplaceAll(name, "/", "_")
	safeName = strings.ReplaceAll(safeName, "\\", "_")
	safeName = strings.ReplaceAll(safeName, "..", "_")
	return filepath.Join(getNotesDir(), safeName+".json")
}

// loadNote loads a note from disk
func loadNote(name string) (*Note, error) {
	path := getNotePath(name)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var note Note
	if err := json.Unmarshal(data, &note); err != nil {
		return nil, err
	}

	return &note, nil
}

// saveNote saves a note to disk
func saveNote(note *Note) error {
	if err := ensureNotesDir(); err != nil {
		return err
	}

	path := getNotePath(note.Name)
	data, err := json.MarshalIndent(note, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// listNotes returns all notes
func listNotes() ([]*Note, error) {
	notesDir := getNotesDir()

	// Check if directory exists
	if _, err := os.Stat(notesDir); os.IsNotExist(err) {
		return []*Note{}, nil
	}

	entries, err := os.ReadDir(notesDir)
	if err != nil {
		return nil, err
	}

	var notes []*Note
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".json")
		note, err := loadNote(name)
		if err != nil {
			continue // Skip invalid notes
		}
		notes = append(notes, note)
	}

	return notes, nil
}

// listNoteNames returns just the names of all notes
func listNoteNames() ([]string, error) {
	notesDir := getNotesDir()

	// Check if directory exists
	if _, err := os.Stat(notesDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	entries, err := os.ReadDir(notesDir)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		names = append(names, strings.TrimSuffix(entry.Name(), ".json"))
	}

	return names, nil
}

// createNote creates a new note
func createNote(argsJSON string) (string, error) {
	var args CreateNoteArgs
	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}

	// Validate name
	if strings.TrimSpace(args.Name) == "" {
		return "", fmt.Errorf("note name cannot be empty")
	}

	// Validate content
	if strings.TrimSpace(args.Content) == "" {
		return "", fmt.Errorf("note content cannot be empty")
	}

	// Check if note already exists
	if _, err := loadNote(args.Name); err == nil {
		return "", fmt.Errorf("note '%s' already exists, use update_note to modify it", args.Name)
	}

	now := time.Now()
	note := &Note{
		Name:      args.Name,
		Content:   args.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := saveNote(note); err != nil {
		return "", err
	}

	return fmt.Sprintf("Note '%s' created successfully.", args.Name), nil
}

// updateNote updates an existing note
func updateNote(argsJSON string) (string, error) {
	var args UpdateNoteArgs
	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}

	// Validate name
	if strings.TrimSpace(args.Name) == "" {
		return "", fmt.Errorf("note name cannot be empty")
	}

	// Validate content
	if strings.TrimSpace(args.Content) == "" {
		return "", fmt.Errorf("note content cannot be empty")
	}

	// Load existing note
	note, err := loadNote(args.Name)
	if err != nil {
		return "", fmt.Errorf("note '%s' not found", args.Name)
	}

	// Update the note
	note.Content = args.Content
	note.UpdatedAt = time.Now()

	if err := saveNote(note); err != nil {
		return "", err
	}

	return fmt.Sprintf("Note '%s' updated successfully.", args.Name), nil
}

// deleteNote deletes a note
func deleteNote(argsJSON string) (string, error) {
	var args DeleteNoteArgs
	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}

	// Validate name
	if strings.TrimSpace(args.Name) == "" {
		return "", fmt.Errorf("note name cannot be empty")
	}

	path := getNotePath(args.Name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("note '%s' not found", args.Name)
	}

	if err := os.Remove(path); err != nil {
		return "", err
	}

	return fmt.Sprintf("Note '%s' deleted successfully.", args.Name), nil
}

// getNotesForSystemPrompt returns all notes formatted for injection into the system prompt
func getNotesForSystemPrompt() string {
	notes, err := listNotes()
	if err != nil || len(notes) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("\n\n=== Agent Notes ===\n")

	for _, note := range notes {
		builder.WriteString(fmt.Sprintf("\n[%s]\n%s\n", note.Name, note.Content))
	}

	return builder.String()
}

// formatNotesList formats the list of notes for display
func formatNotesList() string {
	notes, err := listNotes()
	if err != nil {
		return fmt.Sprintf("Error loading notes: %s", err)
	}

	if len(notes) == 0 {
		return "No notes found."
	}

	var builder strings.Builder
	builder.WriteString("Notes:\n")

	for _, note := range notes {
		builder.WriteString(fmt.Sprintf("  - %s (updated: %s)\n",
			note.Name,
			note.UpdatedAt.Format("2006-01-02 15:04")))
	}

	return builder.String()
}

// formatNoteView formats a single note for display
func formatNoteView(name string) string {
	note, err := loadNote(name)
	if err != nil {
		return fmt.Sprintf("Note '%s' not found.", name)
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("=== %s ===\n", note.Name))
	builder.WriteString(fmt.Sprintf("Created: %s\n", note.CreatedAt.Format("2006-01-02 15:04:05")))
	builder.WriteString(fmt.Sprintf("Updated: %s\n", note.UpdatedAt.Format("2006-01-02 15:04:05")))
	builder.WriteString(fmt.Sprintf("\n%s\n", note.Content))

	return builder.String()
}
