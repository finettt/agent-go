package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// processFileMentions scans the input string for @filename and #note patterns
// and replaces them with the content of the referenced file or note.
func processFileMentions(input string) string {
	// Regex to find @filename patterns.
	// We look for @ followed by valid file name characters (no path separators or dots for traversal).
	// Allowed: alphanumeric, _, -
	// This prevents path traversal and absolute path access.
	fileRe := regexp.MustCompile(`@([\w\-]+)`)

	// Get allowed base directory (current working directory)
	baseDir, err := os.Getwd()
	if err != nil {
		return input
	}
	baseDir = filepath.Clean(baseDir)

	input = fileRe.ReplaceAllStringFunc(input, func(match string) string {
		// match includes the @
		filename := match[1:]

		// Construct full path within base directory
		fullPath := filepath.Join(baseDir, filename)
		fullPath = filepath.Clean(fullPath)

		// Validate path is within base directory (prevent traversal)
		if !strings.HasPrefix(fullPath, baseDir+string(filepath.Separator)) && fullPath != baseDir {
			fmt.Fprintf(os.Stderr, "%sWarning: access denied - path outside allowed directory: %s%s\n", ColorMeta, filename, ColorReset)
			return match
		}

		// Check if file exists and is a file (not directory)
		info, err := os.Stat(fullPath)
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
		content, err := os.ReadFile(fullPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%sWarning: could not read file '%s': %v%s\n", ColorMeta, filename, err, ColorReset)
			return match
		}

		// Return formatted content with just the base filename for display
		return fmt.Sprintf("\nFile '%s':\n```\n%s\n```\n", filepath.Base(filename), string(content))
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
