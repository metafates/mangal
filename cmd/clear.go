package cmd

import (
	"fmt"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/icon"
	"github.com/metafates/mangal/log"
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
		doClearCache := lo.Must(cmd.Flags().GetBool("cache"))
		doClearHistory := lo.Must(cmd.Flags().GetBool("history"))

		if doClearCache {
			anyCleared = true
			e := util.PrintErasable(fmt.Sprintf("%s Clearing cache...", icon.Get(icon.Progress)))
			clearCache()
			e()
			fmt.Printf("%s Cache cleared\n", icon.Get(icon.Success))
		}

		if doClearHistory {
			anyCleared = true
			e := util.PrintErasable(fmt.Sprintf("%s Clearing history...", icon.Get(icon.Progress)))
			clearHistory()
			e()
			fmt.Printf("%s History cleared\n", icon.Get(icon.Success))
		}

		if !anyCleared {
			handleErr(cmd.Help())
		}
	},
}

func clearCache() {
	log.Infof("Clearing cache at %s", where.Cache())
	handleErr(filesystem.Api().RemoveAll(where.Cache()))
}

func clearTemp() {
	log.Infof("Clearing temp files at %s", where.Temp())
	handleErr(filesystem.Api().RemoveAll(where.Temp()))
}

func clearHistory() {
	historyFile := where.History()
	log.Infof("Removing history file at %s", historyFile)
	handleErr(filesystem.Api().Remove(historyFile))
}
