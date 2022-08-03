package plain

import (
	"fmt"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/constants"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type Plain struct{}

func New() *Plain {
	return &Plain{}
}

func (p *Plain) Save(chapter *source.Chapter) (string, error) {
	return p.save(chapter, false)
}

func (p *Plain) SaveTemp(chapter *source.Chapter) (string, error) {
	return p.save(chapter, true)
}

func (_ *Plain) save(chapter *source.Chapter, temp bool) (string, error) {
	var (
		chapterDir string
		err        error
	)

	if temp {
		chapterDir, err = filesystem.Get().TempDir("", constants.TempPrefix)
	} else {
		chapterDir, err = prepareChapterDir(chapter)
	}

	if err != nil {
		return "", err
	}

	wg := sync.WaitGroup{}
	wg.Add(len(chapter.Pages))
	for _, page := range chapter.Pages {
		func(page *source.Page) {
			defer wg.Done()

			if err != nil {
				return
			}

			err = savePage(page, chapterDir)
		}(page)
	}

	wg.Wait()

	if err != nil {
		return "", err
	}

	abs, err := filepath.Abs(chapterDir)
	if err != nil {
		return chapterDir, nil
	}

	return abs, nil
}

// prepareMangaDir will create manga direcotry if it doesn't exist
func prepareChapterDir(chapter *source.Chapter) (chapterDir string, err error) {
	chapterDir = filepath.Join(
		viper.GetString(config.DownloaderPath),
		util.SanitizeFilename(chapter.Manga.Name),
		util.SanitizeFilename(chapter.FormattedName()),
	)

	if err = filesystem.Get().MkdirAll(chapterDir, os.ModePerm); err != nil {
		return "", err
	}

	return chapterDir, nil
}

func savePage(page *source.Page, to string) error {
	pageName := fmt.Sprintf("%d%s", page.Index, page.Extension)
	pageName = util.PadZero(pageName, 10)

	file, err := filesystem.Get().Create(filepath.Join(to, pageName))
	if err != nil {
		return err
	}

	_, err = io.Copy(file, page.Contents)
	if err != nil {
		return err
	}

	_ = file.Close()
	_ = page.Close()

	return nil
}
