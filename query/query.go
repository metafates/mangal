package query

import (
	"github.com/metafates/mangal/cache"
	"github.com/metafates/mangal/where"
	"github.com/samber/mo"
	"time"
)

type queryRecord struct {
	Rank  int    `json:"rank"`
	Query string `json:"query"`
}

var cacher = cache.New[map[string]*queryRecord](
	where.Queries(),
	&cache.Options{
		ExpireEvery: mo.None[time.Duration](),
	},
)
