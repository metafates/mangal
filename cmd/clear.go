package cmd

import (
	"fmt"
	"github.com/metafates/mangal/constants"
	"github.com/metafates/mangal/filesystem"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	rootCmd.AddCommand(clearCmd)
	clearCmd.Flags().BoolP("cache", "c", false, "Clear cache files")
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clears sidelined files",
	Run: func(cmd *cobra.Command, args []string) {
		doClearCache := lo.Must(cmd.Flags().GetBool("cache"))

		if doClearCache {
			clearCache()
		}

		cmd.Println("Cleared")
	},
}

func clearCache() {
	cacheDir := lo.Must(os.UserCacheDir())
	cacheDir = filepath.Join(cacheDir, constants.CachePrefix)
	err := filesystem.Get().RemoveAll(cacheDir)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func clearTemp() {
	tempDir := os.TempDir()

	err := filesystem.Get().Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if strings.HasPrefix(info.Name(), constants.TempPrefix) {
			if info.IsDir() {
				return filesystem.Get().RemoveAll(path)
			} else {
				return filesystem.Get().Remove(path)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
