package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/metafates/mangal/color"
	"os"
	"path/filepath"
	"sort"
	"strconv"

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
)

func errUnknownKey(key string) error {
	closest := lo.MinBy(lo.Keys(config.Default), func(a string, b string) bool {
		return levenshtein.Distance(key, a) < levenshtein.Distance(key, b)
	})
	msg := fmt.Sprintf(
		"unknown key %s, did you mean %s?",
		style.Fg(color.Red)(key),
		style.Fg(color.Yellow)(closest),
	)

	return errors.New(msg)
}

func completionConfigKeys(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return lo.Keys(config.Default), cobra.ShellCompDirectiveNoFileComp
}

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Various config commands",
}

func init() {
	configCmd.AddCommand(configInfoCmd)
	configInfoCmd.Flags().StringSliceP("key", "k", []string{}, "The keys to get info for")
	configInfoCmd.Flags().BoolP("json", "j", false, "Output as JSON")
	_ = configInfoCmd.RegisterFlagCompletionFunc("key", completionConfigKeys)

	configInfoCmd.SetOut(os.Stdout)
}

var configInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show the info for each config field with description",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			keys   = lo.Must(cmd.Flags().GetStringSlice("key"))
			asJson = lo.Must(cmd.Flags().GetBool("json"))
			fields = lo.Values(config.Default)
		)

		if len(keys) > 0 {
			fields = make([]config.Field, 0, len(keys))

			for _, key := range keys {
				if _, ok := config.Default[key]; !ok {
					handleErr(errUnknownKey(key))
				}

				fields = append(fields, config.Default[key])
			}
		}

		sort.Slice(fields, func(i, j int) bool {
			return fields[i].Key < fields[j].Key
		})

		if asJson {
			encoder := json.NewEncoder(cmd.OutOrStdout())
			lo.Must0(encoder.Encode(fields))
			return
		}

		for i, field := range fields {
			fmt.Print(field.Pretty())

			if i < len(fields)-1 {
				fmt.Println()
				fmt.Println()
			}
		}
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
	configSetCmd.Flags().StringP("key", "k", "", "The key to set the value for")
	lo.Must0(configSetCmd.MarkFlagRequired("key"))
	_ = configSetCmd.RegisterFlagCompletionFunc("key", completionConfigKeys)

	configSetCmd.Flags().StringSliceP("value", "v", []string{}, "The value to set")
	lo.Must0(configSetCmd.MarkFlagRequired("value"))

	// deprecated flags for backwards compatibility
	configSetCmd.Flags().BoolP("bool", "b", false, "Set the value type to bool")
	configSetCmd.Flags().IntP("int", "i", 0, "Set the value type to int")
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a config value",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			key   = lo.Must(cmd.Flags().GetString("key"))
			value = lo.Must(cmd.Flags().GetStringSlice("value"))
		)

		if _, ok := config.Default[key]; !ok {
			handleErr(errUnknownKey(key))
		}

		var v any
		switch config.Default[key].Value.(type) {
		case string:
			v = value[0]
		case int:
			parsedInt, err := strconv.ParseInt(value[0], 10, 64)
			if err != nil {
				handleErr(fmt.Errorf("invalid integer value: %s", value))
			}

			v = int(parsedInt)
		case bool:
			parsedBool, err := strconv.ParseBool(value[0])
			if err != nil {
				handleErr(fmt.Errorf("invalid boolean value: %s", value))
			}

			v = parsedBool
		case []string:
			v = value
		}

		viper.Set(key, v)
		switch err := viper.WriteConfig(); err.(type) {
		case viper.ConfigFileNotFoundError:
			handleErr(viper.SafeWriteConfig())
		default:
			handleErr(err)
		}

		fmt.Printf(
			"%s set %s to %s\n",
			style.Fg(color.Green)(icon.Get(icon.Success)),
			style.Fg(color.Purple)(key),
			style.Fg(color.Yellow)(fmt.Sprintf("%v", v)),
		)
	},
}

func init() {
	configCmd.AddCommand(configGetCmd)
	configGetCmd.Flags().StringP("key", "k", "", "The key to get the value for")
	lo.Must0(configGetCmd.MarkFlagRequired("key"))
	_ = configGetCmd.RegisterFlagCompletionFunc("key", completionConfigKeys)
}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a config value",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			key = lo.Must(cmd.Flags().GetString("key"))
		)

		if _, ok := config.Default[key]; !ok {
			handleErr(errUnknownKey(key))
		}

		fmt.Println(viper.Get(key))
	},
}

func init() {
	configCmd.AddCommand(configWriteCmd)
	configWriteCmd.Flags().BoolP("force", "f", false, "Force overwrite of existing config file")
}

var configWriteCmd = &cobra.Command{
	Use:   "write",
	Short: "Write current config to the file",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			force          = lo.Must(cmd.Flags().GetBool("force"))
			configFilePath = filepath.Join(
				where.Config(),
				fmt.Sprintf("%s.%s", constant.Mangal, "toml"),
			)
		)

		if force {
			err := filesystem.
				Api().
				Remove(configFilePath)

			handleErr(err)
		}

		handleErr(viper.SafeWriteConfig())
		fmt.Printf(
			"%s wrote config to %s\n",
			style.Fg(color.Green)(icon.Get(icon.Success)),
			configFilePath,
		)
	},
}

func init() {
	configCmd.AddCommand(configDeleteCmd)
}

var configDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete the config file",
	Aliases: []string{"remove"},
	Run: func(cmd *cobra.Command, args []string) {
		err := filesystem.
			Api().
			Remove(
				filepath.Join(
					where.Config(),
					fmt.Sprintf("%s.%s", constant.Mangal, "toml"),
				),
			)

		handleErr(err)
		fmt.Printf(
			"%s deleted config\n",
			style.Fg(color.Green)(icon.Get(icon.Success)),
		)
	},
}

func init() {
	configCmd.AddCommand(configResetCmd)

	configResetCmd.Flags().StringP("key", "k", "", "The key to reset the value for")
	configResetCmd.Flags().BoolP("all", "a", false, "Reset all config values")
	configResetCmd.MarkFlagsMutuallyExclusive("key", "all")
	_ = configResetCmd.RegisterFlagCompletionFunc("key", completionConfigKeys)
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the config key to default",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("key") && !cmd.Flags().Changed("all") {
			handleErr(fmt.Errorf("either --key or --all must be set"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var (
			key = lo.Must(cmd.Flags().GetString("key"))
			all = lo.Must(cmd.Flags().GetBool("all"))
		)

		if all {
			for key, field := range config.Default {
				viper.Set(key, field.Value)
			}
		} else if _, ok := config.Default[key]; !ok {
			handleErr(errUnknownKey(key))
		} else {
			viper.Set(key, config.Default[key].Value)
		}

		switch err := viper.WriteConfig(); err.(type) {
		case viper.ConfigFileNotFoundError:
			handleErr(viper.SafeWriteConfig())
		default:
			handleErr(err)
		}

		if all {
			fmt.Printf(
				"%s reset all config values\n",
				style.Fg(color.Green)(icon.Get(icon.Success)),
			)
		} else {
			fmt.Printf(
				"%s reset %s to default value %s\n",
				style.Fg(color.Green)(icon.Get(icon.Success)),
				style.Fg(color.Purple)(key),
				style.Fg(color.Yellow)(fmt.Sprintf("%v", config.Default[key].Value)),
			)
		}
	},
}
