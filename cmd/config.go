package cmd

import (
	"fmt"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		fieldTemplate := template.Must(template.New("field").Parse(`
{{.Description}}

Key: {{.Key}}
Value: {{.Value}}
Default: {{.Default}}
`))

		var sb strings.Builder
		for _, field := range config.Fields {
			if err := fieldTemplate.Execute(&sb, field); err != nil {
				return err
			}

			fmt.Print(sb.String())
			sb.Reset()
		}

		return nil
	},
}

func init() {
	configCmd.AddCommand(configWriteCmd)
}

var configWriteCmd = &cobra.Command{
	Use:   "write",
	Short: "Write configuration to disk",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return config.Write()
	},
}

func init() {
	configCmd.AddCommand(configGetCmd)
}

var configGetCmd = &cobra.Command{
	Use:   "get key",
	Short: "Get config value by key",
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		keys := config.Keys()

		filtered := lo.Filter(keys, func(key string, _ int) bool {
			return strings.HasPrefix(key, toComplete)
		})

		return filtered, cobra.ShellCompDirectiveDefault
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		if !config.Exists(key) {
			return fmt.Errorf("config key %q doesn't exist", key)
		}

		fmt.Println(config.Get(key))
		return nil
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
}

var configSetCmd = &cobra.Command{
	Use:   "set key value",
	Short: "Sets value to the config key",
	Args:  cobra.ExactArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		keys := config.Keys()

		filtered := lo.Filter(keys, func(key string, _ int) bool {
			return strings.HasPrefix(key, toComplete)
		})

		return filtered, cobra.ShellCompDirectiveDefault
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		key, value := args[0], args[1]

		var converted any

		switch config.Get(key).(type) {
		case nil:
			return fmt.Errorf("unknown config key %q", key)
		case string:
			converted = value
		case int:
			parsedInt, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}

			converted = int(parsedInt)
		case bool:
			parsedBool, err := strconv.ParseBool(value)

			if err != nil {
				return err
			}

			converted = parsedBool
		default:
			return fmt.Errorf("unknown value type")
		}

		if err := config.Set(key, converted); err != nil {
			return err
		}

		return config.Write()
	},
}
