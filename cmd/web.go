package cmd

import (
	"github.com/mangalorg/mangal/web"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(webCmd)
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Run web UI",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := web.Run(); err != nil {
			errorf(cmd, err.Error())
		}
	},
}
