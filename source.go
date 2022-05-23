package main

type Source struct {
	Base           string
	SearchTemplate string `toml:"search"`

	MangaAnchor string `toml:"manga_anchor"`
	MangaTitle  string `toml:"manga_title"`

	ChapterAnchor string `toml:"chapter_anchor"`
	ChapterTitle  string `toml:"chapter_title"`

	ReaderPage string `toml:"reader_page"`

	RandomDelayMs int `toml:"random_delay_ms"`
}

var DefaultSource = Source{
	Base:           "https://ww5.manganelo.tv",
	SearchTemplate: "https://ww5.manganelo.tv/search/%s",
	MangaAnchor:    ".search-story-item a.item-title",
	MangaTitle:     ".search-story-item a.item-title",
	ChapterAnchor:  "li.a-h a.chapter-name",
	ChapterTitle:   "li.a-h a.chapter-name",
	ReaderPage:     ".container-chapter-reader img",
	RandomDelayMs:  200,
}
