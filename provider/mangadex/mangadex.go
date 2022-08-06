package mangadex

import "github.com/darylhjd/mangodex"

const (
	Name = "Mangadex"
	ID   = Name + " built-in"
)

type Mangadex struct {
	client *mangodex.DexClient
}

func (m *Mangadex) Name() string {
	return Name
}

func (m *Mangadex) ID() string {
	return ID
}

func New() *Mangadex {
	return &Mangadex{
		client: mangodex.NewDexClient(),
	}
}
