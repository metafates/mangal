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

func Default() []*Provider {
	return defaultProviders
}
