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
	handleErr(filesystem.Get().RemoveAll(where.Cache()))
	//
	//cacheDir, err := os.UserCacheDir()
	//handleErr(err)
	//
	//err = filesystem.Get().Walk(cacheDir, func(path string, info fs.FileInfo, err error) error {
	//	if err != nil {
	//		return nil
	//	}
	//
	//	if strings.HasPrefix(info.Name(), constant.CachePrefix) {
	//		if info.IsDir() {
	//			return filesystem.Get().RemoveAll(path)
	//		} else {
	//			return filesystem.Get().Remove(path)
	//		}
	//	}
	//
	//	return nil
	//})
	//
	//handleErr(err)
}

func clearTemp() {
	handleErr(filesystem.Get().RemoveAll(where.Temp()))
	//
	//tempDir := os.TempDir()
	//
	//err := filesystem.Get().Walk(tempDir, func(path string, info os.FileInfo, err error) error {
	//	if err != nil {
	//		return nil
	//	}
	//
	//	if strings.HasPrefix(info.Name(), constant.TempPrefix) {
	//		if info.IsDir() {
	//			return filesystem.Get().RemoveAll(path)
	//		} else {
	//			return filesystem.Get().Remove(path)
	//		}
	//	}
	//
	//	return nil
	//})
	//
	//handleErr(err)
}

func clearHistory() {
	historyFile := where.History()
	handleErr(filesystem.Get().Remove(historyFile))
}
