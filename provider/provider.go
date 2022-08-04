package provider

import (
	"github.com/metafates/mangal/provider/manganelo"
	"github.com/metafates/mangal/source"
)

type Provider struct {
	ID           string
	Name         string
	CreateSource func() source.Source
}

var defaultProviders = []*Provider{
	{
		ID:           manganelo.ID,
		Name:         manganelo.Name,
		CreateSource: manganelo.New,
	},
}

func DefaultProviders() map[string]*Provider {
	providers := make(map[string]*Provider)

	for _, provider := range defaultProviders {
		providers[provider.Name] = provider
	}

	return providers
}
