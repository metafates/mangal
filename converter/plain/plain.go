package plain

import (
	"fmt"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
)

type Plain struct{}

func New() *Plain {
	return &Plain{}
}

func (_ *Plain) Save(chapter *source.Chapter) (string, error) {
	mangaDir, err := prepareMangaDir(chapter.Manga)
	if err != nil {
		return "", err
	}

	chapterDir := filepath.Join(mangaDir, util.SanitizeFilename(chapter.Name))
	exists, err := filesystem.Get().Exists(chapterDir)
	if err != nil {
		return "", err
	}

	if !exists {
		err = filesystem.Get().Mkdir(chapterDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	for _, page := range chapter.Pages {
		pageName := fmt.Sprintf("%d%s", page.Index, page.Extension)
		pageName = util.PadZero(pageName, 10)

		file, err := filesystem.Get().Create(filepath.Join(chapterDir, pageName))
		if err != nil {
			return "", err
		}

		_, err = io.Copy(file, page.Contents)
		if err != nil {
			return "", err
		}

		_ = file.Close()
		_ = page.Close()
	}

	absMangaDir, err := filepath.Abs(mangaDir)
	if err != nil {
		absMangaDir = mangaDir
	}

	return absMangaDir, nil
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
