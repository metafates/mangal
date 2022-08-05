package cbz

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
	"strings"
)

type CBZ struct{}

func New() *CBZ {
	return &CBZ{}
}

func (_ *CBZ) Save(chapter *source.Chapter) (string, error) {
	return save(chapter, false)
}

func (_ *CBZ) SaveTemp(chapter *source.Chapter) (string, error) {
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

	chapterCbz := filepath.Join(mangaDir, util.SanitizeFilename(chapter.FormattedName())+".cbz")
	cbzFile, err := filesystem.Get().Create(chapterCbz)
	if err != nil {
		return "", err
	}

	defer func(cbzFile afero.File) {
		_ = cbzFile.Close()
	}(cbzFile)

	zipWriter := zip.NewWriter(cbzFile)
	defer func(zipWriter *zip.Writer) {
		_ = zipWriter.Close()
	}(zipWriter)

	for _, page := range chapter.Pages {
		pageName := fmt.Sprintf("%d%s", page.Index, page.Extension)
		pageName = util.PadZero(pageName, 10)

		if err = addToZip(zipWriter, page.Contents, pageName); err != nil {
			return "", err
		}
	}

	err = addToZip(zipWriter, strings.NewReader(comicInfo(chapter)), "ComicInfo.xml")
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(chapterCbz)
	if err != nil {
		return chapterCbz, nil
	}

	return absPath, nil
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

func comicInfo(chapter *source.Chapter) string {
	return `
<ComicInfo xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <Title>` + chapter.Name + `</Title>
  <Series>` + chapter.Manga.Name + `</Series>
  <Genre>Web Comic</Genre>
  <Web>` + chapter.Manga.URL + `</Web>
  <Manga>Yes</Manga>
</ComicInfo>`
}

func addToZip(writer *zip.Writer, file io.Reader, name string) error {
	header := &zip.FileHeader{
		Name:   name,
		Method: zip.Store,
	}

	headerWriter, err := writer.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(headerWriter, file)
	return err
}
