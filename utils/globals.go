package utils

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
)

func SetTemplatesFolderPathGlobal() error {
	// Get the directory path of the current file (generator.go)
	mksDir, err := GetExecutablePath()
	if err != nil {
		return fmt.Errorf("failed to get current file path")
	}

	// Save in a global variable the path to templates folder inside MKS
	global.TemplatesFolderPath = filepath.Join(mksDir, config.FOLDER_LIBS, config.FOLDER_TEMPLATES)

	return nil
}

func SetCurrentInstalledTemplates() error {
	// Get installed templates
	addonsPath := path.Join(global.TemplatesFolderPath, config.FOLDER_ADDONS)
	installedTemplates, err := ListDirectories(addonsPath)
	if err != nil {
		return err
	}

	global.InstalledTemplates = installedTemplates
	return nil
}

func SetExecutablePath() error {
	// Get the directory path of the current file (generator.go)
	mksDir, err := GetExecutablePath()
	if err != nil {
		return fmt.Errorf("failed to get current file path")
	}

	// Save in a global variable the path to templates folder inside MKS
	global.ExecutableBasePath = mksDir

	return nil
}
