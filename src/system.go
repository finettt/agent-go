package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

// getCurrentTimeContext returns a formatted string with current time information
// for injection into API requests, providing the AI with temporal awareness.
func getCurrentTimeContext() string {
	now := time.Now()

	// UTC time in ISO 8601 format (matching VS Code extension)
	utcTime := now.UTC().Format(time.RFC3339)

	// Local time in readable format
	localTime := now.Format(time.RFC1123)

	// Get timezone info
	zone, offset := now.Zone()
	offsetHours := offset / 3600
	offsetMinutes := (offset % 3600) / 60

	// Format offset as +HH:MM or -HH:MM
	var offsetStr string
	if offsetMinutes == 0 {
		offsetStr = fmt.Sprintf("UTC%+d:00", offsetHours)
	} else {
		offsetStr = fmt.Sprintf("UTC%+d:%02d", offsetHours, offsetMinutes)
	}

	return fmt.Sprintf("Current Time: %s (UTC) | Local: %s | Timezone: %s (%s)",
		utcTime, localTime, zone, offsetStr)
}

func getSystemInfo() string {
	osName := runtime.GOOS
	arch := runtime.GOARCH
	distro := getDistro()
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "unknown"
	}
	currentTime := time.Now().Format(time.RFC1123)

	return fmt.Sprintf("OS: %s, Architecture: %s, Distribution: %s, CWD: %s, Time: %s", osName, arch, distro, cwd, currentTime)
}

func getDistro() string {
	if runtime.GOOS == "windows" {
		return "Windows"
	}
	if runtime.GOOS != "linux" {
		return "N/A"
	}

	file, err := os.Open("/etc/os-release")
	if err != nil {
		return "Unknown Linux"
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Warning: failed to close file: %v\n", err)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			return strings.Trim(strings.Split(line, "=")[1], `"`)
		}
	}
	return "Unknown Linux"
}
