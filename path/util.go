package path

import (
	"log"

	"github.com/mangalorg/mangal/afs"
)

const (
	permDir  = 0755
	permFile = 0655
)

func createDirIfAbsent(path string) {
	exists, err := afs.Afero.Exists(path)
	if err != nil {
		log.Fatal(err)
		return
	}

	if exists {
		return
	}

	if err := afs.Afero.MkdirAll(path, permDir); err != nil {
		log.Fatal(err)
	}

	return
}
