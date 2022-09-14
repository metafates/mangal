package cmd

import (
	"github.com/metafates/mangal/updater"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.SetOut(os.Stdout)
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update mangal",
	Aliases: []string{"upgrade"},
	Run: func(cmd *cobra.Command, args []string) {
		handleErr(updater.Update())
	},
}
