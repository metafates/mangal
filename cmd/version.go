package cmd

import (
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/style"
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
			constant.AsciiArtLogo,
			style.Italic(constant.Version),
			style.Italic(runtime.GOOS),
			style.Italic(runtime.GOARCH),
			style.Italic(constant.BuiltAt), style.Faint(constant.BuiltBy),
			style.Italic(constant.Revision),
			style.Italic(installedWith),
		)
	},
}
