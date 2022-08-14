package cmd

import (
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/converter"
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

	miniCmd.Flags().StringP("format", "f", "", "output format")
	lo.Must0(miniCmd.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return converter.Available(), cobra.ShellCompDirectiveDefault
	}))
	lo.Must0(viper.BindPFlag(config.FormatsUse, miniCmd.Flags().Lookup("format")))
}

var miniCmd = &cobra.Command{
	Use:   "mini",
	Short: "Launch the in mini mode",
	Long: `Launch mangal the in mini mode.
Will try to mimic ani-cli.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		options := mini.Options{
			Download: lo.Must(cmd.Flags().GetBool("download")),
			Continue: lo.Must(cmd.Flags().GetBool("continue")),
		}
		err := mini.Run(&options)

		if err != nil {
			if err.Error() == "interrupt" {
				return nil
			}

			return err
		}

		return nil
	},
}
