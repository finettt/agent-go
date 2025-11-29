package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// confirmAndExecute checks the execution mode and prompts for confirmation if necessary.
func confirmAndExecute(config *Config, command string) (string, error) {
	if config.ExecutionMode == Ask {
		// The command is already printed as part of the tool call, so we just ask for confirmation.
		fmt.Printf("%s%s%s\nDo you want to execute the above command? [y/N]: ", ColorRed, command, ColorReset)

		var response string
		fmt.Scanln(&response) // This is safer than bufio.NewReader with the readline library.

		if strings.ToLower(strings.TrimSpace(response)) != "y" {
			return "Command not executed by user.", nil
		}
	}
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
