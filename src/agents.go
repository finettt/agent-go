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

// AgentDefinition is a user-created, task-specific agent configuration saved to disk.
type AgentDefinition struct {
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	SystemPrompt string    `json:"system_prompt"`
	Model        string    `json:"model,omitempty"`
	Temperature  *float32  `json:"temperature,omitempty"`
	MaxTokens    *int      `json:"max_tokens,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	// AllowedTools is an optional whitelist of tool function names this agent may use.
	// When non-empty, only the listed tools will be exposed to the model.
	AllowedTools []string `json:"allowed_tools,omitempty"`
	// DeniedTools is an optional blacklist of tool function names this agent may NOT use.
	// Only used when AllowedTools is empty.
	DeniedTools []string `json:"denied_tools,omitempty"`
}

func isBuiltInAgentName(name string) bool {
	return strings.TrimSpace(name) == "default"
}

func builtInDefaultAgentSystemPrompt() string {
	return strings.TrimSpace(`
You are the built-in "default" task-specific agent for Agent-Go.

Behavior:
- Be direct and technical.
- When writing code, prefer correct, minimal changes.
- If command execution is available, propose safe commands and explain briefly.
- Respect tool constraints imposed by the system (Plan vs Build, confirmation mode, etc.).
`)
}

func getBuiltInAgentDefinition(name string) (*AgentDefinition, bool) {
	safe, err := sanitizeAgentName(name)
	if err != nil {
		return nil, false
	}
	if safe != "default" {
		return nil, false
	}
	// The default agent has no tool restrictions - it uses all available tools
	// based on the operation mode (Plan vs Build). The tool availability is
	// controlled by getAvailableTools() which filters based on mode.
	return &AgentDefinition{
		Name:         "default",
		Description:  "Built-in default agent with full tool access",
		SystemPrompt: builtInDefaultAgentSystemPrompt(),
		// No AllowedTools or DeniedTools - allows all tools
	}, true
}

// getAgentsDir returns the path to the agents directory.
func getAgentsDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "agent-go", "agents")
}

func ensureAgentsDir() error {
	return os.MkdirAll(getAgentsDir(), 0755)
}

// sanitizeAgentName converts a user-provided name into a safe identifier.
func sanitizeAgentName(name string) (string, error) {
	s := strings.TrimSpace(name)
	if s == "" {
		return "", fmt.Errorf("agent name cannot be empty")
	}

	// Keep it simple and filesystem-friendly.
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, "\\", "-")
	s = strings.ReplaceAll(s, "..", "-")

	return s, nil
}

func getAgentPath(name string) (string, error) {
	safe, err := sanitizeAgentName(name)
	if err != nil {
		return "", err
	}
	return filepath.Join(getAgentsDir(), safe+".json"), nil
}

func saveAgentDefinition(def *AgentDefinition) error {
	if def == nil {
		return fmt.Errorf("agent definition is nil")
	}
	safe, err := sanitizeAgentName(def.Name)
	if err != nil {
		return err
	}
	def.Name = safe

	if strings.TrimSpace(def.SystemPrompt) == "" {
		return fmt.Errorf("system_prompt cannot be empty")
	}

	// Validate tool policy
	if len(def.AllowedTools) > 0 && len(def.DeniedTools) > 0 {
		fmt.Printf("Warning: Both allowed_tools and denied_tools are set for agent '%s'. Using allowed_tools (whitelist mode).\n", def.Name)
	}

	if err := ensureAgentsDir(); err != nil {
		return err
	}

	now := time.Now()
	if def.CreatedAt.IsZero() {
		def.CreatedAt = now
	}
	def.UpdatedAt = now

	path, err := getAgentPath(def.Name)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(def, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func loadAgentDefinition(name string) (*AgentDefinition, error) {
	// Built-in agent fallback (no file required).
	if def, ok := getBuiltInAgentDefinition(name); ok {
		// If a file exists with the same name, we'll prefer the file (but "default" is reserved).
		// We still attempt file loading below; if it doesn't exist, we return the built-in.
		_ = def
	}

	path, err := getAgentPath(name)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			if def, ok := getBuiltInAgentDefinition(name); ok {
				return def, nil
			}
		}
		return nil, err
	}

	var def AgentDefinition
	if err := json.Unmarshal(data, &def); err != nil {
		return nil, err
	}

	return &def, nil
}

func deleteAgentDefinition(name string) error {
	if isBuiltInAgentName(name) {
		return fmt.Errorf("cannot delete built-in agent '%s'", strings.TrimSpace(name))
	}

	path, err := getAgentPath(name)
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("agent '%s' not found", name)
	}
	return os.Remove(path)
}

func listAgentDefinitions() ([]*AgentDefinition, error) {
	defs := make([]*AgentDefinition, 0)

	// Always include built-in "default".
	if def, ok := getBuiltInAgentDefinition("default"); ok {
		defs = append(defs, def)
	}

	dir := getAgentsDir()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// No directory yet, but built-ins still exist.
		return defs, nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		name := strings.TrimSuffix(entry.Name(), ".json")

		// Avoid duplicates (and reserve the built-in name).
		if isBuiltInAgentName(name) {
			continue
		}

		def, err := loadAgentDefinition(name)
		if err != nil {
			continue
		}
		defs = append(defs, def)
	}

	sort.Slice(defs, func(i, j int) bool {
		iBuiltIn := isBuiltInAgentName(defs[i].Name)
		jBuiltIn := isBuiltInAgentName(defs[j].Name)
		if iBuiltIn != jBuiltIn {
			return iBuiltIn
		}
		return defs[i].UpdatedAt.After(defs[j].UpdatedAt)
	})

	return defs, nil
}

func listAgentNames() ([]string, error) {
	defs, err := listAgentDefinitions()
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(defs))
	for _, d := range defs {
		if d == nil {
			continue
		}
		names = append(names, d.Name)
	}
	sort.Strings(names)
	return names, nil
}

func formatAgentsList() string {
	defs, err := listAgentDefinitions()
	if err != nil {
		return fmt.Sprintf("Error listing agents: %v", err)
	}
	if len(defs) == 0 {
		return "No agents found. Use /agent studio to create one."
	}

	var b strings.Builder
	b.WriteString("Agents:\n")
	for _, d := range defs {
		desc := strings.TrimSpace(d.Description)
		if desc == "" {
			desc = "(no description)"
		}

		if isBuiltInAgentName(d.Name) {
			b.WriteString(fmt.Sprintf("  - %s: %s (built-in)\n", d.Name, desc))
			continue
		}

		b.WriteString(fmt.Sprintf("  - %s: %s (updated: %s)\n", d.Name, desc, d.UpdatedAt.Format("2006-01-02 15:04")))
	}
	return b.String()
}

// getAgentsForSystemPrompt returns a compact list of available task-specific agents for injection into the main system prompt.
// This helps the model choose an agent name when using spawn_agent({agent:"..."}) or guiding the user to /agent commands.
func getAgentsForSystemPrompt() string {
	defs, err := listAgentDefinitions()
	if err != nil || len(defs) == 0 {
		return ""
	}

	const maxAgents = 20

	var b strings.Builder
	b.WriteString("\n\n=== Task-Specific Agents ===\n")
	b.WriteString("Available agent names (for /agent use|view and spawn_agent {agent:\"<name>\"}):\n")

	for i, d := range defs {
		if d == nil {
			continue
		}
		if i >= maxAgents {
			remaining := len(defs) - maxAgents
			if remaining > 0 {
				b.WriteString(fmt.Sprintf("  - ... and %d more\n", remaining))
			}
			break
		}

		desc := strings.TrimSpace(d.Description)
		if desc == "" {
			desc = "(no description)"
		}

		if isBuiltInAgentName(d.Name) {
			b.WriteString(fmt.Sprintf("  - %s (built-in): %s\n", d.Name, desc))
		} else {
			b.WriteString(fmt.Sprintf("  - %s: %s\n", d.Name, desc))
		}
	}

	return b.String()
}

func formatAgentView(name string) string {
	def, err := loadAgentDefinition(name)
	if err != nil {
		return fmt.Sprintf("Agent '%s' not found.", name)
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("=== %s ===\n", def.Name))
	if strings.TrimSpace(def.Description) != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", def.Description))
	}
	if isBuiltInAgentName(def.Name) {
		b.WriteString("Type: built-in\n")
	}
	if def.Model != "" {
		b.WriteString(fmt.Sprintf("Model: %s\n", def.Model))
	}
	if def.Temperature != nil {
		b.WriteString(fmt.Sprintf("Temperature: %v\n", *def.Temperature))
	}
	if def.MaxTokens != nil {
		b.WriteString(fmt.Sprintf("MaxTokens: %d\n", *def.MaxTokens))
	}

	// Display tool policy if configured
	if len(def.AllowedTools) > 0 {
		b.WriteString(fmt.Sprintf("Tool Policy: Whitelist (%d tools)\n", len(def.AllowedTools)))
		b.WriteString("Allowed Tools: " + strings.Join(def.AllowedTools, ", ") + "\n")
	} else if len(def.DeniedTools) > 0 {
		b.WriteString(fmt.Sprintf("Tool Policy: Blacklist (%d tools denied)\n", len(def.DeniedTools)))
		b.WriteString("Denied Tools: " + strings.Join(def.DeniedTools, ", ") + "\n")
	} else {
		b.WriteString("Tool Policy: All tools available\n")
	}

	if !isBuiltInAgentName(def.Name) {
		b.WriteString(fmt.Sprintf("Created: %s\n", def.CreatedAt.Format("2006-01-02 15:04:05")))
		b.WriteString(fmt.Sprintf("Updated: %s\n", def.UpdatedAt.Format("2006-01-02 15:04:05")))
	}
	b.WriteString("\nSystem Prompt:\n")
	b.WriteString(def.SystemPrompt)
	b.WriteString("\n")
	return b.String()
}

// CreateAgentDefinitionArgs represents arguments for create_agent_definition tool.
type CreateAgentDefinitionArgs struct {
	Name         string   `json:"name"`
	Description  string   `json:"description,omitempty"`
	SystemPrompt string   `json:"system_prompt"`
	Model        string   `json:"model,omitempty"`
	Temperature  *float32 `json:"temperature,omitempty"`
	MaxTokens    *int     `json:"max_tokens,omitempty"`
	// Optional tool policy
	AllowedTools []string `json:"allowed_tools,omitempty"`
	DeniedTools  []string `json:"denied_tools,omitempty"`
}

// createAgentDefinition is a tool handler that persists a new agent definition to disk.
func createAgentDefinition(argsJSON string) (string, error) {
	var args CreateAgentDefinitionArgs
	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}

	name, err := sanitizeAgentName(args.Name)
	if err != nil {
		return "", err
	}
	if isBuiltInAgentName(name) {
		return "", fmt.Errorf("agent name '%s' is reserved for a built-in agent", name)
	}

	def := &AgentDefinition{
		Name:         name,
		Description:  strings.TrimSpace(args.Description),
		SystemPrompt: strings.TrimSpace(args.SystemPrompt),
		Model:        strings.TrimSpace(args.Model),
		Temperature:  args.Temperature,
		MaxTokens:    args.MaxTokens,
		AllowedTools: args.AllowedTools,
		DeniedTools:  args.DeniedTools,
	}

	// Prevent silent overwrite: if exists, require user to delete first.
	if _, err := loadAgentDefinition(def.Name); err == nil {
		return "", fmt.Errorf("agent '%s' already exists", def.Name)
	}

	if err := saveAgentDefinition(def); err != nil {
		return "", err
	}

	return fmt.Sprintf("Agent '%s' created successfully.", def.Name), nil
}
