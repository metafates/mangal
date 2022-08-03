package cmd

import (
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/constants"
	"github.com/metafates/mangal/converter"
	"github.com/metafates/mangal/style"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func init() {
	rootCmd.PersistentFlags().StringP("format", "f", "", "output format")
	_ = rootCmd.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return converter.Available(), cobra.ShellCompDirectiveDefault
	})
	_ = viper.BindPFlag(config.FormatsUse, rootCmd.PersistentFlags().Lookup("format"))
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   constants.Mangal,
	Short: "The ultimate manga downloader",
	Long: style.Combined(style.Yellow, style.Bold)(constants.AssciiArtLogo) + "\n" +
		style.Combined(style.HiRed, style.Italic)("    - The ultimate cli manga downloader"),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("TUI is not implemented yet")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	// colored cobra injection
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
