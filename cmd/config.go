package cmd

import (
	"fmt"
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
			err := filesystem.Get().Remove(filepath.Join(where.Config(), "mangal.toml"))
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

		if lo.Must(filesystem.Get().Exists(configPath)) {
			handleErr(filesystem.Get().Remove(configPath))
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
				handleErr(fmt.Errorf("invalid boolean value: %s", v))
			}
		} else if isInt {
			var err error
			value, err = strconv.Atoi(v)
			handleErr(err)
		} else {
			value = v
		}

		if _, ok := config.DefaultValues[key]; !ok {
			handleErr(fmt.Errorf(`unknown key "%s"`, key))
		}

		expectedType := reflect.TypeOf(config.DefaultValues[key])
		actualType := reflect.TypeOf(value)

		if expectedType != actualType {
			handleErr(fmt.Errorf(`expected type "%s" but got "%s"`, expectedType, actualType))
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
