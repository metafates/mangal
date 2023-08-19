package config

import (
	"fmt"

	"github.com/mangalorg/mangal/afs"
	"github.com/mangalorg/mangal/meta"
	"github.com/mangalorg/mangal/path"
	"github.com/spf13/viper"
)

func Load() error {
	viper.SetConfigName(meta.AppName)
	viper.SetConfigType("toml")
	viper.SetFs(afs.Afero.Fs)
	viper.AddConfigPath(path.ConfigDir())
	viper.KeyDelimiter(".")
	viper.SetTypeByDefaultValue(false)

	for _, field := range Fields {
		marshalled, err := field.Marshal(field.Default)
		if err != nil {
			return err
		}
		viper.SetDefault(field.Key, marshalled)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	for _, field := range Fields {
		unmarshalled, err := field.Unmarshal(viper.Get(field.Key))
		if err != nil {
			return err
		}

		if err := field.Validate(unmarshalled); err != nil {
			return fmt.Errorf("%s: %s", field.Key, err)
		}

		if err := field.SetValue(unmarshalled); err != nil {
			return err
		}
	}

	return nil
}
