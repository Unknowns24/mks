package main

import (
	"github.com/unknowns24/mks/cmd"
	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
	"github.com/unknowns24/mks/utils"
)

func main() {
	/***************************
	* SETTING GLOBAL VARIABLES *
	****************************/

	err := utils.SetMksTemplatesFolderPath()
	if err != nil {
		panic(err)
	}

	err = utils.SetUserConfigFolderPath()
	if err != nil {
		panic(err)
	}

	err = utils.SetCurrentInstalledTemplates()
	if err != nil {
		panic(err)
	}

	err = utils.SetCacheFoldersPath()
	if err != nil {
		panic(err)
	}

	err = utils.SetTemporalsPath()
	if err != nil {
		panic(err)
	}

	// Create cache and temporal folders if not exist
	utils.MakeDirectory(global.ZipCachePath, config.FOLDER_PERMISSION)
	utils.MakeDirectory(global.TemplateCachePath, config.FOLDER_PERMISSION)
	utils.MakeDirectory(global.TemporalsPath, config.FOLDER_PERMISSION)
	utils.MakeDirectory(global.UserTemplatesFolderPath, config.FOLDER_PERMISSION)

	/******************
	* EXECUTING COBRA *
	*******************/
	cmd.Execute()
}
