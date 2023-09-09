package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
	"github.com/unknowns24/mks/manager"
)

func InfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Show information about mks application and used paths",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			//Show program information to user
			fmt.Println("MKS - Golang application manager CLI")
			fmt.Println("")
			fmt.Println("Version: \n    ", config.MKS_Info_Version)
			fmt.Println("Authors: \n    ", config.MKS_Info_Author)
			fmt.Println("License: \n    ", config.MKS_Info_License)
			fmt.Println("Repository: \n    ", config.MKS_Info_Repository)
			fmt.Println("")
			fmt.Printf("User data directory: \n    %s\n", global.ConfigFolderPath)
			fmt.Printf("User data temporals: \n    %s\n", global.TemporalsPath)
			fmt.Printf("User data templates zip cache: \n    %s\n", global.ZipCachePath)
			fmt.Printf("User data templates file cache: \n    %s\n", global.TemplateCachePath)
			fmt.Printf("User data templates installed: \n    %s\n", global.UserTemplatesFolderPath)
			fmt.Println("")

			manager.ListTemplate()
			return nil
		},
		SilenceUsage: true, // Suppress printing the usage message
	}

	return cmd
}
