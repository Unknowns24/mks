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
	rootCmd.AddCommand(NewBuildCmd())
	rootCmd.AddCommand(NewAddCmd())
	rootCmd.AddCommand(NewInstallCmd())
	rootCmd.AddCommand(UninstallCmd())
	rootCmd.AddCommand(ListCmd())
}
