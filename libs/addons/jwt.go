package addons

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
	"github.com/unknowns24/mks/utils"
)

func InstallJWT() error {
	// Createng a variable to store the path to jwt templates folder
	jwtTemplatesFolderPath := path.Join(global.TemplatesFolderPath, config.FOLDER_ADDONS, config.FOLDER_JWT)

	if global.Verbose {
		time.Sleep(time.Second / 5) // sleep 200ms
		fmt.Println("[+] Installing JWT addon..")
		fmt.Println("[+] Creating middleware folder if not exists..")
	}

	// Create the middleware folder path
	middlewareFolderPath := filepath.Join(global.BasePath, config.FOLDER_SRC, config.FOLDER_MIDDLEWARES)

	// Create Middlewares folder if not exists
	_, dirErr := os.Stat(middlewareFolderPath)
	if os.IsNotExist(dirErr) {
		if err := os.MkdirAll(middlewareFolderPath, os.ModePerm); err != nil {
			return err
		}
	}

	// Ask if user wants to include generate jwt function
	includeGenerateJWTFunction, err := utils.AskConfirm("Include GenerateJWT function")
	if err != nil {
		return err
	}

	// Create middlewares/auth.go using template
	if global.Verbose {
		time.Sleep(time.Second / 5) // sleep 200ms
		fmt.Println("[+] Creating middlewares/auth.go file..")
	}

	jwtTemplatePath := filepath.Join(jwtTemplatesFolderPath, config.FILE_TEMPLATE_JWT)
	jwtFinalPath := filepath.Join(middlewareFolderPath, config.FILE_GO_JWT)

	err = utils.CreateFileFromTemplate(jwtTemplatePath, global.ServiceName, jwtFinalPath)
	if err != nil {
		return err
	}

	if includeGenerateJWTFunction {
		// Jwt auth.extends file path
		jwtAuthExtendsFilePath := path.Join(jwtTemplatesFolderPath, config.FILE_EXTENDS_AUTH_JWT)

		// Get file content
		extendsFileContent, err := utils.ReadFile(jwtAuthExtendsFilePath)
		if err != nil {
			return err
		}

		utils.ExtendFile(jwtFinalPath, extendsFileContent)
	}

	// Ask JWT Token to the user
	jwtToken, err := utils.AskData("Set JWT Token (this key needs to be the same on all microservices to right authentication)")
	if err != nil {
		return err
	}

	// Map with all placeholders and its values to replace on .env template
	envReplaces := map[string]string{
		config.PLACEHOLDER_JWT_TOKEN: jwtToken,
	}

	/****************************
	* PARSE CONFIGURATION FILES *
	*****************************/

	if global.Verbose {
		fmt.Println("[+] Adding JWT config to env file..")
		time.Sleep(time.Second / 5) // sleep 200ms
	}

	// Jwt .envconfig file path
	jwtEnvConfigPath := path.Join(jwtTemplatesFolderPath, config.FILE_ENVCONFIG_JWT)

	// Get file content with placeholders replaced
	newEnvConfig, err := utils.ReadFileWithCustomReplace(jwtEnvConfigPath, envReplaces)
	if err != nil {
		return err
	}

	// Adding new config at bottom of app.env file
	err = utils.AddEnvConfigFromString(newEnvConfig)
	if err != nil {
		return err
	}

	if global.Verbose {
		fmt.Println("[+] Adding JWT config to config.go file..")
		time.Sleep(time.Second / 5) // sleep 200ms
	}

	// Jwt .goconfig file path
	jwtGoConfigPath := path.Join(jwtTemplatesFolderPath, config.FILE_GOCONFIG_JWT)

	// Get file content
	newGoConfig, err := utils.ReadFile(jwtGoConfigPath)
	if err != nil {
		return err
	}

	// Adding new config inside config struct
	err = utils.AddGoConfigFromString(newGoConfig)
	if err != nil {
		return err
	}

	/**************************
	* INSTALLING ALL PACKAGES *
	***************************/

	if global.Verbose {
		fmt.Println("[+] Installing JWT packages..")
	}

	err = utils.InstallNeededPackages(global.BasePath)
	if err != nil {
		return err
	}

	if global.Verbose {
		time.Sleep(time.Second / 5) // sleep 200ms
		fmt.Println("[+] JWT installed successfully..")
	}

	return nil
}
