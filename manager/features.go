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
	Features []Feature `json:"features"`
}

type Feature struct {
	Feature   string `json:"feature"`
	HasLoad   bool   `json:"hasLoad"`
	HasUnload bool   `json:"hasUnload"`
}

var temporalDirectoryPath string

func getApplicationInstalledFeatures() ([]Feature, error) {
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

func isFeatureInstalled(installedFeatures []Feature, featureName string) bool {
	for _, installedFeature := range installedFeatures {
		if installedFeature.Feature == featureName {
			return true
		}
	}

	return false
}

func IsValidFeature(feature string) bool {
	return utils.FileOrDirectoryExists(path.Join(global.UserTemplatesFolderPath, feature))
}

func hasFeatureLoadFile(feature string) bool {
	return utils.FileOrDirectoryExists(path.Join(global.UserTemplatesFolderPath, feature, config.FILE_ADDON_TEMPLATE_MAIN_LOAD))
}

func hasFeatureUnloadFile(feature string) bool {
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

	err = installFeature(path.Join(global.UserTemplatesFolderPath, feature))
	if err != nil {
		if temporalDirectoryPath != "" {
			utils.DeleteFileOrDirectory(temporalDirectoryPath)
		}

		return err
	}

	return nil
}

func installFeature(templatePath string) error {
	templateName := filepath.Base(templatePath)

	if global.Verbose {
		fmt.Printf("[+] Checking if %s's template is already installed in the application..\n", templateName)
	}

	// Get all installed features inside the application
	installedFeatures, err := getApplicationInstalledFeatures()
	if err != nil {
		return err
	}

	// Check if requested feature is already installed
	if isFeatureInstalled(installedFeatures, templateName) {
		fmt.Printf("%s's template is already installed.\n", templateName)
		// return nil to prevent errors with addAllFeatures
		return nil
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
	temporalDirectoryPath, err = utils.MakeTemporalDirectory()
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
		if !isFeatureInstalled(installedFeatures, dependencyTemplateName) {
			err := importFeatureToApp(path.Join(global.UserTemplatesFolderPath, dependencyTemplateName), applicationTempDir)
			if err != nil {
				return fmt.Errorf(`error on %s's "%s" dependency installation: %s"`, templateName, dependencyTemplateName, err)
			}

			installedFeatures = append(installedFeatures, Feature{
				Feature:   dependencyTemplateName,
				HasLoad:   hasFeatureLoadFile(dependencyTemplateName),
				HasUnload: hasFeatureUnloadFile(dependencyTemplateName),
			})
		}
	}

	if global.Verbose && len(dependenciesInOrder) > 0 {
		fmt.Printf("[+] All %s dependencies installed successfully.\n", templateName)
	}

	// Import main feature to the application
	err = importFeatureToApp(templatePath, applicationTempDir)
	if err != nil {
		return fmt.Errorf(`error on %s installation: %s"`, templateName, err)
	}

	// Add main feature to installed features
	installedFeatures = append(installedFeatures, Feature{
		Feature:   templateName,
		HasLoad:   hasFeatureLoadFile(templateName),
		HasUnload: hasFeatureUnloadFile(templateName),
	})

	if global.Verbose {
		fmt.Println("[+] Setting up mks module manager..")
	}

	// Generate/Regenerate mks_modules files
	err = generateModuleManagerFiles(applicationTempDir, installedFeatures)
	if err != nil {
		return err
	}

	if global.Verbose {
		fmt.Println("[+] Installing missing golang packages..")
	}

	// Install all go packages to go.mod
	err = utils.InstallNeededPackages(applicationTempDir)
	if err != nil {
		return err
	}

	if global.Verbose {
		fmt.Println("[+] Looking all go files to check syntax errors..")
	}

	// Check all go files sintax
	err = utils.CheckAllGoFilesInDirectory(applicationTempDir)
	if err != nil {
		return err
	}

	// Rename user application directory
	if err := os.Rename(global.BasePath, fmt.Sprintf("%s_bkp", global.BasePath)); err != nil {
		return err
	}

	// Copy modified application to BasePath
	if global.Verbose {
		fmt.Println("[+] Copying modified application to old application path..")
	}

	err = utils.CopyFileOrDirectory(applicationTempDir, global.BasePath)
	if err != nil {
		return err
	}

	return nil
}

func importFeatureToApp(templatePath, workingDirectory string) error {
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

func generateModuleManagerFiles(workingDirectory string, installedFeatures []Feature) error {
	// Update/Create the installed_features.json
	if global.Verbose {
		fmt.Println("[+] Updating installed_features.json..")
	}

	// Set struct data
	newInstalledFeatures := InstalledFeaturesFileFormat{
		Features: installedFeatures,
	}

	// Convert the structure to JSON format
	jsonData, err := json.Marshal(newInstalledFeatures)
	if err != nil {
		return err
	}

	// Write the data to the file
	installedFeaturesFilePath := path.Join(workingDirectory, config.FOLDER_MKS_MODULES, config.FILE_MKS_INSTALLED_FEATURES)
	err = os.WriteFile(installedFeaturesFilePath, jsonData, config.FOLDER_PERMISSION)
	if err != nil {
		return err
	}

	//Update modules_manager.go
	if global.Verbose {
		fmt.Println("[+] Updating modules_manager.go..")
	}

	var loadFunctions []string
	var unloadFunctions []string

	for _, installedFeature := range installedFeatures {
		// Check if feature has main.load file
		if installedFeature.HasLoad {
			// Add current feature to loadFunctions string slice variable
			loadFunctions = append(loadFunctions, fmt.Sprintf("%s%s()", config.SPELL_FUNCION_LOAD_PREFIX, installedFeature.Feature))

			// Path to mks_modules/load<feature>.go file
			finaMainLoadFilePath := path.Join(workingDirectory, config.FOLDER_MKS_MODULES, fmt.Sprintf("%s%s%s", config.SPELL_FUNCION_LOAD_PREFIX, installedFeature.Feature, config.FILE_EXTENSION_GO))

			// Check if load file is already installed inside mks_modules folder
			if !utils.FileOrDirectoryExists(finaMainLoadFilePath) {
				// Path of file to copy
				mainLoadFilePath := path.Join(global.UserTemplatesFolderPath, installedFeature.Feature, config.FILE_ADDON_TEMPLATE_MAIN_LOAD)

				// Copy file to mks_modules folder with the correct name
				err = utils.CopyFileOrDirectory(mainLoadFilePath, finaMainLoadFilePath)
				if err != nil {
					return err
				}
			}
		}

		// Check if feature has main.unload file
		if installedFeature.HasUnload {
			// Add current feature to unloadFunctions string slice variable
			unloadFunctions = append(unloadFunctions, fmt.Sprintf("%s%s()", config.SPELL_FUNCION_UNLOAD_PREFIX, installedFeature.Feature))

			// Path to mks_modules/unload<feature>.go file
			finaMainUnloadFilePath := path.Join(workingDirectory, config.FOLDER_MKS_MODULES, fmt.Sprintf("%s%s%s", config.SPELL_FUNCION_UNLOAD_PREFIX, installedFeature.Feature, config.FILE_EXTENSION_GO))

			// Check if unload file is already installed inside mks_modules folder
			if !utils.FileOrDirectoryExists(finaMainUnloadFilePath) {
				// Path of file to copy
				mainUnloadFilePath := path.Join(global.UserTemplatesFolderPath, installedFeature.Feature, config.FILE_ADDON_TEMPLATE_MAIN_UNLOAD)

				// Copy file to mks_modules folder with the correct name
				err = utils.CopyFileOrDirectory(mainUnloadFilePath, finaMainUnloadFilePath)
				if err != nil {
					return err
				}
			}
		}
	}

	// Path to mks_modules/module_manager.go
	moduleManagerFilePath := path.Join(workingDirectory, config.FOLDER_MKS_MODULES, config.FILE_MKS_MODULE_MANAGER)

	// Final function contents as strings
	finalLoadFunctions := fmt.Sprintf("\t%s", strings.Join(loadFunctions, "\n\t"))
	finalUnloadFunctions := fmt.Sprintf("\t%s", strings.Join(unloadFunctions, "\n\t"))

	// Add content to loadModules function
	err = utils.AddContentInsideFunction(moduleManagerFilePath, config.SPELL_FUNCION_LOAD_MODULE_MANAGER, finalLoadFunctions)
	if err != nil {
		return err
	}

	// Add content to unloadModules function
	err = utils.AddContentInsideFunction(moduleManagerFilePath, config.SPELL_FUNCION_UNLOAD_MODULE_MANAGER, finalUnloadFunctions)
	if err != nil {
		return err
	}

	return nil
}
