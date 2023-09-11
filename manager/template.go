package manager

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
	"github.com/unknowns24/mks/utils"
)

type dependsFileFormat struct {
	DependsOn []string `json:"dependsOn"`
}

func ListTemplate() {
	fmt.Println("[+] List of templates installed and availables to use:")
	if len(global.InstalledTemplates) == 0 {
		fmt.Println(" └──── No templates installed")
		return
	}

	for k, currentTemplaName := range global.InstalledTemplates {
		if k == len(global.InstalledTemplates)-1 {
			fmt.Println(" └──── " + currentTemplaName)
		} else {
			fmt.Println(" ├──── " + currentTemplaName)
		}
	}

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

	templateFolderPath := path.Join(global.UserTemplatesFolderPath, template)

	if utils.FileOrDirectoryExists(templateFolderPath) {
		utils.DeleteFileOrDirectory(templateFolderPath) // delete template directory
	}

	if utils.FileOrDirectoryExists(templateFolderPath) {
		return fmt.Errorf("failed to uninstall template: %s (Try delete folder %s manually)", template, templateFolderPath)
	}

	if global.Verbose {
		fmt.Println("[+] Template uninstalled successfully!")
	}

	return nil
}

func installTemplateFiles(templateRootDir string, useFlag []string) error {

	if global.Verbose {
		fmt.Println("[+] Copying template files to addons folder...")
	}

	useTemplates := useFlag
	err := error(nil)

	if len(useTemplates) == 0 {
		useTemplates, err = utils.ListDirectories(templateRootDir)
		if err != nil {
			return err
		}
	}

	for _, currentTemplaName := range useTemplates {
		currentTemplateOriginPath := path.Join(templateRootDir, currentTemplaName)
		currentTemplateDestinationPath := path.Join(global.UserTemplatesFolderPath, currentTemplaName)
		err := utils.CopyFileOrDirectory(currentTemplateOriginPath, currentTemplateDestinationPath)
		if err != nil {
			return fmt.Errorf("failed to install template %s to addons folder: %s", currentTemplaName, err)
		}

		if global.Verbose {
			fmt.Printf("[+] Template %s installed succesfully.\n", currentTemplaName)
		}
	}

	return nil
}

// Check if template is valid and return root directory of template(s), Only checks useFlag templates, if useFlag is empty, check all templates.
func checkTemplateFiles(templateRootDir string, useFlag []string) (string, error) {

	if global.Verbose {
		fmt.Println("[+] Searching for template(s) files...")
	}

	useTemplates := useFlag

	// get folders in template root dir
	dirs, err := utils.ListDirectories(templateRootDir)
	if err != nil {
		return "", fmt.Errorf("failed to detect template(s): %s", err)
	}

	// check if templateRootDir has a folder with -main suffix, and only is one  (it occurs when the template is downloaded from github, by the user or by the program)
	// template must have only one folder on root and this folder has a github branch "-main" suffix, use this as template root dir
	if len(dirs) == 1 && strings.HasSuffix(dirs[0], config.NETWORK_GITHUB_BRANCH_SUFFIX) {
		if global.Verbose {
			fmt.Println("[+] Github template detected, using " + dirs[0] + " as template root dir...")
		}

		templateRootDir = path.Join(templateRootDir, dirs[0])

		// get template(s) name(s) from zip file root folder (could have only one or multiple folders, the name of this folder(s) is the template name(s). Files are ignored.)
		dirs, err = utils.ListDirectories(templateRootDir)
		if err != nil {
			return "", fmt.Errorf("failed to detect template(s): %s", err)
		}
	}

	// template must have only one folder on root, if not, return error
	if len(dirs) < 1 {
		return "", fmt.Errorf("invalid template structure: %s", "almost one template folder is required on root of template")
	}

	// get template(s) name(s) from folder(s) name(s) (the folder inside templates root dir, templaresRootDir/<FOLDER_NAME>, FOLDER_NAME is the template name)
	templatesAvailablesToInstall := dirs

	templatesAlreadyInstalled := []string{}

	// add to templatesAlreadyInstalled templates that are already installed
	// check if templatesAvailablesToInstall has andy element wich contains '-' symbol, if true, return error
	for _, currentTemplaName := range templatesAvailablesToInstall {

		if utils.SliceContainsElement(global.InstalledTemplates, currentTemplaName) {
			templatesAlreadyInstalled = append(templatesAlreadyInstalled, currentTemplaName)
		}

		if strings.Contains(currentTemplaName, "-") {
			return "", fmt.Errorf("invalid template structure: %s template name must not contain '-' symbol, use '_' instead", currentTemplaName)
		}
	}

	if len(templatesAlreadyInstalled) > 1 {
		return "", fmt.Errorf("some templates are already installed: %s", strings.Join(templatesAlreadyInstalled, ", "))
	} else if len(templatesAlreadyInstalled) == 1 {
		return "", fmt.Errorf("template already installed: %s", templatesAlreadyInstalled[0])
	}

	if global.Verbose {
		fmt.Println("[+] Checking if template is installed...")
	}

	templatesNotAvailablesToInstall := []string{}

	if len(useTemplates) > 0 {
		for _, currentTemplaName := range useTemplates {

			// check if template exists in addons folder
			if utils.SliceContainsElement(templatesAvailablesToInstall, currentTemplaName) {
				templatesNotAvailablesToInstall = append(templatesNotAvailablesToInstall, currentTemplaName)
			}
		}
	} else {
		useTemplates = templatesAvailablesToInstall
	}

	if len(templatesNotAvailablesToInstall) > 1 {
		return "", fmt.Errorf("templates are no available to install: %s", strings.Join(templatesNotAvailablesToInstall, ", "))
	} else if len(templatesNotAvailablesToInstall) == 1 {
		return "", fmt.Errorf("template is no available to install: %s", templatesNotAvailablesToInstall[0])
	}

	for _, currentTemplaName := range useTemplates {
		templateFiles, err := utils.ListDirectoriesAndFiles(path.Join(templateRootDir, currentTemplaName))
		if err != nil {
			return "", fmt.Errorf("failed to detect files in template %s: %s", currentTemplaName, err)
		}

		haveTemplateFile := false

		if global.Verbose {
			fmt.Printf("[+] Validating files in %s template...\n", currentTemplaName)
		}

		// Iterate every file on unzipped template folder to search a template file
		for _, templateFile := range templateFiles {
			if strings.HasSuffix(templateFile, config.FILE_EXTENSION_TEMPLATE) {
				valid, err := utils.CheckSyntaxGoFile(path.Join(templateRootDir, currentTemplaName, templateFile))

				if !valid || err != nil {
					return "", fmt.Errorf("invalid template structure: %s", "template file must be a valid go file")
				}

				haveTemplateFile = true
			}
		}

		// If not has a template file returns error
		if !haveTemplateFile {
			return "", fmt.Errorf("invalid template structure: %s", "template must have a template file")
		}

		// Validate that this.load file (if exists) pass code validations
		fileToCheck := path.Join(templateRootDir, currentTemplaName, config.FILE_ADDON_TEMPLATE_MAIN_LOAD)

		err = validateMksModulesFiles(fileToCheck, config.SPELL_FUNCION_LOAD_PREFIX, currentTemplaName)
		if err != nil {
			return "", err
		}

		// Validate that this.unload file (if exists) pass code validations
		fileToCheck = path.Join(templateRootDir, currentTemplaName, config.FILE_ADDON_TEMPLATE_MAIN_UNLOAD)

		err = validateMksModulesFiles(fileToCheck, config.SPELL_FUNCION_UNLOAD_PREFIX, currentTemplaName)
		if err != nil {
			return "", err
		}

		// Validate if has a dependency file and in that case if its dependencies are installed on mks
		dependFile := path.Join(templateRootDir, currentTemplaName, config.FILE_ADDON_TEMPLATE_DEPENDS)

		if utils.FileOrDirectoryExists(dependFile) {

			if global.Verbose {

				fileContent, err := os.ReadFile(dependFile)
				if err != nil {
					return "", err
				}

				var parsedFile dependsFileFormat

				// Parse json file and save data on parsedFile variable
				err = json.Unmarshal([]byte(fileContent), &parsedFile)
				if err != nil {
					return "", err
				}

				fmt.Printf("[*] This template uses dependencies, make sure you have them installed when using it in your project: %s\n", strings.Join(parsedFile.DependsOn, ", "))
			}
		}

	}
	return templateRootDir, nil
}

// download template to zip cache folder
func downloadTemplateToCache(template string) error {

	//do a md5 of template (template is an address) and check if exists in zip cache folder
	templateAddrHash := utils.GetMD5Hash(template)

	// path for this zip file template
	zipCachePath := path.Join(global.ZipCachePath, templateAddrHash)

	// path for this zip file template
	temporalZipCachePath := path.Join(global.TemporalsPath, templateAddrHash+config.FILE_EXTENSION_ZIP)

	// Path to a temporal folder where zip content will be drop
	temporalUnzippedFilesPath := path.Join(global.TemporalsPath, templateAddrHash)

	tempDirPath, err := utils.MakeTemporalDirectory()
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %s", err)
	}

	// make a folder inside temporary directory to unzip the template
	err = utils.MakeDirectory(temporalUnzippedFilesPath, config.FOLDER_PERMISSION)
	if err != nil {
		utils.DeleteFileOrDirectory(tempDirPath) // delete temporary directory
		return fmt.Errorf("failed to create temporary directory: %s", err)
	}

	if global.Verbose {
		fmt.Println("[+] Downloading template from: " + template)
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
	err = utils.DownloadFile(template, temporalZipCachePath)
	if err != nil {
		utils.DeleteFileOrDirectory(temporalZipCachePath) // delete downloaded zip file (if exists)
		utils.DeleteFileOrDirectory(tempDirPath)          // delete temporary directory
		return fmt.Errorf("failed to download template: %s", err)
	}

	err = utils.CheckZipIntegrity(temporalZipCachePath)
	if err != nil {
		utils.DeleteFileOrDirectory(temporalZipCachePath) // delete downloaded zip file (if exists)
		utils.DeleteFileOrDirectory(tempDirPath)          // delete temporary directory
		return err
	}

	err = utils.MoveFileOrDirectory(temporalZipCachePath, zipCachePath)
	if err != nil {
		utils.DeleteFileOrDirectory(zipCachePath)         // delete downloaded zip file (if exists)
		utils.DeleteFileOrDirectory(temporalZipCachePath) // delete downloaded zip file (if exists)
		utils.DeleteFileOrDirectory(tempDirPath)          // delete temporary directory
		return fmt.Errorf("failed to move downloaded zip file to cache directory\n %s -> %s", temporalZipCachePath, zipCachePath)
	}

	return nil
}

// unzip template (cached) to templates cache folder
func unzipTemplateCached(zipCacheName string) error {

	// path for this zip file template
	zipCachePath := path.Join(global.ZipCachePath, zipCacheName)

	// temporal path for this unzipped template
	temporalUnzippedFilesPath := path.Join(global.TemporalsPath, zipCacheName)

	// path for this template folder
	templateCachePath := path.Join(global.TemplateCachePath, zipCacheName)

	// unzip template.zip to template directory inside temporary directory
	err := utils.Unzip(zipCachePath, temporalUnzippedFilesPath)
	if err != nil {
		utils.DeleteFileOrDirectory(temporalUnzippedFilesPath) // delete temporary directory
		return fmt.Errorf("failed to unzip template: %s", err)
	}

	err = utils.MoveFileOrDirectory(temporalUnzippedFilesPath, templateCachePath)
	if err != nil {
		utils.DeleteFileOrDirectory(temporalUnzippedFilesPath) // delete temporary directory
		return fmt.Errorf("failed to move unzipped template: %s", err)
	}

	return nil
}

// unzip template inside temporary directory and return the path to the unzipped template
func unzipTemplateLocaldisk(zipLocalDisk string) (string, error) {
	// temporal path for this unzipped template
	temporalUnzippedFilesPath, err := utils.MakeTemporalDirectory()
	if err != nil {
		return "", err
	}

	// unzip template zip to temporary directory
	err = utils.Unzip(zipLocalDisk, temporalUnzippedFilesPath)
	if err != nil {
		utils.DeleteFileOrDirectory(temporalUnzippedFilesPath) // delete temporary directory
		return "", fmt.Errorf("failed to unzip template: %s", err)
	}

	return temporalUnzippedFilesPath, nil
}

func InstallTemplate(template string, useFlag []string) error {

	fmt.Println("[+] Installing " + template + " template...")

	//do a md5 of template (template is an address) and check if exists in zip cache folder
	templateAddrHash := utils.GetMD5Hash(template)

	// path for this zip file template
	zipCachePath := path.Join(global.ZipCachePath, templateAddrHash)

	// path for this template folder
	templateCachePath := path.Join(global.TemplateCachePath, templateAddrHash)

	var err error = nil

	if utils.FileOrDirectoryExists(templateCachePath) {
		if global.Verbose {
			fmt.Println("[+] Template already downloaded, using cached files...")
			fmt.Println(" └──── if you want a fresh download, delete the cache using 'mks clear'")

			fmt.Println("[+] Installing template...")
		}

		templateRootDir, err := checkTemplateFiles(templateCachePath, useFlag)
		if err != nil {
			return err
		}

		return installTemplateFiles(templateRootDir, useFlag)
	}

	// if zip file not exists in cache try to download it (except if it is a local zip file)
	if !utils.FileOrDirectoryExists(zipCachePath) {
		// Check if is necessary to download the zip from the internet or is a local zip
		if utils.IsGithubUrl(template) || utils.IsUrl(template) {
			err = downloadTemplateToCache(template)
			if err != nil {
				return fmt.Errorf("failed to download template: %s", err)
			}
		} else {
			// if is a local zip file, unzip it and install it
			if global.Verbose {
				fmt.Println("[+] Installing template from local path: " + template)
			}
			// check if template zip file exists, if not return error
			if !utils.FileOrDirectoryExists(template) {
				return fmt.Errorf("template does not exist: %s", template)
			}

			// unzip template.zip to template directory inside temporary directory
			templateLocalDiskPath, err := unzipTemplateLocaldisk(template)
			if err != nil {
				// delete temporary directory
				utils.DeleteFileOrDirectory(templateLocalDiskPath)
				return fmt.Errorf("failed to unzip template: %s", err)
			}

			// check if template "was unzipped correctly" (directory exists)
			if !utils.FileOrDirectoryExists(templateLocalDiskPath) {
				// delete temporary directory
				utils.DeleteFileOrDirectory(templateLocalDiskPath)
				return fmt.Errorf("failed to unzip template: %s", "template on cache does not exist")
			}

			if global.Verbose {
				fmt.Println("[+] Checking if template is valid...")
			}

			templateRootDir, err := checkTemplateFiles(templateLocalDiskPath, useFlag)
			if err != nil {
				// delete temporary directory
				utils.DeleteFileOrDirectory(templateLocalDiskPath)
				return err
			}

			if global.Verbose {
				fmt.Println("[+] Installing template...")
			}

			// install template files (each one of useFlag templates, if useFlag is empty, install all templates presents)
			retValue := installTemplateFiles(templateRootDir, useFlag)

			// delete temporary directory
			utils.DeleteFileOrDirectory(templateLocalDiskPath)

			return retValue
		}
	} else {
		if global.Verbose {
			fmt.Println("[+] Template already downloaded, using cached zip file...")
			fmt.Println(" └──── if you want a fresh download, delete the zip cache using 'mks clear'")
		}
	}

	// if files not exists in cache, but zip file exists, unzip to cache and install it....

	if global.Verbose {
		fmt.Println("[+] Unzipping template...")
	}

	// unzip template to template cache folder
	err = unzipTemplateCached(templateAddrHash)
	if err != nil {
		return fmt.Errorf("failed to unzip template: %s", err)
	}

	// check if template "was unzipped correctly" (directory exists)
	if !utils.FileOrDirectoryExists(templateCachePath) {
		return fmt.Errorf("failed to unzip template: %s", "template on cache does not exist")
	}

	if global.Verbose {
		fmt.Println("[+] Checking if template is valid...")
	}

	templateRootDir, err := checkTemplateFiles(templateCachePath, useFlag)
	if err != nil {
		return err
	}

	if global.Verbose {
		fmt.Println("[+] Installing template...")
	}

	// install template files (each one of useFlag templates, if useFlag is empty, install all templates presents)
	return installTemplateFiles(templateRootDir, useFlag)
}

// Validate mks modules files: load, unload. Validations: oo sintax, package name, function name
func validateMksModulesFiles(fileToCheck, fileType, templateName string) error {
	if utils.FileOrDirectoryExists(fileToCheck) {
		valid, err := utils.CheckSyntaxGoFile(fileToCheck)
		if err != nil || !valid {
			return fmt.Errorf("invalid template structure: %s is not a valid go file", fileToCheck)
		}

		valid, err = utils.CheckPackageNameInFile(fileToCheck, config.SPELL_PACKAGE_MKS_MODULE)
		if err != nil || !valid {
			return fmt.Errorf("invalid template structure: %s must have %s package name", fileToCheck, config.SPELL_PACKAGE_MKS_MODULE)
		}

		valid, err = utils.FunctionExistsInFile(fileToCheck, fmt.Sprintf("%s%s", fileType, templateName))
		if err != nil || !valid {
			return fmt.Errorf("invalid template structure: %s must have %s function", fileToCheck, fmt.Sprintf("%s%s", fileType, templateName))
		}
	}
	return nil
}
