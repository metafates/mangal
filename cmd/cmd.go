package cmd

import (
	"log"

	cc "github.com/ivanpirog/coloredcobra"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/meta"
	"github.com/mangalorg/mangal/provider/manager"
	"github.com/mangalorg/mangal/tui"
	"github.com/mangalorg/mangal/tui/state/providers"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		if rootArgs.Version {
			versionCmd.Run(versionCmd, nil)
			return nil
		}

		loaders, err := manager.Loaders()
		if err != nil {
			return err
		}

		return tui.Run(providers.New(loaders))
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

func Execute() {
	cc.Init(&cc.Config{
		RootCmd:         rootCmd,
		Headings:        cc.HiCyan + cc.Bold + cc.Underline,
		Commands:        cc.HiYellow + cc.Bold,
		Example:         cc.Italic,
		ExecName:        cc.Bold,
		Flags:           cc.Bold,
		FlagsDataType:   cc.Italic + cc.HiBlue,
		NoExtraNewlines: true,
		NoBottomNewline: true,
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
