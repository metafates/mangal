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

	ChaptersReversed bool `toml:"reversed_chapters_order"`
}
