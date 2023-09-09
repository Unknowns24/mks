package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/unknowns24/mks/global"
	"github.com/unknowns24/mks/manager"
)

func InfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Show information abour mks and templates",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			//Show program information to user
			fmt.Println("[+] MKS - Golang application manager CLI")
			fmt.Println(" ├──── Version: ", global.MKS_Info_Version)
			fmt.Println(" ├──── Authors: ", global.MKS_Info_Author)
			fmt.Println(" ├──── License: ", global.MKS_Info_License)
			fmt.Println(" └──── Repository: ", global.MKS_Info_Repository)
			fmt.Println("")
			fmt.Println("[+] Routes:")
			fmt.Println(" └─┬── User data directory: ", global.ConfigFolderPath)
			fmt.Println("   ├── User data temporals: ", global.TemporalsPath)
			fmt.Println("   ├── User data templates zip cache: ", global.ZipCachePath)
			fmt.Println("   ├── User data templates file cache: ", global.TemplateCachePath)
			fmt.Println("   └── User data templates installed: ", global.UserTemplatesFolderPath)
			fmt.Println("")
			manager.ListTemplate()
			return nil
		},
		SilenceUsage: true, // Suppress printing the usage message
	}

	return cmd
}
