package cmd

import (
	"github.com/metafates/mangal/source"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mangal",
	Short: "The ultimate manga downloader",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		sources := lo.Must(source.AvailableSources())

		cmd.Println("Available sources:")
		for _, source := range sources {
			cmd.Println("  " + filepath.Base(source))
		}

		source := lo.Must(source.LoadSource(sources[0]))
		mangas := lo.Must(source.Search("berserk"))

		for _, manga := range mangas {
			cmd.Println(manga)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}

func init() {
}
