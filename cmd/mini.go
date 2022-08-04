package cmd

import (
	"github.com/metafates/mangal/icon"
	"github.com/metafates/mangal/mini"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func init() {
	rootCmd.AddCommand(miniCmd)

	miniCmd.Flags().BoolP("download", "d", false, "download mode")
	miniCmd.Flags().BoolP("continue", "c", false, "continue reading")

	miniCmd.MarkFlagsMutuallyExclusive("download", "continue")
}

var miniCmd = &cobra.Command{
	Use:   "mini",
	Short: "Launch in mini mode",
	Long:  `Launch in mini mode. Will use simple prompts instead of TUI.`,
	Run: func(cmd *cobra.Command, args []string) {
		options := mini.Options{
			Download: lo.Must(cmd.Flags().GetBool("download")),
			Continue: lo.Must(cmd.Flags().GetBool("continue")),
		}
		err := mini.Run(&options)

		if err != nil {
			if err.Error() == "interrupt" {
				os.Exit(0)
			}

			cmd.PrintErrf("%s %s", icon.Get(icon.Fail), strings.Title(err.Error()))
			cmd.Println()
			os.Exit(1)
		}
	},
}
