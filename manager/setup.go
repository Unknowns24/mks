package manager

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
	"github.com/unknowns24/mks/utils"
)

func SetupMks() error {
	// set MksTemplatesFolderPath global variable
	if err := setMksTemplatesFolderPath(); err != nil {
		return err
	}

	// set & create MksDataFolderPath global variable
	if err := setMksDataFolderPath(); err != nil {
		return err
	}

	// Set & create: ExportPath, ZipCachePath, AutoBackupsPath, TemporalsPath, TemplateCachePath global variables
	if err := setMksDataFoldersPath(); err != nil {
		return err
	}

	// Set UserTemplatesFolderPath & InstalledTemplates global variables
	if err := setCurrentInstalledTemplates(); err != nil {
		return err
	}

	return nil
}

func setMksTemplatesFolderPath() error {
	// Get the directory path of the current file (generator.go)
	mksDir, err := utils.GetExecutablePath()
	if err != nil {
		return fmt.Errorf("failed to get current file path")
	}

	// Save in a global variable the path to templates folder inside MKS
	global.MksTemplatesFolderPath = filepath.Join(mksDir, config.FOLDER_TEMPLATES)

	return nil
}

func setMksDataFolderPath() error {
	configPath, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("error happend on config directory: %s", err)
	}

	mksConfigPath := path.Join(configPath, config.FOLDER_MKS)

	if !utils.FileOrDirectoryExists(mksConfigPath) {
		err = os.MkdirAll(mksConfigPath, config.FOLDER_PERMISSION)
		if err != nil {
			return err
		}
	}

	global.MksDataFolderPath = mksConfigPath
	return nil
}

func setCurrentInstalledTemplates() error {
	// As this config uses MksDataFolderPath set it if is not declared
	setMksDataIfNotExist()
	global.UserTemplatesFolderPath = path.Join(global.MksDataFolderPath, config.FOLDER_TEMPLATES)

	// Create templates folder on mks app data directory if not exist
	if !utils.FileOrDirectoryExists(global.UserTemplatesFolderPath) {
		if err := os.MkdirAll(global.UserTemplatesFolderPath, config.FOLDER_PERMISSION); err != nil {
			return err
		}

		if err := InstallTemplate(config.NETWORK_GITHUB_BASE_TEMPLATES_REPO, []string{}); err != nil {
			return err
		}
	}

	// Get installed templates
	installedTemplates, err := utils.ListDirectories(global.UserTemplatesFolderPath)
	if err != nil {
		return err
	}

	global.InstalledTemplates = installedTemplates
	return nil
}

func setMksDataFoldersPath() error {
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
		err := setMksDataFolderPath()
		if err != nil {
			return err
		}
	}

	return nil
}
