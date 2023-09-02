package manager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/utils"
)

func GenerateMicroservice(serviceName string, features []string) error {
	// making serviceName global
	config.ServiceName = serviceName

	if config.Verbose {
		fmt.Println("[+] Creating " + config.ServiceName + " microservice..")
	}

	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	if config.Verbose {
		fmt.Println("[+] Creating root folder..")
	}

	// Create the base path for the microservice
	config.BasePath = filepath.Join(currentDir, config.ServiceName)

	// Check if the microservice already was created
	_, dirErr := os.Stat(config.BasePath)
	if !os.IsNotExist(dirErr) {
		return errors.New(config.ServiceName + " microservice already created!")
	}

	// Create base path directory
	if err := os.MkdirAll(config.BasePath, os.ModePerm); err != nil {
		return err
	}

	if config.Verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
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
		if config.Verbose {
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

		fmt.Printf("[+] Microservice '%s' with features %v generated successfully!\n", config.ServiceName, features)
		return nil
	}

	fmt.Printf("[+] Base Microservice '%s' generated successfully!\n", config.ServiceName)
	return nil
}

func createBaseFiles() error {
	// Get the directory path of the current file (generator.go)
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to get current file path")
	}

	mksDir := filepath.Dir(filename)

	// Create the src path for the microservice
	srcPath := filepath.Join(config.BasePath, config.FOLDER_SRC)

	// Create src path directory
	if err := os.MkdirAll(srcPath, os.ModePerm); err != nil {
		return err
	}

	/************
	* MAIN FILE *
	*************/

	// Create main.go using template
	if config.Verbose {
		fmt.Println("[+] Creating main.go file..")
	}

	mainTemplatePath := filepath.Join(mksDir, "..", config.FOLDER_LIBS, config.FOLDER_TEMPLATES, config.FOLDER_BASE, config.FILE_TEMPLATE_MAIN)
	mainFinalPath := filepath.Join(srcPath, config.FILE_GO_MAIN)

	err := utils.CreateFileFromTemplate(mainTemplatePath, config.ServiceName, mainFinalPath)
	if err != nil {
		return err
	}

	/**************
	* UTILS FILES *
	***************/

	// Create the utils path for the required base files
	if config.Verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating utils files directory..")
	}

	utilsPath := filepath.Join(srcPath, config.FOLDER_UTILS)
	if err := os.MkdirAll(utilsPath, os.ModePerm); err != nil {
		return err
	}

	// Create utils/config.go using template
	if config.Verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating utils/config.go file..")
	}

	configTemplatePath := filepath.Join(mksDir, "..", config.FOLDER_LIBS, config.FOLDER_TEMPLATES, config.FOLDER_BASE, config.FOLDER_UTILS, config.FILE_TEMPLATE_CONFIG)
	configFinalPath := filepath.Join(utilsPath, config.FILE_GO_CONFIG)

	err = utils.CreateFileFromTemplate(configTemplatePath, config.ServiceName, configFinalPath)
	if err != nil {
		return err
	}

	// Create utils/request.go using template
	if config.Verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating utils/request.go file..")
	}

	requestTemplatePath := filepath.Join(mksDir, "..", config.FOLDER_LIBS, config.FOLDER_TEMPLATES, config.FOLDER_BASE, config.FOLDER_UTILS, config.FILE_TEMPLATE_REQUEST)
	requestFinalPath := filepath.Join(utilsPath, config.FILE_GO_REQUEST)

	err = utils.CreateFileFromTemplate(requestTemplatePath, config.ServiceName, requestFinalPath)
	if err != nil {
		return err
	}

	/**************
	* ROUTES FILE *
	**************/

	// Create the routes path for the required base files
	if config.Verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating routes files directory..")
	}

	routesPath := filepath.Join(srcPath, config.FOLDER_ROUTES)
	if err := os.MkdirAll(routesPath, os.ModePerm); err != nil {
		return err
	}

	// Create routes/mainRoutes.go using template
	if config.Verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating routes/mainRoutes.go file..")
	}

	routesTemplatePath := filepath.Join(mksDir, "..", config.FOLDER_LIBS, config.FOLDER_TEMPLATES, config.FOLDER_BASE, config.FOLDER_ROUTES, config.FILE_TEMPLATE_MAINROUTES)
	routesFinalPath := filepath.Join(routesPath, config.FILE_GO_MAINROUTES)

	err = utils.CreateFileFromTemplate(routesTemplatePath, config.ServiceName, routesFinalPath)
	if err != nil {
		return err
	}

	/******************
	* TEST CONTROLLER *
	*******************/

	// Create the routes path for the required base files
	if config.Verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating controllers files directory..")
	}

	controllersPath := filepath.Join(srcPath, config.FOLDER_CONTROLLERS)
	if err := os.MkdirAll(controllersPath, os.ModePerm); err != nil {
		return err
	}

	// Create controllers/testController.go using template
	if config.Verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating controllers/testController.go file..")
	}

	controllerTemplatePath := filepath.Join(mksDir, "..", config.FOLDER_LIBS, config.FOLDER_TEMPLATES, config.FOLDER_BASE, config.FOLDER_CONTROLLERS, config.FILE_TEMPLATE_TESTCONTROLLER)
	controllerFinalPath := filepath.Join(controllersPath, config.FILE_GO_TESTCONTROLLER)

	err = utils.CreateFileFromTemplate(controllerTemplatePath, config.ServiceName, controllerFinalPath)
	if err != nil {
		return err
	}

	/***********
	* ENV FILE *
	************/

	// Create app.env.example using template
	if config.Verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating app.env.example file..")
	}

	appEnvTemplatePath := filepath.Join(mksDir, "..", config.FOLDER_LIBS, config.FOLDER_TEMPLATES, config.FOLDER_BASE, config.FOLDER_OTHERS, config.FILE_ENVEXAMPLE_CONFIG)
	appEnvExampleFinalPath := filepath.Join(config.BasePath, config.FILE_ENVEXAMPLE_CONFIG)

	err = utils.CreateFileFromTemplate(appEnvTemplatePath, config.ServiceName, appEnvExampleFinalPath)
	if err != nil {
		return err
	}

	// Create app.env using template
	if config.Verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating app.env file..")
	}

	appEnvFinalPath := filepath.Join(config.BasePath, config.FILE_ENV_CONFIG)
	err = utils.CreateFileFromTemplate(appEnvTemplatePath, config.ServiceName, appEnvFinalPath)
	if err != nil {
		return err
	}

	/*************
	* GIT IGNORE *
	**************/

	// Create .gitignore file
	if config.Verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating .gitignore file..")
	}

	gitIgnoreTemplatePath := filepath.Join(mksDir, "..", config.FOLDER_LIBS, config.FOLDER_TEMPLATES, config.FOLDER_BASE, config.FOLDER_OTHERS, config.FILE_GITIGNORE)
	gitIgnoreExampleFinalPath := filepath.Join(config.BasePath, config.FILE_GITIGNORE)

	err = utils.CreateFileFromTemplate(gitIgnoreTemplatePath, config.ServiceName, gitIgnoreExampleFinalPath)
	if err != nil {
		return err
	}

	/**************
	* DOCKER FILE *
	***************/

	// Create Dockerfile using template
	if config.Verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating Dockerfile file..")
	}

	dockerTemplatePath := filepath.Join(mksDir, "..", config.FOLDER_LIBS, config.FOLDER_TEMPLATES, config.FOLDER_DOCKER, config.FILE_DOCKER)
	dockerFilePath := filepath.Join(config.BasePath, config.FILE_DOCKER)

	err = utils.CreateFileFromTemplate(dockerTemplatePath, config.ServiceName, dockerFilePath)
	if err != nil {
		return err
	}

	/***********
	* COMMANDS *
	************/

	// Execute "go mod init" command in the basePath directory
	if config.Verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Running go mod init..")
	}

	// Initialice Go modules
	err = utils.InitGoModules(config.ServiceName, config.BasePath)
	if err != nil {
		return err
	}

	if config.Verbose {
		fmt.Println("[+] Base files created successfully.")
	}

	return nil
}

func installBasePackages() error {
	if config.Verbose {
		fmt.Println("[+] Installing base packages..")
	}

	err := utils.InstallNeededPackages(config.BasePath)
	if err != nil {
		return err
	}

	fmt.Println("[+] Base packages installed successfully..")

	return nil
}
