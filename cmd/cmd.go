package cmd

import (
	"fmt"
	"os"

	cc "github.com/ivanpirog/coloredcobra"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/config"
	"github.com/mangalorg/mangal/icon"
	"github.com/mangalorg/mangal/meta"
	"github.com/mangalorg/mangal/provider/manager"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var rootArgs = struct {
	Version bool
}{}

func init() {
	rootCmd.Flags().BoolVarP(&rootArgs.Version, "version", "v", false, "show version information")
}

var rootCmd = &cobra.Command{
	Use:   meta.AppName,
	Short: "The ultimate CLI manga downloader",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if rootArgs.Version {
			versionCmd.Run(versionCmd, nil)
			return
		}

		mode := lo.Must(config.ModeString(config.Config.CLI.DefaultMode.Get()))

		var cmdToExecute *cobra.Command
		switch mode {
		case config.ModeTUI:
			cmdToExecute = tuiCmd
		case config.ModeScript:
			cmdToExecute = scriptCmd
		case config.ModeWeb:
			cmdToExecute = webCmd
		default:
			panic("unreachable")
		}

		if len(args) == 0 {
			if args := config.Config.CLI.DefaultModeArgs.Get(); len(args) != 0 {
				cmdToExecute.SetArgs(args)
			}
		} else {
			cmdToExecute.SetArgs(args)
		}

		_, err := cmdToExecute.ExecuteC()
		if err != nil {
			errorf(cmdToExecute, err.Error())
		}
	},
}

func completionProviderIDs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	loaders, err := manager.Loaders()

	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	IDs := lo.Map(loaders, func(loader libmangal.ProviderLoader, _ int) string {
		return loader.Info().ID
	})

	return IDs, cobra.ShellCompDirectiveDefault
}

func successf(cmd *cobra.Command, format string, a ...any) {
	cmd.Printf(fmt.Sprintf("%s %s\n", icon.Check, format), a...)
}

func errorf(cmd *cobra.Command, format string, a ...any) {
	cmd.PrintErrf(fmt.Sprintf("%s %s\n", icon.Cross, format), a...)
	os.Exit(1)
}

func Execute() {
	cc.Init(&cc.Config{
		RootCmd:         rootCmd,
		Headings:        cc.HiCyan + cc.Bold + cc.Underline,
		Commands:        cc.HiYellow + cc.Bold,
		Example:         cc.Italic,
		ExecName:        cc.Bold,
		Flags:           cc.Bold,
		FlagsDataType:   cc.Italic + cc.HiBlue,
		Aliases:         cc.Italic,
		NoExtraNewlines: true,
		NoBottomNewline: true,
	})

	if err := rootCmd.Execute(); err != nil {
		errorf(rootCmd, err.Error())
	}
}
