package utils

import (
	"fmt"
	"os"
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
	global.TemplatesFolderPath = filepath.Join(mksDir, config.FOLDER_TEMPLATES)

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

func SetUserConfigFolderPath() error {
	configPath, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("error happend on config directory: %s", err)
	}

	mksConfigPath := path.Join(configPath, config.FOLDER_MKS)

	if !FileOrDirectoryExists(mksConfigPath) {
		err = os.MkdirAll(mksConfigPath, config.FOLDER_PERMISSION)
		if err != nil {
			return err
		}
	}

	global.ConfigFolderPath = mksConfigPath
	return nil
}
