package manager

import (
	"fmt"
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

	//addonsPath := path.Join(global.TemplatesFolderPath, "addons")
	//installedTemplates, err := utils.ListDirectories(addonsPath)

	if utils.IsGithubUrl(template) || utils.IsUrl(template) {
		if global.Verbose {
			fmt.Println("[+] Downloading template from url: " + template)
		}
		tempdir, err := utils.MakeTempDirectory()
		if err != nil {
			return fmt.Errorf("failed to create temporary directory: %s", err)
		}
		//defer utils.DeleteFileOrDirectory(tempdir)

		// Convert github project url to zip url (main branch) (use https protocol only)
		// Example: github.com/unknowns24/mks -> https://github.com/unknowns24/mks/archive/refs/heads/main.zip
		if utils.IsGithubUrl(template) {
			template = "https://" + template
			if !strings.HasSuffix(template, "/archive/refs/heads/main.zip") {
				template = template + "/archive/refs/heads/main.zip"
			}
		}

		utils.DownloadFile(template, tempdir+"/template.zip")

		//prints the tempdir path for debugging
		fmt.Printf("tempdir: %s\n", tempdir)

		/*
			TODO:
				in the zip file, the template is in a folder with the same name as the zip file with suffix -main (example: template-main)
				unzip the file to a temporary directory
				check if the template folder exists in the temporary directory
				if exists, check if template has already been installed
				if not has installed, check integrity of the template
				if integrity is ok, install the template in addons folder
		*/

		//TODO: utils.Log("Installing template from url: " + template)
	} else {

		if global.Verbose {
			fmt.Println("[+] Installing template from local path: " + template)
		}

		/*
			TODO:
				Do the same as above, but without downloading the file
		*/

		//TODO: utils.Log("Installing template from local path: " + template)
	}

	return nil
}
