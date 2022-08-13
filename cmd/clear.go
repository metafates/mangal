package cmd

import (
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/where"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	rootCmd.AddCommand(clearCmd)
	clearCmd.Flags().BoolP("cache", "c", false, "Clear cache files")
	clearCmd.Flags().BoolP("history-file", "r", false, "Clear history")
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clears a sidelined files",
	Run: func(cmd *cobra.Command, args []string) {
		doClearCache := lo.Must(cmd.Flags().GetBool("cache"))
		doClearHistory := lo.Must(cmd.Flags().GetBool("history"))

		if doClearCache {
			clearCache()
		}

		if doClearHistory {
			clearHistory()
		}

		cmd.Println("Cleared")
	},
}

func clearCache() {
	cacheDir := lo.Must(os.UserCacheDir())
	cacheDir = filepath.Join(cacheDir, constant.CachePrefix)
	_ = filesystem.Get().RemoveAll(cacheDir)
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
	_ = filesystem.Get().Remove(historyFile)
}
