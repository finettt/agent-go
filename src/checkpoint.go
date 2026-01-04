package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Checkpoint represents a saved state of the agent and environment
type Checkpoint struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at"`
	AgentID       string    `json:"agent_id"`
	Messages      []Message `json:"messages"`
	TotalTokens   int       `json:"total_tokens"`
	GitCommitHash string    `json:"git_commit_hash,omitempty"`
	DockerImageID string    `json:"docker_image_id,omitempty"`
	IsAuto        bool      `json:"is_auto"`
}

// getCheckpointsDir returns the directory for checkpoints
func getCheckpointsDir(agentID string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "agent-go", "checkpoints", "metadata", agentID)
}

// ensureCheckpointsDir creates the checkpoints directory
func ensureCheckpointsDir(agentID string) error {
	dir := getCheckpointsDir(agentID)
	return os.MkdirAll(dir, 0755)
}

// createCheckpoint creates a new checkpoint
func createCheckpoint(agent *Agent, config *Config, name string, isAuto bool) (string, error) {
	if err := ensureCheckpointsDir(agent.ID); err != nil {
		return "", fmt.Errorf("failed to create checkpoint dir: %w", err)
	}

	// 1. Initialize/Commit Shadow Git (Files)
	shadowGit, err := NewShadowGit(agent.ID)
	if err != nil {
		return "", fmt.Errorf("failed to init shadow git: %w", err)
	}
	if err := shadowGit.Init(); err != nil {
		return "", fmt.Errorf("failed to init shadow git repo: %w", err)
	}

	// Determine commit message
	msg := fmt.Sprintf("Checkpoint: %s", name)
	if isAuto {
		// Use "auto" to trigger generation if configured
		msg = "auto"
	}

	commitHash, err := shadowGit.Commit(msg, config)
	if err != nil {
		return "", fmt.Errorf("failed to commit files: %w", err)
	}

	// 2. Docker Commit (System) - Only if in Sandbox
	var dockerImageID string
	if isRunningInDocker() {
		// Try to commit current container
		// We need to know our own container ID.
		// /proc/self/cgroup often contains it, or hostname.
		containerID, err := getCurrentContainerID()
		if err == nil && containerID != "" {
			// Check if we can talk to docker
			if err := checkDockerAccess(); err == nil {
				// Commit container
				// Image name: agent-go-checkpoint-<agentID>-<timestamp>
				timestamp := time.Now().Format("20060102150405")
				imageName := fmt.Sprintf("agent-go-ckpt-%s-%s", agent.ID, timestamp)
				if out, err := exec.Command("docker", "commit", containerID, imageName).Output(); err == nil {
					// Output is usually "sha256:..."
					dockerImageID = strings.TrimSpace(string(out))
				} else {
					fmt.Printf("Warning: Failed to commit docker container: %v\n", err)
				}
			}
		}
	}

	// 3. Create Checkpoint Metadata (Memory)
	checkpointID := time.Now().Format("20060102_150405")
	checkpoint := Checkpoint{
		ID:            checkpointID,
		Name:          name,
		CreatedAt:     time.Now(),
		AgentID:       agent.ID,
		Messages:      make([]Message, len(agent.Messages)),
		TotalTokens:   totalTokens,
		GitCommitHash: commitHash,
		DockerImageID: dockerImageID,
		IsAuto:        isAuto,
	}
	copy(checkpoint.Messages, agent.Messages)

	// Save metadata
	filename := filepath.Join(getCheckpointsDir(agent.ID), checkpointID+".json")
	data, err := json.MarshalIndent(checkpoint, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return "", err
	}

	// Prune old auto-checkpoints (keep last 10)
	if isAuto {
		pruneAutoCheckpoints(agent.ID)
	}

	return checkpointID, nil
}

// listCheckpoints returns a list of checkpoints for an agent
func listCheckpoints(agentID string) ([]Checkpoint, error) {
	dir := getCheckpointsDir(agentID)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return []Checkpoint{}, nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var checkpoints []Checkpoint
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}

		var cp Checkpoint
		if err := json.Unmarshal(data, &cp); err != nil {
			continue
		}
		checkpoints = append(checkpoints, cp)
	}

	// Sort by CreatedAt descending
	sort.Slice(checkpoints, func(i, j int) bool {
		return checkpoints[i].CreatedAt.After(checkpoints[j].CreatedAt)
	})

	return checkpoints, nil
}

// restoreCheckpoint restores a checkpoint
func restoreCheckpoint(agent *Agent, checkpointID string) error {
	// Load checkpoint
	path := filepath.Join(getCheckpointsDir(agent.ID), checkpointID+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("checkpoint not found: %w", err)
	}

	var cp Checkpoint
	if err := json.Unmarshal(data, &cp); err != nil {
		return fmt.Errorf("failed to parse checkpoint: %w", err)
	}

	// 1. Restore Files
	if cp.GitCommitHash != "" {
		shadowGit, err := NewShadowGit(agent.ID)
		if err != nil {
			return fmt.Errorf("failed to init shadow git: %w", err)
		}
		if err := shadowGit.Restore(cp.GitCommitHash); err != nil {
			return fmt.Errorf("failed to restore files: %w", err)
		}
	}

	// 2. Restore System (Docker)
	if cp.DockerImageID != "" {
		// We cannot restore a running container to a previous image from inside.
		// We must warn the user.
		fmt.Printf("\n%sWARNING: This checkpoint includes a system snapshot (Docker Image: %s).%s\n", ColorRed, cp.DockerImageID, ColorReset)
		fmt.Printf("Restoring the file system and memory is done, but installed apps/system changes require a restart.\n")
		fmt.Printf("To fully restore system state, restart the sandbox with:\n")
		fmt.Printf("  docker run ... %s\n\n", cp.DockerImageID)
	} else if isRunningInDocker() {
		fmt.Printf("\n%sWarning: This checkpoint does NOT include a system snapshot. Installed apps will NOT be reverted.%s\n", ColorYellow, ColorReset)
	}

	// 3. Restore Memory
	agent.Messages = make([]Message, len(cp.Messages))
	copy(agent.Messages, cp.Messages)
	totalTokens = cp.TotalTokens

	// If the last message is an assistant message with tool calls, remove it.
	// This prevents "invalid_request_message_order" errors because the tool output is missing (since we rolled back before it happened).
	// By removing it, we revert to the state just before the agent decided to take the action.
	if len(agent.Messages) > 0 {
		lastMsg := agent.Messages[len(agent.Messages)-1]
		if lastMsg.Role == "assistant" && len(lastMsg.ToolCalls) > 0 {
			agent.Messages = agent.Messages[:len(agent.Messages)-1]
			// We might need to adjust totalTokens, but since we just restored exact snapshot,
			// the tokens for that dropped message are technically "lost" or "not spent yet".
			// Ideally we would subtract them, but we don't know exact count for just that message easily without recalculating.
			// However, keeping the token count high is safer than underestimating.
			// Or we could trust the next API call to correct the context usage stats.
			fmt.Println("Note: Reverted pending tool call from conversation history.")
		}
	}

	return nil
}

// deleteCheckpoint deletes a checkpoint
func deleteCheckpoint(agentID, checkpointID string) error {
	path := filepath.Join(getCheckpointsDir(agentID), checkpointID+".json")
	return os.Remove(path)
}

func pruneAutoCheckpoints(agentID string) {
	cps, err := listCheckpoints(agentID)
	if err != nil {
		return
	}

	var autoCPs []Checkpoint
	for _, cp := range cps {
		if cp.IsAuto {
			autoCPs = append(autoCPs, cp)
		}
	}

	// Keep last 10
	if len(autoCPs) > 10 {
		for i := 10; i < len(autoCPs); i++ {
			deleteCheckpoint(agentID, autoCPs[i].ID)
		}
	}
}

// Helper to check if running in Docker
func isRunningInDocker() bool {
	_, err := os.Stat("/.dockerenv")
	return err == nil
}

// Helper to get current container ID (best effort)
func getCurrentContainerID() (string, error) {
	// Hostname is usually the container ID in default docker run
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return hostname, nil
}

// Helper to check docker access
func checkDockerAccess() error {
	cmd := exec.Command("docker", "ps")
	return cmd.Run()
}
