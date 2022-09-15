package provider

import (
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/provider/custom"
	"github.com/metafates/mangal/provider/generic"
	"github.com/metafates/mangal/provider/mangadex"
	"github.com/metafates/mangal/provider/mangakakalot"
	"github.com/metafates/mangal/provider/manganelo"
	"github.com/metafates/mangal/provider/mangapill"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	"github.com/metafates/mangal/where"
	"github.com/samber/lo"
	"os"
	"path/filepath"
)

func init() {
	for _, conf := range []*generic.Configuration{
		manganelo.Config,
		mangakakalot.Config,
		mangapill.Config,
	} {
		defaultProviders = append(defaultProviders, &Provider{
			ID:   conf.ID(),
			Name: conf.Name,
			CreateSource: func() (source.Source, error) {
				return generic.New(conf), nil
			},
		})
	}
}

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
		ID:   mangadex.ID,
		Name: mangadex.Name,
		CreateSource: func() (source.Source, error) {
			return mangadex.New(), nil
		},
	},
}

func DefaultProviders() map[string]*Provider {
	providers := make(map[string]*Provider, len(defaultProviders))

	for _, provider := range defaultProviders {
		providers[provider.Name] = provider
	}

	return providers
}

func CustomProviders() map[string]*Provider {
	files, err := filesystem.Api().ReadDir(where.Sources())

	if err != nil {
		return make(map[string]*Provider)
	}

	paths := lo.FilterMap(files, func(f os.FileInfo, _ int) (string, bool) {
		if filepath.Ext(f.Name()) == customProviderExtension {
			return filepath.Join(where.Sources(), f.Name()), true
		}
		return "", false
	})
	providers := make(map[string]*Provider, len(paths))

	for _, path := range paths {
		name := util.FileStem(path)
		path := path
		providers[name] = &Provider{
			ID:   custom.IDfromName(name),
			Name: name,
			CreateSource: func() (source.Source, error) {
				return custom.LoadSource(path, true)
			},
		}
	}

	return providers
}

func Get(name string) (*Provider, bool) {
	for _, provider := range defaultProviders {
		if provider.Name == name {
			return provider, true
		}
	}

	for _, provider := range CustomProviders() {
		if provider.Name == name {
			return provider, true
		}
	}

	return nil, false
}
