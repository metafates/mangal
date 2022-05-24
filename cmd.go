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
	Short: AppName + " is a manga downloader",
	Long: `A fast and flexible manga downloader built with
love by metafates.`,
	Run: func(cmd *cobra.Command, args []string) {
		showVersion, _ := cmd.Flags().GetBool("version")

		if showVersion {
			fmt.Printf("%s version %s\nBuild %s", AppName, version, build)
			os.Exit(0)
		}

		config, _ := cmd.Flags().GetString("config")
		exists, err := Afero.Exists(config)

		if err != nil {
			log.Fatal(errors.New("can't check if file exists or not"))
		}

		if config != "" {
			config = path.Clean(config)
			if !exists {
				log.Fatal(errors.New(fmt.Sprintf("file at path %s doesn't exist", config)))
			}

			UserConfig = GetConfig(config)
		} else {
			// TODO: replace it with real config
			//UserConfig = GetConfig("")
			UserConfig = DefaultConfig
		}

		program := tea.NewProgram(newBubble(searchState), IfElse[tea.ProgramOption](UserConfig.Fullscreen, tea.WithAltScreen(), nil))

		if err := program.Start(); err != nil {
			log.Fatal(err)
		}
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
	rootCmd.AddCommand(cleanupCmd)
	rootCmd.PersistentFlags().StringP("config", "c", "", "path to config file")
	rootCmd.PersistentFlags().BoolP("version", "v", false, "show version")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
