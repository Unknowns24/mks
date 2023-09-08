package manager

import (
	"fmt"
	"path"
	"strings"

	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
	"github.com/unknowns24/mks/utils"
)

func ListTemplate() error {
	fmt.Println("[+] List of templates installed and availables to use:")

	for _, template := range global.InstalledTemplates {
		fmt.Println(template)
	}
	return nil
}

func UninstallTemplate(template string) error {
	fmt.Println("[+] Uninstalling " + template + " template...")

	if global.Verbose {
		fmt.Println("[+] Checking if template is installed...")
	}

	// check if template exists in addons folder
	if !utils.SliceContainsElement(global.InstalledTemplates, template) {
		return fmt.Errorf("template not installed: %s", template)
	}

	if utils.FileOrDirectoryExists(path.Join(global.TemplatesFolderPath, config.FOLDER_ADDONS, template)) {
		utils.DeleteFileOrDirectory(path.Join(global.TemplatesFolderPath, config.FOLDER_ADDONS, template)) // delete template directory
	}

	if utils.FileOrDirectoryExists(path.Join(global.TemplatesFolderPath, config.FOLDER_ADDONS, template)) {
		return fmt.Errorf("failed to uninstall template: %s (Try delete folder %s manually)", template, path.Join(global.TemplatesFolderPath, config.FOLDER_ADDONS, template))
	}

	if global.Verbose {
		fmt.Println("[+] Template uninstalled successfully!")
	}

	return nil
}

func InstallTemplate(template string) error {
	fmt.Println("[+] Installing " + template + " template...")

	if global.Verbose {
		fmt.Println("[+] Checking if template is valid...")
	}

	addonsPath := path.Join(global.TemplatesFolderPath, config.FOLDER_ADDONS)

	tempDirPath, err := utils.MakeTemporalDirectory()
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %s", err)
	}

	// Final path of the template zip file -> tempDirPath/template.zip
	finalTemplateZipFilePath := path.Join(tempDirPath, config.FILE_ZIP_TEMPLATE)

	// Path to a temporal folder where zip content will be drop -> tempDirPath/template
	temporalUnzippedFilesPath := path.Join(tempDirPath, config.FOLDER_TMP_TEMPLATE)

	// make a folder inside temporary directory to unzip the template
	utils.MakeDirectory(temporalUnzippedFilesPath, config.FOLDER_PERMISSION)
	if err != nil {
		utils.DeleteFileOrDirectory(tempDirPath) // delete temporary directory
		return fmt.Errorf("failed to create temporary directory: %s", err)
	}

	// Check if is necessary to download the zip from the internet or is a local zip
	if utils.IsGithubUrl(template) || utils.IsUrl(template) {
		if global.Verbose {
			fmt.Println("[+] Downloading template from url: " + template)
		}

		// Convert github project url to zip url (main branch) (use https protocol only)
		// Example: github.com/unknowns24/mks -> https://github.com/unknowns24/mks/archive/refs/heads/main.zip
		if utils.IsGithubUrl(template) {
			template = config.NETWORK_HTTPS_PREFIX + template
			if !strings.HasSuffix(template, config.NETWORK_MAIN_ZIP_ROUTE) {
				template = template + config.NETWORK_MAIN_ZIP_ROUTE
			}
		}

		// Download template zip file to temporary directory with name template.zip
		err = utils.DownloadFile(template, finalTemplateZipFilePath)
		if err != nil {
			utils.DeleteFileOrDirectory(tempDirPath) // delete temporary directory
			return fmt.Errorf("failed to download template: %s", err)
		}
	} else {
		if global.Verbose {
			fmt.Println("[+] Installing template from local path: " + template)
		}

		if !utils.FileOrDirectoryExists(template) {
			utils.DeleteFileOrDirectory(tempDirPath) // delete temporary directory
			return fmt.Errorf("template does not exist: %s", template)
		}
	}

	if global.Verbose {
		fmt.Println("[+] Unzipping template...")
	}

	// unzip template.zip to template directory inside temporary directory
	err = utils.Unzip(finalTemplateZipFilePath, temporalUnzippedFilesPath)
	if err != nil {
		utils.DeleteFileOrDirectory(tempDirPath) // delete temporary directory
		return fmt.Errorf("failed to unzip template: %s", err)
	}

	// delete template.zip to save space
	utils.DeleteFileOrDirectory(finalTemplateZipFilePath)
	if err != nil {
		utils.DeleteFileOrDirectory(tempDirPath) // delete temporary directory
		return fmt.Errorf("failed to unzip template: %s", err)
	}

	if global.Verbose {
		fmt.Println("[+] Checking if template is valid...")
	}

	// get template name from zip file root folder (must have only one folder, the name of this folder is the template name)
	dirs, err := utils.ListDirectories(temporalUnzippedFilesPath)
	if err != nil {
		utils.DeleteFileOrDirectory(tempDirPath) // delete temporary directory
		return fmt.Errorf("failed to list directories in template: %s", err)
	}

	// template must have only one folder on root, if not, return error
	if len(dirs) != 1 {
		utils.DeleteFileOrDirectory(tempDirPath) // delete temporary directory
		return fmt.Errorf("invalid template structure: %s", "template must have only one folder on root")
	}

	// get template name from folder name (the folder inside tempDirPath/template, tempDirPath/template/<FOLDER_NAME>, FOLDER_NAME is the template name)
	templateName := dirs[0]

	// check if template name has -main suffix and remove it (it occurs when the template is downloaded from github, by the user or by the program)
	if strings.HasSuffix(dirs[0], config.NETWORK_GITHUB_BRANCH_SUFFIX) {
		// delete -main suffix from template name
		templateName = strings.TrimSuffix(dirs[0], config.NETWORK_GITHUB_BRANCH_SUFFIX)
	}

	if global.Verbose {
		fmt.Println("[+] Checking if template is installed...")
	}

	// check if template exists in addons folder
	if utils.SliceContainsElement(global.InstalledTemplates, templateName) {
		utils.DeleteFileOrDirectory(tempDirPath) // delete temporary directory
		return fmt.Errorf("template already installed (uninstall it first): %s", templateName)
	}

	unzippedTemplateFilesPath := path.Join(temporalUnzippedFilesPath, dirs[0])

	templateFiles, _ := utils.ListDirectoriesAndFiles(unzippedTemplateFilesPath)

	haveTemplateFile := false

	if global.Verbose {
		fmt.Println("[+] Validating template files...")
	}

	// Iterate every file on unzipped template folder to search a template file
	for _, templateFile := range templateFiles {
		if strings.HasSuffix(templateFile, config.FILE_EXTENSION_TEMPLATE) {
			valid, err := utils.CheckSyntaxGoFile(path.Join(unzippedTemplateFilesPath, templateFile))

			fmt.Println(valid, err)
			fmt.Println(tempDirPath)

			if !valid || err != nil {
				utils.DeleteFileOrDirectory(tempDirPath) // delete temporary directory
				return fmt.Errorf("invalid template structure: %s", "template file must be a valid go file")
			}

			haveTemplateFile = true
			break
		}
	}

	// If not has a template file returns error
	if !haveTemplateFile {
		utils.DeleteFileOrDirectory(tempDirPath) // delete temporary directory
		return fmt.Errorf("invalid template structure: %s", "template must have a template file")
	}

	// Validate that this.load file (if exists) pass code validations
	fileToCheck := path.Join(unzippedTemplateFilesPath, config.FILE_ADDON_TEMPLATE_MAIN_LOAD)

	err = validateMksModulesFiles(fileToCheck, config.SPELL_FUNCION_LOAD_PREFIX, templateName, tempDirPath)
	if err != nil {
		return err
	}

	// Validate that this.unload file (if exists) pass code validations
	fileToCheck = path.Join(unzippedTemplateFilesPath, config.FILE_ADDON_TEMPLATE_MAIN_UNLOAD)

	err = validateMksModulesFiles(fileToCheck, config.SPELL_FUNCION_UNLOAD_PREFIX, templateName, tempDirPath)
	if err != nil {
		return err
	}

	// Validate if has a dependency file and in that case if its dependencies are installed on mks
	dependFile := path.Join(unzippedTemplateFilesPath, config.FILE_ADDON_TEMPLATE_DEPENDS)

	if utils.FileOrDirectoryExists(dependFile) {
		dependsOk, dependsMissing, err := utils.ValidateAllDependenciesInstalled(dependFile)

		if err != nil {
			utils.DeleteFileOrDirectory(tempDirPath) // delete temporary directory
			return fmt.Errorf("failed to validate dependencies: %s", err)
		}

		if !dependsOk {
			utils.DeleteFileOrDirectory(tempDirPath) // delete temporary directory
			return fmt.Errorf("template has missing dependencies: %s", strings.Join(dependsMissing, ", "))
		}
	}

	if global.Verbose {
		fmt.Println("[+] Installing template...")
	}

	// move template folder to addons folder and delete folder
	err = utils.MoveFileOrDirectory(unzippedTemplateFilesPath, path.Join(addonsPath, templateName))
	utils.DeleteFileOrDirectory(tempDirPath)

	if err != nil {
		return fmt.Errorf("failed to move template to addons folder: %s", err)
	}

	if global.Verbose {
		fmt.Println("[+] Template installed succesfully.")
	}

	return nil
}

func validateMksModulesFiles(fileToCheck, fileType, templateName, tempDir string) error {
	if utils.FileOrDirectoryExists(fileToCheck) {
		valid, err := utils.CheckSyntaxGoFile(fileToCheck)
		if err != nil || !valid {
			utils.DeleteFileOrDirectory(tempDir) // delete temporary directory
			return fmt.Errorf("invalid template structure: %s is not a valid go file", fileToCheck)
		}

		valid, err = utils.CheckPackageNameInFile(fileToCheck, config.SPELL_PACKAGE_MKS_MODULE)
		if err != nil || !valid {
			utils.DeleteFileOrDirectory(tempDir) // delete temporary directory
			return fmt.Errorf("invalid template structure: %s must have %s package name", fileToCheck, config.SPELL_PACKAGE_MKS_MODULE)
		}

		valid, err = utils.FunctionExistsInFile(fileToCheck, fileType+templateName)
		if err != nil || !valid {
			utils.DeleteFileOrDirectory(tempDir) // delete temporary directory
			return fmt.Errorf("invalid template structure: %s must have %s function", fileToCheck, fileType+templateName)
		}
	}

	return nil
}
