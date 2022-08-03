package cmd

import (
	"github.com/metafates/mangal/constants"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/util"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	rootCmd.AddCommand(clearCmd)
	clearCmd.Flags().BoolP("temp", "t", false, "Clear temporary files")
	clearCmd.Flags().BoolP("cache", "c", false, "Clear cache files")
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clears all useless files (temp, cache)",
	Run: func(cmd *cobra.Command, args []string) {
		var counter uint

		clearTemp := lo.Must(cmd.Flags().GetBool("temp"))
		clearCache := lo.Must(cmd.Flags().GetBool("cache"))

		if clearCache {
			tempDir := os.TempDir()

			err := filesystem.Get().Walk(tempDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}

				if strings.HasPrefix(info.Name(), constants.TempPrefix) {
					counter++

					if info.IsDir() {
						return filesystem.Get().RemoveAll(path)
					} else {
						return filesystem.Get().Remove(path)
					}
				}

				return nil
			})

			if err != nil {
				cmd.PrintErr(err)
				os.Exit(1)
			}
		}

		if clearTemp {
			cacheDir := lo.Must(os.UserCacheDir())
			cacheDir = filepath.Join(cacheDir, constants.CachePrefix)
			err := filesystem.Get().RemoveAll(cacheDir)

			if err != nil {
				cmd.PrintErr(err)
				os.Exit(1)
			}
		}

		cmd.Printf("%s removed\n", util.Quantity(int(counter), "file"))
	},
}
