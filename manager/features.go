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
	templateName := filepath.Base(templatePath)

	if global.Verbose {
		fmt.Printf("[+] Checking if %s's template is already installed in the application..\n", templateName)
	}

	// Get all installed features inside the application
	installedFeatures, err := GetApplicationInstalledFeatures()
	if err != nil {
		return err
	}

	// Check if requested feature is already installed
	if utils.SliceContainsElement(installedFeatures, templateName) {
		return fmt.Errorf("%s's template is already installed", templateName)
	}

	if global.Verbose {
		fmt.Printf("[+] Validating %s's template to install..\n", templateName)
	}

	if global.Verbose {
		fmt.Println("[+] Checking if template has dependency file..")
	}

	dependsFilePath := path.Join(templatePath, config.FILE_ADDON_TEMPLATE_DEPENDS)
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

		// Get the order to install dependencies recursively
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

	applicationTempDir := path.Join(temporalDirectoryPath, filepath.Base(global.BasePath))

	// Copy application on the temporal directory
	err = utils.CopyFileOrDirectory(global.BasePath, applicationTempDir)
	if err != nil {
		utils.DeleteFileOrDirectory(temporalDirectoryPath)
		return err
	}

	// 	Check if mks_modules app is already created on the application (if not exists create it)
	mksModulesFolderPath := path.Join(applicationTempDir, config.FOLDER_MKS_MODULES)

	if !utils.FileOrDirectoryExists(mksModulesFolderPath) {
		return errors.New("this application is not an mks builded application")
	}

	if global.Verbose {
		fmt.Printf("[+] Preparing to install %s dependencies template..\n", templateName)
	}

	// Import all no installed dependencies to the application
	for _, dependencyTemplateName := range dependenciesInOrder {
		if !utils.SliceContainsElement(installedFeatures, dependencyTemplateName) {
			err := ImportFeatureToApp(path.Join(global.UserTemplatesFolderPath, dependencyTemplateName), applicationTempDir)
			if err != nil {
				return fmt.Errorf(`error on %s's "%s" dependency installation: %s"`, templateName, dependencyTemplateName, err)
			}

			installedFeatures = append(installedFeatures, dependencyTemplateName)
		}
	}

	if global.Verbose && len(dependenciesInOrder) > 0 {
		fmt.Printf("[+] All %s dependencies installed successfully.\n", templateName)
	}

	// Import main feature to the application
	err = ImportFeatureToApp(templatePath, applicationTempDir)
	if err != nil {
		return fmt.Errorf(`error on %s installation: %s"`, templateName, err)
	}

	// Add main feature to installed features
	installedFeatures = append(installedFeatures, templateName)

	if global.Verbose {
		fmt.Println("[+] Setting up mks module manager..")
	}

	// Generate/Regenerate mks_modules files
	err = generateModuleManagerFiles(applicationTempDir, installedFeatures)
	if err != nil {
		return err
	}

	// Install all go packages to go.mod
	err = utils.InstallNeededPackages(applicationTempDir)
	if err != nil {
		return err
	}

	return nil
}

func ImportFeatureToApp(templatePath, workingDirectory string) error {
	templateName := filepath.Base(templatePath)

	if global.Verbose {
		fmt.Printf("[+] Preparing %s template to be imported to the application..\n", templateName)
	}

	/******************
	* PROMPTS PARSING *
	*******************/

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

	/**************************
	* LOAD AND UNLOAD PARSING *
	***************************/

	// Check if template has load file
	if global.Verbose {
		fmt.Printf("[+] Check if %s has a load file..\n", templateName)
	}

	// Validate that this.load file (if exists) pass code validations
	mainLoadFilePath := path.Join(templatePath, config.FILE_ADDON_TEMPLATE_MAIN_LOAD)

	if utils.FileOrDirectoryExists(mainLoadFilePath) {
		err := validateMksModulesFiles(mainLoadFilePath, config.SPELL_FUNCION_LOAD_PREFIX, templateName)
		if err != nil {
			return err
		}

		// Install load file
		if global.Verbose {
			fmt.Printf("[+] Installing %s load file..\n", templateName)
		}

		mainLoadFinalPath := path.Join(workingDirectory, config.FOLDER_MKS_MODULES, fmt.Sprintf("%s%s%s", config.SPELL_FUNCION_LOAD_PREFIX, templateName, config.FILE_EXTENSION_GO))
		err = utils.CreateFileFromTemplateWithCustomReplace(mainLoadFilePath, mainLoadFinalPath, placeHoldersToReplace)
		if err != nil {
			return err
		}
	}

	// Check if template has load file
	if global.Verbose {
		fmt.Printf("[+] Check if %s has an unload file..\n", templateName)
	}

	// Validate that this.unload file (if exists) pass code validations
	mainUnloadFilePath := path.Join(templateName, config.FILE_ADDON_TEMPLATE_MAIN_UNLOAD)

	if utils.FileOrDirectoryExists(mainLoadFilePath) {
		err := validateMksModulesFiles(mainUnloadFilePath, config.SPELL_FUNCION_UNLOAD_PREFIX, templateName)
		if err != nil {
			return err
		}

		// Install unload file
		if global.Verbose {
			fmt.Printf("[+] Installing %s unload file..\n", templateName)
		}

		mainUnloadFinalPath := path.Join(workingDirectory, config.FOLDER_MKS_MODULES, fmt.Sprintf("%s%s%s", config.SPELL_FUNCION_UNLOAD_PREFIX, templateName, config.FILE_EXTENSION_GO))
		err = utils.CreateFileFromTemplateWithCustomReplace(mainLoadFilePath, mainUnloadFinalPath, placeHoldersToReplace)
		if err != nil {
			return err
		}
	}

	/***********************
	* PARSING CONFIG FILES *
	************************/

	// template config files path
	envConfigFilePath := path.Join(templatePath, config.FILE_ADDON_TEMPLATE_ENVCONFIG)
	goConfigFilePath := path.Join(templatePath, config.FILE_ADDON_TEMPLATE_GOCONFIG)

	isGoConfigFile := utils.FileOrDirectoryExists(goConfigFilePath)
	isEnvConfigFile := utils.FileOrDirectoryExists(envConfigFilePath)

	// Check if template has load file
	if global.Verbose {
		fmt.Printf("[+] Check if %s has config files to install..\n", templateName)
	}

	// Check if is missing one config file
	if (isGoConfigFile && !isEnvConfigFile) || (isEnvConfigFile && !isGoConfigFile) {
		return fmt.Errorf("%s has a config file missing, if template use config must have %s and %s files", templateName, config.FILE_ADDON_TEMPLATE_GOCONFIG, config.FILE_ADDON_TEMPLATE_ENVCONFIG)
	}

	if isEnvConfigFile {
		if global.Verbose {
			fmt.Printf("[+] Installing %s env config file..\n", templateName)
		}

		// Get file content with placeholders replaced
		newEnvConfig, err := utils.ReadFileWithCustomReplace(envConfigFilePath, placeHoldersToReplace)
		if err != nil {
			return err
		}

		// Adding new config at bottom of app.env file
		err = utils.AddEnvConfigFromString(newEnvConfig, workingDirectory)
		if err != nil {
			return err
		}
	}

	if isGoConfigFile {
		if global.Verbose {
			fmt.Printf("[+] Installing %s go config file..\n", templateName)
		}

		// Get file content with placeholders replaced
		newGoConfig, err := utils.ReadFileWithCustomReplace(goConfigFilePath, placeHoldersToReplace)
		if err != nil {
			return err
		}

		// Adding new config inside config struct
		err = utils.AddGoConfigFromString(newGoConfig, workingDirectory)
		if err != nil {
			return err
		}
	}

	/************************************
	* PARSING EXTENDS & TEMPLATES FILES *
	************************************/

	var ExtendsFiles []string
	var TemplatesFiles []string

	fileList, err := utils.ListFiles(templatePath)
	if err != nil {
		return err
	}

	if global.Verbose {
		fmt.Printf("[+] Check if %s has .extends or .template files to install..\n", templateName)
	}

	// Search templates and extends files
	for _, file := range fileList {
		if strings.HasSuffix(file, config.FILE_EXTENSION_EXTENDS) {
			ExtendsFiles = append(ExtendsFiles, file)
		}

		if strings.HasSuffix(file, config.FILE_EXTENSION_TEMPLATE) {
			TemplatesFiles = append(TemplatesFiles, file)
		}
	}

	if global.Verbose && len(ExtendsFiles) > 0 {
		fmt.Println("[+] Preparing .extends files to install..")
	}

	// Iterate and copy every extend file
	for _, extendFile := range ExtendsFiles {
		if global.Verbose {
			fmt.Printf("[+] Installing %s files..\n", extendFile)
		}

		// Get mks file custom path structure
		filePath := path.Join(templatePath, extendFile)
		filePathStructure, err := utils.ProcessMksCustomFilesPath(filePath)
		if err != nil {
			return err
		}

		// Check if file to extend exists
		finalDirectoriesPath := path.Join(workingDirectory, config.FOLDER_SRC, path.Join(filePathStructure.Folders[:]...))
		finalFileToExtend := path.Join(finalDirectoriesPath, filePathStructure.FileName, config.FILE_EXTENSION_GO)
		if !utils.FileOrDirectoryExists(finalFileToExtend) {
			return fmt.Errorf("%s's %s extend file is trying to extend an unexistent file", templateName, extendFile)
		}

		// Get file content
		extendsFileContent, err := utils.ReadFileWithCustomReplace(filePath, placeHoldersToReplace)
		if err != nil {
			return err
		}

		err = utils.ExtendFile(finalFileToExtend, extendsFileContent)
		if err != nil {
			return err
		}

		if global.Verbose {
			fmt.Printf("[+] %s file installed successfuly..\n", extendFile)
		}
	}

	if global.Verbose && len(TemplatesFiles) > 0 {
		fmt.Println("[+] Preparing .template files to install..")
	}

	// Iterate and copy every template file
	for _, templateFile := range TemplatesFiles {
		if global.Verbose {
			fmt.Printf("[+] Installing %s files..\n", templateFile)
		}

		// Get mks file custom path structure
		filePath := path.Join(templatePath, templateFile)
		filePathStructure, err := utils.ProcessMksCustomFilesPath(templateFile)
		if err != nil {
			return err
		}

		// Check if directory exists (if not create them)
		finalDirectoriesPath := path.Join(workingDirectory, config.FOLDER_SRC, path.Join(filePathStructure.Folders[:]...))
		if !utils.FileOrDirectoryExists(finalDirectoriesPath) {
			err := os.MkdirAll(finalDirectoriesPath, config.FOLDER_PERMISSION)
			if err != nil {
				return err
			}
		}

		finalTemplateFilePath := path.Join(finalDirectoriesPath, fmt.Sprintf("%s%s", filePathStructure.FileName, config.FILE_EXTENSION_GO))
		if utils.FileOrDirectoryExists(finalTemplateFilePath) {
			return fmt.Errorf("%s's %s extend file is trying to create a file that already exists", templateName, templateFile)
		}

		err = utils.CreateFileFromTemplate(filePath, finalTemplateFilePath)
		if err != nil {
			return err
		}

		if global.Verbose {
			fmt.Printf("[+] %s file installed successfuly..\n", templateFile)
		}
	}

	if global.Verbose {
		fmt.Printf("[+] %s installed successfully..\n", templateName)
	}

	return nil
}

func generateModuleManagerFiles(workingDirectory string, installedFeatures []string) error {
	//TODO: Implement generation of loadModules and unloadModules
	//TODO: Implement generation of installed_features.json

	return nil
}
