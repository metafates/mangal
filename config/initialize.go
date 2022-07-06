package config

import (
	"errors"
	"fmt"
	"github.com/metafates/mangal/filesystem"
	"github.com/spf13/afero"
	"log"
	"path"
	"path/filepath"
)

// Initialize initializes the config file
// If the given string is empty, it will use the default config file
func Initialize(configPath string, validate bool) {
	if configPath != "" {
		// check if config is a TOML file
		if filepath.Ext(configPath) != ".toml" {
			log.Fatal("config file must be a TOML file")
		}

		// check if config file exists
		exists, err := afero.Exists(filesystem.Get(), configPath)

		if err != nil {
			log.Fatal(errors.New("access to config file denied"))
		}

		// if config file doesn't exist raise error
		configPath = path.Clean(configPath)
		if !exists {
			log.Fatal(errors.New(fmt.Sprintf("config at path %s doesn't exist", configPath)))
		}

		UserConfig = GetConfig(configPath)
	} else {
		// if config path is empty, use default config file
		UserConfig = GetConfig("")
	}

	if !validate {
		return
	}

	// check if config file is valid
	err := ValidateConfig(UserConfig)

	if err != nil {
		log.Fatal(err)
	}
}
