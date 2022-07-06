package cmd

import (
	"fmt"
	"github.com/metafates/mangal/common"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/util"
	"github.com/spf13/cobra"
	"log"
)

var latestCmd = &cobra.Command{
	Use:   "latest",
	Short: fmt.Sprintf("Check if latest version of %s is used", common.Mangal),
	Long:  "Fetches the latest version from the GitHub and compares it with current version",
	Run: func(cmd *cobra.Command, args []string) {
		const githubReleaseURL = "https://github.com/metafates/mangal/releases/latest"

		latestVersion, err := util.FetchLatestVersion()

		if err != nil || latestVersion == "" {
			log.Fatalf("Can't find latest version\nYou can visit %s to check for updates", githubReleaseURL)
		}

		// check if current version is latest
		if latestVersion <= common.Version {
			fmt.Printf("You are using the latest version of %s\n", common.Mangal)
		} else {
			fmt.Printf("New version of %s is available: %s\n", common.Mangal, style.AccentStyle.Render(latestVersion))
			fmt.Printf("You can download it from %s\n", style.AccentStyle.Render(githubReleaseURL))
			fmt.Println("Or use your package manager to update")
		}
	},
}

func init() {
	mangalCmd.AddCommand(latestCmd)
}
