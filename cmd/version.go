package cmd

import (
	"fmt"
	"github.com/metafates/mangal/common"
	"github.com/metafates/mangal/style"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Long:  fmt.Sprintf("Shows %s versions and build date", common.Mangal),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s version %s\n", common.Mangal, style.Accent.Render(common.Version))
	},
}

func init() {
	mangalCmd.AddCommand(versionCmd)
}
