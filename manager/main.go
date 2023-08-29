package manager

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

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

	mainTemplateContent, err := os.ReadFile(mainTemplatePath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Replace placeholders with the actual service name
	mainContent := strings.ReplaceAll(string(mainTemplateContent), "%PACKAGE_NAME%", serviceName)

	mainPath := filepath.Join(basePath, "main.go")
	if err := os.WriteFile(mainPath, []byte(mainContent), os.ModePerm); err != nil {
		return err
	}

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
	configTemplateContent, err := os.ReadFile(configTemplatePath)
	if err != nil {
		return err
	}

	// Replace placeholders with the actual service name
	configContent := strings.ReplaceAll(string(configTemplateContent), "%PACKAGE_NAME%", serviceName)

	configPath := filepath.Join(utilsPath, "config.go")
	if err := os.WriteFile(configPath, []byte(configContent), os.ModePerm); err != nil {
		return err
	}

	// Create utils/request.go using template
	requestTemplatePath := filepath.Join(mksDir, "..", "libs", "templates", "base", "utils", "request.go.template")
	requestTemplateContent, err := os.ReadFile(requestTemplatePath)
	if err != nil {
		return err
	}

	// Replace placeholders with the actual service name
	requestContent := strings.ReplaceAll(string(requestTemplateContent), "%PACKAGE_NAME%", serviceName)

	requestPath := filepath.Join(utilsPath, "request.go")
	if err := os.WriteFile(requestPath, []byte(requestContent), os.ModePerm); err != nil {
		return err
	}

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
	routesTemplateContent, err := os.ReadFile(routesTemplatePath)
	if err != nil {
		return err
	}

	// Replace placeholders with the actual service name
	routesContent := strings.ReplaceAll(string(routesTemplateContent), "%PACKAGE_NAME%", serviceName)

	routesFilePath := filepath.Join(routesPath, "mainRoutes.go")
	if err := os.WriteFile(routesFilePath, []byte(routesContent), os.ModePerm); err != nil {
		return err
	}

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
