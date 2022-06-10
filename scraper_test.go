package main

import (
	"testing"
)

var testDefaultSource = Source{
	Base:             "https://ww5.manganelo.tv",
	SearchTemplate:   "https://ww5.manganelo.tv/search/%s",
	MangaAnchor:      ".search-story-item a.item-title",
	MangaTitle:       ".search-story-item a.item-title",
	ChapterAnchor:    "li.a-h a.chapter-name",
	ChapterTitle:     "li.a-h a.chapter-name",
	ReaderPage:       ".container-chapter-reader img",
	RandomDelayMs:    700,
	ChaptersReversed: true,
}

func TestMakeSourceScraper(t *testing.T) {
	scraper := MakeSourceScraper(testDefaultSource)

	if scraper == nil {
		t.Failed()
	}

	if scraper.Source == nil || scraper.Source != &testDefaultSource {
		t.Failed()
	}
}

func TestScraper_SearchManga(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping this scraper.SearchManga is too expensive")
	}

	defaultScraper := MakeSourceScraper(testDefaultSource)
	manga, err := defaultScraper.SearchManga(TestQuery)

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
		if m.Scraper != defaultScraper {
			t.Fail()
		}
	}
}

func TestScraper_GetChapters(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping this scraper.GetChapters is too expensive")
	}

	defaultScraper := MakeSourceScraper(testDefaultSource)

	manga, _ := defaultScraper.SearchManga(TestQuery)
	anyManga := manga[0]

	chapters, err := defaultScraper.GetChapters(anyManga)

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
	if testing.Short() {
		t.Skip("skipping this scraper.GetPages is too expensive")
	}

	defaultScraper := MakeSourceScraper(testDefaultSource)

	manga, _ := defaultScraper.SearchManga(TestQuery)
	anyManga := manga[0]

	chapters, _ := defaultScraper.GetChapters(anyManga)
	anyChapter := chapters[0]

	pages, err := defaultScraper.GetPages(anyChapter)

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
	if testing.Short() {
		t.Skip("skipping this scaper.GetFile is too expensive")
	}

	defaultScraper := MakeSourceScraper(testDefaultSource)

	manga, _ := defaultScraper.SearchManga(TestQuery)
	anyManga := manga[0]

	chapters, _ := defaultScraper.GetChapters(anyManga)
	anyChapter := chapters[0]

	pages, _ := defaultScraper.GetPages(anyChapter)
	anyPage := pages[0]

	file, err := defaultScraper.GetFile(anyPage)

	if err != nil || file == nil {
		t.Failed()
	}
}
