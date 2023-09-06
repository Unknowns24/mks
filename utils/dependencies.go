package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
)

type dependsFileFormat struct {
	DependsOn []string `json:"dependsOn"`
}

/*
	This functions will return a boolean true if all dependencies are installed on templates addons folder
 	or false if at least one of that dependencies are missing, the second parameter will be the slice with
	the missing ones. Error indicates that something happend during the function execution.
*/
func ValidateAllDependenciesInstalled(dependsFilePath string) (bool, []string, error) {
	// Read file content
	fileContent, err := ReadFile(dependsFilePath)
	if err != nil {
		return false, []string{}, err
	}

	var parsedFile dependsFileFormat

	// Parse json file and save data on parsedFile variable
	err = json.Unmarshal([]byte(fileContent), &parsedFile)
	if err != nil {
		return false, []string{}, err
	}

	notInstalledTemplates := []string{}

	for _, dependency := range parsedFile.DependsOn {
		templatePath := path.Join(global.TemplatesFolderPath, config.FOLDER_ADDONS, dependency)

		if !FileOrDirectoryExists(templatePath) {
			notInstalledTemplates = append(notInstalledTemplates, dependency)
		}
	}

	if len(notInstalledTemplates) > 0 {
		return false, notInstalledTemplates, nil
	}

	return true, []string{}, nil
}

func loadDependencies(dependsFilePath string) (dependsFileFormat, error) {
	var dependencies dependsFileFormat

	data, err := os.ReadFile(dependsFilePath)
	if err != nil {
		return dependencies, err
	}

	err = json.Unmarshal(data, &dependencies)
	if err != nil {
		return dependencies, err
	}

	return dependencies, nil
}

// deepFirstSearch to perform the topological sort
func deepFirstSearch(node string, visited map[string]bool, stack *[]string, dependencies map[string][]string) bool {
	visited[node] = true
	for _, neighbor := range dependencies[node] {
		if !visited[neighbor] {
			if !deepFirstSearch(neighbor, visited, stack, dependencies) {
				return false // cycle detected
			}
		} else if !SliceContainsElement(*stack, neighbor) {
			// If the neighbor has been visited but is not in the stack, it means there's a cycle
			return false
		}
	}

	// Add the node to the output stack
	*stack = append(*stack, node)
	return true
}

func GetDependenciesInstallationOrder(filePaths []string) ([]string, error) {
	dependencies := make(map[string][]string)

	// Load the dependencies of each file into the dependencies map
	for _, filename := range filePaths {
		t, err := loadDependencies(filename)
		if err != nil {
			return []string{}, fmt.Errorf("error loading dependencies: %s", err)
		}
		dependencies[filename] = t.DependsOn
	}

	visited := make(map[string]bool)
	var stack []string

	for node := range dependencies {
		if !visited[node] {
			if !deepFirstSearch(node, visited, &stack, dependencies) {
				return []string{}, errors.New("a cycle in dependencies has been detected")
			}
		}
	}

	var finalStack []string

	// The stack now contains the order in which the templates should be included
	for i := len(stack) - 1; i >= 0; i-- {
		finalStack = append(finalStack, stack[i])
	}

	return finalStack, nil
}
