package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/manager"
)

func NewBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build [name]",
		Short: "Create a microservice with custom features",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			features, _ := cmd.Flags().GetStringSlice("features")

			// Validate every feature on features string slice
			for _, feature := range features {
				if !manager.IsValidFeature(feature) {
					return fmt.Errorf("unknown feature '%s'. Valid features are: %s", feature, append([]string{config.ALL_FEATURES}, config.Features[:]...))
				}
			}

			return manager.GenerateMicroservice(serviceName, features)
		},
	}

	cmd.Flags().BoolVarP(&manager.Verbose, "verbose", "v", false, "Enable verbose mode")
	cmd.Flags().StringSlice("features", []string{}, "Features required for the microservice")

	return cmd
}
