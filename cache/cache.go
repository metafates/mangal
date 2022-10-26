package cache

import (
	"encoding/json"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/util"
	"github.com/metafates/mangal/where"
	"github.com/samber/mo"
	"io"
	"os"
	"path/filepath"
	"time"
)

type internalData[T any] struct {
	Internal mo.Option[T]         `json:"internal"`
	Time     mo.Option[time.Time] `json:"time"`
}

type Cache[T any] struct {
	data        *internalData[T]
	name        string
	path        string
	expireEvery time.Duration
	initialized bool
}

type Options[T any] struct {
	Initial     T
	ExpireEvery time.Duration
}

func New[T any](name string, options *Options[T]) *Cache[T] {
	return &Cache[T]{
		name: name,
		data: &internalData[T]{
			Internal: mo.Some(options.Initial),
		},
		expireEvery: options.ExpireEvery,
		path:        filepath.Join(where.Cache(), util.SanitizeFilename(name+".json")),
	}
}

func (c *Cache[T]) Init() error {
	if c.initialized {
		return nil
	}

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
		return nil
	}

	var unmarshalled internalData[T]
	err = json.Unmarshal(contents, &unmarshalled)
	if err != nil {
		log.Warn(err)
		return err
	}

	if unmarshalled.Time.IsPresent() {
		// check if timeout
		if time.Since(unmarshalled.Time.MustGet()) > c.expireEvery {
			log.Debugf("%s cache is expired, reseting cache", c.name)
			_ = filesystem.Api().WriteFile(c.path, []byte{}, os.ModePerm)
			return nil
		}
	} else {
		c.data.Time = mo.Some[time.Time](time.Now())
	}

	if unmarshalled.Internal.IsPresent() {
		c.data.Internal = unmarshalled.Internal
	}

	log.Debugf("%s cache file unmarshalled successfully", c.name)
	c.initialized = true
	return nil
}

func (c *Cache[T]) Get() T {
	_ = c.Init()

	return c.data.Internal.MustGet()
}

func (c *Cache[T]) Set(data T) error {
	_ = c.Init()

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
