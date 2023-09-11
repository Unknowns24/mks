package main

import (
	"os"

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

	err = utils.SetExportsPath()
	if err != nil {
		panic(err)
	}

	// Create cache, temporal and exports folders if not exist
	os.MkdirAll(global.ZipCachePath, config.FOLDER_PERMISSION)
	os.MkdirAll(global.TemplateCachePath, config.FOLDER_PERMISSION)
	os.MkdirAll(global.TemporalsPath, config.FOLDER_PERMISSION)
	os.MkdirAll(global.UserTemplatesFolderPath, config.FOLDER_PERMISSION)
	os.MkdirAll(global.ExportPath, config.FOLDER_PERMISSION)

	/******************
	* EXECUTING COBRA *
	*******************/
	cmd.Execute()
}
