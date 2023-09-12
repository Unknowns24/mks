package main

import (
	"github.com/unknowns24/mks/cmd"
	"github.com/unknowns24/mks/manager"
)

func main() {
	/***************************
	* SETTING GLOBAL VARIABLES *
	****************************/

	if err := manager.SetupMks(); err != nil {
		panic(err)
	}

	/******************
	* EXECUTING COBRA *
	*******************/
	cmd.Execute()
}
