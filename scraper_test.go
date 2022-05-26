package main

import (
	"testing"
)

const testQuery = "Death Note"

func TestMakeSourceScraper(t *testing.T) {
	source := &DefaultSource

	scraper := MakeSourceScraper(source)

	if scraper == nil {
		t.Failed()
	}

	if scraper.Source == nil || scraper.Source != source {
		t.Failed()
	}
}

func TestScraper_SearchManga(t *testing.T) {
	scraper := MakeSourceScraper(&DefaultSource)

	manga, err := scraper.SearchManga(testQuery)

	if err != nil {
		t.Failed()
	}

	if len(manga) == 0 {
		t.Fail()
	}

	if !IsUnique(manga) {
		t.Failed()
	}

	for _, m := range manga {
		if m.Scraper != scraper {
			t.Fail()
		}
	}
}

func TestScraper_GetChapters(t *testing.T) {
	scraper := MakeSourceScraper(&DefaultSource)

	manga, _ := scraper.SearchManga(testQuery)
	anyManga := manga[0]

	chapters, err := scraper.GetChapters(anyManga)

	if err != nil {
		t.Failed()
	}

	if len(chapters) == 0 {
		t.Failed()
	}

	if !IsUnique(chapters) {
		t.Failed()
	}

	for _, chapter := range chapters {

		if chapter.Relation != anyManga {
			t.Fail()
		}
	}
}

func TestScraper_GetPages(t *testing.T) {
	scraper := MakeSourceScraper(&DefaultSource)

	manga, _ := scraper.SearchManga(testQuery)
	anyManga := manga[0]

	chapters, _ := scraper.GetChapters(anyManga)
	anyChapter := chapters[0]

	pages, err := scraper.GetPages(anyChapter)

	if err != nil {
		t.Failed()
	}

	if len(pages) == 0 {
		t.Failed()
	}

	if !IsUnique(pages) {
		t.Failed()
	}

	for _, page := range pages {
		if page.Relation != anyChapter {
			t.Fail()
		}
	}
}

func TestScraper_GetFile(t *testing.T) {
	scraper := MakeSourceScraper(&DefaultSource)

	manga, _ := scraper.SearchManga(testQuery)
	anyManga := manga[0]

	chapters, _ := scraper.GetChapters(anyManga)
	anyChapter := chapters[0]

	pages, _ := scraper.GetPages(anyChapter)
	anyPage := pages[0]

	file, err := scraper.GetFile(anyPage)

	if err != nil || file == nil {
		t.Failed()
	}
}