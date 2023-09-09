package main

import (
	"github.com/unknowns24/mks/cmd"
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

	/******************
	* EXECUTING COBRA *
	*******************/
	cmd.Execute()
}
