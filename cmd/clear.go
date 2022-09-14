package cmd

import (
	"fmt"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/icon"
	"github.com/metafates/mangal/util"
	"github.com/metafates/mangal/where"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(clearCmd)
	clearCmd.Flags().Bool("cache", false, "Clear cache files")
	clearCmd.Flags().Bool("history", false, "Clear history")
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clears a sidelined files",
	Run: func(cmd *cobra.Command, args []string) {
		var anyCleared bool

		doClear := func(what string) bool {
			return lo.Must(cmd.Flags().GetBool(what))
		}

		for name, clear := range map[string]func(){
			"cache":   clearCache,
			"history": clearHistory,
		} {
			if doClear(name) {
				anyCleared = true
				e := util.PrintErasable(fmt.Sprintf("%s Clearing %s...", icon.Get(icon.Progress), util.Capitalize(name)))
				clear()
				e()
				fmt.Printf("%s %s cleared\n", icon.Get(icon.Success), util.Capitalize(name))
			}
		}

		if !anyCleared {
			handleErr(cmd.Help())
		}
	},
}

func clearCache() {
	path := where.Cache()
	handleErr(filesystem.Api().RemoveAll(path))
}

func clearTemp() {
	path := where.Temp()
	handleErr(filesystem.Api().RemoveAll(path))
}

func clearHistory() {
	path := where.History()
	handleErr(filesystem.Api().Remove(path))
}
