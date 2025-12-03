package main

import (
	"fmt"
	"os"
	"regexp"
)

// processFileMentions scans the input string for @filename and #note patterns
// and replaces them with the content of the referenced file or note.
func processFileMentions(input string) string {
	// Regex to find @filename patterns.
	// We look for @ followed by valid file path characters.
	// Allowed: alphanumeric, _, -, ., /, \
	fileRe := regexp.MustCompile(`@([\w\-\./\\]+)`)

	input = fileRe.ReplaceAllStringFunc(input, func(match string) string {
		// match includes the @
		filename := match[1:]

		// Check if file exists and is a file (not directory)
		info, err := os.Stat(filename)
		if err != nil {
			// If file not found or error, warn user and keep original text
			// Only warn if it looks like a specific file request, to avoid noise
			// on accidental @ usage, but here explicit @filename implies intent.
			fmt.Fprintf(os.Stderr, "%sWarning: could not find file '%s': %v%s\n", ColorMeta, filename, err, ColorReset)
			return match
		}
		if info.IsDir() {
			fmt.Fprintf(os.Stderr, "%sWarning: '%s' is a directory, not a file%s\n", ColorMeta, filename, ColorReset)
			return match
		}

		// Read the file
		content, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%sWarning: could not read file '%s': %v%s\n", ColorMeta, filename, err, ColorReset)
			return match
		}

		// Return formatted content
		return fmt.Sprintf("\nFile '%s':\n```\n%s\n```\n", filename, string(content))
	})

	// Regex to find #note patterns.
	// We look for # followed by valid note name characters.
	// Allowed: alphanumeric, _, -, .
	noteRe := regexp.MustCompile(`#([\w\-\.]+)`)

	return noteRe.ReplaceAllStringFunc(input, func(match string) string {
		// match includes the #
		noteName := match[1:]

		note, err := loadNote(noteName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%sWarning: could not find note '%s': %v%s\n", ColorMeta, noteName, err, ColorReset)
			return match
		}

		return fmt.Sprintf("\nNote '%s':\n```\n%s\n```\n", note.Name, note.Content)
	})
}
