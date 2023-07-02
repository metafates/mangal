package config

import (
	"bytes"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/mangalorg/mangal/fs"
	"github.com/mangalorg/mangal/path"
	"path/filepath"
)

func Write() error {
	marshalled, err := instance.Marshal(yaml.Parser())
	if err != nil {
		return err
	}

	if !bytes.HasPrefix(marshalled, []byte("---\n")) {
		marshalled = append([]byte("---\n"), marshalled...)
	}

	return fs.Afero.WriteFile(
		filepath.Join(path.ConfigDir(), configFilename),
		marshalled,
		0655,
	)
}
