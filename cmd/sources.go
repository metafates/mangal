package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/icon"
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/where"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
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

func init() {
	sourcesCmd.AddCommand(sourcesRemoveCmd)
}

var sourcesRemoveCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove a custom source",
	Example: "mangal sources remove <name>",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		for _, name := range args {
			path := filepath.Join(where.Sources(), name+provider.CustomProviderExtension)
			handleErr(filesystem.Api().Remove(path))
			fmt.Printf("%s successfully removed %s\n", icon.Get(icon.Success), style.Yellow(name))
		}
	},
}
