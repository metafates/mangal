package cmd

import (
	"fmt"
	"os"
	"strings"

	cc "github.com/ivanpirog/coloredcobra"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/config"
	"github.com/mangalorg/mangal/icon"
	"github.com/mangalorg/mangal/meta"
	"github.com/mangalorg/mangal/provider/manager"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

const groupMode = "mode"

var rootCmd = &cobra.Command{
	Use:  meta.AppName,
	Args: cobra.NoArgs,
}

var subcommands []*cobra.Command

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
	var root *cobra.Command

	switch lo.Must(config.ModeString(config.Config.CLI.DefaultMode.Get())) {
	case config.ModeNone:
		root = rootCmd
	case config.ModeTUI:
		root = tuiCmd
	case config.ModeScript:
		root = scriptCmd
	case config.ModeWeb:
		root = webCmd
	}

	for _, subcommand := range subcommands {
		if subcommand == root {
			continue
		}

		root.AddCommand(subcommand)
	}

	root.Use = strings.Replace(root.Use, root.Name(), rootCmd.Name(), 1)
	root.Long = "The ultimate CLI manga downloader\n\n" + root.Short + " (configured as default)"
	root.AddGroup(&cobra.Group{
		ID:    groupMode,
		Title: "Mode",
	})

	cc.Init(&cc.Config{
		RootCmd:         root,
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

	if err := root.Execute(); err != nil {
		errorf(root, err.Error())
	}
}
