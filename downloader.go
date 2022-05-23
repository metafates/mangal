package main

import (
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/spf13/afero"
	"log"
	"path/filepath"
	"strconv"
	"sync"
)

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

type DownloadProgress struct {
	Stage   DownloaderStage
	Message string
}

func DownloadChapter(chapter *URL) (string, error) {
	mangaTitle := chapter.Relation.Info
	mangaPath := filepath.Join(UserConfig.Path, mangaTitle)

	err := Afero.MkdirAll(mangaPath, 0700)

	if err != nil {
		return "", nil
	}

	//progress <- DownloadProgress{
	//	Stage:   Scraping,
	//	Message: "Getting pages",
	//}

	chapterPath := filepath.Join(mangaPath, chapter.Info+".pdf")
	pages, err := chapter.Scraper.GetPages(chapter)
	pagesCount := len(pages)

	if err != nil {
		return "", err
	}

	//progress <- DownloadProgress{
	//	Stage:   Downloading,
	//	Message: fmt.Sprintf("Downloading %d pages", pagesCount),
	//}

	var (
		tempPaths = make([]string, pagesCount)
		wg        sync.WaitGroup
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
				log.Fatal("Error while downloading page")
				return
			}

			tempPath, err = SaveTemp(data)
			i, _ := strconv.Atoi(p.Info)
			tempPaths[i] = tempPath
		}(page)
	}

	wg.Wait()

	defer chapter.Scraper.CleanupFiles()

	//progress <- DownloadProgress{
	//	Stage:   Converting,
	//	Message: fmt.Sprintf("Converting %d pages to pdf", pagesCount),
	//}

	err = pdfcpu.ImportImagesFile(tempPaths, chapterPath, nil, nil)

	if err != nil {
		return "", err
	}

	//progress <- DownloadProgress{
	//	Stage:   Cleanup,
	//	Message: "Removing temp files",
	//}

	// Cleanup temp files
	err = batchRemove(tempPaths)

	if err != nil {
		return "", err
	}

	//progress <- DownloadProgress{
	//	Stage:   Done,
	//	Message: fmt.Sprintf("Chapter %s downloaded", chapter.Info),
	//}

	return chapterPath, nil
}
