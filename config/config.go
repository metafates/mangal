package config

import "github.com/metafates/mangai/api/scraper"

type Config struct {
	Sources       map[string]scraper.Source
	Use           []string
	Path          string
	Fullscreen    bool
	UsedSources   []*scraper.Source
	AsyncDownload bool `toml:"async_download"`
}

func (c *Config) setUsedSources() {
	for _, name := range c.Use {
		usedSource, found := c.Sources[name]
		if found {
			c.UsedSources = append(c.UsedSources, &usedSource)
		}
	}
}
