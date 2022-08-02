package provider

import (
	"github.com/metafates/mangal/provider/manganelo"
	"github.com/metafates/mangal/source"
)

type Provider struct {
	Name   string
	Create func() source.Source
}

var defaultProviders = []*Provider{
	{
		Name:   "manganelo",
		Create: manganelo.New,
	},
}

func DefaultProviders() map[string]*Provider {
	providers := make(map[string]*Provider)

	for _, provider := range defaultProviders {
		providers[provider.Name] = provider
	}

	return providers
}
