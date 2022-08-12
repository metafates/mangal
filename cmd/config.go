package cmd

import (
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/where"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configInitCmd)
	configInitCmd.Flags().BoolP("force", "f", false, "overwrite existing config")

	configCmd.AddCommand(configRemoveCmd)
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
	RunE: func(cmd *cobra.Command, args []string) error {
		mangalDir := where.Config()
		if !lo.Must(filesystem.Get().Exists(mangalDir)) {
			_ = filesystem.Get().MkdirAll(mangalDir, os.ModePerm)
		}

		return viper.SafeWriteConfig()
	},
}

var configRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		mangalDir := where.Config()
		configPath := filepath.Join(mangalDir, constant.Mangal+".toml")

		if lo.Must(filesystem.Get().Exists(configPath)) {
			return filesystem.Get().Remove(configPath)
		}

		return nil
	},
}
