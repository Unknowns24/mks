package utils

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/unknowns24/mks/config"
)

// Function to add a new configuration before the closing of the Config structure in the source file
func AddGoConfigFromString(newConfig string) error {
	// Path to the source file
	filePath := filepath.Join(config.BasePath, config.FOLDER_SRC, config.FOLDER_UTILS, config.FILE_GO_CONFIG)

	// Read the current content of the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var configClosingBraceIndex int

	// Find the line that defines the Config structure
	for i, line := range lines {
		if strings.Contains(line, "Config") && strings.Contains(line, "struct") {
			configClosingBraceIndex = findClosingBrace(lines, i)
			break
		}
	}

	if configClosingBraceIndex == -1 {
		return fmt.Errorf("could not find closing brace for the Config structure")
	}

	// Insert the new configuration before the closing of the Config structure
	newContent := fmt.Sprintf("%s\n%s\n%s", strings.Join(lines[:configClosingBraceIndex], "\n"), newConfig, strings.Join(lines[configClosingBraceIndex:], "\n"))

	// Write the updated content to the file
	err = ioutil.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

// Function to add a new configuration before the closing of the Config structure in the source file
func AddEnvConfigFromString(newConfig string) error {
	// Path to the source file
	filePath := filepath.Join(config.BasePath, config.FILE_ENV_CONFIG)

	// Read the current content of the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Insert the new configuration before the closing of the Config structure
	newContent := fmt.Sprintf("%s\n%s", content, newConfig)

	// Write the updated content to the file
	err = ioutil.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		return err
	}

	return nil
}
