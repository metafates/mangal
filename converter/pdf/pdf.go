package pdf

import (
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
)

type PDF struct{}

func New() *PDF {
	return &PDF{}
}

func (_ *PDF) Save(chapter *source.Chapter) (string, error) {
	return save(chapter, false)
}

func (_ *PDF) SaveTemp(chapter *source.Chapter) (string, error) {
	return save(chapter, true)
}

func save(chapter *source.Chapter, temp bool) (string, error) {
	var (
		mangaDir string
		err      error
	)

	if temp {
		mangaDir, err = filesystem.Get().TempDir("", constant.TempPrefix)
	} else {
		mangaDir, err = prepareMangaDir(chapter.Manga)
	}

	if err != nil {
		return "", err
	}

	chapterPdf := filepath.Join(mangaDir, util.SanitizeFilename(chapter.FormattedName())+".pdf")
	pdfFile, err := filesystem.Get().Create(chapterPdf)
	if err != nil {
		return "", err
	}

	defer func(pdfFile afero.File) {
		_ = pdfFile.Close()
	}(pdfFile)

	var readers = make([]io.Reader, len(chapter.Pages))
	for i, page := range chapter.Pages {
		readers[i] = page
	}

	err = api.ImportImages(nil, pdfFile, readers, nil, nil)
	if err != nil {
		return "", err
	}

	return chapterPdf, nil
}

// prepareMangaDir will create manga direcotry if it doesn't exist
func prepareMangaDir(manga *source.Manga) (mangaDir string, err error) {
	absDownloaderPath, err := filepath.Abs(viper.GetString(config.DownloaderPath))
	if err != nil {
		return "", err
	}

	mangaDir = filepath.Join(
		absDownloaderPath,
		util.SanitizeFilename(manga.Name),
	)

	if err = filesystem.Get().MkdirAll(mangaDir, os.ModePerm); err != nil {
		return "", err
	}

	return mangaDir, nil
}
