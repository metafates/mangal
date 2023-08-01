package cmd

import (
	"log"

	"github.com/mangalorg/mangal/meta"
	"github.com/mangalorg/mangal/provider/manager"
	"github.com/mangalorg/mangal/tui"
	"github.com/mangalorg/mangal/tui/state/providers"
	"github.com/spf13/cobra"
)

var rootArgs = struct {
	Version bool
}{}

func init() {
	rootCmd.Flags().BoolVarP(&rootArgs.Version, "version", "v", false, "show version information")
}

var rootCmd = &cobra.Command{
	Use:   meta.AppName,
	Short: "The ultimate CLI manga downloader",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if rootArgs.Version {
			versionCmd.Run(versionCmd, nil)
			return nil
		}

		loaders, err := manager.Loaders()
		if err != nil {
			return err
		}

		return tui.Run(providers.New(loaders))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
