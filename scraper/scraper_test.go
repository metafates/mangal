package scraper

import (
	"github.com/metafates/mangal/common"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/util"
	"testing"
)

var testDefaultSource = config.GetConfig("").Scrapers[0].Source

func TestMakeSourceScraper(t *testing.T) {
	scraper := MakeSourceScraper(testDefaultSource)

	if scraper == nil {
		t.Failed()
	}

	if scraper.Source == nil || scraper.Source != testDefaultSource {
		t.Failed()
	}
}

func TestScraper_SearchManga(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping this scraper.SearchManga is too expensive")
	}

	defaultScraper := MakeSourceScraper(testDefaultSource)
	manga, err := defaultScraper.SearchManga(common.TestQuery)

	if err != nil {
		t.Failed()
	}

	if len(manga) == 0 {
		t.Fail()
	}

	if !util.IsUnique(manga) {
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

	manga, _ := defaultScraper.SearchManga(common.TestQuery)
	anyManga := manga[0]

	chapters, err := defaultScraper.GetChapters(anyManga)

	if err != nil {
		t.Failed()
	}

	if len(chapters) == 0 {
		t.Failed()
	}

	if !util.IsUnique(chapters) {
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

	manga, _ := defaultScraper.SearchManga(common.TestQuery)
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

	if !util.IsUnique(pages) {
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

	manga, _ := defaultScraper.SearchManga(common.TestQuery)
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

func TestScraper_ResetFiles(t *testing.T) {
	defaultScraper := MakeSourceScraper(testDefaultSource)

	manga, _ := defaultScraper.SearchManga(common.TestQuery)
	anyManga := manga[0]

	chapters, _ := defaultScraper.GetChapters(anyManga)
	anyChapter := chapters[0]

	pages, _ := defaultScraper.GetPages(anyChapter)
	anyPage := pages[0]

	file, _ := defaultScraper.GetFile(anyPage)

	// check scraper has file
	f, ok := defaultScraper.Files.Get(anyPage.Address)
	if !ok {
		t.Failed()
	}

	// check file is the same as the one we got from scraper
	if f != file {
		t.Failed()
	}

	defaultScraper.ResetFiles()

	// check scraper has no file
	_, ok = defaultScraper.Files.Get(anyPage.Address)
	if ok {
		t.Failed()
	}
}
