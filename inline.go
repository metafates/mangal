package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"sync"
)

// InlineOptions provides all options for inline mode
type InlineOptions struct {
	config     string
	mangaIdx   int
	chapterIdx int
	asJson     bool
	format     FormatType
	showUrls   bool
	asTemp     bool
	doRead     bool
	doOpen     bool
}

// InlineMode provides all functionality of TUI but in inline mode
// TODO: split into subfunctions
func InlineMode(query string, options InlineOptions) (string, error) {
	initConfig(options.config, false)

	if !options.asTemp {
		defer RemoveTemp()
	}

	if options.format != "" {
		UserConfig.Format = options.format
	}

	// Check if config is valid
	if err := ValidateConfig(UserConfig); err != nil {
		return "", err
	}

	var (
		manga []*URL
		wg    sync.WaitGroup
	)

	// Check if query is valid
	if query == "" {
		return "", errors.New("query expected")
	}

	wg.Add(len(UserConfig.Scrapers))

	// Search for manga in all scrapers
	for _, scraper := range UserConfig.Scrapers {
		go func(s *Scraper) {
			defer wg.Done()

			m, err := s.SearchManga(query)

			if err == nil {
				manga = append(manga, m...)
			}
		}(scraper)
	}

	wg.Wait()

	// Check if manga was selected
	if options.mangaIdx >= 0 {
		if options.mangaIdx > len(manga) || options.mangaIdx <= 0 {
			return "", errors.New("index out of range")
		}

		selectedManga := manga[options.mangaIdx-1]

		// Get chapters of selected manga
		chapters, err := selectedManga.Scraper.GetChapters(selectedManga)
		if err != nil {
			return "", errors.New("error while getting chapters")
		}

		// Check if chapter was selected
		if options.chapterIdx >= 0 {

			// Get selected chapter
			selectedChapter, ok := Find(chapters, func(c *URL) bool {
				return c.Index == options.chapterIdx
			})

			if !ok {
				return "", errors.New("index out of range")
			}

			// if option to open chapter is set, open it
			if options.doOpen {
				if err = open.Start(selectedChapter.Address); err != nil {
					return "", errors.New("unexpected error while trying to open url")
				}

				return "", nil
			}

			// if option to read chapter is set download chapter as temp file
			if options.doRead {
				options.asTemp = true
			}

			// Download chapter
			chapterPath, err := DownloadChapter(selectedChapter, nil, options.asTemp)
			if err != nil {
				return "", errors.New("error while downloading chapter")
			}

			// if epub file was used create it
			if EpubFile != nil {
				EpubFile.SetAuthor(selectedManga.Scraper.Source.Base)
				if err := EpubFile.Write(chapterPath); err != nil {
					return "", errors.New("error while making epub file")
				}

				// reset epub file
				EpubFile = nil
			}

			// if options to read chapter is set, read it
			if options.doRead {
				// check if custom reader is set
				if UserConfig.UseCustomReader {
					err = open.StartWith(chapterPath, UserConfig.CustomReader)
				} else {
					err = open.Start(chapterPath)
				}

				if err != nil {
					return "", err
				}

				return "", nil
			}

			return chapterPath, nil
		}

		// if option to print data as json is set, print it as json
		if options.asJson {
			data, err := json.Marshal(chapters)
			if err != nil {
				return "", errors.New("could not get data as json")
			}

			return string(data), nil
		} else if options.doOpen {
			if err = open.Start(selectedManga.Address); err != nil {
				return "", errors.New("unexpected error while trying to open url")
			}
			return "", nil
		} else {
			var chaptersString string

			// print chapters list
			for _, c := range chapters {
				if options.showUrls {
					chaptersString += fmt.Sprintf("[%d] %s %s\n", c.Index, c.Info, c.Address)
				} else {
					chaptersString += fmt.Sprintf("[%d] %s\n", c.Index, c.Info)
				}
			}

			return chaptersString, nil
		}

	} else {
		// if option to print data as json is set, print it as json
		if options.asJson {
			data, err := json.Marshal(manga)
			if err != nil {
				return "", errors.New("could not get data as json")
			}

			return string(data), nil
		} else {
			var mangaString string

			// print manga list
			for i, m := range manga {
				if options.showUrls {
					mangaString += fmt.Sprintf("[%d] %s %s\n", i+1, m.Info, m.Address)
				} else {
					mangaString += fmt.Sprintf("[%d] %s\n", i+1, m.Info)
				}
			}

			return mangaString, nil
		}
	}
}
