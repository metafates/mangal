package cache

import (
	"github.com/metafates/mangal/util"
	"github.com/samber/mo"
	"path/filepath"
	"time"
)

// internalData is a struct that contains the data that is stored in the cache file with time of last update.
// Used to expire the cache.
type internalData[T any] struct {
	Internal mo.Option[T]         `json:"internal"`
	Time     mo.Option[time.Time] `json:"time"`
}

// Cache is a generic cache that can be used to cache any type of data.
// It is used to cache data that is expensive to fetch, such as API responses.
// Cached data is stored in a file, and is automatically expired after a certain amount of time
// (if expiration time is spcified)
type Cache[T any] struct {
	data        *internalData[T]
	name        string
	path        string
	expireEvery mo.Option[time.Duration]
	initialized bool
}

// New creates a new cache with the specified name and path.
// Name will be automatically converted to a valid file name.
func New[T any](path string, options *Options) *Cache[T] {
	name := util.FileStem(filepath.Base(path))
	return &Cache[T]{
		data: &internalData[T]{
			Internal: mo.None[T](),
		},
		expireEvery: options.ExpireEvery,
		name:        name,
		path:        path,
		initialized: false,
	}
}
