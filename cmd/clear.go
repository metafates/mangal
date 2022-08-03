package cmd

import (
	"github.com/metafates/mangal/constants"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/util"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func init() {
	rootCmd.AddCommand(clearCmd)
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all temp files",
	Run: func(cmd *cobra.Command, args []string) {
		var counter uint

		tempDir := os.TempDir()

		err := filesystem.Get().Walk(tempDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if strings.HasPrefix(info.Name(), constants.TempPrefix) {
				counter++
				exists, err := filesystem.Get().Exists(path)

				if !exists || err != nil {
					return nil
				}

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

		cmd.Printf("%s removed\n", util.Quantity(int(counter), "file"))
	},
}
