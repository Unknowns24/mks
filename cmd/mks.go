package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mks",
	Short: "Golang application manager CLI",
}

func Execute() {
	// Start main command
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(BuildCmd())
	rootCmd.AddCommand(AddCmd())
	rootCmd.AddCommand(InstallCmd())
	rootCmd.AddCommand(UninstallCmd())
	rootCmd.AddCommand(ListCmd())
	rootCmd.AddCommand(ClearCacheCmd())
	rootCmd.AddCommand(InfoCmd())
	rootCmd.AddCommand(ExportCmd())
}
