package cmd

import (
	"github.com/mangalorg/libmangal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(formatsCmd)
}

var formatsCmd = &cobra.Command{
	Use:   "formats",
	Short: "Show available formats",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		for _, format := range libmangal.FormatStrings() {
			cmd.Println(format)
		}
	},
}
