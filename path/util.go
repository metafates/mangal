package path

import (
	"github.com/mangalorg/mangal/fs"
	"log"
)

const (
	permDir  = 0755
	permFile = 0655
)

func createDirIfAbsent(path string) {
	exists, err := fs.Afero.Exists(path)
	if err != nil {
		log.Fatal(err)
		return
	}

	if exists {
		return
	}

	if err := fs.Afero.MkdirAll(path, permDir); err != nil {
		log.Fatal(err)
	}

	return
}
