package provider

import (
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/luaprovider"
	"github.com/mangalorg/mangal/fs"
	"github.com/mangalorg/mangal/path"
	"github.com/philippgille/gokv/bbolt"
	"github.com/philippgille/gokv/encoding"
	"net/http"
	"path/filepath"
	"time"
)

func InstalledProviders() ([]libmangal.ProviderLoader, error) {
	var loaders []libmangal.ProviderLoader

	for _, getLoaders := range []func() ([]libmangal.ProviderLoader, error){
		installeLuaProviders,
	} {
		l, err := getLoaders()
		if err != nil {
			return nil, err
		}

		loaders = append(loaders, l...)
	}

	return loaders, nil
}

func installeLuaProviders() ([]libmangal.ProviderLoader, error) {
	dir := path.LuaProvidersDir()
	dirEntries, err := fs.Afero.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var loaders []libmangal.ProviderLoader
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}

		file, err := fs.Afero.ReadFile(filepath.Join(dir, dirEntry.Name()))
		if err != nil {
			return nil, err
		}

		info, err := luaprovider.ExtractInfo(file)
		if err != nil {
			continue
		}

		store, err := bbolt.NewStore(bbolt.Options{
			BucketName: info.Name,
			Path:       filepath.Join(path.CacheDir(), info.Name+".db"),
			Codec:      encoding.Gob,
		})
		if err != nil {
			continue
		}

		options := luaprovider.Options{
			HTTPClient: &http.Client{
				Timeout: time.Minute,
			},
			HTTPStore: store,
		}

		loader, err := luaprovider.NewLoader(file, options)
		if err != nil {
			return nil, err
		}

		loaders = append(loaders, loader)
	}

	return loaders, nil
}
