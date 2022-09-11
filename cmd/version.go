package cmd

import (
	"github.com/metafates/mangal/constant"
	"github.com/spf13/cobra"
	"os"
	"runtime"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of mangal",
	Long:  `All software has versions. This is mangal's`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.SetOut(os.Stdout)
		cmd.Printf(`Mangal - The ultimate manga downloader

Version:    %s
OS:         %s
Arch:       %s
Built:      %s by %s
Git Commit: %s
`, constant.Version, runtime.GOOS, runtime.GOARCH, constant.BuiltAt, constant.BuiltBy, constant.GitCommit)
	},
}
