package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unknowns24/mks/manager"
)

func NewRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [name] [feature]",
		Short: "Remove a feature from a microservice",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			feature := args[1]
			return manager.RemoveFeature(serviceName, feature)
		},
	}

	return cmd
}
