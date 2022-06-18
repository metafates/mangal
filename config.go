package main

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
	"strings"
)

type FormatType string

const (
	PDF   FormatType = "pdf"
	CBZ   FormatType = "cbz"
	Zip   FormatType = "zip"
	Plain FormatType = "plain"
)

type Config struct {
	Scrapers        []*Scraper
	Format          FormatType
	Fullscreen      bool
	Prompt          string
	Title           string
	Placeholder     string
	Mark            string
	UseCustomReader bool
	CustomReader    string
	Path            string
	CacheImages     bool
}

type _tempConfig struct {
	Use             []string
	Format          string
	Fullscreen      bool
	Prompt          string
	Placeholder     string
	Title           string
	Mark            string
	UseCustomReader bool   `toml:"use_custom_reader"`
	CustomReader    string `toml:"custom_reader"`
	Path            string `toml:"download_path"`
	CacheImages     bool   `toml:"cache_images"`
	Sources         map[string]Source
}

func GetConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, strings.ToLower(AppName), "config.toml"), nil
}

func DefaultConfig() *Config {
	conf, _ := ParseConfig(string(DefaultConfigBytes))
	return conf
}

var UserConfig *Config

var DefaultConfigBytes = []byte(`# Which sources to use. You can use several sources, it won't affect perfomance'
use = ['manganelo']

# Available options: pdf, cbz, zip, plain (just images)
format = "pdf"

# If false, then OS default pdf reader will be used
use_custom_reader = false
custom_reader = "zathura"

# Custom download path, can be either relative (to the pwd) or absolute
download_path = '.'

# Fullscreen mode
fullscreen = true

# Input prompt icon
prompt = "🔍"

# Input placeholder
placeholder = "What shall we look for?"

# Selected chapter mark
mark = "▼"

# Search window title
title = "Mangal"

# Add images to cache
# If set to true mangal could crash when trying to redownload something really quickly
# Usually happens on slow machines
cache_images = false

[sources]
    [sources.manganelo]
    # Base url
    base = 'https://ww5.manganelo.tv'

    # Search endpoint. Put %s where the query should be
    search = 'https://ww5.manganelo.tv/search/%s'

    # Selector of entry anchor (<a></a>) on search page
    manga_anchor = '.search-story-item a.item-title'

    # Selector of entry title on search page
    manga_title = '.search-story-item a.item-title'

    # Manga chapters anchors selector
    chapter_anchor = 'li.a-h a.chapter-name'

    # Manga chapters titles selector
    chapter_title = 'li.a-h a.chapter-name'

    # Reader page images selector
    reader_page = '.container-chapter-reader img'
    
    # Random delay between requests
    random_delay_ms = 500 # ms
    
    # Are chapters listed in reversed order on that source?
    # reversed order -> from newest chapter to oldest
    reversed_chapters_order = true
`)

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
		return DefaultConfig()
	}

	configExists, err := Afero.Exists(configPath)
	if err != nil || !configExists {
		return DefaultConfig()
	}

	contents, err := Afero.ReadFile(configPath)
	if err != nil {
		return DefaultConfig()
	}

	config, err := ParseConfig(string(contents))
	if err != nil {
		return DefaultConfig()
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

	conf.CacheImages = tempConf.CacheImages
	// Convert sources to scrapers
	for sourceName, source := range tempConf.Sources {
		if !Contains[string](tempConf.Use, sourceName) {
			continue
		}

		source.Name = sourceName
		scraper := MakeSourceScraper(source)

		if !conf.CacheImages {
			scraper.FilesCollector.CacheDir = ""
		}

		conf.Scrapers = append(conf.Scrapers, scraper)
	}

	conf.Fullscreen = tempConf.Fullscreen
	conf.Mark = tempConf.Mark
	conf.Prompt = tempConf.Prompt
	conf.Placeholder = tempConf.Placeholder
	conf.Path = tempConf.Path
	conf.Title = tempConf.Title
	// Default format is pdf
	conf.Format = IfElse(tempConf.Format == "", PDF, FormatType(tempConf.Format))

	conf.UseCustomReader = tempConf.UseCustomReader
	conf.CustomReader = tempConf.CustomReader

	return &conf, err
}

func ValidateConfig(config *Config) error {
	if config.UseCustomReader && config.CustomReader == "" {
		return errors.New("use_custom_reader is set to true but reader isn't specified")
	}

	if availableFormats := []FormatType{PDF, CBZ, Plain, Zip}; !Contains(availableFormats, config.Format) {
		// I don't like this...
		// TODO: Refactor
		return errors.New(fmt.Sprintf(`unknown format '%s'
available formats: %s`, string(config.Format), strings.Join(Map(availableFormats, func(f FormatType) string {
			return string(f)
		}), ", ")))
	}

	for _, scraper := range config.Scrapers {
		if scraper.Source == nil {
			return errors.New("internal error: scraper source is nil")
		}
		if err := ValidateSource(scraper.Source); err != nil {
			return err
		}
	}

	return nil
}
