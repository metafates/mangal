package cmd

import (
	"github.com/metafates/mangal/mini"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(miniCmd)
}

var miniCmd = &cobra.Command{
	Use:   "mini",
	Short: "Launch in mini mode",
	Long:  `Launch in mini mode. Will use simple prompts instead of TUI.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := mini.Run()

		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}
	},
}
