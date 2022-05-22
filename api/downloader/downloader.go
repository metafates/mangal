package downloader

import (
	"github.com/metafates/mangai/api/scraper"
	"github.com/metafates/mangai/config"
	"github.com/metafates/mangai/shared"
	"github.com/spf13/afero"
	"io"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"sync"

	pdf "github.com/pdfcpu/pdfcpu/pkg/api"
)

// DownloadTemp downloads file to the OS default temp directory.
// Returns temp file path and error
func DownloadTemp(url scraper.URL, prefix string) (string, error) {
	out, err := shared.AferoFS.TempFile("", prefix+"-mangai-*")

	if err != nil {
		return "", err
	}

	defer func(out afero.File) {
		err := out.Close()
		if err != nil {
			log.Fatal("Unexpected error while closing file")
		}
	}(out)

	resp, err := url.Get()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal("Unexpected error while closing http connection")
		}
	}(resp.Body)

	if err != nil {
		return "", err
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return out.Name(), nil
}

func pathOf(mangaName string) string {
	return filepath.Join(config.Get().Path, mangaName)
}

type ChapterDownloadInfo struct {
	PanelsCount     int
	ConvertingToPdf bool
}

// DownloadChapter Downloads Chapters and returns its path
func DownloadChapter(mangaName string, chapter scraper.URL, infoChan chan ChapterDownloadInfo) (string, error) {
	mangaPath := pathOf(mangaName)

	err := shared.AferoFS.MkdirAll(mangaPath, 0700)
	if err != nil {
		return "", err
	}

	path := filepath.Join(mangaPath, chapter.Info+".pdf")

	panels, err := chapter.Source.Panels(chapter)

	if err != nil {
		return "", err
	}

	var (
		temps       []string
		wg          sync.WaitGroup
		panelsCount = len(panels)
	)

	infoChan <- ChapterDownloadInfo{
		PanelsCount:     panelsCount,
		ConvertingToPdf: false,
	}

	wg.Add(panelsCount)

	// Download chapter temp images
	for i, panel := range panels {
		if err != nil {
			return "", nil
		}

		go func(panel *scraper.URL, index int) {
			defer wg.Done()

			var temp string
			temp, err = DownloadTemp(*panel, strconv.Itoa(index))

			temps = append(temps, temp)
		}(panel, i)
	}

	wg.Wait()

	sort.Strings(temps)

	infoChan <- ChapterDownloadInfo{
		PanelsCount:     panelsCount,
		ConvertingToPdf: true,
	}
	// Convert images to pdf
	err = pdf.ImportImagesFile(temps, path, nil, nil)

	if err != nil {
		return "", err
	}

	// Remove temp files
	for _, temp := range temps {
		err := shared.AferoBackend.Remove(temp)
		if err != nil {
			return "", err
		}
	}

	return path, nil
}
