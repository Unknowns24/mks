package manager

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
	"github.com/unknowns24/mks/utils"
)

func GenerateApplication(ApplicationName string, features []string) error {
	fmt.Println("[+] Creating " + ApplicationName + " application..")

	// making ApplicationName global
	global.ApplicationName = ApplicationName

	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	if global.Verbose {
		fmt.Println("[+] Creating root folder..")
	}

	// Create the base path for the application
	global.BasePath = filepath.Join(currentDir, global.ApplicationName)

	// Check if the application already was created
	_, dirErr := os.Stat(global.BasePath)
	if !os.IsNotExist(dirErr) {
		return errors.New(global.ApplicationName + " application already created!")
	}

	// Create base path directory
	if err := os.MkdirAll(global.BasePath, os.ModePerm); err != nil {
		utils.DeleteFileOrDirectory(global.BasePath)
		return err
	}

	if global.Verbose {
		fmt.Println("[+] Creating base files..")
	}

	// Create all base files
	err = createBaseFiles()
	if err != nil {
		utils.DeleteFileOrDirectory(global.BasePath)
		return err
	}

	// Install all base packages
	err = installBasePackages()
	if err != nil {
		return err
	}

	if len(features) >= 1 {
		if global.Verbose {
			fmt.Println("[+] Checking requested features..")
		}

		// Check if "all" features are requested
		if utils.SliceContainsElement(features, config.ALL_FEATURES) {
			// Add all available features to the application
			err := AddAllFeatures()
			if err != nil {
				return err
			}
		} else {
			// Add specified features to the application
			for _, feature := range features {
				err := AddFeature(feature)
				if err != nil {
					return err
				}
			}
		}

		fmt.Printf("[+] application '%s' with features %v generated successfully!\n", global.ApplicationName, features)
		return nil
	}

	fmt.Printf("[+] Base application '%s' generated successfully!\n", global.ApplicationName)
	return nil
}

func createBaseFiles() error {
	// Create the src path for the application
	srcPath := filepath.Join(global.BasePath, config.FOLDER_SRC)

	// Create src path directory
	if err := os.MkdirAll(srcPath, os.ModePerm); err != nil {
		return err
	}

	/***************************************
	* IMPORT MKS BASE FILES TO APPLICATION *
	****************************************/

	mksBaseFolder := path.Join(global.MksTemplatesFolderPath, config.FOLDER_BASE)
	err := utils.ImportBaseContent(mksBaseFolder, global.BasePath)
	if err != nil {
		return err
	}

	/***********
	* ENV FILE *
	************/

	// Create app.env.example using template
	if global.Verbose {
		fmt.Println("[+] Creating app.env.example file..")
	}

	// Map with all placeholders and its values to replace on .env template
	envReplaces := map[string]string{
		config.PLACEHOLDER_APP_NAME: global.ApplicationName,
	}

	appEnvTemplatePath := filepath.Join(global.MksTemplatesFolderPath, config.FOLDER_OTHERS, config.FILE_ENVCONFIG_APP)
	appEnvExampleFinalPath := filepath.Join(global.BasePath, config.FILE_CONFIG_ENVEXAMPLE)

	err = utils.CreateFileFromTemplateWithCustomReplace(appEnvTemplatePath, appEnvExampleFinalPath, envReplaces)
	if err != nil {
		return err
	}

	// Create app.env using template
	if global.Verbose {
		fmt.Println("[+] Creating app.env file..")
	}

	appEnvFinalPath := filepath.Join(global.BasePath, config.FILE_CONFIG_ENV)
	err = utils.CreateFileFromTemplateWithCustomReplace(appEnvTemplatePath, appEnvFinalPath, envReplaces)
	if err != nil {
		return err
	}

	/*************
	* GIT IGNORE *
	**************/

	// Create .gitignore file
	if global.Verbose {
		fmt.Println("[+] Creating .gitignore file..")
	}

	gitIgnoreTemplatePath := filepath.Join(global.MksTemplatesFolderPath, config.FOLDER_OTHERS, config.FILE_GITIGNORE)
	gitIgnoreExampleFinalPath := filepath.Join(global.BasePath, config.FILE_GITIGNORE)

	err = utils.CreateFileFromTemplate(gitIgnoreTemplatePath, gitIgnoreExampleFinalPath)
	if err != nil {
		return err
	}

	/**************
	* DOCKER FILE *
	***************/

	// Create Dockerfile using template
	if global.Verbose {
		fmt.Println("[+] Creating Dockerfile file..")
	}

	dockerTemplatePath := filepath.Join(global.MksTemplatesFolderPath, config.FOLDER_OTHERS, config.FILE_DOCKER)
	dockerFilePath := filepath.Join(global.BasePath, config.FILE_DOCKER)

	err = utils.CreateFileFromTemplate(dockerTemplatePath, dockerFilePath)
	if err != nil {
		return err
	}

	/***********
	* COMMANDS *
	************/

	// Execute "go mod init" command in the basePath directory
	if global.Verbose {
		fmt.Println("[+] Running go mod init..")
	}

	// Initialice Go modules
	err = utils.InitGoModules(global.ApplicationName, global.BasePath)
	if err != nil {
		return err
	}

	if global.Verbose {
		fmt.Println("[+] Base files created successfully.")
	}

	return nil
}

func installBasePackages() error {
	if global.Verbose {
		fmt.Println("[+] Installing base packages..")
	}

	err := utils.InstallNeededPackages(global.BasePath)
	if err != nil {
		return err
	}

	fmt.Println("[+] Base packages installed successfully..")

	return nil
}
