package cmd

import (
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/constants"
	"github.com/metafates/mangal/converter"
	"github.com/metafates/mangal/icon"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/tui"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func init() {
	rootCmd.PersistentFlags().StringP("format", "f", "", "output format")
	lo.Must0(rootCmd.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return converter.Available(), cobra.ShellCompDirectiveDefault
	}))
	lo.Must0(viper.BindPFlag(config.FormatsUse, rootCmd.PersistentFlags().Lookup("format")))

	rootCmd.PersistentFlags().StringP("icons", "i", "", "icons variant")
	lo.Must0(rootCmd.RegisterFlagCompletionFunc("icons", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return icon.AvailableVariants(), cobra.ShellCompDirectiveDefault
	}))
	lo.Must0(viper.BindPFlag(config.IconsVariant, rootCmd.PersistentFlags().Lookup("icons")))

	rootCmd.PersistentFlags().BoolP("history", "H", true, "write history of read chapters")
	lo.Must0(viper.BindPFlag(config.HistorySaveOnRead, rootCmd.PersistentFlags().Lookup("history")))

	rootCmd.Flags().BoolP("continue", "c", false, "continue reading")

	// Clear temporary files on startup
	go clearTemp()
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   constants.Mangal,
	Short: "The ultimate manga downloader",
	Long: style.Combined(style.Yellow, style.Bold)(constants.AssciiArtLogo) + "\n" +
		style.Combined(style.HiRed, style.Italic)("    - The ultimate cli manga downloader"),
	RunE: func(cmd *cobra.Command, args []string) error {
		options := tui.Options{
			Continue: lo.Must(cmd.Flags().GetBool("continue")),
		}
		return tui.Run(&options)
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
