package anilist

import (
	"github.com/metafates/gache"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/where"
	"github.com/samber/mo"
	"path/filepath"
	"time"
)

type cacheData[K comparable, T any] struct {
	Mangas map[K]T `json:"mangas"`
}

type cacher[K comparable, T any] struct {
	internal   *gache.Cache[*cacheData[K, T]]
	keyWrapper func(K) K
}

func (c *cacher[K, T]) Get(key K) mo.Option[T] {
	data, expired, err := c.internal.Get()
	if err != nil || expired || data == nil {
		return mo.None[T]()
	}

	mangas, ok := data.Mangas[c.keyWrapper(key)]
	if ok {
		return mo.Some(mangas)
	}

	return mo.None[T]()
}

func (c *cacher[K, T]) Set(key K, t T) error {
	data, expired, err := c.internal.Get()
	if err != nil {
		return err
	}

	if !expired && data != nil {
		data.Mangas[c.keyWrapper(key)] = t
		return c.internal.Set(data)
	} else {
		internal := &cacheData[K, T]{Mangas: make(map[K]T)}
		internal.Mangas[c.keyWrapper(key)] = t
		return c.internal.Set(internal)
	}
}

func (c *cacher[K, T]) Delete(key K) error {
	data, expired, err := c.internal.Get()
	if err != nil {
		return err
	}

	if !expired {
		delete(data.Mangas, c.keyWrapper(key))
		return c.internal.Set(data)
	}

	return nil
}

var relationCacher = &cacher[string, int]{
	internal: gache.New[*cacheData[string, int]](
		&gache.Options{
			Path:       where.AnilistBinds(),
			FileSystem: &filesystem.GacheFs{},
		},
	),
	keyWrapper: normalizedName,
}

var searchCacher = &cacher[string, []int]{
	internal: gache.New[*cacheData[string, []int]](
		&gache.Options{
			Path:       filepath.Join(where.Cache(), "anilist_search_cache.json"),
			Lifetime:   time.Hour * 24 * 10,
			FileSystem: &filesystem.GacheFs{},
		},
	),
	keyWrapper: normalizedName,
}

var idCacher = &cacher[int, *Manga]{
	internal: gache.New[*cacheData[int, *Manga]](
		&gache.Options{
			Path:       filepath.Join(where.Cache(), "anilist_id_cache.json"),
			Lifetime:   time.Hour * 24 * 2,
			FileSystem: &filesystem.GacheFs{},
		},
	),
	keyWrapper: func(id int) int { return id },
}

var failCacher = &cacher[string, bool]{
	internal: gache.New[*cacheData[string, bool]](
		&gache.Options{
			Path:       filepath.Join(where.Cache(), "anilist_fail_cache.json"),
			Lifetime:   time.Minute,
			FileSystem: &filesystem.GacheFs{},
		},
	),
	keyWrapper: normalizedName,
}
