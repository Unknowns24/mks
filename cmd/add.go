package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unknowns24/mks/manager"
)

func NewAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [feature]",
		Short: "Add a feature to a microservice",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			feature := args[1]
			return manager.AddFeature(feature)
		},
	}

	cmd.Flags().BoolVarP(&manager.Verbose, "verbose", "v", false, "Enable verbose mode")

	return cmd
}
