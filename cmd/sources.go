package cmd

import (
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/style"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(sourcesCmd)

	sourcesCmd.Flags().BoolP("raw", "r", false, "do not print headers")
}

var sourcesCmd = &cobra.Command{
	Use:     "sources",
	Short:   "List an available sources",
	Example: "mangal sources",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.SetOut(os.Stdout)

		printHeader := !lo.Must(cmd.Flags().GetBool("raw"))
		headerStyle := style.Combined(style.Bold, style.HiBlue)

		h := func(s string) {
			if printHeader {
				cmd.Println(headerStyle(s))
			}
		}

		defaultProviders := provider.DefaultProviders()
		customProviders := provider.CustomProviders()

		h("Builtin:")
		for name := range defaultProviders {
			cmd.Println(name)
		}

		if len(customProviders) == 0 {
			return
		}

		h("")

		h("Custom:")
		for name := range provider.CustomProviders() {
			cmd.Println(name)
		}
	},
}
