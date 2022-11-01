package gache

import (
	"os"
	"path/filepath"
)

func (g *Cache[T]) load() error {
	if g.options.Path == "" {
		return nil
	}

	err := g.options.FileSystem.MkdirAll(filepath.Dir(g.options.Path), 0777)
	if err != nil {
		return err
	}

	file, err := g.options.FileSystem.OpenFile(g.options.Path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	err = g.options.Decoder.Decode(file, &g.data)
	if err != nil {
		// if the file is malformed, reset it
		err = g.save()
		if err != nil {
			return err
		}
	}

	// check if the cache has expired
	return g.tryExpire()
}
