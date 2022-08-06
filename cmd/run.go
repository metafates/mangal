package cmd

import (
	"github.com/metafates/mangal/source"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		sourcePath := args[0]

		// LoadSource runs file when it's loaded
		_, err := source.LoadSource(sourcePath, !lo.Must(cmd.Flags().GetBool("lenient")))
		return err
	},
}
