package main

import (
	"fmt"
	"os/exec"
	"runtime"
)

// executeCommand executes a shell command and returns its output
func executeCommand(command string) (string, error) {
	fmt.Printf("%s$ %s%s\n", ColorRed, command, ColorReset)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
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
