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
	Short: "Browse and install custom scrapers",
	Long: `Browse and install custom scrapers from official GitHub repo.
https://github.com/metafates/mangal-scrapers`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.Run(&tui.Options{Install: true})
	},
}
