package cmd

import (
	"context"
	"fmt"
	"net/url"

	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/path"
	"github.com/mangalorg/mangal/provider/info"
	"github.com/mangalorg/mangal/provider/manager"
	"github.com/spf13/cobra"
)

func init() {
	subcommands = append(subcommands, providersCmd)
}

var providersCmd = &cobra.Command{
	Use:     "providers",
	Aliases: []string{"p"},
	Short:   "Providers management",
	Args:    cobra.NoArgs,
}

func init() {
	providersCmd.AddCommand(providersAddCmd)
}

var providersAddCmd = &cobra.Command{
	Use:   "add <url>",
	Short: "Install provider",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		URL, err := url.Parse(args[0])
		if err != nil {
			return err
		}

		return manager.Add(context.Background(), manager.AddOptions{
			URL: URL,
		})
	},
}

func init() {
	providersCmd.AddCommand(providersUpCmd)
}

var providersUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Update providers",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return manager.Update(context.Background(), manager.UpdateOptions{})
	},
}

func init() {
	providersCmd.AddCommand(providersLsCmd)
}

var providersLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List installed providers",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		loaders, err := manager.Loaders()
		if err != nil {
			return err
		}

		for _, loader := range loaders {
			fmt.Println(loader.Info().ID)
		}

		return nil
	},
}

func init() {
	providersCmd.AddCommand(providersRmCmd)
}

var providersRmCmd = &cobra.Command{
	Use:   "rm [tags]",
	Short: "Remove provider",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, tag := range args {
			if err := manager.Remove(tag); err != nil {
				errorf(cmd, err.Error())
			}
		}
	},
}

var providersNewArgs = struct {
	Dir string
}{}

func init() {
	providersCmd.AddCommand(providersNewCmd)

	providersNewCmd.Flags().StringVarP(&providersNewArgs.Dir, "dir", "d", path.ProvidersDir(), "directory inside which create a new provider")

	providersNewCmd.MarkFlagDirname("dir")
}

var providersNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new provider",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		options := manager.NewOptions{
			Dir: providersNewArgs.Dir,
			Info: info.Info{
				Provider: libmangal.ProviderInfo{
					ID:          "test",
					Name:        "test",
					Version:     "0.1.0",
					Description: "Lorem ipsum",
					Website:     "example.com",
				},
				Type: info.TypeLua,
			},
		}

		if err := manager.New(options); err != nil {
			errorf(cmd, err.Error())
		}
	},
}
