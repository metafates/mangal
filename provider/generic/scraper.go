package generic

import (
	"github.com/gocolly/colly/v2"
	"github.com/metafates/mangal/source"
)

// Scraper is a generic scraper downloads html pages and parses them
type Scraper struct {
	mangasCollector   *colly.Collector
	chaptersCollector *colly.Collector
	pagesCollector    *colly.Collector

	mangas   map[string][]*source.Manga
	chapters map[string][]*source.Chapter
	pages    map[string][]*source.Page

	config *Configuration
}

// Name of the scraper
func (s *Scraper) Name() string {
	return s.config.Name
}

// ID of the scraper
func (s *Scraper) ID() string {
	return s.config.ID()
}
