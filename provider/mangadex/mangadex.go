package mangadex

import (
	"github.com/darylhjd/mangodex"
	"github.com/metafates/mangal/source"
)

const (
	Name = "Mangadex"
	ID   = Name + " built-in"
)

type Mangadex struct {
	client *mangodex.DexClient
	cache  struct {
		mangas   *cacher[[]*source.Manga]
		chapters *cacher[[]*source.Chapter]
	}
}

func (*Mangadex) Name() string {
	return Name
}

func (*Mangadex) ID() string {
	return ID
}

func New() *Mangadex {
	dex := &Mangadex{
		client: mangodex.NewDexClient(),
	}

	dex.cache.mangas = newCacher[[]*source.Manga](ID + "_mangas")
	dex.cache.chapters = newCacher[[]*source.Chapter](ID + "_chapters")

	return dex
}
