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
	expireEvery mo.Option[time.Duration]
	initialized bool
}

type Options[T any] struct {
	ExpireEvery mo.Option[time.Duration]
}

func New[T any](name string, options *Options[T]) *Cache[T] {
	return &Cache[T]{
		name: name,
		data: &internalData[T]{
			Internal: mo.None[T](),
		},
		expireEvery: options.ExpireEvery,
		path:        filepath.Join(where.Cache(), util.SanitizeFilename(name+".json")),
		initialized: false,
	}
}

func (c *Cache[T]) Init() error {
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

	if unmarshalled.Time.IsPresent() {
		c.data.Time = unmarshalled.Time
	} else if c.expireEvery.IsPresent() {
		c.data.Time = mo.Some(time.Now())
	}

	if c.expireEvery.IsPresent() &&
		time.Since(unmarshalled.Time.MustGet()) > c.expireEvery.MustGet() {
		log.Debugf("%s cache is expired, reseting cache", c.name)
		_ = filesystem.Api().WriteFile(c.path, []byte{}, os.ModePerm)
		return nil
	}

	if unmarshalled.Internal.IsPresent() {
		c.data.Internal = unmarshalled.Internal
	}

	log.Debugf("%s cache file unmarshalled successfully", c.name)
	return nil
}

func (c *Cache[T]) Get() mo.Option[T] {
	_ = c.Init()

	return c.data.Internal
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
