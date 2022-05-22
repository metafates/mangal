package config

import (
	"github.com/metafates/mangai/shared"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/afero"
)

func getPath() (string, error) {
	configDir, err := os.UserConfigDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, strings.ToLower(shared.Mangai), "config.toml"), nil
}

// Get gets config struct. It tries to read file for the first time
// All other calls are using cached config
var Get = initConfigSupplier()

// initConfigSupplier returns closure function with cached config
func initConfigSupplier() func() Config {
	var cached Config
	has := false

	return func() Config {
		if has {
			return cached
		}

		path, err := getPath()

		if exists, _ := afero.Exists(shared.AferoBackend, path); !exists || err != nil {
			return createDefault()
		}

		contents, err := shared.AferoFS.ReadFile(path)

		if err != nil {
			cached = createDefault()
			has = true
			return cached
		}

		var conf Config
		_, err = toml.Decode(string(contents), &conf)

		if err != nil {
			cached = createDefault()
			has = true

			return cached
		}

		conf.setUsedSources()
		cached = conf
		has = true
		return cached
	}
}
