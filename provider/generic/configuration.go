package generic

import (
	"github.com/PuerkitoBio/goquery"
	"time"
)

// Extractor is responsible for finding specified elements by selector and extracting required data from them
type Extractor struct {
	// Selector CSS selector
	Selector string
	// Name function to get name from element found by selector.
	Name func(*goquery.Selection) string
	// URL function to get URL from element found by selector.
	URL func(*goquery.Selection) string
	// Volume function to get volume from element found by selector. Used by chapters extractor
	Volume func(*goquery.Selection) string
	// Cover function to get cover from element found by selector. Used by manga extractor
	Cover func(*goquery.Selection) string
}

// Configuration is a generic scraper configuration that defines behavior of the scraper
type Configuration struct {
	// Name of the scraper
	Name string
	// Delay between requests
	Delay time.Duration
	// Parallelism of the scraper
	Parallelism uint8

	// ReverseChapters if true, chapters will be shown in reverse order
	ReverseChapters bool

	// BaseURL of the source
	BaseURL string
	// GenerateSearchURL function to create search URL from the query.
	// E.g. "one piece" -> "https://manganelo.com/search/story/one%20piece"
	GenerateSearchURL func(query string) string

	// MangaExtractor is responsible for finding manga elements and extracting required data from them
	MangaExtractor,
	// ChapterExtractor is responsible for finding chapter elements and extracting required data from them
	ChapterExtractor,
	// PageExtractor is responsible for finding page elements and extracting required data from them
	PageExtractor *Extractor
}

func (c *Configuration) ID() string {
	return c.Name + " built-in"
}
