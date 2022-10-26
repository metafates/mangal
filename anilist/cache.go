package anilist

import (
	"github.com/metafates/mangal/cache"
	"github.com/metafates/mangal/constant"
	"time"
)

type cacheData[K comparable, T any] struct {
	Mangas map[K]T `json:"mangas"`
}

type cacher[K comparable, T any] struct {
	internal   *cache.Cache[*cacheData[K, T]]
	keyWrapper func(K) K
}

func (c *cacher[K, T]) Get(key K) (T, bool) {
	mangas, ok := c.internal.Get().Mangas[c.keyWrapper(key)]
	return mangas, ok
}

func (c *cacher[K, T]) Set(key K, t T) error {
	data := c.internal.Get()
	data.Mangas[c.keyWrapper(key)] = t
	return c.internal.Set(data)
}

var relationCacher = &cacher[string, int]{
	internal: cache.New(
		"anilist_relation_cache",
		&cache.Options[*cacheData[string, int]]{
			Initial: &cacheData[string, int]{Mangas: make(map[string]int, 0)},
			// never expire
			ExpireEvery: constant.Forever,
		},
	),
	keyWrapper: normalizedName,
}

var searchCacher = &cacher[string, []int]{
	internal: cache.New(
		"anilist_search_cache",
		&cache.Options[*cacheData[string, []int]]{
			Initial: &cacheData[string, []int]{Mangas: make(map[string][]int, 0)},
			// update ids every 10 days, since new manga are not added that often
			ExpireEvery: time.Hour * 24 * 10,
		},
	),
	keyWrapper: normalizedName,
}

var idCacher = &cacher[int, *Manga]{
	internal: cache.New(
		"anilist_id_cache",
		&cache.Options[*cacheData[int, *Manga]]{
			Initial: &cacheData[int, *Manga]{Mangas: make(map[int]*Manga, 0)},
			// update manga data every day since it can change often
			ExpireEvery: time.Hour * 24,
		},
	),
	keyWrapper: func(id int) int { return id },
}
