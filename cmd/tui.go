package cmd

import (
	"github.com/mangalorg/mangal/provider/manager"
	"github.com/mangalorg/mangal/tui"
	"github.com/mangalorg/mangal/tui/state/providers"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Run mangal in TUI mode",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		loaders, err := manager.Loaders()
		if err != nil {
			errorf(cmd, err.Error())
		}

		if err := tui.Run(providers.New(loaders)); err != nil {
			errorf(cmd, err.Error())
		}
	},
}
