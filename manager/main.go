package manager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
	"github.com/unknowns24/mks/utils"
	"github.com/unknowns24/mks/validators"
)

func GenerateMicroservice(serviceName string, features []string) error {
	fmt.Println("[+] Creating " + serviceName + " microservice..")

	// making serviceName global
	global.ServiceName = serviceName

	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	if global.Verbose {
		fmt.Println("[+] Creating root folder..")
	}

	// Create the base path for the microservice
	global.BasePath = filepath.Join(currentDir, global.ServiceName)

	// Check if the microservice already was created
	_, dirErr := os.Stat(global.BasePath)
	if !os.IsNotExist(dirErr) {
		return errors.New(global.ServiceName + " microservice already created!")
	}

	// Create base path directory
	if err := os.MkdirAll(global.BasePath, os.ModePerm); err != nil {
		return err
	}

	if global.Verbose {
		time.Sleep(time.Second / 5) // sleep 200ms
		fmt.Println("[+] Creating base files..")
	}

	// Create all base files
	err = createBaseFiles()
	if err != nil {
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
			// Add all available features to the microservice
			err := AddAllFeatures()
			if err != nil {
				return err
			}
		} else {
			// Add specified features to the microservice
			for _, feature := range features {
				err := AddFeature(feature)
				if err != nil {
					return err
				}
			}
		}

		fmt.Printf("[+] Microservice '%s' with features %v generated successfully!\n", global.ServiceName, features)
		return nil
	}

	fmt.Printf("[+] Base Microservice '%s' generated successfully!\n", global.ServiceName)
	return nil
}

func createBaseFiles() error {
	// Create the src path for the microservice
	srcPath := filepath.Join(global.BasePath, config.FOLDER_SRC)

	// Create src path directory
	if err := os.MkdirAll(srcPath, os.ModePerm); err != nil {
		return err
	}

	/************
	* MAIN FILE *
	*************/

	// Create main.go using template
	if global.Verbose {
		fmt.Println("[+] Creating main.go file..")
	}

	mainTemplatePath := filepath.Join(global.TemplatesFolderPath, config.FOLDER_BASE, config.FILE_TEMPLATE_MAIN)
	mainFinalPath := filepath.Join(srcPath, config.FILE_GO_MAIN)

	err := utils.CreateFileFromTemplate(mainTemplatePath, global.ServiceName, mainFinalPath)
	if err != nil {
		return err
	}

	/**************
	* UTILS FILES *
	***************/

	// Create the utils path for the required base files
	if global.Verbose {
		time.Sleep(time.Second / 5) // sleep 200ms
		fmt.Println("[+] Creating utils files directory..")
	}

	utilsPath := filepath.Join(srcPath, config.FOLDER_UTILS)
	if err := os.MkdirAll(utilsPath, os.ModePerm); err != nil {
		return err
	}

	// Create utils/config.go using template
	if global.Verbose {
		time.Sleep(time.Second / 5) // sleep 200ms
		fmt.Println("[+] Creating utils/config.go file..")
	}

	configTemplatePath := filepath.Join(global.TemplatesFolderPath, config.FOLDER_BASE, config.FOLDER_UTILS, config.FILE_TEMPLATE_CONFIG)
	configFinalPath := filepath.Join(utilsPath, config.FILE_GO_CONFIG)

	err = utils.CreateFileFromTemplate(configTemplatePath, global.ServiceName, configFinalPath)
	if err != nil {
		return err
	}

	// Create utils/request.go using template
	if global.Verbose {
		time.Sleep(time.Second / 5) // sleep 200ms
		fmt.Println("[+] Creating utils/request.go file..")
	}

	requestTemplatePath := filepath.Join(global.TemplatesFolderPath, config.FOLDER_BASE, config.FOLDER_UTILS, config.FILE_TEMPLATE_REQUEST)
	requestFinalPath := filepath.Join(utilsPath, config.FILE_GO_REQUEST)

	err = utils.CreateFileFromTemplate(requestTemplatePath, global.ServiceName, requestFinalPath)
	if err != nil {
		return err
	}

	/**************
	* ROUTES FILE *
	**************/

	// Create the routes path for the required base files
	if global.Verbose {
		time.Sleep(time.Second / 5) // sleep 200ms
		fmt.Println("[+] Creating routes files directory..")
	}

	routesPath := filepath.Join(srcPath, config.FOLDER_ROUTES)
	if err := os.MkdirAll(routesPath, os.ModePerm); err != nil {
		return err
	}

	// Create routes/mainRoutes.go using template
	if global.Verbose {
		time.Sleep(time.Second / 5) // sleep 200ms
		fmt.Println("[+] Creating routes/mainRoutes.go file..")
	}

	routesTemplatePath := filepath.Join(global.TemplatesFolderPath, config.FOLDER_BASE, config.FOLDER_ROUTES, config.FILE_TEMPLATE_MAINROUTES)
	routesFinalPath := filepath.Join(routesPath, config.FILE_GO_MAINROUTES)

	err = utils.CreateFileFromTemplate(routesTemplatePath, global.ServiceName, routesFinalPath)
	if err != nil {
		return err
	}

	/******************
	* TEST CONTROLLER *
	*******************/

	// Create the routes path for the required base files
	if global.Verbose {
		time.Sleep(time.Second / 5) // sleep 200ms
		fmt.Println("[+] Creating controllers files directory..")
	}

	controllersPath := filepath.Join(srcPath, config.FOLDER_CONTROLLERS)
	if err := os.MkdirAll(controllersPath, os.ModePerm); err != nil {
		return err
	}

	// Create controllers/testController.go using template
	if global.Verbose {
		time.Sleep(time.Second / 5) // sleep 200ms
		fmt.Println("[+] Creating controllers/testController.go file..")
	}

	controllerTemplatePath := filepath.Join(global.TemplatesFolderPath, config.FOLDER_BASE, config.FOLDER_CONTROLLERS, config.FILE_TEMPLATE_TESTCONTROLLER)
	controllerFinalPath := filepath.Join(controllersPath, config.FILE_GO_TESTCONTROLLER)

	err = utils.CreateFileFromTemplate(controllerTemplatePath, global.ServiceName, controllerFinalPath)
	if err != nil {
		return err
	}

	/***********
	* ENV FILE *
	************/

	// Create app.env.example using template
	if global.Verbose {
		time.Sleep(time.Second / 5) // sleep 200ms
		fmt.Println("[+] Creating app.env.example file..")
	}

	// Ask the user for the desired microservice port
	appPort, err := utils.AskDataWithValidation("Set microservice port (1024 to 65535):", validators.ValidatePortRange)
	if err != nil {
		return err
	}

	// Map with all placeholders and its values to replace on .env template
	envReplaces := map[string]string{
		config.PLACEHOLDER_APP_NAME: global.ServiceName,
		config.PLACEHOLDER_APP_PORT: appPort,
	}

	appEnvTemplatePath := filepath.Join(global.TemplatesFolderPath, config.FOLDER_OTHERS, config.FILE_ENVCONFIG_APP)
	appEnvExampleFinalPath := filepath.Join(global.BasePath, config.FILE_CONFIG_ENVEXAMPLE)

	err = utils.CreateFileFromTemplateWithCustomReplace(appEnvTemplatePath, appEnvExampleFinalPath, envReplaces)
	if err != nil {
		return err
	}

	// Create app.env using template
	if global.Verbose {
		time.Sleep(time.Second / 5) // sleep 200ms
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
		time.Sleep(time.Second / 5) // sleep 200ms
		fmt.Println("[+] Creating .gitignore file..")
	}

	gitIgnoreTemplatePath := filepath.Join(global.TemplatesFolderPath, config.FOLDER_OTHERS, config.FILE_GITIGNORE)
	gitIgnoreExampleFinalPath := filepath.Join(global.BasePath, config.FILE_GITIGNORE)

	err = utils.CreateFileFromTemplate(gitIgnoreTemplatePath, global.ServiceName, gitIgnoreExampleFinalPath)
	if err != nil {
		return err
	}

	/**************
	* DOCKER FILE *
	***************/

	// Create Dockerfile using template
	if global.Verbose {
		time.Sleep(time.Second / 5) // sleep 200ms
		fmt.Println("[+] Creating Dockerfile file..")
	}

	dockerTemplatePath := filepath.Join(global.TemplatesFolderPath, config.FOLDER_OTHERS, config.FILE_DOCKER)
	dockerFilePath := filepath.Join(global.BasePath, config.FILE_DOCKER)

	err = utils.CreateFileFromTemplate(dockerTemplatePath, global.ServiceName, dockerFilePath)
	if err != nil {
		return err
	}

	/***********
	* COMMANDS *
	************/

	// Execute "go mod init" command in the basePath directory
	if global.Verbose {
		time.Sleep(time.Second / 5) // sleep 200ms
		fmt.Println("[+] Running go mod init..")
	}

	// Initialice Go modules
	err = InitGoModules(global.ServiceName, global.BasePath)
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

	err := InstallNeededPackages(global.BasePath)
	if err != nil {
		return err
	}

	fmt.Println("[+] Base packages installed successfully..")

	return nil
}
