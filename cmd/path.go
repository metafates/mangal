package cmd

import (
	"fmt"
	"github.com/mangalorg/mangal/path"
	"github.com/spf13/cobra"
)

var pathArgs = struct {
	Config       bool `help:"Path to the config directory"`
	Cache        bool `help:"Path to the cache directory"`
	Temp         bool `help:"Path to a temporary directory"`
	Downloads    bool `help:"Path to the downloads directory"`
	LuaProviders bool `help:"Path to the lua providers directory"`
	Header       bool `help:"Print header" negatable:"" default:"true"`
}{}

func init() {
	rootCmd.AddCommand(pathCmd)

	pathCmd.Flags().BoolVar(&pathArgs.Config, "config", false, "Path to the config directory")
	pathCmd.Flags().BoolVar(&pathArgs.Cache, "cache", false, "Path to the cache directory")
	pathCmd.Flags().BoolVar(&pathArgs.Temp, "temp", false, "Path to a temporary directory")
	pathCmd.Flags().BoolVar(&pathArgs.Downloads, "downloads", false, "Path to the downloads directory")
	pathCmd.Flags().BoolVar(&pathArgs.LuaProviders, "lua-providers", false, "Path to the lua providers directory")
	pathCmd.Flags().BoolVar(&pathArgs.Header, "header", true, "Print header")
}

var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show paths",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		paths := []struct {
			Name  string
			Func  func() string
			Print bool
		}{
			{"Config", path.ConfigDir, pathArgs.Config},
			{"Cache", path.CacheDir, pathArgs.Cache},
			{"Temp", path.TempDir, pathArgs.Temp},
			{"Downloads", path.DownloadsDir, pathArgs.Downloads},
			{"Lua Providers", path.LuaProvidersDir, pathArgs.LuaProviders},
		}

		var anyPrinted bool
		for _, t := range paths {
			if t.Print {
				anyPrinted = true
				if pathArgs.Header {
					fmt.Println(t.Name)
				}
				fmt.Println(t.Func())
			}
		}

		if !anyPrinted {
			for _, t := range paths {
				if pathArgs.Header {
					fmt.Println(t.Name)
				}
				fmt.Println(t.Func())
			}
		}

		return nil
	},
}
