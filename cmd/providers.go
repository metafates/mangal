package cmd

import (
	"context"
	"fmt"
	"net/url"

	"github.com/mangalorg/mangal/provider/manager"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(providersCmd)
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
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, tag := range args {
			if err := manager.Remove(tag); err != nil {
				return err
			}
		}

		return nil
	},
}
