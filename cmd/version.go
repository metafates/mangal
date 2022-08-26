package cmd

import (
	"encoding/json"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/util"
	"github.com/spf13/cobra"
	"net/http"
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
		resp, err := http.Get("https://api.github.com/repos/metafates/mangal/releases/latest")
		handleErr(err)

		defer util.Ignore(resp.Body.Close)

		var release struct {
			TagName string `json:"tag_name"`
		}

		err = json.NewDecoder(resp.Body).Decode(&release)
		handleErr(err)

		// remove the v from the tag name
		latestVersion := release.TagName[1:]

		cmd.Println("mangal latest version is " + latestVersion)
	},
}
