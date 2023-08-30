package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unknowns24/mks/manager"
)

var (
	verbose bool // verbose flag
)

func NewBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build [name]",
		Short: "Create a microservice with custom features",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			features, _ := cmd.Flags().GetStringSlice("features")
			return manager.GenerateMicroservice(serviceName, verbose, features)
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose mode")
	cmd.Flags().StringSlice("features", []string{}, "Features required for the microservice")

	return cmd
}
