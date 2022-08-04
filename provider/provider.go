package provider

import (
	"github.com/metafates/mangal/icon"
	"github.com/metafates/mangal/provider/manganelo"
	"github.com/metafates/mangal/source"
)

type Provider struct {
	Name         string
	CreateSource func() source.Source
}

var defaultProviders = []*Provider{
	{
		Name:         "Manganelo",
		CreateSource: manganelo.New,
	},
}

func DefaultProviders() map[string]*Provider {
	providers := make(map[string]*Provider)

	for _, provider := range defaultProviders {
		provider.Name += " " + icon.Go()
		providers[provider.Name] = provider
	}

	return providers
}
