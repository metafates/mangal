package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/afs"
	"github.com/mangalorg/mangal/anilist"
	"github.com/mangalorg/mangal/client"
	"github.com/mangalorg/mangal/icon"
	"github.com/mangalorg/mangal/provider/manager"
	"github.com/mangalorg/mangal/script"
	"github.com/mangalorg/mangal/script/lib"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	lua "github.com/yuin/gopher-lua"
)

var scriptArgs = struct {
	File      string
	String    string
	Stdin     bool
	Provider  string
	Variables map[string]string
}{}

func init() {
	subcommands = append(subcommands, scriptCmd)

	scriptCmd.Flags().StringVarP(&scriptArgs.File, "file", "f", "", "Read script from file")
	scriptCmd.Flags().StringVarP(&scriptArgs.String, "string", "s", "", "Read script from script")
	scriptCmd.Flags().BoolVarP(&scriptArgs.Stdin, "stdin", "i", false, "Read script from stdin")

	scriptCmd.MarkFlagsMutuallyExclusive("file", "string", "stdin")

	scriptCmd.Flags().StringVarP(&scriptArgs.Provider, "provider", "p", "", "Load provider by tag")

	scriptCmd.Flags().StringToStringVarP(&scriptArgs.Variables, "vars", "v", nil, "Variables to set in the `Vars` table")

	scriptCmd.RegisterFlagCompletionFunc("provider", completionProviderIDs)
}

var scriptCmd = &cobra.Command{
	Use:     "script",
	Short:   "Run mangal in scripting mode",
	GroupID: groupMode,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var reader io.Reader

		switch {
		case cmd.Flag("file").Changed:
			file, err := afs.Afero.OpenFile(
				scriptArgs.File,
				os.O_RDONLY,
				0755,
			)
			if err != nil {
				errorf(cmd, err.Error())
			}

			defer file.Close()

			reader = file
		case cmd.Flag("string").Changed:
			reader = strings.NewReader(scriptArgs.String)
		case cmd.Flag("stdin").Changed:
			reader = os.Stdin
		default:
			errorf(cmd, "either `file`, `string` or `stdin` is required")
		}

		var options script.Options

		options.Variables = scriptArgs.Variables
		options.Anilist = anilist.Client

		if scriptArgs.Provider != "" {
			loaders, err := manager.Loaders()
			if err != nil {
				errorf(cmd, err.Error())
			}

			loader, ok := lo.Find(loaders, func(loader libmangal.ProviderLoader) bool {
				return loader.Info().ID == scriptArgs.Provider
			})

			if !ok {
				errorf(cmd, "provider with ID %q not found", scriptArgs.Provider)
			}

			client, err := client.NewClient(context.Background(), loader)
			if err != nil {
				errorf(cmd, err.Error())
			}

			options.Client = client
		}

		if err := script.Run(context.Background(), reader, options); err != nil {
			errorf(cmd, err.Error())
		}
	},
}

func init() {
	scriptCmd.AddCommand(scriptDocCmd)
}

var scriptDocCmd = &cobra.Command{
	Use:   "doc",
	Short: "Generate documentation for the `mangal` lua library",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		l := lib.Lib(lua.NewState(), lib.Options{})

		filename := fmt.Sprint(l.Name, ".lua")

		err := afs.Afero.WriteFile(filename, []byte(l.LuaDoc()), 0755)

		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "%s Library specs written to %s\n", icon.Mark, filename)

		return nil
	},
}
