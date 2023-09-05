package utils

import (
	"os"
	"strings"
)

// Read file content and return it
func ReadFile(filePath string) (string, error) {
	// Read template content
	templateContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(templateContent), nil
}

func ListDirectories(dirPath string) ([]string, error) {
	var directories []string

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			directories = append(directories, entry.Name())
		}
	}

	return directories, nil
}

func findClosingBrace(lines []string, startIndex int) int {
	braceCount := 0

	for i := startIndex; i < len(lines); i++ {
		line := lines[i]
		braceCount += strings.Count(line, "{")
		braceCount -= strings.Count(line, "}")

		if braceCount == 0 {
			return i
		}
	}

	return -1
}
