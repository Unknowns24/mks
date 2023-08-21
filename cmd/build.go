package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unknowns24/mks/libs/generator"
)

func NewBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build [name]",
		Short: "Create a microservice with custom fetures",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			full, _ := cmd.Flags().GetBool("full")
			return generator.GenerateMicroservice(serviceName, full)
		},
	}

	cmd.Flags().Bool("full", false, "Generate microservice with all features")

	return cmd
}
