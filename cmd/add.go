package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/manager"
)

func NewAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [feature]",
		Short: "Add a feature to a microservice",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			feature := args[0]

			// Validate feature argument
			if !manager.IsValidFeature(feature) {
				return fmt.Errorf("unknown feature '%s'. Valid features are: %s", feature, config.Features[:])
			}

			return manager.AddFeature(feature)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return config.Features[:], cobra.ShellCompDirectiveDefault
		},
		SilenceUsage: true, // Suppress printing the usage message
	}

	cmd.Flags().BoolVarP(&manager.Verbose, "verbose", "v", false, "Enable verbose mode")

	return cmd
}
