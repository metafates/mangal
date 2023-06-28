package path

import (
	"log"
	"path/filepath"

	"github.com/mangalorg/mangal/fs"
)

const (
    permDir = 0755
    permFile = 0655
)

func createIfAbsent(path string) {
    exists, err := fs.FS.Exists(path)
    if err != nil {
        log.Fatal(err)
        return
    }

    if exists {
        return
    }

    isDir, err := fs.FS.IsDir(path)    
    if err != nil {
        log.Fatal(err)
        return
    }

    if isDir {
        if err := fs.FS.MkdirAll(path, permDir); err != nil {
            log.Fatal(err)
        }

        return
    }

    if err = fs.FS.MkdirAll(filepath.Base(path), permDir); err != nil {
        log.Fatal(err)
        return
    }

    file, err := fs.FS.Create(path)
    if err != nil {
        log.Fatal(err)
        return
    }

    file.Close()
}
