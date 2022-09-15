package mangakakalot

import (
	"github.com/gocolly/colly"
	"github.com/metafates/mangal/source"
)

const (
	Name = "Mangakakalot"
	ID   = Name + " built-in"
)

type Mangakakalot struct {
	mangasCollector   *colly.Collector
	chaptersCollector *colly.Collector
	pagesCollector    *colly.Collector

	mangas   map[string][]*source.Manga
	chapters map[string][]*source.Chapter
	pages    map[string][]*source.Page
}

func (*Mangakakalot) Name() string {
	return Name
}

func (*Mangakakalot) ID() string {
	return ID
}
