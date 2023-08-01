package config

import (
	"path/filepath"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/mangalorg/mangal/fs"
	"github.com/mangalorg/mangal/path"
)

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
