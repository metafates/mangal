package scraper

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// Source is a source used to scrape manga
type Source struct {
	Enabled          bool
	Name             string
	Base             string
	SearchTemplate   string `toml:"search"`
	MangaAnchor      string `toml:"manga_anchor"`
	MangaTitle       string `toml:"manga_title"`
	ChapterAnchor    string `toml:"chapter_anchor"`
	ChapterTitle     string `toml:"chapter_title"`
	ReaderPage       string `toml:"reader_page"`
	RandomDelayMs    int    `toml:"random_delay_ms"`
	ChaptersReversed bool   `toml:"reversed_chapters_order"`
	WhitespaceEscape string `toml:"whitespace_escape"`
}

// ValidateSource validates source and returns error if it's invalid
func ValidateSource(source *Source) error {
	// Check if given string is a valid URL
	var isValidURI = func(uri string) bool {
		_, err := url.ParseRequestURI(strings.Replace(uri, "%s", "", 1))
		return err == nil
	}

	type test struct {
		condition    bool
		errorMessage string
	}

	// Tests for source validity
	tests := []test{
		{source.Base != "", "base url is missing"},
		{isValidURI(source.Base), "base url is not a valid url"},
		{source.SearchTemplate != "", "search template is empty"},
		{strings.Contains(source.SearchTemplate, "%s"), "search template is missing query template ('%s')"},
		{isValidURI(source.SearchTemplate), "search template is not a valid url"},
		{source.MangaTitle != "", "manga title selector is empty"},
		{source.MangaAnchor != "", "manga anchor selector is empty"},
		{source.ChapterAnchor != "", "chapter anchor selector is empty"},
		{source.ChapterTitle != "", "chapter title selector is empty"},
		{source.ReaderPage != "", "reader page selector is empty"},
	}

	// Run tests
	for _, t := range tests {
		if !t.condition {
			msg := fmt.Sprintf("[%s] %s", source.Name, t.errorMessage)
			return errors.New(msg)
		}
	}

	return nil
}
