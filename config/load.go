package config

import (
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
	"github.com/mangalorg/mangal/fs"
	"github.com/mangalorg/mangal/meta"
	"github.com/mangalorg/mangal/path"
)

var instance = koanf.NewWithConf(koanf.Conf{
	Delim:       ".",
	StrictMerge: true,
})

const configFilename = "mangal.toml"

func Load() error {
	for _, f := range Fields {
		if err := instance.Set(f.Key(), f.Default()); err != nil {
			return err
		}
	}

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

		if err := instance.Load(rawbytes.Provider(file), toml.Parser()); err != nil {
			return err
		}
	}

	// TODO: fix this
	prefix := strings.ToUpper(meta.AppName)
	if err := instance.Load(env.Provider(prefix, ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, prefix))
	}), nil); err != nil {
		return err
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
