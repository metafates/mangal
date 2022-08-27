package cmd

import (
	"fmt"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/converter"
	"github.com/metafates/mangal/icon"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/tui"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func init() {
	rootCmd.PersistentFlags().StringP("format", "F", "", "output format")
	lo.Must0(rootCmd.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return converter.Available(), cobra.ShellCompDirectiveDefault
	}))
	lo.Must0(viper.BindPFlag(constant.FormatsUse, rootCmd.PersistentFlags().Lookup("format")))

	rootCmd.PersistentFlags().StringP("icons", "I", "", "icons variant")
	lo.Must0(rootCmd.RegisterFlagCompletionFunc("icons", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return icon.AvailableVariants(), cobra.ShellCompDirectiveDefault
	}))
	lo.Must0(viper.BindPFlag(constant.IconsVariant, rootCmd.PersistentFlags().Lookup("icons")))

	rootCmd.PersistentFlags().BoolP("history", "H", true, "write history of the read chapters")
	lo.Must0(viper.BindPFlag(constant.HistorySaveOnRead, rootCmd.PersistentFlags().Lookup("history")))

	rootCmd.PersistentFlags().StringP("source", "S", "", "source")
	lo.Must0(rootCmd.RegisterFlagCompletionFunc("source", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		d := lo.Keys(provider.DefaultProviders())
		d = append(d, lo.Keys(provider.CustomProviders())...)

		return d, cobra.ShellCompDirectiveDefault
	}))
	lo.Must0(viper.BindPFlag(constant.DownloaderDefaultSource, rootCmd.PersistentFlags().Lookup("source")))

	rootCmd.Flags().BoolP("continue", "c", false, "continue reading")

	// Clear temporary files on startup
	go clearTemp()
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     constant.Mangal,
	Version: constant.Version,
	Short:   "The ultimate manga downloader",
	Long: style.Combined(style.Yellow, style.Bold)(constant.AssciiArtLogo) + "\n" +
		style.Combined(style.HiRed, style.Italic)("    - The ultimate cli manga downloader"),
	Run: func(cmd *cobra.Command, args []string) {
		options := tui.Options{
			Continue: lo.Must(cmd.Flags().GetBool("continue")),
		}
		handleErr(tui.Run(&options))
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

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleErr(err error) {
	if err != nil {
		log.Error(err)
		_, _ = fmt.Fprintf(os.Stderr, "%s %s\n", icon.Get(icon.Fail), strings.Trim(err.Error(), " \n"))
		os.Exit(1)
	}
}
