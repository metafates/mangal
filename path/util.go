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
	exists, err := fs.FS.Exists(path)
	if err != nil {
		log.Fatal(err)
		return
	}

	if exists {
		return
	}

	if err := fs.FS.MkdirAll(path, permDir); err != nil {
		log.Fatal(err)
	}

	return
}
