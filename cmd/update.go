package cmd

import (
	"fmt"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/icon"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/updater"
	"github.com/metafates/mangal/util"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update mangal",
	Aliases: []string{"upgrade"},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.SetOut(os.Stdout)

		msg := fmt.Sprintf("%s %s", icon.Get(icon.Progress), "Fetching latest version...")
		erase := util.PrintErasable(msg)

		latestVersion, err := updater.LatestVersion()
		handleErr(err)

		if constant.Version >= latestVersion {
			erase()
			msg := fmt.Sprintf(
				"%s %s %s\n",
				icon.Get(icon.Success),
				style.Green("You are using the latest version"),
				style.Faint("(which is "+constant.Version+")"),
			)
			cmd.Printf(msg)
			return
		} else {
			erase()
			cmd.Printf("%s New version is available: %s\n", icon.Get(icon.Success), latestVersion)
		}

		err = updater.Update()
		handleErr(err)
	},
}
