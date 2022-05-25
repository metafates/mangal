package main

import (
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
)

type Config struct {
	Scrapers   []*Scraper
	Fullscreen bool
	Path       string
}

type _tempConfig struct {
	Use        []string
	Path       string
	Fullscreen bool
	Sources    map[string]Source
}

func GetConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, AppName, "config.toml"), nil
}

var DefaultConfig = &Config{
	Scrapers:   []*Scraper{DefaultScraper},
	Fullscreen: true,
	Path:       ".",
}

var UserConfig *Config

// GetConfig from given path. If path is empty string default config path is used
func GetConfig(path string) *Config {
	var (
		configPath string
		err        error
	)

	if path == "" {
		configPath, err = GetConfigPath()
	} else {
		configPath = path
	}

	if err != nil {
		return DefaultConfig
	}

	configExists, err := Afero.Exists(configPath)
	if err != nil || !configExists {
		return DefaultConfig
	}

	contents, err := Afero.ReadFile(configPath)
	if err != nil {
		return DefaultConfig
	}

	config, err := ParseConfig(string(contents))
	if err != nil {
		return DefaultConfig
	}

	return config
}

func ParseConfig(configString string) (*Config, error) {
	var (
		tempConf _tempConfig
		conf     Config
	)
	_, err := toml.Decode(configString, &tempConf)

	if err != nil {
		return nil, err
	}

	// Handle sources
	for sourceName, source := range tempConf.Sources {
		if !Contains[string](tempConf.Use, sourceName) {
			continue
		}

		scraper := MakeSourceScraper(&source)
		conf.Scrapers = append(conf.Scrapers, scraper)
	}

	conf.Fullscreen = tempConf.Fullscreen
	conf.Path = tempConf.Path

	return &conf, err
}
