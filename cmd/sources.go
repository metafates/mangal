package cmd

import (
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/style"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(sourcesCmd)
}

var sourcesCmd = &cobra.Command{
	Use:     "sources",
	Short:   "List an available sources",
	Example: "mangal sources",
	RunE: func(cmd *cobra.Command, args []string) error {
		headerStyle := style.Combined(style.Bold, style.HiBlue)

		cmd.Println(headerStyle("Builtin:"))
		for name := range provider.DefaultProviders() {
			cmd.Println(name)
		}

		cmd.Println()

		cmd.Println(headerStyle("Custom:"))
		custom := provider.CustomProviders()

		for name := range custom {
			cmd.Println(name)
		}

		return nil
	},
}
