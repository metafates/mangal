package main

import (
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   strings.ToLower(AppName),
	Short: AppName + " - Manga Downloader",
	Long:  `A fast and flexible manga downloader`,
	Run: func(cmd *cobra.Command, args []string) {
		config, _ := cmd.Flags().GetString("config")
		exists, err := Afero.Exists(config)

		if err != nil {
			log.Fatal(errors.New("access to config file denied"))
		}

		if config != "" {
			config = path.Clean(config)
			if !exists {
				log.Fatal(errors.New(fmt.Sprintf("config at path %s doesn't exist", config)))
			}

			UserConfig = GetConfig(config)
		} else {
			UserConfig = GetConfig("") // get config from default config path
		}

		var program *tea.Program

		if UserConfig.Fullscreen {
			program = tea.NewProgram(newBubble(searchState), tea.WithAltScreen())
		} else {
			program = tea.NewProgram(newBubble(searchState))
		}

		if err := program.Start(); err != nil {
			log.Fatal(err)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Long:  fmt.Sprintf("Shows %s versions and build date", AppName),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s version %s\nBuild %s\n", AppName, version, build)
	},
}

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Remove cached and temp files",
	Long:  "Removes cached files produced by scraper and temp files from downloader",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			// counter of removed files
			counter int
			// bytes removed
			bytes int64
		)

		// Cleanup temp files
		tempDir := os.TempDir()
		tempFiles, err := Afero.ReadDir(tempDir)
		if err == nil {
			lowerAppName := strings.ToLower(AppName)
			for _, tempFile := range tempFiles {
				name := tempFile.Name()
				if strings.HasPrefix(name, AppName) || strings.HasPrefix(name, lowerAppName) {
					err = Afero.Remove(filepath.Join(tempDir, name))
					if err == nil {
						bytes += tempFile.Size()
						counter++
					}
				}
			}
		}

		// Cleanup cache files
		cacheDir, err := os.UserCacheDir()
		if err == nil {
			scraperCacheDir := filepath.Join(cacheDir, AppName)
			if exists, err := Afero.Exists(scraperCacheDir); err == nil && exists {
				files, err := Afero.ReadDir(scraperCacheDir)
				if err == nil {
					counter += len(files)
					for _, f := range files {
						bytes += f.Size()
					}
				}

				_ = Afero.RemoveAll(scraperCacheDir)
			}
		}

		fmt.Printf("\U0001F9F9 %d files removed. Cleaned up %.2fMB\n", counter, BytesToMegabytes(bytes))
	},
}

func CmdExecute() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(cleanupCmd)
	rootCmd.PersistentFlags().StringP("config", "c", "", "use config from path")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
