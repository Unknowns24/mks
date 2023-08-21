package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unknowns24/mks/libs/generator"
)

func NewAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [name] [feature]",
		Short: "Add a feature to a microservice",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			feature := args[1]
			return generator.AddFeature(serviceName, feature)
		},
	}

	return cmd
}
