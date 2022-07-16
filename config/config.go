package config

import (
	"errors"
	"fmt"
	"github.com/metafates/mangal/common"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/scraper"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/util"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/afero"
	"golang.org/x/exp/slices"
	"os"
	"strings"
)

type Config struct {
	Scrapers   []*scraper.Scraper
	Formats    *FormatsConfig
	UI         *UIConfig
	Downloader *DownloaderConfig
	Anilist    struct {
		Client         *scraper.AnilistClient
		Enabled        bool
		MarkDownloaded bool
	}
	Reader        *ReaderConfig
	HistoryMode   bool
	IncognitoMode bool
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
		configPath, err = util.UserConfigFile()
	} else {
		configPath = path
	}

	if err != nil {
		return DefaultConfig()
	}

	// If config file doesn't exist then default config will be used

	configExists, err := afero.Exists(filesystem.Get(), configPath)
	if err != nil || !configExists {
		return DefaultConfig()
	}

	// Read config file
	contents, err := afero.ReadFile(filesystem.Get(), configPath)
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
		Formats         FormatsConfig `toml:"formats"`
		UI              UIConfig      `toml:"ui"`
		UseCustomReader bool          `toml:"use_custom_reader"`
		CustomReader    string        `toml:"custom_reader"`
		Sources         []*scraper.Source
		Downloader      DownloaderConfig
		Reader          ReaderConfig
		Anilist         struct {
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

	conf.Downloader = &tempConf.Downloader
	if conf.Downloader.ChapterNameTemplate == "" {
		conf.Downloader.ChapterNameTemplate = "[%d] %s"
	}

	if strings.Contains(conf.Downloader.Path, "$HOME") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		conf.Downloader.Path = strings.ReplaceAll(conf.Downloader.Path, "$HOME", home)
	}

	// Convert sources listed in tempConfig to Scrapers
	for _, source := range tempConf.Sources {
		if source.Enabled {
			// Create scraper for this source
			s := scraper.MakeSourceScraper(source)
			conf.Scrapers = append(conf.Scrapers, s)
		}
	}

	conf.UI = &tempConf.UI
	if tempConf.UI.ChapterNameTemplate == "" {
		tempConf.UI.ChapterNameTemplate = "[%d] %s"
	}

	// Default format is pdf
	conf.Formats = &tempConf.Formats
	if format, ok := os.LookupEnv(common.EnvDefaultFormat); ok {
		conf.Formats.Default = common.FormatType(format)
	} else if conf.Formats.Default == "" {
		conf.Formats.Default = common.PDF
	}

	conf.Reader = &tempConf.Reader
	conf.HistoryMode = false
	conf.IncognitoMode = false

	if customReader := os.Getenv(common.EnvCustomReader); customReader != "" {
		conf.Reader.UseCustomReader = true
		conf.Reader.CustomReader = customReader
	}

	if tempConf.Anilist.Enabled {
		id, secret := tempConf.Anilist.ID, tempConf.Anilist.Secret
		conf.Anilist.Client, err = scraper.NewAnilistClient(id, secret)

		if err != nil {
			return nil, err
		}

		conf.Anilist.Enabled = true
		conf.Anilist.MarkDownloaded = tempConf.Anilist.MarkDownloaded
	}

	if v, ok := os.LookupEnv(common.EnvDownloadPath); ok {
		conf.Downloader.Path = v
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

	err := validateChapterNameTemplate(config.Downloader.ChapterNameTemplate)
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
	if config.Reader.UseCustomReader && config.Reader.CustomReader == "" {
		return errors.New("use_custom_reader is set to true but reader isn't specified")
	}

	// Check if format is valid
	if !slices.Contains(common.AvailableFormats, config.Formats.Default) {
		msg := fmt.Sprintf(
			`unknown format '%s'
type %s to show available formats`,
			string(config.Formats.Default),
			style.Accent.Render(strings.ToLower(common.Mangal)+" formats"),
		)
		return errors.New(msg)
	}

	// Check if scrapers are valid
	for _, s := range config.Scrapers {
		if s.Source == nil {
			return errors.New("internal error: scraper source is nil")
		}
		if err := scraper.ValidateSource(s.Source); err != nil {
			return err
		}
	}

	return nil
}
