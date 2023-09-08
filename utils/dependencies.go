package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
)

type dependsFileFormat struct {
	DependsOn []string `json:"dependsOn"`
}

var processedDependencies = make(map[string]bool)
var currentProcessing = make(map[string]bool)

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
		templatePath := path.Join(global.UserTemplatesFolderPath, dependency)

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

func processDependencies(filePath string, result *[]string) error {
	// Load the dependencies of the current file
	featureDependencies, err := loadDependencies(filePath)
	if err != nil {
		return err
	}

	for _, feature := range featureDependencies.DependsOn {
		// If this dependency is currently being processed, it indicates a cycle
		if currentProcessing[feature] {
			return fmt.Errorf("cyclic redundancy detected on %s", feature)
		}

		// If this dependency has already been processed, move to the next
		if processedDependencies[feature] {
			continue
		}

		// Mark the dependency as currently being processed
		currentProcessing[feature] = true

		// Path to the dependency file of the current feature
		dependencyFilePath := path.Join(global.UserTemplatesFolderPath, feature, config.FILE_ADDON_TEMPLATE_DEPENDS)

		// If a dependency file exists for this feature, process it recursively
		if FileOrDirectoryExists(dependencyFilePath) {
			err := processDependencies(dependencyFilePath, result)
			if err != nil {
				return err
			}
		}

		// Mark the dependency as processed
		processedDependencies[feature] = true
		delete(currentProcessing, feature)

		// Add the current feature to the result
		*result = append(*result, feature)
	}

	return nil
}

func GetDependenciesInstallationOrder(dependencyFilePath string) ([]string, error) {
	var result []string
	fatherFeatureName := filepath.Base(filepath.Dir(dependencyFilePath))

	err := processDependencies(dependencyFilePath, &result)
	if err != nil {
		return nil, err
	}

	if SliceContainsElement(result, fatherFeatureName) {
		return nil, errors.New("cycle dependency detected, posible incompatibility on some of your templates")
	}

	result = append(result, fatherFeatureName)

	return result, nil
}
