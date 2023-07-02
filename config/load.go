package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
	"github.com/mangalorg/mangal/fs"
	"github.com/mangalorg/mangal/path"
	"path/filepath"
)

var instance = koanf.NewWithConf(koanf.Conf{
	Delim:       ".",
	StrictMerge: true,
})

func Load() error {
	for _, f := range Fields {
		if err := instance.Set(f.Key(), f.Default()); err != nil {
			return err
		}
	}

	const configFilename = "mangal.yaml"
	configFilepath := filepath.Join(path.ConfigDir(), configFilename)
	exists, err := fs.Afero.Exists(configFilepath)
	if err != nil {
		return err
	}

	if exists {
		file, err := fs.Afero.ReadFile(configFilepath)
		if err != nil {
			return err
		}

		if err := instance.Load(rawbytes.Provider(file), yaml.Parser()); err != nil {
			return err
		}
	}

	for _, f := range Fields {
		value, err := f.Transform()
		if err != nil {
			return err
		}

		if err := instance.Set(f.Key(), value); err != nil {
			return err
		}

		if err := f.Init(); err != nil {
			return err
		}
	}

	return nil
}
