package downloader

import (
	"github.com/metafates/mangai/api/scraper"
	"github.com/metafates/mangai/config"
	"github.com/metafates/mangai/shared"
	"github.com/spf13/afero"
	"io"
	"log"
	"path/filepath"

	pdf "github.com/pdfcpu/pdfcpu/pkg/api"
)

// DownloadTemp downloads file to the OS default temp directory.
// Returns temp file path and error
func DownloadTemp(url scraper.URL) (string, error) {
	out, err := shared.AferoFS.TempFile("", "mangai-temp-panel-*")

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
	Panel      string
	Percent    float64
	Converting bool
}

// DownloadChapter Downloads Chapters and returns its path
func DownloadChapter(mangaName string, chapter scraper.URL, infoChan chan ChapterDownloadInfo) (string, error) {
	mangaPath := pathOf(mangaName)

	err := shared.AferoFS.MkdirAll(mangaPath, 0700)
	if err != nil {
		return "", err
	}

	path := filepath.Join(mangaPath, chapter.Info+".pdf")

	var temps []string

	panels, err := chapter.Source.Panels(chapter)

	if err != nil {
		return "", err
	}

	var (
		percent   float64
		prevPanel string
	)
	// Download chapter temp images
	for i, panel := range panels {
		infoChan <- ChapterDownloadInfo{
			Panel:      panel.Info,
			Percent:    percent,
			Converting: false,
		}
		temp, err := DownloadTemp(*panel)

		if err != nil {
			return "", err
		}

		temps = append(temps, temp)
		percent = float64(i) / float64(len(panels)-1)
		prevPanel = panel.Info
	}

	infoChan <- ChapterDownloadInfo{
		Panel:      prevPanel,
		Percent:    percent,
		Converting: true,
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
