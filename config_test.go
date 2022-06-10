package main

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func TestParseConfig(t *testing.T) {
	var (
		sourceName    = "manganelo"
		randomDelayMs = 1337
		path          = "."
		fullscreen    = true
		mangaTitle    = ".search-story-item a.item-title"
	)
	configString := fmt.Sprintf(`
use = ['%s']
path = "%s"
fullscreen = %t
[sources]
    [sources.%s]
    base = 'https://ww5.manganelo.tv'
    search = 'https://ww5.manganelo.tv/search/%s'
    manga_anchor = '.search-story-item a.item-title'
    manga_title = '%s'
    chapter_anchor = 'li.a-h a.chapter-name'
    chapter_title = 'li.a-h a.chapter-name'
    reader_page = '.container-chapter-reader img'
	random_delay_ms = %d
`, sourceName, path, fullscreen, sourceName, "", mangaTitle, randomDelayMs)

	config, err := ParseConfig(configString)

	if err != nil {
		t.Fatal(err)
	}

	conditions := []bool{
		config.Fullscreen == fullscreen,
		config.Path == path,
		config.Scrapers[0].Source.RandomDelayMs == randomDelayMs,
		config.Scrapers[0].Source.MangaTitle == mangaTitle,
	}

	for _, condition := range conditions {
		if !condition {
			t.Error()
		}
	}
}

func testConfig(t *testing.T, config *Config) {
	t.Helper()

	var isValidURI = func(uri string) bool {
		_, err := url.ParseRequestURI(strings.Replace(uri, "%s", "", 1))
		return err == nil
	}

	conditions := []bool{
		len(config.Scrapers) > 0,
		config.Scrapers[0].Source != nil,
		config.Scrapers[0].Source.Base != "",
		isValidURI(config.Scrapers[0].Source.Base),
		config.Scrapers[0].Source.SearchTemplate != "",
		strings.Contains(config.Scrapers[0].Source.SearchTemplate, "%s"),
		isValidURI(config.Scrapers[0].Source.SearchTemplate),
		config.Scrapers[0].Source.MangaTitle != "",
		config.Scrapers[0].Source.MangaAnchor != "",
		config.Scrapers[0].Source.ChapterAnchor != "",
		config.Scrapers[0].Source.ChapterTitle != "",
		config.Scrapers[0].Source.ReaderPage != "",
		config.Scrapers[0].Source.RandomDelayMs >= 0,
	}

	for i, condition := range conditions {
		if !condition {
			t.Error(i)
		}
	}
}

func TestGetConfig(t *testing.T) {
	config := GetConfig("")
	testConfig(t, config)

}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	if config == nil {
		t.Fatal("Error while parsing default config file")
	}

	testConfig(t, config)
}
