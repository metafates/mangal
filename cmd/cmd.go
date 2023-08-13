package cmd

import (
	"fmt"
	"os"

	cc "github.com/ivanpirog/coloredcobra"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/icon"
	"github.com/mangalorg/mangal/meta"
	"github.com/mangalorg/mangal/provider/manager"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

const groupMode = "run"

var rootCmd = &cobra.Command{
	Use:   meta.AppName,
	Short: "The ultimate CLI manga downloader",
	Args:  cobra.NoArgs,
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

	rootCmd.AddGroup(&cobra.Group{
		ID:    groupMode,
		Title: "Mode",
	})

	if err := rootCmd.Execute(); err != nil {
		errorf(rootCmd, err.Error())
	}
}
