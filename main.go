package main

import (
	"github.com/unknowns24/mks/cmd"
	"github.com/unknowns24/mks/utils"
)

func main() {
	/***************************
	* SETTING GLOBAL VARIABLES *
	****************************/

	err := utils.SetExecutablePath()
	if err != nil {
		panic(err)
	}

	err = utils.SetTemplatesFolderPathGlobal()
	if err != nil {
		panic(err)
	}

	err = utils.SetCurrentInstalledTemplates()
	if err != nil {
		panic(err)
	}

	err = utils.SetUserConfigFolderPath()
	if err != nil {
		panic(err)
	}

	/******************
	* EXECUTING COBRA *
	*******************/
	cmd.Execute()
}
