package generic

import (
	"github.com/gocolly/colly"
	"github.com/metafates/mangal/source"
)

type Scraper struct {
	mangasCollector   *colly.Collector
	chaptersCollector *colly.Collector
	pagesCollector    *colly.Collector

	mangas   map[string][]*source.Manga
	chapters map[string][]*source.Chapter
	pages    map[string][]*source.Page

	config *Configuration
}

func (s *Scraper) Name() string {
	return s.config.Name
}

func (s *Scraper) ID() string {
	return s.config.ID()
}
