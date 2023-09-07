package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/unknowns24/mks/global"
	"github.com/unknowns24/mks/manager"
)

func NewAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [feature]",
		Short: "Add a feature to a microservice",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			feature := args[0]

			validFeature := manager.IsValidFeature(feature)

			// Validate feature argument
			if !validFeature {
				return fmt.Errorf("unknown feature '%s'. Valid features are: %s", feature, global.InstalledTemplates[:])
			}

			// Main function start
			return manager.AddFeature(feature)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return global.InstalledTemplates[:], cobra.ShellCompDirectiveDefault
		},
		SilenceUsage: true, // Suppress printing the usage message
	}

	cmd.Flags().BoolVarP(&global.Verbose, "verbose", "v", false, "Enable verbose mode")

	return cmd
}
