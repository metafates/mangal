package cmd

import (
	"github.com/metafates/mangal/style"
	"os"
	"runtime"
	"strings"
	"text/template"

	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/updater"
	"github.com/spf13/cobra"
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

		versionInfo := struct {
			Version       string
			InstalledWith string
			OS            string
			Arch          string
			BuiltAt       string
			BuiltBy       string
			Revision      string
			Mangal        string
		}{
			Version:       constant.Version,
			Mangal:        constant.Mangal,
			InstalledWith: installedWith,
			OS:            runtime.GOOS,
			Arch:          runtime.GOARCH,
			BuiltAt:       strings.TrimSpace(constant.BuiltAt),
			BuiltBy:       constant.BuiltBy,
			Revision:      constant.Revision,
		}

		t, err := template.New("version").Funcs(map[string]any{
			"faint":   style.Faint,
			"bold":    style.Bold,
			"magenta": style.Magenta,
			"repeat":  strings.Repeat,
		}).Parse(`{{ magenta "▇▇▇" }} {{ magenta .Mangal }} 

  {{ faint "Version" }}        {{ .Version }}
  {{ faint "Git Commit" }}     {{ .Revision }} 
  {{ faint "Build Date" }} 	 {{ .BuiltAt }}
  {{ faint "Built By" }}       {{ .BuiltBy }}
  {{ faint "Installed With" }} {{ .InstalledWith }} 
  {{ faint "Platform" }}       {{ .OS }}/{{ .Arch }}
`)
		handleErr(err)
		handleErr(t.Execute(cmd.OutOrStdout(), versionInfo))
	},
}
