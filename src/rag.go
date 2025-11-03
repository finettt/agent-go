package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// searchRAGFiles ищет в файлах по указанному пути и возвращает сниппеты.
func searchRAGFiles(path, query string, maxSnippets int) (string, error) {
	var snippets []string
	queryTokens := strings.Fields(strings.ToLower(query))

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return nil // Пропускаем файлы, которые не можем открыть
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if containsAny(strings.ToLower(line), queryTokens) {
					snippets = append(snippets, fmt.Sprintf("- %s (from %s)", line, filepath.Base(filePath)))
					if len(snippets) >= maxSnippets {
						return nil // Достигли лимита
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return strings.Join(snippets, "\n"), nil
}

// containsAny проверяет, содержит ли строка хотя бы один из токенов.
func containsAny(s string, tokens []string) bool {
	for _, token := range tokens {
		if strings.Contains(s, token) {
			return true
		}
	}
	return false
}