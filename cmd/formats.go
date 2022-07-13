package cmd

import (
	"fmt"
	"github.com/metafates/mangal/common"
	"github.com/metafates/mangal/style"
	"github.com/spf13/cobra"
)

var formatsCmd = &cobra.Command{
	Use:   "formats",
	Short: "Information about available formats",
	Long:  "Show information about available formats with quick description of each",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(style.Bold.Render("Available formats") + "\n\n")
		for _, format := range common.AvailableFormats {
			fmt.Printf("%s - %s\n", style.Accent.Render(string(format)), common.FormatsInfo[format])
		}
	},
}

func init() {
	mangalCmd.AddCommand(formatsCmd)
}
