package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mks",
	Short: "Microservice generator CLI",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(NewBuildCmd())
	rootCmd.AddCommand(NewAddCmd())
	rootCmd.AddCommand(NewRemoveCmd())
}
