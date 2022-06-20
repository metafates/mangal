package main

import (
	"testing"
)

func TestParseConfig(t *testing.T) {
	configString := `
use = ['manganelo']
format = "pdf"
use_custom_reader = false
custom_reader = "zathura"
download_path = '.'
cache_images = false
[ui]
fullscreen = true
prompt = "🔍"
placeholder = "What shall we look for?"
mark = "▼"
title = "Mangal"
[sources]
[sources.manganelo]
base = 'https://ww5.manganelo.tv'
search = 'https://ww5.manganelo.tv/search/%s'
manga_anchor = '.search-story-item a.item-title'
manga_title = '.search-story-item a.item-title'
chapter_anchor = 'li.a-h a.chapter-name'
chapter_title = 'li.a-h a.chapter-name'
reader_page = '.container-chapter-reader img'
random_delay_ms = 500 # ms
reversed_chapters_order = true
`

	config, err := ParseConfig(configString)

	if err != nil {
		t.Fatal(err)
	}

	conditions := []bool{
		config.UI.Fullscreen == true,
		config.Path == ".",
		config.Scrapers[0].Source.RandomDelayMs == 500,
		config.Scrapers[0].Source.MangaTitle == ".search-story-item a.item-title",
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
