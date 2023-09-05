package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unknowns24/mks/global"
	"github.com/unknowns24/mks/manager"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List of templates installed and available to use by mks",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return manager.ListTemplate()
		},
		SilenceUsage: true, // Suppress printing the usage message
	}

	cmd.Flags().BoolVarP(&global.Verbose, "verbose", "v", false, "Enable verbose mode")

	return cmd
}
