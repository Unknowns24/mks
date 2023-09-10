package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
	"github.com/unknowns24/mks/manager"
)

func BuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build [name]",
		Short: "Create an application with custom features",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ApplicationName := args[0]
			features, _ := cmd.Flags().GetStringSlice(config.FLAG_FEATURE_LONG)

			// Validate every feature on features string slice
			for _, feature := range features {
				// prevent all features invalid error
				if feature == config.ALL_FEATURES {
					continue
				}

				validFeature := manager.IsValidFeature(feature)

				if !validFeature {
					return fmt.Errorf("unknown feature '%s'. Valid features are: %s", feature, append([]string{config.ALL_FEATURES}, global.InstalledTemplates[:]...))
				}
			}

			// Main function start
			return manager.GenerateApplication(ApplicationName, features)
		},
	}

	cmd.Flags().BoolVarP(&global.Verbose, config.FLAG_VERBOSE_LONG, config.FLAG_VERBOSE_SHORT, false, "Enable verbose mode")
	cmd.Flags().StringSliceP(config.FLAG_FEATURE_LONG, config.FLAG_FEATURE_SHORT, []string{}, "Features required for the application")

	return cmd
}
