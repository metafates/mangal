package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/spf13/afero"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// RemoveIfExists removes file if it exists
func RemoveIfExists(path string) error {
	exists, err := Afero.Exists(path)

	if err != nil {
		return err
	}

	if exists {
		err = Afero.Remove(path)
		if err != nil {
			return err
		}
	}

	return nil
}

// SaveTemp saves file to OS temp dir and returns its path
// It's a caller responsibility to remove created file
func SaveTemp(buffer *bytes.Buffer) (string, error) {
	out, err := Afero.TempFile("", TempPrefix+"*")

	if err != nil {
		return "", err
	}

	defer func(out afero.File) {
		_ = out.Close()
	}(out)

	_, err = buffer.WriteTo(out)
	if err != nil {
		return "", err
	}

	return out.Name(), nil
}

type DownloaderStage int

const (
	Scraping DownloaderStage = iota
	Downloading
	Converting
	Cleanup
	Done
)

type ChaptersDownloadProgress struct {
	Current   *URL
	Done      bool
	Failed    []*URL
	Succeeded []string
	Total     int
	Proceeded int
}

type ChapterDownloadProgress struct {
	Stage   DownloaderStage
	Message string
}

// DownloadChapter downloads chapter from the given url and returns its path
func DownloadChapter(chapter *URL, progress chan ChapterDownloadProgress, temp bool) (string, error) {
	mangaTitle := SanitizeFilename(chapter.Relation.Info)
	var (
		mangaPath string
		err       error
	)

	// Get future path to manga
	if temp {
		mangaPath = os.TempDir()
	} else {
		absPath, err := filepath.Abs(UserConfig.Downloader.Path)

		if err != nil {
			return "", err
		}

		mangaPath = filepath.Join(absPath, mangaTitle)
	}

	showProgress := progress != nil

	if showProgress {
		progress <- ChapterDownloadProgress{
			Stage:   Scraping,
			Message: "Getting pages",
		}
	}

	// replace all placeholders with actual values
	var chapterName string

	// Why pad with 4 zeros? Because there are no manga with more than 9999 chapters
	// Actually, the longest manga has only 1960 chapters (Kochira Katsushika-ku Kameari Kōen-mae Hashutsujo)
	chapterName = strings.ReplaceAll(UserConfig.Downloader.ChapterNameTemplate, "%0d", PadZeros(chapter.Index, 4))
	chapterName = strings.ReplaceAll(chapterName, "%d", strconv.Itoa(chapter.Index))
	chapterName = strings.ReplaceAll(chapterName, "%s", chapter.Info)
	chapterName = SanitizeFilename(chapterName)

	// Get future path to chapter
	var chapterPath string
	if temp {
		chapterPath = filepath.Join(mangaPath, TempPrefix+" "+chapterName)
	} else {
		chapterPath = filepath.Join(mangaPath, chapterName)
	}

	// Get chapter contents
	pages, err := chapter.Scraper.GetPages(chapter)
	pagesCount := len(pages)

	if err != nil {
		return "", err
	}

	if showProgress {
		progress <- ChapterDownloadProgress{
			Stage:   Downloading,
			Message: fmt.Sprintf("Downloading %d pages", pagesCount),
		}
	}

	var (
		buffers          = make([]*bytes.Buffer, pagesCount)
		wg               sync.WaitGroup
		errorEncountered bool
	)

	wg.Add(pagesCount)

	// Download pages in parallel
	for _, page := range pages {
		go func(p *URL) {
			defer wg.Done()

			if errorEncountered {
				return
			}

			var data *bytes.Buffer

			data, err = chapter.Scraper.GetFile(p)

			if err != nil {
				// TODO: use channel
				errorEncountered = true
				return
			}

			buffers[p.Index] = data
		}(page)
	}

	wg.Wait()

	defer chapter.Scraper.ResetFiles()

	if errorEncountered {
		return "", err
	}

	if showProgress {
		progress <- ChapterDownloadProgress{
			Stage:   Converting,
			Message: fmt.Sprintf("Converting %d pages to %s", pagesCount, UserConfig.Formats.Default),
		}
	}

	if len(buffers) == 0 {
		return "", errors.New("pages was not downloaded")
	}

	// Convert pages to desired format
	chapterPath, err = Packers[UserConfig.Formats.Default](buffers, chapterPath, &PackerContext{
		Manga:   chapter.Relation,
		Chapter: chapter,
	})

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	if showProgress {
		progress <- ChapterDownloadProgress{
			Stage:   Cleanup,
			Message: "Removing temp files",
		}
	}

	if err != nil {
		return "", err
	}

	if showProgress {
		progress <- ChapterDownloadProgress{
			Stage:   Done,
			Message: fmt.Sprintf("Chapter %s downloaded", chapter.Info),
		}
	}

	return chapterPath, nil
}
