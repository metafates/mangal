package config

import (
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/where"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var EnvKeyReplacer = strings.NewReplacer(".", "_")

// Setup initializes the configuration
func Setup() error {
	setName()
	setFs()
	setEnvs()
	setDefaults()
	setPaths()

	err := viper.ReadInConfig()

	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			// Use defaults then
			return nil
		default:
			return err
		}
	}

	resolveAliases()
	return nil
}

func setName() {
	viper.SetConfigName(constant.Mangal)
	viper.SetConfigType("toml")
}

func setFs() {
	viper.SetFs(filesystem.Api())
}

// setPaths sets the paths to the config files
func setPaths() {
	viper.AddConfigPath(where.Config())
}

// setEnvs sets the environment variables
func setEnvs() {
	viper.SetEnvPrefix(constant.Mangal)
	viper.SetEnvKeyReplacer(EnvKeyReplacer)

	for _, env := range EnvExposed {
		viper.MustBindEnv(env)
	}
}

// setDefaults sets the default values
func setDefaults() {
	viper.SetTypeByDefaultValue(true)

	for name, field := range DefaultValues {
		viper.SetDefault(name, field.Value)
	}
}

// resolveAliases resolves the aliases for the paths
func resolveAliases() {
	home := lo.Must(os.UserHomeDir())
	path := viper.GetString(constant.DownloaderPath)

	path = strings.ReplaceAll(path, "$HOME", home)
	path = strings.ReplaceAll(path, "~", home)

	viper.Set(constant.DownloaderPath, path)
}
