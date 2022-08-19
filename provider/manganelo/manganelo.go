package manganelo

import (
	"github.com/gocolly/colly"
	"github.com/metafates/mangal/source"
)

const (
	Name = "Manganelo"
	ID   = Name + " built-in"
)

type Manganelo struct {
	mangasCollector   *colly.Collector
	chaptersCollector *colly.Collector
	pagesCollector    *colly.Collector

	mangas   map[string][]*source.Manga
	chapters map[string][]*source.Chapter
	pages    map[string][]*source.Page
}

func (*Manganelo) Name() string {
	return Name
}

func (*Manganelo) ID() string {
	return ID
}
