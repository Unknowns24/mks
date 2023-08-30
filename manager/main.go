package manager

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

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

	createBaseFiles(basePath, serviceName, verbose)

	// If features..
	if len(features) >= 1 {
		if verbose {
			fmt.Println("[+] Checking requested features..")
		}

		// Check if "all" features are requested
		if utils.SliceContainsElement(features, "all") {
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

	/************
	* MAIN FILE *
	*************/

	// Create main.go using template
	if verbose {
		fmt.Println("[+] Creating main.go file..")
	}

	mainTemplatePath := filepath.Join(mksDir, "..", "libs", "templates", "base", "main.go.template")
	mainPath := filepath.Join(basePath, "main.go")

	utils.CreateFileFromTemplate(basePath, mainTemplatePath, serviceName, mainPath)

	/**************
	* UTILS FILES *
	***************/

	// Create the utils path for the required base files
	if verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating utils files directory..")
	}

	utilsPath := filepath.Join(basePath, "utils")
	if err := os.MkdirAll(utilsPath, os.ModePerm); err != nil {
		return err
	}

	// Create utils/config.go using template
	if verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating utils/config.go file..")
	}

	configTemplatePath := filepath.Join(mksDir, "..", "libs", "templates", "base", "utils", "config.go.template")
	configPath := filepath.Join(utilsPath, "config.go")

	utils.CreateFileFromTemplate(basePath, configTemplatePath, serviceName, configPath)

	// Create utils/request.go using template
	if verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating utils/request.go file..")
	}

	requestTemplatePath := filepath.Join(mksDir, "..", "libs", "templates", "base", "utils", "request.go.template")
	requestPath := filepath.Join(utilsPath, "request.go")

	utils.CreateFileFromTemplate(basePath, requestTemplatePath, serviceName, requestPath)

	/**************
	* ROUTES FILE *
	**************/

	// Create the routes path for the required base files
	if verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating routes files directory..")
	}

	routesPath := filepath.Join(basePath, "routes")
	if err := os.MkdirAll(routesPath, os.ModePerm); err != nil {
		return err
	}

	// Create routes/mainRoutes.go using template
	if verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Creating routes/mainRoutes.go file..")
	}

	routesTemplatePath := filepath.Join(mksDir, "..", "libs", "templates", "base", "routes", "mainRoutes.go.template")
	routesFilePath := filepath.Join(routesPath, "mainRoutes.go")

	utils.CreateFileFromTemplate(basePath, routesTemplatePath, serviceName, routesFilePath)

	/***********
	* COMMANDS *
	************/

	// Execute "go mod init" command in the basePath directory
	if verbose {
		time.Sleep(time.Second / 4) // sleep 250ms
		fmt.Println("[+] Running go mod init..")
	}

	initCmd := exec.Command("go", "mod", "init", serviceName)
	initCmd.Dir = basePath
	initOutput, err := initCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running 'go mod init': %s\nOutput: %s", err, initOutput)
	}

	if verbose {
		fmt.Println("[+] Base files created successfully.")
	}

	return nil
}
