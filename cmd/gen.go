package cmd

import (
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().StringP("name", "n", "", "name of the source")
	genCmd.Flags().StringP("url", "u", "", "url of the website")

	lo.Must0(genCmd.MarkFlagRequired("name"))
	lo.Must0(genCmd.MarkFlagRequired("url"))
}

var genCmd = &cobra.Command{
	Use:        "gen",
	Short:      "Generate a new lua source",
	Long:       `Generate a new lua source.`,
	Deprecated: "use `mangal sources gen` instead.",
	Run:        sourcesGenCmd.Run,
}
