package main

import (
	"errors"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"os"
	"path/filepath"
	"strings"
)

// FormatType is type of format used for output
type FormatType string

const (
	PDF   FormatType = "pdf"
	CBZ   FormatType = "cbz"
	Zip   FormatType = "zip"
	Plain FormatType = "plain"
	Epub  FormatType = "epub"
)

type UI struct {
	Fullscreen          bool
	Prompt              string
	Title               string
	Placeholder         string
	Mark                string
	ChapterNameTemplate string `toml:"chapter_name_template"`
}

type Config struct {
	Scrapers []*Scraper
	Format   FormatType
	UI       UI
	Anilist  struct {
		Client         *AnilistClient
		Enabled        bool
		MarkDownloaded bool
	}
	UseCustomReader     bool
	CustomReader        string
	Path                string
	CacheImages         bool
	ChapterNameTemplate string
}

// GetConfigPath returns path to config file
func GetConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, strings.ToLower(Mangal), "config.toml"), nil
}

// DefaultConfig makes default config
func DefaultConfig() *Config {
	conf, _ := ParseConfig(DefaultConfigBytes)
	return conf
}

// UserConfig is a global variable that stores user config
var UserConfig *Config

// GetConfig returns user config or default config if it doesn't exist
// If path is empty string then default config will be returned
func GetConfig(path string) *Config {
	var (
		configPath string
		err        error
	)

	// If path is empty string then default config will be used
	if path == "" {
		configPath, err = GetConfigPath()
	} else {
		configPath = path
	}

	if err != nil {
		return DefaultConfig()
	}

	// If config file doesn't exist then default config will be used
	configExists, err := Afero.Exists(configPath)
	if err != nil || !configExists {
		return DefaultConfig()
	}

	// Read config file
	contents, err := Afero.ReadFile(configPath)
	if err != nil {
		return DefaultConfig()
	}

	// Parse config
	config, err := ParseConfig(contents)
	if err != nil {
		return DefaultConfig()
	}

	return config
}

// ParseConfig parses config from given string
func ParseConfig(configString []byte) (*Config, error) {
	// tempConfig is a temporary config that will be used to store parsed config
	type tempConfig struct {
		Use                 []string
		Format              string
		UI                  UI     `toml:"ui"`
		UseCustomReader     bool   `toml:"use_custom_reader"`
		CustomReader        string `toml:"custom_reader"`
		Path                string `toml:"download_path"`
		CacheImages         bool   `toml:"cache_images"`
		Sources             map[string]Source
		ChapterNameTemplate string `toml:"chapter_name_template"`
		Anilist             struct {
			Enabled        bool   `toml:"enabled"`
			ID             string `toml:"id"`
			Secret         string `toml:"secret"`
			MarkDownloaded bool   `toml:"mark_downloaded"`
		}
	}

	var (
		tempConf tempConfig
		conf     Config
	)
	err := toml.Unmarshal(configString, &tempConf)

	if err != nil {
		return nil, err
	}

	conf.CacheImages = tempConf.CacheImages

	// Convert sources listed in tempConfig to Scrapers
	for sourceName, source := range tempConf.Sources {
		// If source is not listed in Use then skip it
		if !Contains[string](tempConf.Use, sourceName) {
			continue
		}

		// Create scraper
		source.Name = sourceName
		scraper := MakeSourceScraper(&source)

		if !conf.CacheImages {
			scraper.FilesCollector.CacheDir = ""
		}

		conf.Scrapers = append(conf.Scrapers, scraper)
	}

	if tempConf.UI.ChapterNameTemplate == "" {
		tempConf.UI.ChapterNameTemplate = "[%d] %s"
	}

	conf.UI = tempConf.UI

	conf.Path = tempConf.Path

	// Default format is pdf
	conf.Format = IfElse(tempConf.Format == "", PDF, FormatType(tempConf.Format))

	conf.UseCustomReader = tempConf.UseCustomReader
	conf.CustomReader = tempConf.CustomReader

	if template := tempConf.ChapterNameTemplate; template != "" {
		conf.ChapterNameTemplate = template
	} else {
		conf.ChapterNameTemplate = "[%0d] %s"
	}

	if tempConf.Anilist.Enabled {
		id, secret := tempConf.Anilist.ID, tempConf.Anilist.Secret
		conf.Anilist.Client, err = NewAnilistClient(id, secret)

		if err != nil {
			return nil, err
		}

		conf.Anilist.Enabled = true
		conf.Anilist.MarkDownloaded = tempConf.Anilist.MarkDownloaded
	}

	return &conf, err
}

// ValidateConfig checks if config is valid and returns error if it is not
func ValidateConfig(config *Config) error {
	// check if any source is used
	if len(config.Scrapers) == 0 {
		return errors.New("no manga sources listed")
	}

	// check if chapter name template is valid
	// chapter name template should contain %d or %s placeholder
	validateChapterNameTemplate := func(template string) error {
		if !strings.Contains(template, "%d") &&
			!strings.Contains(template, "%s") &&
			!strings.Contains(template, "%0d") {
			return errors.New("chapter name template should contain at least one %d, %0d or %s placeholder")
		}
		return nil
	}

	err := validateChapterNameTemplate(config.ChapterNameTemplate)
	if err != nil {
		return err
	}

	err = validateChapterNameTemplate(config.UI.ChapterNameTemplate)
	if err != nil {
		return err
	}

	// check if anilist is properly configured
	if config.Anilist.Enabled {
		if config.Anilist.Client.ID == "" || config.Anilist.Client.Secret == "" {
			return errors.New("anilist is enabled but id or secret is not set")
		}
	}

	// check if custom reader is properly configured
	if config.UseCustomReader && config.CustomReader == "" {
		return errors.New("use_custom_reader is set to true but reader isn't specified")
	}

	// Check if format is valid
	if !Contains(AvailableFormats, config.Format) {
		msg := fmt.Sprintf(
			`unknown format '%s'
type %s to show available formats`,
			string(config.Format),
			accentStyle.Render(strings.ToLower(Mangal)+" formats"),
		)
		return errors.New(msg)
	}

	// Check if scrapers are valid
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
