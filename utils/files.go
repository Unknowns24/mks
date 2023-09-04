package utils

import (
	"fmt"
	"io/ioutil"
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

// Read file content and return it
func ReadFile(filePath string) (string, error) {
	// Read template content
	templateContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(templateContent), nil
}

// Read file content, replace placeholder with its values in memory and returns modified content
func ReadFileWithCustomReplace(filePath string, replaces map[string]string) (string, error) {
	// Read template content
	templateContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	fileContent := string(templateContent)

	// Replace every placeholders with its value
	for flag, value := range replaces {
		fileContent = strings.ReplaceAll(string(fileContent), flag, value)
	}

	return fileContent, nil
}

// Function to extend file content by adding new content at the bottom of the file
func ExtendFile(filePath, newContent string) error {
	// Read the current content of the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Insert the new configuration before the closing of the Config structure
	finalContent := fmt.Sprintf("%s\n\n%s", content, newContent)

	// Write the updated content to the env file
	err = ioutil.WriteFile(filePath, []byte(finalContent), 0644)
	if err != nil {
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
