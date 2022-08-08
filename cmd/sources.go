package cmd

import (
	"github.com/metafates/mangal/icon"
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

		cmd.Println(headerStyle("Builtin sources:"))
		for name := range provider.DefaultProviders() {
			name = "  " + name + " " + icon.Get(icon.Go)
			cmd.Println(name)
		}

		cmd.Println()

		cmd.Println(headerStyle("Custom sources:"))
		custom, err := provider.CustomProviders()
		if err != nil {
			return err
		}

		for name := range custom {
			name = "  " + name + " " + icon.Get(icon.Lua)
			cmd.Println(name)
		}

		return nil
	},
}
