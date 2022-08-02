package cmd

import (
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/source"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(sourcesCmd)
}

var sourcesCmd = &cobra.Command{
	Use:     "sources",
	Short:   "List available sources",
	Example: "mangal sources",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("Builtin sources:")
		for name := range provider.DefaultProviders() {
			cmd.Println("  " + name)
		}

		cmd.Println()

		cmd.Println("Custom sources:")
		for name := range lo.Must(source.AvailableCustomSources()) {
			cmd.Println("  " + name)
		}
	},
}
