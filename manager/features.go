package manager

import (
	"fmt"
	"os"
	"path"

	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
)

func IsValidFeature(feature string) (bool, error) {
	// Check if requested feature is installed
	for _, template := range global.InstalledFeatures {
		if template == feature {
			return true, nil
		}
	}

	return false, fmt.Errorf("unknown feature: %s", feature)
}

func AddFeature(feature string) error {
	var err error

	// If global variable basePath is empty fill it
	if global.BasePath == "" {
		// Get the current working directory
		global.BasePath, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	// If global variable serviceName is empty fill it
	if global.ServiceName == "" {
		// Get Mircoservice module name
		global.ServiceName, err = GetThisModuleName()
		if err != nil {
			return err
		}
	}

	// Validating feature
	addonsPath := path.Join(global.TemplatesFolderPath, config.FOLDER_ADDONS)
	founded := false

	// Check if requested feature is installed
	for _, template := range global.InstalledFeatures {
		if template != feature {
			continue
		}

		founded = true

		err = InstallFeature(path.Join(addonsPath, template))
		if err != nil {
			return err
		}
	}

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
	// To Implement
	return nil
}
