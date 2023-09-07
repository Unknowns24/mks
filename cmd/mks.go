package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/unknowns24/mks/utils"
)

var rootCmd = &cobra.Command{
	Use:   "mks",
	Short: "Golang application manager CLI",
}

func Execute() {
	// Set global variables
	utils.SetExecutablePath()
	utils.SetTemplatesFolderPathGlobal()
	utils.SetCurrentInstalledTemplates()

	// Start main command
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(NewBuildCmd())
	rootCmd.AddCommand(NewAddCmd())
	rootCmd.AddCommand(NewInstallCmd())
	rootCmd.AddCommand(UninstallCmd())
	rootCmd.AddCommand(ListCmd())
}
