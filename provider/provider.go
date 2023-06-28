package provider

import (
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/luaprovider"
	"github.com/mangalorg/mangal/fs"
	"github.com/mangalorg/mangal/path"
	"path/filepath"
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
	dirEntries, err := fs.FS.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var loaders []libmangal.ProviderLoader
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}

		file, err := fs.FS.ReadFile(filepath.Join(dir, dirEntry.Name()))
		if err != nil {
			return nil, err
		}

		loader, err := luaprovider.NewLoader(file, luaprovider.DefaultOptions())
		if err != nil {
			return nil, err
		}

		loaders = append(loaders, loader)
	}

	return loaders, nil
}
