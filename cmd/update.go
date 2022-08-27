package cmd

import (
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/icon"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/updater"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update mangal",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("%s %s\n", icon.Get(icon.Progress), "Fetching latest version...")

		latestVersion, err := updater.LatestVersion()
		handleErr(err)

		if constant.Version >= latestVersion {
			cmd.Printf("%s %s\n", icon.Get(icon.Success), style.Green("You are using the latest version"))
			return
		} else {
			cmd.Printf("%s New version is available: %s\n", icon.Get(icon.Success), latestVersion)
		}

		err = updater.Update()
		handleErr(err)
	},
}
