package cmd

import (
	"context"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/mangalorg/mangal/fs"
	"github.com/mangalorg/mangal/script"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	rootCmd.AddCommand(scriptCmd)

	scriptCmd.Flags().StringP("file", "f", "", "Read script from file")
	scriptCmd.Flags().StringP("string", "s", "", "Read script from script")
	scriptCmd.Flags().BoolP("stdin", "i", !terminal.IsTerminal(syscall.Stdin), "Read script from stdin")

	scriptCmd.MarkFlagsMutuallyExclusive("file", "string", "stdin")

	scriptCmd.Flags().StringToStringP("vars", "v", nil, "Variables to set in the `Vars` table")
}

var scriptCmd = &cobra.Command{
	Use:   "script",
	Short: "Run mangal in script mode",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var reader io.Reader

		// TODO: this does not work as intended
		switch {
		case cmd.Flag("file").Changed:
			value := lo.Must(cmd.Flags().GetString("file"))
			file, err := fs.Afero.OpenFile(
				value,
				os.O_RDONLY,
				0755,
			)
			if err != nil {
				return err
			}

			defer file.Close()

			reader = file
		case cmd.Flag("string").Changed:
			value := lo.Must(cmd.Flags().GetString("string"))
			reader = strings.NewReader(value)
		case cmd.Flag("stdin").Changed:
			reader = os.Stdin
		default:
			panic("unreachable")
		}

		// TODO: fill other options
		return script.Run(context.Background(), reader, script.Options{
			Variables: lo.Must(cmd.Flags().GetStringToString("vars")),
		})
	},
}
