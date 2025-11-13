package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

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
	if runtime.GOOS != "linux" {
		return "N/A"
	}

	file, err := os.Open("/etc/os-release")
	if err != nil {
		return "Unknown Linux"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			return strings.Trim(strings.Split(line, "=")[1], `"`)
		}
	}
	return "Unknown Linux"
}
