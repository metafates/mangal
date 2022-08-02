package cmd

import (
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/source"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"path/filepath"
)

var sourcesCmd = &cobra.Command{
	Use:   "sources",
	Short: "List available sources",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("Builtin sources:")
		for _, source := range provider.Default() {
			cmd.Println("  " + source.Name)
		}

		cmd.Println()

		cmd.Println("Custom sources:")
		for _, source := range lo.Must(source.AvailableCustomSources()) {
			base := filepath.Base(source)
			cmd.Println("  " + base[:len(base)-4])
		}
	},
}

func init() {
	rootCmd.AddCommand(sourcesCmd)
}
