package cmd

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"

	"github.com/kballard/go-shellquote"
	"github.com/mangalorg/mangal/config"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func init() {
	subcommands = append(subcommands, configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

var configInfoArgs = struct {
	JSON bool
}{}

func init() {
	configCmd.AddCommand(configInfoCmd)

	configInfoCmd.Flags().BoolVarP(&configInfoArgs.JSON, "json", "j", false, "JSON output")
}

var configInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show configuration information",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if configInfoArgs.JSON {
			// TODO
			panic("unimplemented")
		}

		fieldTemplate := template.Must(template.New("field").Funcs(map[string]any{
			"get": func(key string) any {
				return config.Get(key)
			},
		}).Parse(`
{{ .Description }}

Key: {{ .Key }}
Value: {{ get .Key }}
Default: {{ .Default }}
`))

		var out strings.Builder
		for _, field := range config.Fields {
			var sb strings.Builder
			if err := fieldTemplate.Execute(&sb, field); err != nil {
				errorf(cmd, err.Error())
			}

			out.WriteString(sb.String())
		}

		// TODO: ???
		if pager := os.Getenv("PAGER"); pager != "" {
			cmdWithArgs, err := shellquote.Split(pager)
			if err != nil {
				return
			}

			pagerCmd := exec.Command(cmdWithArgs[0], cmdWithArgs[1:]...)
			pagerCmd.Stdin = strings.NewReader(out.String())
			if err := pagerCmd.Run(); err != nil {
				errorf(cmd, err.Error())
			}
		} else {
			cmd.Print(out.String())
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
