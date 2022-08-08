package cmd

import (
	"github.com/metafates/mangal/filesystem"
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
		configDir := lo.Must(os.UserConfigDir())

		mangalDir := filepath.Join(configDir, "mangal")
		if !lo.Must(filesystem.Get().Exists(mangalDir)) {
			_ = filesystem.Get().MkdirAll(mangalDir, os.ModePerm)
		}

		if lo.Must(cmd.Flags().GetBool("force")) {
			err := viper.WriteConfig()
			switch err.(type) {
			case viper.ConfigFileNotFoundError:
				return viper.SafeWriteConfig()
			default:
				return err
			}
		}
		return viper.SafeWriteConfig()
	},
}
