package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

var defaultConfigString = `
    use = ['manganelo']
    path = '.'
    fullscreen = true
	async_download = true
    [sources]
        [sources.manganelo]
        base = 'https://ww5.manganelo.tv'
        search = 'https://ww5.manganelo.tv/search/%s'
        manga_anchor = '.search-story-item a.item-title' manga_title = '.search-story-item a.item-title'
        chapter_anchor = 'li.a-h a.chapter-name'
        chapter_title = 'li.a-h a.chapter-name'
        reader_pages = '.container-chapter-reader img'
    `

func createDefault() Config {
	var conf Config

	if _, err := toml.Decode(defaultConfigString, &conf); err != nil {
		log.Fatal("Unexpected error while loading default config")
	}

	conf.setUsedSources()
	return conf
}
