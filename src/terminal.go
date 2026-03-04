package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/creack/pty"
)

// Terminal session storage
var (
	terminalSessions   = make(map[string]*TerminalSession)
	terminalSessionID  = 1
	terminalSessionMux sync.Mutex
)

// OpenTerminalSessionArgs represents arguments for opening a terminal session
type OpenTerminalSessionArgs struct {
	// Command is optional - if provided, it will be executed after bash starts
	Command string `json:"command,omitempty"`
}

// openTerminalSession starts a new terminal session with bash
// All sessions start with bash by default for stability
// If a command is provided, it will be sent to bash after the session starts
func openTerminalSession(argsJSON string) (string, error) {
	var args OpenTerminalSessionArgs
	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}

	// Security: validate command if provided
	if args.Command != "" {
		if err := validateCommand(args.Command); err != nil {
			return "", fmt.Errorf("command validation failed: %w", err)
		}
	}

	terminalSessionMux.Lock()
	defer terminalSessionMux.Unlock()

	// Check max concurrent sessions (10)
	if len(terminalSessions) >= 10 {
		return "", fmt.Errorf("maximum concurrent sessions (10) reached")
	}

	// Always start with bash for stability - sessions don't close when a command finishes
	var cmd *exec.Cmd
	cmd = exec.Command("bash")

	// Start the command with a PTY
	ptyFile, err := pty.Start(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to start PTY: %w", err)
	}

	// Generate session ID
	sessionID := fmt.Sprintf("sess_%d", terminalSessionID)
	terminalSessionID++

	// Create session
	session := &TerminalSession{
		ID:           sessionID,
		PID:          cmd.Process.Pid,
		StartTime:    time.Now(),
		StdoutBuffer: &bytes.Buffer{},
		InputChan:    make(chan string, 100),
		DoneChan:     make(chan struct{}),
		PTY:          ptyFile,
	}

	// Store session
	terminalSessions[sessionID] = session

	// Read PTY output in background
	go func() {
		defer func() {
			ptyFile.Close()
			close(session.DoneChan)
			terminalSessionMux.Lock()
			delete(terminalSessions, sessionID)
			terminalSessionMux.Unlock()
		}()

		buf := make([]byte, 4096)
		for {
			n, err := ptyFile.Read(buf)
			if n > 0 {
				terminalSessionMux.Lock()
				session.StdoutBuffer.Write(buf[:n])
				terminalSessionMux.Unlock()
			}
			if err != nil {
				if err != io.EOF {
					// Suppress error logging for normal session termination
					// fmt.Fprintf(os.Stderr, "PTY read error: %v\n", err)
				}
				return
			}
		}
	}()

	// Send input to PTY in background
	go func() {
		for input := range session.InputChan {
			_, err := ptyFile.Write([]byte(input))
			if err != nil {
				// Suppress error logging for normal session termination
				// fmt.Fprintf(os.Stderr, "PTY write error: %v\n", err)
			}
		}
	}()

	// If a command was provided, send it to bash after a brief delay to let bash start
	initialCommand := ""
	if args.Command != "" {
		initialCommand = args.Command
		// Give bash a moment to start, then send the command
		go func() {
			time.Sleep(100 * time.Millisecond)
			session.InputChan <- args.Command + "\n"
		}()
	}

	result := map[string]interface{}{
		"session_id": sessionID,
		"pid":        cmd.Process.Pid,
		"shell":      "bash",
	}
	if initialCommand != "" {
		result["initial_command"] = initialCommand
	}
	resultJSON, _ := json.Marshal(result)
	return string(resultJSON), nil
}

// sendTerminalInput sends input to a terminal session and returns new output
func sendTerminalInput(argsJSON string) (string, error) {
	var args TerminalInputArgs
	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}

	if strings.TrimSpace(args.SessionID) == "" {
		return "", fmt.Errorf("session_id cannot be empty")
	}

	if strings.TrimSpace(args.Input) == "" {
		return "", fmt.Errorf("input cannot be empty")
	}

	terminalSessionMux.Lock()
	session, exists := terminalSessions[args.SessionID]
	// Capture the buffer length before sending input
	bufferLenBefore := 0
	if exists {
		bufferLenBefore = session.StdoutBuffer.Len()
	}
	terminalSessionMux.Unlock()

	if !exists {
		return "", fmt.Errorf("session '%s' not found", args.SessionID)
	}

	// Parse keycodes
	input := parseKeycode(args.Input)

	// Check session timeout (5 minutes max)
	if time.Since(session.StartTime) > 5*time.Minute {
		return "", fmt.Errorf("session has expired (5 minute timeout)")
	}

	// Send input
	select {
	case session.InputChan <- input:
		// Wait a brief moment for the terminal to process and generate output
		time.Sleep(150 * time.Millisecond)

		// Read any new output that appeared after sending input
		terminalSessionMux.Lock()
		allOutput := session.StdoutBuffer.String()
		terminalSessionMux.Unlock()

		// Get only the new output (after the position we captured before)
		newOutput := ""
		if len(allOutput) > bufferLenBefore {
			newOutput = allOutput[bufferLenBefore:]
		}

		// Build result with status and new output
		result := map[string]interface{}{
			"status":     "sent",
			"session_id": args.SessionID,
		}
		if newOutput != "" {
			result["new_output"] = newOutput
		}
		resultJSON, _ := json.Marshal(result)
		return string(resultJSON), nil

	case <-time.After(2 * time.Second):
		return "", fmt.Errorf("failed to send input (timeout)")
	}
}

// readTerminalOutput reads output from a terminal session
func readTerminalOutput(argsJSON string) (string, error) {
	var args TerminalReadArgs
	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}

	if strings.TrimSpace(args.SessionID) == "" {
		return "", fmt.Errorf("session_id cannot be empty")
	}

	terminalSessionMux.Lock()
	session, exists := terminalSessions[args.SessionID]
	terminalSessionMux.Unlock()

	if !exists {
		return "", fmt.Errorf("session '%s' not found", args.SessionID)
	}

	var output string
	if args.ReadAll {
		output = session.StdoutBuffer.String()
	} else {
		bytesToRead := args.Bytes
		if bytesToRead <= 0 {
			bytesToRead = 4096
		}
		allOutput := session.StdoutBuffer.String()
		if len(allOutput) > bytesToRead {
			output = allOutput[len(allOutput)-bytesToRead:]
		} else {
			output = allOutput
		}
	}

	return fmt.Sprintf("{\"output\": %q}", output), nil
}

// closeTerminalSession closes a terminal session
func closeTerminalSession(argsJSON string) (string, error) {
	var args TerminalCloseArgs
	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}

	if strings.TrimSpace(args.SessionID) == "" {
		return "", fmt.Errorf("session_id cannot be empty")
	}

	terminalSessionMux.Lock()
	defer terminalSessionMux.Unlock()

	session, exists := terminalSessions[args.SessionID]
	if !exists {
		return "", fmt.Errorf("session '%s' not found", args.SessionID)
	}

	// Close PTY (this will signal the read goroutine to exit)
	if session.PTY != nil {
		session.PTY.Close()
	}

	// Close input channel
	close(session.InputChan)

	return "{\"status\": \"closed\"}", nil
}

// listTerminalSessions returns a list of active terminal sessions
func listTerminalSessions() string {
	terminalSessionMux.Lock()
	defer terminalSessionMux.Unlock()

	if len(terminalSessions) == 0 {
		return "No active terminal sessions."
	}

	var builder strings.Builder
	builder.WriteString("Active Terminal Sessions:\n")
	for id, session := range terminalSessions {
		status := "Running"
		if session.PTY == nil {
			status = "Closed"
		}
		builder.WriteString(fmt.Sprintf("- Session: %s | PID: %d | Command: %s | Duration: %v | Status: %s\n",
			id, session.PID, "", time.Since(session.StartTime), status))
	}
	return builder.String()
}

// parseKeycode converts human-readable key codes to terminal escape sequences
func parseKeycode(input string) string {
	// Check if it's a known key mapping
	if mapped, ok := KeyMappings[input]; ok {
		return mapped
	}

	// If not a special key, treat as plain text
	return input
}

// validateCommand checks if a command is safe to execute
func validateCommand(command string) error {
	// Basic validation
	if strings.TrimSpace(command) == "" {
		return fmt.Errorf("command is empty")
	}

	// Check for dangerous shell constructs (basic safety)
	dangerousPatterns := []string{
		"rm -rf /",
		"rm -rf / ",
		"format c:",
		"format c: /",
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(strings.ToLower(command), pattern) {
			return fmt.Errorf("command contains potentially dangerous pattern: %s", pattern)
		}
	}

	return nil
}
