package main

import (
	"fmt"
	"log"
)

var (
	version    string
	build      string
	UserConfig *Config
)

func main() {
	//UserConfig = GetConfig()
	UserConfig = &DefaultConfig
	fmt.Println("Getting manga")
	manga, err := DefaultScraper.SearchManga("Berserk")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Got Manga, getting chapters")

	m := manga[0]
	chapters, err := DefaultScraper.GetChapters(m)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Got Chapters")

	for _, chapter := range chapters {
		path, err := DownloadChapter(chapter, nil)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(path)
	}

}
