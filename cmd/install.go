package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
	"github.com/unknowns24/mks/manager"
)

func InstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install [template]",
		Short: "Install a template to mks",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			template := args[0]

			useFlag, _ := cmd.Flags().GetStringSlice(config.FLAG_USE_LONG)

			return manager.InstallTemplate(template, useFlag)
		},
		SilenceUsage: true, // Suppress printing the usage message
	}

	cmd.Flags().BoolVarP(&global.Verbose, config.FLAG_VERBOSE_LONG, config.FLAG_VERBOSE_SHORT, false, "Enable verbose mode")
	cmd.Flags().StringSliceP(config.FLAG_USE_LONG, config.FLAG_USE_SHORT, []string{}, "Select one or more templates to use. Example: --use=\"template1,template2,template...\"")

	return cmd
}
