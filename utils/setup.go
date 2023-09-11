package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
)

func SetMksTemplatesFolderPath() error {
	// Get the directory path of the current file (generator.go)
	mksDir, err := GetExecutablePath()
	if err != nil {
		return fmt.Errorf("failed to get current file path")
	}

	// Save in a global variable the path to templates folder inside MKS
	global.MksTemplatesFolderPath = filepath.Join(mksDir, config.FOLDER_TEMPLATES)

	return nil
}

func SetMksDataFolderPath() error {
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

	global.MksDataFolderPath = mksConfigPath
	return nil
}

func SetCurrentInstalledTemplates() error {
	// As this config uses MksDataFolderPath set it if is not declared
	setMksDataIfNotExist()

	userTemplatesPath := path.Join(global.MksDataFolderPath, config.FOLDER_TEMPLATES)

	// Create templates folder on mks app data directory if not exist
	if !FileOrDirectoryExists(userTemplatesPath) {
		err := os.MkdirAll(userTemplatesPath, config.FOLDER_PERMISSION)
		if err != nil {
			return err
		}
	}

	// Get installed templates
	installedTemplates, err := ListDirectories(userTemplatesPath)
	if err != nil {
		return err
	}

	global.InstalledTemplates = installedTemplates
	global.UserTemplatesFolderPath = userTemplatesPath
	return nil
}

func SetMksDataFoldersPath() error {
	// As this config uses MksDataFolderPath set it if is not declared
	setMksDataIfNotExist()

	// Set path for exported files
	global.ExportPath = path.Join(global.MksDataFolderPath, config.FOLDER_EXPORTS)

	// Set cache path for zip files
	global.ZipCachePath = path.Join(global.MksDataFolderPath, config.FOLDER_ZIP_CACHE)

	global.AutoBackupsPath = path.Join(global.MksDataFolderPath, config.FOLDER_MKS_BACKUPS)

	// Set temp path for temporals files
	global.TemporalsPath = path.Join(global.MksDataFolderPath, config.FOLDER_TEMPORALS)

	// Set cache path for templates
	global.TemplateCachePath = path.Join(global.MksDataFolderPath, config.FOLDER_TEMPLATE_CACHE)

	// Create cache, autobackups, temporal and exports folders if not exist
	os.MkdirAll(global.ExportPath, config.FOLDER_PERMISSION)
	os.MkdirAll(global.ZipCachePath, config.FOLDER_PERMISSION)
	os.MkdirAll(global.TemporalsPath, config.FOLDER_PERMISSION)
	os.MkdirAll(global.AutoBackupsPath, config.FOLDER_PERMISSION)
	os.MkdirAll(global.TemplateCachePath, config.FOLDER_PERMISSION)
	os.MkdirAll(global.UserTemplatesFolderPath, config.FOLDER_PERMISSION)

	return nil
}

func setMksDataIfNotExist() error {
	if global.MksDataFolderPath == "" {
		err := SetMksDataFolderPath()
		if err != nil {
			return err
		}
	}

	return nil
}
