package cmd

import (
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/source"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Run source in debug mode",
	Run: func(cmd *cobra.Command, args []string) {
		sourcePath := lo.Must(cmd.Flags().GetString("file"))
		sourceContents, err := filesystem.Get().ReadFile(sourcePath)

		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}

		compiled, err := source.Compile(sourcePath, strings.NewReader(string(sourceContents)))
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}

		// LoadSource runs file when it's loaded
		_, err = source.LoadSource(sourcePath, compiled)
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)

	debugCmd.Flags().StringP("file", "f", "", "source file to run")
	_ = debugCmd.MarkFlagRequired("file")
	_ = debugCmd.MarkFlagFilename("file", "lua")
}
