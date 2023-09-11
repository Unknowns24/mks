package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

func GetThisModuleName() (string, error) {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Create the go.mod file path
	goModPath := filepath.Join(currentDir, "go.mod")

	// Check if go.mod file exists
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return "", errors.New("the go.mod file does not exist")
	}

	// Read the content of the go.mod file
	goModContent, err := os.ReadFile("go.mod")
	if err != nil {
		return "", fmt.Errorf("error reading the go.mod file: %s", err)
	}

	// Use a regular expression to find the module name
	modulePattern := regexp.MustCompile(`module ([^\s]+)`)
	match := modulePattern.FindStringSubmatch(string(goModContent))
	if len(match) < 2 {
		return "", errors.New("could not find the module name in go.mod")
	}

	return match[1], nil
}

func InitGoModules(ApplicationName, basePath string) error {
	initCmd := exec.Command("go", "mod", "init", ApplicationName)
	initCmd.Dir = basePath
	initOutput, err := initCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("¨error running 'go mod init': %s\nOutput: %s", err, initOutput)
	}

	return nil
}

func InstallNeededPackages(basePath string) error {
	initCmd := exec.Command("go", "mod", "tidy")
	initCmd.Dir = basePath
	initOutput, err := initCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running 'go mod tidy': %s\nOutput: %s", err, initOutput)
	}

	return nil
}
