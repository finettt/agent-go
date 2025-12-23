package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

var executionMutex sync.Mutex

type BackgroundProcess struct {
	PID       int
	Command   string
	StartTime time.Time
	Cmd       *exec.Cmd
	Output    *bytes.Buffer
}

var (
	backgroundProcesses = make(map[int]*BackgroundProcess)
	bgMutex             sync.Mutex
	processIDCounter    = 1
)

// confirmAndExecute checks the execution mode and prompts for confirmation if necessary.
//
// Background execution is not agent-controlled. In Ask mode, the user can choose to run the
// command in the foreground or start it as a background process.
func confirmAndExecute(config *Config, command string) (string, error) {
	// Check operation mode first
	if config.OperationMode == Plan {
		return "", fmt.Errorf("command execution is blocked in Plan mode. Switch to Build mode to execute commands")
	}

	// We need to lock here because multiple sub-agents might try to execute commands
	// or ask for confirmation simultaneously, which would mess up the console I/O.
	executionMutex.Lock()
	defer executionMutex.Unlock()

	if config.ExecutionMode == Ask {
		// The command is already printed as part of the tool call, so we just ask for confirmation.
		fmt.Printf("%s$ %s%s\n%s?%s Execute? [y=foreground/b=background/N]: ", ColorCyan, command, ColorReset, ColorHighlight, ColorReset)

		var response string
		fmt.Scanln(&response) // This is safer than bufio.NewReader with the readline library.

		switch strings.ToLower(strings.TrimSpace(response)) {
		case "y", "yes":
			return executeCommand(command)
		case "b", "bg", "background":
			return executeBackgroundCommand(command)
		default:
			return "Command not executed by user.", nil
		}
	}

	// In YOLO mode, commands always execute in the foreground.
	return executeCommand(command)
}

// executeCommand executes a shell command and returns its output
func executeCommand(command string) (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("pwsh", "-CommandWithArgs", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		// Return output even on error - useful for diagnostics
		return outputStr, fmt.Errorf("command execution failed: %w", err)
	}

	return outputStr, nil
}

func executeBackgroundCommand(command string) (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("pwsh", "-CommandWithArgs", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &outBuf

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start background command: %w", err)
	}

	bgMutex.Lock()
	pid := processIDCounter
	processIDCounter++
	backgroundProcesses[pid] = &BackgroundProcess{
		PID:       pid,
		Command:   command,
		StartTime: time.Now(),
		Cmd:       cmd,
		Output:    &outBuf,
	}
	bgMutex.Unlock()

	// Monitor process in a goroutine to clean up when done
	go func(pid int, p *exec.Cmd) {
		p.Wait()
		bgMutex.Lock()
		// We keep the process in the map so we can retrieve logs, but we could mark it as done?
		// For now, let's just keep it. Or maybe we want to remove it?
		// If we remove it, we lose the logs.
		// Let's keep it but maybe mark as finished?
		// The user requirement says "Agent can't stop when it have running backgound commands",
		// so we need to know if it's running.
		// cmd.ProcessState contains information after Wait() returns.
		bgMutex.Unlock()
	}(pid, cmd)

	return fmt.Sprintf("Background command started with PID: %d", pid), nil
}

func killBackgroundCommand(pid int) (string, error) {
	bgMutex.Lock()
	proc, exists := backgroundProcesses[pid]
	bgMutex.Unlock()

	if !exists {
		return "", fmt.Errorf("background process with PID %d not found", pid)
	}

	if proc.Cmd.ProcessState != nil && proc.Cmd.ProcessState.Exited() {
		return "Process already finished", nil
	}

	if err := proc.Cmd.Process.Kill(); err != nil {
		return "", fmt.Errorf("failed to kill process: %w", err)
	}

	return fmt.Sprintf("Process %d killed", pid), nil
}

func getBackgroundLogs(pid int) (string, error) {
	bgMutex.Lock()
	defer bgMutex.Unlock()

	proc, exists := backgroundProcesses[pid]
	if !exists {
		return "", fmt.Errorf("background process with PID %d not found", pid)
	}

	return proc.Output.String(), nil
}

func listBackgroundCommands() string {
	bgMutex.Lock()
	defer bgMutex.Unlock()

	if len(backgroundProcesses) == 0 {
		return "No background commands running."
	}

	var builder strings.Builder
	builder.WriteString("Background Commands:\n")
	for _, proc := range backgroundProcesses {
		status := "Running"
		if proc.Cmd.ProcessState != nil && proc.Cmd.ProcessState.Exited() {
			status = fmt.Sprintf("Finished (Exit Code: %d)", proc.Cmd.ProcessState.ExitCode())
		}
		builder.WriteString(fmt.Sprintf("- PID: %d | Command: %s | Status: %s\n", proc.PID, proc.Command, status))
	}
	return builder.String()
}

func hasRunningBackgroundProcesses() bool {
	bgMutex.Lock()
	defer bgMutex.Unlock()

	for _, proc := range backgroundProcesses {
		if proc.Cmd.ProcessState == nil || !proc.Cmd.ProcessState.Exited() {
			return true
		}
	}
	return false
}

// executeSkill executes a skill command.
// If it's a .sh file, it executes it directly with sh to avoid shell escaping issues.
func executeSkill(command string, argsJSON []byte) (string, error) {
	if strings.HasSuffix(command, ".sh") {
		cmd := exec.Command("sh", command)
		cmd.Env = append(os.Environ(), fmt.Sprintf("SKILL_ARGS=%s", string(argsJSON)))

		var outBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &outBuf

		err := cmd.Run()
		output := outBuf.String()
		if err != nil {
			return output, fmt.Errorf("skill execution failed: %w", err)
		}
		return output, nil
	}

	// Fallback to shell execution
	cmdStr := fmt.Sprintf("export SKILL_ARGS='%s' && %s", string(argsJSON), command)
	return executeCommand(cmdStr)
}
