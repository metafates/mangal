package cmd

import (
	"github.com/mangalorg/mangal/nametemplate/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(templateCmd)
}

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Show available name template functions",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		for k, v := range util.Funcs {
			cmd.Println(k)
			cmd.Println(v.Description)
			cmd.Println()
		}
	},
}
