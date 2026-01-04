package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ShadowGit manages a hidden git repository for checkpoints
type ShadowGit struct {
	RepoDir  string
	WorkTree string
}

// NewShadowGit creates a new ShadowGit instance
func NewShadowGit(agentID string) (*ShadowGit, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Store shadow repos in ~/.config/agent-go/checkpoints/shadow_git/<agent_id>
	repoDir := filepath.Join(home, ".config", "agent-go", "checkpoints", "shadow_git", agentID)

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return &ShadowGit{
		RepoDir:  repoDir,
		WorkTree: cwd,
	}, nil
}

// Init initializes the shadow git repository if it doesn't exist
func (g *ShadowGit) Init() error {
	if err := os.MkdirAll(g.RepoDir, 0755); err != nil {
		return err
	}

	// Check if already initialized
	if _, err := os.Stat(filepath.Join(g.RepoDir, "HEAD")); err == nil {
		return nil
	}

	// For init --bare, we should NOT pass --work-tree, as bare repos don't have one.
	// We call exec directly here instead of using runGit.
	cmd := exec.Command("git", "--git-dir="+g.RepoDir, "init", "--bare")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git init failed: %s: %s", err, string(output))
	}
	return nil
}

// Commit creates a new commit with the current state of the workspace
func (g *ShadowGit) Commit(message string, config *Config) (string, error) {
	// Add all files
	if err := g.runGit("add", "."); err != nil {
		return "", fmt.Errorf("failed to add files: %w", err)
	}

	// Check if there are changes to commit
	if err := g.runGit("diff-index", "--quiet", "HEAD"); err == nil {
		// No changes, get current hash
		return g.getCurrentHash()
	}

	// Auto-generate message if empty or "auto"
	if (message == "" || message == "auto") && config != nil {
		generated, err := g.generateCommitMessage(config)
		if err == nil && generated != "" {
			message = generated
		} else if message == "" {
			message = "Checkpoint: Auto-generated"
		}
	} else if message == "" {
		message = "Checkpoint"
	}

	// Commit
	if err := g.runGit("commit", "--allow-empty", "-m", message); err != nil {
		return "", fmt.Errorf("failed to commit: %w", err)
	}

	return g.getCurrentHash()
}

// generateCommitMessage uses the mini model to generate a commit message from staged changes
func (g *ShadowGit) generateCommitMessage(config *Config) (string, error) {
	// Get diff of staged changes
	// We use --cached (or --staged) to get diff of what is about to be committed
	cmd := exec.Command("git", "--git-dir="+g.RepoDir, "--work-tree="+g.WorkTree, "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get diff: %w", err)
	}

	diff := string(output)
	if len(diff) > 4000 {
		// Truncate diff if too large to avoid hitting token limits
		diff = diff[:4000] + "\n... (truncated)"
	}

	// Construct prompt
	var promptBuilder strings.Builder
	promptBuilder.WriteString("Generate a concise, conventional commit message for the following changes.\n")
	promptBuilder.WriteString("- Format: <type>: <subject>\n")
	promptBuilder.WriteString("- Keep it under 70 characters if possible.\n")
	promptBuilder.WriteString("- Return ONLY the commit message.\n\n")
	promptBuilder.WriteString("Changes:\n")
	promptBuilder.WriteString(diff)

	msg, err := sendMiniLLMRequest(config, []Message{{Role: "user", Content: genericStringPointer(promptBuilder.String())}})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(msg), nil
}

// Restore restores the workspace to a specific commit hash
func (g *ShadowGit) Restore(hash string) error {
	// Force checkout to the specific hash
	// -f throws away local changes
	if err := g.runGit("checkout", "-f", hash); err != nil {
		return fmt.Errorf("failed to checkout %s: %w", hash, err)
	}

	// Clean untracked files/directories
	// -f force, -d directories
	if err := g.runGit("clean", "-fd"); err != nil {
		return fmt.Errorf("failed to clean workspace: %w", err)
	}

	return nil
}

// getCurrentHash returns the current HEAD hash
func (g *ShadowGit) getCurrentHash() (string, error) {
	cmd := exec.Command("git", "--git-dir="+g.RepoDir, "--work-tree="+g.WorkTree, "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		// If HEAD doesn't exist yet (first commit failed or not made), return empty
		return "", nil
	}
	return strings.TrimSpace(string(output)), nil
}

// runGit executes a git command in the shadow repo context
func (g *ShadowGit) runGit(args ...string) error {
	baseArgs := []string{"--git-dir=" + g.RepoDir, "--work-tree=" + g.WorkTree}
	finalArgs := append(baseArgs, args...)

	cmd := exec.Command("git", finalArgs...)
	// We capture stderr for error reporting but discard stdout usually
	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git error: %s - %s", err, stderr.String())
	}
	return nil
}
