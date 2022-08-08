package cmd

import (
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/style"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func init() {
	rootCmd.AddCommand(envCmd)
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Show available environment variables",
	Long:  `Show available environment variables.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, env := range config.EnvExposed {
			env = strings.ToUpper(constant.Mangal + "_" + config.EnvKeyReplacer.Replace(env))
			value, present := os.LookupEnv(env)

			cmd.Println(style.Combined(style.Bold, style.Magenta)(env))

			if present && value != "" {
				cmd.Println(style.Green(value))
			} else {
				cmd.Println(style.Red("unset"))
			}

			cmd.Println()
		}
	},
}
