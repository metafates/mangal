package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:        "install",
	Short:      "Browse and install custom scrapers",
	Deprecated: "use `mangal sources install` instead.",
	Long: `Browse and install custom scrapers from official GitHub repo.
https://github.com/metafates/mangal-scrapers`,
	Run: func(cmd *cobra.Command, args []string) {
		sourcesInstallCmd.Run(sourcesGenCmd, args)
	},
}
