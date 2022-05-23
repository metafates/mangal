package main

import (
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Scrapers    []*Scraper
	Fullscreen  bool
	Path        string
	RandomDelay time.Duration
}

type _tempConfig struct {
	Use         []string
	Path        string
	Fullscreen  bool
	Sources     map[string]Source
	RandomDelay int `toml:"random_delay_ms"`
}

func GetConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, AppName, "config.toml"), nil
}

var DefaultConfig = Config{
	Scrapers:    []*Scraper{DefaultScraper},
	Fullscreen:  true,
	Path:        ".",
	RandomDelay: 700 * time.Millisecond,
}

func GetConfig() *Config {
	// TODO
	configPath, err := GetConfigPath()

	if err != nil {
		return &DefaultConfig
	}

	configExists, err := Afero.Exists(configPath)
	if err != nil || !configExists {
		return &DefaultConfig
	}

	contents, err := Afero.ReadFile(configPath)
	if err != nil {
		return &DefaultConfig
	}

	config, err := ParseConfig(string(contents))
	if err != nil {
		return &DefaultConfig
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
