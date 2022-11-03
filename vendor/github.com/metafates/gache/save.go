package gache

import (
	"os"
	"path/filepath"
	"time"
)

func (g *Cache[T]) save() error {
	// do nothing if we use in-memory caching
	if g.options.Path == "" {
		return nil
	}

	err := g.options.FileSystem.MkdirAll(filepath.Dir(g.options.Path), 0777)
	if err != nil {
		return err
	}

	file, err := g.options.FileSystem.OpenFile(g.options.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	// save the time of the last update
	now := time.Now()
	g.data.Time = &now

	return g.options.Encoder.Encode(file, g.data)
}
