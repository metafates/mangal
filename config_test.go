package main

import "testing"

func TestParseConfig(t *testing.T) {
	configString := `
use = ['manganelo']
path = '.'
fullscreen = true
[sources]
    [sources.manganelo]
    base = 'https://ww5.manganelo.tv'
    search = 'https://ww5.manganelo.tv/search/%s'
    manga_anchor = '.search-story-item a.item-title'
    manga_title = '.search-story-item a.item-title'
    chapter_anchor = 'li.a-h a.chapter-name'
    chapter_title = 'li.a-h a.chapter-name'
    reader_page = '.container-chapter-reader img'
	random_delay_ms = 1337
`

	config, err := ParseConfig(configString)

	if err != nil {
		t.Fatal(err)
	}

	conditions := []bool{
		config.Fullscreen == true,
		config.Path == ".",
		config.Scrapers[0].Source.RandomDelayMs == 1337,
		config.Scrapers[0].Source.MangaTitle == ".search-story-item a.item-title",
	}

	for _, condition := range conditions {
		if !condition {
			t.Fatal()
		}
	}
}

func TestGetConfig(t *testing.T) {
	config := GetConfig("")

	conditions := []bool{
		len(config.Scrapers) > 0,
		config.Scrapers[0].Source != nil,
		config.Scrapers[0].Source.Base != "",
		config.Scrapers[0].Source.SearchTemplate != "",
		config.Scrapers[0].Source.MangaTitle != "",
		config.Scrapers[0].Source.MangaAnchor != "",
		config.Scrapers[0].Source.ChapterAnchor != "",
		config.Scrapers[0].Source.ChapterTitle != "",
		config.Scrapers[0].Source.ReaderPage != "",
	}

	for i, condition := range conditions {
		if !condition {
			t.Log(i)
			t.Fatal()
		}
	}
}
