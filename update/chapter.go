package update

import (
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/log"
	"os"
	"path/filepath"
)

type downloadedChapter struct {
	path   string
	format string
}

func getChapters(manga string) ([]*downloadedChapter, error) {
	log.Infof("getting chapters for %s", manga)
	var chapters []*downloadedChapter

	err := filesystem.Api().Walk(manga, func(path string, info os.FileInfo, err error) error {
		// we will ignore plain chapter (aka folder ones) for the sake of simplicity
		if info.IsDir() {
			return nil
		}

		name := info.Name()
		switch filepath.Ext(name)[1:] {
		case constant.CBZ:
			chapters = append(chapters, &downloadedChapter{path: path, format: constant.CBZ})
		case constant.PDF:
			chapters = append(chapters, &downloadedChapter{path: path, format: constant.PDF})
		case constant.ZIP:
			chapters = append(chapters, &downloadedChapter{path: path, format: constant.ZIP})
		}

		return nil
	})

	return chapters, err
}
