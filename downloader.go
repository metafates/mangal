package main

import (
	"fmt"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/spf13/afero"
	"log"
	"path/filepath"
	"sync"
)

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
func SaveTemp(contents *[]byte) (string, error) {
	out, err := Afero.TempFile("", AppName+"-*")

	if err != nil {
		return "", err
	}

	defer func(out afero.File) {
		err := out.Close()
		if err != nil {
			log.Fatal("Unexpected error while closing file")
		}
	}(out)

	_, err = out.Write(*contents)
	if err != nil {
		return "", err
	}

	return out.Name(), nil
}

// batchRemove removes all files from the path list
func batchRemove(paths []string) error {
	for _, path := range paths {
		err := FileSystem.Remove(path)
		if err != nil {
			return err
		}
	}
	return nil
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
	Failed    []string
	Succeeded []string
	Total     int
	Proceeded int
}

type ChapterDownloadProgress struct {
	Stage   DownloaderStage
	Message string
}

func DownloadChapter(chapter *URL, progress chan ChapterDownloadProgress) (string, error) {
	mangaTitle := chapter.Relation.Info
	mangaPath, err := filepath.Abs(filepath.Join(UserConfig.Path, mangaTitle))

	if err != nil {
		return "", nil
	}

	err = Afero.MkdirAll(mangaPath, 0700)

	if err != nil {
		return "", nil
	}

	showProgress := progress != nil

	if showProgress {
		progress <- ChapterDownloadProgress{
			Stage:   Scraping,
			Message: "Getting pages",
		}
	}

	chapterPath := filepath.Join(mangaPath, fmt.Sprintf("[%d] %s.pdf", chapter.Index, chapter.Info))
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
		tempPaths        = make([]string, pagesCount)
		wg               sync.WaitGroup
		errorEncountered bool
	)

	wg.Add(pagesCount)

	for _, page := range pages {
		go func(p *URL) {
			defer wg.Done()
			var (
				data     *[]byte
				tempPath string
			)

			data, err = chapter.Scraper.GetFile(p)

			if err != nil {
				// TODO: use channel
				errorEncountered = true
				//log.Fatalf("Error while downloading page %s", p.Address)
				return
			}

			tempPath, err = SaveTemp(data)
			tempPaths[p.Index] = tempPath
		}(page)
	}

	wg.Wait()

	defer chapter.Scraper.CleanupFiles()

	if errorEncountered {
		return "", err
	}

	if showProgress {
		progress <- ChapterDownloadProgress{
			Stage:   Converting,
			Message: fmt.Sprintf("Converting %d pages to pdf", pagesCount),
		}
	}

	err = RemoveIfExists(chapterPath)
	if err != nil {
		return "", err
	}

	err = pdfcpu.ImportImagesFile(tempPaths, chapterPath, nil, nil)

	if err != nil {
		return "", err
	}

	if showProgress {
		progress <- ChapterDownloadProgress{
			Stage:   Cleanup,
			Message: "Removing temp files",
		}
	}

	// Cleanup temp files
	err = batchRemove(tempPaths)

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
