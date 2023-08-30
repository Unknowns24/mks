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

func GenerateMicroservice(serviceName string, verbose bool, features []string) error {
	if verbose {
		fmt.Println("[+] Creating " + serviceName + " microservice..")
	}

	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	if verbose {
		fmt.Println("[+] Creating root folder..")
	}

	// Create the base path for the microservice
	basePath := filepath.Join(currentDir, serviceName)

	// Check if the microservice already was created
	_, dirErr := os.Stat(basePath)
	if !os.IsNotExist(dirErr) {
		return errors.New(serviceName + " microservice already created!")
	}

	// Create base path directory
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return err
	}

	if verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating base files..")
	}

	// Create all base files
	err = createBaseFiles(basePath, serviceName, verbose)
	if err != nil {
		return err
	}

	// Install all base packages
	err = installBasePackages(basePath, verbose)
	if err != nil {
		return err
	}

	if len(features) >= 1 {
		if verbose {
			fmt.Println("[+] Checking requested features..")
		}

		// Check if "all" features are requested
		if utils.SliceContainsElement(features, config.ALL_FEATURES) {
			// Add all available features to the microservice
			err := AddAllFeatures(basePath)
			if err != nil {
				return err
			}
		} else {
			// Add specified features to the microservice
			for _, feature := range features {
				err := AddFeature(basePath, feature)
				if err != nil {
					return err
				}
			}
		}

		fmt.Printf("[+] Microservice '%s' with features %v generated successfully!\n", serviceName, features)
		return nil
	}

	fmt.Printf("[+] Base Microservice '%s' generated successfully!\n", serviceName)
	return nil
}

func createBaseFiles(basePath, serviceName string, verbose bool) error {
	// Get the directory path of the current file (generator.go)
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to get current file path")
	}

	mksDir := filepath.Dir(filename)

	// Create the src path for the microservice
	srcPath := filepath.Join(basePath, "src")

	// Create src path directory
	if err := os.MkdirAll(srcPath, os.ModePerm); err != nil {
		return err
	}

	/************
	* MAIN FILE *
	*************/

	// Create main.go using template
	if verbose {
		fmt.Println("[+] Creating main.go file..")
	}

	mainTemplatePath := filepath.Join(mksDir, "..", "libs", "templates", "base", "main.template")
	mainFinalPath := filepath.Join(srcPath, "main.go")

	err := utils.CreateFileFromTemplate(mainTemplatePath, serviceName, mainFinalPath)
	if err != nil {
		return err
	}

	/**************
	* UTILS FILES *
	***************/

	// Create the utils path for the required base files
	if verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating utils files directory..")
	}

	utilsPath := filepath.Join(srcPath, "utils")
	if err := os.MkdirAll(utilsPath, os.ModePerm); err != nil {
		return err
	}

	// Create utils/config.go using template
	if verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating utils/config.go file..")
	}

	configTemplatePath := filepath.Join(mksDir, "..", "libs", "templates", "base", "utils", "config.template")
	configFinalPath := filepath.Join(utilsPath, "config.go")

	err = utils.CreateFileFromTemplate(configTemplatePath, serviceName, configFinalPath)
	if err != nil {
		return err
	}

	// Create utils/request.go using template
	if verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating utils/request.go file..")
	}

	requestTemplatePath := filepath.Join(mksDir, "..", "libs", "templates", "base", "utils", "request.template")
	requestFinalPath := filepath.Join(utilsPath, "request.go")

	err = utils.CreateFileFromTemplate(requestTemplatePath, serviceName, requestFinalPath)
	if err != nil {
		return err
	}

	/**************
	* ROUTES FILE *
	**************/

	// Create the routes path for the required base files
	if verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating routes files directory..")
	}

	routesPath := filepath.Join(srcPath, "routes")
	if err := os.MkdirAll(routesPath, os.ModePerm); err != nil {
		return err
	}

	// Create routes/mainRoutes.go using template
	if verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating routes/mainRoutes.go file..")
	}

	routesTemplatePath := filepath.Join(mksDir, "..", "libs", "templates", "base", "routes", "mainRoutes.template")
	routesFinalPath := filepath.Join(routesPath, "mainRoutes.go")

	err = utils.CreateFileFromTemplate(routesTemplatePath, serviceName, routesFinalPath)
	if err != nil {
		return err
	}

	/******************
	* TEST CONTROLLER *
	*******************/

	// Create the routes path for the required base files
	if verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating controllers files directory..")
	}

	controllersPath := filepath.Join(srcPath, "controllers")
	if err := os.MkdirAll(controllersPath, os.ModePerm); err != nil {
		return err
	}

	// Create controllers/testController.go using template
	if verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating controllers/testController.go file..")
	}

	controllerTemplatePath := filepath.Join(mksDir, "..", "libs", "templates", "base", "controllers", "testController.template")
	controllerFinalPath := filepath.Join(controllersPath, "testController.go")

	err = utils.CreateFileFromTemplate(controllerTemplatePath, serviceName, controllerFinalPath)
	if err != nil {
		return err
	}

	/**************
	* DOCKER FILE *
	***************/

	// Create Dockerfile using template
	if verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating Dockerfile file..")
	}

	dockerTemplatePath := filepath.Join(mksDir, "..", "libs", "templates", "docker", "Dockerfile")
	dockerFilePath := filepath.Join(basePath, "Dockerfile")

	err = utils.CreateFileFromTemplate(dockerTemplatePath, serviceName, dockerFilePath)
	if err != nil {
		return err
	}

	/***********
	* COMMANDS *
	************/

	// Execute "go mod init" command in the basePath directory
	if verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Running go mod init..")
	}

	// Initialice Go modules
	err = utils.InitGoModules(serviceName, basePath)
	if err != nil {
		return err
	}

	if verbose {
		fmt.Println("[+] Base files created successfully.")
	}

	return nil
}

func installBasePackages(basePath string, verbose bool) error {
	if verbose {
		fmt.Println("[+] Installing base packages..")
	}

	err := utils.InstallNeededPackages(basePath)
	if err != nil {
		return err
	}

	fmt.Println("[+] Base packages installed successfully..")

	return nil
}
