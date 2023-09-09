package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
	"github.com/unknowns24/mks/manager"
)

func UninstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall [template]",
		Short: "Uninstall a template from mks",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			template := args[0]
			return manager.UninstallTemplate(template)
		},
		SilenceUsage: true, // Suppress printing the usage message
	}

	cmd.Flags().BoolVarP(&global.Verbose, config.FLAG_VERBOSE_LONG, config.FLAG_VERBOSE_SHORT, false, "Enable verbose mode")

	return cmd
}
