package cmd

import (
	"fmt"
	"github.com/metafates/mangal/icon"
	"github.com/metafates/mangal/source"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().BoolP("lenient", "l", false, "do not warn about missing functions")
}

var runCmd = &cobra.Command{
	Use:   "run [file]",
	Short: "Run lua file",
	Long: `Runs Lua5.1 VM. Useful for debugging.
Or you can use mangal as a standalone lua interpreter.`,
	Args:    cobra.ExactArgs(1),
	Example: "  mangal run ./test.lua",
	Run: func(cmd *cobra.Command, args []string) {
		sourcePath := args[0]

		// LoadSource runs file when it's loaded
		_, err := source.LoadSource(sourcePath, !lo.Must(cmd.Flags().GetBool("lenient")))
		if err != nil {
			fmt.Println(icon.Get(icon.Fail) + " " + err.Error())
			os.Exit(1)
		}
	},
}
