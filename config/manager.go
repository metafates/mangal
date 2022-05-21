package config

import (
	"github.com/metafates/mangai/api/scraper"
	"github.com/metafates/mangai/shared"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/spf13/afero"
)

type Config struct {
	Sources     map[string]scraper.Source
	Use         []string
	Path        string
	Fullscreen  bool
	UsedSources []*scraper.Source
}

func (c *Config) setUsedSources() {
	for _, name := range c.Use {
		usedSource, found := c.Sources[name]
		if found {
			c.UsedSources = append(c.UsedSources, &usedSource)
		}
	}
}

func getPath() (string, error) {
	configDir, err := os.UserConfigDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "mangai", "config.toml"), nil
}

func createDefault() Config {
	var conf Config

	tomlString := `
    use = ['manganelo']
    path = '.'
    fullscreen = true
    [sources]
        [sources.manganelo]
        base = 'https://ww5.manganelo.tv'
        search = 'https://ww5.manganelo.tv/search/%s'
        manga_anchor = '.search-story-item a.item-title' manga_title = '.search-story-item a.item-title'
        chapter_anchor = 'li.a-h a.chapter-name'
        chapter_title = 'li.a-h a.chapter-name'
        chapter_panels = '.container-chapter-reader img'
    `

	if _, err := toml.Decode(tomlString, &conf); err != nil {
		log.Fatal("Unexpected error while loading default config")
	}

	conf.setUsedSources()
	return conf
}

// Get gets config struct. It tries to read file for the first time
// All other calls are using cached config
var Get = initConfigSupplier()

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
