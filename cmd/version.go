package cmd

import (
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/updater"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.AddCommand(versionLatestCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of mangal",
	Long:  `All software has versions. This is mangal's`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("mangal version " + constant.Version)
	},
}

var versionLatestCmd = &cobra.Command{
	Use:   "latest",
	Short: "Print the latest version number of the mangal",
	Long:  `It will fetch the latest version from the github and print it`,
	Run: func(cmd *cobra.Command, args []string) {
		version, err := updater.LatestVersion()
		handleErr(err)

		cmd.Println("mangal latest version is " + version)
	},
}
