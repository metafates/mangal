package cmd

import (
	"fmt"
	"github.com/metafates/mangal/common"
	"github.com/metafates/mangal/style"
	"github.com/spf13/cobra"
	"os"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Show environment variables",
	Long:  "Show environment variables and their values",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(style.BoldStyle.Render("Available environment variables"))
		fmt.Println()

		for envVar, description := range common.AvailableEnvVars {
			value, isSet := os.LookupEnv(envVar)
			fmt.Printf("%s - %s\n", style.AccentStyle.Render(envVar), description)

			if isSet {
				fmt.Printf("%s - %s\n", "Value", value)
			} else {
				fmt.Printf("%s - %s\n", "Value", style.FailStyle.Render("Not set"))
			}

			fmt.Println()
		}
	},
}

func init() {
	mangalCmd.AddCommand(envCmd)
}
