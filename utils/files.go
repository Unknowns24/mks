package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/unknowns24/mks/config"
)

// Create files from templates wich only needs to changes %PACKAGE_NAME%
func CreateFileFromTemplate(templatePath, serviceName, finalPath string) error {
	// Read template content
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Replace placeholders with the actual service name
	fileContent := strings.ReplaceAll(string(templateContent), config.PLACEHOLDER_PACKAGENAME, serviceName)

	// Write the file content in the file
	if err := os.WriteFile(finalPath, []byte(fileContent), os.ModePerm); err != nil {
		return err
	}

	return nil
}

// Create files from templates that has many placeholders to replace
func CreateFileFromTemplateWithCustomReplace(templatePath, finalPath string, replaces map[string]string) error {
	// Read template content
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fileContent := string(templateContent)

	// Replace every placeholders with its value
	for flag, value := range replaces {
		fileContent = strings.ReplaceAll(string(fileContent), flag, value)
	}

	// Write the file content in the file
	if err := os.WriteFile(finalPath, []byte(fileContent), os.ModePerm); err != nil {
		return err
	}

	return nil
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
