package cmd

import (
	"context"
	"net/url"

	"github.com/mangalorg/mangal/provider/manager"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(providersCmd)
}

var providersCmd = &cobra.Command{
	Use:   "providers",
	Short: "Providers management",
	Args:  cobra.NoArgs,
}

func init() {
	providersCmd.AddCommand(providersAddCmd)

	providersAddCmd.Flags().StringP("tag", "t", "", "Tag to use for the provider")
}

var providersAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Install provider",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tag, _ := cmd.Flags().GetString("tag")
		URL, err := url.Parse(args[0])
		if err != nil {
			return err
		}

		return manager.Add(context.Background(), manager.AddOptions{
			Tag: tag,
			URL: URL,
		})
	},
}

func init() {
	providersCmd.AddCommand(providersUpdateCmd)

	providersUpdateCmd.Flags().StringP("tag", "t", "", "Update specific provider by tag")
}

var providersUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update providers",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		tag, _ := cmd.Flags().GetString("tag")

		return manager.Update(context.Background(), manager.UpdateOptions{
			Tag: tag,
		})
	},
}
