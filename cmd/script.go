package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/fs"
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
	rootCmd.AddCommand(scriptCmd)

	scriptCmd.Flags().StringVarP(&scriptArgs.File, "file", "f", "", "Read script from file")
	scriptCmd.Flags().StringVarP(&scriptArgs.String, "string", "s", "", "Read script from script")
	scriptCmd.Flags().BoolVarP(&scriptArgs.Stdin, "stdin", "i", false, "Read script from stdin")

	scriptCmd.MarkFlagsMutuallyExclusive("file", "string", "stdin")

	scriptCmd.Flags().StringVarP(&scriptArgs.Provider, "provider", "p", "", "Load provider by tag")

	scriptCmd.Flags().StringToStringVarP(&scriptArgs.Variables, "vars", "v", nil, "Variables to set in the `Vars` table")

	scriptCmd.RegisterFlagCompletionFunc("provider", completionProviderIDs)
}

var scriptCmd = &cobra.Command{
	Use:   "script",
	Short: "Run mangal in scripting mode",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var reader io.Reader

		switch {
		case cmd.Flag("file").Changed:
			file, err := fs.Afero.OpenFile(
				scriptArgs.File,
				os.O_RDONLY,
				0755,
			)
			if err != nil {
				return err
			}

			defer file.Close()

			reader = file
		case cmd.Flag("string").Changed:
			reader = strings.NewReader(scriptArgs.String)
		case cmd.Flag("stdin").Changed:
			reader = os.Stdin
		default:
			return errors.New("either `file`, `string` or `stdin`")
		}

		var options script.Options

		options.Variables = scriptArgs.Variables

		anilist := libmangal.NewAnilist(libmangal.DefaultAnilistOptions())

		options.Anilist = &anilist

		if scriptArgs.Provider != "" {
			loaders, err := manager.Loaders()
			if err != nil {
				return err
			}

			loader, ok := lo.Find(loaders, func(loader libmangal.ProviderLoader) bool {
				return loader.Info().ID == scriptArgs.Provider
			})

			if !ok {
				return fmt.Errorf("provider with ID %q not found", scriptArgs.Provider)
			}

			clientOptions := libmangal.DefaultClientOptions()

			client, err := libmangal.NewClient(context.Background(), loader, clientOptions)
			if err != nil {
				return err
			}

			options.Client = client
			options.Anilist = client.Anilist()
		}

		return script.Run(context.Background(), reader, options)
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

		err := fs.Afero.WriteFile(filename, []byte(l.LuaDoc()), 0755)

		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "%s Library specs written to %s\n", icon.Mark, filename)

		return nil
	},
}
