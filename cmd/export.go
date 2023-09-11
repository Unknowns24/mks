package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
	"github.com/unknowns24/mks/manager"
)

func ExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export a mks template package",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			useFlag, _ := cmd.Flags().GetStringSlice(config.FLAG_USE_LONG)

			return manager.ExportTemplates(useFlag)
		},
		SilenceUsage: true, // Suppress printing the usage message
	}

	cmd.Flags().BoolVarP(&global.Verbose, config.FLAG_VERBOSE_LONG, config.FLAG_VERBOSE_SHORT, false, "Enable verbose mode")
	cmd.Flags().StringSliceP(config.FLAG_USE_LONG, config.FLAG_USE_SHORT, []string{}, "Select one or more templates to export. Example: --use=\"template1,template2,template...\"")

	return cmd
}
