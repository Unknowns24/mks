package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/unknowns24/mks/utils"
)

var rootCmd = &cobra.Command{
	Use:   "mks",
	Short: "Microservice manager CLI",
}

func Execute() {
	// Set global variables
	utils.SetTemplatesFolderPathGlobal()

	// Start main command
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(NewBuildCmd())
	rootCmd.AddCommand(NewAddCmd())
}
