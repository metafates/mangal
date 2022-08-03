package config

import (
	"github.com/metafates/mangal/constants"
	"github.com/metafates/mangal/filesystem"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

// Setup initializes the configuration
func Setup() error {
	setName()
	setFs()
	setEnvs()
	setDefaults()
	setPaths()

	err := viper.ReadInConfig()

	switch err.(type) {
	case viper.ConfigFileNotFoundError:
		// Use defaults then
		return nil
	default:
		return err
	}
}

func setName() {
	viper.SetConfigName(constants.Mangal)
	viper.SetConfigType("toml")
}

func setFs() {
	viper.SetFs(filesystem.Get())
}

// setPaths sets the paths to the config files
func setPaths() {
	paths := lo.Must(Paths())

	for _, path := range paths {
		viper.AddConfigPath(path)
	}
}

// setEnvs sets the environment variables
func setEnvs() {
	viper.SetEnvPrefix(constants.Mangal)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	for _, env := range envFields {
		viper.MustBindEnv(env)
	}
}

// setDefaults sets the default values
func setDefaults() {
	viper.SetTypeByDefaultValue(true)
	configDir := lo.Must(os.UserConfigDir())

	fields := map[string]any{
		// Downloader
		DownloaderPath:                ".",
		DownloaderChapterNameTemplate: "[%0d] %s",

		// Formats
		FormatsUse: "plain",

		// Sources
		SourcesPath: filepath.Join(configDir, constants.Mangal, "sources"),

		// Mini-mode
		MiniVimMode: false,
	}

	for field, value := range fields {
		viper.SetDefault(field, value)
	}
}

// Paths returns the paths to the config files
func Paths() ([]string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return []string{
		homeDir,
		filepath.Join(configDir, constants.Mangal),
	}, nil
}
