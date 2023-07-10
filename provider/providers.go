package provider

import (
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/provider/lua"
)

func InstalledProviders() ([]libmangal.ProviderLoader, error) {
	var loaders []libmangal.ProviderLoader

	for _, getLoaders := range []func() ([]libmangal.ProviderLoader, error){
		lua.InstalledProviders,
	} {
		l, err := getLoaders()
		if err != nil {
			return nil, err
		}

		loaders = append(loaders, l...)
	}

	return loaders, nil
}
