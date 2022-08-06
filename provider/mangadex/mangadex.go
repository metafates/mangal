package mangadex

import (
	"github.com/darylhjd/mangodex"
	"github.com/metafates/mangal/source"
)

const (
	Name = "Mangadex (Unstable)"
	ID   = Name + " built-in"
)

type Mangadex struct {
	client         *mangodex.DexClient
	cachedMangas   map[string][]*source.Manga
	cachedChapters map[string][]*source.Chapter
}

func (m *Mangadex) Name() string {
	return Name
}

func (m *Mangadex) ID() string {
	return ID
}

func New() *Mangadex {
	return &Mangadex{
		client:         mangodex.NewDexClient(),
		cachedMangas:   make(map[string][]*source.Manga),
		cachedChapters: make(map[string][]*source.Chapter),
	}
}
