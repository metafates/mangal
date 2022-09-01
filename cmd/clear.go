package cmd

import (
	"fmt"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/icon"
	"github.com/metafates/mangal/util"
	"github.com/metafates/mangal/where"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
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
		var any bool
		doClearCache := lo.Must(cmd.Flags().GetBool("cache"))
		doClearHistory := lo.Must(cmd.Flags().GetBool("history"))

		if doClearCache {
			any = true
			e := util.PrintErasable(fmt.Sprintf("%s Clearing cache...", icon.Get(icon.Progress)))
			clearCache()
			e()
			fmt.Printf("%s Cache cleared\n", icon.Get(icon.Success))
		}

		if doClearHistory {
			any = true
			e := util.PrintErasable(fmt.Sprintf("%s Clearing history...", icon.Get(icon.Progress)))
			clearHistory()
			e()
			fmt.Printf("%s History cleared\n", icon.Get(icon.Success))
		}

		if !any {
			handleErr(cmd.Help())
		}
	},
}

func clearCache() {
	cacheDir := lo.Must(os.UserCacheDir())
	cacheDir = filepath.Join(cacheDir, constant.CachePrefix)
	handleErr(filesystem.Get().RemoveAll(cacheDir))
}

func clearTemp() {
	tempDir := os.TempDir()

	_ = filesystem.Get().Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if strings.HasPrefix(info.Name(), constant.TempPrefix) {
			if info.IsDir() {
				return filesystem.Get().RemoveAll(path)
			} else {
				return filesystem.Get().Remove(path)
			}
		}

		return nil
	})
}

func clearHistory() {
	historyFile := where.History()
	handleErr(filesystem.Get().Remove(historyFile))
}
