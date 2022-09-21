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
	sourcesCmd.SetOut(os.Stdout)
}

var sourcesCmd = &cobra.Command{
	Use:     "sources",
	Short:   "List an available sources",
	Example: "mangal sources",
	Run: func(cmd *cobra.Command, args []string) {

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
		for _, p := range defaultProviders {
			cmd.Println(p.Name)
		}

		if len(customProviders) == 0 {
			return
		}

		h("")

		h("Custom:")
		for _, p := range provider.CustomProviders() {
			cmd.Println(p.Name)
		}
	},
}
