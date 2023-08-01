package bundle

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/fs"
	"github.com/mangalorg/mangal/provider/info"
	"github.com/mangalorg/mangal/provider/lua"
	"github.com/pelletier/go-toml"
)

func Loaders(dir string) ([]libmangal.ProviderLoader, error) {
	dirEntries, err := fs.Afero.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var bundleLoaders []libmangal.ProviderLoader
	for _, dirEntry := range dirEntries {
		// skip non directories
		if !dirEntry.IsDir() {
			continue
		}

		loaders, err := getLoaders(filepath.Join(dir, dirEntry.Name()))
		if err != nil {
			return nil, err
		}

		bundleLoaders = append(bundleLoaders, loaders...)
	}

	return bundleLoaders, nil
}

func singletone(loader libmangal.ProviderLoader) []libmangal.ProviderLoader {
	return []libmangal.ProviderLoader{loader}
}

func getLoaders(dir string) ([]libmangal.ProviderLoader, error) {
	infoFile, err := fs.Afero.OpenFile(
		filepath.Join(dir, info.Filename),
		os.O_RDONLY,
		0755,
	)

	if err != nil {
		return nil, err
	}

	defer infoFile.Close()

	decoder := toml.NewDecoder(infoFile)
	decoder.Strict(true)

	var providerInfo info.Info

	if err := decoder.Decode(&providerInfo); err != nil {
		return nil, err
	}

	switch providerInfo.Type {
	case info.TypeLua:
		loader, err := lua.NewLoader(providerInfo.Info, dir)
		if err != nil {
			return nil, err
		}

		return singletone(loader), nil
	case info.TypeBundle:
		return Loaders(dir)
	default:
		return nil, fmt.Errorf("unkown provider type: %#v", providerInfo.Type)
	}
}
