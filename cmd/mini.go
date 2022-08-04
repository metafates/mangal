package cmd

import (
	"github.com/metafates/mangal/mini"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(miniCmd)

	miniCmd.Flags().BoolP("download", "d", false, "download mode")
}

var miniCmd = &cobra.Command{
	Use:   "mini",
	Short: "Launch in mini mode",
	Long:  `Launch in mini mode. Will use simple prompts instead of TUI.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := mini.Run(lo.Must(cmd.Flags().GetBool("download")))

		if err != nil {
			if err.Error() == "interrupt" {
				os.Exit(0)
			}

			cmd.PrintErr(err)
			os.Exit(1)
		}
	},
}
