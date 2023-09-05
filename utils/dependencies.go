package utils

import (
	"encoding/json"
	"path"

	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
)

type dependsFileFormat struct {
	DependsOn []string
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
