package config

import (
	"path/filepath"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/mangalorg/mangal/fs"
	"github.com/mangalorg/mangal/path"
)

func Set(key string, value any) error {
	return instance.Set(key, value)
}

func Get(key string) any {
	return instance.Get(key)
}

func Exists(key string) bool {
	return instance.Exists(key)
}

func Keys() []string {
	return instance.Keys()
}

func Write() error {
	marshalled, err := instance.Marshal(toml.Parser())
	if err != nil {
		return err
	}

	return fs.Afero.WriteFile(
		filepath.Join(path.ConfigDir(), configFilename),
		marshalled,
		0655,
	)
}
