package cmd

import (
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/where"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"os"
	"strings"
)

func init() {
	rootCmd.AddCommand(envCmd)
	envCmd.Flags().BoolP("filter", "f", false, "filter out variables that are not set")
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Show available environment variables",
	Long:  `Show available environment variables.`,
	Run: func(cmd *cobra.Command, args []string) {
		filter := lo.Must(cmd.Flags().GetBool("filter"))

		config.EnvExposed = append(config.EnvExposed, where.EnvConfigPath)
		slices.Sort(config.EnvExposed)
		for _, env := range config.EnvExposed {
			env = strings.ToUpper(constant.Mangal + "_" + config.EnvKeyReplacer.Replace(env))
			value := os.Getenv(env)
			present := value != ""

			if !present && filter {
				continue
			}

			cmd.Print(style.Combined(style.Bold, style.Magenta)(env))
			cmd.Print("=")

			if present {
				cmd.Println(style.Green(value))
			} else {
				cmd.Println(style.Red("unset"))
			}
		}
	},
}
