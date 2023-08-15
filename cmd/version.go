package cmd

import (
	"github.com/mangalorg/mangal/meta"
	"github.com/spf13/cobra"
)

var versionArgs = struct {
	Short bool
}{}

func init() {
	subcommands = append(subcommands, versionCmd)

	versionCmd.Flags().BoolVarP(&versionArgs.Short, "short", "s", false, "just show the version number")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if versionArgs.Short {
			cmd.Println(meta.Version)
			return
		}

		cmd.Println(meta.PrettyVersion())
	},
}
