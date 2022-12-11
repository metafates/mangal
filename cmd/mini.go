package cmd

import (
	"github.com/metafates/mangal/converter"
	"github.com/metafates/mangal/key"
	"github.com/metafates/mangal/mini"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(miniCmd)

	miniCmd.Flags().BoolP("download", "d", false, "download mode")
	miniCmd.Flags().BoolP("continue", "c", false, "continue reading")

	miniCmd.MarkFlagsMutuallyExclusive("download", "continue")
}

var miniCmd = &cobra.Command{
	Use:   "mini",
	Short: "Launch in the mini mode",
	Long: `Launch mangal in the mini mode.
Will try to mimic ani-cli.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if _, err := converter.Get(viper.GetString(key.FormatsUse)); err != nil {
			handleErr(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		options := mini.Options{
			Download: lo.Must(cmd.Flags().GetBool("download")),
			Continue: lo.Must(cmd.Flags().GetBool("continue")),
		}
		err := mini.Run(&options)

		if err != nil && err.Error() != "interrupt" {
			handleErr(err)
		}
	},
}
