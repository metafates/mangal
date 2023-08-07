package cmd

import (
	"strconv"
	"strings"
	"text/template"

	"github.com/mangalorg/mangal/config"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

func init() {
	configCmd.AddCommand(configInfoCmd)
}

var configInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show configuration information",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fieldTemplate := template.Must(template.New("field").Parse(`
{{.Description}}

Key: {{.Key}}
Value: {{.Value}}
Default: {{.Default}}
`))

		var sb strings.Builder
		for _, field := range config.Fields {
			if err := fieldTemplate.Execute(&sb, field); err != nil {
				errorf(cmd, err.Error())
			}

			cmd.Print(sb.String())
			sb.Reset()
		}
	},
}

func init() {
	configCmd.AddCommand(configWriteCmd)
}

var configWriteCmd = &cobra.Command{
	Use:   "write",
	Short: "Write configuration to disk",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.Write(); err != nil {
			errorf(cmd, err.Error())
		}

		successf(cmd, "Wrote config to the file")
	},
}

func init() {
	configCmd.AddCommand(configGetCmd)
}

var configGetCmd = &cobra.Command{
	Use:           "get key",
	Short:         "Get config value by key",
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		keys := config.Keys()

		filtered := lo.Filter(keys, func(key string, _ int) bool {
			return strings.HasPrefix(key, toComplete)
		})

		return filtered, cobra.ShellCompDirectiveDefault
	},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		if !config.Exists(key) {
			errorf(cmd, "config key %q doesn't exist", key)
		}

		cmd.Println(config.Get(key))
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
}

var configSetCmd = &cobra.Command{
	Use:           "set key value",
	Short:         "Sets value to the config key",
	Args:          cobra.ExactArgs(2),
	SilenceErrors: true,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		keys := config.Keys()

		filtered := lo.Filter(keys, func(key string, _ int) bool {
			return strings.HasPrefix(key, toComplete)
		})

		return filtered, cobra.ShellCompDirectiveDefault
	},
	Run: func(cmd *cobra.Command, args []string) {
		key, value := args[0], args[1]

		var converted any

		switch config.Get(key).(type) {
		case nil:
			errorf(cmd, "unknown config key %q", key)
		case string:
			converted = value
		case int:
			parsedInt, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				errorf(cmd, err.Error())
			}

			converted = int(parsedInt)
		case bool:
			parsedBool, err := strconv.ParseBool(value)

			if err != nil {
				errorf(cmd, err.Error())
			}

			converted = parsedBool
		default:
			errorf(cmd, "unknown value type")
		}

		if err := config.Set(key, converted); err != nil {
			errorf(cmd, err.Error())
		}

		if err := config.Write(); err != nil {
			errorf(cmd, err.Error())
		}

		successf(cmd, "Successfully set %q to %v", key, converted)
	},
}
