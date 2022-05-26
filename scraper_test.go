package main

import (
	"testing"
)

func TestMakeSourceScraper(t *testing.T) {
	scraper := MakeSourceScraper(DefaultSource)

	if scraper == nil {
		t.Failed()
	}

	if scraper.Source == nil || scraper.Source != &DefaultSource {
		t.Failed()
	}
}

func TestScraper_SearchManga(t *testing.T) {
	manga, err := DefaultScraper.SearchManga(TestQuery)

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
		if m.Scraper != DefaultScraper {
			t.Fail()
		}
	}
}

func TestScraper_GetChapters(t *testing.T) {
	manga, _ := DefaultScraper.SearchManga(TestQuery)
	anyManga := manga[0]

	chapters, err := DefaultScraper.GetChapters(anyManga)

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
	manga, _ := DefaultScraper.SearchManga(TestQuery)
	anyManga := manga[0]

	chapters, _ := DefaultScraper.GetChapters(anyManga)
	anyChapter := chapters[0]

	pages, err := DefaultScraper.GetPages(anyChapter)

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
	manga, _ := DefaultScraper.SearchManga(TestQuery)
	anyManga := manga[0]

	chapters, _ := DefaultScraper.GetChapters(anyManga)
	anyChapter := chapters[0]

	pages, _ := DefaultScraper.GetPages(anyChapter)
	anyPage := pages[0]

	file, err := DefaultScraper.GetFile(anyPage)

	if err != nil || file == nil {
		t.Failed()
	}
}
