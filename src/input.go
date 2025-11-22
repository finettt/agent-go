package main

import (
	"fmt"
	"os"
	"regexp"
)

// processFileMentions scans the input string for @filename patterns
// and replaces them with the content of the referenced file.
func processFileMentions(input string) string {
	// Regex to find @filename patterns.
	// We look for @ followed by valid file path characters.
	// Allowed: alphanumeric, _, -, ., /, \
	re := regexp.MustCompile(`@([\w\-\./\\]+)`)

	return re.ReplaceAllStringFunc(input, func(match string) string {
		// match includes the @
		filename := match[1:]

		// Check if file exists and is a file (not directory)
		info, err := os.Stat(filename)
		if err != nil {
			// If file not found or error, warn user and keep original text
			// Only warn if it looks like a specific file request, to avoid noise
			// on accidental @ usage, but here explicit @filename implies intent.
			fmt.Fprintf(os.Stderr, "%sWarning: could not find file '%s': %v%s\n", ColorYellow, filename, err, ColorReset)
			return match
		}
		if info.IsDir() {
			fmt.Fprintf(os.Stderr, "%sWarning: '%s' is a directory, not a file%s\n", ColorYellow, filename, ColorReset)
			return match
		}

		// Read the file
		content, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%sWarning: could not read file '%s': %v%s\n", ColorYellow, filename, err, ColorReset)
			return match
		}

		// Return formatted content
		return fmt.Sprintf("\nFile '%s':\n```\n%s\n```\n", filename, string(content))
	})
}
