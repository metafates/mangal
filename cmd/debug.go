package cmd

import (
	"github.com/metafates/mangal/source"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Run source in debug mode",
	Run: func(cmd *cobra.Command, args []string) {
		sourcePath := lo.Must(cmd.Flags().GetString("file"))
		source, err := source.LoadSource(sourcePath)
		if err != nil {
			cmd.PrintErr(err)
		}

		err = source.Debug()
		if err != nil {
			cmd.PrintErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)

	debugCmd.Flags().StringP("file", "f", "", "source file to run")
	_ = debugCmd.MarkFlagRequired("file")
	_ = debugCmd.MarkFlagFilename("file", "lua")
}
