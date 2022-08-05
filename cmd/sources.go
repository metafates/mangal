package cmd

import (
	"github.com/metafates/mangal/icon"
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/style"
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
		headerStyle := style.Combined(style.Bold, style.HiBlue)

		cmd.Println(headerStyle("Builtin sources:"))
		for name := range provider.DefaultProviders() {
			name = "  " + name + " " + icon.Get(icon.Go)
			cmd.Println(name)
		}

		cmd.Println()

		cmd.Println(headerStyle("Custom sources:"))
		for name := range lo.Must(source.AvailableCustomSources()) {
			name = "  " + name + " " + icon.Get(icon.Lua)
			cmd.Println(name)
		}
	},
}
