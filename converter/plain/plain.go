package plain

import (
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/source"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type Plain struct{}

func New() *Plain {
	return &Plain{}
}

func (*Plain) Save(chapter *source.Chapter) (string, error) {
	return save(chapter, false)
}

func (*Plain) SaveTemp(chapter *source.Chapter) (string, error) {
	return save(chapter, true)
}

func save(chapter *source.Chapter, temp bool) (path string, err error) {
	path, err = chapter.Path(temp)
	if err != nil {
		return
	}

	err = filesystem.Api().Mkdir(path, os.ModePerm)
	if err != nil {
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(chapter.Pages))
	for _, page := range chapter.Pages {
		func(page *source.Page) {
			defer wg.Done()

			if err != nil {
				return
			}

			err = savePage(page, path)
		}(page)
	}

	wg.Wait()
	return
}

func savePage(page *source.Page, to string) error {
	file, err := filesystem.Api().Create(filepath.Join(to, page.Filename()))
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
