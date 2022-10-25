package cmd

import (
	"fmt"
	"github.com/metafates/mangal/icon"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/updater"
	"github.com/metafates/mangal/util"
	"os"
	"runtime"
	"strings"
	"text/template"

	"github.com/metafates/mangal/constant"
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
		erase := util.PrintErasable(fmt.Sprintf("%s Checking if new version is available...", icon.Get(icon.Progress)))
		var newVersion string

		version, err := updater.LatestVersion()
		if err == nil {
			comp, err := util.CompareVersions(constant.Version, version)
			if err == nil && comp == -1 {
				newVersion = version
			}
		}

		erase()

		versionInfo := struct {
			Version    string
			OS         string
			Arch       string
			BuiltAt    string
			BuiltBy    string
			Revision   string
			App        string
			NewVersion string
		}{
			Version:    constant.Version,
			App:        constant.Mangal,
			OS:         runtime.GOOS,
			Arch:       runtime.GOARCH,
			BuiltAt:    strings.TrimSpace(constant.BuiltAt),
			BuiltBy:    constant.BuiltBy,
			Revision:   constant.Revision,
			NewVersion: newVersion,
		}

		t, err := template.New("version").Funcs(map[string]any{
			"faint":   style.Faint,
			"bold":    style.Bold,
			"magenta": style.Magenta,
			"green":   style.Green,
			"repeat":  strings.Repeat,
			"concat": func(a, b string) string {
				return a + b
			},
		}).Parse(`{{ magenta "▇▇▇" }} {{ magenta .App }} 

  {{ faint "Version" }}         {{ bold .Version }}
  {{ faint "Git Commit" }}      {{ bold .Revision }} 
  {{ faint "Build Date" }}  	  {{ bold .BuiltAt }}
  {{ faint "Built By" }}        {{ bold .BuiltBy }}
  {{ faint "Platform" }}        {{ bold .OS }}/{{ bold .Arch }}

{{ if not (eq .NewVersion "") }}
{{ green "▇▇▇" }} New version available {{ bold .NewVersion }}
{{ faint (concat "https://github.com/metafates/mangal/releases/tag/v" .NewVersion) }}
{{ end }}`)
		handleErr(err)
		handleErr(t.Execute(cmd.OutOrStdout(), versionInfo))
	},
}
