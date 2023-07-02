package config

import (
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var instance = koanf.New(".")

func Load() error {
	for _, f := range Fields {
		err := instance.Set(f.Key(), f.Default())
		if err != nil {
			return err
		}
	}

	// TODO: handle err
	_ = instance.Load(file.Provider("mangal.toml"), toml.Parser())
	return nil
}
