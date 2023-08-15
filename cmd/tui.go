package cmd

import (
	"github.com/mangalorg/mangal/provider/manager"
	"github.com/mangalorg/mangal/tui"
	"github.com/mangalorg/mangal/tui/state/providers"
	"github.com/spf13/cobra"
)

func init() {
	subcommands = append(subcommands, tuiCmd)
}

var tuiCmd = &cobra.Command{
	Use:     "tui",
	Short:   "Run mangal in TUI mode",
	GroupID: groupMode,
	Args:    cobra.NoArgs,
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
