package cmd

import (
	"github.com/metafates/mangal/color"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/version"
	"github.com/samber/lo"
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
	versionCmd.Flags().BoolP("short", "s", false, "print short version")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of mangal",
	Long:  `All software has versions. This is mangal's`,
	Run: func(cmd *cobra.Command, args []string) {
		if lo.Must(cmd.Flags().GetBool("short")) {
			cmd.Println(constant.Version)
			return
		}

		defer version.Notify()

		versionInfo := struct {
			Version  string
			OS       string
			Arch     string
			BuiltAt  string
			BuiltBy  string
			Revision string
			App      string
		}{
			Version:  constant.Version,
			App:      constant.Mangal,
			OS:       runtime.GOOS,
			Arch:     runtime.GOARCH,
			BuiltAt:  strings.TrimSpace(constant.BuiltAt),
			BuiltBy:  constant.BuiltBy,
			Revision: constant.Revision,
		}

		t, err := template.New("version").Funcs(map[string]any{
			"faint":   style.Faint,
			"bold":    style.Bold,
			"magenta": style.Fg(color.Purple),
			"green":   style.Fg(color.Green),
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
`)
		handleErr(err)
		handleErr(t.Execute(cmd.OutOrStdout(), versionInfo))
	},
}
