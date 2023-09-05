package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unknowns24/mks/global"
	"github.com/unknowns24/mks/manager"
)

func NewInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install [template]",
		Short: "Install a template to mks",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			template := args[0]
			return manager.InstallTemplate(template)
		},
		SilenceUsage: true, // Suppress printing the usage message
	}

	cmd.Flags().BoolVarP(&global.Verbose, "verbose", "v", false, "Enable verbose mode")

	return cmd
}
