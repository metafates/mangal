package scraper

import "testing"

func TestValidateSource(t *testing.T) {
	// create a valid test source
	source := Source{
		Name:             "Test",
		Base:             "https://example.com",
		SearchTemplate:   "https://example.com/search?q=%s",
		MangaAnchor:      "a",
		MangaTitle:       "title",
		ChapterAnchor:    "a",
		ChapterTitle:     "title",
		ReaderPage:       "a",
		RandomDelayMs:    0,
		ChaptersReversed: false,
	}

	// validate source
	err := ValidateSource(&source)
	if err != nil {
		t.Error(err)
	}

	// create an invalid test source
	source.Base = ""

	// validate source
	err = ValidateSource(&source)
	if err == nil {
		t.Error("expected error but got nil")
	}

	// set search template to invalid url
	source.SearchTemplate = "https://example.com/search?q="

	// validate source
	err = ValidateSource(&source)
	if err == nil {
		t.Error("expected error but got nil")
	}
}
