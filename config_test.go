package main

import (
	"fmt"
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

func TestGetConfig(t *testing.T) {
	config := GetConfig("")
	if err := ValidateConfig(config); err != nil {
		t.Error(err)
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	if config == nil {
		t.Fatal("Error while parsing default config file")
	}

	if err := ValidateConfig(config); err != nil {
		t.Error(err)
	}
}
