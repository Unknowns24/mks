package manager

import (
	"fmt"
	"os"
	"path"

	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
	"github.com/unknowns24/mks/utils"
)

func IsValidFeature(feature string) bool {
	// Check if requested feature is installed
	for _, template := range global.InstalledTemplates {
		if template == feature {
			return true
		}
	}

	return false
}

func AddFeature(feature string) error {
	var err error

	// If global variable serviceName is empty fill it
	if global.ServiceName == "" {
		// Get Mircoservice module name
		global.ServiceName, err = GetThisModuleName()
		if err != nil {
			return err
		}
	}

	// If global variable basePath is empty fill it
	if global.BasePath == "" {
		// Get the current working directory
		global.BasePath, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	// Validating feature
	addonsPath := path.Join(global.TemplatesFolderPath, config.FOLDER_ADDONS)
	founded := false

	// Check if requested feature is installed
	for _, template := range global.InstalledTemplates {
		if template != feature {
			continue
		}

		founded = true

		err = InstallFeature(path.Join(addonsPath, template))
		if err != nil {
			return err
		}
	}

	// This should never happen because commands use IsValidFeature
	if !founded {
		return fmt.Errorf("unknown feature: %s", feature)
	}

	return nil
}

func AddAllFeatures() error {
	// To implement
	return nil
}

func InstallFeature(templatePath string) error {
	fmt.Println(templatePath)

	dependsFilePath := path.Join(templatePath, config.FILE_ADDON_TEMPLATE_DEPENDS)

	// Validate if exists a depends file
	if utils.FileOrDirectoryExists(dependsFilePath) {
		dependencies, err := utils.GetDependenciesInstallationOrder(dependsFilePath)
		if err != nil {
			return err
		}

		fmt.Println(dependencies)
	}

	return nil
}

func GetApplicationInstalledFeatures() error {
	return nil
}
