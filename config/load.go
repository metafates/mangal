package config

import (
	"github.com/mangalorg/mangal/fs"
	"github.com/mangalorg/mangal/meta"
	"github.com/mangalorg/mangal/path"
	"github.com/spf13/viper"
)

func Load() error {
	viper.SetConfigName(meta.AppName)
	viper.SetConfigType("toml")
	viper.SetFs(fs.Afero.Fs)
	viper.AddConfigPath(path.ConfigDir())
	viper.KeyDelimiter(".")
	viper.SetTypeByDefaultValue(true)

	for _, f := range Fields {
		viper.SetDefault(f.Key(), f.Default())
	}

	err := viper.ReadInConfig()
	switch err := err.(type) {
	case viper.ConfigFileNotFoundError:
		return nil
	default:
		return err
	}
}
