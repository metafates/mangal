package cmd

import (
	"errors"
	"fmt"
	levenshtein "github.com/ka-weihe/fast-levenshtein"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/icon"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/where"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configInitCmd)
	configInitCmd.Flags().BoolP("force", "f", false, "overwrite existing config")

	configCmd.AddCommand(configRemoveCmd)
	configCmd.AddCommand(configSetCmd)
	configSetCmd.Flags().StringP("key", "k", "", "key to set")
	configSetCmd.Flags().StringP("value", "v", "", "value to set")
	configSetCmd.Flags().BoolP("bool", "b", false, "value is a boolean")
	configSetCmd.Flags().BoolP("int", "i", false, "value is an integer")

	lo.Must0(configSetCmd.MarkFlagRequired("key"))
	configSetCmd.MarkFlagsMutuallyExclusive("bool", "int")
	lo.Must0(configSetCmd.MarkFlagRequired("value"))
	lo.Must0(configSetCmd.RegisterFlagCompletionFunc("key", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return lo.Keys(config.DefaultValues), cobra.ShellCompDirectiveNoFileComp
	}))

	configCmd.AddCommand(configGetCmd)
	configGetCmd.Flags().StringP("key", "k", "", "key to get")
	lo.Must0(configGetCmd.MarkFlagRequired("key"))
	lo.Must0(configGetCmd.RegisterFlagCompletionFunc("key", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return lo.Keys(config.DefaultValues), cobra.ShellCompDirectiveNoFileComp
	}))

	configCmd.AddCommand(configInfoCmd)
	configInfoCmd.Flags().StringP("key", "k", "", "show only this key")
	lo.Must0(configInfoCmd.RegisterFlagCompletionFunc("key", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return lo.Keys(config.DefaultValues), cobra.ShellCompDirectiveNoFileComp
	}))
}

func errUnknownKey(key string) error {
	closest := lo.MinBy(lo.Keys(config.DefaultValues), func(a string, b string) bool {
		return levenshtein.Distance(key, a) < levenshtein.Distance(key, b)
	})
	msg := fmt.Sprintf(`unknown key %s, did you mean %s?`, style.Red(key), style.Yellow(closest))
	return errors.New(msg)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Various config actions",
	Long:  `Various config actions`,
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize config",
	Long:  `Initialize default config`,
	Run: func(cmd *cobra.Command, args []string) {
		force := lo.Must(cmd.Flags().GetBool("force"))
		if force {
			err := filesystem.Api().Remove(filepath.Join(where.Config(), "mangal.toml"))
			handleErr(err)
		}

		handleErr(viper.SafeWriteConfig())
	},
}

var configRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes config file",
	Run: func(cmd *cobra.Command, args []string) {
		mangalDir := where.Config()
		configPath := filepath.Join(mangalDir, constant.Mangal+".toml")

		if lo.Must(filesystem.Api().Exists(configPath)) {
			handleErr(filesystem.Api().Remove(configPath))
		}
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set config value",
	Long: `Set config value. Example:
mangal config set --key "formats.use" --value cbz
`,
	Run: func(cmd *cobra.Command, args []string) {
		key := lo.Must(cmd.Flags().GetString("key"))
		isBool := lo.Must(cmd.Flags().GetBool("bool"))
		isInt := lo.Must(cmd.Flags().GetBool("int"))

		v := lo.Must(cmd.Flags().GetString("value"))

		var value any

		if isBool {
			if v == "true" {
				value = true
			} else if v == "false" {
				value = false
			} else {
				handleErr(fmt.Errorf("invalid boolean value %s", style.Yellow(v)))
			}
		} else if isInt {
			var err error
			value, err = strconv.Atoi(v)
			handleErr(err)
		} else {
			value = v
		}

		if _, ok := config.DefaultValues[key]; !ok {
			handleErr(errUnknownKey(key))
		}

		expectedType := reflect.TypeOf(config.DefaultValues[key].Value)
		actualType := reflect.TypeOf(value)

		if expectedType != actualType {
			handleErr(fmt.Errorf(`expected type %s but got %s`, style.Blue(expectedType.String()), style.Red(actualType.String())))
		}

		viper.Set(key, value)
		switch err := viper.WriteConfig(); err.(type) {
		case viper.ConfigFileNotFoundError:
			handleErr(viper.SafeWriteConfig())
		default:
			handleErr(err)
		}

		fmt.Printf("%s set %s to %s\n", icon.Get(icon.Success), style.Magenta(key), style.Yellow(v))
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get config value",
	Long: `Get config value. Example:
mangal config get --key "formats.use"
`,
	Run: func(cmd *cobra.Command, args []string) {
		key := lo.Must(cmd.Flags().GetString("key"))
		if _, ok := config.DefaultValues[key]; !ok {
			handleErr(errUnknownKey(key))
		}

		fmt.Println(viper.Get(key))
	},
}

var configInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "List all config values with their types and descriptions",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			fields = make([]config.Field, len(config.DefaultValues))
			key    = lo.Must(cmd.Flags().GetString("key"))
		)

		i := 0
		for _, v := range config.DefaultValues {
			fields[i] = v
			i++
		}

		if key != "" {
			if field, ok := config.DefaultValues[key]; ok {
				fmt.Print(field.Pretty())
				return
			}
			handleErr(errUnknownKey(key))
		}

		sort.Slice(fields, func(i, j int) bool {
			return fields[i].Name < fields[j].Name
		})

		for _, field := range fields {
			// extra newline
			fmt.Printf("%s\n", field.Pretty())
		}
	},
}
