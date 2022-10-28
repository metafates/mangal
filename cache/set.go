package cache

import (
	"encoding/json"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/log"
	"github.com/samber/mo"
	"os"
)

// Set sets the cache data.
// May return error if writing to file failed
func (c *Cache[T]) Set(data T) error {
	_ = c.init()

	c.data.Internal = mo.Some(data)
	marshalled, err := json.Marshal(c.data)
	if err != nil {
		log.Warn(err)
		return err
	}

	log.Debugf("Writing %s cache file to %s", c.name, c.path)
	err = filesystem.Api().WriteFile(c.path, marshalled, os.ModePerm)
	if err != nil {
		log.Warn(err)
	}

	return err
}
