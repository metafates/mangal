package generic

import (
	"github.com/PuerkitoBio/goquery"
	"time"
)

type Extractor struct {
	Selector string
	Name     func(*goquery.Selection) string
	URL      func(*goquery.Selection) string
	Volume   func(*goquery.Selection) string
	Cover    func(*goquery.Selection) string
}

type Configuration struct {
	Name        string
	Delay       time.Duration
	Parallelism uint8

	ReverseChapters bool

	BaseURL           string
	GenerateSearchURL func(query string) string

	MangaExtractor,
	ChapterExtractor,
	PageExtractor *Extractor
}

func (c *Configuration) ID() string {
	return c.Name + " built-in"
}
