package manager

import (
	"fmt"
	"path"
	"strings"

	"github.com/unknowns24/mks/global"
	"github.com/unknowns24/mks/utils"
)

func InstallTemplate(template string) error {
	//TODO: Implement template installation

	fmt.Println("[+] Installing " + template + " template...")

	if global.Verbose {
		fmt.Println("[+] Checking if template is valid...")
	}

	addonsPath := path.Join(global.TemplatesFolderPath, "addons")
	installedTemplates, err := utils.ListDirectories(addonsPath)

	if err != nil {
		return fmt.Errorf("failed to list installed templates: %s", err)
	}

	tempdir, err := utils.MakeTempDirectory()
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %s", err)
	}

	// make a directory inside temporary directory to unzip the template
	utils.MakeDirectoryIfNotExists(path.Join(tempdir, "template"), 0755)

	// error at make directory?
	if err != nil {
		// delete temporary directory
		utils.DeleteFileOrDirectory(tempdir)
		// return error
		return fmt.Errorf("failed to create temporary directory: %s", err)
	}

	if utils.IsGithubUrl(template) || utils.IsUrl(template) {
		if global.Verbose {
			fmt.Println("[+] Downloading template from url: " + template)
		}

		// Convert github project url to zip url (main branch) (use https protocol only)
		// Example: github.com/unknowns24/mks -> https://github.com/unknowns24/mks/archive/refs/heads/main.zip
		if utils.IsGithubUrl(template) {
			template = "https://" + template
			if !strings.HasSuffix(template, "/archive/refs/heads/main.zip") {
				template = template + "/archive/refs/heads/main.zip"
			}
		}

		// Download template zip file to temporary directory with name template.zip
		err = utils.DownloadFile(template, path.Join(tempdir, "template.zip"))

		// error at download file?
		if err != nil {
			// delete temporary directory
			utils.DeleteFileOrDirectory(tempdir)
			// return error
			return fmt.Errorf("failed to download template: %s", err)
		}

		if global.Verbose {
			fmt.Println("[+] Unzipping template...")
		}
		// unzip template.zip to template directory inside temporary directory
		err = utils.Unzip(path.Join(tempdir, "template.zip"), path.Join(tempdir, "template"))

		// delete template.zip to save space
		utils.DeleteFileOrDirectory(path.Join(tempdir, "template.zip"))

		// error at unzip?
		if err != nil {
			// delete temporary directory
			utils.DeleteFileOrDirectory(tempdir)
			// return error
			return fmt.Errorf("failed to unzip template: %s", err)
		}
	} else {

		if global.Verbose {
			fmt.Println("[+] Installing template from local path: " + template)
		}

		if !utils.FileOrDirectoryExists(template) {
			// delete temporary directory
			utils.DeleteFileOrDirectory(tempdir)
			// return error
			return fmt.Errorf("template does not exist: %s", template)
		}

		if global.Verbose {
			fmt.Println("[+] Unzipping template...")
		}
		// unzip template.zip to template directory inside temporary directory
		err = utils.Unzip(template, path.Join(tempdir, "template"))

		// error at unzip?
		if err != nil {
			// delete temporary directory
			utils.DeleteFileOrDirectory(tempdir)
			// return error
			return fmt.Errorf("failed to unzip template: %s", err)
		}
	}

	/*
		TODO:
			in the zip file, the template is in a folder with the same name as the zip file with suffix -main (example: template-main)
			unzip the file to a temporary directory
			check if the template folder exists in the temporary directory
			if exists, check if template has already been installed
			if not has installed, check integrity of the template
			if integrity is ok, install the template in addons folder
	*/

	//prints the tempdir path for debugging
	//fmt.Printf("tempdir: %s\n", tempdir)

	if global.Verbose {
		fmt.Println("[+] Checking if template is valid...")
	}
	// get template name from zip file root folder (must have only one folder, the name of this folder is the template name)
	dirs, err := utils.ListDirectories(path.Join(tempdir, "template"))

	// error at list directories?
	if err != nil {
		// delete temporary directory
		utils.DeleteFileOrDirectory(tempdir)
		// return error
		return fmt.Errorf("failed to list directories in template: %s", err)
	}

	// template must have only one folder on root, if not, return error
	if len(dirs) != 1 {
		// delete temporary directory
		utils.DeleteFileOrDirectory(tempdir)
		// return error
		return fmt.Errorf("invalid template structure: %s", "template must have only one folder on root")
	}

	// get template name from folder name (the folder inside tempdir/template, tempdir/template/<FOLDER_NAME>, FOLDER_NAME is the template name)
	tplName := dirs[0]

	// check if template name has -main suffix and remove it (it occurs when the template is downloaded from github, by the user or by the program)
	if strings.HasSuffix(dirs[0], "-main") {
		// delete -main suffix from template name
		tplName = strings.TrimSuffix(dirs[0], "-main")
	}

	// TODO: Other integrity checks (check if template has all required files, etc)
	// TODO: Other integrity checks (check if template has all required files, etc)
	// TODO: Other integrity checks (check if template has all required files, etc)

	if global.Verbose {
		fmt.Println("[+] Checking if template is installed...")
	}
	//check if template exists in addons folder
	if utils.SliceContainsElement(installedTemplates, tplName) {
		// delete temporary directory
		utils.DeleteFileOrDirectory(tempdir)
		// return error
		return fmt.Errorf("template already installed (uninstall it first): %s", tplName)
	} else {

		if global.Verbose {
			fmt.Println("[+] Installing template...")
		}

		// move template folder to addons folder
		err = utils.MoveFileOrDirectory(
			path.Join(tempdir, "template", dirs[0]),
			path.Join(addonsPath, tplName))

		// error at move file or directory?
		if err != nil {
			// delete temporary directory
			utils.DeleteFileOrDirectory(tempdir)
			// return error
			return fmt.Errorf("failed to move template to addons folder: %s", err)
		}
	}

	// delete temporary directory
	utils.DeleteFileOrDirectory(tempdir)

	return nil
}
