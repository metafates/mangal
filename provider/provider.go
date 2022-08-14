package provider

import (
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/provider/mangadex"
	"github.com/metafates/mangal/provider/manganelo"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	"github.com/metafates/mangal/where"
	"github.com/samber/lo"
	"os"
	"path/filepath"
)

type Provider struct {
	ID           string
	Name         string
	CreateSource func() (source.Source, error)
}

func (p Provider) String() string {
	return p.Name
}

var customProviderExtension = ".lua"

var defaultProviders = []*Provider{
	{
		ID:   manganelo.ID,
		Name: manganelo.Name,
		CreateSource: func() (source.Source, error) {
			return manganelo.New(), nil
		},
	},
	{
		ID:   mangadex.ID,
		Name: mangadex.Name,
		CreateSource: func() (source.Source, error) {
			return mangadex.New(), nil
		},
	},
}

func DefaultProviders() map[string]*Provider {
	providers := make(map[string]*Provider)

	for _, provider := range defaultProviders {
		providers[provider.Name] = provider
	}

	return providers
}

func CustomProviders() (map[string]*Provider, error) {
	files, err := filesystem.Get().ReadDir(where.Sources())

	if err != nil {
		return nil, err
	}

	providers := make(map[string]*Provider)
	paths := lo.FilterMap(files, func(f os.FileInfo, _ int) (string, bool) {
		if filepath.Ext(f.Name()) == customProviderExtension {
			return filepath.Join(where.Sources(), f.Name()), true
		}
		return "", false
	})

	for _, path := range paths {
		name := util.FileStem(path)
		path := path
		providers[name] = &Provider{
			ID:   source.IDfromName(name),
			Name: name,
			CreateSource: func() (source.Source, error) {
				return source.LoadSource(path, true)
			},
		}
	}

	return providers, nil
}
