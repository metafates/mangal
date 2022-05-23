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

var DefaultConfig = Config{
	Scrapers:   []*Scraper{DefaultScraper},
	Fullscreen: true,
	Path:       ".",
}

func GetConfig() *Config {
	// TODO
	return nil
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
