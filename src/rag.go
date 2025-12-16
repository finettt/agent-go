package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Supported text file extensions for RAG search
var textFileExtensions = map[string]bool{
	".txt":  true,
	".md":   true,
	".go":   true,
	".py":   true,
	".js":   true,
	".ts":   true,
	".json": true,
	".yaml": true,
	".yml":  true,
	".xml":  true,
	".html": true,
	".css":  true,
	".sh":   true,
	".bat":  true,
}

// searchRAGFiles searches files in the given path and returns relevant snippets
func searchRAGFiles(path, query string, maxSnippets int) (string, error) {
	var snippets []string
	queryTokens := strings.Fields(strings.ToLower(query))

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Stop early if we have enough snippets
		if len(snippets) >= maxSnippets {
			return filepath.SkipAll
		}

		if info.IsDir() {
			return nil
		}

		// Filter by file extension to avoid binary files
		ext := filepath.Ext(filePath)
		if !textFileExtensions[ext] {
			return nil
		}

		file, err := os.Open(filePath)
		if err != nil {
			return nil // Skip files that cannot be opened
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() && len(snippets) < maxSnippets {
			line := scanner.Text()
			if containsAny(strings.ToLower(line), queryTokens) {
				snippets = append(snippets, fmt.Sprintf("- %s (from %s)", line, filepath.Base(filePath)))
			}
		}

		if err := file.Close(); err != nil {
			// Log or handle the error if closing fails
			fmt.Printf("Warning: failed to close file %s: %v\n", filePath, err)
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return strings.Join(snippets, "\n"), nil
}

// containsAny checks if a string contains at least one of the tokens
func containsAny(s string, tokens []string) bool {
	for _, token := range tokens {
		if strings.Contains(s, token) {
			return true
		}
	}
	return false
}
