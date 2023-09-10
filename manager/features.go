package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
	"github.com/unknowns24/mks/utils"
)

type InstalledFeaturesFileFormat struct {
	Features []string `json:"features"`
}

func GetApplicationInstalledFeatures() ([]string, error) {
	// Read file content
	fileContent, err := utils.ReadFile(path.Join(global.BasePath, config.FOLDER_MKS_MODULES, config.FILE_MKS_INSTALLED_FEATURES))
	if err != nil {
		return nil, err
	}

	// Variable to save parsed json data
	var parsedFile InstalledFeaturesFileFormat

	// Parse json file and save data on parsedFile variable
	err = json.Unmarshal([]byte(fileContent), &parsedFile)
	if err != nil {
		return nil, err
	}

	return parsedFile.Features, nil
}

func IsValidFeature(feature string) bool {
	return utils.FileOrDirectoryExists(path.Join(global.UserTemplatesFolderPath, feature))
}

func FeatureHasLoadFile(feature string) bool {
	return utils.FileOrDirectoryExists(path.Join(global.UserTemplatesFolderPath, feature, config.FILE_ADDON_TEMPLATE_MAIN_LOAD))
}

func FeatureHasUnloadFile(feature string) bool {
	return utils.FileOrDirectoryExists(path.Join(global.UserTemplatesFolderPath, feature, config.FILE_ADDON_TEMPLATE_MAIN_UNLOAD))
}

func AddAllFeatures() error {
	for _, feature := range global.InstalledTemplates {
		err := AddFeature(feature)
		if err != nil {
			return err
		}
	}

	return nil
}

func AddFeature(feature string) error {
	var err error

	// If global variable ApplicationName is empty fill it
	if global.ApplicationName == "" {
		// Get Application module name
		global.ApplicationName, err = utils.GetThisModuleName()
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

	// Check if requested feature is installed
	if !IsValidFeature(feature) {
		return fmt.Errorf("unknown feature: %s", feature)
	}

	err = InstallFeature(path.Join(global.UserTemplatesFolderPath, feature))
	if err != nil {
		return err
	}

	return nil
}

func InstallFeature(templatePath string) error {
	if global.Verbose {
		fmt.Printf("[+] Validating %s's template to install..\n", templatePath)
	}

	templateName := filepath.Base(templatePath)
	dependsFilePath := path.Join(templatePath, config.FILE_ADDON_TEMPLATE_DEPENDS)

	if global.Verbose {
		fmt.Println("[+] Checking if template has dependency file..")
	}

	var dependenciesInOrder []string

	// Validate if exists a depends file
	if utils.FileOrDirectoryExists(dependsFilePath) {
		allDependenciesInstalled, missingDependencies, err := utils.ValidateAllDependenciesInstalled(dependsFilePath)
		if err != nil {
			return err
		}

		if !allDependenciesInstalled {
			if len(missingDependencies) > 1 {
				return fmt.Errorf(`missing dependencies on %s: %s are not installed"`, templateName, strings.Join(missingDependencies, ", "))
			}

			return fmt.Errorf(`missing dependency on %s: %s is not installed"`, templateName, missingDependencies)
		}

		if global.Verbose {
			fmt.Println("[+] Parsing dependency file..")
		}

		dependenciesInOrder, err = utils.GetDependenciesInstallationOrder(dependsFilePath)
		if err != nil {
			return err
		}

	}

	if len(dependenciesInOrder) > 0 {
		continueInstallation, err := utils.AskConfirm(fmt.Sprintf("%s has this dependencies: %s. Do you want to install it?", templateName, strings.Join(dependenciesInOrder, ", ")))
		if err != nil {
			return err
		}

		if !continueInstallation {
			return errors.New("installation interrumped by user")
		}
	}

	if global.Verbose {
		fmt.Println("[+] Creating a temporal folder to work..")
	}

	// Create temporal directory to prevent make a mess on the current application
	temporalDirectoryPath, err := utils.MakeTemporalDirectory()
	if err != nil {
		return err
	}

	if global.Verbose {
		fmt.Println("[+] Copying the application to the temporal folder..")
	}

	// Copy application on the temporal directory
	AppFolderContent, err := utils.ListDirectoriesAndFiles(global.BasePath)
	if err != nil {
		utils.DeleteFileOrDirectory(temporalDirectoryPath)
		return err
	}

	for _, fileOrDirectory := range AppFolderContent {
		err := utils.CopyFileOrDirectory(path.Join(global.BasePath, fileOrDirectory), path.Join(temporalDirectoryPath, fileOrDirectory))
		if err != nil {
			utils.DeleteFileOrDirectory(temporalDirectoryPath)
			return err
		}
	}

	// 	Check if mks_modules app is already created on the application (if not exists create it)
	mksModulesFolderPath := path.Join(temporalDirectoryPath, config.FOLDER_MKS_MODULES)
	if !utils.FileOrDirectoryExists(mksModulesFolderPath) {
		if global.Verbose {
			fmt.Println("[+] Copying mks_modules to the application..")
		}

		err := utils.CopyFileOrDirectory(path.Join(global.MksTemplatesFolderPath, config.FOLDER_MKS_MODULES), mksModulesFolderPath)
		if err != nil {
			return err
		}
	}

	if global.Verbose {
		fmt.Printf("[+] Preparing to install %s dependencies template..\n", templateName)
	}

	for _, dependencyTemplateName := range dependenciesInOrder {
		err := ImportFeatureToApp(path.Join(global.UserTemplatesFolderPath, dependencyTemplateName), temporalDirectoryPath)
		if err != nil {
			return fmt.Errorf(`error on %s's "%s" dependency installation: %s"`, templateName, dependencyTemplateName, err)
		}
	}

	err = ImportFeatureToApp(templatePath, temporalDirectoryPath)
	if err != nil {
		return fmt.Errorf(`error on %s installation: %s"`, templateName, err)
	}

	return nil
}

func ImportFeatureToApp(templatePath, workingDirectory string) error {
	templateName := filepath.Base(templatePath)

	appInstalledFeatures, err := GetApplicationInstalledFeatures()
	if err != nil {
		return err
	}

	if global.Verbose {
		fmt.Printf("[+] Preparing %s template to be imported to the application..\n", templateName)
	}

	// if feature already installed pass it
	if utils.SliceContainsElement(appInstalledFeatures, templateName) {
		if global.Verbose {
			fmt.Printf("[+] %s already installed in the application..", templateName)
		}

		return nil
	}

	// Check if template has prompt files
	if global.Verbose {
		fmt.Printf("[+] Check if %s has a prompt file..\n", templateName)
	}

	templatePromptsFile := path.Join(templatePath, config.FILE_ADDON_TEMPLATE_PROMPTS)

	placeHoldersToReplace := map[string]string{
		config.PLACEHOLDER_APP_NAME: global.ApplicationName,
	}

	if utils.FileOrDirectoryExists(templatePromptsFile) {
		err := utils.ParsePromptFile(templatePromptsFile, &placeHoldersToReplace)
		if err != nil {
			return fmt.Errorf(`error parsing %s: %s"`, templatePromptsFile, err)
		}
	}

	return nil
}
