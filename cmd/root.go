package cmd

import (
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/metafates/mangal/constants"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   constants.Mangal,
	Short: "The ultimate manga downloader",
	Long: `                                __
  /\/\   __ _ _ __   __ _  __ _| |
 /    \ / _  | '_ \ / _  |/ _  | |
/ /\/\ \ (_| | | | | (_| | (_| | |
\/    \/\__,_|_| |_|\__, |\__,_|_|
                    |___/

	- The ultimate cli manga downloader`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("Hello, world!")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	cc.Init(&cc.Config{
		RootCmd:       rootCmd,
		Headings:      cc.HiCyan + cc.Bold + cc.Underline,
		Commands:      cc.HiYellow + cc.Bold,
		Example:       cc.Italic,
		ExecName:      cc.Bold,
		Flags:         cc.Bold,
		FlagsDataType: cc.Italic + cc.HiBlue,
	})

	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}
