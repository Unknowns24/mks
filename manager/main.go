package manager

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/unknowns24/mks/utils"
)

func GenerateMicroservice(serviceName string, features []string) error {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Create the base path for the microservice
	basePath := filepath.Join(currentDir, serviceName)
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return err
	}

	createBaseFiles(basePath, serviceName)

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

	fmt.Printf("Microservice '%s' with features %v generated successfully!\n", serviceName, features)
	return nil
}

func createBaseFiles(basePath, serviceName string) error {
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
	mainTemplatePath := filepath.Join(mksDir, "..", "libs", "templates", "base", "main.go.template")
	mainPath := filepath.Join(basePath, "main.go")

	utils.CreateFileFromTemplate(basePath, mainTemplatePath, serviceName, mainPath)

	/**************
	* UTILS FILES *
	***************/

	// Create the utils path for the required base files
	utilsPath := filepath.Join(basePath, "utils")
	if err := os.MkdirAll(utilsPath, os.ModePerm); err != nil {
		return err
	}

	// Create utils/config.go using template
	configTemplatePath := filepath.Join(mksDir, "..", "libs", "templates", "base", "utils", "config.go.template")
	configPath := filepath.Join(utilsPath, "config.go")

	utils.CreateFileFromTemplate(basePath, configTemplatePath, serviceName, configPath)

	// Create utils/request.go using template
	requestTemplatePath := filepath.Join(mksDir, "..", "libs", "templates", "base", "utils", "request.go.template")
	requestPath := filepath.Join(utilsPath, "request.go")

	utils.CreateFileFromTemplate(basePath, requestTemplatePath, serviceName, requestPath)

	/**************
	* ROUTES FILE *
	**************/

	// Create the routes path for the required base files
	routesPath := filepath.Join(basePath, "routes")
	if err := os.MkdirAll(routesPath, os.ModePerm); err != nil {
		return err
	}

	// Create routes/mainRoutes.go using template
	routesTemplatePath := filepath.Join(mksDir, "..", "libs", "templates", "base", "routes", "mainRoutes.go.template")
	routesFilePath := filepath.Join(routesPath, "mainRoutes.go")

	utils.CreateFileFromTemplate(basePath, routesTemplatePath, serviceName, routesFilePath)

	/***********
	* COMMANDS *
	************/

	// Execute "go mod init" command in the basePath directory
	initCmd := exec.Command("go", "mod", "init", serviceName)
	initCmd.Dir = basePath
	initOutput, err := initCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running 'go mod init': %s\nOutput: %s", err, initOutput)
	}

	fmt.Println("Base files created successfully.")
	return nil
}
