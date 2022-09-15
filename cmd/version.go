package cmd

import (
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/updater"
	"github.com/spf13/cobra"
	"os"
	"runtime"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.SetOut(os.Stdout)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of mangal",
	Long:  `All software has versions. This is mangal's`,
	Run: func(cmd *cobra.Command, args []string) {
		var installedWith string

		switch updater.DetectInstallationMethod() {
		case updater.Homebrew:
			installedWith = "Homebrew"
		case updater.Scoop:
			installedWith = "Scoop"
		case updater.Termux:
			installedWith = "Termux"
		case updater.Standalone:
			installedWith = "Standalone"
		case updater.Go:
			installedWith = "From source (" + runtime.Version() + ")"
		default:
			installedWith = "Unknown"
		}

		cmd.Printf(`%s

Version:        %s
OS:             %s
Arch:           %s
Built:          %s by %s
Revision:       %s
Installed With: %s
`,
			constant.AssciiArtLogo,
			constant.Version,
			runtime.GOOS,
			runtime.GOARCH,
			constant.BuiltAt, constant.BuiltBy,
			constant.Revision,
			installedWith,
		)
	},
}
