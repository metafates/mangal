package main

import (
	"fmt"
	"log"
)

// Set on compile time
var (
	version    string
	build      string
	UserConfig *Config
)

func main() {
	//UserConfig = GetConfig()
	UserConfig = &DefaultConfig

	manga, err := DefaultScraper.SearchManga("Attack on titan")

	if err != nil {
		log.Fatal(err)
	}

	m := manga[0]
	chapters, err := DefaultScraper.GetChapters(m)
	if err != nil {
		log.Fatal(err)
	}

	path, err := DownloadChapter(chapters[0])

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(path)
}
