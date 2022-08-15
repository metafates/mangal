package cmd

import (
	"github.com/metafates/mangal/tui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install custom scrapers",
	Long:  `Install custom scrapers from GitHub repo.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.Run(&tui.Options{Install: true})
	},
}
