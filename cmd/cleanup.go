package cmd

import (
	"fmt"
	"github.com/metafates/mangal/cleaner"
	"github.com/metafates/mangal/util"
	"github.com/spf13/cobra"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Remove files created by mangal",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var cleanupTempCmd = &cobra.Command{
	Use:   "temp",
	Short: "Remove temp files",
	Long:  "Removes temp files produced by downloader",
	Run: func(cmd *cobra.Command, args []string) {
		counter, bytes := cleaner.RemoveTemp()
		fmt.Printf("%d temp files removed\nCleaned up %.2fMB\n", counter, util.BytesToMegabytes(bytes))
	},
}

var cleanupCacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Remove cache files",
	Long:  "Removes cache files produced by scraper",
	Run: func(cmd *cobra.Command, args []string) {
		counter, bytes := cleaner.RemoveCache()
		fmt.Printf("%d cache files removed\nCleaned up %.2fMB\n", counter, util.BytesToMegabytes(bytes))
	},
}

var cleanupHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Clear history",
	Long:  "Removes history files produced by reader",
	Run: func(cmd *cobra.Command, args []string) {
		_, bytes := cleaner.RemoveHistory()
		fmt.Printf("History file removed\nCleaned up %.2fMB\n", util.BytesToMegabytes(bytes))
	},
}

var cleanupAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Remove history, cache and temp files",
	Long:  "Removes history files produced by reader, cache files produced by scraper and temp files produced by downloader",
	Run: func(cmd *cobra.Command, args []string) {
		counter, bytes := cleaner.RemoveTemp()
		c, b := cleaner.RemoveCache()
		counter += c
		bytes += b
		c, b = cleaner.RemoveHistory()
		counter += c
		bytes += b
		fmt.Printf("%d files removed\nCleaned up %.2fMB\n", counter, util.BytesToMegabytes(bytes))
	},
}

func init() {
	cleanupCmd.AddCommand(cleanupTempCmd)
	cleanupCmd.AddCommand(cleanupCacheCmd)
	cleanupCmd.AddCommand(cleanupHistoryCmd)
	cleanupCmd.AddCommand(cleanupAllCmd)

	mangalCmd.AddCommand(cleanupCmd)
}
