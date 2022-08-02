package config

import (
	"github.com/metafates/mangal/filesystem"
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
	err := setDefaults()
	if err != nil {
		return err
	}

	err = setPaths()
	if err != nil {
		return err
	}

	err = viper.ReadInConfig()

	switch err.(type) {
	case viper.ConfigFileNotFoundError:
		// Use defaults then
		return nil
	default:
		return err
	}
}

func setName() {
	viper.SetConfigName("mangal")
	viper.SetConfigType("toml")
}

func setFs() {
	viper.SetFs(filesystem.Get())
}

// setPaths sets the paths to the config files
func setPaths() error {
	paths, err := Paths()
	if err != nil {
		return err
	}

	for _, path := range paths {
		viper.AddConfigPath(path)
	}

	return nil
}

// setEnvs sets the environment variables
func setEnvs() {
	viper.SetEnvPrefix("mangal")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	for _, env := range envFields {
		viper.MustBindEnv(env)
	}
}

// setDefaults sets the default values
func setDefaults() error {
	viper.SetTypeByDefaultValue(true)

	viper.SetDefault(DownloaderPath, ".")
	viper.SetDefault(DownloaderChapterNameTemplate, "[%0d] %s")
	viper.SetDefault(FormatsDefault, "pdf")

	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	viper.SetDefault(SourcesPath, filepath.Join(configDir, "mangal", "sources"))

	return nil
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
		".",
		homeDir,
		filepath.Join(configDir, "mangal"),
	}, nil
}
