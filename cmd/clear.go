package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
	"github.com/unknowns24/mks/manager"
)

func ClearCacheCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "clear [" +
			config.ARG_CLEAR_CACHE_ALL + "|" +
			config.ARG_CLEAR_CACHE_FILES + "|" +
			config.ARG_CLEAR_CACHE_ZIP + "|" +
			config.ARG_CLEAR_CACHE_TEMP + "]",
		Short: "Clear mks cache directory. Use to clear specific cache files, by default clear uses " + config.ARG_CLEAR_CACHE_DEFAULT,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 && len(args) == 1 {
				if args[0] == config.ARG_CLEAR_CACHE_FILES {
					manager.ClearCacheFiles()
				} else if args[0] == config.ARG_CLEAR_CACHE_ZIP {
					manager.ClearCacheZip()
				} else if args[0] == config.ARG_CLEAR_CACHE_TEMP {
					manager.ClearCacheTemporals()
				} else if args[0] == config.ARG_CLEAR_CACHE_ALL {
					manager.ClearCacheAll()
				} else {
					cmd.Help()
				}
			} else if len(args) == 0 {
				manager.ClearCacheAll()
			} else {
				cmd.Help()
			}

			return nil
		},
		SilenceUsage: true, // Suppress printing the usage message
	}

	cmd.Flags().BoolVarP(&global.Verbose, config.FLAG_VERBOSE_LONG, config.FLAG_VERBOSE_SHORT, false, "Enable verbose mode")

	return cmd
}
