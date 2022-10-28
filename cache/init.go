package cache

import (
	"encoding/json"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/util"
	"github.com/samber/mo"
	"io"
	"os"
	"time"
)

// init initializes the cache.
func (c *Cache[T]) init() error {
	if c.initialized {
		return nil
	}

	c.initialized = true
	log.Debugf("Initializing %s cacher", c.name)

	log.Debugf("Opening cache file at %s", c.path)
	file, err := filesystem.Api().OpenFile(c.path, os.O_RDONLY|os.O_CREATE, os.ModePerm)

	if err != nil {
		log.Warn(err)
		return err
	}

	defer util.Ignore(file.Close)

	contents, err := io.ReadAll(file)
	if err != nil {
		log.Warn(err)
		return err
	}

	if len(contents) == 0 {
		log.Debugf("%s cache file is empty, skipping unmarshal", c.name)
		if c.expireEvery.IsPresent() {
			c.data.Time = mo.Some(time.Now())
		}
		return nil
	}

	var unmarshalled internalData[T]
	err = json.Unmarshal(contents, &unmarshalled)
	if err != nil {
		log.Warn(err)
		return err
	}

	c.data = &unmarshalled

	if c.expireEvery.IsPresent() &&
		c.data.Time.IsPresent() &&
		time.Since(c.data.Time.MustGet()) >= c.expireEvery.MustGet() {
		log.Debugf("%s cache is expired, reseting cache", c.name)
		c.data.Time = mo.Some(time.Now())
		c.data.Internal = mo.None[T]()
		return filesystem.Api().WriteFile(c.path, []byte{}, os.ModePerm)
	}

	log.Debugf("%s cache file unmarshalled successfully", c.name)
	return nil
}
