package main

import (
	"fmt"
	"os/exec"
	"runtime"
)

// executeCommand выполняет команду оболочки и возвращает ее вывод.
func executeCommand(command string) (string, error) {
	fmt.Printf("\033[31m$ %s\033[0m\n", command)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command execution failed: %s, output: %s", err, string(output))
	}

	return string(output), nil
}