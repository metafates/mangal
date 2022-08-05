package zip

import (
	"archive/zip"
	"fmt"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
	"time"
)

type ZIP struct{}

func New() *ZIP {
	return &ZIP{}
}

func (_ *ZIP) Save(chapter *source.Chapter) (string, error) {
	return save(chapter, false)
}

func (_ *ZIP) SaveTemp(chapter *source.Chapter) (string, error) {
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

	chapterZip := filepath.Join(mangaDir, util.SanitizeFilename(chapter.FormattedName())+".zip")
	zipFile, err := filesystem.Get().Create(chapterZip)
	if err != nil {
		return "", err
	}

	defer func(zipFile afero.File) {
		_ = zipFile.Close()
	}(zipFile)

	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		_ = zipWriter.Close()
	}()

	for _, page := range chapter.Pages {
		pageName := fmt.Sprintf("%d%s", page.Index, page.Extension)
		pageName = util.PadZero(pageName, 10)

		if err = addToZip(zipWriter, page.Contents, pageName); err != nil {
			return "", err
		}
	}

	return chapterZip, nil
}

// prepareMangaDir will create manga direcotry if it doesn't exist
func prepareMangaDir(manga *source.Manga) (mangaDir string, err error) {
	mangaDir = filepath.Join(
		viper.GetString(config.DownloaderPath),
		util.SanitizeFilename(manga.Name),
	)

	if err = filesystem.Get().MkdirAll(mangaDir, os.ModePerm); err != nil {
		return "", err
	}

	return mangaDir, nil
}

func addToZip(writer *zip.Writer, file io.Reader, name string) error {
	header := &zip.FileHeader{
		Name:     name,
		Method:   zip.Deflate,
		Modified: time.Now(),
	}

	headerWriter, err := writer.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(headerWriter, file)
	return err
}
