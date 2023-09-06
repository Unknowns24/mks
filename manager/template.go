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

	tempdir, err := utils.MakeTempDirectory()
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %s", err)
	}

	// current temporal template folder path
	temporalUnzippedPath := path.Join(tempdir, config.FOLDER_TMP_TEMPLATE)

	// current temporal template zip file path
	currentTemplateZipPath := path.Join(tempdir, config.FILE_ZIP_TEMPLATE)

	// make a directory inside temporary directory to unzip the template
	utils.MakeDirectoryIfNotExists(temporalUnzippedPath, 0755)
	if err != nil {

		utils.DeleteFileOrDirectory(tempdir) // delete temporary directory
		return fmt.Errorf("failed to create temporary directory: %s", err)
	}

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
		err = utils.DownloadFile(template, currentTemplateZipPath)
		if err != nil {
			utils.DeleteFileOrDirectory(tempdir) // delete temporary directory
			return fmt.Errorf("failed to download template: %s", err)
		}

		if global.Verbose {
			fmt.Println("[+] Unzipping template...")
		}
		// unzip template.zip to template directory inside temporary directory
		err = utils.Unzip(currentTemplateZipPath, temporalUnzippedPath)

		// delete template.zip to save space
		utils.DeleteFileOrDirectory(currentTemplateZipPath)
		if err != nil {
			utils.DeleteFileOrDirectory(tempdir) // delete temporary directory
			return fmt.Errorf("failed to unzip template: %s", err)
		}
	} else {
		if global.Verbose {
			fmt.Println("[+] Installing template from local path: " + template)
		}

		if !utils.FileOrDirectoryExists(template) {
			utils.DeleteFileOrDirectory(tempdir) // delete temporary directory
			return fmt.Errorf("template does not exist: %s", template)
		}

		if global.Verbose {
			fmt.Println("[+] Unzipping template...")
		}

		// unzip template.zip to template directory inside temporary directory
		err = utils.Unzip(template, temporalUnzippedPath)
		if err != nil {
			utils.DeleteFileOrDirectory(tempdir) // delete temporary directory
			return fmt.Errorf("failed to unzip template: %s", err)
		}
	}

	if global.Verbose {
		fmt.Println("[+] Checking if template is valid...")
	}

	// get template name from zip file root folder (must have only one folder, the name of this folder is the template name)
	dirs, err := utils.ListDirectories(temporalUnzippedPath)
	if err != nil {

		utils.DeleteFileOrDirectory(tempdir) // delete temporary directory
		return fmt.Errorf("failed to list directories in template: %s", err)
	}

	// template must have only one folder on root, if not, return error
	if len(dirs) != 1 {

		utils.DeleteFileOrDirectory(tempdir) // delete temporary directory
		return fmt.Errorf("invalid template structure: %s", "template must have only one folder on root")
	}

	// get template name from folder name (the folder inside tempdir/template, tempdir/template/<FOLDER_NAME>, FOLDER_NAME is the template name)
	templateName := dirs[0]

	// check if template name has -main suffix and remove it (it occurs when the template is downloaded from github, by the user or by the program)
	if strings.HasSuffix(dirs[0], "-main") {
		// delete -main suffix from template name
		templateName = strings.TrimSuffix(dirs[0], "-main")
	}

	if global.Verbose {
		fmt.Println("[+] Checking if template is installed...")
	}

	// check if template exists in addons folder
	if utils.SliceContainsElement(global.InstalledTemplates, templateName) {
		utils.DeleteFileOrDirectory(tempdir) // delete temporary directory
		return fmt.Errorf("template already installed (uninstall it first): %s", templateName)
	}

	temporalUnzippedTemplatePath := path.Join(temporalUnzippedPath, dirs[0])

	templateFiles, _ := utils.ListDirectoriesAndFiles(temporalUnzippedTemplatePath)

	haveTemplateFile := false

	for _, templateFile := range templateFiles {
		if strings.HasSuffix(templateFile, config.FILE_TEMPLATE_EXTENSION) {

			valid, err := utils.IsPseudoValidGoFile(templateFile)

			if !valid || err != nil {
				utils.DeleteFileOrDirectory(tempdir) // delete temporary directory
				return fmt.Errorf("invalid template structure: %s", "template file must be a valid go file")
			}

			haveTemplateFile = true
			break
		}
	}

	fileToCheck := path.Join(temporalUnzippedTemplatePath, config.FILE_ADDON_TEMPLATE_MAIN_LOAD)

	if utils.FileOrDirectoryExists(fileToCheck) {
		valid, err := utils.CheckSyntaxGoFile(fileToCheck)
		if err != nil {
			utils.DeleteFileOrDirectory(tempdir) // delete temporary directory
			return fmt.Errorf("invalid template structure: %s can't be checked", fileToCheck)
		}
		if !valid {
			utils.DeleteFileOrDirectory(tempdir) // delete temporary directory
			return fmt.Errorf("invalid template structure: %s is not a valid go file", fileToCheck)
		}

		valid, err = utils.CheckPackageNameInFile(fileToCheck, config.SPELL_PACKAGE_MKS_MODULE)
		if err != nil {
			utils.DeleteFileOrDirectory(tempdir) // delete temporary directory
			return fmt.Errorf("invalid template structure: can't check package name in %s", fileToCheck)
		}
		if !valid {
			utils.DeleteFileOrDirectory(tempdir) // delete temporary directory
			return fmt.Errorf("invalid template structure: %s must have %s package name", fileToCheck, config.SPELL_PACKAGE_MKS_MODULE)
		}

		valid, err = utils.FunctionExistsInFile(fileToCheck, config.SPELL_FUNCION_LOAD_PREFIX+templateName)
		if err != nil {
			utils.DeleteFileOrDirectory(tempdir) // delete temporary directory
			return fmt.Errorf("invalid template structure: can't check function name in %s", fileToCheck)
		}
		if !valid {
			utils.DeleteFileOrDirectory(tempdir) // delete temporary directory
			return fmt.Errorf("invalid template structure: %s must have %s function", fileToCheck, config.SPELL_FUNCION_LOAD_PREFIX+templateName)
		}

	}

	fileToCheck = path.Join(temporalUnzippedTemplatePath, config.FILE_ADDON_TEMPLATE_MAIN_UNLOAD)

	if utils.FileOrDirectoryExists(fileToCheck) {
		valid, err := utils.CheckSyntaxGoFile(fileToCheck)
		if err != nil {
			utils.DeleteFileOrDirectory(tempdir) // delete temporary directory
			return fmt.Errorf("invalid template structure: %s can't be checked", fileToCheck)
		}
		if !valid {
			utils.DeleteFileOrDirectory(tempdir) // delete temporary directory
			return fmt.Errorf("invalid template structure: %s is not a valid go file", fileToCheck)
		}

		valid, err = utils.CheckPackageNameInFile(fileToCheck, config.SPELL_PACKAGE_MKS_MODULE)
		if err != nil {
			utils.DeleteFileOrDirectory(tempdir) // delete temporary directory
			return fmt.Errorf("invalid template structure: can't check package name in %s", fileToCheck)
		}
		if !valid {
			utils.DeleteFileOrDirectory(tempdir) // delete temporary directory
			return fmt.Errorf("invalid template structure: %s must have %s package name", fileToCheck, config.SPELL_PACKAGE_MKS_MODULE)
		}

		valid, err = utils.FunctionExistsInFile(fileToCheck, config.SPELL_FUNCION_UNLOAD_PREFIX+templateName)
		if err != nil {
			utils.DeleteFileOrDirectory(tempdir) // delete temporary directory
			return fmt.Errorf("invalid template structure: can't check function name in %s", fileToCheck)
		}
		if !valid {
			utils.DeleteFileOrDirectory(tempdir) // delete temporary directory
			return fmt.Errorf("invalid template structure: %s must have %s function", fileToCheck, config.SPELL_FUNCION_UNLOAD_PREFIX+templateName)
		}

	}

	if !haveTemplateFile {
		utils.DeleteFileOrDirectory(tempdir) // delete temporary directory
		return fmt.Errorf("invalid template structure: %s", "template must have a template file")
	}

	dependFile := path.Join(temporalUnzippedTemplatePath, config.FILE_ADDON_TEMPLATE_DEPENDS)

	if utils.FileOrDirectoryExists(dependFile) {
		dependsOk, dependsMissing, err := utils.ValidateAllDependenciesInstalled(dependFile)

		if err != nil {
			utils.DeleteFileOrDirectory(tempdir) // delete temporary directory
			return fmt.Errorf("failed to validate dependencies: %s", err)
		}

		if !dependsOk {
			utils.DeleteFileOrDirectory(tempdir) // delete temporary directory
			return fmt.Errorf("template has missing dependencies: %s", strings.Join(dependsMissing, ", "))
		}
	}

	if global.Verbose {
		fmt.Println("[+] Installing template...")
	}

	// move template folder to addons folder
	err = utils.MoveFileOrDirectory(temporalUnzippedTemplatePath, path.Join(addonsPath, templateName))

	utils.DeleteFileOrDirectory(tempdir) // delete temporary directory

	if err != nil {
		return fmt.Errorf("failed to move template to addons folder: %s", err)
	}

	return nil
}
